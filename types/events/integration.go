package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type IntegrationCreateEvent struct {
	discord.Integration
	GuildID discord.Snowflake `json:"guild_id"`
}

type IntegrationUpdateEvent struct {
	GuildID        discord.Snowflake   `json:"-"`
	NewIntegration discord.Integration `json:"-"`
}

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

func (e IntegrationUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.Integration
		GuildID discord.Snowflake `json:"guild_id"`
	}
	return json.Marshal(wire{e.NewIntegration, e.GuildID})
}

type IntegrationDeleteEvent struct {
	ID            discord.Snowflake  `json:"id"`
	GuildID       discord.Snowflake  `json:"guild_id"`
	ApplicationID *discord.Snowflake `json:"application_id,omitempty"`
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
