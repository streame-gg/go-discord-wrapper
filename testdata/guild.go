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

func NewVoiceState() map[string]interface{} {
	return map[string]interface{}{
		"guild_id":                   discord.RandomSnowflake(),
		"channel_id":                 discord.RandomSnowflake(),
		"user_id":                    discord.RandomSnowflake(),
		"member":                     NewGuildMember(),
		"session_id":                 testutil.RandomString(32),
		"deaf":                       testutil.RandomBool(),
		"mute":                       testutil.RandomBool(),
		"self_deaf":                  testutil.RandomBool(),
		"self_mute":                  testutil.RandomBool(),
		"self_stream":                testutil.RandomBool(),
		"self_video":                 testutil.RandomBool(),
		"suppress":                   testutil.RandomBool(),
		"request_to_speak_timestamp": testutil.RandomTime(),
	}
}

func NewStageInstance() map[string]interface{} {
	return map[string]interface{}{
		"id":         discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"channel_id": discord.RandomSnowflake(),
		"topic":      testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"privacy_level": testutil.RandomItem(
			discord.StageInstancePrivacyLevelPublic,
			discord.StageInstancePrivacyLevelGuildOnly,
		),
		"discoverable_disabled":    testutil.RandomBool(),
		"guild_scheduled_event_id": discord.RandomSnowflake(),
	}
}

func NewSoundboardSound() map[string]interface{} {
	return map[string]interface{}{
		"name":       testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"sound_id":   discord.RandomSnowflake(),
		"volume":     testutil.RandomFloat64InRange(0, 1),
		"emoji_id":   discord.RandomSnowflake(),
		"emoji_name": testutil.RandomString(4),
		"guild_id":   discord.RandomSnowflake(),
		"available":  testutil.RandomBool(),
		"user":       NewUser(),
	}
}

func NewAvailableGuildWithGuildCreateValues() map[string]interface{} {
	withoutValues := NewAvailableGuild()
	withoutValues["joined_at"] = testutil.RandomTime()
	withoutValues["large"] = testutil.RandomBool()
	withoutValues["member_count"] = testutil.RandomIntInRange(2, 10000000)
	withoutValues["voice_states"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewVoiceState())
	})
	withoutValues["members"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewGuildMember())
	})
	withoutValues["channels"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewChannel())
	})
	withoutValues["threads"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewChannel())
	})
	withoutValues["presences"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewPresence())
	})
	withoutValues["stage_instances"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewStageInstance())
	})
	withoutValues["guild_scheduled_events"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewScheduledEvent())
	})
	withoutValues["soundboard_sounds"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]map[string]interface{}) {
		*arrayToFill = append(*arrayToFill, NewSoundboardSound())
	})

	return withoutValues
}

// NewAvailableGuild TODO: payload is still missing some properties
func NewAvailableGuild() map[string]interface{} {
	return map[string]interface{}{
		"unavailable":       false,
		"id":                discord.RandomSnowflake(),
		"name":              testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"icon":              testutil.RandomString(32),
		"icon_hash":         testutil.RandomString(32),
		"splash":            testutil.RandomString(32),
		"discovery_splash":  testutil.RandomString(32),
		"owner":             testutil.RandomBool(),
		"owner_id":          discord.RandomSnowflake(),
		"permissions":       testutil.RandomFlags(testutil.AllPermissions...),
		"region":            testutil.RandomString(10),
		"afk_channel_id":    discord.RandomSnowflake(),
		"afk_timeout":       testutil.RandomIntInRange(0, 3600),
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
		"roles": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 20), func(a *[]map[string]interface{}) {
			*a = append(*a, NewRole())
		}),
		"emojis": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 20), func(a *[]map[string]interface{}) {
			*a = append(*a, NewEmoji())
		}),
		"features": testutil.RandomArray[discord.GuildFeatures](
			testutil.RandomIntInRange(0, 5),
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
		"max_presences":    testutil.RandomIntInRange(0, 25000),
		"max_members":      testutil.RandomIntInRange(1, 500000),
		"vanity_url_code":  testutil.RandomString(testutil.RandomIntInRange(1, 20)),
		"description":      testutil.RandomString(testutil.RandomIntInRange(1, 300)),
		"banner":           testutil.RandomString(32),
		"premium_tier": testutil.RandomItem(
			discord.GuildPremiumTierNone,
			discord.GuildPremiumTierTier1,
			discord.GuildPremiumTierTier2,
			discord.GuildPremiumTierTier3,
		),
		"premium_subscription_count": testutil.RandomIntInRange(0, 100),
		"preferred_locale": testutil.RandomItem(
			discord.LocaleEnglishUS,
			discord.LocaleEnglishUK,
			discord.LocaleGerman,
			discord.LocaleFrench,
			discord.LocaleSpanish,
		),
		"public_updates_channel_id":     discord.RandomSnowflake(),
		"max_video_channel_users":       testutil.RandomIntInRange(0, 25),
		"max_stage_video_channel_users": testutil.RandomIntInRange(0, 50),
		"approximate_member_count":      testutil.RandomIntInRange(1, 100000),
		"approximate_presence_count":    testutil.RandomIntInRange(0, 100000),
		"welcome_screen":                NewWelcomeScreen(),
		"nsfw_level": testutil.RandomItem(
			discord.GuildNSFWLevelDefault,
			discord.GuildNSFWLevelExplicit,
			discord.GuildNSFWLevelSafe,
			discord.GuildNSFWLevelAgeRestricted,
		),
		"stickers": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(a *[]map[string]interface{}) {
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
		"name": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"colors": map[string]interface{}{
			"primary_color":   testutil.RandomIntInRange(0x000000, 0xFFFFFF),
			"secondary_color": testutil.RandomIntInRange(0x000000, 0xFFFFFF),
			"tertiary_color":  testutil.RandomIntInRange(0x000000, 0xFFFFFF),
		},
		"hoist":         testutil.RandomBool(),
		"icon":          testutil.RandomString(32),
		"unicode_emoji": testutil.RandomString(4),
		"position":      testutil.RandomIntInRange(0, 100),
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
		"name":           testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"roles":          testutil.RandomSnowflakeArray(testutil.RandomIntInRange(0, 10)),
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
		"name":        testutil.RandomString(testutil.RandomIntInRange(2, 30)),
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"tags":        testutil.RandomString(testutil.RandomIntInRange(1, 200)),
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
		"sort_value": testutil.RandomIntInRange(0, 100),
	}
}

