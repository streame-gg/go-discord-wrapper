package connection

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// applicationID returns the bot's application ID, or an error if the bot has
// not yet received the READY event and d.User is nil.
func (d *Client) applicationID() (*discord.Snowflake, error) {
	if d.User == nil {
		return nil, errors.New("application ID unavailable: bot is not yet ready")
	}
	return &d.User.ID, nil
}

// ── Application commands ──────────────────────────────────────────────────────

// RegisterCommand registers a single application command.
// Must be called after Login() so that the application ID is available.
func (d *Client) RegisterCommand(ctx context.Context, cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.RegisterCommand(ctx, *appID, cmd)
}

// BulkRegisterCommands overwrites all application commands with the provided list.
// Any existing commands not included will be deleted.
// Must be called after Login() so that the application ID is available.
func (d *Client) BulkRegisterCommands(ctx context.Context, cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.BulkRegisterCommands(ctx, *appID, cmds)
}

// ── Interaction responses ─────────────────────────────────────────────────────

// Reply sends an immediate message response to an interaction.
// Pass withResponse=true to get the created message back.
// Pass optional files to include attachments; when present the request is sent as multipart/form-data.
func (d *Client) Reply(ctx context.Context, i *interactions.Interaction, data *responses.InteractionResponseDataDefault, withResponse bool, files ...discord.MessageFile) (*responses.InteractionCallbackResponse, error) {
	return d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeChannelMessageWithSource,
		Data: data,
	}, withResponse, files...)
}

// DeferReply acknowledges the interaction, telling Discord the bot will respond later.
// Set ephemeral=true to make the eventual follow-up visible only to the invoker.
func (d *Client) DeferReply(ctx context.Context, i *interactions.Interaction, ephemeral bool) error {
	resp := responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeDeferredChannelMessageWithSource,
	}
	if ephemeral {
		resp.Data = &responses.InteractionResponseDataDefault{Flags: discord.MessageFlagEphemeral}
	}
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, resp, false)
	return err
}

// DeferUpdateMessage acknowledges a component interaction without editing the original message.
// Use EditReply afterward to push the actual update.
func (d *Client) DeferUpdateMessage(ctx context.Context, i *interactions.Interaction) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeDeferredUpdateMessage,
	}, false)
	return err
}

// UpdateMessage edits the message that triggered a component interaction.
func (d *Client) UpdateMessage(ctx context.Context, i *interactions.Interaction, data *responses.InteractionResponseDataDefault, files ...discord.MessageFile) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeUpdateMessage,
		Data: data,
	}, false, files...)
	return err
}

// ReplyWithModal responds to an interaction with a modal dialog.
func (d *Client) ReplyWithModal(ctx context.Context, i *interactions.Interaction, modal *components.Modal) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeModal,
		Data: modal,
	}, false)
	return err
}

// LaunchActivity responds to an interaction by launching the app's associated Activity.
func (d *Client) LaunchActivity(ctx context.Context, i *interactions.Interaction) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeLaunchActivity,
	}, false)
	return err
}

// ReplyAutocomplete sends autocomplete choices in response to an autocomplete interaction.
func (d *Client) ReplyAutocomplete(ctx context.Context, i *interactions.Interaction, choices []responses.AutocompleteChoice) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: discord.InteractionCallbackTypeApplicationCommandAutocompleteResult,
		Data: &responses.InteractionResponseDataAutocomplete{Choices: choices},
	}, false)
	return err
}

// GetOriginalResponse fetches the original message sent in response to an interaction.
func (d *Client) GetOriginalResponse(ctx context.Context, i *interactions.Interaction) (*discord.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.GetOriginalInteractionResponse(ctx, *appID, i.Token)
	if err == nil {
		msg.Hydrate(d)

		d.cacheMessage(msg)
	}
	return msg, err
}

