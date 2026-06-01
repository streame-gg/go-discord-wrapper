package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// rateLimitSuite covers the internal, dependency-free rate-limiter helpers:
// route-key normalisation (which drives bucket sharing), major-parameter
// detection, snowflake recognition, Retry-After parsing, and the bucket
// lifecycle (create → wait → purge → close).
type rateLimitSuite struct {
	suite.Suite
}

func TestRateLimitSuite(t *testing.T) {
	suite.Run(t, new(rateLimitSuite))
}

// hdr builds an http.Header from key/value pairs using Set so the keys are
// canonicalised exactly as a real *http.Response would store them (the
// rate-limiter reads them via Header.Get, which canonicalises the lookup key).
func hdr(kv ...string) http.Header {
	h := http.Header{}
	for i := 0; i+1 < len(kv); i += 2 {
		h.Set(kv[i], kv[i+1])
	}
	return h
}

func (s *rateLimitSuite) TestRouteKeyNormalisation() {
	cases := []struct {
		method, path, want string
	}{
		// Non-major snowflakes collapse to {id}; major params are preserved.
		// (routeKey trims the leading slash before joining.)
		{"GET", "/channels/111111111111111111/messages/999999999999999999", "GET:channels/111111111111111111/messages/{id}"},
		{"POST", "/channels/111111111111111111/messages", "POST:channels/111111111111111111/messages"},
		{"GET", "/guilds/222222222222222222/members/333333333333333333", "GET:guilds/222222222222222222/members/{id}"},
		// Webhook id and token are both major params and survive normalisation.
		{"GET", "/webhooks/444444444444444444/tok/messages/555555555555555555", "GET:webhooks/444444444444444444/tok/messages/{id}"},
		// Query string is stripped before keying.
		{"GET", "/channels/111111111111111111/messages?limit=50", "GET:channels/111111111111111111/messages"},
	}
	for _, c := range cases {
		s.Equalf(c.want, routeKey(c.method, c.path), "routeKey(%s %s)", c.method, c.path)
	}
}

func (s *rateLimitSuite) TestIsMajorParamOwner() {
	s.True(isMajorParamOwner("channels"))
	s.True(isMajorParamOwner("guilds"))
	s.True(isMajorParamOwner("webhooks"))
	s.False(isMajorParamOwner("messages"))
	s.False(isMajorParamOwner(""))
}

func (s *rateLimitSuite) TestIsSnowflakeID() {
	s.True(isSnowflakeID("123456789012345678"))     // 18 digits
	s.True(isSnowflakeID("123456789012345"))        // 15 digits (lower bound)
	s.False(isSnowflakeID("12345678901234"))        // 14 digits
	s.False(isSnowflakeID("123456789012345678901")) // 21 digits
	s.False(isSnowflakeID("12345678901234x678"))    // non-digit
	s.False(isSnowflakeID(""))
}

func (s *rateLimitSuite) TestParseRetryAfter() {
	s.Run("seconds float", func() {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"1.5"}}}
		s.Equal(1500*time.Millisecond, parseRetryAfter(resp))
	})
	s.Run("reset-after fallback", func() {
		resp := &http.Response{Header: hdr("X-RateLimit-Reset-After", "2")}
		s.Equal(2*time.Second, parseRetryAfter(resp))
	})
	s.Run("http date in the future", func() {
		future := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)
		resp := &http.Response{Header: http.Header{"Retry-After": []string{future}}}
		d := parseRetryAfter(resp)
		s.Greater(d, time.Duration(0))
		s.LessOrEqual(d, 30*time.Second)
	})
	s.Run("absent headers", func() {
		s.Equal(time.Duration(0), parseRetryAfter(&http.Response{Header: http.Header{}}))
	})
}

func (s *rateLimitSuite) TestNewRateLimiterAndClose() {
	rl := newRateLimiter(0)
	s.Require().NotNil(rl)
	// Unknown route: wait must return immediately (nil bucket → pass-through).
	s.NoError(rl.wait(context.Background(), http.MethodGet, "/channels/111111111111111111/messages"))
	// close is idempotent.
	rl.close()
	rl.close()
}

func (s *rateLimitSuite) TestUpdateLearnsBucketAndWaitRespectsReset() {
	rl := newRateLimiter(0)
	defer rl.close()

	method, path := http.MethodGet, "/channels/111111111111111111/pins"

	// A response that exhausts the bucket for ~40ms.
	resp := &http.Response{Header: hdr(
		"X-RateLimit-Bucket", "bkt1",
		"X-RateLimit-Limit", "1",
		"X-RateLimit-Remaining", "0",
		"X-RateLimit-Reset-After", "0.04",
	)}
	rl.update(method, path, resp)

	// wait should block until the window resets, then return without error.
	start := time.Now()
	s.Require().NoError(rl.wait(context.Background(), method, path))
	s.GreaterOrEqual(time.Since(start), 20*time.Millisecond, "wait should have slept until reset")
}

func (s *rateLimitSuite) TestWaitCancelledContext() {
	rl := newRateLimiter(0)
	defer rl.close()

	method, path := http.MethodGet, "/channels/111111111111111111/typing"
	rl.update(method, path, &http.Response{Header: hdr(
		"X-RateLimit-Bucket", "bkt2",
		"X-RateLimit-Limit", "1",
		"X-RateLimit-Remaining", "0",
		"X-RateLimit-Reset-After", "5", // long window
	)})

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := rl.wait(ctx, method, path)
	s.ErrorIs(err, context.DeadlineExceeded)
}

func (s *rateLimitSuite) TestPurgeStaleBuckets() {
	rl := newRateLimiter(0)
	defer rl.close()

	method, path := http.MethodGet, "/guilds/222222222222222222/bans"
	rl.update(method, path, &http.Response{Header: hdr(
		"X-RateLimit-Bucket", "bkt3",
		"X-RateLimit-Limit", "5",
		"X-RateLimit-Remaining", "5",
	)})

	// Bucket exists now.
	_, ok := rl.buckets.Load("bkt3")
	s.Require().True(ok)

	// Purging with a negative maxAge makes the cutoff in the future, so the
	// just-touched bucket counts as stale and is evicted along with its route map.
	rl.purgeStaleBuckets(-time.Hour)
	_, ok = rl.buckets.Load("bkt3")
	s.False(ok, "stale bucket should be purged")
	_, ok = rl.routeToBucket.Load(routeKey(method, path))
	s.False(ok, "route→bucket entry should be purged with its bucket")
}

// TestGlobalRateLimitParksRequests verifies update() records a global 429 lock
// and that wait() then blocks for the global window.
func (s *rateLimitSuite) TestGlobalRateLimit() {
	rl := newRateLimiter(0)
	defer rl.close()

	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header: hdr(
			"X-RateLimit-Global", "true",
			"Retry-After", "0.04",
		),
	}
	rl.update(http.MethodGet, "/users/@me", resp)

	start := time.Now()
	s.Require().NoError(rl.wait(context.Background(), http.MethodGet, "/users/@me"))
	s.GreaterOrEqual(time.Since(start), 20*time.Millisecond, "global lock should delay wait")
}
