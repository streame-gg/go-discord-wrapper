package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildScheduledEventCreateEvent{}
var _ Event = &GuildScheduledEventUpdateEvent{}
var _ Event = &GuildScheduledEventDeleteEvent{}
var _ Event = &GuildScheduledEventUserAddEvent{}
var _ Event = &GuildScheduledEventUserRemoveEvent{}

func init() {
	RegisterEvent(&GuildScheduledEventCreateEvent{})
	RegisterEvent(&GuildScheduledEventUpdateEvent{})
	RegisterEvent(&GuildScheduledEventDeleteEvent{})
	RegisterEvent(&GuildScheduledEventUserAddEvent{})
	RegisterEvent(&GuildScheduledEventUserRemoveEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-create
type GuildScheduledEventCreateEvent struct {
	discord.GuildScheduledEvent
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-update
type GuildScheduledEventUpdateEvent struct {
	NewScheduledEvent discord.GuildScheduledEvent  `json:"-"`
	OldScheduledEvent *discord.GuildScheduledEvent `json:"-"`
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

	GuildScheduledEvent *discord.GuildScheduledEvent `json:"-"`
	User                *discord.User                `json:"-"`
	Guild               *discord.Guild               `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-scheduled-event-user-remove
type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID discord.Snowflake `json:"guild_scheduled_event_id"`
	UserID                discord.Snowflake `json:"user_id"`
	GuildID               discord.Snowflake `json:"guild_id"`

	GuildScheduledEvent *discord.GuildScheduledEvent `json:"-"`
	User                *discord.User                `json:"-"`
	Guild               *discord.Guild               `json:"-"`
}

func (e *GuildScheduledEventCreateEvent) DesiredEventType() Event {
	return &GuildScheduledEventCreateEvent{}
}
func (e *GuildScheduledEventCreateEvent) Event() EventType { return EventGuildScheduledEventCreate }

func (e *GuildScheduledEventUpdateEvent) DesiredEventType() Event {
	return &GuildScheduledEventUpdateEvent{}
}
func (e *GuildScheduledEventUpdateEvent) Event() EventType { return EventGuildScheduledEventUpdate }
func (e *GuildScheduledEventUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewScheduledEvent)
}

func (e *GuildScheduledEventDeleteEvent) DesiredEventType() Event {
	return &GuildScheduledEventDeleteEvent{}
}
func (e *GuildScheduledEventDeleteEvent) Event() EventType { return EventGuildScheduledEventDelete }

func (e *GuildScheduledEventUserAddEvent) DesiredEventType() Event {
	return &GuildScheduledEventUserAddEvent{}
}
func (e *GuildScheduledEventUserAddEvent) Event() EventType { return EventGuildScheduledEventUserAdd }

func (e *GuildScheduledEventUserRemoveEvent) DesiredEventType() Event {
	return &GuildScheduledEventUserRemoveEvent{}
}
func (e *GuildScheduledEventUserRemoveEvent) Event() EventType {
	return EventGuildScheduledEventUserRemove
}
