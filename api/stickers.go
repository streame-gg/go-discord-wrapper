package api

import (
	"bytes"
	"encoding/json"
	"net/http"

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

// ── Sticker endpoints ─────────────────────────────────────────────────────────

// GetSticker returns the sticker object for the given sticker ID.
func (c *RestClient) GetSticker(stickerID common.Snowflake) (*common.Sticker, error) {
	req, err := c.generateRequest(http.MethodGet, "/stickers/"+stickerID.String(), nil)
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
func (c *RestClient) ListStickerPacks() ([]*StickerPack, error) {
	req, err := c.generateRequest(http.MethodGet, "/sticker-packs", nil)
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
func (c *RestClient) ListGuildStickers(guildID common.Snowflake) ([]*common.Sticker, error) {
	req, err := c.generateRequest(http.MethodGet, "/guilds/"+guildID.String()+"/stickers", nil)
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
func (c *RestClient) GetGuildSticker(guildID, stickerID common.Snowflake) (*common.Sticker, error) {
	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(http.MethodGet, path, nil)
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
func (c *RestClient) ModifyGuildSticker(guildID, stickerID common.Snowflake, params ModifyGuildStickerParams) (*common.Sticker, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
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
func (c *RestClient) DeleteGuildSticker(guildID, stickerID common.Snowflake) error {
	path := "/guilds/" + guildID.String() + "/stickers/" + stickerID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
