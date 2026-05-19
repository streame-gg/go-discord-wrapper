package discord

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
)

// ── Guild sub-manager interfaces ──────────────────────────────────────────────

// MemberManager manages guild members for a single guild.
type MemberManager interface {
	Cache() *collection.Collection[Snowflake, *GuildMember]
	Get(userID Snowflake) (*GuildMember, bool)
	Fetch(ctx context.Context, userID Snowflake) (*GuildMember, error)
	FetchAll(ctx context.Context, opts FetchMembersOptions) (*collection.Collection[Snowflake, *GuildMember], error)
	Resolve(input any) *GuildMember
	Size() int
}

// RoleManager manages roles for a single guild.
type RoleManager interface {
	Cache() *collection.Collection[Snowflake, *Role]
	Get(roleID Snowflake) (*Role, bool)
	Fetch(ctx context.Context, roleID Snowflake) (*Role, error)
	FetchAll(ctx context.Context) (*collection.Collection[Snowflake, *Role], error)
	Create(ctx context.Context, opts RoleCreateOptions) (*Role, error)
	Resolve(input any) *Role
	Size() int
}

// GuildChannelManager manages channels for a single guild.
type GuildChannelManager interface {
	Cache() *collection.Collection[Snowflake, *Channel]
	Get(channelID Snowflake) (*Channel, bool)
	Fetch(ctx context.Context, channelID Snowflake) (*Channel, error)
	FetchAll(ctx context.Context) (*collection.Collection[Snowflake, *Channel], error)
	Create(ctx context.Context, opts ChannelCreateOptions) (*Channel, error)
	Resolve(input any) *Channel
	Size() int
}

// EmojiManager manages emojis for a single guild.
type EmojiManager interface {
	Cache() *collection.Collection[Snowflake, *Emoji]
	Get(emojiID Snowflake) (*Emoji, bool)
	Fetch(ctx context.Context, emojiID Snowflake) (*Emoji, error)
	FetchAll(ctx context.Context) (*collection.Collection[Snowflake, *Emoji], error)
	Create(ctx context.Context, opts EmojiCreateOptions) (*Emoji, error)
	Resolve(input any) *Emoji
	Size() int
}

// StickerManager manages stickers for a single guild.
type StickerManager interface {
	Cache() *collection.Collection[Snowflake, *Sticker]
	Get(stickerID Snowflake) (*Sticker, bool)
	Fetch(ctx context.Context, stickerID Snowflake) (*Sticker, error)
	FetchAll(ctx context.Context) (*collection.Collection[Snowflake, *Sticker], error)
	Create(ctx context.Context, opts StickerCreateOptions) (*Sticker, error)
	Resolve(input any) *Sticker
	Size() int
}

// BanManager manages bans for a single guild. Bans are not cached.
type BanManager interface {
	Fetch(ctx context.Context, userID Snowflake) (*Ban, error)
	FetchAll(ctx context.Context, opts FetchBansOptions) ([]*Ban, error)
	Create(ctx context.Context, userID Snowflake, opts BanOptions) error
	Remove(ctx context.Context, userID Snowflake, reason *string) error
}

// ScheduledEventManager manages scheduled events for a single guild.
type ScheduledEventManager interface {
	Cache() *collection.Collection[Snowflake, *GuildScheduledEvent]
	Get(eventID Snowflake) (*GuildScheduledEvent, bool)
	Fetch(ctx context.Context, eventID Snowflake) (*GuildScheduledEvent, error)
	FetchAll(ctx context.Context) (*collection.Collection[Snowflake, *GuildScheduledEvent], error)
	Create(ctx context.Context, opts ScheduledEventCreateOptions) (*GuildScheduledEvent, error)
	Resolve(input any) *GuildScheduledEvent
	Size() int
}

// StageInstanceManager manages stage instances for a single guild.
// Note: Discord has no REST endpoint to list all stage instances for a guild;
// use Cache() (populated from gateway events) or Fetch(ctx, channelID).
type StageInstanceManager interface {
	Cache() *collection.Collection[Snowflake, *StageInstance]
	Get(instanceID Snowflake) (*StageInstance, bool)
	Fetch(ctx context.Context, channelID Snowflake) (*StageInstance, error)
	Create(ctx context.Context, opts StageCreateOptions) (*StageInstance, error)
	Resolve(input any) *StageInstance
	Size() int
}

