package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewGuildBanAddOrRemove() map[string]interface{} {
	return map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"user":     NewUser(),
	}
}

func NewUnavailableGuild() map[string]interface{} {
	return map[string]interface{}{
		"unavailable": true,
		"id":          discord.RandomSnowflake(),
	}
}

// NewAvailableGuild TODO: payload is still missing some properties
func NewAvailableGuild() map[string]interface{} {
	return map[string]interface{}{
		"unavailable":       false,
		"id":                discord.RandomSnowflake(),
		"name":              testutil.RandomString(testutil.RandomNumberInRange(1, 100)),
		"icon":              testutil.RandomString(32),
		"icon_hash":         testutil.RandomString(32),
		"splash":            testutil.RandomString(32),
		"discovery_splash":  testutil.RandomString(32),
		"owner":             testutil.RandomBool(),
		"owner_id":          discord.RandomSnowflake(),
		"permissions":       testutil.RandomFlags(testutil.AllPermissions...),
		"region":            testutil.RandomString(10),
		"afk_channel_id":    discord.RandomSnowflake(),
		"afk_timeout":       testutil.RandomNumberInRange(0, 3600),
		"widget_enabled":    testutil.RandomBool(),
		"widget_channel_id": discord.RandomSnowflake(),
		"verification_level": testutil.RandomItem(
			discord.GuildVerificationLevelNone,
			discord.GuildVerificationLevelLow,
			discord.GuildVerificationLevelMedium,
			discord.GuildVerificationLevelHigh,
			discord.GuildVerificationLevelVeryHigh,
		),
		"default_message_notifications": testutil.RandomItem(
			discord.DefaultMessageNotificationLevelAllMessages,
			discord.DefaultMessageNotificationLevelOnlyMentions,
		),
		"explicit_content_filter": testutil.RandomItem(
			discord.GuildExplicitContentFilterLevelDisabled,
			discord.GuildExplicitContentFilterLevelMembersWithoutRoles,
			discord.GuildExplicitContentFilterLevelAllMembers,
		),
		"roles": testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 20), func(a *[]map[string]interface{}) {
			*a = append(*a, NewRole())
		}),
		"emojis": testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 20), func(a *[]map[string]interface{}) {
			*a = append(*a, NewEmoji())
		}),
		"features": testutil.RandomArray[discord.GuildFeatures](
			testutil.RandomNumberInRange(0, 5),
			discord.GuildFeatureAnimatedBanner,
			discord.GuildFeatureAnimatedIcon,
			discord.GuildFeatureCommunity,
			discord.GuildFeatureDiscoverable,
			discord.GuildFeatureFeaturable,
			discord.GuildFeatureInviteSplash,
			discord.GuildFeatureNews,
			discord.GuildFeaturePartnered,
			discord.GuildFeatureRoleIcons,
			discord.GuildFeatureVanityURL,
			discord.GuildFeatureVerified,
			discord.GuildFeatureVipRegions,
			discord.GuildFeatureWelcomeScreenEnabled,
		),
		"mfa_level": testutil.RandomItem(
			discord.GuildMFALevelNone,
			discord.GuildMFALevelElevated,
		),
		"application_id":    discord.RandomSnowflake(),
		"system_channel_id": discord.RandomSnowflake(),
		"system_channel_flags": testutil.RandomFlags(
			discord.GuildSystemChannelFlagsSuppressJoinNotifications,
			discord.GuildSystemChannelFlagsSuppressPremiumSubscriptions,
			discord.GuildSystemChannelFlagsSuppressGuildReminderMessages,
			discord.GuildSystemChannelFlagsSuppressJoinNotificationReplies,
			discord.GuildSystemChannelFlagsPurchaseNotifications,
			discord.GuildSystemChannelFlagsPurchaseNotificationReplies,
		),
		"rules_channel_id": discord.RandomSnowflake(),
		"max_presences":    testutil.RandomNumberInRange(0, 25000),
		"max_members":      testutil.RandomNumberInRange(1, 500000),
		"vanity_url_code":  testutil.RandomString(testutil.RandomNumberInRange(1, 20)),
		"description":      testutil.RandomString(testutil.RandomNumberInRange(1, 300)),
		"banner":           testutil.RandomString(32),
		"premium_tier": testutil.RandomItem(
			discord.GuildPremiumTierNone,
			discord.GuildPremiumTierTier1,
			discord.GuildPremiumTierTier2,
			discord.GuildPremiumTierTier3,
		),
		"premium_subscription_count": testutil.RandomNumberInRange(0, 100),
		"preferred_locale": testutil.RandomItem(
			discord.LocaleEnglishUS,
			discord.LocaleEnglishUK,
			discord.LocaleGerman,
			discord.LocaleFrench,
			discord.LocaleSpanish,
		),
		"public_updates_channel_id":     discord.RandomSnowflake(),
		"max_video_channel_users":       testutil.RandomNumberInRange(0, 25),
		"max_stage_video_channel_users": testutil.RandomNumberInRange(0, 50),
		"approximate_member_count":      testutil.RandomNumberInRange(1, 100000),
		"approximate_presence_count":    testutil.RandomNumberInRange(0, 100000),
		"welcome_screen":                NewWelcomeScreen(),
		"nsfw_level": testutil.RandomItem(
			discord.GuildNSFWLevelDefault,
			discord.GuildNSFWLevelExplicit,
			discord.GuildNSFWLevelSafe,
			discord.GuildNSFWLevelAgeRestricted,
		),
		"stickers": testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 10), func(a *[]map[string]interface{}) {
			*a = append(*a, NewSticker())
		}),
		"premium_progress_bar_enabled": testutil.RandomBool(),
		"safety_alerts_channel_id":     discord.RandomSnowflake(),
		"incidents_data":               NewIncidentsData(),
	}
}

