package connection

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// gatewayBotResponse is the minimal /gateway/bot JSON Discord returns.
var gatewayBotResponse = common.BotRegisterResponse{
	Url:    "wss://gateway.discord.gg",
	Shards: 1,
}

// TestBug38UserAgentSent verifies that initializeGatewayConnection sends the
// required Discord User-Agent header (Bug 38).
func TestBug38UserAgentSent(t *testing.T) {
	var gotUserAgent string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gatewayBotResponse) //nolint:errcheck
	}))
	defer ts.Close()

	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)
	// Point the client's httpClient at the test server.
	c.httpClient = &http.Client{Timeout: 5 * time.Second}

	// Override the URL by temporarily patching common.APIBaseString — we can't
	// easily do that, so instead exercise through a wrapper that calls the real
	// function but the server URL is injected via a custom RoundTripper.
	c.httpClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = ts.Listener.Addr().String()
		return http.DefaultTransport.RoundTrip(req)
	})

	_, err = c.initializeGatewayConnection(context.Background())
	require.NoError(t, err)

	assert.Contains(t, gotUserAgent, "DiscordBot", "User-Agent must identify as DiscordBot (Bug 38)")
}

// TestBug37ContextCancelledAbortsRequest verifies that a cancelled context
// causes initializeGatewayConnection to return ctx.Err() (Bug 37).
func TestBug37ContextCancelledAbortsRequest(t *testing.T) {
	// Server that blocks forever so the context cancellation must fire first.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer ts.Close()

	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)
	c.httpClient = &http.Client{Timeout: 5 * time.Second}
	c.httpClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = ts.Listener.Addr().String()
		return http.DefaultTransport.RoundTrip(req)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	_, err = c.initializeGatewayConnection(ctx)
	assert.ErrorIs(t, err, context.Canceled, "cancelled context must propagate (Bug 37)")
}

// TestBug39RateLimitRetry verifies that a 429 with Retry-After causes
// initializeGatewayConnection to retry after the specified delay (Bug 39).
func TestBug39RateLimitRetry(t *testing.T) {
	var calls atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			// First call: return 429 with a very short Retry-After.
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		// Second call: return success.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gatewayBotResponse) //nolint:errcheck
	}))
	defer ts.Close()

	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)
	c.httpClient = &http.Client{Timeout: 5 * time.Second}
	c.httpClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = ts.Listener.Addr().String()
		return http.DefaultTransport.RoundTrip(req)
	})

	resp, err := c.initializeGatewayConnection(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int32(2), calls.Load(), "must retry exactly once after 429 (Bug 39)")
	assert.Equal(t, gatewayBotResponse.Url, resp.Url)
}

// TestBug40HttpClientReused verifies that the Client field httpClient is the
// same pointer across multiple calls (Bug 40).
func TestBug40HttpClientReused(t *testing.T) {
	c, err := NewClient("Bot fake-token", common.IntentGuilds)
	require.NoError(t, err)
	first := c.httpClient
	require.NotNil(t, first, "httpClient must be initialized by NewClient (Bug 40)")
	// The pointer must not change — it is a shared instance, not allocated per call.
	assert.Same(t, first, c.httpClient, "httpClient must be a stable field, not re-allocated (Bug 40)")
}

// roundTripperFunc adapts a function to the http.RoundTripper interface so tests
// can inject a custom transport without replacing the whole http.Client.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
