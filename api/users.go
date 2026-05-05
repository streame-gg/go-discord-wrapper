package api

import (
	"bytes"
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
func (c *RestClient) GetCurrentUser() (*common.User, error) {
	req, err := c.generateRequest(http.MethodGet, "/users/@me", nil)
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
func (c *RestClient) GetUser(userID common.Snowflake) (*common.User, error) {
	req, err := c.generateRequest(http.MethodGet, "/users/"+userID.String(), nil)
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
func (c *RestClient) ModifyCurrentUser(params ModifyCurrentUserParams) (*common.User, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPatch, "/users/@me", bytes.NewReader(body))
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
func (c *RestClient) GetCurrentUserGuilds(params GetCurrentUserGuildsParams) ([]*CurrentUserGuild, error) {
	path := "/users/@me/guilds" + params.toQuery()
	req, err := c.generateRequest(http.MethodGet, path, nil)
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
func (c *RestClient) GetCurrentUserGuildMember(guildID common.Snowflake) (*common.GuildMember, error) {
	path := "/users/@me/guilds/" + guildID.String() + "/member"
	req, err := c.generateRequest(http.MethodGet, path, nil)
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
func (c *RestClient) LeaveGuild(guildID common.Snowflake) error {
	req, err := c.generateRequest(http.MethodDelete, "/users/@me/guilds/"+guildID.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// CreateDM opens a DM channel with the given user and returns it.
func (c *RestClient) CreateDM(recipientID common.Snowflake) (*common.Channel, error) {
	body, err := json.Marshal(map[string]string{"recipient_id": recipientID.String()})
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPost, "/users/@me/channels", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}
