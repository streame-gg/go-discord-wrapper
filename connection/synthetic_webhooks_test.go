package connection

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

func wh(id, name string) *discord.Webhook {
	return &discord.Webhook{ID: discord.Snowflake(mustU64(id)), Name: ptrStr(name)}
}

func mustU64(s string) uint64 {
	var n uint64
	for _, r := range s {
		n = n*10 + uint64(r-'0')
	}
	return n
}

// TestDiffWebhooks covers the pure diff: add, update (name change), delete, and
// no-op for an unchanged set.
func (cs *ConnectionSuite) TestDiffWebhooks() {
	const gid, chid = discord.Snowflake(1), discord.Snowflake(2)

	old := []*discord.Webhook{wh("100", "alpha"), wh("200", "beta")}
	current := []*discord.Webhook{wh("100", "alpha_renamed"), wh("300", "gamma")} // 100 updated, 200 removed, 300 added

	got := diffWebhooks(gid, chid, old, current)

	var creates, updates, deletes int
	for _, e := range got {
		switch ev := e.(type) {
		case *events.WebhookCreateEvent:
			creates++
			cs.Equal(discord.Snowflake(300), ev.Webhook.ID)
			cs.Equal(chid, ev.ChannelID)
		case *events.WebhookUpdateEvent:
			updates++
			cs.Equal(discord.Snowflake(100), ev.NewWebhook.ID)
			cs.Equal("alpha", *ev.OldWebhook.Name, "update must carry the old name")
			cs.Equal("alpha_renamed", *ev.NewWebhook.Name)
		case *events.WebhookDeleteEvent:
			deletes++
			cs.Equal(discord.Snowflake(200), ev.Webhook.ID)
		}
	}
	cs.Equal(1, creates)
	cs.Equal(1, updates)
	cs.Equal(1, deletes)

	// Unchanged set yields nothing.
	cs.Empty(diffWebhooks(gid, chid, old, old))
}

// TestWebhookSynthetic_FullPath drives the whole pipeline: a WEBHOOKS_UPDATE fed
// through internalEventHandler triggers a REST fetch + diff. The first update
// seeds the snapshot (no events); the second produces one Create/Update/Delete
// each, with the Update carrying the old webhook. Also asserts no duplicates.
func (cs *ConnectionSuite) TestWebhookSynthetic_FullPath() {
	t := cs.T()
	const channelID = discord.Snowflake(175928847299117063)
	const guildID = discord.Snowflake(175928847299117001)

	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			_, _ = w.Write([]byte(`[{"id":"500000000000000001","type":1,"name":"a"},{"id":"500000000000000002","type":1,"name":"b"}]`))
		} else {
			_, _ = w.Write([]byte(`[{"id":"500000000000000001","type":1,"name":"a_renamed"},{"id":"500000000000000003","type":1,"name":"c"}]`))
		}
	}))
	defer ts.Close()

	client := newClientWithCacheAndServer(t, ts.URL)
	defer client.Shutdown()

	var created, updated, deleted atomic.Int32
	var mu sync.Mutex
	var oldName string
	client.OnWebhookCreate(func(_ *Client, e *events.WebhookCreateEvent) {
		created.Add(1)
		assert.Equal(t, discord.Snowflake(500000000000000003), e.Webhook.ID)
	})
	client.OnWebhookUpdate(func(_ *Client, e *events.WebhookUpdateEvent) {
		updated.Add(1)
		mu.Lock()
		if e.OldWebhook != nil && e.OldWebhook.Name != nil {
			oldName = *e.OldWebhook.Name
		}
		mu.Unlock()
	})
	client.OnWebhookDelete(func(_ *Client, e *events.WebhookDeleteEvent) {
		deleted.Add(1)
		assert.Equal(t, discord.Snowflake(500000000000000002), e.Webhook.ID)
	})

	feed := func() {
		ev := &events.WebhooksUpdateEvent{GuildID: guildID, ChannelID: channelID}
		client.internalEventHandler(nil, events.EventWebhooksUpdate, ev)
	}

	// First update: cold snapshot -> seed only, no synthetic events.
	feed()
	require.Eventually(t, func() bool {
		_, ok := client.loadWebhookSnapshot(channelID)
		return ok
	}, 5*time.Second, 10*time.Millisecond, "first fetch should seed the snapshot")
	assert.Equal(t, int32(0), created.Load()+updated.Load()+deleted.Load(), "cold snapshot must not fire events")

	// Second update: diff against the seeded snapshot.
	feed()
	require.Eventually(t, func() bool {
		return created.Load() == 1 && updated.Load() == 1 && deleted.Load() == 1
	}, 5*time.Second, 10*time.Millisecond, "expected create=1 update=1 delete=1")

	// Settle and assert no duplicates.
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, int32(1), created.Load())
	assert.Equal(t, int32(1), updated.Load())
	assert.Equal(t, int32(1), deleted.Load())
	mu.Lock()
	assert.Equal(t, "a", oldName, "WebhookUpdate must carry the previously cached name")
	mu.Unlock()
}

