package events

import (
	"encoding/json"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &ChannelCreateEvent{}
var _ Event = &ChannelUpdateEvent{}
var _ Event = &ChannelDeleteEvent{}
var _ Event = &ChannelPinsUpdateEvent{}

func init() {
	RegisterEvent(ChannelCreateEvent{})
	RegisterEvent(ChannelDeleteEvent{})
	RegisterEvent(ChannelUpdateEvent{})
	RegisterEvent(ChannelPinsUpdateEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#channel-create
type ChannelCreateEvent struct {
	Channel discord.Channel
}

// https://docs.discord.com/developers/events/gateway-events#channel-delete
type ChannelDeleteEvent struct {
	Channel discord.Channel
}

// https://docs.discord.com/developers/events/gateway-events#channel-pins-update
type ChannelPinsUpdateEvent struct {
	GuildID          *discord.Snowflake `json:"guild_id"`
	ChannelID        discord.Snowflake  `json:"channel_id"`
	LastPinTimestamp *time.Time         `json:"last_pin_timestamp,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#channel-update
type ChannelUpdateEvent struct {
	NewChannel discord.Channel  `json:"-"`
	OldChannel *discord.Channel `json:"-"`
}

func (e ChannelCreateEvent) DesiredEventType() Event {
	return &ChannelCreateEvent{}
}
func (e ChannelCreateEvent) Event() EventType {
	return EventChannelCreate
}
func (e ChannelCreateEvent) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &e.Channel)
}

func (e *ChannelUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewChannel)
}
func (e ChannelUpdateEvent) DesiredEventType() Event { return &ChannelUpdateEvent{} }
func (e ChannelUpdateEvent) Event() EventType        { return EventChannelUpdate }

func (e ChannelDeleteEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Channel)
}
func (e ChannelDeleteEvent) DesiredEventType() Event {
	return &ChannelDeleteEvent{}
}
func (e ChannelDeleteEvent) Event() EventType {
	return EventChannelDelete
}

func (e ChannelPinsUpdateEvent) DesiredEventType() Event { return &ChannelPinsUpdateEvent{} }
func (e ChannelPinsUpdateEvent) Event() EventType        { return EventChannelPinsUpdate }
