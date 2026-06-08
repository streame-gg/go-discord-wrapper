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

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/guild#list-guild-members
type GetGuildMembersParams struct {
	After *discord.Snowflake
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

// https://docs.discord.com/developers/resources/guild#search-guild-members
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

// https://docs.discord.com/developers/resources/guild#modify-guild-member
type ModifyGuildMemberParams struct {
	Nick                       *string              `json:"nick,omitempty"`
	Roles                      *[]discord.Snowflake `json:"roles,omitempty"`
	Mute                       *bool                `json:"mute,omitempty"`
	Deaf                       *bool                `json:"deaf,omitempty"`
	ChannelID                  *discord.Snowflake   `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *string              `json:"communication_disabled_until,omitempty"`
	Flags                      *int                 `json:"flags,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/guild#remove-guild-member
type KickGuildMemberOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#add-guild-member-role
type AddGuildMemberRoleOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#remove-guild-member-role
type RemoveGuildMemberRoleOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#modify-current-member
type ModifyCurrentMemberParams struct {
	Nick *string `json:"nick,omitempty"`
}

// ── Guild member endpoints ────────────────────────────────────────────────────

// GetGuildMember returns the guild member object for a user in a guild.
// https://docs.discord.com/developers/resources/guild#get-guild-member
func (c *RestClient) GetGuildMember(ctx context.Context, guildID, userID discord.Snowflake) (*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := userID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListGuildMembers returns a paginated list of members in a guild (max 1000 per request).
// https://docs.discord.com/developers/resources/guild#list-guild-members
func (c *RestClient) ListGuildMembers(ctx context.Context, guildID discord.Snowflake, params GetGuildMembersParams) ([]*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// SearchGuildMembers returns members whose username or nickname starts with the given query string.
// https://docs.discord.com/developers/resources/guild#search-guild-members
func (c *RestClient) SearchGuildMembers(ctx context.Context, guildID discord.Snowflake, params SearchGuildMembersParams) ([]*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/search" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildMember updates attributes of a guild member.
// https://docs.discord.com/developers/resources/guild#modify-guild-member
func (c *RestClient) ModifyGuildMember(ctx context.Context, guildID, userID discord.Snowflake, params ModifyGuildMemberParams) (*discord.GuildMember, error) {
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

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyCurrentMember updates the current user's member object in a guild (e.g. nickname).
// https://docs.discord.com/developers/resources/guild#modify-current-member
func (c *RestClient) ModifyCurrentMember(ctx context.Context, guildID discord.Snowflake, params ModifyCurrentMemberParams) (*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/@me"
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// AddGuildMemberRole grants a role to a guild member. Requires MANAGE_ROLES.
// https://docs.discord.com/developers/resources/guild#add-guild-member-role
func (c *RestClient) AddGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, opts *AddGuildMemberRoleOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	if err := roleID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization(), WithAuditLogReason(&opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// RemoveGuildMemberRole removes a role from a guild member. Requires MANAGE_ROLES.
// https://docs.discord.com/developers/resources/guild#remove-guild-member-role
func (c *RestClient) RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, opts *RemoveGuildMemberRoleOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	if err := roleID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(&opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// KickGuildMember kicks a member from a guild. Requires KICK_MEMBERS.
// https://docs.discord.com/developers/resources/guild#remove-guild-member
func (c *RestClient) KickGuildMember(ctx context.Context, guildID, userID discord.Snowflake, opts *KickGuildMemberOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(&opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}
