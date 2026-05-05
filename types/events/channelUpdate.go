package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type ChannelUpdateEvent struct {
	common.Channel
}

func init() { RegisterEvent(ChannelUpdateEvent{}) }

func (e ChannelUpdateEvent) DesiredEventType() Event { return &ChannelUpdateEvent{} }
func (e ChannelUpdateEvent) Event() EventType        { return EventChannelUpdate }
