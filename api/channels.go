package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// GetChannel fetches a channel by ID.
func (c *RestClient) GetChannel(ctx context.Context, channelID common.Snowflake) (*common.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, "GET", "/channels/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if err := c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusOK,
		Out:    &channel,
	}); err != nil {
		return nil, err
	}

	return &channel, nil
}

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyChannelParams struct {
	Name                          *string                              `json:"name,omitempty"`
	Type                          *common.ChannelType                  `json:"type,omitempty"`
	Position                      *int                                 `json:"position,omitempty"`
	Topic                         *string                              `json:"topic,omitempty"`
	NSFW                          *bool                                `json:"nsfw,omitempty"`
	RateLimitPerUser              *int                                 `json:"rate_limit_per_user,omitempty"`
	Bitrate                       *int                                 `json:"bitrate,omitempty"`
	UserLimit                     *int                                 `json:"user_limit,omitempty"`
	PermissionOverwrites          []common.ChannelPermissionOverwrites `json:"permission_overwrites,omitempty"`
	ParentID                      *common.Snowflake                    `json:"parent_id,omitempty"`
	RTCRegion                     *string                              `json:"rtc_region,omitempty"`
	VideoQualityMode              *common.VideoQualityMode             `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                                 `json:"default_auto_archive_duration,omitempty"`
	Flags                         *common.ChannelFlags                 `json:"flags,omitempty"`
	AvailableTags                 []common.ChannelTag                  `json:"available_tags,omitempty"`
	DefaultReactionEmoji          *common.DefaultReactionEmoji         `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                                 `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *common.DefaultSortOrder             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *common.ChannelForumLayoutType       `json:"default_forum_layout,omitempty"`
}

type EditChannelPermissionsParams struct {
	Allow *string                                `json:"allow,omitempty"`
	Deny  *string                                `json:"deny,omitempty"`
	Type  common.ChannelPermissionOverwritesType `json:"type"`
}

type CreateChannelInviteParams struct {
	MaxAge              *int              `json:"max_age,omitempty"`
	MaxUses             *int              `json:"max_uses,omitempty"`
	Temporary           *bool             `json:"temporary,omitempty"`
	Unique              *bool             `json:"unique,omitempty"`
	TargetType          *int              `json:"target_type,omitempty"`
	TargetUserID        *common.Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID *common.Snowflake `json:"target_application_id,omitempty"`
}

// ── Channel endpoints ─────────────────────────────────────────────────────────

// ModifyChannel updates settings for a channel. Returns the updated channel.
func (c *RestClient) ModifyChannel(ctx context.Context, channelID common.Snowflake, params ModifyChannelParams) (*common.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/channels/"+channelID.String(), bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if err := c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusOK,
		Out:    &channel,
	}); err != nil {
		return nil, err
	}

	return &channel, nil
}

// DeleteChannel deletes a channel by ID. For guild channels requires MANAGE_CHANNELS.
// Returns the deleted channel object.
func (c *RestClient) DeleteChannel(ctx context.Context, channelID common.Snowflake) (*common.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if err := c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusOK,
		Out:    &channel,
	}); err != nil {
		return nil, err
	}

	return &channel, nil
}

// GetChannelInvites returns all invites for a channel. Requires MANAGE_CHANNELS.
func (c *RestClient) GetChannelInvites(ctx context.Context, channelID common.Snowflake) ([]*Invite, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/channels/"+channelID.String()+"/invites", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var invites []*Invite
	if err := c.do(req, SuccessReturn[[]*Invite]{
		status: http.StatusOK,
		Out:    &invites,
	}); err != nil {
		return nil, err
	}

	return invites, nil
}

// CreateChannelInvite creates a new invite for a channel. Requires CREATE_INSTANT_INVITE.
func (c *RestClient) CreateChannelInvite(ctx context.Context, channelID common.Snowflake, params CreateChannelInviteParams) (*Invite, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/invites", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var invite Invite
	if err := c.do(req, SuccessReturn[Invite]{
		status: http.StatusOK,
		Out:    &invite,
	}); err != nil {
		return nil, err
	}

	return &invite, nil
}

// EditChannelPermissions creates or updates a permission overwrite for a user or role in a channel.
// Requires MANAGE_ROLES. Use params.Type to specify whether overwriteID is a role or member.
func (c *RestClient) EditChannelPermissions(ctx context.Context, channelID, overwriteID common.Snowflake, params EditChannelPermissionsParams) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := overwriteID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// DeleteChannelPermission removes a permission overwrite for a user or role from a channel.
// Requires MANAGE_ROLES.
func (c *RestClient) DeleteChannelPermission(ctx context.Context, channelID, overwriteID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := overwriteID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// TriggerTypingIndicator posts a typing indicator to the channel for ~10 seconds.
func (c *RestClient) TriggerTypingIndicator(ctx context.Context, channelID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/typing", nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// ── Additional channel endpoints ──────────────────────────────────────────────

// FollowedChannel is returned when following an announcement channel.
type FollowedChannel struct {
	ChannelID common.Snowflake `json:"channel_id"`
	WebhookID common.Snowflake `json:"webhook_id"`
}

// AddGroupDMRecipientParams holds params for adding a user to a Group DM.
type AddGroupDMRecipientParams struct {
	AccessToken string `json:"access_token"`
	Nick        string `json:"nick"`
}

// SetVoiceChannelStatus sets the status message of a voice channel.
func (c *RestClient) SetVoiceChannelStatus(ctx context.Context, channelID common.Snowflake, status *string) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(map[string]interface{}{"status": status})
	if err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/channels/"+channelID.String()+"/voice-status", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// FollowAnnouncementChannel follows an announcement channel, publishing messages to webhookChannelID.
func (c *RestClient) FollowAnnouncementChannel(ctx context.Context, channelID, webhookChannelID common.Snowflake) (*FollowedChannel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := webhookChannelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]string{"webhook_channel_id": webhookChannelID.String()})
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/followers", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var result FollowedChannel
	if err := c.do(req, SuccessReturn[FollowedChannel]{
		status: http.StatusOK,
		Out:    &result,
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddGroupDMRecipient adds a user to a Group DM using their OAuth2 access token.
func (c *RestClient) AddGroupDMRecipient(ctx context.Context, channelID, userID common.Snowflake, params AddGroupDMRecipientParams) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/recipients/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}

// RemoveGroupDMRecipient removes a user from a Group DM.
func (c *RestClient) RemoveGroupDMRecipient(ctx context.Context, channelID, userID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/recipients/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return c.do(req, SuccessReturn[common.Channel]{
		status: http.StatusNoContent,
		Out:    nil,
	})
}
