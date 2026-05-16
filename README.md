# go-discord-wrapper

A Go library for the Discord gateway and REST API.

> **Status:** This library is very new and being checked for bugs. I am currently updating a lot so please do not use yet. :)

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
| Guild Integrations Update | `events.EventIntegrationUpdate` |
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

## API reference

Full API documentation is available at [pkg.go.dev/github.com/streame-gg/go-discord-wrapper](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

Join the Discord server: https://discord.gg/nfBuZVejqp

## License

Apache 2.0 — see [LICENSE](LICENSE).
