package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#user-select
type UserSelectMenuComponent struct {
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

func (u *UserSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (u *UserSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type Alias UserSelectMenuComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*u),
		Type:  discord.ComponentTypeUserSelect,
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

	if raw.Alias == nil {
		return nil
	}
	*u = UserSelectMenuComponent(*raw.Alias)
	return nil
}

func (u *UserSelectMenuComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeUserSelect
}

func (u *UserSelectMenuComponent) IsAnyLabelComponent() {
}

// https://docs.discord.com/developers/components/reference#user-select
type UserSelectComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []discord.Snowflake   `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id,omitempty"`
	ComponentType discord.ComponentType `json:"component_type"`
	Resolved      *discord.ResolvedData `json:"resolved,omitempty"`
}

func (u *UserSelectComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (u *UserSelectComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias UserSelectComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type          discord.ComponentType `json:"type"`
		ComponentType discord.ComponentType `json:"component_type"`
	}{
		Alias:         Alias(*u),
		Type:          discord.ComponentTypeUserSelect,
		ComponentType: discord.ComponentTypeUserSelect,
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

	if raw.Alias == nil {
		return nil
	}
	*u = UserSelectComponentInteractionResponse(*raw.Alias)
	return nil
}
