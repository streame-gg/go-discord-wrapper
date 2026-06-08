package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/emoji#create-guild-emoji
type CreateGuildEmojiParams struct {
	Name string `json:"name"`
	// Image is a base64-encoded image data URI (e.g. "data:image/png;base64,...").
	Image string              `json:"image"`
	Roles []discord.Snowflake `json:"roles,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/emoji#modify-guild-emoji
type ModifyGuildEmojiParams struct {
	Name  *string             `json:"name,omitempty"`
	Roles []discord.Snowflake `json:"roles,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/emoji#delete-guild-emoji
type DeleteGuildEmojiOptions struct {
	Reason string
}

// ── Emoji endpoints ───────────────────────────────────────────────────────────

// ListGuildEmojis returns all emojis for a guild.
// https://docs.discord.com/developers/resources/emoji#list-guild-emojis
func (c *RestClient) ListGuildEmojis(ctx context.Context, guildID discord.Snowflake) ([]*discord.Emoji, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/emojis", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Emoji](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildEmoji returns a specific emoji from a guild.
// https://docs.discord.com/developers/resources/emoji#get-guild-emoji
func (c *RestClient) GetGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake) (*discord.Emoji, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := emojiID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Emoji](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildEmoji creates a new emoji in a guild.
// Image must be a base64-encoded data URI, max 256 KB.
// https://docs.discord.com/developers/resources/emoji#create-guild-emoji
func (c *RestClient) CreateGuildEmoji(ctx context.Context, guildID discord.Snowflake, params CreateGuildEmojiParams) (*discord.Emoji, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/emojis", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Emoji](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// ModifyGuildEmoji updates the name or allowed roles for a guild emoji.
// https://docs.discord.com/developers/resources/emoji#modify-guild-emoji
func (c *RestClient) ModifyGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, params ModifyGuildEmojiParams) (*discord.Emoji, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := emojiID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Emoji](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildEmoji deletes the given guild emoji.
// https://docs.discord.com/developers/resources/emoji#delete-guild-emoji
func (c *RestClient) DeleteGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, opts *DeleteGuildEmojiOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := emojiID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/emojis/" + emojiID.String()

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

// ── Application emoji endpoints ───────────────────────────────────────────────

// CreateEmojiParams holds params for creating an application emoji.
// https://docs.discord.com/developers/resources/emoji#create-application-emoji
type CreateEmojiParams struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// ModifyEmojiParams holds params for modifying an application emoji.
// https://docs.discord.com/developers/resources/emoji#modify-application-emoji
type ModifyEmojiParams struct {
	Name *string `json:"name,omitempty"`
}

// https://docs.discord.com/developers/resources/emoji#list-application-emojis
type ListApplicationEmojisResponse struct {
	Items []*discord.Emoji `json:"items"`
}

// ListApplicationEmojis returns all emojis for an application.
// https://docs.discord.com/developers/resources/emoji#list-application-emojis
func (c *RestClient) ListApplicationEmojis(ctx context.Context, appID discord.Snowflake) (*ListApplicationEmojisResponse, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/"+appID.String()+"/emojis", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ListApplicationEmojisResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetApplicationEmoji returns a specific application emoji.
// https://docs.discord.com/developers/resources/emoji#get-application-emoji
func (c *RestClient) GetApplicationEmoji(ctx context.Context, appID, emojiID discord.Snowflake) (*discord.Emoji, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := emojiID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Emoji](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateApplicationEmoji creates a new emoji for an application.
// https://docs.discord.com/developers/resources/emoji#create-application-emoji
func (c *RestClient) CreateApplicationEmoji(ctx context.Context, appID discord.Snowflake, params CreateEmojiParams) (*discord.Emoji, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/applications/"+appID.String()+"/emojis", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Emoji](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// ModifyApplicationEmoji updates an application emoji.
// https://docs.discord.com/developers/resources/emoji#modify-application-emoji
func (c *RestClient) ModifyApplicationEmoji(ctx context.Context, appID, emojiID discord.Snowflake, params ModifyEmojiParams) (*discord.Emoji, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := emojiID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Emoji](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteApplicationEmoji deletes an application emoji.
// https://docs.discord.com/developers/resources/emoji#delete-application-emoji
func (c *RestClient) DeleteApplicationEmoji(ctx context.Context, appID, emojiID discord.Snowflake) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	if err := emojiID.Validate(); err != nil {
		return err
	}

	path := "/applications/" + appID.String() + "/emojis/" + emojiID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}
