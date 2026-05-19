package connection

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// Guilds returns the client-level GuildManager backed by the local cache.
// Returns nil when no cache is configured.
func (d *Client) Guilds() discord.GuildManager {
	if d.Cache == nil {
		return nil
	}
	return managers.NewClientGuildManager(d)
}

// Users returns the client-level UserManager backed by the local cache.
// Returns nil when no cache is configured.
func (d *Client) Users() discord.UserManager {
	if d.Cache == nil {
		return nil
	}
	return managers.NewClientUserManager(d)
}

// Channels returns the client-level ChannelManager backed by the local cache.
// Returns nil when no cache is configured.
func (d *Client) Channels() discord.ChannelManager {
	if d.Cache == nil {
		return nil
	}
	return managers.NewClientChannelManager(d)
}
