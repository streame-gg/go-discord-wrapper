package connection

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// Compile-time assertion that *Client satisfies discord.EntityClient.
var _ discord.EntityClient = (*Client)(nil)

// ── Message ────────────────────────────────────────────────────────────────────

func (d *Client) EditMessage(ctx context.Context, channelID, messageID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error) {
	params := api.EditMessageParams{
		Content:         opts.Content,
		Embeds:          opts.Embeds,
		Flags:           opts.Flags,
		AllowedMentions: opts.AllowedMentions,
		Components:      opts.Components,
		Attachments:     opts.Attachments,
		Files:           opts.Files,
	}
	msg, err := d.RestClient.EditMessage(ctx, channelID, messageID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) CreateMessage(ctx context.Context, channelID discord.Snowflake, opts discord.MessageCreateOptions) (*discord.Message, error) {
	params := api.CreateMessageParams{
		Content:          opts.Content,
		TTS:              opts.TTS,
		Embeds:           opts.Embeds,
		AllowedMentions:  opts.AllowedMentions,
		MessageReference: opts.MessageReference,
		Components:       opts.Components,
		StickerIDs:       opts.StickerIDs,
		Flags:            opts.Flags,
		Files:            opts.Files,
	}
	msg, err := d.RestClient.CreateMessage(ctx, channelID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) GetChannelMessages(ctx context.Context, channelID discord.Snowflake, opts discord.FetchMessagesOptions) ([]*discord.Message, error) {
	params := api.GetMessagesParams{
		Around: opts.Around,
		Before: opts.Before,
		After:  opts.After,
		Limit:  opts.Limit,
	}
	msgs, err := d.RestClient.GetMessages(ctx, channelID, params)
	if err == nil {
		d.cacheMessages(msgs)
	}
	return msgs, err
}

// ── Channel ────────────────────────────────────────────────────────────────────

func (d *Client) ModifyChannel(ctx context.Context, channelID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error) {
	params := api.ModifyChannelParams{
		Name:                          opts.Name,
		Type:                          opts.Type,
		Position:                      opts.Position,
		Topic:                         opts.Topic,
		NSFW:                          opts.NSFW,
		RateLimitPerUser:              opts.RateLimitPerUser,
		Bitrate:                       opts.Bitrate,
		UserLimit:                     opts.UserLimit,
		PermissionOverwrites:          opts.PermissionOverwrites,
		ParentID:                      opts.ParentID,
		RTCRegion:                     opts.RTCRegion,
		VideoQualityMode:              opts.VideoQualityMode,
		DefaultAutoArchiveDuration:    opts.DefaultAutoArchiveDuration,
		Flags:                         opts.Flags,
		AvailableTags:                 opts.AvailableTags,
		DefaultReactionEmoji:          opts.DefaultReactionEmoji,
		DefaultThreadRateLimitPerUser: opts.DefaultThreadRateLimitPerUser,
		DefaultSortOrder:              opts.DefaultSortOrder,
		DefaultForumLayout:            opts.DefaultForumLayout,
	}
	var apiOpts *api.ModifyChannelOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyChannelOptions{Reason: *opts.AuditLogReason}
	}
	ch, err := d.RestClient.ModifyChannel(ctx, channelID, params, apiOpts)
	if err == nil {
		d.cacheChannel(ch)
	}
	return ch, err
}

func (d *Client) CreateChannelInvite(ctx context.Context, channelID discord.Snowflake, opts discord.InviteCreateOptions) (*discord.Invite, error) {
	params := api.CreateChannelInviteParams{
		MaxAge:              opts.MaxAge,
		MaxUses:             opts.MaxUses,
		Temporary:           opts.Temporary,
		Unique:              opts.Unique,
		TargetType:          opts.TargetType,
		TargetUserID:        opts.TargetUserID,
		TargetApplicationID: opts.TargetApplicationID,
	}
	var apiOpts *api.CreateChannelInviteOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateChannelInviteOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.CreateChannelInvite(ctx, channelID, params, apiOpts)
}

// ── Guild ──────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuild(ctx context.Context, guildID discord.Snowflake, opts discord.GuildEditOptions) (*discord.Guild, error) {
	params := api.ModifyGuildParams{
		Name:                        opts.Name,
		VerificationLevel:           opts.VerificationLevel,
		DefaultMessageNotifications: opts.DefaultMessageNotifications,
		ExplicitContentFilter:       opts.ExplicitContentFilter,
		AFKChannelID:                opts.AFKChannelID,
		AFKTimeout:                  opts.AFKTimeout,
		Icon:                        opts.Icon,
		OwnerID:                     opts.OwnerID,
		Splash:                      opts.Splash,
		DiscoverySplash:             opts.DiscoverySplash,
		Banner:                      opts.Banner,
		SystemChannelID:             opts.SystemChannelID,
		SystemChannelFlags:          opts.SystemChannelFlags,
		RulesChannelID:              opts.RulesChannelID,
		PublicUpdatesChannelID:      opts.PublicUpdatesChannelID,
		PreferredLocale:             opts.PreferredLocale,
		Features:                    opts.Features,
		Description:                 opts.Description,
		PremiumProgressBarEnabled:   opts.PremiumProgressBarEnabled,
		SafetyAlertsChannelID:       opts.SafetyAlertsChannelID,
	}
	var apiOpts *api.ModifyGuildOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildOptions{Reason: *opts.AuditLogReason}
	}
	guild, err := d.RestClient.ModifyGuild(ctx, guildID, params, apiOpts)
	if err == nil {
		d.cacheGuild(guild)
	}
	return guild, err
}

