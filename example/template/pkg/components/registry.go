// Package components provides the shared core for handling message-component and
// modal interactions: a Handler interface and a reloadable Registry keyed by
// custom ID.
//
// The handlers themselves live in three sibling packages, one per interaction
// kind, each owning its own Registry:
//
//	components/buttons      — button clicks
//	components/selectmenus  — select-menu choices
//	components/modals       — modal submissions
//
// Splitting them means a button and a modal can safely share a custom ID, and
// the bot router can dispatch by interaction type to the right registry. Each
// handler file self-registers from an init(); Reload rebuilds the lookup table
// without a restart (see /dev reload).
package components

import (
	"sync"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// Handler handles one component or modal interaction, matched by custom ID.
type Handler interface {
	// CustomID is the custom_id to match: set it on the button, menu, or modal
	// you send, and incoming interactions are routed to the matching handler.
	CustomID() string
	// Handle runs when an interaction carrying this custom ID arrives.
	Handle(c *connection.Client, ev *events.InteractionCreateEvent)
}

// Registry is a reloadable set of handlers keyed by custom ID. Each interaction
// kind (buttons, select menus, modals) owns one, so the same custom ID may be
// reused across kinds without colliding.
type Registry struct {
	mu      sync.RWMutex
	sources []Handler
	live    map[string]Handler
}

// NewRegistry returns an empty registry. Call Reload once before Lookup to build
// the lookup table from whatever has registered.
func NewRegistry() *Registry {
	return &Registry{live: map[string]Handler{}}
}

// Register adds h to the set. Call it from an init() in each handler file.
func (r *Registry) Register(h Handler) {
	r.sources = append(r.sources, h)
}

// Lookup returns the handler registered under id, if any.
func (r *Registry) Lookup(id string) (Handler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.live[id]
	return h, ok
}

// Reload rebuilds the lookup table from the registered handlers. Panics on a
// duplicate custom ID so collisions surface immediately instead of silently
// shadowing a handler.
func (r *Registry) Reload() {
	m := make(map[string]Handler, len(r.sources))
	for _, h := range r.sources {
		id := h.CustomID()
		if _, dup := m[id]; dup {
			panic("components: duplicate custom ID " + id)
		}
		m[id] = h
	}
	r.mu.Lock()
	r.live = m
	r.mu.Unlock()
}
