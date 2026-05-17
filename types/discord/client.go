package discord

import "github.com/streame-gg/go-discord-wrapper/collection"

// RestClient is the subset of *connection.Client needed by Guild, Channel, and
// Message methods. Pass the client received in your event handler directly —
// *connection.Client satisfies this interface.
//
// Operations that require complex parameter structs (e.g. CreateGuildBan,
// ModifyChannel) are intentionally omitted; call those directly on the client.
type RestClient interface {
	// Guild operations
	GetGuildChannels(guildID Snowflake) (*collection.Collection[Snowflake, *Channel], error)
	GetGuildRoles(guildID Snowflake) (*collection.Collection[Snowflake, *Role], error)
	GetGuildRoleMemberCounts(guildID Snowflake) (map[string]int, error)
	DeleteGuild(guildID Snowflake) error
	RemoveGuildBan(guildID, userID Snowflake) error
	RemoveGuildMember(guildID, userID Snowflake) error
	AddGuildMemberRole(guildID, userID, roleID Snowflake) error
	RemoveGuildMemberRole(guildID, userID, roleID Snowflake) error

	// Channel operations
	DeleteChannel(channelID Snowflake) (*Channel, error)
	TriggerTypingIndicator(channelID Snowflake) error
	GetPinnedMessages(channelID Snowflake) (*collection.Collection[Snowflake, *Message], error)

	// Message operations
	DeleteMessage(channelID, messageID Snowflake) error
	PinMessage(channelID, messageID Snowflake) error
	UnpinMessage(channelID, messageID Snowflake) error
	AddReaction(channelID, messageID Snowflake, emoji string) error
	DeleteOwnReaction(channelID, messageID Snowflake, emoji string) error
	DeleteAllReactions(channelID, messageID Snowflake) error
	CrosspostMessage(channelID, messageID Snowflake) (*Message, error)
}

// ── Guild methods ─────────────────────────────────────────────────────────────

// GetChannels returns all channels in the guild, keyed by channel ID.
func (g *Guild) GetChannels(client RestClient) (*collection.Collection[Snowflake, *Channel], error) {
	return client.GetGuildChannels(g.ID)
}

// GetRoles returns all roles in the guild, keyed by role ID.
func (g *Guild) GetRoles(client RestClient) (*collection.Collection[Snowflake, *Role], error) {
	return client.GetGuildRoles(g.ID)
}

// Delete deletes the guild. The bot must be the owner.
func (g *Guild) Delete(client RestClient) error {
	return client.DeleteGuild(g.ID)
}

// Unban removes a ban for the given user.
func (g *Guild) Unban(client RestClient, userID Snowflake) error {
	return client.RemoveGuildBan(g.ID, userID)
}

// Kick removes a member from the guild. Requires KICK_MEMBERS.
func (g *Guild) Kick(client RestClient, userID Snowflake) error {
	return client.RemoveGuildMember(g.ID, userID)
}

// AddMemberRole assigns a role to a guild member. Requires MANAGE_ROLES.
func (g *Guild) AddMemberRole(client RestClient, userID, roleID Snowflake) error {
	return client.AddGuildMemberRole(g.ID, userID, roleID)
}

// RemoveMemberRole removes a role from a guild member. Requires MANAGE_ROLES.
func (g *Guild) RemoveMemberRole(client RestClient, userID, roleID Snowflake) error {
	return client.RemoveGuildMemberRole(g.ID, userID, roleID)
}

func (g *Guild) GetGuildRoleMemberCounts(client RestClient) (map[string]int, error) {
	return client.GetGuildRoleMemberCounts(g.ID)
}

// ── Channel methods ───────────────────────────────────────────────────────────

// Delete deletes the channel. Returns the deleted channel object.
func (c *Channel) Delete(client RestClient) (*Channel, error) {
	return client.DeleteChannel(c.ID)
}

// TriggerTyping posts a typing indicator for ~10 seconds.
func (c *Channel) TriggerTyping(client RestClient) error {
	return client.TriggerTypingIndicator(c.ID)
}

// GetPinnedMessages returns all pinned messages in the channel, keyed by message ID.
func (c *Channel) GetPinnedMessages(client RestClient) (*collection.Collection[Snowflake, *Message], error) {
	return client.GetPinnedMessages(c.ID)
}

// ── Message methods ───────────────────────────────────────────────────────────

// Delete deletes the message.
func (m *Message) Delete(client RestClient) error {
	return client.DeleteMessage(m.ChannelID, m.ID)
}

// Pin pins the message in its channel.
func (m *Message) Pin(client RestClient) error {
	return client.PinMessage(m.ChannelID, m.ID)
}

// Unpin removes the message from the channel's pins.
func (m *Message) Unpin(client RestClient) error {
	return client.UnpinMessage(m.ChannelID, m.ID)
}

// AddReaction adds a reaction to the message. emoji should be a Unicode emoji
// or "name:id" for a custom emoji.
func (m *Message) AddReaction(client RestClient, emoji string) error {
	return client.AddReaction(m.ChannelID, m.ID, emoji)
}

// DeleteOwnReaction removes the bot's reaction from the message.
func (m *Message) DeleteOwnReaction(client RestClient, emoji string) error {
	return client.DeleteOwnReaction(m.ChannelID, m.ID, emoji)
}

// DeleteAllReactions removes every reaction from the message. Requires MANAGE_MESSAGES.
func (m *Message) DeleteAllReactions(client RestClient) error {
	return client.DeleteAllReactions(m.ChannelID, m.ID)
}

// Crosspost publishes the message to follower channels. Only valid in announcement channels.
func (m *Message) Crosspost(client RestClient) (*Message, error) {
	return client.CrosspostMessage(m.ChannelID, m.ID)
}
