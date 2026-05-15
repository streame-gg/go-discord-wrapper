package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildDeleteEvent struct {
	discord.UnavailableGuild
}

func init() { RegisterEvent(GuildDeleteEvent{}) }

func (g GuildDeleteEvent) DesiredEventType() Event {
	return &GuildDeleteEvent{}
}

func (g GuildDeleteEvent) Event() EventType {
	return EventGuildDelete
}
