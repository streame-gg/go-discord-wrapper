package events

import (
	"encoding/json"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-create
type GuildCreateEvent struct {
	Unavailable bool              `json:"-"`
	ID          discord.Snowflake `json:"-"`

	JoinedAt             time.Time                     `json:"-"`
	Large                bool                          `json:"-"`
	MemberCount          int                           `json:"-"`
	VoiceStates          []discord.VoiceState          `json:"-"`
	Members              []discord.GuildMember         `json:"-"`
	Channels             []discord.Channel             `json:"-"`
	Threads              []discord.Channel             `json:"-"`
	Presences            []discord.Presence            `json:"-"`
	StageInstances       []discord.StageInstance       `json:"-"`
	GuildScheduledEvents []discord.GuildScheduledEvent `json:"-"`
	SoundboardSounds     []discord.SoundboardSound     `json:"-"`

	Guild *discord.Guild `json:"-"`
}

// UnmarshalJSON overrides the promoted GatewayGuildWrapper.UnmarshalJSON so
// that the outer GUILD_CREATE-only fields (large, unavailable, member_count)
// are also decoded. Without this override the embedded method consumes the
// entire payload and the outer fields are never populated.
func (e *GuildCreateEvent) UnmarshalJSON(data []byte) error {
	var outer struct {
		Unavailable *bool             `json:"unavailable"`
		ID          discord.Snowflake `json:"id"`

		JoinedAt             time.Time                     `json:"joined_at"`
		Large                bool                          `json:"large"`
		MemberCount          int                           `json:"member_count"`
		VoiceStates          []discord.VoiceState          `json:"voice_states"`
		Members              []discord.GuildMember         `json:"members"`
		Channels             []discord.Channel             `json:"channels"`
		Threads              []discord.Channel             `json:"threads"`
		Presences            []discord.Presence            `json:"presences"`
		StageInstances       []discord.StageInstance       `json:"stage_instances"`
		GuildScheduledEvents []discord.GuildScheduledEvent `json:"guild_scheduled_events"`
		SoundboardSounds     []discord.SoundboardSound     `json:"soundboard_sounds"`
	}
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}

	if outer.Unavailable != nil && *outer.Unavailable {
		e.Unavailable = true
		e.ID = outer.ID

		return nil
	}

	e.Unavailable = false
	e.ID = outer.ID
	e.JoinedAt = outer.JoinedAt
	e.Large = outer.Large
	e.MemberCount = outer.MemberCount
	e.VoiceStates = outer.VoiceStates
	e.Members = outer.Members
	e.Channels = outer.Channels
	e.Threads = outer.Threads
	e.Presences = outer.Presences
	e.StageInstances = outer.StageInstances
	e.GuildScheduledEvents = outer.GuildScheduledEvents
	e.SoundboardSounds = outer.SoundboardSounds

	var g discord.Guild
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
	e.Guild = &g

	return nil
}

func init() { RegisterEvent(GuildCreateEvent{}) }

func (e GuildCreateEvent) DesiredEventType() Event {
	return &GuildCreateEvent{}
}

func (e GuildCreateEvent) Event() EventType {
	return EventGuildCreate
}
