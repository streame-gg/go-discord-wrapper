package discord

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object-keyword-preset-types
type KeywordPresetType int

const (
	KeywordPresetTypeProfanity     KeywordPresetType = 1
	KeywordPresetTypeSexualContent KeywordPresetType = 2
	KeywordPresetTypeSlurs         KeywordPresetType = 3
)

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object-trigger-types
type AutoModerationTriggerType int

const (
	AutoModerationTriggerTypeKeyword       AutoModerationTriggerType = 1
	AutoModerationTriggerTypeSpam          AutoModerationTriggerType = 3
	AutoModerationTriggerTypeKeywordPreset AutoModerationTriggerType = 4
	AutoModerationTriggerTypeMentionSpam   AutoModerationTriggerType = 5
	AutoModerationTriggerTypeMemberProfile AutoModerationTriggerType = 6
)

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-action-object-action-types
type AutoModerationActionType int

const (
	AutoModerationActionTypeBlockMessage           AutoModerationActionType = 1
	AutoModerationActionTypeSendAlertMessage       AutoModerationActionType = 2
	AutoModerationActionTypeTimeout                AutoModerationActionType = 3
	AutoModerationActionTypeBlockMemberInteraction AutoModerationActionType = 4
)

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-action-object-action-metadata
type AutoModerationActionMetadata struct {
	ChannelID       *Snowflake `json:"channel_id,omitempty"`
	DurationSeconds *int       `json:"duration_seconds,omitempty"`
	CustomMessage   *string    `json:"custom_message,omitempty"`
}

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-action-object
type AutoModerationAction struct {
	Type     AutoModerationActionType      `json:"type"`
	Metadata *AutoModerationActionMetadata `json:"metadata,omitempty"`
}

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object-trigger-metadata
type AutoModerationTriggerMetadata struct {
	KeywordFilter                []string            `json:"keyword_filter,omitempty"`
	RegexPatterns                []string            `json:"regex_patterns,omitempty"`
	Presets                      []KeywordPresetType `json:"presets,omitempty"`
	AllowList                    []string            `json:"allow_list,omitempty"`
	MentionTotalLimit            *int                `json:"mention_total_limit,omitempty"`
	MentionRaidProtectionEnabled bool                `json:"mention_raid_protection_enabled,omitempty"`
}

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object-event-types
type AutoModerationEventType int

const (
	AutoModerationEventTypeMessageSend  AutoModerationEventType = 1
	AutoModerationEventTypeMemberUpdate AutoModerationEventType = 2
)

// https://docs.discord.com/developers/resources/auto-moderation#auto-moderation-rule-object
type AutoModerationRule struct {
	hClient EntityClient

	ID              Snowflake                     `json:"id"`
	GuildID         Snowflake                     `json:"guild_id"`
	Name            string                        `json:"name"`
	CreatorID       Snowflake                     `json:"creator_id"`
	EventType       AutoModerationEventType       `json:"event_type"`
	TriggerType     AutoModerationTriggerType     `json:"trigger_type"`
	TriggerMetadata AutoModerationTriggerMetadata `json:"trigger_metadata"`
	Actions         []AutoModerationAction        `json:"actions"`
	Enabled         bool                          `json:"enabled"`
	ExemptRoles     []Snowflake                   `json:"exempt_roles"`
	ExemptChannels  []Snowflake                   `json:"exempt_channels"`
}
