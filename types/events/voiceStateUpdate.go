package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type VoiceStateUpdateEvent struct {
	discord.VoiceState
	// OldState is the user's previous voice state. Nil if the user was not in
	// any voice channel before this event (e.g. they just joined for the first time).
	OldState *discord.VoiceState
}

// VoiceServerUpdateEvent is dispatched when a guild's voice server is updated.
// Use the Token and Endpoint to connect to the voice websocket.
type VoiceServerUpdateEvent struct {
	Token    string            `json:"token"`
	GuildID  discord.Snowflake `json:"guild_id"`
	Endpoint *string           `json:"endpoint,omitempty"` // null when the guild has no voice server
}

// VoiceChannelEffectSendEvent is dispatched when a user sends an effect (soundboard,
// emoji, etc.) in a voice channel.
type VoiceChannelEffectSendEvent struct {
	ChannelID     discord.Snowflake  `json:"channel_id"`
	GuildID       discord.Snowflake  `json:"guild_id"`
	UserID        discord.Snowflake  `json:"user_id"`
	Emoji         *discord.Emoji     `json:"emoji,omitempty"`
	AnimationType *int               `json:"animation_type,omitempty"`
	AnimationID   *int               `json:"animation_id,omitempty"`
	SoundID       *discord.Snowflake `json:"sound_id,omitempty"`
	SoundVolume   *float64           `json:"sound_volume,omitempty"`
}

func init() {
	RegisterEvent(VoiceStateUpdateEvent{})
	RegisterEvent(VoiceServerUpdateEvent{})
	RegisterEvent(VoiceChannelEffectSendEvent{})
}

func (e VoiceStateUpdateEvent) DesiredEventType() Event { return &VoiceStateUpdateEvent{} }
func (e VoiceStateUpdateEvent) Event() EventType        { return EventVoiceStateUpdate }

func (e VoiceServerUpdateEvent) DesiredEventType() Event { return &VoiceServerUpdateEvent{} }
func (e VoiceServerUpdateEvent) Event() EventType        { return EventVoiceServerUpdate }

func (e VoiceChannelEffectSendEvent) DesiredEventType() Event { return &VoiceChannelEffectSendEvent{} }
func (e VoiceChannelEffectSendEvent) Event() EventType        { return EventVoiceChannelEffectSend }
