// sharding demonstrates running multiple gateway connections in one process
// with a ShardManager and a local in-process coordinator.
//
// A ShardManager:
//   - creates one Client per shard via your factory function,
//   - starts shards in IDENTIFY buckets, respecting Discord's max_concurrency,
//   - and routes cross-shard messages through the coordinator.
//
// For multi-machine deployments you do not need a ShardManager: run one Client
// per machine and point options.WithCoordinator at an external broker
// (Redis, NATS, …) that implements options.ShardCoordinator.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/sharding"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

func main() {
	const totalShards = 4
	token := os.Getenv("DISCORD_TOKEN")

	// The coordinator wires the shards together so a handler on one shard can
	// broadcast to, or message, another shard in the same process.
	coord := sharding.NewLocalCoordinator(totalShards)

	// The factory is called once per shard ID. It must return a configured but
	// not-yet-logged-in Client. WithSharding tells Discord which slice of guilds
	// this connection is responsible for; WithCoordinator joins it to the group.
	factory := func(shardID, total int) (*connection.Client, error) {
		bot, err := connection.NewClient(token, discord.IntentGuilds|discord.IntentGuildMessages,
			options.WithSharding(total, shardID),
			options.WithCoordinator(coord),
		)
		if err != nil {
			return nil, err
		}

		bot.OnReady(func(c *connection.Client, e *events.ReadyEvent) {
			c.Logger.Info("shard ready", slog.Int("shard", shardID), slog.Int("of", total))
		})
		return bot, nil
	}

	mgr := sharding.NewShardManager(coord, factory)

	// Start blocks until every shard has IDENTIFY'd and received READY. If any
	// shard fails, already-started shards are shut down before Start returns.
	if err := mgr.Start(); err != nil {
		slog.Error("failed to start shards", "err", err)
		os.Exit(1)
	}
	slog.Info("all shards online", slog.Int("total", len(mgr.Shards())))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	// Shutdown disconnects every shard and closes the coordinator.
	if err := mgr.Shutdown(); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}
