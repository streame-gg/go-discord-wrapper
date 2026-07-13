package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#checkbox
type Checkbox struct {
	Type     discord.ComponentType `json:"type"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id"`
	Default  bool                  `json:"default"`
}

func (c *Checkbox) MarshalJSON() ([]byte, error) {
	type Alias Checkbox
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*c),
		Type:  discord.ComponentTypeCheckbox,
	})
}

func (c *Checkbox) UnmarshalJSON(data []byte) error {
	type Alias Checkbox
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = Checkbox(*raw.Alias)
	return nil
}

func (c *Checkbox) GetType() discord.ComponentType {
	return discord.ComponentTypeCheckbox
}

func (c *Checkbox) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#checkbox
type CheckboxComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	Value    bool                  `json:"value"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id,omitempty"`
}

func (c *CheckboxComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (c *CheckboxComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias CheckboxComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*c),
		Type:  discord.ComponentTypeCheckbox,
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
