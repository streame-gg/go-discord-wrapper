package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type GuildAuditLogEntryCreateEvent struct {
	*common.AuditLogEntry
	GuildID common.Snowflake `json:"guild_id"`
}

func (g GuildAuditLogEntryCreateEvent) DesiredEventType() Event {
	return &GuildAuditLogEntryCreateEvent{}
}

func (g GuildAuditLogEntryCreateEvent) Event() EventType {
	return EventGuildAuditLogEntryCreate
}
