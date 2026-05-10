# Changelog

All notable changes to this project will be documented in this file.

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
