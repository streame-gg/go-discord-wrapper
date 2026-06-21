package connection

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// Compile-time assertion that *Client satisfies discord.EntityClient.
var _ discord.EntityClient = (*Client)(nil)

// ── Message ────────────────────────────────────────────────────────────────────

func (d *Client) EditMessage(ctx context.Context, channelID, messageID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error) {
	params := api.EditMessageParams{
		Content:         optPtr(opts.Content),
		Embeds:          optPtr(opts.Embeds),
		Flags:           optPtr(opts.Flags),
		AllowedMentions: optPtr(opts.AllowedMentions),
		Components:      opts.Components,
		Attachments:     optPtr(opts.Attachments),
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
		Content:          optVal(opts.Content),
		TTS:              optVal(opts.TTS),
		Embeds:           optSlice(opts.Embeds),
		AllowedMentions:  optPtr(opts.AllowedMentions),
		MessageReference: optPtr(opts.MessageReference),
		Components:       opts.Components,
		StickerIDs:       optSlice(opts.StickerIDs),
		Flags:            optVal(opts.Flags),
		Files:            opts.Files,
	}
	msg, err := d.RestClient.CreateMessage(ctx, channelID, params)
	if err == nil {
		d.cacheMessage(msg)
	}
	return msg, err
}

func (d *Client) ListChannelMessages(ctx context.Context, channelID discord.Snowflake, opts discord.FetchMessagesOptions) ([]*discord.Message, error) {
	params := api.GetMessagesParams{
		Around: opts.Around,
		Before: opts.Before,
		After:  opts.After,
		Limit:  opts.Limit,
	}
	msgs, err := d.RestClient.ListMessages(ctx, channelID, params)
	if err == nil {
		d.cacheMessages(msgs)
	}
	return msgs, err
}

// ── Channel ────────────────────────────────────────────────────────────────────

func (d *Client) ModifyChannel(ctx context.Context, channelID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error) {
	params := api.ModifyChannelParams{
		Name:                          optPtr(opts.Name),
		Type:                          optPtr(opts.Type),
		Position:                      optPtr(opts.Position),
		Topic:                         optPtr(opts.Topic),
		NSFW:                          optPtr(opts.NSFW),
		RateLimitPerUser:              optPtr(opts.RateLimitPerUser),
		Bitrate:                       optPtr(opts.Bitrate),
		UserLimit:                     optPtr(opts.UserLimit),
		PermissionOverwrites:          optSlice(opts.PermissionOverwrites),
		ParentID:                      optPtr(opts.ParentID),
		RTCRegion:                     optPtr(opts.RTCRegion),
		VideoQualityMode:              optPtr(opts.VideoQualityMode),
		DefaultAutoArchiveDuration:    optPtr(opts.DefaultAutoArchiveDuration),
		Flags:                         optPtr(opts.Flags),
		AvailableTags:                 optSlice(opts.AvailableTags),
		DefaultReactionEmoji:          optPtr(opts.DefaultReactionEmoji),
		DefaultThreadRateLimitPerUser: optPtr(opts.DefaultThreadRateLimitPerUser),
		DefaultSortOrder:              optPtr(opts.DefaultSortOrder),
		DefaultForumLayout:            optPtr(opts.DefaultForumLayout),
	}
	params.AuditLogReason = opts.AuditLogReason
	ch, err := d.RestClient.ModifyChannel(ctx, channelID, params)
	if err == nil {
		d.cacheChannel(ch)
	}
	return ch, err
}

func (d *Client) CreateChannelInvite(ctx context.Context, channelID discord.Snowflake, opts discord.InviteCreateOptions) (*discord.Invite, error) {
	params := api.CreateChannelInviteParams{
		MaxAge:              optPtr(opts.MaxAge),
		MaxUses:             optPtr(opts.MaxUses),
		Temporary:           optPtr(opts.Temporary),
		Unique:              optPtr(opts.Unique),
		TargetType:          optPtr(opts.TargetType),
		TargetUserID:        optPtr(opts.TargetUserID),
		TargetApplicationID: optPtr(opts.TargetApplicationID),
	}
	params.AuditLogReason = opts.AuditLogReason
	invite, err := d.RestClient.CreateChannelInvite(ctx, channelID, params)
	if err == nil {
		invite.Hydrate(d)
		d.cacheInvite(invite)
	}
	return invite, err
}

