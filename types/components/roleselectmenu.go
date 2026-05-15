package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type RoleSelectMenuComponent struct {
	Type          discord.ComponentType `json:"type"`
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
	r.Type = discord.ComponentTypeRoleSelectMenu
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

func (r *RoleSelectMenuComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeRoleSelectMenu
}

func (r *RoleSelectMenuComponent) IsAnyLabelComponent() {

}

type RoleComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []discord.Snowflake   `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id,omitempty"`
	ComponentType discord.ComponentType `json:"component_type"`
	Resolved      *discord.ResolvedData `json:"resolved,omitempty"`
}

func (r *RoleComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (r *RoleComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	r.ComponentType = discord.ComponentTypeRoleSelectMenu
	r.Type = discord.ComponentTypeRoleSelectMenu

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
