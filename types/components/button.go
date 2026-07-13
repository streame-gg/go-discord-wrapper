package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#button-button-styles
type ButtonStyle int

const (
	ButtonStylePrimary   ButtonStyle = 1
	ButtonStyleSecondary ButtonStyle = 2
	ButtonStyleSuccess   ButtonStyle = 3
	ButtonStyleDanger    ButtonStyle = 4
	ButtonStyleLink      ButtonStyle = 5
	ButtonStylePremium   ButtonStyle = 6
)

// https://docs.discord.com/developers/components/reference#button
type Button struct {
	Type     discord.ComponentType `json:"type"`
	ID       int                   `json:"id,omitempty"`
	Style    ButtonStyle           `json:"style"`
	Label    string                `json:"label,omitempty"`
	Emoji    *discord.Emoji        `json:"emoji,omitempty"`
	CustomID string                `json:"custom_id,omitempty"`
	SkuID    discord.Snowflake     `json:"sku_id,omitempty"`
	URL      string                `json:"url,omitempty"`
	Disabled bool                  `json:"disabled"`
}

func (b *Button) UnmarshalJSON(data []byte) error {
	type Alias Button
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*b = Button(*raw.Alias)
	return nil
}

func (b *Button) MarshalJSON() ([]byte, error) {
	type Alias Button
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*b),
		Type:  discord.ComponentTypeButton,
	})
}

func (b *Button) GetType() discord.ComponentType {
	return discord.ComponentTypeButton
}

func (b *Button) IsAnySectionAccessory() {}

func (b *Button) IsAnyContainerAccessory() {
}
