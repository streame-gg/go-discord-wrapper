package interactions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
	"github.com/stretchr/testify/suite"
)

type functionSuite struct {
	suite.Suite
}

func TestFunctionSuite(t *testing.T) {
	suite.Run(t, new(functionSuite))
}

func (s *functionSuite) TestGetSubcommands() {
	interaction := Interaction{
		Type: discord.InteractionTypeApplicationCommand,
		Data: &responses.InteractionDataApplicationCommand{
			Name: "name",
			Type: discord.ApplicationCommandTypeChatInput,
			Options: []responses.ApplicationCommandInteractionDataOption[interface{}]{
				{
					Name: "groupsubcommand",
					Type: discord.ApplicationCommandOptionTypeSubCommandGroup,
					Options: []responses.ApplicationCommandInteractionDataOption[interface{}]{
						{
							Name: "subcommand",
							Type: discord.ApplicationCommandOptionTypeSubCommand,
						},
					},
				},
			},
		},
	}

	s.EqualValues(interaction.GetSubCommand(), "subcommand")
	s.EqualValues(interaction.GetSubCommandGroup(), "groupsubcommand")
	s.EqualValues(interaction.GetCommandName(), "name")
	s.EqualValues(interaction.GetFullCommand(), "name groupsubcommand subcommand")

	interaction = Interaction{
		Type: discord.InteractionTypeApplicationCommand,
		Data: &responses.InteractionDataApplicationCommand{
			Name: "name",
			Type: discord.ApplicationCommandTypeChatInput,
			Options: []responses.ApplicationCommandInteractionDataOption[interface{}]{
				{
					Name: "subcommand",
					Type: discord.ApplicationCommandOptionTypeSubCommand,
				},
			},
		},
	}

	s.EqualValues(interaction.GetSubCommand(), "subcommand")
	s.EqualValues(interaction.GetSubCommandGroup(), "")
	s.EqualValues(interaction.GetCommandName(), "name")
	s.EqualValues(interaction.GetFullCommand(), "name subcommand")

	interaction = Interaction{
		Type: discord.InteractionTypeApplicationCommand,
		Data: &responses.InteractionDataApplicationCommand{
			Name: "name",
			Type: discord.ApplicationCommandTypeChatInput,
		},
	}

	s.EqualValues(interaction.GetSubCommand(), "")
	s.EqualValues(interaction.GetSubCommandGroup(), "")
	s.EqualValues(interaction.GetCommandName(), "name")
	s.EqualValues(interaction.GetFullCommand(), "name")
}

func (s *functionSuite) TestGetSubcommands_InvalidOptions() {
	interaction := Interaction{
		Type: discord.InteractionTypePing,
		Data: nil,
	}

	s.EqualValues(interaction.GetSubCommand(), "")
	s.EqualValues(interaction.GetSubCommandGroup(), "")
	s.EqualValues(interaction.GetCommandName(), "")
	s.EqualValues(interaction.GetFullCommand(), "")
}

func (s *functionSuite) TestCommandUtil() {
	interaction := Interaction{
		Type: discord.InteractionTypePing,
		Data: nil,
		Member: &discord.GuildMember{
			User: &discord.User{
				ID: discord.Snowflake(123123123123123123),
			},
		},
	}

	s.EqualValues(interaction.InvokingUser().ID, discord.Snowflake(123123123123123123))
	s.EqualValues(interaction.UserID(), discord.Snowflake(123123123123123123))

	interaction = Interaction{
		Type: discord.InteractionTypePing,
		Data: nil,
		User: &discord.User{
			ID: discord.Snowflake(213123123123123123),
		},
		Member: &discord.GuildMember{},
	}

	s.EqualValues(interaction.InvokingUser().ID, discord.Snowflake(213123123123123123))
	s.EqualValues(interaction.UserID(), discord.Snowflake(213123123123123123))

	interaction = Interaction{
		Type:    discord.InteractionTypePing,
		Data:    nil,
		GuildID: util.PointerOf(discord.RandomSnowflake()),
	}

	s.True(interaction.InGuild())

	interaction = Interaction{
		Type:    discord.InteractionTypePing,
		Data:    nil,
		GuildID: nil,
	}

	s.False(interaction.InGuild())
}

func (s *functionSuite) TestComponentUtil() {
	interaction := Interaction{
		Data: &responses.InteractionDataMessageComponent{
			CustomID: "lol",
		},
	}

	s.EqualValues(interaction.GetCustomID(), "lol")
}

func (s *functionSuite) TestModalUtil() {
	interaction := Interaction{
		Data: &responses.InteractionDataModalSubmit{
			CustomID: "lol",
		},
	}

	s.EqualValues(interaction.GetCustomID(), "lol")
}

// TODO check if this really works in prod environment (if the resolved data is really sending the params I expect)
func (s *functionSuite) TestOptions() {
	user := testdata.NewUser()
	role := testdata.NewRole()
	channel := testdata.NewChannel()
	attachment := testdata.NewAttachment()

	resolvedUser := mustConvert[discord.User](s.T(), user)
	resolvedRole := mustConvert[discord.Role](s.T(), role)
	resolvedChannel := mustConvert[discord.Channel](s.T(), channel)
	resolvedAttachment := mustConvert[discord.Attachment](s.T(), attachment)

	interaction := Interaction{
		Type: discord.InteractionTypeApplicationCommand,
		Data: &responses.InteractionDataApplicationCommand{
			Resolved: &discord.ResolvedData{
				Users: map[discord.Snowflake]discord.User{
					resolvedUser.ID: resolvedUser,
				},
				Channels: map[discord.Snowflake]discord.Channel{
					resolvedChannel.ID: resolvedChannel,
				},
				Roles: map[discord.Snowflake]discord.Role{
					resolvedRole.ID: resolvedRole,
				},
				Attachments: map[discord.Snowflake]discord.Attachment{
					resolvedAttachment.ID: resolvedAttachment,
				},
			},
			Name: "name",
			Type: discord.ApplicationCommandTypeChatInput,
			Options: []responses.ApplicationCommandInteractionDataOption[interface{}]{
				{
					Name:  "string",
					Type:  discord.ApplicationCommandOptionTypeString,
					Value: util.PointerOf(any("test!")),
				},
				{
					Name:  "integer",
					Type:  discord.ApplicationCommandOptionTypeInteger,
					Value: util.PointerOf(any(1)),
				},
				{
					Name:  "boolean",
					Type:  discord.ApplicationCommandOptionTypeBoolean,
					Value: util.PointerOf(any(true)),
				},
				{
					Name:  "user",
					Type:  discord.ApplicationCommandOptionTypeUser,
					Value: util.PointerOf(user["id"]),
				},
				{
					Name:  "channel",
					Type:  discord.ApplicationCommandOptionTypeChannel,
					Value: util.PointerOf(channel["id"]),
				},
				{
					Name:  "role",
					Type:  discord.ApplicationCommandOptionTypeRole,
					Value: util.PointerOf(role["id"]),
				},
				{
					Name:  "mentionable_user",
					Type:  discord.ApplicationCommandOptionTypeMentionable,
					Value: util.PointerOf(user["id"]),
				},
				{
					Name:  "mentionable_role",
					Type:  discord.ApplicationCommandOptionTypeMentionable,
					Value: util.PointerOf(role["id"]),
				},
				{
					Name:  "number",
					Type:  discord.ApplicationCommandOptionTypeNumber,
					Value: util.PointerOf(any(1.0)),
				},
				{
					Name:  "attachment",
					Type:  discord.ApplicationCommandOptionTypeAttachment,
					Value: util.PointerOf(attachment["id"]),
				},
			},
		},
	}

	s.Equal(interaction.GetStringOption("string"), "test!")
	s.Equal(interaction.GetIntegerOption("integer"), 1)
	s.Equal(interaction.GetBoolOption("boolean"), true)
	s.Equal(interaction.GetNumberOption("number"), 1.0)
	s.compareUser(user, *(interaction.GetUserOption("user")))
	s.compareAttachment(attachment, *(interaction.GetAttachmentOption("attachment")))
	s.compareRole(role, *(interaction.GetRoleOption("role")))
	s.compareChannel(channel, *(interaction.GetChannelOption("channel")))

	res := interaction.GetMentionableOption("mentionable_role")
	s.Require().NotNil(res)
	s.compareRole(role, *res.Role)
	s.Equal(MentionableOptionRole, res.Type)

	res = interaction.GetMentionableOption("mentionable_user")
	s.Require().NotNil(res)
	s.compareUser(user, *res.User)
	s.Equal(MentionableOptionUser, res.Type)
}

