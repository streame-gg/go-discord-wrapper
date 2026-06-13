package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#separator-spacing-size
type SeparatorComponentSpacing int

const (
	SeparatorComponentSpacingSmall SeparatorComponentSpacing = 1
	SeparatorComponentSpacingLarge SeparatorComponentSpacing = 2
)

// https://docs.discord.com/developers/components/reference#separator
type SeparatorComponent struct {
	Type                      discord.ComponentType     `json:"type"`
	ID                        *int                      `json:"id,omitempty"`
	Divider                   bool                      `json:"divider"`
	SeparatorComponentSpacing SeparatorComponentSpacing `json:"spacing,omitempty"`
}

func (s *SeparatorComponent) UnmarshalJSON(data []byte) error {
	type Alias SeparatorComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*s = SeparatorComponent(*raw.Alias)
	return nil
}

func (s *SeparatorComponent) MarshalJSON() ([]byte, error) {
	type Alias SeparatorComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*s),
		Type:  discord.ComponentTypeSeparator,
	})
}

func (s *SeparatorComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeSeparator
}

func (s *SeparatorComponent) IsAnyContainerComponent() {}
