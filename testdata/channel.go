package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewChannel() map[string]interface{} {
	return map[string]interface{}{
		"id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.ChannelTypeGuildText,
			discord.ChannelTypeDM,
			discord.ChannelTypeGuildVoice,
			discord.ChannelTypeGuildCategory,
			discord.ChannelTypeGuildAnnouncement,
			discord.ChannelTypeAnnouncementThread,
			discord.ChannelTypePublicThread,
			discord.ChannelTypePrivateThread,
			discord.ChannelTypeGuildStageVoice,
			discord.ChannelTypeGuildDirectory,
			discord.ChannelTypeGuildForum,
			discord.ChannelTypeGuildMedia,
		),
		"guild_id": discord.RandomSnowflake(),
		"position": testutil.RandomIntInRange(0, 500),
		"permission_overwrites": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 3), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewPermissionOverwrite())
		}),
		"name":                testutil.RandomString(testutil.RandomIntInRange(1, 100)),
		"topic":               testutil.RandomString(testutil.RandomIntInRange(1, 4096)),
		"nsfw":                testutil.RandomBool(),
		"last_message_id":     discord.RandomSnowflake(),
		"bitrate":             testutil.RandomIntInRange(1, 384),
		"user_limit":          testutil.RandomIntInRange(0, 100),
		"rate_limit_per_user": testutil.RandomIntInRange(1, 21600),
		"recipients": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(0, 3), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewUser())
		}),
		"icon":               testutil.RandomString(32),
		"owner_id":           discord.RandomSnowflake(),
		"application_id":     discord.RandomSnowflake(),
		"managed":            testutil.RandomBool(),
		"parent_id":          discord.RandomSnowflake(),
		"last_pin_timestamp": testutil.RandomTime(),
		"rtc_region":         testutil.RandomString(2),
		"video_quality_mode": testutil.RandomItem(discord.VideoQualityModeAuto, discord.VideoQualityModeFull),
		"message_count":      testutil.RandomIntInRange(1, 100000),
		"member_count":       testutil.RandomIntInRange(1, 100000),
		"thread_metadata": map[string]interface{}{
			"archived":              testutil.RandomBool(),
			"auto_archive_duration": testutil.RandomItem(60, 1440, 4320, 10080),
			"archive_timestamp":     testutil.RandomTime(),
			"created_timestamp":     testutil.RandomTime(),
			"locked":                testutil.RandomBool(),
			"invitable":             testutil.RandomBool(),
		},
		"member":                        NewThreadMember(),
		"default_auto_archive_duration": testutil.RandomItem(60, 1440, 4320, 10080),
		"permissions":                   testutil.RandomFlags(testutil.AllPermissions...),
		"flags":                         testutil.RandomFlags(discord.ChannelFlagPinned, discord.ChannelFlagRequireTag, discord.ChannelFlagHideMediaDownloadOptions),
		"total_message_sent":            testutil.RandomIntInRange(1, 100000),
		"available_tags": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 3), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewChannelTag())
		}),
		"applied_tags": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 25)),
		"default_reaction_emoji": map[string]interface{}{
			"emoji_id":   discord.RandomSnowflake(),
			"emoji_name": testutil.RandomString(4),
		},
		"default_thread_rate_limit_per_user": testutil.RandomIntInRange(1, 21600),
		"default_sort_order":                 testutil.RandomItem(discord.DefaultSortOrderLatestActivity, discord.DefaultSortOrderCreationDate),
		"default_forum_layout":               testutil.RandomItem(discord.ChannelForumLayoutTypeNotSet, discord.ChannelForumLayoutTypeListView, discord.ChannelForumLayoutTypeGalleryView),
	}
}

func NewChannelTag() map[string]interface{} {
	return map[string]interface{}{
		"id":         discord.RandomSnowflake(),
		"name":       testutil.RandomString(testutil.RandomIntInRange(0, 20)),
		"moderated":  testutil.RandomBool(),
		"emoji_id":   discord.RandomSnowflake(),
		"emoji_name": testutil.RandomString(4),
	}
}

func NewPermissionOverwrite() map[string]interface{} {
	return map[string]interface{}{
		"id":    discord.RandomSnowflake(),
		"type":  testutil.RandomItem(discord.PermissionOverwriteTypeRole, discord.PermissionOverwriteTypeUser),
		"allow": testutil.RandomFlags(testutil.AllPermissions...),
		"deny":  testutil.RandomFlags(testutil.AllPermissions...),
	}
}
