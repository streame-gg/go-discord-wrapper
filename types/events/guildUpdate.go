package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildUpdateEvent{}

func init() { RegisterEvent(&GuildUpdateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#guild-update
type GuildUpdateEvent struct {
	NewGuild discord.Guild  `json:"-"`
	OldGuild *discord.Guild `json:"-"`
}

func (e *GuildUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewGuild)
}
func (e *GuildUpdateEvent) DesiredEventType() Event { return &GuildUpdateEvent{} }
func (e *GuildUpdateEvent) Event() EventType        { return EventGuildUpdate }
