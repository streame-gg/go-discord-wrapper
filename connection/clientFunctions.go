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

// ── Application commands ──────────────────────────────────────────────────────

// RegisterCommand registers a single application command.
// Must be called after Login() so that the application ID is available.
func (d *Client) RegisterCommand(cmd *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	return d.RestClient.RegisterCommand(context.Background(), d.User.ID, cmd)
}

// BulkRegisterCommands overwrites all application commands with the provided list.
// Any existing commands not included will be deleted.
// Must be called after Login() so that the application ID is available.
func (d *Client) BulkRegisterCommands(cmds []*commands.ApplicationCommand) ([]*commands.ApplicationCommand, error) {
	return d.RestClient.BulkRegisterCommands(context.Background(), d.User.ID, cmds)
}

// ── Interaction responses ─────────────────────────────────────────────────────

// Reply sends an immediate message response to an interaction.
// Pass withResponse=true to get the created message back.
func (d *Client) Reply(i *interactions.Interaction, data *responses.InteractionResponseDataDefault, withResponse bool) (*responses.InteractionCallbackResponse, error) {
	return d.RestClient.CreateInteractionResponse(context.Background(), i.ID, i.Token, responses.InteractionResponse{
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
	_, err := d.RestClient.CreateInteractionResponse(context.Background(), i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeDeferredChannelMessageWithSource,
		Data: data,
	}, false)
	return err
}

// DeferUpdateMessage acknowledges a component interaction without editing the original message.
// Use EditReply afterwards to push the actual update.
func (d *Client) DeferUpdateMessage(i *interactions.Interaction) error {
	_, err := d.RestClient.CreateInteractionResponse(context.Background(), i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeDeferredUpdateMessage,
	}, false)
	return err
}

// UpdateMessage edits the message that triggered a component interaction.
func (d *Client) UpdateMessage(i *interactions.Interaction, data *responses.InteractionResponseDataDefault) error {
	_, err := d.RestClient.CreateInteractionResponse(context.Background(), i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeUpdateMessage,
		Data: data,
	}, false)
	return err
}

// ReplyWithModal responds to an interaction with a modal dialog.
func (d *Client) ReplyWithModal(i *interactions.Interaction, modal *components.Modal) error {
	_, err := d.RestClient.CreateInteractionResponse(context.Background(), i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeModal,
		Data: modal,
	}, false)
	return err
}

// ReplyAutocomplete sends autocomplete choices in response to an autocomplete interaction.
func (d *Client) ReplyAutocomplete(i *interactions.Interaction, choices []responses.AutocompleteChoice) error {
	_, err := d.RestClient.CreateInteractionResponse(context.Background(), i.ID, i.Token, responses.InteractionResponse{
		Type: common.InteractionCallbackTypeApplicationCommandAutocompleteResult,
		Data: &responses.InteractionResponseDataAutocomplete{Choices: choices},
	}, false)
	return err
}

// GetOriginalResponse fetches the original message sent in response to an interaction.
func (d *Client) GetOriginalResponse(i *interactions.Interaction) (*common.Message, error) {
	return d.RestClient.GetOriginalInteractionResponse(context.Background(), d.User.ID, i.Token)
}

// EditReply edits the original interaction response.
func (d *Client) EditReply(i *interactions.Interaction, params api.EditMessageParams) (*common.Message, error) {
	return d.RestClient.EditOriginalInteractionResponse(context.Background(), d.User.ID, i.Token, params)
}

// DeleteReply deletes the original interaction response.
func (d *Client) DeleteReply(i *interactions.Interaction) error {
	return d.RestClient.DeleteOriginalInteractionResponse(context.Background(), d.User.ID, i.Token)
}

// CreateFollowup sends a follow-up message to an interaction (up to 15 minutes after the initial response).
func (d *Client) CreateFollowup(i *interactions.Interaction, params api.CreateMessageParams) (*common.Message, error) {
	return d.RestClient.CreateFollowupMessage(context.Background(), d.User.ID, i.Token, params)
}

