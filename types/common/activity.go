package common

type ActivityLocation struct {
	ID                string     `json:"id"`
	Kind              string     `json:"kind"`
	ChannelID         Snowflake  `json:"channel_id"`
	GuildID           *Snowflake `json:"guild_id,omitempty"`
	CurrentUsersCount int        `json:"current_users_count"`
}

type ActivityInstance struct {
	ApplicationID Snowflake        `json:"application_id"`
	InstanceID    string           `json:"instance_id"`
	LaunchID      string           `json:"launch_id"`
	Location      ActivityLocation `json:"location"`
	Users         []Snowflake      `json:"users"`
}
