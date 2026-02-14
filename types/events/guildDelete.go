package events

import (
	"go-discord-wrapper/types/common"
)

type GuildDeleteEvent struct {
	common.UnavailableGuild
}

func (g GuildDeleteEvent) DesiredEventType() Event {
	return &GuildDeleteEvent{}
}

func (g GuildDeleteEvent) Event() EventType {
	return EventGuildDelete
}
