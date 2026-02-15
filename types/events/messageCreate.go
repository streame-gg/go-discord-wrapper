package events

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type MessageCreateEvent struct {
	common.Message
	GuildID  *string             `json:"guild_id"`
	Member   *common.GuildMember `json:"member,omitempty"`
	Mentions *[]common.User      `json:"mentions"`
}

func (e MessageCreateEvent) DesiredEventType() Event {
	return &MessageCreateEvent{}
}

func (e MessageCreateEvent) Event() EventType {
	return EventMessageCreate
}
