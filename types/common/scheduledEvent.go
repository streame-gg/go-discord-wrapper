package common

import "time"

type GuildScheduledEventPrivacyLevel int

const GuildScheduledEventPrivacyLevelGuildOnly GuildScheduledEventPrivacyLevel = 2

type GuildScheduledEventStatus int

const (
	GuildScheduledEventStatusScheduled GuildScheduledEventStatus = 1
	GuildScheduledEventStatusActive    GuildScheduledEventStatus = 2
	GuildScheduledEventStatusCompleted GuildScheduledEventStatus = 3
	GuildScheduledEventStatusCanceled  GuildScheduledEventStatus = 4
)

type GuildScheduledEventEntityType int

const (
	GuildScheduledEventEntityTypeStageInstance GuildScheduledEventEntityType = 1
	GuildScheduledEventEntityTypeVoice         GuildScheduledEventEntityType = 2
	GuildScheduledEventEntityTypeExternal      GuildScheduledEventEntityType = 3
)

type GuildScheduledEventEntityMetadata struct {
	Location *string `json:"location,omitempty"`
}

type GuildScheduledEvent struct {
	ID                 Snowflake                          `json:"id"`
	GuildID            Snowflake                          `json:"guild_id"`
	ChannelID          *Snowflake                         `json:"channel_id,omitempty"`
	CreatorID          *Snowflake                         `json:"creator_id,omitempty"`
	Name               string                             `json:"name"`
	Description        *string                            `json:"description,omitempty"`
	ScheduledStartTime time.Time                          `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                         `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       GuildScheduledEventPrivacyLevel    `json:"privacy_level"`
	Status             GuildScheduledEventStatus          `json:"status"`
	EntityType         GuildScheduledEventEntityType      `json:"entity_type"`
	EntityID           *Snowflake                         `json:"entity_id,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Creator            *User                              `json:"creator,omitempty"`
	UserCount          *int                               `json:"user_count,omitempty"`
	Image              *string                            `json:"image,omitempty"`
}
