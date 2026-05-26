package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildCreateEvent struct {
	discord.GatewayGuildWrapper
	Large       bool  `json:"large"`
	Unavailable *bool `json:"unavailable"`
	MemberCount int   `json:"member_count"`
}

// UnmarshalJSON overrides the promoted GatewayGuildWrapper.UnmarshalJSON so
// that the outer GUILD_CREATE-only fields (large, unavailable, member_count)
// are also decoded. Without this override the embedded method consumes the
// entire payload and the outer fields are never populated.
func (e *GuildCreateEvent) UnmarshalJSON(data []byte) error {
	if err := e.GatewayGuildWrapper.UnmarshalJSON(data); err != nil {
		return err
	}
	var outer struct {
		Large       bool  `json:"large"`
		Unavailable *bool `json:"unavailable"`
		MemberCount int   `json:"member_count"`
	}
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}
	e.Large = outer.Large
	e.Unavailable = outer.Unavailable
	e.MemberCount = outer.MemberCount
	return nil
}

func init() { RegisterEvent(GuildCreateEvent{}) }

func (e GuildCreateEvent) DesiredEventType() Event {
	return &GuildCreateEvent{}
}

func (e GuildCreateEvent) Event() EventType {
	return EventGuildCreate
}