// GetFollowup fetches a follow-up message.
func (d *Client) GetFollowup(i *interactions.Interaction, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.GetFollowupMessage(context.Background(), d.User.ID, i.Token, messageID)
}

// EditFollowup edits a follow-up message.
func (d *Client) EditFollowup(i *interactions.Interaction, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	return d.RestClient.EditFollowupMessage(context.Background(), d.User.ID, i.Token, messageID, params)
}

// DeleteFollowup deletes a follow-up message.
func (d *Client) DeleteFollowup(i *interactions.Interaction, messageID common.Snowflake) error {
	return d.RestClient.DeleteFollowupMessage(context.Background(), d.User.ID, i.Token, messageID)
}

// ── Message methods ───────────────────────────────────────────────────────────

func (d *Client) GetMessages(channelID common.Snowflake, params api.GetMessagesParams) ([]*common.Message, error) {
	return d.RestClient.GetMessages(context.Background(), channelID, params)
}

func (d *Client) GetMessage(channelID, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.GetMessage(context.Background(), channelID, messageID)
}

func (d *Client) SendMessage(channelID common.Snowflake, params api.CreateMessageParams) (*common.Message, error) {
	return d.RestClient.CreateMessage(context.Background(), channelID, params)
}

func (d *Client) EditMessage(channelID, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	return d.RestClient.EditMessage(context.Background(), channelID, messageID, params)
}

func (d *Client) DeleteMessage(channelID, messageID common.Snowflake) error {
	return d.RestClient.DeleteMessage(context.Background(), channelID, messageID)
}

func (d *Client) BulkDeleteMessages(channelID common.Snowflake, messageIDs []common.Snowflake) error {
	return d.RestClient.BulkDeleteMessages(context.Background(), channelID, messageIDs)
}

func (d *Client) CrosspostMessage(channelID, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.CrosspostMessage(context.Background(), channelID, messageID)
}

func (d *Client) GetPinnedMessages(channelID common.Snowflake) ([]*common.Message, error) {
	return d.RestClient.GetPinnedMessages(context.Background(), channelID)
}

func (d *Client) PinMessage(channelID, messageID common.Snowflake) error {
	return d.RestClient.PinMessage(context.Background(), channelID, messageID)
}

func (d *Client) UnpinMessage(channelID, messageID common.Snowflake) error {
	return d.RestClient.UnpinMessage(context.Background(), channelID, messageID)
}

func (d *Client) AddReaction(channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.AddReaction(context.Background(), channelID, messageID, emoji)
}

func (d *Client) DeleteOwnReaction(channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.DeleteOwnReaction(context.Background(), channelID, messageID, emoji)
}

func (d *Client) DeleteUserReaction(channelID, messageID common.Snowflake, emoji string, userID common.Snowflake) error {
	return d.RestClient.DeleteUserReaction(context.Background(), channelID, messageID, emoji, userID)
}

func (d *Client) GetReactions(channelID, messageID common.Snowflake, emoji string, params api.GetReactionsParams) ([]*common.User, error) {
	return d.RestClient.GetReactions(context.Background(), channelID, messageID, emoji, params)
}

func (d *Client) DeleteAllReactions(channelID, messageID common.Snowflake) error {
	return d.RestClient.DeleteAllReactions(context.Background(), channelID, messageID)
}

func (d *Client) DeleteAllReactionsForEmoji(channelID, messageID common.Snowflake, emoji string) error {
	return d.RestClient.DeleteAllReactionsForEmoji(context.Background(), channelID, messageID, emoji)
}

// ── Channel methods ───────────────────────────────────────────────────────────

func (d *Client) GetChannel(channelID common.Snowflake) (*common.Channel, error) {
	return d.RestClient.GetChannel(context.Background(), channelID)
}

