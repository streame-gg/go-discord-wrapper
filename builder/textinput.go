package builder

import (
	"github.com/streame-gg/go-discord-wrapper/types/components"
)

// TextInputBuilder builds a TextInputComponent using a fluent API.
//
//	input := builder.NewTextInput().
//	    SetCustomID("name").
//	    SetStyle(components.TextInputStyleShort).
//	    SetPlaceholder("Your name").
//	    SetRequired(true).
//	    Build()
type TextInputBuilder struct {
	input components.TextInputComponent
}

func NewTextInput() *TextInputBuilder { return &TextInputBuilder{} }

func (b *TextInputBuilder) SetCustomID(id string) *TextInputBuilder {
	b.input.CustomID = id
	return b
}

func (b *TextInputBuilder) SetStyle(style components.TextInputStyle) *TextInputBuilder {
	b.input.Style = style
	return b
}

func (b *TextInputBuilder) SetMinLength(min int) *TextInputBuilder {
	b.input.MinLength = &min
	return b
}

func (b *TextInputBuilder) SetMaxLength(max int) *TextInputBuilder {
	b.input.MaxLength = &max
	return b
}

func (b *TextInputBuilder) SetRequired(required bool) *TextInputBuilder {
	b.input.Required = &required
	return b
}

// SetValue pre-fills the input with the given value.
func (b *TextInputBuilder) SetValue(value string) *TextInputBuilder {
	b.input.Value = value
	return b
}

func (b *TextInputBuilder) SetPlaceholder(placeholder string) *TextInputBuilder {
	b.input.Placeholder = placeholder
	return b
}

func (b *TextInputBuilder) Build() *components.TextInputComponent {
	return &b.input
}

// ── Label (Components v2 form wrapper) ───────────────────────────────────────

// LabelBuilder builds a LabelComponent (Components v2) that wraps an input
// with a label and optional description.
//
//	label := builder.NewLabel().
//	    SetLabel("Your name").
//	    SetDescription("Enter your full name").
//	    SetComponent(builder.NewTextInput().SetCustomID("name").SetStyle(components.TextInputStyleShort).Build()).
//	    Build()
type LabelBuilder struct {
	lbl components.LabelComponent
}

func NewLabel() *LabelBuilder { return &LabelBuilder{} }

func (b *LabelBuilder) SetLabel(label string) *LabelBuilder {
	b.lbl.Label = label
	return b
}

func (b *LabelBuilder) SetDescription(desc string) *LabelBuilder {
	b.lbl.Description = desc
	return b
}

func (b *LabelBuilder) SetComponent(c components.AnyChildComponent) *LabelBuilder {
	b.lbl.Component = c
	return b
}

func (b *LabelBuilder) Build() *components.LabelComponent {
	return &b.lbl
}
