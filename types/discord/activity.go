// Package discord contains Discord data types shared across the gateway,
// REST API, interactions, and cache packages.
package discord

// https://docs.discord.com/developers/resources/application#activity-location-object
type ActivityLocation struct {
	ID                string     `json:"id"`
	Kind              string     `json:"kind"`
	ChannelID         Snowflake  `json:"channel_id"`
	GuildID           *Snowflake `json:"guild_id,omitempty"`
	CurrentUsersCount int        `json:"current_users_count"`
}

// https://docs.discord.com/developers/resources/application#activity-instance-object
type ActivityInstance struct {
	ApplicationID Snowflake        `json:"application_id"`
	InstanceID    string           `json:"instance_id"`
	LaunchID      string           `json:"launch_id"`
	Location      ActivityLocation `json:"location"`
	Users         []Snowflake      `json:"users"`
}
