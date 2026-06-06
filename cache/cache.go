// Package cache defines the caching interface for Discord entities and provides
// an in-process implementation via [NewMemoryCache].
//
// # Interface design
//
// [Cache] is the top-level interface with typed sub-stores:
//   - [GuildStore]         — discord guilds
//   - [ChannelStore]       — discord channels (including threads)
//   - [UserStore]          — discord users
//   - [MemberStore]        — guild members keyed by (guildID, userID)
//   - [MessageStore]       — per-channel message history with bounded ring buffers
//   - [RoleStore]          — guild roles
//   - [VoiceStateStore]    — guild voice states keyed by (guildID, userID)
//   - [SoundboardStore]    — guild soundboard sounds
//   - [ScheduledEventStore]— guild scheduled events
//   - [StageInstanceStore] — guild stage instances
//   - [EmojiStore]         — guild emojis
//   - [StickerStore]       — guild stickers
//   - [PresenceStore]      — guild presences (optional, high volume)
//
// Implement [Cache] for any external backend (Redis, MongoDB, etc.) and pass
// it to options.WithCache; the gateway client populates all stores automatically
// from gateway events.
//
// # Capacity and eviction
//
// [MemoryCache] keeps memory bounded through layered controls:
//
//  1. [Options.TTL] + [Options.EvictBehavior]: background sweeper removes
//     expired or access-idle entries on every [Options.SweepInterval] tick.
//
//  2. [Options.Limits]: hard per-category caps (MaxGuilds, MaxChannels, …) and
//     global caps (MaxEntries, MaxSizeMB). Violations trigger [Options.OnOverflow].
//
//  3. [Options.Messages.MaxPerChannel]: hard per-channel ring buffer bound.
//     An active channel can never accumulate more than this many messages.
//
// # Choosing an eviction strategy
//
// Use [Options.OnOverflow].[OverflowPolicy.ClearBy] to control which entries
// are selected for removal when a capacity limit is hit:
//
//   - [ClearByLastUsed] (default) — least-recently-accessed entries go first.
//   - [ClearByAge]                — oldest-inserted entries go first (FIFO).
//   - [ClearByFrequency]          — least-frequently-accessed entries go first (LFU).
//
// Use [OverflowPolicy.Target] (a bitmask of [OverflowCategory] values) to
// restrict overflow eviction to specific stores — e.g. evict only messages and
// users, leaving guilds untouched.
//
// Always call [Cache.Close] on shutdown to stop the background goroutine and
// allow all cached memory to be GC'd.
package cache

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Eviction strategy ─────────────────────────────────────────────────────────

// EvictBehavior controls which entries the background sweeper considers for
// removal on every [Options.SweepInterval] tick.
type EvictBehavior int

const (
	// EvictExpired removes only entries whose TTL has elapsed. Default.
	EvictExpired EvictBehavior = iota

	// EvictUnused removes TTL-expired entries AND entries that have not been
	// accessed within [Options.UnusedWindow], even if their TTL has not elapsed.
	// Effective for shedding cold data faster than TTL alone would.
	EvictUnused
)

// ClearBy determines how entries are ranked when selecting which to evict upon
// a capacity limit breach. Lower-ranked entries are removed first.
type ClearBy int

const (
	// ClearByLastUsed evicts the least-recently accessed entries first (LRU).
	// Default: hot entities stay, cold ones go.
	ClearByLastUsed ClearBy = iota

	// ClearByAge evicts entries by insertion time, oldest first (FIFO).
	// Useful when freshness matters more than access frequency — e.g. you always
	// want the most recently seen guilds cached, not necessarily the busiest ones.
	ClearByAge

	// ClearByFrequency evicts entries with the fewest total accesses first (LFU).
	// Useful when some entities are structurally rare (e.g. a guild the bot rarely
	// processes events for) and you want to prefer discarding them over busy ones.
	ClearByFrequency
)

// ── Overflow targeting ────────────────────────────────────────────────────────

