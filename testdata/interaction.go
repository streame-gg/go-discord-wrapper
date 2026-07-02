package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// In reality, some of the structs in NewInteraction are only partial, but it is not really specified what fields of those
// are omitted and which not, so we are using the full objects here
func NewInteraction() map[string]interface{} {
	return map[string]interface{}{
		"id":             discord.RandomSnowflake(),
		"application_id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.InteractionTypePing,
			discord.InteractionTypeApplicationCommand,
			discord.InteractionTypeMessageComponent,
			discord.InteractionTypeApplicationCommandAutocomplete,
			discord.InteractionTypeModalSubmit,
		),
		"data": testutil.RandomItem(
			NewInteractionDataApplicationCommand(),
			NewInteractionDataMessageComponent(),
			NewModalSubmitData(),
		),
		"guild":           NewAvailableGuild(),
		"guild_id":        discord.RandomSnowflake(),
		"channel":         NewChannel(),
		"channel_id":      discord.RandomSnowflake(),
		"member":          NewGuildMember(),
		"user":            NewUser(),
		"token":           testutil.RandomString(32),
		"version":         "1",
		"message":         nil,
		"app_permissions": testutil.RandomFlags(testutil.AllPermissions...),
		"locale":          testutil.RandomString(2),
		"guild_locale":    testutil.RandomString(2),
		"entitlements": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 5), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewEntitlement())
		}),
		"authorizing_integration_owners": map[discord.ApplicationIntegrationType]discord.Snowflake{
			discord.ApplicationIntegrationTypeGuildInstall: discord.RandomSnowflake(),
			discord.ApplicationIntegrationTypeUserInstall:  discord.RandomSnowflake(),
		},
		"context": testutil.RandomItem(
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypeBotDM,
			discord.InteractionContextTypePrivateChannel,
		),
		"attachment_size_limit": testutil.RandomIntInRange(1, 1000000),
	}
}

func NewModalSubmitData() map[string]interface{} {
	return map[string]interface{}{
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		//TODO
		"components": nil,
		"resolved":   NewResolvedData(),
	}
}

func NewResolvedData() map[string]interface{} {
	return map[string]interface{}{
		"users": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewUser())
		}),
		"members": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewGuildMember())
		}),
		"roles": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewRole())
		}),
		"channels": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewChannel())
		}),
		//TODO
		"messages": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{})
		}),
		//TODO
		"attachments": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{})
		}),
	}
}

func NewInteractionDataMessageComponent() map[string]interface{} {
	return map[string]interface{}{
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"component_type": testutil.RandomItem(
			discord.ComponentTypeActionRow,
			discord.ComponentTypeButton,
			discord.ComponentTypeStringSelect,
			discord.ComponentTypeTextInput,
			discord.ComponentTypeUserSelect,
			discord.ComponentTypeRoleSelect,
			discord.ComponentTypeMentionableSelect,
			discord.ComponentTypeChannelSelect,
			discord.ComponentTypeSection,
			discord.ComponentTypeTextDisplay,
			discord.ComponentTypeThumbnail,
			discord.ComponentTypeMediaGallery,
			discord.ComponentTypeFileDisplay,
			discord.ComponentTypeSeparator,
			discord.ComponentTypeContainer,
			discord.ComponentTypeLabel,
			discord.ComponentTypeFileUpload,
			discord.ComponentTypeRadioGroup,
			discord.ComponentTypeCheckboxGroup,
			discord.ComponentTypeCheckbox,
		),
		"resolved": NewResolvedData(),
		"values": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 15), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewSelectOption())
		}),
	}
}

func NewSelectOption() map[string]interface{} {
	return map[string]interface{}{
		"label":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"value":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"emoji": map[string]interface{}{
			"id":       discord.RandomSnowflake(),
			"name":     testutil.RandomString(testutil.RandomIntInRange(1, 100)),
			"animated": testutil.RandomBool(),
		},
		"default": testutil.RandomBool(),
	}
}

