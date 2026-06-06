package connection

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// launchClientWithCache is launchClient with a memory cache attached, so
// cache-dependent paths (notably synthetic expression events) are exercised.
func launchClientWithCache(t require.TestingT, wsURL string, intents discord.Intent) (c *Client, stop func()) {
	return launchClientWithCacheAndServer(t, wsURL, "", intents)
}

// launchClientWithCacheAndServer is launchClientWithCache that also points the
// client's REST calls at baseURL. It is needed for synthetic events whose
// derivation performs a REST fetch (notably WEBHOOKS_UPDATE). An empty baseURL
// leaves the default Discord endpoint in place.
func launchClientWithCacheAndServer(t require.TestingT, wsURL, baseURL string, intents discord.Intent) (c *Client, stop func()) {
	mc := cache.NewMemoryCache(cache.Options{})
	opts := []options.Option{options.WithCache(mc)}
	if baseURL != "" {
		opts = append(opts, options.WithBaseURL(baseURL))
	}
	var err error
	c, err = NewClient("Bot fake-token", intents, opts...)
	require.NoError(t, err)
	require.NoError(t, c.connectWebsocket(wsURL, false, nil, nil))
	go func() { _ = c.listenWebsocket() }()
	select {
	case <-c.Ready():
	case <-time.After(5 * time.Second):
		t.Errorf("timeout waiting for READY")
		t.FailNow()
	}
	return c, func() { c.Shutdown() }
}

// pushGateway is mockGateway with on-demand packet delivery: instead of firing a
// fixed list of extras right after READY, it returns a send func the test calls
// to push dispatch packets one at a time. This lets a test observe the effect of
// one packet (e.g. a seeded webhook snapshot) before sending the next.
func pushGateway(t *testing.T) (wsURL string, send func(map[string]interface{}), closeFn func()) {
	t.Helper()

	pkts := make(chan map[string]interface{}, 8)
	done := make(chan struct{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// OP 10 HELLO
		_ = conn.WriteJSON(map[string]interface{}{
			"op": 10,
			"d":  map[string]interface{}{"heartbeat_interval": 60000},
		})

		// Drain IDENTIFY (OP 2).
		var ignore interface{}
		_ = conn.ReadJSON(&ignore)

		// OP 0 READY
		_ = conn.WriteJSON(map[string]interface{}{
			"op": 0, "t": "READY", "s": 1,
			"d": map[string]interface{}{
				"v":                  10,
				"user":               map[string]interface{}{"id": "1", "username": "testbot", "discriminator": "0000"},
				"guilds":             []interface{}{},
				"session_id":         "session_test",
				"resume_gateway_url": "wss://fake",
			},
		})

		// Drain client reads (heartbeats / close) so the socket stays healthy.
		connClosed := make(chan struct{})
		go func() {
			defer close(connClosed)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()

		for {
			select {
			case pkt := <-pkts:
				if err := conn.WriteJSON(pkt); err != nil {
					return
				}
			case <-done:
				return
			case <-connClosed:
				return
			}
		}
	}))

	var once sync.Once
	return "ws" + ts.URL[len("http"):],
		func(pkt map[string]interface{}) { pkts <- pkt },
		func() { once.Do(func() { close(done) }); ts.Close() }
}