// OverflowCategory is a bitmask that selects which entity stores are eligible
// for eviction when a global [Limits] value is breached.
//
// Per-category limits (MaxGuilds, MaxChannels, …) always target their own store
// regardless of this setting; OverflowCategory only applies to global limits.
type OverflowCategory uint

const (
	CategoryGuilds          OverflowCategory = 1 << 0
	CategoryChannels        OverflowCategory = 1 << 1
	CategoryUsers           OverflowCategory = 1 << 2
	CategoryMembers         OverflowCategory = 1 << 3
	CategoryMessages        OverflowCategory = 1 << 4
	CategoryRoles           OverflowCategory = 1 << 5
	CategoryVoiceStates     OverflowCategory = 1 << 6
	CategorySoundboard      OverflowCategory = 1 << 7
	CategoryScheduledEvents OverflowCategory = 1 << 8
	CategoryStageInstances  OverflowCategory = 1 << 9
	CategoryEmojis          OverflowCategory = 1 << 10
	CategoryStickers        OverflowCategory = 1 << 11
	CategoryPresences       OverflowCategory = 1 << 12
	CategoryBans            OverflowCategory = 1 << 13
	CategoryAutoModRules    OverflowCategory = 1 << 14
	CategoryInvites         OverflowCategory = 1 << 15

	// CategoryAll targets every entity store. This is the default.
	CategoryAll OverflowCategory = (1 << 16) - 1
)

// OverflowPolicy configures what happens when a [Limits] value is exceeded.
// Applied after every Set and during each sweep.
type OverflowPolicy struct {
	// ClearBy determines the ordering used when selecting entries to evict.
	// Default: ClearByLastUsed.
	ClearBy ClearBy

	// Target is a bitmask of which stores are eligible for eviction when a
	// global limit (MaxEntries or MaxSizeMB) is breached.
	// Default: CategoryAll.
	//
	// Example — spare guilds; evict only from volatile stores on overflow:
	//
	//   cache.CategoryMessages | cache.CategoryUsers | cache.CategoryMembers
	Target OverflowCategory
}

// ── Capacity limits ───────────────────────────────────────────────────────────

// Limits sets hard capacity bounds. Zero values mean "unlimited".
//
// Limits are enforced both immediately on each Set and on every sweep tick, so
// they act as soft real-time ceilings rather than deferred batch cuts.
type Limits struct {
	// MaxEntries caps the total entry count across all entity stores combined.
	// When exceeded, [OverflowPolicy] is applied to bring the total back down.
	MaxEntries int

	// MaxSizeMB caps the approximate total payload size in megabytes.
	// Size is estimated by JSON-marshalling each value on Set; this reflects
	// data payload size, not exact Go heap usage, but is a consistent bound.
	// Enabling this adds a json.Marshal call per Set; leave at 0 if throughput
	// is the priority and entry-count limits are sufficient.
	MaxSizeMB float64

	// Per-category entry caps — each enforced independently within its own store.
	// When a store exceeds its cap, [OverflowPolicy.ClearBy] selects which of
	// its own entries to evict (Target is ignored for per-category enforcement).
	MaxGuilds   int // max Guild entries
	MaxChannels int // max Channel entries
	MaxUsers    int // max User entries
	MaxMembers  int // max GuildMember entries (summed across all guilds)
	MaxMessages int // max total Message entries (summed across all channels)
	MaxRoles    int // max Role entries (summed across all guilds)
	MaxStickers int // max Sticker entries (summed across all guilds)
}

// ── Message options ───────────────────────────────────────────────────────────

// MessageOptions configures per-channel message caching behaviour.
type MessageOptions struct {
	// MaxPerChannel is the maximum number of messages cached per channel.
	// When the ring is full the oldest message is silently evicted to make room.
	// Default: 100. Set explicitly to 0 to disable message caching.
	MaxPerChannel int

	// TTL overrides [Options.TTL] for messages specifically.
	// Zero inherits the global TTL (which may itself be zero = no expiry).
	TTL time.Duration
}

// ── Top-level options ─────────────────────────────────────────────────────────

