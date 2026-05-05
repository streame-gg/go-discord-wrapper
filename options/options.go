// Package options provides functional options for configuring the Discord client
// and REST client. Import this package to access all With* configuration functions.
package options

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// Option is a functional option for configuring a Client or RestClient.
// All With* functions in this package return an Option.
type Option func(*Config)

// Config holds configuration for both Client and RestClient.
// Fields that do not apply to a given type are ignored.
type Config struct {
	// Client-only
	Sharding    *Sharding
	Logger      *slog.Logger
	Coordinator ShardCoordinator
	Cache       cache.Cache

	// Shared between Client and its embedded RestClient
	APIVersion common.APIVersion

	// REST client options
	BaseURL            string
	HTTPClient         *http.Client
	Retry              RetryOptions
	MinRequestInterval time.Duration
	RateLimiter        *RateLimiterOptions
}

// RateLimiterOptions configures proactive rate limit throttling.
//
// The rate limiter tracks Discord's per-route and global rate limit headers and
// blocks outgoing requests before they would hit a limit, rather than waiting
// for a 429 to arrive. It is enabled by default.
type RateLimiterOptions struct {
	// Disabled turns off proactive rate limiting entirely.
	// When disabled, only reactive 429 handling (via RetryOptions) applies.
	Disabled bool

	// SafetyMargin is the number of remaining requests at which the client
	// proactively waits for the current window to expire before sending.
	// Increase this when many goroutines or shards share the same bucket and
	// you want extra headroom. Default: 0 (block only when remaining hits 0).
	SafetyMargin int
}

// Sharding holds bot sharding configuration.
type Sharding struct {
	TotalShards int
	ShardID     int
}

// RetryOptions configures retry behaviour for the REST client.
type RetryOptions struct {
	MaxRetries          int
	BaseBackoff         time.Duration
	MaxBackoff          time.Duration
	RetryOnRateLimit    bool
	RetryOnServerErrors bool
}

// Build applies opts on top of defaults and returns the resulting Config.
func Build(defaults Config, opts []Option) Config {
	for _, o := range opts {
		o(&defaults)
	}
	return defaults
}

// Ptr returns a pointer to v. Useful for passing literal values to fields that
// require a pointer (e.g. Required, MinValues, MaxValues on command options).
//
//	options.Ptr(true)   // *bool
//	options.Ptr(42)     // *int
func Ptr[T any](v T) *T { return &v }

// WithSharding configures bot sharding.
//
//	options.WithSharding(1, 0)  // 1 total shard, shard ID 0
func WithSharding(totalShards, shardID int) Option {
	return func(c *Config) {
		c.Sharding = &Sharding{TotalShards: totalShards, ShardID: shardID}
	}
}

// WithAPIVersion sets the Discord API version (default: APIVersion10).
func WithAPIVersion(v common.APIVersion) Option {
	return func(c *Config) { c.APIVersion = v }
}

// WithLogger sets a custom slog logger on the Client.
func WithLogger(l *slog.Logger) Option {
	return func(c *Config) { c.Logger = l }
}

// WithHTTPClient sets a custom *http.Client used for all REST requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) { c.HTTPClient = client }
}

// WithRetry configures retry behaviour for the REST client.
//
//	options.WithRetry(options.RetryOptions{
//	    MaxRetries:       5,
//	    RetryOnRateLimit: true,
//	})
func WithRetry(r RetryOptions) Option {
	return func(c *Config) { c.Retry = r }
}

// WithMinRequestInterval sets the minimum time between consecutive REST requests,
// acting as a global rate-limiter across all calls.
func WithMinRequestInterval(d time.Duration) Option {
	return func(c *Config) { c.MinRequestInterval = d }
}

// WithBaseURL overrides the Discord REST API base URL.
// Defaults to "https://discord.com/api". Useful for testing or API proxies.
func WithBaseURL(url string) Option {
	return func(c *Config) { c.BaseURL = url }
}

// WithCoordinator attaches a ShardCoordinator to the client.
// The coordinator enables cross-shard messaging (SendToShard, BroadcastToShards,
// and the RequestAll helper in the sharding package).
//
// For single-machine bots use sharding.NewLocalCoordinator.
// For distributed bots supply your own ShardCoordinator implementation backed
// by a message broker (Redis, NATS, gRPC, etc.).
func WithCoordinator(c ShardCoordinator) Option {
	return func(cfg *Config) { cfg.Coordinator = c }
}

// WithRateLimiting configures proactive rate limit throttling.
// Proactive throttling is enabled by default; use this option to tune it or
// to disable it if you have your own rate-limit handling.
//
//	// Disable proactive rate limiting:
//	options.WithRateLimiting(options.RateLimiterOptions{Disabled: true})
//
//	// Leave 2 requests as headroom (useful with many concurrent goroutines):
//	options.WithRateLimiting(options.RateLimiterOptions{SafetyMargin: 2})
func WithRateLimiting(r RateLimiterOptions) Option {
	return func(c *Config) { c.RateLimiter = &r }
}

// WithCache attaches a Cache to the client.
// The gateway client automatically populates all stores from gateway events
// (GUILD_CREATE/DELETE, CHANNEL_CREATE/DELETE, MESSAGE_CREATE/UPDATE/DELETE).
//
// Use cache.NewMemoryCache for an in-process cache.
// For Redis, MongoDB, or other backends, implement the cache.Cache interface.
//
//	c := cache.NewMemoryCache(cache.Options{
//	    TTL:           10 * time.Minute,
//	    EvictBehavior: cache.EvictUnused,
//	    UnusedWindow:  5 * time.Minute,
//	})
//	connection.NewClient(token, intents, options.WithCache(c))
func WithCache(ca cache.Cache) Option {
	return func(c *Config) { c.Cache = ca }
}
