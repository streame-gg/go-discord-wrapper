package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#mentionable-select
type MentionableSelectMenu struct {
	Type          discord.ComponentType `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required"`
	Disabled      bool                  `json:"disabled"`
	DefaultValues []SelectDefaultValue  `json:"default_values,omitempty"`
}

func (m *MentionableSelectMenu) IsAnyContainerAccessory() bool {
	return true
}

func (m *MentionableSelectMenu) UnmarshalJSON(data []byte) error {
	type Alias MentionableSelectMenu
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*m = MentionableSelectMenu(*raw.Alias)
	return nil
}

func (m *MentionableSelectMenu) IsAnyLabelComponent() {

}

func (m *MentionableSelectMenu) MarshalJSON() ([]byte, error) {
	type Alias MentionableSelectMenu
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*m),
		Type:  discord.ComponentTypeMentionableSelect,
	})
}

func (m *MentionableSelectMenu) GetType() discord.ComponentType {
	return discord.ComponentTypeMentionableSelect
}

// https://docs.discord.com/developers/components/reference#mentionable-select
type MentionableComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []discord.Snowflake   `json:"values"`
	ID            int                   `json:"id"`
	CustomID      string                `json:"custom_id"`
	ComponentType discord.ComponentType `json:"component_type"`
	Resolved      discord.ResolvedData  `json:"resolved"`
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