// Options configures a [Cache] backend.
type Options struct {
	// TTL is the default time-to-live for cache entries.
	// Zero means entries never expire from TTL alone.
	TTL time.Duration

	// SweepInterval controls how often the background cleaner runs.
	// Default: 30 seconds.
	SweepInterval time.Duration

	// EvictBehavior selects the TTL-sweep strategy. Default: EvictExpired.
	EvictBehavior EvictBehavior

	// UnusedWindow is used with EvictUnused. Entries not accessed within this
	// duration become eviction candidates regardless of their remaining TTL.
	// Zero disables access-idle eviction.
	UnusedWindow time.Duration

	// Limits sets capacity ceilings for entries and approximate payload size.
	Limits Limits

	// OnOverflow configures which entries are evicted and from which stores
	// when a global Limits value is breached.
	OnOverflow OverflowPolicy

	// Messages configures message-specific caching behaviour.
	Messages MessageOptions

	// WriteTimeout is the per-call timeout applied to every write operation
	// (Set, Delete, SetAll, Add, etc.) in remote cache backends (MongoDB, Redis).
	// Zero means no per-call timeout; the lifetime context is used directly.
	// Recommended value for production deployments: 2–5 seconds.
	WriteTimeout time.Duration

	// WriteQueueSize is the capacity of the async write queue used by remote
	// cache backends (MongoDB, Redis). Writes that arrive when the queue is full
	// are silently dropped. Zero or negative values use the default of 512.
	WriteQueueSize int
}

// ── Entity store interfaces ───────────────────────────────────────────────────

// GuildStore is a thread-safe cache for [discord.Guild] objects, keyed by guild ID.
// https://docs.discord.com/developers/resources/guild#guild-object
type GuildStore interface {
	Set(guild *discord.Guild)
	Get(id discord.Snowflake) (*discord.Guild, bool)
	Delete(id discord.Snowflake)
	// All returns a snapshot Collection of every cached guild, keyed by guild ID.
	All() *collection.Collection[discord.Snowflake, *discord.Guild]
	Size() int
}

// ChannelStore is a thread-safe cache for [discord.Channel] objects, keyed by channel ID.
// https://docs.discord.com/developers/resources/channel#channel-object
type ChannelStore interface {
	Set(channel *discord.Channel)
	Get(id discord.Snowflake) (*discord.Channel, bool)
	Delete(id discord.Snowflake)
	// All returns a snapshot Collection of every cached channel, keyed by channel ID.
	All() *collection.Collection[discord.Snowflake, *discord.Channel]
	Size() int
}

// UserStore is a thread-safe cache for [discord.User] objects, keyed by user ID.
// https://docs.discord.com/developers/resources/user#user-object
type UserStore interface {
	Set(user *discord.User)
	Get(id discord.Snowflake) (*discord.User, bool)
	Delete(id discord.Snowflake)
	// All returns a snapshot Collection of every cached user, keyed by user ID.
	All() *collection.Collection[discord.Snowflake, *discord.User]
	Size() int
}

// MemberStore is a thread-safe cache for [discord.GuildMember] objects,
// keyed by the composite (guildID, userID) pair.
// https://docs.discord.com/developers/resources/guild#guild-member-object
type MemberStore interface {
	Set(guildID discord.Snowflake, member *discord.GuildMember)
	Get(guildID, userID discord.Snowflake) (*discord.GuildMember, bool)
	Delete(guildID, userID discord.Snowflake)
	// DeleteGuild removes every member entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	// AllInGuild returns a snapshot Collection of all members for guildID,
	// keyed by user ID. The returned Collection is a copy.
	AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildMember]
	Size() int
}

// RoleStore is a thread-safe cache for [discord.Role] objects, keyed by role ID.
// Roles can also be looked up or deleted by guild ID.
// https://docs.discord.com/developers/topics/permissions#role-object
type RoleStore interface {
	Set(guildID discord.Snowflake, role *discord.Role)
	Get(roleID discord.Snowflake) (*discord.Role, bool)
	// GetByGuild returns a snapshot Collection of all roles for guildID,
	// keyed by role ID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Role]
	Delete(roleID discord.Snowflake)
	// DeleteGuild removes every role entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	// All returns a snapshot Collection of all roles across all guilds,
	// keyed by role ID. The returned Collection is a copy.
	All() *collection.Collection[discord.Snowflake, *discord.Role]
	Size() int
}

