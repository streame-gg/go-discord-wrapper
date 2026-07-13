package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#radio-group
type RadioGroup struct {
	Type     discord.ComponentType       `json:"type"`
	ID       *int                        `json:"id,omitempty"`
	CustomID string                      `json:"custom_id"`
	Options  []RadioGroupComponentOption `json:"options"`
	Required bool                        `json:"required"`
}

// https://docs.discord.com/developers/components/reference#radio-group-option-structure
type RadioGroupComponentOption struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default"`
}

func (r *RadioGroup) MarshalJSON() ([]byte, error) {
	type Alias RadioGroup
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*r),
		Type:  discord.ComponentTypeRadioGroup,
	})
}

func (r *RadioGroup) UnmarshalJSON(data []byte) error {
	type Alias RadioGroup
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*r = RadioGroup(*raw.Alias)
	return nil
}

func (r *RadioGroup) GetType() discord.ComponentType {
	return discord.ComponentTypeRadioGroup
}

func (r *RadioGroup) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#radio-group
type RadioGroupComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id"`
	Value    *string               `json:"value"`
}

func (r *RadioGroupComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (r *RadioGroupComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias RadioGroupComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*r),
		Type:  discord.ComponentTypeRadioGroup,
	})
}

func (r *RadioGroupComponentInteractionResponse) UnmarshalJSON(bytes []byte) error {
	type Alias RadioGroupComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*r = RadioGroupComponentInteractionResponse(*raw.Alias)
	return nil
}
