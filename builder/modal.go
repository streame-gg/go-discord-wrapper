package builder

import "github.com/streame-gg/go-discord-wrapper/types/components"

// ModalBuilder builds a Modal using a fluent API.
//
//	modal := builder.NewModal().
//	    SetCustomID("my_modal").
//	    SetTitle("User Info").
//	    AddComponents(
//	        builder.NewLabel().
//	            SetLabel("Name").
//	            SetComponent(builder.NewTextInput().SetCustomID("name").SetStyle(components.TextInputStyleShort).Build()).
//	            Build(),
//	    ).
//	    Build()
//
// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-modal
type ModalBuilder struct {
	modal components.Modal
}

func NewModal() *ModalBuilder {
	comps := []components.Label{}
	return &ModalBuilder{modal: components.Modal{Components: comps}}
}

func (b *ModalBuilder) SetCustomID(id string) *ModalBuilder {
	b.modal.CustomID = id
	return b
}

func (b *ModalBuilder) SetTitle(title string) *ModalBuilder {
	b.modal.Title = title
	return b
}

func (b *ModalBuilder) AddComponents(c ...*components.Label) *ModalBuilder {
	for _, comp := range c {
		b.modal.Components = append(b.modal.Components, *comp)
	}
	return b
}

func (b *ModalBuilder) Build() *components.Modal {
	modal := b.modal
	if b.modal.Components != nil {
		comps := make([]components.Label, len(b.modal.Components))
		copy(comps, b.modal.Components)
		modal.Components = comps
	}
	return &modal
}
