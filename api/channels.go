package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// GetChannel fetches a channel by ID.
func (c *RestClient) GetChannel(ctx context.Context, id common.Snowflake) (*common.Channel, error) {
	req, err := c.generateRequest(ctx, "GET", "/channels/"+id.String(), nil)
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
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
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/channels/"+channelID.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// DeleteChannel deletes a channel by ID. For guild channels requires MANAGE_CHANNELS.
// Returns the deleted channel object.
func (c *RestClient) DeleteChannel(ctx context.Context, channelID common.Snowflake) (*common.Channel, error) {
	req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String(), nil)
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// GetChannelInvites returns all invites for a channel. Requires MANAGE_CHANNELS.
func (c *RestClient) GetChannelInvites(ctx context.Context, channelID common.Snowflake) ([]*Invite, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/channels/"+channelID.String()+"/invites", nil)
	if err != nil {
		return nil, err
	}

	var invites []*Invite
	if _, err := c.do(req, http.StatusOK, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

// CreateChannelInvite creates a new invite for a channel. Requires CREATE_INSTANT_INVITE.
func (c *RestClient) CreateChannelInvite(ctx context.Context, channelID common.Snowflake, params CreateChannelInviteParams) (*Invite, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/invites", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var invite Invite
	if _, err := c.do(req, http.StatusOK, &invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

// EditChannelPermissions creates or updates a permission overwrite for a user or role in a channel.
// Requires MANAGE_ROLES. Use params.Type to specify whether overwriteID is a role or member.
func (c *RestClient) EditChannelPermissions(ctx context.Context, channelID, overwriteID common.Snowflake, params EditChannelPermissionsParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// DeleteChannelPermission removes a permission overwrite for a user or role from a channel.
// Requires MANAGE_ROLES.
func (c *RestClient) DeleteChannelPermission(ctx context.Context, channelID, overwriteID common.Snowflake) error {
	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// TriggerTypingIndicator posts a typing indicator to the channel for ~10 seconds.
func (c *RestClient) TriggerTypingIndicator(ctx context.Context, channelID common.Snowflake) error {
	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/typing", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
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
	body, err := json.Marshal(map[string]interface{}{"status": status})
	if err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/channels/"+channelID.String()+"/voice-status", bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// FollowAnnouncementChannel follows an announcement channel, publishing messages to webhookChannelID.
func (c *RestClient) FollowAnnouncementChannel(ctx context.Context, channelID, webhookChannelID common.Snowflake) (*FollowedChannel, error) {
	body, err := json.Marshal(map[string]string{"webhook_channel_id": webhookChannelID.String()})
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/followers", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result FollowedChannel
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddGroupDMRecipient adds a user to a Group DM using their OAuth2 access token.
func (c *RestClient) AddGroupDMRecipient(ctx context.Context, channelID, userID common.Snowflake, params AddGroupDMRecipientParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/recipients/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// RemoveGroupDMRecipient removes a user from a Group DM.
func (c *RestClient) RemoveGroupDMRecipient(ctx context.Context, channelID, userID common.Snowflake) error {
	path := "/channels/" + channelID.String() + "/recipients/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