// MessageStore caches per-channel message history in bounded ring buffers.
// Each channel is bounded by Options.Messages.MaxPerChannel.
// https://docs.discord.com/developers/resources/message#message-object
type MessageStore interface {
	Add(msg *discord.Message)
	Get(channelID, messageID discord.Snowflake) (*discord.Message, bool)
	// Update replaces the stored copy of msg in its channel ring. If the
	// channel ring does not exist or msg.ID is absent from it, Update is a
	// silent no-op — the message is not inserted. Use Add to populate the
	// cache for the first time.
	Update(msg *discord.Message)
	Delete(channelID, messageID discord.Snowflake)
	DeleteBulk(channelID discord.Snowflake, ids []discord.Snowflake)
	// Channel returns a snapshot Collection of cached messages for channelID,
	// ordered newest-first, keyed by message ID. Returns an empty Collection
	// for unknown or disabled channels.
	Channel(channelID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Message]
	// DeleteChannel drops the entire ring for channelID. Call on CHANNEL_DELETE.
	DeleteChannel(channelID discord.Snowflake)
	Size() int
}

// VoiceStateStore is a thread-safe cache for [discord.VoiceState] objects,
// keyed by the composite (guildID, userID) pair.
// https://docs.discord.com/developers/resources/voice#voice-state-object
type VoiceStateStore interface {
	Set(guildID discord.Snowflake, state *discord.VoiceState)
	Get(guildID, userID discord.Snowflake) (*discord.VoiceState, bool)
	Delete(guildID, userID discord.Snowflake)
	// DeleteGuild removes every voice state entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	// AllInGuild returns a snapshot Collection of all voice states for guildID,
	// keyed by UserID. The returned Collection is a copy.
	AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.VoiceState]
	Size() int
}

// SoundboardStore is a thread-safe cache for [discord.SoundboardSound] objects.
// https://docs.discord.com/developers/resources/soundboard#soundboard-sound-object
type SoundboardStore interface {
	Set(guildID discord.Snowflake, sound *discord.SoundboardSound)
	Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool)
	// GetByGuild returns a snapshot Collection of all soundboard sounds for guildID,
	// keyed by SoundID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.SoundboardSound]
	SetAll(guildID discord.Snowflake, sounds []*discord.SoundboardSound)
	Delete(soundID discord.Snowflake)
	// DeleteGuild removes every soundboard sound for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// ScheduledEventStore is a thread-safe cache for [discord.GuildScheduledEvent] objects.
// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object
type ScheduledEventStore interface {
	Set(event *discord.GuildScheduledEvent)
	Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool)
	// GetByGuild returns a snapshot Collection of all scheduled events for guildID,
	// keyed by event ID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent]
	Delete(eventID discord.Snowflake)
	// DeleteGuild removes every scheduled event for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// StageInstanceStore is a thread-safe cache for [discord.StageInstance] objects.
