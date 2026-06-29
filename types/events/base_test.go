package events

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/suite"
)

type eventSuite struct {
	suite.Suite
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(eventSuite))
}

func (s *eventSuite) TestAllCommandsRegistered() {
	s.Equal(77, len(eventFactories))
}

func (s *eventSuite) compareUser(expected map[string]interface{}, actual discord.User) {
	s.T().Helper()

	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["username"], actual.Username)
	s.EqualValues(expected["discriminator"], actual.Discriminator)
	s.EqualValues(expected["avatar"], *actual.Avatar)
	s.EqualValues(expected["bot"], *actual.Bot)
	s.EqualValues(expected["system"], *actual.System)
	s.EqualValues(expected["mfa_enabled"], *actual.MFAEnabled)
	s.EqualValues(expected["banner"], *actual.Banner)
	s.EqualValues(expected["accent_color"], *actual.AccentColor)
	s.EqualValues(expected["locale"], *actual.Locale)
	s.EqualValues(expected["verified"], *actual.Verified)
	s.EqualValues(expected["email"], *actual.Email)
	s.EqualValues(expected["flags"], actual.Flags)
	s.EqualValues(expected["premium_type"], *actual.PremiumType)
	s.EqualValues(expected["public_flags"], actual.PublicFlags)
	s.EqualValues(expected["global_name"], *actual.GlobalName)

	nameplate := expected["collectibles"].(map[string]interface{})["nameplate"].(map[string]interface{})
	s.EqualValues(nameplate["sku_id"], actual.Collectibles.Nameplate.SkuID)
	s.EqualValues(nameplate["asset"], actual.Collectibles.Nameplate.Asset)
	s.EqualValues(nameplate["label"], actual.Collectibles.Nameplate.Label)
	s.EqualValues(nameplate["palette"], actual.Collectibles.Nameplate.Palette)

	avatarDecorationData := expected["avatar_decoration_data"].(map[string]interface{})
	s.EqualValues(avatarDecorationData["asset"], actual.AvatarDecorationData.Asset)
	s.EqualValues(avatarDecorationData["sku_id"], actual.AvatarDecorationData.SkuID)

	primaryGuild := expected["primary_guild"].(map[string]interface{})
	s.EqualValues(primaryGuild["identity_guild_id"], *actual.PrimaryGuild.IdentityGuildID)
	s.EqualValues(primaryGuild["badge"], *actual.PrimaryGuild.Badge)
	s.EqualValues(primaryGuild["identity_enabled"], *actual.PrimaryGuild.IdentityEnabled)
	s.EqualValues(primaryGuild["tag"], *actual.PrimaryGuild.Tag)
}

