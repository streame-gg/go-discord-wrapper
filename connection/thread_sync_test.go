package connection

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// seedThread adds a thread to both the channel cache and the threadsByParent index.
func seedThread(t *testing.T, c *Client, guildID, parentID, threadID common.Snowflake) {
	t.Helper()
	ch := &common.Channel{
		ID:       threadID,
		GuildID:  &guildID,
		ParentID: &parentID,
		Type:     common.ChannelTypePublicThread,
	}
	c.cacheChannel(ch)
	c.trackThread(ch)
}

// TestBug43ThreadListSync_EvictsStaleThreads verifies that THREAD_LIST_SYNC
// removes threads for the named parent channels that are absent from the
// server's authoritative list (Bug 43).
func TestBug43ThreadListSync_EvictsStaleThreads(t *testing.T) {
	c := newClientWithCache(t)

	guildID := common.Snowflake("1000")
	parentA := common.Snowflake("2000")
	parentB := common.Snowflake("3000")

	staleA := common.Snowflake("4000") // should be evicted — not in the sync list
	keepA := common.Snowflake("4001")  // should be kept — included in the sync list
	keepB := common.Snowflake("5000")  // should be untouched — different parent, not in scope

	// Seed the cache.
	seedThread(t, c, guildID, parentA, staleA)
	seedThread(t, c, guildID, parentA, keepA)
	seedThread(t, c, guildID, parentB, keepB)

	// Dispatch THREAD_LIST_SYNC scoped to parentA only; keepA is included but staleA is not.
	ev := events.ThreadListSyncEvent{
		GuildID:    guildID,
		ChannelIDs: []common.Snowflake{parentA},
		Threads: []common.Channel{
			{ID: keepA, GuildID: &guildID, ParentID: &parentA, Type: common.ChannelTypePublicThread},
		},
	}
	raw, err := json.Marshal(ev)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventThreadListSync, nil)

	_, staleInCache := c.Cache.Channels().Get(staleA)
	assert.False(t, staleInCache, "stale thread in synced parent must be evicted (Bug 43)")

	_, keepAInCache := c.Cache.Channels().Get(keepA)
	assert.True(t, keepAInCache, "thread present in THREAD_LIST_SYNC must remain in cache (Bug 43)")

	_, keepBInCache := c.Cache.Channels().Get(keepB)
	assert.True(t, keepBInCache, "thread in un-synced parent must not be evicted (Bug 43)")
}

// TestBug43ThreadListSync_NoChannelIDsEvictsGuild verifies that when ChannelIDs
// is absent, THREAD_LIST_SYNC evicts stale threads from all parents in the guild.
func TestBug43ThreadListSync_NoChannelIDsEvictsGuild(t *testing.T) {
	c := newClientWithCache(t)

	guildID := common.Snowflake("1001")
	parentA := common.Snowflake("2001")
	parentB := common.Snowflake("3001")

	staleA := common.Snowflake("4100")
	staleB := common.Snowflake("5100")
	keep := common.Snowflake("4101")

	// Seed parents and threads.
	c.cacheChannel(&common.Channel{ID: parentA, GuildID: &guildID})
	c.cacheChannel(&common.Channel{ID: parentB, GuildID: &guildID})
	seedThread(t, c, guildID, parentA, staleA)
	seedThread(t, c, guildID, parentB, staleB)
	seedThread(t, c, guildID, parentA, keep)

	// Dispatch a guild-wide THREAD_LIST_SYNC with only `keep` in the list.
	ev := events.ThreadListSyncEvent{
		GuildID: guildID, // no ChannelIDs → guild-wide
		Threads: []common.Channel{
			{ID: keep, GuildID: &guildID, ParentID: &parentA, Type: common.ChannelTypePublicThread},
		},
	}
	raw, err := json.Marshal(ev)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventThreadListSync, nil)

	_, okA := c.Cache.Channels().Get(staleA)
	assert.False(t, okA, "staleA must be evicted in guild-wide THREAD_LIST_SYNC (Bug 43)")

	_, okB := c.Cache.Channels().Get(staleB)
	assert.False(t, okB, "staleB must be evicted in guild-wide THREAD_LIST_SYNC (Bug 43)")

	_, okKeep := c.Cache.Channels().Get(keep)
	assert.True(t, okKeep, "thread present in THREAD_LIST_SYNC must remain in cache (Bug 43)")
}

// TestBug43ThreadIndexUpdatesOnCreateAndDelete verifies that the threadsByParent
// index is maintained correctly by THREAD_CREATE and THREAD_DELETE events.
func TestBug43ThreadIndexUpdatesOnCreateAndDelete(t *testing.T) {
	c := newClientWithCache(t)

	guildID := common.Snowflake("1002")
	parentID := common.Snowflake("2002")
	threadID := common.Snowflake("3002")

	// THREAD_CREATE: thread should appear in threadsByParent.
	createPayload := map[string]any{
		"id": string(threadID), "guild_id": string(guildID),
		"parent_id": string(parentID), "type": common.ChannelTypePublicThread,
	}
	raw, err := json.Marshal(createPayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventThreadCreate, nil)

	c.threadIndexMu.RLock()
	_, tracked := c.threadsByParent[parentID][threadID]
	c.threadIndexMu.RUnlock()
	assert.True(t, tracked, "thread must be indexed after THREAD_CREATE (Bug 43)")

	// THREAD_DELETE: thread must be removed from the index.
	deletePayload := map[string]any{
		"id": string(threadID), "guild_id": string(guildID),
		"parent_id": string(parentID), "type": common.ChannelTypePublicThread,
	}
	raw, err = json.Marshal(deletePayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventThreadDelete, nil)

	c.threadIndexMu.RLock()
	_, stillTracked := c.threadsByParent[parentID][threadID]
	c.threadIndexMu.RUnlock()
	assert.False(t, stillTracked, "thread must be removed from index after THREAD_DELETE (Bug 43)")
}
