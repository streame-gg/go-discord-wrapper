# Cache Behavior

## What gets cached on connect/reconnect

When the bot connects (or reconnects via Identify), Discord sends one
`GUILD_CREATE` event per guild the bot belongs to. Each event contains:

- Full guild metadata (name, icon, settings)
- All roles
- All emojis and stickers
- All channels and active threads
- Voice states of members currently connected to voice channels
- Presences (requires `GUILD_PRESENCES` intent)
- **A subset of members** — typically only:
  - Members currently in voice channels
  - The bot itself
  - Up to `large_threshold` members (default 250, max 250)

For guilds with more than ~250 members, **the member cache will be
incomplete** after `GUILD_CREATE`. This is a Discord protocol limitation,
not a library bug.

## Working with incomplete member caches

Options to populate members:

1. **Lazy via gateway events** — as members interact (send messages, join
   voice, get role changes), they are added to the cache. Eventually
   consistent.

2. **Explicit fetch** — `guild.Members().Fetch(ctx, userID)` does a REST
   call for a single member and adds the result to the cache.

3. **Bulk fetch** — `guild.Members().FetchAll(ctx, opts)` requests all
   members via REST. Subject to rate limits; can be slow for large guilds.

4. **Gateway request** — `client.RequestGuildMembers(guildID, params)` sends
   opcode 8. Members arrive via `GUILD_MEMBERS_CHUNK` events and populate
   the cache asynchronously. Requires the `GUILD_MEMBERS` intent. This call
   does not take a `context.Context` — it only enqueues the request on the
   gateway socket; watch for the chunks via `OnGuildMembersChunk`.

## Avoiding nil panics on manager methods

Every entity fetched from the cache has all sub-managers pre-populated by
the library. For example, a `*discord.Guild` returned by
`cache.Guilds().Get(id)` or passed to `OnGuildCreate` always has working
`Members()`, `Roles()`, `Channels()`, etc. managers — no extra setup needed.

This is true even when no cache is configured: the client-level
`bot.Guilds()`, `bot.Users()`, and `bot.Channels()` managers always return
a non-nil manager. Cache-backed methods (`Get`, `Cache`, `Size`) return
sensible zero values; REST-backed methods (`Fetch`, `Create`) work normally.

## Behavior after reconnect

- **Resume** (short downtime, <~3 min): cache is preserved; missed events
  are replayed via the resume sequence.
- **Identify** (longer downtime, session expired): Discord sends a fresh
  `GUILD_CREATE` storm. Per-guild caches (members, roles, channels, etc.)
  are replaced with current state. **Events that occurred between disconnect
  and reconnect are not delivered** — the cache silently jumps to the
  current Discord state.

This matches discord.js behavior. Bots that need audit-grade event delivery
should not treat the cache as a source of truth for what happened during an
outage.

## What gets cached from incremental gateway events

Beyond the `GUILD_CREATE` snapshot, the cache is kept current from the
event stream. Each of these stores is populated as the matching gateway
event arrives, and is reachable from the top-level cache:

| Store | Accessor | Populated by |
|-------|----------|--------------|
| Bans | `bot.Cache.Bans()` | `GUILD_BAN_ADD` / `GUILD_BAN_REMOVE` |
| Invites | `bot.Cache.Invites()` | `INVITE_CREATE` / `INVITE_DELETE` |
| Auto-mod rules | `bot.Cache.AutoModRules()` | `AUTO_MODERATION_RULE_CREATE` / `_UPDATE` / `_DELETE` |

> **Changed in a recent release.** Bans, invites, and auto-mod rules used to
> be REST-only. They are now cached from their gateway events. The guild-level
> `guild.Bans()`, `guild.Invites()`, and `guild.AutoModRules()` managers still
> expose REST operations (`Fetch`, `FetchAll`, `Create`, …); the cached copies
> live on the top-level `bot.Cache`.

## What is NOT cached

The library does not cache the following — there is no gateway event stream
for them, so they are REST-only:

| Entity | How to fetch |
|--------|--------------|
| Webhooks | `guild.Webhooks().FetchAll(ctx)` |
| Integrations | `guild.Integrations().FetchAll(ctx)` |

These managers exist on every guild but have no `Cache()` or `Get()`
methods — they expose only REST operations.

## Selecting which stores are populated

By default every store is auto-populated. To trim memory or focus on the
entities you actually read, narrow the set with `options.WithCacheStores`
(a bitmask of `cache.Category*` values), or subtract individual stores with
`options.WithDisableCacheStore`:

```go
// Only guilds, channels, and members.
options.WithCacheStores(cache.CategoryGuilds | cache.CategoryChannels | cache.CategoryMembers)

// Everything except presences and messages.
options.WithDisableCacheStore(cache.CategoryPresences | cache.CategoryMessages)
```

To turn off auto-population entirely and manage the cache yourself, use
`options.WithDisableCacheAutoPopulation()`.

See [docs/CONFIGURATION.md](CONFIGURATION.md) for the full option reference and
[docs/EVENTS.md](EVENTS.md) for which events feed which stores.
