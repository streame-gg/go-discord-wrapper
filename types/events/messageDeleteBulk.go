package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#message-delete-bulk
type MessageDeleteBulkEvent struct {
	IDs       []discord.Snowflake `json:"ids"`
	ChannelID discord.Snowflake   `json:"channel_id"`
	GuildID   *discord.Snowflake  `json:"guild_id,omitempty"`
}

func init() { RegisterEvent(MessageDeleteBulkEvent{}) }

func (m MessageDeleteBulkEvent) DesiredEventType() Event {
	return &MessageDeleteBulkEvent{}
}

func (m MessageDeleteBulkEvent) Event() EventType {
	return EventMessageDeleteBulk
}
