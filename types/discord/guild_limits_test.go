package discord

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// guildLimitsSuite covers the tier-derived capacity helpers and feature
// predicates on Guild.
type guildLimitsSuite struct {
	suite.Suite
}

func TestGuildLimitsSuite(t *testing.T) {
	suite.Run(t, new(guildLimitsSuite))
}

func (s *guildLimitsSuite) TestMaxEmojis() {
	s.Equal(50, (&Guild{}).MaxEmojis())
	s.Equal(100, (&Guild{PremiumTier: GuildPremiumTierTier1}).MaxEmojis())
	s.Equal(150, (&Guild{PremiumTier: GuildPremiumTierTier2}).MaxEmojis())
	s.Equal(250, (&Guild{PremiumTier: GuildPremiumTierTier3}).MaxEmojis())
}

func (s *guildLimitsSuite) TestMaxStickers() {
	s.Equal(5, (&Guild{}).MaxStickers())
	s.Equal(15, (&Guild{Features: []GuildFeatures{GuildFeatureMoreStickers}}).MaxStickers())
	s.Equal(15, (&Guild{PremiumTier: GuildPremiumTierTier1}).MaxStickers())
	s.Equal(30, (&Guild{PremiumTier: GuildPremiumTierTier2}).MaxStickers())
	s.Equal(60, (&Guild{PremiumTier: GuildPremiumTierTier3}).MaxStickers())
}

func (s *guildLimitsSuite) TestMaxBitrate() {
	s.Equal(96_000, (&Guild{}).MaxBitrate())
	s.Equal(128_000, (&Guild{PremiumTier: GuildPremiumTierTier1}).MaxBitrate())
	s.Equal(256_000, (&Guild{PremiumTier: GuildPremiumTierTier2}).MaxBitrate())
	s.Equal(384_000, (&Guild{PremiumTier: GuildPremiumTierTier3}).MaxBitrate())
	// VIP_REGIONS unlocks the max regardless of tier.
	s.Equal(384_000, (&Guild{Features: []GuildFeatures{GuildFeatureVipRegions}}).MaxBitrate())
}

func (s *guildLimitsSuite) TestMaxSoundboardSounds() {
	s.Equal(8, (&Guild{}).MaxSoundboardSounds())
	s.Equal(24, (&Guild{PremiumTier: GuildPremiumTierTier1}).MaxSoundboardSounds())
	s.Equal(36, (&Guild{PremiumTier: GuildPremiumTierTier2}).MaxSoundboardSounds())
	s.Equal(48, (&Guild{PremiumTier: GuildPremiumTierTier3}).MaxSoundboardSounds())
}

func (s *guildLimitsSuite) TestFeaturePredicates() {
	s.False((*Guild)(nil).HasFeature(GuildFeaturePartnered))
	g := &Guild{Features: []GuildFeatures{GuildFeaturePartnered, GuildFeatureCommunity}}
	s.True(g.Partnered())
	s.True(g.Community())
	s.False(g.Verified())
	s.True(g.HasFeature(GuildFeaturePartnered))
	s.False(g.HasFeature(GuildFeatureVerified))
}

func (s *guildLimitsSuite) TestNameAcronym() {
	s.Equal("", (*Guild)(nil).NameAcronym())
	s.Equal("MCS", (&Guild{Name: "Mike's Cool Server"}).NameAcronym())
	s.Equal("GDW", (&Guild{Name: "Go Discord Wrapper"}).NameAcronym())
	s.Equal("S", (&Guild{Name: "Solo"}).NameAcronym())
	s.Equal("", (&Guild{Name: "   "}).NameAcronym())
}

func (s *guildLimitsSuite) TestVanityURL() {
	s.Equal("", (*Guild)(nil).VanityURL())
	s.Equal("", (&Guild{}).VanityURL())
	s.Equal("", (&Guild{VanityUrlCode: strptr("")}).VanityURL())
	s.Equal("https://discord.gg/golang", (&Guild{VanityUrlCode: strptr("golang")}).VanityURL())
}
