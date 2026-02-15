package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type TextInputStyle int

const (
	TextInputStyleShort     TextInputStyle = 1
	TextInputStyleParagraph TextInputStyle = 2
)

type TextInputComponent struct {
	Type        common.ComponentType `json:"type"`
	ID          *int                 `json:"id,omitempty"`
	CustomID    string               `json:"custom_id"`
	Style       TextInputStyle       `json:"style"`
	MinLength   *int                 `json:"min_length,omitempty"`
	MaxLength   *int                 `json:"max_length,omitempty"`
	Required    *bool                `json:"required,omitempty"`
	Value       string               `json:"value,omitempty"`
	Placeholder string               `json:"placeholder,omitempty"`
}

func (t *TextInputComponent) MarshalJSON() ([]byte, error) {
	t.Type = common.ComponentTypeTextInput
	type Alias TextInputComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextInputComponent) UnmarshalJSON(data []byte) error {
	type Alias TextInputComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextInputComponent(*raw.Alias)
	return nil
}

func (t *TextInputComponent) GetType() common.ComponentType {
	return common.ComponentTypeTextInput
}

func (t *TextInputComponent) IsAnyContainerComponent() {

}

func (t *TextInputComponent) IsAnyLabelComponent() {

}

type TextInputComponentInteractionResponse struct {
	Type     common.ComponentType `json:"type"`
	Value    string               `json:"value"`
	ID       *int                 `json:"id,omitempty"`
	CustomID string               `json:"custom_id"`
}

func (t *TextInputComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (t *TextInputComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	t.Type = common.ComponentTypeTextInput

	type Alias TextInputComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextInputComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias TextInputComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextInputComponentInteractionResponse(*raw.Alias)
	return nil
}