func NewInteractionDataApplicationCommand() map[string]interface{} {
	return map[string]interface{}{
		"id":   discord.RandomSnowflake(),
		"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"type": testutil.RandomItem(
			discord.ApplicationCommandTypeChatInput,
			discord.ApplicationCommandTypeUser,
			discord.ApplicationCommandTypeMessage,
			discord.ApplicationCommandTypePrimaryEndpoint,
		),
		"guild_id": discord.RandomSnowflake(),
		"options": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 32), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewApplicationCommandInteractionDataOptionStructure(true))
		}),
		"target_id": discord.RandomSnowflake(),
		"resolved":  NewResolvedData(),
	}
}

func NewApplicationCommandInteractionDataOptionStructure(shouldHaveOptions bool) map[string]interface{} {
	obj := map[string]interface{}{
		"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"type": testutil.RandomItem(
			discord.ApplicationCommandOptionTypeSubCommand,
			discord.ApplicationCommandOptionTypeSubCommandGroup,
			discord.ApplicationCommandOptionTypeString,
			discord.ApplicationCommandOptionTypeInteger,
			discord.ApplicationCommandOptionTypeBoolean,
			discord.ApplicationCommandOptionTypeUser,
			discord.ApplicationCommandOptionTypeChannel,
			discord.ApplicationCommandOptionTypeRole,
			discord.ApplicationCommandOptionTypeMentionable,
			discord.ApplicationCommandOptionTypeNumber,
			discord.ApplicationCommandOptionTypeAttachment,
		),
		"value":   testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"focused": testutil.RandomBool(),
	}

	if shouldHaveOptions {
		obj["options"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 32), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewApplicationCommandInteractionDataOptionStructure(false))
		})
	}

	return obj
}

func NewApplicationCommandOption(addChoices bool) map[string]interface{} {
	obj := map[string]interface{}{
		"type": testutil.RandomItem(
			discord.ApplicationCommandOptionTypeSubCommand,
			discord.ApplicationCommandOptionTypeSubCommandGroup,
			discord.ApplicationCommandOptionTypeString,
			discord.ApplicationCommandOptionTypeInteger,
			discord.ApplicationCommandOptionTypeBoolean,
			discord.ApplicationCommandOptionTypeUser,
			discord.ApplicationCommandOptionTypeChannel,
			discord.ApplicationCommandOptionTypeRole,
			discord.ApplicationCommandOptionTypeMentionable,
			discord.ApplicationCommandOptionTypeNumber,
			discord.ApplicationCommandOptionTypeAttachment,
		),
		"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"name_localizations": map[string]interface{}{
			string(discord.LocaleGerman): testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			string(discord.LocaleDanish): testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		},
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"description_localizations": map[string]interface{}{
			string(discord.LocaleGerman): testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			string(discord.LocaleDanish): testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		},
		"required": testutil.RandomBool(),
		"choices": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 2), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{
				"name": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
				"name_localizations": map[string]interface{}{
					string(discord.LocaleGerman): testutil.RandomString(testutil.RandomIntInRange(1, 32)),
					string(discord.LocaleDanish): testutil.RandomString(testutil.RandomIntInRange(1, 32)),
				},
				"value": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
			})
		}),
		"channel_types": testutil.RandomItem(
			discord.ChannelTypeGuildText,
			discord.ChannelTypeDM,
			discord.ChannelTypeGuildVoice,
			discord.ChannelTypeGroupDM,
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
		"min_value":    testutil.RandomIntInRange(1, 100000),
		"max_value":    testutil.RandomIntInRange(100000, 100000000),
		"min_length":   testutil.RandomIntInRange(1, 100000),
		"max_length":   testutil.RandomIntInRange(100000, 100000000),
		"autocomplete": testutil.RandomBool(),
	}

	if addChoices {
		obj["options"] = testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewApplicationCommandOption(false))
		})
	}

	return obj
}
