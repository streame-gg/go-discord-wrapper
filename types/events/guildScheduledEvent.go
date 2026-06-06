package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-create
type GuildScheduledEventCreateEvent struct {
	discord.GuildScheduledEvent
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-update
type GuildScheduledEventUpdateEvent struct {
	NewEvent discord.GuildScheduledEvent  `json:"-"`
	OldEvent *discord.GuildScheduledEvent `json:"-"`
}

func (e *GuildScheduledEventUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewEvent)
}

func (e GuildScheduledEventUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.GuildScheduledEvent
		OldEvent *discord.GuildScheduledEvent `json:"old_event,omitempty"`
	}
	return json.Marshal(wire{e.NewEvent, e.OldEvent})
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-delete
type GuildScheduledEventDeleteEvent struct {
	discord.GuildScheduledEvent
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-user-add
type GuildScheduledEventUserAddEvent struct {
	GuildScheduledEventID discord.Snowflake `json:"guild_scheduled_event_id"`
	UserID                discord.Snowflake `json:"user_id"`
	GuildID               discord.Snowflake `json:"guild_id"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-user-remove
type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID discord.Snowflake `json:"guild_scheduled_event_id"`
	UserID                discord.Snowflake `json:"user_id"`
	GuildID               discord.Snowflake `json:"guild_id"`
}

func init() {
	RegisterEvent(GuildScheduledEventCreateEvent{})
	RegisterEvent(GuildScheduledEventUpdateEvent{})
	RegisterEvent(GuildScheduledEventDeleteEvent{})
	RegisterEvent(GuildScheduledEventUserAddEvent{})
	RegisterEvent(GuildScheduledEventUserRemoveEvent{})
}

func (e GuildScheduledEventCreateEvent) DesiredEventType() Event {
	return &GuildScheduledEventCreateEvent{}
}
func (e GuildScheduledEventCreateEvent) Event() EventType { return EventGuildScheduledEventCreate }

func (e GuildScheduledEventUpdateEvent) DesiredEventType() Event {
	return &GuildScheduledEventUpdateEvent{}
}
func (e GuildScheduledEventUpdateEvent) Event() EventType { return EventGuildScheduledEventUpdate }

func (e GuildScheduledEventDeleteEvent) DesiredEventType() Event {
	return &GuildScheduledEventDeleteEvent{}
}
func (e GuildScheduledEventDeleteEvent) Event() EventType { return EventGuildScheduledEventDelete }

func (e GuildScheduledEventUserAddEvent) DesiredEventType() Event {
	return &GuildScheduledEventUserAddEvent{}
}
func (e GuildScheduledEventUserAddEvent) Event() EventType { return EventGuildScheduledEventUserAdd }

func (e GuildScheduledEventUserRemoveEvent) DesiredEventType() Event {
	return &GuildScheduledEventUserRemoveEvent{}
}
func (e GuildScheduledEventUserRemoveEvent) Event() EventType {
	return EventGuildScheduledEventUserRemove
}
