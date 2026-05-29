// Package sharding provides a local shard coordinator that manages multiple
// gateway connections and routes cross-shard messages between them.
package sharding

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
)

// identifyInterval is the minimum time Discord requires between IDENTIFY buckets.
// Violating this results in a 4008 (rate limited) close.
const identifyInterval = 5 * time.Second

const defaultLoginTimeout = 30 * time.Second

// ShardManagerOption configures a ShardManager.
type ShardManagerOption func(*ShardManager)

// WithLoginTimeout sets how long startOneShard waits for a single shard's Login
// call before giving up. Defaults to 30 seconds.
func WithLoginTimeout(d time.Duration) ShardManagerOption {
	return func(m *ShardManager) { m.loginTimeout = d }
}

// ShardManager orchestrates multiple Discord gateway connections in a single process.
// It creates one Client per shard using a factory function, respects Discord's
// max_concurrency IDENTIFY bucket rule during startup, and provides convenience accessors.
//
// For multi-machine deployments you do not need a ShardManager — run one Client
// per machine with options.WithCoordinator pointing at your external broker.
type ShardManager struct {
	total        int
	coordinator  options.ShardCoordinator
	factory      func(shardID, totalShards int) (*connection.Client, error)
	rest         *api.RestClient
	clients      []*connection.Client
	loginTimeout time.Duration
}

// shardErr pairs a shard ID with the error returned during startup.
type shardErr struct {
	id  int
	err error
}

// NewShardManager creates a manager that will start total shards using coord.
//
// factory is called once per shard (0 … total-1) during Start. It receives the
// shard ID and total shard count and must return a fully configured but not-yet-
// logged-in Client or an error. Typically you pass options.WithSharding and
// options.WithCoordinator(coord) inside factory:
//
//	coord := sharding.NewLocalCoordinator(4)
//	mgr := sharding.NewShardManager(coord, func(shardID, total int) (*connection.Client, error) {
//	    return connection.NewClient(token, intents,
//	        options.WithSharding(total, shardID),
//	        options.WithCoordinator(coord),
//	    )
//	})
//
// This constructor defaults to max_concurrency=1 (sequential start). For faster
// startup on larger bots, use NewShardManagerWithRest and pass a RestClient so
// that Start can fetch the actual max_concurrency from Discord's /gateway/bot.
func NewShardManager(
	coord options.ShardCoordinator,
	factory func(shardID, totalShards int) (*connection.Client, error),
	opts ...ShardManagerOption,
) *ShardManager {
	return NewShardManagerWithRest(coord, factory, nil, opts...)
}

// NewShardManagerWithRest creates a manager that fetches Discord's max_concurrency
// from /gateway/bot before starting shards.
//
// When rest is non-nil, Start calls rest.GetBotGateway to determine how many shards
// may be started in parallel per IDENTIFY bucket (e.g. 16 for large bots). If the
// fetch fails, Start falls back to max_concurrency=1. When rest is nil, Start always
// uses max_concurrency=1, which is identical to NewShardManager behaviour.
//
// Migration note: to unlock parallel bucket start on large bots, replace:
//
//	sharding.NewShardManager(coord, factory)
//
// with:
//
//	rest, _ := api.NewRestClient(token)
//	sharding.NewShardManagerWithRest(coord, factory, rest)
func NewShardManagerWithRest(
	coord options.ShardCoordinator,
	factory func(shardID, totalShards int) (*connection.Client, error),
	rest *api.RestClient,
	opts ...ShardManagerOption,
) *ShardManager {
	total := coord.TotalShards()
	m := &ShardManager{
		total:        total,
		coordinator:  coord,
		factory:      factory,
		rest:         rest,
		clients:      make([]*connection.Client, total),
		loginTimeout: defaultLoginTimeout,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Start connects all shards to the Discord gateway using Discord's bucket-based
// IDENTIFY rate limit.
//
// Shards within the same max_concurrency bucket are started concurrently; a
// 5-second gap is enforced between buckets per the Discord spec. When
// NewShardManagerWithRest was used with a non-nil RestClient, Start fetches
// max_concurrency from Discord's /gateway/bot endpoint before the loop begins.
// Without a RestClient it defaults to max_concurrency=1 (sequential, one shard
// per bucket), which is identical to the previous behaviour.
//
// Start blocks until all shards have sent their IDENTIFY and received READY.
// If any shard in a bucket fails, Start shuts down all already-started shards
// before returning the error so no shards are left orphaned.
func (m *ShardManager) Start() error {
	maxConcurrency := m.fetchMaxConcurrency()

	for bucketStart := 0; bucketStart < m.total; bucketStart += maxConcurrency {
		bucketEnd := bucketStart + maxConcurrency
		if bucketEnd > m.total {
			bucketEnd = m.total
		}

		errCh := make(chan shardErr, bucketEnd-bucketStart)
		var wg sync.WaitGroup
		for id := bucketStart; id < bucketEnd; id++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				if err := m.startOneShard(id); err != nil {
					errCh <- shardErr{id: id, err: err}
				}
			}(id)
		}
		wg.Wait()
		close(errCh)

		var bucketErrs []error
		for se := range errCh {
			bucketErrs = append(bucketErrs, fmt.Errorf("shard %d: %w", se.id, se.err))
		}
		if len(bucketErrs) > 0 {
			for _, c := range m.clients {
				if c != nil {
					_ = c.Shutdown()
				}
			}
			return fmt.Errorf("sharding: bucket starting at shard %d failed: %w",
				bucketStart, errors.Join(bucketErrs...))
		}

		if bucketEnd < m.total {
			time.Sleep(identifyInterval)
		}
	}
	return nil
}

// startOneShard runs factory + Login for a single shard and stores the result.
// Safe to call concurrently for different shard IDs (each goroutine writes to a
// distinct slice index, so no mutex is needed for m.clients).
func (m *ShardManager) startOneShard(id int) error {
	client, err := m.factory(id, m.total)
	if err != nil {
		return fmt.Errorf("factory: %w", err)
	}
	if client == nil {
		return errors.New("factory returned nil client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.loginTimeout)
	defer cancel()
	if err := client.Login(ctx); err != nil {
		return fmt.Errorf("login: %w", err)
	}
	m.clients[id] = client
	return nil
}

// fetchMaxConcurrency asks Discord what IDENTIFY bucket size to use.
// Falls back to 1 (safe — slower start but never violates spec) if rest is nil
// or the request fails.
func (m *ShardManager) fetchMaxConcurrency() int {
	if m.rest == nil {
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := m.rest.GetBotGateway(ctx)
	if err != nil {
		return 1
	}
	mc := resp.SessionStartLimit.MaxConcurrency
	if mc < 1 {
		mc = 1
	}
	return mc
}

// Shard returns the Client for the given shard ID.
// Returns nil if shardID is out of range or that shard has not yet started.
func (m *ShardManager) Shard(shardID int) *connection.Client {
	if shardID < 0 || shardID >= m.total {
		return nil
	}
	return m.clients[shardID]
}

// Shards returns a snapshot of all shard clients in ascending shard-ID order.
// Entries are nil for shards that have not yet been started.
func (m *ShardManager) Shards() []*connection.Client {
	out := make([]*connection.Client, m.total)
	copy(out, m.clients)
	return out
}

// Shutdown disconnects every shard and closes the coordinator.
// All shards are shut down even if one returns an error; all errors are joined.
func (m *ShardManager) Shutdown() error {
	var errs []error
	for _, c := range m.clients {
		if c != nil {
			if err := c.Shutdown(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if err := m.coordinator.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
