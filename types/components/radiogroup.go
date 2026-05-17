package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type RadioGroupComponent struct {
	Type     discord.ComponentType        `json:"type"`
	ID       *int                         `json:"id,omitempty"`
	CustomID string                       `json:"custom_id"`
	Options  *[]RadioGroupComponentOption `json:"options"`
	Required *bool                        `json:"required,omitempty"`
}

type RadioGroupComponentOption struct {
	Value       string  `json:"value"`
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	Default     *bool   `json:"default,omitempty"`
}

func (r *RadioGroupComponent) MarshalJSON() ([]byte, error) {
	r.Type = discord.ComponentTypeRadioGroup
	type Alias RadioGroupComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RadioGroupComponent) UnmarshalJSON(data []byte) error {
	type Alias RadioGroupComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*r = RadioGroupComponent(*raw.Alias)
	return nil
}

func (r *RadioGroupComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeRadioGroup
}

func (r *RadioGroupComponent) IsAnyLabelComponent() {

}

type RadioGroupComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id,omitempty"`
	Value    *string               `json:"value,omitempty"`
}

func (r *RadioGroupComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (r *RadioGroupComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	r.Type = discord.ComponentTypeRadioGroup

	type Alias RadioGroupComponentInteractionResponse
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
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
