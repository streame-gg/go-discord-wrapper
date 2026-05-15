package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

type EntitlementCreateEvent struct {
	discord.Entitlement
}

type EntitlementUpdateEvent struct {
	discord.Entitlement
}

type EntitlementDeleteEvent struct {
	discord.Entitlement
}

func init() {
	RegisterEvent(EntitlementCreateEvent{})
	RegisterEvent(EntitlementUpdateEvent{})
	RegisterEvent(EntitlementDeleteEvent{})
}

func (e EntitlementCreateEvent) DesiredEventType() Event { return &EntitlementCreateEvent{} }
func (e EntitlementCreateEvent) Event() EventType        { return EventEntitlementCreate }

func (e EntitlementUpdateEvent) DesiredEventType() Event { return &EntitlementUpdateEvent{} }
func (e EntitlementUpdateEvent) Event() EventType        { return EventEntitlementUpdate }

func (e EntitlementDeleteEvent) DesiredEventType() Event { return &EntitlementDeleteEvent{} }
func (e EntitlementDeleteEvent) Event() EventType        { return EventEntitlementDelete }
