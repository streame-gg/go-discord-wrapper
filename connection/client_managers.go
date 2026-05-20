package connection

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// Guilds returns the client-level GuildManager. Cache-backed methods return empty
// results when no cache is configured; Fetch/Create always work via REST.
func (d *Client) Guilds() discord.GuildManager {
	return managers.NewClientGuildManager(d)
}

// Users returns the client-level UserManager. Cache-backed methods return empty
// results when no cache is configured; Fetch always works via REST.
func (d *Client) Users() discord.UserManager {
	return managers.NewClientUserManager(d)
}

// Channels returns the client-level ChannelManager. Cache-backed methods return
// empty results when no cache is configured; Fetch/Create always work via REST.
func (d *Client) Channels() discord.ChannelManager {
	return managers.NewClientChannelManager(d)
}
