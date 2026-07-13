package components

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#section
type Section struct {
	Type       discord.ComponentType `json:"type"`
	ID         *int                  `json:"id,omitempty"`
	Components []AnySectionComponent `json:"components"`
	Accessory  AnySectionAccessory   `json:"accessory,omitempty"`
}

func (s *Section) IsAnyContainerComponent() {}

func (s *Section) UnmarshalJSON(data []byte) error {
	type Alias Section

	var raw struct {
		Alias
		Accessory  json.RawMessage   `json:"accessory"`
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = Section(raw.Alias)

	var probe struct {
		Type discord.ComponentType `json:"type"`
	}

	if err := json.Unmarshal(raw.Accessory, &probe); err != nil {
		return err
	}

	switch probe.Type {
	case discord.ComponentTypeButton:
		var t *Button
		if err := json.Unmarshal(raw.Accessory, &t); err != nil {
			return err
		}
		s.Accessory = t
	case discord.ComponentTypeThumbnail:
		var t *Thumbnail
		if err := json.Unmarshal(raw.Accessory, &t); err != nil {
			return err
		}
		s.Accessory = t
	default:
		return fmt.Errorf("unknown section component type: %d", probe.Type)
	}

	for _, c := range raw.Components {
		var probe struct {
			Type discord.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case discord.ComponentTypeTextDisplay:
			var t *TextDisplay
			if err := json.Unmarshal(c, &t); err != nil {
				return err
			}
			s.Components = append(s.Components, t)
		default:
			return fmt.Errorf("unknown section component type: %d", probe.Type)
		}
	}

	return nil
}

func (s *Section) MarshalJSON() ([]byte, error) {
	type Alias Section
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*s),
		Type:  discord.ComponentTypeSection,
	})
}

func (s *Section) GetType() discord.ComponentType {
	return discord.ComponentTypeSection
}
