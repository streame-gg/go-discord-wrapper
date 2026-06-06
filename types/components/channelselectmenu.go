package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#channel-select
type ChannelSelectMenuComponent struct {
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

func (c *ChannelSelectMenuComponent) IsAnyContainerAccessory() bool {
	return true
}

func (c *ChannelSelectMenuComponent) UnmarshalJSON(data []byte) error {
	type Alias ChannelSelectMenuComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = ChannelSelectMenuComponent(*raw.Alias)
	return nil
}

func (c *ChannelSelectMenuComponent) MarshalJSON() ([]byte, error) {
	type Alias ChannelSelectMenuComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*c),
		Type:  discord.ComponentTypeChannelSelect,
	})
}

func (c *ChannelSelectMenuComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeChannelSelect
}

func (c *ChannelSelectMenuComponent) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#channel-select
type ChannelComponentInteractionResponse struct {
	Type          discord.ComponentType `json:"type"`
	Values        []discord.Snowflake   `json:"values"`
	ID            *int                  `json:"id,omitempty"`
	CustomID      string                `json:"custom_id,omitempty"`
	ComponentType discord.ComponentType `json:"component_type"`
	Resolved      *discord.ResolvedData `json:"resolved,omitempty"`
}

func (c *ChannelComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (c *ChannelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias ChannelComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type          discord.ComponentType `json:"type"`
		ComponentType discord.ComponentType `json:"component_type"`
	}{
		Alias:         Alias(*c),
		Type:          discord.ComponentTypeChannelSelect,
		ComponentType: discord.ComponentTypeChannelSelect,
	})
}

func (c *ChannelComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias ChannelComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*c = ChannelComponentInteractionResponse(*raw.Alias)
	return nil
}
