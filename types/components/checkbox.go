package components

import (
	"encoding/json"
	"go-discord-wrapper/types/common"
)

type CheckboxComponent struct {
	Type     common.ComponentType `json:"type"`
	ID       *int                 `json:"id,omitempty"`
	CustomID string               `json:"custom_id"`
	Default  *bool                `json:"default,omitempty"`
}

func (c *CheckboxComponent) MarshalJSON() ([]byte, error) {
	c.Type = common.ComponentTypeCheckbox
	type Alias CheckboxComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *CheckboxComponent) UnmarshalJSON(data []byte) error {
	type Alias CheckboxComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = CheckboxComponent(*raw.Alias)
	return nil
}

func (c *CheckboxComponent) GetType() common.ComponentType {
	return common.ComponentTypeCheckbox
}

func (c *CheckboxComponent) IsAnyLabelComponent() bool {
	return true
}
