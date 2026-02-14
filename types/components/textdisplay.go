package components

import (
	"encoding/json"
	"go-discord-wrapper/types/common"
)

type TextDisplayComponent struct {
	Type    common.ComponentType `json:"type"`
	ID      *int                 `json:"id,omitempty"`
	Content string               `json:"content"`
}

func (t *TextDisplayComponent) UnmarshalJSON(data []byte) error {
	type Alias TextDisplayComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextDisplayComponent(*raw.Alias)
	return nil
}

func (t *TextDisplayComponent) IsAnyContainerComponent() bool {
	return true
}

func (t *TextDisplayComponent) GetType() common.ComponentType {
	return common.ComponentTypeTextDisplay
}

func (t *TextDisplayComponent) MarshalJSON() ([]byte, error) {
	t.Type = common.ComponentTypeTextDisplay
	type Alias TextDisplayComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextDisplayComponent) IsAnySectionComponent() bool {
	return true
}

type TextDisplayComponentInteractionResponse struct {
	Type common.ComponentType `json:"type"`
	ID   *int                 `json:"id,omitempty"`
}

func (t *TextDisplayComponentInteractionResponse) IsInteractionResponseDataComponent() bool {
	return true
}

func (t *TextDisplayComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	t.Type = common.ComponentTypeTextDisplay

	type Alias TextDisplayComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextDisplayComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias TextDisplayComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = TextDisplayComponentInteractionResponse(*raw.Alias)
	return nil
}
