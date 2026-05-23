package builder

import "github.com/streame-gg/go-discord-wrapper/types/components"

// SectionBuilder builds a Section (Components v2) using a fluent API.
//
//	section := builder.NewSection().
//	    AddComponents(builder.NewTextDisplay().SetContent("Hello").Build()).
//	    SetAccessory(builder.NewThumbnail().SetURL("https://example.com/img.png").Build()).
//	    Build()
type SectionBuilder struct {
	s components.Section
}

func NewSection() *SectionBuilder {
	return &SectionBuilder{}
}

func (b *SectionBuilder) AddComponents(c ...components.AnySectionComponent) *SectionBuilder {
	b.s.Components = append(b.s.Components, c...)
	return b
}

func (b *SectionBuilder) SetAccessory(a components.AnySectionAccessory) *SectionBuilder {
	b.s.Accessory = a
	return b
}

func (b *SectionBuilder) Build() *components.Section {
	return &b.s
}

// ── Thumbnail ─────────────────────────────────────────────────────────────────

// ThumbnailBuilder builds a ThumbnailComponent (Components v2 section accessory).
//
//	thumb := builder.NewThumbnail().SetURL("https://example.com/img.png").Build()
type ThumbnailBuilder struct {
	t components.ThumbnailComponent
}

func NewThumbnail() *ThumbnailBuilder { return &ThumbnailBuilder{} }

func (b *ThumbnailBuilder) SetURL(url string) *ThumbnailBuilder {
	b.t.Media = &components.UnfurledMediaItem{URL: url}
	return b
}

func (b *ThumbnailBuilder) SetDescription(desc string) *ThumbnailBuilder {
	b.t.Description = desc
	return b
}

func (b *ThumbnailBuilder) SetSpoiler(spoiler bool) *ThumbnailBuilder {
	b.t.Spoiler = spoiler
	return b
}

func (b *ThumbnailBuilder) Build() *components.ThumbnailComponent {
	return &b.t
}
