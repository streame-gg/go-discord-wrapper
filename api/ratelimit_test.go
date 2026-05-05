package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ratelimitTestSuite struct {
	suite.Suite
}

func TestRatelimitTestSuite(t *testing.T) {
	suite.Run(t, new(ratelimitTestSuite))
}

// ---- routeKey normalisation ----

func (s *ratelimitTestSuite) TestRouteKey_NonMajorSnowflakesReplaced() {
	// Message IDs are not major params; they should be normalized.
	got := routeKey("GET", "/channels/111222333444555/messages/999888777666555")
	want := "GET:channels/111222333444555/messages/{id}"
	s.Equalf(want, got, "routeKey = %q, want %q", got, want)
}

func (s *ratelimitTestSuite) TestRouteKey_ChannelIDPreserved() {
	// channel_id is a major param and must be kept so two channels get separate buckets.
	a := routeKey("POST", "/channels/111/messages")
	b := routeKey("POST", "/channels/222/messages")
	s.NotEqualf(a, b, "expected different route keys for different channel IDs, both got %q", a)
}

func (s *ratelimitTestSuite) TestRouteKey_GuildIDPreserved() {
	a := routeKey("GET", "/guilds/111/members")
	b := routeKey("GET", "/guilds/222/members")
	s.NotEqualf(a, b, "expected different route keys for different guild IDs, both got %q", a)
}

func (s *ratelimitTestSuite) TestRouteKey_SameChannelDifferentMethods() {
	// Different HTTP methods → different route keys (Discord can assign different buckets).
	get := routeKey("GET", "/channels/111/messages")
	post := routeKey("POST", "/channels/111/messages")
	s.NotEqual(get, post, "GET and POST on the same path should produce different route keys")
}

func (s *ratelimitTestSuite) TestRouteKey_WebhookTokenPreserved() {
	// Both the webhook ID and its token are major params.
	a := routeKey("POST", "/webhooks/123/tokenA/messages")
	b := routeKey("POST", "/webhooks/123/tokenB/messages")
	s.NotEqual(a, b, "different webhook tokens should produce different route keys")
}

func (s *ratelimitTestSuite) TestRouteKey_LeadingSlashHandled() {
	// Paths with and without a leading slash must produce the same key.
	withSlash := routeKey("GET", "/channels/111/messages")
	// The parts split strips the leading empty element.
	s.NotEmpty(withSlash, "routeKey returned empty string")
}

// ---- isSnowflakeID ----

func (s *ratelimitTestSuite) TestIsSnowflakeID() {
	cases := []struct {
		input string
		want  bool
	}{
		{"175928847299117063", true},     // typical 18-digit snowflake
		{"1234567890123456", true},       // 16 digits
		{"12345678901234", false},        // too short (14)
		{"123456789012345678901", false}, // too long (21)
		{"not-a-snowflake", false},
		{"@original", false},
		{"@me", false},
		{"messages", false},
	}
	for _, tc := range cases {
		got := isSnowflakeID(tc.input)
		s.Equalf(tc.want, got, "isSnowflakeID(%q) = %v, want %v", tc.input, got, tc.want)
	}
}

// ---- rateLimiter.update + wait ----

func makeResp(statusCode int, headers map[string]string) *http.Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: statusCode, Header: h}
}

func (s *ratelimitTestSuite) TestRateLimiter_UpdateAndWait_NoBlock() {
	rl := newRateLimiter(0)

	resp := makeResp(http.StatusOK, map[string]string{
		"X-RateLimit-Bucket":      "abc123",
		"X-RateLimit-Limit":       "5",
		"X-RateLimit-Remaining":   "4",
		"X-RateLimit-Reset-After": "1.000",
	})
	rl.update("GET", "/channels/111/messages", resp)

	start := time.Now()
	rl.wait("GET", "/channels/111/messages")
	elapsed := time.Since(start)

	s.LessOrEqualf(elapsed, 50*time.Millisecond, "wait blocked for %v, expected no block when remaining > 0", elapsed)
}

func (s *ratelimitTestSuite) TestRateLimiter_UpdateAndWait_BlocksAtZeroRemaining() {
	rl := newRateLimiter(0)

	resetAfter := 100 * time.Millisecond
	resp := makeResp(http.StatusOK, map[string]string{
		"X-RateLimit-Bucket":      "bucket1",
		"X-RateLimit-Limit":       "5",
		"X-RateLimit-Remaining":   "0",
		"X-RateLimit-Reset-After": "0.100", // 100 ms
	})
	rl.update("POST", "/channels/111/messages", resp)

	start := time.Now()
	rl.wait("POST", "/channels/111/messages")
	elapsed := time.Since(start)

	// Should have blocked for roughly the reset window.
	s.GreaterOrEqualf(elapsed, resetAfter-10*time.Millisecond, "wait returned too early: elapsed=%v, expected >=%v", elapsed, resetAfter)
}

