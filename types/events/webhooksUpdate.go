package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhooksUpdateEvent struct {
	GuildID   discord.Snowflake `json:"guild_id"`
	ChannelID discord.Snowflake `json:"channel_id"`
}

func init() { RegisterEvent(WebhooksUpdateEvent{}) }

func (e WebhooksUpdateEvent) DesiredEventType() Event { return &WebhooksUpdateEvent{} }
func (e WebhooksUpdateEvent) Event() EventType        { return EventWebhooksUpdate }
