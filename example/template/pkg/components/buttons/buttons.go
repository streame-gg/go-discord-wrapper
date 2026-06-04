// Package buttons holds button interaction handlers, matched by custom ID.
//
// Each button lives in its own file and self-registers from an init():
//
//	func init() { buttons.Register(pingAgain{}) }
//
// Send a button carrying the same custom ID from a command (see commands/ping.go),
// and the bot routes clicks here. Reload rebuilds the table (see /dev reload).
package buttons

import "github.com/streame-gg/go-discord-wrapper/example/template/pkg/components"

var registry = components.NewRegistry()

// Register adds a button handler. Call it from an init() in each handler file.
func Register(h components.Handler) { registry.Register(h) }

// Lookup returns the button handler for the given custom ID, if any.
func Lookup(id string) (components.Handler, bool) { return registry.Lookup(id) }

// Reload rebuilds the lookup table. Called at startup and by /dev reload.
func Reload() { registry.Reload() }
