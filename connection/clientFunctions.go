package connection

import (
	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// ── Application commands ──────────────────────────────────────────────────────

// RegisterCommand registers a single application command.
// Must be called after Login() so that the application ID is available.
func (d *Client) RegisterCommand(cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	return d.RestClient.RegisterCommand(d.User.ID, cmd)
}

// BulkRegisterCommands overwrites all application commands with the provided list.
// Any existing commands not included will be deleted.
// Must be called after Login() so that the application ID is available.
func (d *Client) BulkRegisterCommands(cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
	return d.RestClient.BulkRegisterCommands(d.User.ID, cmds)
}

// ── Interaction responses ─────────────────────────────────────────────────────

// Reply sends an immediate message response to an interaction.
// Pass withResponse=true to get the created message back.
func (d *Client) Reply(i *interactions.Interaction, data *responses.InteractionResponseDataDefault, withResponse bool) (*responses.InteractionCallbackResponse, error) {
	return d.RestClient.CreateInteractionResponse(i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeChannelMessageWithSource,
		Data: data,
	}, withResponse)
}

// DeferReply acknowledges the interaction, telling Discord the bot will respond later.
// Set ephemeral=true to make the eventual follow-up visible only to the invoker.
func (d *Client) DeferReply(i *interactions.Interaction, ephemeral bool) error {
	var data *responses.InteractionResponseDataDefault
	if ephemeral {
		data = &responses.InteractionResponseDataDefault{Flags: common.MessageFlagEphemeral}
	}
	_, err := d.RestClient.CreateInteractionResponse(i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeDeferredChannelMessageWithSource,
		Data: data,
	}, false)
	return err
}

// DeferUpdateMessage acknowledges a component interaction without editing the original message.
// Use EditReply afterwards to push the actual update.
func (d *Client) DeferUpdateMessage(i *interactions.Interaction) error {
	_, err := d.RestClient.CreateInteractionResponse(i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeDeferredUpdateMessage,
	}, false)
	return err
}

// UpdateMessage edits the message that triggered a component interaction.
func (d *Client) UpdateMessage(i *interactions.Interaction, data *responses.InteractionResponseDataDefault) error {
	_, err := d.RestClient.CreateInteractionResponse(i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeUpdateMessage,
		Data: data,
	}, false)
	return err
}

// ReplyWithModal responds to an interaction with a modal dialog.
func (d *Client) ReplyWithModal(i *interactions.Interaction, modal *components.Modal) error {
	_, err := d.RestClient.CreateInteractionResponse(i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeModal,
		Data: modal,
	}, false)
	return err
}

// ReplyAutocomplete sends autocomplete choices in response to an autocomplete interaction.
func (d *Client) ReplyAutocomplete(i *interactions.Interaction, choices []responses.AutocompleteChoice) error {
	_, err := d.RestClient.CreateInteractionResponse(i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeApplicationCommandAutocompleteResult,
		Data: &responses.InteractionResponseDataAutocomplete{Choices: choices},
	}, false)
	return err
}

// GetOriginalResponse fetches the original message sent in response to an interaction.
func (d *Client) GetOriginalResponse(i *interactions.Interaction) (*common.Message, error) {
	return d.RestClient.GetOriginalInteractionResponse(d.User.ID, i.Token)
}

// EditReply edits the original interaction response.
func (d *Client) EditReply(i *interactions.Interaction, params api.EditMessageParams) (*common.Message, error) {
	return d.RestClient.EditOriginalInteractionResponse(d.User.ID, i.Token, params)
}

// DeleteReply deletes the original interaction response.
func (d *Client) DeleteReply(i *interactions.Interaction) error {
	return d.RestClient.DeleteOriginalInteractionResponse(d.User.ID, i.Token)
}

// CreateFollowup sends a follow-up message to an interaction (up to 15 minutes after the initial response).
func (d *Client) CreateFollowup(i *interactions.Interaction, params api.CreateMessageParams) (*common.Message, error) {
	return d.RestClient.CreateFollowupMessage(d.User.ID, i.Token, params)
}

// GetFollowup fetches a follow-up message.
func (d *Client) GetFollowup(i *interactions.Interaction, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.GetFollowupMessage(d.User.ID, i.Token, messageID)
}

// EditFollowup edits a follow-up message.
func (d *Client) EditFollowup(i *interactions.Interaction, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	return d.RestClient.EditFollowupMessage(d.User.ID, i.Token, messageID, params)
}

// DeleteFollowup deletes a follow-up message.
func (d *Client) DeleteFollowup(i *interactions.Interaction, messageID common.Snowflake) error {
	return d.RestClient.DeleteFollowupMessage(d.User.ID, i.Token, messageID)
}

// ── Message methods ───────────────────────────────────────────────────────────

func (d *Client) GetMessages(channelID common.Snowflake, params api.GetMessagesParams) ([]*common.Message, error) {
	return d.RestClient.GetMessages(channelID, params)
}

