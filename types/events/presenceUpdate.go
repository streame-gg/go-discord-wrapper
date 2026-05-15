package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type PresenceUpdateEvent struct {
	User         discord.PartialPresenceUser `json:"user"`
	GuildID      discord.Snowflake           `json:"guild_id"`
	Status       discord.PresenceStatus      `json:"status"`
	Activities   []discord.FullActivity      `json:"activities"`
	ClientStatus discord.ClientStatus        `json:"client_status"`
}

func init() { RegisterEvent(PresenceUpdateEvent{}) }

func (e PresenceUpdateEvent) DesiredEventType() Event { return &PresenceUpdateEvent{} }
func (e PresenceUpdateEvent) Event() EventType        { return EventPresenceUpdate }
