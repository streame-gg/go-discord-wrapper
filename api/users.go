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

// ── Param / response types ────────────────────────────────────────────────────

type ModifyCurrentUserParams struct {
	Username *string `json:"username,omitempty"`
	// Avatar is a base64-encoded image data URI (e.g. "data:image/png;base64,...").
	Avatar *string `json:"avatar,omitempty"`
}

type GetCurrentUserGuildsParams struct {
	Before     *common.Snowflake
	After      *common.Snowflake
	Limit      *int
	WithCounts *bool
}

func (p GetCurrentUserGuildsParams) toQuery() string {
	q := url.Values{}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.WithCounts != nil && *p.WithCounts {
		q.Set("with_counts", "true")
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// CurrentUserGuild is the partial guild object returned by GetCurrentUserGuilds.
type CurrentUserGuild struct {
	ID                       common.Snowflake       `json:"id"`
	Name                     string                 `json:"name"`
	IconHash                 *string                `json:"icon,omitempty"`
	Owner                    bool                   `json:"owner"`
	Permissions              string                 `json:"permissions"`
	Features                 []common.GuildFeatures `json:"features"`
	ApproximateMemberCount   *int                   `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount *int                   `json:"approximate_presence_count,omitempty"`
}

// ── User endpoints ────────────────────────────────────────────────────────────

// GetCurrentUser returns the bot user associated with the current token.
func (c *RestClient) GetCurrentUser(ctx context.Context) (*common.User, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/users/@me", nil)
	if err != nil {
		return nil, err
	}

	var user common.User
	if _, err := c.do(req, http.StatusOK, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUser returns the user object for the given user ID.
func (c *RestClient) GetUser(ctx context.Context, userID common.Snowflake) (*common.User, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/users/"+userID.String(), nil)
	if err != nil {
		return nil, err
	}

	var user common.User
	if _, err := c.do(req, http.StatusOK, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ModifyCurrentUser updates the bot user's username or avatar.
func (c *RestClient) ModifyCurrentUser(ctx context.Context, params ModifyCurrentUserParams) (*common.User, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/users/@me", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var user common.User
	if _, err := c.do(req, http.StatusOK, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetCurrentUserGuilds returns the guilds the current user is a member of.
func (c *RestClient) GetCurrentUserGuilds(ctx context.Context, params GetCurrentUserGuildsParams) ([]*CurrentUserGuild, error) {
	path := "/users/@me/guilds" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var guilds []*CurrentUserGuild
	if _, err := c.do(req, http.StatusOK, &guilds); err != nil {
		return nil, err
	}

	return guilds, nil
}

// GetCurrentUserGuildMember returns the guild member object for the current user in the given guild.
func (c *RestClient) GetCurrentUserGuildMember(ctx context.Context, guildID common.Snowflake) (*common.GuildMember, error) {
	path := "/users/@me/guilds/" + guildID.String() + "/member"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var member common.GuildMember
	if _, err := c.do(req, http.StatusOK, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// LeaveGuild makes the current user leave the given guild.
func (c *RestClient) LeaveGuild(ctx context.Context, guildID common.Snowflake) error {
	req, err := c.generateRequest(ctx, http.MethodDelete, "/users/@me/guilds/"+guildID.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// CreateDM opens a DM channel with the given user and returns it.
func (c *RestClient) CreateDM(ctx context.Context, recipientID common.Snowflake) (*common.Channel, error) {
	body, err := json.Marshal(map[string]string{"recipient_id": recipientID.String()})
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/users/@me/channels", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// ── Additional user endpoints ─────────────────────────────────────────────────

// UserConnection represents an external account linked to a Discord user.
type UserConnection struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Type         string                `json:"type"`
	Revoked      *bool                 `json:"revoked,omitempty"`
	Integrations []*common.Integration `json:"integrations,omitempty"`
	Verified     bool                  `json:"verified"`
	FriendSync   bool                  `json:"friend_sync"`
	ShowActivity bool                  `json:"show_activity"`
	TwoWayLink   bool                  `json:"two_way_link"`
	Visibility   int                   `json:"visibility"`
}

// ApplicationRoleConnection represents the user's role connection for an application.
type ApplicationRoleConnection struct {
	PlatformName     *string           `json:"platform_name,omitempty"`
	PlatformUsername *string           `json:"platform_username,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// UpdateRoleConnectionParams holds params for updating the user's application role connection.
type UpdateRoleConnectionParams struct {
	PlatformName     *string           `json:"platform_name,omitempty"`
	PlatformUsername *string           `json:"platform_username,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// CreateGroupDMParams holds params for creating a Group DM.
type CreateGroupDMParams struct {
	AccessTokens []string          `json:"access_tokens"`
	Nicks        map[string]string `json:"nicks"`
}

// GetCurrentUserConnections returns the connections linked to the current user's account.
func (c *RestClient) GetCurrentUserConnections(ctx context.Context) ([]*UserConnection, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/users/@me/connections", nil)
	if err != nil {
		return nil, err
	}

	var result []*UserConnection
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCurrentUserApplicationRoleConnection returns the application role connection for the current user.
func (c *RestClient) GetCurrentUserApplicationRoleConnection(ctx context.Context, appID common.Snowflake) (*ApplicationRoleConnection, error) {
	path := "/users/@me/applications/" + appID.String() + "/role-connection"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ApplicationRoleConnection
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateCurrentUserApplicationRoleConnection updates the application role connection for the current user.
func (c *RestClient) UpdateCurrentUserApplicationRoleConnection(ctx context.Context, appID common.Snowflake, params UpdateRoleConnectionParams) (*ApplicationRoleConnection, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/users/@me/applications/" + appID.String() + "/role-connection"
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result ApplicationRoleConnection
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateGroupDM creates a new Group DM channel.
func (c *RestClient) CreateGroupDM(ctx context.Context, params CreateGroupDMParams) (*common.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/users/@me/channels", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}
