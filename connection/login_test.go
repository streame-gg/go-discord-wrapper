package connection

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mockGatewayNoReady starts a WebSocket server that completes the handshake
// (sends OP 10, drains IDENTIFY) but never sends a READY event.
func mockGatewayNoReady(t *testing.T) (wsURL string, closeFn func()) {
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

		// Drain IDENTIFY (OP 2) but never send READY.
		var ignore interface{}
		_ = conn.ReadJSON(&ignore)

		// Hold open until the client disconnects.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	return "ws" + ts.URL[len("http"):], ts.Close
}

// Bug 1: When the context is cancelled before a READY event arrives, Login
// must return ctx.Err() instead of blocking forever.
func TestBug1_LoginReturnsOnContextCancel(t *testing.T) {
	wsURL, closeWS := mockGatewayNoReady(t)
	defer closeWS()

	// Stub /gateway/bot to return our no-READY WS URL so Login can proceed.
	gatewayStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"url":    wsURL,
			"shards": 1,
			"session_start_limit": map[string]interface{}{
				"total": 1000, "remaining": 999,
				"reset_after": 14400000, "max_concurrency": 1,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer gatewayStub.Close()

	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	require.NoError(t, err)

	// Point the client's internal HTTP client at our gateway stub so that
	// initializeGatewayConnection hits it instead of Discord's real API.
	c.httpClient = gatewayStub.Client()
	// Repoint the gateway URL: we monkey-patch the http transport to redirect.
	c.httpClient = &http.Client{
		Transport: &gatewayStubTransport{real: gatewayStub},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	loginErr := make(chan error, 1)
	go func() {
		loginErr <- c.Login(ctx)
	}()

	select {
	case err := <-loginErr:
		assert.Error(t, err, "Login must return an error when context is cancelled (Bug 1)")
		assert.ErrorIs(t, err, context.DeadlineExceeded,
			"Login must return DeadlineExceeded after ctx times out without READY (Bug 1)")
	case <-time.After(5 * time.Second):
		t.Fatal("Login blocked forever — Bug 1 select fix not working")
	}
}

// gatewayStubTransport redirects all requests to the stub server.
type gatewayStubTransport struct {
	real *httptest.Server
}

func (t *gatewayStubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Host = t.real.Listener.Addr().String()
	req2.URL.Scheme = "http"
	return http.DefaultTransport.RoundTrip(req2)
}

// Bug 2: SessionID and ReconnectURL must be written under wsMu.Lock so that
// concurrent reads (as in reconnect()) are race-free. This test verifies the
// race detector reports no data race.
func TestBug2_SessionIDWriteUnderLock(t *testing.T) {
	wsURL, closeServer := mockGateway(t)
	defer closeServer()

	c, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	// Continuously read SessionID/ReconnectURL under wsMu.RLock, mirroring
	// what reconnect() does (websocket.go).  With -race this detects any
	// unsynchronised write in internalEventHandler.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 500; i++ {
			c.wsMu.RLock()
			_ = c.Websocket.SessionID
			_ = c.Websocket.ReconnectURL
			c.wsMu.RUnlock()
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("race reader goroutine did not finish")
	}
}
