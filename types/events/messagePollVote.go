package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#message-poll-vote-add
type MessagePollVoteAddEvent struct {
	UserID    discord.Snowflake  `json:"user_id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int                `json:"answer_id"`
}

// https://docs.discord.com/developers/events/gateway-events#message-poll-vote-remove
type MessagePollVoteRemoveEvent struct {
	UserID    discord.Snowflake  `json:"user_id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int                `json:"answer_id"`
}

func init() {
	RegisterEvent(MessagePollVoteAddEvent{})
	RegisterEvent(MessagePollVoteRemoveEvent{})
}

func (e MessagePollVoteAddEvent) DesiredEventType() Event { return &MessagePollVoteAddEvent{} }
func (e MessagePollVoteAddEvent) Event() EventType        { return EventMessagePollVoteAdd }

func (e MessagePollVoteRemoveEvent) DesiredEventType() Event { return &MessagePollVoteRemoveEvent{} }
func (e MessagePollVoteRemoveEvent) Event() EventType        { return EventMessagePollVoteRemove }
