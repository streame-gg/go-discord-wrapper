package discord

import "context"

// EntityClient is the subset of *connection.Client required by the entity
// convenience methods (e.g. msg.Edit, member.Ban, guild.CreateRole).
//
// *connection.Client satisfies this interface — pass the client received in
// your event handler directly.
//
// Entity methods accept ctx per call; the entity itself does not store a
// context (unlike Interaction which captures the dispatch context).
type EntityClient interface {
	// ── Message ───────────────────────────────────────────────────────────
	EditMessage(ctx context.Context, channelID, messageID Snowflake, opts MessageEditOptions) (*Message, error)
	DeleteMessage(ctx context.Context, channelID, messageID Snowflake, reason *string) error
	CreateMessage(ctx context.Context, channelID Snowflake, opts MessageCreateOptions) (*Message, error)
	PinMessage(ctx context.Context, channelID, messageID Snowflake, reason *string) error
	UnpinMessage(ctx context.Context, channelID, messageID Snowflake, reason *string) error
	CrosspostMessage(ctx context.Context, channelID, messageID Snowflake) (*Message, error)
	AddReaction(ctx context.Context, channelID, messageID Snowflake, emoji string) error

	// ── Channel ───────────────────────────────────────────────────────────
	ModifyChannel(ctx context.Context, channelID Snowflake, opts ChannelEditOptions) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID Snowflake, reason *string) (*Channel, error)
	BulkDeleteMessages(ctx context.Context, channelID Snowflake, messageIDs []Snowflake, reason *string) error
	GetChannelMessages(ctx context.Context, channelID Snowflake, opts FetchMessagesOptions) ([]*Message, error)
	TriggerTypingIndicator(ctx context.Context, channelID Snowflake) error
	SetVoiceChannelStatus(ctx context.Context, channelID Snowflake, status *string) error
	CreateChannelInvite(ctx context.Context, channelID Snowflake, opts InviteCreateOptions) (*Invite, error)

	// ── Guild ─────────────────────────────────────────────────────────────
	ModifyGuild(ctx context.Context, guildID Snowflake, opts GuildEditOptions) (*Guild, error)
	LeaveGuild(ctx context.Context, guildID Snowflake) error
	ListGuildMembers(ctx context.Context, guildID Snowflake, opts FetchMembersOptions) ([]*GuildMember, error)
	CreateGuildRole(ctx context.Context, guildID Snowflake, opts RoleCreateOptions) (*Role, error)
	CreateGuildChannel(ctx context.Context, guildID Snowflake, opts ChannelCreateOptions) (*Channel, error)
	CreateGuildEmoji(ctx context.Context, guildID Snowflake, opts EmojiCreateOptions) (*Emoji, error)
	CreateGuildSticker(ctx context.Context, guildID Snowflake, opts StickerCreateOptions) (*Sticker, error)
	GetGuildAuditLog(ctx context.Context, guildID Snowflake, opts AuditLogOptions) (*AuditLog, error)
	CreateGuildScheduledEvent(ctx context.Context, guildID Snowflake, opts ScheduledEventCreateOptions) (*GuildScheduledEvent, error)

	// ── GuildMember ───────────────────────────────────────────────────────
	ModifyGuildMember(ctx context.Context, guildID, userID Snowflake, opts MemberEditOptions) (*GuildMember, error)
	KickGuildMember(ctx context.Context, guildID, userID Snowflake, reason *string) error
	CreateGuildBan(ctx context.Context, guildID, userID Snowflake, opts BanOptions) error
	AddGuildMemberRole(ctx context.Context, guildID, userID, roleID Snowflake, reason *string) error
	RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID Snowflake, reason *string) error

	// ── Role ──────────────────────────────────────────────────────────────
	ModifyGuildRole(ctx context.Context, guildID, roleID Snowflake, opts RoleEditOptions) (*Role, error)
	DeleteGuildRole(ctx context.Context, guildID, roleID Snowflake, reason *string) error
	ModifyGuildRolePositions(ctx context.Context, guildID Snowflake, opts RolePositionOptions) ([]*Role, error)

	// ── User ──────────────────────────────────────────────────────────────
	GetUser(ctx context.Context, userID Snowflake) (*User, error)
	CreateDM(ctx context.Context, recipientID Snowflake) (*Channel, error)

	// ── Emoji ─────────────────────────────────────────────────────────────
	ModifyGuildEmoji(ctx context.Context, guildID, emojiID Snowflake, opts EmojiEditOptions) (*Emoji, error)
	DeleteGuildEmoji(ctx context.Context, guildID, emojiID Snowflake, reason *string) error

	// ── Webhook ───────────────────────────────────────────────────────────
	ModifyWebhook(ctx context.Context, webhookID Snowflake, opts WebhookEditOptions) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID Snowflake, reason *string) error
	ExecuteWebhook(ctx context.Context, webhookID Snowflake, token string, opts WebhookExecuteOptions) (*Message, error)
	GetWebhookMessage(ctx context.Context, webhookID Snowflake, token string, messageID Snowflake) (*Message, error)

	// ── Invite ────────────────────────────────────────────────────────────
	DeleteInvite(ctx context.Context, code string, reason *string) (*Invite, error)

	// ── StageInstance ─────────────────────────────────────────────────────
	ModifyStageInstance(ctx context.Context, channelID Snowflake, opts StageEditOptions) (*StageInstance, error)
	DeleteStageInstance(ctx context.Context, channelID Snowflake, reason *string) error

	// ── GuildScheduledEvent ───────────────────────────────────────────────
	ModifyGuildScheduledEvent(ctx context.Context, guildID, eventID Snowflake, opts ScheduledEventEditOptions) (*GuildScheduledEvent, error)
	DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID Snowflake) error
	GetGuildScheduledEventUsers(ctx context.Context, guildID, eventID Snowflake, opts FetchUsersOptions) ([]*GuildScheduledEventUser, error)

	// ── Sticker ───────────────────────────────────────────────────────────
	ModifyGuildSticker(ctx context.Context, guildID, stickerID Snowflake, opts StickerEditOptions) (*Sticker, error)
	DeleteGuildSticker(ctx context.Context, guildID, stickerID Snowflake, reason *string) error

	// ── AutoModerationRule ────────────────────────────────────────────────
	ModifyAutoModerationRule(ctx context.Context, guildID, ruleID Snowflake, opts RuleEditOptions) (*AutoModerationRule, error)
	DeleteAutoModerationRule(ctx context.Context, guildID, ruleID Snowflake) error

	// ── SoundboardSound ───────────────────────────────────────────────────
	ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID Snowflake, opts SoundEditOptions) (*SoundboardSound, error)
	DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID Snowflake, reason *string) error

	// ── Integration ───────────────────────────────────────────────────────
	DeleteGuildIntegration(ctx context.Context, guildID, integrationID Snowflake, reason *string) error
}
