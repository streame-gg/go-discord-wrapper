package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type IntegrationCreateEvent struct {
	common.Integration
	GuildID common.Snowflake `json:"guild_id"`
}

type IntegrationUpdateEvent struct {
	common.Integration
	GuildID common.Snowflake `json:"guild_id"`
}

type IntegrationDeleteEvent struct {
	ID            common.Snowflake  `json:"id"`
	GuildID       common.Snowflake  `json:"guild_id"`
	ApplicationID *common.Snowflake `json:"application_id,omitempty"`
}

func init() {
	RegisterEvent(IntegrationCreateEvent{})
	RegisterEvent(IntegrationUpdateEvent{})
	RegisterEvent(IntegrationDeleteEvent{})
}

func (e IntegrationCreateEvent) DesiredEventType() Event { return &IntegrationCreateEvent{} }
func (e IntegrationCreateEvent) Event() EventType        { return EventIntegrationCreate }

func (e IntegrationUpdateEvent) DesiredEventType() Event { return &IntegrationUpdateEvent{} }
func (e IntegrationUpdateEvent) Event() EventType        { return EventIntegrationUpdate }

func (e IntegrationDeleteEvent) DesiredEventType() Event { return &IntegrationDeleteEvent{} }
func (e IntegrationDeleteEvent) Event() EventType        { return EventIntegrationDelete }
