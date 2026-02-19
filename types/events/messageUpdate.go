package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type MessageUpdateEvent struct {
	common.Message
	GuildID  *common.Snowflake   `json:"guild_id"`
	Member   *common.GuildMember `json:"member,omitempty"`
	Mentions *[]common.User      `json:"mentions"`
}

func (e MessageUpdateEvent) DesiredEventType() Event {
	return &MessageUpdateEvent{}
}

func (e MessageUpdateEvent) Event() EventType {
	return EventMessageUpdate
}