// ── Guild ──────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuild(ctx context.Context, guildID discord.Snowflake, opts discord.GuildEditOptions) (*discord.Guild, error) {
	params := api.ModifyGuildParams{
		Name:                        optPtr(opts.Name),
		VerificationLevel:           optPtr(opts.VerificationLevel),
		DefaultMessageNotifications: optPtr(opts.DefaultMessageNotifications),
		ExplicitContentFilter:       optPtr(opts.ExplicitContentFilter),
		AFKChannelID:                optPtr(opts.AFKChannelID),
		AFKTimeout:                  optPtr(opts.AFKTimeout),
		Icon:                        optPtr(opts.Icon),
		OwnerID:                     optPtr(opts.OwnerID),
		Splash:                      optPtr(opts.Splash),
		DiscoverySplash:             optPtr(opts.DiscoverySplash),
		Banner:                      optPtr(opts.Banner),
		SystemChannelID:             optPtr(opts.SystemChannelID),
		SystemChannelFlags:          optPtr(opts.SystemChannelFlags),
		RulesChannelID:              optPtr(opts.RulesChannelID),
		PublicUpdatesChannelID:      optPtr(opts.PublicUpdatesChannelID),
		PreferredLocale:             optPtr(opts.PreferredLocale),
		Features:                    optSlice(opts.Features),
		Description:                 optPtr(opts.Description),
		PremiumProgressBarEnabled:   optPtr(opts.PremiumProgressBarEnabled),
		SafetyAlertsChannelID:       optPtr(opts.SafetyAlertsChannelID),
	}
	params.AuditLogReason = opts.AuditLogReason
	guild, err := d.RestClient.ModifyGuild(ctx, guildID, params)
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
		Name:         optPtr(opts.Name),
		Permissions:  optPtr(opts.Permissions),
		Color:        optPtr(opts.Color),
		Hoist:        optPtr(opts.Hoist),
		Icon:         optPtr(opts.Icon),
		UnicodeEmoji: optPtr(opts.UnicodeEmoji),
		Mentionable:  optPtr(opts.Mentionable),
	}
	params.AuditLogReason = opts.AuditLogReason
	role, err := d.RestClient.CreateGuildRole(ctx, guildID, params)
	if err == nil {
		d.cacheRole(guildID, role)
	}
	return role, err
}

func (d *Client) CreateGuildChannel(ctx context.Context, guildID discord.Snowflake, opts discord.ChannelCreateOptions) (*discord.Channel, error) {
	params := api.CreateGuildChannelParams{
		Name:                          opts.Name,
		Type:                          optPtr(opts.Type),
		Topic:                         optPtr(opts.Topic),
		Bitrate:                       optPtr(opts.Bitrate),
		UserLimit:                     optPtr(opts.UserLimit),
		RateLimitPerUser:              optPtr(opts.RateLimitPerUser),
		Position:                      optPtr(opts.Position),
		PermissionOverwrites:          optSlice(opts.PermissionOverwrites),
		ParentID:                      optPtr(opts.ParentID),
		NSFW:                          optPtr(opts.NSFW),
		RTCRegion:                     optPtr(opts.RTCRegion),
		VideoQualityMode:              optPtr(opts.VideoQualityMode),
		DefaultAutoArchiveDuration:    optPtr(opts.DefaultAutoArchiveDuration),
		DefaultReactionEmoji:          optPtr(opts.DefaultReactionEmoji),
		AvailableTags:                 optSlice(opts.AvailableTags),
		DefaultSortOrder:              optPtr(opts.DefaultSortOrder),
		DefaultForumLayout:            optPtr(opts.DefaultForumLayout),
		DefaultThreadRateLimitPerUser: optPtr(opts.DefaultThreadRateLimitPerUser),
	}
	params.AuditLogReason = opts.AuditLogReason
	ch, err := d.RestClient.CreateGuildChannel(ctx, guildID, params)
	if err == nil {
		d.cacheChannel(ch)
	}
	return ch, err
}

func (d *Client) CreateGuildEmoji(ctx context.Context, guildID discord.Snowflake, opts discord.EmojiCreateOptions) (*discord.Emoji, error) {
	params := api.CreateGuildEmojiParams{
		Name:  opts.Name,
		Image: opts.Image,
		Roles: optSlice(opts.Roles),
	}
	params.AuditLogReason = opts.AuditLogReason
	return d.RestClient.CreateGuildEmoji(ctx, guildID, params)
}

