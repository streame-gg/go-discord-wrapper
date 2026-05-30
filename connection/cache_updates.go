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
		role.GuildID = guild.ID
		role.Hydrate(d)
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

func (d *Client) cacheBan(guildID discord.Snowflake, ban *discord.Ban) {
	if ban == nil {
		return
	}
	ban.User.Hydrate(d)
	if d.cacheStoreEnabled(cache.CategoryBans) {
		d.Cache.Bans().Set(guildID, ban)
	}
	if d.cacheStoreEnabled(cache.CategoryUsers) {
		d.cacheUser(&ban.User)
	}
}

func (d *Client) cacheInviteForGuild(guildID discord.Snowflake, invite *discord.Invite) {
	if invite == nil || !d.cacheStoreEnabled(cache.CategoryInvites) {
		return
	}
	d.Cache.Invites().SetWithGuild(guildID, invite)
}

func (d *Client) cacheInvite(invite *discord.Invite) {
	if invite == nil || !d.cacheStoreEnabled(cache.CategoryInvites) {
		return
	}
	d.Cache.Invites().Set(invite)
}

func (d *Client) cacheAutoModRule(guildID discord.Snowflake, rule *discord.AutoModerationRule) {
	if rule == nil || !d.cacheStoreEnabled(cache.CategoryAutoModRules) {
		return
	}
	rule.Hydrate(d)
	d.Cache.AutoModRules().Set(guildID, rule)
}

func (d *Client) cacheUser(user *discord.User) {
	if d.Cache == nil || user == nil {
		return
	}

	user.Hydrate(d)
	d.Cache.Users().Set(user)
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
	d.Cache.Bans().DeleteGuild(guildID)
	d.Cache.AutoModRules().DeleteGuild(guildID)
	d.Cache.Invites().DeleteGuild(guildID)

	for _, channelID := range d.drainGuildChannelIDs(guildID) {
		d.Cache.Channels().Delete(channelID)
		d.Cache.Messages().DeleteChannel(channelID)
		// Bug 83: evict threadsByParent entries for any channel that was a parent.
		d.drainParentThreadIDs(channelID)
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

// reactionEmojiMatches returns true when a and b refer to the same emoji.
// Custom emojis are matched by ID; unicode emojis are matched by Name.
func reactionEmojiMatches(a, b discord.Emoji) bool {
	if a.ID != 0 && b.ID != 0 {
		return a.ID == b.ID
	}
	return a.Name != "" && a.Name == b.Name
}

// appendOrIncrementReaction returns a new slice with count for emoji incremented
// by one, or a new Reaction appended if the emoji is not yet present.
func appendOrIncrementReaction(existing *[]discord.Reaction, emoji discord.Emoji) []discord.Reaction {
	var reactions []discord.Reaction
	if existing != nil {
		reactions = make([]discord.Reaction, len(*existing))
		copy(reactions, *existing)
	}
	for i := range reactions {
		if reactionEmojiMatches(reactions[i].Emoji, emoji) {
			reactions[i].Count++
			return reactions
		}
	}
	return append(reactions, discord.Reaction{Count: 1, Emoji: emoji})
}

// decrementOrRemoveReaction returns a new slice with the count for emoji
// decremented by one; the entry is dropped when the count reaches zero.
func decrementOrRemoveReaction(existing *[]discord.Reaction, emoji discord.Emoji) []discord.Reaction {
	if existing == nil {
		return nil
	}
	reactions := make([]discord.Reaction, 0, len(*existing))
	for _, r := range *existing {
		if !reactionEmojiMatches(r.Emoji, emoji) {
			reactions = append(reactions, r)
			continue
		}
		if r.Count > 1 {
			r.Count--
			reactions = append(reactions, r)
		}
	}
	return reactions
}

// removeEmojiReaction returns a new slice with all reactions for the given
// emoji removed, regardless of count.
func removeEmojiReaction(existing *[]discord.Reaction, emoji discord.Emoji) []discord.Reaction {
	if existing == nil {
		return nil
	}
	reactions := make([]discord.Reaction, 0, len(*existing))
	for _, r := range *existing {
		if !reactionEmojiMatches(r.Emoji, emoji) {
			reactions = append(reactions, r)
		}
	}
	return reactions
}

func (d *Client) removeRoleFromCache(guildID, roleID discord.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Roles().Delete(roleID)

	if d.cacheStoreEnabled(cache.CategoryMembers) {
		roleStr := roleID.String()
		for _, m := range d.Cache.Members().AllInGuild(guildID).Values() {
			found := false
			for _, r := range m.Roles {
				if r.String() == roleStr {
					found = true
					break
				}
			}
			if !found {
				continue
			}
			updated := *m
			filtered := make([]discord.Snowflake, 0, len(m.Roles)-1)
			for _, r := range m.Roles {
				if r.String() != roleStr {
					filtered = append(filtered, r)
				}
			}
			updated.Roles = filtered
			d.Cache.Members().Set(guildID, &updated)
		}
	}
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
