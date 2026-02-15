package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type RadioGroupComponent struct {
	Type     common.ComponentType         `json:"type"`
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
	r.Type = common.ComponentTypeRadioGroup
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

	*r = RadioGroupComponent(*raw.Alias)
	return nil
}

func (r *RadioGroupComponent) GetType() common.ComponentType {
	return common.ComponentTypeRadioGroup
}

func (r *RadioGroupComponent) IsAnyLabelComponent() {

}

type RadioGroupComponentInteractionResponse struct {
	Type     common.ComponentType `json:"type"`
	ID       *int                 `json:"id,omitempty"`
	CustomID string               `json:"custom_id,omitempty"`
	Value    *string              `json:"value"`
}

func (r *RadioGroupComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (r *RadioGroupComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	r.Type = common.ComponentTypeRadioGroup

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

	*r = RadioGroupComponentInteractionResponse(*raw.Alias)
	return nil
}
