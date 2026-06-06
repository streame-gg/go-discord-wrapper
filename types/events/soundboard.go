package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

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
}

func (e *GuildSoundboardSoundUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewSound)
}

func (e GuildSoundboardSoundUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.SoundboardSound
		OldSound *discord.SoundboardSound `json:"old_sound,omitempty"`
	}
	return json.Marshal(wire{e.NewSound, e.OldSound})
}

// GuildSoundboardSoundDeleteEvent is dispatched when a soundboard sound is deleted.
// https://docs.discord.com/developers/events/gateway-events#guild-soundboard-sound-delete
type GuildSoundboardSoundDeleteEvent struct {
	SoundID discord.Snowflake `json:"sound_id"`
	GuildID discord.Snowflake `json:"guild_id"`
}

// GuildSoundboardSoundsUpdateEvent is dispatched when multiple guild soundboard sounds are updated.
// https://docs.discord.com/developers/events/gateway-events#guild-soundboard-sound-update
type GuildSoundboardSoundsUpdateEvent struct {
	GuildID             discord.Snowflake          `json:"-"`
	NewSoundboardSounds []*discord.SoundboardSound `json:"-"`
}

func (e *GuildSoundboardSoundsUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID          discord.Snowflake          `json:"guild_id"`
		SoundboardSounds []*discord.SoundboardSound `json:"soundboard_sounds"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewSoundboardSounds = w.SoundboardSounds
	return nil
}

func (e GuildSoundboardSoundsUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		GuildID          discord.Snowflake          `json:"guild_id"`
		SoundboardSounds []*discord.SoundboardSound `json:"soundboard_sounds"`
	}
	return json.Marshal(wire{e.GuildID, e.NewSoundboardSounds})
}

func init() {
	RegisterEvent(GuildSoundboardSoundCreateEvent{})
	RegisterEvent(GuildSoundboardSoundUpdateEvent{})
	RegisterEvent(GuildSoundboardSoundDeleteEvent{})
	RegisterEvent(GuildSoundboardSoundsUpdateEvent{})
}

func (e GuildSoundboardSoundCreateEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundCreateEvent{}
}
func (e GuildSoundboardSoundCreateEvent) Event() EventType { return EventGuildSoundboardSoundCreate }

func (e GuildSoundboardSoundUpdateEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundUpdateEvent{}
}
func (e GuildSoundboardSoundUpdateEvent) Event() EventType { return EventGuildSoundboardSoundUpdate }

func (e GuildSoundboardSoundDeleteEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundDeleteEvent{}
}
func (e GuildSoundboardSoundDeleteEvent) Event() EventType { return EventGuildSoundboardSoundDelete }

func (e GuildSoundboardSoundsUpdateEvent) DesiredEventType() Event {
	return &GuildSoundboardSoundsUpdateEvent{}
}
func (e GuildSoundboardSoundsUpdateEvent) Event() EventType { return EventGuildSoundboardSoundsUpdate }
