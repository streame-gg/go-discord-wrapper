# go-discord-wrapper

A Go library for the Discord gateway and REST API.

> **Status:** Alpha — actively developed, APIs may change before the first tagged release.

## Features

- Gateway (WebSocket) client with automatic reconnection and resume
- Full REST API coverage for guilds, channels, messages, roles, interactions, webhooks, and more
- Proactive rate limiter — respects per-route and global Discord rate limits before sending
- In-process memory cache (guild, channel, member, message, role, voice state, emoji, sticker, …) with TTL and LRU eviction
- Pluggable cache backends (Redis and MongoDB drivers included)
- **21 high-level synthetic events** derived from raw gateway events (voice transitions, role diffs, nick changes, boost detection, emoji/sticker diffs, …)
- Slash command and interaction support (reply, defer, modal, component)
- Horizontal sharding with a built-in shard coordinator
- Middleware support for event handlers

> **Voice support:** Voice state and server events are tracked, but voice-connection
> establishment (UDP, Opus encoding/decoding, sending or receiving audio) is **not supported**.
> If you need voice capabilities, consider using [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
> alongside this library for voice operations only.

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

Most synthetic events require the **cache to be enabled** (via `options.WithCache`).
When the cache is cold (no prior state), the event is silently skipped rather than firing
with incomplete data.

```go
mc := cache.NewMemoryCache(cache.Options{})
bot, err := connection.NewClient(token, intents, options.WithCache(mc))

// Fires once per role added — never need to diff role slices yourself.
bot.OnGuildMemberRoleAdd(func(c *connection.Client, ev *events.GuildMemberRoleAddEvent) {
    fmt.Printf("user %s gained role %s\n", ev.UserID, ev.RoleID)
})

// Fires when a user joins a voice channel for the first time.
bot.OnVoiceMemberJoin(func(c *connection.Client, ev *events.VoiceMemberJoinEvent) {
    fmt.Printf("user %s joined channel %s\n", ev.UserID, ev.ChannelID)
})
```

Raw events are dispatched/enqueued before any derived synthetic events. If dispatch is
configured to run serially (single-threaded), handlers observe that same order; in
concurrent modes, synthetic handlers may run before the corresponding raw handler has
finished.

### Voice events (source: `VOICE_STATE_UPDATE`)

No cache required — voice state is tracked internally.

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnVoiceMemberJoin` | `VoiceMemberJoinEvent` | User enters a voice channel (was not in any channel before) |
| `OnVoiceMemberLeave` | `VoiceMemberLeaveEvent` | User leaves all voice channels |
| `OnVoiceMemberMove` | `VoiceMemberMoveEvent` | User switches from one channel to another |
| `OnVoiceMemberUpdate` | `VoiceMemberUpdateEvent` | User's state changes while staying in the same channel (mute, deafen, …) |

### Member events (source: `GUILD_MEMBER_UPDATE`, cache required)

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnGuildMemberRoleAdd` | `GuildMemberRoleAddEvent` | A role is added to a member (one event per role) |
| `OnGuildMemberRoleRemove` | `GuildMemberRoleRemoveEvent` | A role is removed from a member (one event per role) |
| `OnGuildMemberNickChange` | `GuildMemberNickChangeEvent` | A member's nickname is set, changed, or cleared |
| `OnGuildMemberTimeout` | `GuildMemberTimeoutEvent` | A timeout is applied or extended |
| `OnGuildMemberBoostStart` | `GuildMemberBoostStartEvent` | A member starts boosting the server |
| `OnGuildMemberBoostEnd` | `GuildMemberBoostEndEvent` | A member stops boosting the server |

### Presence events (source: `PRESENCE_UPDATE`)

Status transitions require cache. `OnUserProfileUpdate` fires on any non-ID field in the
partial user payload (no cache needed).

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnUserOnline` | `UserOnlineEvent` | User transitions from offline to any active status |
| `OnUserOffline` | `UserOfflineEvent` | User transitions to offline |
| `OnUserActivityChange` | `UserActivityChangeEvent` | User's activity set changes |
| `OnUserProfileUpdate` | `UserProfileUpdateEvent` | User's profile fields change (username, avatar, …) |

### Emoji events (source: `GUILD_EMOJIS_UPDATE`, cache required)

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnGuildEmojiAdd` | `GuildEmojiAddEvent` | An emoji is created in a guild |
| `OnGuildEmojiRemove` | `GuildEmojiRemoveEvent` | An emoji is deleted from a guild |
| `OnGuildEmojiUpdate` | `GuildEmojiUpdateEvent` | An emoji's name changes |

### Sticker events (source: `GUILD_STICKERS_UPDATE`, cache required)

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnGuildStickerAdd` | `GuildStickerAddEvent` | A sticker is created in a guild |
| `OnGuildStickerRemove` | `GuildStickerRemoveEvent` | A sticker is deleted from a guild |
| `OnGuildStickerUpdate` | `GuildStickerUpdateEvent` | A sticker's name changes |

### Role events (source: `GUILD_ROLE_UPDATE`, cache required)

| Helper | Struct | Fires when |
|--------|--------|------------|
| `OnGuildRolePermissionsChange` | `GuildRolePermissionsChangeEvent` | A role's permission bitfield changes |

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
