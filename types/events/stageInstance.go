package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &StageInstanceCreateEvent{}
var _ Event = &StageInstanceUpdateEvent{}
var _ Event = &StageInstanceDeleteEvent{}

func init() {
	RegisterEvent(&StageInstanceCreateEvent{})
	RegisterEvent(&StageInstanceUpdateEvent{})
	RegisterEvent(&StageInstanceDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#stage-instance-create
type StageInstanceCreateEvent struct {
	discord.StageInstance
}

// https://docs.discord.com/developers/events/gateway-events#stage-instance-update
type StageInstanceUpdateEvent struct {
	NewStage discord.StageInstance  `json:"-"`
	OldStage *discord.StageInstance `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#stage-instance-delete
type StageInstanceDeleteEvent struct {
	discord.StageInstance
}

func (e *StageInstanceCreateEvent) DesiredEventType() Event { return &StageInstanceCreateEvent{} }
func (e *StageInstanceCreateEvent) Event() EventType        { return EventStageInstanceCreate }

func (e *StageInstanceUpdateEvent) DesiredEventType() Event { return &StageInstanceUpdateEvent{} }
func (e *StageInstanceUpdateEvent) Event() EventType        { return EventStageInstanceUpdate }
func (e *StageInstanceUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewStage)
}

func (e *StageInstanceDeleteEvent) DesiredEventType() Event { return &StageInstanceDeleteEvent{} }
func (e *StageInstanceDeleteEvent) Event() EventType        { return EventStageInstanceDelete }
