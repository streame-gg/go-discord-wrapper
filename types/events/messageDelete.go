package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type MessageDeleteEvent struct {
	ID        common.Snowflake  `json:"id"`
	ChannelID common.Snowflake  `json:"channel_id"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
}

func init() { RegisterEvent(MessageDeleteEvent{}) }

func (m MessageDeleteEvent) DesiredEventType() Event {
	return &MessageDeleteEvent{}
}

func (m MessageDeleteEvent) Event() EventType {
	return EventMessageDelete
}