func (d *Client) LeaveGuild(ctx context.Context, guildID discord.Snowflake) error {
	return d.RestClient.LeaveGuild(ctx, guildID)
}

func (d *Client) ListGuildMembers(ctx context.Context, guildID discord.Snowflake, opts discord.FetchMembersOptions) ([]*discord.GuildMember, error) {
	params := api.GetGuildMembersParams{
		After: opts.After,
		Limit: opts.Limit,
	}
	members, err := d.RestClient.ListGuildMembers(ctx, guildID, params)
	if err == nil {
		d.cacheMembers(guildID, members)
	}
	return members, err
}

func (d *Client) CreateGuildRole(ctx context.Context, guildID discord.Snowflake, opts discord.RoleCreateOptions) (*discord.Role, error) {
	params := api.CreateGuildRoleParams{
		Name:         opts.Name,
		Permissions:  opts.Permissions,
		Color:        opts.Color,
		Hoist:        opts.Hoist,
		Icon:         opts.Icon,
		UnicodeEmoji: opts.UnicodeEmoji,
		Mentionable:  opts.Mentionable,
	}
	var apiOpts *api.CreateGuildRoleOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateGuildRoleOptions{Reason: *opts.AuditLogReason}
	}
	role, err := d.RestClient.CreateGuildRole(ctx, guildID, params, apiOpts)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) CreateGuildChannel(ctx context.Context, guildID discord.Snowflake, opts discord.ChannelCreateOptions) (*discord.Channel, error) {
	params := api.CreateGuildChannelParams{
		Name:                          opts.Name,
		Type:                          opts.Type,
		Topic:                         opts.Topic,
		Bitrate:                       opts.Bitrate,
		UserLimit:                     opts.UserLimit,
		RateLimitPerUser:              opts.RateLimitPerUser,
		Position:                      opts.Position,
		PermissionOverwrites:          opts.PermissionOverwrites,
		ParentID:                      opts.ParentID,
		NSFW:                          opts.NSFW,
		RTCRegion:                     opts.RTCRegion,
		VideoQualityMode:              opts.VideoQualityMode,
		DefaultAutoArchiveDuration:    opts.DefaultAutoArchiveDuration,
		DefaultReactionEmoji:          opts.DefaultReactionEmoji,
		AvailableTags:                 opts.AvailableTags,
		DefaultSortOrder:              opts.DefaultSortOrder,
		DefaultForumLayout:            opts.DefaultForumLayout,
		DefaultThreadRateLimitPerUser: opts.DefaultThreadRateLimitPerUser,
	}
	var apiOpts *api.CreateGuildChannelOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateGuildChannelOptions{Reason: *opts.AuditLogReason}
	}
	ch, err := d.RestClient.CreateGuildChannel(ctx, guildID, params, apiOpts)
	if err == nil {
		d.cacheChannel(ch)
	}
	return ch, err
}

func (d *Client) CreateGuildEmoji(ctx context.Context, guildID discord.Snowflake, opts discord.EmojiCreateOptions) (*discord.Emoji, error) {
	params := api.CreateGuildEmojiParams{
		Name:  opts.Name,
		Image: opts.Image,
		Roles: opts.Roles,
	}
	var apiOpts *api.CreateGuildEmojiOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateGuildEmojiOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.CreateGuildEmoji(ctx, guildID, params, apiOpts)
}

