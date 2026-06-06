# Configuration

Both `connection.NewClient` and `api.NewRestClient` are configured with
functional options from the [`options`](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper/options)
package. Pass any number of `options.With…` values after the required
arguments; each is validated, and `NewClient` / `NewRestClient` return an error
for an inconsistent configuration (bad shard range, negative timeouts, …).

```go
bot, err := connection.NewClient(token, intents,
    options.WithCache(mc),
    options.WithRetry(options.RetryOptions{MaxRetries: 5}),
    options.WithLogLevel(slog.LevelDebug),
)
```

## Intents

The second argument to `NewClient` is an intent bitmask. Discord only sends
events for the intents you request, so combine exactly the ones you need:

```go
discord.IntentGuilds | discord.IntentGuildMessages | discord.IntentMessageContent
```

`GUILD_MEMBERS`, `GUILD_PRESENCES`, and `MESSAGE_CONTENT` are **privileged** —
enable them in the Developer Portal before requesting them, or the gateway will
reject the connection.

## Option reference

### Client & gateway

| Option | Effect | Default |
|--------|--------|---------|
| `WithCache(c)` | Attach a cache backend; the gateway auto-populates it | none |
| `WithCacheStores(mask)` | Limit auto-population to specific `cache.Category*` stores | all stores |
| `WithDisableCacheStore(mask)` | Auto-populate everything except these stores | — |
| `WithDisableCacheAutoPopulation()` | Never write to the cache from events (you manage it) | off |
| `WithSharding(total, id)` | Run as shard `id` of `total` | 1 shard |
| `WithCoordinator(c)` | Join a shard coordinator for cross-shard messaging | none |
| `WithMaxReconnectRetries(n)` | Gateway reconnect attempts before giving up (`-1` = infinite) | 3 |
| `WithMaxConcurrentEvents(n)` | Cap on events processed concurrently; excess is dropped | 64 |
| `WithLogger(l)` | Use a custom `*slog.Logger` | default text logger |
| `WithLogLevel(level)` | Minimum level for the default logger (ignored if `WithLogger` set) | Info |

### REST client

| Option | Effect | Default |
|--------|--------|---------|
| `WithAPIVersion(v)` | Discord API version | v10 |
| `WithBaseURL(url)` | Override the API base URL (handy for `httptest` mock servers) | Discord |
| `WithHTTPClient(c)` | Provide your own `*http.Client` | shared default |
| `WithRetry(opts)` | Retry policy (see below) | sensible defaults |
| `WithRateLimiting(opts)` | Proactive rate-limiter tuning (see below) | enabled |
| `WithMinRequestInterval(d)` | Global minimum delay between REST requests | 0 |
| `WithMaxResponseBodySize(n)` | Cap bytes read from a response (`-1` = unlimited) | 50 MiB |

## Retries

`RetryOptions` controls how the REST client reacts to transient failures:

```go
options.WithRetry(options.RetryOptions{
    MaxRetries:          5,
    BaseBackoff:         200 * time.Millisecond,
    MaxBackoff:          5 * time.Second,
    RetryOnRateLimit:    true, // retry after a 429 (reactive)
    RetryOnServerErrors: true, // retry 5xx
})
```

Backoff grows exponentially from `BaseBackoff`, capped at `MaxBackoff`.

## Rate limiting

The client ships a **proactive** rate limiter (enabled by default). It reads
Discord's per-route and global rate-limit headers and waits *before* a request
would breach a bucket, so you rarely see a 429 at all.

```go
options.WithRateLimiting(options.RateLimiterOptions{
    // Block once this many requests remain in a bucket, leaving headroom for
    // other goroutines/shards sharing it. 0 = block only when it hits zero.
    SafetyMargin: 2,
})

// Turn it off and rely only on reactive 429 handling via WithRetry:
options.WithRateLimiting(options.RateLimiterOptions{Disabled: true})
```

`WithMinRequestInterval` adds a flat floor between *all* requests, independent
of buckets — useful as a blunt global throttle.

## A note on `options.Ptr`

Many Discord fields are optional and therefore pointers (`*bool`, `*int`).
`options.Ptr` wraps a literal so you can set them inline:

```go
Required: options.Ptr(true),
MaxValue: options.Ptr[int64](100),
```

## Testing without Discord

Combine `WithBaseURL` with `net/http/httptest` to run integration tests against
a mock server:

```go
ts := httptest.NewServer(myHandler)
defer ts.Close()

client, _ := api.NewRestClient("test-token", options.WithBaseURL(ts.URL))
```

## See also

- [docs/CACHE.md](CACHE.md) — cache options in depth (`cache.Options`).
- [docs/SHARDING.md](SHARDING.md) — `WithSharding` and `WithCoordinator`.
- [docs/EVENTS.md](EVENTS.md) — `WithMaxConcurrentEvents` and dispatch.
