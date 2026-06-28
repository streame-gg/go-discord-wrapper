package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &VoiceStateUpdateEvent{}
var _ Event = &VoiceServerUpdateEvent{}
var _ Event = &VoiceStateUpdateEvent{}

func init() {
	RegisterEvent(VoiceStateUpdateEvent{})
	RegisterEvent(VoiceServerUpdateEvent{})
	RegisterEvent(VoiceChannelEffectSendEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#voice-state-update
type VoiceStateUpdateEvent struct {
	NewState *discord.VoiceState
	OldState *discord.VoiceState
}

// https://docs.discord.com/developers/events/gateway-events#voice-server-update
type VoiceServerUpdateEvent struct {
	Token    string            `json:"token"`
	GuildID  discord.Snowflake `json:"guild_id"`
	Endpoint *string           `json:"endpoint"`
}

// https://docs.discord.com/developers/events/gateway-events#voice-channel-effect-send
type VoiceChannelEffectSendEvent struct {
	ChannelID     discord.Snowflake                  `json:"channel_id"`
	GuildID       discord.Snowflake                  `json:"guild_id"`
	UserID        discord.Snowflake                  `json:"user_id"`
	Emoji         *discord.Emoji                     `json:"emoji,omitempty"`
	AnimationType *discord.VoiceChannelAnimationType `json:"animation_type,omitempty"`
	AnimationID   *int                               `json:"animation_id,omitempty"`
	SoundID       *discord.Snowflake                 `json:"sound_id,omitempty"`
	SoundVolume   *float64                           `json:"sound_volume,omitempty"`

	Channel *discord.Channel         `json:"-"`
	Guild   *discord.Guild           `json:"-"`
	User    *discord.User            `json:"-"`
	Sound   *discord.SoundboardSound `json:"-"`
}

func (e VoiceStateUpdateEvent) DesiredEventType() Event { return &VoiceStateUpdateEvent{} }
func (e VoiceStateUpdateEvent) Event() EventType        { return EventVoiceStateUpdate }
func (e VoiceStateUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewState)
}

func (e VoiceServerUpdateEvent) DesiredEventType() Event { return &VoiceServerUpdateEvent{} }
func (e VoiceServerUpdateEvent) Event() EventType        { return EventVoiceServerUpdate }

func (e VoiceChannelEffectSendEvent) DesiredEventType() Event { return &VoiceChannelEffectSendEvent{} }
func (e VoiceChannelEffectSendEvent) Event() EventType        { return EventVoiceChannelEffectSend }
