package builder

import (
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ButtonBuilder builds a Button using a fluent API.
//
//	btn := builder.NewButton().
//	    SetStyle(components.ButtonStylePrimary).
//	    SetLabel("Click me").
//	    SetCustomID("my_button").
//	    Build()
//
// https://docs.discord.com/developers/components/reference#button
type ButtonBuilder struct {
	btn components.Button
}

func NewButton() *ButtonBuilder {
	return &ButtonBuilder{
		btn: components.Button{
			Type: discord.ComponentTypeButton,
		},
	}
}

func (b *ButtonBuilder) SetStyle(style components.ButtonStyle) *ButtonBuilder {
	b.btn.Style = style
	return b
}

func (b *ButtonBuilder) SetLabel(label string) *ButtonBuilder {
	b.btn.Label = label
	return b
}

func (b *ButtonBuilder) SetCustomID(id string) *ButtonBuilder {
	b.btn.CustomID = id
	return b
}

// SetURL sets the URL for link-style buttons (ButtonStyleLink).
func (b *ButtonBuilder) SetURL(url string) *ButtonBuilder {
	b.btn.URL = url
	return b
}

func (b *ButtonBuilder) SetEmoji(emoji *discord.Emoji) *ButtonBuilder {
	b.btn.Emoji = emoji
	return b
}

// SetSKUID sets the SKU ID for premium-style buttons (ButtonStylePremium).
func (b *ButtonBuilder) SetSKUID(id discord.Snowflake) *ButtonBuilder {
	b.btn.SkuID = id
	return b
}

func (b *ButtonBuilder) SetDisabled(disabled bool) *ButtonBuilder {
	b.btn.Disabled = disabled
	return b
}

func (b *ButtonBuilder) Build() *components.Button {
	return &b.btn
}
