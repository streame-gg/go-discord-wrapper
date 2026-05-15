package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ReactionType int

const (
	ReactionTypeNormal ReactionType = 0
	ReactionTypeBurst  ReactionType = 1
)

type MessageReactionAddEvent struct {
	UserID          discord.Snowflake    `json:"user_id"`
	ChannelID       discord.Snowflake    `json:"channel_id"`
	MessageID       discord.Snowflake    `json:"message_id"`
	GuildID         *discord.Snowflake   `json:"guild_id,omitempty"`
	Member          *discord.GuildMember `json:"member,omitempty"`
	Emoji           discord.Emoji        `json:"emoji"`
	MessageAuthorID *discord.Snowflake   `json:"message_author_id,omitempty"`
	Burst           bool                 `json:"burst"`
	BurstColors     []string             `json:"burst_colors,omitempty"`
	Type            ReactionType         `json:"type"`
}

type MessageReactionRemoveEvent struct {
	UserID    discord.Snowflake  `json:"user_id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	Emoji     discord.Emoji      `json:"emoji"`
	Burst     bool               `json:"burst"`
	Type      ReactionType       `json:"type"`
}

type MessageReactionRemoveAllEvent struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
}

type MessageReactionRemoveEmojiEvent struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	MessageID discord.Snowflake  `json:"message_id"`
	Emoji     discord.Emoji      `json:"emoji"`
}

func init() {
	RegisterEvent(MessageReactionAddEvent{})
	RegisterEvent(MessageReactionRemoveEvent{})
	RegisterEvent(MessageReactionRemoveAllEvent{})
	RegisterEvent(MessageReactionRemoveEmojiEvent{})
}

func (e MessageReactionAddEvent) DesiredEventType() Event { return &MessageReactionAddEvent{} }
func (e MessageReactionAddEvent) Event() EventType        { return EventMessageReactionAdd }

func (e MessageReactionRemoveEvent) DesiredEventType() Event { return &MessageReactionRemoveEvent{} }
func (e MessageReactionRemoveEvent) Event() EventType        { return EventMessageReactionRemove }

func (e MessageReactionRemoveAllEvent) DesiredEventType() Event {
	return &MessageReactionRemoveAllEvent{}
}
func (e MessageReactionRemoveAllEvent) Event() EventType { return EventMessageReactionRemoveAll }

func (e MessageReactionRemoveEmojiEvent) DesiredEventType() Event {
	return &MessageReactionRemoveEmojiEvent{}
}
func (e MessageReactionRemoveEmojiEvent) Event() EventType { return EventMessageReactionRemoveEmoji }
