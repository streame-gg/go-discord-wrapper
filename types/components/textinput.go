package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#text-input-text-input-styles
type TextInputStyle int

const (
	TextInputStyleShort     TextInputStyle = 1
	TextInputStyleParagraph TextInputStyle = 2
)

// https://docs.discord.com/developers/components/reference#text-input
type TextInput struct {
	Type        discord.ComponentType `json:"type"`
	ID          *int                  `json:"id,omitempty"`
	CustomID    string                `json:"custom_id"`
	Style       TextInputStyle        `json:"style"`
	MinLength   *int                  `json:"min_length,omitempty"`
	MaxLength   *int                  `json:"max_length,omitempty"`
	Required    bool                  `json:"required"`
	Value       string                `json:"value,omitempty"`
	Placeholder string                `json:"placeholder,omitempty"`
}

func (t *TextInput) MarshalJSON() ([]byte, error) {
	type Alias TextInput
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*t),
		Type:  discord.ComponentTypeTextInput,
	})
}

func (t *TextInput) UnmarshalJSON(data []byte) error {
	type Alias TextInput
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*t = TextInput(*raw.Alias)
	return nil
}

func (t *TextInput) GetType() discord.ComponentType {
	return discord.ComponentTypeTextInput
}

func (t *TextInput) IsAnyContainerComponent() {

}

func (t *TextInput) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#text-input
type TextInputComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	Value    string                `json:"value"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id"`
}

func (t *TextInputComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (t *TextInputComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias TextInputComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*t),
		Type:  discord.ComponentTypeTextInput,
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

	if raw.Alias == nil {
		return nil
	}
	*t = TextInputComponentInteractionResponse(*raw.Alias)
	return nil
}
