# Testing event dispatch (integration suite)

This is a **contributor** guide. It documents the in-process harness that drives
events end-to-end through the real gateway pipeline — JSON decode, dispatch,
cache side-effects, synthetic-event derivation, and handler invocation — without
a live token or a network connection to Discord.

Use it when you add or change a gateway event, a cache side-effect, or a
synthetic event, and want a test that proves the whole path works (not just the
struct decode).

All of this lives in the `connection` package's `*_test.go` files. The synthetic
event tests are in `connection/event_dispatch_integration_test.go` and
`connection/synthetic_webhooks_test.go`.

## The harness

Tests are methods on `ConnectionSuite` (`connection/connection_suite_test.go`),
run via `testify/suite`. To add a test, write a new
`func (cs *ConnectionSuite) TestXxx()` method — it is picked up automatically.

The harness provides an in-process WebSocket "gateway" plus client launchers:

| Helper | Purpose |
|--------|---------|
| `mockGateway(t, extras...)` | Starts a WS server that sends `HELLO` + `READY`, then fires each `extras` packet once, in order, right after `READY`. Returns `wsURL, closeFn`. |
| `pushGateway(t)` | Like `mockGateway` but with **on-demand** delivery: returns `wsURL, send, closeFn`. Call `send(pkt)` to push one packet when you want it. Use this when packet *N+1* depends on the observable effect of packet *N*. |
| `dispatchPacket(eventType, data)` | Builds an OP-0 dispatch frame (`{"op":0,"t":eventType,"d":data}`) suitable for `mockGateway`/`send`. `data` is any value that marshals to the event's JSON. |
| `launchClient(t, wsURL, intents)` | Connects a plain client, starts the listener, and blocks until `READY`. Returns `client, stop`. |
| `launchClientWithCache(t, wsURL, intents)` | As above, with a `MemoryCache` attached. Required for any event with a cache side-effect or a **cache-derived** synthetic event (emoji/sticker). |
| `launchClientWithCacheAndServer(t, wsURL, baseURL, intents)` | As above, and points the REST client at `baseURL`. Required for **REST-derived** synthetic events (webhooks), whose derivation fetches over REST. |

Always `defer stop()` and `defer closeServer()`.

### Asserting "fired exactly once" / no cross-talk

Synthetic derivation runs on background goroutines, so tests are
eventually-consistent. The standard pattern:

1. Count invocations with `atomic.Int32` per handler (collect IDs under a
   `sync.Mutex` if you need to assert *which* entities fired).
2. `require.Eventually(...)` until the expected counts are reached.
3. `time.Sleep(150 * time.Millisecond)` to let any *erroneous* extra dispatch
   land.
4. `assert.Equal` the exact counts — this is what catches duplicates.
5. Register a handler for a *sibling* event (e.g. stickers while testing emojis)
   and assert it stayed at `0` to catch cross-talk between dispatch buckets.

## Recipe 1 — a plain (raw) gateway event

For an event that is just decoded and handed to handlers (optionally with a
cache side-effect). Example: `GUILD_INTEGRATIONS_UPDATE`
(`connection/guild_integrations_update_test.go`).

```go
func (cs *ConnectionSuite) TestMyEventDispatch() {
	t := cs.T()

	packet := dispatchPacket("MY_EVENT", map[string]any{"guild_id": "123"})

	wsURL, closeServer := mockGateway(t, packet)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	got := make(chan *events.MyEvent, 1)
	client.OnMyEvent(func(_ *Client, ev *events.MyEvent) { got <- ev })

	select {
	case ev := <-got:
		assert.Equal(t, discord.Snowflake(123), ev.GuildID)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for MY_EVENT handler")
	}
}
```

If the event mutates the cache, use `launchClientWithCache` and assert on the
cached entity after the handler fires.

## Recipe 2 — a cache-derived synthetic event (emoji / sticker)

Discord sends the coarse `GUILD_EMOJIS_UPDATE` / `GUILD_STICKERS_UPDATE` (the
whole set). The wrapper diffs it against the **cached** previous set to emit
granular `WRAPPER_GUILD_EMOJI_ADD/REMOVE/UPDATE` events. So the test needs a
cache, two packets (warm, then mutate), and the cold-cache assertion.
See `TestEmojiSyntheticDispatch_FullPath_NoDuplicates`.

```go
func (cs *ConnectionSuite) TestExprSynthetic() {
	t := cs.T()
	const gid = "100000000000000001"

	p1 := dispatchPacket("GUILD_EMOJIS_UPDATE", map[string]any{ /* {alpha, beta} */ })
	p2 := dispatchPacket("GUILD_EMOJIS_UPDATE", map[string]any{ /* {alpha_renamed, gamma} */ })

	wsURL, closeServer := mockGateway(t, p1, p2) // both fire after READY
	defer closeServer()

	client, stop := launchClientWithCache(t, wsURL, discord.IntentGuilds)
	defer stop()

	var raw, add, rem, upd, stickerCross atomic.Int32
	client.OnGuildEmojisUpdate(func(_ *Client, _ *events.GuildEmojisUpdateEvent) { raw.Add(1) })
	client.OnGuildEmojiAdd(func(_ *Client, _ *events.GuildEmojiAddEvent) { add.Add(1) })
	client.OnGuildEmojiRemove(func(_ *Client, _ *events.GuildEmojiRemoveEvent) { rem.Add(1) })
	client.OnGuildEmojiUpdate(func(_ *Client, _ *events.GuildEmojiUpdateEvent) { upd.Add(1) })
	// Cross-talk guard: an emoji update must never reach sticker handlers.
	client.OnGuildStickerAdd(func(_ *Client, _ *events.GuildStickerAddEvent) { stickerCross.Add(1) })

	require.Eventually(t, func() bool {
		return raw.Load() == 2 && add.Load() == 1 && rem.Load() == 1 && upd.Load() == 1
	}, 5*time.Second, 10*time.Millisecond)

	time.Sleep(150 * time.Millisecond) // settle, then prove no duplicates
	assert.Equal(t, int32(0), stickerCross.Load())
}
```

