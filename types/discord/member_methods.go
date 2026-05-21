package discord

import (
	"context"
	"errors"
	"time"
)

// errMemberNoGuildID is returned when GuildID is not set on a GuildMember.
var errMemberNoGuildID = errors.New("discord: GuildMember.GuildID is not set — ensure the member was received via a gateway event or cache")

// errMemberNoUserID is returned when UserID is not set on a GuildMember.
var errMemberNoUserID = errors.New("discord: GuildMember.UserID is not set — ensure the member was received via a gateway event or cache")

// ── Options ───────────────────────────────────────────────────────────────────

// MemberEditOptions configures GuildMember.Edit.
type MemberEditOptions struct {
	Nick                       *string
	Roles                      []Snowflake
	Mute                       *bool
	Deaf                       *bool
	ChannelID                  *Snowflake
	CommunicationDisabledUntil *string
	Flags                      *int
	AuditLogReason             *string
}

// BanOptions configures GuildMember.Ban.
type BanOptions struct {
	DeleteMessageSeconds *int
	AuditLogReason       *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client so that the convenience methods can be
// called without passing ctx and client explicitly. Also derives UserID from
// User.ID when User is available.
func (m *GuildMember) Hydrate(c EntityClient) {
	m.hClient = c
	if m.User != nil && m.UserID.IsEmpty() {
		m.UserID = m.User.ID
		m.User.Hydrate(c)
	}
}

// WithClient returns a shallow copy of the GuildMember with the client replaced.
func (m *GuildMember) WithClient(c EntityClient) *GuildMember {
	cp := *m
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the member has an associated client.
func (m *GuildMember) IsHydrated() bool { return m.hClient != nil }

// ids returns the guildID and userID needed by most methods, or errors if
// they are not set.
func (m *GuildMember) ids() (guildID, userID Snowflake, err error) {
	if m.GuildID.IsEmpty() {
		return 0, 0, errMemberNoGuildID
	}
	uid := m.UserID
	if uid.IsEmpty() && m.User != nil {
		uid = m.User.ID
	}
	if uid.IsEmpty() {
		return 0, 0, errMemberNoUserID
	}
	return m.GuildID, uid, nil
}

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this member's attributes. Requires appropriate permissions.
func (m *GuildMember) Edit(ctx context.Context, opts MemberEditOptions) (*GuildMember, error) {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return nil, err
	}
	guildID, userID, err := m.ids()
	if err != nil {
		return nil, err
	}
	return c.ModifyGuildMember(ctx, guildID, userID, opts)
}

// Kick removes this member from the guild. Requires KICK_MEMBERS.
func (m *GuildMember) Kick(ctx context.Context, reason *string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	guildID, userID, err := m.ids()
	if err != nil {
		return err
	}
	return c.KickGuildMember(ctx, guildID, userID, reason)
}

// Ban bans this member from the guild. Requires BAN_MEMBERS.
func (m *GuildMember) Ban(ctx context.Context, opts BanOptions) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	guildID, userID, err := m.ids()
	if err != nil {
		return err
	}
	return c.CreateGuildBan(ctx, guildID, userID, opts)
}

// AddRole assigns a role to this member. Requires MANAGE_ROLES.
func (m *GuildMember) AddRole(ctx context.Context, roleID Snowflake, reason *string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	guildID, userID, err := m.ids()
	if err != nil {
		return err
	}
	return c.AddGuildMemberRole(ctx, guildID, userID, roleID, reason)
}

// RemoveRole removes a role from this member. Requires MANAGE_ROLES.
func (m *GuildMember) RemoveRole(ctx context.Context, roleID Snowflake, reason *string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	guildID, userID, err := m.ids()
	if err != nil {
		return err
	}
	return c.RemoveGuildMemberRole(ctx, guildID, userID, roleID, reason)
}

// Timeout applies a communication timeout until the given time.
// Pass a zero time.Time or the zero value to remove an active timeout.
// Requires MODERATE_MEMBERS.
func (m *GuildMember) Timeout(ctx context.Context, until *time.Time, reason *string) (*GuildMember, error) {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return nil, err
	}
	guildID, userID, err := m.ids()
	if err != nil {
		return nil, err
	}
	var disabledUntil *string
	if until != nil && !until.IsZero() {
		s := until.UTC().Format(time.RFC3339)
		disabledUntil = &s
	} else {
		empty := ""
		disabledUntil = &empty
	}
	return c.ModifyGuildMember(ctx, guildID, userID, MemberEditOptions{
		CommunicationDisabledUntil: disabledUntil,
		AuditLogReason:             reason,
	})
}
