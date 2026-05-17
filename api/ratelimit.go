package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// rateLimitBucket holds proactive rate limit state for a single Discord bucket.
type rateLimitBucket struct {
	mu               sync.Mutex
	limit            int       // X-RateLimit-Limit
	remaining        int       // X-RateLimit-Remaining (locally tracked)
	resetAt          time.Time // when this bucket's window resets
	lastUsed         time.Time // updated on every wait/update; drives eviction
	windowRestoredAt time.Time // when remaining was last proactively restored
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
	closeOnce    sync.Once

	// Global rate limit. All requests park here when the global limit fires.
	globalMu    sync.RWMutex
	globalUntil time.Time

	// routeToBucket maps a normalised route key (e.g. "GET:/channels/{id}/messages:111")
	// to the Discord-assigned bucket hash string.
	routeToBucket sync.Map // map[string]string

	// buckets maps a bucket hash string to its live state.
	buckets sync.Map // map[string]*rateLimitBucket

	// done is closed by close() to stop the background cleanup goroutine.
	done chan struct{}
}

func newRateLimiter(safetyMargin int) *rateLimiter {
	rl := &rateLimiter{
		safetyMargin: safetyMargin,
		done:         make(chan struct{}),
	}
	go rl.cleanupLoop()
	return rl
}

// close stops the background cleanup goroutine. Call it when the owning
// RestClient is no longer needed to prevent a goroutine leak.
// Safe to call more than once.
func (r *rateLimiter) close() {
	r.closeOnce.Do(func() { close(r.done) })
}

// cleanupLoop runs purgeStaleBuckets on a fixed cadence until close() is called.
func (r *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.purgeStaleBuckets(10 * time.Minute)
		case <-r.done:
			return
		}
	}
}

// purgeStaleBuckets deletes any bucket that has not been touched within maxAge.
// It is called automatically by the cleanup goroutine; it is also exported for
// testing.
func (r *rateLimiter) purgeStaleBuckets(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)
	r.buckets.Range(func(key, val any) bool {
		b := val.(*rateLimitBucket)
		b.mu.Lock()
		stale := !b.lastUsed.IsZero() && b.lastUsed.Before(cutoff)
		b.mu.Unlock()
		if stale {
			r.buckets.Delete(key)
			// Also remove all routeToBucket entries pointing at this bucket so
			// routeToBucket does not grow unboundedly (Bug 22).
			bucketID := key.(string)
			r.routeToBucket.Range(func(rk, rv any) bool {
				if rv.(string) == bucketID {
					r.routeToBucket.Delete(rk)
				}
				return true
			})
		}
		return true
	})
}

// wait blocks until both the global rate limit and the per-route bucket allow
// the request to proceed, then proactively decrements the remaining counter.
// It returns ctx.Err() immediately if the context is cancelled while waiting.
func (r *rateLimiter) wait(ctx context.Context, method, path string) error {
	// 1. Global rate limit — loop in case it is extended while we sleep.
	for {
		r.globalMu.RLock()
		until := r.globalUntil
		r.globalMu.RUnlock()

		wait := time.Until(until)
		if wait <= 0 {
			break
		}
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// 2. Per-route bucket.
	key := routeKey(method, path)

	bucketIDVal, ok := r.routeToBucket.Load(key)
	if !ok {
		return nil // No bucket known yet; let the request through.
	}

	bucketVal, ok := r.buckets.Load(bucketIDVal.(string))
	if !ok {
		return nil
	}

	b := bucketVal.(*rateLimitBucket)
	b.mu.Lock()
	b.lastUsed = time.Now()

	// Loop so that if the bucket is still exhausted after sleeping (e.g. the
	// window didn't actually reset yet) we sleep again rather than over-sending.
	for b.remaining <= r.safetyMargin && !b.resetAt.IsZero() {
		wait := time.Until(b.resetAt)
		if wait <= 0 {
			// Window has elapsed; restore remaining exactly once per window.
			// windowRestoredAt tracks when we last restored so that multiple
			// goroutines that all slept for the same window do not each
			// independently set remaining = limit (TOCTOU race).
			if b.windowRestoredAt.Before(b.resetAt) {
				if b.limit > 0 {
					b.remaining = b.limit
				} else {
					// Unknown limit (e.g. 429 without X-RateLimit-Limit header).
					// Let one request through to re-learn the bucket.
					b.remaining = 1
				}
				b.windowRestoredAt = time.Now()
			}
			break
		}
		b.mu.Unlock()
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
		b.mu.Lock()
		if !time.Now().Before(b.resetAt) && b.windowRestoredAt.Before(b.resetAt) {
			if b.limit > 0 {
				b.remaining = b.limit
			} else {
				// Unknown limit (e.g. 429 without X-RateLimit-Limit header).
				// Let one request through to re-learn the bucket.
				b.remaining = 1
			}
			b.windowRestoredAt = time.Now()
		}
	}

	if b.remaining > 0 {
		b.remaining--
	}
	b.mu.Unlock()
	return nil
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
	b.lastUsed = time.Now()
	if !resetAt.IsZero() {
		// Two responses within the same rate-limit window will have X-RateLimit-Reset-After
		// values that differ by at most the round-trip latency. Using 1 s as the tolerance
		// treats responses with nearly identical reset times as the same window.
		const windowTolerance = time.Second
		if resetAt.After(b.resetAt.Add(windowTolerance)) {
			// Clearly a newer rate-limit window — accept all values from the response.
			b.remaining = remaining
			b.resetAt = resetAt
		} else if !resetAt.Before(b.resetAt.Add(-windowTolerance)) {
			// Same window (within tolerance) — only take the more conservative remaining.
			if remaining < b.remaining {
				b.remaining = remaining
			}
		}
		// Older window: ignore remaining/resetAt entirely.
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
	if i := strings.IndexAny(path, "?#"); i >= 0 {
		path = path[:i]
	}
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
