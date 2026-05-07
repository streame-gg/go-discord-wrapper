package util

import "sync"

type handlerEntry[H any] struct {
	id      int
	handler H
	once    bool
}

type EventEmitter[K comparable, H any] struct {
	mu       sync.RWMutex
	handlers map[K][]handlerEntry[H]
	nextID   int
}

func NewEventEmitter[K comparable, H any]() *EventEmitter[K, H] {
	return &EventEmitter[K, H]{
		handlers: make(map[K][]handlerEntry[H]),
	}
}

// On registers handler for event and returns a unique ID that can be passed to Off.
func (e *EventEmitter[K, H]) On(event K, handler H) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	id := e.nextID
	e.nextID++
	e.handlers[event] = append(e.handlers[event], handlerEntry[H]{id: id, handler: handler})
	return id
}

// Once registers a one-shot handler that is automatically removed after its first invocation.
// Returns a unique ID that can be passed to Off to cancel it before it fires.
func (e *EventEmitter[K, H]) Once(event K, handler H) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	id := e.nextID
	e.nextID++
	e.handlers[event] = append(e.handlers[event], handlerEntry[H]{id: id, handler: handler, once: true})
	return id
}

// Off removes the handler registered under event with the given id.
// Returns true if a handler was found and removed.
func (e *EventEmitter[K, H]) Off(event K, id int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	entries := e.handlers[event]
	for i, entry := range entries {
		if entry.id == id {
			e.handlers[event] = append(entries[:i], entries[i+1:]...)
			return true
		}
	}
	return false
}

// Handlers returns all handlers registered for event. One-shot handlers (registered
// via Once) are removed from the emitter atomically so they fire exactly once.
func (e *EventEmitter[K, H]) Handlers(event K) []H {
	e.mu.Lock()
	defer e.mu.Unlock()

	entries := e.handlers[event]
	if len(entries) == 0 {
		return nil
	}

	out := make([]H, len(entries))
	keep := entries[:0]
	for i, entry := range entries {
		out[i] = entry.handler
		if !entry.once {
			keep = append(keep, entry)
		}
	}
	e.handlers[event] = keep
	return out
}
