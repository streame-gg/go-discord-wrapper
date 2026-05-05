package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type ReactionType int

const (
	ReactionTypeNormal ReactionType = 0
	ReactionTypeBurst  ReactionType = 1
)

type MessageReactionAddEvent struct {
	UserID          common.Snowflake    `json:"user_id"`
	ChannelID       common.Snowflake    `json:"channel_id"`
	MessageID       common.Snowflake    `json:"message_id"`
	GuildID         *common.Snowflake   `json:"guild_id,omitempty"`
	Member          *common.GuildMember `json:"member,omitempty"`
	Emoji           common.Emoji        `json:"emoji"`
	MessageAuthorID *common.Snowflake   `json:"message_author_id,omitempty"`
	Burst           bool                `json:"burst"`
	BurstColors     []string            `json:"burst_colors,omitempty"`
	Type            ReactionType        `json:"type"`
}

type MessageReactionRemoveEvent struct {
	UserID    common.Snowflake  `json:"user_id"`
	ChannelID common.Snowflake  `json:"channel_id"`
	MessageID common.Snowflake  `json:"message_id"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
	Emoji     common.Emoji      `json:"emoji"`
	Burst     bool              `json:"burst"`
	Type      ReactionType      `json:"type"`
}

type MessageReactionRemoveAllEvent struct {
	ChannelID common.Snowflake  `json:"channel_id"`
	MessageID common.Snowflake  `json:"message_id"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
}

type MessageReactionRemoveEmojiEvent struct {
	ChannelID common.Snowflake  `json:"channel_id"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
	MessageID common.Snowflake  `json:"message_id"`
	Emoji     common.Emoji      `json:"emoji"`
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
