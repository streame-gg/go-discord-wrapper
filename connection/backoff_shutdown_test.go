package connection

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mockWSAlwaysFail starts a WebSocket server that upgrades the connection and
// immediately closes it so every NewWebsocket attempt returns an error.
func mockWSAlwaysFail(t *testing.T) (wsURL string, closeFn func()) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close()
	}))
	return "ws" + ts.URL[len("http"):], ts.Close
}

// TestP0_4_ShutdownDuringBackoff verifies that calling Shutdown while
// reconnect is sleeping in its exponential-backoff wait returns in well under
// 100 ms — not after the full backoff duration (≥ 1 s without the fix).
//
// The sequence is:
//  1. Establish a real initial connection via mockGateway so the Client has a
//     valid Websocket struct with a ReconnectURL.
//  2. Override ReconnectURL to point to a server that always fails so every
//     reconnect attempt returns an error immediately.
//  3. Call reconnect in a goroutine with unlimited retries (-1).
//  4. Wait 200 ms — enough for the first attempt to fail and the second to
//     enter the backoff select (backoff ≥ 1 s without the fix).
//  5. Call Shutdown and assert it returns in < 100 ms.
func TestP0_4_ShutdownDuringBackoff(t *testing.T) {
	// Step 1: initial good gateway so we have a valid Websocket.
	goodURL, closeGood := mockGateway(t)
	defer closeGood()

	// Step 2: always-failing gateway that triggers repeated reconnect failures.
	failURL, closeFail := mockWSAlwaysFail(t)
	defer closeFail()

	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	require.NoError(t, err)

	require.NoError(t, c.connectWebsocket(goodURL, false, nil, nil))
	go func() { _ = c.listenWebsocket() }()

	// Wait for READY so the Websocket struct is fully set up.
	select {
	case <-c.Ready():
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for initial READY")
	}

	// Point ReconnectURL at the always-failing server.
	c.wsMu.Lock()
	c.Websocket.ReconnectURL = &failURL
	c.wsMu.Unlock()

	// Unlimited retries so reconnect keeps looping and always enters backoff.
	c.maxReconnectRetries = -1

	var reconnectDone atomic.Bool
	go func() {
		_ = c.reconnect(false)
		reconnectDone.Store(true)
	}()

	// Allow the first attempt to fail and the second to enter the backoff sleep.
	time.Sleep(300 * time.Millisecond)

	start := time.Now()
	_ = c.Shutdown()
	elapsed := time.Since(start)

	assert.Less(t, elapsed, 100*time.Millisecond,
		"Shutdown must return in < 100 ms while reconnect is in backoff (got %v)", elapsed)

	// Reconnect goroutine must exit after Shutdown closes shutdownCh.
	time.Sleep(50 * time.Millisecond)
	assert.True(t, reconnectDone.Load(), "reconnect goroutine must exit after Shutdown")
}
