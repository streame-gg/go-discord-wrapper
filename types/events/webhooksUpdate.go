package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &WebhooksUpdateEvent{}

func init() { RegisterEvent(&WebhooksUpdateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhooksUpdateEvent struct {
	GuildID   discord.Snowflake `json:"guild_id"`
	ChannelID discord.Snowflake `json:"channel_id"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
}

func (e *WebhooksUpdateEvent) DesiredEventType() Event { return &WebhooksUpdateEvent{} }
func (e *WebhooksUpdateEvent) Event() EventType        { return EventWebhooksUpdate }
