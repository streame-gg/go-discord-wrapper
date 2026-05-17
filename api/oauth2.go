package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// UpdateRoleConnectionParams holds params for updating the user's application role connection.
type UpdateRoleConnectionParams struct {
	PlatformName     *string           `json:"platform_name,omitempty"`
	PlatformUsername *string           `json:"platform_username,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// GetCurrentAuthorizationInformation returns the authorization info for the given bearer token.
// Requires an OAuth2 bearer token; bot tokens will receive a 401.
func (c *RestClient) GetCurrentAuthorizationInformation(ctx context.Context, userToken string) (*discord.OAuth2Authorization, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/oauth2/@me", nil, WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.OAuth2Authorization](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteCurrentUserApplicationRoleConnection removes the role connection for the current user.
// Requires an OAuth2 bearer token with the role_connections.write scope.
func (c *RestClient) DeleteCurrentUserApplicationRoleConnection(ctx context.Context, appID discord.Snowflake, userToken string) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	path := "/users/@me/applications/" + appID.String() + "/role-connection"
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, WithUserAuthorization(userToken))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// GetCurrentUserConnections returns the connections linked to the current user's account.
// Requires an OAuth2 bearer token with the connections scope; bot tokens will receive a 401.
func (c *RestClient) GetCurrentUserConnections(ctx context.Context, userToken string) ([]*discord.UserConnection, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/users/@me/connections", nil, WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.UserConnection](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetCurrentUserApplicationRoleConnection returns the application role connection for the current user.
// Requires an OAuth2 bearer token with the role_connections.write scope; bot tokens will receive a 401.
func (c *RestClient) GetCurrentUserApplicationRoleConnection(ctx context.Context, appID discord.Snowflake, userToken string) (*discord.ApplicationRoleConnection, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/users/@me/applications/" + appID.String() + "/role-connection"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationRoleConnection](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// UpdateCurrentUserApplicationRoleConnection updates the application role connection for the current user.
// Requires an OAuth2 bearer token with the role_connections.write scope; bot tokens will receive a 401.
func (c *RestClient) UpdateCurrentUserApplicationRoleConnection(ctx context.Context, appID discord.Snowflake, userToken string, params UpdateRoleConnectionParams) (*discord.ApplicationRoleConnection, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/users/@me/applications/" + appID.String() + "/role-connection"
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ApplicationRoleConnection](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// AddGuildMember adds a user to a guild using their OAuth2 access token.
func (c *RestClient) AddGuildMember(ctx context.Context, guildID, userID discord.Snowflake, params AddGuildMemberParams) (*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := userID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildMember](c, req, map[int]bool{
		http.StatusCreated:   true,
		http.StatusNoContent: false,
	})
}
