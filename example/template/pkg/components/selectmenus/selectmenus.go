// Package selectmenus holds select-menu interaction handlers, matched by custom
// ID. Works like the buttons package: each menu self-registers from an init(),
// the bot routes choices here by custom ID, and Reload rebuilds the table.
package selectmenus

import (
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

var registry = components.NewRegistry()

// Register adds a select-menu handler. Call it from an init() in each file.
func Register(h components.Handler) { registry.Register(h) }

// Lookup returns the select-menu handler for the given custom ID, if any.
func Lookup(id string) (components.Handler, bool) { return registry.Lookup(id) }

// Reload rebuilds the lookup table, returning an error on a duplicate custom ID.
// Called at startup and by /dev reload.
func Reload() error { return registry.Reload() }

// SelectedValues returns the values the user chose in a select-menu interaction.
func SelectedValues(i *interactions.Interaction) []string {
	if data, ok := interactions.As[*responses.InteractionDataMessageComponent](i.Data); ok {
		return data.Values
	}
	return nil
}