func (d *Client) CreateGuildSticker(ctx context.Context, guildID discord.Snowflake, opts discord.StickerCreateOptions) (*discord.Sticker, error) {
	params := api.CreateGuildStickerParams{
		Name:        opts.Name,
		Description: opts.Description,
		Tags:        opts.Tags,
		File:        opts.File,
		ContentType: opts.ContentType,
	}
	var apiOpts *api.CreateGuildStickerOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateGuildStickerOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.CreateGuildSticker(ctx, guildID, params, apiOpts)
}

func (d *Client) GetGuildAuditLog(ctx context.Context, guildID discord.Snowflake, opts discord.AuditLogOptions) (*discord.AuditLog, error) {
	params := api.GetGuildAuditLogParams{
		UserID:     opts.UserID,
		ActionType: opts.ActionType,
		Before:     opts.Before,
		After:      opts.After,
		Limit:      opts.Limit,
	}
	return d.RestClient.GetGuildAuditLog(ctx, guildID, params)
}

func (d *Client) CreateGuildScheduledEvent(ctx context.Context, guildID discord.Snowflake, opts discord.ScheduledEventCreateOptions) (*discord.GuildScheduledEvent, error) {
	params := api.CreateGuildScheduledEventParams{
		ChannelID:          opts.ChannelID,
		EntityMetadata:     opts.EntityMetadata,
		Name:               opts.Name,
		PrivacyLevel:       opts.PrivacyLevel,
		ScheduledStartTime: opts.ScheduledStartTime,
		ScheduledEndTime:   opts.ScheduledEndTime,
		Description:        opts.Description,
		EntityType:         opts.EntityType,
		Image:              opts.Image,
	}
	return d.RestClient.CreateGuildScheduledEvent(ctx, guildID, params)
}

// ── GuildMember ────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildMember(ctx context.Context, guildID, userID discord.Snowflake, opts discord.MemberEditOptions) (*discord.GuildMember, error) {
	params := api.ModifyGuildMemberParams{
		Nick:                       opts.Nick,
		Roles:                      opts.Roles,
		Mute:                       opts.Mute,
		Deaf:                       opts.Deaf,
		ChannelID:                  opts.ChannelID,
		CommunicationDisabledUntil: opts.CommunicationDisabledUntil,
		Flags:                      opts.Flags,
	}
	var apiOpts *api.ModifyGuildMemberOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildMemberOptions{Reason: *opts.AuditLogReason}
	}
	member, err := d.RestClient.ModifyGuildMember(ctx, guildID, userID, params, apiOpts)
	if err == nil {
		d.cacheMember(guildID, member)
	}
	return member, err
}

func (d *Client) KickGuildMember(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error {
	var opts *api.KickGuildMemberOptions
	if reason != nil {
		opts = &api.KickGuildMemberOptions{Reason: *reason}
	}
	if err := d.RestClient.KickGuildMember(ctx, guildID, userID, opts); err != nil {
		return err
	}
	d.removeGuildMemberFromCache(guildID, userID)
	return nil
}

func (d *Client) CreateGuildBan(ctx context.Context, guildID, userID discord.Snowflake, opts discord.BanOptions) error {
	params := api.CreateGuildBanParams{
		DeleteMessageSeconds: opts.DeleteMessageSeconds,
	}
	var apiOpts *api.CreateGuildBanOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateGuildBanOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.CreateGuildBan(ctx, guildID, userID, params, apiOpts)
}

// ── Role ───────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, opts discord.RoleEditOptions) (*discord.Role, error) {
	params := api.ModifyGuildRoleParams{
		Name:         opts.Name,
		Permissions:  opts.Permissions,
		Color:        opts.Color,
		Hoist:        opts.Hoist,
		Icon:         opts.Icon,
		UnicodeEmoji: opts.UnicodeEmoji,
		Mentionable:  opts.Mentionable,
	}
	var apiOpts *api.ModifyGuildRoleOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildRoleOptions{Reason: *opts.AuditLogReason}
	}
	role, err := d.RestClient.ModifyGuildRole(ctx, guildID, roleID, params, apiOpts)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) ModifyGuildRolePositions(ctx context.Context, guildID discord.Snowflake, opts discord.RolePositionOptions) ([]*discord.Role, error) {
	entries := make([]api.ModifyGuildRolePositionsEntry, len(opts.Entries))
	for i, e := range opts.Entries {
		entries[i] = api.ModifyGuildRolePositionsEntry{
			ID:       e.ID,
			Position: e.Position,
		}
	}
	var apiOpts *api.ModifyGuildRolePositionsOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildRolePositionsOptions{Reason: *opts.AuditLogReason}
	}
	roles, err := d.RestClient.ModifyGuildRolePositions(ctx, guildID, entries, apiOpts)
	if err == nil {
		if d.Cache != nil {
			d.Cache.Roles().DeleteGuild(guildID)
		}
		d.cacheRoles(guildID, roles)
	}
	return roles, err
}

