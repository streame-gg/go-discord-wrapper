package discord

import (
	"errors"

	"github.com/streame-gg/go-discord-wrapper/collection"
)

var ErrNotConvertable = errors.New("unable to resolve from input")

// ── Read-only store interfaces ────────────────────────────────────────────────

// GuildStore is the read-only view of cached guilds exposed to entity managers.
// https://docs.discord.com/developers/resources/guild#guild-object
type GuildStore interface {
	Get(id Snowflake) (*Guild, bool)
	All() *collection.Collection[Snowflake, *Guild]
	Size() int
}

// ChannelStore is the read-only view of cached channels exposed to entity managers.
// https://docs.discord.com/developers/resources/channel#channel-object
type ChannelStore interface {
	Get(id Snowflake) (*Channel, bool)
	All() *collection.Collection[Snowflake, *Channel]
	Size() int
}

// UserStore is the read-only view of cached users exposed to entity managers.
// https://docs.discord.com/developers/resources/user#user-object
type UserStore interface {
	Get(id Snowflake) (*User, bool)
	All() *collection.Collection[Snowflake, *User]
	Size() int
}

// MemberStore is the read-only view of cached guild members exposed to entity managers.
// https://docs.discord.com/developers/resources/guild#guild-member-object
type MemberStore interface {
	Get(guildID, userID Snowflake) (*GuildMember, bool)
	AllInGuild(guildID Snowflake) *collection.Collection[Snowflake, *GuildMember]
	Size() int
}

// MessageStore is the read-only view of cached messages exposed to entity managers.
// https://docs.discord.com/developers/resources/message#message-object
type MessageStore interface {
	Get(channelID, messageID Snowflake) (*Message, bool)
	Channel(channelID Snowflake) *collection.Collection[Snowflake, *Message]
	Size() int
}

// RoleStore is the read-only view of cached roles exposed to entity managers.
// https://docs.discord.com/developers/topics/permissions#role-object
type RoleStore interface {
	Get(roleID Snowflake) (*Role, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *Role]
	All() *collection.Collection[Snowflake, *Role]
	Size() int
}

// VoiceStateStore is the read-only view of cached voice states exposed to entity managers.
// https://docs.discord.com/developers/resources/voice#voice-state-object
type VoiceStateStore interface {
	Get(guildID, userID Snowflake) (*VoiceState, bool)
	AllInGuild(guildID Snowflake) *collection.Collection[Snowflake, *VoiceState]
	Size() int
}

// SoundboardStore is the read-only view of cached soundboard sounds exposed to entity managers.
// https://docs.discord.com/developers/resources/soundboard#soundboard-sound-object
type SoundboardStore interface {
	Get(soundID Snowflake) (*SoundboardSound, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *SoundboardSound]
	Size() int
}

// ScheduledEventStore is the read-only view of cached scheduled events exposed to entity managers.
// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object
type ScheduledEventStore interface {
	Get(eventID Snowflake) (*GuildScheduledEvent, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *GuildScheduledEvent]
	Size() int
}

// StageInstanceStore is the read-only view of cached stage instances exposed to entity managers.
// https://docs.discord.com/developers/resources/stage-instance#stage-instance-object
type StageInstanceStore interface {
	Get(instanceID Snowflake) (*StageInstance, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *StageInstance]
	Size() int
}

// EmojiStore is the read-only view of cached emojis exposed to entity managers.
// https://docs.discord.com/developers/resources/emoji#emoji-object
type EmojiStore interface {
	Get(emojiID Snowflake) (*Emoji, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, Emoji]
	Size() int
}

// StickerStore is the read-only view of cached stickers exposed to entity managers.
// https://docs.discord.com/developers/resources/sticker#sticker-object
type StickerStore interface {
	Get(stickerID Snowflake) (*Sticker, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, Sticker]
	Size() int
}

// BanStore is the read-only view of cached guild bans exposed to entity managers.
// https://docs.discord.com/developers/resources/guild#ban-object
type BanStore interface {
	Get(guildID, userID Snowflake) (*Ban, bool)
	AllInGuild(guildID Snowflake) *collection.Collection[Snowflake, Ban]
	Size() int
}

// AutoModerationRuleStore is the read-only view of cached automod rules exposed to entity managers.
// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object
type AutoModerationRuleStore interface {
	Get(ruleID Snowflake) (*AutoModerationRule, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, AutoModerationRule]
	Size() int
}

// InviteStore is the read-only view of cached invites exposed to entity managers.
// https://docs.discord.com/developers/resources/invite#invite-object
type InviteStore interface {
	Get(code string) (*Invite, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[string, Invite]
	Size() int
}

// ── Top-level cache interface ─────────────────────────────────────────────────

// Cache is the read-only cache interface accessible to entity managers via
// [EntityClient.ClientCache]. It is a strict read subset of the full cache.Cache
// interface — mutation methods (Set, Delete, Close) are intentionally absent.
// Returns nil from EntityClient.ClientCache if no cache is configured.
// Caches the entity objects delivered by Discord gateway events:
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
	Bans() BanStore
	AutoModRules() AutoModerationRuleStore
	Invites() InviteStore
}
