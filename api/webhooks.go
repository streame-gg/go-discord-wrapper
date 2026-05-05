package api

import (
	"bytes"
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
func (c *RestClient) CreateWebhook(channelID common.Snowflake, params CreateWebhookParams) (*common.Webhook, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPost, "/channels/"+channelID.String()+"/webhooks", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetChannelWebhooks returns all webhooks for a channel.
func (c *RestClient) GetChannelWebhooks(channelID common.Snowflake) ([]*common.Webhook, error) {
	req, err := c.generateRequest(http.MethodGet, "/channels/"+channelID.String()+"/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var result []*common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGuildWebhooks returns all webhooks in a guild.
func (c *RestClient) GetGuildWebhooks(guildID common.Snowflake) ([]*common.Webhook, error) {
	req, err := c.generateRequest(http.MethodGet, "/guilds/"+guildID.String()+"/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var result []*common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetWebhook returns the webhook object for the given ID.
func (c *RestClient) GetWebhook(webhookID common.Snowflake) (*common.Webhook, error) {
	req, err := c.generateRequest(http.MethodGet, "/webhooks/"+webhookID.String(), nil)
	if err != nil {
		return nil, err
	}

	var result common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWebhookWithToken returns a webhook without requiring authentication, using its token.
func (c *RestClient) GetWebhookWithToken(webhookID common.Snowflake, token string) (*common.Webhook, error) {
	req, err := c.generateRequest(http.MethodGet, "/webhooks/"+webhookID.String()+"/"+token, nil)
	if err != nil {
		return nil, err
	}

	var result common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyWebhook updates a webhook's name, avatar, or channel.
func (c *RestClient) ModifyWebhook(webhookID common.Snowflake, params ModifyWebhookParams) (*common.Webhook, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPatch, "/webhooks/"+webhookID.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ModifyWebhookWithToken updates a webhook using its token (no bot authentication required).
func (c *RestClient) ModifyWebhookWithToken(webhookID common.Snowflake, token string, params ModifyWebhookParams) (*common.Webhook, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPatch, "/webhooks/"+webhookID.String()+"/"+token, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result common.Webhook
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteWebhook deletes a webhook permanently.
func (c *RestClient) DeleteWebhook(webhookID common.Snowflake) error {
	req, err := c.generateRequest(http.MethodDelete, "/webhooks/"+webhookID.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// DeleteWebhookWithToken deletes a webhook using its token (no bot authentication required).
func (c *RestClient) DeleteWebhookWithToken(webhookID common.Snowflake, token string) error {
	req, err := c.generateRequest(http.MethodDelete, "/webhooks/"+webhookID.String()+"/"+token, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ExecuteWebhook sends a message via a webhook. Returns the created message when query.Wait is true.
func (c *RestClient) ExecuteWebhook(webhookID common.Snowflake, token string, params ExecuteWebhookParams, query ExecuteWebhookQueryParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + token + query.toQuery()
	req, err := c.generateRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if query.Wait != nil && *query.Wait {
		var result common.Message
		if _, err := c.do(req, http.StatusOK, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	if _, err := c.do(req, http.StatusNoContent, nil); err != nil {
		return nil, err
	}
	return nil, nil
}

// GetWebhookMessage returns a previously sent webhook message.
func (c *RestClient) GetWebhookMessage(webhookID common.Snowflake, token string, messageID common.Snowflake) (*common.Message, error) {
	path := "/webhooks/" + webhookID.String() + "/" + token + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result common.Message
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditWebhookMessage edits a previously sent webhook message.
func (c *RestClient) EditWebhookMessage(webhookID common.Snowflake, token string, messageID common.Snowflake, params EditMessageParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + token + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result common.Message
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteWebhookMessage deletes a previously sent webhook message.
func (c *RestClient) DeleteWebhookMessage(webhookID common.Snowflake, token string, messageID common.Snowflake) error {
	path := "/webhooks/" + webhookID.String() + "/" + token + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
