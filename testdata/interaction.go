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
		"message":         NewMessage(true),
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

func NewApplicationCommandInteractionMetadata() map[string]interface{} {
	return map[string]interface{}{
		"id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.InteractionTypePing,
			discord.InteractionTypeApplicationCommand,
			discord.InteractionTypeMessageComponent,
			discord.InteractionTypeApplicationCommandAutocomplete,
			discord.InteractionTypeModalSubmit,
		),
		"user": NewUser(),
		"authorizing_integration_owners": map[discord.ApplicationIntegrationType]discord.Snowflake{
			discord.ApplicationIntegrationTypeGuildInstall: discord.RandomSnowflake(),
			discord.ApplicationIntegrationTypeUserInstall:  discord.RandomSnowflake(),
		},
		"original_response_message_id": discord.RandomSnowflake(),
		"target_user":                  NewUser(),
		"target_message_id":            discord.RandomSnowflake(),
	}
}

func NewMessageComponentInteractionMetadata() map[string]interface{} {
	return map[string]interface{}{
		"id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.InteractionTypePing,
			discord.InteractionTypeApplicationCommand,
			discord.InteractionTypeMessageComponent,
			discord.InteractionTypeApplicationCommandAutocomplete,
			discord.InteractionTypeModalSubmit,
		),
		"user": NewUser(),
		"authorizing_integration_owners": map[discord.ApplicationIntegrationType]discord.Snowflake{
			discord.ApplicationIntegrationTypeGuildInstall: discord.RandomSnowflake(),
			discord.ApplicationIntegrationTypeUserInstall:  discord.RandomSnowflake(),
		},
		"original_response_message_id": discord.RandomSnowflake(),
		"interacted_message_id":        discord.RandomSnowflake(),
	}
}

func NewModalSubmitInteractionMetadata() map[string]interface{} {
	return map[string]interface{}{
		"id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.InteractionTypePing,
			discord.InteractionTypeApplicationCommand,
			discord.InteractionTypeMessageComponent,
			discord.InteractionTypeApplicationCommandAutocomplete,
			discord.InteractionTypeModalSubmit,
		),
		"user": NewUser(),
		"authorizing_integration_owners": map[discord.ApplicationIntegrationType]discord.Snowflake{
			discord.ApplicationIntegrationTypeGuildInstall: discord.RandomSnowflake(),
			discord.ApplicationIntegrationTypeUserInstall:  discord.RandomSnowflake(),
		},
		"original_response_message_id": discord.RandomSnowflake(),
		"triggering_interaction_metadata": testutil.RandomItem(
			NewApplicationCommandInteractionMetadata(),
			NewMessageComponentInteractionMetadata(),
		),
	}
}

func NewModalSubmitData() map[string]interface{} {
	return map[string]interface{}{
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"components": testutil.RandomItem(
			NewStringSelectMenuData(),
			NewTextInputData(),
			NewUserSelectMenuData(),
			NewRoleSelectMenuData(),
			NewMentionableSelectMenuData(),
			NewChannelSelectMenuData(),
			NewTextDisplayData(),
			NewLabelData(),
			NewFileUploadData(),
			NewRadioGroupData(),
			NewCheckboxGroupData(),
			NewCheckboxData(),
		),
		"resolved": NewResolvedData(),
	}
}

func NewStringSelectMenuData() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeStringSelect,
		"component_type": discord.ComponentTypeStringSelect,
		"id":             testutil.RandomIntInRange(1, 100),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"values":         testutil.RandomStringArray(testutil.RandomIntInRange(1, 25), testutil.RandomIntInRange(1, 32), testutil.RandomIntInRange(1, 32)),
	}
}

func NewTextInputData() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeTextInput,
		"id":        testutil.RandomIntInRange(1, 100),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"value":     testutil.RandomString(testutil.RandomIntInRange(1, 32)),
	}
}

func NewUserSelectMenuData() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeUserSelect,
		"component_type": discord.ComponentTypeUserSelect,
		"id":             testutil.RandomIntInRange(1, 100),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"resolved":       NewResolvedData(),
		"values":         testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewRoleSelectMenuData() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeRoleSelect,
		"component_type": discord.ComponentTypeRoleSelect,
		"id":             testutil.RandomIntInRange(1, 100),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"resolved":       NewResolvedData(),
		"values":         testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewMentionableSelectMenuData() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeMentionableSelect,
		"component_type": discord.ComponentTypeMentionableSelect,
		"id":             testutil.RandomIntInRange(1, 100),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"resolved":       NewResolvedData(),
		"values":         testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewChannelSelectMenuData() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeChannelSelect,
		"component_type": discord.ComponentTypeChannelSelect,
		"id":             testutil.RandomIntInRange(1, 100),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"resolved":       NewResolvedData(),
		"values":         testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewTextDisplayData() map[string]interface{} {
	return map[string]interface{}{
		"type": discord.ComponentTypeTextDisplay,
		"id":   testutil.RandomIntInRange(1, 100),
	}
}

func NewLabelData() map[string]interface{} {
	return map[string]interface{}{
		"type": discord.ComponentTypeLabel,
		"id":   testutil.RandomIntInRange(1, 100),
		"component": testutil.RandomItem(
			NewStringSelectMenuData(),
			NewTextInputData(),
			NewUserSelectMenuData(),
			NewRoleSelectMenuData(),
			NewMentionableSelectMenuData(),
			NewChannelSelectMenuData(),
			NewFileUploadData(),
			NewRadioGroupData(),
			NewCheckboxGroupData(),
			NewCheckboxData(),
		),
	}
}

func NewFileUploadData() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeFileUpload,
		"id":        testutil.RandomIntInRange(1, 100),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"values":    testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewRadioGroupData() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeRadioGroup,
		"id":        testutil.RandomIntInRange(1, 100),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"value":     testutil.RandomString(testutil.RandomIntInRange(1, 32)),
	}
}

func NewCheckboxGroupData() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeCheckboxGroup,
		"id":        testutil.RandomIntInRange(1, 100),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"values":    testutil.RandomStringArray(testutil.RandomIntInRange(1, 25), testutil.RandomIntInRange(1, 32), testutil.RandomIntInRange(1, 32)),
	}
}

func NewCheckboxData() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeCheckbox,
		"id":        testutil.RandomIntInRange(1, 100),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"value":     testutil.RandomString(testutil.RandomIntInRange(1, 25)),
	}
}

func NewResolvedData() map[string]interface{} {
	return map[string]interface{}{
		"users": map[discord.Snowflake]interface{}{
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
			discord.RandomSnowflake(): NewUser(),
		},
		"members": map[discord.Snowflake]interface{}{
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
			discord.RandomSnowflake(): NewGuildMember(),
		},
		"roles": map[discord.Snowflake]interface{}{
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
			discord.RandomSnowflake(): NewRole(),
		},
		"channels": map[discord.Snowflake]interface{}{
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
			discord.RandomSnowflake(): NewChannel(),
		},
		"messages": map[discord.Snowflake]interface{}{
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
			discord.RandomSnowflake(): NewMessage(true),
		},
		"attachments": map[discord.Snowflake]interface{}{
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
			discord.RandomSnowflake(): NewAttachment(),
		},
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