func (d *Client) GetMessage(channelID, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.GetMessage(channelID, messageID)
}

func (d *Client) SendMessage(channelID common.Snowflake, params api.CreateMessageParams) (*common.Message, error) {
	return d.RestClient.CreateMessage(channelID, params)
}

func (d *Client) EditMessage(channelID, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	return d.RestClient.EditMessage(channelID, messageID, params)
}

func (d *Client) DeleteMessage(channelID, messageID common.Snowflake) error {
	return d.RestClient.DeleteMessage(channelID, messageID)
}

func (d *Client) BulkDeleteMessages(channelID common.Snowflake, messageIDs []common.Snowflake) error {
	return d.RestClient.BulkDeleteMessages(channelID, messageIDs)
}

func (d *Client) CrosspostMessage(channelID, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.CrosspostMessage(channelID, messageID)
}

func (d *Client) GetPinnedMessages(channelID common.Snowflake) ([]*common.Message, error) {
	return d.RestClient.GetPinnedMessages(channelID)
}

func (d *Client) PinMessage(channelID, messageID common.Snowflake) error {
	return d.RestClient.PinMessage(channelID, messageID)
}

func (d *Client) UnpinMessage(channelID, messageID common.Snowflake) error {
	return d.RestClient.UnpinMessage(channelID, messageID)
}

func (d *Client) AddReaction(channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.AddReaction(channelID, messageID, emoji)
}

func (d *Client) DeleteOwnReaction(channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.DeleteOwnReaction(channelID, messageID, emoji)
}

func (d *Client) DeleteUserReaction(channelID, messageID common.Snowflake, emoji string, userID common.Snowflake) error {
	return d.RestClient.DeleteUserReaction(channelID, messageID, emoji, userID)
}

func (d *Client) GetReactions(channelID, messageID common.Snowflake, emoji string, params api.GetReactionsParams) ([]*common.User, error) {
	return d.RestClient.GetReactions(channelID, messageID, emoji, params)
}

func (d *Client) DeleteAllReactions(channelID, messageID common.Snowflake) error {
	return d.RestClient.DeleteAllReactions(channelID, messageID)
}

func (d *Client) DeleteAllReactionsForEmoji(channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.DeleteAllReactionsForEmoji(channelID, messageID, emoji)
}

// ── Channel methods ───────────────────────────────────────────────────────────

func (d *Client) GetChannel(channelID common.Snowflake) (*common.Channel, error) {
	return d.RestClient.GetChannel(channelID)
}

func (d *Client) ModifyChannel(channelID common.Snowflake, params api.ModifyChannelParams) (*common.Channel, error) {
	return d.RestClient.ModifyChannel(channelID, params)
}

func (d *Client) DeleteChannel(channelID common.Snowflake) (*common.Channel, error) {
	return d.RestClient.DeleteChannel(channelID)
}

func (d *Client) GetChannelInvites(channelID common.Snowflake) ([]*api.Invite, error) {
	return d.RestClient.GetChannelInvites(channelID)
}

func (d *Client) CreateChannelInvite(channelID common.Snowflake, params api.CreateChannelInviteParams) (*api.Invite, error) {
	return d.RestClient.CreateChannelInvite(channelID, params)
}

func (d *Client) EditChannelPermissions(channelID, overwriteID common.Snowflake, params api.EditChannelPermissionsParams) error {
	return d.RestClient.EditChannelPermissions(channelID, overwriteID, params)
}

func (d *Client) DeleteChannelPermission(channelID, overwriteID common.Snowflake) error {
	return d.RestClient.DeleteChannelPermission(channelID, overwriteID)
}

func (d *Client) TriggerTypingIndicator(channelID common.Snowflake) error {
	return d.RestClient.TriggerTypingIndicator(channelID)
}

// ── Guild methods ─────────────────────────────────────────────────────────────

func (d *Client) GetGuild(guildID common.Snowflake, withCounts bool) (*common.Guild, error) {
	return d.RestClient.GetGuild(guildID, withCounts)
}

func (d *Client) GetGuildPreview(guildID common.Snowflake) (*common.Guild, error) {
	return d.RestClient.GetGuildPreview(guildID)
}

func (d *Client) ModifyGuild(guildID common.Snowflake, params api.ModifyGuildParams) (*common.Guild, error) {
	return d.RestClient.ModifyGuild(guildID, params)
}

func (d *Client) DeleteGuild(guildID common.Snowflake) error {
	return d.RestClient.DeleteGuild(guildID)
}

func (d *Client) GetGuildChannels(guildID common.Snowflake) ([]*common.Channel, error) {
	return d.RestClient.GetGuildChannels(guildID)
}

func (d *Client) CreateGuildChannel(guildID common.Snowflake, params api.CreateGuildChannelParams) (*common.Channel, error) {
	return d.RestClient.CreateGuildChannel(guildID, params)
}

