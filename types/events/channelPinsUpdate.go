package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type ChannelPinsUpdateEvent struct {
	GuildID          *common.Snowflake `json:"guild_id,omitempty"`
	ChannelID        common.Snowflake  `json:"channel_id"`
	LastPinTimestamp *time.Time        `json:"last_pin_timestamp,omitempty"`
}

func init() { RegisterEvent(ChannelPinsUpdateEvent{}) }

func (e ChannelPinsUpdateEvent) DesiredEventType() Event { return &ChannelPinsUpdateEvent{} }
func (e ChannelPinsUpdateEvent) Event() EventType        { return EventChannelPinsUpdate }
