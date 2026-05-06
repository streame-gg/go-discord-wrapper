package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type CreateGuildSoundboardSoundParams struct {
	Name      string            `json:"name"`
	Sound     string            `json:"sound"` // base64 data URI
	Volume    *float64          `json:"volume,omitempty"`
	EmojiID   *common.Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string           `json:"emoji_name,omitempty"`
}

type ModifyGuildSoundboardSoundParams struct {
	Name      *string           `json:"name,omitempty"`
	Volume    *float64          `json:"volume,omitempty"`
	EmojiID   *common.Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string           `json:"emoji_name,omitempty"`
}

type SendSoundboardSoundParams struct {
	SoundID       common.Snowflake  `json:"sound_id"`
	SourceGuildID *common.Snowflake `json:"source_guild_id,omitempty"`
}

// ── Soundboard endpoints ──────────────────────────────────────────────────────

// ListDefaultSoundboardSounds returns the list of default sounds available to all users.
func (c *RestClient) ListDefaultSoundboardSounds(ctx context.Context) ([]*common.SoundboardSound, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/soundboard-default-sounds", nil)
	if err != nil {
		return nil, err
	}

	var sounds []*common.SoundboardSound
	if _, err := c.do(req, http.StatusOK, &sounds); err != nil {
		return nil, err
	}

	return sounds, nil
}

// ListGuildSoundboardSounds returns all soundboard sounds for a guild.
func (c *RestClient) ListGuildSoundboardSounds(ctx context.Context, guildID common.Snowflake) ([]*common.SoundboardSound, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/soundboard-sounds", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []*common.SoundboardSound `json:"items"`
	}
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// GetGuildSoundboardSound returns a single soundboard sound for a guild.
func (c *RestClient) GetGuildSoundboardSound(ctx context.Context, guildID, soundID common.Snowflake) (*common.SoundboardSound, error) {
	path := "/guilds/" + guildID.String() + "/soundboard-sounds/" + soundID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var sound common.SoundboardSound
	if _, err := c.do(req, http.StatusOK, &sound); err != nil {
		return nil, err
	}

	return &sound, nil
}

// CreateGuildSoundboardSound creates a new soundboard sound in a guild. Requires MANAGE_GUILD_EXPRESSIONS.
func (c *RestClient) CreateGuildSoundboardSound(ctx context.Context, guildID common.Snowflake, params CreateGuildSoundboardSoundParams) (*common.SoundboardSound, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/soundboard-sounds", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var sound common.SoundboardSound
	if _, err := c.do(req, http.StatusOK, &sound); err != nil {
		return nil, err
	}

	return &sound, nil
}

// ModifyGuildSoundboardSound edits a soundboard sound. Requires MANAGE_GUILD_EXPRESSIONS.
func (c *RestClient) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID common.Snowflake, params ModifyGuildSoundboardSoundParams) (*common.SoundboardSound, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/soundboard-sounds/" + soundID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var sound common.SoundboardSound
	if _, err := c.do(req, http.StatusOK, &sound); err != nil {
		return nil, err
	}

	return &sound, nil
}

// DeleteGuildSoundboardSound deletes a soundboard sound. Requires MANAGE_GUILD_EXPRESSIONS.
func (c *RestClient) DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/soundboard-sounds/" + soundID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// SendSoundboardSound sends a soundboard sound to a voice channel the user is connected to.
func (c *RestClient) SendSoundboardSound(ctx context.Context, channelID common.Snowflake, params SendSoundboardSoundParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/send-soundboard-sound"
	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
