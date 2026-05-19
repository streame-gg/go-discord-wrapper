package discord

import "context"

// ── Options ───────────────────────────────────────────────────────────────────

// WebhookEditOptions configures Webhook.Edit.
type WebhookEditOptions struct {
	Name           *string
	Avatar         *string
	ChannelID      *Snowflake
	AuditLogReason *string
}

// WebhookExecuteOptions configures Webhook.Execute.
type WebhookExecuteOptions struct {
	Content         string
	Username        *string
	AvatarURL       *string
	TTS             bool
	Embeds          []Embed
	AllowedMentions *AllowedMentions
	Components      []AnyComponent
	Flags           MessageFlag
	ThreadName      *string
	Files           []MessageFile
	Wait            bool
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the webhook.
func (w *Webhook) Hydrate(c EntityClient) { w.hClient = c }

// WithClient returns a shallow copy of the Webhook with the client replaced.
func (w *Webhook) WithClient(c EntityClient) *Webhook {
	cp := *w
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the webhook has an associated client.
func (w *Webhook) IsHydrated() bool { return w.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this webhook. Requires MANAGE_WEBHOOKS.
func (w *Webhook) Edit(ctx context.Context, opts WebhookEditOptions) (*Webhook, error) {
	c, err := ensureClient(w.hClient)
	if err != nil {
		return nil, err
	}
	return c.ModifyWebhook(ctx, w.ID, opts)
}

// Delete deletes this webhook. Requires MANAGE_WEBHOOKS.
func (w *Webhook) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(w.hClient)
	if err != nil {
		return err
	}
	return c.DeleteWebhook(ctx, w.ID, reason)
}

// Execute sends a message via this webhook. Token must be non-nil.
func (w *Webhook) Execute(ctx context.Context, opts WebhookExecuteOptions) (*Message, error) {
	c, err := ensureClient(w.hClient)
	if err != nil {
		return nil, err
	}
	if w.Token == nil {
		return nil, errWebhookNoToken
	}
	return c.ExecuteWebhook(ctx, w.ID, *w.Token, opts)
}

// FetchMessage retrieves a previously sent webhook message.
func (w *Webhook) FetchMessage(ctx context.Context, messageID Snowflake) (*Message, error) {
	c, err := ensureClient(w.hClient)
	if err != nil {
		return nil, err
	}
	if w.Token == nil {
		return nil, errWebhookNoToken
	}
	return c.GetWebhookMessage(ctx, w.ID, *w.Token, messageID)
}

var errWebhookNoToken = errorf("discord: Webhook.Token is not set — only incoming webhooks have tokens")

func errorf(s string) error {
	return &webhookErr{s}
}

type webhookErr struct{ msg string }

func (e *webhookErr) Error() string { return e.msg }
