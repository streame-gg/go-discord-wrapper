package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/stage-instance#create-stage-instance
type CreateStageInstanceParams struct {
	ChannelID             discord.Snowflake                  `json:"channel_id"`
	Topic                 string                             `json:"topic"`
	PrivacyLevel          *discord.StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	SendStartNotification *bool                              `json:"send_start_notification,omitempty"`
	GuildScheduledEventID *discord.Snowflake                 `json:"guild_scheduled_event_id,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/stage-instance#modify-stage-instance
type ModifyStageInstanceParams struct {
	Topic        *string                            `json:"topic,omitempty"`
	PrivacyLevel *discord.StageInstancePrivacyLevel `json:"privacy_level,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/stage-instance#delete-stage-instance
type DeleteStageInstanceOptions struct {
	Reason string
}

// ── Stage instance endpoints ──────────────────────────────────────────────────

// CreateStageInstance creates a new stage instance in a stage voice channel.
// https://docs.discord.com/developers/resources/stage-instance#create-stage-instance
func (c *RestClient) CreateStageInstance(ctx context.Context, params CreateStageInstanceParams) (*discord.StageInstance, error) {
	if err := params.ChannelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/stage-instances", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.StageInstance](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// GetStageInstance returns the stage instance for the given stage channel.
// https://docs.discord.com/developers/resources/stage-instance#get-stage-instance
func (c *RestClient) GetStageInstance(ctx context.Context, channelID discord.Snowflake) (*discord.StageInstance, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/stage-instances/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.StageInstance](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyStageInstance updates fields on an existing stage instance.
// https://docs.discord.com/developers/resources/stage-instance#modify-stage-instance
func (c *RestClient) ModifyStageInstance(ctx context.Context, channelID discord.Snowflake, params ModifyStageInstanceParams) (*discord.StageInstance, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/stage-instances/"+channelID.String(), bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.StageInstance](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteStageInstance deletes the stage instance for the given stage channel.
// https://docs.discord.com/developers/resources/stage-instance#delete-stage-instance
func (c *RestClient) DeleteStageInstance(ctx context.Context, channelID discord.Snowflake, opts *DeleteStageInstanceOptions) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, "/stage-instances/"+channelID.String(), nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/stage-instances/"+channelID.String(), nil, c.WithBotAuthorization(), WithAuditLogReason(&opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}
