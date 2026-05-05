package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type EntitlementCreateEvent struct {
	common.Entitlement
}

type EntitlementUpdateEvent struct {
	common.Entitlement
}

type EntitlementDeleteEvent struct {
	common.Entitlement
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
