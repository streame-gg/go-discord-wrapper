package types

import (
	"encoding/json"
)

type DiscordComponentType int

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
	DiscordComponentTypeTextDisplay      DiscordComponentType = 10
	DiscordComponentTypeThumbnail        DiscordComponentType = 11
	DiscordComponentTypeMediaGallery     DiscordComponentType = 12
	DiscordComponentTypeFileDisplay      DiscordComponentType = 13
	DiscordComponentTypeSeparator        DiscordComponentType = 14
	DiscordComponentTypeContainer        DiscordComponentType = 17
	DiscordComponentTypeLabel            DiscordComponentType = 18
	DiscordComponentTypeFileUpload       DiscordComponentType = 19
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
	MarshalJSON() ([]byte, error)
}

type ActionRow struct {
	Type       DiscordComponentType `json:"type"`
	ID         *int                 `json:"id"`
	Components []AnyComponent       `json:"components"`
}

func (a *ActionRow) GetType() DiscordComponentType {
	return DiscordComponentTypeActionRow
}

func NewActionRow(id *int, components ...AnyComponent) *ActionRow {
	return &ActionRow{
		Type:       DiscordComponentTypeActionRow,
		Components: components,
		ID:         id,
	}
}

func (a *ActionRow) MarshalJSON() ([]byte, error) {
	a.Type = DiscordComponentTypeActionRow
	type Alias ActionRow
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}

func (a *ActionRow) IsAnyContainerComponent() bool {
	return true
}

func (a *ActionRow) UnmarshalJSON(data []byte) error {
	type Alias ActionRow

	var raw struct {
		Alias
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*a = ActionRow(raw.Alias)

	for _, c := range raw.Components {
		var probe struct {
			Type DiscordComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case DiscordComponentTypeButton:
			var b ButtonComponent
			if err := json.Unmarshal(c, &b); err != nil {
				return err
			}
			a.Components = append(a.Components, &b)
		}
	}

	return nil
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
	Label    string               `json:"label,omitempty"`
	Emoji    *DiscordEmoji        `json:"emoji,omitempty"`
	CustomID string               `json:"custom_id,omitempty"`
	SkuID    *DiscordSnowflake    `json:"sku_id,omitempty"`
	URL      string               `json:"url,omitempty"`
	Disabled bool                 `json:"disabled,omitempty"`
}

func (b ButtonComponent) MarshalJSON() ([]byte, error) {
	b.Type = DiscordComponentTypeButton
	type Alias ButtonComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(b),
	})
}

func (b ButtonComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeButton
}

func (b ButtonComponent) IsAnySectionAccessory() bool {
	return true
}

type StringSelectMenuComponent struct {
	Type        DiscordComponentType               `json:"type"`
	ID          *int                               `json:"id,omitempty"`
	CustomID    string                             `json:"custom_id"`
	Placeholder string                             `json:"placeholder,omitempty"`
	MinValues   *int                               `json:"min_values,omitempty"`
	MaxValues   *int                               `json:"max_values,omitempty"`
	Required    bool                               `json:"required,omitempty"`
	Options     *[]StringSelectMenuComponentOption `json:"options"`
	Disabled    bool                               `json:"disabled,omitempty"`
}

func (s StringSelectMenuComponent) MarshalJSON() ([]byte, error) {
	s.Type = DiscordComponentTypeStringSelectMenu
	type Alias StringSelectMenuComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(s),
	})
}

func (s StringSelectMenuComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeStringSelectMenu
}

type StringSelectMenuComponentOption struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description,omitempty"`
	Emoji       *DiscordEmoji `json:"emoji,omitempty"`
	Default     bool          `json:"default,omitempty"`
}

type UserSelectMenuComponent struct {
	Type          DiscordComponentType  `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (u UserSelectMenuComponent) MarshalJSON() ([]byte, error) {
	u.Type = DiscordComponentTypeUserSelectMenu
	type Alias UserSelectMenuComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(u),
	})
}

func (u UserSelectMenuComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeUserSelectMenu
}

type RoleSelectMenuComponent struct {
	Type          DiscordComponentType  `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (r RoleSelectMenuComponent) MarshalJSON() ([]byte, error) {
	r.Type = DiscordComponentTypeRoleSelectMenu
	type Alias UserSelectMenuComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(r),
	})
}

func (r RoleSelectMenuComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeRoleSelectMenu
}

