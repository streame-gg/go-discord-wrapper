package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// RegisterCommand registers a single application command for the given application ID.
func (c *RestClient) RegisterCommand(ctx context.Context, appID common.Snowflake, cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
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

	var registered commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[commands.ApplicationCommand]{
		status: http.StatusCreated,
		Out:    &registered,
	}); err != nil {
		return nil, err
	}
	return &registered, nil
}

// BulkRegisterCommands overwrites all global application commands for the given application ID.
// Any commands not included in cmds are deleted.
func (c *RestClient) BulkRegisterCommands(ctx context.Context, appID common.Snowflake, cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
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

	var registered []*commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[[]*commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &registered,
	}); err != nil {
		return nil, err
	}

	return registered, nil
}

// ── Global command management ─────────────────────────────────────────────────

// GetGlobalApplicationCommands returns all global application commands for the given application ID.
// Set withLocalizations to true to include localization dictionaries.
func (c *RestClient) GetGlobalApplicationCommands(ctx context.Context, appID common.Snowflake, withLocalizations bool) ([]*commands.ApplicationCommand, error) {
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

	var cmds []*commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[[]*commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &cmds,
	}); err != nil {
		return nil, err
	}

	return cmds, nil
}

// GetGlobalApplicationCommand returns a single global application command.
func (c *RestClient) GetGlobalApplicationCommand(ctx context.Context, appID, cmdID common.Snowflake) (*commands.ApplicationCommand, error) {
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

	var cmd commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &cmd,
	}); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// EditGlobalApplicationCommand updates a global application command.
func (c *RestClient) EditGlobalApplicationCommand(ctx context.Context, appID, cmdID common.Snowflake, params *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
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

	var cmd commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &cmd,
	}); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// DeleteGlobalApplicationCommand deletes a global application command.
func (c *RestClient) DeleteGlobalApplicationCommand(ctx context.Context, appID, cmdID common.Snowflake) error {
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

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// ── Guild command management ──────────────────────────────────────────────────

// GetGuildApplicationCommands returns all application commands registered to a specific guild.
func (c *RestClient) GetGuildApplicationCommands(ctx context.Context, appID, guildID common.Snowflake, withLocalizations bool) ([]*commands.ApplicationCommand, error) {
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

	var cmds []*commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[[]*commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &cmds,
	}); err != nil {
		return nil, err
	}

	return cmds, nil
}

// CreateGuildApplicationCommand registers a command in a specific guild.
func (c *RestClient) CreateGuildApplicationCommand(ctx context.Context, appID, guildID common.Snowflake, cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
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

	var registered commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &registered,
	}); err != nil {
		return nil, err
	}

	return &registered, nil
}

// GetGuildApplicationCommand returns a single guild-specific application command.
func (c *RestClient) GetGuildApplicationCommand(ctx context.Context, appID, guildID, cmdID common.Snowflake) (*commands.ApplicationCommand, error) {
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

	var cmd commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &cmd,
	}); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// EditGuildApplicationCommand updates a guild-specific application command.
func (c *RestClient) EditGuildApplicationCommand(ctx context.Context, appID, guildID, cmdID common.Snowflake, params *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
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

	var cmd commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &cmd,
	}); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// DeleteGuildApplicationCommand deletes a guild-specific application command.
func (c *RestClient) DeleteGuildApplicationCommand(ctx context.Context, appID, guildID, cmdID common.Snowflake) error {
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

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// BulkOverwriteGuildApplicationCommands overwrites all guild-specific commands for the given guild.
// Any commands not included in cmds are deleted.
func (c *RestClient) BulkOverwriteGuildApplicationCommands(ctx context.Context, appID, guildID common.Snowflake, cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
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

	var registered []*commands.ApplicationCommand
	if err := c.do(req, SuccessReturn[[]*commands.ApplicationCommand]{
		status: http.StatusOK,
		Out:    &registered,
	}); err != nil {
		return nil, err
	}

	return registered, nil
}

// ── Command permissions ───────────────────────────────────────────────────────

// GetGuildApplicationCommandPermissions returns all permission overrides for every command in a guild.
func (c *RestClient) GetGuildApplicationCommandPermissions(ctx context.Context, appID, guildID common.Snowflake) ([]*common.GuildApplicationCommandPermissions, error) {
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

	var perms []*common.GuildApplicationCommandPermissions
	if err := c.do(req, SuccessReturn[[]*common.GuildApplicationCommandPermissions]{
		status: http.StatusOK,
		Out:    &perms,
	}); err != nil {
		return nil, err
	}

	return perms, nil
}

// GetApplicationCommandPermissions returns the permission overrides for a specific command in a guild.
func (c *RestClient) GetApplicationCommandPermissions(ctx context.Context, appID, guildID, cmdID common.Snowflake) (*common.GuildApplicationCommandPermissions, error) {
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

	var perms common.GuildApplicationCommandPermissions
	if err := c.do(req, SuccessReturn[common.GuildApplicationCommandPermissions]{
		status: http.StatusOK,
		Out:    &perms,
	}); err != nil {
		return nil, err
	}

	return &perms, nil
}

// EditApplicationCommandPermissions overwrites the permission overrides for a specific command in a guild.
// Requires a Bearer token with applications.commands.permissions.update scope; bot tokens cannot use this endpoint.
func (c *RestClient) EditApplicationCommandPermissions(ctx context.Context, appID, guildID, cmdID common.Snowflake, permissions []common.ApplicationCommandPermission) (*common.GuildApplicationCommandPermissions, error) {
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

	var perms common.GuildApplicationCommandPermissions
	if err := c.do(req, SuccessReturn[common.GuildApplicationCommandPermissions]{
		status: http.StatusOK,
		Out:    &perms,
	}); err != nil {
		return nil, err
	}

	return &perms, nil
}
