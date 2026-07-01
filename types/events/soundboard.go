package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildSoundboardSoundCreateEvent{}
var _ Event = &GuildSoundboardSoundUpdateEvent{}
var _ Event = &GuildSoundboardSoundDeleteEvent{}
var _ Event = &GuildSoundboardSoundsUpdateEvent{}
var _ Event = &SoundboardSoundsEvent{}

func init() {
	RegisterEvent(&GuildSoundboardSoundCreateEvent{})
	RegisterEvent(&GuildSoundboardSoundUpdateEvent{})
	RegisterEvent(&GuildSoundboardSoundDeleteEvent{})
	RegisterEvent(&GuildSoundboardSoundsUpdateEvent{})
	RegisterEvent(&SoundboardSoundsEvent{})
}

// GuildSoundboardSoundCreateEvent is dispatched when a soundboard sound is created in a guild.
// https://docs.discord.com/developers/events/gateway-events#guild-soundboard-sound-create
type GuildSoundboardSoundCreateEvent struct {
	discord.SoundboardSound
}

// GuildSoundboardSoundUpdateEvent is dispatched when a soundboard sound is updated.
// https://docs.discord.com/developers/events/gateway-events#guild-soundboard-sound-update
type GuildSoundboardSoundUpdateEvent struct {
	NewSound discord.SoundboardSound  `json:"-"`
	OldSound *discord.SoundboardSound `json:"-"`

	Guild *discord.Guild `json:"-"`
}

// GuildSoundboardSoundDeleteEvent is dispatched when a soundboard sound is deleted.
// https://docs.discord.com/developers/events/gateway-events#guild-soundboard-sound-delete
type GuildSoundboardSoundDeleteEvent struct {
	SoundID discord.Snowflake `json:"sound_id"`
	GuildID discord.Snowflake `json:"guild_id"`

	Guild *discord.Guild `json:"-"`
}

// GuildSoundboardSoundsUpdateEvent is dispatched when multiple guild soundboard sounds are updated.
// https://docs.discord.com/developers/events/gateway-events#guild-soundboard-sound-update
type GuildSoundboardSoundsUpdateEvent struct {
	GuildID             discord.Snowflake         `json:"guild_id"`
	NewSoundboardSounds []discord.SoundboardSound `json:"soundboard_sounds"`

	Guild *discord.Guild `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#soundboard-sounds
type SoundboardSoundsEvent struct {
	GuildID          discord.Snowflake         `json:"guild_id"`
	SoundboardSounds []discord.SoundboardSound `json:"soundboard_sounds"`
}

func (e *GuildSoundboardSoundCreateEvent) Event() EventType { return EventGuildSoundboardSoundCreate }
func (e *GuildSoundboardSoundCreateEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundCreateEvent{}
}

func (e *GuildSoundboardSoundUpdateEvent) Event() EventType { return EventGuildSoundboardSoundUpdate }
func (e *GuildSoundboardSoundUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewSound)
}
func (e *GuildSoundboardSoundUpdateEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundUpdateEvent{}
}

func (e *GuildSoundboardSoundDeleteEvent) Event() EventType { return EventGuildSoundboardSoundDelete }
func (e *GuildSoundboardSoundDeleteEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundDeleteEvent{}
}

func (e *GuildSoundboardSoundsUpdateEvent) Event() EventType { return EventGuildSoundboardSoundsUpdate }
func (e *GuildSoundboardSoundsUpdateEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundsUpdateEvent{}
}

func (e *SoundboardSoundsEvent) DesiredEventType() Event { return &SoundboardSoundsEvent{} }
func (e *SoundboardSoundsEvent) Event() EventType        { return EventSoundboardSounds }
