package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestChannelCreate() {
	s.T().Log("Testing Channel Create Unmarshal Logic")

	sub := testutil.InitSub[ChannelCreateEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelCreateEvent) {
				s.EqualValues(channel["id"], got.Channel.ID)
				s.EqualValues(channel["type"], got.Channel.Type)
				s.EqualValues(channel["guild_id"], *got.Channel.GuildID)
				s.EqualValues(channel["position"], got.Channel.Position)
				s.EqualValues(channel["name"], *got.Channel.Name)
				s.EqualValues(channel["topic"], *got.Channel.Topic)
				s.EqualValues(channel["nsfw"], *got.Channel.NSFW)
				s.EqualValues(channel["last_message_id"], *got.Channel.LastMessageID)
				s.EqualValues(channel["bitrate"], got.Channel.Bitrate)
				s.EqualValues(channel["user_limit"], got.Channel.UserLimit)
				s.EqualValues(channel["rate_limit_per_user"], got.Channel.RateLimitPerUser)
				s.EqualValues(channel["icon"], *got.Channel.Icon)
				s.EqualValues(channel["owner_id"], *got.Channel.OwnerID)
				s.EqualValues(channel["application_id"], *got.Channel.ApplicationID)
				s.EqualValues(channel["managed"], got.Channel.Managed)
				s.EqualValues(channel["parent_id"], *got.Channel.ParentID)
				s.EqualValues(channel["last_pin_timestamp"], *got.Channel.LastPinTimestamp)
				s.EqualValues(channel["rtc_region"], *got.Channel.RtcRegion)
				s.EqualValues(channel["video_quality_mode"], got.Channel.VideoQualityMode)
				s.EqualValues(channel["message_count"], got.Channel.MessageCount)
				s.EqualValues(channel["member_count"], got.Channel.MemberCount)
				s.EqualValues(channel["default_auto_archive_duration"], got.Channel.DefaultAutoArchiveDuration)
				s.EqualValues(channel["permissions"], got.Channel.Permissions)
				s.EqualValues(channel["flags"], got.Channel.Flags)
				s.EqualValues(channel["total_message_sent"], got.Channel.TotalMessageSent)
				s.EqualValues(channel["applied_tags"], got.Channel.AppliedTags)
				s.EqualValues(channel["default_thread_rate_limit_per_user"], got.Channel.DefaultThreadRateLimitPerUser)
				s.EqualValues(channel["default_sort_order"], *got.Channel.DefaultSortOrder)
				s.EqualValues(channel["default_forum_layout"], got.Channel.DefaultForumLayout)

				threadMetadata := channel["thread_metadata"].(map[string]interface{})
				s.EqualValues(threadMetadata["archived"], got.Channel.ThreadMetadata.Archived)
				s.EqualValues(threadMetadata["auto_archive_duration"], got.Channel.ThreadMetadata.AutoArchiveDuration)
				s.EqualValues(threadMetadata["archive_timestamp"], got.Channel.ThreadMetadata.ArchiveTimestamp)
				s.EqualValues(threadMetadata["created_timestamp"], *got.Channel.ThreadMetadata.CreatedTimestamp)
				s.EqualValues(threadMetadata["locked"], got.Channel.ThreadMetadata.Locked)
				s.EqualValues(threadMetadata["invitable"], got.Channel.ThreadMetadata.Invitable)

				defaultReactionEmoji := channel["default_reaction_emoji"].(map[string]interface{})
				s.EqualValues(defaultReactionEmoji["emoji_id"], *got.Channel.DefaultReactionEmoji.EmojiID)
				s.EqualValues(defaultReactionEmoji["emoji_name"], *got.Channel.DefaultReactionEmoji.EmojiName)

				permissionOverwrites := channel["permission_overwrites"].([]map[string]interface{})
				for i, po := range permissionOverwrites {
					s.EqualValues(po["id"], got.Channel.PermissionOverwrites[i].ID)
					s.EqualValues(po["type"], got.Channel.PermissionOverwrites[i].Type)
					s.EqualValues(po["allow"], got.Channel.PermissionOverwrites[i].Allow)
					s.EqualValues(po["deny"], got.Channel.PermissionOverwrites[i].Deny)
				}

				availableTags := channel["available_tags"].([]map[string]interface{})
				for i, tag := range availableTags {
					s.EqualValues(tag["id"], got.Channel.AvailableTags[i].ID)
					s.EqualValues(tag["name"], got.Channel.AvailableTags[i].Name)
					s.EqualValues(tag["moderated"], got.Channel.AvailableTags[i].Moderated)
					s.EqualValues(tag["emoji_id"], *got.Channel.AvailableTags[i].EmojiID)
					s.EqualValues(tag["emoji_name"], *got.Channel.AvailableTags[i].EmojiName)
				}

				member := channel["member"].(map[string]interface{})
				s.EqualValues(member["id"], *got.Channel.Member.ID)
				s.EqualValues(member["user_id"], *got.Channel.Member.UserID)
				s.EqualValues(member["join_timestamp"], got.Channel.Member.JoinTimestamp)
				s.EqualValues(member["flags"], got.Channel.Member.Flags)

				guildMember := member["member"].(map[string]interface{})
				s.EqualValues(guildMember["avatar"], *got.Channel.Member.Member.Avatar)
				s.EqualValues(guildMember["banner"], *got.Channel.Member.Member.Banner)
				s.EqualValues(guildMember["communication_disabled_until"], *got.Channel.Member.Member.CommunicationDisabledUntil)
				s.EqualValues(guildMember["deaf"], got.Channel.Member.Member.Deaf)
				s.EqualValues(guildMember["flags"], got.Channel.Member.Member.Flags)
				s.EqualValues(guildMember["joined_at"], *got.Channel.Member.Member.JoinedAt)
				s.EqualValues(guildMember["mute"], got.Channel.Member.Member.Mute)
				s.EqualValues(guildMember["nick"], *got.Channel.Member.Member.Nick)
				s.EqualValues(guildMember["pending"], got.Channel.Member.Member.Pending)
				s.EqualValues(guildMember["premium_since"], *got.Channel.Member.Member.PremiumSince)
				s.EqualValues(guildMember["roles"], got.Channel.Member.Member.Roles)
				s.EqualValues(guildMember["permissions"], *got.Channel.Member.Member.Permissions)

				guildMemberUser := guildMember["user"].(map[string]interface{})
				s.EqualValues(guildMemberUser["id"], got.Channel.Member.Member.User.ID)
				s.EqualValues(guildMemberUser["username"], got.Channel.Member.Member.User.Username)
				s.EqualValues(guildMemberUser["discriminator"], got.Channel.Member.Member.User.Discriminator)
				s.EqualValues(guildMemberUser["avatar"], *got.Channel.Member.Member.User.Avatar)
				s.EqualValues(guildMemberUser["bot"], *got.Channel.Member.Member.User.Bot)
				s.EqualValues(guildMemberUser["system"], *got.Channel.Member.Member.User.System)
				s.EqualValues(guildMemberUser["mfa_enabled"], *got.Channel.Member.Member.User.MFAEnabled)
				s.EqualValues(guildMemberUser["banner"], *got.Channel.Member.Member.User.Banner)
				s.EqualValues(guildMemberUser["accent_color"], *got.Channel.Member.Member.User.AccentColor)
				s.EqualValues(guildMemberUser["locale"], *got.Channel.Member.Member.User.Locale)
				s.EqualValues(guildMemberUser["verified"], *got.Channel.Member.Member.User.Verified)
				s.EqualValues(guildMemberUser["email"], *got.Channel.Member.Member.User.Email)
				s.EqualValues(guildMemberUser["flags"], got.Channel.Member.Member.User.Flags)
				s.EqualValues(guildMemberUser["premium_type"], *got.Channel.Member.Member.User.PremiumType)
				s.EqualValues(guildMemberUser["public_flags"], got.Channel.Member.Member.User.PublicFlags)
				s.EqualValues(guildMemberUser["global_name"], *got.Channel.Member.Member.User.GlobalName)

				avatarDecoration := guildMemberUser["avatar_decoration"].(map[string]interface{})
				s.EqualValues(avatarDecoration["asset"], got.Channel.Member.Member.User.AvatarDecorationData.Asset)
				s.EqualValues(avatarDecoration["sku_id"], got.Channel.Member.Member.User.AvatarDecorationData.SkuID)

				recipients := channel["recipients"].([]map[string]interface{})
				for i, r := range recipients {
					s.EqualValues(r["id"], got.Channel.Recipients[i].ID)
					s.EqualValues(r["username"], got.Channel.Recipients[i].Username)
					s.EqualValues(r["discriminator"], got.Channel.Recipients[i].Discriminator)
					s.EqualValues(r["avatar"], *got.Channel.Recipients[i].Avatar)
					s.EqualValues(r["bot"], *got.Channel.Recipients[i].Bot)
					s.EqualValues(r["system"], *got.Channel.Recipients[i].System)
					s.EqualValues(r["mfa_enabled"], *got.Channel.Recipients[i].MFAEnabled)
					s.EqualValues(r["banner"], *got.Channel.Recipients[i].Banner)
					s.EqualValues(r["accent_color"], *got.Channel.Recipients[i].AccentColor)
					s.EqualValues(r["locale"], *got.Channel.Recipients[i].Locale)
					s.EqualValues(r["verified"], *got.Channel.Recipients[i].Verified)
					s.EqualValues(r["email"], *got.Channel.Recipients[i].Email)
					s.EqualValues(r["flags"], got.Channel.Recipients[i].Flags)
					s.EqualValues(r["premium_type"], *got.Channel.Recipients[i].PremiumType)
					s.EqualValues(r["public_flags"], got.Channel.Recipients[i].PublicFlags)
					s.EqualValues(r["global_name"], *got.Channel.Recipients[i].GlobalName)

					recipientAvatarDecoration := r["avatar_decoration"].(map[string]interface{})
					s.EqualValues(recipientAvatarDecoration["asset"], got.Channel.Recipients[i].AvatarDecorationData.Asset)
					s.EqualValues(recipientAvatarDecoration["sku_id"], got.Channel.Recipients[i].AvatarDecorationData.SkuID)
				}
			},
		},
	})
}