func NewRole() map[string]interface{} {
	return map[string]interface{}{
		"id":   discord.RandomSnowflake(),
		"name": testutil.RandomString(testutil.RandomNumberInRange(1, 100)),
		"colors": map[string]interface{}{
			"primary_color":   testutil.RandomNumberInRange(0x000000, 0xFFFFFF),
			"secondary_color": testutil.RandomNumberInRange(0x000000, 0xFFFFFF),
			"tertiary_color":  testutil.RandomNumberInRange(0x000000, 0xFFFFFF),
		},
		"hoist":         testutil.RandomBool(),
		"icon":          testutil.RandomString(32),
		"unicode_emoji": testutil.RandomString(4),
		"position":      testutil.RandomNumberInRange(0, 100),
		"permissions":   testutil.RandomFlags(testutil.AllPermissions...),
		"managed":       testutil.RandomBool(),
		"mentionable":   testutil.RandomBool(),
		"tags": map[string]interface{}{
			"bot_id":                  discord.RandomSnowflake(),
			"integration_id":          discord.RandomSnowflake(),
			"subscription_listing_id": discord.RandomSnowflake(),
		},
		"flags": testutil.RandomFlags(discord.RoleFlagsInPrompt),
	}
}

func NewEmoji() map[string]interface{} {
	return map[string]interface{}{
		"id":             discord.RandomSnowflake(),
		"name":           testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		"roles":          testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(0, 10)),
		"user":           NewUser(),
		"require_colons": testutil.RandomBool(),
		"managed":        testutil.RandomBool(),
		"animated":       testutil.RandomBool(),
		"available":      testutil.RandomBool(),
	}
}

func NewSticker() map[string]interface{} {
	return map[string]interface{}{
		"id":          discord.RandomSnowflake(),
		"pack_id":     discord.RandomSnowflake(),
		"name":        testutil.RandomString(testutil.RandomNumberInRange(2, 30)),
		"description": testutil.RandomString(testutil.RandomNumberInRange(1, 100)),
		"tags":        testutil.RandomString(testutil.RandomNumberInRange(1, 200)),
		"type": testutil.RandomItem(
			discord.StickerTypeStandard,
			discord.StickerTypeGuild,
		),
		"format_type": testutil.RandomItem(
			discord.StickerFormatTypePNG,
			discord.StickerFormatTypeAPNG,
			discord.StickerFormatTypeLottie,
			discord.StickerFormatTypeGIF,
		),
		"available":  testutil.RandomBool(),
		"guild_id":   discord.RandomSnowflake(),
		"user":       NewUser(),
		"sort_value": testutil.RandomNumberInRange(0, 100),
	}
}

func NewWelcomeScreen() map[string]interface{} {
	return map[string]interface{}{
		"description": testutil.RandomString(testutil.RandomNumberInRange(1, 140)),
		"welcome_channels": testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 5), func(a *[]map[string]interface{}) {
			*a = append(*a, NewWelcomeScreenChannel())
		}),
	}
}

func NewWelcomeScreenChannel() map[string]interface{} {
	return map[string]interface{}{
		"channel_id":  discord.RandomSnowflake(),
		"description": testutil.RandomString(testutil.RandomNumberInRange(1, 50)),
		"emoji_id":    discord.RandomSnowflake(),
		"emoji_name":  testutil.RandomString(4),
	}
}

func NewIncidentsData() map[string]interface{} {
	return map[string]interface{}{
		"invites_disabled_until": testutil.RandomTime(),
		"dms_disabled_until":     testutil.RandomTime(),
		"dm_spam_detected_at":    testutil.RandomTime(),
		"raid_detected_at":       testutil.RandomTime(),
	}
}
