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

// ── Response types ────────────────────────────────────────────────────────────

// StickerPack is a pack of standard stickers sold in the Discord store.
type StickerPack struct {
	ID             discord.Snowflake  `json:"id"`
	Stickers       []discord.Sticker  `json:"stickers"`
	Name           string             `json:"name"`
	SKUId          discord.Snowflake  `json:"sku_id"`
	CoverStickerID *discord.Snowflake `json:"cover_sticker_id,omitempty"`
	Description    string             `json:"description"`
	BannerAssetID  *discord.Snowflake `json:"banner_asset_id,omitempty"`
}

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyGuildStickerParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
}

type CreateGuildStickerParams struct {
	Name        string
	Description string
	Tags        string
	File        []byte
	ContentType string
}

// ── Sticker endpoints ─────────────────────────────────────────────────────────

// GetSticker returns the sticker object for the given sticker ID.
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

type ListStickerPacksResponse struct {
	StickerPacks []*StickerPack `json:"sticker_packs"`
}

// ListStickerPacks returns the list of sticker packs available to Nitro subscribers.
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
func (c *RestClient) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) (*[]*discord.Sticker, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/stickers", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildSticker returns a specific sticker from a guild.
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
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildSticker deletes a sticker from a guild.
func (c *RestClient) DeleteGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake) error {
	if err := stickerID.Validate(); err != nil {
		return err
	}

	if err := guildID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// GetStickerPack returns the sticker pack object for the given pack ID.
func (c *RestClient) GetStickerPack(ctx context.Context, packID discord.Snowflake) (*StickerPack, error) {
	if err := packID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/sticker-packs/"+packID.String(), nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[StickerPack](c, req, map[int]bool{
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

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/stickers", &buf, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}
