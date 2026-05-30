package connection

// Tests for Bug 82 (THREAD_UPDATE fallback OldThread), Bug 83 (GUILD_DELETE
// threadsByParent leak), and Bug 84 (THREAD_UPDATE ParentID reparent leak).

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── Bug 82: THREAD_UPDATE fallback path sets OldThread ───────────────────────

func TestBug82_ThreadUpdate_FallbackSetsOldThread(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	parentID := discord.Snowflake(2000)
	threadID := discord.Snowflake(3000)

	// Pre-seed thread in cache.
	c.cacheChannel(&discord.Channel{
		ID:       threadID,
		GuildID:  &guildID,
		ParentID: &parentID,
		Name:     "old-name",
	})

	payload := map[string]any{
		"id":        threadID.String(),
		"guild_id":  guildID.String(),
		"parent_id": parentID.String(),
		"name":      "new-name",
		"type":      11,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	// Pre-unmarshal so ev.NewThread.ID is populated (mirrors the websocket loop).
	var ev events.ThreadUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventThreadUpdate, &ev)

	// OldThread must have been populated from cache before overwrite.
	require.NotNil(t, ev.OldThread, "OldThread must be set from cache (Bug 82)")
	assert.Equal(t, "old-name", ev.OldThread.Name)
	assert.Equal(t, "new-name", ev.NewThread.Name)
}

func TestBug82_ThreadUpdate_FallbackWithNilEvent_UpdatesCache(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	parentID := discord.Snowflake(2000)
	threadID := discord.Snowflake(3000)

	c.cacheChannel(&discord.Channel{
		ID:       threadID,
		GuildID:  &guildID,
		ParentID: &parentID,
		Name:     "old-name",
	})

	payload := map[string]any{
		"id":        threadID.String(),
		"guild_id":  guildID.String(),
		"parent_id": parentID.String(),
		"name":      "updated-name",
		"type":      11,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	// Pass nil — exercises the fallback unmarshal path.
	c.internalEventHandler(json.RawMessage(raw), events.EventThreadUpdate, nil)

	got, ok := c.Cache.Channels().Get(threadID)
	require.True(t, ok)
	assert.Equal(t, "updated-name", got.Name, "cache must reflect updated thread name after fallback path (Bug 82)")
}

// ── Bug 83: GUILD_DELETE cleans up threadsByParent index ─────────────────────

func TestBug83_GuildDelete_CleansThreadsByParentIndex(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	parentID := discord.Snowflake(2000)
	threadID := discord.Snowflake(3000)

	// Register a parent channel and a thread under it.
	c.cacheChannel(&discord.Channel{ID: parentID, GuildID: &guildID})
	c.cacheChannel(&discord.Channel{ID: threadID, GuildID: &guildID, ParentID: &parentID})
	c.trackThread(&discord.Channel{ID: threadID, GuildID: &guildID, ParentID: &parentID})

	// Verify the thread is indexed.
	c.threadIndexMu.RLock()
	threads := c.threadsByParent[parentID]
	c.threadIndexMu.RUnlock()
	require.Len(t, threads, 1, "thread must be indexed before GUILD_DELETE")

	dispatchGuildDelete(t, c, guildID)

	// After GUILD_DELETE, threadsByParent must have no entry for parentID.
	c.threadIndexMu.RLock()
	threadsAfter := c.threadsByParent[parentID]
	c.threadIndexMu.RUnlock()
	assert.Empty(t, threadsAfter, "threadsByParent must be cleaned up on GUILD_DELETE (Bug 83)")
}

// ── Bug 84: THREAD_UPDATE with changed ParentID removes old index entry ───────

func TestBug84_ThreadUpdate_ReparentCleansOldIndex(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	oldParentID := discord.Snowflake(2000)
	newParentID := discord.Snowflake(3000)
	threadID := discord.Snowflake(4000)

	// Seed thread under oldParent.
	c.cacheChannel(&discord.Channel{ID: threadID, GuildID: &guildID, ParentID: &oldParentID})
	c.trackThread(&discord.Channel{ID: threadID, GuildID: &guildID, ParentID: &oldParentID})

	c.threadIndexMu.RLock()
	assert.Contains(t, c.threadsByParent[oldParentID], threadID, "thread must be in old parent index")
	c.threadIndexMu.RUnlock()

	// THREAD_UPDATE reparents the thread.
	payload := map[string]any{
		"id":        threadID.String(),
		"guild_id":  guildID.String(),
		"parent_id": newParentID.String(),
		"name":      "thread",
		"type":      11,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	// Pre-unmarshal so ev.NewThread.ID is populated — same as the websocket loop.
	var ev events.ThreadUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventThreadUpdate, &ev)

	c.threadIndexMu.RLock()
	defer c.threadIndexMu.RUnlock()
	assert.NotContains(t, c.threadsByParent[oldParentID], threadID,
		"thread must be removed from old parent index on reparent (Bug 84)")
	assert.Contains(t, c.threadsByParent[newParentID], threadID,
		"thread must be added to new parent index on reparent (Bug 84)")
}

func TestBug84_ThreadUpdate_SameParent_IndexUnchanged(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	parentID := discord.Snowflake(2000)
	threadID := discord.Snowflake(3000)

	c.cacheChannel(&discord.Channel{ID: threadID, GuildID: &guildID, ParentID: &parentID})
	c.trackThread(&discord.Channel{ID: threadID, GuildID: &guildID, ParentID: &parentID})

	// Same parent — no reparent.
	payload := map[string]any{
		"id":        threadID.String(),
		"guild_id":  guildID.String(),
		"parent_id": parentID.String(),
		"name":      "updated",
		"type":      11,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.ThreadUpdateEvent
	c.internalEventHandler(json.RawMessage(raw), events.EventThreadUpdate, &ev)

	c.threadIndexMu.RLock()
	defer c.threadIndexMu.RUnlock()
	assert.Contains(t, c.threadsByParent[parentID], threadID,
		"thread must remain in same parent index when parent did not change (Bug 84)")
}
