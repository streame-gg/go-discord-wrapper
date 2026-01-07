package types

type HelloPayloadData struct {
	HeartbeatInterval float64 `json:"heartbeat_interval"`
}

type DiscordMessageCreateEvent struct {
	DiscordMessage
	GuildID  *string        `json:"guild_id"`
	Member   *GuildMember   `json:"member,omitempty"`
	Mentions *[]DiscordUser `json:"mentions"`
}
