package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/guild#modify-guild
type ModifyGuildParams struct {
	Name                        *string                                  `json:"name,omitempty"`
	VerificationLevel           *discord.GuildVerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications *discord.DefaultMessageNotificationLevel `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *discord.GuildExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *discord.Snowflake                       `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int                                     `json:"afk_timeout,omitempty"`
	Icon                        *string                                  `json:"icon,omitempty"`
	OwnerID                     *discord.Snowflake                       `json:"owner_id,omitempty"`
	Splash                      *string                                  `json:"splash,omitempty"`
	DiscoverySplash             *string                                  `json:"discovery_splash,omitempty"`
	Banner                      *string                                  `json:"banner,omitempty"`
	SystemChannelID             *discord.Snowflake                       `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *discord.GuildSystemChannelFlags         `json:"system_channel_flags,omitempty"`
	RulesChannelID              *discord.Snowflake                       `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *discord.Snowflake                       `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *discord.Locale                          `json:"preferred_locale,omitempty"`
	Features                    []discord.GuildFeatures                  `json:"features,omitempty"`
	Description                 *string                                  `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                                    `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *discord.Snowflake                       `json:"safety_alerts_channel_id,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#modify-guild
type ModifyGuildOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#create-guild-channel
type CreateGuildChannelParams struct {
	Name                          string                               `json:"name"`
	Type                          *discord.ChannelType                 `json:"type,omitempty"`
	Topic                         *string                              `json:"topic,omitempty"`
	Bitrate                       *int                                 `json:"bitrate,omitempty"`
	UserLimit                     *int                                 `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                                 `json:"rate_limit_per_user,omitempty"`
	Position                      *int                                 `json:"position,omitempty"`
	PermissionOverwrites          []discord.ChannelPermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      *discord.Snowflake                   `json:"parent_id,omitempty"`
	NSFW                          *bool                                `json:"nsfw,omitempty"`
	RTCRegion                     *string                              `json:"rtc_region,omitempty"`
	VideoQualityMode              *discord.VideoQualityMode            `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                                 `json:"default_auto_archive_duration,omitempty"`
	DefaultReactionEmoji          *discord.DefaultReactionEmoji        `json:"default_reaction_emoji,omitempty"`
	AvailableTags                 []discord.ChannelTag                 `json:"available_tags,omitempty"`
	DefaultSortOrder              *discord.DefaultSortOrder            `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *discord.ChannelForumLayoutType      `json:"default_forum_layout,omitempty"`
	DefaultThreadRateLimitPerUser *int                                 `json:"default_thread_rate_limit_per_user,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#create-guild-channel
type CreateGuildChannelOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositionsOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositionsEntry struct {
	ID              discord.Snowflake  `json:"id"`
	Position        *int               `json:"position,omitempty"`
	LockPermissions *bool              `json:"lock_permissions,omitempty"`
	ParentID        *discord.Snowflake `json:"parent_id,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#create-guild-role
type CreateGuildRoleParams struct {
	Name         *string             `json:"name,omitempty"`
	Permissions  *discord.Permission `json:"permissions,omitempty"`
	Color        *int                `json:"color,omitempty"`
	Hoist        *bool               `json:"hoist,omitempty"`
	Icon         *string             `json:"icon,omitempty"`
	UnicodeEmoji *string             `json:"unicode_emoji,omitempty"`
	Mentionable  *bool               `json:"mentionable,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#create-guild-role
type CreateGuildRoleOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#modify-guild-role
type ModifyGuildRoleParams struct {
	Name         *string             `json:"name,omitempty"`
	Permissions  *discord.Permission `json:"permissions,omitempty"`
	Color        *int                `json:"color,omitempty"`
	Hoist        *bool               `json:"hoist,omitempty"`
	Icon         *string             `json:"icon,omitempty"`
	UnicodeEmoji *string             `json:"unicode_emoji,omitempty"`
	Mentionable  *bool               `json:"mentionable,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#modify-guild-role
type ModifyGuildRoleOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositionsEntry struct {
	ID       discord.Snowflake `json:"id"`
	Position *int              `json:"position,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositionsOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#delete-guild-role
type DeleteGuildRoleOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#get-guild-bans
type GetGuildBansParams struct {
	Before *discord.Snowflake
	After  *discord.Snowflake
	Limit  *int
}

func (p GetGuildBansParams) toQuery() string {
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
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// https://docs.discord.com/developers/resources/guild#create-guild-ban
type CreateGuildBanParams struct {
	DeleteMessageSeconds *int `json:"delete_message_seconds,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#create-guild-ban
type CreateGuildBanOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#remove-guild-ban
type RemoveGuildBanOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#get-guild-prune-count
type GetGuildPruneCountParams struct {
	Days         *int
	IncludeRoles []discord.Snowflake
}

func (p GetGuildPruneCountParams) toQuery() string {
	q := url.Values{}
	if p.Days != nil {
		q.Set("days", strconv.Itoa(*p.Days))
	}
	if len(p.IncludeRoles) > 0 {
		ids := make([]string, len(p.IncludeRoles))
		for i, id := range p.IncludeRoles {
			ids[i] = id.String()
		}
		q.Set("include_roles", strings.Join(ids, ","))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// https://docs.discord.com/developers/resources/guild#begin-guild-prune
type BeginGuildPruneParams struct {
	Days              *int                `json:"days,omitempty"`
	ComputePruneCount *bool               `json:"compute_prune_count,omitempty"`
	IncludeRoles      []discord.Snowflake `json:"include_roles,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#begin-guild-prune
type BeginGuildPruneOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/audit-log#get-guild-audit-log
type GetGuildAuditLogParams struct {
	UserID     *discord.Snowflake
	ActionType *discord.AuditLogActionType
	Before     *discord.Snowflake
	After      *discord.Snowflake
	Limit      *int
}

func (p GetGuildAuditLogParams) toQuery() string {
	q := url.Values{}
	if p.UserID != nil {
		q.Set("user_id", p.UserID.String())
	}
	if p.ActionType != nil {
		q.Set("action_type", strconv.Itoa(int(*p.ActionType)))
	}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
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

// https://docs.discord.com/developers/resources/guild#modify-guild-widget
type ModifyGuildWidgetParams struct {
	Enabled   *bool              `json:"enabled,omitempty"`
	ChannelID *discord.Snowflake `json:"channel_id,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreenParams struct {
	Enabled         *bool                               `json:"enabled,omitempty"`
	WelcomeChannels []discord.GuildWelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                             `json:"description,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreenOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#modify-guild-onboarding
type ModifyGuildOnboardingParams struct {
	Prompts           []discord.OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []discord.Snowflake        `json:"default_channel_ids"`
	Enabled           bool                       `json:"enabled"`
	Mode              discord.OnboardingMode     `json:"mode"`
}

// https://docs.discord.com/developers/resources/guild#modify-guild-onboarding
type ModifyGuildOnboardingOptions struct {
	Reason string
}

// ── Guild endpoints ───────────────────────────────────────────────────────────

// GetGuild returns the guild object for the given guild ID.
// Set withCounts to true to include approximate member and presence counts.
func (c *RestClient) GetGuild(ctx context.Context, guildID discord.Snowflake, withCounts bool) (*discord.Guild, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String()
	if withCounts {
		path += "?with_counts=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Guild](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildPreview returns the preview object for a public guild.
func (c *RestClient) GetGuildPreview(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/preview", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Guild](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuild updates settings for a guild. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuild(ctx context.Context, guildID discord.Snowflake, params ModifyGuildParams, opts *ModifyGuildOptions) (*discord.Guild, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String(), bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Guild](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String(), bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Guild](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuild permanently deletes a guild. The current user must be the owner.
func (c *RestClient) DeleteGuild(ctx context.Context, guildID discord.Snowflake) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/guilds/"+guildID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ── Guild channel endpoints ───────────────────────────────────────────────────

// ListGuildChannels returns all channels in a guild.
func (c *RestClient) ListGuildChannels(ctx context.Context, guildID discord.Snowflake) ([]*discord.Channel, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/channels", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildChannel creates a new channel in a guild. Requires MANAGE_CHANNELS.
func (c *RestClient) CreateGuildChannel(ctx context.Context, guildID discord.Snowflake, params CreateGuildChannelParams, opts *CreateGuildChannelOptions) (*discord.Channel, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/channels", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Channel](c, req, map[int]bool{
			http.StatusCreated: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/channels", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// ModifyGuildChannelPositions bulk-updates channel positions. Requires MANAGE_CHANNELS.
func (c *RestClient) ModifyGuildChannelPositions(ctx context.Context, guildID discord.Snowflake, entries []ModifyGuildChannelPositionsEntry, opts *ModifyGuildChannelPositionsOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/channels", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/channels", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ── Guild role endpoints ──────────────────────────────────────────────────────

// ListGuildRoles returns all roles in a guild.
func (c *RestClient) ListGuildRoles(ctx context.Context, guildID discord.Snowflake) ([]*discord.Role, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/roles", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Role](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildRoleMemberCounts returns the role.
// https://docs.discord.com/developers/resources/guild#get-guild-role-member-counts
func (c *RestClient) GetGuildRoleMemberCounts(ctx context.Context, guildID discord.Snowflake) (map[string]int64, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/roles/member-counts", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	counts, err := doRequest[map[string]int64](c, req, map[int]bool{
		http.StatusOK: true,
	})
	if err != nil {
		return nil, err
	}
	if counts == nil {
		return nil, nil
	}
	return *counts, nil
}

// GetGuildRole returns a single role in a guild.
func (c *RestClient) GetGuildRole(ctx context.Context, guildID, roleID discord.Snowflake) (*discord.Role, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := roleID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/roles/" + roleID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Role](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildRole creates a new role in a guild. Requires MANAGE_ROLES.
func (c *RestClient) CreateGuildRole(ctx context.Context, guildID discord.Snowflake, params CreateGuildRoleParams, opts *CreateGuildRoleOptions) (*discord.Role, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/roles", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Role](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/roles", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Role](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildRolePositions bulk-updates the positions of roles. Requires MANAGE_ROLES.
func (c *RestClient) ModifyGuildRolePositions(ctx context.Context, guildID discord.Snowflake, entries []ModifyGuildRolePositionsEntry, opts *ModifyGuildRolePositionsOptions) ([]*discord.Role, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if opts == nil {
		req, err = c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/roles", bytes.NewReader(body), c.WithBotAuthorization())
	} else {
		req, err = c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/roles", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	}
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Role](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildRole updates a role's settings. Requires MANAGE_ROLES.
func (c *RestClient) ModifyGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, params ModifyGuildRoleParams, opts *ModifyGuildRoleOptions) (*discord.Role, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := roleID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/roles/" + roleID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Role](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Role](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildRole deletes a role from a guild. Requires MANAGE_ROLES.
func (c *RestClient) DeleteGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, opts *DeleteGuildRoleOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := roleID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/roles/" + roleID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ListGuildBans returns a paginated list of bans in a guild. Requires BAN_MEMBERS.
func (c *RestClient) ListGuildBans(ctx context.Context, guildID discord.Snowflake, params GetGuildBansParams) ([]*discord.Ban, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/bans" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Ban](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildBan returns the ban record for a specific user. Requires BAN_MEMBERS.
func (c *RestClient) GetGuildBan(ctx context.Context, guildID, userID discord.Snowflake) (*discord.Ban, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := userID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/bans/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Ban](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildBan bans a user from a guild. Requires BAN_MEMBERS.
func (c *RestClient) CreateGuildBan(ctx context.Context, guildID, userID discord.Snowflake, params CreateGuildBanParams, opts *CreateGuildBanOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/bans/" + userID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// RemoveGuildBan lifts a ban from a user. Requires BAN_MEMBERS.
func (c *RestClient) RemoveGuildBan(ctx context.Context, guildID, userID discord.Snowflake, opts *RemoveGuildBanOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/bans/" + userID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// GetGuildPruneCount returns the number of members that would be pruned with the given criteria.
func (c *RestClient) GetGuildPruneCount(ctx context.Context, guildID discord.Snowflake, params GetGuildPruneCountParams) (*discord.GuildPruneCountResult, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/prune" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildPruneCountResult](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// BeginGuildPrune kicks inactive members. Requires KICK_MEMBERS.
func (c *RestClient) BeginGuildPrune(ctx context.Context, guildID discord.Snowflake, params BeginGuildPruneParams, opts *BeginGuildPruneOptions) (*discord.GuildPruneCountResult, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/prune", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.GuildPruneCountResult](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/prune", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildPruneCountResult](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Invite and misc endpoints ─────────────────────────────────────────────────

// ListGuildInvites returns all active invites for a guild. Requires MANAGE_GUILD.
func (c *RestClient) ListGuildInvites(ctx context.Context, guildID discord.Snowflake) ([]*discord.Invite, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/invites", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildVanityURL returns the vanity invite code and use count for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildVanityURL(ctx context.Context, guildID discord.Snowflake) (*discord.GuildVanityURL, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/vanity-url", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildVanityURL](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildAuditLog returns audit log entries for a guild. Requires VIEW_AUDIT_LOG.
func (c *RestClient) GetGuildAuditLog(ctx context.Context, guildID discord.Snowflake, params GetGuildAuditLogParams) (*discord.AuditLog, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/audit-logs" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.AuditLog](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Widget ────────────────────────────────────────────────────────────────────

// GetGuildWidgetSettings returns the widget settings for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildWidgetSettings(ctx context.Context, guildID discord.Snowflake) (*discord.GuildWidgetSettings, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/widget", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildWidgetSettings](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildWidgetSettings updates the widget settings for a guild. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuildWidgetSettings(ctx context.Context, guildID discord.Snowflake, params ModifyGuildWidgetParams) (*discord.GuildWidgetSettings, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/widget", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildWidgetSettings](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildWidgetImage returns the PNG widget image for a guild as raw bytes.
// style is an optional widget style ("shield", "banner1"–"banner4"); pass "" for the default.
func (c *RestClient) GetGuildWidgetImage(ctx context.Context, guildID discord.Snowflake, style string) ([]byte, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/widget.png"
	if style != "" {
		path += "?style=" + style
	}
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	// Widget image is PNG — use doRawBytes so the request goes through the
	// full rate-limit / retry pipeline instead of bypassing it.
	return doRawBytes(c, req)
}

// GetGuildWidget returns the public JSON widget for a guild (no authentication required).
func (c *RestClient) GetGuildWidget(ctx context.Context, guildID discord.Snowflake) (*discord.GuildWidget, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/widget.json", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildWidget](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Welcome screen ────────────────────────────────────────────────────────────

// GetGuildWelcomeScreen returns the welcome screen for a community guild.
func (c *RestClient) GetGuildWelcomeScreen(ctx context.Context, guildID discord.Snowflake) (*discord.GuildWelcomeScreen, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/welcome-screen", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildWelcomeScreen](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildWelcomeScreen updates the welcome screen for a community guild. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuildWelcomeScreen(ctx context.Context, guildID discord.Snowflake, params ModifyGuildWelcomeScreenParams, opts *ModifyGuildWelcomeScreenOptions) (*discord.GuildWelcomeScreen, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/welcome-screen", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.GuildWelcomeScreen](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/welcome-screen", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildWelcomeScreen](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Onboarding ────────────────────────────────────────────────────────────────

// GetGuildOnboarding returns the onboarding configuration for a guild.
func (c *RestClient) GetGuildOnboarding(ctx context.Context, guildID discord.Snowflake) (*discord.GuildOnboarding, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/onboarding", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildOnboarding](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildOnboarding updates the onboarding configuration for a guild. Requires MANAGE_GUILD and MANAGE_ROLES.
func (c *RestClient) ModifyGuildOnboarding(ctx context.Context, guildID discord.Snowflake, params ModifyGuildOnboardingParams, opts *ModifyGuildOnboardingOptions) (*discord.GuildOnboarding, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPut, "/guilds/"+guildID.String()+"/onboarding", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.GuildOnboarding](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/guilds/"+guildID.String()+"/onboarding", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildOnboarding](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Additional guild endpoints ─────────────────────────────────────────────────

// CreateGuildParams holds params for creating a new guild.
// https://docs.discord.com/developers/resources/guild#create-guild
type CreateGuildParams struct {
	Name                        string             `json:"name"`
	Icon                        *string            `json:"icon,omitempty"`
	VerificationLevel           *int               `json:"verification_level,omitempty"`
	DefaultMessageNotifications *int               `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *int               `json:"explicit_content_filter,omitempty"`
	Roles                       []discord.Role     `json:"roles,omitempty"`
	Channels                    []discord.Channel  `json:"channels,omitempty"`
	AFKChannelID                *discord.Snowflake `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int               `json:"afk_timeout,omitempty"`
	SystemChannelID             *discord.Snowflake `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *int               `json:"system_channel_flags,omitempty"`
}

// AddGuildMemberParams holds params for adding a member to a guild via OAuth2.
// https://docs.discord.com/developers/resources/guild#add-guild-member
type AddGuildMemberParams struct {
	AccessToken string              `json:"access_token"`
	Nick        *string             `json:"nick,omitempty"`
	Roles       []discord.Snowflake `json:"roles,omitempty"`
	Mute        *bool               `json:"mute,omitempty"`
	Deaf        *bool               `json:"deaf,omitempty"`
}

// BulkBanParams holds params for bulk banning guild members.
// https://docs.discord.com/developers/resources/guild#bulk-guild-ban
type BulkBanParams struct {
	UserIDs              []discord.Snowflake `json:"user_ids"`
	DeleteMessageSeconds *int                `json:"delete_message_seconds,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#bulk-guild-ban
type BulkBanOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/guild#delete-guild-integration
type DeleteGuildIntegrationOptions struct {
	Reason string
}

// BulkBanResult is the result of a bulk ban operation.
// https://docs.discord.com/developers/resources/guild#bulk-guild-ban
type BulkBanResult struct {
	BannedUsers []discord.Snowflake `json:"banned_users"`
	FailedUsers []discord.Snowflake `json:"failed_users"`
}

// CreateGuild creates a new guild. Only available for bots in fewer than 10 guilds.
func (c *RestClient) CreateGuild(ctx context.Context, params CreateGuildParams) (*discord.Guild, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Guild](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// BulkBanGuildMembers bans up to 200 members from a guild at once. Requires BAN_MEMBERS and MANAGE_GUILD.
func (c *RestClient) BulkBanGuildMembers(ctx context.Context, guildID discord.Snowflake, params BulkBanParams, opts *BulkBanOptions) (*BulkBanResult, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/bulk-ban", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[BulkBanResult](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/bulk-ban", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[BulkBanResult](c, req, map[int]bool{
		http.StatusOK:        true,
		http.StatusNoContent: false,
	})
}

// ── Join requests ─────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/guild
type GetGuildJoinRequestsParams struct {
	Status *discord.GuildJoinRequestApplicationStatus
	Limit  *int
	Before *discord.Snowflake
	After  *discord.Snowflake
}

func (p GetGuildJoinRequestsParams) toQuery() string {
	q := url.Values{}
	if p.Status != nil {
		q.Set("status", string(*p.Status))
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// https://docs.discord.com/developers/resources/guild
type GuildJoinRequestsResponse struct {
	Total             int                         `json:"total"`
	GuildJoinRequests []*discord.GuildJoinRequest `json:"guild_join_requests"`
}

// https://docs.discord.com/developers/resources/guild
type ActionGuildJoinRequestParams struct {
	Action          discord.GuildJoinRequestApplicationStatus `json:"action"`
	RejectionReason *string                                   `json:"rejection_reason,omitempty"`
}

// GetGuildJoinRequests returns the pending join requests for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildJoinRequests(ctx context.Context, guildID discord.Snowflake, params GetGuildJoinRequestsParams) (*GuildJoinRequestsResponse, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/requests" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[GuildJoinRequestsResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ActionGuildJoinRequest approves or rejects a guild join request. Requires MANAGE_GUILD.
func (c *RestClient) ActionGuildJoinRequest(ctx context.Context, guildID, requestID discord.Snowflake, params ActionGuildJoinRequestParams) (*discord.GuildJoinRequest, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := requestID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/requests/" + requestID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildJoinRequest](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildNewMemberWelcome returns the new member welcome configuration for a guild.
func (c *RestClient) GetGuildNewMemberWelcome(ctx context.Context, guildID discord.Snowflake) (*discord.GuildNewMemberWelcome, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/new-member-welcome"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildNewMemberWelcome](c, req, map[int]bool{
		http.StatusOK:        true,
		http.StatusNoContent: false,
	})
}

// ListGuildIntegrations returns all integrations for a guild. Requires MANAGE_GUILD.
func (c *RestClient) ListGuildIntegrations(ctx context.Context, guildID discord.Snowflake) ([]*discord.Integration, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/integrations", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Integration](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildIntegration deletes an integration from a guild. Requires MANAGE_GUILD.
func (c *RestClient) DeleteGuildIntegration(ctx context.Context, guildID, integrationID discord.Snowflake, opts *DeleteGuildIntegrationOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := integrationID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/integrations/" + integrationID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ModifyGuildIncidentActionsOptions configures Modify Guild Incident Actions.
// Each field is a timestamp up to 24 hours in the future, or nil to clear that
// pause. Because the endpoint replaces the guild's incident actions, a nil
// field is sent explicitly as JSON null (re-enabling invites or DMs).
// https://docs.discord.com/developers/resources/guild#modify-guild-incident-actions
type ModifyGuildIncidentActionsOptions struct {
	// InvitesDisabledUntil pauses invites until this time, or nil to re-enable invites.
	InvitesDisabledUntil *time.Time
	// DMsDisabledUntil pauses DMs from non-friends until this time, or nil to re-enable them.
	DMsDisabledUntil *time.Time
}

// ModifyGuildIncidentActions sets the guild's incident actions — temporarily
// pausing invites and/or DMs (each at most 24 hours out). Requires MANAGE_GUILD.
// A nil field clears the corresponding pause. Returns the updated incidents data.
//
// See https://docs.discord.com/developers/resources/guild#modify-guild-incident-actions.
func (c *RestClient) ModifyGuildIncidentActions(ctx context.Context, guildID discord.Snowflake, opts *ModifyGuildIncidentActionsOptions) (*discord.GuildIncidentsData, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}
	if opts == nil {
		opts = &ModifyGuildIncidentActionsOptions{}
	}

	// PUT replaces the incident actions, so nil fields are sent as null rather
	// than omitted (omitting would be indistinguishable from clearing here).
	body, err := json.Marshal(struct {
		InvitesDisabledUntil *time.Time `json:"invites_disabled_until"`
		DMsDisabledUntil     *time.Time `json:"dms_disabled_until"`
	}{
		InvitesDisabledUntil: opts.InvitesDisabledUntil,
		DMsDisabledUntil:     opts.DMsDisabledUntil,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/guilds/"+guildID.String()+"/incident-actions", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildIncidentsData](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
