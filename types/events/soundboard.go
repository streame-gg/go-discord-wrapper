package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

// GuildSoundboardSoundCreateEvent is dispatched when a soundboard sound is created in a guild.
type GuildSoundboardSoundCreateEvent struct {
	common.SoundboardSound
}

// GuildSoundboardSoundUpdateEvent is dispatched when a soundboard sound is updated.
type GuildSoundboardSoundUpdateEvent struct {
	common.SoundboardSound
}

// GuildSoundboardSoundDeleteEvent is dispatched when a soundboard sound is deleted.
type GuildSoundboardSoundDeleteEvent struct {
	SoundID common.Snowflake `json:"sound_id"`
	GuildID common.Snowflake `json:"guild_id"`
}

// GuildSoundboardSoundsUpdateEvent is dispatched when multiple guild soundboard sounds are updated.
type GuildSoundboardSoundsUpdateEvent struct {
	GuildID          common.Snowflake         `json:"guild_id"`
	SoundboardSounds []common.SoundboardSound `json:"soundboard_sounds"`
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
