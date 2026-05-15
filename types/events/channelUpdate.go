package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

type ChannelUpdateEvent struct {
	discord.Channel
}

func init() { RegisterEvent(ChannelUpdateEvent{}) }

func (e ChannelUpdateEvent) DesiredEventType() Event { return &ChannelUpdateEvent{} }
func (e ChannelUpdateEvent) Event() EventType        { return EventChannelUpdate }
