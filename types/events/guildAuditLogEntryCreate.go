package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-audit-log-entry-create
type GuildAuditLogEntryCreateEvent struct {
	*discord.AuditLogEntry
	GuildID discord.Snowflake `json:"guild_id"`
}

func init() { RegisterEvent(GuildAuditLogEntryCreateEvent{}) }

func (g GuildAuditLogEntryCreateEvent) DesiredEventType() Event {
	return &GuildAuditLogEntryCreateEvent{}
}

func (g GuildAuditLogEntryCreateEvent) Event() EventType {
	return EventGuildAuditLogEntryCreate
}
