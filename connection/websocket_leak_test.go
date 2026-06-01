package connection

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mockWSCloseAfterUpgrade starts a WebSocket server that upgrades the
// connection and then immediately closes it, triggering the first error path
// in newWSConn (ReadMessage fails).
func mockWSCloseAfterUpgrade(t *testing.T) (wsURL string, closeFn func()) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.Close() // close immediately, no Hello
	}))
	return "ws" + ts.URL[len("http"):], ts.Close
}

// TestP0_3_NewWebsocketClosesConnOnError verifies that newWSConn closes
// the underlying TCP connection when it encounters an error after a
// successful dial (e.g. server closes without sending OP 10).
//
// Without the fix, the gorilla/websocket connection leaks because the four
// early-return paths in newWSConn returned without calling c.Close().
func TestP0_3_NewWebsocketClosesConnOnError(t *testing.T) {
	wsURL, closeFn := mockWSCloseAfterUpgrade(t)
	defer closeFn()

	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	require.NoError(t, err)

	before := runtime.NumGoroutine()

	_, wsErr := newWSConn(c, wsURL, false, nil, nil)
	assert.Error(t, wsErr, "newWSConn must return an error when server closes without Hello")

	// Allow any cleanup goroutines spawned by the gorilla/websocket library
	// to finish — they should exit quickly once the connection is closed.
	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	// Goroutine count must not grow; allow a small margin for scheduler jitter.
	assert.LessOrEqual(t, after, before+2,
		"goroutine count must not grow after newWSConn error: before=%d after=%d", before, after)
}
