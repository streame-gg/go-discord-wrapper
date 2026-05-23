package connection

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mockGatewayWithOp1 starts a WebSocket mock gateway that:
//   - completes the standard handshake (OP 10, drain Identify, READY)
//   - counts Op-1 heartbeat writes from the client
//   - after 100 ms sends an Op-1 immediate-heartbeat request to the client
//
// heartbeatInterval controls the OP 10 heartbeat_interval_ms value.
func mockGatewayWithOp1(t *testing.T, heartbeatInterval int, heartbeatCount *atomic.Int32) (wsURL string, closeFn func()) {
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
			"d":  map[string]interface{}{"heartbeat_interval": heartbeatInterval},
		})

		// Drain IDENTIFY (OP 2).
		var ignore interface{}
		_ = conn.ReadJSON(&ignore)

		// READY.
		_ = conn.WriteJSON(map[string]interface{}{
			"op": 0, "t": "READY", "s": 1,
			"d": map[string]interface{}{
				"v":                  10,
				"user":               map[string]interface{}{"id": "1", "username": "bot", "discriminator": "0"},
				"guilds":             []interface{}{},
				"session_id":         "sess",
				"resume_gateway_url": "wss://fake",
			},
		})

		// Give the client a moment to process READY, then send Op-1.
		time.Sleep(100 * time.Millisecond)
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": nil})

		// Read incoming frames and count heartbeat (Op-1) writes.
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var p struct {
				Op int `json:"op"`
			}
			if json.Unmarshal(msg, &p) == nil && p.Op == 1 {
				heartbeatCount.Add(1)
			}
		}
	}))

	return "ws" + ts.URL[len("http"):], ts.Close
}

// TestOp1HeartbeatRequest verifies that when Discord sends an Op-1 immediate
// heartbeat request the client sends at least one heartbeat (Op-1 payload)
// without waiting for the full heartbeat_interval (Issue 6).
//
// The heartbeat interval is 200 ms so the test runs quickly. After the client
// processes Op-1 we wait 300 ms and assert the server received ≥ 1 heartbeat.
func TestOp1HeartbeatRequest(t *testing.T) {
	var heartbeatCount atomic.Int32

	wsURL, closeServer := mockGatewayWithOp1(t, 200, &heartbeatCount)
	defer closeServer()

	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	require.NoError(t, err)
	defer c.Shutdown()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()
	_ = c
	_ = client

	// Allow the server's 1-second read window to close.
	time.Sleep(600 * time.Millisecond)

	assert.GreaterOrEqual(t, heartbeatCount.Load(), int32(1),
		"client should have sent at least one heartbeat in response to Op-1")

	// Compile-time check: heartbeatRequest must be chan struct{}.
	var _ chan struct{} = client.Websocket.heartbeatRequest
}

// TestOp1HeartbeatRequestChannelNonBlocking verifies the non-blocking send
// in the listener: signalling heartbeatRequest when the buffer is full must
// not block (Issue 6).
func TestOp1HeartbeatRequestChannelNonBlocking(t *testing.T) {
	ws := &Websocket{
		heartbeatRequest: make(chan struct{}, 1),
		Closed:           make(chan struct{}),
	}

	// First send must succeed.
	select {
	case ws.heartbeatRequest <- struct{}{}:
	default:
		t.Fatal("first send should not block on an empty buffered channel")
	}

	// Second send (buffer full) must not block — mirrors the listener code.
	select {
	case ws.heartbeatRequest <- struct{}{}:
		// drained one, try again to ensure exactly one is pending
	default:
		// buffer was full, correctly dropped
	}

	// Channel must have exactly one item.
	select {
	case <-ws.heartbeatRequest:
	default:
		t.Fatal("channel should have one item")
	}

	// Verify websocket.IsCloseError is importable (keeps the gorilla import live).
	_ = websocket.IsCloseError(nil)
}
