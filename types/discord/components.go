package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/components/reference#component-object-component-types
type ComponentType int

const (
	ComponentTypeActionRow         ComponentType = 1
	ComponentTypeButton            ComponentType = 2
	ComponentTypeStringSelect      ComponentType = 3
	ComponentTypeTextInput         ComponentType = 4
	ComponentTypeUserSelect        ComponentType = 5
	ComponentTypeRoleSelect        ComponentType = 6
	ComponentTypeMentionableSelect ComponentType = 7
	ComponentTypeChannelSelect     ComponentType = 8
	ComponentTypeSection           ComponentType = 9
	ComponentTypeTextDisplay       ComponentType = 10
	ComponentTypeThumbnail         ComponentType = 11
	ComponentTypeMediaGallery      ComponentType = 12
	ComponentTypeFileDisplay       ComponentType = 13
	ComponentTypeSeparator         ComponentType = 14
	ComponentTypeContainer         ComponentType = 17
	ComponentTypeLabel             ComponentType = 18
	ComponentTypeFileUpload        ComponentType = 19
	ComponentTypeRadioGroup        ComponentType = 21
	ComponentTypeCheckboxGroup     ComponentType = 22
	ComponentTypeCheckbox          ComponentType = 23
)

func (c ComponentType) IsAnySelectMenu() bool {
	return c == ComponentTypeStringSelect ||
		c == ComponentTypeUserSelect ||
		c == ComponentTypeRoleSelect ||
		c == ComponentTypeMentionableSelect ||
		c == ComponentTypeChannelSelect
}

// https://docs.discord.com/developers/components/reference#component-object
type AnyComponent interface {
	GetType() ComponentType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// https://docs.discord.com/developers/components/reference#component-object
type ComponentWrapper struct {
	Component AnyComponent
}

func (c *ComponentWrapper) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type ComponentType `json:"type"`
	}

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	c.Component = &RawComponent{
		Type: probe.Type,
		Data: data,
	}

	return nil
}

func (c *ComponentWrapper) MarshalJSON() ([]byte, error) {
	if c.Component == nil {
		return []byte("null"), nil
	}
	return c.Component.MarshalJSON()
}

// https://docs.discord.com/developers/components/reference#component-object
type RawComponent struct {
	Type ComponentType
	Data json.RawMessage
}

func (r *RawComponent) GetType() ComponentType {
	return r.Type
}

func (r *RawComponent) MarshalJSON() ([]byte, error) {
	return r.Data, nil
}

func (r *RawComponent) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type ComponentType `json:"type"`
	}

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	r.Type = probe.Type
	r.Data = data
	return nil
}
