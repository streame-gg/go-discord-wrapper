package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#role-select
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
	type Alias RoleSelectMenuComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*r),
		Type:  discord.ComponentTypeRoleSelect,
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

	if raw.Alias == nil {
		return nil
	}
	*r = RoleSelectMenuComponent(*raw.Alias)
	return nil
}

func (r *RoleSelectMenuComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeRoleSelect
}

func (r *RoleSelectMenuComponent) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#role-select
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
	type Alias RoleComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type          discord.ComponentType `json:"type"`
		ComponentType discord.ComponentType `json:"component_type"`
	}{
		Alias:         Alias(*r),
		Type:          discord.ComponentTypeRoleSelect,
		ComponentType: discord.ComponentTypeRoleSelect,
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

	if raw.Alias == nil {
		return nil
	}
	*r = RoleComponentInteractionResponse(*raw.Alias)
	return nil
}
