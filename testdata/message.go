package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewMessage(allowMessageSnapshots bool) map[string]interface{} {
	obj := map[string]interface{}{
		"id":               discord.RandomSnowflake(),
		"channel_id":       discord.RandomSnowflake(),
		"author":           NewUser(),
		"content":          testutil.RandomString(testutil.RandomIntInRange(1, 4000)),
		"timestamp":        testutil.RandomTime(),
		"edited_timestamp": discord.RandomSnowflake(),
		"tts":              testutil.RandomBool(),
		"mention_everyone": testutil.RandomBool(),
		"mentions": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewUser())
		}),
		"mention_roles": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 50)),
		"mention_channels": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewChannelMention())
		}),
		"attachments": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewAttachment())
		}),
		"embeds": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewEmbed())
		}),
		"reactions": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewReaction())
		}),
		"nonce":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"pinned":     testutil.RandomBool(),
		"webhook_id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.MessageTypeDefault,
			discord.MessageTypeRecipientAdd,
			discord.MessageTypeRecipientRemove,
			discord.MessageTypeCall,
			discord.MessageTypeChannelNameChange,
			discord.MessageTypeChannelIconChange,
			discord.MessageTypeChannelPinnedMessage,
			discord.MessageTypeGuildUserJoin,
			discord.MessageTypeGuildBoost,
			discord.MessageTypeGuildBoostTier1,
			discord.MessageTypeGuildBoostTier2,
			discord.MessageTypeGuildBoostTier3,
			discord.MessageTypeChannelFollowAdd,
			discord.MessageTypeGuildDiscoveryDisqualified,
			discord.MessageTypeGuildDiscoveryRequalified,
			discord.MessageTypeGuildDiscoveryGracePeriodInitialWarning,
			discord.MessageTypeGuildDiscoveryGracePeriodFinalWarning,
			discord.MessageTypeThreadCreated,
			discord.MessageTypeReply,
			discord.MessageTypeChatInputCommand,
			discord.MessageTypeThreadStarterMessage,
			discord.MessageTypeGuildInviteReminder,
			discord.MessageTypeContextMenuCommand,
			discord.MessageTypeAutoModerationAction,
			discord.MessageTypeRoleSubscriptionPurchase,
			discord.MessageTypeInteractionPremiumUpsell,
			discord.MessageTypeStageStart,
			discord.MessageTypeStageEnd,
			discord.MessageTypeStageSpeaker,
			discord.MessageTypeStageTopic,
			discord.MessageTypeGuildApplicationPremiumSubscription,
			discord.MessageTypeGuildIncidentAlertModeEnabled,
			discord.MessageTypeGuildIncidentAlertModeDisabled,
			discord.MessageTypeReportRaid,
			discord.MessageTypeReportFalseAlarm,
			discord.MessageTypePurchaseNotification,
			discord.MessageTypePollResult,
		),
		"activity": map[string]interface{}{
			"type": testutil.RandomItem(
				discord.MessageActivityTypeJoin,
				discord.MessageActivityTypeSpectate,
				discord.MessageActivityTypeListen,
				discord.MessageActivityTypeJoinRequest,
			),
			"party_id": testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		},
		"application":    nil,
		"application_id": discord.RandomSnowflake(),
		"flags": testutil.RandomFlags(
			discord.MessageFlagCrossposted,
			discord.MessageFlagIsCrosspost,
			discord.MessageFlagSuppressEmbeds,
			discord.MessageFlagSourceMessageDeleted,
			discord.MessageFlagUrgent,
			discord.MessageFlagHasThread,
			discord.MessageFlagEphemeral,
			discord.MessageFlagLoading,
			discord.MessageFlagFailedToMentionSomeRolesInThread,
			discord.MessageFlagSuppressNotification,
			discord.MessageFlagIsVoiceMessage,
			discord.MessageFlagHasSnapshot,
			discord.MessageFlagIsComponentsV2,
		),
		"message_reference": map[string]interface{}{
			"type": testutil.RandomItem(
				discord.MessageReferenceTypeDefault,
				discord.MessageReferenceTypeForward,
			),
			"message_id":         discord.RandomSnowflake(),
			"channel_id":         discord.RandomSnowflake(),
			"guild_id":           discord.RandomSnowflake(),
			"fail_if_not_exists": testutil.RandomBool(),
		},
		"referenced_message": NewMessage(true),
		"interaction_metadata": testutil.RandomItem(
			NewApplicationCommandInteractionMetadata(),
			NewMessageComponentInteractionMetadata(),
			NewModalSubmitInteractionMetadata(),
		),
		"interaction": map[string]interface{}{
			"id": discord.RandomSnowflake(),
			"type": testutil.RandomItem(
				discord.InteractionTypePing,
				discord.InteractionTypeApplicationCommand,
				discord.InteractionTypeMessageComponent,
				discord.InteractionTypeApplicationCommandAutocomplete,
				discord.InteractionTypeModalSubmit,
			),
			"name":   testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			"user":   NewUser(),
			"member": NewGuildMember(),
		},
		"thread": NewChannel(),
		"components": testutil.RandomItem(
			NewComponentActionRow(),
			NewComponentButton(),
			NewComponentStringSelect(),
			NewComponentTextInput(),
			NewComponentUserSelect(),
			NewComponentRoleSelect(),
			NewComponentMentionableSelect(),
			NewComponentChannelSelect(),
			NewComponentSection(),
			NewComponentTextDisplay(),
			NewComponentThumbnail(),
			NewComponentMediaGallery(),
			NewComponentFile(),
			NewComponentSeparator(),
			NewComponentContainer(),
			NewComponentLabel(nil),
			NewComponentFileUpload(),
			NewComponentRadioGroup(),
			NewComponentCheckboxGroup(),
			NewComponentCheckbox(),
		),
		"sticker_items": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewStickerItem())
		}),
		"stickers": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewSticker())
		}),
		"position": testutil.RandomIntInRange(1, 1000),
		"role_subscription_data": map[string]interface{}{
			"role_subscription_listing_id": discord.RandomSnowflake(),
			"tier_name":                    testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			"total_months_subscribed":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			"is_renewal":                   testutil.RandomBool(),
		},
		"resolved": NewResolvedData(),
		"poll":     NewPoll(),
		"call": map[string]interface{}{
			"participants":    testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 10)),
			"ended_timestamp": testutil.RandomTime(),
		},
		"shared_client_theme": map[string]interface{}{
			"colors":         testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 7, 7),
			"gradient_angle": testutil.RandomIntInRange(1, 360),
			"base_mix":       testutil.RandomIntInRange(1, 100),
			"base_theme": testutil.RandomItem(
				discord.BaseThemeUnset,
				discord.BaseThemeDark,
				discord.BaseThemeLight,
				discord.BaseThemeDarker,
				discord.BaseThemeMidnight,
			),
		},
	}

	if allowMessageSnapshots {
		obj["message_snapshots"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewMessageSnapshot())
		})
	}

	return obj
}

