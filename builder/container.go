package builder

import "github.com/streame-gg/go-discord-wrapper/types/components"

// ContainerBuilder builds a Container (Components v2) using a fluent API.
//
//	c := builder.NewContainer().
//	    SetAccentColor(0x5865F2).
//	    AddComponents(
//	        builder.NewTextDisplay().SetContent("Hello!").Build(),
//	    ).
//	    Build()
type ContainerBuilder struct {
	c components.Container
}

func NewContainer() *ContainerBuilder {
	return &ContainerBuilder{
		c: components.Container{
			Components: make([]components.AnyContainerComponent, 0),
		},
	}
}

func (b *ContainerBuilder) SetAccentColor(color int) *ContainerBuilder {
	b.c.AccentColor = color
	return b
}

func (b *ContainerBuilder) SetSpoiler(spoiler bool) *ContainerBuilder {
	b.c.Spoiler = spoiler
	return b
}

func (b *ContainerBuilder) AddComponents(c ...components.AnyContainerComponent) *ContainerBuilder {
	b.c.Components = append(b.c.Components, c...)
	return b
}

func (b *ContainerBuilder) Build() *components.Container {
	return &b.c
}

// ── Text display ──────────────────────────────────────────────────────────────

// TextDisplayBuilder builds a TextDisplayComponent (Components v2).
//
//	td := builder.NewTextDisplay().SetContent("## Hello world").Build()
type TextDisplayBuilder struct {
	td components.TextDisplayComponent
}

func NewTextDisplay() *TextDisplayBuilder { return &TextDisplayBuilder{} }

func (b *TextDisplayBuilder) SetContent(content string) *TextDisplayBuilder {
	b.td.Content = content
	return b
}

func (b *TextDisplayBuilder) Build() *components.TextDisplayComponent {
	return &b.td
}

// ── Separator ─────────────────────────────────────────────────────────────────

// SeparatorBuilder builds a SeparatorComponent (Components v2).
//
//	sep := builder.NewSeparator().SetDivider(true).SetSpacing(components.SeparatorComponentSpacingLarge).Build()
type SeparatorBuilder struct {
	s components.SeparatorComponent
}

func NewSeparator() *SeparatorBuilder { return &SeparatorBuilder{} }

func (b *SeparatorBuilder) SetDivider(divider bool) *SeparatorBuilder {
	b.s.Divider = divider
	return b
}

func (b *SeparatorBuilder) SetSpacing(spacing components.SeparatorComponentSpacing) *SeparatorBuilder {
	b.s.SeparatorComponentSpacing = spacing
	return b
}

func (b *SeparatorBuilder) Build() *components.SeparatorComponent {
	return &b.s
}
