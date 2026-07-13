package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#string-select
type StringSelectMenu struct {
	Type        discord.ComponentType             `json:"type"`
	ID          *int                              `json:"id,omitempty"`
	CustomID    string                            `json:"custom_id"`
	Placeholder string                            `json:"placeholder,omitempty"`
	MinValues   *int                              `json:"min_values,omitempty"`
	MaxValues   *int                              `json:"max_values,omitempty"`
	Required    bool                              `json:"required"`
	Options     []StringSelectMenuComponentOption `json:"options"`
	Disabled    bool                              `json:"disabled"`
}

func (s *StringSelectMenu) IsAnyContainerAccessory() bool {
	return true
}

func (s *StringSelectMenu) MarshalJSON() ([]byte, error) {
	type Alias StringSelectMenu
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*s),
		Type:  discord.ComponentTypeStringSelect,
	})
}

func (s *StringSelectMenu) GetType() discord.ComponentType {
	return discord.ComponentTypeStringSelect
}

func (s *StringSelectMenu) UnmarshalJSON(data []byte) error {
	type Alias StringSelectMenu
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*s = StringSelectMenu(*raw.Alias)
	return nil
}

// https://docs.discord.com/developers/components/reference#string-select-select-option-structure
type StringSelectMenuComponentOption struct {
	Label       string         `json:"label"`
	Value       string         `json:"value"`
	Description string         `json:"description,omitempty"`
	Emoji       *discord.Emoji `json:"emoji,omitempty"`
	Default     bool           `json:"default"`
}

func (s *StringSelectMenu) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#string-select
type StringSelectComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []string              `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	ComponentType discord.ComponentType `json:"component_type"`
}

func (s *StringSelectComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (s *StringSelectComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias StringSelectComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type          discord.ComponentType `json:"type"`
		ComponentType discord.ComponentType `json:"component_type"`
	}{
		Alias:         Alias(*s),
		Type:          discord.ComponentTypeStringSelect,
		ComponentType: discord.ComponentTypeStringSelect,
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

	if raw.Alias == nil {
		return nil
	}
	*s = StringSelectComponentInteractionResponse(*raw.Alias)
	return nil
}
