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
		code      int
		wantRetry bool
	}{
		{http.StatusInternalServerError, true},
		{http.StatusBadGateway, true},
		{http.StatusTooManyRequests, false},
	}

	for _, tc := range cases {
		got := c.shouldRetry(tc.code)
		assert.Equalf(t, tc.wantRetry, got,
			"shouldRetry(%d) = %v, want %v", tc.code, got, tc.wantRetry)
	}
}
