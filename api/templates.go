package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/guild-template#create-guild-template
type CreateGuildTemplateParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-template#modify-guild-template
type ModifyGuildTemplateParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-template#create-guild-from-guild-template
type CreateGuildFromTemplateParams struct {
	Name string  `json:"name"`
	Icon *string `json:"icon,omitempty"`
}

// ── Template endpoints ────────────────────────────────────────────────────────

// GetTemplate returns the guild template for the given template code.
// https://docs.discord.com/developers/resources/guild-template#get-guild-template
func (c *RestClient) GetTemplate(ctx context.Context, templateCode string) (*discord.GuildTemplate, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/templates/"+url.PathEscape(templateCode), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildTemplate](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildFromTemplate creates a new guild from a template. Only available for bots in fewer than 10 guilds.
// https://docs.discord.com/developers/resources/guild-template#create-guild-from-guild-template
func (c *RestClient) CreateGuildFromTemplate(ctx context.Context, templateCode string, params CreateGuildFromTemplateParams) (*discord.Guild, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/templates/"+url.PathEscape(templateCode), bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Guild](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// ListGuildTemplates returns the templates for a guild. Requires MANAGE_GUILD.
// https://docs.discord.com/developers/resources/guild-template#get-guild-templates
func (c *RestClient) ListGuildTemplates(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/templates", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.GuildTemplate](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildTemplate creates a template for a guild. Requires MANAGE_GUILD.
// https://docs.discord.com/developers/resources/guild-template#create-guild-template
func (c *RestClient) CreateGuildTemplate(ctx context.Context, guildID discord.Snowflake, params CreateGuildTemplateParams) (*discord.GuildTemplate, error) {
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

	return doRequest[discord.GuildTemplate](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// SyncGuildTemplate syncs the template to the guild's current state. Requires MANAGE_GUILD.
// https://docs.discord.com/developers/resources/guild-template#sync-guild-template
func (c *RestClient) SyncGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/templates/" + url.PathEscape(templateCode)
	req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildTemplate](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildTemplate modifies a guild template. Requires MANAGE_GUILD.
// https://docs.discord.com/developers/resources/guild-template#modify-guild-template
func (c *RestClient) ModifyGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string, params ModifyGuildTemplateParams) (*discord.GuildTemplate, error) {
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

	return doRequest[discord.GuildTemplate](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildTemplate deletes a guild template. Requires MANAGE_GUILD.
// https://docs.discord.com/developers/resources/guild-template#delete-guild-template
func (c *RestClient) DeleteGuildTemplate(ctx context.Context, guildID discord.Snowflake, templateCode string) (*discord.GuildTemplate, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/templates/" + url.PathEscape(templateCode)
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildTemplate](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