// https://docs.discord.com/developers/resources/stage-instance#stage-instance-object
type StageInstanceStore interface {
	Set(instance *discord.StageInstance)
	Get(instanceID discord.Snowflake) (*discord.StageInstance, bool)
	// GetByGuild returns a snapshot Collection of all stage instances for guildID,
	// keyed by instance ID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.StageInstance]
	Delete(instanceID discord.Snowflake)
	// DeleteGuild removes every stage instance for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// EmojiStore is a thread-safe cache for [discord.Emoji] objects.
// https://docs.discord.com/developers/resources/emoji#emoji-object
type EmojiStore interface {
	Set(guildID discord.Snowflake, emoji *discord.Emoji)
	Get(emojiID discord.Snowflake) (*discord.Emoji, bool)
	// GetByGuild returns a snapshot Collection of all emojis for guildID,
	// keyed by emoji ID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Emoji]
	SetAll(guildID discord.Snowflake, emojis []*discord.Emoji)
	Delete(emojiID discord.Snowflake)
	// DeleteGuild removes every emoji for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// StickerStore is a thread-safe cache for [discord.Sticker] objects.
// https://docs.discord.com/developers/resources/sticker#sticker-object
type StickerStore interface {
	Set(guildID discord.Snowflake, sticker *discord.Sticker)
	Get(stickerID discord.Snowflake) (*discord.Sticker, bool)
	// GetByGuild returns a snapshot Collection of all stickers for guildID,
	// keyed by sticker ID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Sticker]
	SetAll(guildID discord.Snowflake, stickers []*discord.Sticker)
	Delete(stickerID discord.Snowflake)
	// DeleteGuild removes every sticker for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// PresenceStore is a thread-safe cache for [discord.Presence] objects.
// https://docs.discord.com/developers/events/gateway-events#presence-update
type PresenceStore interface {
	Set(presence *discord.Presence)
	Get(guildID, userID discord.Snowflake) (*discord.Presence, bool)
	// GetByGuild returns a snapshot Collection of all presences for guildID.
	// The returned Collection is a copy — mutating it does not affect the cache.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Presence]
	Delete(guildID, userID discord.Snowflake)
	// DeleteGuild removes every presence entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// BanStore is a thread-safe cache for [discord.Ban] objects,
// keyed by the composite (guildID, userID) pair.
// https://docs.discord.com/developers/resources/guild#ban-object
type BanStore interface {
	Set(guildID discord.Snowflake, ban *discord.Ban)
	Get(guildID, userID discord.Snowflake) (*discord.Ban, bool)
	Delete(guildID, userID discord.Snowflake)
	// DeleteGuild removes every ban entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	// AllInGuild returns a snapshot Collection of all bans for guildID,
	// keyed by user ID. The returned Collection is a copy.
	AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Ban]
	Size() int
}

// AutoModerationRuleStore is a thread-safe cache for [discord.AutoModerationRule] objects.
// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object
type AutoModerationRuleStore interface {
	Set(guildID discord.Snowflake, rule *discord.AutoModerationRule)
	Get(ruleID discord.Snowflake) (*discord.AutoModerationRule, bool)
	// GetByGuild returns a snapshot Collection of all automod rules for guildID,
	// keyed by rule ID. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.AutoModerationRule]
	Delete(ruleID discord.Snowflake)
	// DeleteGuild removes every automod rule for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// InviteStore is a thread-safe cache for [discord.Invite] objects, keyed by invite code.
// https://docs.discord.com/developers/resources/invite#invite-object
type InviteStore interface {
	Set(invite *discord.Invite)
	// SetWithGuild stores an invite and associates it with the given guild ID,
	// for cases where the invite's Guild field is not populated (e.g. gateway events).
	SetWithGuild(guildID discord.Snowflake, invite *discord.Invite)
	Get(code string) (*discord.Invite, bool)
	// GetByGuild returns a snapshot Collection of all invites for guildID,
	// keyed by code. The returned Collection is a copy.
	GetByGuild(guildID discord.Snowflake) *collection.Collection[string, *discord.Invite]
	Delete(code string)
	// DeleteGuild removes every invite for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID discord.Snowflake)
	Size() int
}

// Cache is the top-level caching interface for Discord entities.
// All methods and the stores they return are safe for concurrent use.
// Close stops background goroutines; the cache must not be used after it returns.
// https://docs.discord.com/developers/events/gateway-events#receive-events
type Cache interface {
	Guilds() GuildStore
	Channels() ChannelStore
	Users() UserStore
	Members() MemberStore
	Messages() MessageStore
	Roles() RoleStore
	VoiceStates() VoiceStateStore
	Soundboard() SoundboardStore
	ScheduledEvents() ScheduledEventStore
	StageInstances() StageInstanceStore
	Emojis() EmojiStore
	Stickers() StickerStore
	Presences() PresenceStore
	Bans() BanStore
	AutoModRules() AutoModerationRuleStore
	Invites() InviteStore
	Close() error
}