func (s *functionSuite) compareUser(expected map[string]interface{}, actual discord.User) {
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

func (s *functionSuite) compareChannel(expected map[string]interface{}, actual discord.Channel) {
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

func (s *functionSuite) compareMember(expected map[string]interface{}, actual discord.GuildMember) {
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

func (s *functionSuite) compareRole(expected map[string]interface{}, actual discord.Role) {
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

func (s *functionSuite) resolveNullFlag(input interface{}) interface{} {
	if input == nil {
		return discord.NullFlag(false)
	}

	return discord.NullFlag(true)
}

func (s *functionSuite) compareSticker(expected map[string]interface{}, actual discord.Sticker) {
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

func (s *functionSuite) compareEmoji(expected map[string]interface{}, actual discord.Emoji) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["roles"], *actual.Roles)
	s.EqualValues(expected["require_colons"], *actual.RequireColons)
	s.EqualValues(expected["animated"], *actual.Animated)
	s.EqualValues(expected["available"], *actual.Available)
	s.EqualValues(expected["managed"], *actual.Managed)

	s.compareUser(expected["user"].(map[string]interface{}), *actual.User)
}

func (s *functionSuite) comparePresence(expected map[string]interface{}, actual discord.Presence) {
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

func (s *functionSuite) compareGuildScheduledEvent(expected map[string]interface{}, actual discord.GuildScheduledEvent) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["guild_id"], actual.GuildID)
	s.EqualValues(expected["channel_id"], *actual.ChannelID)
	s.EqualValues(expected["creator_id"], *actual.CreatorID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["description"], *actual.Description)
	s.EqualValues(expected["scheduled_start_time"], actual.ScheduledStartTime)
	s.EqualValues(expected["scheduled_end_time"], *actual.ScheduledEndTime)
	s.EqualValues(expected["privacy_level"], actual.PrivacyLevel)
	s.EqualValues(expected["status"], actual.Status)
	s.EqualValues(expected["entity_type"], actual.EntityType)
	s.EqualValues(expected["user_count"], *actual.UserCount)
	s.EqualValues(expected["image"], *actual.Image)

	s.compareUser(expected["creator"].(map[string]interface{}), *actual.Creator)

	entityMetadata := expected["entity_metadata"].(map[string]interface{})
	s.EqualValues(entityMetadata["location"], *actual.EntityMetadata.Location)

	recurrenceRule := expected["recurrence_rule"].(map[string]interface{})
	s.EqualValues(recurrenceRule["start"], actual.RecurrenceRule.Start)
	s.EqualValues(recurrenceRule["end"], *actual.RecurrenceRule.End)
	s.EqualValues(recurrenceRule["frequency"], actual.RecurrenceRule.Frequency)
	s.EqualValues(recurrenceRule["interval"], actual.RecurrenceRule.Interval)
	s.EqualValues(recurrenceRule["count"], *actual.RecurrenceRule.Count)

	expectedByWeekday := recurrenceRule["by_weekday"].([]discord.GuildScheduledEventRecurrenceRuleWeekday)
	s.Require().Len(actual.RecurrenceRule.ByWeekday, len(expectedByWeekday))
	for i, exp := range expectedByWeekday {
		s.EqualValues(exp, actual.RecurrenceRule.ByWeekday[i])
	}

	expectedByNWeekday := recurrenceRule["by_n_weekday"].([]map[string]interface{})
	s.Require().Len(actual.RecurrenceRule.ByNWeekday, len(expectedByNWeekday))
	for i, exp := range expectedByNWeekday {
		s.EqualValues(exp["n"], actual.RecurrenceRule.ByNWeekday[i].N)
		s.EqualValues(exp["day"], actual.RecurrenceRule.ByNWeekday[i].Day)
	}

	expectedByMonth := recurrenceRule["by_month"].([]discord.GuildScheduledEventRecurrenceRuleMonth)
	s.Require().Len(actual.RecurrenceRule.ByMonth, len(expectedByMonth))
	for i, exp := range expectedByMonth {
		s.EqualValues(exp, actual.RecurrenceRule.ByMonth[i])
	}

	expectedByMonthDay := recurrenceRule["by_month_day"].([]int)
	s.Require().Len(actual.RecurrenceRule.ByMonthDay, len(expectedByMonthDay))
	for i, exp := range expectedByMonthDay {
		s.EqualValues(exp, actual.RecurrenceRule.ByMonthDay[i])
	}

	expectedByYearDay := recurrenceRule["by_year_day"].([]int)
	s.Require().Len(actual.RecurrenceRule.ByYearDay, len(expectedByYearDay))
	for i, exp := range expectedByYearDay {
		s.EqualValues(exp, actual.RecurrenceRule.ByYearDay[i])
	}
}

func (s *functionSuite) compareGuild(expected map[string]interface{}, actual discord.Guild) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["icon"], *actual.Icon)
	s.EqualValues(expected["icon_hash"], *actual.IconHash)
	s.EqualValues(expected["splash"], *actual.Splash)
	s.EqualValues(expected["discovery_splash"], *actual.DiscoverySplash)
	s.EqualValues(expected["owner"], actual.Owner)
	s.EqualValues(expected["owner_id"], actual.OwnerID)
	s.EqualValues(expected["permissions"], actual.Permissions)
	s.EqualValues(expected["region"], *actual.Region)
	s.EqualValues(expected["afk_channel_id"], *actual.AfkChannelID)
	s.EqualValues(expected["afk_timeout"], actual.AfkTimeout)
	s.EqualValues(expected["widget_enabled"], actual.WidgetEnabled)
	s.EqualValues(expected["widget_channel_id"], *actual.WidgetChannelID)
	s.EqualValues(expected["verification_level"], actual.VerificationLevel)
	s.EqualValues(expected["default_message_notifications"], actual.DefaultMessageNotifications)
	s.EqualValues(expected["explicit_content_filter"], actual.ExplicitContentFilter)
	s.EqualValues(expected["mfa_level"], actual.MfaLevel)
	s.EqualValues(expected["application_id"], *actual.ApplicationID)
	s.EqualValues(expected["system_channel_id"], *actual.SystemChannelID)
	s.EqualValues(expected["system_channel_flags"], actual.SystemChannelFlags)
	s.EqualValues(expected["rules_channel_id"], *actual.RulesChannelID)
	s.EqualValues(expected["max_presences"], *actual.MaxPresences)
	s.EqualValues(expected["max_members"], actual.MaxMembers)
	s.EqualValues(expected["vanity_url_code"], *actual.VanityUrlCode)
	s.EqualValues(expected["description"], *actual.Description)
	s.EqualValues(expected["banner"], *actual.Banner)
	s.EqualValues(expected["premium_tier"], actual.PremiumTier)
	s.EqualValues(expected["premium_subscription_count"], actual.PremiumSubscriptionCount)
	s.EqualValues(expected["preferred_locale"], actual.PreferredLocale)
	s.EqualValues(expected["public_updates_channel_id"], *actual.PublicUpdatesChannelID)
	s.EqualValues(expected["max_video_channel_users"], actual.MaxVideoChannelUsers)
	s.EqualValues(expected["max_stage_video_channel_users"], actual.MaxStageVideoChannelUsers)
	s.EqualValues(expected["approximate_member_count"], actual.ApproximateMemberCount)
	s.EqualValues(expected["approximate_presence_count"], actual.ApproximatePresenceCount)
	s.EqualValues(expected["nsfw_level"], actual.NSFWLevel)
	s.EqualValues(expected["premium_progress_bar_enabled"], actual.PremiumProgressBarEnabled)
	s.EqualValues(expected["safety_alerts_channel_id"], *actual.SafetyAlertsChannelID)

	welcomeScreen := expected["welcome_screen"].(map[string]interface{})
	s.EqualValues(welcomeScreen["description"], *actual.WelcomeScreen.Description)
	welcomeChannels := welcomeScreen["welcome_channels"].([]map[string]interface{})
	for i, wc := range welcomeChannels {
		s.EqualValues(wc["channel_id"], actual.WelcomeScreen.WelcomeChannels[i].ChannelID)
		s.EqualValues(wc["description"], actual.WelcomeScreen.WelcomeChannels[i].Description)
		s.EqualValues(wc["emoji_id"], *actual.WelcomeScreen.WelcomeChannels[i].EmojiID)
		s.EqualValues(wc["emoji_name"], *actual.WelcomeScreen.WelcomeChannels[i].EmojiName)
	}

	incidentsData := expected["incidents_data"].(map[string]interface{})
	s.EqualValues(incidentsData["invites_disabled_until"], *actual.IncidentsData.InvitesDisabledUntil)
	s.EqualValues(incidentsData["dms_disabled_until"], *actual.IncidentsData.DmsDisabledUntil)
	s.EqualValues(incidentsData["dm_spam_detected_at"], *actual.IncidentsData.DmSpanDetectedAt)
	s.EqualValues(incidentsData["raid_detected_at"], *actual.IncidentsData.RaidDetectedAt)

	roles := expected["roles"].([]map[string]interface{})
	s.Len(actual.RawRoles, len(roles))

	for i, r := range roles {
		s.compareRole(r, actual.RawRoles[i])
	}

	stickers := expected["stickers"].([]map[string]interface{})
	s.Len(actual.RawStickers, len(stickers))

	for i, st := range stickers {
		s.compareSticker(st, actual.RawStickers[i])
	}

	emojis := expected["emojis"].([]map[string]interface{})
	s.Len(actual.RawEmojis, len(emojis))

	for i, em := range emojis {
		s.compareEmoji(em, actual.RawEmojis[i])
	}
}

