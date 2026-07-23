package options

import (
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/suite"
)

type optionsSuite struct{ suite.Suite }

func TestOptionsSuite(t *testing.T) { suite.Run(t, new(optionsSuite)) }

func (s *optionsSuite) TestValidate() {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with sharding",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 0},
			},
			wantErr: false,
		},
		{
			name: "shard ID equals total shards",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 4},
			},
			wantErr: true,
		},
		{
			name: "shard ID greater than total shards",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 5},
			},
			wantErr: true,
		},
		{
			name: "total shards zero",
			config: Config{
				Sharding: &Sharding{TotalShards: 0, ShardID: 0},
			},
			wantErr: true,
		},
		{
			name: "total shards negative",
			config: Config{
				Sharding: &Sharding{TotalShards: -1, ShardID: 0},
			},
			wantErr: true,
		},
		{
			name: "shard ID negative",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: -1},
			},
			wantErr: true,
		},
		{
			name: "negative MaxRetries",
			config: Config{
				Retry: RetryOptions{MaxRetries: -1},
			},
			wantErr: true,
		},
		{
			name: "negative MinRequestInterval",
			config: Config{
				MinRequestInterval: -1 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "negative SafetyMargin",
			config: Config{
				RateLimiter: &RateLimiterOptions{SafetyMargin: -1},
			},
			wantErr: true,
		},
		{
			name: "nil sharding with negative MaxRetries",
			config: Config{
				Sharding: nil,
				Retry:    RetryOptions{MaxRetries: -1},
			},
			wantErr: true,
		},
		{
			name:    "all zero values",
			config:  Config{},
			wantErr: false,
		},
		{
			name: "zero MaxRetries is valid",
			config: Config{
				Retry: RetryOptions{MaxRetries: 0},
			},
			wantErr: false,
		},
		{
			name: "zero MinRequestInterval is valid",
			config: Config{
				MinRequestInterval: 0,
			},
			wantErr: false,
		},
		{
			name: "safety margin of zero is valid",
			config: Config{
				RateLimiter: &RateLimiterOptions{SafetyMargin: 0},
			},
			wantErr: false,
		},
		{
			name: "last valid shard ID",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 3},
			},
			wantErr: false,
		},
		{
			name: "negative number is rejected",
			config: Config{
				MaxReconnectRetries: -2,
			},
			wantErr: true,
		},
		{
			name: "-1 is infinite retries",
			config: Config{
				MaxReconnectRetries: -1,
			},
			wantErr: false,
		},
		{
			name: "0 applies default",
			config: Config{
				MaxReconnectRetries: 0,
			},
			wantErr: false,
		},
		{
			name: "accept positive number",
			config: Config{
				MaxReconnectRetries: 10,
			},
			wantErr: false,
		},
		{
			name: "reject negative number",
			config: Config{
				MaxReconnectRetries: -1,
			},
			wantErr: true,
		},
		{
			name: "apply default on 0",
			config: Config{
				MaxReconnectRetries: 0,
			},
			wantErr: false,
		},
		{
			name: "accept positive number",
			config: Config{
				MaxReconnectRetries: 64,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			if tc.wantErr {
				s.Error(tc.config.Validate())
			} else {
				s.NoError(tc.config.Validate())
			}
		})
	}
}

// build is a small helper that applies opts onto a zero Config.
func build(opts ...Option) Config {
	return Build(Config{}, opts)
}

func (s *optionsSuite) TestWithSharding() {
	c := build(WithSharding(4, 2))
	s.NotNil(c.Sharding)
	s.Equal(4, c.Sharding.TotalShards)
	s.Equal(2, c.Sharding.ShardID)
}

func (s *optionsSuite) TestWithAPIVersion() {
	c := build(WithAPIVersion(discord.APIVersion(9)))
	s.Equal(discord.APIVersion(9), c.APIVersion)
}

func (s *optionsSuite) TestWithLogger() {
	l := slog.Default()
	c := build(WithLogger(l))
	s.Equal(l, c.Logger)
}

