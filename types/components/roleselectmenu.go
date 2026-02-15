package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type RoleSelectMenuComponent struct {
	Type          common.ComponentType  `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required,omitempty"`
	Disabled      bool                  `json:"disabled,omitempty"`
	DefaultValues *[]SelectDefaultValue `json:"default_values,omitempty"`
}

func (r *RoleSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (r *RoleSelectMenuComponent) MarshalJSON() ([]byte, error) {
	r.Type = common.ComponentTypeRoleSelectMenu
	type Alias UserSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RoleSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias RoleSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*r = RoleSelectMenuComponent(*raw.Alias)
	return nil
}

func (r *RoleSelectMenuComponent) GetType() common.ComponentType {
	return common.ComponentTypeRoleSelectMenu
}

func (r *RoleSelectMenuComponent) IsAnyLabelComponent() {

}

type RoleComponentInteractionResponse struct {
	Type          common.ComponentType `json:"type"`
	Values        []common.Snowflake   `json:"values"`
	ID            *int                 `json:"id,omitempty"`
	CustomID      string               `json:"custom_id,omitempty"`
	ComponentType common.ComponentType `json:"component_type"`
	Resolved      *common.ResolvedData `json:"resolved,omitempty"`
}

func (r *RoleComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (r *RoleComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	r.ComponentType = common.ComponentTypeRoleSelectMenu
	r.Type = common.ComponentTypeRoleSelectMenu

	type Alias RoleComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *RoleComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias RoleComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*r = RoleComponentInteractionResponse(*raw.Alias)
	return nil
}
