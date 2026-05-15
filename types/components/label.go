package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type LabelComponent struct {
	Type        discord.ComponentType `json:"type"`
	ID          *int                  `json:"id,omitempty"`
	Label       string                `json:"label"`
	Description string                `json:"description,omitempty"`
	Component   AnyChildComponent     `json:"component,omitempty"`
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
			Type discord.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(raw.Component, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case discord.ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(raw.Component, &t); err != nil {
				return err
			}
			l.Component = t
		case discord.ComponentTypeFileUpload:
			var f *FileUploadComponent
			if err := json.Unmarshal(raw.Component, &f); err != nil {
				return err
			}
			l.Component = f
		case discord.ComponentTypeStringSelectMenu:
			var s *StringSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &s); err != nil {
				return err
			}
			l.Component = s
		case discord.ComponentTypeUserSelectMenu:
			var u *UserSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &u); err != nil {
				return err
			}
			l.Component = u
		case discord.ComponentTypeRoleSelectMenu:
			var r *RoleSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &r); err != nil {
				return err
			}
			l.Component = r
		case discord.ComponentTypeMentionableMenu:
			var m *MentionableSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &m); err != nil {
				return err
			}
			l.Component = m
		case discord.ComponentTypeChannelSelect:
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
	l.Type = discord.ComponentTypeLabel
	type Alias LabelComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

func (l *LabelComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeLabel
}

type LabelComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	Value    string                `json:"values"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id,omitempty"`
}

func (l *LabelComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (l *LabelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	l.Type = discord.ComponentTypeLabel

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

type ComponentLabelComponent struct {
	Type        discord.ComponentType            `json:"type"`
	ID          *int                             `json:"id,omitempty"`
	Label       *string                          `json:"label"`
	Description *string                          `json:"description,omitempty"`
	Component   *AnyComponentInteractionResponse `json:"component,omitempty"`
}

func (l *ComponentLabelComponent) UnmarshalJSON(data []byte) error {
	type Alias ComponentLabelComponent
	raw := &struct {
		*Alias
		Component *json.RawMessage `json:"component,omitempty"`
	}{
		Alias: (*Alias)(l),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Component == nil {
		return nil
	}

	var probe struct {
		Type discord.ComponentType `json:"type"`
	}
	if err := json.Unmarshal(*raw.Component, &probe); err != nil {
		return err
	}

	var c AnyComponentInteractionResponse

	switch probe.Type {
	case discord.ComponentTypeUserSelectMenu:
		c = &UserSelectComponentInteractionResponse{}
	case discord.ComponentTypeRoleSelectMenu:
		c = &RoleComponentInteractionResponse{}
	case discord.ComponentTypeStringSelectMenu:
		c = &StringSelectComponentInteractionResponse{}
	case discord.ComponentTypeChannelSelect:
		c = &ChannelComponentInteractionResponse{}
	case discord.ComponentTypeMentionableMenu:
		c = &MentionableComponentInteractionResponse{}
	case discord.ComponentTypeTextDisplay:
		c = &TextDisplayComponentInteractionResponse{}
	case discord.ComponentTypeTextInput:
		c = &TextInputComponentInteractionResponse{}
	case discord.ComponentTypeFileUpload:
		c = &FileUploadComponentInteractionResponse{}
	case discord.ComponentTypeLabel:
		c = &LabelComponentInteractionResponse{}
	case discord.ComponentTypeRadioGroup:
		c = &RadioGroupComponentInteractionResponse{}
	case discord.ComponentTypeCheckboxGroup:
		c = &CheckboxGroupComponentInteractionResponse{}
	case discord.ComponentTypeCheckbox:
		c = &CheckboxComponentInteractionResponse{}

	default:
		return fmt.Errorf("unknown interaction component type: %d", probe.Type)
	}

	if err := json.Unmarshal(*raw.Component, c); err != nil {
		return err
	}

	l.Component = &c

	return nil
}
