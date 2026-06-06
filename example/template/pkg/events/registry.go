// Package events holds the bot's gateway event handlers plus a small registry
// that makes them reloadable at runtime.
//
// A handler file registers itself from an init() with On:
//
//	func init() {
//	    events.On(devents.EventMessageCreate,
//	        func(c *connection.Client, e *devents.MessageCreateEvent) { … })
//	}
//
// Drop a new file in this folder and it wires itself up — no central list to
// edit. Attach installs the handlers on the client once at startup; Reload
// re-applies the registry without restarting the bot or double-registering on
// the client (see the /dev reload command).
package events

import (
	"sync"

	"github.com/streame-gg/go-discord-wrapper/connection"
	devents "github.com/streame-gg/go-discord-wrapper/types/events"
)

// Handler is a gateway event handler already bound to its concrete event type.
type Handler func(*connection.Client, devents.Event)

type registration struct {
	event   devents.EventType
	handler Handler
}

// registrations is the static set of handlers built by the init()s in this
// package. Go is compiled, so this never changes at runtime — it is the source
// of truth that Attach and Reload apply to the live client.
var registrations []registration

// On registers a typed handler for a gateway event. Call it from an init().
// T is the concrete event pointer (e.g. *devents.MessageCreateEvent); the
// dispatched event is asserted to T before your handler runs, so handlers stay
// fully typed.
func On[T devents.Event](event devents.EventType, h func(*connection.Client, T)) {
	registrations = append(registrations, registration{
		event: event,
		handler: func(c *connection.Client, ev devents.Event) {
			if typed, ok := ev.(T); ok {
				h(c, typed)
			}
		},
	})
}

var (
	mu       sync.RWMutex
	live     map[devents.EventType][]Handler
	attached = map[devents.EventType]bool{}
)

// Attach wires every registered handler onto c. Call it once at startup.
//
// It installs exactly one stable dispatcher per event type, and each dispatcher
// reads the current ("live") handler set under a lock. That indirection is what
// makes Reload safe: rebuilding the live set swaps behaviour instantly without
// ever adding a second handler to the client (the wrapper exposes no removal
// API, so re-attaching directly would double-fire every handler).
func Attach(c *connection.Client) {
	rebuild()

	mu.RLock()
	types := make([]devents.EventType, 0, len(live))
	for et := range live {
		types = append(types, et)
	}
	mu.RUnlock()

	for _, et := range types {
		if attached[et] {
			continue
		}
		attached[et] = true

		et := et
		_ = c.OnEvent(et, connection.EventHandler(func(c *connection.Client, ev devents.Event) {
			mu.RLock()
			hs := live[et]
			mu.RUnlock()
			for _, h := range hs {
				h(c, ev)
			}
		}))
	}
}

// Reload rebuilds the live handler set from the registry. Handlers take effect
// immediately because the dispatchers installed by Attach read the live set.
//
// In a compiled bot the handler code itself cannot change without a rebuild, so
// Reload re-applies the same handlers; its value is as the seam a real bot
// extends — e.g. handlers gated on config or feature flags re-read them here.
func Reload() { rebuild() }

// rebuild regenerates the live handler map from the static registrations.
func rebuild() {
	m := make(map[devents.EventType][]Handler)
	for _, r := range registrations {
		m[r.event] = append(m[r.event], r.handler)
	}
	mu.Lock()
	live = m
	mu.Unlock()
}
