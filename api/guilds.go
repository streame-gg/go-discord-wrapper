package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyGuildParams struct {
	Name                        *string                                 `json:"name,omitempty"`
	VerificationLevel           *common.GuildVerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications *common.DefaultMessageNotificationLevel `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *common.GuildExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *common.Snowflake                       `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int                                    `json:"afk_timeout,omitempty"`
	Icon                        *string                                 `json:"icon,omitempty"`
	OwnerID                     *common.Snowflake                       `json:"owner_id,omitempty"`
	Splash                      *string                                 `json:"splash,omitempty"`
	DiscoverySplash             *string                                 `json:"discovery_splash,omitempty"`
	Banner                      *string                                 `json:"banner,omitempty"`
	SystemChannelID             *common.Snowflake                       `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *common.GuildSystemChannelFlags         `json:"system_channel_flags,omitempty"`
	RulesChannelID              *common.Snowflake                       `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *common.Snowflake                       `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *common.Locale                          `json:"preferred_locale,omitempty"`
	Features                    []common.GuildFeatures                  `json:"features,omitempty"`
	Description                 *string                                 `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                                   `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *common.Snowflake                       `json:"safety_alerts_channel_id,omitempty"`
}

type CreateGuildChannelParams struct {
	Name                          string                               `json:"name"`
	Type                          *common.ChannelType                  `json:"type,omitempty"`
	Topic                         *string                              `json:"topic,omitempty"`
	Bitrate                       *int                                 `json:"bitrate,omitempty"`
	UserLimit                     *int                                 `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                                 `json:"rate_limit_per_user,omitempty"`
	Position                      *int                                 `json:"position,omitempty"`
	PermissionOverwrites          []common.ChannelPermissionOverwrites `json:"permission_overwrites,omitempty"`
	ParentID                      *common.Snowflake                    `json:"parent_id,omitempty"`
	NSFW                          *bool                                `json:"nsfw,omitempty"`
	RTCRegion                     *string                              `json:"rtc_region,omitempty"`
	VideoQualityMode              *common.VideoQualityMode             `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                                 `json:"default_auto_archive_duration,omitempty"`
	DefaultReactionEmoji          *common.DefaultReactionEmoji         `json:"default_reaction_emoji,omitempty"`
	AvailableTags                 []common.ChannelTag                  `json:"available_tags,omitempty"`
	DefaultSortOrder              *common.DefaultSortOrder             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *common.ChannelForumLayoutType       `json:"default_forum_layout,omitempty"`
	DefaultThreadRateLimitPerUser *int                                 `json:"default_thread_rate_limit_per_user,omitempty"`
}

type ModifyGuildChannelPositionsEntry struct {
	ID              common.Snowflake  `json:"id"`
	Position        *int              `json:"position,omitempty"`
	LockPermissions *bool             `json:"lock_permissions,omitempty"`
	ParentID        *common.Snowflake `json:"parent_id,omitempty"`
}

type CreateGuildRoleParams struct {
	Name         *string `json:"name,omitempty"`
	Permissions  *string `json:"permissions,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        *bool   `json:"hoist,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  *bool   `json:"mentionable,omitempty"`
}

type ModifyGuildRoleParams struct {
	Name         *string `json:"name,omitempty"`
	Permissions  *string `json:"permissions,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        *bool   `json:"hoist,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  *bool   `json:"mentionable,omitempty"`
}

type ModifyGuildRolePositionsEntry struct {
	ID       common.Snowflake `json:"id"`
	Position *int             `json:"position,omitempty"`
}

