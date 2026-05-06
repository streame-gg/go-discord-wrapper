package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type ReadyEvent struct {
	User             common.User              `json:"user"`
	SessionID        string                   `json:"session_id"`
	ResumeGatewayURL string                   `json:"resume_gateway_url"`
	Shard            []int                    `json:"shard,omitempty"`
	Guilds           []common.AnyGuildWrapper `json:"guilds"`
}

// ResumedEvent is dispatched when a session is successfully resumed.
type ResumedEvent struct{}

func init() {
	RegisterEvent(ReadyEvent{})
	RegisterEvent(ResumedEvent{})
}

func (e ReadyEvent) DesiredEventType() Event {
	return &ReadyEvent{}
}

func (e ReadyEvent) Event() EventType { return EventReady }

func (e ResumedEvent) DesiredEventType() Event { return &ResumedEvent{} }
func (e ResumedEvent) Event() EventType        { return EventResumed }
