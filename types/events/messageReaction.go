package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &MessageReactionAddEvent{}
var _ Event = &MessageReactionRemoveEvent{}
var _ Event = &MessageReactionRemoveAllEvent{}
var _ Event = &MessageReactionRemoveEmojiEvent{}

func init() {
	RegisterEvent(MessageReactionAddEvent{})
	RegisterEvent(MessageReactionRemoveEvent{})
	RegisterEvent(MessageReactionRemoveAllEvent{})
	RegisterEvent(MessageReactionRemoveEmojiEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#message-reaction-add
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
	Type            discord.ReactionType `json:"type"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
	Message *discord.Message `json:"-"`
	User    *discord.User    `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-reaction-remove
type MessageReactionRemoveEvent struct {
	UserID    discord.Snowflake    `json:"user_id"`
	ChannelID discord.Snowflake    `json:"channel_id"`
	MessageID discord.Snowflake    `json:"message_id"`
	GuildID   *discord.Snowflake   `json:"guild_id,omitempty"`
	Emoji     discord.Emoji        `json:"emoji"`
	Burst     bool                 `json:"burst"`
	Type      discord.ReactionType `json:"type"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
	Message *discord.Message `json:"-"`
	User    *discord.User    `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-reaction-remove-all
type MessageReactionRemoveAllEvent struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
	Message *discord.Message `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-reaction-remove-emoji
type MessageReactionRemoveEmojiEvent struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	MessageID discord.Snowflake  `json:"message_id"`
	Emoji     discord.Emoji      `json:"emoji"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
	Message *discord.Message `json:"-"`
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
