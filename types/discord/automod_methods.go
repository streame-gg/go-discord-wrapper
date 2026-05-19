package discord

import "context"

// ── Options ───────────────────────────────────────────────────────────────────

// RuleEditOptions configures AutoModerationRule.Edit.
type RuleEditOptions struct {
	Name            *string
	EventType       *AutoModerationEventType
	TriggerMetadata *AutoModerationTriggerMetadata
	Actions         []AutoModerationAction
	Enabled         *bool
	ExemptRoles     []Snowflake
	ExemptChannels  []Snowflake
	AuditLogReason  *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the auto moderation rule.
func (r *AutoModerationRule) Hydrate(c EntityClient) { r.hClient = c }

// WithClient returns a shallow copy of the AutoModerationRule with the client replaced.
func (r *AutoModerationRule) WithClient(c EntityClient) *AutoModerationRule {
	cp := *r
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the rule has an associated client.
func (r *AutoModerationRule) IsHydrated() bool { return r.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this auto moderation rule. Requires MANAGE_GUILD.
func (r *AutoModerationRule) Edit(ctx context.Context, opts RuleEditOptions) (*AutoModerationRule, error) {
	c, err := ensureClient(r.hClient)
	if err != nil {
		return nil, err
	}
	return c.ModifyAutoModerationRule(ctx, r.GuildID, r.ID, opts)
}

// Delete deletes this auto moderation rule. Requires MANAGE_GUILD.
func (r *AutoModerationRule) Delete(ctx context.Context) error {
	c, err := ensureClient(r.hClient)
	if err != nil {
		return err
	}
	return c.DeleteAutoModerationRule(ctx, r.GuildID, r.ID)
}
