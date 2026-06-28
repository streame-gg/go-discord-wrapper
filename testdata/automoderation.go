package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewAutomoderationRule() map[string]interface{} {
	actions := make([]map[string]interface{}, 0, 5)
	for i := 0; i < 5; i++ {
		channelID := discord.RandomSnowflake()
		customMessage := testutil.RandomString(testutil.RandomNumberInRange(0, 150))
		durationSeconds := testutil.RandomNumberInRange(1, 2419200)

		actions = append(actions, map[string]interface{}{
			"type": testutil.RandomItem(
				discord.AutoModerationActionTypeBlockMessage,
				discord.AutoModerationActionTypeSendAlertMessage,
				discord.AutoModerationActionTypeTimeout,
				discord.AutoModerationActionTypeBlockMemberInteraction,
			),
			"metadata": map[string]interface{}{
				"channel_id":       channelID,
				"duration_seconds": durationSeconds,
				"custom_message":   customMessage,
			},
		})
	}

	return map[string]interface{}{
		"id":         discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"name":       testutil.RandomString(testutil.RandomNumberInRange(1, 100)),
		"creator_id": discord.RandomSnowflake(),
		"event_type": testutil.RandomItem(
			discord.AutoModerationEventTypeMessageSend,
			discord.AutoModerationEventTypeMemberUpdate,
		),
		"trigger_type": testutil.RandomItem(
			discord.AutoModerationTriggerTypeKeyword,
			discord.AutoModerationTriggerTypeSpam,
			discord.AutoModerationTriggerTypeKeywordPreset,
			discord.AutoModerationTriggerTypeMentionSpam,
			discord.AutoModerationTriggerTypeMemberProfile,
		),
		"trigger_metadata": map[string]interface{}{
			"keyword_filter": testutil.RandomStringArray(testutil.RandomNumberInRange(0, 1000), 1, 10),
			"regex_patterns": testutil.RandomStringArray(testutil.RandomNumberInRange(0, 10), 1, 10),
			"presets": testutil.RandomArray[discord.KeywordPresetType](
				testutil.RandomNumberInRange(1, 3),
				discord.KeywordPresetTypeProfanity,
				discord.KeywordPresetTypeSexualContent,
				discord.KeywordPresetTypeSlurs,
			),
			"allow_list":                      testutil.RandomStringArray(testutil.RandomNumberInRange(0, 100), 1, 10),
			"mention_total_limit":             testutil.RandomNumberInRange(1, 50),
			"mention_raid_protection_enabled": testutil.RandomBool(),
		},
		"actions":         testutil.RandomArray(testutil.RandomNumberInRange(1, 5), actions...),
		"enabled":         testutil.RandomBool(),
		"exempt_roles":    testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(0, 20)),
		"exempt_channels": testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(0, 20)),
	}
}