func (d *Client) CreateGuildSticker(ctx context.Context, guildID discord.Snowflake, opts discord.StickerCreateOptions) (*discord.Sticker, error) {
	params := api.CreateGuildStickerParams{
		Name:        opts.Name,
		Description: opts.Description,
		Tags:        opts.Tags,
		File:        opts.File,
		ContentType: opts.ContentType,
	}
	params.AuditLogReason = opts.AuditLogReason
	sticker, err := d.RestClient.CreateGuildSticker(ctx, guildID, params)
	if err == nil {
		sticker.GuildID = guildID
		sticker.Hydrate(d)
		if d.Cache != nil {
			d.Cache.Stickers().Set(guildID, sticker)
		}
	}
	return sticker, err
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
		ChannelID:          optPtr(opts.ChannelID),
		EntityMetadata:     optPtr(opts.EntityMetadata),
		Name:               opts.Name,
		PrivacyLevel:       opts.PrivacyLevel,
		ScheduledStartTime: opts.ScheduledStartTime,
		ScheduledEndTime:   optPtr(opts.ScheduledEndTime),
		Description:        optPtr(opts.Description),
		EntityType:         opts.EntityType,
		Image:              optPtr(opts.Image),
	}
	return d.RestClient.CreateGuildScheduledEvent(ctx, guildID, params)
}

// ── GuildMember ────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildMember(ctx context.Context, guildID, userID discord.Snowflake, opts discord.MemberEditOptions) (*discord.GuildMember, error) {
	params := api.ModifyGuildMemberParams{
		Nick:                       optPtr(opts.Nick),
		Roles:                      optPtr(opts.Roles),
		Mute:                       optPtr(opts.Mute),
		Deaf:                       optPtr(opts.Deaf),
		ChannelID:                  optPtr(opts.ChannelID),
		CommunicationDisabledUntil: optPtr(opts.CommunicationDisabledUntil),
		Flags:                      optPtr(opts.Flags),
	}
	params.AuditLogReason = opts.AuditLogReason
	member, err := d.RestClient.ModifyGuildMember(ctx, guildID, userID, params)
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
		DeleteMessageSeconds: optPtr(opts.DeleteMessageSeconds),
	}
	params.AuditLogReason = opts.AuditLogReason
	return d.RestClient.CreateGuildBan(ctx, guildID, userID, params)
}

// ── Role ───────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, opts discord.RoleEditOptions) (*discord.Role, error) {
	params := api.ModifyGuildRoleParams{
		Name:         optPtr(opts.Name),
		Permissions:  optPtr(opts.Permissions),
		Color:        optPtr(opts.Color),
		Hoist:        optPtr(opts.Hoist),
		Icon:         optPtr(opts.Icon),
		UnicodeEmoji: optPtr(opts.UnicodeEmoji),
		Mentionable:  optPtr(opts.Mentionable),
	}
	params.AuditLogReason = opts.AuditLogReason
	role, err := d.RestClient.ModifyGuildRole(ctx, guildID, roleID, params)
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
			Position: optPtr(e.Position),
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
		Name:  optPtr(opts.Name),
		Roles: optSlice(opts.Roles),
	}
	params.AuditLogReason = opts.AuditLogReason
	return d.RestClient.ModifyGuildEmoji(ctx, guildID, emojiID, params)
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
		Name:      optPtr(opts.Name),
		Avatar:    optPtr(opts.Avatar),
		ChannelID: optPtr(opts.ChannelID),
	}
	params.AuditLogReason = opts.AuditLogReason
	return d.RestClient.ModifyWebhook(ctx, webhookID, params)
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
		Content:         optVal(opts.Content),
		Username:        optPtr(opts.Username),
		AvatarURL:       optPtr(opts.AvatarURL),
		TTS:             optVal(opts.TTS),
		Embeds:          optSlice(opts.Embeds),
		AllowedMentions: optPtr(opts.AllowedMentions),
		Components:      opts.Components,
		Flags:           optVal(opts.Flags),
		ThreadName:      optPtr(opts.ThreadName),
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
		Topic:        optPtr(opts.Topic),
		PrivacyLevel: optPtr(opts.PrivacyLevel),
	}
	params.AuditLogReason = opts.AuditLogReason
	return d.RestClient.ModifyStageInstance(ctx, channelID, params)
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
		ChannelID:          optPtr(opts.ChannelID),
		EntityMetadata:     optPtr(opts.EntityMetadata),
		Name:               optPtr(opts.Name),
		PrivacyLevel:       optPtr(opts.PrivacyLevel),
		ScheduledStartTime: optPtr(opts.ScheduledStartTime),
		ScheduledEndTime:   optPtr(opts.ScheduledEndTime),
		Description:        optPtr(opts.Description),
		EntityType:         optPtr(opts.EntityType),
		Status:             optPtr(opts.Status),
		Image:              optPtr(opts.Image),
	}
	return d.RestClient.ModifyGuildScheduledEvent(ctx, guildID, eventID, params)
}

