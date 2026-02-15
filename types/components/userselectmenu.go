package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type UserSelectMenuComponent struct {
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

func (u *UserSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (u *UserSelectMenuComponent) MarshalJSON() ([]byte, error) {
	u.Type = common.ComponentTypeUserSelectMenu
	type Alias UserSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

func (u *UserSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias UserSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*u = UserSelectMenuComponent(*raw.Alias)
	return nil
}

func (u *UserSelectMenuComponent) GetType() common.ComponentType {
	return common.ComponentTypeUserSelectMenu
}

func (u *UserSelectMenuComponent) IsAnyLabelComponent() {
}

type UserSelectComponentInteractionResponse struct {
	Type          common.ComponentType `json:"type"`
	Values        []common.Snowflake   `json:"values"`
	ID            *int                 `json:"id,omitempty"`
	CustomID      string               `json:"custom_id,omitempty"`
	ComponentType common.ComponentType `json:"component_type"`
	Resolved      *common.ResolvedData `json:"resolved,omitempty"`
}

func (u *UserSelectComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (u *UserSelectComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	u.ComponentType = common.ComponentTypeRoleSelectMenu
	u.Type = common.ComponentTypeRoleSelectMenu

	type Alias UserSelectComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

func (u *UserSelectComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias UserSelectComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*u = UserSelectComponentInteractionResponse(*raw.Alias)
	return nil
}
