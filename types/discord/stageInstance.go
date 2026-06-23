package discord

// https://docs.discord.com/developers/resources/stage-instance#stage-instance-object-privacy-level
type StageInstancePrivacyLevel int

const (
	StageInstancePrivacyLevelPublic    StageInstancePrivacyLevel = 1
	StageInstancePrivacyLevelGuildOnly StageInstancePrivacyLevel = 2
)

// https://docs.discord.com/developers/resources/stage-instance#stage-instance-object
type StageInstance struct {
	hClient EntityClient

	ID                    Snowflake                 `json:"id"`
	GuildID               Snowflake                 `json:"guild_id"`
	ChannelID             Snowflake                 `json:"channel_id"`
	Topic                 string                    `json:"topic"`
	PrivacyLevel          StageInstancePrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled  bool                      `json:"discoverable_disabled"`
	GuildScheduledEventID *Snowflake                `json:"guild_scheduled_event_id,omitempty"`
}