// TestEmojiSyntheticDispatch_FullPath_NoDuplicates feeds two GUILD_EMOJIS_UPDATE
// packets through the real gateway pipeline. The first warms the cache (cold ->
// stored, no synthetic); the second mutates the set (one added, one removed, one
// renamed). It asserts the raw event fires once per packet, each synthetic event
// fires exactly once for the right emoji, no synthetic event duplicates, and the
// emoji update does not leak into sticker handlers.
func (cs *ConnectionSuite) TestEmojiSyntheticDispatch_FullPath_NoDuplicates() {
	t := cs.T()
	const gid = "100000000000000001"
	emoji := func(id, name string) map[string]any {
		return map[string]any{"id": id, "name": name, "available": true}
	}

	p1 := dispatchPacket("GUILD_EMOJIS_UPDATE", map[string]any{
		"guild_id": gid,
		"emojis": []any{
			emoji("200000000000000001", "alpha"),
			emoji("200000000000000002", "beta"),
		},
	})
	p2 := dispatchPacket("GUILD_EMOJIS_UPDATE", map[string]any{
		"guild_id": gid,
		"emojis": []any{
			emoji("200000000000000001", "alpha_renamed"), // updated
			emoji("200000000000000003", "gamma"),         // added; beta removed
		},
	})

	wsURL, closeServer := mockGateway(t, p1, p2)
	defer closeServer()

	client, stop := launchClientWithCache(t, wsURL, discord.IntentGuilds)
	defer stop()

	var raw, add, rem, upd, stickerCross atomic.Int32
	var mu sync.Mutex
	var addedIDs, removedIDs, updatedIDs []string

	client.OnGuildEmojisUpdate(func(_ *Client, _ *events.GuildEmojisUpdateEvent) { raw.Add(1) })
	client.OnGuildEmojiAdd(func(_ *Client, e *events.GuildEmojiAddEvent) {
		add.Add(1)
		mu.Lock()
		addedIDs = append(addedIDs, e.Emoji.ID.String())
		mu.Unlock()
	})
	client.OnGuildEmojiRemove(func(_ *Client, e *events.GuildEmojiRemoveEvent) {
		rem.Add(1)
		mu.Lock()
		removedIDs = append(removedIDs, e.Emoji.ID.String())
		mu.Unlock()
	})
	client.OnGuildEmojiUpdate(func(_ *Client, e *events.GuildEmojiUpdateEvent) {
		upd.Add(1)
		mu.Lock()
		updatedIDs = append(updatedIDs, e.NewEmoji.ID.String())
		mu.Unlock()
	})
	// Cross-talk guard: an emoji update must never invoke sticker handlers.
	client.OnGuildStickerAdd(func(_ *Client, _ *events.GuildStickerAddEvent) { stickerCross.Add(1) })
	client.OnGuildStickerUpdate(func(_ *Client, _ *events.GuildStickerUpdateEvent) { stickerCross.Add(1) })
	client.OnGuildStickerRemove(func(_ *Client, _ *events.GuildStickerRemoveEvent) { stickerCross.Add(1) })

	require.Eventually(t, func() bool {
		return raw.Load() == 2 && add.Load() == 1 && rem.Load() == 1 && upd.Load() == 1
	}, 5*time.Second, 10*time.Millisecond, "expected raw=2 add=1 rem=1 upd=1")

	// Settle, then assert nothing fired twice.
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, int32(2), raw.Load(), "raw GUILD_EMOJIS_UPDATE fires once per packet")
	assert.Equal(t, int32(1), add.Load(), "GuildEmojiAdd must fire exactly once")
	assert.Equal(t, int32(1), rem.Load(), "GuildEmojiRemove must fire exactly once")
	assert.Equal(t, int32(1), upd.Load(), "GuildEmojiUpdate must fire exactly once")
	assert.Equal(t, int32(0), stickerCross.Load(), "emoji update must not trigger sticker handlers")

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"200000000000000003"}, addedIDs)
	assert.Equal(t, []string{"200000000000000002"}, removedIDs)
	assert.Equal(t, []string{"200000000000000001"}, updatedIDs)
}

