package components

import (
	"encoding/json"
	"errors"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type Container struct {
	Type        common.ComponentType     `json:"type"`
	ID          *int                     `json:"id,omitempty"`
	Components  *[]AnyContainerComponent `json:"components"`
	AccentColor int                      `json:"accent_color,omitempty"`
	Spoiler     bool                     `json:"spoiler,omitempty"`
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

	*c = Container(*raw.Alias)

	for _, comp := range raw.Components {
		var probe struct {
			Type common.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(comp, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case common.ComponentTypeMediaGallery:
			var m *MediaGalleryComponent
			if err := json.Unmarshal(comp, &m); err != nil {
				return err
			}
			*c.Components = append(*c.Components, m)
		case common.ComponentTypeFileDisplay:
			var f *FileComponent
			if err := json.Unmarshal(comp, &f); err != nil {
				return err
			}
			*c.Components = append(*c.Components, f)
		case common.ComponentTypeSeparator:
			var s *SeparatorComponent
			if err := json.Unmarshal(comp, &s); err != nil {
				return err
			}
			*c.Components = append(*c.Components, s)
		case common.ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(comp, &t); err != nil {
				return err
			}
			*c.Components = append(*c.Components, t)
		case common.ComponentTypeActionRow:
			var a *ActionRow
			if err := json.Unmarshal(comp, &a); err != nil {
				return err
			}
			*c.Components = append(*c.Components, a)
		case common.ComponentTypeTextDisplay:
			var t *TextDisplayComponent
			if err := json.Unmarshal(comp, &t); err != nil {
				return err
			}
			*c.Components = append(*c.Components, t)
		case common.ComponentTypeSection:
			var s *Section
			if err := json.Unmarshal(comp, &s); err != nil {
				return err
			}
			*c.Components = append(*c.Components, s)
		default:
			return errors.New("unknown container component type" + string(rune(probe.Type)))
		}
	}

	return nil
}

func (c *Container) GetType() common.ComponentType {
	return common.ComponentTypeContainer
}

func (c *Container) MarshalJSON() ([]byte, error) {
	c.Type = common.ComponentTypeContainer
	type Alias Container
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}
