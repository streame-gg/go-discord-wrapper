# Events & handlers

Everything your bot reacts to — messages, member joins, interactions, voice
state changes — arrives as a **gateway event**. This page covers how to
register handlers, the middleware model, concurrency, and the *synthetic*
events the library derives for you.

## Registering handlers

There are two styles. Prefer the typed helpers; fall back to `OnEvent` for the
handful of events without one.

### Typed helpers (recommended)

Each common event has an `On…` method that gives you a correctly typed event
struct — no casting:

```go
bot.OnMessageCreate(func(c *connection.Client, e *events.MessageCreateEvent) {
    if e.Author == nil || e.Author.Bot {
        return
    }
    c.Logger.Info("message", "author", e.Author.Username, "content", e.Content)
})

bot.OnGuildMemberAdd(func(c *connection.Client, e *events.GuildMemberAddEvent) {
    c.Logger.Info("member joined", "guild", e.GuildID.String())
})
```

There are typed helpers for the gateway lifecycle (`OnConnect`, `OnDisconnect`,
`OnReconnect`, `OnReady`, `OnResumed`, `OnPacketError`) and for the vast
majority of Discord events. See the
[full event list](../README.md#supported-gateway-events) in the README.

### Generic `OnEvent`

For any event, register against its `events.Event…` constant. The handler is a
typed function and is validated at registration time:

```go
err := bot.OnEvent(events.EventInteractionCreate,
    func(c *connection.Client, e *events.InteractionCreateEvent) {
        // …
    })
```

`OnEvent` returns an error if the handler signature doesn't match the event —
check it.

## Middleware

`Use` wraps **every handler registered after it** with cross-cutting behaviour:
logging, metrics, auth checks, panic recovery. A middleware receives the next
handler and must call it to continue the chain — skipping the call
short-circuits the event.

```go
bot.Use(func(next connection.EventHandler) connection.EventHandler {
    return func(c *connection.Client, ev events.Event) {
        start := time.Now()
        next(c, ev)
        c.Logger.Debug("handled", "type", ev.Event(), "took", time.Since(start))
    }
})
```

Middleware run **in registration order**: the first registered is the outermost
wrapper. Register middleware before the handlers they should apply to.

## Concurrency & dropped events

Events are dispatched concurrently across a worker pool. The pool size is capped
by `options.WithMaxConcurrentEvents(n)` (default 64). When the pool is
saturated, **further events are dropped with a warning log** rather than queued
unbounded — so handlers should be fast and offload slow work to their own
goroutines.

Increase the cap for bots that take dense bursts (e.g. large-guild
`GUILD_MEMBERS_CHUNK` storms):

```go
connection.NewClient(token, intents, options.WithMaxConcurrentEvents(256))
```

Within a single handler, treat the event as read-only shared state; if you fan
work out to goroutines, copy what you need first.

## Synthetic events

Some high-level events are **computed by the library**, not sent by Discord.
Discord delivers a coarse `GUILD_EMOJIS_UPDATE` containing the *entire* new
emoji list; the wrapper diffs it against the cached list and fires a precise
add/remove/update event per change.

Because they are derived from cache state, **synthetic events require a cache**
(`options.WithCache`). When the cache is cold (no prior list to diff against),
the event is skipped rather than fired with incomplete data.

| Helper | Fires when |
|--------|------------|
| `OnGuildEmojiAdd` / `Remove` / `Update` | an emoji is added / removed / changed |
| `OnGuildStickerAdd` / `Remove` / `Update` | a sticker is added / removed / changed |

```go
bot.OnGuildEmojiAdd(func(c *connection.Client, ev *events.GuildEmojiAddEvent) {
    fmt.Printf("emoji %s added to guild %s\n", ev.Emoji.Name, ev.GuildID)
})
```

## Events and the cache

The same events that you handle also feed the cache (when one is attached). For
example, `GUILD_BAN_ADD` populates `bot.Cache.Bans()` and `MESSAGE_CREATE`
appends to the per-channel message ring. You can narrow which stores are written
with `options.WithCacheStores`, independently of which events you handle. See
[docs/CACHE.md](CACHE.md) for the store-by-store breakdown.

## See also

- [docs/COMMANDS.md](COMMANDS.md) — interaction events and replying.
- [docs/CONFIGURATION.md](CONFIGURATION.md) — `WithMaxConcurrentEvents`,
  reconnect tuning, and intents.
- [docs/SHARDING.md](SHARDING.md) — how events are split across shards.
