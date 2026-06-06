package options

import (
	"github.com/stretchr/testify/suite"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// build is a small helper that applies opts onto a zero Config.
func build(opts ...Option) Config {
	return Build(Config{}, opts)
}

func (su *withOptionsSuite) TestPtr() {
	t := su.T()
	b := Ptr(true)
	if b == nil || *b != true {
		t.Fatalf("Ptr(true) = %v, want pointer to true", b)
	}
	n := Ptr(42)
	if n == nil || *n != 42 {
		t.Fatalf("Ptr(42) = %v, want pointer to 42", n)
	}
}

func (su *withOptionsSuite) TestWithSharding() {
	t := su.T()
	c := build(WithSharding(4, 2))
	if c.Sharding == nil {
		t.Fatal("Sharding not set")
	}
	if c.Sharding.TotalShards != 4 || c.Sharding.ShardID != 2 {
		t.Fatalf("got %+v", *c.Sharding)
	}
}

func (su *withOptionsSuite) TestWithAPIVersion() {
	t := su.T()
	c := build(WithAPIVersion(discord.APIVersion(9)))
	if c.APIVersion != discord.APIVersion(9) {
		t.Fatalf("got %v", c.APIVersion)
	}
}

func (su *withOptionsSuite) TestWithLogger() {
	t := su.T()
	l := slog.Default()
	c := build(WithLogger(l))
	if c.Logger != l {
		t.Fatal("Logger not set")
	}
}

func (su *withOptionsSuite) TestWithHTTPClient() {
	t := su.T()
	hc := &http.Client{Timeout: time.Second}
	c := build(WithHTTPClient(hc))
	if c.HTTPClient != hc {
		t.Fatal("HTTPClient not set")
	}
}

func (su *withOptionsSuite) TestWithRetry() {
	t := su.T()
	r := RetryOptions{MaxRetries: 5, RetryOnRateLimit: true}
	c := build(WithRetry(r))
	if c.Retry != r {
		t.Fatalf("got %+v", c.Retry)
	}
}

func (su *withOptionsSuite) TestWithMinRequestInterval() {
	t := su.T()
	c := build(WithMinRequestInterval(250 * time.Millisecond))
	if c.MinRequestInterval != 250*time.Millisecond {
		t.Fatalf("got %v", c.MinRequestInterval)
	}
}

func (su *withOptionsSuite) TestWithBaseURL() {
	t := su.T()
	c := build(WithBaseURL("https://proxy.example.com"))
	if c.BaseURL != "https://proxy.example.com" {
		t.Fatalf("got %q", c.BaseURL)
	}
}

func (su *withOptionsSuite) TestWithMaxResponseBodySize() {
	t := su.T()
	c := build(WithMaxResponseBodySize(-1))
	if c.MaxResponseBodySize != -1 {
		t.Fatalf("got %d", c.MaxResponseBodySize)
	}
}

func (su *withOptionsSuite) TestWithRateLimiting() {
	t := su.T()
	c := build(WithRateLimiting(RateLimiterOptions{SafetyMargin: 2}))
	if c.RateLimiter == nil || c.RateLimiter.SafetyMargin != 2 {
		t.Fatalf("got %+v", c.RateLimiter)
	}
}

func (su *withOptionsSuite) TestWithMaxReconnectRetries() {
	t := su.T()
	c := build(WithMaxReconnectRetries(-1))
	if c.MaxReconnectRetries != -1 {
		t.Fatalf("got %d", c.MaxReconnectRetries)
	}
}

func (su *withOptionsSuite) TestWithMaxConcurrentEvents() {
	t := su.T()
	c := build(WithMaxConcurrentEvents(128))
	if c.MaxConcurrentEvents != 128 {
		t.Fatalf("got %d", c.MaxConcurrentEvents)
	}
}

func (su *withOptionsSuite) TestWithLogLevel() {
	t := su.T()
	c := build(WithLogLevel(slog.LevelDebug))
	if c.LogLevel == nil || *c.LogLevel != slog.LevelDebug {
		t.Fatalf("got %v", c.LogLevel)
	}
}

func (su *withOptionsSuite) TestWithCache() {
	t := su.T()
	mc := cache.NewMemoryCache(cache.Options{})
	c := build(WithCache(mc))
	if c.Cache != mc {
		t.Fatal("Cache not set")
	}
}

// stubCoordinator is a no-op ShardCoordinator for option testing.
type stubCoordinator struct{}

func (stubCoordinator) Register(int, func(ShardMessage)) error { return nil }
func (stubCoordinator) Send(ShardMessage) error                { return nil }
func (stubCoordinator) Broadcast(ShardMessage) error           { return nil }
func (stubCoordinator) TotalShards() int                       { return 1 }
func (stubCoordinator) Close() error                           { return nil }

func (su *withOptionsSuite) TestWithCoordinator() {
	t := su.T()
	co := stubCoordinator{}
	c := build(WithCoordinator(co))
	if c.Coordinator != co {
		t.Fatal("Coordinator not set")
	}
}

// ── Cache-store auto-population interplay ─────────────────────────────────────

func (su *withOptionsSuite) TestWithDisableCacheAutoPopulation() {
	t := su.T()
	c := build(WithCacheStores(cache.CategoryGuilds), WithDisableCacheAutoPopulation())
	if !c.DisableCacheAutoPopulation {
		t.Fatal("expected DisableCacheAutoPopulation=true")
	}
	if c.CacheStores != 0 {
		t.Fatalf("expected CacheStores reset to 0, got %d", c.CacheStores)
	}
}

func (su *withOptionsSuite) TestWithCacheStores_ReenablesAutoPopulation() {
	t := su.T()
	// WithCacheStores after a disable re-enables auto-population.
	c := build(WithDisableCacheAutoPopulation(), WithCacheStores(cache.CategoryGuilds|cache.CategoryChannels))
	if c.DisableCacheAutoPopulation {
		t.Fatal("expected WithCacheStores to clear DisableCacheAutoPopulation")
	}
	if c.CacheStores != cache.CategoryGuilds|cache.CategoryChannels {
		t.Fatalf("got %d", c.CacheStores)
	}
}

func (su *withOptionsSuite) TestWithDisableCacheStore_FromDefault() {
	t := su.T()
	// CacheStores starts at 0, so the option seeds CategoryAll then clears bits.
	c := build(WithDisableCacheStore(cache.CategoryMessages))
	if c.CacheStores&cache.CategoryMessages != 0 {
		t.Fatal("expected CategoryMessages cleared")
	}
	if c.CacheStores&cache.CategoryGuilds == 0 {
		t.Fatal("expected other categories (Guilds) still present")
	}
}

func (su *withOptionsSuite) TestWithDisableCacheStore_FromExistingStores() {
	t := su.T()
	// When CacheStores is already non-zero, the option clears bits without
	// re-seeding CategoryAll.
	c := build(
		WithCacheStores(cache.CategoryGuilds|cache.CategoryChannels),
		WithDisableCacheStore(cache.CategoryChannels),
	)
	if c.CacheStores != cache.CategoryGuilds {
		t.Fatalf("expected only Guilds to remain, got %d", c.CacheStores)
	}
}

func (su *withOptionsSuite) TestBuild_AppliesInOrder() {
	t := su.T()
	// Later options override earlier ones.
	c := Build(Config{}, []Option{
		WithBaseURL("first"),
		WithBaseURL("second"),
	})
	if c.BaseURL != "second" {
		t.Fatalf("got %q", c.BaseURL)
	}
}

type withOptionsSuite struct{ suite.Suite }

func TestWithOptionsSuite(t *testing.T) { suite.Run(t, new(withOptionsSuite)) }
