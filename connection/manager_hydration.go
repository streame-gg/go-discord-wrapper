package connection

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// setGuildManagers injects all sub-managers into a hydrated Guild.
// Called by the connection layer after guild.Hydrate(d).
func (d *Client) setGuildManagers(g *discord.Guild) {
	g.SetMembersManager(managers.NewMemberManager(g.ID, d))
	g.SetRolesManager(managers.NewRoleManager(g.ID, d))
	g.SetChannelsManager(managers.NewGuildChannelManager(g.ID, d))
	g.SetEmojisManager(managers.NewEmojiManager(g.ID, d))
	g.SetStickersManager(managers.NewStickerManager(g.ID, d))
	g.SetBansManager(managers.NewBanManager(g.ID, d))
	g.SetScheduledEventsManager(managers.NewScheduledEventManager(g.ID, d))
	g.SetStageInstancesManager(managers.NewStageInstanceManager(g.ID, d))
	g.SetSoundboardManager(managers.NewSoundboardManager(g.ID, d))
	g.SetInvitesManager(managers.NewGuildInviteManager(g.ID, d))
	g.SetVoiceStatesManager(managers.NewVoiceStateManager(g.ID, d))
	g.SetAutoModRulesManager(managers.NewAutoModRuleManager(g.ID, d))
	g.SetWebhooksManager(managers.NewGuildWebhookManager(g.ID, d))
	g.SetIntegrationsManager(managers.NewIntegrationManager(g.ID, d))
}

// setChannelManagers injects all sub-managers into a hydrated Channel.
// Called by the connection layer after ch.Hydrate(d).
func (d *Client) setChannelManagers(ch *discord.Channel) {
	ch.SetMessagesManager(managers.NewMessageManager(ch.ID, d))
	ch.SetThreadsManager(managers.NewThreadManager(ch.ID, d))
}