func (s *eventSuite) compareChannel(expected map[string]interface{}, actual discord.Channel) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["type"], actual.Type)
	s.EqualValues(expected["guild_id"], *actual.GuildID)
	s.EqualValues(expected["position"], actual.Position)
	s.EqualValues(expected["name"], *actual.Name)
	s.EqualValues(expected["topic"], *actual.Topic)
	s.EqualValues(expected["nsfw"], *actual.NSFW)
	s.EqualValues(expected["last_message_id"], *actual.LastMessageID)
	s.EqualValues(expected["bitrate"], actual.Bitrate)
	s.EqualValues(expected["user_limit"], actual.UserLimit)
	s.EqualValues(expected["rate_limit_per_user"], actual.RateLimitPerUser)
	s.EqualValues(expected["icon"], *actual.Icon)
	s.EqualValues(expected["owner_id"], *actual.OwnerID)
	s.EqualValues(expected["application_id"], *actual.ApplicationID)
	s.EqualValues(expected["managed"], actual.Managed)
	s.EqualValues(expected["parent_id"], *actual.ParentID)
	s.EqualValues(expected["last_pin_timestamp"], *actual.LastPinTimestamp)
	s.EqualValues(expected["rtc_region"], *actual.RtcRegion)
	s.EqualValues(expected["video_quality_mode"], actual.VideoQualityMode)
	s.EqualValues(expected["message_count"], actual.MessageCount)
	s.EqualValues(expected["member_count"], actual.MemberCount)
	s.EqualValues(expected["default_auto_archive_duration"], actual.DefaultAutoArchiveDuration)
	s.EqualValues(expected["permissions"], actual.Permissions)
	s.EqualValues(expected["flags"], actual.Flags)
	s.EqualValues(expected["total_message_sent"], actual.TotalMessageSent)
	s.EqualValues(expected["applied_tags"], actual.AppliedTags)
	s.EqualValues(expected["default_thread_rate_limit_per_user"], actual.DefaultThreadRateLimitPerUser)
	s.EqualValues(expected["default_sort_order"], *actual.DefaultSortOrder)
	s.EqualValues(expected["default_forum_layout"], actual.DefaultForumLayout)

	threadMetadata := expected["thread_metadata"].(map[string]interface{})
	s.EqualValues(threadMetadata["archived"], actual.ThreadMetadata.Archived)
	s.EqualValues(threadMetadata["auto_archive_duration"], actual.ThreadMetadata.AutoArchiveDuration)
	s.EqualValues(threadMetadata["archive_timestamp"], actual.ThreadMetadata.ArchiveTimestamp)
	s.EqualValues(threadMetadata["created_timestamp"], *actual.ThreadMetadata.CreatedTimestamp)
	s.EqualValues(threadMetadata["locked"], actual.ThreadMetadata.Locked)
	s.EqualValues(threadMetadata["invitable"], actual.ThreadMetadata.Invitable)

	defaultReactionEmoji := expected["default_reaction_emoji"].(map[string]interface{})
	s.EqualValues(defaultReactionEmoji["emoji_id"], *actual.DefaultReactionEmoji.EmojiID)
	s.EqualValues(defaultReactionEmoji["emoji_name"], *actual.DefaultReactionEmoji.EmojiName)

	permissionOverwrites := expected["permission_overwrites"].([]map[string]interface{})
	s.Len(actual.PermissionOverwrites, len(permissionOverwrites))

	for i, po := range permissionOverwrites {
		s.EqualValues(po["id"], actual.PermissionOverwrites[i].ID)
		s.EqualValues(po["type"], actual.PermissionOverwrites[i].Type)
		s.EqualValues(po["allow"], actual.PermissionOverwrites[i].Allow)
		s.EqualValues(po["deny"], actual.PermissionOverwrites[i].Deny)
	}

	availableTags := expected["available_tags"].([]map[string]interface{})
	s.Len(actual.AvailableTags, len(availableTags))

	for i, tag := range availableTags {
		s.EqualValues(tag["id"], actual.AvailableTags[i].ID)
		s.EqualValues(tag["name"], actual.AvailableTags[i].Name)
		s.EqualValues(tag["moderated"], actual.AvailableTags[i].Moderated)
		s.EqualValues(tag["emoji_id"], *actual.AvailableTags[i].EmojiID)
		s.EqualValues(tag["emoji_name"], *actual.AvailableTags[i].EmojiName)
	}

	member := expected["member"].(map[string]interface{})
	s.EqualValues(member["id"], *actual.Member.ID)
	s.EqualValues(member["user_id"], *actual.Member.UserID)
	s.EqualValues(member["join_timestamp"], actual.Member.JoinTimestamp)
	s.EqualValues(member["flags"], actual.Member.Flags)
	s.compareMember(member["member"].(map[string]interface{}), *actual.Member.Member)

	recipients := expected["recipients"].([]map[string]interface{})
	s.Len(actual.Recipients, len(recipients))
	for i, r := range recipients {
		s.compareUser(r, actual.Recipients[i])
	}
}

