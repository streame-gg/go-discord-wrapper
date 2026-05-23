package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/options"
)

// newRetryClient returns a RestClient with server-error retry enabled.
func newRetryClient(t *testing.T) *RestClient {
	t.Helper()
	rc, err := NewRestClient("test-token",
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{
			MaxRetries:          3,
			RetryOnServerErrors: true,
			RetryOnRateLimit:    false,
		}),
	)
	require.NoError(t, err)
	return rc
}

// TestShouldRetryMethodFilter verifies that POST and PATCH are never retried
// on 5xx responses, while idempotent methods (GET, PUT, DELETE) are (Issue 10).
func TestShouldRetryMethodFilter(t *testing.T) {
	c := newRetryClient(t)

	cases := []struct {
		method    string
		code      int
		wantRetry bool
	}{
		{http.MethodGet, http.StatusInternalServerError, true},
		{http.MethodGet, http.StatusBadGateway, true},
		{http.MethodPut, http.StatusInternalServerError, true},
		{http.MethodDelete, http.StatusInternalServerError, true},
		// Non-idempotent: must never retry on 5xx
		{http.MethodPost, http.StatusInternalServerError, false},
		{http.MethodPost, http.StatusBadGateway, false},
		{http.MethodPatch, http.StatusInternalServerError, false},
		{http.MethodPatch, http.StatusBadGateway, false},
		// Rate limit is controlled by RetryOnRateLimit, not method
		{http.MethodPost, http.StatusTooManyRequests, false},
		{http.MethodGet, http.StatusTooManyRequests, false},
	}

	for _, tc := range cases {
		got := c.shouldRetry(tc.code, tc.method)
		assert.Equalf(t, tc.wantRetry, got,
			"shouldRetry(%d, %q) = %v, want %v", tc.code, tc.method, got, tc.wantRetry)
	}
}
