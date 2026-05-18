package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyCurrentUserVoiceStateParams struct {
	ChannelID               *discord.Snowflake `json:"channel_id,omitempty"`
	Suppress                *bool              `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *time.Time         `json:"request_to_speak_timestamp,omitempty"`
}

type ModifyUserVoiceStateParams struct {
	ChannelID *discord.Snowflake `json:"channel_id,omitempty"`
	Suppress  *bool              `json:"suppress,omitempty"`
}

// ── Voice endpoints ───────────────────────────────────────────────────────────

// ListVoiceRegions returns all available voice regions, keyed by region ID.
func (c *RestClient) ListVoiceRegions(ctx context.Context) (*collection.Collection[string, *discord.VoiceRegion], error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/voice/regions", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	regions, err := doRequestSlice[discord.VoiceRegion](c, req, map[int]bool{
		http.StatusOK: true,
	})
	if err != nil {
		return nil, err
	}

	coll := collection.NewWithCapacity[string, *discord.VoiceRegion](len(regions))
	for _, r := range regions {
		coll.Set(r.ID, r)
	}
	return coll, nil
}

// ListGuildVoiceRegions returns voice regions available for a guild, keyed by region ID, including VIP regions if applicable.
func (c *RestClient) ListGuildVoiceRegions(ctx context.Context, guildID discord.Snowflake) (*collection.Collection[string, *discord.VoiceRegion], error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/regions", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	regions, err := doRequestSlice[discord.VoiceRegion](c, req, map[int]bool{
		http.StatusOK: true,
	})
	if err != nil {
		return nil, err
	}

	coll := collection.NewWithCapacity[string, *discord.VoiceRegion](len(regions))
	for _, r := range regions {
		coll.Set(r.ID, r)
	}
	return coll, nil
}

// ModifyCurrentUserVoiceState updates the bot's voice state in a guild stage channel.
func (c *RestClient) ModifyCurrentUserVoiceState(ctx context.Context, guildID discord.Snowflake, params ModifyCurrentUserVoiceStateParams) error {
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

	return doRequestWithoutResponse(c, req)
}

// ModifyUserVoiceState updates another user's voice state in a guild stage channel. Requires MUTE_MEMBERS.
func (c *RestClient) ModifyUserVoiceState(ctx context.Context, guildID, userID discord.Snowflake, params ModifyUserVoiceStateParams) error {
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

	return doRequestWithoutResponse(c, req)
}

// GetCurrentUserVoiceState returns the current user's voice state in a guild.
func (c *RestClient) GetCurrentUserVoiceState(ctx context.Context, guildID discord.Snowflake) (*discord.VoiceState, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/voice-states/@me", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.VoiceState](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetUserVoiceState returns a specific user's voice state in a guild.
func (c *RestClient) GetUserVoiceState(ctx context.Context, guildID, userID discord.Snowflake) (*discord.VoiceState, error) {
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

	return doRequest[discord.VoiceState](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
