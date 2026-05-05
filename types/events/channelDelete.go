package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type ChannelDeleteEvent struct {
	common.Channel
}

func init() { RegisterEvent(ChannelDeleteEvent{}) }

func (e ChannelDeleteEvent) DesiredEventType() Event {
	return &ChannelDeleteEvent{}
}

func (e ChannelDeleteEvent) Event() EventType {
	return EventChannelDelete
}