// TestStickerSyntheticDispatch_FullPath_NoDuplicates is the sticker analogue of
// the emoji test above.
func (cs *ConnectionSuite) TestStickerSyntheticDispatch_FullPath_NoDuplicates() {
	t := cs.T()
	const gid = "100000000000000001"
	sticker := func(id, name string) map[string]any {
		return map[string]any{"id": id, "name": name, "tags": "tag", "type": 2, "format_type": 1, "available": true}
	}

	p1 := dispatchPacket("GUILD_STICKERS_UPDATE", map[string]any{
		"guild_id": gid,
		"stickers": []any{
			sticker("300000000000000001", "one"),
			sticker("300000000000000002", "two"),
		},
	})
	p2 := dispatchPacket("GUILD_STICKERS_UPDATE", map[string]any{
		"guild_id": gid,
		"stickers": []any{
			sticker("300000000000000001", "one_renamed"), // updated
			sticker("300000000000000003", "three"),       // added; two removed
		},
	})

	wsURL, closeServer := mockGateway(t, p1, p2)
	defer closeServer()

	client, stop := launchClientWithCache(t, wsURL, discord.IntentGuilds)
	defer stop()

	var raw, add, rem, upd atomic.Int32
	var mu sync.Mutex
	var addedIDs, removedIDs, updatedIDs []string

	client.OnGuildStickersUpdate(func(_ *Client, _ *events.GuildStickersUpdateEvent) { raw.Add(1) })
	client.OnGuildStickerAdd(func(_ *Client, e *events.GuildStickerAddEvent) {
		add.Add(1)
		mu.Lock()
		addedIDs = append(addedIDs, e.Sticker.ID.String())
		mu.Unlock()
	})
	client.OnGuildStickerRemove(func(_ *Client, e *events.GuildStickerRemoveEvent) {
		rem.Add(1)
		mu.Lock()
		removedIDs = append(removedIDs, e.Sticker.ID.String())
		mu.Unlock()
	})
	client.OnGuildStickerUpdate(func(_ *Client, e *events.GuildStickerUpdateEvent) {
		upd.Add(1)
		mu.Lock()
		updatedIDs = append(updatedIDs, e.NewSticker.ID.String())
		mu.Unlock()
	})

	require.Eventually(t, func() bool {
		return raw.Load() == 2 && add.Load() == 1 && rem.Load() == 1 && upd.Load() == 1
	}, 5*time.Second, 10*time.Millisecond, "expected raw=2 add=1 rem=1 upd=1")

	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, int32(2), raw.Load())
	assert.Equal(t, int32(1), add.Load())
	assert.Equal(t, int32(1), rem.Load())
	assert.Equal(t, int32(1), upd.Load())

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"300000000000000003"}, addedIDs)
	assert.Equal(t, []string{"300000000000000002"}, removedIDs)
	assert.Equal(t, []string{"300000000000000001"}, updatedIDs)
}

// TestColdCacheSyntheticDispatch_NoSyntheticEvents verifies that the first
// GUILD_EMOJIS_UPDATE on a cold cache fires only the raw event — no synthetic
// add/remove/update — since there is no prior state to diff against.
func (cs *ConnectionSuite) TestColdCacheSyntheticDispatch_NoSyntheticEvents() {
	t := cs.T()
	const gid = "100000000000000009"

	p1 := dispatchPacket("GUILD_EMOJIS_UPDATE", map[string]any{
		"guild_id": gid,
		"emojis":   []any{map[string]any{"id": "200000000000000099", "name": "solo", "available": true}},
	})

	wsURL, closeServer := mockGateway(t, p1)
	defer closeServer()

	client, stop := launchClientWithCache(t, wsURL, discord.IntentGuilds)
	defer stop()

	var raw, synthetic atomic.Int32
	client.OnGuildEmojisUpdate(func(_ *Client, _ *events.GuildEmojisUpdateEvent) { raw.Add(1) })
	client.OnGuildEmojiAdd(func(_ *Client, _ *events.GuildEmojiAddEvent) { synthetic.Add(1) })
	client.OnGuildEmojiRemove(func(_ *Client, _ *events.GuildEmojiRemoveEvent) { synthetic.Add(1) })
	client.OnGuildEmojiUpdate(func(_ *Client, _ *events.GuildEmojiUpdateEvent) { synthetic.Add(1) })

	require.Eventually(t, func() bool { return raw.Load() == 1 }, 5*time.Second, 10*time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, int32(1), raw.Load(), "raw event fires on cold cache")
	assert.Equal(t, int32(0), synthetic.Load(), "no synthetic events on a cold cache")
}