func (s *functionSuite) compareSoundboardSound(expected map[string]interface{}, actual discord.SoundboardSound) {
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["sound_id"], actual.SoundID)
	s.EqualValues(expected["volume"], actual.Volume)
	s.EqualValues(expected["emoji_id"], *actual.EmojiID)
	s.EqualValues(expected["emoji_name"], *actual.EmojiName)
	s.EqualValues(expected["guild_id"], *actual.GuildID)
	s.EqualValues(expected["available"], actual.Available)
	s.compareUser(expected["user"].(map[string]interface{}), *actual.User)
}

func (s *functionSuite) compareStageInstance(expected map[string]interface{}, actual discord.StageInstance) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["guild_id"], actual.GuildID)
	s.EqualValues(expected["channel_id"], actual.ChannelID)
	s.EqualValues(expected["topic"], actual.Topic)
	s.EqualValues(expected["privacy_level"], actual.PrivacyLevel)
	s.EqualValues(expected["discoverable_disabled"], actual.DiscoverableDisabled)
	s.EqualValues(expected["guild_scheduled_event_id"], *actual.GuildScheduledEventID)
}

func (s *functionSuite) compareVoiceState(expected map[string]interface{}, actual discord.VoiceState) {
	s.EqualValues(expected["guild_id"], *actual.GuildID)
	s.EqualValues(expected["channel_id"], *actual.ChannelID)
	s.EqualValues(expected["user_id"], actual.UserID)
	s.EqualValues(expected["session_id"], actual.SessionID)
	s.EqualValues(expected["deaf"], actual.Deaf)
	s.EqualValues(expected["mute"], actual.Mute)
	s.EqualValues(expected["self_deaf"], actual.SelfDeaf)
	s.EqualValues(expected["self_mute"], actual.SelfMute)
	s.EqualValues(expected["self_stream"], *actual.SelfStream)
	s.EqualValues(expected["self_video"], actual.SelfVideo)
	s.EqualValues(expected["suppress"], actual.Suppress)
	s.EqualValues(expected["request_to_speak_timestamp"], *actual.RequestToSpeakTimestamp)

	s.compareMember(expected["member"].(map[string]interface{}), *actual.Member)
}

func (s *functionSuite) compareIntegration(expected map[string]interface{}, actual discord.Integration) {
	s.Require().NotNil(actual)
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["type"], actual.Type)
	s.EqualValues(expected["enabled"], actual.Enabled)
	s.EqualValues(expected["syncing"], actual.Syncing)
	s.EqualValues(expected["role_id"], *actual.RoleID)
	s.EqualValues(expected["enable_emoticons"], actual.EnableEmoticons)
	s.EqualValues(expected["expire_behavior"], *actual.ExpireBehavior)
	s.EqualValues(expected["expire_grace_period"], actual.ExpireGracePeriod)
	s.EqualValues(expected["synced_at"], *actual.SyncedAt)
	s.EqualValues(expected["subscriber_count"], actual.SubscriberCount)
	s.EqualValues(expected["revoked"], actual.Revoked)
	s.EqualValues(expected["scopes"], actual.Scopes)

	s.compareUser(expected["user"].(map[string]interface{}), *actual.User)

	account := expected["account"].(map[string]interface{})
	s.Require().NotNil(account)
	s.EqualValues(account["id"], actual.Account.ID)
	s.EqualValues(account["name"], actual.Account.Name)

	application := expected["application"].(map[string]interface{})
	s.Require().NotNil(application)
	s.EqualValues(application["id"], actual.Application.ID)
	s.EqualValues(application["name"], actual.Application.Name)
	s.EqualValues(application["description"], actual.Application.Description)
	s.EqualValues(application["icon"], *actual.Application.Icon)
	s.compareUser(application["bot"].(map[string]interface{}), *actual.Application.Bot)
}

func (s *functionSuite) compareThreadMember(expected map[string]interface{}, actual discord.ThreadMember) {
	s.EqualValues(expected["id"], *actual.ID)
	s.EqualValues(expected["user_id"], *actual.UserID)
	s.EqualValues(expected["join_timestamp"], actual.JoinTimestamp)
	s.EqualValues(expected["flags"], actual.Flags)
	s.compareMember(expected["member"].(map[string]interface{}), *actual.Member)
}

func (s *functionSuite) compareSubscription(expected map[string]interface{}, actual discord.Subscription) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["user_id"], actual.UserID)
	s.EqualValues(expected["sku_ids"], actual.SKUIDs)
	s.EqualValues(expected["entitlement_ids"], actual.EntitlementIDs)
	s.EqualValues(expected["renewal_sku_ids"], actual.RenewalSKUIDs)
	s.EqualValues(expected["current_period_start"], actual.CurrentPeriodStart)
	s.EqualValues(expected["current_period_end"], actual.CurrentPeriodEnd)
	s.EqualValues(expected["status"], actual.Status)
	s.EqualValues(expected["canceled_at"], *actual.CanceledAt)
	s.EqualValues(expected["country"], actual.Country)
}

func (s *functionSuite) compareEntitlement(expected map[string]interface{}, actual discord.Entitlement) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["sku_id"], actual.SkuID)
	s.EqualValues(expected["application_id"], actual.ApplicationID)
	s.EqualValues(expected["user_id"], *actual.UserID)
	s.EqualValues(expected["type"], actual.Type)
	s.EqualValues(expected["deleted"], actual.Deleted)
	s.EqualValues(expected["starts_at"], *actual.StartsAt)
	s.EqualValues(expected["ends_at"], *actual.EndsAt)
	s.EqualValues(expected["guild_id"], *actual.GuildID)
	s.EqualValues(expected["consumed"], *actual.Consumed)
}

