package discord

// ── Guild sub-manager setters (called by connection layer) ────────────────────

func (g *Guild) SetMembersManager(m MemberManager)         { g.members = m }
func (g *Guild) SetRolesManager(m RoleManager)             { g.roles = m }
func (g *Guild) SetChannelsManager(m GuildChannelManager)  { g.channels = m }
func (g *Guild) SetEmojisManager(m EmojiManager)           { g.emojis = m }
func (g *Guild) SetStickersManager(m StickerManager)       { g.stickers = m }
func (g *Guild) SetBansManager(m BanManager)               { g.bans = m }
func (g *Guild) SetScheduledEventsManager(m ScheduledEventManager) {
	g.scheduledEvents = m
}
func (g *Guild) SetStageInstancesManager(m StageInstanceManager) {
	g.stageInstances = m
}
func (g *Guild) SetSoundboardManager(m SoundboardManager)      { g.soundboardSounds = m }
func (g *Guild) SetInvitesManager(m GuildInviteManager)        { g.invites = m }
func (g *Guild) SetVoiceStatesManager(m VoiceStateManager)     { g.voiceStates = m }
func (g *Guild) SetAutoModRulesManager(m AutoModRuleManager)   { g.autoModRules = m }
func (g *Guild) SetWebhooksManager(m GuildWebhookManager)      { g.webhooks = m }
func (g *Guild) SetIntegrationsManager(m IntegrationManager)   { g.integrations = m }

// ── Guild sub-manager getters ─────────────────────────────────────────────────

// Members returns the manager for this guild's members.
func (g *Guild) Members() MemberManager { return g.members }

// Roles returns the manager for this guild's roles.
func (g *Guild) Roles() RoleManager { return g.roles }

// Channels returns the manager for this guild's channels.
func (g *Guild) Channels() GuildChannelManager { return g.channels }

// Emojis returns the manager for this guild's emojis.
func (g *Guild) Emojis() EmojiManager { return g.emojis }

// Stickers returns the manager for this guild's stickers.
func (g *Guild) Stickers() StickerManager { return g.stickers }

// Bans returns the manager for this guild's bans.
func (g *Guild) Bans() BanManager { return g.bans }

// ScheduledEvents returns the manager for this guild's scheduled events.
func (g *Guild) ScheduledEvents() ScheduledEventManager { return g.scheduledEvents }

// StageInstances returns the manager for this guild's stage instances.
func (g *Guild) StageInstances() StageInstanceManager { return g.stageInstances }

// Soundboard returns the manager for this guild's soundboard sounds.
func (g *Guild) Soundboard() SoundboardManager { return g.soundboardSounds }

// Invites returns the manager for this guild's invites.
func (g *Guild) Invites() GuildInviteManager { return g.invites }

// VoiceStates returns the manager for this guild's voice states.
func (g *Guild) VoiceStates() VoiceStateManager { return g.voiceStates }

// AutoModRules returns the manager for this guild's auto moderation rules.
func (g *Guild) AutoModRules() AutoModRuleManager { return g.autoModRules }

// Webhooks returns the manager for this guild's webhooks.
func (g *Guild) Webhooks() GuildWebhookManager { return g.webhooks }

// Integrations returns the manager for this guild's integrations.
func (g *Guild) Integrations() IntegrationManager { return g.integrations }
