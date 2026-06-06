package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#message-delete
type MessageDeleteEvent struct {
	ID        discord.Snowflake  `json:"id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
}

func init() { RegisterEvent(MessageDeleteEvent{}) }

func (m MessageDeleteEvent) DesiredEventType() Event {
	return &MessageDeleteEvent{}
}

func (m MessageDeleteEvent) Event() EventType {
	return EventMessageDelete
}
