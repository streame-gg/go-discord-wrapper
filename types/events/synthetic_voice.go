package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// VoiceMemberJoinEvent fires when a user enters a voice channel for the first time
// in this session (OldState was nil). Derived from VOICE_STATE_UPDATE.
type VoiceMemberJoinEvent struct {
	GuildID   discord.Snowflake
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	State     *discord.VoiceState
}

func (e VoiceMemberJoinEvent) Event() EventType        { return EventWrapperVoiceMemberJoin }
func (e VoiceMemberJoinEvent) DesiredEventType() Event { return &VoiceMemberJoinEvent{} }

// VoiceMemberLeaveEvent fires when a user leaves all voice channels (ChannelID
// becomes nil). Derived from VOICE_STATE_UPDATE.
type VoiceMemberLeaveEvent struct {
	GuildID      discord.Snowflake
	UserID       discord.Snowflake
	OldChannelID discord.Snowflake
	OldState     *discord.VoiceState
}

func (e VoiceMemberLeaveEvent) Event() EventType        { return EventWrapperVoiceMemberLeave }
func (e VoiceMemberLeaveEvent) DesiredEventType() Event { return &VoiceMemberLeaveEvent{} }

// VoiceMemberMoveEvent fires when a user switches from one voice channel to
// another. Derived from VOICE_STATE_UPDATE.
type VoiceMemberMoveEvent struct {
	GuildID      discord.Snowflake
	UserID       discord.Snowflake
	OldChannelID discord.Snowflake
	NewChannelID discord.Snowflake
	OldState     *discord.VoiceState
	NewState     *discord.VoiceState
}

func (e VoiceMemberMoveEvent) Event() EventType        { return EventWrapperVoiceMemberMove }
func (e VoiceMemberMoveEvent) DesiredEventType() Event { return &VoiceMemberMoveEvent{} }

// VoiceMemberUpdateEvent fires when a user's voice state changes while remaining
// in the same channel (e.g. mute, deafen, suppress). Derived from VOICE_STATE_UPDATE.
type VoiceMemberUpdateEvent struct {
	GuildID   discord.Snowflake
	UserID    discord.Snowflake
	ChannelID discord.Snowflake
	OldState  *discord.VoiceState
	NewState  *discord.VoiceState
}

func (e VoiceMemberUpdateEvent) Event() EventType        { return EventWrapperVoiceMemberUpdate }
func (e VoiceMemberUpdateEvent) DesiredEventType() Event { return &VoiceMemberUpdateEvent{} }
