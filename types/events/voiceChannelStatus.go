package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &VoiceChannelStatusUpdateEvent{}
var _ Event = &VoiceChannelStartTimeUpdateEvent{}

func init() {
	RegisterEvent(&VoiceChannelStatusUpdateEvent{})
	RegisterEvent(&VoiceChannelStartTimeUpdateEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#voice-channel-status-update
type VoiceChannelStatusUpdateEvent struct {
	ChannelID discord.Snowflake `json:"id"`
	GuildID   discord.Snowflake `json:"guild_id"`
	Status    *string           `json:"status"`

	Guild   *discord.Guild   `json:"guild"`
	Channel *discord.Channel `json:"channel"`
}

// https://docs.discord.com/developers/events/gateway-events#voice-channel-start-time-update
type VoiceChannelStartTimeUpdateEvent struct {
	ChannelID      discord.Snowflake `json:"id"`
	GuildID        discord.Snowflake `json:"guild_id"`
	VoiceStartTime int64             `json:"voice_start_time"`

	Guild   *discord.Guild   `json:"guild"`
	Channel *discord.Channel `json:"channel"`
}

func (e *VoiceChannelStatusUpdateEvent) DesiredEventType() Event {
	return &VoiceChannelStatusUpdateEvent{}
}
func (e *VoiceChannelStatusUpdateEvent) Event() EventType { return EventVoiceChannelStatusUpdate }

func (e *VoiceChannelStartTimeUpdateEvent) DesiredEventType() Event {
	return &VoiceChannelStartTimeUpdateEvent{}
}
func (e *VoiceChannelStartTimeUpdateEvent) Event() EventType { return EventVoiceChannelStartTimeUpdate }
