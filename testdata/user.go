package testdata

import (
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewGuildMember() map[string]interface{} {
	return map[string]interface{}{
		"avatar":                       testutil.RandomString(32),
		"banner":                       testutil.RandomString(32),
		"communication_disabled_until": testutil.RandomTime(),
		"deaf":                         testutil.RandomBool(),
		"flags": testutil.RandomFlags(
			discord.GuildMemberFlagDidRejoin,
			discord.GuildMemberFlagCompletedOnboarding,
			discord.GuildMemberFlagBypassesVerification,
			discord.GuildMemberFlagStartedOnboarding,
			discord.GuildMemberFlagIsGuest,
			discord.GuildMemberFlagStartedHomeActions,
			discord.GuildMemberFlagCompletedHomeActions,
			discord.GuildMemberFlagAutomodQuarantinedUsername,
			discord.GuildMemberFlagDmSettingsUpsellAcked,
			discord.GuildMemberFlagAutomodQuarantinedGuildTag,
		),
		"joined_at":              testutil.RandomTime(),
		"mute":                   testutil.RandomBool(),
		"nick":                   testutil.RandomString(32),
		"pending":                testutil.RandomBool(),
		"premium_since":          testutil.RandomTime(),
		"roles":                  testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(1, 500)),
		"user":                   NewUser(),
		"permissions":            testutil.RandomFlags(testutil.AllPermissions...),
		"avatar_decoration_data": NewAvatarDecorationData(),
		"collectibles":           NewCollectible(),
	}
}

func NewUser() map[string]interface{} {
	return map[string]interface{}{
		"accent_color":      testutil.RandomNumberInRange(0x000000, 0xFFFFFF),
		"avatar":            testutil.RandomString(32),
		"avatar_decoration": NewAvatarDecorationData(),
		"bot":               testutil.RandomBool(),
		"discriminator":     strconv.Itoa(testutil.RandomNumberInRange(0000, 9999)),
		"flags": testutil.RandomFlags(
			discord.UserFlagStaff,
			discord.UserFlagPartner,
			discord.UserFlagHypeSquad,
			discord.UserFlagBugHunterLevel1,
			discord.UserFlagHypeSquadHouseBravery,
			discord.UserFlagHypeSquadHouseBrillia,
			discord.UserFlagHypeSquadHouseBalance,
			discord.UserFlagPremiumEarlySupporter,
			discord.UserFlagTeamPseudoUser,
			discord.UserFlagBugHunterLevel2,
			discord.UserFlagVerifiedBot,
			discord.UserFlagVerifiedDeveloper,
			discord.UserFlagCertifiedModerator,
			discord.UserFlagBotHTTPInteractions,
		),
		"global_name": testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		"id":          discord.RandomSnowflake(),
		"locale": testutil.RandomItem(
			discord.LocaleIndonesian,
			discord.LocaleDanish,
			discord.LocaleGerman,
			discord.LocaleEnglishUS,
			discord.LocaleEnglishUK,
			discord.LocaleSpanish,
			discord.LocaleSpanishLATAM,
			discord.LocaleFrench,
			discord.LocaleCroatian,
			discord.LocaleItalian,
			discord.LocaleLithuanian,
			discord.LocaleHungarian,
			discord.LocaleDutch,
			discord.LocaleNorwegian,
			discord.LocalePolish,
			discord.LocalePortuguese,
			discord.LocaleRomanian,
			discord.LocaleFinnish,
			discord.LocaleSwedish,
			discord.LocaleVietnamese,
			discord.LocaleTurkish,
			discord.LocaleCzech,
			discord.LocaleGreek,
			discord.LocaleBulgarian,
			discord.LocaleRussian,
			discord.LocaleUkrainian,
			discord.LocaleHindi,
			discord.LocaleThai,
			discord.LocaleChineseCN,
			discord.LocaleJapanese,
			discord.LocaleChineseTW,
			discord.LocaleKorean,
		),
		"mfa_enabled": testutil.RandomBool(),
		"primary_guild": map[string]interface{}{
			"badge":             testutil.RandomString(32),
			"identity_enabled":  testutil.RandomBool(),
			"identity_guild_id": discord.RandomSnowflake(),
			"tag":               testutil.RandomString(testutil.RandomNumberInRange(1, 4)),
		},
		"public_flags": testutil.RandomFlags(
			discord.UserFlagStaff,
			discord.UserFlagPartner,
			discord.UserFlagHypeSquad,
			discord.UserFlagBugHunterLevel1,
			discord.UserFlagHypeSquadHouseBravery,
			discord.UserFlagHypeSquadHouseBrillia,
			discord.UserFlagHypeSquadHouseBalance,
			discord.UserFlagPremiumEarlySupporter,
			discord.UserFlagTeamPseudoUser,
			discord.UserFlagBugHunterLevel2,
			discord.UserFlagVerifiedBot,
			discord.UserFlagVerifiedDeveloper,
			discord.UserFlagCertifiedModerator,
			discord.UserFlagBotHTTPInteractions,
		),
		"system":   testutil.RandomBool(),
		"username": testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		"banner":   testutil.RandomString(32),
		"verified": testutil.RandomBool(),
		"email":    testutil.RandomString(10) + "@streame.gg",
		"premium_type": testutil.RandomItem(
			discord.UserPremiumTypeNone,
			discord.UserPremiumTypeNitroClassic,
			discord.UserPremiumTypeNitro,
			discord.UserPremiumTypeNitroBasic,
		),
		"collectibles": NewCollectible(),
	}
}

func NewCollectible() map[string]interface{} {
	return map[string]interface{}{
		"nameplate": map[string]interface{}{
			"sku_id": discord.RandomSnowflake(),
			"asset":  testutil.RandomString(32),
			"label":  testutil.RandomString(32),
			"palette": testutil.RandomItem(
				discord.NameplatePaletteCrimson,
				discord.NameplatePaletteBerry,
				discord.NameplatePaletteSky,
				discord.NameplatePaletteTeal,
				discord.NameplatePaletteForest,
				discord.NameplatePaletteBubbleGum,
				discord.NameplatePaletteViolet,
				discord.NameplatePaletteCobalt,
				discord.NameplatePaletteClover,
				discord.NameplatePaletteLemon,
				discord.NameplatePaletteWhite,
			),
		},
	}
}

func NewAvatarDecorationData() map[string]interface{} {
	return map[string]interface{}{
		"asset":  testutil.RandomString(32),
		"sku_id": discord.RandomSnowflake(),
	}
}
