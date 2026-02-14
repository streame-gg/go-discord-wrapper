package components

import (
	"encoding/json"
	"go-discord-wrapper/types/common"
)

type CheckboxGroupComponent struct {
	Type      common.ComponentType            `json:"type"`
	ID        *int                            `json:"id,omitempty"`
	CustomID  string                          `json:"custom_id"`
	Options   *[]CheckboxGroupComponentOption `json:"options"`
	MinValues *int                            `json:"min_values,omitempty"`
	MaxValues *int                            `json:"max_values,omitempty"`
	Required  *bool                           `json:"required,omitempty"`
}

type CheckboxGroupComponentOption struct {
	Value       string  `json:"value"`
	Label       string  `json:"label"`
	Description *string `json:"description,omitempty"`
	Default     *bool   `json:"default,omitempty"`
}

func (c *CheckboxGroupComponent) MarshalJSON() ([]byte, error) {
	c.Type = common.ComponentTypeCheckboxGroup
	type Alias CheckboxGroupComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *CheckboxGroupComponent) UnmarshalJSON(data []byte) error {
	type Alias CheckboxGroupComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*c = CheckboxGroupComponent(*raw.Alias)
	return nil
}

func (c *CheckboxGroupComponent) GetType() common.ComponentType {
	return common.ComponentTypeCheckboxGroup
}

func (c *CheckboxGroupComponent) IsAnyLabelComponent() bool {
	return true
}
