# go-discord-wrapper

A Go library for the Discord gateway and REST API.

> **Status:** The library is in active development toward a stable v0.1.0 release. Public APIs may still change. Bug reports and pull requests are welcome.

## Installation

```bash
go get github.com/streame-gg/go-discord-wrapper
```

Requires Go 1.21 or later.

## Quickstart

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/streame-gg/go-discord-wrapper/connection"
    "github.com/streame-gg/go-discord-wrapper/options"
    "github.com/streame-gg/go-discord-wrapper/types/common"
    "github.com/streame-gg/go-discord-wrapper/types/events"
    "github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

func main() {
    bot := connection.NewClient(
        os.Getenv("DISCORD_TOKEN"),
        common.IntentGuilds|common.IntentGuildMessages,
    )

    bot.OnEvent(events.EventMessageCreate, func(client *connection.Client, e *events.MessageCreateEvent) {
        slog.Info("message", "author", e.Author.Username, "content", e.Content)
    })

    bot.OnEvent(events.EventInteractionCreate, func(client *connection.Client, e *events.InteractionCreateEvent) {
        if e.IsCommand() {
            _, err := client.Reply(&e.Interaction, &responses.InteractionResponseDataDefault{
                Content: "Hello, " + e.Member.User.DisplayName() + "!",
            }, false)
            if err != nil {
                slog.Error("reply failed", "err", err)
            }
        }
    })

    if err := bot.Login(); err != nil {
        slog.Error("login failed", "err", err)
        os.Exit(1)
    }

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    <-ctx.Done()

    bot.Shutdown()
}
```

See [`example/main.go`](example/main.go) for a more complete example covering slash commands, buttons, modals, and cache usage.

## Supported Gateway Events

All standard Discord gateway events are supported. Register a handler with `client.OnEvent`:

| Constant | Discord Event |
|---|---|
| `events.EventReady` | `READY` |
| `events.EventGuildCreate` | `GUILD_CREATE` |
| `events.EventGuildUpdate` | `GUILD_UPDATE` |
| `events.EventGuildDelete` | `GUILD_DELETE` |
| `events.EventGuildBanAdd` | `GUILD_BAN_ADD` |
| `events.EventGuildBanRemove` | `GUILD_BAN_REMOVE` |
| `events.EventGuildEmojisUpdate` | `GUILD_EMOJIS_UPDATE` |
| `events.EventGuildStickersUpdate` | `GUILD_STICKERS_UPDATE` |
| `events.EventGuildIntegrationsUpdate` | `GUILD_INTEGRATIONS_UPDATE` |
| `events.EventGuildMemberAdd` | `GUILD_MEMBER_ADD` |
| `events.EventGuildMemberRemove` | `GUILD_MEMBER_REMOVE` |
| `events.EventGuildMemberUpdate` | `GUILD_MEMBER_UPDATE` |
| `events.EventGuildRoleCreate` | `GUILD_ROLE_CREATE` |
| `events.EventGuildRoleUpdate` | `GUILD_ROLE_UPDATE` |
| `events.EventGuildRoleDelete` | `GUILD_ROLE_DELETE` |
| `events.EventGuildScheduledEventCreate` | `GUILD_SCHEDULED_EVENT_CREATE` |
| `events.EventGuildScheduledEventUpdate` | `GUILD_SCHEDULED_EVENT_UPDATE` |
| `events.EventGuildScheduledEventDelete` | `GUILD_SCHEDULED_EVENT_DELETE` |
| `events.EventGuildScheduledEventUserAdd` | `GUILD_SCHEDULED_EVENT_USER_ADD` |
| `events.EventGuildScheduledEventUserRemove` | `GUILD_SCHEDULED_EVENT_USER_REMOVE` |
| `events.EventGuildAuditLogEntryCreate` | `GUILD_AUDIT_LOG_ENTRY_CREATE` |
| `events.EventChannelCreate` | `CHANNEL_CREATE` |
| `events.EventChannelUpdate` | `CHANNEL_UPDATE` |
| `events.EventChannelDelete` | `CHANNEL_DELETE` |
| `events.EventChannelPinsUpdate` | `CHANNEL_PINS_UPDATE` |
| `events.EventThreadCreate` | `THREAD_CREATE` |
| `events.EventThreadUpdate` | `THREAD_UPDATE` |
| `events.EventThreadDelete` | `THREAD_DELETE` |
| `events.EventThreadListSync` | `THREAD_LIST_SYNC` |
| `events.EventMessageCreate` | `MESSAGE_CREATE` |
| `events.EventMessageUpdate` | `MESSAGE_UPDATE` |
| `events.EventMessageDelete` | `MESSAGE_DELETE` |
| `events.EventMessageDeleteBulk` | `MESSAGE_DELETE_BULK` |
| `events.EventMessageReactionAdd` | `MESSAGE_REACTION_ADD` |
| `events.EventMessageReactionRemove` | `MESSAGE_REACTION_REMOVE` |
| `events.EventMessageReactionRemoveAll` | `MESSAGE_REACTION_REMOVE_ALL` |
| `events.EventMessageReactionRemoveEmoji` | `MESSAGE_REACTION_REMOVE_EMOJI` |
| `events.EventMessagePollVoteAdd` | `MESSAGE_POLL_VOTE_ADD` |
| `events.EventMessagePollVoteRemove` | `MESSAGE_POLL_VOTE_REMOVE` |
| `events.EventInteractionCreate` | `INTERACTION_CREATE` |
| `events.EventInviteCreate` | `INVITE_CREATE` |
| `events.EventInviteDelete` | `INVITE_DELETE` |
| `events.EventPresenceUpdate` | `PRESENCE_UPDATE` |
| `events.EventTypingStart` | `TYPING_START` |
| `events.EventUserUpdate` | `USER_UPDATE` |
| `events.EventVoiceStateUpdate` | `VOICE_STATE_UPDATE` |
| `events.EventStageInstanceCreate` | `STAGE_INSTANCE_CREATE` |
| `events.EventStageInstanceUpdate` | `STAGE_INSTANCE_UPDATE` |
| `events.EventStageInstanceDelete` | `STAGE_INSTANCE_DELETE` |
| `events.EventIntegrationCreate` | `INTEGRATION_CREATE` |
| `events.EventIntegrationUpdate` | `INTEGRATION_UPDATE` |
| `events.EventIntegrationDelete` | `INTEGRATION_DELETE` |
| `events.EventWebhooksUpdate` | `WEBHOOKS_UPDATE` |
| `events.EventAutoModerationRuleCreate` | `AUTO_MODERATION_RULE_CREATE` |
| `events.EventAutoModerationRuleUpdate` | `AUTO_MODERATION_RULE_UPDATE` |
| `events.EventAutoModerationRuleDelete` | `AUTO_MODERATION_RULE_DELETE` |
| `events.EventAutoModerationActionExecution` | `AUTO_MODERATION_ACTION_EXECUTION` |
| `events.EventEntitlementCreate` | `ENTITLEMENT_CREATE` |
| `events.EventEntitlementUpdate` | `ENTITLEMENT_UPDATE` |
| `events.EventEntitlementDelete` | `ENTITLEMENT_DELETE` |

The following events are automatically handled internally and also dispatched to your handlers:

- `GUILD_CREATE / DELETE / UPDATE` — populate the guild, member count, and unavailable-guild caches
- `CHANNEL_CREATE / UPDATE / DELETE` — populate the channel cache
- `GUILD_MEMBER_ADD / REMOVE / UPDATE` — populate the member cache and member counts
- `MESSAGE_CREATE / UPDATE / DELETE / DELETE_BULK` — populate the message cache
- `USER_UPDATE` — updates the user cache

## Cache

Pass a cache to `options.WithCache` and the gateway populates it automatically from events. Three backends are provided:

```go
import "github.com/streame-gg/go-discord-wrapper/cache"