func NewStickerItem() map[string]interface{} {
	return map[string]interface{}{
		"id":   discord.RandomSnowflake(),
		"name": testutil.RandomString(testutil.RandomIntInRange(2, 30)),
		"format_type": testutil.RandomItem(
			discord.StickerFormatTypePNG,
			discord.StickerFormatTypeAPNG,
			discord.StickerFormatTypeLottie,
			discord.StickerFormatTypeGIF,
		),
	}
}

func NewPoll() map[string]interface{} {
	return map[string]interface{}{
		"question": NewPollMediaItem(false),
		"answers": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{
				"answer_id":  testutil.RandomIntInRange(1, 10),
				"poll_media": NewPollMediaItem(testutil.RandomBool()),
			})
		}),
		"expiry":            testutil.RandomTime(),
		"allow_multiselect": testutil.RandomBool(),
		"layout_type": testutil.RandomItem(
			discord.PollLayoutTypeDefault,
		),
		"results": map[string]interface{}{
			"is_finalized": testutil.RandomBool(),
			"answer_counts": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
				*arrayToFill = append(*arrayToFill, map[string]interface{}{
					"id":       testutil.RandomIntInRange(1, 10),
					"count":    testutil.RandomIntInRange(1, 1000),
					"me_voted": testutil.RandomBool(),
				})
			}),
		},
	}
}

func NewPollMediaItem(hasEmoji bool) map[string]interface{} {
	obj := map[string]interface{}{
		"text": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
	}

	if hasEmoji {
		obj["emoji"] = NewEmoji()
	}

	return obj
}

func NewMessageSnapshot() map[string]interface{} {
	return map[string]interface{}{
		"message": NewMessage(false),
	}
}

func NewReaction() map[string]interface{} {
	return map[string]interface{}{
		"count": testutil.RandomIntInRange(1, 5000),
		"count_details": map[string]interface{}{
			"burst":  testutil.RandomIntInRange(1, 5000),
			"normal": testutil.RandomIntInRange(1, 5000),
		},
		"me":           testutil.RandomBool(),
		"me_burst":     testutil.RandomBool(),
		"emoji":        NewEmoji(),
		"burst_colors": testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 7, 7),
	}
}