func (s *functionSuite) compareResolved(expected map[string]interface{}, actual discord.ResolvedData) {
	users := expected["users"].(map[discord.Snowflake]interface{})
	s.Len(actual.Users, len(users))

	for i, user := range actual.Users {
		s.compareUser(users[i].(map[string]interface{}), user)
	}

	members := expected["members"].(map[discord.Snowflake]interface{})
	s.Len(actual.Members, len(members))

	for i, member := range actual.Members {
		s.compareMember(members[i].(map[string]interface{}), member)
	}

	messages := expected["messages"].(map[discord.Snowflake]interface{})
	s.Len(actual.Messages, len(messages))

	for i, message := range actual.Messages {
		s.compareMessage(messages[i].(map[string]interface{}), message)
	}

	channels := expected["channels"].(map[discord.Snowflake]interface{})
	s.Len(actual.Channels, len(channels))

	for i, channel := range actual.Channels {
		s.compareChannel(channels[i].(map[string]interface{}), channel)
	}

	roles := expected["roles"].(map[discord.Snowflake]interface{})
	s.Len(actual.Roles, len(roles))

	for i, role := range actual.Roles {
		s.compareRole(roles[i].(map[string]interface{}), role)
	}

	attachments := expected["attachments"].(map[discord.Snowflake]interface{})
	s.Len(actual.Attachments, len(attachments))

	for i, attachment := range actual.Attachments {
		s.compareAttachment(attachments[i].(map[string]interface{}), attachment)
	}
}

func (s *functionSuite) compareAttachment(expected map[string]interface{}, actual discord.Attachment) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["filename"], actual.Filename)
	s.EqualValues(expected["title"], actual.Title)
	s.EqualValues(expected["description"], actual.Description)
	s.EqualValues(expected["content_type"], actual.ContentType)
	s.EqualValues(expected["size"], actual.Size)
	s.EqualValues(expected["url"], actual.URL)
	s.EqualValues(expected["proxy_url"], actual.ProxyURL)
	s.EqualValues(expected["height"], *actual.Height)
	s.EqualValues(expected["width"], *actual.Width)
	s.EqualValues(expected["placeholder"], actual.Placeholder)
	s.EqualValues(expected["placeholder_version"], actual.PlaceholderVersion)
	s.EqualValues(expected["ephemeral"], actual.Ephemeral)
	s.EqualValues(expected["duration_secs"], actual.DurationSecs)
	s.EqualValues(expected["waveform"], actual.Waveform)
	s.EqualValues(expected["flags"], actual.Flags)
	s.EqualValues(expected["clip_created_at"], actual.ClipCreatedAt)

	users := expected["clip_participants"].([]map[string]interface{})
	s.Len(actual.ClipParticipants, len(users))

	for i, user := range actual.ClipParticipants {
		s.compareUser(users[i], user)
	}

	s.compareApplication(expected["application"].(map[string]interface{}), *actual.Application)
}

func (s *functionSuite) compareApplication(expected map[string]interface{}, actual discord.Application) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["icon"], *actual.Icon)
	s.EqualValues(expected["description"], actual.Description)
	s.EqualValues(expected["rpc_origins"], actual.RpcOrigins)
	s.EqualValues(expected["bot_public"], actual.BotPublic)
	s.EqualValues(expected["bot_require_code_grant"], actual.BotRequireCodeGrant)
	s.EqualValues(expected["terms_of_service_url"], actual.TermsOfServiceURL)
	s.EqualValues(expected["privacy_policy_url"], actual.PrivacyPolicyURL)
	s.EqualValues(expected["verify_key"], actual.VerifyKey)
	s.EqualValues(expected["guild_id"], *actual.GuildID)
	s.EqualValues(expected["primary_sku_id"], *actual.PrimarySKUID)
	s.EqualValues(expected["slug"], actual.Slug)
	s.EqualValues(expected["cover_image"], actual.CoverImage)
	s.EqualValues(expected["flags"], *actual.Flags)
	s.EqualValues(expected["flags_new"], actual.FlagsNew.String())
	s.EqualValues(expected["approximate_guild_count"], *actual.ApproximateGuildCount)
	s.EqualValues(expected["approximate_user_install_count"], *actual.ApproximateUserInstallCount)
	s.EqualValues(expected["approximate_user_authorization_count"], *actual.ApproximateUserAuthorizationCount)
	s.EqualValues(expected["redirect_uris"], actual.RedirectURIs)
	s.EqualValues(expected["interactions_endpoint_url"], *actual.InteractionsEndpointURL)
	s.EqualValues(expected["role_connections_verification_url"], *actual.RoleConnectionsVerificationURL)
	s.EqualValues(expected["event_webhooks_url"], *actual.EventWebhooksURL)
	s.EqualValues(expected["event_webhooks_status"], *actual.EventWebhooksStatus)
	s.EqualValues(expected["event_webhooks_types"], actual.EventWebhooksTypes)
	s.EqualValues(expected["tags"], actual.Tags)
	s.EqualValues(expected["custom_install_url"], actual.CustomInstallURL)

	installParams := expected["install_params"].(map[string]interface{})
	s.EqualValues(installParams["permissions"], actual.InstallParams.Permissions)
	s.EqualValues(installParams["scopes"], actual.InstallParams.Scopes)

	integrationTypesConfig := expected["integration_types_config"].(map[discord.ApplicationIntegrationType]map[string]interface{})
	s.EqualValues(integrationTypesConfig["0"]["oauth2_install_params"].(map[string]interface{})["permissions"], actual.IntegrationTypesConfig[discord.ApplicationIntegrationTypeGuildInstall].OAuth2InstallParams.Permissions)
	s.EqualValues(integrationTypesConfig["0"]["oauth2_install_params"].(map[string]interface{})["scopes"], actual.IntegrationTypesConfig[discord.ApplicationIntegrationTypeGuildInstall].OAuth2InstallParams.Scopes)
	s.EqualValues(integrationTypesConfig["1"]["oauth2_install_params"].(map[string]interface{})["permissions"], actual.IntegrationTypesConfig[discord.ApplicationIntegrationTypeUserInstall].OAuth2InstallParams.Permissions)
	s.EqualValues(integrationTypesConfig["1"]["oauth2_install_params"].(map[string]interface{})["scopes"], actual.IntegrationTypesConfig[discord.ApplicationIntegrationTypeUserInstall].OAuth2InstallParams.Scopes)

	team := expected["team"].(map[string]interface{})
	s.EqualValues(team["id"], actual.Team.ID)
	s.EqualValues(team["icon"], *actual.Team.Icon)
	s.EqualValues(team["name"], actual.Team.Name)
	s.EqualValues(team["owner_user_id"], actual.Team.OwnerUserID)

	teamMembers := team["members"].([]map[string]interface{})
	s.Len(actual.Team.Members, len(teamMembers))

	for i, teamMember := range actual.Team.Members {
		s.EqualValues(teamMembers[i]["team_id"], teamMember.TeamID)
		s.EqualValues(teamMembers[i]["membership_state"], teamMember.MembershipState)
		s.EqualValues(teamMembers[i]["role"], teamMember.Role)

		s.compareUser(teamMembers[i]["user"].(map[string]interface{}), teamMember.User)
	}

	s.compareUser(expected["bot"].(map[string]interface{}), *actual.Bot)
	s.compareUser(expected["owner"].(map[string]interface{}), *actual.Owner)
	s.compareGuild(expected["guild"].(map[string]interface{}), *actual.Guild)
}

