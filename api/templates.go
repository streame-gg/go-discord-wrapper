package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type CreateGuildTemplateParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type ModifyGuildTemplateParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateGuildFromTemplateParams struct {
	Name string  `json:"name"`
	Icon *string `json:"icon,omitempty"`
}

// ── Template endpoints ────────────────────────────────────────────────────────

// GetTemplate returns the guild template for the given template code.
func (c *RestClient) GetTemplate(ctx context.Context, templateCode string) (*common.GuildTemplate, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/templates/"+url.PathEscape(templateCode), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var tmpl common.GuildTemplate
	if _, err := c.do(req, []int{http.StatusOK}, &tmpl); err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// CreateGuildFromTemplate creates a new guild from a template. Only available for bots in fewer than 10 guilds.
func (c *RestClient) CreateGuildFromTemplate(ctx context.Context, templateCode string, params CreateGuildFromTemplateParams) (*common.Guild, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/templates/"+url.PathEscape(templateCode), bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var guild common.Guild
	if _, err := c.do(req, []int{http.StatusCreated}, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

// GetGuildTemplates returns the templates for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildTemplates(ctx context.Context, guildID common.Snowflake) ([]*common.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/templates", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var templates []*common.GuildTemplate
	if _, err := c.do(req, []int{http.StatusOK}, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// CreateGuildTemplate creates a template for a guild. Requires MANAGE_GUILD.
func (c *RestClient) CreateGuildTemplate(ctx context.Context, guildID common.Snowflake, params CreateGuildTemplateParams) (*common.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/templates", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var tmpl common.GuildTemplate
	if _, err := c.do(req, []int{http.StatusCreated}, &tmpl); err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// SyncGuildTemplate syncs the template to the guild's current state. Requires MANAGE_GUILD.
func (c *RestClient) SyncGuildTemplate(ctx context.Context, guildID common.Snowflake, templateCode string) (*common.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/templates/" + url.PathEscape(templateCode)
	req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var tmpl common.GuildTemplate
	if _, err := c.do(req, []int{http.StatusOK}, &tmpl); err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// ModifyGuildTemplate modifies a guild template. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuildTemplate(ctx context.Context, guildID common.Snowflake, templateCode string, params ModifyGuildTemplateParams) (*common.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/templates/" + url.PathEscape(templateCode)
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var tmpl common.GuildTemplate
	if _, err := c.do(req, []int{http.StatusOK}, &tmpl); err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// DeleteGuildTemplate deletes a guild template. Requires MANAGE_GUILD.
func (c *RestClient) DeleteGuildTemplate(ctx context.Context, guildID common.Snowflake, templateCode string) (*common.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/templates/" + url.PathEscape(templateCode)
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var tmpl common.GuildTemplate
	if _, err := c.do(req, []int{http.StatusOK}, &tmpl); err != nil {
		return nil, err
	}

	return &tmpl, nil
}
