# Getting started

This walks you from an empty folder to a running bot that answers a slash
command. If you'd rather start from a structured project, copy
[`example/template`](../example/template) instead — it has the same pieces wired
up with auto-loading command/event folders.

## Prerequisites

- Go 1.26 or later.
- A bot application and token from the
  [Discord Developer Portal](https://discord.com/developers/applications).
- The bot invited to a server with the `applications.commands` scope.

## 1. Create the project

```sh
mkdir mybot && cd mybot
go mod init example.com/mybot
go get github.com/streame-gg/go-discord-wrapper
```

## 2. Connect and listen

Create `main.go`:

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/streame-gg/go-discord-wrapper/connection"
    "github.com/streame-gg/go-discord-wrapper/types/discord"
    "github.com/streame-gg/go-discord-wrapper/types/events"
)

func main() {
    bot, err := connection.NewClient(
        os.Getenv("DISCORD_TOKEN"),
        discord.IntentGuilds|discord.IntentGuildMessages,
    )
    if err != nil {
        slog.Error("create client", "err", err)
        os.Exit(1)
    }

    bot.OnReady(func(c *connection.Client, e *events.ReadyEvent) {
        c.Logger.Info("logged in", "as", e.User.Username)
    })

    if err := bot.Login(context.Background()); err != nil {
        slog.Error("login", "err", err)
        os.Exit(1)
    }

    // Block until Ctrl-C, then shut down cleanly.
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    <-ctx.Done()
    _ = bot.Shutdown()
}
```

Run it:

```sh
DISCORD_TOKEN=your-token go run .
```

You should see `logged in` in the logs.

## 3. Add a slash command

Register a command after login and handle its interaction:

```go
bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
    if ev.IsCommand() && ev.GetFullCommand() == "ping" {
        _, _ = ev.Reply(interactions.ReplyOptions{Content: "Pong!"})
    }
})

// …after bot.Login(…):
_, _ = bot.BulkRegisterCommands(context.Background(), []*commands.ApplicationCommand{
    {Name: "ping", Description: "Replies with Pong!", Type: discord.ApplicationCommandTypeChatInput},
})
```

(Add the `commands` and `interactions` imports.) Global commands can take up to
an hour to appear — register to a test guild during development, see
[COMMAND_MANAGEMENT.md](COMMAND_MANAGEMENT.md).

## 4. Turn on the cache (optional)

Pass a cache and the gateway populates it automatically, so you can read guilds,
channels, and members without REST calls:

```go
import (
    "time"
    "github.com/streame-gg/go-discord-wrapper/cache"
    "github.com/streame-gg/go-discord-wrapper/options"
)

mc := cache.NewMemoryCache(cache.Options{TTL: 30 * time.Minute})
bot, err := connection.NewClient(token, intents, options.WithCache(mc))
```

See [CACHE.md](CACHE.md) for what's stored and when.

## Intents cheat-sheet

You only receive events for the intents you request. Common ones:

| Intent | Unlocks |
|--------|---------|
| `IntentGuilds` | guild/channel/role create/update/delete |
| `IntentGuildMessages` | message events in servers |
| `IntentMessageContent` | message `Content` (privileged) |
| `IntentGuildMembers` | member join/leave/update (privileged) |
| `IntentGuildPresences` | presence updates (privileged) |

Privileged intents must be enabled in the Developer Portal first.

## Where to next

| You want to… | Read |
|--------------|------|
| Structure a real project | [`example/template`](../example/template) |
| Add command options & subcommands | [COMMANDS.md](COMMANDS.md) |
| Manage/edit/delete commands | [COMMAND_MANAGEMENT.md](COMMAND_MANAGEMENT.md) |
| Add buttons & select menus | [COMPONENTS.md](COMPONENTS.md) |
| Show a form | [MODALS.md](MODALS.md) |
| Send embeds | [EMBEDS.md](EMBEDS.md) |
| Tune the client | [CONFIGURATION.md](CONFIGURATION.md) |
| Scale past ~2,500 guilds | [SHARDING.md](SHARDING.md) |