// SoundboardManager manages soundboard sounds for a single guild.
type SoundboardManager interface {
	Cache() *collection.Collection[Snowflake, *SoundboardSound]
	Get(soundID Snowflake) (*SoundboardSound, bool)
	Fetch(ctx context.Context, soundID Snowflake) (*SoundboardSound, error)
	FetchAll(ctx context.Context) (*collection.Collection[Snowflake, *SoundboardSound], error)
	Resolve(input any) *SoundboardSound
	Size() int
}

// GuildInviteManager manages invites for a single guild. Invites are not cached.
type GuildInviteManager interface {
	FetchAll(ctx context.Context) ([]*Invite, error)
}

// VoiceStateManager manages voice states for a single guild.
// Voice states are gateway-only; there is no REST fetch endpoint.
type VoiceStateManager interface {
	Cache() *collection.Collection[Snowflake, *VoiceState]
	Get(userID Snowflake) (*VoiceState, bool)
	Resolve(input any) *VoiceState
	Size() int
}

// AutoModRuleManager manages auto moderation rules for a single guild.
// Rules are not cached.
type AutoModRuleManager interface {
	Fetch(ctx context.Context, ruleID Snowflake) (*AutoModerationRule, error)
	FetchAll(ctx context.Context) ([]*AutoModerationRule, error)
	Create(ctx context.Context, opts RuleCreateOptions) (*AutoModerationRule, error)
}

// GuildWebhookManager manages webhooks for a single guild. Webhooks are not cached.
type GuildWebhookManager interface {
	FetchAll(ctx context.Context) ([]*Webhook, error)
}

// IntegrationManager manages integrations for a single guild. Integrations are not cached.
type IntegrationManager interface {
	FetchAll(ctx context.Context) ([]*Integration, error)
	Remove(ctx context.Context, integrationID Snowflake, reason *string) error
}

// ── Channel sub-manager interfaces ────────────────────────────────────────────

// MessageManager manages messages for a single channel.
type MessageManager interface {
	Cache() *collection.Collection[Snowflake, *Message]
	Get(messageID Snowflake) (*Message, bool)
	Fetch(ctx context.Context, messageID Snowflake) (*Message, error)
	FetchAll(ctx context.Context, opts FetchMessagesOptions) (*collection.Collection[Snowflake, *Message], error)
	Create(ctx context.Context, opts MessageCreateOptions) (*Message, error)
	Resolve(input any) *Message
	Size() int
}

// ThreadManager manages threads for a single channel.
// Threads are channels cached in the channel store; there is no REST endpoint
// to list all threads for a channel without additional filtering.
type ThreadManager interface {
	Cache() *collection.Collection[Snowflake, *Channel]
	Get(threadID Snowflake) (*Channel, bool)
	Fetch(ctx context.Context, threadID Snowflake) (*Channel, error)
	Resolve(input any) *Channel
	Size() int
}

// ── Client-level manager interfaces ──────────────────────────────────────────

// GuildManager is the client-level manager for guilds.
type GuildManager interface {
	Cache() *collection.Collection[Snowflake, *Guild]
	Get(guildID Snowflake) (*Guild, bool)
	Fetch(ctx context.Context, guildID Snowflake) (*Guild, error)
	Resolve(input any) *Guild
	Size() int
}

// UserManager is the client-level manager for users.
type UserManager interface {
	Cache() *collection.Collection[Snowflake, *User]
	Get(userID Snowflake) (*User, bool)
	Fetch(ctx context.Context, userID Snowflake) (*User, error)
	Resolve(input any) *User
	Size() int
}

// ChannelManager is the client-level manager for channels.
type ChannelManager interface {
	Cache() *collection.Collection[Snowflake, *Channel]
	Get(channelID Snowflake) (*Channel, bool)
	Fetch(ctx context.Context, channelID Snowflake) (*Channel, error)
	Resolve(input any) *Channel
	Size() int
}
