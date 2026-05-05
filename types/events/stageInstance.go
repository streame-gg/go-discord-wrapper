package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type StageInstanceCreateEvent struct {
	common.StageInstance
}

type StageInstanceUpdateEvent struct {
	common.StageInstance
}

type StageInstanceDeleteEvent struct {
	common.StageInstance
}

func init() {
	RegisterEvent(StageInstanceCreateEvent{})
	RegisterEvent(StageInstanceUpdateEvent{})
	RegisterEvent(StageInstanceDeleteEvent{})
}

func (e StageInstanceCreateEvent) DesiredEventType() Event { return &StageInstanceCreateEvent{} }
func (e StageInstanceCreateEvent) Event() EventType        { return EventStageInstanceCreate }

func (e StageInstanceUpdateEvent) DesiredEventType() Event { return &StageInstanceUpdateEvent{} }
func (e StageInstanceUpdateEvent) Event() EventType        { return EventStageInstanceUpdate }

func (e StageInstanceDeleteEvent) DesiredEventType() Event { return &StageInstanceDeleteEvent{} }
func (e StageInstanceDeleteEvent) Event() EventType        { return EventStageInstanceDelete }
