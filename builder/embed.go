// Package builder provides fluent constructors for Discord UI components:
// embeds, buttons, select menus, modals, action rows, and text inputs.
package builder

import (
	"fmt"
	"time"
	"unicode/utf8"

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

// Validate checks whether the embed satisfies Discord's per-field and total
// character limits. Returns a descriptive error for the first violation found.
//
// Limits enforced:
//   - Title       ≤ 256 characters
//   - Description ≤ 4 096 characters
//   - Fields      ≤ 25
//   - Field.Name  ≤ 256 characters each
//   - Field.Value ≤ 1 024 characters each
//   - Footer.Text ≤ 2 048 characters
//   - Author.Name ≤ 256 characters
//   - Total (sum of all text) ≤ 6 000 characters
func (b *EmbedBuilder) Validate() error {
	e := &b.embed
	total := 0

	if e.Title != nil {
		n := utf8.RuneCountInString(*e.Title)
		if n > 256 {
			return fmt.Errorf("embed title exceeds 256 characters (%d)", n)
		}
		total += n
	}

	if e.Description != nil {
		n := utf8.RuneCountInString(*e.Description)
		if n > 4096 {
			return fmt.Errorf("embed description exceeds 4096 characters (%d)", n)
		}
		total += n
	}

	if len(e.Fields) > 25 {
		return fmt.Errorf("embed has %d fields; maximum is 25", len(e.Fields))
	}

	for i, f := range e.Fields {
		nameLen := utf8.RuneCountInString(f.Name)
		valueLen := utf8.RuneCountInString(f.Value)
		if nameLen > 256 {
			return fmt.Errorf("embed field %d name exceeds 256 characters (%d)", i, nameLen)
		}
		if valueLen > 1024 {
			return fmt.Errorf("embed field %d value exceeds 1024 characters (%d)", i, valueLen)
		}
		total += nameLen + valueLen
	}

	if e.Footer != nil {
		n := utf8.RuneCountInString(e.Footer.Text)
		if n > 2048 {
			return fmt.Errorf("embed footer text exceeds 2048 characters (%d)", n)
		}
		total += n
	}

	if e.Author != nil {
		n := utf8.RuneCountInString(e.Author.Name)
		if n > 256 {
			return fmt.Errorf("embed author name exceeds 256 characters (%d)", n)
		}
		total += n
	}

	if total > 6000 {
		return fmt.Errorf("embed total character count %d exceeds 6000", total)
	}

	return nil
}

// BuildChecked calls Validate before building. Returns an error if any Discord
// limit is violated; returns a zero-value Embed on error.
// Build() is unaffected and continues to return without validation.
func (b *EmbedBuilder) BuildChecked() (discord.Embed, error) {
	if err := b.Validate(); err != nil {
		return discord.Embed{}, err
	}
	return b.Build(), nil
}