func NewEmbed() map[string]interface{} {
	return map[string]interface{}{
		"title":       testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		"type":        discord.EmbedTypeRich,
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 4096)),
		"url":         testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		"timestamp":   testutil.RandomTime(),
		"color":       testutil.RandomIntInRange(0x000000, 0xFFFFFF),
		"footer": map[string]interface{}{
			"text":           testutil.RandomString(testutil.RandomIntInRange(1, 256)),
			"icon_url":       testutil.RandomString(testutil.RandomIntInRange(1, 256)),
			"proxy_icon_url": testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		},
		"image":     NewEmbedImage(),
		"thumbnail": NewEmbedImage(),
		"video":     NewEmbedImage(),
		"provider": map[string]interface{}{
			"url":  testutil.RandomString(testutil.RandomIntInRange(1, 256)),
			"name": testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		},
		"author": map[string]interface{}{
			"name":           testutil.RandomString(testutil.RandomIntInRange(1, 256)),
			"icon_url":       testutil.RandomString(testutil.RandomIntInRange(1, 256)),
			"url":            testutil.RandomString(testutil.RandomIntInRange(1, 256)),
			"proxy_icon_url": testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		},
		"fields": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewEmbedField())
		}),
		"flags": testutil.RandomFlags(discord.EmbedFlagIsContentInventoryEntry),
	}
}

func NewEmbedField() map[string]interface{} {
	return map[string]interface{}{
		"name":   testutil.RandomString(testutil.RandomIntInRange(1, 1024)),
		"value":  testutil.RandomString(testutil.RandomIntInRange(1, 1024)),
		"inline": testutil.RandomBool(),
	}
}

func NewEmbedImage() map[string]interface{} {
	return map[string]interface{}{
		"url":                 testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		"proxy_url":           testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		"height":              testutil.RandomIntInRange(1, 4096),
		"width":               testutil.RandomIntInRange(1, 4096),
		"content_type":        testutil.RandomString(testutil.RandomIntInRange(1, 16)),
		"placeholder":         testutil.RandomString(testutil.RandomIntInRange(1, 16)),
		"placeholder_version": testutil.RandomIntInRange(1, 100),
		"description":         testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"flags":               testutil.RandomFlags(discord.EmbedMediaFlagIsAnimated),
	}
}

func NewAttachment() map[string]interface{} {
	return map[string]interface{}{
		"id":                  discord.RandomSnowflake(),
		"name":                testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"title":               testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"description":         testutil.RandomString(testutil.RandomIntInRange(1, 1024)),
		"content_type":        testutil.RandomString(testutil.RandomIntInRange(1, 5)),
		"size":                testutil.RandomIntInRange(1, 1000000),
		"url":                 testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"proxy_url":           testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"height":              testutil.RandomIntInRange(1, 4096),
		"width":               testutil.RandomIntInRange(1, 4096),
		"placeholder":         testutil.RandomString(testutil.RandomIntInRange(1, 50)),
		"placeholder_version": testutil.RandomIntInRange(1, 50),
		"ephemeral":           testutil.RandomBool(),
		"duration_secs":       testutil.RandomFloat64InRange(1, 300),
		"waveform":            testutil.RandomString(testutil.RandomIntInRange(1, 1000)),
		"flags": testutil.RandomFlags(
			discord.AttachmentFlagIsClip,
			discord.AttachmentFlagIsThumbnail,
			discord.AttachmentFlagIsRemix,
			discord.AttachmentFlagIsSpoiler,
			discord.AttachmentFlagIsAnimated,
		),
		"clip_participants": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewUser())
		}),
		"clip_created_at": testutil.RandomTime(),
		"application":     nil,
	}
}

func NewChannelMention() map[string]interface{} {
	return map[string]interface{}{
		"id":       discord.RandomSnowflake(),
		"guild_id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.ChannelTypeGuildText,
			discord.ChannelTypeDM,
			discord.ChannelTypeGuildVoice,
			discord.ChannelTypeGuildCategory,
			discord.ChannelTypeGuildAnnouncement,
			discord.ChannelTypeAnnouncementThread,
			discord.ChannelTypePublicThread,
			discord.ChannelTypePrivateThread,
			discord.ChannelTypeGuildStageVoice,
			discord.ChannelTypeGuildDirectory,
			discord.ChannelTypeGuildForum,
			discord.ChannelTypeGuildMedia,
		),
		"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
	}
}
