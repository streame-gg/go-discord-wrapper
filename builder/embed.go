package builder

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// EmbedBuilder builds a common.Embed using a fluent API.
//
//	embed := builder.NewEmbed().
//	    SetTitle("Hello").
//	    SetDescription("World").
//	    SetColor(0x5865F2).
//	    Build()
type EmbedBuilder struct {
	embed common.Embed
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
	b.embed.Footer = &common.EmbedFooter{Text: text, IconURL: iconURL}
	return b
}

func (b *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	b.embed.Image = &common.EmbedImage{URL: url}
	return b
}

func (b *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	b.embed.Thumbnail = &common.EmbedThumbnail{URL: url}
	return b
}

func (b *EmbedBuilder) SetAuthor(name, url, iconURL string) *EmbedBuilder {
	a := &common.EmbedAuthor{Name: name}
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
func (b *EmbedBuilder) AddFields(fields ...common.EmbedFields) *EmbedBuilder {
	b.embed.Fields = append(b.embed.Fields, fields...)
	return b
}

func (b *EmbedBuilder) Build() common.Embed {
	return b.embed
}
