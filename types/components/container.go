package components

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type Container struct {
	Type        discord.ComponentType   `json:"type"`
	ID          *int                    `json:"id,omitempty"`
	Components  []AnyContainerComponent `json:"components"`
	AccentColor int                     `json:"accent_color,omitempty"`
	Spoiler     bool                    `json:"spoiler,omitempty"`
}

func (c *Container) UnmarshalJSON(data []byte) error {
	type Alias Container

	var raw struct {
		*Alias
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = Container(*raw.Alias)

	for _, comp := range raw.Components {
		var probe struct {
			Type discord.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(comp, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case discord.ComponentTypeMediaGallery:
			var m *MediaGalleryComponent
			if err := json.Unmarshal(comp, &m); err != nil {
				return err
			}
			c.Components = append(c.Components, m)
		case discord.ComponentTypeFileDisplay:
			var f *FileComponent
			if err := json.Unmarshal(comp, &f); err != nil {
				return err
			}
			c.Components = append(c.Components, f)
		case discord.ComponentTypeSeparator:
			var s *SeparatorComponent
			if err := json.Unmarshal(comp, &s); err != nil {
				return err
			}
			c.Components = append(c.Components, s)
		case discord.ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(comp, &t); err != nil {
				return err
			}
			c.Components = append(c.Components, t)
		case discord.ComponentTypeActionRow:
			var a *ActionRow
			if err := json.Unmarshal(comp, &a); err != nil {
				return err
			}
			c.Components = append(c.Components, a)
		case discord.ComponentTypeTextDisplay:
			var t *TextDisplayComponent
			if err := json.Unmarshal(comp, &t); err != nil {
				return err
			}
			c.Components = append(c.Components, t)
		case discord.ComponentTypeSection:
			var s *Section
			if err := json.Unmarshal(comp, &s); err != nil {
				return err
			}
			c.Components = append(c.Components, s)
		default:
			return fmt.Errorf("unknown container component type: %d", probe.Type)
		}
	}

	return nil
}

func (c *Container) GetType() discord.ComponentType {
	return discord.ComponentTypeContainer
}

func (c *Container) MarshalJSON() ([]byte, error) {
	c.Type = discord.ComponentTypeContainer
	type Alias Container
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}
