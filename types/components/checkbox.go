package components

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type CheckboxComponent struct {
	Type     discord.ComponentType `json:"type"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id"`
	Default  *bool                 `json:"default,omitempty"`
}

func (c *CheckboxComponent) MarshalJSON() ([]byte, error) {
	c.Type = discord.ComponentTypeCheckbox
	type Alias CheckboxComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *CheckboxComponent) UnmarshalJSON(data []byte) error {
	type Alias CheckboxComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = CheckboxComponent(*raw.Alias)
	return nil
}

func (c *CheckboxComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeCheckbox
}

func (c *CheckboxComponent) IsAnyLabelComponent() {

}

type CheckboxComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	Value    bool                  `json:"value"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id,omitempty"`
}

func (c *CheckboxComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (c *CheckboxComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	c.Type = discord.ComponentTypeCheckbox
	type Alias CheckboxComponentInteractionResponse
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *CheckboxComponentInteractionResponse) UnmarshalJSON(bytes []byte) error {
	type Alias CheckboxComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = CheckboxComponentInteractionResponse(*raw.Alias)
	return nil
}
