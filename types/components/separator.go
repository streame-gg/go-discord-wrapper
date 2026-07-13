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
type Separator struct {
	Type    discord.ComponentType     `json:"type"`
	ID      *int                      `json:"id,omitempty"`
	Divider bool                      `json:"divider"`
	Spacing SeparatorComponentSpacing `json:"spacing,omitempty"`
}

func (s *Separator) UnmarshalJSON(data []byte) error {
	type Alias Separator
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*s = Separator(*raw.Alias)
	return nil
}

func (s *Separator) MarshalJSON() ([]byte, error) {
	type Alias Separator
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*s),
		Type:  discord.ComponentTypeSeparator,
	})
}

func (s *Separator) GetType() discord.ComponentType {
	return discord.ComponentTypeSeparator
}

func (s *Separator) IsAnyContainerComponent() {}
