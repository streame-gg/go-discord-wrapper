package discord

import "errors"

// ErrEntityNotHydrated is returned by entity convenience methods when the
// entity was not obtained through the gateway or cache and therefore has no
// associated client. Call Hydrate on the entity, or obtain it via a gateway
// event handler / cache lookup.
var ErrEntityNotHydrated = errors.New("discord: entity not hydrated — call Hydrate or obtain it via gateway event / cache")

// ensureClient checks that c is non-nil and returns it unchanged.
// Returns ErrEntityNotHydrated if c is nil.
func ensureClient(c EntityClient) (EntityClient, error) {
	if c == nil {
		return nil, ErrEntityNotHydrated
	}
	return c, nil
}
