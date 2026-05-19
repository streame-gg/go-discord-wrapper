package discord

import "context"

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the user.
func (u *User) Hydrate(c EntityClient) { u.hClient = c }

// WithClient returns a shallow copy of the User with the client replaced.
func (u *User) WithClient(c EntityClient) *User {
	cp := *u
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the user has an associated client.
func (u *User) IsHydrated() bool { return u.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// CreateDM opens a direct message channel with this user.
func (u *User) CreateDM(ctx context.Context) (*Channel, error) {
	c, err := ensureClient(u.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateDM(ctx, u.ID)
}

// FetchProfile re-fetches the user's public profile from the API.
func (u *User) FetchProfile(ctx context.Context) (*User, error) {
	c, err := ensureClient(u.hClient)
	if err != nil {
		return nil, err
	}
	return c.GetUser(ctx, u.ID)
}
