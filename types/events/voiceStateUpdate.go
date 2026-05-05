package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type VoiceStateUpdateEvent struct {
	common.VoiceState
}

func init() { RegisterEvent(VoiceStateUpdateEvent{}) }

func (e VoiceStateUpdateEvent) DesiredEventType() Event { return &VoiceStateUpdateEvent{} }
func (e VoiceStateUpdateEvent) Event() EventType        { return EventVoiceStateUpdate }
