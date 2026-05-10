package connection

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// applicationID returns the bot's application ID, or an error if the bot has
// not yet received the READY event and d.User is nil.
func (d *Client) applicationID() (common.Snowflake, error) {
	if d.User == nil {
		return "", errors.New("application ID unavailable: bot is not yet ready")
	}
	return d.User.ID, nil
}

// ── Application commands ──────────────────────────────────────────────────────

// RegisterCommand registers a single application command.
// Must be called after Login() so that the application ID is available.
func (d *Client) RegisterCommand(ctx context.Context, cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.RegisterCommand(ctx, appID, cmd)
}

// BulkRegisterCommands overwrites all application commands with the provided list.
// Any existing commands not included will be deleted.
// Must be called after Login() so that the application ID is available.
func (d *Client) BulkRegisterCommands(ctx context.Context, cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.BulkRegisterCommands(ctx, appID, cmds)
}

// ── Interaction responses ─────────────────────────────────────────────────────

// Reply sends an immediate message response to an interaction.
// Pass withResponse=true to get the created message back.
// Pass optional files to include attachments; when present the request is sent as multipart/form-data.
func (d *Client) Reply(ctx context.Context, i *interactions.Interaction, data *responses.InteractionResponseDataDefault, withResponse bool, files ...api.MessageFile) (*responses.InteractionCallbackResponse, error) {
	return d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeChannelMessageWithSource,
		Data: data,
	}, withResponse, files...)
}

// DeferReply acknowledges the interaction, telling Discord the bot will respond later.
// Set ephemeral=true to make the eventual follow-up visible only to the invoker.
func (d *Client) DeferReply(ctx context.Context, i *interactions.Interaction, ephemeral bool) error {
	var data *responses.InteractionResponseDataDefault
	if ephemeral {
		data = &responses.InteractionResponseDataDefault{Flags: common.MessageFlagEphemeral}
	}
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeDeferredChannelMessageWithSource,
		Data: data,
	}, false)
	return err
}

// DeferUpdateMessage acknowledges a component interaction without editing the original message.
// Use EditReply afterwards to push the actual update.
func (d *Client) DeferUpdateMessage(ctx context.Context, i *interactions.Interaction) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeDeferredUpdateMessage,
	}, false)
	return err
}

// UpdateMessage edits the message that triggered a component interaction.
func (d *Client) UpdateMessage(ctx context.Context, i *interactions.Interaction, data *responses.InteractionResponseDataDefault) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeUpdateMessage,
		Data: data,
	}, false)
	return err
}

// ReplyWithModal responds to an interaction with a modal dialog.
func (d *Client) ReplyWithModal(ctx context.Context, i *interactions.Interaction, modal *components.Modal) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeModal,
		Data: modal,
	}, false)
	return err
}

// LaunchActivity responds to an interaction by launching the app's associated Activity.
func (d *Client) LaunchActivity(ctx context.Context, i *interactions.Interaction) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeLaunchActivity,
	}, false)
	return err
}

// ReplyAutocomplete sends autocomplete choices in response to an autocomplete interaction.
func (d *Client) ReplyAutocomplete(ctx context.Context, i *interactions.Interaction, choices []responses.AutocompleteChoice) error {
	_, err := d.RestClient.CreateInteractionResponse(ctx, i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeApplicationCommandAutocompleteResult,
		Data: &responses.InteractionResponseDataAutocomplete{Choices: choices},
	}, false)
	return err
}