func (s *eventSuite) TestChannelUpdate() {
	s.T().Log("Testing Channel Update Unmarshal Logic")

	sub := testutil.InitSub[ChannelUpdateEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelUpdateEvent) {
				s.EqualValues(channel["id"], got.NewChannel.ID)
				s.EqualValues(channel["type"], got.NewChannel.Type)
				s.EqualValues(channel["guild_id"], *got.NewChannel.GuildID)
				s.EqualValues(channel["position"], got.NewChannel.Position)
				s.EqualValues(channel["name"], *got.NewChannel.Name)
				s.EqualValues(channel["topic"], *got.NewChannel.Topic)
				s.EqualValues(channel["nsfw"], *got.NewChannel.NSFW)
				s.EqualValues(channel["last_message_id"], *got.NewChannel.LastMessageID)
				s.EqualValues(channel["bitrate"], got.NewChannel.Bitrate)
				s.EqualValues(channel["user_limit"], got.NewChannel.UserLimit)
				s.EqualValues(channel["rate_limit_per_user"], got.NewChannel.RateLimitPerUser)
				s.EqualValues(channel["icon"], *got.NewChannel.Icon)
				s.EqualValues(channel["owner_id"], *got.NewChannel.OwnerID)
				s.EqualValues(channel["application_id"], *got.NewChannel.ApplicationID)
				s.EqualValues(channel["managed"], got.NewChannel.Managed)
				s.EqualValues(channel["parent_id"], *got.NewChannel.ParentID)
				s.EqualValues(channel["last_pin_timestamp"], *got.NewChannel.LastPinTimestamp)
				s.EqualValues(channel["rtc_region"], *got.NewChannel.RtcRegion)
				s.EqualValues(channel["video_quality_mode"], got.NewChannel.VideoQualityMode)
				s.EqualValues(channel["message_count"], got.NewChannel.MessageCount)
				s.EqualValues(channel["member_count"], got.NewChannel.MemberCount)
				s.EqualValues(channel["default_auto_archive_duration"], got.NewChannel.DefaultAutoArchiveDuration)
				s.EqualValues(channel["permissions"], got.NewChannel.Permissions)
				s.EqualValues(channel["flags"], got.NewChannel.Flags)
				s.EqualValues(channel["total_message_sent"], got.NewChannel.TotalMessageSent)
				s.EqualValues(channel["applied_tags"], got.NewChannel.AppliedTags)
				s.EqualValues(channel["default_thread_rate_limit_per_user"], got.NewChannel.DefaultThreadRateLimitPerUser)
				s.EqualValues(channel["default_sort_order"], *got.NewChannel.DefaultSortOrder)
				s.EqualValues(channel["default_forum_layout"], got.NewChannel.DefaultForumLayout)

				threadMetadata := channel["thread_metadata"].(map[string]interface{})
				s.EqualValues(threadMetadata["archived"], got.NewChannel.ThreadMetadata.Archived)
				s.EqualValues(threadMetadata["auto_archive_duration"], got.NewChannel.ThreadMetadata.AutoArchiveDuration)
				s.EqualValues(threadMetadata["archive_timestamp"], got.NewChannel.ThreadMetadata.ArchiveTimestamp)
				s.EqualValues(threadMetadata["created_timestamp"], *got.NewChannel.ThreadMetadata.CreatedTimestamp)
				s.EqualValues(threadMetadata["locked"], got.NewChannel.ThreadMetadata.Locked)
				s.EqualValues(threadMetadata["invitable"], got.NewChannel.ThreadMetadata.Invitable)

				defaultReactionEmoji := channel["default_reaction_emoji"].(map[string]interface{})
				s.EqualValues(defaultReactionEmoji["emoji_id"], *got.NewChannel.DefaultReactionEmoji.EmojiID)
				s.EqualValues(defaultReactionEmoji["emoji_name"], *got.NewChannel.DefaultReactionEmoji.EmojiName)

				permissionOverwrites := channel["permission_overwrites"].([]map[string]interface{})
				for i, po := range permissionOverwrites {
					s.EqualValues(po["id"], got.NewChannel.PermissionOverwrites[i].ID)
					s.EqualValues(po["type"], got.NewChannel.PermissionOverwrites[i].Type)
					s.EqualValues(po["allow"], got.NewChannel.PermissionOverwrites[i].Allow)
					s.EqualValues(po["deny"], got.NewChannel.PermissionOverwrites[i].Deny)
				}

				availableTags := channel["available_tags"].([]map[string]interface{})
				for i, tag := range availableTags {
					s.EqualValues(tag["id"], got.NewChannel.AvailableTags[i].ID)
					s.EqualValues(tag["name"], got.NewChannel.AvailableTags[i].Name)
					s.EqualValues(tag["moderated"], got.NewChannel.AvailableTags[i].Moderated)
					s.EqualValues(tag["emoji_id"], *got.NewChannel.AvailableTags[i].EmojiID)
					s.EqualValues(tag["emoji_name"], *got.NewChannel.AvailableTags[i].EmojiName)
				}

				member := channel["member"].(map[string]interface{})
				s.EqualValues(member["id"], *got.NewChannel.Member.ID)
				s.EqualValues(member["user_id"], *got.NewChannel.Member.UserID)
				s.EqualValues(member["join_timestamp"], got.NewChannel.Member.JoinTimestamp)
				s.EqualValues(member["flags"], got.NewChannel.Member.Flags)

				guildMember := member["member"].(map[string]interface{})
				s.EqualValues(guildMember["avatar"], *got.NewChannel.Member.Member.Avatar)
				s.EqualValues(guildMember["banner"], *got.NewChannel.Member.Member.Banner)
				s.EqualValues(guildMember["communication_disabled_until"], *got.NewChannel.Member.Member.CommunicationDisabledUntil)
				s.EqualValues(guildMember["deaf"], got.NewChannel.Member.Member.Deaf)
				s.EqualValues(guildMember["flags"], got.NewChannel.Member.Member.Flags)
				s.EqualValues(guildMember["joined_at"], *got.NewChannel.Member.Member.JoinedAt)
				s.EqualValues(guildMember["mute"], got.NewChannel.Member.Member.Mute)
				s.EqualValues(guildMember["nick"], *got.NewChannel.Member.Member.Nick)
				s.EqualValues(guildMember["pending"], got.NewChannel.Member.Member.Pending)
				s.EqualValues(guildMember["premium_since"], *got.NewChannel.Member.Member.PremiumSince)
				s.EqualValues(guildMember["roles"], got.NewChannel.Member.Member.Roles)
				s.EqualValues(guildMember["permissions"], *got.NewChannel.Member.Member.Permissions)

				guildMemberUser := guildMember["user"].(map[string]interface{})
				s.EqualValues(guildMemberUser["id"], got.NewChannel.Member.Member.User.ID)
				s.EqualValues(guildMemberUser["username"], got.NewChannel.Member.Member.User.Username)
				s.EqualValues(guildMemberUser["discriminator"], got.NewChannel.Member.Member.User.Discriminator)
				s.EqualValues(guildMemberUser["avatar"], *got.NewChannel.Member.Member.User.Avatar)
				s.EqualValues(guildMemberUser["bot"], *got.NewChannel.Member.Member.User.Bot)
				s.EqualValues(guildMemberUser["system"], *got.NewChannel.Member.Member.User.System)
				s.EqualValues(guildMemberUser["mfa_enabled"], *got.NewChannel.Member.Member.User.MFAEnabled)
				s.EqualValues(guildMemberUser["banner"], *got.NewChannel.Member.Member.User.Banner)
				s.EqualValues(guildMemberUser["accent_color"], *got.NewChannel.Member.Member.User.AccentColor)
				s.EqualValues(guildMemberUser["locale"], *got.NewChannel.Member.Member.User.Locale)
				s.EqualValues(guildMemberUser["verified"], *got.NewChannel.Member.Member.User.Verified)
				s.EqualValues(guildMemberUser["email"], *got.NewChannel.Member.Member.User.Email)
				s.EqualValues(guildMemberUser["flags"], got.NewChannel.Member.Member.User.Flags)
				s.EqualValues(guildMemberUser["premium_type"], *got.NewChannel.Member.Member.User.PremiumType)
				s.EqualValues(guildMemberUser["public_flags"], got.NewChannel.Member.Member.User.PublicFlags)
				s.EqualValues(guildMemberUser["global_name"], *got.NewChannel.Member.Member.User.GlobalName)

				avatarDecoration := guildMemberUser["avatar_decoration"].(map[string]interface{})
				s.EqualValues(avatarDecoration["asset"], got.NewChannel.Member.Member.User.AvatarDecorationData.Asset)
				s.EqualValues(avatarDecoration["sku_id"], got.NewChannel.Member.Member.User.AvatarDecorationData.SkuID)

				recipients := channel["recipients"].([]map[string]interface{})
				for i, r := range recipients {
					s.EqualValues(r["id"], got.NewChannel.Recipients[i].ID)
					s.EqualValues(r["username"], got.NewChannel.Recipients[i].Username)
					s.EqualValues(r["discriminator"], got.NewChannel.Recipients[i].Discriminator)
					s.EqualValues(r["avatar"], *got.NewChannel.Recipients[i].Avatar)
					s.EqualValues(r["bot"], *got.NewChannel.Recipients[i].Bot)
					s.EqualValues(r["system"], *got.NewChannel.Recipients[i].System)
					s.EqualValues(r["mfa_enabled"], *got.NewChannel.Recipients[i].MFAEnabled)
					s.EqualValues(r["banner"], *got.NewChannel.Recipients[i].Banner)
					s.EqualValues(r["accent_color"], *got.NewChannel.Recipients[i].AccentColor)
					s.EqualValues(r["locale"], *got.NewChannel.Recipients[i].Locale)
					s.EqualValues(r["verified"], *got.NewChannel.Recipients[i].Verified)
					s.EqualValues(r["email"], *got.NewChannel.Recipients[i].Email)
					s.EqualValues(r["flags"], got.NewChannel.Recipients[i].Flags)
					s.EqualValues(r["premium_type"], *got.NewChannel.Recipients[i].PremiumType)
					s.EqualValues(r["public_flags"], got.NewChannel.Recipients[i].PublicFlags)
					s.EqualValues(r["global_name"], *got.NewChannel.Recipients[i].GlobalName)

					recipientAvatarDecoration := r["avatar_decoration"].(map[string]interface{})
					s.EqualValues(recipientAvatarDecoration["asset"], got.NewChannel.Recipients[i].AvatarDecorationData.Asset)
					s.EqualValues(recipientAvatarDecoration["sku_id"], got.NewChannel.Recipients[i].AvatarDecorationData.SkuID)
				}

				s.Nil(got.OldChannel)
			},
		},
	})
}