`mockGateway(t, p1, p2)` is fine here because the cache write for `p1` happens on
the dispatch path *before* `p2` is processed — packet ordering on the socket is
enough. (Compare Recipe 3, which is not.)

For the cold path, send a single packet and assert the synthetic counters stay at
`0` — the first set for a guild only seeds the cache. See
`TestColdCacheSyntheticDispatch_NoSyntheticEvents`.

## Recipe 3 — a REST-derived synthetic event (webhooks)

`WEBHOOKS_UPDATE` carries only the channel, so the wrapper **fetches** the
channel's webhooks over REST and diffs them against the previous snapshot to emit
`WRAPPER_WEBHOOK_CREATE/UPDATE/DELETE`. This adds two requirements over Recipe 2:

- A mock **REST** server, wired via `launchClientWithCacheAndServer(..., baseURL, ...)`.
- `pushGateway` instead of `mockGateway`: the first `WEBHOOKS_UPDATE` only seeds
  the snapshot, and that seeding completes on a background goroutine. The second
  packet must not be sent until the snapshot exists, or the diff has nothing to
  compare against. `pushGateway` lets the test wait between packets.

See `TestWebhookSyntheticDispatch_FullGatewayPath_NoDuplicates`.

```go
func (cs *ConnectionSuite) TestWebhookSyntheticFullPath() {
	t := cs.T()
	const channelID = "175928847299117063"

	// REST: first fetch seeds {a,b}; second returns {a_renamed,c} (b removed).
	var calls int32
	rest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if atomic.AddInt32(&calls, 1) == 1 {
			_, _ = w.Write([]byte(`[{"id":"500...01","type":1,"name":"a"},{"id":"500...02","type":1,"name":"b"}]`))
		} else {
			_, _ = w.Write([]byte(`[{"id":"500...01","type":1,"name":"a_renamed"},{"id":"500...03","type":1,"name":"c"}]`))
		}
	}))
	defer rest.Close()

	wsURL, send, closeServer := pushGateway(t)
	defer closeServer()

	client, stop := launchClientWithCacheAndServer(t, wsURL, rest.URL, discord.IntentGuilds)
	defer stop()

	var created, updated, deleted atomic.Int32
	client.OnWebhookCreate(func(_ *Client, _ *events.WebhookCreateEvent) { created.Add(1) })
	client.OnWebhookUpdate(func(_ *Client, _ *events.WebhookUpdateEvent) { updated.Add(1) })
	client.OnWebhookDelete(func(_ *Client, _ *events.WebhookDeleteEvent) { deleted.Add(1) })

	pkt := dispatchPacket("WEBHOOKS_UPDATE", map[string]any{"channel_id": channelID})

	// 1) Cold snapshot — seed only, no synthetic events.
	send(pkt)
	require.Eventually(t, func() bool {
		_, ok := client.loadWebhookSnapshot(discord.Snowflake(mustU64(channelID)))
		return ok
	}, 5*time.Second, 10*time.Millisecond)
	assert.Equal(t, int32(0), created.Load()+updated.Load()+deleted.Load())

	// 2) Diff against the seeded snapshot.
	send(pkt)
	require.Eventually(t, func() bool {
		return created.Load() == 1 && updated.Load() == 1 && deleted.Load() == 1
	}, 5*time.Second, 10*time.Millisecond)

	time.Sleep(150 * time.Millisecond) // no duplicates
	assert.Equal(t, int32(1), created.Load())
}
```

The opt-in guard is worth a separate, smaller test: with no webhook handler
registered, a `WEBHOOKS_UPDATE` must trigger **no** REST fetch (assert the mock
server's call count stayed `0`). See `TestWebhookSynthetic_NoHandlers_NoFetch`.

## Pure-diff unit tests

The diff logic underneath the synthetic events is a pure function and should also
be tested directly, without the gateway — fast and exhaustive on edge cases
(add / update / delete / no-op). See `TestDiffWebhooks` for the pattern.

## Running

```sh
# the whole integration suite
go test ./connection/

# one test
go test ./connection/ -run 'ConnectionSuite/TestWebhookSyntheticDispatch_FullGatewayPath_NoDuplicates' -v

# always run synthetic/concurrent tests under the race detector before a PR
go test ./connection/ -race
```

## Checklist for a new event

- [ ] Decode/roundtrip unit test in `types/events` (struct ↔ JSON).
- [ ] Dispatch integration test here (Recipe 1, 2, or 3) proving handlers fire.
- [ ] If it derives synthetic events: cold-state assertion + exact-count
      (no-duplicate) assertion + a sibling cross-talk guard.
- [ ] If it derives from a diff: a pure-diff unit test.
- [ ] `go test ./connection/ -race` and `go vet ./...` are clean.
