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

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	Location *string `json:"location,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-object
type GuildScheduledEvent struct {
	hClient EntityClient

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
	RecurrenceRule     *GuildScheduledEventRecurrenceRule `json:"recurrence_rule,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-recurrence-rule-object
type GuildScheduledEventRecurrenceRule struct {
	Start      time.Time                                   `json:"start"`
	End        *time.Time                                  `json:"end"`
	Frequency  GuildScheduledEventRecurrenceRuleFrequency  `json:"frequency"`
	Interval   int                                         `json:"interval"`
	ByWeekday  []GuildScheduledEventRecurrenceRuleWeekday  `json:"by_weekday,omitempty"`
	ByNWeekday []GuildScheduledEventRecurrenceRuleNWeekday `json:"by_n_weekday,omitempty"`
	ByMonth    []GuildScheduledEventRecurrenceRuleMonth    `json:"by_month,omitempty"`
	ByMonthDay []int                                       `json:"by_month_day,omitempty"`
	ByYearDay  []int                                       `json:"by_year_day,omitempty"`
	Count      *int                                        `json:"count,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-recurrence-rule-object-guild-scheduled-event-recurrence-rule-month
type GuildScheduledEventRecurrenceRuleMonth uint8

const (
	GuildScheduledEventRecurrenceRuleMonthJanuary   GuildScheduledEventRecurrenceRuleMonth = 1
	GuildScheduledEventRecurrenceRuleMonthFebruary  GuildScheduledEventRecurrenceRuleMonth = 2
	GuildScheduledEventRecurrenceRuleMonthMarch     GuildScheduledEventRecurrenceRuleMonth = 3
	GuildScheduledEventRecurrenceRuleMonthApril     GuildScheduledEventRecurrenceRuleMonth = 4
	GuildScheduledEventRecurrenceRuleMonthMay       GuildScheduledEventRecurrenceRuleMonth = 5
	GuildScheduledEventRecurrenceRuleMonthJune      GuildScheduledEventRecurrenceRuleMonth = 6
	GuildScheduledEventRecurrenceRuleMonthJuly      GuildScheduledEventRecurrenceRuleMonth = 7
	GuildScheduledEventRecurrenceRuleMonthAugust    GuildScheduledEventRecurrenceRuleMonth = 8
	GuildScheduledEventRecurrenceRuleMonthSeptember GuildScheduledEventRecurrenceRuleMonth = 9
	GuildScheduledEventRecurrenceRuleMonthOctober   GuildScheduledEventRecurrenceRuleMonth = 10
	GuildScheduledEventRecurrenceRuleMonthNovember  GuildScheduledEventRecurrenceRuleMonth = 11
	GuildScheduledEventRecurrenceRuleMonthDecember  GuildScheduledEventRecurrenceRuleMonth = 12
)

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-recurrence-rule-object-guild-scheduled-event-recurrence-rule-nweekday-structure
type GuildScheduledEventRecurrenceRuleNWeekday struct {
	N   int                                      `json:"n"`
	Day GuildScheduledEventRecurrenceRuleWeekday `json:"day"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-recurrence-rule-object-guild-scheduled-event-recurrence-rule-weekday
type GuildScheduledEventRecurrenceRuleWeekday uint8

const (
	GuildScheduledEventRecurrenceRuleWeekdayMonday    GuildScheduledEventRecurrenceRuleWeekday = 0
	GuildScheduledEventRecurrenceRuleWeekdayTuesday   GuildScheduledEventRecurrenceRuleWeekday = 1
	GuildScheduledEventRecurrenceRuleWeekdayWednesday GuildScheduledEventRecurrenceRuleWeekday = 2
	GuildScheduledEventRecurrenceRuleWeekdayThursday  GuildScheduledEventRecurrenceRuleWeekday = 3
	GuildScheduledEventRecurrenceRuleWeekdayFriday    GuildScheduledEventRecurrenceRuleWeekday = 4
	GuildScheduledEventRecurrenceRuleWeekdaySaturday  GuildScheduledEventRecurrenceRuleWeekday = 5
	GuildScheduledEventRecurrenceRuleWeekdaySunday    GuildScheduledEventRecurrenceRuleWeekday = 6
)

// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-recurrence-rule-object-guild-scheduled-event-recurrence-rule-frequency
type GuildScheduledEventRecurrenceRuleFrequency uint8

const (
	GuildScheduledEventRecurrenceRuleFrequencyYearly  GuildScheduledEventRecurrenceRuleFrequency = 0
	GuildScheduledEventRecurrenceRuleFrequencyMonthly GuildScheduledEventRecurrenceRuleFrequency = 1
	GuildScheduledEventRecurrenceRuleFrequencyWeekly  GuildScheduledEventRecurrenceRuleFrequency = 2
	GuildScheduledEventRecurrenceRuleFrequencyDaily   GuildScheduledEventRecurrenceRuleFrequency = 3
)

// GuildScheduledEventUser is an entry in the list returned by ListGuildScheduledEventUsers.
//
// https://docs.discord.com/developers/resources/guild-scheduled-event#guild-scheduled-event-user-object
type GuildScheduledEventUser struct {
	GuildScheduledEventID Snowflake    `json:"guild_scheduled_event_id"`
	User                  User         `json:"user"`
	Member                *GuildMember `json:"member,omitempty"`
}