// ── User ───────────────────────────────────────────────────────────────────────

func (d *Client) GetUser(ctx context.Context, userID discord.Snowflake) (*discord.User, error) {
	user, err := d.RestClient.GetUser(ctx, userID)
	if err == nil {
		d.cacheUser(user)
	}
	return user, err
}

func (d *Client) CreateDM(ctx context.Context, recipientID discord.Snowflake) (*discord.Channel, error) {
	ch, err := d.RestClient.CreateDM(ctx, recipientID)
	if err == nil {
		d.cacheChannel(ch)
	}
	return ch, err
}

// ── Emoji ──────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, opts discord.EmojiEditOptions) (*discord.Emoji, error) {
	params := api.ModifyGuildEmojiParams{
		Name:  opts.Name,
		Roles: opts.Roles,
	}
	var apiOpts *api.ModifyGuildEmojiOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildEmojiOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.ModifyGuildEmoji(ctx, guildID, emojiID, params, apiOpts)
}

func (d *Client) DeleteGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, reason *string) error {
	var opts *api.DeleteGuildEmojiOptions
	if reason != nil {
		opts = &api.DeleteGuildEmojiOptions{Reason: *reason}
	}
	return d.RestClient.DeleteGuildEmoji(ctx, guildID, emojiID, opts)
}

// ── Webhook ────────────────────────────────────────────────────────────────────

func (d *Client) ModifyWebhook(ctx context.Context, webhookID discord.Snowflake, opts discord.WebhookEditOptions) (*discord.Webhook, error) {
	params := api.ModifyWebhookParams{
		Name:      opts.Name,
		Avatar:    opts.Avatar,
		ChannelID: opts.ChannelID,
	}
	var apiOpts *api.ModifyWebhookOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyWebhookOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.ModifyWebhook(ctx, webhookID, params, apiOpts)
}

func (d *Client) DeleteWebhook(ctx context.Context, webhookID discord.Snowflake, reason *string) error {
	var opts *api.DeleteWebhookOptions
	if reason != nil {
		opts = &api.DeleteWebhookOptions{Reason: *reason}
	}
	return d.RestClient.DeleteWebhook(ctx, webhookID, opts)
}

func (d *Client) ExecuteWebhook(ctx context.Context, webhookID discord.Snowflake, token string, opts discord.WebhookExecuteOptions) (*discord.Message, error) {
	params := api.ExecuteWebhookParams{
		Content:         opts.Content,
		Username:        opts.Username,
		AvatarURL:       opts.AvatarURL,
		TTS:             opts.TTS,
		Embeds:          opts.Embeds,
		AllowedMentions: opts.AllowedMentions,
		Components:      opts.Components,
		Flags:           opts.Flags,
		ThreadName:      opts.ThreadName,
		Files:           opts.Files,
	}
	wait := opts.Wait
	query := api.ExecuteWebhookQueryParams{Wait: &wait}
	return d.RestClient.ExecuteWebhook(ctx, webhookID, token, params, query)
}

func (d *Client) GetWebhookMessage(ctx context.Context, webhookID discord.Snowflake, token string, messageID discord.Snowflake) (*discord.Message, error) {
	return d.RestClient.GetWebhookMessage(ctx, webhookID, token, messageID)
}

// ── Invite ─────────────────────────────────────────────────────────────────────

// (DeleteInvite is already defined in clientFunctions.go with reason *string)

// ── StageInstance ──────────────────────────────────────────────────────────────

