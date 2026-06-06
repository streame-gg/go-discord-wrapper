package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#entitlement-create
type EntitlementCreateEvent struct {
	discord.Entitlement
}

// https://docs.discord.com/developers/events/gateway-events#entitlement-update
type EntitlementUpdateEvent struct {
	NewEntitlement discord.Entitlement `json:"-"`
	// OldEntitlement is the state before the update, populated by the gateway
	// client from its in-memory entitlement cache. Nil if this entitlement was
	// not seen before (e.g. the client missed the preceding ENTITLEMENT_CREATE).
	OldEntitlement *discord.Entitlement `json:"-"`
}

func (e *EntitlementUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewEntitlement)
}

func (e EntitlementUpdateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&e.NewEntitlement)
}

// https://docs.discord.com/developers/events/gateway-events#entitlement-delete
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