func (s *functionSuite) compareMessage(expected map[string]interface{}, actual discord.Message) {
	s.EqualValues(expected["id"], actual.ID)
	s.EqualValues(expected["channel_id"], actual.ChannelID)
	s.EqualValues(expected["content"], actual.Content)
	s.EqualValues(expected["timestamp"], actual.Timestamp)
	s.EqualValues(expected["edited_timestamp"], *actual.EditedTimestamp)
	s.EqualValues(expected["tts"], actual.TTS)
	s.EqualValues(expected["mention_everyone"], actual.MentionEveryone)
	s.EqualValues(expected["mention_roles"], actual.MentionRoles)
	s.EqualValues(expected["nonce"], actual.Nonce)
	s.EqualValues(expected["pinned"], actual.Pinned)
	s.EqualValues(expected["webhook_id"], *actual.WebhookID)
	s.EqualValues(expected["type"], actual.Type)
	s.EqualValues(expected["application_id"], *actual.ApplicationID)
	s.EqualValues(expected["flags"], actual.Flags)
	s.EqualValues(expected["position"], *actual.Position)

	interaction, ok := expected["interaction"].(map[string]interface{})
	if !ok {
		s.Failf("Conversion Error", "expected interaction to be map[string]interface{}, got %T with data %+v", expected["interaction"], expected["interaction"])
		return
	}

	s.EqualValues(interaction["id"], actual.Interaction.ID)
	s.EqualValues(interaction["type"], actual.Interaction.Type)
	s.EqualValues(interaction["name"], actual.Interaction.Name)
	s.compareUser(interaction["user"].(map[string]interface{}), actual.Interaction.User)
	s.compareMember(interaction["member"].(map[string]interface{}), *actual.Interaction.Member)

	messageReference := expected["message_reference"].(map[string]interface{})
	s.EqualValues(messageReference["type"], *actual.MessageReference.Type)
	s.EqualValues(messageReference["message_id"], *actual.MessageReference.MessageID)
	s.EqualValues(messageReference["channel_id"], *actual.MessageReference.ChannelID)
	s.EqualValues(messageReference["guild_id"], *actual.MessageReference.GuildID)
	s.EqualValues(messageReference["fail_if_not_exists"], *actual.MessageReference.FailIfNotExists)

	activity := expected["activity"].(map[string]interface{})
	s.EqualValues(activity["type"], actual.Activity.Type)
	s.EqualValues(activity["party_id"], actual.Activity.PartyID)

	roleSubscriptionData := expected["role_subscription_data"].(map[string]interface{})
	s.EqualValues(roleSubscriptionData["role_subscription_listing_id"], actual.RoleSubscriptionData.RoleSubscriptionListingID)
	s.EqualValues(roleSubscriptionData["tier_name"], actual.RoleSubscriptionData.TierName)
	s.EqualValues(roleSubscriptionData["total_months_subscribed"], actual.RoleSubscriptionData.TotalMonthsSubscribed)
	s.EqualValues(roleSubscriptionData["is_renewal"], actual.RoleSubscriptionData.IsRenewal)

	poll := expected["poll"].(map[string]interface{})

	question := poll["question"].(map[string]interface{})
	s.EqualValues(question["text"], actual.Poll.Question.Text)

	results := poll["results"].(map[string]interface{})
	s.EqualValues(results["is_finalized"], actual.Poll.Results.IsFinalized)

	pollResultAnswers := results["answer_counts"].([]map[string]interface{})
	s.Len(actual.Poll.Results.AnswerCounts, len(pollResultAnswers))

	for i, answer := range actual.Poll.Results.AnswerCounts {
		s.EqualValues(pollResultAnswers[i]["id"], answer.ID)
		s.EqualValues(pollResultAnswers[i]["count"], answer.Count)
		s.EqualValues(pollResultAnswers[i]["me_voted"], answer.MeVoted)
	}

	pollAnswers := poll["answers"].([]map[string]interface{})
	s.Len(actual.Poll.Answers, len(pollAnswers))

	for i, answer := range actual.Poll.Answers {
		s.EqualValues(pollAnswers[i]["answer_id"], *answer.AnswerID)

		media := pollAnswers[i]["poll_media"].(map[string]interface{})
		s.EqualValues(media["text"], answer.PollMedia.Text)

		val, ok := media["emoji"]
		if ok {
			s.compareEmoji(val.(map[string]interface{}), *answer.PollMedia.Emoji)
		}
	}

	s.EqualValues(poll["expiry"], *actual.Poll.Expiry)
	s.EqualValues(poll["allow_multiselect"], actual.Poll.AllowMultiselect)
	s.EqualValues(poll["layout_type"], actual.Poll.LayoutType)

	call := expected["call"].(map[string]interface{})
	s.EqualValues(call["ended_timestamp"], *actual.Call.EndedTimestamp)
	s.EqualValues(call["participants"], actual.Call.Participants)

	sharedClientTheme := expected["shared_client_theme"].(map[string]interface{})
	s.EqualValues(sharedClientTheme["colors"], actual.SharedClientTheme.Colors)
	s.EqualValues(sharedClientTheme["gradient_angle"], actual.SharedClientTheme.GradientAngle)
	s.EqualValues(sharedClientTheme["base_mix"], actual.SharedClientTheme.BaseMix)
	s.EqualValues(sharedClientTheme["base_theme"], actual.SharedClientTheme.BaseTheme)

	messageSnapshots, ok := expected["message_snapshots"].([]map[string]interface{})
	if !ok && expected["message_snapshots"] != nil {
		s.Failf("Conversion Error", "expected message_snapshots to be []map[string]interface{}, got %T with data %+v", expected["message_snapshots"], expected["message_snapshots"])
		return
	}

	s.Len(actual.MessageSnapshots, len(messageSnapshots))

	for i, snapshot := range actual.MessageSnapshots {
		s.compareMessage(messageSnapshots[i]["message"].(map[string]interface{}), snapshot.Message)
	}

	mentions := expected["mentions"].([]map[string]interface{})
	s.Len(actual.Mentions, len(mentions))

	for i, mention := range actual.Mentions {
		s.compareUser(mentions[i], mention)
	}

	channelMentions := expected["mention_channels"].([]map[string]interface{})
	s.Len(actual.MentionChannels, len(channelMentions))

	for i, mention := range actual.MentionChannels {
		s.EqualValues(channelMentions[i]["id"], mention.ID)
		s.EqualValues(channelMentions[i]["type"], mention.Type)
		s.EqualValues(channelMentions[i]["guild_id"], mention.GuildID)
		s.EqualValues(channelMentions[i]["name"], mention.Name)
	}

	attachments := expected["attachments"].([]map[string]interface{})
	s.Len(actual.Attachments, len(attachments))

	for i, attachment := range actual.Attachments {
		s.compareAttachment(attachments[i], attachment)
	}

	embeds := expected["embeds"].([]map[string]interface{})
	s.Len(actual.Embeds, len(embeds))

	for i, embed := range actual.Embeds {
		s.compareEmbed(embeds[i], embed)
	}

	reactions := expected["reactions"].([]map[string]interface{})
	s.Len(actual.Reactions, len(reactions))

	for i, reaction := range actual.Reactions {
		s.compareReaction(reactions[i], reaction)
	}

	messageComponents := expected["components"].([]map[string]interface{})
	s.Len(actual.Components, len(messageComponents))

	for i, component := range actual.Components {
		s.compareComponent(messageComponents[i], component)
	}

	stickerItems := expected["sticker_items"].([]map[string]interface{})
	s.Len(actual.StickerItems, len(stickerItems))

	for i, sticker := range actual.StickerItems {
		s.EqualValues(stickerItems[i]["name"], sticker.Name)
		s.EqualValues(stickerItems[i]["id"], sticker.ID)
		s.EqualValues(stickerItems[i]["format_type"], sticker.FormatType)
	}

	stickers := expected["stickers"].([]map[string]interface{})
	s.Len(actual.Stickers, len(stickers))

	for i, sticker := range actual.Stickers {
		s.compareSticker(stickers[i], sticker)
	}

	s.compareUser(expected["author"].(map[string]interface{}), actual.Author)
	s.compareApplication(expected["application"].(map[string]interface{}), *actual.Application)

	referencedMessage, ok := expected["referenced_message"].(map[string]interface{})
	if ok {
		s.compareMessage(referencedMessage, *actual.ReferencedMessage)
	}
	s.compareChannel(expected["thread"].(map[string]interface{}), *actual.Thread)

	resolved, ok := expected["resolved"].(map[string]interface{})
	if ok {
		s.compareResolved(resolved, *actual.Resolved)
	}

	s.compareMessageInteractionMetadata(expected["interaction_metadata"].(map[string]interface{}), actual.InteractionMetadata)
}