func (d *Client) DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) error {
	return d.RestClient.DeleteGuildScheduledEvent(ctx, guildID, eventID)
}

func (d *Client) ListGuildScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.FetchUsersOptions) ([]*discord.GuildScheduledEventUser, error) {
	params := api.GetGuildScheduledEventUsersParams{
		Limit:      opts.Limit,
		WithMember: opts.WithMember,
		Before:     opts.Before,
		After:      opts.After,
	}
	return d.RestClient.ListGuildScheduledEventUsers(ctx, guildID, eventID, params)
}

// ── Sticker ────────────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, opts discord.StickerEditOptions) (*discord.Sticker, error) {
	params := api.ModifyGuildStickerParams{
		Name:        optPtr(opts.Name),
		Description: optPtr(opts.Description),
		Tags:        optPtr(opts.Tags),
	}
	params.AuditLogReason = opts.AuditLogReason
	sticker, err := d.RestClient.ModifyGuildSticker(ctx, guildID, stickerID, params)
	if err == nil {
		sticker.GuildID = guildID
		sticker.Hydrate(d)
		if d.Cache != nil {
			d.Cache.Stickers().Set(guildID, sticker)
		}
	}
	return sticker, err
}

func (d *Client) DeleteGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, reason *string) error {
	var opts *api.DeleteGuildStickerOptions
	if reason != nil {
		opts = &api.DeleteGuildStickerOptions{Reason: *reason}
	}
	if err := d.RestClient.DeleteGuildSticker(ctx, guildID, stickerID, opts); err != nil {
		return err
	}
	if d.Cache != nil {
		d.Cache.Stickers().Delete(stickerID)
	}
	return nil
}

// ── AutoModerationRule ─────────────────────────────────────────────────────────

func (d *Client) ModifyAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake, opts discord.RuleEditOptions) (*discord.AutoModerationRule, error) {
	params := api.ModifyAutoModerationRuleParams{
		Name:            optPtr(opts.Name),
		EventType:       optPtr(opts.EventType),
		TriggerMetadata: optPtr(opts.TriggerMetadata),
		Actions:         optSlice(opts.Actions),
		Enabled:         optPtr(opts.Enabled),
		ExemptRoles:     optSlice(opts.ExemptRoles),
		ExemptChannels:  optSlice(opts.ExemptChannels),
	}
	params.AuditLogReason = opts.AuditLogReason
	rule, err := d.RestClient.ModifyAutoModerationRule(ctx, guildID, ruleID, params)
	if err == nil {
		d.cacheAutoModRule(guildID, rule)
	}
	return rule, err
}

func (d *Client) DeleteAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) error {
	if err := d.RestClient.DeleteAutoModerationRule(ctx, guildID, ruleID, nil); err != nil {
		return err
	}
	if d.cacheStoreEnabled(cache.CategoryAutoModRules) {
		d.Cache.AutoModRules().Delete(ruleID)
	}
	return nil
}

// ── SoundboardSound ────────────────────────────────────────────────────────────

func (d *Client) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, opts discord.SoundEditOptions) (*discord.SoundboardSound, error) {
	params := api.ModifyGuildSoundboardSoundParams{
		Name:      optPtr(opts.Name),
		Volume:    optPtr(opts.Volume),
		EmojiID:   optPtr(opts.EmojiID),
		EmojiName: optPtr(opts.EmojiName),
	}
	params.AuditLogReason = opts.AuditLogReason
	sound, err := d.RestClient.ModifyGuildSoundboardSound(ctx, guildID, soundID, params)
	if err == nil {
		sound.GuildID = &guildID
		sound.Hydrate(d)
		if d.Cache != nil {
			d.Cache.Soundboard().Set(guildID, sound)
		}
	}
	return sound, err
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
	emoji, err := d.RestClient.GetGuildEmoji(ctx, guildID, emojiID)
	if err == nil {
		emoji.GuildID = guildID
		emoji.Hydrate(d)

		if d.Cache != nil {
			d.Cache.Emojis().Set(guildID, emoji)
		}
	}

	return emoji, err
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
	sticker, err := d.RestClient.GetGuildSticker(ctx, guildID, stickerID)
	if err == nil {
		sticker.Hydrate(d)
		sticker.GuildID = guildID

		if d.Cache != nil {
			d.Cache.Stickers().Set(guildID, sticker)
		}
	}

	return sticker, err
}

