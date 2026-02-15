package events

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ReadyEvent struct {
	User             common.User              `json:"user"`
	SessionID        string                   `json:"session_id"`
	ResumeGatewayURL string                   `json:"resume_gateway_url"`
	Shard            []int                    `json:"shard,omitempty"`
	Guilds           []common.AnyGuildWrapper `json:"guilds"`
}

func (e ReadyEvent) DesiredEventType() Event {
	return &ReadyEvent{}
}

func (e ReadyEvent) Event() EventType {
	return EventReady
}