// GetOriginalResponse fetches the original message sent in response to an interaction.
func (d *Client) GetOriginalResponse(ctx context.Context, i *interactions.Interaction) (*common.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.GetOriginalInteractionResponse(ctx, appID, i.Token)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// EditReply edits the original interaction response.
func (d *Client) EditReply(ctx context.Context, i *interactions.Interaction, params api.EditMessageParams) (*common.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.EditOriginalInteractionResponse(ctx, appID, i.Token, params)
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
	return d.RestClient.DeleteOriginalInteractionResponse(ctx, appID, i.Token)
}

// CreateFollowup sends a follow-up message to an interaction (up to 15 minutes after the initial response).
func (d *Client) CreateFollowup(ctx context.Context, i *interactions.Interaction, params api.CreateMessageParams) (*common.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.CreateFollowupMessage(ctx, appID, i.Token, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// GetFollowup fetches a follow-up message.
func (d *Client) GetFollowup(ctx context.Context, i *interactions.Interaction, messageID common.Snowflake) (*common.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.GetFollowupMessage(ctx, appID, i.Token, messageID)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// EditFollowup edits a follow-up message.
func (d *Client) EditFollowup(ctx context.Context, i *interactions.Interaction, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	msg, err := d.RestClient.EditFollowupMessage(ctx, appID, i.Token, messageID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// DeleteFollowup deletes a follow-up message.
func (d *Client) DeleteFollowup(ctx context.Context, i *interactions.Interaction, messageID common.Snowflake) error {
	appID, err := d.applicationID()
	if err != nil {
		return err
	}
	return d.RestClient.DeleteFollowupMessage(ctx, appID, i.Token, messageID)
}

// ── Message methods ───────────────────────────────────────────────────────────

func (d *Client) GetMessages(ctx context.Context, channelID common.Snowflake, params api.GetMessagesParams) ([]*common.Message, error) {
	msgs, err := d.RestClient.GetMessages(ctx, channelID, params)
	if err == nil {
		d.cacheMessages(msgs)
	}
	return msgs, err
}

func (d *Client) GetMessage(ctx context.Context, channelID, messageID common.Snowflake) (*common.Message, error) {
	msg, err := d.RestClient.GetMessage(ctx, channelID, messageID)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) SendMessage(ctx context.Context, channelID common.Snowflake, params api.CreateMessageParams) (*common.Message, error) {
	msg, err := d.RestClient.CreateMessage(ctx, channelID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) EditMessage(ctx context.Context, channelID, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	msg, err := d.RestClient.EditMessage(ctx, channelID, messageID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) DeleteMessage(ctx context.Context, channelID, messageID common.Snowflake) error {
	if err := d.RestClient.DeleteMessage(ctx, channelID, messageID); err != nil {
		return err
	}
	d.removeMessageFromCache(channelID, messageID)
	return nil
}

func (d *Client) BulkDeleteMessages(ctx context.Context, channelID common.Snowflake, messageIDs []common.Snowflake) error {
	if err := d.RestClient.BulkDeleteMessages(ctx, channelID, messageIDs); err != nil {
		return err
	}
	d.removeMessagesFromCache(channelID, messageIDs)
	return nil
}

func (d *Client) CrosspostMessage(ctx context.Context, channelID, messageID common.Snowflake) (*common.Message, error) {
	msg, err := d.RestClient.CrosspostMessage(ctx, channelID, messageID)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) GetPinnedMessages(ctx context.Context, channelID common.Snowflake) ([]*common.Message, error) {
	msgs, err := d.RestClient.GetPinnedMessages(ctx, channelID)
	if err == nil {
		d.cacheMessages(msgs)
	}
	return msgs, err
}

func (d *Client) PinMessage(ctx context.Context, channelID, messageID common.Snowflake) error {
	return d.RestClient.PinMessage(ctx, channelID, messageID)
}

func (d *Client) UnpinMessage(ctx context.Context, channelID, messageID common.Snowflake) error {
	return d.RestClient.UnpinMessage(ctx, channelID, messageID)
}

func (d *Client) AddReaction(ctx context.Context, channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.AddReaction(ctx, channelID, messageID, emoji)
}

func (d *Client) DeleteOwnReaction(ctx context.Context, channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.DeleteOwnReaction(ctx, channelID, messageID, emoji)
}

func (d *Client) DeleteUserReaction(ctx context.Context, channelID, messageID common.Snowflake, emoji string, userID common.Snowflake) error {
	return d.RestClient.DeleteUserReaction(ctx, channelID, messageID, emoji, userID)
}

func (d *Client) GetReactions(ctx context.Context, channelID, messageID common.Snowflake, emoji string, params api.GetReactionsParams) ([]*common.User, error) {
	users, err := d.RestClient.GetReactions(ctx, channelID, messageID, emoji, params)
	if err == nil {
		d.cacheUsers(users)
	}
	return users, err
}

func (d *Client) DeleteAllReactions(ctx context.Context, channelID, messageID common.Snowflake) error {
	return d.RestClient.DeleteAllReactions(ctx, channelID, messageID)
}

func (d *Client) DeleteAllReactionsForEmoji(ctx context.Context, channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.DeleteAllReactionsForEmoji(ctx, channelID, messageID, emoji)
}

// ── Channel methods ───────────────────────────────────────────────────────────

func (d *Client) GetChannel(ctx context.Context, channelID common.Snowflake) (*common.Channel, error) {
	channel, err := d.RestClient.GetChannel(ctx, channelID)
	if err == nil {
		d.cacheChannel(channel)
	}
	return channel, err
}

func (d *Client) ModifyChannel(ctx context.Context, channelID common.Snowflake, params api.ModifyChannelParams) (*common.Channel, error) {
	channel, err := d.RestClient.ModifyChannel(ctx, channelID, params)
	if err == nil {
		d.cacheChannel(channel)
	}
	return channel, err
}

func (d *Client) DeleteChannel(ctx context.Context, channelID common.Snowflake) (*common.Channel, error) {
	channel, err := d.RestClient.DeleteChannel(ctx, channelID)
	if err == nil {
		d.removeChannelFromCache(channelID)
	}
	return channel, err
}

func (d *Client) GetChannelInvites(ctx context.Context, channelID common.Snowflake) ([]*api.Invite, error) {
	return d.RestClient.GetChannelInvites(ctx, channelID)
}

func (d *Client) CreateChannelInvite(ctx context.Context, channelID common.Snowflake, params api.CreateChannelInviteParams) (*api.Invite, error) {
	return d.RestClient.CreateChannelInvite(ctx, channelID, params)
}

func (d *Client) EditChannelPermissions(ctx context.Context, channelID, overwriteID common.Snowflake, params api.EditChannelPermissionsParams) error {
	return d.RestClient.EditChannelPermissions(ctx, channelID, overwriteID, params)
}

func (d *Client) DeleteChannelPermission(ctx context.Context, channelID, overwriteID common.Snowflake) error {
	return d.RestClient.DeleteChannelPermission(ctx, channelID, overwriteID)
}

func (d *Client) TriggerTypingIndicator(ctx context.Context, channelID common.Snowflake) error {
	return d.RestClient.TriggerTypingIndicator(ctx, channelID)
}

// ── Guild methods ─────────────────────────────────────────────────────────────

func (d *Client) GetGuild(ctx context.Context, guildID common.Snowflake, withCounts bool) (*common.Guild, error) {
	guild, err := d.RestClient.GetGuild(ctx, guildID, withCounts)
	if err == nil {
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) GetGuildPreview(ctx context.Context, guildID common.Snowflake) (*common.Guild, error) {
	return d.RestClient.GetGuildPreview(ctx, guildID)
}

func (d *Client) ModifyGuild(ctx context.Context, guildID common.Snowflake, params api.ModifyGuildParams) (*common.Guild, error) {
	guild, err := d.RestClient.ModifyGuild(ctx, guildID, params)
	if err == nil {
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) DeleteGuild(ctx context.Context, guildID common.Snowflake) error {
	if err := d.RestClient.DeleteGuild(ctx, guildID); err != nil {
		return err
	}
	d.removeGuildFromCache(guildID)
	return nil
}

func (d *Client) GetGuildChannels(ctx context.Context, guildID common.Snowflake) ([]*common.Channel, error) {
	channels, err := d.RestClient.GetGuildChannels(ctx, guildID)
	if err == nil {
		d.cacheChannels(channels)
	}
	return channels, err
}

func (d *Client) CreateGuildChannel(ctx context.Context, guildID common.Snowflake, params api.CreateGuildChannelParams) (*common.Channel, error) {
	channel, err := d.RestClient.CreateGuildChannel(ctx, guildID, params)
	if err == nil {
		d.cacheChannel(channel)
	}
	return channel, err
}

func (d *Client) ModifyGuildChannelPositions(ctx context.Context, guildID common.Snowflake, entries []api.ModifyGuildChannelPositionsEntry) error {
	return d.RestClient.ModifyGuildChannelPositions(ctx, guildID, entries)
}

func (d *Client) GetGuildRoles(ctx context.Context, guildID common.Snowflake) ([]*common.Role, error) {
	roles, err := d.RestClient.GetGuildRoles(ctx, guildID)
	if err == nil {
		d.cacheRoles(guildID, roles)
	}
	return roles, err
}

func (d *Client) GetGuildRole(ctx context.Context, guildID, roleID common.Snowflake) (*common.Role, error) {
	role, err := d.RestClient.GetGuildRole(ctx, guildID, roleID)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) CreateGuildRole(ctx context.Context, guildID common.Snowflake, params api.CreateGuildRoleParams) (*common.Role, error) {
	role, err := d.RestClient.CreateGuildRole(ctx, guildID, params)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) ModifyGuildRolePositions(ctx context.Context, guildID common.Snowflake, entries []api.ModifyGuildRolePositionsEntry) ([]*common.Role, error) {
	roles, err := d.RestClient.ModifyGuildRolePositions(ctx, guildID, entries)
	if err == nil {
		d.cacheRoles(guildID, roles)
	}
	return roles, err
}

func (d *Client) ModifyGuildRole(ctx context.Context, guildID, roleID common.Snowflake, params api.ModifyGuildRoleParams) (*common.Role, error) {
	role, err := d.RestClient.ModifyGuildRole(ctx, guildID, roleID, params)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) DeleteGuildRole(ctx context.Context, guildID, roleID common.Snowflake) error {
	if err := d.RestClient.DeleteGuildRole(ctx, guildID, roleID); err != nil {
		return err
	}
	d.removeRoleFromCache(roleID)
	return nil
}

func (d *Client) GetGuildBans(ctx context.Context, guildID common.Snowflake, params api.GetGuildBansParams) ([]*api.Ban, error) {
	return d.RestClient.GetGuildBans(ctx, guildID, params)
}

func (d *Client) GetGuildBan(ctx context.Context, guildID, userID common.Snowflake) (*api.Ban, error) {
	return d.RestClient.GetGuildBan(ctx, guildID, userID)
}

func (d *Client) CreateGuildBan(ctx context.Context, guildID, userID common.Snowflake, params api.CreateGuildBanParams) error {
	return d.RestClient.CreateGuildBan(ctx, guildID, userID, params)
}

func (d *Client) RemoveGuildBan(ctx context.Context, guildID, userID common.Snowflake) error {
	return d.RestClient.RemoveGuildBan(ctx, guildID, userID)
}

func (d *Client) GetGuildPruneCount(ctx context.Context, guildID common.Snowflake, params api.GetGuildPruneCountParams) (*api.GuildPruneCountResult, error) {
	return d.RestClient.GetGuildPruneCount(ctx, guildID, params)
}

func (d *Client) BeginGuildPrune(ctx context.Context, guildID common.Snowflake, params api.BeginGuildPruneParams) (*api.GuildPruneCountResult, error) {
	return d.RestClient.BeginGuildPrune(ctx, guildID, params)
}

func (d *Client) GetGuildInvites(ctx context.Context, guildID common.Snowflake) ([]*api.Invite, error) {
	return d.RestClient.GetGuildInvites(ctx, guildID)
}

func (d *Client) GetGuildVanityURL(ctx context.Context, guildID common.Snowflake) (*api.GuildVanityURL, error) {
	return d.RestClient.GetGuildVanityURL(ctx, guildID)
}

func (d *Client) GetGuildAuditLog(ctx context.Context, guildID common.Snowflake, params api.GetGuildAuditLogParams) (*api.AuditLog, error) {
	return d.RestClient.GetGuildAuditLog(ctx, guildID, params)
}

// GetGuildRoleMemberCounts returns a map of role ID → member count for every role in the guild.
func (d *Client) GetGuildRoleMemberCounts(ctx context.Context, guildID common.Snowflake) (map[string]int64, error) {
	return d.RestClient.GetGuildRoleMemberCounts(ctx, guildID)
}

// ── Member methods ────────────────────────────────────────────────────────────

func (d *Client) GetGuildMember(ctx context.Context, guildID, userID common.Snowflake) (*common.GuildMember, error) {
	member, err := d.RestClient.GetGuildMember(ctx, guildID, userID)
	if err == nil {
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) ListGuildMembers(ctx context.Context, guildID common.Snowflake, params api.GetGuildMembersParams) ([]*common.GuildMember, error) {
	members, err := d.RestClient.ListGuildMembers(ctx, guildID, params)
	if err == nil {
		d.cacheMembers(guildID, members)
	}
	return members, err
}

func (d *Client) SearchGuildMembers(ctx context.Context, guildID common.Snowflake, params api.SearchGuildMembersParams) ([]*common.GuildMember, error) {
	members, err := d.RestClient.SearchGuildMembers(ctx, guildID, params)
	if err == nil {
		d.cacheMembers(guildID, members)
	}
	return members, err
}

func (d *Client) ModifyGuildMember(ctx context.Context, guildID, userID common.Snowflake, params api.ModifyGuildMemberParams) (*common.GuildMember, error) {
	member, err := d.RestClient.ModifyGuildMember(ctx, guildID, userID, params)
	if err == nil {
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) ModifyCurrentMember(ctx context.Context, guildID common.Snowflake, params api.ModifyCurrentMemberParams) (*common.GuildMember, error) {
	member, err := d.RestClient.ModifyCurrentMember(ctx, guildID, params)
	if err == nil {
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) AddGuildMemberRole(ctx context.Context, guildID, userID, roleID common.Snowflake) error {
	return d.RestClient.AddGuildMemberRole(ctx, guildID, userID, roleID)
}

func (d *Client) RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID common.Snowflake) error {
	return d.RestClient.RemoveGuildMemberRole(ctx, guildID, userID, roleID)
}

// RemoveGuildMember removes a member from a guild. Requires KICK_MEMBERS.
func (d *Client) RemoveGuildMember(ctx context.Context, guildID, userID common.Snowflake) error {
	if err := d.RestClient.RemoveGuildMember(ctx, guildID, userID); err != nil {
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
func (d *Client) FetchAllGuildMembers(ctx context.Context, guildID common.Snowflake) ([]*common.GuildMember, error) {
	members, err := d.RestClient.FetchAllGuildMembers(ctx, guildID)
	if err == nil {
		d.cacheMembers(guildID, members)
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
	UserIDs []common.Snowflake
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
func (d *Client) RequestGuildMembers(guildID common.Snowflake, params RequestGuildMembersParams) error {
	if d.Websocket == nil || d.Websocket.Connection == nil {
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
	return d.Websocket.writeJSON(map[string]interface{}{
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
	Activities []common.FullActivity
	// Status is the bot's status ("online", "dnd", "idle", "invisible", "offline").
	Status common.PresenceStatus
	// AFK indicates whether the bot is AFK.
	AFK bool
}

// UpdatePresence sends an OP 3 Presence Update to the Discord gateway,
// updating the bot's displayed status and activity.
func (d *Client) UpdatePresence(params UpdatePresenceParams) error {
	if d.Websocket == nil || d.Websocket.Connection == nil {
		return errors.New("not connected to gateway")
	}

	activities := params.Activities
	if activities == nil {
		activities = []common.FullActivity{}
	}

	data := map[string]interface{}{
		"since":      params.Since,
		"activities": activities,
		"status":     string(params.Status),
		"afk":        params.AFK,
	}
	return d.Websocket.writeJSON(map[string]interface{}{
		"op": 3,
		"d":  data,
	})
}

// ── Guild widget ──────────────────────────────────────────────────────────────

func (d *Client) GetGuildWidgetSettings(ctx context.Context, guildID common.Snowflake) (*common.GuildWidgetSettings, error) {
	return d.RestClient.GetGuildWidgetSettings(ctx, guildID)
}

func (d *Client) ModifyGuildWidgetSettings(ctx context.Context, guildID common.Snowflake, params api.ModifyGuildWidgetParams) (*common.GuildWidgetSettings, error) {
	return d.RestClient.ModifyGuildWidgetSettings(ctx, guildID, params)
}

func (d *Client) GetGuildWidget(ctx context.Context, guildID common.Snowflake) (*common.GuildWidget, error) {
	return d.RestClient.GetGuildWidget(ctx, guildID)
}

// GetGuildWidgetImage returns the PNG widget image for a guild as raw bytes.
// style is an optional widget style ("shield", "banner1"–"banner4"); pass "" for the default.
func (d *Client) GetGuildWidgetImage(ctx context.Context, guildID common.Snowflake, style string) ([]byte, error) {
	return d.RestClient.GetGuildWidgetImage(ctx, guildID, style)
}

// ModifyGuildMFALevel sets the required MFA level for moderators in a guild.
// Requires guild ownership. level: 0 = none, 1 = elevated.
func (d *Client) ModifyGuildMFALevel(ctx context.Context, guildID common.Snowflake, level int) (int, error) {
	return d.RestClient.ModifyGuildMFALevel(ctx, guildID, level)
}

// ── Welcome screen ────────────────────────────────────────────────────────────

func (d *Client) GetGuildWelcomeScreen(ctx context.Context, guildID common.Snowflake) (*common.GuildWelcomeScreen, error) {
	return d.RestClient.GetGuildWelcomeScreen(ctx, guildID)
}

func (d *Client) ModifyGuildWelcomeScreen(ctx context.Context, guildID common.Snowflake, params api.ModifyGuildWelcomeScreenParams) (*common.GuildWelcomeScreen, error) {
	return d.RestClient.ModifyGuildWelcomeScreen(ctx, guildID, params)
}

// ── Guild onboarding ──────────────────────────────────────────────────────────

func (d *Client) GetGuildOnboarding(ctx context.Context, guildID common.Snowflake) (*common.GuildOnboarding, error) {
	return d.RestClient.GetGuildOnboarding(ctx, guildID)
}

func (d *Client) ModifyGuildOnboarding(ctx context.Context, guildID common.Snowflake, params api.ModifyGuildOnboardingParams) (*common.GuildOnboarding, error) {
	return d.RestClient.ModifyGuildOnboarding(ctx, guildID, params)
}

// ── Voice ─────────────────────────────────────────────────────────────────────

func (d *Client) ListVoiceRegions(ctx context.Context) ([]*common.VoiceRegion, error) {
	return d.RestClient.ListVoiceRegions(ctx)
}

func (d *Client) ListGuildVoiceRegions(ctx context.Context, guildID common.Snowflake) ([]*common.VoiceRegion, error) {
	return d.RestClient.ListGuildVoiceRegions(ctx, guildID)
}

func (d *Client) ModifyCurrentUserVoiceState(ctx context.Context, guildID common.Snowflake, params api.ModifyCurrentUserVoiceStateParams) error {
	return d.RestClient.ModifyCurrentUserVoiceState(ctx, guildID, params)
}

func (d *Client) ModifyUserVoiceState(ctx context.Context, guildID, userID common.Snowflake, params api.ModifyUserVoiceStateParams) error {
	return d.RestClient.ModifyUserVoiceState(ctx, guildID, userID, params)
}

// ── Soundboard ────────────────────────────────────────────────────────────────

func (d *Client) ListDefaultSoundboardSounds(ctx context.Context) ([]*common.SoundboardSound, error) {
	return d.RestClient.ListDefaultSoundboardSounds(ctx)
}

func (d *Client) ListGuildSoundboardSounds(ctx context.Context, guildID common.Snowflake) ([]*common.SoundboardSound, error) {
	return d.RestClient.ListGuildSoundboardSounds(ctx, guildID)
}

func (d *Client) GetGuildSoundboardSound(ctx context.Context, guildID, soundID common.Snowflake) (*common.SoundboardSound, error) {
	return d.RestClient.GetGuildSoundboardSound(ctx, guildID, soundID)
}

func (d *Client) CreateGuildSoundboardSound(ctx context.Context, guildID common.Snowflake, params api.CreateGuildSoundboardSoundParams) (*common.SoundboardSound, error) {
	return d.RestClient.CreateGuildSoundboardSound(ctx, guildID, params)
}

func (d *Client) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID common.Snowflake, params api.ModifyGuildSoundboardSoundParams) (*common.SoundboardSound, error) {
	return d.RestClient.ModifyGuildSoundboardSound(ctx, guildID, soundID, params)
}

func (d *Client) DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID common.Snowflake) error {
	return d.RestClient.DeleteGuildSoundboardSound(ctx, guildID, soundID)
}

func (d *Client) SendSoundboardSound(ctx context.Context, channelID common.Snowflake, params api.SendSoundboardSoundParams) error {
	return d.RestClient.SendSoundboardSound(ctx, channelID, params)
}

// ── Application ───────────────────────────────────────────────────────────────

func (d *Client) GetCurrentApplication(ctx context.Context) (*common.Application, error) {
	return d.RestClient.GetCurrentApplication(ctx)
}

func (d *Client) ModifyCurrentApplication(ctx context.Context, params api.ModifyCurrentApplicationParams) (*common.Application, error) {
	return d.RestClient.ModifyCurrentApplication(ctx, params)
}

// ── Command permissions ───────────────────────────────────────────────────────

// GetGuildApplicationCommandPermissions returns all permission overrides for every command in a guild.
// Uses the bot's own application ID automatically.
func (d *Client) GetGuildApplicationCommandPermissions(ctx context.Context, guildID common.Snowflake) ([]*common.GuildApplicationCommandPermissions, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.GetGuildApplicationCommandPermissions(ctx, appID, guildID)
}

// GetApplicationCommandPermissions returns the permission overrides for a specific command.
func (d *Client) GetApplicationCommandPermissions(ctx context.Context, guildID, cmdID common.Snowflake) (*common.GuildApplicationCommandPermissions, error) {
	appID, err := d.applicationID()
	if err != nil {
		return nil, err
	}
	return d.RestClient.GetApplicationCommandPermissions(ctx, appID, guildID, cmdID)
}

// ── Application emojis ────────────────────────────────────────────────────────

func (d *Client) ListApplicationEmojis(ctx context.Context, appID common.Snowflake) ([]*common.Emoji, error) {
	return d.RestClient.ListApplicationEmojis(ctx, appID)
}

func (d *Client) GetApplicationEmoji(ctx context.Context, appID, emojiID common.Snowflake) (*common.Emoji, error) {
	return d.RestClient.GetApplicationEmoji(ctx, appID, emojiID)
}

func (d *Client) CreateApplicationEmoji(ctx context.Context, appID common.Snowflake, params api.CreateEmojiParams) (*common.Emoji, error) {
	return d.RestClient.CreateApplicationEmoji(ctx, appID, params)
}

func (d *Client) ModifyApplicationEmoji(ctx context.Context, appID, emojiID common.Snowflake, params api.ModifyEmojiParams) (*common.Emoji, error) {
	return d.RestClient.ModifyApplicationEmoji(ctx, appID, emojiID, params)
}

func (d *Client) DeleteApplicationEmoji(ctx context.Context, appID, emojiID common.Snowflake) error {
	return d.RestClient.DeleteApplicationEmoji(ctx, appID, emojiID)
}

// ── Entitlements ──────────────────────────────────────────────────────────────

func (d *Client) ListEntitlements(ctx context.Context, appID common.Snowflake, params api.ListEntitlementsParams) ([]*common.Entitlement, error) {
	return d.RestClient.ListEntitlements(ctx, appID, params)
}

func (d *Client) GetEntitlement(ctx context.Context, appID, entitlementID common.Snowflake) (*common.Entitlement, error) {
	return d.RestClient.GetEntitlement(ctx, appID, entitlementID)
}

func (d *Client) CreateTestEntitlement(ctx context.Context, appID common.Snowflake, params api.CreateTestEntitlementParams) (*common.Entitlement, error) {
	return d.RestClient.CreateTestEntitlement(ctx, appID, params)
}

func (d *Client) ConsumeEntitlement(ctx context.Context, appID, entitlementID common.Snowflake) error {
	return d.RestClient.ConsumeEntitlement(ctx, appID, entitlementID)
}

func (d *Client) DeleteTestEntitlement(ctx context.Context, appID, entitlementID common.Snowflake) error {
	return d.RestClient.DeleteTestEntitlement(ctx, appID, entitlementID)
}

// ── SKUs ──────────────────────────────────────────────────────────────────────

func (d *Client) ListSKUs(ctx context.Context, appID common.Snowflake) ([]*common.SKU, error) {
	return d.RestClient.ListSKUs(ctx, appID)
}

// ── Subscriptions ─────────────────────────────────────────────────────────────

func (d *Client) ListSKUSubscriptions(ctx context.Context, skuID common.Snowflake, params api.ListSKUSubscriptionsParams) ([]*common.Subscription, error) {
	return d.RestClient.ListSKUSubscriptions(ctx, skuID, params)
}

func (d *Client) GetSKUSubscription(ctx context.Context, skuID, subscriptionID common.Snowflake) (*common.Subscription, error) {
	return d.RestClient.GetSKUSubscription(ctx, skuID, subscriptionID)
}

// ── Poll ──────────────────────────────────────────────────────────────────────

func (d *Client) GetPollAnswerVoters(ctx context.Context, channelID, messageID common.Snowflake, answerID int, params api.GetPollAnswerVotersParams) ([]*common.User, error) {
	users, err := d.RestClient.GetPollAnswerVoters(ctx, channelID, messageID, answerID, params)
	if err == nil {
		d.cacheUsers(users)
	}
	return users, err
}

func (d *Client) EndPoll(ctx context.Context, channelID, messageID common.Snowflake) (*common.Message, error) {
	msg, err := d.RestClient.EndPoll(ctx, channelID, messageID)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

// ── Activity ──────────────────────────────────────────────────────────────────

func (d *Client) GetActivityInstance(ctx context.Context, appID common.Snowflake, instanceID string) (*common.ActivityInstance, error) {
	return d.RestClient.GetActivityInstance(ctx, appID, instanceID)
}

// ── Additional channel methods ────────────────────────────────────────────────

func (d *Client) SetVoiceChannelStatus(ctx context.Context, channelID common.Snowflake, status *string) error {
	return d.RestClient.SetVoiceChannelStatus(ctx, channelID, status)
}

func (d *Client) FollowAnnouncementChannel(ctx context.Context, channelID, webhookChannelID common.Snowflake) (*api.FollowedChannel, error) {
	return d.RestClient.FollowAnnouncementChannel(ctx, channelID, webhookChannelID)
}

func (d *Client) AddGroupDMRecipient(ctx context.Context, channelID, userID common.Snowflake, params api.AddGroupDMRecipientParams) error {
	return d.RestClient.AddGroupDMRecipient(ctx, channelID, userID, params)
}

func (d *Client) RemoveGroupDMRecipient(ctx context.Context, channelID, userID common.Snowflake) error {
	return d.RestClient.RemoveGroupDMRecipient(ctx, channelID, userID)
}

// ── Additional guild methods ──────────────────────────────────────────────────

func (d *Client) CreateGuild(ctx context.Context, params api.CreateGuildParams) (*common.Guild, error) {
	guild, err := d.RestClient.CreateGuild(ctx, params)
	if err == nil {
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) AddGuildMember(ctx context.Context, guildID, userID common.Snowflake, params api.AddGuildMemberParams) (*common.GuildMember, error) {
	member, err := d.RestClient.AddGuildMember(ctx, guildID, userID, params)
	if err == nil {
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) BulkBanGuildMembers(ctx context.Context, guildID common.Snowflake, params api.BulkBanParams) (*api.BulkBanResult, error) {
	return d.RestClient.BulkBanGuildMembers(ctx, guildID, params)
}

func (d *Client) GetGuildIntegrations(ctx context.Context, guildID common.Snowflake) ([]*common.Integration, error) {
	return d.RestClient.GetGuildIntegrations(ctx, guildID)
}

func (d *Client) DeleteGuildIntegration(ctx context.Context, guildID, integrationID common.Snowflake) error {
	return d.RestClient.DeleteGuildIntegration(ctx, guildID, integrationID)
}

// ── Additional voice methods ──────────────────────────────────────────────────

func (d *Client) GetCurrentUserVoiceState(ctx context.Context, guildID common.Snowflake) (*common.VoiceState, error) {
	return d.RestClient.GetCurrentUserVoiceState(ctx, guildID)
}

func (d *Client) GetUserVoiceState(ctx context.Context, guildID, userID common.Snowflake) (*common.VoiceState, error) {
	return d.RestClient.GetUserVoiceState(ctx, guildID, userID)
}

// ── Additional sticker methods ────────────────────────────────────────────────

func (d *Client) GetStickerPack(ctx context.Context, packID common.Snowflake) (*api.StickerPack, error) {
	return d.RestClient.GetStickerPack(ctx, packID)
}

func (d *Client) CreateGuildSticker(ctx context.Context, guildID common.Snowflake, params api.CreateGuildStickerParams) (*common.Sticker, error) {
	return d.RestClient.CreateGuildSticker(ctx, guildID, params)
}

// ── Invite methods ────────────────────────────────────────────────────────────

func (d *Client) GetInvite(ctx context.Context, code string, params api.GetInviteParams) (*api.Invite, error) {
	return d.RestClient.GetInvite(ctx, code, params)
}

func (d *Client) DeleteInvite(ctx context.Context, code string) (*api.Invite, error) {
	return d.RestClient.DeleteInvite(ctx, code)
}

// ── Additional user methods ───────────────────────────────────────────────────

func (d *Client) GetCurrentUserConnections(ctx context.Context) ([]*api.UserConnection, error) {
	return d.RestClient.GetCurrentUserConnections(ctx)
}

func (d *Client) GetCurrentUserApplicationRoleConnection(ctx context.Context, appID common.Snowflake) (*api.ApplicationRoleConnection, error) {
	return d.RestClient.GetCurrentUserApplicationRoleConnection(ctx, appID)
}

func (d *Client) UpdateCurrentUserApplicationRoleConnection(ctx context.Context, appID common.Snowflake, params api.UpdateRoleConnectionParams) (*api.ApplicationRoleConnection, error) {
	return d.RestClient.UpdateCurrentUserApplicationRoleConnection(ctx, appID, params)
}

func (d *Client) CreateGroupDM(ctx context.Context, params api.CreateGroupDMParams) (*common.Channel, error) {
	channel, err := d.RestClient.CreateGroupDM(ctx, params)
	if err == nil {
		d.cacheChannel(channel)
	}
	return channel, err
}

// ── Additional webhook methods ────────────────────────────────────────────────

func (d *Client) ExecuteSlackWebhook(ctx context.Context, webhookID common.Snowflake, token string, wait bool, body json.RawMessage) error {
	return d.RestClient.ExecuteSlackWebhook(ctx, webhookID, token, wait, body)
}

func (d *Client) ExecuteGitHubWebhook(ctx context.Context, webhookID common.Snowflake, token string, wait bool, body json.RawMessage) error {
	return d.RestClient.ExecuteGitHubWebhook(ctx, webhookID, token, wait, body)
}