func (d *Client) ModifyChannel(channelID common.Snowflake, params api.ModifyChannelParams) (*common.Channel, error) {
	return d.RestClient.ModifyChannel(context.Background(), channelID, params)
}

func (d *Client) DeleteChannel(channelID common.Snowflake) (*common.Channel, error) {
	return d.RestClient.DeleteChannel(context.Background(), channelID)
}

func (d *Client) GetChannelInvites(channelID common.Snowflake) ([]*api.Invite, error) {
	return d.RestClient.GetChannelInvites(context.Background(), channelID)
}

func (d *Client) CreateChannelInvite(channelID common.Snowflake, params api.CreateChannelInviteParams) (*api.Invite, error) {
	return d.RestClient.CreateChannelInvite(context.Background(), channelID, params)
}

func (d *Client) EditChannelPermissions(channelID, overwriteID common.Snowflake, params api.EditChannelPermissionsParams) error {
	return d.RestClient.EditChannelPermissions(context.Background(), channelID, overwriteID, params)
}

func (d *Client) DeleteChannelPermission(channelID, overwriteID common.Snowflake) error {
	return d.RestClient.DeleteChannelPermission(context.Background(), channelID, overwriteID)
}

func (d *Client) TriggerTypingIndicator(channelID common.Snowflake) error {
	return d.RestClient.TriggerTypingIndicator(context.Background(), channelID)
}

// ── Guild methods ─────────────────────────────────────────────────────────────

func (d *Client) GetGuild(guildID common.Snowflake, withCounts bool) (*common.Guild, error) {
	return d.RestClient.GetGuild(context.Background(), guildID, withCounts)
}

func (d *Client) GetGuildPreview(guildID common.Snowflake) (*common.Guild, error) {
	return d.RestClient.GetGuildPreview(context.Background(), guildID)
}

func (d *Client) ModifyGuild(guildID common.Snowflake, params api.ModifyGuildParams) (*common.Guild, error) {
	return d.RestClient.ModifyGuild(context.Background(), guildID, params)
}

func (d *Client) DeleteGuild(guildID common.Snowflake) error {
	return d.RestClient.DeleteGuild(context.Background(), guildID)
}

func (d *Client) GetGuildChannels(guildID common.Snowflake) ([]*common.Channel, error) {
	return d.RestClient.GetGuildChannels(context.Background(), guildID)
}

func (d *Client) CreateGuildChannel(guildID common.Snowflake, params api.CreateGuildChannelParams) (*common.Channel, error) {
	return d.RestClient.CreateGuildChannel(context.Background(), guildID, params)
}

func (d *Client) ModifyGuildChannelPositions(guildID common.Snowflake, entries []api.ModifyGuildChannelPositionsEntry) error {
	return d.RestClient.ModifyGuildChannelPositions(context.Background(), guildID, entries)
}

func (d *Client) GetGuildRoles(guildID common.Snowflake) ([]*common.Role, error) {
	return d.RestClient.GetGuildRoles(context.Background(), guildID)
}

func (d *Client) GetGuildRole(guildID, roleID common.Snowflake) (*common.Role, error) {
	return d.RestClient.GetGuildRole(context.Background(), guildID, roleID)
}

func (d *Client) CreateGuildRole(guildID common.Snowflake, params api.CreateGuildRoleParams) (*common.Role, error) {
	return d.RestClient.CreateGuildRole(context.Background(), guildID, params)
}

func (d *Client) ModifyGuildRolePositions(guildID common.Snowflake, entries []api.ModifyGuildRolePositionsEntry) ([]*common.Role, error) {
	return d.RestClient.ModifyGuildRolePositions(context.Background(), guildID, entries)
}

func (d *Client) ModifyGuildRole(guildID, roleID common.Snowflake, params api.ModifyGuildRoleParams) (*common.Role, error) {
	return d.RestClient.ModifyGuildRole(context.Background(), guildID, roleID, params)
}

