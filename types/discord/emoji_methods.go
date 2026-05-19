package discord

import (
	"context"
	"errors"
)

var errEmojiNoGuildID = errors.New("discord: Emoji.GuildID is not set — ensure the emoji was received via a gateway event or cache")

// ── Options ───────────────────────────────────────────────────────────────────

// EmojiCreateOptions configures Guild.CreateEmoji.
type EmojiCreateOptions struct {
	Name string
	// Image is a base64-encoded image data URI (e.g. "data:image/png;base64,...").
	Image          string
	Roles          []Snowflake
	AuditLogReason *string
}

// EmojiEditOptions configures Emoji.Edit.
type EmojiEditOptions struct {
	Name           *string
	Roles          []Snowflake
	AuditLogReason *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the emoji.
func (e *Emoji) Hydrate(c EntityClient) { e.hClient = c }

// WithClient returns a shallow copy of the Emoji with the client replaced.
func (e *Emoji) WithClient(c EntityClient) *Emoji {
	cp := *e
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the emoji has an associated client.
func (e *Emoji) IsHydrated() bool { return e.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this emoji's name or role restrictions.
// Requires MANAGE_GUILD_EXPRESSIONS.
func (e *Emoji) Edit(ctx context.Context, opts EmojiEditOptions) (*Emoji, error) {
	c, err := ensureClient(e.hClient)
	if err != nil {
		return nil, err
	}
	if e.GuildID == "" {
		return nil, errEmojiNoGuildID
	}
	return c.ModifyGuildEmoji(ctx, e.GuildID, e.ID, opts)
}

// Delete deletes this emoji. Requires MANAGE_GUILD_EXPRESSIONS.
func (e *Emoji) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(e.hClient)
	if err != nil {
		return err
	}
	if e.GuildID == "" {
		return errEmojiNoGuildID
	}
	return c.DeleteGuildEmoji(ctx, e.GuildID, e.ID, reason)
}
