package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

type ChannelDeleteEvent struct {
	discord.Channel
}

func init() { RegisterEvent(ChannelDeleteEvent{}) }

func (e ChannelDeleteEvent) DesiredEventType() Event {
	return &ChannelDeleteEvent{}
}

func (e ChannelDeleteEvent) Event() EventType {
	return EventChannelDelete
}
