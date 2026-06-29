package testdata

import (
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewAuditLogEntryChange() map[string]interface{} {
	return map[string]interface{}{
		"key":       testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		"new_value": testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		"old_value": testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
	}
}

type AuditLogEntryPayload struct {
	TargetID   *discord.Snowflake            `json:"target_id"`
	Changes    []discord.AuditLogEntryChange `json:"changes"`
	UserID     *discord.Snowflake            `json:"user_id"`
	ID         discord.Snowflake             `json:"id"`
	ActionType discord.AuditLogActionType    `json:"action_type"`
	Options    *discord.AuditLogEntryOptions `json:"options,omitempty"`
	Reason     string                        `json:"reason,omitempty"`
	GuildID    discord.Snowflake             `json:"guild_id"`
}

func NewAuditLogEntry() map[string]interface{} {
	return map[string]interface{}{
		"guild_id":  discord.RandomSnowflake(),
		"target_id": discord.RandomSnowflake(),
		"changes": testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 32), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewAuditLogEntryChange())
		}),
		"user_id": discord.RandomSnowflake(),
		"id":      discord.RandomSnowflake(),
		"action_type": testutil.RandomItem(
			discord.AuditLogActionTypeGuildUpdate,
			discord.AuditLogActionTypeChannelCreate,
			discord.AuditLogActionTypeChannelUpdate,
			discord.AuditLogActionTypeChannelDelete,
			discord.AuditLogActionTypeChannelOverwriteCreate,
			discord.AuditLogActionTypeChannelOverwriteUpdate,
			discord.AuditLogActionTypeChannelOverwriteDelete,
			discord.AuditLogActionTypeMemberKick,
			discord.AuditLogActionTypeMemberPrune,
			discord.AuditLogActionTypeMemberBanAdd,
			discord.AuditLogActionTypeMemberBanRemove,
			discord.AuditLogActionTypeMemberUpdate,
			discord.AuditLogActionTypeMemberRoleUpdate,
			discord.AuditLogActionTypeMemberMove,
			discord.AuditLogActionTypeMemberDisconnect,
			discord.AuditLogActionTypeMemberBotAdd,
			discord.AuditLogActionTypeRoleCreate,
			discord.AuditLogActionTypeRoleUpdate,
			discord.AuditLogActionTypeRoleDelete,
			discord.AuditLogActionTypeInviteCreate,
			discord.AuditLogActionTypeInviteUpdate,
			discord.AuditLogActionTypeInviteDelete,
			discord.AuditLogActionTypeWebhookCreate,
			discord.AuditLogActionTypeWebhookUpdate,
			discord.AuditLogActionTypeWebhookDelete,
			discord.AuditLogActionTypeEmojiCreate,
			discord.AuditLogActionTypeEmojiUpdate,
			discord.AuditLogActionTypeEmojiDelete,
			discord.AuditLogActionTypeMessageDelete,
			discord.AuditLogActionTypeMessageBulkDelete,
			discord.AuditLogActionTypeMessagePin,
			discord.AuditLogActionTypeMessageUnpin,
			discord.AuditLogActionTypeIntegrationCreate,
			discord.AuditLogActionTypeIntegrationUpdate,
			discord.AuditLogActionTypeIntegrationDelete,
			discord.AuditLogActionTypeStageInstanceCreate,
			discord.AuditLogActionTypeStageInstanceUpdate,
			discord.AuditLogActionTypeStageInstanceDelete,
			discord.AuditLogActionTypeStickerCreate,
			discord.AuditLogActionTypeStickerUpdate,
			discord.AuditLogActionTypeStickerDelete,
			discord.AuditLogActionTypeGuildScheduledEventCreate,
			discord.AuditLogActionTypeGuildScheduledEventUpdate,
			discord.AuditLogActionTypeGuildScheduledEventDelete,
			discord.AuditLogActionTypeThreadCreate,
			discord.AuditLogActionTypeThreadUpdate,
			discord.AuditLogActionTypeThreadDelete,
			discord.AuditLogActionTypeApplicationCommandPermissionUpdate,
			discord.AuditLogActionTypeSoundboardSoundCreate,
			discord.AuditLogActionTypeSoundboardSoundUpdate,
			discord.AuditLogActionTypeSoundboardSoundDelete,
			discord.AuditLogActionTypeAutoModerationRuleCreate,
			discord.AuditLogActionTypeAutoModerationRuleUpdate,
			discord.AuditLogActionTypeAutoModerationRuleDelete,
			discord.AuditLogActionTypeAutoModerationBlockMessage,
			discord.AuditLogActionTypeAutoModerationFlagToChannel,
			discord.AuditLogActionTypeAutoModerationUserCommunicationDisabled,
			discord.AuditLogActionTypeAutoModerationUserQuarantined,
			discord.AuditLogActionTypeCreatorMonetizationRequestCreated,
			discord.AuditLogActionTypeCreatorMonetizationRequestAccepted,
			discord.AuditLogActionTypeOnboardingPromptCreate,
			discord.AuditLogActionTypeOnboardingPromptUpdate,
			discord.AuditLogActionTypeOnboardingPromptDelete,
			discord.AuditLogActionTypeOnboardingCreate,
			discord.AuditLogActionTypeOnboardingUpdate,
			discord.AuditLogActionTypeHomeSettingsCreate,
			discord.AuditLogActionTypeHomeSettingsUpdate,
			discord.AuditLogActionTypeVoiceChannelStatusUpdate,
			discord.AuditLogActionTypeVoiceChannelStatusDelete,
		),
		"options": map[string]interface{}{
			"application_id":            discord.RandomSnowflake(),
			"auto_moderation_rule_name": testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
			"auto_moderation_rule_trigger_type": testutil.RandomItem(
				discord.AutoModerationTriggerTypeKeyword,
				discord.AutoModerationTriggerTypeSpam,
				discord.AutoModerationTriggerTypeKeywordPreset,
				discord.AutoModerationTriggerTypeMentionSpam,
				discord.AutoModerationTriggerTypeMemberProfile,
			),
			"channel_id":         discord.RandomSnowflake(),
			"count":              strconv.Itoa(testutil.RandomNumberInRange(1, 100)),
			"delete_member_days": strconv.Itoa(testutil.RandomNumberInRange(1, 365)),
			"id":                 discord.RandomSnowflake(),
			"members_removed":    strconv.Itoa(testutil.RandomNumberInRange(1, 1000)),
			"message_id":         discord.RandomSnowflake(),
			"role_name":          testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
			"type": testutil.RandomItem(
				discord.PermissionOverwriteTypeRole,
				discord.PermissionOverwriteTypeUser,
			),
			"integration_type": testutil.RandomItem(
				discord.InteractionApplicationIntegrationTypeGuildInstall,
				discord.InteractionApplicationIntegrationTypeUserInstall,
			),
			"status": testutil.RandomString(testutil.RandomNumberInRange(1, 32)),
		},
		"reason": testutil.RandomString(testutil.RandomNumberInRange(1, 512)),
	}
}