// In-process (zero dependencies)
c := cache.NewMemoryCache(cache.Options{
    TTL:           30 * time.Minute,
    SweepInterval: 5 * time.Minute,
    Limits: cache.Limits{
        MaxMessages: 10_000,
    },
    Messages: cache.MessageOptions{
        MaxPerChannel: 200,
    },
})
defer c.Close()

bot := connection.NewClient(token, intents, options.WithCache(c))
```

```go
// Redis (requires go-redis)
import "github.com/streame-gg/go-discord-wrapper/cache/rediscache"

rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
c := rediscache.NewRedisCache(rdb, cache.Options{TTL: 30 * time.Minute})
```

```go
// MongoDB
import "github.com/streame-gg/go-discord-wrapper/cache/mongocache"

c := mongocache.NewMongoDBCache(db, cache.Options{TTL: 30 * time.Minute})
```

Read from the cache directly on your client:

```go
// Prefer cache; fall back to REST on a miss
if ch, ok := bot.Cache.Channels().Get(channelID); ok {
    fmt.Println("from cache:", ch.Name)
} else {
    ch, _ := bot.GetChannel(channelID)
    fmt.Println("from REST:", ch.Name)
}
```

## Sharding

Use `options.WithSharding` for single-process sharding and `sharding.NewShardManager` to orchestrate multiple shards:

```go
import (
    "github.com/streame-gg/go-discord-wrapper/connection"
    "github.com/streame-gg/go-discord-wrapper/options"
    "github.com/streame-gg/go-discord-wrapper/sharding"
)

