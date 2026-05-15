package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ChannelPinsUpdateEvent struct {
	GuildID          *discord.Snowflake `json:"guild_id,omitempty"`
	ChannelID        discord.Snowflake  `json:"channel_id"`
	LastPinTimestamp *time.Time         `json:"last_pin_timestamp,omitempty"`
}

func init() { RegisterEvent(ChannelPinsUpdateEvent{}) }

func (e ChannelPinsUpdateEvent) DesiredEventType() Event { return &ChannelPinsUpdateEvent{} }
func (e ChannelPinsUpdateEvent) Event() EventType        { return EventChannelPinsUpdate }