// EditReply edits the original interaction response.
func (d *Client) EditReply(ctx context.Context, i *interactions.Interaction, params api.EditMessageParams) (*discord.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.EditOriginalInteractionResponse(ctx, *appID, i.Token, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// DeleteReply deletes the original interaction response.
func (d *Client) DeleteReply(ctx context.Context, i *interactions.Interaction) error {
	appID, err := d.applicationID()
	if err != nil {
		return err
	}
	return d.RestClient.DeleteOriginalInteractionResponse(ctx, *appID, i.Token)
}

// CreateFollowup sends a follow-up message to an interaction (up to 15 minutes after the initial response).
func (d *Client) CreateFollowup(ctx context.Context, i *interactions.Interaction, params api.CreateMessageParams) (*discord.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.CreateFollowupMessage(ctx, *appID, i.Token, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// GetFollowup fetches a follow-up message.
func (d *Client) GetFollowup(ctx context.Context, i *interactions.Interaction, messageID discord.Snowflake) (*discord.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.GetFollowupMessage(ctx, *appID, i.Token, messageID)
	if err == nil {
		msg.Hydrate(d)

		d.cacheMessage(msg)
	}
	return msg, err
}

// EditFollowup edits a follow-up message.
func (d *Client) EditFollowup(ctx context.Context, i *interactions.Interaction, messageID discord.Snowflake, params api.EditMessageParams) (*discord.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.EditFollowupMessage(ctx, *appID, i.Token, messageID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// DeleteFollowup deletes a follow-up message.
func (d *Client) DeleteFollowup(ctx context.Context, i *interactions.Interaction, messageID discord.Snowflake) error {
	appID, err := d.applicationID()
	if err != nil {
		return err
	}
	return d.RestClient.DeleteFollowupMessage(ctx, *appID, i.Token, messageID)
}

// ── Message methods ───────────────────────────────────────────────────────────

func (d *Client) GetMessages(ctx context.Context, channelID discord.Snowflake, params api.GetMessagesParams) ([]*discord.Message, error) {
	msgs, err := d.RestClient.GetMessages(ctx, channelID, params)
	if err == nil {
		for _, msg := range msgs {
			msg.Hydrate(d)
			d.cacheMessage(msg)
		}
	}
	return msgs, err
}

func (d *Client) GetMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	msg, err := d.RestClient.GetMessage(ctx, channelID, messageID)
	if err == nil {
		msg.Hydrate(d)
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) SendMessage(ctx context.Context, channelID discord.Snowflake, params api.CreateMessageParams) (*discord.Message, error) {
	msg, err := d.RestClient.CreateMessage(ctx, channelID, params)
	if err == nil {
		msg.Hydrate(d)
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) EditMessageRaw(ctx context.Context, channelID, messageID discord.Snowflake, params api.EditMessageParams) (*discord.Message, error) {
	msg, err := d.RestClient.EditMessage(ctx, channelID, messageID, params)
	if err == nil {
		msg.Hydrate(d)
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) DeleteMessage(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error {
	var opts *api.DeleteMessageOptions
	if reason != nil {
		opts = &api.DeleteMessageOptions{Reason: *reason}
	}
	if err := d.RestClient.DeleteMessage(ctx, channelID, messageID, opts); err != nil {
		return err
	}
	d.removeMessageFromCache(channelID, messageID)
	return nil
}

func (d *Client) BulkDeleteMessages(ctx context.Context, channelID discord.Snowflake, messageIDs []discord.Snowflake, reason *string) error {
	var opts *api.BulkDeleteMessagesOptions
	if reason != nil {
		opts = &api.BulkDeleteMessagesOptions{Reason: *reason}
	}
	if err := d.RestClient.BulkDeleteMessages(ctx, channelID, messageIDs, opts); err != nil {
		return err
	}
	d.removeMessagesFromCache(channelID, messageIDs)
	return nil
}

func (d *Client) CrosspostMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	msg, err := d.RestClient.CrosspostMessage(ctx, channelID, messageID)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) GetPinnedMessages(ctx context.Context, channelID discord.Snowflake) ([]*discord.Message, error) {
	msgs, err := d.RestClient.GetPinnedMessages(ctx, channelID)
	if err == nil {
		for _, msg := range msgs {
			msg.Hydrate(d)
			d.cacheMessage(msg)
		}
	}
	return msgs, err
}

func (d *Client) PinMessage(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error {
	var opts *api.PinMessageOptions
	if reason != nil {
		opts = &api.PinMessageOptions{Reason: *reason}
	}
	return d.RestClient.PinMessage(ctx, channelID, messageID, opts)
}

func (d *Client) UnpinMessage(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error {
	var opts *api.UnpinMessageOptions
	if reason != nil {
		opts = &api.UnpinMessageOptions{Reason: *reason}
	}
	return d.RestClient.UnpinMessage(ctx, channelID, messageID, opts)
}

func (d *Client) AddReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	return d.RestClient.AddReaction(ctx, channelID, messageID, emoji)
}

func (d *Client) DeleteOwnReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	return d.RestClient.DeleteOwnReaction(ctx, channelID, messageID, emoji)
}

func (d *Client) DeleteUserReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string, userID discord.Snowflake) error {
	return d.RestClient.DeleteUserReaction(ctx, channelID, messageID, emoji, userID)
}

func (d *Client) GetReactions(ctx context.Context, channelID, messageID discord.Snowflake, emoji string, params api.GetReactionsParams) ([]*discord.User, error) {
	users, err := d.RestClient.GetReactions(ctx, channelID, messageID, emoji, params)
	if err == nil {
		for _, user := range users {
			user.Hydrate(d)
			d.cacheUser(user)
		}
	}
	return users, err
}

func (d *Client) DeleteAllReactions(ctx context.Context, channelID, messageID discord.Snowflake) error {
	return d.RestClient.DeleteAllReactions(ctx, channelID, messageID)
}

func (d *Client) DeleteAllReactionsForEmoji(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	return d.RestClient.DeleteAllReactionsForEmoji(ctx, channelID, messageID, emoji)
}

// ── Channel methods ───────────────────────────────────────────────────────────

func (d *Client) GetChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	channel, err := d.RestClient.GetChannel(ctx, channelID)
	if err == nil {
		channel.Hydrate(d)
		d.cacheChannel(channel)
	}
	return channel, err
}

func (d *Client) ModifyChannelRaw(ctx context.Context, channelID discord.Snowflake, params api.ModifyChannelParams) (*discord.Channel, error) {
	channel, err := d.RestClient.ModifyChannel(ctx, channelID, params, nil)
	if err == nil {
		d.cacheChannel(channel)
	}
	return channel, err
}

func (d *Client) DeleteChannel(ctx context.Context, channelID discord.Snowflake, reason *string) (*discord.Channel, error) {
	var opts *api.DeleteChannelOptions
	if reason != nil {
		opts = &api.DeleteChannelOptions{Reason: *reason}
	}
	channel, err := d.RestClient.DeleteChannel(ctx, channelID, opts)
	if err == nil {
		d.removeChannelFromCache(channelID)
	}
	return channel, err
}

func (d *Client) GetChannelInvites(ctx context.Context, channelID discord.Snowflake) ([]*discord.Invite, error) {
	return d.RestClient.GetChannelInvites(ctx, channelID)
}

func (d *Client) CreateChannelInviteRaw(ctx context.Context, channelID discord.Snowflake, params api.CreateChannelInviteParams) (*discord.Invite, error) {
	return d.RestClient.CreateChannelInvite(ctx, channelID, params, nil)
}

func (d *Client) EditChannelPermissions(ctx context.Context, channelID, overwriteID discord.Snowflake, params api.EditChannelPermissionsParams) error {
	return d.RestClient.EditChannelPermissions(ctx, channelID, overwriteID, params, nil)
}

func (d *Client) DeleteChannelPermission(ctx context.Context, channelID, overwriteID discord.Snowflake) error {
	return d.RestClient.DeleteChannelPermission(ctx, channelID, overwriteID, nil)
}

func (d *Client) TriggerTypingIndicator(ctx context.Context, channelID discord.Snowflake) error {
	return d.RestClient.TriggerTypingIndicator(ctx, channelID)
}

// ── Guild methods ─────────────────────────────────────────────────────────────

func (d *Client) GetGuild(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	guild, err := d.RestClient.GetGuild(ctx, guildID, false)
	if err == nil {
		guild.Hydrate(d)
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) GetGuildPreview(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	guild, err := d.RestClient.GetGuildPreview(ctx, guildID)
	if err == nil {
		guild.Hydrate(d)
		d.cacheGuild(guild)
	}

	return guild, err
}

func (d *Client) ModifyGuildRaw(ctx context.Context, guildID discord.Snowflake, params api.ModifyGuildParams) (*discord.Guild, error) {
	guild, err := d.RestClient.ModifyGuild(ctx, guildID, params, nil)
	if err == nil {
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) DeleteGuild(ctx context.Context, guildID discord.Snowflake) error {
	if err := d.RestClient.DeleteGuild(ctx, guildID); err != nil {
		return err
	}
	d.removeGuildFromCache(guildID)
	return nil
}

func (d *Client) GetGuildChannels(ctx context.Context, guildID discord.Snowflake) ([]*discord.Channel, error) {
	channels, err := d.RestClient.GetGuildChannels(ctx, guildID)
	if err == nil {
		if d.Cache != nil {
			// Remove stale channels deleted since the last cache fill (Bug 25).
			for _, oldID := range d.drainGuildChannelIDs(guildID) {
				d.Cache.Channels().Delete(oldID)
				d.Cache.Messages().DeleteChannel(oldID)
			}
		}
		for _, channel := range channels {
			channel.Hydrate(d)
			d.cacheChannel(channel)
		}
	}
	return channels, err
}

func (d *Client) CreateGuildChannelRaw(ctx context.Context, guildID discord.Snowflake, params api.CreateGuildChannelParams) (*discord.Channel, error) {
	channel, err := d.RestClient.CreateGuildChannel(ctx, guildID, params, nil)
	if err == nil {
		d.cacheChannel(channel)
	}
	return channel, err
}

func (d *Client) ModifyGuildChannelPositions(ctx context.Context, guildID discord.Snowflake, entries []api.ModifyGuildChannelPositionsEntry) error {
	return d.RestClient.ModifyGuildChannelPositions(ctx, guildID, entries, nil)
}

func (d *Client) GetGuildRoles(ctx context.Context, guildID discord.Snowflake) ([]*discord.Role, error) {
	roles, err := d.RestClient.GetGuildRoles(ctx, guildID)
	if err == nil {
		if d.Cache != nil {
			// Remove stale roles deleted since the last cache fill (Bug 25).
			d.Cache.Roles().DeleteGuild(guildID)
		}

		for _, role := range roles {
			role.Hydrate(d)
			d.cacheRole(guildID, role)
		}
	}

	return roles, err
}

func (d *Client) GetGuildRole(ctx context.Context, guildID, roleID discord.Snowflake) (*discord.Role, error) {
	role, err := d.RestClient.GetGuildRole(ctx, guildID, roleID)
	if err == nil {
		role.Hydrate(d)
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) CreateGuildRoleRaw(ctx context.Context, guildID discord.Snowflake, params api.CreateGuildRoleParams) (*discord.Role, error) {
	role, err := d.RestClient.CreateGuildRole(ctx, guildID, params, nil)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) ModifyGuildRolePositionsRaw(ctx context.Context, guildID discord.Snowflake, entries []api.ModifyGuildRolePositionsEntry) ([]*discord.Role, error) {
	roles, err := d.RestClient.ModifyGuildRolePositions(ctx, guildID, entries, nil)
	if err == nil {
		if d.Cache != nil {
			// Remove stale roles deleted since the last cache fill (Bug 25).
			d.Cache.Roles().DeleteGuild(guildID)
		}
		d.cacheRoles(guildID, roles)
	}
	return roles, err
}

func (d *Client) ModifyGuildRoleRaw(ctx context.Context, guildID, roleID discord.Snowflake, params api.ModifyGuildRoleParams) (*discord.Role, error) {
	role, err := d.RestClient.ModifyGuildRole(ctx, guildID, roleID, params, nil)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) DeleteGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, reason *string) error {
	var opts *api.DeleteGuildRoleOptions
	if reason != nil {
		opts = &api.DeleteGuildRoleOptions{Reason: *reason}
	}
	if err := d.RestClient.DeleteGuildRole(ctx, guildID, roleID, opts); err != nil {
		return err
	}
	d.removeRoleFromCache(roleID)
	return nil
}

func (d *Client) GetGuildBans(ctx context.Context, guildID discord.Snowflake, opts discord.FetchBansOptions) ([]*discord.Ban, error) {
	bans, err := d.RestClient.GetGuildBans(ctx, guildID, api.GetGuildBansParams{
		Before: opts.Before,
		After:  opts.After,
		Limit:  opts.Limit,
	})
	if err == nil {
		for _, ban := range bans {
			ban.User.Hydrate(d)
		}
	}
	return bans, err
}

func (d *Client) GetGuildBan(ctx context.Context, guildID, userID discord.Snowflake) (*discord.Ban, error) {
	return d.RestClient.GetGuildBan(ctx, guildID, userID)
}

func (d *Client) CreateGuildBanRaw(ctx context.Context, guildID, userID discord.Snowflake, params api.CreateGuildBanParams) error {
	return d.RestClient.CreateGuildBan(ctx, guildID, userID, params, nil)
}

func (d *Client) RemoveGuildBan(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error {
	var opts *api.RemoveGuildBanOptions
	if reason != nil {
		opts = &api.RemoveGuildBanOptions{Reason: *reason}
	}
	return d.RestClient.RemoveGuildBan(ctx, guildID, userID, opts)
}

func (d *Client) GetGuildPruneCount(ctx context.Context, guildID discord.Snowflake, params api.GetGuildPruneCountParams) (*discord.GuildPruneCountResult, error) {
	return d.RestClient.GetGuildPruneCount(ctx, guildID, params)
}

func (d *Client) BeginGuildPrune(ctx context.Context, guildID discord.Snowflake, params api.BeginGuildPruneParams) (*discord.GuildPruneCountResult, error) {
	return d.RestClient.BeginGuildPrune(ctx, guildID, params, nil)
}

func (d *Client) GetGuildInvites(ctx context.Context, guildID discord.Snowflake) ([]*discord.Invite, error) {
	invites, err := d.RestClient.GetGuildInvites(ctx, guildID)
	if err == nil {
		for _, i := range invites {
			i.Hydrate(d)
		}
	}
	return invites, err
}

func (d *Client) GetGuildVanityURL(ctx context.Context, guildID discord.Snowflake) (*discord.GuildVanityURL, error) {
	return d.RestClient.GetGuildVanityURL(ctx, guildID)
}

func (d *Client) GetGuildAuditLogRaw(ctx context.Context, guildID discord.Snowflake, params api.GetGuildAuditLogParams) (*discord.AuditLog, error) {
	return d.RestClient.GetGuildAuditLog(ctx, guildID, params)
}

// GetGuildRoleMemberCounts returns a map of role ID → member count for every role in the guild.
func (d *Client) GetGuildRoleMemberCounts(ctx context.Context, guildID discord.Snowflake) (*map[string]int64, error) {
	return d.RestClient.GetGuildRoleMemberCounts(ctx, guildID)
}

// ── Member methods ────────────────────────────────────────────────────────────

func (d *Client) GetGuildMember(ctx context.Context, guildID, userID discord.Snowflake) (*discord.GuildMember, error) {
	member, err := d.RestClient.GetGuildMember(ctx, guildID, userID)
	if err == nil {
		member.Hydrate(d)
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) ListGuildMembersRaw(ctx context.Context, guildID discord.Snowflake, params api.GetGuildMembersParams) ([]*discord.GuildMember, error) {
	members, err := d.RestClient.ListGuildMembers(ctx, guildID, params)
	if err == nil {
		for _, member := range members {
			member.Hydrate(d)
			d.cacheMember(guildID, member)
		}
	}
	return members, err
}

func (d *Client) SearchGuildMembers(ctx context.Context, guildID discord.Snowflake, params api.SearchGuildMembersParams) ([]*discord.GuildMember, error) {
	members, err := d.RestClient.SearchGuildMembers(ctx, guildID, params)
	if err == nil {
		for _, member := range members {
			member.Hydrate(d)
			d.cacheMember(guildID, member)
		}
	}
	return members, err
}

func (d *Client) ModifyGuildMemberRaw(ctx context.Context, guildID, userID discord.Snowflake, params api.ModifyGuildMemberParams) (*discord.GuildMember, error) {
	member, err := d.RestClient.ModifyGuildMember(ctx, guildID, userID, params, nil)
	if err == nil {
		member.Hydrate(d)
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) ModifyCurrentMember(ctx context.Context, guildID discord.Snowflake, params api.ModifyCurrentMemberParams) (*discord.GuildMember, error) {
	member, err := d.RestClient.ModifyCurrentMember(ctx, guildID, params)
	if err == nil {
		member.Hydrate(d)
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) AddGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, reason *string) error {
	var opts *api.AddGuildMemberRoleOptions
	if reason != nil {
		opts = &api.AddGuildMemberRoleOptions{Reason: *reason}
	}
	return d.RestClient.AddGuildMemberRole(ctx, guildID, userID, roleID, opts)
}

func (d *Client) RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, reason *string) error {
	var opts *api.RemoveGuildMemberRoleOptions
	if reason != nil {
		opts = &api.RemoveGuildMemberRoleOptions{Reason: *reason}
	}
	return d.RestClient.RemoveGuildMemberRole(ctx, guildID, userID, roleID, opts)
}

// RemoveGuildMember removes a member from a guild. Requires KICK_MEMBERS.
func (d *Client) RemoveGuildMember(ctx context.Context, guildID, userID discord.Snowflake) error {
	if err := d.RestClient.KickGuildMember(ctx, guildID, userID, nil); err != nil {
		return err
	}
	d.removeGuildMemberFromCache(guildID, userID)
	return nil
}

// FetchAllGuildMembers retrieves every member in a guild by paginating the REST
// API with 1 000-member pages until all members are returned.
//
// Requires the GUILD_MEMBERS privileged intent to be enabled both in the Discord
// developer portal and in the intents passed to NewClient.
func (d *Client) FetchAllGuildMembers(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildMember, error) {
	members, err := d.RestClient.FetchAllGuildMembers(ctx, guildID)
	if err == nil {
		for _, member := range members {
			member.Hydrate(d)
			d.cacheMember(guildID, member)
		}
	}
	return members, err
}

// RequestGuildMembersParams controls what the OP 8 Request Guild Members
// gateway command fetches.
type RequestGuildMembersParams struct {
	// Query is a username prefix. Use "" to request all members (requires the
	// GUILD_MEMBERS privileged intent).
	Query string
	// Limit caps the number of members returned per guild. 0 = no limit.
	Limit int
	// Presences includes presence objects in the GUILD_MEMBERS_CHUNK events
	// (requires the GUILD_PRESENCES intent).
	Presences bool
	// UserIDs fetches specific members by user ID. Overrides Query when set.
	UserIDs []discord.Snowflake
	// Nonce is echoed back in every GUILD_MEMBERS_CHUNK event, useful for
	// correlating a request to its responses.
	Nonce *string
}

// RequestGuildMembers sends an OP 8 Request Guild Members command to the
// Discord gateway. Discord responds with one or more GUILD_MEMBERS_CHUNK
// events that can be received via OnGuildMembersChunk.
//
// Use FetchAllGuildMembers for a simpler REST-based alternative that does not
// require subscribing to gateway events.
func (d *Client) RequestGuildMembers(guildID discord.Snowflake, params RequestGuildMembersParams) error {
	d.wsMu.RLock()
	ws := d.Websocket
	d.wsMu.RUnlock()
	if ws == nil || ws.Connection == nil {
		return errors.New("not connected to gateway")
	}
	data := map[string]interface{}{
		"guild_id": guildID,
		"query":    params.Query,
		"limit":    params.Limit,
	}
	if params.Presences {
		data["presences"] = true
	}
	if len(params.UserIDs) > 0 {
		data["user_ids"] = params.UserIDs
	}
	if params.Nonce != nil {
		data["nonce"] = *params.Nonce
	}
	return ws.writeJSON(map[string]interface{}{
		"op": 8,
		"d":  data,
	})
}

// UpdatePresenceParams controls the bot's displayed presence in Discord.
type UpdatePresenceParams struct {
	// Since is the unix time (in milliseconds) of when the client went idle.
	// Set to nil if the client is not idle.
	Since *int64
	// Activities is the list of activities the bot is performing (e.g. "Playing X").
	// Pass nil or an empty slice to clear all activities.
	Activities []discord.FullActivity
	// Status is the bot's status ("online", "dnd", "idle", "invisible", "offline").
	Status discord.PresenceStatus
	// AFK indicates whether the bot is AFK.
	AFK bool
}

// UpdatePresence sends an OP 3 Presence Update to the Discord gateway,
// updating the bot's displayed status and activity.
func (d *Client) UpdatePresence(params UpdatePresenceParams) error {
	d.wsMu.RLock()
	ws := d.Websocket
	d.wsMu.RUnlock()
	if ws == nil || ws.Connection == nil {
		return errors.New("not connected to gateway")
	}

	activities := params.Activities
	if activities == nil {
		activities = []discord.FullActivity{}
	}

	data := map[string]interface{}{
		"since":      params.Since,
		"activities": activities,
		"status":     string(params.Status),
		"afk":        params.AFK,
	}
	return ws.writeJSON(map[string]interface{}{
		"op": 3,
		"d":  data,
	})
}

// ── Guild widget ──────────────────────────────────────────────────────────────

func (d *Client) GetGuildWidgetSettings(ctx context.Context, guildID discord.Snowflake) (*discord.GuildWidgetSettings, error) {
	return d.RestClient.GetGuildWidgetSettings(ctx, guildID)
}

func (d *Client) ModifyGuildWidgetSettings(ctx context.Context, guildID discord.Snowflake, params api.ModifyGuildWidgetParams) (*discord.GuildWidgetSettings, error) {
	return d.RestClient.ModifyGuildWidgetSettings(ctx, guildID, params)
}

func (d *Client) GetGuildWidget(ctx context.Context, guildID discord.Snowflake) (*discord.GuildWidget, error) {
	return d.RestClient.GetGuildWidget(ctx, guildID)
}

// GetGuildWidgetImage returns the PNG widget image for a guild as raw bytes.
// style is an optional widget style ("shield", "banner1"–"banner4"); pass "" for the default.
func (d *Client) GetGuildWidgetImage(ctx context.Context, guildID discord.Snowflake, style string) ([]byte, error) {
	return d.RestClient.GetGuildWidgetImage(ctx, guildID, style)
}

// ── Welcome screen ────────────────────────────────────────────────────────────

func (d *Client) GetGuildWelcomeScreen(ctx context.Context, guildID discord.Snowflake) (*discord.GuildWelcomeScreen, error) {
	return d.RestClient.GetGuildWelcomeScreen(ctx, guildID)
}

func (d *Client) ModifyGuildWelcomeScreen(ctx context.Context, guildID discord.Snowflake, params api.ModifyGuildWelcomeScreenParams) (*discord.GuildWelcomeScreen, error) {
	return d.RestClient.ModifyGuildWelcomeScreen(ctx, guildID, params, nil)
}

// ── Guild onboarding ──────────────────────────────────────────────────────────

func (d *Client) GetGuildOnboarding(ctx context.Context, guildID discord.Snowflake) (*discord.GuildOnboarding, error) {
	return d.RestClient.GetGuildOnboarding(ctx, guildID)
}

func (d *Client) ModifyGuildOnboarding(ctx context.Context, guildID discord.Snowflake, params api.ModifyGuildOnboardingParams) (*discord.GuildOnboarding, error) {
	return d.RestClient.ModifyGuildOnboarding(ctx, guildID, params, nil)
}

// ── Voice ─────────────────────────────────────────────────────────────────────

func (d *Client) ListVoiceRegions(ctx context.Context) ([]*discord.VoiceRegion, error) {
	return d.RestClient.ListVoiceRegions(ctx)
}

func (d *Client) ListGuildVoiceRegions(ctx context.Context, guildID discord.Snowflake) ([]*discord.VoiceRegion, error) {
	return d.RestClient.ListGuildVoiceRegions(ctx, guildID)
}

func (d *Client) ModifyCurrentUserVoiceState(ctx context.Context, guildID discord.Snowflake, params api.ModifyCurrentUserVoiceStateParams) error {
	return d.RestClient.ModifyCurrentUserVoiceState(ctx, guildID, params)
}

func (d *Client) ModifyUserVoiceState(ctx context.Context, guildID, userID discord.Snowflake, params api.ModifyUserVoiceStateParams) error {
	return d.RestClient.ModifyUserVoiceState(ctx, guildID, userID, params)
}

// ── Soundboard ────────────────────────────────────────────────────────────────

func (d *Client) ListDefaultSoundboardSounds(ctx context.Context) ([]*discord.SoundboardSound, error) {
	return d.RestClient.ListDefaultSoundboardSounds(ctx)
}

func (d *Client) ListGuildSoundboardSounds(ctx context.Context, guildID discord.Snowflake) ([]*discord.SoundboardSound, error) {
	resp, err := d.RestClient.ListGuildSoundboardSounds(ctx, guildID)
	if err != nil {
		return nil, err
	}
	sounds := resp.Items
	for _, s := range sounds {
		s.Hydrate(d)
	}
	if d.Cache != nil {
		d.Cache.Soundboard().SetAll(guildID, sounds)
	}
	return sounds, nil
}

func (d *Client) GetGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake) (*discord.SoundboardSound, error) {
	soundboardSounds, err := d.RestClient.GetGuildSoundboardSound(ctx, guildID, soundID)
	if err == nil {
		soundboardSounds.Hydrate(d)
		soundboardSounds.GuildID = &guildID

		if d.Cache != nil {
			d.Cache.Soundboard().Set(guildID, soundboardSounds)
		}
	}

	return soundboardSounds, err
}

func (d *Client) CreateGuildSoundboardSound(ctx context.Context, guildID discord.Snowflake, params api.CreateGuildSoundboardSoundParams) (*discord.SoundboardSound, error) {
	return d.RestClient.CreateGuildSoundboardSound(ctx, guildID, params, nil)
}

func (d *Client) ModifyGuildSoundboardSoundRaw(ctx context.Context, guildID, soundID discord.Snowflake, params api.ModifyGuildSoundboardSoundParams) (*discord.SoundboardSound, error) {
	return d.RestClient.ModifyGuildSoundboardSound(ctx, guildID, soundID, params, nil)
}

func (d *Client) DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, reason *string) error {
	var opts *api.DeleteGuildSoundboardSoundOptions
	if reason != nil {
		opts = &api.DeleteGuildSoundboardSoundOptions{Reason: *reason}
	}
	return d.RestClient.DeleteGuildSoundboardSound(ctx, guildID, soundID, opts)
}

func (d *Client) SendSoundboardSound(ctx context.Context, channelID discord.Snowflake, params api.SendSoundboardSoundParams) error {
	return d.RestClient.SendSoundboardSound(ctx, channelID, params)
}

// ── Application ───────────────────────────────────────────────────────────────

func (d *Client) GetCurrentApplication(ctx context.Context) (*discord.Application, error) {
	return d.RestClient.GetCurrentApplication(ctx)
}

func (d *Client) ModifyCurrentApplication(ctx context.Context, params api.ModifyCurrentApplicationParams) (*discord.Application, error) {
	return d.RestClient.ModifyCurrentApplication(ctx, params)
}

// ── Command permissions ───────────────────────────────────────────────────────

// GetGuildApplicationCommandPermissions returns all permission overrides for every command in a guild.
// Uses the bot's own application ID automatically.
func (d *Client) GetGuildApplicationCommandPermissions(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildApplicationCommandPermissions, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.GetGuildApplicationCommandPermissions(ctx, *appID, guildID)
}

// GetApplicationCommandPermissions returns the permission overrides for a specific command.
func (d *Client) GetApplicationCommandPermissions(ctx context.Context, guildID, cmdID discord.Snowflake) (*discord.GuildApplicationCommandPermissions, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.GetApplicationCommandPermissions(ctx, *appID, guildID, cmdID)
}

// ── Application emojis ────────────────────────────────────────────────────────

func (d *Client) ListApplicationEmojis(ctx context.Context, appID discord.Snowflake) (*api.ListApplicationEmojisResponse, error) {
	return d.RestClient.ListApplicationEmojis(ctx, appID)
}

func (d *Client) GetApplicationEmoji(ctx context.Context, appID, emojiID discord.Snowflake) (*discord.Emoji, error) {
	return d.RestClient.GetApplicationEmoji(ctx, appID, emojiID)
}

func (d *Client) CreateApplicationEmoji(ctx context.Context, appID discord.Snowflake, params api.CreateEmojiParams) (*discord.Emoji, error) {
	return d.RestClient.CreateApplicationEmoji(ctx, appID, params)
}

func (d *Client) ModifyApplicationEmoji(ctx context.Context, appID, emojiID discord.Snowflake, params api.ModifyEmojiParams) (*discord.Emoji, error) {
	return d.RestClient.ModifyApplicationEmoji(ctx, appID, emojiID, params)
}

func (d *Client) DeleteApplicationEmoji(ctx context.Context, appID, emojiID discord.Snowflake) error {
	return d.RestClient.DeleteApplicationEmoji(ctx, appID, emojiID)
}

// ── Entitlements ──────────────────────────────────────────────────────────────

func (d *Client) ListEntitlements(ctx context.Context, appID discord.Snowflake, params api.ListEntitlementsParams) ([]*discord.Entitlement, error) {
	return d.RestClient.ListEntitlements(ctx, appID, params)
}

func (d *Client) GetEntitlement(ctx context.Context, appID, entitlementID discord.Snowflake) (*discord.Entitlement, error) {
	return d.RestClient.GetEntitlement(ctx, appID, entitlementID)
}

func (d *Client) CreateTestEntitlement(ctx context.Context, appID discord.Snowflake, params api.CreateTestEntitlementParams) (*discord.Entitlement, error) {
	return d.RestClient.CreateTestEntitlement(ctx, appID, params)
}

func (d *Client) ConsumeEntitlement(ctx context.Context, appID, entitlementID discord.Snowflake) error {
	return d.RestClient.ConsumeEntitlement(ctx, appID, entitlementID)
}

func (d *Client) DeleteTestEntitlement(ctx context.Context, appID, entitlementID discord.Snowflake) error {
	return d.RestClient.DeleteTestEntitlement(ctx, appID, entitlementID)
}

// ── SKUs ──────────────────────────────────────────────────────────────────────

func (d *Client) ListSKUs(ctx context.Context, appID discord.Snowflake) ([]*discord.SKU, error) {
	return d.RestClient.ListSKUs(ctx, appID)
}

// ── Subscriptions ─────────────────────────────────────────────────────────────

func (d *Client) ListSKUSubscriptions(ctx context.Context, skuID discord.Snowflake, params api.ListSKUSubscriptionsParams) ([]*discord.Subscription, error) {
	return d.RestClient.ListSKUSubscriptions(ctx, skuID, params)
}

func (d *Client) GetSKUSubscription(ctx context.Context, skuID, subscriptionID discord.Snowflake) (*discord.Subscription, error) {
	return d.RestClient.GetSKUSubscription(ctx, skuID, subscriptionID)
}

// ── Poll ──────────────────────────────────────────────────────────────────────

func (d *Client) GetPollAnswerVoters(ctx context.Context, channelID, messageID discord.Snowflake, answerID int, params api.GetPollAnswerVotersParams) (*api.GetPollAnswerVotersResponse, error) {
	users, err := d.RestClient.GetPollAnswerVoters(ctx, channelID, messageID, answerID, params)
	if err == nil {
		for _, user := range users.Users {
			user.Hydrate(d)
			d.cacheUser(user)
		}
	}
	return users, err
}

func (d *Client) EndPoll(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	msg, err := d.RestClient.EndPoll(ctx, channelID, messageID)
	if err == nil {
		msg.Hydrate(d)
		d.cacheMessage(msg)
	}
	return msg, err
}

// ── Activity ──────────────────────────────────────────────────────────────────

func (d *Client) GetActivityInstance(ctx context.Context, appID discord.Snowflake, instanceID string) (*discord.ActivityInstance, error) {
	return d.RestClient.GetActivityInstance(ctx, appID, instanceID)
}

// ── Additional channel methods ────────────────────────────────────────────────

func (d *Client) SetVoiceChannelStatus(ctx context.Context, channelID discord.Snowflake, status *string) error {
	return d.RestClient.SetVoiceChannelStatus(ctx, channelID, status)
}

func (d *Client) FollowAnnouncementChannel(ctx context.Context, channelID, webhookChannelID discord.Snowflake) (*discord.FollowedChannel, error) {
	return d.RestClient.FollowAnnouncementChannel(ctx, channelID, webhookChannelID, nil)
}

func (d *Client) AddGroupDMRecipient(ctx context.Context, channelID, userID discord.Snowflake, params api.AddGroupDMRecipientParams) error {
	return d.RestClient.AddGroupDMRecipient(ctx, channelID, userID, params)
}

func (d *Client) RemoveGroupDMRecipient(ctx context.Context, channelID, userID discord.Snowflake) error {
	return d.RestClient.RemoveGroupDMRecipient(ctx, channelID, userID)
}

// ── Additional guild methods ──────────────────────────────────────────────────

func (d *Client) CreateGuild(ctx context.Context, params api.CreateGuildParams) (*discord.Guild, error) {
	guild, err := d.RestClient.CreateGuild(ctx, params)
	if err == nil {
		guild.Hydrate(d)
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) AddGuildMember(ctx context.Context, guildID, userID discord.Snowflake, params api.AddGuildMemberParams) (*discord.GuildMember, error) {
	member, err := d.RestClient.AddGuildMember(ctx, guildID, userID, params)
	if err == nil {
		member.Hydrate(d)
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) BulkBanGuildMembers(ctx context.Context, guildID discord.Snowflake, params api.BulkBanParams) (*api.BulkBanResult, error) {
	return d.RestClient.BulkBanGuildMembers(ctx, guildID, params, nil)
}

func (d *Client) GetGuildIntegrations(ctx context.Context, guildID discord.Snowflake) ([]*discord.Integration, error) {
	integrations, err := d.RestClient.GetGuildIntegrations(ctx, guildID)
	if err == nil {
		for _, i := range integrations {
			i.GuildID = guildID
			i.Hydrate(d)
		}
	}
	return integrations, err
}

func (d *Client) DeleteGuildIntegration(ctx context.Context, guildID, integrationID discord.Snowflake, reason *string) error {
	var opts *api.DeleteGuildIntegrationOptions
	if reason != nil {
		opts = &api.DeleteGuildIntegrationOptions{Reason: *reason}
	}
	return d.RestClient.DeleteGuildIntegration(ctx, guildID, integrationID, opts)
}

// ── Additional voice methods ──────────────────────────────────────────────────

func (d *Client) GetCurrentUserVoiceState(ctx context.Context, guildID discord.Snowflake) (*discord.VoiceState, error) {
	state, err := d.RestClient.GetCurrentUserVoiceState(ctx, guildID)
	if err == nil {
		if d.Cache != nil {
			d.Cache.VoiceStates().Set(guildID, state)
		}
	}
	return state, err
}

func (d *Client) GetUserVoiceState(ctx context.Context, guildID, userID discord.Snowflake) (*discord.VoiceState, error) {
	state, err := d.RestClient.GetUserVoiceState(ctx, guildID, userID)
	if err == nil {
		if d.Cache != nil {
			d.Cache.VoiceStates().Set(guildID, state)
		}
	}
	return state, err
}

// ── Additional sticker methods ────────────────────────────────────────────────

func (d *Client) GetStickerPack(ctx context.Context, packID discord.Snowflake) (*discord.StickerPack, error) {
	return d.RestClient.GetStickerPack(ctx, packID)
}

func (d *Client) CreateGuildStickerRaw(ctx context.Context, guildID discord.Snowflake, params api.CreateGuildStickerParams) (*discord.Sticker, error) {
	return d.RestClient.CreateGuildSticker(ctx, guildID, params, nil)
}

// ── Invite methods ────────────────────────────────────────────────────────────

func (d *Client) GetInvite(ctx context.Context, code string, params api.GetInviteParams) (*discord.Invite, error) {
	invite, err := d.RestClient.GetInvite(ctx, code, params)
	if err == nil {
		invite.Hydrate(d)
	}
	return invite, err
}

func (d *Client) DeleteInvite(ctx context.Context, code string, reason *string) (*discord.Invite, error) {
	var opts *api.DeleteInviteOptions
	if reason != nil {
		opts = &api.DeleteInviteOptions{Reason: *reason}
	}
	return d.RestClient.DeleteInvite(ctx, code, opts)
}

// ── Additional webhook methods ────────────────────────────────────────────────

func (d *Client) ExecuteSlackWebhook(ctx context.Context, webhookID discord.Snowflake, token string, wait bool, body json.RawMessage) error {
	return d.RestClient.ExecuteSlackWebhook(ctx, webhookID, token, wait, body)
}

func (d *Client) ExecuteGitHubWebhook(ctx context.Context, webhookID discord.Snowflake, token string, wait bool, body json.RawMessage) error {
	return d.RestClient.ExecuteGitHubWebhook(ctx, webhookID, token, wait, body)
}
