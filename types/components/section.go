package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type Section struct {
	Type       common.ComponentType   `json:"type"`
	ID         *int                   `json:"id,omitempty"`
	Components *[]AnySectionComponent `json:"components"`
	Accessory  AnySectionAccessory    `json:"accessory,omitempty"`
}

func (s *Section) IsAnyContainerComponent() {

}

func (s *Section) UnmarshalJSON(data []byte) error {
	type Alias Section

	var raw struct {
		Alias
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = Section(raw.Alias)

	for _, c := range raw.Components {
		var probe struct {
			Type common.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case common.ComponentTypeTextDisplay:
			var t *TextDisplayComponent
			if err := json.Unmarshal(c, &t); err != nil {
				return err
			}
			*s.Components = append(*s.Components, t)
		}
	}

	return nil
}

func (s *Section) MarshalJSON() ([]byte, error) {
	s.Type = common.ComponentTypeSection
	type Alias Section
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *Section) GetType() common.ComponentType {
	return common.ComponentTypeSection
}
