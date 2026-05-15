package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type CreateWebhookParams struct {
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
}

type ModifyWebhookParams struct {
	Name      *string           `json:"name,omitempty"`
	Avatar    *string           `json:"avatar,omitempty"`
	ChannelID *common.Snowflake `json:"channel_id,omitempty"`
}

type ExecuteWebhookParams struct {
	Content         string                  `json:"content,omitempty"`
	Username        *string                 `json:"username,omitempty"`
	AvatarURL       *string                 `json:"avatar_url,omitempty"`
	TTS             bool                    `json:"tts,omitempty"`
	Embeds          []common.Embed          `json:"embeds,omitempty"`
	AllowedMentions *common.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []common.AnyComponent   `json:"-"`
	Flags           common.MessageFlag      `json:"flags,omitempty"`
	ThreadName      *string                 `json:"thread_name,omitempty"`
	// Files are binary attachments sent via multipart/form-data.
	// When set, the request is encoded as multipart rather than JSON.
	Files []MessageFile `json:"-"`
}

func (p ExecuteWebhookParams) MarshalJSON() ([]byte, error) {
	type Alias ExecuteWebhookParams
	return json.Marshal(&struct {
		Components []common.AnyComponent `json:"components,omitempty"`
		Alias
	}{
		Components: p.Components,
		Alias:      (Alias)(p),
	})
}

type ExecuteWebhookQueryParams struct {
	Wait     *bool
	ThreadID *common.Snowflake
}

func (p ExecuteWebhookQueryParams) toQuery() string {
	q := url.Values{}
	if p.Wait != nil {
		q.Set("wait", strconv.FormatBool(*p.Wait))
	}
	if p.ThreadID != nil {
		q.Set("thread_id", p.ThreadID.String())
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// CreateWebhook creates a new webhook for a channel. Requires MANAGE_WEBHOOKS.
func (c *RestClient) CreateWebhook(ctx context.Context, channelID common.Snowflake, params CreateWebhookParams) (*common.Webhook, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/webhooks", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetChannelWebhooks returns all webhooks for a channel.
func (c *RestClient) GetChannelWebhooks(ctx context.Context, channelID common.Snowflake) (*[]*common.Webhook, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/channels/"+channelID.String()+"/webhooks", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildWebhooks returns all webhooks in a guild.
func (c *RestClient) GetGuildWebhooks(ctx context.Context, guildID common.Snowflake) (*[]*common.Webhook, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/webhooks", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetWebhook returns the webhook object for the given ID.
func (c *RestClient) GetWebhook(ctx context.Context, webhookID common.Snowflake) (*common.Webhook, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/webhooks/"+webhookID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetWebhookWithToken returns a webhook without requiring authentication, using its token.
func (c *RestClient) GetWebhookWithToken(ctx context.Context, webhookID common.Snowflake, token string) (*common.Webhook, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/webhooks/"+webhookID.String()+"/"+url.PathEscape(token), nil)
	if err != nil {
		return nil, err
	}

	return doRequest[common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyWebhook updates a webhook's name, avatar, or channel.
func (c *RestClient) ModifyWebhook(ctx context.Context, webhookID common.Snowflake, params ModifyWebhookParams) (*common.Webhook, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/webhooks/"+webhookID.String(), bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyWebhookWithToken updates a webhook using its token (no bot authentication required).
func (c *RestClient) ModifyWebhookWithToken(ctx context.Context, webhookID common.Snowflake, token string, params ModifyWebhookParams) (*common.Webhook, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/webhooks/"+webhookID.String()+"/"+url.PathEscape(token), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doRequest[common.Webhook](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteWebhook deletes a webhook permanently.
func (c *RestClient) DeleteWebhook(ctx context.Context, webhookID common.Snowflake) error {
	if err := webhookID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/webhooks/"+webhookID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// DeleteWebhookWithToken deletes a webhook using its token (no bot authentication required).
func (c *RestClient) DeleteWebhookWithToken(ctx context.Context, webhookID common.Snowflake, token string) error {
	if err := webhookID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/webhooks/"+webhookID.String()+"/"+url.PathEscape(token), nil)
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ExecuteWebhook sends a message via a webhook. Returns the created message when query.Wait is true.
// When params.Files is non-empty the request is sent as multipart/form-data.
func (c *RestClient) ExecuteWebhook(ctx context.Context, webhookID common.Snowflake, token string, params ExecuteWebhookParams, query ExecuteWebhookQueryParams) (*common.Message, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + query.toQuery()

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if len(params.Files) > 0 {
		buf, ct, err := buildMultipartMessage(jsonBody, params.Files)
		if err != nil {
			return nil, err
		}
		req, err = c.generateRequest(ctx, http.MethodPost, path, buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", ct)
	} else {
		req, err = c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(jsonBody))
		if err != nil {
			return nil, err
		}
	}

	if query.Wait != nil && *query.Wait {
		return doRequest[common.Message](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	return nil, doRequestWithoutResponse(c, req)
}

// GetWebhookMessage returns a previously sent webhook message.
func (c *RestClient) GetWebhookMessage(ctx context.Context, webhookID common.Snowflake, token string, messageID common.Snowflake) (*common.Message, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/messages/" + messageID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditWebhookMessage edits a previously sent webhook message.
func (c *RestClient) EditWebhookMessage(ctx context.Context, webhookID common.Snowflake, token string, messageID common.Snowflake, params EditMessageParams) (*common.Message, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/messages/" + messageID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doRequest[common.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteWebhookMessage deletes a previously sent webhook message.
func (c *RestClient) DeleteWebhookMessage(ctx context.Context, webhookID common.Snowflake, token string, messageID common.Snowflake) error {
	if err := webhookID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/messages/" + messageID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ExecuteSlackWebhook sends a Slack-compatible payload via a webhook.
func (c *RestClient) ExecuteSlackWebhook(ctx context.Context, webhookID common.Snowflake, token string, wait bool, body json.RawMessage) error {
	if err := webhookID.Validate(); err != nil {
		return err
	}

	q := ""
	if wait {
		q = "?wait=true"
	}
	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/slack" + q
	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	if wait {
		return doRequestWithoutResponse(c, req, http.StatusOK)
	}

	return doRequestWithoutResponse(c, req, http.StatusNoContent)
}

// ExecuteGitHubWebhook sends a GitHub-compatible payload via a webhook.
func (c *RestClient) ExecuteGitHubWebhook(ctx context.Context, webhookID common.Snowflake, token string, wait bool, body json.RawMessage) error {
	if err := webhookID.Validate(); err != nil {
		return err
	}

	q := ""
	if wait {
		q = "?wait=true"
	}
	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/github" + q
	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	if wait {
		return doRequestWithoutResponse(c, req, http.StatusOK)
	}

	return doRequestWithoutResponse(c, req, http.StatusNoContent)
}
