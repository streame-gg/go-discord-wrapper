package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type MentionableSelectMenuComponent struct {
	Type          common.ComponentType  `json:"type"`
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

func (m *MentionableSelectMenuComponent) IsAnyLabelComponent() {

}

func (m *MentionableSelectMenuComponent) MarshalJSON() ([]byte, error) {
	m.Type = common.ComponentTypeMentionableMenu
	type Alias MentionableSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

func (m *MentionableSelectMenuComponent) GetType() common.ComponentType {
	return common.ComponentTypeMentionableMenu
}

type MentionableComponentInteractionResponse struct {
	Type          common.ComponentType `json:"type"`
	Values        []common.Snowflake   `json:"values"`
	ID            *int                 `json:"id,omitempty"`
	CustomID      string               `json:"custom_id,omitempty"`
	ComponentType common.ComponentType `json:"component_type"`
	Resolved      *common.ResolvedData `json:"resolved,omitempty"`
}

func (m *MentionableComponentInteractionResponse) IsInteractionResponseDataComponent() {
}

func (m *MentionableComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	m.ComponentType = common.ComponentTypeMentionableMenu
	m.Type = common.ComponentTypeMentionableMenu

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
