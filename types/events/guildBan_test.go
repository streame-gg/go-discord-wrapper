package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildBanAdd() {
	s.T().Log("Testing Guild Ban Add Unmarshal Logic")

	sub := testutil.InitSub[GuildBanAddEvent](s)

	sub.RunCommonEdgeCases()

	ban := testdata.NewGuildBanAddOrRemove()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildBanAddEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(ban),
			Validate: func(got GuildBanAddEvent) {
				s.EqualValues(ban["guild_id"], got.GuildID)
				guildMemberUser := ban["user"].(map[string]interface{})
				s.EqualValues(guildMemberUser["id"], got.User.ID)
				s.EqualValues(guildMemberUser["username"], got.User.Username)
				s.EqualValues(guildMemberUser["discriminator"], got.User.Discriminator)
				s.EqualValues(guildMemberUser["avatar"], *got.User.Avatar)
				s.EqualValues(guildMemberUser["bot"], *got.User.Bot)
				s.EqualValues(guildMemberUser["system"], *got.User.System)
				s.EqualValues(guildMemberUser["mfa_enabled"], *got.User.MFAEnabled)
				s.EqualValues(guildMemberUser["banner"], *got.User.Banner)
				s.EqualValues(guildMemberUser["accent_color"], *got.User.AccentColor)
				s.EqualValues(guildMemberUser["locale"], *got.User.Locale)
				s.EqualValues(guildMemberUser["verified"], *got.User.Verified)
				s.EqualValues(guildMemberUser["email"], *got.User.Email)
				s.EqualValues(guildMemberUser["flags"], got.User.Flags)
				s.EqualValues(guildMemberUser["premium_type"], *got.User.PremiumType)
				s.EqualValues(guildMemberUser["public_flags"], got.User.PublicFlags)
				s.EqualValues(guildMemberUser["global_name"], *got.User.GlobalName)
			},
		},
	})
}

func (s *eventSuite) TestGuildBanRemove() {
	s.T().Log("Testing Guild Ban Remove Unmarshal Logic")

	sub := testutil.InitSub[GuildBanRemoveEvent](s)

	sub.RunCommonEdgeCases()

	ban := testdata.NewGuildBanAddOrRemove()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildBanRemoveEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(ban),
			Validate: func(got GuildBanRemoveEvent) {
				s.EqualValues(ban["guild_id"], got.GuildID)
				guildMemberUser := ban["user"].(map[string]interface{})
				s.EqualValues(guildMemberUser["id"], got.User.ID)
				s.EqualValues(guildMemberUser["username"], got.User.Username)
				s.EqualValues(guildMemberUser["discriminator"], got.User.Discriminator)
				s.EqualValues(guildMemberUser["avatar"], *got.User.Avatar)
				s.EqualValues(guildMemberUser["bot"], *got.User.Bot)
				s.EqualValues(guildMemberUser["system"], *got.User.System)
				s.EqualValues(guildMemberUser["mfa_enabled"], *got.User.MFAEnabled)
				s.EqualValues(guildMemberUser["banner"], *got.User.Banner)
				s.EqualValues(guildMemberUser["accent_color"], *got.User.AccentColor)
				s.EqualValues(guildMemberUser["locale"], *got.User.Locale)
				s.EqualValues(guildMemberUser["verified"], *got.User.Verified)
				s.EqualValues(guildMemberUser["email"], *got.User.Email)
				s.EqualValues(guildMemberUser["flags"], got.User.Flags)
				s.EqualValues(guildMemberUser["premium_type"], *got.User.PremiumType)
				s.EqualValues(guildMemberUser["public_flags"], got.User.PublicFlags)
				s.EqualValues(guildMemberUser["global_name"], *got.User.GlobalName)
			},
		},
	})
}
