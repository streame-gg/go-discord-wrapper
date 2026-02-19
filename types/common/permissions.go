package common

type PermissionOverwriteType int

const (
	PermissionOverwriteTypeRole PermissionOverwriteType = 0
	PermissionOverwriteTypeUser PermissionOverwriteType = 1
)

type Permission int64

const (
	PermissionCreateInstantInvite              Permission = 1 << 0
	PermissionKickMembers                      Permission = 1 << 1
	PermissionBanMembers                       Permission = 1 << 2
	PermissionAdministrator                    Permission = 1 << 3
	PermissionManageChannels                   Permission = 1 << 4
	PermissionManageGuild                      Permission = 1 << 5
	PermissionAddReactions                     Permission = 1 << 6
	PermissionViewAuditLog                     Permission = 1 << 7
	PermissionPrioritySpeaker                  Permission = 1 << 8
	PermissionStream                           Permission = 1 << 9
	PermissionViewChannel                      Permission = 1 << 10
	PermissionSendMessages                     Permission = 1 << 11
	PermissionSendTTSMessages                  Permission = 1 << 12
	PermissionManageMessages                   Permission = 1 << 13
	PermissionEmbedLinks                       Permission = 1 << 14
	PermissionAttachFiles                      Permission = 1 << 15
	PermissionReadMessageHistory               Permission = 1 << 16
	PermissionMentionEveryone                  Permission = 1 << 17
	PermissionUseExternalEmojis                Permission = 1 << 18
	PermissionViewGuildInsights                Permission = 1 << 19
	PermissionConnect                          Permission = 1 << 20
	PermissionSpeak                            Permission = 1 << 21
	PermissionMuteMembers                      Permission = 1 << 22
	PermissionDeafenMembers                    Permission = 1 << 23
	PermissionMoveMembers                      Permission = 1 << 24
	PermissionUseVAD                           Permission = 1 << 25
	PermissionChangeNickname                   Permission = 1 << 26
	PermissionManageNicknames                  Permission = 1 << 27
	PermissionManageRoles                      Permission = 1 << 28
	PermissionManageWebhooks                   Permission = 1 << 29
	PermissionManageGuildExpressions           Permission = 1 << 30
	PermissionUseApplicationCommands           Permission = 1 << 31
	PermissionRequestToSpeak                   Permission = 1 << 32
	PermissionManageEvents                     Permission = 1 << 33
	PermissionManageThreads                    Permission = 1 << 34
	PermissionCreatePublicThreads              Permission = 1 << 35
	PermissionCreatePrivateThreads             Permission = 1 << 36
	PermissionUseExternalStickers              Permission = 1 << 37
	PermissionSendMessagesInThreads            Permission = 1 << 38
	PermissionUseEmbeddedActivities            Permission = 1 << 39
	PermissionModerateMembers                  Permission = 1 << 40
	PermissionViewCreatorMonetizationAnalytics Permission = 1 << 41
	PermissionUseSoundboard                    Permission = 1 << 42
	PermissionCreateGuildExpressions           Permission = 1 << 43
	PermissionCreateEvents                     Permission = 1 << 44
	PermissionUseExternalSounds                Permission = 1 << 45
	PermissionSendVoiceMessages                Permission = 1 << 46
	PermissionSendPolls                        Permission = 1 << 49
	PermissionUseExternalApps                  Permission = 1 << 50
	PermissionPinMessages                      Permission = 1 << 51
	PermissionBypassSlowmode                   Permission = 1 << 52
)
