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

	UserToken *string
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
func (c *RestClient) GetCurrentUser(ctx context.Context, userToken *string) (*common.User, error) {
	var authOption func(req *http.Request)
	if userToken != nil {
		authOption = WithUserAuthorization(*userToken)
	} else {
		authOption = c.WithBotAuthorization()
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/users/@me", nil, authOption)
	if err != nil {
		return nil, err
	}

	return doRequest[common.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetUser returns the user object for the given user ID.
func (c *RestClient) GetUser(ctx context.Context, userID common.Snowflake) (*common.User, error) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/users/"+userID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyCurrentUser updates the bot user's username or avatar.
func (c *RestClient) ModifyCurrentUser(ctx context.Context, params ModifyCurrentUserParams) (*common.User, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/users/@me", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetCurrentUserGuilds returns the guilds the current user is a member of.
func (c *RestClient) GetCurrentUserGuilds(ctx context.Context, params GetCurrentUserGuildsParams) (*[]*CurrentUserGuild, error) {
	var authOption func(req *http.Request)
	if params.UserToken != nil {
		authOption = WithUserAuthorization(*params.UserToken)
	} else {
		authOption = c.WithBotAuthorization()
	}

	path := "/users/@me/guilds" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, authOption)
	if err != nil {
		return nil, err
	}

	return doRequest[[]*CurrentUserGuild](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetCurrentUserGuildMember returns the guild member object for the current user in the given guild.
func (c *RestClient) GetCurrentUserGuildMember(ctx context.Context, guildID common.Snowflake, userAccessToken *string) (*common.GuildMember, error) {
	var authOption func(req *http.Request)
	if userAccessToken != nil {
		authOption = WithUserAuthorization(*userAccessToken)
	} else {
		authOption = c.WithBotAuthorization()
	}

	path := "/users/@me/guilds/" + guildID.String() + "/member"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, authOption)
	if err != nil {
		return nil, err
	}

	return doRequest[common.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// LeaveGuild makes the current user leave the given guild.
func (c *RestClient) LeaveGuild(ctx context.Context, guildID common.Snowflake) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/users/@me/guilds/"+guildID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// CreateDM opens a DM channel with the given user and returns it.
func (c *RestClient) CreateDM(ctx context.Context, recipientID common.Snowflake) (*common.Channel, error) {
	if err := recipientID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]string{"recipient_id": recipientID.String()})
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/users/@me/channels", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

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

// CreateGroupDMParams holds params for creating a Group DM.
type CreateGroupDMParams struct {
	AccessTokens []string          `json:"access_tokens"`
	Nicks        map[string]string `json:"nicks"`

	UserAccessToken *string
}

// CreateGroupDM creates a new Group DM channel.
func (c *RestClient) CreateGroupDM(ctx context.Context, params CreateGroupDMParams) (*common.Channel, error) {
	var authOption func(req *http.Request)
	if params.UserAccessToken != nil {
		authOption = WithUserAuthorization(*params.UserAccessToken)
	} else {
		authOption = c.WithBotAuthorization()
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/users/@me/channels", bytes.NewReader(body), authOption)
	if err != nil {
		return nil, err
	}

	return doRequest[common.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