totalShards := 4
coord := sharding.NewLocalCoordinator(totalShards)

manager := sharding.NewShardManager(coord, func(shardID int) *connection.Client {
    return connection.NewClient(
        token,
        intents,
        options.WithSharding(totalShards, shardID),
        options.WithCoordinator(coord),
    )
})

if err := manager.Start(); err != nil {
    log.Fatal(err)
}
```

Cross-shard aggregate statistics:

```go
ctx := context.Background()
totalServers, _ := sharding.GetTotalServers(ctx, manager)
totalMembers, _ := sharding.GetTotalUsers(ctx, manager)
```

## REST API

All REST methods are on `*connection.Client` (which embeds `*api.RestClient`) or directly on `api.NewRestClient`. The full list is available at [pkg.go.dev/github.com/streame-gg/go-discord-wrapper](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper).

Covered resource groups: channels, guilds, members, roles, bans, messages, reactions, pins, invites, webhooks, threads, emojis, stickers, auto moderation, stage instances, scheduled events, interactions, application commands, users.

## REST-only usage (no gateway)

```go
import "github.com/streame-gg/go-discord-wrapper/api"

client := api.NewRestClient(os.Getenv("DISCORD_TOKEN"))
msg, err := client.CreateMessage(context.Background(), channelID, api.CreateMessageParams{
    Content: "Hello!",
})
```

## Error Handling

All REST methods return `*api.APIError` on failure. Use the sentinel errors for common HTTP status checks:

```go
import "github.com/streame-gg/go-discord-wrapper/api"

msg, err := client.CreateMessage(ctx, channelID, params)
if err != nil {
    if errors.Is(err, api.ErrNotFound) {
        // channel was deleted
    } else if errors.Is(err, api.ErrForbidden) {
        // bot lacks permission
    } else if errors.Is(err, api.ErrRateLimited) {
        // rate limited (normally retried automatically)
    }
}
```

To access the full Discord error detail (JSON error code, message, field errors):

```go
var apiErr *api.APIError
if errors.As(err, &apiErr) {
    log.Printf("discord code %d: %s (http %d)", apiErr.Code, apiErr.Message, apiErr.HTTPStatus)
}
```

## Context and Timeouts

Every REST method accepts a `context.Context` as its first argument. Pass a context with a deadline to enforce per-request timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

guild, err := client.GetGuild(ctx, guildID, false)
```

Cancel a context to abort an in-flight request or any pending retry/rate-limit sleep:

```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(2 * time.Second)
    cancel() // abort the request
}()

msg, err := client.CreateMessage(ctx, channelID, params)
if errors.Is(err, context.Canceled) {
    // request was cancelled
}
```

The `connection.Client` wrapper methods (e.g. `bot.GetChannel`, `bot.CreateMessage`) pass `context.Background()` internally. Use `api.NewRestClient` directly when you need per-call context control.

## Documentation

Full API reference: [pkg.go.dev/github.com/streame-gg/go-discord-wrapper](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper)

## Troubleshooting

**Bot connects but receives no events**
Check that you have the correct intents. `IntentGuildMessages` is needed for `MESSAGE_CREATE`; `IntentGuildMembers` and `IntentPresences` are privileged and must be enabled in the Discord Developer Portal.

**4014 Disallowed Intents**
A privileged intent (`GuildMembers`, `GuildPresences`, `MessageContent`) was requested but not enabled in the Developer Portal. Go to your application's "Bot" page and toggle the intent on.

**WebSocket close code 4004 / Authentication failed**
The bot token is invalid or was regenerated. Ensure `DISCORD_TOKEN` is set to the current token from the Developer Portal.

**WebSocket close code 4011 / Sharding required**
Your bot is in too many servers for a single shard. Use `options.WithSharding` and increase `TotalShards`.

**REST returns 429 even with rate limiting enabled**
Global rate limits apply across all bots sharing an IP. Consider `options.WithRateLimiting(options.RateLimiterOptions{SafetyMargin: 2})` for extra headroom when running many concurrent goroutines.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

Join the Discord server: https://discord.gg/nfBuZVejqp

## License

Apache 2.0 — see [LICENSE](LICENSE).