func (s *functionSuite) compareMessageInteractionMetadata(expected map[string]interface{}, actual discord.AnyMessageInteractionMetadata) {
	switch actual.(type) {
	case discord.MessageInteractionMetadataApplicationCommand:
		actualMetadata := actual.(discord.MessageInteractionMetadataApplicationCommand)
		s.EqualValues(expected["id"], actualMetadata.ID)
		s.EqualValues(expected["type"], actualMetadata.Type)
		s.EqualValues(expected["authorizing_integration_owners"], actualMetadata.AuthorizingIntegrationOwners)
		s.EqualValues(expected["original_response_message_id"], actualMetadata.OriginalResponseMessageID)
		s.EqualValues(expected["target_message_id"], actualMetadata.TargetMessageID)

		s.compareUser(expected["target_user"].(map[string]interface{}), *actualMetadata.TargetUser)
		s.compareUser(expected["user"].(map[string]interface{}), actualMetadata.User)
	case discord.MessageInteractionMetadataMessageComponent:
		actualMetadata := actual.(discord.MessageInteractionMetadataMessageComponent)
		s.EqualValues(expected["id"], actualMetadata.ID)
		s.EqualValues(expected["type"], actualMetadata.Type)
		s.EqualValues(expected["authorizing_integration_owners"], actualMetadata.AuthorizingIntegrationOwners)
		s.EqualValues(expected["original_response_message_id"], actualMetadata.OriginalResponseMessageID)
		s.EqualValues(expected["interacted_message_id"], actualMetadata.InteractedMessageID)

		s.compareUser(expected["user"].(map[string]interface{}), actualMetadata.User)
	case discord.MessageInteractionMetadataModalSubmit:
		actualMetadata := actual.(discord.MessageInteractionMetadataModalSubmit)
		s.EqualValues(expected["id"], actualMetadata.ID)
		s.EqualValues(expected["type"], actualMetadata.Type)
		s.EqualValues(expected["authorizing_integration_owners"], actualMetadata.AuthorizingIntegrationOwners)
		s.EqualValues(expected["original_response_message_id"], actualMetadata.OriginalResponseMessageID)

		s.compareMessageInteractionMetadata(expected["triggering_interaction_metadata"].(map[string]interface{}), actualMetadata.TriggeringInteractionMetadata)
		s.compareUser(expected["user"].(map[string]interface{}), actualMetadata.User)
	default:
		s.Error(fmt.Errorf("unknown interaction metadata type of Type %T", actual))
	}
}

func (s *functionSuite) expectMap(expected map[string]interface{}, key string) map[string]interface{} {
	s.T().Helper()
	v, ok := expected[key].(map[string]interface{})
	s.Require().True(ok, `expected[%q] must be a map[string]interface{}, got %T: %v`, key, expected[key], expected[key])
	return v
}

func (s *functionSuite) compareInteractionData(expected map[string]interface{}, raw discord.InteractionData) {
	switch raw.GetType() {
	case discord.InteractionTypeApplicationCommand:
		actual := raw.(*responses.InteractionDataApplicationCommand)

		s.EqualValues(expected["id"], actual.ID)
		s.EqualValues(expected["name"], actual.Name)
		s.EqualValues(expected["guild_id"], *actual.GuildID)
		s.EqualValues(expected["target_id"], *actual.TargetID)
		s.EqualValues(expected["type"], actual.Type)

		s.compareResolved(s.expectMap(expected, "resolved"), *actual.Resolved)

		options := expected["options"].([]map[string]interface{})
		s.Len(actual.Options, len(options))

		for i, option := range actual.Options {
			s.compareApplicationCommandInteractionDataOption(options[i], option)
		}
	case discord.InteractionTypeMessageComponent:
		actual := raw.(*responses.InteractionDataMessageComponent)

		s.EqualValues(expected["custom_id"], actual.CustomID)
		s.EqualValues(expected["component_type"], actual.ComponentType)

		values := expected["values"].([]map[string]interface{})
		s.Len(actual.Values, len(values))

		for i, value := range actual.Values {
			s.Equal(values[i]["value"], value.Value)
			s.Equal(values[i]["label"], value.Label)
			s.Equal(values[i]["description"], value.Description)
			s.Equal(values[i]["default"], value.Default)

			emoji := values[i]["emoji"].(map[string]interface{})
			s.EqualValues(emoji["id"], value.Emoji.ID)
			s.EqualValues(emoji["name"], value.Emoji.Name)
			s.EqualValues(emoji["animated"], *value.Emoji.Animated)
		}

		s.compareResolved(expected["resolved"].(map[string]interface{}), *actual.Resolved)

	case discord.InteractionTypeApplicationCommandAutocomplete:
		actual := raw.(*responses.InteractionDataAutocomplete)

		s.EqualValues(expected["id"], actual.ID)
		s.EqualValues(expected["name"], actual.Name)
		s.EqualValues(expected["type"], actual.Type)
		s.EqualValues(expected["target_id"], *actual.TargetID)
		s.EqualValues(expected["guild_id"], *actual.GuildID)

		s.compareResolved(expected["resolved"].(map[string]interface{}), *actual.Resolved)

		options := expected["options"].([]map[string]interface{})
		s.Len(actual.Options, len(options))

		for i, option := range actual.Options {
			s.compareApplicationCommandInteractionDataOption(options[i], option)
		}

	case discord.InteractionTypeModalSubmit:
		actual := raw.(*responses.InteractionDataModalSubmit)

		s.EqualValues(expected["custom_id"], actual.CustomID)

		cmp, ok := expected["components"].([]map[string]interface{})
		s.Require().True(ok, `expected["cmp"] must be a []map[string]interface{} (e.g. from testdata.NewModalSubmitData()/NewModalSubmitDataWithComponents()), got %T: %v`, expected["cmp"], expected["cmp"])
		s.Len(actual.Components, len(cmp))

		for i, component := range actual.Components {
			s.EqualValues(cmp[i]["type"], component.Type)
			s.EqualValues(cmp[i]["id"], *component.ID)
			s.EqualValues(cmp[i]["label"], component.Label)
			s.EqualValues(cmp[i]["description"], component.Description)

			child, ok := cmp[i]["component"].(map[string]interface{})
			s.Require().True(ok, `expected cmp[%d]["component"] must be a map[string]interface{}, got %T: %v`, i, cmp[i]["component"], cmp[i]["component"])
			s.compareComponentLabelChild(child, component.Component)
		}

		s.compareResolved(expected["resolved"].(map[string]interface{}), *actual.Resolved)
	}
}

func (s *functionSuite) compareApplicationCommandInteractionDataOption(expected map[string]interface{}, actual responses.ApplicationCommandInteractionDataOption[interface{}]) {
	s.EqualValues(expected["name"], actual.Name)
	s.EqualValues(expected["type"], actual.Type)

	if actual.Type == discord.ApplicationCommandOptionTypeSubCommand || actual.Type == discord.ApplicationCommandOptionTypeSubCommandGroup {
		options := expected["options"].([]map[string]interface{})
		s.Len(actual.Options, len(options))

		if len(actual.Options) != 0 {
			for i, option := range actual.Options {
				s.compareApplicationCommandInteractionDataOption(options[i], option)
			}
		}
	} else {
		s.EqualValues(expected["value"], *actual.Value)
		s.EqualValues(expected["focused"], *actual.Focused)
	}
}

