# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## v0.2.0 Breaking Changes

- **P2-25 — `AuditLog` fields typed**: `AutoModerationRules`, `GuildScheduledEvents`, `Integrations`, and `Webhooks` are now typed slices instead of `[]any`. `ApplicationCommands` remains `[]any` due to an import-cycle constraint (the `ApplicationCommand` struct lives in `types/commands` which imports `types/discord`); unmarshal elements into `*commands.ApplicationCommand` individually. `AuditLogEntryChange.NewValue`/`OldValue` also remain `any` because Discord audit-log change values are genuinely polymorphic.

  ```go
  // Before:
  for _, raw := range log.AutoModerationRules {
      rule := raw.(map[string]any) // fragile
  }

  // After:
  for _, rule := range log.AutoModerationRules {
      // rule is *discord.AutoModerationRule — fully typed
      fmt.Println(rule.Name)
  }
  ```

- **P2-28 — `ComponentType` select-menu constants renamed for consistency**: All select-menu component types now use the `...Select` suffix, matching Discord's canonical naming. Deprecated aliases with the old names are provided for one release cycle.

  | Old name | New name |
  |---|---|
  | `ComponentTypeStringSelectMenu` | `ComponentTypeStringSelect` |
  | `ComponentTypeUserSelectMenu` | `ComponentTypeUserSelect` |
  | `ComponentTypeRoleSelectMenu` | `ComponentTypeRoleSelect` |
  | `ComponentTypeMentionableMenu` | `ComponentTypeMentionableSelect` |
  | `ComponentTypeChannelSelect` | *(unchanged)* |

  ```go
  // Before:
  discord.ComponentTypeStringSelectMenu
  discord.ComponentTypeMentionableMenu

  // After:
  discord.ComponentTypeStringSelect
  discord.ComponentTypeMentionableSelect
  ```

- **P2-31 — `RedisCache.WithKeyPrefix` derivative now has an independent lifecycle**: Previously, the returned `*RedisCache` shared `ctx`, `cancel`, and `stopOnce` with the original. Calling `Close()` on the derivative cancelled the original's context (and vice versa), causing all subsequent Redis operations on the surviving instance to fail. Each instance created by `WithKeyPrefix` now holds its own `context.WithCancel` root and `sync.Once`, so `Close()` affects only that instance.

- **P2-36 — `WithAuditLogReason` now validates and encodes the reason**: The reason string is now URL percent-encoded (`url.PathEscape`) so non-ASCII characters (Umlauts, emoji, CJK) are transmitted safely in the `X-Audit-Log-Reason` header. Reasons longer than Discord's 512-character limit are silently truncated. Empty reasons still set no header. No API-signature change; the behaviour change is invisible to callers that only pass short ASCII reasons.

- **P2-34 — `DeferAndFollowup` send-callback now takes an explicit `ctx` parameter**: The closure returned by `DeferAndFollowup` previously captured the `ctx` passed to the defer call. If that context was cancelled by the time the follow-up was ready to send, `CreateFollowup` would fail. The send-function now accepts its own `context.Context` as the first argument.

  ```go
  // Before:
  send, err := i.DeferAndFollowup(ctx, client, false)
  // ...
  send(api.CreateMessageParams{Content: "Done!"})

  // After:
  send, err := i.DeferAndFollowup(ctx, client, false)
  // ...
  send(ctx, api.CreateMessageParams{Content: "Done!"})
  // or use a fresh context for the follow-up:
  send(context.Background(), api.CreateMessageParams{Content: "Done!"})
  ```

- **P2-33 — `EmbedBuilder.Build()` no longer aliases the builder's `Fields` slice**: `Build()` now returns an embed whose `Fields` slice is backed by an independent array. Calling `AddFields` after `Build()` will not mutate the previously returned `Embed`. Only `Fields` was affected; all other `Embed` pointer fields (`Title`, `Description`, `Footer`, etc.) were already independent because each setter stores a fresh pointer.

- **P2-32 — `ModalBuilder.Build()` no longer aliases the builder's internal state**: `Build()` now returns a deep copy of the modal. Calling `AddComponents` after `Build()` will not mutate the previously returned `*Modal`. Code that relied on post-`Build()` mutations observing the returned value must be updated (such usage was a bug).

- **P2-27 — `AutoModerationTriggerMetadata.Presets` typed enum**: `Presets []int` is now `Presets []KeywordPresetType`. Use the new constants `KeywordPresetTypeProfanity` (1), `KeywordPresetTypeSexualContent` (2), `KeywordPresetTypeSlurs` (3) instead of raw integer literals.

  ```go
  // Before:
  meta := discord.AutoModerationTriggerMetadata{
      Presets: []int{1, 3},
  }

  // After:
  meta := discord.AutoModerationTriggerMetadata{
      Presets: []discord.KeywordPresetType{
          discord.KeywordPresetTypeProfanity,
          discord.KeywordPresetTypeSlurs,
      },
  }
  ```

### Added

