package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-member-add
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

// https://docs.discord.com/developers/events/gateway-events#guild-member-remove
type GuildMemberRemoveEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

func (e GuildMemberRemoveEvent) DesiredEventType() Event { return &GuildMemberRemoveEvent{} }
func (e GuildMemberRemoveEvent) Event() EventType        { return EventGuildMemberRemove }

// https://docs.discord.com/developers/events/gateway-events#guild-member-update
type GuildMemberUpdateEvent struct {
	GuildID   discord.Snowflake    `json:"-"`
	NewMember discord.GuildMember  `json:"-"`
	OldMember *discord.GuildMember `json:"-"`
}

func (e *GuildMemberUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID                    discord.Snowflake         `json:"guild_id"`
		Roles                      []discord.Snowflake       `json:"roles"`
		User                       discord.User              `json:"user"`
		Nick                       *string                   `json:"nick,omitempty"`
		AvatarHash                 *string                   `json:"avatar,omitempty"`
		JoinedAt                   *time.Time                `json:"joined_at,omitempty"`
		PremiumSince               *time.Time                `json:"premium_since,omitempty"`
		Deaf                       *bool                     `json:"deaf,omitempty"`
		Mute                       *bool                     `json:"mute,omitempty"`
		Pending                    *bool                     `json:"pending,omitempty"`
		CommunicationDisabledUntil *string                   `json:"communication_disabled_until,omitempty"`
		Flags                      *discord.GuildMemberFlags `json:"flags,omitempty"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewMember.GuildID = w.GuildID
	e.NewMember.User = &w.User
	e.NewMember.UserID = w.User.ID
	e.NewMember.Nick = w.Nick
	e.NewMember.AvatarHash = w.AvatarHash
	e.NewMember.PremiumSince = w.PremiumSince
	if w.CommunicationDisabledUntil != nil {
		t, err := time.Parse(time.RFC3339, *w.CommunicationDisabledUntil)
		if err != nil {
			return fmt.Errorf("communication_disabled_until: %w", err)
		}
		e.NewMember.CommunicationDisabledUntil = &t
	}
	e.NewMember.Roles = make([]discord.Snowflake, len(w.Roles))
	copy(e.NewMember.Roles, w.Roles)
	if w.JoinedAt != nil {
		e.NewMember.JoinedAt = *w.JoinedAt
	}
	e.NewMember.Deaf = w.Deaf
	e.NewMember.Mute = w.Mute
	if w.Pending != nil {
		e.NewMember.Pending = *w.Pending
	}
	if w.Flags != nil {
		e.NewMember.Flags = *w.Flags
	}
	return nil
}

func (e GuildMemberUpdateEvent) MarshalJSON() ([]byte, error) {
	m := e.NewMember
	roles := make([]discord.Snowflake, len(m.Roles))
	copy(roles, m.Roles)
	type wire struct {
		GuildID                    discord.Snowflake        `json:"guild_id"`
		Roles                      []discord.Snowflake      `json:"roles"`
		User                       *discord.User            `json:"user,omitempty"`
		Nick                       *string                  `json:"nick,omitempty"`
		AvatarHash                 *string                  `json:"avatar,omitempty"`
		JoinedAt                   *time.Time               `json:"joined_at,omitempty"`
		PremiumSince               *time.Time               `json:"premium_since,omitempty"`
		CommunicationDisabledUntil *time.Time               `json:"communication_disabled_until,omitempty"`
		Deaf                       *bool                    `json:"deaf,omitempty"`
		Mute                       *bool                    `json:"mute,omitempty"`
		Pending                    bool                     `json:"pending,omitempty"`
		Flags                      discord.GuildMemberFlags `json:"flags,omitempty"`
		OldMember                  *discord.GuildMember     `json:"old_member,omitempty"`
	}
	var jat *time.Time
	if !m.JoinedAt.IsZero() {
		jat = &m.JoinedAt
	}
	return json.Marshal(wire{
		GuildID:                    e.GuildID,
		Roles:                      roles,
		User:                       m.User,
		Nick:                       m.Nick,
		AvatarHash:                 m.AvatarHash,
		JoinedAt:                   jat,
		PremiumSince:               m.PremiumSince,
		CommunicationDisabledUntil: m.CommunicationDisabledUntil,
		Deaf:                       m.Deaf,
		Mute:                       m.Mute,
		Pending:                    m.Pending,
		Flags:                      m.Flags,
		OldMember:                  e.OldMember,
	})
}

func (e GuildMemberUpdateEvent) DesiredEventType() Event { return &GuildMemberUpdateEvent{} }
func (e GuildMemberUpdateEvent) Event() EventType        { return EventGuildMemberUpdate }
