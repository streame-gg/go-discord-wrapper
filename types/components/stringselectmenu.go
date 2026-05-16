package components

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type StringSelectMenuComponent struct {
	Type        discord.ComponentType              `json:"type"`
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
	s.Type = discord.ComponentTypeStringSelectMenu
	type Alias StringSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *StringSelectMenuComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeStringSelectMenu
}

func (s *StringSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias StringSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*s = StringSelectMenuComponent(*raw.Alias)
	return nil
}

type StringSelectMenuComponentOption struct {
	Label       string         `json:"label"`
	Value       string         `json:"value"`
	Description string         `json:"description,omitempty"`
	Emoji       *discord.Emoji `json:"emoji,omitempty"`
	Default     bool           `json:"default,omitempty"`
}

func (s *StringSelectMenuComponent) IsAnyLabelComponent() {

}

type StringSelectComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []string              `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id,omitempty"`
	ComponentType discord.ComponentType `json:"component_type"`
}

func (s *StringSelectComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (s *StringSelectComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	s.ComponentType = discord.ComponentTypeStringSelectMenu
	s.Type = discord.ComponentTypeStringSelectMenu

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

	if raw.Alias == nil {
		return nil
	}
	*s = StringSelectComponentInteractionResponse(*raw.Alias)
	return nil
}
