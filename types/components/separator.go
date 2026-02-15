package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type SeparatorComponentSpacing int

const (
	SeparatorComponentSpacingSmall SeparatorComponentSpacing = 1
	SeparatorComponentSpacingLarge SeparatorComponentSpacing = 2
)

type SeparatorComponent struct {
	Type                      common.ComponentType      `json:"type"`
	ID                        *int                      `json:"id,omitempty"`
	Divider                   bool                      `json:"divider,omitempty"`
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

	*s = SeparatorComponent(*raw.Alias)
	return nil
}

func (s *SeparatorComponent) MarshalJSON() ([]byte, error) {
	s.Type = common.ComponentTypeSeparator
	type Alias SeparatorComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *SeparatorComponent) GetType() common.ComponentType {
	return common.ComponentTypeSeparator
}

func (s *SeparatorComponent) IsAnyContainerComponent() {}
