package discord

import (
	"errors"

	"github.com/streame-gg/go-discord-wrapper/collection"
)

var ErrNotConvertable = errors.New("unable to resolve from input")

// ── Read-only store interfaces ────────────────────────────────────────────────

// GuildStore is the read-only view of cached guilds exposed to entity managers.
type GuildStore interface {
	Get(id Snowflake) (*Guild, bool)
	All() *collection.Collection[Snowflake, *Guild]
	Size() int
}

// ChannelStore is the read-only view of cached channels exposed to entity managers.
type ChannelStore interface {
	Get(id Snowflake) (*Channel, bool)
	All() *collection.Collection[Snowflake, *Channel]
	Size() int
}

// UserStore is the read-only view of cached users exposed to entity managers.
type UserStore interface {
	Get(id Snowflake) (*User, bool)
	All() *collection.Collection[Snowflake, *User]
	Size() int
}

// MemberStore is the read-only view of cached guild members exposed to entity managers.
type MemberStore interface {
	Get(guildID, userID Snowflake) (*GuildMember, bool)
	AllInGuild(guildID Snowflake) *collection.Collection[Snowflake, *GuildMember]
	Size() int
}

// MessageStore is the read-only view of cached messages exposed to entity managers.
type MessageStore interface {
	Get(channelID, messageID Snowflake) (*Message, bool)
	Channel(channelID Snowflake) *collection.Collection[Snowflake, *Message]
	Size() int
}

// RoleStore is the read-only view of cached roles exposed to entity managers.
type RoleStore interface {
	Get(roleID Snowflake) (*Role, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *Role]
	All() *collection.Collection[Snowflake, *Role]
	Size() int
}

// VoiceStateStore is the read-only view of cached voice states exposed to entity managers.
type VoiceStateStore interface {
	Get(guildID, userID Snowflake) (*VoiceState, bool)
	AllInGuild(guildID Snowflake) *collection.Collection[Snowflake, *VoiceState]
	Size() int
}

// SoundboardStore is the read-only view of cached soundboard sounds exposed to entity managers.
type SoundboardStore interface {
	Get(soundID Snowflake) (*SoundboardSound, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *SoundboardSound]
	Size() int
}

// ScheduledEventStore is the read-only view of cached scheduled events exposed to entity managers.
type ScheduledEventStore interface {
	Get(eventID Snowflake) (*GuildScheduledEvent, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *GuildScheduledEvent]
	Size() int
}

// StageInstanceStore is the read-only view of cached stage instances exposed to entity managers.
type StageInstanceStore interface {
	Get(instanceID Snowflake) (*StageInstance, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *StageInstance]
	Size() int
}

// EmojiStore is the read-only view of cached emojis exposed to entity managers.
type EmojiStore interface {
	Get(emojiID Snowflake) (*Emoji, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *Emoji]
	Size() int
}

// StickerStore is the read-only view of cached stickers exposed to entity managers.
type StickerStore interface {
	Get(stickerID Snowflake) (*Sticker, bool)
	GetByGuild(guildID Snowflake) *collection.Collection[Snowflake, *Sticker]
	Size() int
}

// ── Top-level cache interface ─────────────────────────────────────────────────

// Cache is the read-only cache interface accessible to entity managers via
// [EntityClient.ClientCache]. It is a strict read subset of the full cache.Cache
// interface — mutation methods (Set, Delete, Close) are intentionally absent.
// Returns nil from EntityClient.ClientCache if no cache is configured.
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
}