func (s *functionSuite) compareComponentLabelChild(expected map[string]interface{}, actual components.AnyComponentInteractionResponse) {
	switch v := actual.(type) {
	case *components.StringSelectComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["component_type"], v.ComponentType)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)

	case *components.UserSelectComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["component_type"], v.ComponentType)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)
		s.compareResolved(s.expectMap(expected, "resolved"), v.Resolved)

	case *components.RoleComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["component_type"], v.ComponentType)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)
		s.compareResolved(s.expectMap(expected, "resolved"), *v.Resolved)

	case *components.MentionableComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["component_type"], v.ComponentType)
		s.EqualValues(expected["id"], v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)
		s.compareResolved(s.expectMap(expected, "resolved"), v.Resolved)

	case *components.ChannelComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["component_type"], v.ComponentType)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)
		s.compareResolved(s.expectMap(expected, "resolved"), *v.Resolved)

	case *components.TextInputComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["value"], v.Value)

	case *components.FileUploadComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)

	case *components.RadioGroupComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["value"], *v.Value)

	case *components.CheckboxGroupComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["values"], v.Values)

	case *components.CheckboxComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)
		s.EqualValues(expected["custom_id"], v.CustomID)
		s.EqualValues(expected["value"], v.Value)

	case *components.TextDisplayComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)

	case *components.LabelComponentInteractionResponse:
		s.EqualValues(expected["type"], v.Type)
		s.EqualValues(expected["id"], *v.ID)
		s.compareComponentLabelChild(expected["component"].(map[string]interface{}), v.Component)

	default:
		s.Failf("unhandled component response type in compareComponentLabelChild", "%T", actual)
	}
}

func (s *functionSuite) compareReaction(expected map[string]interface{}, actual discord.Reaction) {
	s.EqualValues(expected["count"], actual.Count)
	s.EqualValues(expected["me"], actual.Me)
	s.EqualValues(expected["me_burst"], actual.MeBurst)
	s.EqualValues(expected["burst_colors"], actual.BurstColors)

	countDetails := expected["count_details"].(map[string]interface{})
	s.EqualValues(countDetails["normal"], actual.CountDetails.Normal)
	s.EqualValues(countDetails["burst"], actual.CountDetails.Burst)

	s.compareEmoji(expected["emoji"].(map[string]interface{}), actual.Emoji)
}

func (s *functionSuite) compareEmbed(expected map[string]interface{}, actual discord.Embed) {
	s.EqualValues(expected["title"], actual.Title)
	s.EqualValues(expected["type"], actual.Type)
	s.EqualValues(expected["description"], actual.Description)
	s.EqualValues(expected["url"], actual.URL)
	s.EqualValues(expected["timestamp"], *actual.Timestamp)
	s.EqualValues(expected["color"], *actual.Color)
	s.EqualValues(expected["flags"], actual.Flags)

	footer := expected["footer"].(map[string]interface{})
	s.EqualValues(footer["text"], actual.Footer.Text)
	s.EqualValues(footer["icon_url"], actual.Footer.IconURL)
	s.EqualValues(footer["proxy_icon_url"], actual.Footer.ProxyIconURL)

	author := expected["author"].(map[string]interface{})
	s.EqualValues(author["name"], actual.Author.Name)
	s.EqualValues(author["icon_url"], actual.Author.IconURL)
	s.EqualValues(author["proxy_icon_url"], actual.Author.ProxyIconURL)
	s.EqualValues(author["url"], actual.Author.URL)

	image := expected["image"].(map[string]interface{})
	s.EqualValues(image["url"], actual.Image.URL)
	s.EqualValues(image["proxy_url"], actual.Image.ProxyURL)
	s.EqualValues(image["height"], *actual.Image.Height)
	s.EqualValues(image["width"], *actual.Image.Width)
	s.EqualValues(image["content_type"], actual.Image.ContentType)
	s.EqualValues(image["placeholder"], actual.Image.Placeholder)
	s.EqualValues(image["placeholder_version"], *actual.Image.PlaceholderVersion)
	s.EqualValues(image["description"], actual.Image.Description)
	s.EqualValues(image["flags"], actual.Image.Flags)

	thumbnail := expected["thumbnail"].(map[string]interface{})
	s.EqualValues(thumbnail["url"], actual.Thumbnail.URL)
	s.EqualValues(thumbnail["proxy_url"], actual.Thumbnail.ProxyURL)
	s.EqualValues(thumbnail["height"], *actual.Thumbnail.Height)
	s.EqualValues(thumbnail["width"], *actual.Thumbnail.Width)
	s.EqualValues(thumbnail["content_type"], actual.Thumbnail.ContentType)
	s.EqualValues(thumbnail["placeholder"], actual.Thumbnail.Placeholder)
	s.EqualValues(thumbnail["placeholder_version"], *actual.Thumbnail.PlaceholderVersion)
	s.EqualValues(thumbnail["description"], actual.Thumbnail.Description)
	s.EqualValues(thumbnail["flags"], actual.Thumbnail.Flags)

	video := expected["video"].(map[string]interface{})
	s.EqualValues(video["url"], actual.Video.URL)
	s.EqualValues(video["proxy_url"], actual.Video.ProxyURL)
	s.EqualValues(video["height"], *actual.Video.Height)
	s.EqualValues(video["width"], *actual.Video.Width)
	s.EqualValues(video["content_type"], actual.Video.ContentType)
	s.EqualValues(video["placeholder"], actual.Video.Placeholder)
	s.EqualValues(video["placeholder_version"], *actual.Video.PlaceholderVersion)
	s.EqualValues(video["description"], actual.Video.Description)
	s.EqualValues(video["flags"], actual.Video.Flags)

	provider := expected["provider"].(map[string]interface{})
	s.EqualValues(provider["name"], actual.Provider.Name)
	s.EqualValues(provider["url"], actual.Provider.URL)

	fields := expected["fields"].([]map[string]interface{})
	s.Len(actual.Fields, len(fields))

	for i, field := range actual.Fields {
		s.EqualValues(fields[i]["name"], field.Name)
		s.EqualValues(fields[i]["value"], field.Value)
		s.EqualValues(fields[i]["inline"], field.Inline)
	}
}

func mustUnmarshalInto[V any](t *testing.T, data []byte) V {
	t.Helper()

	var marshalInto V
	if err := json.Unmarshal(data, &marshalInto); err != nil {
		t.Errorf("got error %+v with data %+v", err, data)
	}

	return marshalInto
}

func mustConvert[V any](t *testing.T, m map[string]interface{}) V {
	t.Helper()

	data, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got error %+v marshaling %+v", err, m)
	}

	return mustUnmarshalInto[V](t, data)
}

