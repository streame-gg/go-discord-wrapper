package discord

import (
	"time"
)

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
type GuildScheduledEventPrivacyLevel int

const GuildScheduledEventPrivacyLevelGuildOnly GuildScheduledEventPrivacyLevel = 2

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
type GuildScheduledEventStatus int

const (
	GuildScheduledEventStatusScheduled GuildScheduledEventStatus = 1
	GuildScheduledEventStatusActive    GuildScheduledEventStatus = 2
	GuildScheduledEventStatusCompleted GuildScheduledEventStatus = 3
	GuildScheduledEventStatusCanceled  GuildScheduledEventStatus = 4
)

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
type GuildScheduledEventEntityType int

const (
	GuildScheduledEventEntityTypeStageInstance GuildScheduledEventEntityType = 1
	GuildScheduledEventEntityTypeVoice         GuildScheduledEventEntityType = 2
	GuildScheduledEventEntityTypeExternal      GuildScheduledEventEntityType = 3
)

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	Location *string `json:"location,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object
type GuildScheduledEvent struct {
	hClient            EntityClient
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

// GuildScheduledEventUser is an entry in the list returned by ListGuildScheduledEventUsers.
//
// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-user-object
type GuildScheduledEventUser struct {
	GuildScheduledEventID Snowflake    `json:"guild_scheduled_event_id"`
	User                  User         `json:"user"`
	Member                *GuildMember `json:"member,omitempty"`
}
