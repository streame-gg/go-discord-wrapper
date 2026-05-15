package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildCreateEvent struct {
	discord.GatewayGuildWrapper
	Large       bool  `json:"large"`
	Unavailable *bool `json:"unavailable"`
	MemberCount int   `json:"member_count"`
}

func init() { RegisterEvent(GuildCreateEvent{}) }

func (e GuildCreateEvent) DesiredEventType() Event {
	return &GuildCreateEvent{}
}

func (e GuildCreateEvent) Event() EventType {
	return EventGuildCreate
}
