package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// UpdateRoleConnectionParams holds params for updating the user's application role connection.
type UpdateRoleConnectionParams struct {
	PlatformName     *string           `json:"platform_name,omitempty"`
	PlatformUsername *string           `json:"platform_username,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// GetCurrentUserConnections returns the connections linked to the current user's account.
// Requires an OAuth2 bearer token with the connections scope; bot tokens will receive a 401.
func (c *RestClient) GetCurrentUserConnections(ctx context.Context, userToken string) ([]*UserConnection, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/users/@me/connections", nil, WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	var result []*UserConnection
	if _, err := c.do(req, []int{http.StatusOK}, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCurrentUserApplicationRoleConnection returns the application role connection for the current user.
// Requires an OAuth2 bearer token with the role_connections.write scope; bot tokens will receive a 401.
func (c *RestClient) GetCurrentUserApplicationRoleConnection(ctx context.Context, appID common.Snowflake, userToken string) (*ApplicationRoleConnection, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/users/@me/applications/" + appID.String() + "/role-connection"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	var result ApplicationRoleConnection
	if _, err := c.do(req, []int{http.StatusOK}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateCurrentUserApplicationRoleConnection updates the application role connection for the current user.
// Requires an OAuth2 bearer token with the role_connections.write scope; bot tokens will receive a 401.
func (c *RestClient) UpdateCurrentUserApplicationRoleConnection(ctx context.Context, appID common.Snowflake, userToken string, params UpdateRoleConnectionParams) (*ApplicationRoleConnection, error) {
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

	var result ApplicationRoleConnection
	if _, err := c.do(req, []int{http.StatusOK}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