- **CI pipeline** (P0-12): Added `.github/workflows/ci.yml` with three jobs: `build` (`go build` + `go vet`), `test-race` (`go test -race -timeout 5m`), and `doc-examples` (builds and runs `cmd/doc-check`). Triggers on push and pull request for all branches.
- **`cmd/doc-check` tool** (P0-12): New CLI that extracts fenced ` ```go ` blocks from `README.md` and `//\t`-indented godoc examples from source files, compiles each in isolation, and exits 1 on any compile error. Surfaced two pre-existing documentation bugs (P0-1, P0-11) at CI time.

### Changed

- **BREAKING — `ShardManager` factory returns `(*Client, error)`** (P0-9): The factory function passed to `NewShardManager` previously returned `*connection.Client`, making it impossible to propagate errors from `connection.NewClient`. The signature is now `func(shardID, totalShards int) (*connection.Client, error)`. `Start()` returns the factory error as `"sharding: shard N factory failed: ..."` and shuts down any already-started shards.

- **BREAKING — `events.EventFactories` is now unexported** (P0-6): The global `events.EventFactories` map has been renamed to `eventFactories` (unexported). Use the new `events.GetEventFactory(t EventType) (func() Event, bool)` for read access, and the existing `events.RegisterEvent` for writes. This prevents user code from accidentally overwriting entries or racing on the map at runtime.

- **BREAKING — `Client.OnEvent` now returns `error`** (P0-7): `OnEvent` previously logged a warning and silently dropped the handler when the handler had the wrong type signature. It now returns a non-nil error so the caller can surface the bug at registration time. The typed helpers (`OnMessageCreate`, `OnInteractionCreate`, etc.) are unaffected — they guarantee correct types at compile time and internally discard the (always-nil) error.

- **`WithMaxConcurrentEvents` semantics** (Bug 9): In pool mode (`MaxConcurrentEvents > 0`), event handlers now run **serially inside each worker goroutine** rather than spawning an unbounded number of goroutines per event. `MaxConcurrentEvents` is now the true upper bound on the number of handlers executing concurrently at any moment. Bots that relied on concurrent handler execution within a single event (i.e. multiple `On*` handlers firing simultaneously for the same event) will now see those handlers run in registration order, one after the other. If you need concurrent execution, register a single handler that fans out internally.

### Fixed

- **P0-1 — README import paths referenced removed package `types/common`**: Updated all code examples to import `types/discord` and use `discord.IntentGuilds`, `discord.IntentGuildMessages`. Fixed the cache-backends table to use the correct module paths `cache/rediscache` and `cache/mongocache`. Updated `OnMessageCreate` handler to use the typed helper signature.
- **P0-2 — Login-Reconnect Race: `Login` could block forever after reconnect**: `Websocket.Ready` was a per-connection channel that was replaced on each reconnect; `Login` held a reference to the old channel and blocked forever if READY arrived on the new connection first. Moved the signal to a client-level `readyCh chan struct{}` closed exactly once via `sync.Once`. Added `Client.Ready() <-chan struct{}` as a stable, reconnect-safe signal for callers.
- **P0-3 — `NewWebsocket` leaked the WebSocket connection on error**: If any step after `dialer.Dial` (reading hello frame, sending identify, spawning heartbeat) failed, the underlying `*websocket.Conn` was never closed. Added a `dialSucceeded` guard: the deferred closer skips cleanup only when the function returns successfully.
- **P0-4 — Reconnect backoff ignored shutdown signal**: The reconnect loop used `time.Sleep(backoff)`, blocking up to 5 minutes and preventing `Shutdown()` from returning promptly. Replaced with `select { case <-time.After(backoff): case <-shutdownCh: return nil }` so shutdown is always honoured within one scheduler tick.
- **P0-5 — Zero-limit rate-limit bucket hung indefinitely**: When Discord returns `X-RateLimit-Limit: 0` (undocumented routes), `bucket.limit` was 0 and `bucket.remaining` was never restored after reset, causing every subsequent call to sleep forever. Both restore sites now fall back to `remaining = 1` when `limit == 0`, letting one request through to re-learn the actual limit.
- **P0-10 — `CreateForumThread` silently dropped `Message.Files` attachments**: The method always serialised params with `json.Marshal`, but `CreateMessageParams.Files` is tagged `json:"-"` and therefore never included. When `params.Message.Files` is non-empty, the request is now sent as `multipart/form-data` with a `payload_json` part (matching the `CreateMessage` behaviour). The `opts.Reason` header is forwarded in both code paths.
- **P0-11 — `DeferAndFollowup` godoc example called function with wrong arity**: The example omitted the required `ctx context.Context` first argument, causing `go vet` / doc-check failures. Fixed the call to `i.DeferAndFollowup(ctx, client, false)` and simplified the example body to avoid undefined placeholder symbols.

- **P1-20 — `NewRestClient` accepted malformed tokens silently**: `NewRestClient` now strips leading/trailing whitespace and `"Bot "`/`"Bearer "` prefixes before storing the token, and returns an error immediately when the result is empty. Callers that accidentally pass `"Bot <secret>"` (the full Authorization header value) now get the correct raw secret without any extra stripping in `generateRequest`.

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
