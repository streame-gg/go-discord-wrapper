package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type WebhooksUpdateEvent struct {
	GuildID   common.Snowflake `json:"guild_id"`
	ChannelID common.Snowflake `json:"channel_id"`
}

func init() { RegisterEvent(WebhooksUpdateEvent{}) }

func (e WebhooksUpdateEvent) DesiredEventType() Event { return &WebhooksUpdateEvent{} }
func (e WebhooksUpdateEvent) Event() EventType        { return EventWebhooksUpdate }
