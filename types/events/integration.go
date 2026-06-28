package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &IntegrationCreateEvent{}
var _ Event = &IntegrationUpdateEvent{}
var _ Event = &IntegrationDeleteEvent{}

func init() {
	RegisterEvent(&IntegrationCreateEvent{})
	RegisterEvent(&IntegrationUpdateEvent{})
	RegisterEvent(&IntegrationDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#integration-create
type IntegrationCreateEvent struct {
	Integration discord.Integration
	GuildID     discord.Snowflake `json:"guild_id"`
}

// https://docs.discord.com/developers/events/gateway-events#integration-update
type IntegrationUpdateEvent struct {
	GuildID        discord.Snowflake   `json:"-"`
	NewIntegration discord.Integration `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#integration-delete
type IntegrationDeleteEvent struct {
	ID            discord.Snowflake  `json:"id"`
	GuildID       discord.Snowflake  `json:"guild_id"`
	ApplicationID *discord.Snowflake `json:"application_id,omitempty"`
}

func (e *IntegrationCreateEvent) DesiredEventType() Event { return &IntegrationCreateEvent{} }
func (e *IntegrationCreateEvent) Event() EventType        { return EventIntegrationCreate }
func (e *IntegrationCreateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		discord.Integration
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.Integration = w.Integration
	return nil
}

func (e *IntegrationUpdateEvent) DesiredEventType() Event { return &IntegrationUpdateEvent{} }
func (e *IntegrationUpdateEvent) Event() EventType        { return EventIntegrationUpdate }
func (e *IntegrationUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		discord.Integration
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewIntegration = w.Integration
	return nil
}

func (e *IntegrationDeleteEvent) DesiredEventType() Event { return &IntegrationDeleteEvent{} }
func (e *IntegrationDeleteEvent) Event() EventType        { return EventIntegrationDelete }
