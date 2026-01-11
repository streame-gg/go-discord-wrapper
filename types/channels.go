package types

import "time"

type DiscordChannelType int

const (
	DiscordChannelTypeGuildText          DiscordChannelType = 0
	DiscordChannelTypeDM                 DiscordChannelType = 1
	DiscordChannelTypeGuildVoice         DiscordChannelType = 2
	DiscordChannelTypeGroupDM            DiscordChannelType = 3
	DiscordChannelTypeGuildCategory      DiscordChannelType = 4
	DiscordChannelTypeGuildAnnouncement  DiscordChannelType = 5
	DiscordChannelTypeAnnouncementThread DiscordChannelType = 10
	DiscordChannelTypePublicThread       DiscordChannelType = 11
	DiscordChannelTypePrivateThread      DiscordChannelType = 12
	DiscordChannelTypeGuildStageVoice    DiscordChannelType = 13
	DiscordChannelTypeGuildDirectory     DiscordChannelType = 14
	DiscordChannelTypeGuildForum         DiscordChannelType = 15
	DiscordChannelTypeGuildMedia         DiscordChannelType = 16
)

type DiscordVideoQualityMode int

const (
	DiscordVideoQualityModeAuto DiscordVideoQualityMode = 1
	DiscordVideoQualityModeFull DiscordVideoQualityMode = 2
)

type DiscordChannelTag struct {
	ID        DiscordSnowflake  `json:"id"`
	Name      string            `json:"name"`
	Moderated bool              `json:"moderated"`
	EmojiID   *DiscordSnowflake `json:"emoji_id,omitempty"`
	EmojiName *string           `json:"emoji_name,omitempty"`
}

type DiscordChannel struct {
	ID                            DiscordSnowflake                     `json:"id"`
	Type                          DiscordChannelType                   `json:"type"`
	GuildID                       *DiscordSnowflake                    `json:"guild_id,omitempty"`
	Position                      *int                                 `json:"position,omitempty"`
	PermissionOverwrites          []DiscordChannelPermissionOverwrites `json:"permission_overwrites,omitempty"`
	Name                          *string                              `json:"name,omitempty"`
	Topic                         *string                              `json:"topic,omitempty"`
	NSFW                          *bool                                `json:"nsfw,omitempty"`
	LastMessageID                 *DiscordSnowflake                    `json:"last_message_id,omitempty"`
	Bitrate                       *int                                 `json:"bitrate,omitempty"`
	UserLimit                     *int                                 `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                                 `json:"rate_limit_per_user,omitempty"`
	Recipients                    *[]DiscordUser                       `json:"recipients,omitempty"`
	IconHash                      *string                              `json:"icon,omitempty"`
	OwnerID                       *DiscordSnowflake                    `json:"owner_id,omitempty"`
	ApplicationID                 *DiscordSnowflake                    `json:"application_id,omitempty"`
	ParentID                      *DiscordSnowflake                    `json:"parent_id,omitempty"`
	LastPinTimestamp              *time.Time                           `json:"last_pin_timestamp,omitempty"`
	RtcRegion                     *string                              `json:"rtc_region,omitempty"`
	VideoQualityMode              *DiscordVideoQualityMode             `json:"video_quality_mode,omitempty"`
	MessageCount                  *int                                 `json:"message_count,omitempty"`
	MemberCount                   *int                                 `json:"member_count,omitempty"`
	ThreadMetadata                *DiscordThreadMetadata               `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember                        `json:"member,omitempty"`
	DefaultAutoArchiveDuration    *int                                 `json:"default_auto_archive_duration,omitempty"`
	Permissions                   *string                              `json:"permissions,omitempty"`
	Flags                         *DiscordChannelFlags                 `json:"flags,omitempty"`
	TotalMessageSent              *int                                 `json:"total_message_sent,omitempty"`
	AvailableTags                 *[]DiscordChannelTag                 `json:"available_tags,omitempty"`
	AppliedTags                   *[]DiscordSnowflake                  `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          *DiscordDefaultReactionEmoji         `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                                 `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *DiscordDefaultSortOrder             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *DiscordChannelForumLayoutType       `json:"default_forum_layout,omitempty"`
}

type DiscordChannelPermissionOverwritesType int

const (
	DiscordChannelPermissionOverwritesTypeRole DiscordChannelPermissionOverwritesType = 0
	DiscordChannelPermissionOverwritesTypeUser DiscordChannelPermissionOverwritesType = 1
)

type DiscordChannelPermissionOverwrites struct {
	ID    DiscordSnowflake                       `json:"id"`
	Type  DiscordChannelPermissionOverwritesType `json:"type"`
	Allow string                                 `json:"allow"`
	Deny  string                                 `json:"deny"`
}

type DiscordChannelFlags int

const (
	DiscordChannelFlagPinned                   DiscordChannelFlags = 1 << 1
	DiscordChannelFlagRequireTag               DiscordChannelFlags = 1 << 4
	DiscordChannelFlagHideMediaDownloadOptions DiscordChannelFlags = 1 << 15
)

type DiscordDefaultSortOrder int

const (
	DiscordDefaultSortOrderLatestActivity DiscordDefaultSortOrder = 0
	DiscordDefaultSortOrderCreationDate   DiscordDefaultSortOrder = 1
)

type DiscordChannelForumLayoutType int

const (
	DiscordChannelForumLayoutTypeNotSet      DiscordChannelForumLayoutType = 0
	DiscordChannelForumLayoutTypeListView    DiscordChannelForumLayoutType = 1
	DiscordChannelForumLayoutTypeGalleryView DiscordChannelForumLayoutType = 2
)

type DiscordThreadMetadata struct {
	Archived            bool       `json:"archived"`
	AutoArchiveDuration *int       `json:"auto_archive_duration,omitempty"`
	ArchiveTimestamp    time.Time  `json:"archive_timestamp"`
	CreateTimestanp     *time.Time `json:"create_timestamp,omitempty"`
	Locked              *bool      `json:"locked,omitempty"`
	Invitable           *bool      `json:"invitable,omitempty"`
}

type DiscordDefaultReactionEmoji struct {
	EmojiID   *DiscordSnowflake `json:"emoji_id,omitempty"`
	EmojiName *string           `json:"emoji_name,omitempty"`
}
