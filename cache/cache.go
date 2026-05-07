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

	"github.com/streame-gg/go-discord-wrapper/types/common"
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
	CategoryGuilds   OverflowCategory = 1 << 0
	CategoryChannels OverflowCategory = 1 << 1
	CategoryUsers    OverflowCategory = 1 << 2
	CategoryMembers  OverflowCategory = 1 << 3
	CategoryMessages OverflowCategory = 1 << 4
	CategoryRoles    OverflowCategory = 1 << 5
	CategoryVoiceStates     OverflowCategory = 1 << 6
	CategorySoundboard      OverflowCategory = 1 << 7
	CategoryScheduledEvents OverflowCategory = 1 << 8
	CategoryStageInstances  OverflowCategory = 1 << 9
	CategoryEmojis          OverflowCategory = 1 << 10
	CategoryPresences       OverflowCategory = 1 << 11

	// CategoryAll targets every entity store. This is the default.
	CategoryAll OverflowCategory = (1 << 12) - 1
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
}

// ── Entity store interfaces ───────────────────────────────────────────────────

// GuildStore is a thread-safe cache for [common.Guild] objects.
type GuildStore interface {
	Set(guild *common.Guild)
	Get(id common.Snowflake) (*common.Guild, bool)
	Delete(id common.Snowflake)
	All() []*common.Guild
	Size() int
}

// ChannelStore is a thread-safe cache for [common.Channel] objects.
type ChannelStore interface {
	Set(channel *common.Channel)
	Get(id common.Snowflake) (*common.Channel, bool)
	Delete(id common.Snowflake)
	All() []*common.Channel
	Size() int
}

// UserStore is a thread-safe cache for [common.User] objects.
type UserStore interface {
	Set(user *common.User)
	Get(id common.Snowflake) (*common.User, bool)
	Delete(id common.Snowflake)
	All() []*common.User
	Size() int
}

// MemberStore is a thread-safe cache for [common.GuildMember] objects,
// keyed by the composite (guildID, userID) pair.
type MemberStore interface {
	Set(guildID common.Snowflake, member *common.GuildMember)
	Get(guildID, userID common.Snowflake) (*common.GuildMember, bool)
	Delete(guildID, userID common.Snowflake)
	// DeleteGuild removes every member entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	AllInGuild(guildID common.Snowflake) []*common.GuildMember
	Size() int
}

// RoleStore is a thread-safe cache for [common.Role] objects, keyed by role ID.
// Roles can also be looked up or deleted by guild ID.
type RoleStore interface {
	Set(guildID common.Snowflake, role *common.Role)
	Get(roleID common.Snowflake) (*common.Role, bool)
	GetByGuild(guildID common.Snowflake) []*common.Role
	Delete(roleID common.Snowflake)
	// DeleteGuild removes every role entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	All() []*common.Role
	Size() int
}

// MessageStore caches per-channel message history in bounded ring buffers.
// Each channel is bounded by Options.Messages.MaxPerChannel.
type MessageStore interface {
	Add(msg *common.Message)
	Get(channelID, messageID common.Snowflake) (*common.Message, bool)
	Update(msg *common.Message)
	Delete(channelID, messageID common.Snowflake)
	DeleteBulk(channelID common.Snowflake, ids []common.Snowflake)
	// Channel returns cached messages for channelID newest-first,
	// excluding TTL-expired entries. Returns nil for unknown channels.
	Channel(channelID common.Snowflake) []*common.Message
	// DeleteChannel drops the entire ring for channelID. Call on CHANNEL_DELETE.
	DeleteChannel(channelID common.Snowflake)
	Size() int
}

// VoiceStateStore is a thread-safe cache for [common.VoiceState] objects,
// keyed by the composite (guildID, userID) pair.
type VoiceStateStore interface {
	Set(guildID common.Snowflake, state *common.VoiceState)
	Get(guildID, userID common.Snowflake) (*common.VoiceState, bool)
	Delete(guildID, userID common.Snowflake)
	// DeleteGuild removes every voice state entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	AllInGuild(guildID common.Snowflake) []*common.VoiceState
	Size() int
}

// SoundboardStore is a thread-safe cache for [common.SoundboardSound] objects.
type SoundboardStore interface {
	Set(guildID common.Snowflake, sound *common.SoundboardSound)
	Get(soundID common.Snowflake) (*common.SoundboardSound, bool)
	GetByGuild(guildID common.Snowflake) []*common.SoundboardSound
	SetAll(guildID common.Snowflake, sounds []*common.SoundboardSound)
	Delete(soundID common.Snowflake)
	// DeleteGuild removes every soundboard sound for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	Size() int
}

// ScheduledEventStore is a thread-safe cache for [common.GuildScheduledEvent] objects.
type ScheduledEventStore interface {
	Set(event *common.GuildScheduledEvent)
	Get(eventID common.Snowflake) (*common.GuildScheduledEvent, bool)
	GetByGuild(guildID common.Snowflake) []*common.GuildScheduledEvent
	Delete(eventID common.Snowflake)
	// DeleteGuild removes every scheduled event for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	Size() int
}

// StageInstanceStore is a thread-safe cache for [common.StageInstance] objects.
type StageInstanceStore interface {
	Set(instance *common.StageInstance)
	Get(instanceID common.Snowflake) (*common.StageInstance, bool)
	GetByGuild(guildID common.Snowflake) []*common.StageInstance
	Delete(instanceID common.Snowflake)
	// DeleteGuild removes every stage instance for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	Size() int
}

// EmojiStore is a thread-safe cache for [common.Emoji] objects.
type EmojiStore interface {
	Set(guildID common.Snowflake, emoji *common.Emoji)
	Get(emojiID common.Snowflake) (*common.Emoji, bool)
	GetByGuild(guildID common.Snowflake) []*common.Emoji
	SetAll(guildID common.Snowflake, emojis []*common.Emoji)
	Delete(emojiID common.Snowflake)
	// DeleteGuild removes every emoji for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	Size() int
}

// PresenceStore is a thread-safe cache for [common.Presence] objects.
type PresenceStore interface {
	Set(presence *common.Presence)
	Get(guildID, userID common.Snowflake) (*common.Presence, bool)
	GetByGuild(guildID common.Snowflake) []*common.Presence
	Delete(guildID, userID common.Snowflake)
	// DeleteGuild removes every presence entry for guildID. Call on GUILD_DELETE.
	DeleteGuild(guildID common.Snowflake)
	Size() int
}

// Cache is the top-level caching interface for Discord entities.
// All methods and the stores they return are safe for concurrent use.
// Close stops background goroutines; the cache must not be used after it returns.
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
	Presences() PresenceStore
	Close() error
}
