package discord

import (
	"context"
	"errors"
)

var errIntegrationNoGuildID = errors.New("discord: Integration.GuildID is not set — ensure the integration was received with guild context")

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the integration.
func (i *Integration) Hydrate(c EntityClient) { i.hClient = c }

// WithClient returns a shallow copy of the Integration with the client replaced.
func (i *Integration) WithClient(c EntityClient) *Integration {
	cp := *i
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the integration has an associated client.
func (i *Integration) IsHydrated() bool { return i.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Delete removes this integration from its guild. Requires MANAGE_GUILD.
func (i *Integration) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(i.hClient)
	if err != nil {
		return err
	}
	if i.GuildID == "" {
		return errIntegrationNoGuildID
	}
	return c.DeleteGuildIntegration(ctx, i.GuildID, i.ID, reason)
}
