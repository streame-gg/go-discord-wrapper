package discord

// AuditLog is the response wrapper returned by the Get Guild Audit Log endpoint.
//
// ApplicationCommands is typed as []any because the ApplicationCommand struct
// lives in types/commands, which imports this package. Unmarshal elements into
// *commands.ApplicationCommand as needed.
type AuditLog struct {
	ApplicationCommands  []any                 `json:"application_commands"`
	AuditLogEntries      []AuditLogEntry       `json:"audit_log_entries"`
	AutoModerationRules  []AutoModerationRule  `json:"auto_moderation_rules"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	Integrations         []Integration         `json:"integrations"`
	Threads              []Channel             `json:"threads"`
	Users                []User                `json:"users"`
	Webhooks             []Webhook             `json:"webhooks"`
}

// AuditLogEntry // https://docs.discord.com/developers/resources/audit-log#audit-log-entry-object
type AuditLogEntry struct {
	TargetID   Snowflake             `json:"target_id"`
	Changes    []AuditLogEntryChange `json:"changes"`
	UserID     *Snowflake            `json:"user_id,omitempty"`
	ID         Snowflake             `json:"id"`
	ActionType AuditLogActionType    `json:"action_type"`
	Options    *AuditLogEntryOptions `json:"options,omitempty"`
	Reason     *string               `json:"reason,omitempty"`
}

type AuditLogEntryChange struct {
	Key      string `json:"key"`
	NewValue any    `json:"new_value,omitempty"`
	OldValue any    `json:"old_value,omitempty"`
}

// AuditLogEntryOptions // https://docs.discord.com/developers/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogEntryOptions struct {
	ApplicationID                 *Snowflake                             `json:"application_id,omitempty"`
	AutoModerationRuleName        *string                                `json:"auto_moderation_rule_name,omitempty"`
	AutoModerationRuleTriggerType *AutoModerationTriggerType             `json:"auto_moderation_rule_trigger_type,omitempty"`
	ChannelID                     *Snowflake                             `json:"channel_id,omitempty"`
	Count                         *int                                   `json:"count,omitempty"`
	DeleteMemberDays              *string                                `json:"delete_member_days,omitempty"`
	ID                            *Snowflake                             `json:"id,omitempty"`
	MembersRemoved                *string                                `json:"members_removed,omitempty"`
	MessageID                     *Snowflake                             `json:"message_id,omitempty"`
	RoleName                      *string                                `json:"role_name,omitempty"`
	Type                          *PermissionOverwriteType               `json:"type,omitempty"`
	IntegrationType               *InteractionApplicationIntegrationType `json:"integration_type,omitempty"`
}

// AuditLogActionType // https://docs.discord.com/developers/resources/audit-log#audit-log-entry-object-audit-log-events
type AuditLogActionType int

const (
	AuditLogActionTypeGuildUpdate            AuditLogActionType = 1
	AuditLogActionTypeChannelCreate          AuditLogActionType = 10
	AuditLogActionTypeChannelUpdate          AuditLogActionType = 11
	AuditLogActionTypeChannelDelete          AuditLogActionType = 12
	AuditLogActionTypeChannelOverwriteCreate AuditLogActionType = 13
	AuditLogActionTypeChannelOverwriteUpdate AuditLogActionType = 14
	AuditLogActionTypeChannelOverwriteDelete AuditLogActionType = 15

	AuditLogActionTypeMemberKick       AuditLogActionType = 20
	AuditLogActionTypeMemberPrune      AuditLogActionType = 21
	AuditLogActionTypeMemberBan        AuditLogActionType = 22
	AuditLogActionTypeMemberUnban      AuditLogActionType = 23
	AuditLogActionTypeMemberUpdate     AuditLogActionType = 24
	AuditLogActionTypeMemberRoleUpdate AuditLogActionType = 25
	AuditLogActionTypeMemberMove       AuditLogActionType = 26
	AuditLogActionTypeMemberDisconnect AuditLogActionType = 27
	AuditLogActionTypeMemberBotAdd     AuditLogActionType = 28

	AuditLogActionTypeRoleCreate AuditLogActionType = 30
	AuditLogActionTypeRoleUpdate AuditLogActionType = 31
	AuditLogActionTypeRoleDelete AuditLogActionType = 32

	AuditLogActionTypeInviteCreate AuditLogActionType = 40
	AuditLogActionTypeInviteUpdate AuditLogActionType = 41
	AuditLogActionTypeInviteDelete AuditLogActionType = 42

	AuditLogActionTypeWebhookCreate AuditLogActionType = 50
	AuditLogActionTypeWebhookUpdate AuditLogActionType = 51
	AuditLogActionTypeWebhookDelete AuditLogActionType = 52

	AuditLogActionTypeEmojiCreate AuditLogActionType = 60
	AuditLogActionTypeEmojiUpdate AuditLogActionType = 61
	AuditLogActionTypeEmojiDelete AuditLogActionType = 62

	AuditLogActionTypeMessageDelete     AuditLogActionType = 72
	AuditLogActionTypeMessageBulkDelete AuditLogActionType = 73
	AuditLogActionTypeMessagePin        AuditLogActionType = 74
	AuditLogActionTypeMessageUnpin      AuditLogActionType = 75

	AuditLogActionTypeIntegrationCreate AuditLogActionType = 80
	AuditLogActionTypeIntegrationUpdate AuditLogActionType = 81
	AuditLogActionTypeIntegrationDelete AuditLogActionType = 82

	AuditLogActionTypeStageInstanceCreate AuditLogActionType = 85
	AuditLogActionTypeStageInstanceUpdate AuditLogActionType = 86
	AuditLogActionTypeStageInstanceDelete AuditLogActionType = 87

	AuditLogActionTypeStickerCreate AuditLogActionType = 90
	AuditLogActionTypeStickerUpdate AuditLogActionType = 91
	AuditLogActionTypeStickerDelete AuditLogActionType = 92

	AuditLogActionTypeGuildScheduledEventCreate AuditLogActionType = 100
	AuditLogActionTypeGuildScheduledEventUpdate AuditLogActionType = 101
	AuditLogActionTypeGuildScheduledEventDelete AuditLogActionType = 102

	AuditLogActionTypeThreadCreate AuditLogActionType = 110
	AuditLogActionTypeThreadUpdate AuditLogActionType = 111
	AuditLogActionTypeThreadDelete AuditLogActionType = 112

	AuditLogActionTypeApplicationCommandPermissionUpdate AuditLogActionType = 121

	AuditLogActionTypeSoundboardSoundCreate AuditLogActionType = 130
	AuditLogActionTypeSoundboardSoundUpdate AuditLogActionType = 131
	AuditLogActionTypeSoundboardSoundDelete AuditLogActionType = 132

	AuditLogActionTypeAutoModerationRuleCreate                AuditLogActionType = 140
	AuditLogActionTypeAutoModerationRuleUpdate                AuditLogActionType = 141
	AuditLogActionTypeAutoModerationRuleDelete                AuditLogActionType = 142
	AuditLogActionTypeAutoModerationBlockMessage              AuditLogActionType = 143
	AuditLogActionTypeAutoModerationFlagToChannel             AuditLogActionType = 144
	AuditLogActionTypeAutoModerationUserCommunicationDisabled AuditLogActionType = 145
	AuditLogActionTypeAutoModerationUserQuarantined           AuditLogActionType = 146

	AuditLogActionTypeCreatorMonetizationRequestCreated  AuditLogActionType = 150
	AuditLogActionTypeCreatorMonetizationRequestAccepted AuditLogActionType = 151

	AuditLogActionTypeOnboardingPromptCreate AuditLogActionType = 163
	AuditLogActionTypeOnboardingPromptUpdate AuditLogActionType = 164
	AuditLogActionTypeOnboardingPromptDelete AuditLogActionType = 165
	AuditLogActionTypeOnboardingCreate       AuditLogActionType = 166
	AuditLogActionTypeOnboardingUpdate       AuditLogActionType = 167

	AuditLogActionTypeHomeSettingsCreate AuditLogActionType = 190
	AuditLogActionTypeHomeSettingsUpdate AuditLogActionType = 191
)
