package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildMemberAddEvent struct {
	discord.GuildMember
	GuildID discord.Snowflake `json:"guild_id"`
}

func init() {
	RegisterEvent(GuildMemberAddEvent{})
	RegisterEvent(GuildMemberRemoveEvent{})
	RegisterEvent(GuildMemberUpdateEvent{})
}

func (e GuildMemberAddEvent) DesiredEventType() Event { return &GuildMemberAddEvent{} }
func (e GuildMemberAddEvent) Event() EventType        { return EventGuildMemberAdd }

type GuildMemberRemoveEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

func (e GuildMemberRemoveEvent) DesiredEventType() Event { return &GuildMemberRemoveEvent{} }
func (e GuildMemberRemoveEvent) Event() EventType        { return EventGuildMemberRemove }

type GuildMemberUpdateEvent struct {
	GuildID                    discord.Snowflake    `json:"guild_id"`
	Roles                      []discord.Snowflake  `json:"roles"`
	User                       discord.User         `json:"user"`
	Nick                       *string              `json:"nick,omitempty"`
	AvatarHash                 *string              `json:"avatar,omitempty"`
	JoinedAt                   *time.Time           `json:"joined_at,omitempty"`
	PremiumSince               *time.Time           `json:"premium_since,omitempty"`
	Deaf                       *bool                `json:"deaf,omitempty"`
	Mute                       *bool                `json:"mute,omitempty"`
	Pending                    *bool                `json:"pending,omitempty"`
	CommunicationDisabledUntil *string              `json:"communication_disabled_until,omitempty"`
	Flags                      *int                 `json:"flags,omitempty"`
	OldMember                  *discord.GuildMember `json:"old_member,omitempty"`
}

func (e GuildMemberUpdateEvent) DesiredEventType() Event { return &GuildMemberUpdateEvent{} }
func (e GuildMemberUpdateEvent) Event() EventType        { return EventGuildMemberUpdate }
