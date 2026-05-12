package api

import (
	"bytes"
	"context"
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
func (c *RestClient) ListVoiceRegions(ctx context.Context) ([]*common.VoiceRegion, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/voice/regions", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var regions []*common.VoiceRegion
	if _, err := c.do(req, []int{http.StatusOK}, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

// ListGuildVoiceRegions returns voice regions available for a guild, including VIP regions if applicable.
func (c *RestClient) ListGuildVoiceRegions(ctx context.Context, guildID common.Snowflake) ([]*common.VoiceRegion, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/regions", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var regions []*common.VoiceRegion
	if _, err := c.do(req, []int{http.StatusOK}, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

// ModifyCurrentUserVoiceState updates the bot's voice state in a guild stage channel.
func (c *RestClient) ModifyCurrentUserVoiceState(ctx context.Context, guildID common.Snowflake, params ModifyCurrentUserVoiceStateParams) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/guilds/"+guildID.String()+"/voice-states/@me", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	_, err = c.do(req, []int{http.StatusNoContent}, nil)
	return err
}

// ModifyUserVoiceState updates another user's voice state in a guild stage channel. Requires MUTE_MEMBERS.
func (c *RestClient) ModifyUserVoiceState(ctx context.Context, guildID, userID common.Snowflake, params ModifyUserVoiceStateParams) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/voice-states/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	_, err = c.do(req, []int{http.StatusNoContent}, nil)
	return err
}

// GetCurrentUserVoiceState returns the current user's voice state in a guild.
func (c *RestClient) GetCurrentUserVoiceState(ctx context.Context, guildID common.Snowflake) (*common.VoiceState, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/voice-states/@me", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var state common.VoiceState
	if _, err := c.do(req, []int{http.StatusOK}, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// GetUserVoiceState returns a specific user's voice state in a guild.
func (c *RestClient) GetUserVoiceState(ctx context.Context, guildID, userID common.Snowflake) (*common.VoiceState, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := userID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/voice-states/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var state common.VoiceState
	if _, err := c.do(req, []int{http.StatusOK}, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
