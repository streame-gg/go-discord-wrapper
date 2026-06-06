package components

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#label
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

	if raw.Alias == nil {
		return nil
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
		case discord.ComponentTypeStringSelect:
			var s *StringSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &s); err != nil {
				return err
			}
			l.Component = s
		case discord.ComponentTypeUserSelect:
			var u *UserSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &u); err != nil {
				return err
			}
			l.Component = u
		case discord.ComponentTypeRoleSelect:
			var r *RoleSelectMenuComponent
			if err := json.Unmarshal(raw.Component, &r); err != nil {
				return err
			}
			l.Component = r
		case discord.ComponentTypeMentionableSelect:
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
			return fmt.Errorf("unknown component type: %d", probe.Type)
		}
	}

	return nil
}

func (l *LabelComponent) MarshalJSON() ([]byte, error) {
	type Alias LabelComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*l),
		Type:  discord.ComponentTypeLabel,
	})
}

func (l *LabelComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeLabel
}

// LabelComponentInteractionResponse https://docs.discord.com/developers/components/reference#label-label-interaction-response-structure
type LabelComponentInteractionResponse struct {
	Type      discord.ComponentType           `json:"type"`
	ID        *int                            `json:"id,omitempty"`
	Component AnyComponentInteractionResponse `json:"component,omitempty"`
}

func (l *LabelComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (l *LabelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias LabelComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*l),
		Type:  discord.ComponentTypeLabel,
	})
}

func (l *LabelComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type      discord.ComponentType `json:"type"`
		ID        *int                  `json:"id,omitempty"`
		Component *json.RawMessage      `json:"component,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.Type = raw.Type
	l.ID = raw.ID

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
	case discord.ComponentTypeTextInput:
		c = &TextInputComponentInteractionResponse{}
	case discord.ComponentTypeStringSelect:
		c = &StringSelectComponentInteractionResponse{}
	case discord.ComponentTypeUserSelect:
		c = &UserSelectComponentInteractionResponse{}
	case discord.ComponentTypeRoleSelect:
		c = &RoleComponentInteractionResponse{}
	case discord.ComponentTypeMentionableSelect:
		c = &MentionableComponentInteractionResponse{}
	case discord.ComponentTypeChannelSelect:
		c = &ChannelComponentInteractionResponse{}
	case discord.ComponentTypeFileUpload:
		c = &FileUploadComponentInteractionResponse{}
	case discord.ComponentTypeRadioGroup:
		c = &RadioGroupComponentInteractionResponse{}
	case discord.ComponentTypeCheckboxGroup:
		c = &CheckboxGroupComponentInteractionResponse{}
	case discord.ComponentTypeCheckbox:
		c = &CheckboxComponentInteractionResponse{}
	default:
		return fmt.Errorf("unknown label child component type: %d", probe.Type)
	}

	if err := json.Unmarshal(*raw.Component, c); err != nil {
		return err
	}
	l.Component = c
	return nil
}

// https://docs.discord.com/developers/components/reference#label
type ComponentLabelComponent struct {
	Type        discord.ComponentType           `json:"type"`
	ID          *int                            `json:"id,omitempty"`
	Label       *string                         `json:"label"`
	Description *string                         `json:"description,omitempty"`
	Component   AnyComponentInteractionResponse `json:"component,omitempty"`
}

func (l *ComponentLabelComponent) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type        discord.ComponentType `json:"type"`
		ID          *int                  `json:"id,omitempty"`
		Label       *string               `json:"label"`
		Description *string               `json:"description,omitempty"`
		Component   *json.RawMessage      `json:"component,omitempty"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.Type = raw.Type
	l.ID = raw.ID
	l.Label = raw.Label
	l.Description = raw.Description

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
	case discord.ComponentTypeUserSelect:
		c = &UserSelectComponentInteractionResponse{}
	case discord.ComponentTypeRoleSelect:
		c = &RoleComponentInteractionResponse{}
	case discord.ComponentTypeStringSelect:
		c = &StringSelectComponentInteractionResponse{}
	case discord.ComponentTypeChannelSelect:
		c = &ChannelComponentInteractionResponse{}
	case discord.ComponentTypeMentionableSelect:
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

	l.Component = c
	return nil
}
