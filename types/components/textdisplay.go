package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#text-display
type TextDisplay struct {
	Type    discord.ComponentType `json:"type"`
	ID      *int                  `json:"id,omitempty"`
	Content string                `json:"content"`
}

func (t *TextDisplay) UnmarshalJSON(data []byte) error {
	type Alias TextDisplay
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*t = TextDisplay(*raw.Alias)
	return nil
}

func (t *TextDisplay) IsAnyContainerComponent() {}

func (t *TextDisplay) GetType() discord.ComponentType {
	return discord.ComponentTypeTextDisplay
}

func (t *TextDisplay) MarshalJSON() ([]byte, error) {
	type Alias TextDisplay
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*t),
		Type:  discord.ComponentTypeTextDisplay,
	})
}

func (t *TextDisplay) IsAnySectionComponent() {}

// https://docs.discord.com/developers/components/reference#text-display
type TextDisplayComponentInteractionResponse struct {
	Type discord.ComponentType `json:"type"`
	ID   *int                  `json:"id,omitempty"`
}

func (t *TextDisplayComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (t *TextDisplayComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias TextDisplayComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*t),
		Type:  discord.ComponentTypeTextDisplay,
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
