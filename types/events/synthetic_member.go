package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GuildMemberRoleAddEvent fires when a role is added to a member.
// One event fires per added role in a single GUILD_MEMBER_UPDATE.
// Derived from GUILD_MEMBER_UPDATE; requires cache to be enabled.
type GuildMemberRoleAddEvent struct {
	GuildID   discord.Snowflake
	UserID    discord.Snowflake
	RoleID    discord.Snowflake
	OldMember *discord.GuildMember
	NewMember *discord.GuildMember
}

func (e GuildMemberRoleAddEvent) Event() EventType        { return EventWrapperGuildMemberRoleAdd }
func (e GuildMemberRoleAddEvent) DesiredEventType() Event { return &GuildMemberRoleAddEvent{} }

// GuildMemberRoleRemoveEvent fires when a role is removed from a member.
// One event fires per removed role in a single GUILD_MEMBER_UPDATE.
type GuildMemberRoleRemoveEvent struct {
	GuildID   discord.Snowflake
	UserID    discord.Snowflake
	RoleID    discord.Snowflake
	OldMember *discord.GuildMember
	NewMember *discord.GuildMember
}

func (e GuildMemberRoleRemoveEvent) Event() EventType        { return EventWrapperGuildMemberRoleRemove }
func (e GuildMemberRoleRemoveEvent) DesiredEventType() Event { return &GuildMemberRoleRemoveEvent{} }

// GuildMemberNickChangeEvent fires when a member's nickname is set, changed, or cleared.
type GuildMemberNickChangeEvent struct {
	GuildID   discord.Snowflake
	UserID    discord.Snowflake
	OldNick   *string
	NewNick   *string
	OldMember *discord.GuildMember
	NewMember *discord.GuildMember
}

func (e GuildMemberNickChangeEvent) Event() EventType        { return EventWrapperGuildMemberNickChange }
func (e GuildMemberNickChangeEvent) DesiredEventType() Event { return &GuildMemberNickChangeEvent{} }

// GuildMemberTimeoutEvent fires when a communication timeout is applied or extended.
// It does not fire when a timeout is removed (CommunicationDisabledUntil cleared).
type GuildMemberTimeoutEvent struct {
	GuildID                    discord.Snowflake
	UserID                     discord.Snowflake
	CommunicationDisabledUntil time.Time
	OldMember                  *discord.GuildMember
	NewMember                  *discord.GuildMember
}

func (e GuildMemberTimeoutEvent) Event() EventType        { return EventWrapperGuildMemberTimeout }
func (e GuildMemberTimeoutEvent) DesiredEventType() Event { return &GuildMemberTimeoutEvent{} }

// GuildMemberBoostStartEvent fires when a member begins boosting the server
// (PremiumSince transitions from nil to non-nil).
type GuildMemberBoostStartEvent struct {
	GuildID      discord.Snowflake
	UserID       discord.Snowflake
	PremiumSince time.Time
	OldMember    *discord.GuildMember
	NewMember    *discord.GuildMember
}

func (e GuildMemberBoostStartEvent) Event() EventType        { return EventWrapperGuildMemberBoostStart }
func (e GuildMemberBoostStartEvent) DesiredEventType() Event { return &GuildMemberBoostStartEvent{} }

// GuildMemberBoostEndEvent fires when a member stops boosting the server
// (PremiumSince transitions from non-nil to nil).
type GuildMemberBoostEndEvent struct {
	GuildID   discord.Snowflake
	UserID    discord.Snowflake
	OldMember *discord.GuildMember
	NewMember *discord.GuildMember
}

func (e GuildMemberBoostEndEvent) Event() EventType        { return EventWrapperGuildMemberBoostEnd }
func (e GuildMemberBoostEndEvent) DesiredEventType() Event { return &GuildMemberBoostEndEvent{} }
