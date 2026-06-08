package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param / response types ────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/user#modify-current-user
type ModifyCurrentUserParams struct {
	Username *string `json:"username,omitempty"`
	// Avatar is a base64-encoded image data URI (e.g. "data:image/png;base64,...").
	Avatar *string `json:"avatar,omitempty"`
}

// https://docs.discord.com/developers/resources/user#get-current-user-guilds
type GetCurrentUserGuildsParams struct {
	Before     *discord.Snowflake
	After      *discord.Snowflake
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

// ── User endpoints ────────────────────────────────────────────────────────────

// GetCurrentUser returns the user associated with the current token.
// Pass nil for userToken to use the bot token; pass a non-nil OAuth2 bearer
// token to fetch the user that token belongs to.
// https://docs.discord.com/developers/resources/user#get-current-user
func (c *RestClient) GetCurrentUser(ctx context.Context, userToken *string) (*discord.User, error) {
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

	return doRequest[discord.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetUser returns the user object for the given user ID.
// https://docs.discord.com/developers/resources/user#get-user
func (c *RestClient) GetUser(ctx context.Context, userID discord.Snowflake) (*discord.User, error) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/users/"+userID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyCurrentUser updates the bot user's username or avatar.
// https://docs.discord.com/developers/resources/user#modify-current-user
func (c *RestClient) ModifyCurrentUser(ctx context.Context, params ModifyCurrentUserParams) (*discord.User, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/users/@me", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListCurrentUserGuilds returns the guilds the current user is a member of.
// https://docs.discord.com/developers/resources/user#get-current-user-guilds
func (c *RestClient) ListCurrentUserGuilds(ctx context.Context, params GetCurrentUserGuildsParams) ([]*discord.CurrentUserGuild, error) {
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

	return doRequestSlice[discord.CurrentUserGuild](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetCurrentUserGuildMember returns the guild member object for the current user in the given guild.
// Pass nil for userAccessToken to use the bot token; pass a non-nil OAuth2 bearer token
// (with the guilds.members.read scope) to fetch the member for that token's user.
// https://docs.discord.com/developers/resources/user#get-current-user-guild-member
func (c *RestClient) GetCurrentUserGuildMember(ctx context.Context, guildID discord.Snowflake, userAccessToken *string) (*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}
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

	return doRequest[discord.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// LeaveGuild makes the current user leave the given guild.
// https://docs.discord.com/developers/resources/user#leave-guild
func (c *RestClient) LeaveGuild(ctx context.Context, guildID discord.Snowflake) error {
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
// https://docs.discord.com/developers/resources/user#create-dm
func (c *RestClient) CreateDM(ctx context.Context, recipientID discord.Snowflake) (*discord.Channel, error) {
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

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGroupDMParams holds params for creating a Group DM.
// https://docs.discord.com/developers/resources/user#create-group-dm
type CreateGroupDMParams struct {
	AccessTokens []string          `json:"access_tokens"`
	Nicks        map[string]string `json:"nicks"`

	UserAccessToken *string
}

// CreateGroupDM creates a new Group DM channel.
// https://docs.discord.com/developers/resources/user#create-group-dm
func (c *RestClient) CreateGroupDM(ctx context.Context, params CreateGroupDMParams) (*discord.Channel, error) {
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

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
