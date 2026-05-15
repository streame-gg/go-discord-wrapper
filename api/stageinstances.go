package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type CreateStageInstanceParams struct {
	ChannelID             common.Snowflake                  `json:"channel_id"`
	Topic                 string                            `json:"topic"`
	PrivacyLevel          *common.StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
	SendStartNotification *bool                             `json:"send_start_notification,omitempty"`
	GuildScheduledEventID *common.Snowflake                 `json:"guild_scheduled_event_id,omitempty"`
}

type ModifyStageInstanceParams struct {
	Topic        *string                           `json:"topic,omitempty"`
	PrivacyLevel *common.StageInstancePrivacyLevel `json:"privacy_level,omitempty"`
}

// ── Stage instance endpoints ──────────────────────────────────────────────────

// CreateStageInstance creates a new stage instance in a stage voice channel.
func (c *RestClient) CreateStageInstance(ctx context.Context, params CreateStageInstanceParams) (*common.StageInstance, error) {
	if err := params.ChannelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/stage-instances", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.StageInstance](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// GetStageInstance returns the stage instance for the given stage channel.
func (c *RestClient) GetStageInstance(ctx context.Context, channelID common.Snowflake) (*common.StageInstance, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/stage-instances/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.StageInstance](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyStageInstance updates fields on an existing stage instance.
func (c *RestClient) ModifyStageInstance(ctx context.Context, channelID common.Snowflake, params ModifyStageInstanceParams) (*common.StageInstance, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/stage-instances/"+channelID.String(), bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.StageInstance](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteStageInstance deletes the stage instance for the given stage channel.
func (c *RestClient) DeleteStageInstance(ctx context.Context, channelID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/stage-instances/"+channelID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}
