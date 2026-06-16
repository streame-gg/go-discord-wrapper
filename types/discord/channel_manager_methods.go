package discord

// ── Channel sub-manager setters (called by connection layer) ──────────────────

func (c *Channel) SetMessagesManager(m MessageManager) { c.messages = m }
func (c *Channel) SetThreadsManager(m ThreadManager)   { c.threads = m }

// ── Channel sub-manager getters ───────────────────────────────────────────────

// Messages returns the manager for this channel's messages.
// Returns nil when the channel was not obtained from the gateway cache
// (e.g. a channel returned from a REST call or constructed directly).
// Only channels hydrated by the gateway client via GUILD_CREATE / CHANNEL_CREATE
// or retrieved from the cache have this manager populated.
func (c *Channel) Messages() MessageManager { return c.messages }

// Threads returns the manager for this channel's active threads.
// Returns nil when the channel was not obtained from the gateway cache.
// Only channels hydrated by the gateway client via GUILD_CREATE / CHANNEL_CREATE
// or retrieved from the cache have this manager populated.
func (c *Channel) Threads() ThreadManager { return c.threads }
