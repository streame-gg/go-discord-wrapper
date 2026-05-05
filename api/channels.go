package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// GetChannel fetches a channel by ID.
func (c *RestClient) GetChannel(id common.Snowflake) (*common.Channel, error) {
	req, err := c.generateRequest("GET", "/channels/"+id.String(), nil)
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
func (c *RestClient) ModifyChannel(channelID common.Snowflake, params ModifyChannelParams) (*common.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPatch, "/channels/"+channelID.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *RestClient) DeleteChannel(channelID common.Snowflake) (*common.Channel, error) {
	req, err := c.generateRequest(http.MethodDelete, "/channels/"+channelID.String(), nil)
	if err != nil {
		return nil, err
	}

	var channel common.Channel
	if _, err := c.do(req, http.StatusOK, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *RestClient) GetChannelInvites(channelID common.Snowflake) ([]*Invite, error) {
	req, err := c.generateRequest(http.MethodGet, "/channels/"+channelID.String()+"/invites", nil)
	if err != nil {
		return nil, err
	}

	var invites []*Invite
	if _, err := c.do(req, http.StatusOK, &invites); err != nil {
		return nil, err
	}

	return invites, nil
}

func (c *RestClient) CreateChannelInvite(channelID common.Snowflake, params CreateChannelInviteParams) (*Invite, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPost, "/channels/"+channelID.String()+"/invites", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var invite Invite
	if _, err := c.do(req, http.StatusOK, &invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

func (c *RestClient) EditChannelPermissions(channelID, overwriteID common.Snowflake, params EditChannelPermissionsParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()
	req, err := c.generateRequest(http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

func (c *RestClient) DeleteChannelPermission(channelID, overwriteID common.Snowflake) error {
	path := "/channels/" + channelID.String() + "/permissions/" + overwriteID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

func (c *RestClient) TriggerTypingIndicator(channelID common.Snowflake) error {
	req, err := c.generateRequest(http.MethodPost, "/channels/"+channelID.String()+"/typing", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
