package discord

// ── Channel sub-manager setters (called by connection layer) ──────────────────

func (ch *Channel) SetMessagesManager(m MessageManager) { ch.messages = m }
func (ch *Channel) SetThreadsManager(m ThreadManager)   { ch.threads = m }

// ── Channel sub-manager getters ───────────────────────────────────────────────

// Messages returns the manager for this channel's messages.
func (ch *Channel) Messages() MessageManager { return ch.messages }

// Threads returns the manager for this channel's threads.
func (ch *Channel) Threads() ThreadManager { return ch.threads }
