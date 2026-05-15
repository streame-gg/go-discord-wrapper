// Package sharding provides multi-shard orchestration for Discord bots.
//
// For single-machine deployments, use LocalCoordinator and ShardManager:
//
//	coord := sharding.NewLocalCoordinator(4)
//	mgr := sharding.NewShardManager(coord, func(shardID, total int) *connection.Client {
//	    return connection.NewClient(token, intents,
//	        options.WithSharding(total, shardID),
//	        options.WithCoordinator(coord),
//	    )
//	})
//	if err := mgr.Start(); err != nil { panic(err) }
//
// For multi-machine deployments, implement options.ShardCoordinator with your
// preferred message broker and pass it via options.WithCoordinator on each
// machine independently — no ShardManager is needed across machines.
package sharding

import (
	"fmt"
	"sync"

	"github.com/streame-gg/go-discord-wrapper/options"
)

// LocalCoordinator is an in-process implementation of options.ShardCoordinator.
//
// All shards must run within the same Go process. Messages are delivered
// asynchronously via goroutines, so handlers must be safe for concurrent use.
//
// For distributed (multi-machine) deployments, implement options.ShardCoordinator
// with a message broker such as Redis Pub/Sub, NATS, or a gRPC stream instead.
type LocalCoordinator struct {
	total int

	mu       sync.RWMutex
	handlers map[int]func(options.ShardMessage)
	closed   bool
	wg       sync.WaitGroup // tracks in-flight handler goroutines (Bug 27)
}

// NewLocalCoordinator creates a coordinator for totalShards in-process shards.
func NewLocalCoordinator(totalShards int) *LocalCoordinator {
	return &LocalCoordinator{
		total:    totalShards,
		handlers: make(map[int]func(options.ShardMessage), totalShards),
	}
}

// Register installs the message handler for shardID.
// Returns an error if shardID is out of range or the coordinator is closed.
func (c *LocalCoordinator) Register(shardID int, handler func(options.ShardMessage)) error {
	if shardID < 0 || shardID >= c.total {
		return fmt.Errorf("sharding: shard ID %d out of range [0, %d)", shardID, c.total)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return fmt.Errorf("sharding: coordinator is closed")
	}
	c.handlers[shardID] = handler
	return nil
}

// Send delivers msg to the shard identified by msg.To.
// Delivery is asynchronous; the goroutine is started before Send returns.
func (c *LocalCoordinator) Send(msg options.ShardMessage) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return fmt.Errorf("sharding: coordinator is closed")
	}
	h, ok := c.handlers[msg.To]
	if ok {
		c.wg.Add(1) // increment while holding RLock so Close cannot race wg.Wait
	}
	c.mu.RUnlock()
	if !ok {
		return fmt.Errorf("sharding: shard %d is not registered", msg.To)
	}
	go func() {
		defer c.wg.Done()
		h(msg)
	}()
	return nil
}

// Broadcast delivers msg to every registered shard, including the sender.
// Each handler is called in its own goroutine.
func (c *LocalCoordinator) Broadcast(msg options.ShardMessage) error {
	msg.To = options.BroadcastAll

	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return fmt.Errorf("sharding: coordinator is closed")
	}
	hs := make([]func(options.ShardMessage), 0, len(c.handlers))
	for _, h := range c.handlers {
		hs = append(hs, h)
	}
	c.wg.Add(len(hs)) // increment while holding RLock so Close cannot race wg.Wait
	c.mu.RUnlock()

	for _, h := range hs {
		h := h
		go func() {
			defer c.wg.Done()
			h(msg)
		}()
	}
	return nil
}

// TotalShards returns the number of shards this coordinator was created with.
func (c *LocalCoordinator) TotalShards() int { return c.total }

// Close removes all registered handlers, waits for in-flight handler goroutines
// to finish, and prevents further message delivery.
func (c *LocalCoordinator) Close() error {
	c.mu.Lock()
	c.closed = true
	c.handlers = nil
	c.mu.Unlock()
	c.wg.Wait()
	return nil
}
