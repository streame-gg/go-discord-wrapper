package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type MessageDeleteBulkEvent struct {
	IDs       []common.Snowflake `json:"ids"`
	ChannelID common.Snowflake   `json:"channel_id"`
	GuildID   *common.Snowflake  `json:"guild_id,omitempty"`
}

func init() { RegisterEvent(MessageDeleteBulkEvent{}) }

func (m MessageDeleteBulkEvent) DesiredEventType() Event {
	return &MessageDeleteBulkEvent{}
}

func (m MessageDeleteBulkEvent) Event() EventType {
	return EventMessageDeleteBulk
}
