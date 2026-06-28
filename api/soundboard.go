package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/soundboard#create-guild-soundboard-sound
type CreateGuildSoundboardSoundParams struct {
	Name      string                            `json:"name"`
	Sound     string                            `json:"sound"` // base64 data URI
	Volume    discord.Option[float64]           `json:"volume,omitempty"`
	EmojiID   discord.Option[discord.Snowflake] `json:"emoji_id,omitempty"`
	EmojiName discord.Option[string]            `json:"emoji_name,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/soundboard#modify-guild-soundboard-sound
type ModifyGuildSoundboardSoundParams struct {
	Name      discord.Option[string]            `json:"name,omitempty"`
	Volume    discord.Option[float64]           `json:"volume,omitempty"`
	EmojiID   discord.Option[discord.Snowflake] `json:"emoji_id,omitempty"`
	EmojiName discord.Option[string]            `json:"emoji_name,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/soundboard#delete-guild-soundboard-sound
type DeleteGuildSoundboardSoundOptions struct {
	Reason string
}

// https://docs.discord.com/developers/resources/soundboard#send-soundboard-sound
type SendSoundboardSoundParams struct {
	SoundID       discord.Snowflake                 `json:"sound_id"`
	SourceGuildID discord.Option[discord.Snowflake] `json:"source_guild_id,omitempty"`
}

// ── Soundboard endpoints ──────────────────────────────────────────────────────

// ListDefaultSoundboardSounds returns the list of default sounds available to all users.
// https://docs.discord.com/developers/resources/soundboard#list-default-soundboard-sounds
func (c *RestClient) ListDefaultSoundboardSounds(ctx context.Context) ([]*discord.SoundboardSound, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/soundboard-default-sounds", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.SoundboardSound](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// https://docs.discord.com/developers/resources/soundboard#list-guild-soundboard-sounds
type ListGuildSoundboardSoundsResponse struct {
	Items []discord.SoundboardSound `json:"items"`
}

// ListGuildSoundboardSounds returns all soundboard sounds for a guild.
// https://docs.discord.com/developers/resources/soundboard#list-guild-soundboard-sounds
func (c *RestClient) ListGuildSoundboardSounds(ctx context.Context, guildID discord.Snowflake) (*ListGuildSoundboardSoundsResponse, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/soundboard-sounds", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ListGuildSoundboardSoundsResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildSoundboardSound returns a single soundboard sound for a guild.
// https://docs.discord.com/developers/resources/soundboard#get-guild-soundboard-sound
func (c *RestClient) GetGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake) (*discord.SoundboardSound, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := soundID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/soundboard-sounds/" + soundID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.SoundboardSound](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildSoundboardSound creates a new soundboard sound in a guild. Requires MANAGE_GUILD_EXPRESSIONS.
// https://docs.discord.com/developers/resources/soundboard#create-guild-soundboard-sound
func (c *RestClient) CreateGuildSoundboardSound(ctx context.Context, guildID discord.Snowflake, params CreateGuildSoundboardSoundParams) (*discord.SoundboardSound, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/soundboard-sounds", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.SoundboardSound](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildSoundboardSound edits a soundboard sound. Requires MANAGE_GUILD_EXPRESSIONS.
// https://docs.discord.com/developers/resources/soundboard#modify-guild-soundboard-sound
func (c *RestClient) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, params ModifyGuildSoundboardSoundParams) (*discord.SoundboardSound, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := soundID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/soundboard-sounds/" + soundID.String()

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.SoundboardSound](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildSoundboardSound deletes a soundboard sound. Requires MANAGE_GUILD_EXPRESSIONS.
// https://docs.discord.com/developers/resources/soundboard#delete-guild-soundboard-sound
func (c *RestClient) DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, opts *DeleteGuildSoundboardSoundOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := soundID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/soundboard-sounds/" + soundID.String()

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

// SendSoundboardSound sends a soundboard sound to a voice channel the user is connected to.
// https://docs.discord.com/developers/resources/soundboard#send-soundboard-sound
func (c *RestClient) SendSoundboardSound(ctx context.Context, channelID discord.Snowflake, params SendSoundboardSoundParams) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/send-soundboard-sound"
	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}
