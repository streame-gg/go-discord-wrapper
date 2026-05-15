package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GetChannel fetches a channel by ID.
func (c *RestClient) GetChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, "GET", "/channels/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyChannelParams struct {
	Name                          *string                               `json:"name,omitempty"`
	Type                          *discord.ChannelType                  `json:"type,omitempty"`
	Position                      *int                                  `json:"position,omitempty"`
	Topic                         *string                               `json:"topic,omitempty"`
	NSFW                          *bool                                 `json:"nsfw,omitempty"`
	RateLimitPerUser              *int                                  `json:"rate_limit_per_user,omitempty"`
	Bitrate                       *int                                  `json:"bitrate,omitempty"`
	UserLimit                     *int                                  `json:"user_limit,omitempty"`
	PermissionOverwrites          []discord.ChannelPermissionOverwrites `json:"permission_overwrites,omitempty"`
	ParentID                      *discord.Snowflake                    `json:"parent_id,omitempty"`
	RTCRegion                     *string                               `json:"rtc_region,omitempty"`
	VideoQualityMode              *discord.VideoQualityMode             `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                                  `json:"default_auto_archive_duration,omitempty"`
	Flags                         *discord.ChannelFlags                 `json:"flags,omitempty"`
	AvailableTags                 []discord.ChannelTag                  `json:"available_tags,omitempty"`
	DefaultReactionEmoji          *discord.DefaultReactionEmoji         `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                                  `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *discord.DefaultSortOrder             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *discord.ChannelForumLayoutType       `json:"default_forum_layout,omitempty"`
}

type EditChannelPermissionsParams struct {
	Allow *string                                 `json:"allow,omitempty"`
	Deny  *string                                 `json:"deny,omitempty"`
	Type  discord.ChannelPermissionOverwritesType `json:"type"`
}

type CreateChannelInviteParams struct {
	MaxAge              *int               `json:"max_age,omitempty"`
	MaxUses             *int               `json:"max_uses,omitempty"`
	Temporary           *bool              `json:"temporary,omitempty"`
	Unique              *bool              `json:"unique,omitempty"`
	TargetType          *int               `json:"target_type,omitempty"`
	TargetUserID        *discord.Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID *discord.Snowflake `json:"target_application_id,omitempty"`
}

// ── Channel endpoints ─────────────────────────────────────────────────────────

// ModifyChannel updates settings for a channel. Returns the updated channel.
func (c *RestClient) ModifyChannel(ctx context.Context, channelID discord.Snowflake, params ModifyChannelParams) (*discord.Channel, error) {
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

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteChannel deletes a channel by ID. For guild channels requires MANAGE_CHANNELS.
// Returns the deleted channel object.
func (c *RestClient) DeleteChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetChannelInvites returns all invites for a channel. Requires MANAGE_CHANNELS.
func (c *RestClient) GetChannelInvites(ctx context.Context, channelID discord.Snowflake) (*[]*Invite, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/channels/"+channelID.String()+"/invites", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateChannelInvite creates a new invite for a channel. Requires CREATE_INSTANT_INVITE.
func (c *RestClient) CreateChannelInvite(ctx context.Context, channelID discord.Snowflake, params CreateChannelInviteParams) (*Invite, error) {
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

	return doRequest[Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditChannelPermissions creates or updates a permission overwrite for a user or role in a channel.
// Requires MANAGE_ROLES. Use params.Type to specify whether overwriteID is a role or member.
func (c *RestClient) EditChannelPermissions(ctx context.Context, channelID, overwriteID discord.Snowflake, params EditChannelPermissionsParams) error {
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

	return doRequestWithoutResponse(c, req)
}

// DeleteChannelPermission removes a permission overwrite for a user or role from a channel.
// Requires MANAGE_ROLES.
func (c *RestClient) DeleteChannelPermission(ctx context.Context, channelID, overwriteID discord.Snowflake) error {
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

	return doRequestWithoutResponse(c, req)
}

// TriggerTypingIndicator posts a typing indicator to the channel for ~10 seconds.
func (c *RestClient) TriggerTypingIndicator(ctx context.Context, channelID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/typing", nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ── Additional channel endpoints ──────────────────────────────────────────────

// FollowedChannel is returned when following an announcement channel.
type FollowedChannel struct {
	ChannelID discord.Snowflake `json:"channel_id"`
	WebhookID discord.Snowflake `json:"webhook_id"`
}

// AddGroupDMRecipientParams holds params for adding a user to a Group DM.
type AddGroupDMRecipientParams struct {
	AccessToken string `json:"access_token"`
	Nick        string `json:"nick"`
}

// SetVoiceChannelStatus sets the status message of a voice channel.
func (c *RestClient) SetVoiceChannelStatus(ctx context.Context, channelID discord.Snowflake, status *string) error {
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

	_, err = doRequest[NoReturnData](c, req, map[int]bool{})

	return doRequestWithoutResponse(c, req)
}

// FollowAnnouncementChannel follows an announcement channel, publishing messages to webhookChannelID.
func (c *RestClient) FollowAnnouncementChannel(ctx context.Context, channelID, webhookChannelID discord.Snowflake) (*FollowedChannel, error) {
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

	return doRequest[FollowedChannel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// AddGroupDMRecipient adds a user to a Group DM using their OAuth2 access token.
func (c *RestClient) AddGroupDMRecipient(ctx context.Context, channelID, userID discord.Snowflake, params AddGroupDMRecipientParams) error {
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

	return doRequestWithoutResponse(c, req)
}

// RemoveGroupDMRecipient removes a user from a Group DM.
func (c *RestClient) RemoveGroupDMRecipient(ctx context.Context, channelID, userID discord.Snowflake) error {
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

	return doRequestWithoutResponse(c, req)
}