func (s *optionsSuite) TestWithHTTPClient() {
	hc := &http.Client{Timeout: time.Second}
	c := build(WithHTTPClient(hc))
	s.Equal(hc, c.HTTPClient)
}

func (s *optionsSuite) TestWithRetry() {
	r := RetryOptions{MaxRetries: 5, RetryOnRateLimit: true}
	c := build(WithRetry(r))
	s.Equal(r, c.Retry)
}

func (s *optionsSuite) TestWithMinRequestInterval() {
	c := build(WithMinRequestInterval(250 * time.Millisecond))
	s.Equal(250*time.Millisecond, c.MinRequestInterval)
}

func (s *optionsSuite) TestWithBaseURL() {
	c := build(WithBaseURL("https://proxy.example.com"))
	s.Equal("https://proxy.example.com", c.BaseURL)
}

func (s *optionsSuite) TestWithMaxResponseBodySize() {
	c := build(WithMaxResponseBodySize(-1))
	s.Equal(-1, c.MaxResponseBodySize)
}

func (s *optionsSuite) TestWithRateLimiting() {
	c := build(WithRateLimiting(RateLimiterOptions{SafetyMargin: 2}))
	s.NotNil(c.RateLimiter)
	s.Equal(2, c.RateLimiter.SafetyMargin)
}

func (s *optionsSuite) TestWithMaxReconnectRetries() {
	c := build(WithMaxReconnectRetries(-1))
	s.Equal(-1, c.MaxReconnectRetries)
}

func (s *optionsSuite) TestWithMaxConcurrentEvents() {
	c := build(WithMaxConcurrentEvents(128))
	s.Equal(128, c.MaxConcurrentEvents)
}

func (s *optionsSuite) TestWithLogLevel() {
	c := build(WithLogLevel(slog.LevelDebug))
	s.NotNil(c.LogLevel)
	s.Equal(slog.LevelDebug, *c.LogLevel)
}

func (s *optionsSuite) TestWithCache() {
	mc := cache.NewMemoryCache(cache.Options{})
	c := build(WithCache(mc))
	s.Equal(mc, c.Cache)
}

type stubCoordinator struct{}

func (stubCoordinator) Register(int, func(ShardMessage)) error { return nil }
func (stubCoordinator) Send(ShardMessage) error                { return nil }
func (stubCoordinator) Broadcast(ShardMessage) error           { return nil }
func (stubCoordinator) TotalShards() int                       { return 1 }
func (stubCoordinator) Close() error                           { return nil }

func (s *optionsSuite) TestWithCoordinator() {
	co := stubCoordinator{}
	c := build(WithCoordinator(co))
	s.Equal(co, c.Coordinator)
}

func (s *optionsSuite) TestWithDisableCacheAutoPopulation() {
	c := build(WithCacheStores(cache.CategoryGuilds), WithDisableCacheAutoPopulation())
	s.True(c.DisableCacheAutoPopulation)
	s.Equal(0, c.CacheStores)
}

func (s *optionsSuite) TestWithCacheStores_ReenablesAutoPopulation() {
	c := build(WithDisableCacheAutoPopulation(), WithCacheStores(cache.CategoryGuilds|cache.CategoryChannels))
	s.False(c.DisableCacheAutoPopulation)
	s.Equal(cache.CategoryGuilds|cache.CategoryChannels, c.CacheStores)
}

func (s *optionsSuite) TestWithDisableCacheStore_FromDefault() {
	c := build(WithDisableCacheStore(cache.CategoryMessages))
	s.False(c.CacheStores&cache.CategoryMessages != 0)
	s.True(c.CacheStores&cache.CategoryGuilds != 0)
}

func (s *optionsSuite) TestWithDisableCacheStore_FromExistingStores() {
	c := build(
		WithCacheStores(cache.CategoryGuilds|cache.CategoryChannels),
		WithDisableCacheStore(cache.CategoryChannels),
	)
	s.Equal(cache.CategoryGuilds, c.CacheStores)
}

func (s *optionsSuite) TestBuild_AppliesInOrder() {
	c := Build(Config{}, []Option{
		WithBaseURL("first"),
		WithBaseURL("second"),
	})
	s.Equal("second", c.BaseURL)
}
