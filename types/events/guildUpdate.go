package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-update
type GuildUpdateEvent struct {
	NewGuild discord.Guild  `json:"-"`
	OldGuild *discord.Guild `json:"-"`
}

func (e *GuildUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewGuild)
}

func (e GuildUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.Guild
		OldGuild *discord.Guild `json:"old_guild,omitempty"`
	}
	return json.Marshal(wire{e.NewGuild, e.OldGuild})
}

func init() { RegisterEvent(GuildUpdateEvent{}) }

func (e GuildUpdateEvent) DesiredEventType() Event { return &GuildUpdateEvent{} }
func (e GuildUpdateEvent) Event() EventType        { return EventGuildUpdate }
