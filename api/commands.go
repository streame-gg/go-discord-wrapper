package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// RegisterCommand registers a single application command for the given application ID.
// https://docs.discord.com/developers/interactions/application-commands#create-global-application-command
func (c *RestClient) RegisterCommand(ctx context.Context, appID discord.Snowflake, cmd *discord.ApplicationCommand) (*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, "POST", "/applications/"+appID.String()+"/commands", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK:      true,
		http.StatusCreated: true,
	})
}

// BulkRegisterCommands overwrites all global application commands for the given application ID.
// Any commands not included in cmds are deleted.
// https://docs.discord.com/developers/interactions/application-commands#bulk-overwrite-global-application-commands
func (c *RestClient) BulkRegisterCommands(ctx context.Context, appID discord.Snowflake, cmds []*discord.ApplicationCommand) ([]*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, "PUT", "/applications/"+appID.String()+"/commands", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Global command management ─────────────────────────────────────────────────

// ListGlobalApplicationCommands returns all global application commands for the given application ID.
// Set withLocalizations to true to include localization dictionaries.
// https://docs.discord.com/developers/interactions/application-commands#get-global-application-commands
func (c *RestClient) ListGlobalApplicationCommands(ctx context.Context, appID discord.Snowflake, withLocalizations bool) ([]*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/commands"
	if withLocalizations {
		path += "?with_localizations=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGlobalApplicationCommand returns a single global application command.
// https://docs.discord.com/developers/interactions/application-commands#get-global-application-command
func (c *RestClient) GetGlobalApplicationCommand(ctx context.Context, appID, cmdID discord.Snowflake) (*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := cmdID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditGlobalApplicationCommand updates a global application command.
// https://docs.discord.com/developers/interactions/application-commands#edit-global-application-command
func (c *RestClient) EditGlobalApplicationCommand(ctx context.Context, appID, cmdID discord.Snowflake, params *discord.ApplicationCommand) (*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := cmdID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGlobalApplicationCommand deletes a global application command.
// https://docs.discord.com/developers/interactions/application-commands#delete-global-application-command
func (c *RestClient) DeleteGlobalApplicationCommand(ctx context.Context, appID, cmdID discord.Snowflake) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	if err := cmdID.Validate(); err != nil {
		return err
	}

	path := "/applications/" + appID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ── Guild command management ──────────────────────────────────────────────────

// ListGuildApplicationCommands returns all application commands registered to a specific guild.
// https://docs.discord.com/developers/interactions/application-commands#get-guild-application-commands
func (c *RestClient) ListGuildApplicationCommands(ctx context.Context, appID, guildID discord.Snowflake, withLocalizations bool) ([]*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands"
	if withLocalizations {
		path += "?with_localizations=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildApplicationCommand registers a command in a specific guild.
// https://docs.discord.com/developers/interactions/application-commands#create-guild-application-command
func (c *RestClient) CreateGuildApplicationCommand(ctx context.Context, appID, guildID discord.Snowflake, cmd *discord.ApplicationCommand) (*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands"
	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK:      true,
		http.StatusCreated: true,
	})
}

// GetGuildApplicationCommand returns a single guild-specific application command.
// https://docs.discord.com/developers/interactions/application-commands#get-guild-application-command
func (c *RestClient) GetGuildApplicationCommand(ctx context.Context, appID, guildID, cmdID discord.Snowflake) (*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := cmdID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditGuildApplicationCommand updates a guild-specific application command.
// https://docs.discord.com/developers/interactions/application-commands#edit-guild-application-command
func (c *RestClient) EditGuildApplicationCommand(ctx context.Context, appID, guildID, cmdID discord.Snowflake, params *discord.ApplicationCommand) (*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := cmdID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildApplicationCommand deletes a guild-specific application command.
// https://docs.discord.com/developers/interactions/application-commands#delete-guild-application-command
func (c *RestClient) DeleteGuildApplicationCommand(ctx context.Context, appID, guildID, cmdID discord.Snowflake) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := cmdID.Validate(); err != nil {
		return err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// BulkOverwriteGuildApplicationCommands overwrites all guild-specific commands for the given guild.
// Any commands not included in cmds are deleted.
// https://docs.discord.com/developers/interactions/application-commands#bulk-overwrite-guild-application-commands
func (c *RestClient) BulkOverwriteGuildApplicationCommands(ctx context.Context, appID, guildID discord.Snowflake, cmds []*discord.ApplicationCommand) ([]*discord.ApplicationCommand, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands"
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.ApplicationCommand](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Command permissions ───────────────────────────────────────────────────────

// ListGuildApplicationCommandPermissions returns all permission overrides for every command in a guild.
// https://docs.discord.com/developers/interactions/application-commands#get-guild-application-command-permissions
func (c *RestClient) ListGuildApplicationCommandPermissions(ctx context.Context, appID, guildID discord.Snowflake) ([]*discord.GuildApplicationCommandPermissions, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/permissions"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.GuildApplicationCommandPermissions](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetApplicationCommandPermissions returns the permission overrides for a specific command in a guild.
// https://docs.discord.com/developers/interactions/application-commands#get-application-command-permissions
func (c *RestClient) GetApplicationCommandPermissions(ctx context.Context, appID, guildID, cmdID discord.Snowflake) (*discord.GuildApplicationCommandPermissions, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := cmdID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String() + "/permissions"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildApplicationCommandPermissions](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditApplicationCommandPermissions overwrites the permission overrides for a specific command in a guild.
// Accepts either a bot token or an OAuth2 Bearer token with the applications.commands.permissions.update scope.
// https://docs.discord.com/developers/interactions/application-commands#edit-application-command-permissions
func (c *RestClient) EditApplicationCommandPermissions(ctx context.Context, appID, guildID, cmdID discord.Snowflake, permissions []discord.ApplicationCommandPermission) (*discord.GuildApplicationCommandPermissions, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := cmdID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]any{"permissions": permissions})
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String() + "/permissions"
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildApplicationCommandPermissions](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
