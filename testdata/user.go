package testdata

import (
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewGuildMemberWithGuildID() map[string]interface{} {
	return map[string]interface{}{
		"guild_id":                     discord.RandomSnowflake(),
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
		"roles":                  testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 500)),
		"user":                   NewUser(),
		"permissions":            testutil.RandomFlags(testutil.AllPermissions...),
		"avatar_decoration_data": NewAvatarDecorationData(),
		"collectibles":           NewCollectible(),
	}
}

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
		"roles":                  testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 500)),
		"user":                   NewUser(),
		"permissions":            testutil.RandomFlags(testutil.AllPermissions...),
		"avatar_decoration_data": NewAvatarDecorationData(),
		"collectibles":           NewCollectible(),
	}
}

func NewUser() map[string]interface{} {
	return map[string]interface{}{
		"accent_color":           testutil.RandomIntInRange(0x000000, 0xFFFFFF),
		"avatar":                 testutil.RandomString(32),
		"avatar_decoration_data": NewAvatarDecorationData(),
		"bot":                    testutil.RandomBool(),
		"discriminator":          strconv.Itoa(testutil.RandomIntInRange(0000, 9999)),
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
		"global_name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
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
			"tag":               testutil.RandomString(testutil.RandomIntInRange(1, 4)),
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
		"username": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
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

func NewUserWithGuildID() map[string]interface{} {
	return map[string]interface{}{
		"guild_id":               discord.RandomSnowflake(),
		"accent_color":           testutil.RandomIntInRange(0x000000, 0xFFFFFF),
		"avatar":                 testutil.RandomString(32),
		"avatar_decoration_data": NewAvatarDecorationData(),
		"bot":                    testutil.RandomBool(),
		"discriminator":          strconv.Itoa(testutil.RandomIntInRange(0000, 9999)),
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
		"global_name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
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
			"tag":               testutil.RandomString(testutil.RandomIntInRange(1, 4)),
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
		"username": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
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

func NewGuildMembersChunkEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"guild_id":    discord.RandomSnowflake(),
		"chunk_index": testutil.RandomIntInRange(1, 100),
		"chunk_count": testutil.RandomIntInRange(1, 100),
		"not_found": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]interface{}) {
			*arrayToFill = append(*arrayToFill, testutil.RandomString(testutil.RandomIntInRange(1, 5)))
		}),
		"presences": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]interface{}) {
			*arrayToFill = append(*arrayToFill, NewPresence())
		}),
		"nonce": testutil.RandomString(32),
		"members": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 100), func(arrayToFill *[]interface{}) {
			*arrayToFill = append(*arrayToFill, NewGuildMember())
		}),
	}
}

func NewPresenceActivity() map[string]interface{} {
	return map[string]interface{}{
		"type": testutil.RandomItem(
			discord.ActivityTypePlaying,
			discord.ActivityTypeStreaming,
			discord.ActivityTypeListening,
			discord.ActivityTypeWatching,
			discord.ActivityTypeCustom,
			discord.ActivityTypeCompeting,
		),
		"name":       testutil.RandomString(32),
		"url":        testutil.RandomString(32),
		"created_at": testutil.RandomTime().Unix(),
		"timestamps": map[string]interface{}{
			"start": testutil.RandomTime().Unix(),
			"end":   testutil.RandomTime().Unix(),
		},
		"application_id": discord.RandomSnowflake(),
		"status_display_type": testutil.RandomItem(
			discord.StatusDisplayTypeName,
			discord.StatusDisplayTypeState,
			discord.StatusDisplayTypeDetails,
		),
		"details":     testutil.RandomString(32),
		"details_url": testutil.RandomString(32),
		"state":       testutil.RandomString(32),
		"state_url":   testutil.RandomString(32),
		"emoji":       NewEmoji(),
		"party": map[string]interface{}{
			"id":   testutil.RandomString(32),
			"size": []int{testutil.RandomIntInRange(1, 100), testutil.RandomIntInRange(1, 100)},
		},
		"assets": map[string]interface{}{
			"large_image":        testutil.RandomString(32),
			"large_text":         testutil.RandomString(32),
			"large_url":          testutil.RandomString(32),
			"small_image":        testutil.RandomString(32),
			"small_text":         testutil.RandomString(32),
			"small_url":          testutil.RandomString(32),
			"invite_cover_image": testutil.RandomString(32),
		},
		"secrets": map[string]interface{}{
			"join":     testutil.RandomString(32),
			"spectate": testutil.RandomString(32),
			"match":    testutil.RandomString(32),
		},
		"instance": testutil.RandomBool(),
		"flags": testutil.RandomFlags(
			discord.ActivityFlagInstance,
			discord.ActivityFlagJoin,
			discord.ActivityFlagSpectate,
			discord.ActivityFlagJoinRequest,
			discord.ActivityFlagSync,
			discord.ActivityFlagPlay,
			discord.ActivityFlagPartyPrivacyFriends,
			discord.ActivityFlagPartyPrivacyVoiceChannel,
			discord.ActivityFlagEmbedded,
		),
		"buttons": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 2), func(arrayToFill *[]interface{}) {
			*arrayToFill = append(*arrayToFill, NewPresenceActivityButton())
		}),
	}
}