func (d *Client) ModifyStageInstance(ctx context.Context, channelID discord.Snowflake, opts discord.StageEditOptions) (*discord.StageInstance, error) {
	params := api.ModifyStageInstanceParams{
		Topic:        opts.Topic,
		PrivacyLevel: opts.PrivacyLevel,
	}
	var apiOpts *api.ModifyStageInstanceOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyStageInstanceOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.ModifyStageInstance(ctx, channelID, params, apiOpts)
}

func (d *Client) DeleteStageInstance(ctx context.Context, channelID discord.Snowflake, reason *string) error {
	var opts *api.DeleteStageInstanceOptions
	if reason != nil {
		opts = &api.DeleteStageInstanceOptions{Reason: *reason}
	}
	return d.RestClient.DeleteStageInstance(ctx, channelID, opts)
}

// ── GuildScheduledEvent ────────────────────────────────────────────────────────

func (d *Client) ModifyGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.ScheduledEventEditOptions) (*discord.GuildScheduledEvent, error) {
	params := api.ModifyGuildScheduledEventParams{
		ChannelID:          opts.ChannelID,
		EntityMetadata:     opts.EntityMetadata,
		Name:               opts.Name,
		PrivacyLevel:       opts.PrivacyLevel,
		ScheduledStartTime: opts.ScheduledStartTime,
		ScheduledEndTime:   opts.ScheduledEndTime,
		Description:        opts.Description,
		EntityType:         opts.EntityType,
		Status:             opts.Status,
		Image:              opts.Image,
	}
	return d.RestClient.ModifyGuildScheduledEvent(ctx, guildID, eventID, params)
}

func (d *Client) DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) error {
	return d.RestClient.DeleteGuildScheduledEvent(ctx, guildID, eventID)
}

func (d *Client) GetGuildScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.FetchUsersOptions) ([]*discord.GuildScheduledEventUser, error) {
	params := api.GetGuildScheduledEventUsersParams{
		Limit:      opts.Limit,
		WithMember: opts.WithMember,
		Before:     opts.Before,
		After:      opts.After,
	}
	return d.RestClient.GetGuildScheduledEventUsers(ctx, guildID, eventID, params)
}

// ── Sticker ────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, opts discord.StickerEditOptions) (*discord.Sticker, error) {
	params := api.ModifyGuildStickerParams{
		Name:        opts.Name,
		Description: opts.Description,
		Tags:        opts.Tags,
	}
	var apiOpts *api.ModifyGuildStickerOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildStickerOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.ModifyGuildSticker(ctx, guildID, stickerID, params, apiOpts)
}

func (d *Client) DeleteGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, reason *string) error {
	var opts *api.DeleteGuildStickerOptions
	if reason != nil {
		opts = &api.DeleteGuildStickerOptions{Reason: *reason}
	}
	return d.RestClient.DeleteGuildSticker(ctx, guildID, stickerID, opts)
}

// ── AutoModerationRule ─────────────────────────────────────────────────────────

func (d *Client) ModifyAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake, opts discord.RuleEditOptions) (*discord.AutoModerationRule, error) {
	params := api.ModifyAutoModerationRuleParams{
		Name:            opts.Name,
		EventType:       opts.EventType,
		TriggerMetadata: opts.TriggerMetadata,
		Actions:         opts.Actions,
		Enabled:         opts.Enabled,
		ExemptRoles:     opts.ExemptRoles,
		ExemptChannels:  opts.ExemptChannels,
	}
	var apiOpts *api.ModifyAutoModerationRuleOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyAutoModerationRuleOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.ModifyAutoModerationRule(ctx, guildID, ruleID, params, apiOpts)
}

func (d *Client) DeleteAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) error {
	return d.RestClient.DeleteAutoModerationRule(ctx, guildID, ruleID, nil)
}

// ── SoundboardSound ────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, opts discord.SoundEditOptions) (*discord.SoundboardSound, error) {
	params := api.ModifyGuildSoundboardSoundParams{
		Name:      opts.Name,
		Volume:    opts.Volume,
		EmojiID:   opts.EmojiID,
		EmojiName: opts.EmojiName,
	}
	var apiOpts *api.ModifyGuildSoundboardSoundOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.ModifyGuildSoundboardSoundOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.ModifyGuildSoundboardSound(ctx, guildID, soundID, params, apiOpts)
}

// ── Cache access ───────────────────────────────────────────────────────────────

func (d *Client) ClientCache() discord.Cache {
	if d.Cache == nil {
		return nil
	}
	return cacheAdapter{d.Cache}
}

