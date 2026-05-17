// Package builder provides fluent constructors for Discord UI components:
// embeds, buttons, select menus, modals, action rows, and text inputs.
package builder

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// EmbedBuilder builds a discord.Embed using a fluent API.
//
//	embed := builder.NewEmbed().
//	    SetTitle("Hello").
//	    SetDescription("World").
//	    SetColor(0x5865F2).
//	    Build()
type EmbedBuilder struct {
	embed discord.Embed
}

func NewEmbed() *EmbedBuilder { return &EmbedBuilder{} }

func (b *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	b.embed.Title = &title
	return b
}

func (b *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	b.embed.Description = &description
	return b
}

func (b *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	b.embed.URL = &url
	return b
}

func (b *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	b.embed.Color = &color
	return b
}

func (b *EmbedBuilder) SetTimestamp(t time.Time) *EmbedBuilder {
	b.embed.Timestamp = &t
	return b
}

func (b *EmbedBuilder) SetFooter(text, iconURL string) *EmbedBuilder {
	b.embed.Footer = &discord.EmbedFooter{Text: text, IconURL: iconURL}
	return b
}

func (b *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	b.embed.Image = &discord.EmbedImage{URL: url}
	return b
}

func (b *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	b.embed.Thumbnail = &discord.EmbedThumbnail{URL: url}
	return b
}

func (b *EmbedBuilder) SetAuthor(name, url, iconURL string) *EmbedBuilder {
	a := &discord.EmbedAuthor{Name: name}
	if url != "" {
		a.URL = &url
	}
	if iconURL != "" {
		a.IconURL = &iconURL
	}
	b.embed.Author = a
	return b
}

// AddFields appends one or more fields to the embed.
func (b *EmbedBuilder) AddFields(fields ...discord.EmbedFields) *EmbedBuilder {
	b.embed.Fields = append(b.embed.Fields, fields...)
	return b
}

func (b *EmbedBuilder) Build() discord.Embed {
	embed := b.embed
	if b.embed.Fields != nil {
		fields := make([]discord.EmbedFields, len(b.embed.Fields))
		copy(fields, b.embed.Fields)
		embed.Fields = fields
	}
	return embed
}
