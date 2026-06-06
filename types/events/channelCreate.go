package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// https://docs.discord.com/developers/events/gateway-events#channel-create
type ChannelCreateEvent struct {
	discord.Channel
}

func init() { RegisterEvent(ChannelCreateEvent{}) }

func (e ChannelCreateEvent) DesiredEventType() Event {
	return &ChannelCreateEvent{}
}

func (e ChannelCreateEvent) Event() EventType {
	return EventChannelCreate
}
