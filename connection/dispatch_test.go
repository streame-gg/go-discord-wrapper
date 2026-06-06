package connection

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── Mock gateway server ───────────────────────────────────────────────────────

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// mockGateway starts an in-process WebSocket server that sends a READY event
// followed by any extra packets, then holds the connection open.
//
// It returns the "ws://" URL and a close function.
func mockGateway(t *testing.T, extras ...map[string]interface{}) (wsURL string, closeFn func()) {
	t.Helper()

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

		// Drain IDENTIFY (OP 2) sent by the client.
		var ignore interface{}
		_ = conn.ReadJSON(&ignore)

		// OP 0 READY
		ready := map[string]interface{}{
			"op": 0, "t": "READY", "s": 1,
			"d": map[string]interface{}{
				"v":                  10,
				"user":               map[string]interface{}{"id": "1", "username": "testbot", "discriminator": "0000"},
				"guilds":             []interface{}{},
				"session_id":         "session_test",
				"resume_gateway_url": "wss://fake",
			},
		}
		_ = conn.WriteJSON(ready)

		// Brief pause so the test goroutine can register handlers before the
		// extra packets arrive — the mock server and test client run concurrently.
		time.Sleep(50 * time.Millisecond)

		// Extra test dispatch packets.
		for _, pkt := range extras {
			_ = conn.WriteJSON(pkt)
		}

		// Keep alive until client disconnects.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))

	return "ws" + ts.URL[len("http"):], ts.Close
}

// dispatchPacket builds an OP-0 dispatch packet for the given event type and data.
func dispatchPacket(eventType string, data interface{}) map[string]interface{} {
	raw, _ := json.Marshal(data)
	var d interface{}
	_ = json.Unmarshal(raw, &d)
	return map[string]interface{}{
		"op": 0,
		"t":  eventType,
		"s":  2,
		"d":  d,
	}
}

// launchClient connects to wsURL, starts the listener goroutine, and waits
// for the READY event before returning.  Call stop() to clean up.
func launchClient(t *testing.T, wsURL string, intents discord.Intent) (c *Client, stop func()) {
	t.Helper()

	var err error
	c, err = NewClient("Bot fake-token", intents)
	require.NoError(t, err)

	err = c.connectWebsocket(wsURL, false, nil, nil)
	require.NoError(t, err)

	// Start the listener first — it's what processes READY and closes Ready ch.
	go func() { _ = c.listenWebsocket() }()

	// Now wait for READY to be processed.
	select {
	case <-c.Ready():
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for READY")
	}

	return c, func() { c.Shutdown() }
}

// ── Event dispatch tests ──────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestDispatchFiresRegisteredHandler() {
	t := cs.T()
	msgPacket := dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
		"id": "123", "channel_id": "456",
		"author":    map[string]interface{}{"id": "1", "username": "u", "discriminator": "0"},
		"content":   "hello",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	wsURL, closeServer := mockGateway(t, msgPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var called atomic.Bool
	done := make(chan struct{})

	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		called.Store(true)
		close(done)
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for MESSAGE_CREATE handler")
	}

	assert.True(t, called.Load())
}

func (cs *ConnectionSuite) TestDispatchCallsAllRegisteredHandlers() {
	t := cs.T()
	msgPacket := dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "hi",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	wsURL, closeServer := mockGateway(t, msgPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var count atomic.Int32
	var wg sync.WaitGroup
	wg.Add(2)

	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		count.Add(1)
		wg.Done()
	})
	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		count.Add(1)
		wg.Done()
	})

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for both handlers")
	}

	assert.Equal(t, int32(2), count.Load())
}

// ── Middleware tests ──────────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestMiddlewareWrapsHandler() {
	t := cs.T()
	msgPacket := dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "hi",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	wsURL, closeServer := mockGateway(t, msgPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var mu sync.Mutex
	var order []string
	done := make(chan struct{})

	client.Use(func(next EventHandler) EventHandler {
		return func(c *Client, ev events.Event) {
			mu.Lock()
			order = append(order, "before")
			mu.Unlock()
			next(c, ev)
			mu.Lock()
			order = append(order, "after")
			mu.Unlock()
		}
	})

	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		mu.Lock()
		order = append(order, "handler")
		mu.Unlock()
		close(done)
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
	time.Sleep(20 * time.Millisecond) // let the "after" line execute

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, order, 3)
	assert.Equal(t, "before", order[0])
	assert.Equal(t, "handler", order[1])
	assert.Equal(t, "after", order[2])
}

func (cs *ConnectionSuite) TestMiddlewareCanShortCircuit() {
	t := cs.T()
	msgPacket := dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "hi",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	wsURL, closeServer := mockGateway(t, msgPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var handlerCalled atomic.Bool
	mwFired := make(chan struct{})

	client.Use(func(next EventHandler) EventHandler {
		return func(c *Client, ev events.Event) {
			// intentionally skip next — simulates a permission guard
			close(mwFired)
		}
	})

	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		handlerCalled.Store(true)
	})

	select {
	case <-mwFired:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for middleware")
	}

	time.Sleep(20 * time.Millisecond)
	assert.False(t, handlerCalled.Load(), "handler should not run when middleware short-circuits")
}

func (cs *ConnectionSuite) TestHandlerRegisteredBeforeUseIsNotWrapped() {
	t := cs.T()
	msgPacket := dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "hi",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	wsURL, closeServer := mockGateway(t, msgPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var mwCalled atomic.Bool
	done := make(chan struct{})

	// Handler registered BEFORE Use — must not be wrapped.
	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		close(done)
	})

	client.Use(func(next EventHandler) EventHandler {
		return func(c *Client, ev events.Event) {
			mwCalled.Store(true)
			next(c, ev)
		}
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	time.Sleep(10 * time.Millisecond)
	assert.False(t, mwCalled.Load(), "middleware added after handler must not wrap it")
}

// ── Benchmarks ────────────────────────────────────────────────────────────────

// BenchmarkDispatch measures raw event throughput through the dispatch loop,
// bypassing the WebSocket layer.
func BenchmarkDispatch(b *testing.B) {
	client, err := NewClient("Bot fake-token", discord.IntentGuilds)
	if err != nil {
		b.Fatal(err)
	}
	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {})

	raw, _ := json.Marshal(map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "bench",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	factory, _ := events.GetEventFactory(events.EventMessageCreate)

	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ev := factory()
		if err := json.Unmarshal(raw, ev); err != nil {
			b.Fatal(err)
		}
		wg.Add(1)
		go func(e events.Event) {
			defer wg.Done()
			client.dispatch(e)
		}(ev)
	}
	wg.Wait()
}

// BenchmarkDispatchWithMiddleware measures the overhead of a single pass-through middleware.
func BenchmarkDispatchWithMiddleware(b *testing.B) {
	client, err := NewClient("Bot fake-token", discord.IntentGuilds)
	if err != nil {
		b.Fatal(err)
	}
	client.Use(func(next EventHandler) EventHandler {
		return func(c *Client, ev events.Event) { next(c, ev) }
	})
	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {})

	raw, _ := json.Marshal(map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "bench",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	factory, _ := events.GetEventFactory(events.EventMessageCreate)

	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ev := factory()
		if err := json.Unmarshal(raw, ev); err != nil {
			b.Fatal(err)
		}
		wg.Add(1)
		go func(e events.Event) {
			defer wg.Done()
			client.dispatch(e)
		}(ev)
	}
	wg.Wait()
}
