package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// RegisterCommand registers a single application command for the given application ID.
func (c *RestClient) RegisterCommand(appID common.Snowflake, cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest("POST", "/applications/"+appID.String()+"/commands", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var registered commands.ApplicationCommand
	if _, err := c.do(req, http.StatusCreated, &registered); err != nil {
		return nil, err
	}

	return &registered, nil
}

// BulkRegisterCommands overwrites all global application commands for the given application ID.
// Any commands not included in cmds are deleted.
func (c *RestClient) BulkRegisterCommands(appID common.Snowflake, cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
	body, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest("PUT", "/applications/"+appID.String()+"/commands", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var registered []*commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &registered); err != nil {
		return nil, err
	}

	return registered, nil
}

// ── Global command management ─────────────────────────────────────────────────

// GetGlobalApplicationCommands returns all global application commands for the given application ID.
// Set withLocalizations to true to include localization dictionaries.
func (c *RestClient) GetGlobalApplicationCommands(appID common.Snowflake, withLocalizations bool) ([]*commands.ApplicationCommand, error) {
	path := "/applications/" + appID.String() + "/commands"
	if withLocalizations {
		path += "?with_localizations=true"
	}

	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var cmds []*commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &cmds); err != nil {
		return nil, err
	}

	return cmds, nil
}

// GetGlobalApplicationCommand returns a single global application command.
func (c *RestClient) GetGlobalApplicationCommand(appID, cmdID common.Snowflake) (*commands.ApplicationCommand, error) {
	path := "/applications/" + appID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var cmd commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &cmd); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// EditGlobalApplicationCommand updates a global application command.
func (c *RestClient) EditGlobalApplicationCommand(appID, cmdID common.Snowflake, params *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var cmd commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &cmd); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// DeleteGlobalApplicationCommand deletes a global application command.
func (c *RestClient) DeleteGlobalApplicationCommand(appID, cmdID common.Snowflake) error {
	path := "/applications/" + appID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Guild command management ──────────────────────────────────────────────────

// GetGuildApplicationCommands returns all application commands registered to a specific guild.
func (c *RestClient) GetGuildApplicationCommands(appID, guildID common.Snowflake, withLocalizations bool) ([]*commands.ApplicationCommand, error) {
	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands"
	if withLocalizations {
		path += "?with_localizations=true"
	}

	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var cmds []*commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &cmds); err != nil {
		return nil, err
	}

	return cmds, nil
}

// CreateGuildApplicationCommand registers a command in a specific guild.
func (c *RestClient) CreateGuildApplicationCommand(appID, guildID common.Snowflake, cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands"
	req, err := c.generateRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var registered commands.ApplicationCommand
	if _, err := c.do(req, http.StatusCreated, &registered); err != nil {
		return nil, err
	}

	return &registered, nil
}

// GetGuildApplicationCommand returns a single guild-specific application command.
func (c *RestClient) GetGuildApplicationCommand(appID, guildID, cmdID common.Snowflake) (*commands.ApplicationCommand, error) {
	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var cmd commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &cmd); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// EditGuildApplicationCommand updates a guild-specific application command.
func (c *RestClient) EditGuildApplicationCommand(appID, guildID, cmdID common.Snowflake, params *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var cmd commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &cmd); err != nil {
		return nil, err
	}

	return &cmd, nil
}

// DeleteGuildApplicationCommand deletes a guild-specific application command.
func (c *RestClient) DeleteGuildApplicationCommand(appID, guildID, cmdID common.Snowflake) error {
	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands/" + cmdID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// BulkOverwriteGuildApplicationCommands overwrites all guild-specific commands for the given guild.
// Any commands not included in cmds are deleted.
func (c *RestClient) BulkOverwriteGuildApplicationCommands(appID, guildID common.Snowflake, cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
	body, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/guilds/" + guildID.String() + "/commands"
	req, err := c.generateRequest(http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var registered []*commands.ApplicationCommand
	if _, err := c.do(req, http.StatusOK, &registered); err != nil {
		return nil, err
	}

	return registered, nil
}
