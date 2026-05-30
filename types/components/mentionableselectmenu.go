package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type MentionableSelectMenuComponent struct {
	Type          discord.ComponentType `json:"type"`
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

	if raw.Alias == nil {
		return nil
	}
	*m = MentionableSelectMenuComponent(*raw.Alias)
	return nil
}

func (m *MentionableSelectMenuComponent) IsAnyLabelComponent() {

}

func (m *MentionableSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type Alias MentionableSelectMenuComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*m),
		Type:  discord.ComponentTypeMentionableSelect,
	})
}

func (m *MentionableSelectMenuComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeMentionableSelect
}

type MentionableComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []discord.Snowflake   `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id,omitempty"`
	ComponentType discord.ComponentType `json:"component_type"`
	Resolved      *discord.ResolvedData `json:"resolved,omitempty"`
}

func (m *MentionableComponentInteractionResponse) IsInteractionResponseDataComponent() {
}

func (m *MentionableComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias MentionableComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type          discord.ComponentType `json:"type"`
		ComponentType discord.ComponentType `json:"component_type"`
	}{
		Alias:         Alias(*m),
		Type:          discord.ComponentTypeMentionableSelect,
		ComponentType: discord.ComponentTypeMentionableSelect,
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

	if raw.Alias == nil {
		return nil
	}
	*m = MentionableComponentInteractionResponse(*raw.Alias)
	return nil
}
