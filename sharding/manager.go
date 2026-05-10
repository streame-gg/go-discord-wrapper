// Package sharding provides a local shard coordinator that manages multiple
// gateway connections and routes cross-shard messages between them.
package sharding

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
)

// identifyInterval is the minimum time Discord requires between IDENTIFY payloads
// on the same bot token. Violating this results in a 4008 (rate limited) close.
const identifyInterval = 5 * time.Second

// ShardManager orchestrates multiple Discord gateway connections in a single process.
// It creates one Client per shard using a factory function, enforces Discord's
// 5-second IDENTIFY rate limit during startup, and provides convenience accessors.
//
// For multi-machine deployments you do not need a ShardManager — run one Client
// per machine with options.WithCoordinator pointing at your external broker.
type ShardManager struct {
	total       int
	coordinator options.ShardCoordinator
	factory     func(shardID, totalShards int) *connection.Client
	clients     []*connection.Client
}

// NewShardManager creates a manager that will start total shards using coord.
//
// factory is called once per shard (0 … total-1) during Start. It receives the
// shard ID and total shard count and must return a fully configured but not-yet-
// logged-in Client. Typically you pass options.WithSharding and
// options.WithCoordinator(coord) inside factory:
//
//	coord := sharding.NewLocalCoordinator(4)
//	mgr := sharding.NewShardManager(coord, func(shardID, total int) *connection.Client {
//	    return connection.NewClient(token, intents,
//	        options.WithSharding(total, shardID),
//	        options.WithCoordinator(coord),
//	    )
//	})
func NewShardManager(
	coord options.ShardCoordinator,
	factory func(shardID, totalShards int) *connection.Client,
) *ShardManager {
	total := coord.TotalShards()
	return &ShardManager{
		total:       total,
		coordinator: coord,
		factory:     factory,
		clients:     make([]*connection.Client, total),
	}
}

// Start connects all shards to the Discord gateway sequentially.
//
// Discord mandates at least 5 seconds between IDENTIFY payloads for the same
// bot token (across all shards). Start enforces this automatically, so the
// full startup time for N shards is at least (N-1)*5 seconds.
//
// Start blocks until all shards have sent their IDENTIFY and received READY.
// If any shard fails to connect, Start shuts down all already-started shards
// before returning the error so no shards are left orphaned.
func (m *ShardManager) Start() error {
	for id := 0; id < m.total; id++ {
		client := m.factory(id, m.total)

		if err := client.Login(context.Background()); err != nil {
			// Shut down all shards that started successfully before this failure.
			for _, c := range m.clients[:id] {
				if c != nil {
					_ = c.Shutdown()
				}
			}
			return fmt.Errorf("sharding: shard %d failed to connect: %w", id, err)
		}

		m.clients[id] = client

		// Respect Discord's rate limit between IDENTIFYs.
		if id < m.total-1 {
			time.Sleep(identifyInterval)
		}
	}
	return nil
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
