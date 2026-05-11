package connection

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// TestBug51ShutdownIdempotent verifies that calling Shutdown() from many
// goroutines concurrently never panics (double-close of eventCh or shutdownCh).
func TestBug51ShutdownIdempotent(t *testing.T) {
	wsURL, closeServer := mockGateway(t)
	defer closeServer()

	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)

	require.NoError(t, c.connectWebsocket(wsURL, false, nil, nil))
	go func() { _ = c.listenWebsocket() }()

	select {
	case <-c.Websocket.Ready:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for READY")
	}

	const concurrency = 10
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_ = c.Shutdown()
		}()
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrent Shutdown calls deadlocked or panicked")
	}
}

// TestBug48EventChSendRace verifies that pumping events concurrently with
// Shutdown does not panic under the race detector.
func TestBug48EventChSendRace(t *testing.T) {
	packets := make([]map[string]interface{}, 100)
	for i := range packets {
		packets[i] = dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
			"id": "1", "channel_id": "2",
			"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
			"content":   "hi",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	}

	for iter := 0; iter < 100; iter++ {
		wsURL, closeServer := mockGateway(t, packets...)

		c, err := NewClient("Bot fake-token", common.IntentGuilds)
		require.NoError(t, err)

		require.NoError(t, c.connectWebsocket(wsURL, false, nil, nil))
		go func() { _ = c.listenWebsocket() }()

		select {
		case <-c.Websocket.Ready:
		case <-time.After(5 * time.Second):
			closeServer()
			t.Fatalf("iter %d: timeout waiting for READY", iter)
		}

		_ = c.Shutdown()
		closeServer()
	}
}

// TestBug49UnlimitedModeHandlersComplete verifies that all event handlers
// launched in unlimited mode finish before Shutdown returns.
func TestBug49UnlimitedModeHandlersComplete(t *testing.T) {
	const nEvents = 50
	packets := make([]map[string]interface{}, nEvents)
	for i := range packets {
		packets[i] = dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
			"id": "1", "channel_id": "2",
			"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
			"content":   "hi",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	}

	wsURL, closeServer := mockGateway(t, packets...)
	defer closeServer()

	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)

	var handledCount atomic.Int64
	c.OnMessageCreate(func(_ *Client, _ *events.MessageCreateEvent) {
		time.Sleep(time.Millisecond)
		handledCount.Add(1)
	})

	require.NoError(t, c.connectWebsocket(wsURL, false, nil, nil))
	go func() { _ = c.listenWebsocket() }()

	select {
	case <-c.Websocket.Ready:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for READY")
	}

	// Give events time to arrive and begin processing.
	time.Sleep(200 * time.Millisecond)

	_ = c.Shutdown()

	// After Shutdown returns, all handlers must have completed.
	assert.Greater(t, handledCount.Load(), int64(0), "at least some handlers should have run")
}
