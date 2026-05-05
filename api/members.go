package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type GetGuildMembersParams struct {
	After *common.Snowflake
	Limit *int
}

func (p GetGuildMembersParams) toQuery() string {
	q := url.Values{}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

type SearchGuildMembersParams struct {
	Query string
	Limit *int
}

func (p SearchGuildMembersParams) toQuery() string {
	q := url.Values{}
	q.Set("query", p.Query)
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	return "?" + q.Encode()
}

type ModifyGuildMemberParams struct {
	Nick                       *string            `json:"nick,omitempty"`
	Roles                      []common.Snowflake `json:"roles,omitempty"`
	Mute                       *bool              `json:"mute,omitempty"`
	Deaf                       *bool              `json:"deaf,omitempty"`
	ChannelID                  *common.Snowflake  `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *string            `json:"communication_disabled_until,omitempty"`
	Flags                      *int               `json:"flags,omitempty"`
}

type ModifyCurrentMemberParams struct {
	Nick *string `json:"nick,omitempty"`
}

// ── Guild member endpoints ────────────────────────────────────────────────────

// GetGuildMember returns the guild member object for a user in a guild.
func (c *RestClient) GetGuildMember(guildID, userID common.Snowflake) (*common.GuildMember, error) {
	path := "/guilds/" + guildID.String() + "/members/" + userID.String()
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

// ListGuildMembers returns a paginated list of members in a guild (max 1000 per request).
func (c *RestClient) ListGuildMembers(guildID common.Snowflake, params GetGuildMembersParams) ([]*common.GuildMember, error) {
	path := "/guilds/" + guildID.String() + "/members" + params.toQuery()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var members []*common.GuildMember
	if _, err := c.do(req, http.StatusOK, &members); err != nil {
		return nil, err
	}

	return members, nil
}

// SearchGuildMembers returns members whose username or nickname starts with the given query string.
func (c *RestClient) SearchGuildMembers(guildID common.Snowflake, params SearchGuildMembersParams) ([]*common.GuildMember, error) {
	path := "/guilds/" + guildID.String() + "/members/search" + params.toQuery()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var members []*common.GuildMember
	if _, err := c.do(req, http.StatusOK, &members); err != nil {
		return nil, err
	}

	return members, nil
}

// ModifyGuildMember updates attributes of a guild member.
func (c *RestClient) ModifyGuildMember(guildID, userID common.Snowflake, params ModifyGuildMemberParams) (*common.GuildMember, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var member common.GuildMember
	if _, err := c.do(req, http.StatusOK, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// ModifyCurrentMember updates the current user's member object in a guild (e.g. nickname).
func (c *RestClient) ModifyCurrentMember(guildID common.Snowflake, params ModifyCurrentMemberParams) (*common.GuildMember, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/@me"
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var member common.GuildMember
	if _, err := c.do(req, http.StatusOK, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// AddGuildMemberRole grants a role to a guild member. Requires MANAGE_ROLES.
func (c *RestClient) AddGuildMemberRole(guildID, userID, roleID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String()
	req, err := c.generateRequest(http.MethodPut, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// RemoveGuildMemberRole removes a role from a guild member. Requires MANAGE_ROLES.
func (c *RestClient) RemoveGuildMemberRole(guildID, userID, roleID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// RemoveGuildMember kicks a member from a guild. Requires KICK_MEMBERS.
func (c *RestClient) RemoveGuildMember(guildID, userID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/members/" + userID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
