package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewChannel() discord.Channel {
	return discord.Channel{
		ID: discord.RandomSnowflake(),
		Type: testutil.RandomItem(
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
		GuildID:  util.PointerOf(discord.RandomSnowflake()),
		Position: testutil.RandomNumberInRange(0, 500),
		PermissionOverwrites: testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 100), func(arrayToFill *[]discord.ChannelPermissionOverwrite) {
			*arrayToFill = append(*arrayToFill, NewPermissionOverwrite())
		}),
		Name:             util.PointerOf(testutil.RandomString(testutil.RandomNumberInRange(1, 100))),
		Topic:            util.PointerOf(testutil.RandomString(testutil.RandomNumberInRange(1, 4096))),
		NSFW:             util.PointerOf(testutil.RandomBool()),
		LastMessageID:    util.PointerOf(discord.RandomSnowflake()),
		Bitrate:          testutil.RandomNumberInRange(1, 384),
		UserLimit:        testutil.RandomNumberInRange(0, 100),
		RateLimitPerUser: testutil.RandomNumberInRange(1, 21600),
		Recipients: testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(0, 10), func(arrayToFill *[]discord.User) {
			*arrayToFill = append(*arrayToFill, NewUser())
		}),
		Icon:             util.PointerOf(testutil.RandomString(32)),
		OwnerID:          util.PointerOf(discord.RandomSnowflake()),
		ApplicationID:    util.PointerOf(discord.RandomSnowflake()),
		Managed:          testutil.RandomBool(),
		ParentID:         util.PointerOf(discord.RandomSnowflake()),
		LastPinTimestamp: util.PointerOf(testutil.RandomTime()),
		RtcRegion:        util.PointerOf(testutil.RandomString(2)),
		VideoQualityMode: testutil.RandomItem(discord.VideoQualityModeAuto, discord.VideoQualityModeFull),
		MessageCount:     testutil.RandomNumberInRange(1, 100000),
		MemberCount:      testutil.RandomNumberInRange(1, 100000),
		ThreadMetadata: &discord.ThreadMetadata{
			Archived:            testutil.RandomBool(),
			AutoArchiveDuration: testutil.RandomItem(60, 1440, 4320, 10080),
			ArchiveTimestamp:    testutil.RandomTime(),
			CreatedTimestamp:    util.PointerOf(testutil.RandomTime()),
			Locked:              testutil.RandomBool(),
			Invitable:           testutil.RandomBool(),
		},
		Member: &discord.ThreadMember{
			ID:            util.PointerOf(discord.RandomSnowflake()),
			UserID:        util.PointerOf(discord.RandomSnowflake()),
			JoinTimestamp: testutil.RandomTime(),
			Flags:         testutil.RandomNumberInRange(0, 1),
			Member:        util.PointerOf(NewGuildMember()),
		},
		DefaultAutoArchiveDuration: testutil.RandomItem(60, 1440, 4320, 10080),
		Permissions:                testutil.RandomFlags(testutil.AllPermissions...),
		Flags:                      testutil.RandomFlags(discord.ChannelFlagPinned, discord.ChannelFlagRequireTag, discord.ChannelFlagHideMediaDownloadOptions),
		TotalMessageSent:           testutil.RandomNumberInRange(1, 100000),
		AvailableTags: testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 25), func(arrayToFill *[]discord.ChannelTag) {
			*arrayToFill = append(*arrayToFill, NewChannelTag())
		}),
		AppliedTags: testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(1, 25)),
		DefaultReactionEmoji: &discord.DefaultReactionEmoji{
			EmojiID:   util.PointerOf(discord.RandomSnowflake()),
			EmojiName: util.PointerOf(testutil.RandomString(4)),
		},
		DefaultThreadRateLimitPerUser: testutil.RandomNumberInRange(1, 21600),
		DefaultSortOrder:              util.PointerOf(testutil.RandomItem(discord.DefaultSortOrderLatestActivity, discord.DefaultSortOrderCreationDate)),
		DefaultForumLayout:            testutil.RandomItem(discord.ChannelForumLayoutTypeNotSet, discord.ChannelForumLayoutTypeListView, discord.ChannelForumLayoutTypeGalleryView),
	}
}

func NewChannelTag() discord.ChannelTag {
	return discord.ChannelTag{
		ID:        discord.RandomSnowflake(),
		Name:      testutil.RandomString(testutil.RandomNumberInRange(0, 20)),
		Moderated: testutil.RandomBool(),
		EmojiID:   util.PointerOf(discord.RandomSnowflake()),
		EmojiName: util.PointerOf(testutil.RandomString(4)),
	}
}

func NewPermissionOverwrite() discord.ChannelPermissionOverwrite {
	return discord.ChannelPermissionOverwrite{
		ID:    discord.RandomSnowflake(),
		Type:  testutil.RandomItem(discord.PermissionOverwriteTypeRole, discord.PermissionOverwriteTypeUser),
		Allow: testutil.RandomFlags(testutil.AllPermissions...),
		Deny:  testutil.RandomFlags(testutil.AllPermissions...),
	}
}
