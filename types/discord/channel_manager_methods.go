package discord

// ── Channel sub-manager setters (called by connection layer) ──────────────────

func (ch *Channel) SetMessagesManager(m MessageManager) { ch.messages = m }
func (ch *Channel) SetThreadsManager(m ThreadManager)   { ch.threads = m }

// ── Channel sub-manager getters ───────────────────────────────────────────────

// Messages returns the manager for this channel's messages.
// Returns nil when the channel was not obtained from the gateway cache
// (e.g. a channel returned from a REST call or constructed directly).
// Only channels hydrated by the gateway client via GUILD_CREATE / CHANNEL_CREATE
// or retrieved from the cache have this manager populated.
func (ch *Channel) Messages() MessageManager { return ch.messages }

// Threads returns the manager for this channel's active threads.
// Returns nil when the channel was not obtained from the gateway cache.
// Only channels hydrated by the gateway client via GUILD_CREATE / CHANNEL_CREATE
// or retrieved from the cache have this manager populated.
func (ch *Channel) Threads() ThreadManager { return ch.threads }