type MentionableSelectMenuComponent struct {
	Type          DiscordComponentType  `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (m MentionableSelectMenuComponent) MarshalJSON() ([]byte, error) {
	m.Type = DiscordComponentTypeMentionableMenu
	type Alias MentionableSelectMenuComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(m),
	})
}

func (m MentionableSelectMenuComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeMentionableMenu
}

type ChannelSelectMenuComponent struct {
	Type          DiscordComponentType  `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (c ChannelSelectMenuComponent) RoleSelectMenuComponent() DiscordComponentType {
	return DiscordComponentTypeRoleSelectMenu
}

type SelectDefaultValueType string

const (
	SelectDefaultValueTypeUser    SelectDefaultValueType = "user"
	SelectDefaultValueTypeRole    SelectDefaultValueType = "role"
	SelectDefaultValueTypeChannel SelectDefaultValueType = "channel"
)

type Section struct {
	Type       DiscordComponentType   `json:"type"`
	ID         *int                   `json:"id,omitempty"`
	Components *[]AnySectionComponent `json:"components"`
	Accessory  AnySectionAccessory    `json:"accessory,omitempty"`
}

func (s Section) IsAnyContainerComponent() bool {
	return true
}

type AnySectionComponent interface {
	MarshalJSON() ([]byte, error)
	GetType() DiscordComponentType
	IsAnySectionComponent() bool
}

type TextDisplayComponent struct {
	Type    DiscordComponentType `json:"type"`
	ID      *int                 `json:"id,omitempty"`
	Content string               `json:"content"`
}

func (t TextDisplayComponent) IsAnyContainerComponent() bool {
	return true
}

func (t TextDisplayComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeTextDisplay
}

func (t TextDisplayComponent) MarshalJSON() ([]byte, error) {
	t.Type = DiscordComponentTypeTextDisplay
	type Alias TextDisplayComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(t),
	})
}

func (t TextDisplayComponent) IsAnySectionComponent() bool {
	return true
}

type AnySectionAccessory interface {
	MarshalJSON() ([]byte, error)
	GetType() DiscordComponentType
	IsAnySectionAccessory() bool
}

type UnfurledMediaItem struct {
	URL          string            `json:"url"`
	ProxyURL     string            `json:"proxy_url,omitempty"`
	Height       int               `json:"height,omitempty"`
	Width        int               `json:"width,omitempty"`
	ContentType  string            `json:"content_type,omitempty"`
	AttachmentID *DiscordSnowflake `json:"attachment_id,omitempty"`
}

type ThumbnailComponent struct {
	Type        DiscordComponentType `json:"type"`
	ID          *int                 `json:"id,omitempty"`
	Description string               `json:"description,omitempty"`
	Spoiler     bool                 `json:"spoiler,omitempty"`
	Media       *UnfurledMediaItem   `json:"media,omitempty"`
}

func (t ThumbnailComponent) MarshalJSON() ([]byte, error) {
	t.Type = DiscordComponentTypeThumbnail
	type Alias ThumbnailComponent
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(t),
	})
}

func (t ThumbnailComponent) GetType() DiscordComponentType {
	return DiscordComponentTypeThumbnail
}

func (t ThumbnailComponent) IsAnySectionAccessory() bool {
	return true
}

type ImageDisplayComponent struct {
	Type   DiscordComponentType `json:"type"`
	ID     *int                 `json:"id,omitempty"`
	Source string               `json:"source"`
	Alt    *string              `json:"alt,omitempty"`
}

func (s Section) MarshalJSON() ([]byte, error) {
	s.Type = DiscordComponentTypeSection
	type Alias Section
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(s),
	})
}

func (s Section) GetType() DiscordComponentType {
	return DiscordComponentTypeSection
}

type Container struct {
	Type        DiscordComponentType     `json:"type"`
	ID          *int                     `json:"id,omitempty"`
	Components  *[]AnyContainerComponent `json:"components"`
	AccentColor int                      `json:"accent_color,omitempty"`
	Spoiler     bool                     `json:"spoiler,omitempty"`
}

func (c Container) GetType() DiscordComponentType {
	return DiscordComponentTypeContainer
}

func (c Container) MarshalJSON() ([]byte, error) {
	c.Type = DiscordComponentTypeContainer
	type Alias Container
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(c),
	})
}

type AnyContainerComponent interface {
	MarshalJSON() ([]byte, error)
	GetType() DiscordComponentType
	IsAnyContainerComponent() bool
}

type SelectDefaultValue struct {
	ID   DiscordSnowflake       `json:"id"`
	Type SelectDefaultValueType `json:"type"`
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