func (d *Client) DeleteGuildRole(guildID, roleID common.Snowflake) error {
	return d.RestClient.DeleteGuildRole(context.Background(), guildID, roleID)
}

func (d *Client) GetGuildBans(guildID common.Snowflake, params api.GetGuildBansParams) ([]*api.Ban, error) {
	return d.RestClient.GetGuildBans(context.Background(), guildID, params)
}

func (d *Client) GetGuildBan(guildID, userID common.Snowflake) (*api.Ban, error) {
	return d.RestClient.GetGuildBan(context.Background(), guildID, userID)
}

func (d *Client) CreateGuildBan(guildID, userID common.Snowflake, params api.CreateGuildBanParams) error {
	return d.RestClient.CreateGuildBan(context.Background(), guildID, userID, params)
}

func (d *Client) RemoveGuildBan(guildID, userID common.Snowflake) error {
	return d.RestClient.RemoveGuildBan(context.Background(), guildID, userID)
}

func (d *Client) GetGuildPruneCount(guildID common.Snowflake, params api.GetGuildPruneCountParams) (*api.GuildPruneCountResult, error) {
	return d.RestClient.GetGuildPruneCount(context.Background(), guildID, params)
}

func (d *Client) BeginGuildPrune(guildID common.Snowflake, params api.BeginGuildPruneParams) (*api.GuildPruneCountResult, error) {
	return d.RestClient.BeginGuildPrune(context.Background(), guildID, params)
}

func (d *Client) GetGuildInvites(guildID common.Snowflake) ([]*api.Invite, error) {
	return d.RestClient.GetGuildInvites(context.Background(), guildID)
}

func (d *Client) GetGuildVanityURL(guildID common.Snowflake) (*api.GuildVanityURL, error) {
	return d.RestClient.GetGuildVanityURL(context.Background(), guildID)
}

func (d *Client) GetGuildAuditLog(guildID common.Snowflake, params api.GetGuildAuditLogParams) (*api.AuditLog, error) {
	return d.RestClient.GetGuildAuditLog(context.Background(), guildID, params)
}

// ── Member methods ────────────────────────────────────────────────────────────

func (d *Client) GetGuildMember(guildID, userID common.Snowflake) (*common.GuildMember, error) {
	return d.RestClient.GetGuildMember(context.Background(), guildID, userID)
}

func (d *Client) ListGuildMembers(guildID common.Snowflake, params api.GetGuildMembersParams) ([]*common.GuildMember, error) {
	return d.RestClient.ListGuildMembers(context.Background(), guildID, params)
}

func (d *Client) SearchGuildMembers(guildID common.Snowflake, params api.SearchGuildMembersParams) ([]*common.GuildMember, error) {
	return d.RestClient.SearchGuildMembers(context.Background(), guildID, params)
}

func (d *Client) ModifyGuildMember(guildID, userID common.Snowflake, params api.ModifyGuildMemberParams) (*common.GuildMember, error) {
	return d.RestClient.ModifyGuildMember(context.Background(), guildID, userID, params)
}

func (d *Client) ModifyCurrentMember(guildID common.Snowflake, params api.ModifyCurrentMemberParams) (*common.GuildMember, error) {
	return d.RestClient.ModifyCurrentMember(context.Background(), guildID, params)
}

func (d *Client) AddGuildMemberRole(guildID, userID, roleID common.Snowflake) error {
	return d.RestClient.AddGuildMemberRole(context.Background(), guildID, userID, roleID)
}

func (d *Client) RemoveGuildMemberRole(guildID, userID, roleID common.Snowflake) error {
	return d.RestClient.RemoveGuildMemberRole(context.Background(), guildID, userID, roleID)
}

// RemoveGuildMember removes a member from a guild. Requires KICK_MEMBERS.
func (d *Client) RemoveGuildMember(guildID, userID common.Snowflake) error {
	return d.RestClient.RemoveGuildMember(context.Background(), guildID, userID)
}

