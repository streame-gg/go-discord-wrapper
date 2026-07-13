package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#checkbox-group
type CheckboxGroup struct {
	Type      discord.ComponentType          `json:"type"`
	ID        *int                           `json:"id,omitempty"`
	CustomID  string                         `json:"custom_id"`
	Options   []CheckboxGroupComponentOption `json:"options"`
	MinValues *int                           `json:"min_values,omitempty"`
	MaxValues *int                           `json:"max_values,omitempty"`
	Required  bool                           `json:"required"`
}

// https://docs.discord.com/developers/components/reference#checkbox-group-option-structure
type CheckboxGroupComponentOption struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default"`
}

func (c *CheckboxGroup) MarshalJSON() ([]byte, error) {
	type Alias CheckboxGroup
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*c),
		Type:  discord.ComponentTypeCheckboxGroup,
	})
}

func (c *CheckboxGroup) UnmarshalJSON(data []byte) error {
	type Alias CheckboxGroup
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = CheckboxGroup(*raw.Alias)
	return nil
}

func (c *CheckboxGroup) GetType() discord.ComponentType {
	return discord.ComponentTypeCheckboxGroup
}

func (c *CheckboxGroup) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#checkbox-group-interaction-response-structure
type CheckboxGroupComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	Values   []string              `json:"values"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id"`
}

func (c *CheckboxGroupComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (c *CheckboxGroupComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias CheckboxGroupComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*c),
		Type:  discord.ComponentTypeCheckboxGroup,
	})
}

func (c *CheckboxGroupComponentInteractionResponse) UnmarshalJSON(bytes []byte) error {
	type Alias CheckboxGroupComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = CheckboxGroupComponentInteractionResponse(*raw.Alias)
	return nil
}
