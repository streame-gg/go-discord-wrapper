package connection

// Tests for Bug 78: GUILD_SCHEDULED_EVENT_USER_ADD/REMOVE must update UserCount
// on the cached GuildScheduledEvent.

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

func seedScheduledEvent(t *testing.T, c *Client, guildID, eventID discord.Snowflake, userCount int) {
	t.Helper()
	c.Cache.ScheduledEvents().Set(&discord.GuildScheduledEvent{
		ID:        eventID,
		GuildID:   guildID,
		UserCount: &userCount,
	})
}

func (cs *ConnectionSuite) TestBug78_ScheduledEventUserAdd_IncrementsUserCount() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100)
	eventID := discord.Snowflake(200)
	seedScheduledEvent(t, c, guildID, eventID, 5)

	payload := map[string]any{
		"guild_scheduled_event_id": eventID.String(),
		"user_id":                  "999",
		"guild_id":                 guildID.String(),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildScheduledEventUserAdd, nil)

	got, ok := c.Cache.ScheduledEvents().Get(eventID)
	require.True(t, ok)
	require.NotNil(t, got.UserCount)
	assert.Equal(t, 6, *got.UserCount, "UserCount must be incremented to 6 (Bug 78)")
}

func (cs *ConnectionSuite) TestBug78_ScheduledEventUserRemove_DecrementsUserCount() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100)
	eventID := discord.Snowflake(200)
	seedScheduledEvent(t, c, guildID, eventID, 5)

	payload := map[string]any{
		"guild_scheduled_event_id": eventID.String(),
		"user_id":                  "999",
		"guild_id":                 guildID.String(),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildScheduledEventUserRemove, nil)

	got, ok := c.Cache.ScheduledEvents().Get(eventID)
	require.True(t, ok)
	require.NotNil(t, got.UserCount)
	assert.Equal(t, 4, *got.UserCount, "UserCount must be decremented to 4 (Bug 78)")
}

func (cs *ConnectionSuite) TestBug78_ScheduledEventUserRemove_DoesNotGoBelowZero() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100)
	eventID := discord.Snowflake(200)
	seedScheduledEvent(t, c, guildID, eventID, 0)

	payload := map[string]any{
		"guild_scheduled_event_id": eventID.String(),
		"user_id":                  "999",
		"guild_id":                 guildID.String(),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildScheduledEventUserRemove, nil)

	got, ok := c.Cache.ScheduledEvents().Get(eventID)
	require.True(t, ok)
	require.NotNil(t, got.UserCount)
	assert.Equal(t, 0, *got.UserCount, "UserCount must not go below 0 (Bug 78)")
}

func (cs *ConnectionSuite) TestBug78_ScheduledEventUserAdd_NilUserCount_IsNoop() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100)
	eventID := discord.Snowflake(200)

	// Seed with nil UserCount (Discord sometimes omits this field).
	c.Cache.ScheduledEvents().Set(&discord.GuildScheduledEvent{
		ID:        eventID,
		GuildID:   guildID,
		UserCount: nil,
	})

	payload := map[string]any{
		"guild_scheduled_event_id": eventID.String(),
		"user_id":                  "999",
		"guild_id":                 guildID.String(),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildScheduledEventUserAdd, nil)

	got, ok := c.Cache.ScheduledEvents().Get(eventID)
	require.True(t, ok)
	// When UserCount was nil we don't blindly write 1 (we don't know the real count).
	assert.Nil(t, got.UserCount, "UserCount must remain nil when it was not known (Bug 78)")
}

func (cs *ConnectionSuite) TestBug78_ScheduledEventUserAdd_UncachedEvent_IsNoop() {
	t := cs.T()
	c := newClientWithCache(t)

	payload := map[string]any{
		"guild_scheduled_event_id": "999999",
		"user_id":                  "111",
		"guild_id":                 "222",
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	// Should not panic or insert phantom entry.
	c.internalEventHandler(json.RawMessage(raw), events.EventGuildScheduledEventUserAdd, nil)

	assert.Equal(t, 0, c.Cache.ScheduledEvents().Size(),
		"no phantom entry must appear for uncached scheduled event (Bug 78)")
}
