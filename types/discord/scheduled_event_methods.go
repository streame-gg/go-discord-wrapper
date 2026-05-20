package discord

import (
	"context"
	"time"
)

// ── Options ───────────────────────────────────────────────────────────────────

// ScheduledEventEditOptions configures GuildScheduledEvent.Edit.
type ScheduledEventEditOptions struct {
	ChannelID          *Snowflake
	EntityMetadata     *GuildScheduledEventEntityMetadata
	Name               *string
	PrivacyLevel       *GuildScheduledEventPrivacyLevel
	ScheduledStartTime *time.Time
	ScheduledEndTime   *time.Time
	Description        *string
	EntityType         *GuildScheduledEventEntityType
	Status             *GuildScheduledEventStatus
	Image              *string
}

// FetchUsersOptions configures GuildScheduledEvent.FetchUsers.
type FetchUsersOptions struct {
	Limit      *int
	WithMember *bool
	Before     *Snowflake
	After      *Snowflake
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client into the scheduled event.
func (e *GuildScheduledEvent) Hydrate(c EntityClient) {
	e.hClient = c
	if e.Creator != nil {
		e.Creator.Hydrate(c)
	}
}

// WithClient returns a shallow copy of the GuildScheduledEvent with the client replaced.
func (e *GuildScheduledEvent) WithClient(c EntityClient) *GuildScheduledEvent {
	cp := *e
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the scheduled event has an associated client.
func (e *GuildScheduledEvent) IsHydrated() bool { return e.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this scheduled event. Requires MANAGE_EVENTS.
func (e *GuildScheduledEvent) Edit(ctx context.Context, opts ScheduledEventEditOptions) (*GuildScheduledEvent, error) {
	c, err := ensureClient(e.hClient)
	if err != nil {
		return nil, err
	}
	return c.ModifyGuildScheduledEvent(ctx, e.GuildID, e.ID, opts)
}

// Delete deletes this scheduled event. Requires MANAGE_EVENTS.
func (e *GuildScheduledEvent) Delete(ctx context.Context) error {
	c, err := ensureClient(e.hClient)
	if err != nil {
		return err
	}
	return c.DeleteGuildScheduledEvent(ctx, e.GuildID, e.ID)
}

// FetchUsers retrieves users who have subscribed to this scheduled event.
func (e *GuildScheduledEvent) FetchUsers(ctx context.Context, opts FetchUsersOptions) ([]*GuildScheduledEventUser, error) {
	c, err := ensureClient(e.hClient)
	if err != nil {
		return nil, err
	}
	return c.GetGuildScheduledEventUsers(ctx, e.GuildID, e.ID, opts)
}
