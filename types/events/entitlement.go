package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &EntitlementCreateEvent{}
var _ Event = &EntitlementUpdateEvent{}
var _ Event = &EntitlementDeleteEvent{}

func init() {
	RegisterEvent(&EntitlementCreateEvent{})
	RegisterEvent(&EntitlementUpdateEvent{})
	RegisterEvent(&EntitlementDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#entitlement-create
type EntitlementCreateEvent struct {
	Entitlement discord.Entitlement
}

// https://docs.discord.com/developers/events/gateway-events#entitlement-update
type EntitlementUpdateEvent struct {
	NewEntitlement discord.Entitlement `json:"-"`
	// OldEntitlement is the state before the update, populated by the gateway
	// client from its in-memory entitlement cache. Nil if this entitlement was
	// not seen before (e.g. the client missed the preceding ENTITLEMENT_CREATE).
	OldEntitlement *discord.Entitlement `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#entitlement-delete
type EntitlementDeleteEvent struct {
	Entitlement discord.Entitlement
}

func (e *EntitlementCreateEvent) DesiredEventType() Event { return &EntitlementCreateEvent{} }
func (e *EntitlementCreateEvent) Event() EventType        { return EventEntitlementCreate }
func (e *EntitlementCreateEvent) UnmarshalJSON(byte []byte) error {
	return json.Unmarshal(byte, &e.Entitlement)
}

func (e *EntitlementUpdateEvent) DesiredEventType() Event { return &EntitlementUpdateEvent{} }
func (e *EntitlementUpdateEvent) Event() EventType        { return EventEntitlementUpdate }
func (e *EntitlementUpdateEvent) UnmarshalJSON(byte []byte) error {
	return json.Unmarshal(byte, &e.NewEntitlement)
}

func (e *EntitlementDeleteEvent) DesiredEventType() Event { return &EntitlementDeleteEvent{} }
func (e *EntitlementDeleteEvent) Event() EventType        { return EventEntitlementDelete }
func (e *EntitlementDeleteEvent) UnmarshalJSON(byte []byte) error {
	return json.Unmarshal(byte, &e.Entitlement)
}
