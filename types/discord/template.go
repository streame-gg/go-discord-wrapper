package discord

import (
	"time"
)

// GuildTemplate represents a code that when used, creates a guild based on a snapshot of an existing guild.
type GuildTemplate struct {
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	Description           *string   `json:"description,omitempty"`
	UsageCount            int       `json:"usage_count"`
	CreatorID             Snowflake `json:"creator_id"`
	Creator               User      `json:"creator"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	SourceGuildID         Snowflake `json:"source_guild_id"`
	SerializedSourceGuild Guild     `json:"serialized_source_guild"`
	IsDirty               *bool     `json:"is_dirty,omitempty"`
}