// KickGuildMember removes a member from a guild. Requires KICK_MEMBERS.
func (d *Client) KickGuildMember(guildID, userID common.Snowflake) error {
	return d.RestClient.RemoveGuildMember(context.Background(), guildID, userID)
}

// FetchAllGuildMembers retrieves every member in a guild by paginating the REST
// API with 1 000-member pages until all members are returned.
//
// Requires the GUILD_MEMBERS privileged intent to be enabled both in the Discord
// developer portal and in the intents passed to NewClient.
func (d *Client) FetchAllGuildMembers(guildID common.Snowflake) ([]*common.GuildMember, error) {
	const pageSize = 1000
	limit := pageSize
	var all []*common.GuildMember
	var after *common.Snowflake

	for {
		batch, err := d.RestClient.ListGuildMembers(context.Background(), guildID, api.GetGuildMembersParams{
			Limit: &limit,
			After: after,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < pageSize {
			break
		}
		last := batch[len(batch)-1]
		if last.User == nil {
			break
		}
		after = &last.User.ID
	}
	return all, nil
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
	return d.Websocket.Connection.WriteJSON(map[string]interface{}{
		"op": 8,
		"d":  data,
	})
}

// ── Guild widget ──────────────────────────────────────────────────────────────

func (d *Client) GetGuildWidgetSettings(guildID common.Snowflake) (*common.GuildWidgetSettings, error) {
	return d.RestClient.GetGuildWidgetSettings(context.Background(), guildID)
}

func (d *Client) ModifyGuildWidgetSettings(guildID common.Snowflake, params api.ModifyGuildWidgetParams) (*common.GuildWidgetSettings, error) {
	return d.RestClient.ModifyGuildWidgetSettings(context.Background(), guildID, params)
}

func (d *Client) GetGuildWidget(guildID common.Snowflake) (*common.GuildWidget, error) {
	return d.RestClient.GetGuildWidget(context.Background(), guildID)
}

// ── Welcome screen ────────────────────────────────────────────────────────────

func (d *Client) GetGuildWelcomeScreen(guildID common.Snowflake) (*common.GuildWelcomeScreen, error) {
	return d.RestClient.GetGuildWelcomeScreen(context.Background(), guildID)
}

func (d *Client) ModifyGuildWelcomeScreen(guildID common.Snowflake, params api.ModifyGuildWelcomeScreenParams) (*common.GuildWelcomeScreen, error) {
	return d.RestClient.ModifyGuildWelcomeScreen(context.Background(), guildID, params)
}

// ── Guild onboarding ──────────────────────────────────────────────────────────

func (d *Client) GetGuildOnboarding(guildID common.Snowflake) (*common.GuildOnboarding, error) {
	return d.RestClient.GetGuildOnboarding(context.Background(), guildID)
}

func (d *Client) ModifyGuildOnboarding(guildID common.Snowflake, params api.ModifyGuildOnboardingParams) (*common.GuildOnboarding, error) {
	return d.RestClient.ModifyGuildOnboarding(context.Background(), guildID, params)
}

// ── Voice ─────────────────────────────────────────────────────────────────────

func (d *Client) ListVoiceRegions() ([]*common.VoiceRegion, error) {
	return d.RestClient.ListVoiceRegions(context.Background())
}

func (d *Client) ListGuildVoiceRegions(guildID common.Snowflake) ([]*common.VoiceRegion, error) {
	return d.RestClient.ListGuildVoiceRegions(context.Background(), guildID)
}

func (d *Client) ModifyCurrentUserVoiceState(guildID common.Snowflake, params api.ModifyCurrentUserVoiceStateParams) error {
	return d.RestClient.ModifyCurrentUserVoiceState(context.Background(), guildID, params)
}

func (d *Client) ModifyUserVoiceState(guildID, userID common.Snowflake, params api.ModifyUserVoiceStateParams) error {
	return d.RestClient.ModifyUserVoiceState(context.Background(), guildID, userID, params)
}

// ── Soundboard ────────────────────────────────────────────────────────────────

func (d *Client) ListDefaultSoundboardSounds() ([]*common.SoundboardSound, error) {
	return d.RestClient.ListDefaultSoundboardSounds(context.Background())
}

func (d *Client) ListGuildSoundboardSounds(guildID common.Snowflake) ([]*common.SoundboardSound, error) {
	return d.RestClient.ListGuildSoundboardSounds(context.Background(), guildID)
}

func (d *Client) GetGuildSoundboardSound(guildID, soundID common.Snowflake) (*common.SoundboardSound, error) {
	return d.RestClient.GetGuildSoundboardSound(context.Background(), guildID, soundID)
}

func (d *Client) CreateGuildSoundboardSound(guildID common.Snowflake, params api.CreateGuildSoundboardSoundParams) (*common.SoundboardSound, error) {
	return d.RestClient.CreateGuildSoundboardSound(context.Background(), guildID, params)
}

func (d *Client) ModifyGuildSoundboardSound(guildID, soundID common.Snowflake, params api.ModifyGuildSoundboardSoundParams) (*common.SoundboardSound, error) {
	return d.RestClient.ModifyGuildSoundboardSound(context.Background(), guildID, soundID, params)
}

func (d *Client) DeleteGuildSoundboardSound(guildID, soundID common.Snowflake) error {
	return d.RestClient.DeleteGuildSoundboardSound(context.Background(), guildID, soundID)
}

func (d *Client) SendSoundboardSound(channelID common.Snowflake, params api.SendSoundboardSoundParams) error {
	return d.RestClient.SendSoundboardSound(context.Background(), channelID, params)
}

// ── Application ───────────────────────────────────────────────────────────────

func (d *Client) GetCurrentApplication() (*common.Application, error) {
	return d.RestClient.GetCurrentApplication(context.Background())
}

func (d *Client) ModifyCurrentApplication(params api.ModifyCurrentApplicationParams) (*common.Application, error) {
	return d.RestClient.ModifyCurrentApplication(context.Background(), params)
}

// ── Command permissions ───────────────────────────────────────────────────────

// GetGuildApplicationCommandPermissions returns all permission overrides for every command in a guild.
// Uses the bot's own application ID automatically.
func (d *Client) GetGuildApplicationCommandPermissions(guildID common.Snowflake) ([]*common.GuildApplicationCommandPermissions, error) {
	return d.RestClient.GetGuildApplicationCommandPermissions(context.Background(), d.User.ID, guildID)
}

// GetApplicationCommandPermissions returns the permission overrides for a specific command.
func (d *Client) GetApplicationCommandPermissions(guildID, cmdID common.Snowflake) (*common.GuildApplicationCommandPermissions, error) {
	return d.RestClient.GetApplicationCommandPermissions(context.Background(), d.User.ID, guildID, cmdID)
}

// ── Application emojis ────────────────────────────────────────────────────────

func (d *Client) ListApplicationEmojis(appID common.Snowflake) ([]*common.Emoji, error) {
	return d.RestClient.ListApplicationEmojis(context.Background(), appID)
}

func (d *Client) GetApplicationEmoji(appID, emojiID common.Snowflake) (*common.Emoji, error) {
	return d.RestClient.GetApplicationEmoji(context.Background(), appID, emojiID)
}

func (d *Client) CreateApplicationEmoji(appID common.Snowflake, params api.CreateEmojiParams) (*common.Emoji, error) {
	return d.RestClient.CreateApplicationEmoji(context.Background(), appID, params)
}

func (d *Client) ModifyApplicationEmoji(appID, emojiID common.Snowflake, params api.ModifyEmojiParams) (*common.Emoji, error) {
	return d.RestClient.ModifyApplicationEmoji(context.Background(), appID, emojiID, params)
}

func (d *Client) DeleteApplicationEmoji(appID, emojiID common.Snowflake) error {
	return d.RestClient.DeleteApplicationEmoji(context.Background(), appID, emojiID)
}

// ── Entitlements ──────────────────────────────────────────────────────────────

func (d *Client) ListEntitlements(appID common.Snowflake, params api.ListEntitlementsParams) ([]*common.Entitlement, error) {
	return d.RestClient.ListEntitlements(context.Background(), appID, params)
}

func (d *Client) GetEntitlement(appID, entitlementID common.Snowflake) (*common.Entitlement, error) {
	return d.RestClient.GetEntitlement(context.Background(), appID, entitlementID)
}

func (d *Client) CreateTestEntitlement(appID common.Snowflake, params api.CreateTestEntitlementParams) (*common.Entitlement, error) {
	return d.RestClient.CreateTestEntitlement(context.Background(), appID, params)
}

func (d *Client) ConsumeEntitlement(appID, entitlementID common.Snowflake) error {
	return d.RestClient.ConsumeEntitlement(context.Background(), appID, entitlementID)
}

func (d *Client) DeleteTestEntitlement(appID, entitlementID common.Snowflake) error {
	return d.RestClient.DeleteTestEntitlement(context.Background(), appID, entitlementID)
}

// ── SKUs ──────────────────────────────────────────────────────────────────────

func (d *Client) ListSKUs(appID common.Snowflake) ([]*common.SKU, error) {
	return d.RestClient.ListSKUs(context.Background(), appID)
}

// ── Subscriptions ─────────────────────────────────────────────────────────────

func (d *Client) ListSKUSubscriptions(skuID common.Snowflake, params api.ListSKUSubscriptionsParams) ([]*common.Subscription, error) {
	return d.RestClient.ListSKUSubscriptions(context.Background(), skuID, params)
}

func (d *Client) GetSKUSubscription(skuID, subscriptionID common.Snowflake) (*common.Subscription, error) {
	return d.RestClient.GetSKUSubscription(context.Background(), skuID, subscriptionID)
}

// ── Poll ──────────────────────────────────────────────────────────────────────

func (d *Client) GetPollAnswerVoters(channelID, messageID common.Snowflake, answerID int, params api.GetPollAnswerVotersParams) ([]*common.User, error) {
	return d.RestClient.GetPollAnswerVoters(context.Background(), channelID, messageID, answerID, params)
}

func (d *Client) EndPoll(channelID, messageID common.Snowflake) (*common.Message, error) {
	return d.RestClient.EndPoll(context.Background(), channelID, messageID)
}

// ── Activity ──────────────────────────────────────────────────────────────────

func (d *Client) GetActivityInstance(appID common.Snowflake, instanceID string) (*common.ActivityInstance, error) {
	return d.RestClient.GetActivityInstance(context.Background(), appID, instanceID)
}

// ── Additional channel methods ────────────────────────────────────────────────

func (d *Client) SetVoiceChannelStatus(channelID common.Snowflake, status *string) error {
	return d.RestClient.SetVoiceChannelStatus(context.Background(), channelID, status)
}

func (d *Client) FollowAnnouncementChannel(channelID, webhookChannelID common.Snowflake) (*api.FollowedChannel, error) {
	return d.RestClient.FollowAnnouncementChannel(context.Background(), channelID, webhookChannelID)
}

func (d *Client) AddGroupDMRecipient(channelID, userID common.Snowflake, params api.AddGroupDMRecipientParams) error {
	return d.RestClient.AddGroupDMRecipient(context.Background(), channelID, userID, params)
}

func (d *Client) RemoveGroupDMRecipient(channelID, userID common.Snowflake) error {
	return d.RestClient.RemoveGroupDMRecipient(context.Background(), channelID, userID)
}

// ── Additional guild methods ──────────────────────────────────────────────────

func (d *Client) CreateGuild(params api.CreateGuildParams) (*common.Guild, error) {
	return d.RestClient.CreateGuild(context.Background(), params)
}

func (d *Client) AddGuildMember(guildID, userID common.Snowflake, params api.AddGuildMemberParams) (*common.GuildMember, error) {
	return d.RestClient.AddGuildMember(context.Background(), guildID, userID, params)
}

func (d *Client) BulkBanGuildMembers(guildID common.Snowflake, params api.BulkBanParams) (*api.BulkBanResult, error) {
	return d.RestClient.BulkBanGuildMembers(context.Background(), guildID, params)
}

func (d *Client) GetGuildIntegrations(guildID common.Snowflake) ([]*common.Integration, error) {
	return d.RestClient.GetGuildIntegrations(context.Background(), guildID)
}

func (d *Client) DeleteGuildIntegration(guildID, integrationID common.Snowflake) error {
	return d.RestClient.DeleteGuildIntegration(context.Background(), guildID, integrationID)
}

// ── Additional voice methods ──────────────────────────────────────────────────

func (d *Client) GetCurrentUserVoiceState(guildID common.Snowflake) (*common.VoiceState, error) {
	return d.RestClient.GetCurrentUserVoiceState(context.Background(), guildID)
}

func (d *Client) GetUserVoiceState(guildID, userID common.Snowflake) (*common.VoiceState, error) {
	return d.RestClient.GetUserVoiceState(context.Background(), guildID, userID)
}

// ── Additional sticker methods ────────────────────────────────────────────────

func (d *Client) GetStickerPack(packID common.Snowflake) (*api.StickerPack, error) {
	return d.RestClient.GetStickerPack(context.Background(), packID)
}

func (d *Client) CreateGuildSticker(guildID common.Snowflake, params api.CreateGuildStickerParams) (*common.Sticker, error) {
	return d.RestClient.CreateGuildSticker(context.Background(), guildID, params)
}

// ── Invite methods ────────────────────────────────────────────────────────────

func (d *Client) GetInvite(code string, params api.GetInviteParams) (*api.Invite, error) {
	return d.RestClient.GetInvite(context.Background(), code, params)
}

func (d *Client) DeleteInvite(code string) (*api.Invite, error) {
	return d.RestClient.DeleteInvite(context.Background(), code)
}

// ── Additional user methods ───────────────────────────────────────────────────

func (d *Client) GetCurrentUserConnections() ([]*api.UserConnection, error) {
	return d.RestClient.GetCurrentUserConnections(context.Background())
}

func (d *Client) GetCurrentUserApplicationRoleConnection(appID common.Snowflake) (*api.ApplicationRoleConnection, error) {
	return d.RestClient.GetCurrentUserApplicationRoleConnection(context.Background(), appID)
}

func (d *Client) UpdateCurrentUserApplicationRoleConnection(appID common.Snowflake, params api.UpdateRoleConnectionParams) (*api.ApplicationRoleConnection, error) {
	return d.RestClient.UpdateCurrentUserApplicationRoleConnection(context.Background(), appID, params)
}

func (d *Client) CreateGroupDM(params api.CreateGroupDMParams) (*common.Channel, error) {
	return d.RestClient.CreateGroupDM(context.Background(), params)
}

// ── Additional webhook methods ────────────────────────────────────────────────

func (d *Client) ExecuteSlackWebhook(webhookID common.Snowflake, token string, wait bool, body json.RawMessage) error {
	return d.RestClient.ExecuteSlackWebhook(context.Background(), webhookID, token, wait, body)
}

func (d *Client) ExecuteGitHubWebhook(webhookID common.Snowflake, token string, wait bool, body json.RawMessage) error {
	return d.RestClient.ExecuteGitHubWebhook(context.Background(), webhookID, token, wait, body)
}
