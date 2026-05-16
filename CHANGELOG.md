# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Changed

- **BREAKING — `ShardManager` factory returns `(*Client, error)`** (P0-9): The factory function passed to `NewShardManager` previously returned `*connection.Client`, making it impossible to propagate errors from `connection.NewClient`. The signature is now `func(shardID, totalShards int) (*connection.Client, error)`. `Start()` returns the factory error as `"sharding: shard N factory failed: ..."` and shuts down any already-started shards.

- **BREAKING — `events.EventFactories` is now unexported** (P0-6): The global `events.EventFactories` map has been renamed to `eventFactories` (unexported). Use the new `events.GetEventFactory(t EventType) (func() Event, bool)` for read access, and the existing `events.RegisterEvent` for writes. This prevents user code from accidentally overwriting entries or racing on the map at runtime.

- **BREAKING — `Client.OnEvent` now returns `error`** (P0-7): `OnEvent` previously logged a warning and silently dropped the handler when the handler had the wrong type signature. It now returns a non-nil error so the caller can surface the bug at registration time. The typed helpers (`OnMessageCreate`, `OnInteractionCreate`, etc.) are unaffected — they guarantee correct types at compile time and internally discard the (always-nil) error.

- **`WithMaxConcurrentEvents` semantics** (Bug 9): In pool mode (`MaxConcurrentEvents > 0`), event handlers now run **serially inside each worker goroutine** rather than spawning an unbounded number of goroutines per event. `MaxConcurrentEvents` is now the true upper bound on the number of handlers executing concurrently at any moment. Bots that relied on concurrent handler execution within a single event (i.e. multiple `On*` handlers firing simultaneously for the same event) will now see those handlers run in registration order, one after the other. If you need concurrent execution, register a single handler that fans out internally.

### Fixed

- **Bug 1 — `Login` blocked forever on context cancel without READY**: `Login` now returns `ctx.Err()` when the context is cancelled before the Discord READY event arrives, and returns an error when the gateway connection is closed before READY.
- **Bug 2 — Data race on `SessionID` and `ReconnectURL`**: Writes to `Websocket.SessionID` and `Websocket.ReconnectURL` in the READY handler now hold `wsMu.Lock()`, eliminating a data race detectable by `go test -race`.
- **Bug 3a — REST event handler could observe a closed response body**: The response body is now buffered once with `io.ReadAll`; event handlers receive a fresh `io.NopCloser` snapshot so reading or closing the body in a handler no longer corrupts the decoder.
- **Bug 3b — Decode error was silently swallowed**: A JSON decode failure for a successful HTTP response now returns `(nil, error)` instead of the previous `(nil, nil)`.
- **Bug 4 — Rate-limit bucket key included query string**: `routeKey` now strips query parameters and fragments from the path before building the key. `GetMessages` calls with different `?before=` / `?after=` parameters now share the same rate-limit bucket as Discord intends.
- **Bug 5 — `cacheMember` ignored `cacheStoreEnabled`**: `cacheMember` now checks `cacheStoreEnabled(CategoryMembers)` before writing to the member store, and checks `cacheStoreEnabled(CategoryUsers)` before writing the embedded user. A `GUILD_MEMBER_ADD` event will no longer populate the member cache when the `CategoryMembers` store is disabled via `WithCacheStores`.
- **Bug 6 — `RestClient.Close()` panicked on second call**: `rateLimiter.close()` is now idempotent via `sync.Once`; calling `RestClient.Close()` more than once no longer panics.
- **Bug 7 — Message store race with `DeleteChannel`**: `memMessageStore.Get`, `Update`, `Delete`, `DeleteBulk`, and `Channel` now hold the write lock for the entire operation, eliminating the TOCTOU window where `DeleteChannel` could remove the ring between the read-lock and write-lock phases.
- **Bug 8 — Cache `Set` methods panicked on nil**: `Set` on `memGuildStore`, `memChannelStore`, `memUserStore`, `memRoleStore`, `memVoiceStateStore`, `memScheduledEventStore`, `memStageInstanceStore`, `memEmojiStore`, and `memStickerStore` now silently ignores nil arguments instead of dereferencing a nil pointer.
- **Bug 10 — `LocalCoordinator.Register` panicked after `Close`**: `Register`, `Send`, and `Broadcast` now return an error when called on a closed coordinator instead of writing to a nil map.
- **Bug 11 — `MaxConcurrentEvents < 0` was accepted**: `options.Config.Validate()` now returns an error when `MaxConcurrentEvents` is negative.

## [0.1.0] - 2026-05-10

Initial release.

### Added

**Gateway / WebSocket**
- Persistent WebSocket connection to the Discord gateway with automatic reconnect and session resume.
- Configurable identify payload: intents, presence, shard ID.
- Panic recovery in event handlers so one bad handler cannot crash the gateway loop.
- `OnConnect`, `OnDisconnect`, `OnReconnect`, and `OnEvent` hooks.
- Middleware support via `Use()` — runs before every event handler.

**REST API**
- Full coverage of Discord REST API v10: guilds, channels, messages, members, roles, webhooks, voice, stage instances, invites, reactions, polls, threads, auto-moderation, scheduled events, entitlements, subscriptions, soundboard, and more (135+ exported methods).
- Proactive per-route and global rate limiter — waits for the correct reset window before sending a request, eliminating 429 responses.
- Configurable retry policy with exponential back-off for transient server errors.
- Configurable minimum interval between requests for extra safety margin.

**Interactions**
- `Reply`, `DeferReply`, `UpdateMessage`, `DeferUpdateMessage` — all four response types.
- `ReplyWithModal` for modal pop-ups.
- `ReplyAutocomplete` for autocomplete suggestions.
- `LaunchActivity` for embedded activity launch responses.
- Follow-up helpers: `GetOriginalResponse`, `EditReply`, `DeleteReply`, `CreateFollowup`, `EditFollowup`, `DeleteFollowup`.
- `RegisterCommand` and `BulkRegisterCommands` for slash command registration.

**Cache**
- In-process memory cache with TTL expiry and configurable LRU limits per entity type (guilds, channels, members, messages, roles, voice states, emojis, stickers, users, threads).
- Automatic cache population from gateway events (no extra calls needed).
- Pluggable cache interface — bring your own backend.
- Redis cache backend (`cache/rediscache`).
- MongoDB cache backend (`cache/mongocache`).

**Sharding**
- `options.WithSharding(total, shardID)` to run a single shard.
- Built-in local shard coordinator (`sharding.Manager`) for multi-shard processes.
- `BroadcastToShards`, `SendToShard`, `SubscribeShardResponse` for cross-shard messaging.

**Builders**
- Fluent builders for embeds, buttons, select menus (string, user, role, channel, mentionable), modals, action rows, and text inputs.

**Gateway events (76 events)**
- All standard Discord gateway events through API v10: ready, guild lifecycle, channel/thread lifecycle, messages, reactions, polls, interactions, voice, presence, moderation, integrations, webhooks, entitlements, subscriptions, soundboard, and scheduled events.

**Configuration**
- `options.WithCache`, `options.WithSharding`, `options.WithLogger`, `options.WithRetry`, `options.WithRateLimiting`, `options.WithMinRequestInterval`, `options.WithAPIVersion`.
