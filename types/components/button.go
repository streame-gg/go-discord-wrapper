package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ButtonStyle int

const (
	ButtonStylePrimary   ButtonStyle = 1
	ButtonStyleSecondary ButtonStyle = 2
	ButtonStyleSuccess   ButtonStyle = 3
	ButtonStyleDanger    ButtonStyle = 4
	ButtonStyleLink      ButtonStyle = 5
	ButtonStylePremium   ButtonStyle = 6
)

type ButtonComponent struct {
	Type     discord.ComponentType `json:"type"`
	ID       *int                  `json:"id,omitempty"`
	Style    ButtonStyle           `json:"style"`
	Label    string                `json:"label,omitempty"`
	Emoji    *discord.Emoji        `json:"emoji,omitempty"`
	CustomID string                `json:"custom_id,omitempty"`
	SkuID    *discord.Snowflake    `json:"sku_id,omitempty"`
	URL      string                `json:"url,omitempty"`
	Disabled bool                  `json:"disabled,omitempty"`
}

func (b *ButtonComponent) UnmarshalJSON(data []byte) error {
	type Alias ButtonComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*b = ButtonComponent(*raw.Alias)
	return nil
}

func (b *ButtonComponent) MarshalJSON() ([]byte, error) {
	b.Type = discord.ComponentTypeButton
	type Alias ButtonComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}

func (b *ButtonComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeButton
}

func (b *ButtonComponent) IsAnySectionAccessory() {}

func (b *ButtonComponent) IsAnyContainerAccessory() {
}
