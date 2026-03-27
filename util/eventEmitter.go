package util

import "sync"

type EventEmitter[K comparable, H any] struct {
	mu       sync.RWMutex
	handlers map[K][]H
}

func NewEventEmitter[K comparable, H any]() *EventEmitter[K, H] {
	return &EventEmitter[K, H]{
		handlers: make(map[K][]H),
	}
}

func (e *EventEmitter[K, H]) On(event K, handler H) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlers[event] = append(e.handlers[event], handler)
}

func (e *EventEmitter[K, H]) Handlers(event K) []H {
	e.mu.RLock()
	defer e.mu.RUnlock()

	handlers := e.handlers[event]
	if len(handlers) == 0 {
		return nil
	}

	cloned := make([]H, len(handlers))
	copy(cloned, handlers)
	return cloned
}
