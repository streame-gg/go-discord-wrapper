package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// WebhookCreateEvent fires when a webhook is added to a channel.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhookCreateEvent struct {
	GuildID   discord.Snowflake
	ChannelID discord.Snowflake
	Webhook   *discord.Webhook
}

// WebhookDeleteEvent fires when a webhook is removed from a channel.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhookDeleteEvent struct {
	GuildID   discord.Snowflake
	ChannelID discord.Snowflake
	Webhook   *discord.Webhook
}

// WebhookUpdateEvent fires when a webhook's name, avatar, or channel changes.
// OldWebhook is the previously cached version; NewWebhook is the current one.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhookUpdateEvent struct {
	GuildID    discord.Snowflake
	ChannelID  discord.Snowflake
	OldWebhook *discord.Webhook
	NewWebhook *discord.Webhook
}

func (e WebhookCreateEvent) Event() EventType        { return EventWrapperWebhookCreate }
func (e WebhookCreateEvent) DesiredEventType() Event { return &WebhookCreateEvent{} }

func (e WebhookUpdateEvent) Event() EventType        { return EventWrapperWebhookUpdate }
func (e WebhookUpdateEvent) DesiredEventType() Event { return &WebhookUpdateEvent{} }

func (e WebhookDeleteEvent) Event() EventType        { return EventWrapperWebhookDelete }
func (e WebhookDeleteEvent) DesiredEventType() Event { return &WebhookDeleteEvent{} }