func (s *eventSuite) compareMember(expected map[string]interface{}, actual discord.GuildMember) {
	s.EqualValues(expected["avatar"], *actual.Avatar)
	s.EqualValues(expected["banner"], *actual.Banner)
	s.EqualValues(expected["communication_disabled_until"], *actual.CommunicationDisabledUntil)
	s.EqualValues(expected["deaf"], actual.Deaf)
	s.EqualValues(expected["flags"], actual.Flags)
	s.EqualValues(expected["joined_at"], *actual.JoinedAt)
	s.EqualValues(expected["mute"], actual.Mute)
	s.EqualValues(expected["nick"], *actual.Nick)
	s.EqualValues(expected["pending"], actual.Pending)
	s.EqualValues(expected["premium_since"], *actual.PremiumSince)
	s.EqualValues(expected["roles"], actual.Roles)
	s.EqualValues(expected["permissions"], *actual.Permissions)
	s.compareUser(expected["user"].(map[string]interface{}), *actual.User)

	nameplate := expected["collectibles"].(map[string]interface{})["nameplate"].(map[string]interface{})
	s.EqualValues(nameplate["sku_id"], actual.Collectibles.Nameplate.SkuID)
	s.EqualValues(nameplate["asset"], actual.Collectibles.Nameplate.Asset)
	s.EqualValues(nameplate["label"], actual.Collectibles.Nameplate.Label)
	s.EqualValues(nameplate["palette"], actual.Collectibles.Nameplate.Palette)

	avatarDecorationData := expected["avatar_decoration_data"].(map[string]interface{})
	s.EqualValues(avatarDecorationData["asset"], actual.AvatarDecorationData.Asset)
	s.EqualValues(avatarDecorationData["sku_id"], actual.AvatarDecorationData.SkuID)
}

func (s *eventSuite) compareRole(expected map[string]interface{}, actual discord.Role) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["hoist"], actual.Hoist)
	s.EqualValues(expected["icon"], *actual.Icon)
	s.EqualValues(expected["unicode_emoji"], *actual.UnicodeEmoji)
	s.EqualValues(expected["position"], actual.Position)
	s.EqualValues(expected["permissions"], actual.Permissions)
	s.EqualValues(expected["managed"], actual.Managed)
	s.EqualValues(expected["mentionable"], actual.Mentionable)
	s.Equal(expected["flags"], *actual.Flags)

	colors := expected["colors"].(map[string]interface{})
	s.EqualValues(colors["primary_color"], actual.Colors.PrimaryColor)
	s.EqualValues(colors["secondary_color"], *actual.Colors.SecondaryColor)
	s.EqualValues(colors["tertiary_color"], *actual.Colors.TertiaryColor)

	tags := expected["tags"].(map[string]interface{})
	s.EqualValues(tags["bot_id"], *actual.Tags.BotID)
	s.EqualValues(tags["integration_id"], *actual.Tags.IntegrationID)
	s.EqualValues(tags["subscription_listing_id"], *actual.Tags.SubscriptionListingID)
	s.EqualValues(s.resolveNullFlag(tags["premium_subscriber"]), actual.Tags.PremiumSubscriber)
	s.EqualValues(s.resolveNullFlag(tags["available_for_purchase"]), actual.Tags.AvailableForPurchase)
	s.EqualValues(s.resolveNullFlag(tags["guild_connections"]), actual.Tags.GuildConnections)
}

func (s *eventSuite) resolveNullFlag(input interface{}) interface{} {
	if input == nil {
		return discord.NullFlag(false)
	}

	return discord.NullFlag(true)
}

func (s *eventSuite) compareSticker(expected map[string]interface{}, actual discord.Sticker) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["pack_id"], *actual.PackID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["description"], *actual.Description)
	s.EqualValues(expected["tags"], actual.Tags)
	s.EqualValues(expected["type"], actual.Type)
	s.EqualValues(expected["format_type"], actual.FormatType)
	s.EqualValues(expected["available"], *actual.Available)
	s.EqualValues(expected["guild_id"], *actual.GuildID)
	s.EqualValues(expected["sort_value"], actual.SortValue)
	s.compareUser(expected["user"].(map[string]interface{}), *actual.User)
}

