package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type GuildScheduledEventCreateEvent struct {
	common.GuildScheduledEvent
}

type GuildScheduledEventUpdateEvent struct {
	common.GuildScheduledEvent
}

type GuildScheduledEventDeleteEvent struct {
	common.GuildScheduledEvent
}

type GuildScheduledEventUserAddEvent struct {
	GuildScheduledEventID common.Snowflake `json:"guild_scheduled_event_id"`
	UserID                common.Snowflake `json:"user_id"`
	GuildID               common.Snowflake `json:"guild_id"`
}

type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID common.Snowflake `json:"guild_scheduled_event_id"`
	UserID                common.Snowflake `json:"user_id"`
	GuildID               common.Snowflake `json:"guild_id"`
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
