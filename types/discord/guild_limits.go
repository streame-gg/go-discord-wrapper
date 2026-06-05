package discord

import "strings"

// Tier-derived guild capacity helpers. Discord raises several per-guild limits
// with the boost (premium) tier, and a few features (e.g. MORE_STICKERS,
// VIP_REGIONS) override the tier baseline. These mirror discord.js' maximumX
// getters so callers can validate uploads or surface "X / max" counters without
// hard-coding the tables. See https://discord.com/developers/docs/resources/guild.

// MaxEmojis returns the maximum number of custom emoji this guild may have for
// each of the static and animated buckets. Guilds with the MORE_EMOJI feature
// get a flat raised cap regardless of tier.
func (g *Guild) MaxEmojis() int {
	if g.HasFeature(GuildFeatures("MORE_EMOJI")) {
		return 200
	}
	switch g.PremiumTier {
	case GuildPremiumTierTier1:
		return 100
	case GuildPremiumTierTier2:
		return 150
	case GuildPremiumTierTier3:
		return 250
	default:
		return 50
	}
}

// MaxStickers returns the maximum number of custom stickers this guild may have.
// The MORE_STICKERS feature raises tier-0 guilds to the tier-1 cap.
func (g *Guild) MaxStickers() int {
	switch g.PremiumTier {
	case GuildPremiumTierTier1:
		return 15
	case GuildPremiumTierTier2:
		return 30
	case GuildPremiumTierTier3:
		return 60
	default:
		if g.HasFeature(GuildFeatureMoreStickers) {
			return 15
		}
		return 5
	}
}

// MaxBitrate returns the highest voice channel bitrate (in bits per second)
// allowed in this guild. The VIP_REGIONS feature unlocks the maximum bitrate
// regardless of tier.
func (g *Guild) MaxBitrate() int {
	if g.HasFeature(GuildFeatureVipRegions) {
		return 384_000
	}
	switch g.PremiumTier {
	case GuildPremiumTierTier1:
		return 128_000
	case GuildPremiumTierTier2:
		return 256_000
	case GuildPremiumTierTier3:
		return 384_000
	default:
		return 96_000
	}
}

// MaxSoundboardSounds returns the maximum number of custom soundboard sounds
// this guild may have.
func (g *Guild) MaxSoundboardSounds() int {
	switch g.PremiumTier {
	case GuildPremiumTierTier1:
		return 24
	case GuildPremiumTierTier2:
		return 36
	case GuildPremiumTierTier3:
		return 48
	default:
		return 8
	}
}

// ── Feature predicates ──────────────────────────────────────────────────────

// HasFeature reports whether the guild has the given feature flag enabled.
func (g *Guild) HasFeature(f GuildFeatures) bool {
	if g == nil {
		return false
	}
	for _, have := range g.Features {
		if have == f {
			return true
		}
	}
	return false
}

// Partnered reports whether this is a Discord Partner server.
func (g *Guild) Partnered() bool { return g.HasFeature(GuildFeaturePartnered) }

// Verified reports whether this guild is verified.
func (g *Guild) Verified() bool { return g.HasFeature(GuildFeatureVerified) }

// Community reports whether Community features are enabled for this guild.
func (g *Guild) Community() bool { return g.HasFeature(GuildFeatureCommunity) }

// ── Misc ────────────────────────────────────────────────────────────────────

// NameAcronym returns the initials Discord shows in place of a missing guild
// icon — the first letter of each word, with a leading "'s" possessive folded
// into the preceding word ("Mike's Cool Server" → "MCS").
func (g *Guild) NameAcronym() string {
	if g == nil {
		return ""
	}
	name := strings.ReplaceAll(g.Name, "'s ", " ")
	var b strings.Builder
	for _, word := range strings.Fields(name) {
		b.WriteString(string([]rune(word)[:1]))
	}
	return b.String()
}

// VanityURL returns the full https://discord.gg/<code> invite for this guild's
// vanity URL, or "" when the guild has no vanity code configured.
func (g *Guild) VanityURL() string {
	if g == nil || g.VanityUrlCode == nil || *g.VanityUrlCode == "" {
		return ""
	}
	return "https://discord.gg/" + *g.VanityUrlCode
}
