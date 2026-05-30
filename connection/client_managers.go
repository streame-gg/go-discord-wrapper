package connection

import (
	"github.com/streame-gg/go-discord-wrapper/collection"
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

// ThreadsForParent returns a snapshot of all cached threads whose parent channel
// is parentID. Uses the internal threadsByParent index — O(threads_in_parent)
// rather than O(total_channels). Returns an empty collection when no cache is
// configured. Satisfies the discord.EntityClient interface.
func (d *Client) ThreadsForParent(parentID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Channel] {
	result := collection.New[discord.Snowflake, *discord.Channel]()
	if d.Cache == nil {
		return result
	}
	d.threadIndexMu.RLock()
	ids := make([]discord.Snowflake, 0, len(d.threadsByParent[parentID]))
	for id := range d.threadsByParent[parentID] {
		ids = append(ids, id)
	}
	d.threadIndexMu.RUnlock()

	for _, id := range ids {
		if ch, ok := d.Cache.Channels().Get(id); ok {
			result.Set(id, ch)
		}
	}
	return result
}
