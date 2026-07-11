package testutil

import (
	"math/rand"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var AllPermissions = []discord.Permission{
	discord.PermissionCreateInstantInvite,
	discord.PermissionKickMembers,
	discord.PermissionBanMembers,
	discord.PermissionAdministrator,
	discord.PermissionManageChannels,
	discord.PermissionManageGuild,
	discord.PermissionAddReactions,
	discord.PermissionViewAuditLog,
	discord.PermissionPrioritySpeaker,
	discord.PermissionStream,
	discord.PermissionViewChannel,
	discord.PermissionSendMessages,
	discord.PermissionSendTTSMessages,
	discord.PermissionManageMessages,
	discord.PermissionEmbedLinks,
	discord.PermissionAttachFiles,
	discord.PermissionReadMessageHistory,
	discord.PermissionMentionEveryone,
	discord.PermissionUseExternalEmojis,
	discord.PermissionViewGuildInsights,
	discord.PermissionConnect,
	discord.PermissionSpeak,
	discord.PermissionMuteMembers,
	discord.PermissionDeafenMembers,
	discord.PermissionMoveMembers,
	discord.PermissionUseVAD,
	discord.PermissionChangeNickname,
	discord.PermissionManageNicknames,
	discord.PermissionManageRoles,
	discord.PermissionManageWebhooks,
	discord.PermissionManageGuildExpressions,
	discord.PermissionUseApplicationCommands,
	discord.PermissionRequestToSpeak,
	discord.PermissionManageEvents,
	discord.PermissionManageThreads,
	discord.PermissionCreatePublicThreads,
	discord.PermissionCreatePrivateThreads,
	discord.PermissionUseExternalStickers,
	discord.PermissionSendMessagesInThreads,
	discord.PermissionUseEmbeddedActivities,
	discord.PermissionModerateMembers,
	discord.PermissionViewCreatorMonetizationAnalytics,
	discord.PermissionUseSoundboard,
	discord.PermissionCreateGuildExpressions,
	discord.PermissionCreateEvents,
	discord.PermissionUseExternalSounds,
	discord.PermissionSendVoiceMessages,
	discord.PermissionSetVoiceChannelStatus,
	discord.PermissionSendPolls,
	discord.PermissionUseExternalApps,
	discord.PermissionPinMessages,
	discord.PermissionBypassSlowmode,
}

func RandomFlags[
	V discord.Permission | discord.ChannelFlags | discord.UserFlags |
		discord.GuildMemberFlags | discord.RoleFlags | discord.GuildSystemChannelFlags |
		discord.ActivityFlags | discord.ApplicationFlags | discord.AttachmentFlag | discord.EmbedMediaFlags |
		discord.EmbedFlags | discord.MessageFlag | components.UnfurledMediaItemFlags,
](flags ...V) V {
	var result V
	for _, p := range flags {
		if rand.Intn(2) == 1 {
			result |= p
		}
	}
	return result
}

func RandomFlagsAsString[
	V discord.Permission | discord.ChannelFlags | discord.UserFlags |
		discord.GuildMemberFlags | discord.RoleFlags | discord.GuildSystemChannelFlags |
		discord.ActivityFlags | discord.ApplicationFlags | discord.AttachmentFlag | discord.EmbedMediaFlags |
		discord.EmbedFlags | discord.MessageFlag | components.UnfurledMediaItemFlags,
](flags ...V) string {
	var result V
	for _, p := range flags {
		if rand.Intn(2) == 1 {
			result |= p
		}
	}
	return strconv.Itoa(int(result))
}
