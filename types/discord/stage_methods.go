package discord

import "context"

// ── Options ───────────────────────────────────────────────────────────────────

// StageEditOptions configures StageInstance.Edit.
type StageEditOptions struct {
	Topic          *string
	PrivacyLevel   *StageInstancePrivacyLevel
	AuditLogReason *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the stage instance.
func (s *StageInstance) Hydrate(c EntityClient) { s.hClient = c }

// WithClient returns a shallow copy of the StageInstance with the client replaced.
func (s *StageInstance) WithClient(c EntityClient) *StageInstance {
	cp := *s
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the stage instance has an associated client.
func (s *StageInstance) IsHydrated() bool { return s.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this stage instance. Requires MANAGE_CHANNELS.
func (s *StageInstance) Edit(ctx context.Context, opts StageEditOptions) (*StageInstance, error) {
	c, err := ensureClient(s.hClient)
	if err != nil {
		return nil, err
	}
	return c.ModifyStageInstance(ctx, s.ChannelID, opts)
}

// Delete closes this stage instance. Requires MANAGE_CHANNELS.
func (s *StageInstance) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(s.hClient)
	if err != nil {
		return err
	}
	return c.DeleteStageInstance(ctx, s.ChannelID, reason)
}