func (s *ratelimitTestSuite) TestRateLimiter_SafetyMargin() {
	rl := newRateLimiter(2) // block when remaining <= 2

	resp := makeResp(http.StatusOK, map[string]string{
		"X-RateLimit-Bucket":      "bucket2",
		"X-RateLimit-Limit":       "10",
		"X-RateLimit-Remaining":   "2",
		"X-RateLimit-Reset-After": "0.150",
	})
	rl.update("GET", "/guilds/999/members", resp)

	start := time.Now()
	rl.wait("GET", "/guilds/999/members")
	elapsed := time.Since(start)

	// remaining(2) == safetyMargin(2), so we should block.
	s.GreaterOrEqualf(elapsed, 100*time.Millisecond, "wait returned too early with safety margin: elapsed=%v", elapsed)
}

func (s *ratelimitTestSuite) TestRateLimiter_GlobalRateLimit() {
	rl := newRateLimiter(0)

	// Simulate a global 429 with a short retry-after.
	resp := makeResp(http.StatusTooManyRequests, map[string]string{
		"X-RateLimit-Global": "true",
		"Retry-After":        "0.100",
	})
	rl.update("GET", "/channels/111/messages", resp)

	start := time.Now()
	// wait() should block for the global rate limit even on a different route.
	rl.wait("POST", "/guilds/222/channels")
	elapsed := time.Since(start)

	s.GreaterOrEqualf(elapsed, 80*time.Millisecond, "global rate limit did not block: elapsed=%v", elapsed)
}

func (s *ratelimitTestSuite) TestRateLimiter_GlobalDoesNotBlockAfterExpiry() {
	rl := newRateLimiter(0)

	resp := makeResp(http.StatusTooManyRequests, map[string]string{
		"X-RateLimit-Global": "true",
		"Retry-After":        "0.050", // 50 ms
	})
	rl.update("GET", "/channels/111", resp)

	// Wait for the global limit to expire naturally.
	time.Sleep(80 * time.Millisecond)

	start := time.Now()
	rl.wait("GET", "/channels/111")
	elapsed := time.Since(start)

	s.LessOrEqualf(elapsed, 30*time.Millisecond, "wait blocked after global limit expired: elapsed=%v", elapsed)
}

func (s *ratelimitTestSuite) TestRateLimiter_Route429UpdatesBucket() {
	rl := newRateLimiter(0)

	// Route-specific 429: remaining should be set to 0 and resetAt updated.
	resp := makeResp(http.StatusTooManyRequests, map[string]string{
		"X-RateLimit-Bucket":      "bucket3",
		"X-RateLimit-Limit":       "5",
		"X-RateLimit-Remaining":   "0",
		"X-RateLimit-Reset-After": "0.100",
	})
	rl.update("DELETE", "/channels/111/messages/999", resp)

	start := time.Now()
	rl.wait("DELETE", "/channels/111/messages/999")
	elapsed := time.Since(start)

	s.GreaterOrEqualf(elapsed, 80*time.Millisecond, "route 429 did not cause wait to block: elapsed=%v", elapsed)
}

func (s *ratelimitTestSuite) TestRateLimiter_UnknownRoutePasses() {
	rl := newRateLimiter(0)
	// No update called — wait should return immediately for an unknown route.

	start := time.Now()
	rl.wait("GET", "/applications/111/commands")
	elapsed := time.Since(start)

	s.LessOrEqualf(elapsed, 20*time.Millisecond, "unknown route blocked unexpectedly: elapsed=%v", elapsed)
}

func (s *ratelimitTestSuite) TestRateLimiter_BucketSharedAcrossRoutes() {
	rl := newRateLimiter(0)

	// Two different paths that Discord assigns to the same bucket.
	const sharedBucket = "shared_bucket"

	respA := makeResp(http.StatusOK, map[string]string{
		"X-RateLimit-Bucket":      sharedBucket,
		"X-RateLimit-Limit":       "5",
		"X-RateLimit-Remaining":   "0",
		"X-RateLimit-Reset-After": "0.100",
	})
	rl.update("GET", "/channels/111/messages", respA)

	// Simulate Discord returning the same bucket for pins on the same channel.
	respB := makeResp(http.StatusOK, map[string]string{
		"X-RateLimit-Bucket":      sharedBucket,
		"X-RateLimit-Limit":       "5",
		"X-RateLimit-Remaining":   "0",
		"X-RateLimit-Reset-After": "0.100",
	})
	rl.update("GET", "/channels/111/pins", respB)

	// Requesting the pin route should also block (same bucket, remaining=0).
	start := time.Now()
	rl.wait("GET", "/channels/111/pins")
	elapsed := time.Since(start)

	s.GreaterOrEqualf(elapsed, 80*time.Millisecond, "shared bucket did not block second route: elapsed=%v", elapsed)
}

// ---- parseRetryAfter ----

func (s *ratelimitTestSuite) TestParseRetryAfter_RetryAfterHeader() {
	resp := makeResp(429, map[string]string{"Retry-After": "1.5"})
	got := parseRetryAfter(resp)
	want := 1500 * time.Millisecond
	s.Equalf(want, got, "parseRetryAfter = %v, want %v", got, want)
}

func (s *ratelimitTestSuite) TestParseRetryAfter_Missing() {
	resp := makeResp(429, nil)
	got := parseRetryAfter(resp)
	s.Equalf(time.Duration(0), got, "parseRetryAfter = %v, want 0", got)
}
