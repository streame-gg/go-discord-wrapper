package connection

import (
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (d *Client) cacheChannel(channel *discord.Channel) {
	if d.Cache == nil || channel == nil {
		return
	}

	channel.Hydrate(d)
	d.setChannelManagers(channel)
	d.trackChannel(channel)
	d.Cache.Channels().Set(channel)
	if channel.Recipients == nil {
		return
	}

	recipients := *channel.Recipients
	for i := range recipients {
		u := recipients[i]
		d.cacheUser(&u)
	}
}

func (d *Client) cacheChannels(channels []*discord.Channel) {
	for _, channel := range channels {
		d.cacheChannel(channel)
	}
}

func (d *Client) cacheGuild(guild *discord.Guild) {
	if d.Cache == nil || guild == nil {
		return
	}

	guild.Hydrate(d)
	d.setGuildManagers(guild)
	d.Cache.Guilds().Set(guild)
	for i := range guild.RawRoles {
		role := guild.RawRoles[i]
		d.Cache.Roles().Set(guild.ID, &role)
	}
}

func (d *Client) cacheMember(guildID discord.Snowflake, member *discord.GuildMember) {
	if member == nil || member.User == nil {
		return
	}
	member.GuildID = guildID
	member.UserID = member.User.ID
	member.Hydrate(d)
	if d.cacheStoreEnabled(cache.CategoryMembers) {
		d.Cache.Members().Set(guildID, member)
	}
	if d.cacheStoreEnabled(cache.CategoryUsers) {
		d.cacheUser(member.User)
	}
}

func (d *Client) cacheMembers(guildID discord.Snowflake, members []*discord.GuildMember) {
	for _, member := range members {
		d.cacheMember(guildID, member)
	}
}

func (d *Client) cacheMessage(msg *discord.Message) {
	if d.Cache == nil || msg == nil {
		return
	}

	msg.Hydrate(d)
	d.Cache.Messages().Add(msg)
	if msg.Author != nil {
		d.cacheUser(msg.Author)
	}
}

func (d *Client) cacheMessages(messages []*discord.Message) {
	for _, msg := range messages {
		d.cacheMessage(msg)
	}
}

func (d *Client) cacheRole(guildID discord.Snowflake, role *discord.Role) {
	if d.Cache == nil || role == nil {
		return
	}

	role.GuildID = guildID
	role.Hydrate(d)
	d.Cache.Roles().Set(guildID, role)
}

func (d *Client) cacheRoles(guildID discord.Snowflake, roles []*discord.Role) {
	for _, role := range roles {
		d.cacheRole(guildID, role)
	}
}

func (d *Client) cacheUser(user *discord.User) {
	if d.Cache == nil || user == nil {
		return
	}

	user.Hydrate(d)
	d.Cache.Users().Set(user)
}

func (d *Client) cacheUsers(users []*discord.User) {
	for _, user := range users {
		d.cacheUser(user)
	}
}

func (d *Client) removeChannelFromCache(channelID discord.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.untrackChannel(channelID)
	d.Cache.Channels().Delete(channelID)
	d.Cache.Messages().DeleteChannel(channelID)

	// If channelID is a parent channel, evict all threads it hosted too.
	for _, threadID := range d.drainParentThreadIDs(channelID) {
		d.untrackChannel(threadID)
		d.Cache.Channels().Delete(threadID)
		d.Cache.Messages().DeleteChannel(threadID)
	}
}

func (d *Client) removeGuildFromCache(guildID discord.Snowflake) {
	d.guildMu.Lock()
	delete(d.guildMemberCounts, guildID)
	d.guildMu.Unlock()

	if d.Cache == nil {
		return
	}

	d.Cache.Guilds().Delete(guildID)
	d.Cache.Members().DeleteGuild(guildID)
	d.Cache.Roles().DeleteGuild(guildID)
	d.Cache.VoiceStates().DeleteGuild(guildID)
	d.Cache.Presences().DeleteGuild(guildID)
	d.Cache.Soundboard().DeleteGuild(guildID)
	d.Cache.ScheduledEvents().DeleteGuild(guildID)
	d.Cache.StageInstances().DeleteGuild(guildID)
	d.Cache.Emojis().DeleteGuild(guildID)
	d.Cache.Stickers().DeleteGuild(guildID)

	for _, channelID := range d.drainGuildChannelIDs(guildID) {
		d.Cache.Channels().Delete(channelID)
		d.Cache.Messages().DeleteChannel(channelID)
	}
}

func (d *Client) removeGuildMemberFromCache(guildID, userID discord.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Members().Delete(guildID, userID)
}

func (d *Client) removeMessageFromCache(channelID, messageID discord.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Messages().Delete(channelID, messageID)
}

func (d *Client) removeMessagesFromCache(channelID discord.Snowflake, messageIDs []discord.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Messages().DeleteBulk(channelID, messageIDs)
}

func (d *Client) removeRoleFromCache(roleID discord.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Roles().Delete(roleID)
}

// trackChannel records the channel→guild association in the bidirectional
// index so that all channels belonging to a guild can be found efficiently.
// If the channel has no GuildID it is a DM channel and is not indexed.
// Safe for concurrent use; acquires channelIndexMu internally.
func (d *Client) trackChannel(channel *discord.Channel) {
	if channel == nil {
		return
	}

	d.channelIndexMu.Lock()
	defer d.channelIndexMu.Unlock()

	if oldGuildID, ok := d.guildByChannel[channel.ID]; ok {
		if channel.GuildID == nil || oldGuildID != *channel.GuildID {
			if set := d.channelsByGuild[oldGuildID]; set != nil {
				delete(set, channel.ID)
				if len(set) == 0 {
					delete(d.channelsByGuild, oldGuildID)
				}
			}
			delete(d.guildByChannel, channel.ID)
		}
	}

	if channel.GuildID == nil {
		return
	}

	guildID := *channel.GuildID
	set := d.channelsByGuild[guildID]
	if set == nil {
		set = make(map[discord.Snowflake]struct{})
		d.channelsByGuild[guildID] = set
	}
	set[channel.ID] = struct{}{}
	d.guildByChannel[channel.ID] = guildID
}

// untrackChannel removes a channel from the bidirectional index. It is called
// when a channel is deleted (CHANNEL_DELETE) or when the entire guild is
// evicted from the cache.
// Safe for concurrent use; acquires channelIndexMu internally.
func (d *Client) untrackChannel(channelID discord.Snowflake) {
	d.channelIndexMu.Lock()
	defer d.channelIndexMu.Unlock()

	guildID, ok := d.guildByChannel[channelID]
	if !ok {
		return
	}

	if set := d.channelsByGuild[guildID]; set != nil {
		delete(set, channelID)
		if len(set) == 0 {
			delete(d.channelsByGuild, guildID)
		}
	}
	delete(d.guildByChannel, channelID)
}

// drainGuildChannelIDs atomically removes and returns all channel IDs
// associated with guildID from the bidirectional index. It is used during
// GUILD_DELETE processing to collect every channel that must be evicted from
// the cache.
// Safe for concurrent use; acquires channelIndexMu internally.
func (d *Client) drainGuildChannelIDs(guildID discord.Snowflake) []discord.Snowflake {
	d.channelIndexMu.Lock()
	defer d.channelIndexMu.Unlock()

	set := d.channelsByGuild[guildID]
	if len(set) == 0 {
		delete(d.channelsByGuild, guildID)
		return nil
	}

	ids := make([]discord.Snowflake, 0, len(set))
	for channelID := range set {
		ids = append(ids, channelID)
		delete(d.guildByChannel, channelID)
	}
	delete(d.channelsByGuild, guildID)

	return ids
}

// trackThread records a thread's ID under its parent channel in threadsByParent.
// Safe for concurrent use; acquires threadIndexMu internally.
func (d *Client) trackThread(thread *discord.Channel) {
	if thread == nil || thread.ParentID == nil {
		return
	}
	d.threadIndexMu.Lock()
	defer d.threadIndexMu.Unlock()
	parentID := *thread.ParentID
	set := d.threadsByParent[parentID]
	if set == nil {
		set = make(map[discord.Snowflake]struct{})
		d.threadsByParent[parentID] = set
	}
	set[thread.ID] = struct{}{}
}

// untrackThread removes a thread from threadsByParent.
// Safe for concurrent use; acquires threadIndexMu internally.
func (d *Client) untrackThread(threadID, parentID discord.Snowflake) {
	d.threadIndexMu.Lock()
	defer d.threadIndexMu.Unlock()
	set := d.threadsByParent[parentID]
	if set == nil {
		return
	}
	delete(set, threadID)
	if len(set) == 0 {
		delete(d.threadsByParent, parentID)
	}
}

// drainParentThreadIDs atomically removes and returns all thread IDs associated
// with parentID from the threadsByParent index.
// Safe for concurrent use; acquires threadIndexMu internally.
func (d *Client) drainParentThreadIDs(parentID discord.Snowflake) []discord.Snowflake {
	d.threadIndexMu.Lock()
	defer d.threadIndexMu.Unlock()
	set := d.threadsByParent[parentID]
	if len(set) == 0 {
		delete(d.threadsByParent, parentID)
		return nil
	}
	ids := make([]discord.Snowflake, 0, len(set))
	for threadID := range set {
		ids = append(ids, threadID)
	}
	delete(d.threadsByParent, parentID)
	return ids
}
