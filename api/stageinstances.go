package api

import (
	"bytes"
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
func (c *RestClient) CreateStageInstance(params CreateStageInstanceParams) (*common.StageInstance, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPost, "/stage-instances", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var instance common.StageInstance
	if _, err := c.do(req, http.StatusCreated, &instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

// GetStageInstance returns the stage instance for the given stage channel.
func (c *RestClient) GetStageInstance(channelID common.Snowflake) (*common.StageInstance, error) {
	req, err := c.generateRequest(http.MethodGet, "/stage-instances/"+channelID.String(), nil)
	if err != nil {
		return nil, err
	}

	var instance common.StageInstance
	if _, err := c.do(req, http.StatusOK, &instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

// ModifyStageInstance updates fields on an existing stage instance.
func (c *RestClient) ModifyStageInstance(channelID common.Snowflake, params ModifyStageInstanceParams) (*common.StageInstance, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPatch, "/stage-instances/"+channelID.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var instance common.StageInstance
	if _, err := c.do(req, http.StatusOK, &instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

// DeleteStageInstance deletes the stage instance for the given stage channel.
func (c *RestClient) DeleteStageInstance(channelID common.Snowflake) error {
	req, err := c.generateRequest(http.MethodDelete, "/stage-instances/"+channelID.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
