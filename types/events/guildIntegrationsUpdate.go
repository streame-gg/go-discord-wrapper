package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GuildIntegrationsUpdateEvent is sent when a guild integration is created,
// updated, or removed. Discord only includes the guild ID — it carries no
// integration data.
//
// For the actual integration that changed, handle the granular
// INTEGRATION_CREATE / INTEGRATION_UPDATE / INTEGRATION_DELETE events
// (OnIntegrationCreate / OnIntegrationUpdate / OnIntegrationDelete) instead;
// Discord dispatches them alongside this event under the same
// IntentGuildIntegrations gateway intent. Use this coarse event when you only
// need to know that *something* about a guild's integrations changed — for
// example, to invalidate a cached integration list and re-fetch with
// ListGuildIntegrations.
// https://docs.discord.com/developers/events/gateway-events#guild-integrations-update
type GuildIntegrationsUpdateEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
}

func init() { RegisterEvent(GuildIntegrationsUpdateEvent{}) }

func (e GuildIntegrationsUpdateEvent) DesiredEventType() Event {
	return &GuildIntegrationsUpdateEvent{}
}
func (e GuildIntegrationsUpdateEvent) Event() EventType { return EventGuildIntegrationsUpdate }