func (d *Client) ModifyGuildChannelPositions(guildID common.Snowflake, entries []api.ModifyGuildChannelPositionsEntry) error {
	return d.RestClient.ModifyGuildChannelPositions(guildID, entries)
}

func (d *Client) GetGuildRoles(guildID common.Snowflake) ([]*common.Role, error) {
	return d.RestClient.GetGuildRoles(guildID)
}

func (d *Client) GetGuildRole(guildID, roleID common.Snowflake) (*common.Role, error) {
	return d.RestClient.GetGuildRole(guildID, roleID)
}

func (d *Client) CreateGuildRole(guildID common.Snowflake, params api.CreateGuildRoleParams) (*common.Role, error) {
	return d.RestClient.CreateGuildRole(guildID, params)
}

func (d *Client) ModifyGuildRolePositions(guildID common.Snowflake, entries []api.ModifyGuildRolePositionsEntry) ([]*common.Role, error) {
	return d.RestClient.ModifyGuildRolePositions(guildID, entries)
}

func (d *Client) ModifyGuildRole(guildID, roleID common.Snowflake, params api.ModifyGuildRoleParams) (*common.Role, error) {
	return d.RestClient.ModifyGuildRole(guildID, roleID, params)
}

func (d *Client) DeleteGuildRole(guildID, roleID common.Snowflake) error {
	return d.RestClient.DeleteGuildRole(guildID, roleID)
}

func (d *Client) GetGuildBans(guildID common.Snowflake, params api.GetGuildBansParams) ([]*api.Ban, error) {
	return d.RestClient.GetGuildBans(guildID, params)
}

func (d *Client) GetGuildBan(guildID, userID common.Snowflake) (*api.Ban, error) {
	return d.RestClient.GetGuildBan(guildID, userID)
}

func (d *Client) CreateGuildBan(guildID, userID common.Snowflake, params api.CreateGuildBanParams) error {
	return d.RestClient.CreateGuildBan(guildID, userID, params)
}

func (d *Client) RemoveGuildBan(guildID, userID common.Snowflake) error {
	return d.RestClient.RemoveGuildBan(guildID, userID)
}

func (d *Client) GetGuildPruneCount(guildID common.Snowflake, params api.GetGuildPruneCountParams) (*api.GuildPruneCountResult, error) {
	return d.RestClient.GetGuildPruneCount(guildID, params)
}

func (d *Client) BeginGuildPrune(guildID common.Snowflake, params api.BeginGuildPruneParams) (*api.GuildPruneCountResult, error) {
	return d.RestClient.BeginGuildPrune(guildID, params)
}

func (d *Client) GetGuildInvites(guildID common.Snowflake) ([]*api.Invite, error) {
	return d.RestClient.GetGuildInvites(guildID)
}

func (d *Client) GetGuildVanityURL(guildID common.Snowflake) (*api.GuildVanityURL, error) {
	return d.RestClient.GetGuildVanityURL(guildID)
}

func (d *Client) GetGuildAuditLog(guildID common.Snowflake, params api.GetGuildAuditLogParams) (*api.AuditLog, error) {
	return d.RestClient.GetGuildAuditLog(guildID, params)
}

// ── Member methods ────────────────────────────────────────────────────────────

func (d *Client) GetGuildMember(guildID, userID common.Snowflake) (*common.GuildMember, error) {
	return d.RestClient.GetGuildMember(guildID, userID)
}

func (d *Client) ListGuildMembers(guildID common.Snowflake, params api.GetGuildMembersParams) ([]*common.GuildMember, error) {
	return d.RestClient.ListGuildMembers(guildID, params)
}

func (d *Client) SearchGuildMembers(guildID common.Snowflake, params api.SearchGuildMembersParams) ([]*common.GuildMember, error) {
	return d.RestClient.SearchGuildMembers(guildID, params)
}

func (d *Client) ModifyGuildMember(guildID, userID common.Snowflake, params api.ModifyGuildMemberParams) (*common.GuildMember, error) {
	return d.RestClient.ModifyGuildMember(guildID, userID, params)
}

func (d *Client) ModifyCurrentMember(guildID common.Snowflake, params api.ModifyCurrentMemberParams) (*common.GuildMember, error) {
	return d.RestClient.ModifyCurrentMember(guildID, params)
}

func (d *Client) AddGuildMemberRole(guildID, userID, roleID common.Snowflake) error {
	return d.RestClient.AddGuildMemberRole(guildID, userID, roleID)
}

func (d *Client) RemoveGuildMemberRole(guildID, userID, roleID common.Snowflake) error {
	return d.RestClient.RemoveGuildMemberRole(guildID, userID, roleID)
}

// KickGuildMember removes a member from a guild. Requires KICK_MEMBERS.
func (d *Client) KickGuildMember(guildID, userID common.Snowflake) error {
	return d.RestClient.RemoveGuildMember(guildID, userID)
}
