# Sharding

Discord requires bots in many guilds to split their gateway connection across
**shards**. Each shard is an independent WebSocket connection that handles a
deterministic slice of your guilds, chosen by Discord from
`(guild_id >> 22) % total_shards`. A bot under ~2,500 guilds can run on a single
shard; past that, Discord mandates more.

There are two ways to shard with this library, depending on whether all shards
run in one process or are spread across machines.

## Single shard

If you do not pass any sharding options, the client runs as shard 0 of 1 — a
normal single connection. Nothing to configure.

## Multiple shards in one process: `ShardManager`

For a bot that runs every shard inside one binary, use the `sharding` package.
A `ShardManager` builds one `Client` per shard, starts them in the order and
concurrency Discord allows, and wires them together through a coordinator.

```go
coord := sharding.NewLocalCoordinator(4) // 4 shards, in-process message bus

factory := func(shardID, total int) (*connection.Client, error) {
    return connection.NewClient(token, intents,
        options.WithSharding(total, shardID),
        options.WithCoordinator(coord),
    )
}

mgr := sharding.NewShardManager(coord, factory)
if err := mgr.Start(); err != nil { /* … */ }
defer mgr.Shutdown()
```

- **`factory`** is called once per shard ID and must return a configured but
  **not-yet-logged-in** `Client`. `Start` performs the `Login`.
- **`Start`** blocks until every shard has IDENTIFY'd and received READY. If any
  shard fails, the shards already started are shut down before it returns the
  error, so you are never left with a half-started fleet.
- **`Shutdown`** disconnects all shards and closes the coordinator. Every shard
  is shut down even if one errors; the errors are joined.

A complete runnable program is in
[`example/sharding`](../example/sharding/main.go).

### IDENTIFY buckets and `max_concurrency`

Discord rate-limits how fast you may IDENTIFY: shards are grouped into buckets
of `max_concurrency`, and a 5-second gap is enforced between buckets.

`NewShardManager` defaults to `max_concurrency = 1` (one shard at a time —
safe, but slow for large fleets). To start whole buckets in parallel, give the
manager a `RestClient` so it can read the real `max_concurrency` from
`/gateway/bot`:

```go
rest, _ := api.NewRestClient(token)
mgr := sharding.NewShardManagerWithRest(coord, factory, rest)
```

If the fetch fails, it falls back to `max_concurrency = 1`.

### Accessing individual shards

| Method | Returns |
|--------|---------|
| `mgr.Shard(id)` | the `*connection.Client` for one shard (nil if out of range / not started) |
| `mgr.Shards()` | a snapshot slice of all shard clients, by ascending shard ID |

## Multiple machines: external coordinator

You do **not** need a `ShardManager` for multi-machine deployments. Run one
`Client` per machine (or per group of shards) and point each at an external
broker by implementing `options.ShardCoordinator` (backed by Redis, NATS,
etc.):

```go
bot, _ := connection.NewClient(token, intents,
    options.WithSharding(totalShards, thisShardID),
    options.WithCoordinator(myRedisCoordinator),
)
```

The coordinator interface is small — `Register`, `Send`, `Broadcast`,
`TotalShards`, and `Close` — and `sharding.LocalCoordinator` is a reference
implementation for the in-process case.

## Cross-shard messaging

Because each shard only knows about its own guilds, answering a question like
"how many guilds are we in total?" requires asking every shard. The coordinator
provides a small message bus for exactly this.

Each client exposes:

| Method | Does |
|--------|------|
| `client.OnShardMessage(fn)` | register a handler for messages from other shards |
| `client.SendToShard(id, type, payload)` | message one shard |
| `client.BroadcastToShards(type, payload)` | message every shard (including self) |
| `client.ReplyToShard(orig, type, payload)` | reply, echoing the correlation ID |

Messages carry a user-defined `Type` string and a JSON `Payload`. The typical
pattern is request/response: broadcast a request, have each shard reply with its
local answer, and aggregate.

`sharding.RequestAll` wraps that whole dance — broadcast, collect one typed
response per shard, return when all have answered or the context expires:

```go
// On EVERY shard: answer the request with this shard's local count.
client.OnShardMessage(func(msg options.ShardMessage) {
    if msg.Type == "GUILD_COUNT_REQ" {
        _ = client.ReplyToShard(msg, "GUILD_COUNT_RESP", client.Cache.Guilds().Size())
    }
})

// On the shard that needs the total:
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
counts, err := sharding.RequestAll[int](ctx, client, "GUILD_COUNT_REQ", nil, "GUILD_COUNT_RESP")
// counts holds one int per shard that replied; sum them.
```

`RequestAll` returns partial results on timeout — compare `len(counts)` against
`coord.TotalShards()` to detect a shard that didn't answer.

## Choosing a shard count

- Start with what Discord recommends: `GetBotGateway(ctx).Shards`.
- Shard count must stay constant for a session; changing it requires a fresh
  IDENTIFY for every shard.
- More shards means more connections to manage but smaller per-shard event
  volume. Scale up when a single shard can't keep pace with its event stream
  (watch for dropped events — see `options.WithMaxConcurrentEvents`).

## See also

- [docs/CONFIGURATION.md](CONFIGURATION.md) — `WithSharding`, `WithCoordinator`,
  and event-throughput options.
- [docs/EVENTS.md](EVENTS.md) — how events are dispatched per shard.