func (s *eventSuite) compareEmoji(expected map[string]interface{}, actual discord.Emoji) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["roles"], actual.Roles)
	s.EqualValues(expected["require_colons"], *actual.RequireColons)
	s.EqualValues(expected["animated"], *actual.Animated)
	s.EqualValues(expected["available"], *actual.Available)
	s.EqualValues(expected["managed"], *actual.Managed)

	s.compareUser(expected["user"].(map[string]interface{}), *actual.User)
}

func (s *eventSuite) comparePresence(expected map[string]interface{}, actual discord.Presence) {
	s.compareUser(expected["user"].(map[string]interface{}), actual.User)

	s.EqualValues(expected["guild_id"], actual.GuildID)
	s.EqualValues(expected["status"], actual.Status)

	clientStatus := expected["client_status"].(map[string]interface{})
	s.EqualValues(clientStatus["desktop"], actual.ClientStatus.Desktop)
	s.EqualValues(clientStatus["mobile"], actual.ClientStatus.Mobile)
	s.EqualValues(clientStatus["web"], actual.ClientStatus.Web)

	expectedActivities := expected["activities"].([]interface{})
	s.Len(actual.Activities, len(expectedActivities))

	for i, activityRaw := range expectedActivities {
		expectedActivity := activityRaw.(map[string]interface{})
		actualActivity := actual.Activities[i]

		s.EqualValues(expectedActivity["type"], actualActivity.Type)
		s.EqualValues(expectedActivity["name"], actualActivity.Name)
		s.EqualValues(expectedActivity["url"], *actualActivity.URL)
		s.EqualValues(
			expectedActivity["created_at"],
			actualActivity.CreatedAt,
		)

		s.EqualValues(
			expectedActivity["application_id"],
			*actualActivity.ApplicationID,
		)

		s.EqualValues(
			expectedActivity["status_display_type"],
			*actualActivity.StatusDisplayType,
		)

		s.EqualValues(expectedActivity["details"], *actualActivity.Details)
		s.EqualValues(expectedActivity["state"], *actualActivity.State)
		s.EqualValues(expectedActivity["instance"], *actualActivity.Instance)
		s.EqualValues(expectedActivity["flags"], *actualActivity.Flags)

		timestamps := expectedActivity["timestamps"].(map[string]interface{})
		s.EqualValues(
			timestamps["start"],
			actualActivity.Timestamps.Start,
		)
		s.EqualValues(
			timestamps["end"],
			actualActivity.Timestamps.End,
		)

		party := expectedActivity["party"].(map[string]interface{})
		s.EqualValues(party["id"], *actualActivity.Party.ID)

		expectedSize := party["size"].([]int)
		s.EqualValues(expectedSize[0], actualActivity.Party.Size[0])
		s.EqualValues(expectedSize[1], actualActivity.Party.Size[1])

		assets := expectedActivity["assets"].(map[string]interface{})
		s.EqualValues(assets["large_image"], *actualActivity.Assets.LargeImage)
		s.EqualValues(assets["large_text"], *actualActivity.Assets.LargeText)
		s.EqualValues(assets["large_url"], *actualActivity.Assets.LargeURL)
		s.EqualValues(assets["small_image"], *actualActivity.Assets.SmallImage)
		s.EqualValues(assets["small_text"], *actualActivity.Assets.SmallText)
		s.EqualValues(assets["small_url"], *actualActivity.Assets.SmallURL)
		s.EqualValues(assets["invite_cover_image"], *actualActivity.Assets.InviteCoverImage)

		secrets := expectedActivity["secrets"].(map[string]interface{})
		s.EqualValues(secrets["join"], *actualActivity.Secrets.Join)
		s.EqualValues(secrets["spectate"], *actualActivity.Secrets.Spectate)
		s.EqualValues(secrets["match"], *actualActivity.Secrets.Match)

		emoji := expectedActivity["emoji"].(map[string]interface{})
		s.EqualValues(emoji["id"], *actualActivity.Emoji.ID)
		s.EqualValues(emoji["name"], actualActivity.Emoji.Name)
		s.EqualValues(emoji["animated"], *actualActivity.Emoji.Animated)
	}
}
