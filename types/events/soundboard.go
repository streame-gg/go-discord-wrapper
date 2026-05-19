package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GuildSoundboardSoundCreateEvent is dispatched when a soundboard sound is created in a guild.
type GuildSoundboardSoundCreateEvent struct {
	discord.SoundboardSound
}

// GuildSoundboardSoundUpdateEvent is dispatched when a soundboard sound is updated.
type GuildSoundboardSoundUpdateEvent struct {
	discord.SoundboardSound
	OldSound *discord.SoundboardSound `json:"old_sound,omitempty"`
}

// GuildSoundboardSoundDeleteEvent is dispatched when a soundboard sound is deleted.
type GuildSoundboardSoundDeleteEvent struct {
	SoundID discord.Snowflake `json:"sound_id"`
	GuildID discord.Snowflake `json:"guild_id"`
}

// GuildSoundboardSoundsUpdateEvent is dispatched when multiple guild soundboard sounds are updated.
type GuildSoundboardSoundsUpdateEvent struct {
	GuildID          discord.Snowflake         `json:"guild_id"`
	SoundboardSounds []discord.SoundboardSound `json:"soundboard_sounds"`
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
