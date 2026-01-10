package types

import (
	"encoding/json"
)

type DiscordMessageCreateEvent struct {
	DiscordMessage
	GuildID  *string        `json:"guild_id"`
	Member   *GuildMember   `json:"member,omitempty"`
	Mentions *[]DiscordUser `json:"mentions"`
}

func (e DiscordMessageCreateEvent) Unmarshal(data []byte) (DiscordEvent, error) {
	var event DiscordMessageCreateEvent
	err := json.Unmarshal(data, &event)
	return event, err
}

type DiscordReadyEvent struct {
	User   DiscordUser       `json:"user"`
	Guilds []AnyGuildWrapper `json:"guilds"`
}

func (e DiscordReadyEvent) Unmarshal(data []byte) (DiscordEvent, error) {
	var event DiscordReadyEvent
	err := json.Unmarshal(data, &event)
	return event, err
}

type DiscordGuildCreateEvent struct {
	AnyGuildWrapper
	Large       bool  `json:"large"`
	Unavailable *bool `json:"unavailable"`
	MemberCount int   `json:"member_count"`
}
