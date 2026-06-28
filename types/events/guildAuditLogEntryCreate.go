package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildAuditLogEntryCreateEvent{}

func init() { RegisterEvent(&GuildAuditLogEntryCreateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#guild-audit-log-entry-create
type GuildAuditLogEntryCreateEvent struct {
	AuditLogEntry *discord.AuditLogEntry `json:"-"`
	GuildID       discord.Snowflake      `json:"guild_id"`
}

func (g *GuildAuditLogEntryCreateEvent) DesiredEventType() Event {
	return &GuildAuditLogEntryCreateEvent{}
}
func (g *GuildAuditLogEntryCreateEvent) Event() EventType {
	return EventGuildAuditLogEntryCreate
}
func (g *GuildAuditLogEntryCreateEvent) UnmarshalJSON(byte []byte) error {
	var wire struct {
		discord.AuditLogEntry
		GuildID discord.Snowflake `json:"guild_id"`
	}

	if err := json.Unmarshal(byte, &wire); err != nil {
		return err
	}

	g.AuditLogEntry = &wire.AuditLogEntry
	g.GuildID = wire.GuildID

	return nil
}
