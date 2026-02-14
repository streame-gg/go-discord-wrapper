package types

import (
	"encoding/json"
	"errors"
)

type ComponentType int

const (
	ComponentTypeActionRow        ComponentType = 1
	ComponentTypeButton           ComponentType = 2
	ComponentTypeStringSelectMenu ComponentType = 3
	ComponentTypeTextInput        ComponentType = 4
	ComponentTypeUserSelectMenu   ComponentType = 5
	ComponentTypeRoleSelectMenu   ComponentType = 6
	ComponentTypeMentionableMenu  ComponentType = 7
	ComponentTypeChannelSelect    ComponentType = 8
	ComponentTypeSection          ComponentType = 9
	ComponentTypeTextDisplay      ComponentType = 10
	ComponentTypeThumbnail        ComponentType = 11
	ComponentTypeMediaGallery     ComponentType = 12
	ComponentTypeFileDisplay      ComponentType = 13
	ComponentTypeSeparator        ComponentType = 14
	ComponentTypeContainer        ComponentType = 17
	ComponentTypeLabel            ComponentType = 18
	ComponentTypeFileUpload       ComponentType = 19
	ComponentTypeRadioGroup       ComponentType = 21
	ComponentTypeCheckboxGroup    ComponentType = 22
	ComponentTypeCheckbox         ComponentType = 23
)

func (c ComponentType) IsAnySelectMenu() bool {
	return c == ComponentTypeStringSelectMenu ||
		c == ComponentTypeUserSelectMenu ||
		c == ComponentTypeRoleSelectMenu ||
		c == ComponentTypeMentionableMenu ||
		c == ComponentTypeChannelSelect
}

