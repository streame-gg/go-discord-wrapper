package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewAutomoderationRule() discord.AutoModerationRule {
	actions := make([]discord.AutoModerationAction, 5)
	for i := 0; i < 5; i++ {
		actions = append(actions, discord.AutoModerationAction{
			Type: testutil.RandomItem(
				discord.AutoModerationActionTypeBlockMessage,
				discord.AutoModerationActionTypeSendAlertMessage,
				discord.AutoModerationActionTypeTimeout,
				discord.AutoModerationActionTypeBlockMemberInteraction,
			),
			Metadata: &discord.AutoModerationActionMetadata{
				ChannelID:       util.PointerOf(discord.RandomSnowflake()),
				DurationSeconds: util.PointerOf(testutil.RandomNumberInRange(1, 2419200)),
				CustomMessage:   util.PointerOf(testutil.RandString(testutil.RandomNumberInRange(0, 150))),
			},
		})
	}

	return discord.AutoModerationRule{
		ID:        discord.RandomSnowflake(),
		GuildID:   discord.RandomSnowflake(),
		Name:      testutil.RandString(testutil.RandomNumberInRange(1, 100)),
		CreatorID: discord.RandomSnowflake(),
		EventType: testutil.RandomItem(discord.AutoModerationEventTypeMessageSend, discord.AutoModerationEventTypeMemberUpdate),
		TriggerType: testutil.RandomItem(
			discord.AutoModerationTriggerTypeKeyword,
			discord.AutoModerationTriggerTypeSpam,
			discord.AutoModerationTriggerTypeKeywordPreset,
			discord.AutoModerationTriggerTypeMentionSpam,
			discord.AutoModerationTriggerTypeMemberProfile,
		),
		TriggerMetadata: discord.AutoModerationTriggerMetadata{
			KeywordFilter: testutil.RandomStringArray(testutil.RandomNumberInRange(0, 1000), 1, 10),
			RegexPatterns: testutil.RandomStringArray(testutil.RandomNumberInRange(0, 10), 1, 10),
			Presets: testutil.RandomArray[discord.KeywordPresetType](
				testutil.RandomNumberInRange(1, 3),
				discord.KeywordPresetTypeProfanity,
				discord.KeywordPresetTypeSexualContent,
				discord.KeywordPresetTypeSlurs,
			),
			AllowList:                    testutil.RandomStringArray(testutil.RandomNumberInRange(0, 100), 1, 10),
			MentionTotalLimit:            util.PointerOf(testutil.RandomNumberInRange(1, 50)),
			MentionRaidProtectionEnabled: testutil.RandomBool(),
		},
		Actions:        testutil.RandomArray(testutil.RandomNumberInRange(1, 5), actions...),
		Enabled:        testutil.RandomBool(),
		ExemptRoles:    testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(0, 20)),
		ExemptChannels: testutil.RandomSnowflakeArray(testutil.RandomNumberInRange(0, 50)),
	}
}
