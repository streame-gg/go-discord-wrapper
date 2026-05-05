package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type ChannelCreateEvent struct {
	common.Channel
}

func init() { RegisterEvent(ChannelCreateEvent{}) }

func (e ChannelCreateEvent) DesiredEventType() Event {
	return &ChannelCreateEvent{}
}

func (e ChannelCreateEvent) Event() EventType {
	return EventChannelCreate
}
