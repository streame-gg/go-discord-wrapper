package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type MessagePollVoteAddEvent struct {
	UserID    common.Snowflake  `json:"user_id"`
	ChannelID common.Snowflake  `json:"channel_id"`
	MessageID common.Snowflake  `json:"message_id"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int               `json:"answer_id"`
}

type MessagePollVoteRemoveEvent struct {
	UserID    common.Snowflake  `json:"user_id"`
	ChannelID common.Snowflake  `json:"channel_id"`
	MessageID common.Snowflake  `json:"message_id"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int               `json:"answer_id"`
}

func init() {
	RegisterEvent(MessagePollVoteAddEvent{})
	RegisterEvent(MessagePollVoteRemoveEvent{})
}

func (e MessagePollVoteAddEvent) DesiredEventType() Event { return &MessagePollVoteAddEvent{} }
func (e MessagePollVoteAddEvent) Event() EventType        { return EventMessagePollVoteAdd }

func (e MessagePollVoteRemoveEvent) DesiredEventType() Event { return &MessagePollVoteRemoveEvent{} }
func (e MessagePollVoteRemoveEvent) Event() EventType        { return EventMessagePollVoteRemove }
