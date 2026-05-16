package components

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type TextDisplayComponent struct {
	Type    discord.ComponentType `json:"type"`
	ID      *int                  `json:"id,omitempty"`
	Content string                `json:"content"`
}

func (t *TextDisplayComponent) UnmarshalJSON(data []byte) error {
	type Alias TextDisplayComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*t = TextDisplayComponent(*raw.Alias)
	return nil
}

func (t *TextDisplayComponent) IsAnyContainerComponent() {}

func (t *TextDisplayComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeTextDisplay
}

func (t *TextDisplayComponent) MarshalJSON() ([]byte, error) {
	t.Type = discord.ComponentTypeTextDisplay
	type Alias TextDisplayComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *TextDisplayComponent) IsAnySectionComponent() {}

type TextDisplayComponentInteractionResponse struct {
	Type discord.ComponentType `json:"type"`
	ID   *int                  `json:"id,omitempty"`
}

func (t *TextDisplayComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (t *TextDisplayComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	t.Type = discord.ComponentTypeTextDisplay

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

	if raw.Alias == nil {
		return nil
	}
	*t = TextDisplayComponentInteractionResponse(*raw.Alias)
	return nil
}
