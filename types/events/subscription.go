package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func init() {
	RegisterEvent(&SubscriptionCreateEvent{})
	RegisterEvent(&SubscriptionUpdateEvent{})
	RegisterEvent(&SubscriptionDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#subscription-create
type SubscriptionCreateEvent struct {
	discord.Subscription
}

// https://docs.discord.com/developers/events/gateway-events#subscription-update
type SubscriptionUpdateEvent struct {
	NewSubscription discord.Subscription `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#subscription-delete
type SubscriptionDeleteEvent struct {
	discord.Subscription
}

func (e *SubscriptionCreateEvent) DesiredEventType() Event { return &SubscriptionCreateEvent{} }
func (e *SubscriptionCreateEvent) Event() EventType        { return EventSubscriptionCreate }

func (e *SubscriptionUpdateEvent) DesiredEventType() Event { return &SubscriptionUpdateEvent{} }
func (e *SubscriptionUpdateEvent) Event() EventType        { return EventSubscriptionUpdate }
func (e *SubscriptionUpdateEvent) UnmarshalJSON(data []byte) (err error) {
	return json.Unmarshal(data, &e.NewSubscription)
}

func (e *SubscriptionDeleteEvent) DesiredEventType() Event { return &SubscriptionDeleteEvent{} }
func (e *SubscriptionDeleteEvent) Event() EventType        { return EventSubscriptionDelete }