// TestWebhookSyntheticDispatch_FullGatewayPath_NoDuplicates is the webhook
// analogue of the emoji/sticker tests above, driven through the *real* gateway
// pipeline (decode + dispatch) rather than internalEventHandler. Discord's
// WEBHOOKS_UPDATE carries only the channel, so each packet triggers a REST fetch
// of the channel's webhooks which is diffed against the previous snapshot. The
// first packet warms the snapshot (cold -> seed, no synthetic); the second
// mutates the set (one added, one removed, one renamed). It asserts each
// synthetic event fires exactly once with no duplicates, the renamed webhook's
// Update carries the previous name, and the coarse handler fires once per packet.
func (cs *ConnectionSuite) TestWebhookSyntheticDispatch_FullGatewayPath_NoDuplicates() {
	t := cs.T()
	const guildID = "175928847299117001"
	const channelID = "175928847299117063"

	// REST: first fetch seeds {a, b}; second mutates to {a_renamed, c} (b removed).
	var calls int32
	rest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		if n == 1 {
			_, _ = w.Write([]byte(`[{"id":"500000000000000001","type":1,"name":"a"},{"id":"500000000000000002","type":1,"name":"b"}]`))
		} else {
			_, _ = w.Write([]byte(`[{"id":"500000000000000001","type":1,"name":"a_renamed"},{"id":"500000000000000003","type":1,"name":"c"}]`))
		}
	}))
	defer rest.Close()

	wsURL, send, closeServer := pushGateway(t)
	defer closeServer()

	client, stop := launchClientWithCacheAndServer(t, wsURL, rest.URL, discord.IntentGuilds)
	defer stop()

	var coarse, created, updated, deleted atomic.Int32
	var mu sync.Mutex
	var createdIDs, deletedIDs []string
	var oldName string

	client.OnWebhooksUpdate(func(_ *Client, _ *events.WebhooksUpdateEvent) { coarse.Add(1) })
	client.OnWebhookCreate(func(_ *Client, e *events.WebhookCreateEvent) {
		created.Add(1)
		mu.Lock()
		createdIDs = append(createdIDs, e.Webhook.ID.String())
		mu.Unlock()
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
		mu.Lock()
		deletedIDs = append(deletedIDs, e.Webhook.ID.String())
		mu.Unlock()
	})

	pkt := dispatchPacket("WEBHOOKS_UPDATE", map[string]any{
		"guild_id":   guildID,
		"channel_id": channelID,
	})

	// First packet: cold snapshot -> seed only, no synthetic events.
	send(pkt)
	require.Eventually(t, func() bool {
		_, ok := client.loadWebhookSnapshot(discord.Snowflake(mustU64(channelID)))
		return ok
	}, 5*time.Second, 10*time.Millisecond, "first fetch should seed the snapshot")
	assert.Equal(t, int32(0), created.Load()+updated.Load()+deleted.Load(), "cold snapshot must not fire synthetic events")

	// Second packet: diff against the seeded snapshot.
	send(pkt)
	require.Eventually(t, func() bool {
		return created.Load() == 1 && updated.Load() == 1 && deleted.Load() == 1
	}, 5*time.Second, 10*time.Millisecond, "expected create=1 update=1 delete=1")

	// Settle, then assert nothing fired twice.
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, int32(2), coarse.Load(), "raw WEBHOOKS_UPDATE fires once per packet")
	assert.Equal(t, int32(1), created.Load(), "WebhookCreate must fire exactly once")
	assert.Equal(t, int32(1), updated.Load(), "WebhookUpdate must fire exactly once")
	assert.Equal(t, int32(1), deleted.Load(), "WebhookDelete must fire exactly once")

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"500000000000000003"}, createdIDs)
	assert.Equal(t, []string{"500000000000000002"}, deletedIDs)
	assert.Equal(t, "a", oldName, "WebhookUpdate must carry the previously cached name")
}