// ── Fetch (read) — used by sub-managers, not yet in clientFunctions.go ──────

func (d *Client) GetGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake) (*discord.Emoji, error) {
	return d.RestClient.GetGuildEmoji(ctx, guildID, emojiID)
}

func (d *Client) ListGuildEmojis(ctx context.Context, guildID discord.Snowflake) ([]*discord.Emoji, error) {
	emojis, err := d.RestClient.ListGuildEmojis(ctx, guildID)
	if err == nil {
		for _, e := range emojis {
			e.GuildID = guildID
			e.Hydrate(d)
		}
		if d.Cache != nil {
			d.Cache.Emojis().SetAll(guildID, emojis)
		}
	}
	return emojis, err
}

func (d *Client) GetGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake) (*discord.Sticker, error) {
	return d.RestClient.GetGuildSticker(ctx, guildID, stickerID)
}

func (d *Client) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) ([]*discord.Sticker, error) {
	stickers, err := d.RestClient.ListGuildStickers(ctx, guildID)
	if err == nil {
		for _, s := range stickers {
			s.GuildID = &guildID
			s.Hydrate(d)
		}
		if d.Cache != nil {
			d.Cache.Stickers().SetAll(guildID, stickers)
		}
	}
	return stickers, err
}

func (d *Client) GetGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	return d.RestClient.GetGuildScheduledEvent(ctx, guildID, eventID, false)
}

func (d *Client) ListGuildScheduledEvents(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildScheduledEvent, error) {
	events, err := d.RestClient.ListGuildScheduledEvents(ctx, guildID, false)
	if err == nil {
		for _, e := range events {
			e.Hydrate(d)
			if d.Cache != nil {
				d.Cache.ScheduledEvents().Set(e)
			}
		}
	}
	return events, err
}

func (d *Client) GetStageInstance(ctx context.Context, channelID discord.Snowflake) (*discord.StageInstance, error) {
	return d.RestClient.GetStageInstance(ctx, channelID)
}

func (d *Client) CreateStageInstance(ctx context.Context, opts discord.StageCreateOptions) (*discord.StageInstance, error) {
	params := api.CreateStageInstanceParams{
		ChannelID:             opts.ChannelID,
		Topic:                 opts.Topic,
		PrivacyLevel:          opts.PrivacyLevel,
		GuildScheduledEventID: opts.GuildScheduledEventID,
		SendStartNotification: opts.SendStartNotification,
	}
	var apiOpts *api.CreateStageInstanceOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateStageInstanceOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.CreateStageInstance(ctx, params, apiOpts)
}

func (d *Client) GetGuildWebhooks(ctx context.Context, guildID discord.Snowflake) ([]*discord.Webhook, error) {
	webhooks, err := d.RestClient.GetGuildWebhooks(ctx, guildID)
	if err == nil {
		for _, w := range webhooks {
			w.Hydrate(d)
		}
	}
	return webhooks, err
}

func (d *Client) GetAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) (*discord.AutoModerationRule, error) {
	return d.RestClient.GetAutoModerationRule(ctx, guildID, ruleID)
}

func (d *Client) ListAutoModerationRules(ctx context.Context, guildID discord.Snowflake) ([]*discord.AutoModerationRule, error) {
	rules, err := d.RestClient.ListAutoModerationRules(ctx, guildID)
	if err == nil {
		for _, r := range rules {
			r.Hydrate(d)
		}
	}
	return rules, err
}

func (d *Client) CreateAutoModerationRule(ctx context.Context, guildID discord.Snowflake, opts discord.RuleCreateOptions) (*discord.AutoModerationRule, error) {
	params := api.CreateAutoModerationRuleParams{
		Name:            opts.Name,
		EventType:       opts.EventType,
		TriggerType:     opts.TriggerType,
		TriggerMetadata: opts.TriggerMetadata,
		Actions:         opts.Actions,
		Enabled:         opts.Enabled,
		ExemptRoles:     opts.ExemptRoles,
		ExemptChannels:  opts.ExemptChannels,
	}
	var apiOpts *api.CreateAutoModerationRuleOptions
	if opts.AuditLogReason != nil {
		apiOpts = &api.CreateAutoModerationRuleOptions{Reason: *opts.AuditLogReason}
	}
	return d.RestClient.CreateAutoModerationRule(ctx, guildID, params, apiOpts)
}
