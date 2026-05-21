package discord

import (
	"context"
	"errors"
)

var errRoleNoGuildID = errors.New("discord: Role.GuildID is not set — ensure the role was received via a gateway event or cache")

// ── Options ───────────────────────────────────────────────────────────────────

// RoleCreateOptions configures Guild.CreateRole.
type RoleCreateOptions struct {
	Name           *string
	Permissions    *string
	Color          *int
	Hoist          *bool
	Icon           *string
	UnicodeEmoji   *string
	Mentionable    *bool
	AuditLogReason *string
}

// RoleEditOptions configures Role.Edit.
type RoleEditOptions struct {
	Name           *string
	Permissions    *string
	Color          *int
	Hoist          *bool
	Icon           *string
	UnicodeEmoji   *string
	Mentionable    *bool
	AuditLogReason *string
}

// RolePositionEntry is a single role position change.
type RolePositionEntry struct {
	ID       Snowflake
	Position *int
}

// RolePositionOptions configures Role.SetPosition.
type RolePositionOptions struct {
	Entries        []RolePositionEntry
	AuditLogReason *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the role.
func (r *Role) Hydrate(c EntityClient) { r.hClient = c }

// WithClient returns a shallow copy of the Role with the client replaced.
func (r *Role) WithClient(c EntityClient) *Role {
	cp := *r
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the role has an associated client.
func (r *Role) IsHydrated() bool { return r.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this role. Requires MANAGE_ROLES.
func (r *Role) Edit(ctx context.Context, opts RoleEditOptions) (*Role, error) {
	c, err := ensureClient(r.hClient)
	if err != nil {
		return nil, err
	}
	if r.GuildID.IsEmpty() {
		return nil, errRoleNoGuildID
	}
	return c.ModifyGuildRole(ctx, r.GuildID, r.ID, opts)
}

// Delete deletes this role. Requires MANAGE_ROLES.
func (r *Role) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(r.hClient)
	if err != nil {
		return err
	}
	if r.GuildID.IsEmpty() {
		return errRoleNoGuildID
	}
	return c.DeleteGuildRole(ctx, r.GuildID, r.ID, reason)
}

// SetPosition changes this role's position in the guild's role list.
// Returns the full updated role list. Requires MANAGE_ROLES.
func (r *Role) SetPosition(ctx context.Context, opts RolePositionOptions) ([]*Role, error) {
	c, err := ensureClient(r.hClient)
	if err != nil {
		return nil, err
	}
	if r.GuildID.IsEmpty() {
		return nil, errRoleNoGuildID
	}
	return c.ModifyGuildRolePositions(ctx, r.GuildID, opts)
}
