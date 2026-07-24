# go-discord-wrapper

# IMPORTANT
The project is currently in a state where I manually have to rewrite many parts of it. This is taking a very long time and therefore, a first stable release will take some time to finish. 

[![CI](https://github.com/streame-gg/go-discord-wrapper/actions/workflows/go.yml/badge.svg)](https://github.com/streame-gg/go-discord-wrapper/actions/workflows/go.yml)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/streame-gg/go-discord-wrapper/master/.github/badges/coverage.json)](https://github.com/streame-gg/go-discord-wrapper/blob/master/COVERAGE.md)
[![Version](https://img.shields.io/github/v/tag/streame-gg/go-discord-wrapper?sort=semver&label=version)](https://github.com/streame-gg/go-discord-wrapper/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/streame-gg/go-discord-wrapper.svg)](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/streame-gg/go-discord-wrapper)](https://goreportcard.com/report/github.com/streame-gg/go-discord-wrapper)

A Go library for the Discord gateway and REST API.

## Features

- Gateway (WebSocket) client with automatic reconnection and resume
- Near-complete REST API coverage for guilds, channels, messages, roles, interactions, webhooks, and more (~224 endpoints; see [known gaps](docs/COVERAGE_GAPS.md))
- Proactive rate limiter — respects per-route and global Discord rate limits before sending
- In-process memory cache (guild, channel, member, message, role, voice state, emoji, sticker, …) with TTL and LRU eviction
- Pluggable cache backends (Redis and MongoDB drivers included)
- Synthetic emoji and sticker events derived from raw gateway payloads
- Slash command and interaction support (reply, defer, modal, component)
- Horizontal sharding with a built-in shard coordinator
- Middleware support for event handlers

Also see Features in the documentation: [docs/README.md](docs/README.md).

## Project status & v1.0 scope

**v1.0 — first stable release.** The public API is frozen and follows
[semantic versioning](https://semver.org/): no breaking changes within the v1
line. That said, this is a young library — expect to hit bugs in less-travelled
corners. Please [open an issue](https://github.com/streame-gg/go-discord-wrapper/issues)
when you do; reports against v1.0 are very welcome.

### How this library was built

The original foundation of this library was written by hand. Most of the surface
area beyond that — the breadth of REST endpoints, event handling, and supporting
code — was built with substantial AI assistance and has **not yet been reviewed
line-by-line by a human**. Correctness is instead guarded by an extensive
automated test suite (see [COVERAGE.md](COVERAGE.md)); the tests, not a manual
read-through, are the primary correctness guarantee today. We're being upfront
about this so you can weigh it for your use case. Human review is ongoing, and —
as above — bug reports are the fastest way to harden the library.

## Installation

```sh
go get github.com/streame-gg/go-discord-wrapper
```

Requires Go 1.26 or later.

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
    "github.com/streame-gg/go-discord-wrapper/types/discord"
    "github.com/streame-gg/go-discord-wrapper/types/events"
    "github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

func main() {
    bot, err := connection.NewClient(os.Getenv("DISCORD_TOKEN"), discord.IntentGuilds|discord.IntentGuildMessages)
    if err != nil {
        slog.Error("failed to create client", "err", err)
        os.Exit(1)
    }

    bot.OnMessageCreate(func(c *connection.Client, e *events.MessageCreateEvent) {
        if e.Author == nil {
            return
        }
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

## Working with Collections

Most cache and listing operations return a `*collection.Collection[K, V]` —
a map with insertion order, plus 30+ utility methods inspired by
[@discordjs/collection](https://github.com/discordjs/discord.js/tree/main/packages/collection).

### Basic operations

```go
import "github.com/streame-gg/go-discord-wrapper/collection"

users := bot.Cache.Users().All()  // returns *Collection[Snowflake, *User]

count := users.Len()
admin, ok := users.Get(adminID)
humanCount := users.Filter(func(u *discord.User) bool { return !u.Bot }).Len()
```

### Search and partition

```go
bots, humans := users.Partition(func(u *discord.User) bool { return u.Bot })
recent := messages.Filter(func(m *discord.Message) bool {
    return time.Since(m.Created()) < time.Hour
})

oldest, _ := messages.Sorted(func(a, b *discord.Message) bool {
    return a.Created().Before(b.Created())
}).First()
```

### Transform

```go
// Method-based: returns a new Collection
sorted := users.Sorted(func(a, b *discord.User) bool {
    return a.Username < b.Username
})

// Top-level: transforms to a different value-type
names := collection.Map(users, func(u *discord.User) string {
    return u.Username
})
// names is []string

// Iterate (Go 1.23+ range-over-func)
for id, user := range users.All() {
    fmt.Printf("%s -> %s\n", id, user.Username)
}
```

Full API: see [godoc](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper/collection).

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
| Redis | `cache/rediscache` package |
| MongoDB | `cache/mongocache` package |

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
        Type:        discord.ApplicationCommandTypeChatInput,
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
| `options.WithBaseURL(url)` | Override the Discord API base URL (useful for mock servers in tests) |

## Testing with a mock server

Use `WithBaseURL` together with Go's `net/http/httptest` to run integration tests without
hitting Discord:

```go
import (
    "net/http/httptest"
    "github.com/streame-gg/go-discord-wrapper/api"
    "github.com/streame-gg/go-discord-wrapper/options"
)

ts := httptest.NewServer(myHandler)
defer ts.Close()

client, err := api.NewRestClient("test-token",
    options.WithBaseURL(ts.URL),
)
```

## Synthetic events

Synthetic events are high-level events derived by the wrapper from raw Discord gateway events.
They are never sent by Discord directly — the library computes them by diffing state across
consecutive gateway payloads.

Synthetic events require the **cache to be enabled** (via `options.WithCache`). When the cache
is cold (no prior state), the event is silently skipped rather than firing with incomplete data.

```go
mc := cache.NewMemoryCache(cache.Options{})
bot, err := connection.NewClient(token, intents, options.WithCache(mc))

// Fires once per emoji added to a guild.
bot.OnGuildEmojiAdd(func(c *connection.Client, ev *events.GuildEmojiAddEvent) {
    fmt.Printf("emoji %s added to guild %s\n", ev.Emoji.Name, ev.GuildID)
})

// Fires once per sticker removed from a guild.
bot.OnGuildStickerRemove(func(c *connection.Client, ev *events.GuildStickerRemoveEvent) {
    fmt.Printf("sticker %s removed from guild %s\n", ev.Sticker.Name, ev.GuildID)
})
```

### Emoji events (source: `GUILD_EMOJIS_UPDATE`, cache required)

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnGuildEmojiAdd` | `GuildEmojiAddEvent` | An emoji is created in a guild |
| `OnGuildEmojiRemove` | `GuildEmojiRemoveEvent` | An emoji is deleted from a guild |
| `OnGuildEmojiUpdate` | `GuildEmojiUpdateEvent` | An emoji's name, roles, or availability changes |

### Sticker events (source: `GUILD_STICKERS_UPDATE`, cache required)

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnGuildStickerAdd` | `GuildStickerAddEvent` | A sticker is created in a guild |
| `OnGuildStickerRemove` | `GuildStickerRemoveEvent` | A sticker is deleted from a guild |
| `OnGuildStickerUpdate` | `GuildStickerUpdateEvent` | A sticker's name, tags, or availability changes |

## Supported gateway events

| Event | Constant |
|-------|----------|
| Ready | `events.EventReady` |
| Resumed | `events.EventResumed` |
| Application Command Permissions Update | `events.EventApplicationCommandPermissionsUpdate` |
| Auto Moderation Rule Create | `events.EventAutoModerationRuleCreate` |
| Auto Moderation Rule Update | `events.EventAutoModerationRuleUpdate` |
| Auto Moderation Rule Delete | `events.EventAutoModerationRuleDelete` |
| Auto Moderation Action Execution | `events.EventAutoModerationActionExecution` |
| Channel Create | `events.EventChannelCreate` |
| Channel Update | `events.EventChannelUpdate` |
| Channel Delete | `events.EventChannelDelete` |
| Channel Pins Update | `events.EventChannelPinsUpdate` |
| Thread Create | `events.EventThreadCreate` |
| Thread Update | `events.EventThreadUpdate` |
| Thread Delete | `events.EventThreadDelete` |
| Thread List Sync | `events.EventThreadListSync` |
| Thread Member Update | `events.EventThreadMemberUpdate` |
| Thread Members Update | `events.EventThreadMembersUpdate` |
| Entitlement Create | `events.EventEntitlementCreate` |
| Entitlement Update | `events.EventEntitlementUpdate` |
| Entitlement Delete | `events.EventEntitlementDelete` |
| Guild Create | `events.EventGuildCreate` |
| Guild Update | `events.EventGuildUpdate` |
| Guild Delete | `events.EventGuildDelete` |
| Guild Audit Log Entry Create | `events.EventGuildAuditLogEntryCreate` |
| Guild Ban Add | `events.EventGuildBanAdd` |
| Guild Ban Remove | `events.EventGuildBanRemove` |
| Guild Emojis Update | `events.EventGuildEmojisUpdate` |
| Guild Stickers Update | `events.EventGuildStickersUpdate` |
| Guild Integrations Update | `events.EventGuildIntegrationsUpdate` |
| Guild Member Add | `events.EventGuildMemberAdd` |
| Guild Member Update | `events.EventGuildMemberUpdate` |
| Guild Member Remove | `events.EventGuildMemberRemove` |
| Guild Members Chunk | `events.EventGuildMembersChunk` |
| Guild Role Create | `events.EventGuildRoleCreate` |
| Guild Role Update | `events.EventGuildRoleUpdate` |
| Guild Role Delete | `events.EventGuildRoleDelete` |
| Guild Scheduled Event Create | `events.EventGuildScheduledEventCreate` |
| Guild Scheduled Event Update | `events.EventGuildScheduledEventUpdate` |
| Guild Scheduled Event Delete | `events.EventGuildScheduledEventDelete` |
| Guild Scheduled Event User Add | `events.EventGuildScheduledEventUserAdd` |
| Guild Scheduled Event User Remove | `events.EventGuildScheduledEventUserRemove` |
| Guild Soundboard Sound Create | `events.EventGuildSoundboardSoundCreate` |
| Guild Soundboard Sound Update | `events.EventGuildSoundboardSoundUpdate` |
| Guild Soundboard Sound Delete | `events.EventGuildSoundboardSoundDelete` |
| Guild Soundboard Sounds Update | `events.EventGuildSoundboardSoundsUpdate` |
| Soundboard Sounds | `events.EventSoundboardSounds` |
| Integration Create | `events.EventIntegrationCreate` |
| Integration Update | `events.EventIntegrationUpdate` |
| Integration Delete | `events.EventIntegrationDelete` |
| Interaction Create | `events.EventInteractionCreate` |
| Invite Create | `events.EventInviteCreate` |
| Invite Delete | `events.EventInviteDelete` |
| Message Create | `events.EventMessageCreate` |
| Message Update | `events.EventMessageUpdate` |
| Message Delete | `events.EventMessageDelete` |
| Message Delete Bulk | `events.EventMessageDeleteBulk` |
| Message Reaction Add | `events.EventMessageReactionAdd` |
| Message Reaction Remove | `events.EventMessageReactionRemove` |
| Message Reaction Remove All | `events.EventMessageReactionRemoveAll` |
| Message Reaction Remove Emoji | `events.EventMessageReactionRemoveEmoji` |
| Message Poll Vote Add | `events.EventMessagePollVoteAdd` |
| Message Poll Vote Remove | `events.EventMessagePollVoteRemove` |
| Presence Update | `events.EventPresenceUpdate` |
| Stage Instance Create | `events.EventStageInstanceCreate` |
| Stage Instance Update | `events.EventStageInstanceUpdate` |
| Stage Instance Delete | `events.EventStageInstanceDelete` |
| Subscription Create | `events.EventSubscriptionCreate` |
| Subscription Update | `events.EventSubscriptionUpdate` |
| Subscription Delete | `events.EventSubscriptionDelete` |
| Typing Start | `events.EventTypingStart` |
| User Update | `events.EventUserUpdate` |
| Voice State Update | `events.EventVoiceStateUpdate` |
| Voice Server Update | `events.EventVoiceServerUpdate` |
| Voice Channel Status Update | `events.EventVoiceChannelStatusUpdate` |
| Voice Channel Start Time Update | `events.EventVoiceChannelStartTimeUpdate` |
| Voice Channel Effect Send | `events.EventVoiceChannelEffectSend` |
| Webhooks Update | `events.EventWebhooksUpdate` |

> **Voice channel status ≠ channel update.** When a voice channel's *status*
> text changes, Discord sends a dedicated `VOICE_CHANNEL_STATUS_UPDATE` — **not**
> `CHANNEL_UPDATE`. Subscribe with `OnVoiceChannelStatusUpdate`, not
> `OnChannelUpdate`. (When the cache is enabled, the wrapper also patches the
> cached channel's `Status` field for you.)

## Documentation

In-depth guides live in [`docs/`](docs/README.md):

| Guide | Covers |
|-------|--------|
| [Getting started](docs/GETTING_STARTED.md) | From empty folder to a running bot |
| [Configuration](docs/CONFIGURATION.md) | Intents and every `options.With…` option |
| [Events & handlers](docs/EVENTS.md) | Handlers, middleware, concurrency, synthetic events |
| [Slash commands](docs/COMMANDS.md) | Registering commands, reading options, replying |
| [Command management](docs/COMMAND_MANAGEMENT.md) | Create, edit, delete, scope commands; permissions |
| [Components](docs/COMPONENTS.md) | Buttons, select menus, Components V2 |
| [Modals](docs/MODALS.md) | Modal forms: build, show, read |
| [Messages](docs/MESSAGES.md) | Sending, editing, attachments, mentions |
| [Embeds](docs/EMBEDS.md) | Building and validating embeds |
| [Cache behavior](docs/CACHE.md) | What is cached, when, and incomplete member caches |
| [Sharding](docs/SHARDING.md) | `ShardManager`, local coordinator, cross-shard messaging |
| [REST client](docs/REST.md) | Standalone REST usage and typed errors |

New to the library? Start with [Getting started](docs/GETTING_STARTED.md), then
copy [`example/template`](example/template) as your project skeleton. Other
runnable programs under [`example/`](example): `caching`, `commands`,
`sharding`, and `slash_with_defer`.

## API reference

Full API documentation is available at [pkg.go.dev/github.com/streame-gg/go-discord-wrapper](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

Join the Discord server: https://discord.gg/nfBuZVejqp

## License

Apache 2.0 — see [LICENSE](LICENSE).
