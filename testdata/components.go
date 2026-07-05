package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewComponentActionRow(predefinedComponents ...map[string]interface{}) map[string]interface{} {
	if len(predefinedComponents) == 0 {
		return map[string]interface{}{
			"type": discord.ComponentTypeActionRow,
			"id":   testutil.RandomIntInRange(1, 25),
			"components": testutil.RandomItem(
				NewComponentButton(),
				NewComponentStringSelect(),
				NewComponentUserSelect(),
				NewComponentRoleSelect(),
				NewComponentMentionableSelect(),
				NewComponentChannelSelect(),
			),
		}
	}

	return map[string]interface{}{
		"type":       discord.ComponentTypeActionRow,
		"id":         testutil.RandomIntInRange(1, 25),
		"components": predefinedComponents,
	}
}

func NewComponentButton() map[string]interface{} {
	return map[string]interface{}{
		"type": discord.ComponentTypeButton,
		"id":   testutil.RandomIntInRange(1, 25),
		"style": testutil.RandomItem(
			components.ButtonStylePrimary,
			components.ButtonStyleSecondary,
			components.ButtonStyleSuccess,
			components.ButtonStyleDanger,
			components.ButtonStyleLink,
			components.ButtonStylePremium,
		),
		"label":     testutil.RandomString(testutil.RandomIntInRange(1, 80)),
		"emoji":     NewEmoji(),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"sku_id":    discord.RandomSnowflake(),
		"url":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"disabled":  testutil.RandomBool(),
	}
}

func NewComponentStringSelect() map[string]interface{} {
	return map[string]interface{}{
		"type":        discord.ComponentTypeStringSelect,
		"id":          testutil.RandomIntInRange(1, 25),
		"custom_id":   testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"placeholder": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"min_values":  testutil.RandomIntInRange(0, 100),
		"max_values":  testutil.RandomIntInRange(0, 100),
		"required":    testutil.RandomBool(),
		"disabled":    testutil.RandomBool(),
		"options": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{
				"label":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"value":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"emoji":       NewEmoji(),
				"default":     testutil.RandomBool(),
			})
		}),
	}
}

func NewComponentTextInput() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeTextInput,
		"id":        testutil.RandomIntInRange(1, 25),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"style": testutil.RandomItem(
			components.TextInputStyleShort,
			components.TextInputStyleParagraph,
		),
		"min_length":  testutil.RandomIntInRange(0, 100),
		"max_length":  testutil.RandomIntInRange(0, 100),
		"required":    testutil.RandomBool(),
		"value":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"placeholder": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
	}
}

func NewComponentUserSelect() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeUserSelect,
		"id":             testutil.RandomIntInRange(1, 25),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"placeholder":    testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"min_values":     testutil.RandomIntInRange(0, 100),
		"max_values":     testutil.RandomIntInRange(0, 100),
		"required":       testutil.RandomBool(),
		"disabled":       testutil.RandomBool(),
		"default_values": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewComponentRoleSelect() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeRoleSelect,
		"id":             testutil.RandomIntInRange(1, 25),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"placeholder":    testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"min_values":     testutil.RandomIntInRange(0, 100),
		"max_values":     testutil.RandomIntInRange(0, 100),
		"required":       testutil.RandomBool(),
		"disabled":       testutil.RandomBool(),
		"default_values": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewComponentMentionableSelect() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeMentionableSelect,
		"id":             testutil.RandomIntInRange(1, 25),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"placeholder":    testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"min_values":     testutil.RandomIntInRange(0, 100),
		"max_values":     testutil.RandomIntInRange(0, 100),
		"required":       testutil.RandomBool(),
		"disabled":       testutil.RandomBool(),
		"default_values": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
	}
}

func NewComponentChannelSelect() map[string]interface{} {
	return map[string]interface{}{
		"type":           discord.ComponentTypeChannelSelect,
		"id":             testutil.RandomIntInRange(1, 25),
		"custom_id":      testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"placeholder":    testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"min_values":     testutil.RandomIntInRange(0, 100),
		"max_values":     testutil.RandomIntInRange(0, 100),
		"required":       testutil.RandomBool(),
		"disabled":       testutil.RandomBool(),
		"default_values": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
		"channel_types": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]discord.ChannelType) {
			*arrayToFill = append(*arrayToFill, testutil.RandomItem(
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
			))
		}),
	}
}

func NewComponentSection() map[string]interface{} {
	return map[string]interface{}{
		"type": discord.ComponentTypeSection,
		"id":   testutil.RandomIntInRange(1, 25),
		"components": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewComponentTextDisplay())
		}),
		"accessory": testutil.RandomItem(
			NewComponentButton(),
			NewComponentThumbnail(),
		),
	}
}

func NewComponentTextDisplay() map[string]interface{} {
	return map[string]interface{}{
		"type":    discord.ComponentTypeTextDisplay,
		"id":      testutil.RandomIntInRange(1, 25),
		"content": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
	}
}

func NewComponentThumbnail() map[string]interface{} {
	return map[string]interface{}{
		"type":        discord.ComponentTypeThumbnail,
		"id":          testutil.RandomIntInRange(1, 25),
		"media":       NewUnfurledMediaItem(),
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"spoiler":     testutil.RandomBool(),
	}
}

