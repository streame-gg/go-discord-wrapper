package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/sticker#modify-guild-sticker
type ModifyGuildStickerParams struct {
	Name        discord.Option[string] `json:"name,omitempty"`
	Description discord.Option[string] `json:"description,omitempty"`
	Tags        discord.Option[string] `json:"tags,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/sticker#create-guild-sticker
type CreateGuildStickerParams struct {
	Name        string
	Description string
	Tags        string
	File        []byte
	ContentType string

	AuditLogReason *string
}

// https://docs.discord.com/developers/resources/sticker#delete-guild-sticker
type DeleteGuildStickerOptions struct {
	Reason string
}

// ── Sticker endpoints ─────────────────────────────────────────────────────────

// GetSticker returns the sticker object for the given sticker ID.
// https://docs.discord.com/developers/resources/sticker#get-sticker
func (c *RestClient) GetSticker(ctx context.Context, stickerID discord.Snowflake) (*discord.Sticker, error) {
	if err := stickerID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/stickers/"+stickerID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// https://docs.discord.com/developers/resources/sticker#list-sticker-packs
type ListStickerPacksResponse struct {
	StickerPacks []*discord.StickerPack `json:"sticker_packs"`
}

// ListStickerPacks returns the list of sticker packs available to Nitro subscribers.
// https://docs.discord.com/developers/resources/sticker#list-sticker-packs
func (c *RestClient) ListStickerPacks(ctx context.Context) (*ListStickerPacksResponse, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/sticker-packs", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ListStickerPacksResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListGuildStickers returns all stickers for the given guild.
// https://docs.discord.com/developers/resources/sticker#list-guild-stickers
func (c *RestClient) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) ([]*discord.Sticker, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/stickers", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildSticker returns a specific sticker from a guild.
// https://docs.discord.com/developers/resources/sticker#get-guild-sticker
func (c *RestClient) GetGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake) (*discord.Sticker, error) {
	if err := stickerID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildSticker updates the name, description, or tags of a guild sticker.
// https://docs.discord.com/developers/resources/sticker#modify-guild-sticker
func (c *RestClient) ModifyGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, params ModifyGuildStickerParams) (*discord.Sticker, error) {
	if err := stickerID.Validate(); err != nil {
		return nil, err
	}

	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildSticker deletes a sticker from a guild.
// https://docs.discord.com/developers/resources/sticker#delete-guild-sticker
func (c *RestClient) DeleteGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, opts *DeleteGuildStickerOptions) error {
	if err := stickerID.Validate(); err != nil {
		return err
	}

	if err := guildID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()

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

// GetStickerPack returns the sticker pack object for the given pack ID.
// https://docs.discord.com/developers/resources/sticker#get-sticker-pack
func (c *RestClient) GetStickerPack(ctx context.Context, packID discord.Snowflake) (*discord.StickerPack, error) {
	if err := packID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/sticker-packs/"+packID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.StickerPack](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

var validStickerContentTypes = map[string]bool{
	"image/png":        true,
	"image/apng":       true,
	"image/gif":        true,
	"application/json": true, // Lottie
}

// CreateGuildSticker uploads a new sticker to a guild using multipart form encoding.
// https://docs.discord.com/developers/resources/sticker#create-guild-sticker
func (c *RestClient) CreateGuildSticker(ctx context.Context, guildID discord.Snowflake, params CreateGuildStickerParams) (*discord.Sticker, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if !validStickerContentTypes[params.ContentType] {
		return nil, fmt.Errorf("invalid sticker content type %q: must be image/png, image/apng, image/gif, or application/json", params.ContentType)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("name", params.Name)
	_ = w.WriteField("description", params.Description)
	_ = w.WriteField("tags", params.Tags)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="sticker"`)
	h.Set("Content-Type", params.ContentType)
	fw, err := w.CreatePart(h)
	if err != nil {
		return nil, err
	}
	if _, err := fw.Write(params.File); err != nil {
		return nil, err
	}
	w.Close()

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/stickers", &buf, c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}
