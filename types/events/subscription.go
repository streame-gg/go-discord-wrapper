package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

type SubscriptionCreateEvent struct {
	discord.Subscription
}

type SubscriptionUpdateEvent struct {
	discord.Subscription
}

type SubscriptionDeleteEvent struct {
	discord.Subscription
}

func init() {
	RegisterEvent(SubscriptionCreateEvent{})
	RegisterEvent(SubscriptionUpdateEvent{})
	RegisterEvent(SubscriptionDeleteEvent{})
}

func (e SubscriptionCreateEvent) DesiredEventType() Event { return &SubscriptionCreateEvent{} }
func (e SubscriptionCreateEvent) Event() EventType        { return EventSubscriptionCreate }

func (e SubscriptionUpdateEvent) DesiredEventType() Event { return &SubscriptionUpdateEvent{} }
func (e SubscriptionUpdateEvent) Event() EventType        { return EventSubscriptionUpdate }

func (e SubscriptionDeleteEvent) DesiredEventType() Event { return &SubscriptionDeleteEvent{} }
func (e SubscriptionDeleteEvent) Event() EventType        { return EventSubscriptionDelete }
