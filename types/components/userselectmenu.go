package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#user-select
type UserSelectMenu struct {
	Type          discord.ComponentType `json:"type"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	Placeholder   string                `json:"placeholder,omitempty"`
	MinValues     *int                  `json:"min_values,omitempty"`
	MaxValues     *int                  `json:"max_values,omitempty"`
	Required      bool                  `json:"required"`
	Disabled      bool                  `json:"disabled"`
	DefaultValues []SelectDefaultValue  `json:"default_values,omitempty"`
}

func (u *UserSelectMenu) IsAnyContainerAccessory() bool {
	return true
}

func (u *UserSelectMenu) MarshalJSON() ([]byte, error) {
	type Alias UserSelectMenu
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*u),
		Type:  discord.ComponentTypeUserSelect,
	})
}

func (u *UserSelectMenu) UnmarshalJSON(data []byte) error {
	type Alias UserSelectMenu
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*u = UserSelectMenu(*raw.Alias)
	return nil
}

func (u *UserSelectMenu) GetType() discord.ComponentType {
	return discord.ComponentTypeUserSelect
}

func (u *UserSelectMenu) IsAnyLabelComponent() {
}

// https://docs.discord.com/developers/components/reference#user-select
type UserSelectComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []discord.Snowflake   `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id"`
	ComponentType discord.ComponentType `json:"component_type"`
	Resolved      discord.ResolvedData  `json:"resolved"`
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
