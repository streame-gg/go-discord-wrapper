package events

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type GuildCreateEvent struct {
	common.AnyGuildWrapper
	Large       bool  `json:"large"`
	Unavailable *bool `json:"unavailable"`
	MemberCount int   `json:"member_count"`
}

func (e GuildCreateEvent) DesiredEventType() Event {
	return &GuildCreateEvent{}
}

func (e GuildCreateEvent) Event() EventType {
	return EventGuildCreate
}