type AnyComponent interface {
	GetType() ComponentType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type ActionRow struct {
	Type       ComponentType  `json:"type"`
	ID         *int           `json:"id"`
	Components []AnyComponent `json:"components"`
}

func (a *ActionRow) GetType() ComponentType {
	return ComponentTypeActionRow
}

func NewActionRow(id *int, components ...AnyComponent) *ActionRow {
	return &ActionRow{
		Type:       ComponentTypeActionRow,
		Components: components,
		ID:         id,
	}
}

func (a *ActionRow) MarshalJSON() ([]byte, error) {
	a.Type = ComponentTypeActionRow
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
			Type ComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case ComponentTypeButton:
			var b *ButtonComponent
			if err := json.Unmarshal(c, &b); err != nil {
				return err
			}
			a.Components = append(a.Components, b)
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
	Type     ComponentType `json:"type"`
	ID       *int          `json:"id,omitempty"`
	Style    ButtonStyle   `json:"style"`
	Label    string        `json:"label,omitempty"`
	Emoji    *Emoji        `json:"emoji,omitempty"`
	CustomID string        `json:"custom_id,omitempty"`
	SkuID    *Snowflake    `json:"sku_id,omitempty"`
	URL      string        `json:"url,omitempty"`
	Disabled bool          `json:"disabled,omitempty"`
}

func (b *ButtonComponent) UnmarshalJSON(data []byte) error {
	type Alias ButtonComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*b = ButtonComponent(*raw.Alias)
	return nil
}

func (b *ButtonComponent) MarshalJSON() ([]byte, error) {
	b.Type = ComponentTypeButton
	type Alias ButtonComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}

func (b *ButtonComponent) GetType() ComponentType {
	return ComponentTypeButton
}

func (b *ButtonComponent) IsAnySectionAccessory() bool {
	return true
}

func (b *ButtonComponent) IsAnyContainerAccessory() bool {
	return true
}

type StringSelectMenuComponent struct {
	Type        ComponentType                      `json:"type"`
	ID          *int                               `json:"id,omitempty"`
	CustomID    string                             `json:"custom_id"`
	Placeholder string                             `json:"placeholder,omitempty"`
	MinValues   *int                               `json:"min_values,omitempty"`
	MaxValues   *int                               `json:"max_values,omitempty"`
	Required    bool                               `json:"required,omitempty"`
	Options     *[]StringSelectMenuComponentOption `json:"options"`
	Disabled    bool                               `json:"disabled,omitempty"`
}

func (s *StringSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (s *StringSelectMenuComponent) MarshalJSON() ([]byte, error) {
	s.Type = ComponentTypeStringSelectMenu
	type Alias StringSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *StringSelectMenuComponent) GetType() ComponentType {
	return ComponentTypeStringSelectMenu
}

func (s *StringSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias StringSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = StringSelectMenuComponent(*raw.Alias)
	return nil
}

type StringSelectMenuComponentOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

func (s *StringSelectMenuComponent) IsAnyLabelComponent() bool {
	return true
}

type UserSelectMenuComponent struct {
	Type          ComponentType         `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (u *UserSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (u *UserSelectMenuComponent) MarshalJSON() ([]byte, error) {
	u.Type = ComponentTypeUserSelectMenu
	type Alias UserSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

func (u *UserSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias UserSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*u = UserSelectMenuComponent(*raw.Alias)
	return nil
}

func (u *UserSelectMenuComponent) GetType() ComponentType {
	return ComponentTypeUserSelectMenu
}

func (u *UserSelectMenuComponent) IsAnyLabelComponent() bool {
	return true
}

type RoleSelectMenuComponent struct {
	Type          ComponentType         `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (r *RoleSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (r *RoleSelectMenuComponent) MarshalJSON() ([]byte, error) {
	r.Type = ComponentTypeRoleSelectMenu
	type Alias UserSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RoleSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias RoleSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*r = RoleSelectMenuComponent(*raw.Alias)
	return nil
}

func (r *RoleSelectMenuComponent) GetType() ComponentType {
	return ComponentTypeRoleSelectMenu
}

func (r *RoleSelectMenuComponent) IsAnyLabelComponent() bool {
	return true
}

type MentionableSelectMenuComponent struct {
	Type          ComponentType         `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (m *MentionableSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (m *MentionableSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias MentionableSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*m = MentionableSelectMenuComponent(*raw.Alias)
	return nil
}

func (m *MentionableSelectMenuComponent) IsAnyLabelComponent() bool {
	return true
}

func (m *MentionableSelectMenuComponent) MarshalJSON() ([]byte, error) {
	m.Type = ComponentTypeMentionableMenu
	type Alias MentionableSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

func (m *MentionableSelectMenuComponent) GetType() ComponentType {
	return ComponentTypeMentionableMenu
}

type ChannelSelectMenuComponent struct {
	Type          ComponentType         `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (c *ChannelSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (c *ChannelSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias ChannelSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = ChannelSelectMenuComponent(*raw.Alias)
	return nil
}

func (c *ChannelSelectMenuComponent) MarshalJSON() ([]byte, error) {
	c.Type = ComponentTypeChannelSelect
	type Alias ChannelSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *ChannelSelectMenuComponent) GetType() ComponentType {
	return ComponentTypeChannelSelect
}

func (c *ChannelSelectMenuComponent) IsAnyLabelComponent() bool {
	return true
}

type SelectDefaultValueType string

const (
	SelectDefaultValueTypeUser    SelectDefaultValueType = "user"
	SelectDefaultValueTypeRole    SelectDefaultValueType = "role"
	SelectDefaultValueTypeChannel SelectDefaultValueType = "channel"
)

type Section struct {
	Type       ComponentType          `json:"type"`
	ID         *int                   `json:"id,omitempty"`
	Components *[]AnySectionComponent `json:"components"`
	Accessory  AnySectionAccessory    `json:"accessory,omitempty"`
}

func (s *Section) IsAnyContainerComponent() bool {
	return true
}

func (s *Section) UnmarshalJSON(data []byte) error {
	type Alias Section

	var raw struct {
		Alias
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = Section(raw.Alias)

	for _, c := range raw.Components {
		var probe struct {
			Type ComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case ComponentTypeTextDisplay:
			var t *TextDisplayComponent
			if err := json.Unmarshal(c, &t); err != nil {
				return err
			}
			*s.Components = append(*s.Components, t)
		}
	}

	return nil
}

type AnySectionComponent interface {
	IsAnySectionComponent() bool
}

type TextDisplayComponent struct {
	Type    ComponentType `json:"type"`
	ID      *int          `json:"id,omitempty"`
	Content string        `json:"content"`
}

func (t *TextDisplayComponent) UnmarshalJSON(data []byte) error {
	type Alias TextDisplayComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextDisplayComponent(*raw.Alias)
	return nil
}

func (t *TextDisplayComponent) IsAnyContainerComponent() bool {
	return true
}

func (t *TextDisplayComponent) GetType() ComponentType {
	return ComponentTypeTextDisplay
}

func (t *TextDisplayComponent) MarshalJSON() ([]byte, error) {
	t.Type = ComponentTypeTextDisplay
	type Alias TextDisplayComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextDisplayComponent) IsAnySectionComponent() bool {
	return true
}

type AnySectionAccessory interface {
	IsAnySectionAccessory() bool
}

type UnfurledMediaItem struct {
	URL          string     `json:"url"`
	ProxyURL     string     `json:"proxy_url,omitempty"`
	Height       int        `json:"height,omitempty"`
	Width        int        `json:"width,omitempty"`
	ContentType  string     `json:"content_type,omitempty"`
	AttachmentID *Snowflake `json:"attachment_id,omitempty"`
}

type ThumbnailComponent struct {
	Type        ComponentType      `json:"type"`
	ID          *int               `json:"id,omitempty"`
	Description string             `json:"description,omitempty"`
	Spoiler     bool               `json:"spoiler,omitempty"`
	Media       *UnfurledMediaItem `json:"media,omitempty"`
}

func (t *ThumbnailComponent) UnmarshalJSON(data []byte) error {
	type Alias ThumbnailComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = ThumbnailComponent(*raw.Alias)
	return nil
}

func (t *ThumbnailComponent) MarshalJSON() ([]byte, error) {
	t.Type = ComponentTypeThumbnail
	type Alias ThumbnailComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *ThumbnailComponent) GetType() ComponentType {
	return ComponentTypeThumbnail
}

func (t *ThumbnailComponent) IsAnySectionAccessory() bool {
	return true
}

type ImageDisplayComponent struct {
	Type   ComponentType `json:"type"`
	ID     *int          `json:"id,omitempty"`
	Source string        `json:"source"`
	Alt    *string       `json:"alt,omitempty"`
}

func (s *Section) MarshalJSON() ([]byte, error) {
	s.Type = ComponentTypeSection
	type Alias Section
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *Section) GetType() ComponentType {
	return ComponentTypeSection
}

type Container struct {
	Type        ComponentType            `json:"type"`
	ID          *int                     `json:"id,omitempty"`
	Components  *[]AnyContainerComponent `json:"components"`
	AccentColor int                      `json:"accent_color,omitempty"`
	Spoiler     bool                     `json:"spoiler,omitempty"`
}

func (c *Container) UnmarshalJSON(data []byte) error {
	type Alias Container

	var raw struct {
		*Alias
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = Container(*raw.Alias)

	for _, comp := range raw.Components {
		var probe struct {
			Type ComponentType `json:"type"`
		}

		if err := json.Unmarshal(comp, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case ComponentTypeMediaGallery:
			var m *MediaGalleryComponent
			if err := json.Unmarshal(comp, &m); err != nil {
				return err
			}
			*c.Components = append(*c.Components, m)
		case ComponentTypeFileDisplay:
			var f *FileComponent
			if err := json.Unmarshal(comp, &f); err != nil {
				return err
			}
			*c.Components = append(*c.Components, f)
		case ComponentTypeSeparator:
			var s *SeparatorComponent
			if err := json.Unmarshal(comp, &s); err != nil {
				return err
			}
			*c.Components = append(*c.Components, s)
		case ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(comp, &t); err != nil {
				return err
			}
			*c.Components = append(*c.Components, t)
		case ComponentTypeActionRow:
			var a *ActionRow
			if err := json.Unmarshal(comp, &a); err != nil {
				return err
			}
			*c.Components = append(*c.Components, a)
		case ComponentTypeTextDisplay:
			var t *TextDisplayComponent
			if err := json.Unmarshal(comp, &t); err != nil {
				return err
			}
			*c.Components = append(*c.Components, t)
		case ComponentTypeSection:
			var s *Section
			if err := json.Unmarshal(comp, &s); err != nil {
				return err
			}
			*c.Components = append(*c.Components, s)
		default:
			return errors.New("unknown container component type" + string(rune(probe.Type)))
		}
	}

	return nil
}

func (c *Container) GetType() ComponentType {
	return ComponentTypeContainer
}

func (c *Container) MarshalJSON() ([]byte, error) {
	c.Type = ComponentTypeContainer
	type Alias Container
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

type AnyContainerComponent interface {
	MarshalJSON() ([]byte, error)
	GetType() ComponentType
	IsAnyContainerComponent() bool
}

type MediaGalleryComponent struct {
	Type  ComponentType       `json:"type"`
	ID    *int                `json:"id,omitempty"`
	Items *[]MediaGalleryItem `json:"items"`
}

func (m *MediaGalleryComponent) UnmarshalJSON(data []byte) error {
	type Alias MediaGalleryComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*m = MediaGalleryComponent(*raw.Alias)
	return nil
}

func (m *MediaGalleryComponent) MarshalJSON() ([]byte, error) {
	m.Type = ComponentTypeMediaGallery
	type Alias MediaGalleryComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

func (m *MediaGalleryComponent) GetType() ComponentType {
	return ComponentTypeMediaGallery
}

func (m *MediaGalleryComponent) IsAnyContainerComponent() bool {
	return true
}

type MediaGalleryItem struct {
	Media       *UnfurledMediaItem `json:"media"`
	Description string             `json:"description,omitempty"`
	Spoiler     bool               `json:"spoiler,omitempty"`
}

type FileComponent struct {
	Type    ComponentType      `json:"type"`
	ID      *int               `json:"id,omitempty"`
	Spoiler bool               `json:"spoiler,omitempty"`
	Name    string             `json:"name,omitempty"`
	Size    int                `json:"size,omitempty"`
	File    *UnfurledMediaItem `json:"file,omitempty"`
}

func (f *FileComponent) UnmarshalJSON(data []byte) error {
	type Alias FileComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*f = FileComponent(*raw.Alias)
	return nil
}

func (f *FileComponent) MarshalJSON() ([]byte, error) {
	f.Type = ComponentTypeFileDisplay
	type Alias FileComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileComponent) GetType() ComponentType {
	return ComponentTypeFileDisplay
}

func (f *FileComponent) IsAnyContainerComponent() bool {
	return true
}

type SeparatorComponentSpacing int

const (
	SeparatorComponentSpacingSmall SeparatorComponentSpacing = 1
	SeparatorComponentSpacingLarge SeparatorComponentSpacing = 2
)

type SeparatorComponent struct {
	Type                      ComponentType             `json:"type"`
	ID                        *int                      `json:"id,omitempty"`
	Divider                   bool                      `json:"divider,omitempty"`
	SeparatorComponentSpacing SeparatorComponentSpacing `json:"spacing,omitempty"`
}

func (s *SeparatorComponent) UnmarshalJSON(data []byte) error {
	type Alias SeparatorComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = SeparatorComponent(*raw.Alias)
	return nil
}

func (s *SeparatorComponent) MarshalJSON() ([]byte, error) {
	s.Type = ComponentTypeSeparator
	type Alias SeparatorComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *SeparatorComponent) GetType() ComponentType {
	return ComponentTypeSeparator
}

func (s *SeparatorComponent) IsAnyContainerComponent() bool {
	return true
}

type LabelComponent struct {
	Type        ComponentType     `json:"type"`
	ID          *int              `json:"id,omitempty"`
	Label       string            `json:"label"`
	Description string            `json:"description,omitempty"`
	Component   AnyChildComponent `json:"component,omitempty"`
}

func (l *LabelComponent) UnmarshalJSON(data []byte) error {
	type Alias LabelComponent
	var raw struct {
		*Alias
		Component json.RawMessage `json:"component,omitempty"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*l = LabelComponent(*raw.Alias)

	if raw.Component != nil {
		var probe struct {
			Type ComponentType `json:"type"`
		}

		if err := json.Unmarshal(raw.Component, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(raw.Component, &t); err != nil {
				return err
			}
			l.Component = t
		case ComponentTypeFileUpload:
			var f *FileUploadComponent
			if err := json.Unmarshal(raw.Component, &f); err != nil {
				return err
			}
			l.Component = f
		case ComponentTypeStringSelectMenu:
			var s *StringSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &s); err != nil {
				return err
			}
			l.Component = s
		case ComponentTypeUserSelectMenu:
			var u *UserSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &u); err != nil {
				return err
			}
			l.Component = u
		case ComponentTypeRoleSelectMenu:
			var r *RoleSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &r); err != nil {
				return err
			}
			l.Component = r
		case ComponentTypeMentionableMenu:
			var m *MentionableSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &m); err != nil {
				return err
			}
			l.Component = m
		case ComponentTypeChannelSelect:
			var c *ChannelSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &c); err != nil {
				return err
			}
			l.Component = c
		default:
			return errors.New("unknown component type" + string(rune(probe.Type)))
		}
	}

	return nil
}

func (l *LabelComponent) MarshalJSON() ([]byte, error) {
	l.Type = ComponentTypeLabel
	type Alias LabelComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

func (l *LabelComponent) GetType() ComponentType {
	return ComponentTypeLabel
}

type AnyChildComponent interface {
	IsAnyLabelComponent() bool
}

type Modal struct {
	Title      string            `json:"title"`
	CustomID   string            `json:"custom_id"`
	Components *[]LabelComponent `json:"components"`
}

type AnyComponentInteractionResponse interface {
	IsInteractionResponseDataComponent() bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type StringSelectComponentInteractionResponse struct {
	Type          ComponentType `json:"type"`
	Values        []string      `json:"values"`
	ID            *int          `json:"id,omitempty"`
	CustomID      string        `json:"custom_id,omitempty"`
	ComponentType ComponentType `json:"component_type"`
}

func (s *StringSelectComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (s *StringSelectComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	s.ComponentType = ComponentTypeStringSelectMenu
	s.Type = ComponentTypeStringSelectMenu

	type Alias StringSelectComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *StringSelectComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias StringSelectComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = StringSelectComponentInteractionResponse(*raw.Alias)
	return nil
}

type TextInputComponentInteractionResponse struct {
	Type     ComponentType `json:"type"`
	Value    string        `json:"value"`
	ID       *int          `json:"id,omitempty"`
	CustomID string        `json:"custom_id"`
}

func (t *TextInputComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (t *TextInputComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	t.Type = ComponentTypeTextInput

	type Alias TextInputComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextInputComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias TextInputComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextInputComponentInteractionResponse(*raw.Alias)
	return nil
}

type UserSelectComponentInteractionResponse struct {
	Type          ComponentType `json:"type"`
	Values        []Snowflake   `json:"values"`
	ID            *int          `json:"id,omitempty"`
	CustomID      string        `json:"custom_id,omitempty"`
	ComponentType ComponentType `json:"component_type"`
	Resolved      *ResolvedData `json:"resolved,omitempty"`
}

func (u *UserSelectComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (u *UserSelectComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	u.ComponentType = ComponentTypeRoleSelectMenu
	u.Type = ComponentTypeRoleSelectMenu

	type Alias UserSelectComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

func (u *UserSelectComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias UserSelectComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*u = UserSelectComponentInteractionResponse(*raw.Alias)
	return nil
}

type RoleComponentInteractionResponse struct {
	Type          ComponentType `json:"type"`
	Values        []Snowflake   `json:"values"`
	ID            *int          `json:"id,omitempty"`
	CustomID      string        `json:"custom_id,omitempty"`
	ComponentType ComponentType `json:"component_type"`
	Resolved      *ResolvedData `json:"resolved,omitempty"`
}

func (r *RoleComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (r *RoleComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	r.ComponentType = ComponentTypeRoleSelectMenu
	r.Type = ComponentTypeRoleSelectMenu

	type Alias RoleComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RoleComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias RoleComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*r = RoleComponentInteractionResponse(*raw.Alias)
	return nil
}

type MentionableComponentInteractionResponse struct {
	Type          ComponentType `json:"type"`
	Values        []Snowflake   `json:"values"`
	ID            *int          `json:"id,omitempty"`
	CustomID      string        `json:"custom_id,omitempty"`
	ComponentType ComponentType `json:"component_type"`
	Resolved      *ResolvedData `json:"resolved,omitempty"`
}

func (m *MentionableComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (m *MentionableComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	m.ComponentType = ComponentTypeMentionableMenu
	m.Type = ComponentTypeMentionableMenu

	type Alias MentionableComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

func (m *MentionableComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias MentionableComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*m = MentionableComponentInteractionResponse(*raw.Alias)
	return nil
}

type LabelComponentInteractionResponse struct {
	Type     ComponentType `json:"type"`
	Value    string        `json:"values"`
	ID       *int          `json:"id,omitempty"`
	CustomID string        `json:"custom_id,omitempty"`
}

func (l *LabelComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (l *LabelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	l.Type = ComponentTypeLabel

	type Alias LabelComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

func (l *LabelComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias LabelComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*l = LabelComponentInteractionResponse(*raw.Alias)
	return nil
}

type TextDisplayComponentInteractionResponse struct {
	Type ComponentType `json:"type"`
	ID   *int          `json:"id,omitempty"`
}

func (t *TextDisplayComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (t *TextDisplayComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	t.Type = ComponentTypeTextDisplay

	type Alias TextDisplayComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextDisplayComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias TextDisplayComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextDisplayComponentInteractionResponse(*raw.Alias)
	return nil
}

type FileUploadComponentInteractionResponse struct {
	Type     ComponentType `json:"type"`
	ID       *int          `json:"id,omitempty"`
	CustomID string        `json:"custom_id"`
	Values   []Snowflake   `json:"values"`
}

func (f *FileUploadComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (f *FileUploadComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	f.Type = ComponentTypeFileUpload

	type Alias FileUploadComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileUploadComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias FileUploadComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*f = FileUploadComponentInteractionResponse(*raw.Alias)
	return nil
}

type ChannelComponentInteractionResponse struct {
	Type          ComponentType `json:"type"`
	Values        []Snowflake   `json:"values"`
	ID            *int          `json:"id,omitempty"`
	CustomID      string        `json:"custom_id,omitempty"`
	ComponentType ComponentType `json:"component_type"`
	Resolved      *ResolvedData `json:"resolved,omitempty"`
}

func (c *ChannelComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (c *ChannelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	c.ComponentType = ComponentTypeChannelSelect
	c.Type = ComponentTypeChannelSelect

	type Alias ChannelComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *ChannelComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias ChannelComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = ChannelComponentInteractionResponse(*raw.Alias)
	return nil
}

func (m Modal) IsInteractionResponseData() bool {
	return true
}

func (m Modal) MarshalJSON() ([]byte, error) {
	type Alias Modal
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

type TextInputStyle int

const (
	TextInputStyleShort     TextInputStyle = 1
	TextInputStyleParagraph TextInputStyle = 2
)

type TextInputComponent struct {
	Type        ComponentType  `json:"type"`
	ID          *int           `json:"id,omitempty"`
	CustomID    string         `json:"custom_id"`
	Style       TextInputStyle `json:"style"`
	MinLength   *int           `json:"min_length,omitempty"`
	MaxLength   *int           `json:"max_length,omitempty"`
	Required    *bool          `json:"required,omitempty"`
	Value       string         `json:"value,omitempty"`
	Placeholder string         `json:"placeholder,omitempty"`
}

func (t *TextInputComponent) MarshalJSON() ([]byte, error) {
	t.Type = ComponentTypeTextInput
	type Alias TextInputComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextInputComponent) UnmarshalJSON(data []byte) error {
	type Alias TextInputComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextInputComponent(*raw.Alias)
	return nil
}

func (t *TextInputComponent) GetType() ComponentType {
	return ComponentTypeTextInput
}

func (t *TextInputComponent) IsAnyContainerComponent() bool {
	return true
}

func (t *TextInputComponent) IsAnyLabelComponent() bool {
	return true
}

type FileUploadComponent struct {
	Type      ComponentType `json:"type"`
	ID        *int          `json:"id,omitempty"`
	CustomID  string        `json:"custom_id"`
	Required  *bool         `json:"required,omitempty"`
	MinValues *int          `json:"min_values,omitempty"`
	MaxValues *int          `json:"max_values,omitempty"`
}

func (f *FileUploadComponent) MarshalJSON() ([]byte, error) {
	f.Type = ComponentTypeFileUpload
	type Alias FileUploadComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileUploadComponent) UnmarshalJSON(data []byte) error {
	type Alias FileUploadComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*f = FileUploadComponent(*raw.Alias)
	return nil
}

func (f *FileUploadComponent) GetType() ComponentType {
	return ComponentTypeFileUpload
}

func (f *FileUploadComponent) IsAnyLabelComponent() bool {
	return true
}

type RadioGroupComponent struct {
	Type     ComponentType                `json:"type"`
	ID       *int                         `json:"id,omitempty"`
	CustomID string                       `json:"custom_id"`
	Options  *[]RadioGroupComponentOption `json:"options"`
	Required *bool                        `json:"required,omitempty"`
}

type RadioGroupComponentOption struct {
	Value       string  `json:"value"`
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	Default     *bool   `json:"default,omitempty"`
}

func (r *RadioGroupComponent) MarshalJSON() ([]byte, error) {
	r.Type = ComponentTypeRadioGroup
	type Alias RadioGroupComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RadioGroupComponent) UnmarshalJSON(data []byte) error {
	type Alias RadioGroupComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*r = RadioGroupComponent(*raw.Alias)
	return nil
}

func (r *RadioGroupComponent) GetType() ComponentType {
	return ComponentTypeRadioGroup
}

func (r *RadioGroupComponent) IsAnyLabelComponent() bool {
	return true
}

type CheckboxGroupComponent struct {
	Type      ComponentType                   `json:"type"`
	ID        *int                            `json:"id,omitempty"`
	CustomID  string                          `json:"custom_id"`
	Options   *[]CheckboxGroupComponentOption `json:"options"`
	MinValues *int                            `json:"min_values,omitempty"`
	MaxValues *int                            `json:"max_values,omitempty"`
	Required  *bool                           `json:"required,omitempty"`
}

type CheckboxGroupComponentOption struct {
	Value       string  `json:"value"`
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	Default     *bool   `json:"default,omitempty"`
}

func (c *CheckboxGroupComponent) MarshalJSON() ([]byte, error) {
	c.Type = ComponentTypeCheckboxGroup
	type Alias CheckboxGroupComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *CheckboxGroupComponent) UnmarshalJSON(data []byte) error {
	type Alias CheckboxGroupComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = CheckboxGroupComponent(*raw.Alias)
	return nil
}

func (c *CheckboxGroupComponent) GetType() ComponentType {
	return ComponentTypeCheckboxGroup
}

func (c *CheckboxGroupComponent) IsAnyLabelComponent() bool {
	return true
}

type CheckboxComponent struct {
	Type     ComponentType `json:"type"`
	ID       *int          `json:"id,omitempty"`
	CustomID string        `json:"custom_id"`
	Default  *bool         `json:"default,omitempty"`
}

func (c *CheckboxComponent) MarshalJSON() ([]byte, error) {
	c.Type = ComponentTypeCheckbox
	type Alias CheckboxComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *CheckboxComponent) UnmarshalJSON(data []byte) error {
	type Alias CheckboxComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = CheckboxComponent(*raw.Alias)
	return nil
}

func (c *CheckboxComponent) GetType() ComponentType {
	return ComponentTypeCheckbox
}

func (c *CheckboxComponent) IsAnyLabelComponent() bool {
	return true
}

type SelectDefaultValue struct {
	ID   Snowflake              `json:"id"`
	Type SelectDefaultValueType `json:"type"`
}

type ApplicationCommandInteractionOptionType int

const (
	ApplicationCommandInteractionOptionTypeSubCommand      ApplicationCommandInteractionOptionType = 1
	ApplicationCommandInteractionOptionTypeSubCommandGroup ApplicationCommandInteractionOptionType = 2
	ApplicationCommandInteractionOptionTypeString          ApplicationCommandInteractionOptionType = 3
	ApplicationCommandInteractionOptionTypeInteger         ApplicationCommandInteractionOptionType = 4
	ApplicationCommandInteractionOptionTypeBoolean         ApplicationCommandInteractionOptionType = 5
	ApplicationCommandInteractionOptionTypeUser            ApplicationCommandInteractionOptionType = 6
	ApplicationCommandInteractionOptionTypeChannel         ApplicationCommandInteractionOptionType = 7
	ApplicationCommandInteractionOptionTypeRole            ApplicationCommandInteractionOptionType = 8
	ApplicationCommandInteractionOptionTypeMentionable     ApplicationCommandInteractionOptionType = 9
	ApplicationCommandInteractionOptionTypeNumber          ApplicationCommandInteractionOptionType = 10
	ApplicationCommandInteractionOptionTypeAttachment      ApplicationCommandInteractionOptionType = 11
)
