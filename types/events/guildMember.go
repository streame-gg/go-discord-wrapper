package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type GuildMemberAddEvent struct {
	common.GuildMember
	GuildID common.Snowflake `json:"guild_id"`
}

func init() {
	RegisterEvent(GuildMemberAddEvent{})
	RegisterEvent(GuildMemberRemoveEvent{})
	RegisterEvent(GuildMemberUpdateEvent{})
}

func (e GuildMemberAddEvent) DesiredEventType() Event { return &GuildMemberAddEvent{} }
func (e GuildMemberAddEvent) Event() EventType        { return EventGuildMemberAdd }

type GuildMemberRemoveEvent struct {
	GuildID common.Snowflake `json:"guild_id"`
	User    common.User      `json:"user"`
}

func (e GuildMemberRemoveEvent) DesiredEventType() Event { return &GuildMemberRemoveEvent{} }
func (e GuildMemberRemoveEvent) Event() EventType        { return EventGuildMemberRemove }

type GuildMemberUpdateEvent struct {
	GuildID                    common.Snowflake   `json:"guild_id"`
	Roles                      []common.Snowflake `json:"roles"`
	User                       common.User        `json:"user"`
	Nick                       *string            `json:"nick,omitempty"`
	AvatarHash                 *string            `json:"avatar,omitempty"`
	JoinedAt                   *time.Time         `json:"joined_at,omitempty"`
	PremiumSince               *time.Time         `json:"premium_since,omitempty"`
	Deaf                       *bool              `json:"deaf,omitempty"`
	Mute                       *bool              `json:"mute,omitempty"`
	Pending                    *bool              `json:"pending,omitempty"`
	CommunicationDisabledUntil *string            `json:"communication_disabled_until,omitempty"`
	Flags                      *int               `json:"flags,omitempty"`
}

func (e GuildMemberUpdateEvent) DesiredEventType() Event { return &GuildMemberUpdateEvent{} }
func (e GuildMemberUpdateEvent) Event() EventType        { return EventGuildMemberUpdate }
