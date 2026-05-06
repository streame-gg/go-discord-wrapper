package api

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Response types ────────────────────────────────────────────────────────────

// StickerPack is a pack of standard stickers sold in the Discord store.
type StickerPack struct {
	ID             common.Snowflake  `json:"id"`
	Stickers       []common.Sticker  `json:"stickers"`
	Name           string            `json:"name"`
	SKUId          common.Snowflake  `json:"sku_id"`
	CoverStickerID *common.Snowflake `json:"cover_sticker_id,omitempty"`
	Description    string            `json:"description"`
	BannerAssetID  *common.Snowflake `json:"banner_asset_id,omitempty"`
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
func (c *RestClient) GetSticker(ctx context.Context, stickerID common.Snowflake) (*common.Sticker, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/stickers/"+stickerID.String(), nil)
	if err != nil {
		return nil, err
	}

	var sticker common.Sticker
	if _, err := c.do(req, http.StatusOK, &sticker); err != nil {
		return nil, err
	}

	return &sticker, nil
}

// ListStickerPacks returns the list of sticker packs available to Nitro subscribers.
func (c *RestClient) ListStickerPacks(ctx context.Context) ([]*StickerPack, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/sticker-packs", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		StickerPacks []*StickerPack `json:"sticker_packs"`
	}
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.StickerPacks, nil
}

// ListGuildStickers returns all stickers for the given guild.
func (c *RestClient) ListGuildStickers(ctx context.Context, guildID common.Snowflake) ([]*common.Sticker, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/stickers", nil)
	if err != nil {
		return nil, err
	}

	var stickers []*common.Sticker
	if _, err := c.do(req, http.StatusOK, &stickers); err != nil {
		return nil, err
	}

	return stickers, nil
}

// GetGuildSticker returns a specific sticker from a guild.
func (c *RestClient) GetGuildSticker(ctx context.Context, guildID, stickerID common.Snowflake) (*common.Sticker, error) {
	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var sticker common.Sticker
	if _, err := c.do(req, http.StatusOK, &sticker); err != nil {
		return nil, err
	}

	return &sticker, nil
}

// ModifyGuildSticker updates the name, description, or tags of a guild sticker.
func (c *RestClient) ModifyGuildSticker(ctx context.Context, guildID, stickerID common.Snowflake, params ModifyGuildStickerParams) (*common.Sticker, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var sticker common.Sticker
	if _, err := c.do(req, http.StatusOK, &sticker); err != nil {
		return nil, err
	}

	return &sticker, nil
}

// DeleteGuildSticker deletes a sticker from a guild.
func (c *RestClient) DeleteGuildSticker(ctx context.Context, guildID, stickerID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// GetStickerPack returns the sticker pack object for the given pack ID.
func (c *RestClient) GetStickerPack(ctx context.Context, packID common.Snowflake) (*StickerPack, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/sticker-packs/"+packID.String(), nil)
	if err != nil {
		return nil, err
	}

	var pack StickerPack
	if _, err := c.do(req, http.StatusOK, &pack); err != nil {
		return nil, err
	}

	return &pack, nil
}

// CreateGuildSticker uploads a new sticker to a guild using multipart form encoding.
func (c *RestClient) CreateGuildSticker(ctx context.Context, guildID common.Snowflake, params CreateGuildStickerParams) (*common.Sticker, error) {
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

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/stickers", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	var sticker common.Sticker
	if _, err := c.do(req, http.StatusCreated, &sticker); err != nil {
		return nil, err
	}

	return &sticker, nil
}