// TestWebhookSynthetic_NoHandlers_NoFetch verifies the opt-in guard: with no
// webhook handler registered, a WEBHOOKS_UPDATE triggers no REST fetch.
func (cs *ConnectionSuite) TestWebhookSynthetic_NoHandlers_NoFetch() {
	t := cs.T()
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	client := newClientWithCacheAndServer(t, ts.URL)
	defer client.Shutdown()

	// Only the coarse handler is registered — not the synthetic ones.
	client.OnWebhooksUpdate(func(_ *Client, _ *events.WebhooksUpdateEvent) {})

	ev := &events.WebhooksUpdateEvent{GuildID: discord.Snowflake(175928847299117001), ChannelID: discord.Snowflake(175928847299117063)}
	client.internalEventHandler(nil, events.EventWebhooksUpdate, ev)

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, int32(0), atomic.LoadInt32(&calls), "no webhook handler => no REST fetch")
}

// TestWebhookSynthetic_ListError_NoEventsNoSnapshot covers the REST-error branch
// of fetchDiffEmitWebhooks: when ListChannelWebhooks fails, no synthetic event is
// dispatched and the snapshot is left unseeded (so a later successful fetch still
// only seeds the baseline rather than diffing against partial state).
func (cs *ConnectionSuite) TestWebhookSynthetic_ListError_NoEventsNoSnapshot() {
	t := cs.T()
	const channelID = discord.Snowflake(175928847299117063)
	const guildID = discord.Snowflake(175928847299117001)

	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.Error(w, `{"message":"boom","code":0}`, http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := newClientWithCacheAndServer(t, ts.URL)
	defer client.Shutdown()

	var fired atomic.Int32
	client.OnWebhookCreate(func(_ *Client, _ *events.WebhookCreateEvent) { fired.Add(1) })
	client.OnWebhookUpdate(func(_ *Client, _ *events.WebhookUpdateEvent) { fired.Add(1) })
	client.OnWebhookDelete(func(_ *Client, _ *events.WebhookDeleteEvent) { fired.Add(1) })

	ev := &events.WebhooksUpdateEvent{GuildID: guildID, ChannelID: channelID}
	client.internalEventHandler(nil, events.EventWebhooksUpdate, ev)

	// The fetch is attempted (handler registered) but fails; give it time to run.
	require.Eventually(t, func() bool { return atomic.LoadInt32(&calls) >= 1 }, 5*time.Second, 10*time.Millisecond,
		"a registered handler should trigger the REST fetch")
	time.Sleep(150 * time.Millisecond)

	assert.Equal(t, int32(0), fired.Load(), "a failed webhook fetch must not dispatch synthetic events")
	_, ok := client.loadWebhookSnapshot(channelID)
	assert.False(t, ok, "a failed fetch must not seed the snapshot")
}
