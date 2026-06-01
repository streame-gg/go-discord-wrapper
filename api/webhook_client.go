package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"sync"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// WebhookClient is a standalone client for a single Discord webhook, identified
// by its ID and token. It executes and manages the webhook through the
// token-authenticated endpoints, so it needs neither a bot token nor a gateway
// connection.
//
// It shares the same request pipeline as RestClient (proactive rate limiting,
// retries, lifecycle events). Call Close when finished to release the
// rate-limiter's background goroutine.
type WebhookClient struct {
	rest  *RestClient
	id    discord.Snowflake
	token string

	// mu guards the persistent send defaults, which may be mutated via Set/Reset
	// while Send reads them concurrently from another goroutine.
	mu               sync.RWMutex
	defaultUsername  *string
	defaultAvatarURL *string
}

// NewWebhookClient creates a WebhookClient from a webhook ID and token. The same
// options accepted by NewRestClient (retry, rate limiting, HTTP client, …) apply.
func NewWebhookClient(id discord.Snowflake, token string, opts ...options.Option) (*WebhookClient, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("go-discord-wrapper: webhook token must not be empty")
	}

	rest, err := newRestClient("", opts...)
	if err != nil {
		return nil, err
	}

	return &WebhookClient{rest: rest, id: id, token: token}, nil
}

// NewWebhookClientFromURL creates a WebhookClient from a full webhook URL of the
// form https://discord.com/api/webhooks/{id}/{token} (an optional /vN API
// version segment and any query string are ignored).
func NewWebhookClientFromURL(webhookURL string, opts ...options.Option) (*WebhookClient, error) {
	id, token, err := parseWebhookURL(webhookURL)
	if err != nil {
		return nil, err
	}
	return NewWebhookClient(id, token, opts...)
}

// parseWebhookURL extracts the webhook ID and token from a webhook URL.
func parseWebhookURL(webhookURL string) (discord.Snowflake, string, error) {
	u, err := url.Parse(strings.TrimSpace(webhookURL))
	if err != nil {
		return 0, "", err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	for i, p := range parts {
		if p == "webhooks" && i+2 < len(parts) {
			id, err := discord.ParseSnowflake(parts[i+1])
			if err != nil {
				return 0, "", err
			}
			token := parts[i+2]
			if token == "" {
				return 0, "", errors.New("go-discord-wrapper: webhook URL is missing a token")
			}
			return *id, token, nil
		}
	}
	return 0, "", errors.New("go-discord-wrapper: not a valid webhook URL")
}

// ID returns the webhook's ID.
func (w *WebhookClient) ID() discord.Snowflake { return w.id }

// ── Persistent send defaults ──────────────────────────────────────────────────
//
// Set* methods store a default that Send applies to every message that does not
// already specify that field; Reset* methods clear them. They are safe for
// concurrent use and return the client so calls can be chained.

// SetUsername sets the default username override applied to subsequent Sends.
func (w *WebhookClient) SetUsername(username string) *WebhookClient {
	w.mu.Lock()
	w.defaultUsername = &username
	w.mu.Unlock()
	return w
}

// SetAvatarURL sets the default avatar URL override applied to subsequent Sends.
func (w *WebhookClient) SetAvatarURL(avatarURL string) *WebhookClient {
	w.mu.Lock()
	w.defaultAvatarURL = &avatarURL
	w.mu.Unlock()
	return w
}

// ResetUsername clears the default username override.
func (w *WebhookClient) ResetUsername() *WebhookClient {
	w.mu.Lock()
	w.defaultUsername = nil
	w.mu.Unlock()
	return w
}

// ResetAvatarURL clears the default avatar URL override.
func (w *WebhookClient) ResetAvatarURL() *WebhookClient {
	w.mu.Lock()
	w.defaultAvatarURL = nil
	w.mu.Unlock()
	return w
}

// ResetDefaults clears all persistent send defaults.
func (w *WebhookClient) ResetDefaults() *WebhookClient {
	w.mu.Lock()
	w.defaultUsername = nil
	w.defaultAvatarURL = nil
	w.mu.Unlock()
	return w
}

// applyDefaults fills in p's Username/AvatarURL from the stored defaults when p
// does not already set them.
func (w *WebhookClient) applyDefaults(p *ExecuteWebhookParams) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if p.Username == nil && w.defaultUsername != nil {
		v := *w.defaultUsername
		p.Username = &v
	}
	if p.AvatarURL == nil && w.defaultAvatarURL != nil {
		v := *w.defaultAvatarURL
		p.AvatarURL = &v
	}
}

// Close releases the underlying rate limiter's background goroutine. Call it
// when the WebhookClient is no longer needed.
func (w *WebhookClient) Close() { w.rest.Close() }

// Fetch returns the webhook's metadata.
func (w *WebhookClient) Fetch(ctx context.Context) (*discord.Webhook, error) {
	return w.rest.GetWebhookWithToken(ctx, w.id, w.token)
}

// Edit modifies the webhook's name and/or avatar. Channel changes require bot
// authentication and are not available through the token, so ModifyWebhookParams
// fields other than Name and Avatar are ignored by Discord here.
func (w *WebhookClient) Edit(ctx context.Context, params ModifyWebhookParams) (*discord.Webhook, error) {
	return w.rest.ModifyWebhookWithToken(ctx, w.id, w.token, params)
}

// Delete permanently deletes the webhook.
func (w *WebhookClient) Delete(ctx context.Context) error {
	return w.rest.DeleteWebhookWithToken(ctx, w.id, w.token)
}

// Send executes the webhook, posting a message. Any persistent defaults set via
// SetUsername/SetAvatarURL are applied when params does not specify those
// fields. The returned message is non-nil only when query.Wait is true (Discord
// otherwise returns no body).
func (w *WebhookClient) Send(ctx context.Context, params ExecuteWebhookParams, query ExecuteWebhookQueryParams) (*discord.Message, error) {
	w.applyDefaults(&params)
	return w.rest.ExecuteWebhook(ctx, w.id, w.token, params, query)
}

// SendSlack posts a Slack-compatible payload via the webhook's /slack endpoint.
func (w *WebhookClient) SendSlack(ctx context.Context, wait bool, body json.RawMessage) error {
	return w.rest.ExecuteSlackWebhook(ctx, w.id, w.token, wait, body)
}

// SendGitHub posts a GitHub-compatible payload via the webhook's /github endpoint.
func (w *WebhookClient) SendGitHub(ctx context.Context, wait bool, body json.RawMessage) error {
	return w.rest.ExecuteGitHubWebhook(ctx, w.id, w.token, wait, body)
}

// FetchMessage returns a message previously sent by this webhook.
func (w *WebhookClient) FetchMessage(ctx context.Context, messageID discord.Snowflake) (*discord.Message, error) {
	return w.rest.GetWebhookMessage(ctx, w.id, w.token, messageID)
}

// EditMessage edits a message previously sent by this webhook.
func (w *WebhookClient) EditMessage(ctx context.Context, messageID discord.Snowflake, params EditMessageParams) (*discord.Message, error) {
	return w.rest.EditWebhookMessage(ctx, w.id, w.token, messageID, params)
}

// DeleteMessage deletes a message previously sent by this webhook.
func (w *WebhookClient) DeleteMessage(ctx context.Context, messageID discord.Snowflake) error {
	return w.rest.DeleteWebhookMessage(ctx, w.id, w.token, messageID)
}
