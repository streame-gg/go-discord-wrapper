package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// Webhook synthetic events are derived by go-discord-wrapper, not sent by
// Discord. Discord only emits the coarse WEBHOOKS_UPDATE (guild_id + channel_id);
// the wrapper fetches the channel's webhooks and diffs them against the previous
// snapshot to produce these granular Create/Update/Delete events. They require a
// webhook handler to be registered (the fetch is skipped otherwise) and, like the
// emoji/sticker synthetic events, fire nothing on a cold snapshot (the first
// WEBHOOKS_UPDATE for a channel only establishes the baseline).

// WebhookCreateEvent fires when a webhook is added to a channel.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhookCreateEvent struct {
	GuildID   discord.Snowflake
	ChannelID discord.Snowflake
	Webhook   *discord.Webhook
}

func (e WebhookCreateEvent) Event() EventType        { return EventWrapperWebhookCreate }
func (e WebhookCreateEvent) DesiredEventType() Event { return &WebhookCreateEvent{} }

// WebhookDeleteEvent fires when a webhook is removed from a channel.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhookDeleteEvent struct {
	GuildID   discord.Snowflake
	ChannelID discord.Snowflake
	Webhook   *discord.Webhook
}

func (e WebhookDeleteEvent) Event() EventType        { return EventWrapperWebhookDelete }
func (e WebhookDeleteEvent) DesiredEventType() Event { return &WebhookDeleteEvent{} }

// WebhookUpdateEvent fires when a webhook's name, avatar, or channel changes.
// OldWebhook is the previously cached version; NewWebhook is the current one.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#webhooks-update
type WebhookUpdateEvent struct {
	GuildID    discord.Snowflake
	ChannelID  discord.Snowflake
	OldWebhook *discord.Webhook
	NewWebhook *discord.Webhook
}

func (e WebhookUpdateEvent) Event() EventType        { return EventWrapperWebhookUpdate }
func (e WebhookUpdateEvent) DesiredEventType() Event { return &WebhookUpdateEvent{} }