func NewWelcomeScreen() map[string]interface{} {
	return map[string]interface{}{
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 140)),
		"welcome_channels": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 5), func(a *[]map[string]interface{}) {
			*a = append(*a, NewWelcomeScreenChannel())
		}),
	}
}

func NewWelcomeScreenChannel() map[string]interface{} {
	return map[string]interface{}{
		"channel_id":  discord.RandomSnowflake(),
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 50)),
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

func NewMemberUpdateEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"guild_id":                     discord.RandomSnowflake(),
		"roles":                        testutil.RandomSnowflakeArray(testutil.RandomIntInRange(0, 500)),
		"user":                         NewUser(),
		"nick":                         testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"avatar":                       testutil.RandomString(32),
		"banner":                       testutil.RandomString(32),
		"joined_at":                    testutil.RandomTime(),
		"premium_since":                testutil.RandomTime(),
		"deaf":                         testutil.RandomBool(),
		"mute":                         testutil.RandomBool(),
		"pending":                      testutil.RandomBool(),
		"communication_disabled_until": testutil.RandomTime(),
		"avatar_decoration_data":       NewAvatarDecorationData(),
		"collectibles":                 NewCollectible(),
	}
}

func NewScheduledEvent() map[string]interface{} {
	return map[string]interface{}{
		"id":                   discord.RandomSnowflake(),
		"guild_id":             discord.RandomSnowflake(),
		"channel_id":           discord.RandomSnowflake(),
		"creator_id":           discord.RandomSnowflake(),
		"name":                 testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"description":          testutil.RandomString(testutil.RandomIntInRange(1, 1000)),
		"scheduled_start_time": testutil.RandomTime(),
		"scheduled_end_time":   testutil.RandomTime(),
		"privacy_level": testutil.RandomItem(
			discord.GuildScheduledEventPrivacyLevelGuildOnly,
		),
		"status": testutil.RandomItem(
			discord.GuildScheduledEventStatusScheduled,
			discord.GuildScheduledEventStatusActive,
			discord.GuildScheduledEventStatusCompleted,
			discord.GuildScheduledEventStatusCanceled,
		),
		"entity_type": testutil.RandomItem(
			discord.GuildScheduledEventEntityTypeStageInstance,
			discord.GuildScheduledEventEntityTypeVoice,
			discord.GuildScheduledEventEntityTypeExternal,
		),
		"entity_id": discord.RandomSnowflake(),
		"entity_metadata": map[string]interface{}{
			"location": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		},
		"creator":    NewUser(),
		"user_count": testutil.RandomIntInRange(1, 1000),
		"image":      testutil.RandomString(32),
		"recurrence_rule": map[string]interface{}{
			"start": testutil.RandomTime(),
			"end":   testutil.RandomTime(),
			"frequency": testutil.RandomItem(
				discord.GuildScheduledEventRecurrenceRuleFrequencyDaily,
				discord.GuildScheduledEventRecurrenceRuleFrequencyWeekly,
				discord.GuildScheduledEventRecurrenceRuleFrequencyMonthly,
				discord.GuildScheduledEventRecurrenceRuleFrequencyYearly,
			),
			"interval": testutil.RandomIntInRange(1, 4),
			"by_weekday": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 7), func(a *[]discord.GuildScheduledEventRecurrenceRuleWeekday) {
				*a = append(*a, testutil.RandomItem(
					discord.GuildScheduledEventRecurrenceRuleWeekdayMonday,
					discord.GuildScheduledEventRecurrenceRuleWeekdayTuesday,
					discord.GuildScheduledEventRecurrenceRuleWeekdayWednesday,
					discord.GuildScheduledEventRecurrenceRuleWeekdayThursday,
					discord.GuildScheduledEventRecurrenceRuleWeekdayFriday,
					discord.GuildScheduledEventRecurrenceRuleWeekdaySaturday,
					discord.GuildScheduledEventRecurrenceRuleWeekdaySunday,
				))
			}),
			"by_n_weekday": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 7), func(arrayToFill *[]map[string]interface{}) {
				*arrayToFill = append(*arrayToFill, map[string]interface{}{
					"n":   testutil.RandomIntInRange(1, 5),
					"day": testutil.RandomIntInRange(1, 7),
				})
			}),
			"by_month": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 7), func(a *[]discord.GuildScheduledEventRecurrenceRuleMonth) {
				*a = append(*a, testutil.RandomItem(
					discord.GuildScheduledEventRecurrenceRuleMonthJanuary,
					discord.GuildScheduledEventRecurrenceRuleMonthFebruary,
					discord.GuildScheduledEventRecurrenceRuleMonthMarch,
					discord.GuildScheduledEventRecurrenceRuleMonthApril,
					discord.GuildScheduledEventRecurrenceRuleMonthMay,
					discord.GuildScheduledEventRecurrenceRuleMonthJune,
					discord.GuildScheduledEventRecurrenceRuleMonthJuly,
					discord.GuildScheduledEventRecurrenceRuleMonthAugust,
					discord.GuildScheduledEventRecurrenceRuleMonthSeptember,
					discord.GuildScheduledEventRecurrenceRuleMonthOctober,
					discord.GuildScheduledEventRecurrenceRuleMonthNovember,
					discord.GuildScheduledEventRecurrenceRuleMonthDecember,
				))
			}),
			"by_month_day": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 31), func(a *[]int) {
				*a = append(*a, testutil.RandomIntInRange(1, 31))
			}),
			"by_year_day": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 365), func(a *[]int) {
				*a = append(*a, testutil.RandomIntInRange(1, 365))
			}),
			"count": testutil.RandomIntInRange(1, 5),
		},
	}
}

