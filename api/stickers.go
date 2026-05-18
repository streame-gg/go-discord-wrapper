package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

type ModifyGuildStickerParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
}

type ModifyGuildStickerOptions struct {
	Reason string
}

type CreateGuildStickerParams struct {
	Name        string
	Description string
	Tags        string
	File        []byte
	ContentType string
}

type CreateGuildStickerOptions struct {
	Reason string
}

type DeleteGuildStickerOptions struct {
	Reason string
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
	StickerPacks []*discord.StickerPack `json:"sticker_packs"`
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

// ListGuildStickers returns all stickers for the given guild, keyed by sticker ID.
func (c *RestClient) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) (*collection.Collection[discord.Snowflake, *discord.Sticker], error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/stickers", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	stickers, err := doRequestSlice[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
	if err != nil {
		return nil, err
	}

	coll := collection.NewWithCapacity[discord.Snowflake, *discord.Sticker](len(stickers))
	for _, s := range stickers {
		coll.Set(s.ID, s)
	}
	return coll, nil
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
func (c *RestClient) ModifyGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, params ModifyGuildStickerParams, opts *ModifyGuildStickerOptions) (*discord.Sticker, error) {
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

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Sticker](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildSticker deletes a sticker from a guild.
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

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// GetStickerPack returns the sticker pack object for the given pack ID.
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
func (c *RestClient) CreateGuildSticker(ctx context.Context, guildID discord.Snowflake, params CreateGuildStickerParams, opts *CreateGuildStickerOptions) (*discord.Sticker, error) {
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

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/stickers", &buf, c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		return doRequest[discord.Sticker](c, req, map[int]bool{
			http.StatusCreated: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/stickers", &buf, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return doRequest[discord.Sticker](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}
