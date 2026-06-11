package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GetChannel fetches a channel by ID.
// https://docs.discord.com/developers/resources/channel#get-channel
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

// https://docs.discord.com/developers/resources/channel#modify-channel
type ModifyChannelParams struct {
	Name                          discord.Option[string]                               `json:"name,omitempty"`
	Type                          discord.Option[discord.ChannelType]                  `json:"type,omitempty"`
	Position                      discord.Option[int]                                  `json:"position,omitempty"`
	Topic                         discord.Option[string]                               `json:"topic,omitempty"`
	NSFW                          discord.Option[bool]                                 `json:"nsfw,omitempty"`
	RateLimitPerUser              discord.Option[int]                                  `json:"rate_limit_per_user,omitempty"`
	Bitrate                       discord.Option[int]                                  `json:"bitrate,omitempty"`
	UserLimit                     discord.Option[int]                                  `json:"user_limit,omitempty"`
	PermissionOverwrites          discord.Option[[]discord.ChannelPermissionOverwrite] `json:"permission_overwrites,omitempty"`
	ParentID                      discord.Option[discord.Snowflake]                    `json:"parent_id,omitempty"`
	RTCRegion                     discord.Option[string]                               `json:"rtc_region,omitempty"`
	VideoQualityMode              discord.Option[discord.VideoQualityMode]             `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    discord.Option[int]                                  `json:"default_auto_archive_duration,omitempty"`
	Flags                         discord.Option[discord.ChannelFlags]                 `json:"flags,omitempty"`
	AvailableTags                 discord.Option[[]discord.ChannelTag]                 `json:"available_tags,omitempty"`
	DefaultReactionEmoji          discord.Option[discord.DefaultReactionEmoji]         `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser discord.Option[int]                                  `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              discord.Option[discord.DefaultSortOrder]             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            discord.Option[discord.ChannelForumLayoutType]       `json:"default_forum_layout,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/channel#edit-channel-permissions
type EditChannelPermissionsParams struct {
	Allow discord.Option[string]          `json:"allow,omitempty"`
	Deny  discord.Option[string]          `json:"deny,omitempty"`
	Type  discord.PermissionOverwriteType `json:"type"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/channel#create-channel-invite
type CreateChannelInviteParams struct {
	MaxAge              discord.Option[int]               `json:"max_age,omitempty"`
	MaxUses             discord.Option[int]               `json:"max_uses,omitempty"`
	Temporary           discord.Option[bool]              `json:"temporary,omitempty"`
	Unique              discord.Option[bool]              `json:"unique,omitempty"`
	TargetType          discord.Option[int]               `json:"target_type,omitempty"`
	TargetUserID        discord.Option[discord.Snowflake] `json:"target_user_id,omitempty"`
	TargetApplicationID discord.Option[discord.Snowflake] `json:"target_application_id,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/channel#deleteclose-channel
type DeleteChannelOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/channel#delete-channel-permission
type DeleteChannelPermissionOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/channel#follow-announcement-channel
type FollowAnnouncementChannelOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/channel#pin-message
type PinMessageOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/channel#unpin-message
type UnpinMessageOptions struct {
	Reason string
}

// ── Channel endpoints ─────────────────────────────────────────────────────────

// ModifyChannel updates settings for a channel. Returns the updated channel.
// https://docs.discord.com/developers/resources/channel#modify-channel
func (c *RestClient) ModifyChannel(ctx context.Context, channelID discord.Snowflake, params ModifyChannelParams) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/channels/"+channelID.String(), bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteChannel deletes a channel by ID. For guild channels requires MANAGE_CHANNELS.
// Returns the deleted channel object.
// https://docs.discord.com/developers/resources/channel#deleteclose-channel
func (c *RestClient) DeleteChannel(ctx context.Context, channelID discord.Snowflake, opts *DeleteChannelOptions) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String(), nil, c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Channel](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String(), nil, c.WithBotAuthorization(), WithAuditLogReason(&opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListChannelInvites returns all invites for a channel. Requires MANAGE_CHANNELS.
// https://docs.discord.com/developers/resources/channel#get-channel-invites
func (c *RestClient) ListChannelInvites(ctx context.Context, channelID discord.Snowflake) ([]*discord.Invite, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/channels/"+channelID.String()+"/invites", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateChannelInvite creates a new invite for a channel. Requires CREATE_INSTANT_INVITE.
// https://docs.discord.com/developers/resources/channel#create-channel-invite
func (c *RestClient) CreateChannelInvite(ctx context.Context, channelID discord.Snowflake, params CreateChannelInviteParams) (*discord.Invite, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/invites", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditChannelPermissions creates or updates a permission overwrite for a user or role in a channel.
// Requires MANAGE_ROLES. Use params.Type to specify whether overwriteID is a role or member.
// https://docs.discord.com/developers/resources/channel#edit-channel-permissions
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

	req, err := c.generateRequest(ctx, http.MethodPut, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// DeleteChannelPermission removes a permission overwrite for a user or role from a channel.
// Requires MANAGE_ROLES.
// https://docs.discord.com/developers/resources/channel#delete-channel-permission
func (c *RestClient) DeleteChannelPermission(ctx context.Context, channelID, overwriteID discord.Snowflake, opts *DeleteChannelPermissionOptions) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := overwriteID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()

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

// TriggerTypingIndicator posts a typing indicator to the channel for ~10 seconds.
// https://docs.discord.com/developers/resources/channel#trigger-typing-indicator
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

// AddGroupDMRecipientParams holds params for adding a user to a Group DM.
// https://docs.discord.com/developers/resources/channel#group-dm-add-recipient
type AddGroupDMRecipientParams struct {
	AccessToken string `json:"access_token"`
	Nick        string `json:"nick"`
}

// SetVoiceChannelStatus sets the status message of a voice channel.
// https://docs.discord.com/developers/resources/channel#set-voice-channel-status
func (c *RestClient) SetVoiceChannelStatus(ctx context.Context, channelID discord.Snowflake, status *string) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	// status is nullable: a nil pointer must marshal to JSON null (which clears
	// the voice channel status), not "". Keep the pointer so encoding/json emits
	// null for nil and the string value otherwise.
	body, err := json.Marshal(map[string]*string{"status": status})
	if err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/channels/"+channelID.String()+"/voice-status", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// FollowAnnouncementChannel follows an announcement channel, publishing messages to webhookChannelID.
// https://docs.discord.com/developers/resources/channel#follow-announcement-channel
func (c *RestClient) FollowAnnouncementChannel(ctx context.Context, channelID, webhookChannelID discord.Snowflake, opts *FollowAnnouncementChannelOptions) (*discord.FollowedChannel, error) {
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

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/followers", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.FollowedChannel](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/followers", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(&opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.FollowedChannel](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// AddGroupDMRecipient adds a user to a Group DM using their OAuth2 access token.
// https://docs.discord.com/developers/resources/channel#group-dm-add-recipient
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
// https://docs.discord.com/developers/resources/channel#group-dm-remove-recipient
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
