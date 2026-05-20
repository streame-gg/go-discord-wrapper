package discord

import "context"

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the invite.
func (i *Invite) Hydrate(c EntityClient) { i.hClient = c }

// WithClient returns a shallow copy of the Invite with the client replaced.
func (i *Invite) WithClient(c EntityClient) *Invite {
	cp := *i
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the invite has an associated client.
func (i *Invite) IsHydrated() bool { return i.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Delete revokes this invite. Requires MANAGE_CHANNELS or MANAGE_GUILD.
func (i *Invite) Delete(ctx context.Context, reason *string) (*Invite, error) {
	c, err := ensureClient(i.hClient)
	if err != nil {
		return nil, err
	}
	return c.DeleteInvite(ctx, i.Code, reason)
}
