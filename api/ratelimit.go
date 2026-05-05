package api

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// rateLimitBucket holds proactive rate limit state for a single Discord bucket.
type rateLimitBucket struct {
	mu        sync.Mutex
	limit     int       // X-RateLimit-Limit
	remaining int       // X-RateLimit-Remaining (locally tracked)
	resetAt   time.Time // when this bucket's window resets
}

// rateLimiter implements proactive, header-driven rate limiting for the Discord API.
//
// It tracks two levels:
//   - Global rate limit (50 req/s shared across all routes)
//   - Per-route bucket limits (keyed by X-RateLimit-Bucket value)
//
// Route → bucket mapping is learned from X-RateLimit-Bucket response headers and
// persisted for future requests on the same route.
type rateLimiter struct {
	safetyMargin int // block when remaining <= safetyMargin

	// Global rate limit. All requests park here when the global limit fires.
	globalMu    sync.RWMutex
	globalUntil time.Time

	// routeToBucket maps a normalised route key (e.g. "GET:/channels/{id}/messages:111")
	// to the Discord-assigned bucket hash string.
	routeToBucket sync.Map // map[string]string

	// buckets maps a bucket hash string to its live state.
	buckets sync.Map // map[string]*rateLimitBucket
}

func newRateLimiter(safetyMargin int) *rateLimiter {
	return &rateLimiter{safetyMargin: safetyMargin}
}

// wait blocks until both the global rate limit and the per-route bucket allow
// the request to proceed, then proactively decrements the remaining counter.
func (r *rateLimiter) wait(method, path string) {
	// 1. Global rate limit — loop in case it is extended while we sleep.
	for {
		r.globalMu.RLock()
		until := r.globalUntil
		r.globalMu.RUnlock()

		wait := time.Until(until)
		if wait <= 0 {
			break
		}
		time.Sleep(wait)
	}

	// 2. Per-route bucket.
	key := routeKey(method, path)

	bucketIDVal, ok := r.routeToBucket.Load(key)
	if !ok {
		return // No bucket known yet; let the request through.
	}

	bucketVal, ok := r.buckets.Load(bucketIDVal.(string))
	if !ok {
		return
	}

	b := bucketVal.(*rateLimitBucket)
	b.mu.Lock()

	if b.remaining <= r.safetyMargin && !b.resetAt.IsZero() {
		wait := time.Until(b.resetAt)
		b.mu.Unlock()

		if wait > 0 {
			time.Sleep(wait)
		}

		// Re-acquire and optimistically restore remaining so other goroutines
		// do not all sleep again after this window has passed.
		b.mu.Lock()
		if b.limit > 0 && !time.Now().Before(b.resetAt) {
			b.remaining = b.limit
		}
	}

	if b.remaining > 0 {
		b.remaining--
	}
	b.mu.Unlock()
}

// update reads Discord rate-limit headers from resp and updates internal state.
// It must be called after every HTTP response, including 429s.
func (r *rateLimiter) update(method, path string, resp *http.Response) {
	if resp == nil {
		return
	}

	isGlobal := resp.Header.Get("X-RateLimit-Global") == "true"

	// Handle global rate limit lock first.
	if resp.StatusCode == http.StatusTooManyRequests && isGlobal {
		wait := parseRetryAfter(resp)
		if wait > 0 {
			r.globalMu.Lock()
			proposed := time.Now().Add(wait)
			if proposed.After(r.globalUntil) {
				r.globalUntil = proposed
			}
			r.globalMu.Unlock()
		}
		return
	}

	bucketID := resp.Header.Get("X-RateLimit-Bucket")
	if bucketID == "" {
		return // Route is not rate-limited by Discord (e.g. some gateway endpoints).
	}

	// Persist route → bucket mapping for future proactive checks.
	key := routeKey(method, path)
	r.routeToBucket.Store(key, bucketID)

	// Parse headers.
	limit, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	resetAfterSecs, _ := strconv.ParseFloat(
		strings.TrimSpace(resp.Header.Get("X-RateLimit-Reset-After")), 64,
	)

	var resetAt time.Time
	if resetAfterSecs > 0 {
		resetAt = time.Now().Add(time.Duration(resetAfterSecs * float64(time.Second)))
	}

	// On a route-specific 429, also lock this bucket.
	if resp.StatusCode == http.StatusTooManyRequests {
		wait := parseRetryAfter(resp)
		if wait > 0 && resetAt.IsZero() {
			resetAt = time.Now().Add(wait)
		}
		remaining = 0
	}

	// Update (or create) the bucket entry.
	newBucket := &rateLimitBucket{}
	actual, _ := r.buckets.LoadOrStore(bucketID, newBucket)
	b := actual.(*rateLimitBucket)

	b.mu.Lock()
	b.limit = limit
	b.remaining = remaining
	if !resetAt.IsZero() {
		b.resetAt = resetAt
	}
	b.mu.Unlock()
}

// routeKey returns a stable string key for a Discord API route.
//
// Snowflake IDs that are not major parameters are replaced with "{id}" so that
// requests to the same logical route (e.g. different messages in the same channel)
// share one bucket entry. Major parameters (channel_id, guild_id, webhook_id) are
// preserved because Discord scopes its buckets by those values.
//
// Examples:
//
//	GET  /channels/111/messages/999  → "GET:/channels/111/messages/{id}"
//	POST /channels/111/messages      → "POST:/channels/111/messages"
//	GET  /guilds/222/members/333     → "GET:/guilds/222/members/{id}"
//	GET  /webhooks/444/tok/messages  → "GET:/webhooks/444/tok/messages"
func routeKey(method, path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	for i, part := range parts {
		// The segment immediately after a major-parameter owner is a major param;
		// keep it as-is so each major-param value gets its own bucket.
		if i > 0 && isMajorParamOwner(parts[i-1]) {
			continue
		}
		// Webhook tokens (position 2 after "webhooks") are also major params.
		if i > 1 && parts[i-2] == "webhooks" {
			continue
		}
		if isSnowflakeID(part) {
			parts[i] = "{id}"
		}
	}

	return method + ":" + strings.Join(parts, "/")
}

// isMajorParamOwner reports whether the path segment prev immediately precedes
// a major parameter value according to Discord's rate-limit documentation.
func isMajorParamOwner(prev string) bool {
	return prev == "channels" || prev == "guilds" || prev == "webhooks"
}

// isSnowflakeID reports whether s looks like a Discord Snowflake (15–20 decimal digits).
func isSnowflakeID(s string) bool {
	if len(s) < 15 || len(s) > 20 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// parseRetryAfter reads the Retry-After or X-RateLimit-Reset-After header and
// returns it as a duration. Returns 0 if neither header is present or parseable.
func parseRetryAfter(resp *http.Response) time.Duration {
	for _, h := range []string{"Retry-After", "X-RateLimit-Reset-After"} {
		val := strings.TrimSpace(resp.Header.Get(h))
		if val == "" {
			continue
		}
		secs, err := strconv.ParseFloat(val, 64)
		if err != nil || secs <= 0 {
			continue
		}
		return time.Duration(secs * float64(time.Second))
	}
	return 0
}
