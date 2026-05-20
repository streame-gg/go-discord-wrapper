package discord

import (
	"context"
	"errors"
)

var errStickerNoGuildID = errors.New("discord: Sticker.GuildID is not set — ensure the sticker was received via a gateway event or cache")

// ── Options ───────────────────────────────────────────────────────────────────

// StickerCreateOptions configures Guild.CreateSticker.
type StickerCreateOptions struct {
	Name        string
	Description string
	Tags        string
	// File is the raw image bytes. ContentType is its MIME type (e.g. "image/png").
	File           []byte
	ContentType    string
	AuditLogReason *string
}

// StickerEditOptions configures Sticker.Edit.
type StickerEditOptions struct {
	Name           *string
	Description    *string
	Tags           *string
	AuditLogReason *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the sticker.
func (s *Sticker) Hydrate(c EntityClient) { s.hClient = c }

// WithClient returns a shallow copy of the Sticker with the client replaced.
func (s *Sticker) WithClient(c EntityClient) *Sticker {
	cp := *s
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the sticker has an associated client.
func (s *Sticker) IsHydrated() bool { return s.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this sticker. Requires MANAGE_GUILD_EXPRESSIONS.
func (s *Sticker) Edit(ctx context.Context, opts StickerEditOptions) (*Sticker, error) {
	c, err := ensureClient(s.hClient)
	if err != nil {
		return nil, err
	}
	if s.GuildID == nil || *s.GuildID == "" {
		return nil, errStickerNoGuildID
	}
	return c.ModifyGuildSticker(ctx, *s.GuildID, s.ID, opts)
}

// Delete deletes this sticker. Requires MANAGE_GUILD_EXPRESSIONS.
func (s *Sticker) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(s.hClient)
	if err != nil {
		return err
	}
	if s.GuildID == nil || *s.GuildID == "" {
		return errStickerNoGuildID
	}
	return c.DeleteGuildSticker(ctx, *s.GuildID, s.ID, reason)
}