func (s *eventSuite) TestChannelDelete() {
	s.T().Log("Testing Channel Delete Unmarshal Logic")

	sub := testutil.InitSub[ChannelDeleteEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelDeleteEvent) {
				s.EqualValues(channel["id"], got.Channel.ID)
				s.EqualValues(channel["type"], got.Channel.Type)
				s.EqualValues(channel["guild_id"], *got.Channel.GuildID)
				s.EqualValues(channel["position"], got.Channel.Position)
				s.EqualValues(channel["name"], *got.Channel.Name)
				s.EqualValues(channel["topic"], *got.Channel.Topic)
				s.EqualValues(channel["nsfw"], *got.Channel.NSFW)
				s.EqualValues(channel["last_message_id"], *got.Channel.LastMessageID)
				s.EqualValues(channel["bitrate"], got.Channel.Bitrate)
				s.EqualValues(channel["user_limit"], got.Channel.UserLimit)
				s.EqualValues(channel["rate_limit_per_user"], got.Channel.RateLimitPerUser)
				s.EqualValues(channel["icon"], *got.Channel.Icon)
				s.EqualValues(channel["owner_id"], *got.Channel.OwnerID)
				s.EqualValues(channel["application_id"], *got.Channel.ApplicationID)
				s.EqualValues(channel["managed"], got.Channel.Managed)
				s.EqualValues(channel["parent_id"], *got.Channel.ParentID)
				s.EqualValues(channel["last_pin_timestamp"], *got.Channel.LastPinTimestamp)
				s.EqualValues(channel["rtc_region"], *got.Channel.RtcRegion)
				s.EqualValues(channel["video_quality_mode"], got.Channel.VideoQualityMode)
				s.EqualValues(channel["message_count"], got.Channel.MessageCount)
				s.EqualValues(channel["member_count"], got.Channel.MemberCount)
				s.EqualValues(channel["default_auto_archive_duration"], got.Channel.DefaultAutoArchiveDuration)
				s.EqualValues(channel["permissions"], got.Channel.Permissions)
				s.EqualValues(channel["flags"], got.Channel.Flags)
				s.EqualValues(channel["total_message_sent"], got.Channel.TotalMessageSent)
				s.EqualValues(channel["applied_tags"], got.Channel.AppliedTags)
				s.EqualValues(channel["default_thread_rate_limit_per_user"], got.Channel.DefaultThreadRateLimitPerUser)
				s.EqualValues(channel["default_sort_order"], *got.Channel.DefaultSortOrder)
				s.EqualValues(channel["default_forum_layout"], got.Channel.DefaultForumLayout)

				threadMetadata := channel["thread_metadata"].(map[string]interface{})
				s.EqualValues(threadMetadata["archived"], got.Channel.ThreadMetadata.Archived)
				s.EqualValues(threadMetadata["auto_archive_duration"], got.Channel.ThreadMetadata.AutoArchiveDuration)
				s.EqualValues(threadMetadata["archive_timestamp"], got.Channel.ThreadMetadata.ArchiveTimestamp)
				s.EqualValues(threadMetadata["created_timestamp"], *got.Channel.ThreadMetadata.CreatedTimestamp)
				s.EqualValues(threadMetadata["locked"], got.Channel.ThreadMetadata.Locked)
				s.EqualValues(threadMetadata["invitable"], got.Channel.ThreadMetadata.Invitable)

				defaultReactionEmoji := channel["default_reaction_emoji"].(map[string]interface{})
				s.EqualValues(defaultReactionEmoji["emoji_id"], *got.Channel.DefaultReactionEmoji.EmojiID)
				s.EqualValues(defaultReactionEmoji["emoji_name"], *got.Channel.DefaultReactionEmoji.EmojiName)

				permissionOverwrites := channel["permission_overwrites"].([]map[string]interface{})
				for i, po := range permissionOverwrites {
					s.EqualValues(po["id"], got.Channel.PermissionOverwrites[i].ID)
					s.EqualValues(po["type"], got.Channel.PermissionOverwrites[i].Type)
					s.EqualValues(po["allow"], got.Channel.PermissionOverwrites[i].Allow)
					s.EqualValues(po["deny"], got.Channel.PermissionOverwrites[i].Deny)
				}

				availableTags := channel["available_tags"].([]map[string]interface{})
				for i, tag := range availableTags {
					s.EqualValues(tag["id"], got.Channel.AvailableTags[i].ID)
					s.EqualValues(tag["name"], got.Channel.AvailableTags[i].Name)
					s.EqualValues(tag["moderated"], got.Channel.AvailableTags[i].Moderated)
					s.EqualValues(tag["emoji_id"], *got.Channel.AvailableTags[i].EmojiID)
					s.EqualValues(tag["emoji_name"], *got.Channel.AvailableTags[i].EmojiName)
				}

				member := channel["member"].(map[string]interface{})
				s.EqualValues(member["id"], *got.Channel.Member.ID)
				s.EqualValues(member["user_id"], *got.Channel.Member.UserID)
				s.EqualValues(member["join_timestamp"], got.Channel.Member.JoinTimestamp)
				s.EqualValues(member["flags"], got.Channel.Member.Flags)

				guildMember := member["member"].(map[string]interface{})
				s.EqualValues(guildMember["avatar"], *got.Channel.Member.Member.Avatar)
				s.EqualValues(guildMember["banner"], *got.Channel.Member.Member.Banner)
				s.EqualValues(guildMember["communication_disabled_until"], *got.Channel.Member.Member.CommunicationDisabledUntil)
				s.EqualValues(guildMember["deaf"], got.Channel.Member.Member.Deaf)
				s.EqualValues(guildMember["flags"], got.Channel.Member.Member.Flags)
				s.EqualValues(guildMember["joined_at"], *got.Channel.Member.Member.JoinedAt)
				s.EqualValues(guildMember["mute"], got.Channel.Member.Member.Mute)
				s.EqualValues(guildMember["nick"], *got.Channel.Member.Member.Nick)
				s.EqualValues(guildMember["pending"], got.Channel.Member.Member.Pending)
				s.EqualValues(guildMember["premium_since"], *got.Channel.Member.Member.PremiumSince)
				s.EqualValues(guildMember["roles"], got.Channel.Member.Member.Roles)
				s.EqualValues(guildMember["permissions"], *got.Channel.Member.Member.Permissions)

				guildMemberUser := guildMember["user"].(map[string]interface{})
				s.EqualValues(guildMemberUser["id"], got.Channel.Member.Member.User.ID)
				s.EqualValues(guildMemberUser["username"], got.Channel.Member.Member.User.Username)
				s.EqualValues(guildMemberUser["discriminator"], got.Channel.Member.Member.User.Discriminator)
				s.EqualValues(guildMemberUser["avatar"], *got.Channel.Member.Member.User.Avatar)
				s.EqualValues(guildMemberUser["bot"], *got.Channel.Member.Member.User.Bot)
				s.EqualValues(guildMemberUser["system"], *got.Channel.Member.Member.User.System)
				s.EqualValues(guildMemberUser["mfa_enabled"], *got.Channel.Member.Member.User.MFAEnabled)
				s.EqualValues(guildMemberUser["banner"], *got.Channel.Member.Member.User.Banner)
				s.EqualValues(guildMemberUser["accent_color"], *got.Channel.Member.Member.User.AccentColor)
				s.EqualValues(guildMemberUser["locale"], *got.Channel.Member.Member.User.Locale)
				s.EqualValues(guildMemberUser["verified"], *got.Channel.Member.Member.User.Verified)
				s.EqualValues(guildMemberUser["email"], *got.Channel.Member.Member.User.Email)
				s.EqualValues(guildMemberUser["flags"], got.Channel.Member.Member.User.Flags)
				s.EqualValues(guildMemberUser["premium_type"], *got.Channel.Member.Member.User.PremiumType)
				s.EqualValues(guildMemberUser["public_flags"], got.Channel.Member.Member.User.PublicFlags)
				s.EqualValues(guildMemberUser["global_name"], *got.Channel.Member.Member.User.GlobalName)

				avatarDecoration := guildMemberUser["avatar_decoration"].(map[string]interface{})
				s.EqualValues(avatarDecoration["asset"], got.Channel.Member.Member.User.AvatarDecorationData.Asset)
				s.EqualValues(avatarDecoration["sku_id"], got.Channel.Member.Member.User.AvatarDecorationData.SkuID)

				recipients := channel["recipients"].([]map[string]interface{})
				for i, r := range recipients {
					s.EqualValues(r["id"], got.Channel.Recipients[i].ID)
					s.EqualValues(r["username"], got.Channel.Recipients[i].Username)
					s.EqualValues(r["discriminator"], got.Channel.Recipients[i].Discriminator)
					s.EqualValues(r["avatar"], *got.Channel.Recipients[i].Avatar)
					s.EqualValues(r["bot"], *got.Channel.Recipients[i].Bot)
					s.EqualValues(r["system"], *got.Channel.Recipients[i].System)
					s.EqualValues(r["mfa_enabled"], *got.Channel.Recipients[i].MFAEnabled)
					s.EqualValues(r["banner"], *got.Channel.Recipients[i].Banner)
					s.EqualValues(r["accent_color"], *got.Channel.Recipients[i].AccentColor)
					s.EqualValues(r["locale"], *got.Channel.Recipients[i].Locale)
					s.EqualValues(r["verified"], *got.Channel.Recipients[i].Verified)
					s.EqualValues(r["email"], *got.Channel.Recipients[i].Email)
					s.EqualValues(r["flags"], got.Channel.Recipients[i].Flags)
					s.EqualValues(r["premium_type"], *got.Channel.Recipients[i].PremiumType)
					s.EqualValues(r["public_flags"], got.Channel.Recipients[i].PublicFlags)
					s.EqualValues(r["global_name"], *got.Channel.Recipients[i].GlobalName)

					recipientAvatarDecoration := r["avatar_decoration"].(map[string]interface{})
					s.EqualValues(recipientAvatarDecoration["asset"], got.Channel.Recipients[i].AvatarDecorationData.Asset)
					s.EqualValues(recipientAvatarDecoration["sku_id"], got.Channel.Recipients[i].AvatarDecorationData.SkuID)
				}
			},
		},
	})
}

