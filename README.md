# go-discord-wrapper

A Go library for the Discord gateway and REST API.

> **Status:** In active development toward a stable v0.1.0 release. Public APIs may still change. Bug reports and pull requests are welcome.

## Features

- Gateway (WebSocket) client with automatic reconnection and resume
- Full REST API coverage for guilds, channels, messages, roles, interactions, webhooks, and more
- Proactive rate limiter — respects per-route and global Discord rate limits before sending
- In-process memory cache (guild, channel, member, message, role, voice state, emoji, sticker, …) with TTL and LRU eviction
- Pluggable cache backends (Redis and MongoDB drivers included)
- Slash command and interaction support (reply, defer, modal, component)
- Horizontal sharding with a built-in shard coordinator
- Middleware support for event handlers

## Installation

```sh
go get github.com/streame-gg/go-discord-wrapper
```

Requires Go 1.21 or later.

## Quick start

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/streame-gg/go-discord-wrapper/connection"
    "github.com/streame-gg/go-discord-wrapper/types/common"
    "github.com/streame-gg/go-discord-wrapper/types/events"
    "github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

func main() {
    bot, err := connection.NewClient(os.Getenv("DISCORD_TOKEN"), common.IntentGuilds|common.IntentGuildMessages)
    if err != nil {
        slog.Error("failed to create client", "err", err)
        os.Exit(1)
    }

    bot.OnEvent(events.EventMessageCreate, func(c *connection.Client, e *events.MessageCreateEvent) {
        slog.Info("new message", "author", e.Author.Username, "content", e.Content)
    })

    bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, e *events.InteractionCreateEvent) {
        if e.IsCommand() && e.GetFullCommand() == "ping" {
            _, _ = c.Reply(context.Background(), &e.Interaction, &responses.InteractionResponseDataDefault{
                Content: "Pong!",
            }, false)
        }
    })

    if err := bot.Login(context.Background()); err != nil {
        slog.Error("login failed", "err", err)
        os.Exit(1)
    }

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    <-ctx.Done()

    _ = bot.Shutdown()
}
```

## Cache

Pass a cache to `NewClient` and the gateway will populate it automatically:

```go
import "github.com/streame-gg/go-discord-wrapper/cache"

mc := cache.NewMemoryCache(cache.Options{
    TTL:           30 * time.Minute,
    SweepInterval: 5 * time.Minute,
    Limits: cache.Limits{
        MaxMessages: 10_000,
        MaxUsers:    5_000,
    },
    Messages: cache.MessageOptions{
        MaxPerChannel: 200,
    },
})

bot, err := connection.NewClient(token, intents, options.WithCache(mc))
```

Look up entities without hitting the REST API:

```go
if ch, ok := bot.Cache.Channels().Get(channelID); ok {
    fmt.Println(ch.Name)
}
```

### Cache backends

| Backend | Import |
|---------|--------|
| In-process memory (default) | `cache.NewMemoryCache(...)` |
| Redis | `cache/redis` package |
| MongoDB | `cache/mongo` package |

## Sharding

```go
// Shard 0 of 4 total shards.
bot, err := connection.NewClient(token, intents, options.WithSharding(4, 0))
```

## Slash commands

```go
_, err := bot.BulkRegisterCommands(ctx, []*commands.ApplicationCommand{
    {
        Name:        "ping",
        Description: "Replies with Pong!",
        Type:        common.ApplicationCommandTypeChatInput,
    },
})
```

## Configuration reference

| Option | Description |
|--------|-------------|
| `options.WithCache(c)` | Attach a cache backend |
| `options.WithSharding(total, shardID)` | Enable sharding |
| `options.WithLogger(l)` | Custom `*slog.Logger` |
| `options.WithRetry(opts)` | Retry policy for REST requests |
| `options.WithRateLimiting(opts)` | Rate-limiter tuning (safety margin, disable) |
| `options.WithMinRequestInterval(d)` | Global minimum delay between REST requests |
| `options.WithAPIVersion(v)` | Discord API version (default: v10) |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

Join the Discord server: https://discord.gg/nfBuZVejqp

## License

Apache 2.0 — see [LICENSE](LICENSE).
