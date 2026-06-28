package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildIntegrationsUpdateEvent{}

func init() { RegisterEvent(&GuildIntegrationsUpdateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#guild-integrations-update
type GuildIntegrationsUpdateEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
}

func (e *GuildIntegrationsUpdateEvent) DesiredEventType() Event {
	return &GuildIntegrationsUpdateEvent{}
}
func (e *GuildIntegrationsUpdateEvent) Event() EventType { return EventGuildIntegrationsUpdate }