func NewIntegrationWithGuildID() map[string]interface{} {
	integration := NewIntegration()
	integration["guild_id"] = discord.RandomSnowflake()
	return integration
}

func NewIntegration() map[string]interface{} {
	return map[string]interface{}{
		"id":   discord.RandomSnowflake(),
		"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"type": testutil.RandomItem(
			discord.IntegrationTypeTwitch,
			discord.IntegrationTypeDiscord,
			discord.IntegrationTypeTwitch,
			discord.IntegrationTypeGuildSubscription,
		),
		"enabled":          testutil.RandomBool(),
		"syncing":          testutil.RandomBool(),
		"role_id":          discord.RandomSnowflake(),
		"enable_emoticons": testutil.RandomBool(),
		"expire_behavior": testutil.RandomItem(
			discord.IntegrationExpireBehaviorRemoveRole,
			discord.IntegrationExpireBehaviorKick,
		),
		"expire_grace_period": testutil.RandomIntInRange(1, 30),
		"user":                NewUser(),
		"application": map[string]interface{}{
			"id":          discord.RandomSnowflake(),
			"name":        testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			"icon":        testutil.RandomString(32),
			"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
			"bot":         NewUser(),
		},
		"synced_at":        testutil.RandomTime(),
		"subscriber_count": testutil.RandomIntInRange(1, 100),
		"revoked":          testutil.RandomBool(),
		"account": map[string]interface{}{
			"id":   testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		},
		"scopes": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 32), func(arrayToFill *[]discord.Scope) {
			*arrayToFill = append(*arrayToFill, testutil.RandomItem(
				discord.ScopeActivitiesRead,
				discord.ScopeActivitiesWrite,
				discord.ScopeApplicationsBuildRead,
				discord.ScopeApplicationsBuildUpload,
				discord.ScopeApplicationCommands,
				discord.ScopeApplicationCommandsUpdate,
				discord.ScopeApplicationCommandsPermissionsUpdate,
				discord.ScopeApplicationEntitlements,
				discord.ScopeApplicationsStoreUpdate,
				discord.ScopeBot,
				discord.ScopeConnections,
				discord.ScopeDMChannelsRead,
				discord.ScopeEmail,
				discord.ScopeGdmJoin,
				discord.ScopeGuilds,
				discord.ScopeGuildsJoin,
				discord.ScopeGuildMembersRead,
				discord.ScopeIdentify,
				discord.ScopeIdentifyPremium,
				discord.ScopeMessagesRead,
				discord.ScopeRelationshipsRead,
				discord.ScopeRoleConnectionsWrite,
				discord.ScopeRPC,
				discord.ScopeRPCActivitiesWrite,
				discord.ScopeRPCNotificationsRead,
				discord.ScopeRPCVoiceRead,
				discord.ScopeRPCVoiceWrite,
				discord.ScopeVoice,
				discord.ScopeWebhookIncoming,
			))
		}),
	}
}
