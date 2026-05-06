package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyCurrentUserVoiceStateParams struct {
	ChannelID               *common.Snowflake `json:"channel_id,omitempty"`
	Suppress                *bool             `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *time.Time        `json:"request_to_speak_timestamp,omitempty"`
}

type ModifyUserVoiceStateParams struct {
	ChannelID *common.Snowflake `json:"channel_id,omitempty"`
	Suppress  *bool             `json:"suppress,omitempty"`
}

// ── Voice endpoints ───────────────────────────────────────────────────────────

// ListVoiceRegions returns all available voice regions.
func (c *RestClient) ListVoiceRegions() ([]*common.VoiceRegion, error) {
	req, err := c.generateRequest(http.MethodGet, "/voice/regions", nil)
	if err != nil {
		return nil, err
	}

	var regions []*common.VoiceRegion
	if _, err := c.do(req, http.StatusOK, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

// ListGuildVoiceRegions returns voice regions available for a guild, including VIP regions if applicable.
func (c *RestClient) ListGuildVoiceRegions(guildID common.Snowflake) ([]*common.VoiceRegion, error) {
	req, err := c.generateRequest(http.MethodGet, "/guilds/"+guildID.String()+"/regions", nil)
	if err != nil {
		return nil, err
	}

	var regions []*common.VoiceRegion
	if _, err := c.do(req, http.StatusOK, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

// ModifyCurrentUserVoiceState updates the bot's voice state in a guild stage channel.
func (c *RestClient) ModifyCurrentUserVoiceState(guildID common.Snowflake, params ModifyCurrentUserVoiceStateParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := c.generateRequest(http.MethodPatch, "/guilds/"+guildID.String()+"/voice-states/@me", bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ModifyUserVoiceState updates another user's voice state in a guild stage channel. Requires MUTE_MEMBERS.
func (c *RestClient) ModifyUserVoiceState(guildID, userID common.Snowflake, params ModifyUserVoiceStateParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/voice-states/" + userID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