func NewPresenceActivityButton() map[string]interface{} {
	return map[string]interface{}{
		"label": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"url":   testutil.RandomString(testutil.RandomIntInRange(1, 512)),
	}
}

func NewPresence() map[string]interface{} {
	return map[string]interface{}{
		"user":     NewUser(),
		"guild_id": discord.RandomSnowflake(),
		"status": testutil.RandomItem(
			discord.PresenceStatusIdle,
			discord.PresenceStatusDND,
			discord.PresenceStatusOnline,
			discord.PresenceStatusOffline,
		),
		"activities": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]interface{}) {
			*arrayToFill = append(*arrayToFill, NewPresenceActivity())
		}),
		"client_status": map[string]interface{}{
			"desktop": testutil.RandomItem(
				discord.PresenceStatusIdle,
				discord.PresenceStatusDND,
				discord.PresenceStatusOnline,
				discord.PresenceStatusOffline,
			),
			"mobile": testutil.RandomItem(
				discord.PresenceStatusIdle,
				discord.PresenceStatusDND,
				discord.PresenceStatusOnline,
				discord.PresenceStatusOffline,
			),
			"web": testutil.RandomItem(
				discord.PresenceStatusIdle,
				discord.PresenceStatusDND,
				discord.PresenceStatusOnline,
				discord.PresenceStatusOffline,
			),
		},
	}
}

func NewThreadMember() map[string]interface{} {
	return map[string]interface{}{
		"id":             discord.RandomSnowflake(),
		"user_id":        discord.RandomSnowflake(),
		"join_timestamp": testutil.RandomTime(),
		"flags":          testutil.RandomIntInRange(0, 1),
		"member":         NewGuildMember(),
	}
}

func NewApplication() map[string]interface{} {
	return map[string]interface{}{
		"id":                     discord.RandomSnowflake(),
		"name":                   testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"icon":                   testutil.RandomString(32),
		"description":            testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"rpc_origins":            testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 1, testutil.RandomIntInRange(10, 100)),
		"bot_public":             testutil.RandomBool(),
		"bot_require_code_grant": testutil.RandomBool(),
		"bot":                    NewUser(),
		"terms_of_service_url":   testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"privacy_policy_url":     testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"owner":                  NewUser(),
		"verify_key":             testutil.RandomString(32),
		"team": map[string]interface{}{
			"icon":          testutil.RandomString(32),
			"name":          testutil.RandomString(testutil.RandomIntInRange(1, 50)),
			"owner_user_id": discord.RandomSnowflake(),
			"id":            discord.RandomSnowflake(),
			"members": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
				*arrayToFill = append(*arrayToFill, map[string]interface{}{
					"membership_state": testutil.RandomItem(
						discord.ApplicationTeamMemberMembershipStateInvited,
						discord.ApplicationTeamMemberMembershipStateAccepted,
					),
					"team_id": discord.RandomSnowflake(),
					"user":    NewUser(),
					"role":    testutil.RandomItem("admin", "developer", "read_only"),
				})
			}),
		},
		"guild_id":       discord.RandomSnowflake(),
		"guild":          NewAvailableGuild(),
		"primary_sku_id": discord.RandomSnowflake(),
		"slug":           testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"cover_image":    testutil.RandomString(32),
		"flags": testutil.RandomFlags(
			discord.ApplicationFlagApplicationAutoModerationRuleCreateBadge,
			discord.ApplicationFlagGatewayPresence,
			discord.ApplicationFlagGatewayPresenceLimited,
			discord.ApplicationFlagGatewayGuildMembers,
			discord.ApplicationFlagGatewayGuildMembersLimited,
			discord.ApplicationFlagVerificationPendingGuildLimit,
			discord.ApplicationFlagEmbedded,
			discord.ApplicationFlagGatewayMessageContent,
			discord.ApplicationFlagGatewayMessageContentLimited,
			discord.ApplicationFlagApplicationCommandBadge,
		),
		"approximate_guild_count":              testutil.RandomIntInRange(1, 1000000),
		"approximate_user_authorization_count": testutil.RandomIntInRange(1, 1000000),
		"redirect_uris":                        testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 1, testutil.RandomIntInRange(10, 100)),
		"interactions_endpoint_url":            testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"role_connections_verification_url":    testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"event_webhooks_url":                   testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"event_webhooks_status": testutil.RandomItem(
			discord.ApplicationEventWebhookStatusDisabled,
			discord.ApplicationEventWebhookStatusEnabled,
			discord.ApplicationEventWebhookStatusDisabledByDiscord,
		),
		"event_webhooks_types": testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 1, testutil.RandomIntInRange(10, 100)),
		"tags":                 testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 1, testutil.RandomIntInRange(10, 100)),
		"install_params": map[string]interface{}{
			"scopes": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]discord.Scope) {
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
			"permissions": testutil.RandomFlags(testutil.AllPermissions...),
		},
		"integration_types_config": map[discord.ApplicationIntegrationType]discord.Snowflake{
			discord.ApplicationIntegrationTypeGuildInstall: discord.RandomSnowflake(),
			discord.ApplicationIntegrationTypeUserInstall:  discord.RandomSnowflake(),
		},
		"custom_install_url": testutil.RandomString(testutil.RandomIntInRange(1, 50)),
	}
}
