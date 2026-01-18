package types

import "encoding/json"

type DiscordComponentType int

type Components []AnyComponent

const (
	DiscordComponentTypeActionRow        DiscordComponentType = 1
	DiscordComponentTypeButton           DiscordComponentType = 2
	DiscordComponentTypeStringSelectMenu DiscordComponentType = 3
	DiscordComponentTypeTextInput        DiscordComponentType = 4
	DiscordComponentTypeUserSelectMenu   DiscordComponentType = 5
	DiscordComponentTypeRoleSelectMenu   DiscordComponentType = 6
	DiscordComponentTypeMentionableMenu  DiscordComponentType = 7
	DiscordComponentTypeChannelSelect    DiscordComponentType = 8
	DiscordComponentTypeSection          DiscordComponentType = 9
	//TODO
)

func (c DiscordComponentType) IsAnySelectMenu() bool {
	return c == DiscordComponentTypeStringSelectMenu ||
		c == DiscordComponentTypeUserSelectMenu ||
		c == DiscordComponentTypeRoleSelectMenu ||
		c == DiscordComponentTypeMentionableMenu ||
		c == DiscordComponentTypeChannelSelect
}

type AnyComponent interface {
	GetType() DiscordComponentType
}

type ComponentBase struct {
	Type DiscordComponentType `json:"type"`
}

func (b ComponentBase) GetType() DiscordComponentType {
	return b.Type
}

func (c *Components) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, r := range raw {
		var base ComponentBase
		if err := json.Unmarshal(r, &base); err != nil {
			return err
		}

		var comp AnyComponent

		switch base.Type {
		case 2:
			var btn ButtonComponent
			if err := json.Unmarshal(r, &btn); err != nil {
				return err
			}
			comp = btn
		default:
			comp = base
		}

		*c = append(*c, comp)
	}

	return nil
}

type ActionRow struct {
	ComponentBase
	ID         *int       `json:"id"`
	Components Components `json:"components"`
}

type ButtonStyle int

const (
	ButtonStylePrimary   ButtonStyle = 1
	ButtonStyleSecondary ButtonStyle = 2
	ButtonStyleSuccess   ButtonStyle = 3
	ButtonStyleDanger    ButtonStyle = 4
	ButtonStyleLink      ButtonStyle = 5
	ButtonStylePremium   ButtonStyle = 6
)

type ButtonComponent struct {
	Type     DiscordComponentType `json:"type"`
	ID       *int                 `json:"id,omitempty"`
	Style    ButtonStyle          `json:"style"`
	Label    *string              `json:"label,omitempty"`
	Emoji    *DiscordEmoji        `json:"emoji,omitempty"`
	CustomID *string              `json:"custom_id,omitempty"`
	SkuID    *DiscordSnowflake    `json:"sku_id,omitempty"`
	URL      *string              `json:"url,omitempty"`
	Disabled *bool                `json:"disabled,omitempty"`
}

func (b ButtonComponent) GetType() DiscordComponentType {
	return b.Type
}

type DiscordApplicationCommandInteractionOptionType int

const (
	DiscordApplicationCommandInteractionOptionTypeSubCommand      DiscordApplicationCommandInteractionOptionType = 1
	DiscordApplicationCommandInteractionOptionTypeSubCommandGroup DiscordApplicationCommandInteractionOptionType = 2
	DiscordApplicationCommandInteractionOptionTypeString          DiscordApplicationCommandInteractionOptionType = 3
	DiscordApplicationCommandInteractionOptionTypeInteger         DiscordApplicationCommandInteractionOptionType = 4
	DiscordApplicationCommandInteractionOptionTypeBoolean         DiscordApplicationCommandInteractionOptionType = 5
	DiscordApplicationCommandInteractionOptionTypeUser            DiscordApplicationCommandInteractionOptionType = 6
	DiscordApplicationCommandInteractionOptionTypeChannel         DiscordApplicationCommandInteractionOptionType = 7
	DiscordApplicationCommandInteractionOptionTypeRole            DiscordApplicationCommandInteractionOptionType = 8
	DiscordApplicationCommandInteractionOptionTypeMentionable     DiscordApplicationCommandInteractionOptionType = 9
	DiscordApplicationCommandInteractionOptionTypeNumber          DiscordApplicationCommandInteractionOptionType = 10
	DiscordApplicationCommandInteractionOptionTypeAttachment      DiscordApplicationCommandInteractionOptionType = 11
)

type DiscordSelectOptionValue struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description *string       `json:"description,omitempty"`
	Emoji       *DiscordEmoji `json:"emoji,omitempty"`
	Default     *bool         `json:"default,omitempty"`
}
