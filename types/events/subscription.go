package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func init() {
	RegisterEvent(SubscriptionCreateEvent{})
	RegisterEvent(SubscriptionUpdateEvent{})
	RegisterEvent(SubscriptionDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#subscription-create
type SubscriptionCreateEvent struct {
	Subscription discord.Subscription `json:"-"`

	User *discord.User `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#subscription-update
type SubscriptionUpdateEvent struct {
	Subscription discord.Subscription `json:"-"`

	User *discord.User `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#subscription-delete
type SubscriptionDeleteEvent struct {
	Subscription discord.Subscription `json:"-"`

	User *discord.User `json:"-"`
}

func (e SubscriptionCreateEvent) DesiredEventType() Event { return &SubscriptionCreateEvent{} }
func (e SubscriptionCreateEvent) Event() EventType        { return EventSubscriptionCreate }
func (e SubscriptionCreateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Subscription)
}

func (e SubscriptionUpdateEvent) DesiredEventType() Event { return &SubscriptionUpdateEvent{} }
func (e SubscriptionUpdateEvent) Event() EventType        { return EventSubscriptionUpdate }
func (e SubscriptionUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Subscription)
}

func (e SubscriptionDeleteEvent) DesiredEventType() Event { return &SubscriptionDeleteEvent{} }
func (e SubscriptionDeleteEvent) Event() EventType        { return EventSubscriptionDelete }
func (e SubscriptionDeleteEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Subscription)
}
