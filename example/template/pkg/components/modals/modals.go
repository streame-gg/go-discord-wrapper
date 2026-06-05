// Package modals holds modal-submit handlers, matched by the modal's custom ID.
// Works like the buttons package: each modal handler self-registers from an
// init(), the bot routes submissions here by custom ID, and Reload rebuilds the
// table. Open a modal from a button or command (see components/buttons/open_form.go).
package modals

import (
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components"
	dcomponents "github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

var registry = components.NewRegistry()

// Register adds a modal handler. Call it from an init() in each file.
func Register(h components.Handler) { registry.Register(h) }

// Lookup returns the modal handler for the given custom ID, if any.
func Lookup(id string) (components.Handler, bool) { return registry.Lookup(id) }

// Reload rebuilds the lookup table, returning an error on a duplicate custom ID.
// Called at startup and by /dev reload.
func Reload() error { return registry.Reload() }

// TextValue returns the submitted value of the text input with the given custom
// ID, or "" if the modal carried no such field.
func TextValue(i *interactions.Interaction, customID string) string {
	data, ok := interactions.As[*responses.InteractionDataModalSubmit](i.Data)
	if !ok || data.Components == nil {
		return ""
	}
	for _, label := range *data.Components {
		input, ok := label.Component.(*dcomponents.TextInputComponentInteractionResponse)
		if ok && input.CustomID == customID {
			return input.Value
		}
	}
	return ""
}
