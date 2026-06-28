package testdata

import (
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewGuildMember() discord.GuildMember {
	return discord.GuildMember{
		GuildID:                    discord.RandomSnowflake(),
		Avatar:                     util.PointerOf(testutil.RandomString(32)),
		Banner:                     util.PointerOf(testutil.RandomString(32)),
		CommunicationDisabledUntil: util.PointerOf(testutil.RandomTime()),
		Deaf:                       testutil.RandomBool(),
		Flags: testutil.RandomFlags(
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
		JoinedAt:     util.PointerOf(testutil.RandomTime()),
		Mute:         testutil.RandomBool(),
		Nick:         util.PointerOf(testutil.RandomString(32)),
		Pending:      testutil.RandomBool(),
		PremiumSince: util.PointerOf(testutil.RandomTime()),
		Roles:        testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(1, 500)),
		User:         util.PointerOf(NewUser()),
		Permissions:  util.PointerOf(testutil.RandomFlags(testutil.AllPermissions...)),
		AvatarDecorationData: &discord.AvatarDecorationData{
			Asset: testutil.RandomString(32),
			SkuID: discord.RandomSnowflake(),
		},
		Collectibles: util.PointerOf(NewCollectible()),
	}
}

func NewUser() discord.User {
	return discord.User{
		AccentColor:          util.PointerOf(testutil.RandomNumberInRange(0x000000, 0xFFFFFF)),
		Avatar:               util.PointerOf(testutil.RandomString(32)),
		AvatarDecorationData: util.PointerOf(NewAvatarDecorationData()),
		Bot:                  util.PointerOf(testutil.RandomBool()),
		Discriminator:        strconv.Itoa(testutil.RandomNumberInRange(0000, 9999)),
		Flags: testutil.RandomFlags(
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
		GlobalName: util.PointerOf(testutil.RandomString(testutil.RandomNumberInRange(1, 32))),
		ID:         discord.RandomSnowflake(),
		Locale: util.PointerOf(testutil.RandomItem(
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
		)),
		MFAEnabled: util.PointerOf(testutil.RandomBool()),
		PrimaryGuild: &discord.PrimaryGuild{
			Badge:           util.PointerOf(testutil.RandomString(32)),
			IdentityEnabled: util.PointerOf(testutil.RandomBool()),
			IdentityGuildID: util.PointerOf(discord.RandomSnowflake()),
			Tag:             util.PointerOf(testutil.RandomString(testutil.RandomNumberInRange(1, 4))),
		},
		PublicFlags: testutil.RandomFlags(
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
		System:   util.PointerOf(testutil.RandomBool()),
		Username: testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		Banner:   util.PointerOf(testutil.RandomString(32)),
		Verified: util.PointerOf(testutil.RandomBool()),
		Email:    util.PointerOf(testutil.RandomString(10) + "@streame.gg"),
		PremiumType: util.PointerOf(testutil.RandomItem(
			discord.UserPremiumTypeNone,
			discord.UserPremiumTypeNitroClassic,
			discord.UserPremiumTypeNitro,
			discord.UserPremiumTypeNitroBasic,
		)),
		Collectibles: util.PointerOf(NewCollectible()),
	}
}

func NewCollectible() discord.Collectible {
	return discord.Collectible{
		Nameplate: &discord.Nameplate{
			SkuID: discord.RandomSnowflake(),
			Asset: testutil.RandomString(32),
			Label: testutil.RandomString(32),
			Palette: testutil.RandomItem(
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

func NewAvatarDecorationData() discord.AvatarDecorationData {
	return discord.AvatarDecorationData{
		Asset: testutil.RandomString(32),
		SkuID: discord.RandomSnowflake(),
	}
}
