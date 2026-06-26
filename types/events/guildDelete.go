package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-delete
type GuildDeleteEvent struct {
	Unavailable bool              `json:"-"`
	ID          discord.Snowflake `json:"-"`

	Guild *discord.Guild `json:"-"`
}

func (g GuildDeleteEvent) UnmarshalJSON(bytes []byte) error {
	var outer struct {
		Unavailable bool              `json:"unavailable"`
		ID          discord.Snowflake `json:"id"`
	}

	if err := json.Unmarshal(bytes, &outer); err != nil {
		return err
	}

	g.ID = outer.ID
	g.Unavailable = outer.Unavailable

	return nil
}

func init() { RegisterEvent(GuildDeleteEvent{}) }

func (g GuildDeleteEvent) DesiredEventType() Event {
	return &GuildDeleteEvent{}
}

func (g GuildDeleteEvent) Event() EventType {
	return EventGuildDelete
}
