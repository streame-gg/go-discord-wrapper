package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ChannelSelectMenuComponent struct {
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

	*c = ChannelSelectMenuComponent(*raw.Alias)
	return nil
}

func (c *ChannelSelectMenuComponent) MarshalJSON() ([]byte, error) {
	c.Type = common.ComponentTypeChannelSelect
	type Alias ChannelSelectMenuComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

func (c *ChannelSelectMenuComponent) GetType() common.ComponentType {
	return common.ComponentTypeChannelSelect
}

func (c *ChannelSelectMenuComponent) IsAnyLabelComponent() {

}

type ChannelComponentInteractionResponse struct {
	Type          common.ComponentType `json:"type"`
	Values        []common.Snowflake   `json:"values"`
	ID            *int                 `json:"id,omitempty"`
	CustomID      string               `json:"custom_id,omitempty"`
	ComponentType common.ComponentType `json:"component_type"`
	Resolved      *common.ResolvedData `json:"resolved,omitempty"`
}

func (c *ChannelComponentInteractionResponse) IsInteractionResponseDataComponent() {}

func (c *ChannelComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	c.ComponentType = common.ComponentTypeChannelSelect
	c.Type = common.ComponentTypeChannelSelect

	type Alias ChannelComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
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

	*c = ChannelComponentInteractionResponse(*raw.Alias)
	return nil
}
