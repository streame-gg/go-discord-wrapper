package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type VoiceStateUpdateEvent struct {
	common.VoiceState
}

// VoiceServerUpdateEvent is dispatched when a guild's voice server is updated.
// Use the Token and Endpoint to connect to the voice websocket.
type VoiceServerUpdateEvent struct {
	Token    string           `json:"token"`
	GuildID  common.Snowflake `json:"guild_id"`
	Endpoint *string          `json:"endpoint"` // null when the guild has no voice server
}

// VoiceChannelEffectSendEvent is dispatched when a user sends an effect (soundboard,
// emoji, etc.) in a voice channel.
type VoiceChannelEffectSendEvent struct {
	ChannelID     common.Snowflake  `json:"channel_id"`
	GuildID       common.Snowflake  `json:"guild_id"`
	UserID        common.Snowflake  `json:"user_id"`
	Emoji         *common.Emoji     `json:"emoji,omitempty"`
	AnimationType *int              `json:"animation_type,omitempty"`
	AnimationID   *int              `json:"animation_id,omitempty"`
	SoundID       *common.Snowflake `json:"sound_id,omitempty"`
	SoundVolume   *float64          `json:"sound_volume,omitempty"`
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
