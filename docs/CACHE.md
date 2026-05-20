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

4. **Gateway request** — `client.RequestGuildMembers(ctx, ...)` sends
   opcode 8. Members arrive via `GUILD_MEMBERS_CHUNK` events and populate
   the cache asynchronously. Requires the `GUILD_MEMBERS` intent.

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

## What is NOT cached

The library intentionally does not cache:

| Entity | Reason |
|--------|--------|
| Bans | No gateway event stream; REST-only via `guild.Bans().FetchAll()` |
| Invites | No gateway event stream; REST-only via `guild.Invites().FetchAll()` |
| Webhooks | No gateway event stream; REST-only via `guild.Webhooks().FetchAll()` |
| Integrations | REST-only via `guild.Integrations().FetchAll()` |
| Auto-mod rules | REST-only via `guild.AutoModRules().FetchAll()` |

These managers exist on every guild but have no `Cache()` or `Get()`
methods — they expose only REST operations.
