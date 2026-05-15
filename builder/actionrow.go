package builder

import (
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ActionRowBuilder builds an ActionRow using a fluent API.
//
//	row := builder.NewActionRow().
//	    AddComponents(btn1, btn2).
//	    Build()
type ActionRowBuilder struct {
	row components.ActionRow
}

func NewActionRow() *ActionRowBuilder { return &ActionRowBuilder{} }

func (b *ActionRowBuilder) AddComponents(c ...discord.AnyComponent) *ActionRowBuilder {
	b.row.Components = append(b.row.Components, c...)
	return b
}

func (b *ActionRowBuilder) Build() *components.ActionRow {
	return &b.row
}
