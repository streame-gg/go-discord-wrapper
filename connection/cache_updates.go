package connection

import "github.com/streame-gg/go-discord-wrapper/types/common"

func (d *Client) cacheChannel(channel *common.Channel) {
	if d.Cache == nil || channel == nil {
		return
	}

	d.Cache.Channels().Set(channel)
	if channel.Recipients == nil {
		return
	}

	recipients := *channel.Recipients
	for i := range recipients {
		user := recipients[i]
		d.cacheUser(&user)
	}
}

func (d *Client) cacheChannels(channels []*common.Channel) {
	for _, channel := range channels {
		d.cacheChannel(channel)
	}
}

func (d *Client) cacheGuild(guild *common.Guild) {
	if d.Cache == nil || guild == nil {
		return
	}

	d.Cache.Guilds().Set(guild)
	for i := range guild.Roles {
		role := guild.Roles[i]
		d.Cache.Roles().Set(guild.ID, &role)
	}
}

func (d *Client) cacheMember(guildID common.Snowflake, member *common.GuildMember) {
	if d.Cache == nil || member == nil || member.User == nil {
		return
	}

	d.Cache.Members().Set(guildID, member)
	d.cacheUser(member.User)
}

func (d *Client) cacheMembers(guildID common.Snowflake, members []*common.GuildMember) {
	for _, member := range members {
		d.cacheMember(guildID, member)
	}
}

func (d *Client) cacheMessage(msg *common.Message) {
	if d.Cache == nil || msg == nil {
		return
	}

	d.Cache.Messages().Add(msg)
	if msg.Author != nil {
		d.cacheUser(msg.Author)
	}
}

func (d *Client) cacheMessages(messages []*common.Message) {
	for _, msg := range messages {
		d.cacheMessage(msg)
	}
}

func (d *Client) cacheRole(guildID common.Snowflake, role *common.Role) {
	if d.Cache == nil || role == nil {
		return
	}

	d.Cache.Roles().Set(guildID, role)
}

func (d *Client) cacheRoles(guildID common.Snowflake, roles []*common.Role) {
	for _, role := range roles {
		d.cacheRole(guildID, role)
	}
}

func (d *Client) cacheUser(user *common.User) {
	if d.Cache == nil || user == nil {
		return
	}

	d.Cache.Users().Set(user)
}

func (d *Client) cacheUsers(users []*common.User) {
	for _, user := range users {
		d.cacheUser(user)
	}
}

func (d *Client) removeChannelFromCache(channelID common.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Channels().Delete(channelID)
	d.Cache.Messages().DeleteChannel(channelID)
}

func (d *Client) removeGuildFromCache(guildID common.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Guilds().Delete(guildID)
	d.Cache.Members().DeleteGuild(guildID)
	d.Cache.Roles().DeleteGuild(guildID)

	for _, channel := range d.Cache.Channels().All() {
		if channel.GuildID != nil && *channel.GuildID == guildID {
			d.Cache.Channels().Delete(channel.ID)
			d.Cache.Messages().DeleteChannel(channel.ID)
		}
	}
}

func (d *Client) removeGuildMemberFromCache(guildID, userID common.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Members().Delete(guildID, userID)
}

func (d *Client) removeMessageFromCache(channelID, messageID common.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Messages().Delete(channelID, messageID)
}

func (d *Client) removeMessagesFromCache(channelID common.Snowflake, messageIDs []common.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Messages().DeleteBulk(channelID, messageIDs)
}

func (d *Client) removeRoleFromCache(roleID common.Snowflake) {
	if d.Cache == nil {
		return
	}

	d.Cache.Roles().Delete(roleID)
}