func (s *eventSuite) TestChannelPinsUpdate() {
	s.T().Log("Testing Channel Pings Update Unmarshal Logic")

	sub := testutil.InitSub[ChannelPinsUpdateEvent](s)

	sub.RunCommonEdgeCases()

	guildId := discord.RandomSnowflake()
	channelId := discord.RandomSnowflake()
	timestamp := testutil.RandomTime()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelPinsUpdateEvent]{
		{
			Name: "valid payload with guild_id and timestamp",
			Input: sub.MustMarshal(map[string]interface{}{
				"guild_id":           guildId,
				"channel_id":         channelId,
				"last_pin_timestamp": timestamp.Format(time.RFC3339),
			}),
			Validate: func(got ChannelPinsUpdateEvent) {
				s.EqualValues(guildId, *got.GuildID)
				s.EqualValues(channelId, got.ChannelID)
				s.EqualValues(timestamp, *got.LastPinTimestamp)
			},
		},
		{
			Name: "valid payload without guild_id and timestamp",
			Input: sub.MustMarshal(map[string]interface{}{
				"guild_id":           nil,
				"channel_id":         channelId,
				"last_pin_timestamp": nil,
			}),
			Validate: func(got ChannelPinsUpdateEvent) {
				s.Nil(got.GuildID)
				s.EqualValues(channelId, got.ChannelID)
				s.Nil(got.LastPinTimestamp)
			},
		},
	})
}
