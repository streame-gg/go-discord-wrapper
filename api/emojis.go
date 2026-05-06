package api

import (
	"bytes"
	"context"
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
func (c *RestClient) ListGuildEmojis(ctx context.Context, guildID common.Snowflake) ([]*common.Emoji, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/emojis", nil)
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
func (c *RestClient) GetGuildEmoji(ctx context.Context, guildID, emojiID common.Snowflake) (*common.Emoji, error) {
	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
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
func (c *RestClient) CreateGuildEmoji(ctx context.Context, guildID common.Snowflake, params CreateGuildEmojiParams) (*common.Emoji, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/emojis", bytes.NewReader(body))
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
func (c *RestClient) ModifyGuildEmoji(ctx context.Context, guildID, emojiID common.Snowflake, params ModifyGuildEmojiParams) (*common.Emoji, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
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
func (c *RestClient) DeleteGuildEmoji(ctx context.Context, guildID, emojiID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Application emoji endpoints ───────────────────────────────────────────────

// CreateEmojiParams holds params for creating an application emoji.
type CreateEmojiParams struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// ModifyEmojiParams holds params for modifying an application emoji.
type ModifyEmojiParams struct {
	Name *string `json:"name,omitempty"`
}

// ListApplicationEmojis returns all emojis for an application.
func (c *RestClient) ListApplicationEmojis(ctx context.Context, appID common.Snowflake) ([]*common.Emoji, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/"+appID.String()+"/emojis", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []*common.Emoji `json:"items"`
	}
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// GetApplicationEmoji returns a specific application emoji.
func (c *RestClient) GetApplicationEmoji(ctx context.Context, appID, emojiID common.Snowflake) (*common.Emoji, error) {
	path := "/applications/" + appID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var emoji common.Emoji
	if _, err := c.do(req, http.StatusOK, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// CreateApplicationEmoji creates a new emoji for an application.
func (c *RestClient) CreateApplicationEmoji(ctx context.Context, appID common.Snowflake, params CreateEmojiParams) (*common.Emoji, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/applications/"+appID.String()+"/emojis", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var emoji common.Emoji
	if _, err := c.do(req, http.StatusOK, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// ModifyApplicationEmoji updates an application emoji.
func (c *RestClient) ModifyApplicationEmoji(ctx context.Context, appID, emojiID common.Snowflake, params ModifyEmojiParams) (*common.Emoji, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var emoji common.Emoji
	if _, err := c.do(req, http.StatusOK, &emoji); err != nil {
		return nil, err
	}

	return &emoji, nil
}

// DeleteApplicationEmoji deletes an application emoji.
func (c *RestClient) DeleteApplicationEmoji(ctx context.Context, appID, emojiID common.Snowflake) error {
	path := "/applications/" + appID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
