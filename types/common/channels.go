package common

import (
	"time"
)

type ChannelType int

const (
	ChannelTypeGuildText          ChannelType = 0
	ChannelTypeDM                 ChannelType = 1
	ChannelTypeGuildVoice         ChannelType = 2
	ChannelTypeGroupDM            ChannelType = 3
	ChannelTypeGuildCategory      ChannelType = 4
	ChannelTypeGuildAnnouncement  ChannelType = 5
	ChannelTypeAnnouncementThread ChannelType = 10
	ChannelTypePublicThread       ChannelType = 11
	ChannelTypePrivateThread      ChannelType = 12
	ChannelTypeGuildStageVoice    ChannelType = 13
	ChannelTypeGuildDirectory     ChannelType = 14
	ChannelTypeGuildForum         ChannelType = 15
	ChannelTypeGuildMedia         ChannelType = 16
)

type VideoQualityMode int

const (
	VideoQualityModeAuto VideoQualityMode = 1
	VideoQualityModeFull VideoQualityMode = 2
)

type ChannelTag struct {
	ID        Snowflake  `json:"id"`
	Name      string     `json:"name"`
	Moderated bool       `json:"moderated"`
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

type Channel struct {
	ID                            Snowflake                     `json:"id"`
	Type                          ChannelType                   `json:"type"`
	GuildID                       *Snowflake                    `json:"guild_id,omitempty"`
	Position                      *int                          `json:"position,omitempty"`
	PermissionOverwrites          []ChannelPermissionOverwrites `json:"permission_overwrites,omitempty"`
	Name                          string                        `json:"name,omitempty"`
	Topic                         *string                       `json:"topic,omitempty"`
	NSFW                          *bool                         `json:"nsfw,omitempty"`
	LastMessageID                 *Snowflake                    `json:"last_message_id,omitempty"`
	Bitrate                       *int                          `json:"bitrate,omitempty"`
	UserLimit                     *int                          `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                          `json:"rate_limit_per_user,omitempty"`
	Recipients                    *[]User                       `json:"recipients,omitempty"`
	IconHash                      *string                       `json:"icon,omitempty"`
	OwnerID                       *Snowflake                    `json:"owner_id,omitempty"`
	ApplicationID                 *Snowflake                    `json:"application_id,omitempty"`
	ParentID                      *Snowflake                    `json:"parent_id,omitempty"`
	LastPinTimestamp              *time.Time                    `json:"last_pin_timestamp,omitempty"`
	RtcRegion                     *string                       `json:"rtc_region,omitempty"`
	VideoQualityMode              *VideoQualityMode             `json:"video_quality_mode,omitempty"`
	MessageCount                  *int                          `json:"message_count,omitempty"`
	MemberCount                   *int                          `json:"member_count,omitempty"`
	ThreadMetadata                *ThreadMetadata               `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember                 `json:"member,omitempty"`
	DefaultAutoArchiveDuration    *int                          `json:"default_auto_archive_duration,omitempty"`
	Permissions                   *string                       `json:"permissions,omitempty"`
	Flags                         *ChannelFlags                 `json:"flags,omitempty"`
	TotalMessageSent              *int                          `json:"total_message_sent,omitempty"`
	AvailableTags                 *[]ChannelTag                 `json:"available_tags,omitempty"`
	AppliedTags                   *[]Snowflake                  `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          *DefaultReactionEmoji         `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                          `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *DefaultSortOrder             `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *ChannelForumLayoutType       `json:"default_forum_layout,omitempty"`
}

type ChannelPermissionOverwritesType int

const (
	ChannelPermissionOverwritesTypeRole ChannelPermissionOverwritesType = 0
	ChannelPermissionOverwritesTypeUser ChannelPermissionOverwritesType = 1
)

type ChannelPermissionOverwrites struct {
	ID    Snowflake                       `json:"id"`
	Type  ChannelPermissionOverwritesType `json:"type"`
	Allow string                          `json:"allow"`
	Deny  string                          `json:"deny"`
}

type ChannelFlags int

const (
	ChannelFlagPinned                   ChannelFlags = 1 << 1
	ChannelFlagRequireTag               ChannelFlags = 1 << 4
	ChannelFlagHideMediaDownloadOptions ChannelFlags = 1 << 15
)

type DefaultSortOrder int

const (
	DefaultSortOrderLatestActivity DefaultSortOrder = 0
	DefaultSortOrderCreationDate   DefaultSortOrder = 1
)

type ChannelForumLayoutType int

const (
	ChannelForumLayoutTypeNotSet      ChannelForumLayoutType = 0
	ChannelForumLayoutTypeListView    ChannelForumLayoutType = 1
	ChannelForumLayoutTypeGalleryView ChannelForumLayoutType = 2
)

type DefaultReactionEmoji struct {
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

type ThreadMetadata struct {
	Archived            bool       `json:"archived"`
	AutoArchiveDuration *int       `json:"auto_archive_duration,omitempty"`
	ArchiveTimestamp    time.Time  `json:"archive_timestamp"`
	CreatedTimestamp    *time.Time `json:"created_timestamp,omitempty"`
	Locked              *bool      `json:"locked,omitempty"`
	Invitable           *bool      `json:"invitable,omitempty"`
}