func (s *functionSuite) compareComponent(expected map[string]interface{}, actual discord.AnyComponent) {
	s.T().Helper()

	var componentType discord.ComponentType
	var rawData json.RawMessage

	switch v := actual.(type) {
	case *discord.RawComponent:
		componentType = v.GetType()
		rawData = v.Data
	default:
		if actual == nil {
			s.Fail("Conversion Error", "expected a component, got nil")
			return
		}
		componentType = actual.GetType()
		data, err := actual.MarshalJSON()
		if err != nil {
			s.Failf("Marshal Error", "failed to marshal concrete component %T back to JSON: %v", actual, err)
			return
		}
		rawData = data
	}

	s.EqualValues(expected["type"], componentType)

	switch actual.GetType() {
	case discord.ComponentTypeActionRow:
		component := mustUnmarshalInto[components.ActionRow](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)

		underLyingComponents, ok := expected["components"].([]map[string]interface{})
		if !ok {
			s.Failf("Conversion Error", "expected components to be []map[string]interface{}, got %T with data %+v", expected["components"], expected["components"])
			return
		}

		s.Len(component.Components, len(underLyingComponents))

		for i, underLyingComponent := range component.Components {
			s.compareComponent(underLyingComponents[i], underLyingComponent)
		}
	case discord.ComponentTypeButton:
		component := mustUnmarshalInto[components.Button](s.T(), rawData)

		s.EqualValues(expected["id"], component.ID)
		s.EqualValues(expected["style"], component.Style)
		s.EqualValues(expected["label"], component.Label)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["sku_id"], component.SkuID)
		s.EqualValues(expected["url"], component.URL)
		s.EqualValues(expected["disabled"], component.Disabled)

		s.compareEmoji(expected["emoji"].(map[string]interface{}), *component.Emoji)
	case discord.ComponentTypeStringSelect:
		component := mustUnmarshalInto[components.StringSelectMenu](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["placeholder"], component.Placeholder)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["required"], component.Required)
		s.EqualValues(expected["disabled"], component.Disabled)

		options := expected["options"].([]map[string]interface{})
		s.Len(component.Options, len(options))

		for i, option := range component.Options {
			s.EqualValues(options[i]["value"], option.Value)
			s.EqualValues(options[i]["label"], option.Label)
			s.EqualValues(options[i]["default"], option.Default)
			s.EqualValues(options[i]["description"], option.Description)
			s.compareEmoji(options[i]["emoji"].(map[string]interface{}), *option.Emoji)
		}

	case discord.ComponentTypeTextInput:
		component := mustUnmarshalInto[components.TextInput](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["style"], component.Style)
		s.EqualValues(expected["min_length"], *component.MinLength)
		s.EqualValues(expected["max_length"], *component.MaxLength)
		s.EqualValues(expected["required"], component.Required)
		s.EqualValues(expected["value"], component.Value)
		s.EqualValues(expected["placeholder"], component.Placeholder)

	case discord.ComponentTypeUserSelect:
		component := mustUnmarshalInto[components.UserSelectMenu](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["placeholder"], component.Placeholder)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["required"], component.Required)
		s.EqualValues(expected["disabled"], component.Disabled)

		s.compareDefaultValues(expected["default_values"].([]map[string]interface{}), component.DefaultValues)
	case discord.ComponentTypeRoleSelect:
		component := mustUnmarshalInto[components.RoleSelectMenu](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["placeholder"], component.Placeholder)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["required"], component.Required)
		s.EqualValues(expected["disabled"], component.Disabled)

		s.compareDefaultValues(expected["default_values"].([]map[string]interface{}), component.DefaultValues)
	case discord.ComponentTypeMentionableSelect:
		component := mustUnmarshalInto[components.MentionableSelectMenu](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["placeholder"], component.Placeholder)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["required"], component.Required)
		s.EqualValues(expected["disabled"], component.Disabled)

		s.compareDefaultValues(expected["default_values"].([]map[string]interface{}), component.DefaultValues)
	case discord.ComponentTypeChannelSelect:
		component := mustUnmarshalInto[components.ChannelSelectMenu](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["placeholder"], component.Placeholder)
		s.EqualValues(expected["channel_types"], component.ChannelTypes)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["required"], component.Required)
		s.EqualValues(expected["disabled"], component.Disabled)

		s.compareDefaultValues(expected["default_values"].([]map[string]interface{}), component.DefaultValues)
	case discord.ComponentTypeSection:
		component := mustUnmarshalInto[components.Section](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)

		underLyingComponents := expected["components"].([]map[string]interface{})
		s.Len(component.Components, len(underLyingComponents))

		for i, underLyingComponent := range component.Components {
			s.compareComponent(underLyingComponents[i], underLyingComponent)
		}

		s.compareComponent(expected["accessory"].(map[string]interface{}), component.Accessory)
	case discord.ComponentTypeTextDisplay:
		component := mustUnmarshalInto[components.TextDisplay](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["content"], component.Content)
	case discord.ComponentTypeThumbnail:
		component := mustUnmarshalInto[components.Thumbnail](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["description"], *component.Description)
		s.EqualValues(expected["spoiler"], component.Spoiler)
		s.compareUnfurledMediaItem(expected["media"].(map[string]interface{}), component.Media)
	case discord.ComponentTypeMediaGallery:
		component := mustUnmarshalInto[components.MediaGallery](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)

		items := expected["items"].([]map[string]interface{})
		s.Len(component.Items, len(items))

		for i, item := range component.Items {
			s.EqualValues(items[i]["description"], *item.Description)
			s.EqualValues(items[i]["spoiler"], item.Spoiler)
			s.compareUnfurledMediaItem(items[i]["media"].(map[string]interface{}), item.Media)
		}
	case discord.ComponentTypeFile:
		component := mustUnmarshalInto[components.File](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["spoiler"], component.Spoiler)
		s.EqualValues(expected["name"], component.Name)
		s.EqualValues(expected["size"], *component.Size)
		s.compareUnfurledMediaItem(expected["file"].(map[string]interface{}), *component.File)
	case discord.ComponentTypeSeparator:
		component := mustUnmarshalInto[components.Separator](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["divider"], component.Divider)
		s.EqualValues(expected["spacing"], component.Spacing)
	case discord.ComponentTypeContainer:
		component := mustUnmarshalInto[components.Container](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["accent_color"], *component.AccentColor)
		s.EqualValues(expected["spoiler"], component.Spoiler)

		underLyingComponents := expected["components"].([]map[string]interface{})
		s.Len(component.Components, len(underLyingComponents))

		for i, underLyingComponent := range component.Components {
			s.compareComponent(underLyingComponents[i], underLyingComponent)
		}
	case discord.ComponentTypeLabel:
		component := mustUnmarshalInto[components.Label](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["label"], component.Label)
		s.EqualValues(expected["description"], component.Description)
		s.Require().NotNil(component.Component)

		underlyingComponent, ok := expected["component"].(map[string]interface{})
		if !ok {
			s.Failf("Conversion Error", "expected component to be map[string]interface{}, got %T with data %+v", expected["component"], expected["component"])
			return
		}

		s.compareComponent(underlyingComponent, component.Component)
	case discord.ComponentTypeFileUpload:
		component := mustUnmarshalInto[components.FileUpload](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["required"], component.Required)
	case discord.ComponentTypeRadioGroup:
		component := mustUnmarshalInto[components.RadioGroup](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["required"], component.Required)

		options := expected["options"].([]map[string]interface{})
		s.Len(component.Options, len(options))

		for i, option := range component.Options {
			s.EqualValues(options[i]["label"], option.Label)
			s.EqualValues(options[i]["value"], option.Value)
			s.EqualValues(options[i]["description"], option.Description)
			s.EqualValues(options[i]["default"], option.Default)
		}
	case discord.ComponentTypeCheckboxGroup:
		component := mustUnmarshalInto[components.CheckboxGroup](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["min_values"], *component.MinValues)
		s.EqualValues(expected["max_values"], *component.MaxValues)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["required"], component.Required)

		options, ok := expected["options"].([]map[string]interface{})
		if !ok {
			s.Failf("Conversion Error", "expected options to be map[string]interface{}, got %T with data %+v", expected["options"], expected["options"])
			return
		}
		s.Len(component.Options, len(options))

		for i, option := range component.Options {
			s.EqualValues(options[i]["label"], option.Label)
			s.EqualValues(options[i]["value"], option.Value)
			s.EqualValues(options[i]["description"], option.Description)
			s.EqualValues(options[i]["default"], option.Default)
		}
	case discord.ComponentTypeCheckbox:
		component := mustUnmarshalInto[components.Checkbox](s.T(), rawData)

		s.EqualValues(expected["id"], *component.ID)
		s.EqualValues(expected["custom_id"], component.CustomID)
		s.EqualValues(expected["default"], component.Default)
	}
}

func (s *functionSuite) compareUnfurledMediaItem(expected map[string]interface{}, actual components.UnfurledMediaItem) {
	s.EqualValues(expected["url"], actual.URL)
	s.EqualValues(expected["proxy_url"], actual.ProxyURL)
	s.EqualValues(expected["height"], *actual.Height)
	s.EqualValues(expected["width"], *actual.Width)
	s.EqualValues(expected["content_type"], actual.ContentType)
	s.EqualValues(expected["placeholder"], actual.Placeholder)
	s.EqualValues(expected["placeholder_version"], actual.PlaceholderVersion)
	s.EqualValues(expected["flags"], actual.Flags)
	s.EqualValues(expected["attachment_id"], *actual.AttachmentID)
}

func (s *functionSuite) compareDefaultValues(expected []map[string]interface{}, actual []components.SelectDefaultValue) {
	s.Len(actual, len(expected))

	for i, option := range actual {
		s.EqualValues(expected[i]["type"], option.Type)
		s.EqualValues(expected[i]["id"], option.ID)
	}
}