func (d *Client) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) ([]*discord.Sticker, error) {
	stickers, err := d.RestClient.ListGuildStickers(ctx, guildID)
	if err == nil {
		for _, s := range stickers {
			s.GuildID = guildID
			s.Hydrate(d)
		}
		if d.Cache != nil {
			d.Cache.Stickers().SetAll(guildID, stickers)
		}
	}
	return stickers, err
}

func (d *Client) GetGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	event, err := d.RestClient.GetGuildScheduledEvent(ctx, guildID, eventID, false)

	if err == nil {
		event.Hydrate(d)
		if d.Cache.ScheduledEvents() != nil {
			d.Cache.ScheduledEvents().Set(event)
		}
	}

	return event, err
}

func (d *Client) ListGuildScheduledEvents(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildScheduledEvent, error) {
	events, err := d.RestClient.ListGuildScheduledEvents(ctx, guildID, false)
	if err == nil {
		for _, e := range events {
			e.Hydrate(d)
			if d.Cache.ScheduledEvents() != nil {
				d.Cache.ScheduledEvents().Set(e)
			}
		}
	}
	return events, err
}

func (d *Client) GetStageInstance(ctx context.Context, channelID discord.Snowflake) (*discord.StageInstance, error) {
	stageInstance, err := d.RestClient.GetStageInstance(ctx, channelID)
	if err == nil {
		stageInstance.Hydrate(d)
		if d.Cache != nil {
			d.Cache.StageInstances().Set(stageInstance)
		}
	}

	return stageInstance, err
}

func (d *Client) CreateStageInstance(ctx context.Context, opts discord.StageCreateOptions) (*discord.StageInstance, error) {
	params := api.CreateStageInstanceParams{
		ChannelID:             opts.ChannelID,
		Topic:                 opts.Topic,
		PrivacyLevel:          optPtr(opts.PrivacyLevel),
		GuildScheduledEventID: optPtr(opts.GuildScheduledEventID),
		SendStartNotification: optPtr(opts.SendStartNotification),
	}
	params.AuditLogReason = opts.AuditLogReason
	return d.RestClient.CreateStageInstance(ctx, params)
}

func (d *Client) ListGuildWebhooks(ctx context.Context, guildID discord.Snowflake) ([]*discord.Webhook, error) {
	webhooks, err := d.RestClient.ListGuildWebhooks(ctx, guildID)
	if err == nil {
		for _, w := range webhooks {
			w.Hydrate(d)
		}
	}
	return webhooks, err
}

func (d *Client) GetAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) (*discord.AutoModerationRule, error) {
	rule, err := d.RestClient.GetAutoModerationRule(ctx, guildID, ruleID)
	if err == nil {
		d.cacheAutoModRule(guildID, rule)
	}
	return rule, err
}

func (d *Client) ListAutoModerationRules(ctx context.Context, guildID discord.Snowflake) ([]*discord.AutoModerationRule, error) {
	rules, err := d.RestClient.ListAutoModerationRules(ctx, guildID)
	if err == nil {
		for _, r := range rules {
			d.cacheAutoModRule(guildID, r)
		}
	}
	return rules, err
}

func (d *Client) CreateAutoModerationRule(ctx context.Context, guildID discord.Snowflake, opts discord.RuleCreateOptions) (*discord.AutoModerationRule, error) {
	params := api.CreateAutoModerationRuleParams{
		Name:            opts.Name,
		EventType:       opts.EventType,
		TriggerType:     opts.TriggerType,
		TriggerMetadata: optPtr(opts.TriggerMetadata),
		Actions:         opts.Actions,
		Enabled:         optPtr(opts.Enabled),
		ExemptRoles:     optSlice(opts.ExemptRoles),
		ExemptChannels:  optSlice(opts.ExemptChannels),
	}
	params.AuditLogReason = opts.AuditLogReason
	rule, err := d.RestClient.CreateAutoModerationRule(ctx, guildID, params)
	if err == nil {
		d.cacheAutoModRule(guildID, rule)
	}
	return rule, err
}
