package connection

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// mockGatewayThenClose starts a WS server that sends HELLO → reads IDENTIFY →
// sends READY → pauses pauseMs → sends a close frame with closeCode.
func mockGatewayThenClose(t *testing.T, closeCode int) (wsURL string, closeFn func()) {
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

		// Drain IDENTIFY
		var ignore interface{}
		_ = conn.ReadJSON(&ignore)

		// OP 0 READY
		_ = conn.WriteJSON(map[string]interface{}{
			"op": 0, "t": "READY", "s": 1,
			"d": map[string]interface{}{
				"v":                  10,
				"user":               map[string]interface{}{"id": "1", "username": "bot", "discriminator": "0000"},
				"guilds":             []interface{}{},
				"session_id":         "sess",
				"resume_gateway_url": "wss://fake",
			},
		})

		// Give client time to process READY before closing.
		time.Sleep(80 * time.Millisecond)

		_ = conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(closeCode, ""))
		// Hold connection briefly so the client can read the close frame.
		time.Sleep(200 * time.Millisecond)
	}))
	return "ws" + ts.URL[len("http"):], ts.Close
}

// TestBug17NonRecoverableCloseCodesExitListenerLoop verifies that Discord close
// codes 4004, 4010, 4011, 4012, 4013, 4014 cause listenWebsocket to return a
// correctly-typed error, so the Login goroutine can classify them as terminal
// rather than triggering a reconnect loop.
func TestBug17NonRecoverableCloseCodesExitListenerLoop(t *testing.T) {
	nonRecoverable := []int{4004, 4010, 4011, 4012, 4013, 4014}

	for _, code := range nonRecoverable {
		code := code
		t.Run(http.StatusText(0)+"_"+string(rune('0'+code/1000)), func(t *testing.T) {
			wsURL, closeServer := mockGatewayThenClose(t, code)
			defer closeServer()

			c, err := NewClient("Bot fake-token", common.IntentGuilds)
			require.NoError(t, err)

			require.NoError(t, c.connectWebsocket(wsURL, false, nil, nil))

			// listenWebsocket closes the Ready channel internally when it processes
			// the READY event, so start it first.
			listenerErr := make(chan error, 1)
			go func() {
				listenerErr <- c.listenWebsocket()
			}()

			// Wait for READY before asserting.
			select {
			case <-c.Websocket.Ready:
			case <-time.After(5 * time.Second):
				t.Fatalf("code %d: timeout waiting for READY", code)
			}

			// After the server closes with the non-recoverable code, listenWebsocket
			// must return within a short timeout.
			select {
			case err := <-listenerErr:
				// The returned error must be classified as non-recoverable.
				assert.True(t,
					websocket.IsCloseError(err, 4004, 4010, 4011, 4012, 4013, 4014),
					"code %d: error should be classified as non-recoverable, got: %v", code, err)
			case <-time.After(3 * time.Second):
				t.Fatalf("code %d: listenWebsocket did not exit after non-recoverable close", code)
			}

			c.Shutdown()
		})
	}
}

// TestBug18CorruptReadyPayloadDoesNotCloseReadyChan verifies that a corrupted
// READY payload causes internalEventHandler to return false without closing
// the Ready channel or modifying the cache.
func TestBug18CorruptReadyPayloadDoesNotCloseReadyChan(t *testing.T) {
	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)

	// Construct a minimal websocket so we can call internalEventHandler directly.
	c.Websocket = &Websocket{
		Ready:  make(chan struct{}),
		Closed: make(chan struct{}),
	}

	corruptPayload := []byte(`{invalid json`)

	result := c.internalEventHandler(corruptPayload, "READY", nil)
	assert.False(t, result, "internalEventHandler must return false on READY unmarshal failure")

	// Ready channel must NOT be closed.
	select {
	case <-c.Websocket.Ready:
		t.Fatal("Ready channel was closed despite corrupt READY payload")
	default:
		// correct — channel is still open
	}
}

// TestBug16NoHeartbeatLeakOnWriteFailure verifies that if the Identify/Resume
// write fails, no goroutine is left running. The heartbeat goroutine must not
// start until after a successful write.
func TestBug16NoHeartbeatLeakOnWriteFailure(t *testing.T) {
	// Server: sends HELLO then immediately closes the connection before IDENTIFY
	// can be written, so the WriteJSON in NewWebsocket will fail.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// Send HELLO
		_ = conn.WriteJSON(map[string]interface{}{
			"op": 10,
			"d":  map[string]interface{}{"heartbeat_interval": 60000},
		})
		// Close immediately — the client's WriteJSON(IDENTIFY) will fail.
		conn.Close()
	}))
	defer ts.Close()

	wsURL := "ws" + ts.URL[len("http"):]

	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)

	before := runtime.NumGoroutine()

	// NewWebsocket (called inside connectWebsocket) should fail cleanly.
	_ = c.connectWebsocket(wsURL, false, nil, nil)

	// Give any stray goroutines time to start (and then they shouldn't).
	time.Sleep(200 * time.Millisecond)

	after := runtime.NumGoroutine()

	// Allow a small tolerance for background goroutines from the test runtime.
	const tolerance = 3
	assert.LessOrEqual(t, after, before+tolerance,
		"goroutine count grew by %d after failed NewWebsocket — heartbeat leak suspected", after-before)
}

