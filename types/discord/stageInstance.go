package discord

type StageInstancePrivacyLevel int

const (
	StageInstancePrivacyLevelPublic    StageInstancePrivacyLevel = 1
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

type StageInstance struct {
	ID                    Snowflake                 `json:"id"`
	GuildID               Snowflake                 `json:"guild_id"`
	ChannelID             Snowflake                 `json:"channel_id"`
	Topic                 string                    `json:"topic"`
	PrivacyLevel          StageInstancePrivacyLevel `json:"privacy_level"`
	GuildScheduledEventID *Snowflake                `json:"guild_scheduled_event_id,omitempty"`
}