type GetGuildBansParams struct {
	Before *common.Snowflake
	After  *common.Snowflake
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

type CreateGuildBanParams struct {
	DeleteMessageSeconds *int `json:"delete_message_seconds,omitempty"`
}

type GetGuildPruneCountParams struct {
	Days         *int
	IncludeRoles []common.Snowflake
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

type BeginGuildPruneParams struct {
	Days              *int               `json:"days,omitempty"`
	ComputePruneCount *bool              `json:"compute_prune_count,omitempty"`
	IncludeRoles      []common.Snowflake `json:"include_roles,omitempty"`
}

type GetGuildAuditLogParams struct {
	UserID     *common.Snowflake
	ActionType *common.AuditLogActionType
	Before     *common.Snowflake
	After      *common.Snowflake
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

// ── Response types ────────────────────────────────────────────────────────────

type GuildPruneCountResult struct {
	Pruned *int `json:"pruned"`
}

type ModifyGuildWidgetParams struct {
	Enabled   *bool             `json:"enabled,omitempty"`
	ChannelID *common.Snowflake `json:"channel_id,omitempty"`
}

type ModifyGuildWelcomeScreenParams struct {
	Enabled         *bool                              `json:"enabled,omitempty"`
	WelcomeChannels []common.GuildWelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                            `json:"description,omitempty"`
}

type ModifyGuildOnboardingParams struct {
	Prompts           []common.OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []common.Snowflake        `json:"default_channel_ids"`
	Enabled           bool                      `json:"enabled"`
	Mode              common.OnboardingMode     `json:"mode"`
}

type AuditLog struct {
	ApplicationCommands  []any                  `json:"application_commands"`
	AuditLogEntries      []common.AuditLogEntry `json:"audit_log_entries"`
	AutoModerationRules  []any                  `json:"auto_moderation_rules"`
	GuildScheduledEvents []any                  `json:"guild_scheduled_events"`
	Integrations         []any                  `json:"integrations"`
	Threads              []common.Channel       `json:"threads"`
	Users                []common.User          `json:"users"`
	Webhooks             []any                  `json:"webhooks"`
}

type GuildVanityURL struct {
	Code *string `json:"code"`
	Uses int     `json:"uses"`
}

type Ban struct {
	Reason *string     `json:"reason"`
	User   common.User `json:"user"`
}

type Invite struct {
	Code                     string                       `json:"code"`
	Guild                    *common.Guild                `json:"guild,omitempty"`
	Channel                  *common.Channel              `json:"channel,omitempty"`
	Inviter                  *common.User                 `json:"inviter,omitempty"`
	TargetType               *common.InviteTargetUserType `json:"target_type,omitempty"`
	TargetUser               *common.User                 `json:"target_user,omitempty"`
	ApproximatePresenceCount *int                         `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   *int                         `json:"approximate_member_count,omitempty"`
	ExpiresAt                *time.Time                   `json:"expires_at,omitempty"`
	Flags                    *common.InviteFlags          `json:"flags,omitempty"`
	Uses                     *int                         `json:"uses,omitempty"`
	MaxUses                  *int                         `json:"max_uses,omitempty"`
	MaxAge                   *int                         `json:"max_age,omitempty"`
	Temporary                *bool                        `json:"temporary,omitempty"`
	CreatedAt                *time.Time                   `json:"created_at,omitempty"`
}

// ── Guild endpoints ───────────────────────────────────────────────────────────

// GetGuild returns the guild object for the given guild ID.
// Set withCounts to true to include approximate member and presence counts.
func (c *RestClient) GetGuild(ctx context.Context, guildID common.Snowflake, withCounts bool) (*common.Guild, error) {
	path := "/guilds/" + guildID.String()
	if withCounts {
		path += "?with_counts=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var guild common.Guild
	if _, err := c.do(req, http.StatusOK, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

// GetGuildPreview returns the preview object for a public guild.
func (c *RestClient) GetGuildPreview(ctx context.Context, guildID common.Snowflake) (*common.Guild, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/preview", nil)
	if err != nil {
		return nil, err
	}

	var guild common.Guild
	if _, err := c.do(req, http.StatusOK, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

// ModifyGuild updates settings for a guild. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuild(ctx context.Context, guildID common.Snowflake, params ModifyGuildParams) (*common.Guild, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var guild common.Guild
	if _, err := c.do(req, http.StatusOK, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

// DeleteGuild permanently deletes a guild. The current user must be the owner.
func (c *RestClient) DeleteGuild(ctx context.Context, guildID common.Snowflake) error {
	req, err := c.generateRequest(ctx, http.MethodDelete, "/guilds/"+guildID.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Guild channel endpoints ───────────────────────────────────────────────────

// GetGuildChannels returns all channels in a guild.
func (c *RestClient) GetGuildChannels(ctx context.Context, guildID common.Snowflake) ([]*common.Channel, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/channels", nil)
	if err != nil {
		return nil, err
	}

	var channels []*common.Channel
	if _, err := c.do(req, http.StatusOK, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

// CreateGuildChannel creates a new channel in a guild. Requires MANAGE_CHANNELS.
func (c *RestClient) CreateGuildChannel(ctx context.Context, guildID common.Snowflake, params CreateGuildChannelParams) (*common.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/channels", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// ModifyGuildChannelPositions bulk-updates channel positions. Requires MANAGE_CHANNELS.
func (c *RestClient) ModifyGuildChannelPositions(ctx context.Context, guildID common.Snowflake, entries []ModifyGuildChannelPositionsEntry) error {
	body, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/channels", bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Guild role endpoints ──────────────────────────────────────────────────────

// GetGuildRoles returns all roles in a guild.
func (c *RestClient) GetGuildRoles(ctx context.Context, guildID common.Snowflake) ([]*common.Role, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/roles", nil)
	if err != nil {
		return nil, err
	}

	var roles []*common.Role
	if _, err := c.do(req, http.StatusOK, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetGuildRole returns a single role in a guild.
func (c *RestClient) GetGuildRole(ctx context.Context, guildID, roleID common.Snowflake) (*common.Role, error) {
	path := "/guilds/" + guildID.String() + "/roles/" + roleID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var role common.Role
	if _, err := c.do(req, http.StatusOK, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// CreateGuildRole creates a new role in a guild. Requires MANAGE_ROLES.
func (c *RestClient) CreateGuildRole(ctx context.Context, guildID common.Snowflake, params CreateGuildRoleParams) (*common.Role, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/roles", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var role common.Role
	if _, err := c.do(req, http.StatusOK, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// ModifyGuildRolePositions bulk-updates the positions of roles. Requires MANAGE_ROLES.
func (c *RestClient) ModifyGuildRolePositions(ctx context.Context, guildID common.Snowflake, entries []ModifyGuildRolePositionsEntry) ([]*common.Role, error) {
	body, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/roles", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var roles []*common.Role
	if _, err := c.do(req, http.StatusOK, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// ModifyGuildRole updates a role's settings. Requires MANAGE_ROLES.
func (c *RestClient) ModifyGuildRole(ctx context.Context, guildID, roleID common.Snowflake, params ModifyGuildRoleParams) (*common.Role, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/roles/" + roleID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var role common.Role
	if _, err := c.do(req, http.StatusOK, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// DeleteGuildRole deletes a role from a guild. Requires MANAGE_ROLES.
func (c *RestClient) DeleteGuildRole(ctx context.Context, guildID, roleID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/roles/" + roleID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Ban endpoints ─────────────────────────────────────────────────────────────

// GetGuildBans returns a paginated list of bans in a guild. Requires BAN_MEMBERS.
func (c *RestClient) GetGuildBans(ctx context.Context, guildID common.Snowflake, params GetGuildBansParams) ([]*Ban, error) {
	path := "/guilds/" + guildID.String() + "/bans" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var bans []*Ban
	if _, err := c.do(req, http.StatusOK, &bans); err != nil {
		return nil, err
	}

	return bans, nil
}

// GetGuildBan returns the ban record for a specific user. Requires BAN_MEMBERS.
func (c *RestClient) GetGuildBan(ctx context.Context, guildID, userID common.Snowflake) (*Ban, error) {
	path := "/guilds/" + guildID.String() + "/bans/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var ban Ban
	if _, err := c.do(req, http.StatusOK, &ban); err != nil {
		return nil, err
	}

	return &ban, nil
}

// CreateGuildBan bans a user from a guild. Requires BAN_MEMBERS.
func (c *RestClient) CreateGuildBan(ctx context.Context, guildID, userID common.Snowflake, params CreateGuildBanParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/bans/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// RemoveGuildBan lifts a ban from a user. Requires BAN_MEMBERS.
func (c *RestClient) RemoveGuildBan(ctx context.Context, guildID, userID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/bans/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Prune endpoints ───────────────────────────────────────────────────────────

// GetGuildPruneCount returns the number of members that would be pruned with the given criteria.
func (c *RestClient) GetGuildPruneCount(ctx context.Context, guildID common.Snowflake, params GetGuildPruneCountParams) (*GuildPruneCountResult, error) {
	path := "/guilds/" + guildID.String() + "/prune" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result GuildPruneCountResult
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// BeginGuildPrune kicks inactive members. Requires KICK_MEMBERS.
func (c *RestClient) BeginGuildPrune(ctx context.Context, guildID common.Snowflake, params BeginGuildPruneParams) (*GuildPruneCountResult, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/prune", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result GuildPruneCountResult
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ── Invite and misc endpoints ─────────────────────────────────────────────────

// GetGuildInvites returns all active invites for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildInvites(ctx context.Context, guildID common.Snowflake) ([]*Invite, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/invites", nil)
	if err != nil {
		return nil, err
	}

	var invites []*Invite
	if _, err := c.do(req, http.StatusOK, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

// GetGuildVanityURL returns the vanity invite code and use count for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildVanityURL(ctx context.Context, guildID common.Snowflake) (*GuildVanityURL, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/vanity-url", nil)
	if err != nil {
		return nil, err
	}

	var vanity GuildVanityURL
	if _, err := c.do(req, http.StatusOK, &vanity); err != nil {
		return nil, err
	}

	return &vanity, nil
}

// GetGuildAuditLog returns audit log entries for a guild. Requires VIEW_AUDIT_LOG.
func (c *RestClient) GetGuildAuditLog(ctx context.Context, guildID common.Snowflake, params GetGuildAuditLogParams) (*AuditLog, error) {
	path := "/guilds/" + guildID.String() + "/audit-logs" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var log AuditLog
	if _, err := c.do(req, http.StatusOK, &log); err != nil {
		return nil, err
	}

	return &log, nil
}

// ── Widget ────────────────────────────────────────────────────────────────────

// GetGuildWidgetSettings returns the widget settings for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildWidgetSettings(ctx context.Context, guildID common.Snowflake) (*common.GuildWidgetSettings, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/widget", nil)
	if err != nil {
		return nil, err
	}

	var settings common.GuildWidgetSettings
	if _, err := c.do(req, http.StatusOK, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// ModifyGuildWidgetSettings updates the widget settings for a guild. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuildWidgetSettings(ctx context.Context, guildID common.Snowflake, params ModifyGuildWidgetParams) (*common.GuildWidgetSettings, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/widget", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var settings common.GuildWidgetSettings
	if _, err := c.do(req, http.StatusOK, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// GetGuildWidgetImage returns the PNG widget image for a guild as raw bytes.
// style is an optional widget style ("shield", "banner1"–"banner4"); pass "" for the default.
func (c *RestClient) GetGuildWidgetImage(ctx context.Context, guildID common.Snowflake, style string) ([]byte, error) {
	path := "/guilds/" + guildID.String() + "/widget.png"
	if style != "" {
		path += "?style=" + style
	}
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Widget image is PNG, not JSON — bypass the normal do() JSON decoder.
	resp, reqErr := c.httpClient.Do(req)
	if reqErr != nil {
		return nil, reqErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeGatewayError(resp)
	}
	return io.ReadAll(resp.Body)
}

// GetGuildWidget returns the public JSON widget for a guild (no authentication required).
func (c *RestClient) GetGuildWidget(ctx context.Context, guildID common.Snowflake) (*common.GuildWidget, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/widget.json", nil)
	if err != nil {
		return nil, err
	}

	var widget common.GuildWidget
	if _, err := c.do(req, http.StatusOK, &widget); err != nil {
		return nil, err
	}

	return &widget, nil
}

// ── Welcome screen ────────────────────────────────────────────────────────────

// GetGuildWelcomeScreen returns the welcome screen for a community guild.
func (c *RestClient) GetGuildWelcomeScreen(ctx context.Context, guildID common.Snowflake) (*common.GuildWelcomeScreen, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/welcome-screen", nil)
	if err != nil {
		return nil, err
	}

	var screen common.GuildWelcomeScreen
	if _, err := c.do(req, http.StatusOK, &screen); err != nil {
		return nil, err
	}

	return &screen, nil
}

// ModifyGuildWelcomeScreen updates the welcome screen for a community guild. Requires MANAGE_GUILD.
func (c *RestClient) ModifyGuildWelcomeScreen(ctx context.Context, guildID common.Snowflake, params ModifyGuildWelcomeScreenParams) (*common.GuildWelcomeScreen, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/welcome-screen", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var screen common.GuildWelcomeScreen
	if _, err := c.do(req, http.StatusOK, &screen); err != nil {
		return nil, err
	}

	return &screen, nil
}

// ── Onboarding ────────────────────────────────────────────────────────────────

// GetGuildOnboarding returns the onboarding configuration for a guild.
func (c *RestClient) GetGuildOnboarding(ctx context.Context, guildID common.Snowflake) (*common.GuildOnboarding, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/onboarding", nil)
	if err != nil {
		return nil, err
	}

	var onboarding common.GuildOnboarding
	if _, err := c.do(req, http.StatusOK, &onboarding); err != nil {
		return nil, err
	}

	return &onboarding, nil
}

// ModifyGuildOnboarding updates the onboarding configuration for a guild. Requires MANAGE_GUILD and MANAGE_ROLES.
func (c *RestClient) ModifyGuildOnboarding(ctx context.Context, guildID common.Snowflake, params ModifyGuildOnboardingParams) (*common.GuildOnboarding, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/guilds/"+guildID.String()+"/onboarding", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var onboarding common.GuildOnboarding
	if _, err := c.do(req, http.StatusOK, &onboarding); err != nil {
		return nil, err
	}

	return &onboarding, nil
}

// ── Additional guild endpoints ─────────────────────────────────────────────────

// CreateGuildParams holds params for creating a new guild.
type CreateGuildParams struct {
	Name                        string            `json:"name"`
	Icon                        *string           `json:"icon,omitempty"`
	VerificationLevel           *int              `json:"verification_level,omitempty"`
	DefaultMessageNotifications *int              `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *int              `json:"explicit_content_filter,omitempty"`
	Roles                       []common.Role     `json:"roles,omitempty"`
	Channels                    []common.Channel  `json:"channels,omitempty"`
	AFKChannelID                *common.Snowflake `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int              `json:"afk_timeout,omitempty"`
	SystemChannelID             *common.Snowflake `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *int              `json:"system_channel_flags,omitempty"`
}

// AddGuildMemberParams holds params for adding a member to a guild via OAuth2.
type AddGuildMemberParams struct {
	AccessToken string             `json:"access_token"`
	Nick        *string            `json:"nick,omitempty"`
	Roles       []common.Snowflake `json:"roles,omitempty"`
	Mute        *bool              `json:"mute,omitempty"`
	Deaf        *bool              `json:"deaf,omitempty"`
}

// BulkBanParams holds params for bulk banning guild members.
type BulkBanParams struct {
	UserIDs              []common.Snowflake `json:"user_ids"`
	DeleteMessageSeconds *int               `json:"delete_message_seconds,omitempty"`
}

// BulkBanResult is the result of a bulk ban operation.
type BulkBanResult struct {
	BannedUsers []common.Snowflake `json:"banned_users"`
	FailedUsers []common.Snowflake `json:"failed_users"`
}

// CreateGuild creates a new guild. Only available for bots in fewer than 10 guilds.
func (c *RestClient) CreateGuild(ctx context.Context, params CreateGuildParams) (*common.Guild, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var guild common.Guild
	if _, err := c.do(req, http.StatusCreated, &guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

// AddGuildMember adds a user to a guild using their OAuth2 access token.
func (c *RestClient) AddGuildMember(ctx context.Context, guildID, userID common.Snowflake, params AddGuildMemberParams) (*common.GuildMember, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var member common.GuildMember
	if _, err := c.do(req, http.StatusCreated, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// BulkBanGuildMembers bans up to 200 members from a guild at once. Requires BAN_MEMBERS and MANAGE_GUILD.
func (c *RestClient) BulkBanGuildMembers(ctx context.Context, guildID common.Snowflake, params BulkBanParams) (*BulkBanResult, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/bulk-ban", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result BulkBanResult
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetGuildIntegrations returns all integrations for a guild. Requires MANAGE_GUILD.
func (c *RestClient) GetGuildIntegrations(ctx context.Context, guildID common.Snowflake) ([]*common.Integration, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/integrations", nil)
	if err != nil {
		return nil, err
	}

	var result []*common.Integration
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ModifyGuildMFALevel sets the required MFA level for moderators in a guild.
// Requires guild ownership. level: 0 = none, 1 = elevated.
func (c *RestClient) ModifyGuildMFALevel(ctx context.Context, guildID common.Snowflake, level int) (int, error) {
	body, err := json.Marshal(map[string]int{"level": level})
	if err != nil {
		return 0, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/mfa", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}

	var result int
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// DeleteGuildIntegration deletes an integration from a guild. Requires MANAGE_GUILD.
func (c *RestClient) DeleteGuildIntegration(ctx context.Context, guildID, integrationID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/integrations/" + integrationID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
