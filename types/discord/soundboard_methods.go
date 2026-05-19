package discord

import (
	"context"
	"errors"
)

var errSoundNoGuildID = errors.New("discord: SoundboardSound.GuildID is not set — ensure the sound was received via a gateway event or cache")

// ── Options ───────────────────────────────────────────────────────────────────

// SoundEditOptions configures SoundboardSound.Edit.
type SoundEditOptions struct {
	Name           *string
	Volume         *float64
	EmojiID        *Snowflake
	EmojiName      *string
	AuditLogReason *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the soundboard sound.
func (s *SoundboardSound) Hydrate(c EntityClient) {
	s.hClient = c
	if s.User != nil {
		s.User.Hydrate(c)
	}
}

// WithClient returns a shallow copy of the SoundboardSound with the client replaced.
func (s *SoundboardSound) WithClient(c EntityClient) *SoundboardSound {
	cp := *s
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the sound has an associated client.
func (s *SoundboardSound) IsHydrated() bool { return s.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this soundboard sound. Requires MANAGE_GUILD_EXPRESSIONS.
func (s *SoundboardSound) Edit(ctx context.Context, opts SoundEditOptions) (*SoundboardSound, error) {
	c, err := ensureClient(s.hClient)
	if err != nil {
		return nil, err
	}
	if s.GuildID == nil || *s.GuildID == "" {
		return nil, errSoundNoGuildID
	}
	return c.ModifyGuildSoundboardSound(ctx, *s.GuildID, s.SoundID, opts)
}

// Delete deletes this soundboard sound. Requires MANAGE_GUILD_EXPRESSIONS.
func (s *SoundboardSound) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(s.hClient)
	if err != nil {
		return err
	}
	if s.GuildID == nil || *s.GuildID == "" {
		return errSoundNoGuildID
	}
	return c.DeleteGuildSoundboardSound(ctx, *s.GuildID, s.SoundID, reason)
}
