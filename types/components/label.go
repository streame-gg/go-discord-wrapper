package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type LabelComponent struct {
	Type        common.ComponentType `json:"type"`
	ID          *int                 `json:"id,omitempty"`
	Label       string               `json:"label"`
	Description string               `json:"description,omitempty"`
	Component   AnyChildComponent    `json:"component,omitempty"`
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
			Type common.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(raw.Component, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case common.ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(raw.Component, &t); err != nil {
				return err
			}
			l.Component = t
		case common.ComponentTypeFileUpload:
			var f *FileUploadComponent
			if err := json.Unmarshal(raw.Component, &f); err != nil {
				return err
			}
			l.Component = f
		case common.ComponentTypeStringSelectMenu:
			var s *StringSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &s); err != nil {
				return err
			}
			l.Component = s
		case common.ComponentTypeUserSelectMenu:
			var u *UserSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &u); err != nil {
				return err
			}
			l.Component = u
		case common.ComponentTypeRoleSelectMenu:
			var r *RoleSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &r); err != nil {
				return err
			}
			l.Component = r
		case common.ComponentTypeMentionableMenu:
			var m *MentionableSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &m); err != nil {
				return err
			}
			l.Component = m
		case common.ComponentTypeChannelSelect:
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
	l.Type = common.ComponentTypeLabel
	type Alias LabelComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

func (l *LabelComponent) GetType() common.ComponentType {
	return common.ComponentTypeLabel
}

type LabelComponentInteractionResponse struct {
	Type     common.ComponentType `json:"type"`
	Value    string               `json:"values"`
	ID       *int                 `json:"id,omitempty"`
	CustomID string               `json:"custom_id,omitempty"`
}

func (l *LabelComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (l *LabelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	l.Type = common.ComponentTypeLabel

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
	Type        common.ComponentType             `json:"type"`
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
		Type common.ComponentType `json:"type"`
	}
	if err := json.Unmarshal(*raw.Component, &probe); err != nil {
		return err
	}

	var c AnyComponentInteractionResponse

	switch probe.Type {
	case common.ComponentTypeUserSelectMenu:
		c = &UserSelectComponentInteractionResponse{}
	case common.ComponentTypeRoleSelectMenu:
		c = &RoleComponentInteractionResponse{}
	case common.ComponentTypeStringSelectMenu:
		c = &StringSelectComponentInteractionResponse{}
	case common.ComponentTypeChannelSelect:
		c = &ChannelComponentInteractionResponse{}
	case common.ComponentTypeMentionableMenu:
		c = &MentionableComponentInteractionResponse{}
	case common.ComponentTypeTextDisplay:
		c = &TextDisplayComponentInteractionResponse{}
	case common.ComponentTypeTextInput:
		c = &TextInputComponentInteractionResponse{}
	case common.ComponentTypeFileUpload:
		c = &FileUploadComponentInteractionResponse{}
	case common.ComponentTypeLabel:
		c = &LabelComponentInteractionResponse{}
	case common.ComponentTypeRadioGroup:
		c = &RadioGroupComponentInteractionResponse{}
	case common.ComponentTypeCheckboxGroup:
		c = &CheckboxGroupComponentInteractionResponse{}
	case common.ComponentTypeCheckbox:
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
