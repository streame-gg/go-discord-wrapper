package events

import (
	"encoding/json"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildCreateEvent{}
var _ Event = &GuildDeleteEvent{}

func init() {
	RegisterEvent(&GuildCreateEvent{})
	RegisterEvent(&GuildDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#guild-create
type GuildCreateEvent struct {
	Unavailable bool              `json:"unavailable"`
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

	Guild *discord.Guild `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-delete
type GuildDeleteEvent struct {
	Unavailable bool              `json:"unavailable"`
	ID          discord.Snowflake `json:"id"`

	Guild *discord.Guild `json:"-"`
}

func (e *GuildCreateEvent) DesiredEventType() Event {
	return &GuildCreateEvent{}
}
func (e *GuildCreateEvent) Event() EventType {
	return EventGuildCreate
}
func (e *GuildCreateEvent) UnmarshalJSON(data []byte) error {
	var wire struct {
		Unavailable bool `json:"unavailable"`

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

		discord.Guild
	}

	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	e.ID = wire.ID
	e.Unavailable = wire.Unavailable

	if wire.Unavailable {
		return nil
	}

	e.JoinedAt = wire.JoinedAt
	e.Large = wire.Large
	e.MemberCount = wire.MemberCount
	e.VoiceStates = wire.VoiceStates
	e.Members = wire.Members
	e.Channels = wire.Channels
	e.Threads = wire.Threads
	e.Presences = wire.Presences
	e.StageInstances = wire.StageInstances
	e.GuildScheduledEvents = wire.GuildScheduledEvents
	e.SoundboardSounds = wire.SoundboardSounds

	e.Guild = &wire.Guild

	return nil
}

func (g *GuildDeleteEvent) DesiredEventType() Event {
	return &GuildDeleteEvent{}
}
func (g *GuildDeleteEvent) Event() EventType {
	return EventGuildDelete
}
