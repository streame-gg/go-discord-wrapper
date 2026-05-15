package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type VoiceChannelStatusUpdateEvent struct {
	ChannelID discord.Snowflake `json:"id"`
	GuildID   discord.Snowflake `json:"guild_id"`
	Status    *string           `json:"status"`
}

type VoiceChannelStartTimeUpdateEvent struct {
	ChannelID              discord.Snowflake `json:"id"`
	GuildID                discord.Snowflake `json:"guild_id"`
	VoiceSessionsStartedAt time.Time         `json:"voice_sessions_started_at"`
}

type SoundboardSoundsEvent struct {
	GuildID          discord.Snowflake         `json:"guild_id"`
	SoundboardSounds []discord.SoundboardSound `json:"soundboard_sounds"`
}

func init() {
	RegisterEvent(VoiceChannelStatusUpdateEvent{})
	RegisterEvent(VoiceChannelStartTimeUpdateEvent{})
	RegisterEvent(SoundboardSoundsEvent{})
}

func (e VoiceChannelStatusUpdateEvent) DesiredEventType() Event {
	return &VoiceChannelStatusUpdateEvent{}
}
func (e VoiceChannelStatusUpdateEvent) Event() EventType { return EventVoiceChannelStatusUpdate }

func (e VoiceChannelStartTimeUpdateEvent) DesiredEventType() Event {
	return &VoiceChannelStartTimeUpdateEvent{}
}
func (e VoiceChannelStartTimeUpdateEvent) Event() EventType { return EventVoiceChannelStartTimeUpdate }

func (e SoundboardSoundsEvent) DesiredEventType() Event { return &SoundboardSoundsEvent{} }
func (e SoundboardSoundsEvent) Event() EventType        { return EventSoundboardSounds }
