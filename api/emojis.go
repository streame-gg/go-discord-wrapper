package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type CreateGuildEmojiParams struct {
	Name string `json:"name"`
	// Image is a base64-encoded image data URI (e.g. "data:image/png;base64,...").
	Image string             `json:"image"`
	Roles []common.Snowflake `json:"roles,omitempty"`
}

type ModifyGuildEmojiParams struct {
	Name  *string            `json:"name,omitempty"`
	Roles []common.Snowflake `json:"roles,omitempty"`
}

// ── Emoji endpoints ───────────────────────────────────────────────────────────

// ListGuildEmojis returns all emojis for a guild.
func (c *RestClient) ListGuildEmojis(guildID common.Snowflake) ([]*common.Emoji, error) {
	req, err := c.generateRequest(http.MethodGet, "/guilds/"+guildID.String()+"/emojis", nil)
	if err != nil {
		return nil, err
	}

	var emojis []*common.Emoji
	if _, err := c.do(req, http.StatusOK, &emojis); err != nil {
		return nil, err
	}

	return emojis, nil
}

// GetGuildEmoji returns a specific emoji from a guild.
func (c *RestClient) GetGuildEmoji(guildID, emojiID common.Snowflake) (*common.Emoji, error) {
	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var emoji common.Emoji
	if _, err := c.do(req, http.StatusOK, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// CreateGuildEmoji creates a new emoji in a guild.
// Image must be a base64-encoded data URI, max 256 KB.
func (c *RestClient) CreateGuildEmoji(guildID common.Snowflake, params CreateGuildEmojiParams) (*common.Emoji, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPost, "/guilds/"+guildID.String()+"/emojis", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var emoji common.Emoji
	if _, err := c.do(req, http.StatusOK, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// ModifyGuildEmoji updates the name or allowed roles for a guild emoji.
func (c *RestClient) ModifyGuildEmoji(guildID, emojiID common.Snowflake, params ModifyGuildEmojiParams) (*common.Emoji, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var emoji common.Emoji
	if _, err := c.do(req, http.StatusOK, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// DeleteGuildEmoji deletes the given guild emoji.
func (c *RestClient) DeleteGuildEmoji(guildID, emojiID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