func NewUnfurledMediaItem() map[string]interface{} {
	return map[string]interface{}{
		"url":                 testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		"proxy_url":           testutil.RandomString(testutil.RandomIntInRange(1, 256)),
		"height":              testutil.RandomIntInRange(1, 4096),
		"width":               testutil.RandomIntInRange(1, 4096),
		"content_type":        testutil.RandomString(testutil.RandomIntInRange(1, 16)),
		"placeholder":         testutil.RandomString(testutil.RandomIntInRange(1, 16)),
		"placeholder_version": testutil.RandomIntInRange(1, 100),
		"flags":               testutil.RandomFlags(components.UnfurledMediaItemFlagIsAnimated),
		"attachment_id":       discord.RandomSnowflake(),
	}
}

func NewComponentMediaGallery() map[string]interface{} {
	return map[string]interface{}{
		"type": discord.ComponentTypeMediaGallery,
		"id":   testutil.RandomIntInRange(1, 25),
		"items": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{
				"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"spoiler":     testutil.RandomBool(),
				"media":       NewUnfurledMediaItem(),
			})
		}),
	}
}

func NewComponentFile() map[string]interface{} {
	return map[string]interface{}{
		"type":    discord.ComponentTypeFileDisplay,
		"id":      testutil.RandomIntInRange(1, 25),
		"file":    NewUnfurledMediaItem(),
		"spoiler": testutil.RandomBool(),
		"name":    testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"size":    testutil.RandomIntInRange(1, 4096),
	}
}

func NewComponentSeparator() map[string]interface{} {
	return map[string]interface{}{
		"type":    discord.ComponentTypeSeparator,
		"id":      testutil.RandomIntInRange(1, 25),
		"divider": testutil.RandomBool(),
		"spacing": testutil.RandomItem(
			components.SeparatorComponentSpacingSmall,
			components.SeparatorComponentSpacingLarge,
		),
	}
}

func NewComponentContainer(predefinedComponents ...map[string]interface{}) map[string]interface{} {
	if len(predefinedComponents) == 0 {
		return map[string]interface{}{
			"type": discord.ComponentTypeContainer,
			"id":   testutil.RandomIntInRange(1, 25),
			"components": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
				*arrayToFill = append(*arrayToFill, testutil.RandomItem(
					NewComponentActionRow(),
					NewComponentTextDisplay(),
					NewComponentSection(),
					NewComponentMediaGallery(),
					NewComponentSeparator(),
					NewComponentFile(),
				))
			}),
			"accent_color": testutil.RandomIntInRange(0x000000, 0xFFFFFF),
			"spoiler":      testutil.RandomBool(),
		}
	}

	return map[string]interface{}{
		"type":         discord.ComponentTypeContainer,
		"id":           testutil.RandomIntInRange(1, 25),
		"components":   predefinedComponents,
		"accent_color": testutil.RandomIntInRange(0x000000, 0xFFFFFF),
		"spoiler":      testutil.RandomBool(),
	}
}

func NewComponentLabel(predefinedComponent *map[string]interface{}) map[string]interface{} {
	if predefinedComponent == nil {
		return map[string]interface{}{
			"type":        discord.ComponentTypeLabel,
			"id":          testutil.RandomIntInRange(1, 25),
			"label":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
			"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
			"component": testutil.RandomItem(
				NewComponentTextInput(),
				NewComponentStringSelect(),
				NewComponentUserSelect(),
				NewComponentRoleSelect(),
				NewComponentMentionableSelect(),
				NewComponentChannelSelect(),
				NewComponentFileUpload(),
				NewComponentRadioGroup(),
				NewComponentCheckboxGroup(),
				NewComponentCheckbox(),
			),
		}
	}

	return map[string]interface{}{
		"type":        discord.ComponentTypeLabel,
		"id":          testutil.RandomIntInRange(1, 25),
		"label":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"component":   *predefinedComponent,
	}
}

func NewComponentFileUpload() map[string]interface{} {
	return map[string]interface{}{
		"type":       discord.ComponentTypeFileUpload,
		"id":         testutil.RandomIntInRange(1, 25),
		"custom_id":  testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"min_values": testutil.RandomIntInRange(0, 100),
		"max_values": testutil.RandomIntInRange(0, 100),
		"required":   testutil.RandomBool(),
	}
}

func NewComponentRadioGroup() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeRadioGroup,
		"id":        testutil.RandomIntInRange(1, 25),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"required":  testutil.RandomBool(),
		"options": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{
				"label":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"value":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"default":     testutil.RandomBool(),
			})
		}),
	}
}

func NewComponentCheckboxGroup() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeCheckboxGroup,
		"id":        testutil.RandomIntInRange(1, 25),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"options": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 25), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, map[string]interface{}{
				"label":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"value":       testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
				"default":     testutil.RandomBool(),
			})
		}),
		"min_values": testutil.RandomIntInRange(0, 100),
		"max_values": testutil.RandomIntInRange(0, 100),
		"required":   testutil.RandomBool(),
	}
}

func NewComponentCheckbox() map[string]interface{} {
	return map[string]interface{}{
		"type":      discord.ComponentTypeCheckbox,
		"id":        testutil.RandomIntInRange(1, 25),
		"custom_id": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"default":   testutil.RandomBool(),
	}
}
