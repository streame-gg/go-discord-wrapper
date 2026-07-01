package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestAutomoderationRuleCreate() {
	s.T().Log("Testing Automoderation Rule Create Unmarshal Logic")

	sub := testutil.InitSub[AutoModerationRuleCreateEvent](s)

	sub.RunCommonEdgeCases()

	randomAutomoderationRule := testdata.NewAutomoderationRule()
	randomAutomoderationRule2 := testdata.NewAutomoderationRule()

	randomAutomoderationRule2LogChannelID := discord.RandomSnowflake()

	testActions := []map[string]interface{}{
		{
			"type": discord.AutoModerationActionTypeBlockMessage,
			"metadata": map[string]interface{}{
				"custom_message": nil,
			},
		},
		{
			"type": discord.AutoModerationActionTypeBlockMessage,
			"metadata": map[string]interface{}{
				"custom_message": "lol!",
			},
		},
		{
			"type": discord.AutoModerationActionTypeSendAlertMessage,
			"metadata": map[string]interface{}{
				"channel_id": randomAutomoderationRule2LogChannelID,
			},
		},
		{
			"type": discord.AutoModerationActionTypeTimeout,
			"metadata": map[string]interface{}{
				"duration_seconds": 3600,
			},
		},
		{
			"type": discord.AutoModerationActionTypeBlockMemberInteraction,
		},
	}

	randomAutomoderationRule2["actions"] = testActions

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(e AutoModerationRuleCreateEvent) {
				s.EqualValues(randomAutomoderationRule["id"], e.ID)
				s.EqualValues(randomAutomoderationRule["guild_id"], e.GuildID)
				s.EqualValues(randomAutomoderationRule["name"], e.Name)
				s.EqualValues(randomAutomoderationRule["creator_id"], e.CreatorID)
				s.EqualValues(randomAutomoderationRule["event_type"], e.EventType)
				s.EqualValues(randomAutomoderationRule["trigger_type"], e.TriggerType)
				s.EqualValues(randomAutomoderationRule["enabled"], e.Enabled)
				s.EqualValues(randomAutomoderationRule["exempt_roles"], e.ExemptRoles)
				s.EqualValues(randomAutomoderationRule["exempt_channels"], e.ExemptChannels)

				triggerMetadata := randomAutomoderationRule["trigger_metadata"].(map[string]interface{})
				s.EqualValues(triggerMetadata["keyword_filter"], e.TriggerMetadata.KeywordFilter)
				s.EqualValues(triggerMetadata["regex_patterns"], e.TriggerMetadata.RegexPatterns)
				s.EqualValues(triggerMetadata["presets"], e.TriggerMetadata.Presets)
				s.EqualValues(triggerMetadata["allow_list"], e.TriggerMetadata.AllowList)
				s.EqualValues(triggerMetadata["mention_total_limit"], *e.TriggerMetadata.MentionTotalLimit)
				s.EqualValues(triggerMetadata["mention_raid_protection_enabled"], e.TriggerMetadata.MentionRaidProtectionEnabled)

				actions := randomAutomoderationRule["actions"].([]map[string]interface{})
				for i, action := range actions {
					s.EqualValues(action["type"], e.Actions[i].Type)
					if action["metadata"] == nil {
						s.Nil(e.Actions[i].Metadata)
					} else {
						metadata := action["metadata"].(map[string]interface{})
						if metadata["channel_id"] != nil {
							s.EqualValues(metadata["channel_id"], *e.Actions[i].Metadata.ChannelID)
						}
						s.EqualValues(metadata["duration_seconds"], e.Actions[i].Metadata.DurationSeconds)
						if metadata["custom_message"] != nil {
							s.EqualValues(metadata["custom_message"], *e.Actions[i].Metadata.CustomMessage)
						} else {
							s.Nil(e.Actions[i].Metadata.CustomMessage)
						}
					}
				}
			},
		},
		{
			Name:  "full payload with all actions",
			Input: sub.MustMarshal(randomAutomoderationRule2),
			Validate: func(e AutoModerationRuleCreateEvent) {
				s.EqualValues(randomAutomoderationRule2["id"], e.ID)
				s.EqualValues(randomAutomoderationRule2["guild_id"], e.GuildID)
				s.EqualValues(randomAutomoderationRule2["name"], e.Name)
				s.EqualValues(randomAutomoderationRule2["creator_id"], e.CreatorID)
				s.EqualValues(randomAutomoderationRule2["event_type"], e.EventType)
				s.EqualValues(randomAutomoderationRule2["trigger_type"], e.TriggerType)
				s.EqualValues(randomAutomoderationRule2["enabled"], e.Enabled)
				s.EqualValues(randomAutomoderationRule2["exempt_roles"], e.ExemptRoles)
				s.EqualValues(randomAutomoderationRule2["exempt_channels"], e.ExemptChannels)

				metadata := randomAutomoderationRule2["trigger_metadata"].(map[string]interface{})
				s.EqualValues(metadata["keyword_filter"], e.TriggerMetadata.KeywordFilter)
				s.EqualValues(metadata["regex_patterns"], e.TriggerMetadata.RegexPatterns)
				s.EqualValues(metadata["presets"], e.TriggerMetadata.Presets)
				s.EqualValues(metadata["allow_list"], e.TriggerMetadata.AllowList)
				s.EqualValues(metadata["mention_total_limit"], *e.TriggerMetadata.MentionTotalLimit)
				s.Equal(metadata["mention_raid_protection_enabled"], e.TriggerMetadata.MentionRaidProtectionEnabled)

				s.Len(e.Actions, len(randomAutomoderationRule2["actions"].([]map[string]interface{})))

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[0]["type"])
				s.Nil(e.Actions[0].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[1]["type"])
				s.Equal("lol!", *e.Actions[1].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeSendAlertMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[2]["type"])
				s.EqualValues(randomAutomoderationRule2LogChannelID, *e.Actions[2].Metadata.ChannelID)

				s.Equal(discord.AutoModerationActionTypeTimeout, randomAutomoderationRule2["actions"].([]map[string]interface{})[3]["type"])
				s.Equal(3600, e.Actions[3].Metadata.DurationSeconds)

				s.Equal(discord.AutoModerationActionTypeBlockMemberInteraction, randomAutomoderationRule2["actions"].([]map[string]interface{})[4]["type"])
				s.Nil(e.Actions[4].Metadata)
			},
		},
	})
}

func (s *eventSuite) TestAutomoderationRuleUpdate() {
	s.T().Log("Testing Automoderation Rule Update Unmarshal Logic")

	sub := testutil.InitSub[AutoModerationRuleUpdateEvent](s)

	sub.RunCommonEdgeCases()

	randomAutomoderationRule := testdata.NewAutomoderationRule()
	randomAutomoderationRule2 := testdata.NewAutomoderationRule()

	randomAutomoderationRule2LogChannelID := discord.RandomSnowflake()

	testActions := []map[string]interface{}{
		{
			"type": discord.AutoModerationActionTypeBlockMessage,
			"metadata": map[string]interface{}{
				"custom_message": nil,
			},
		},
		{
			"type": discord.AutoModerationActionTypeBlockMessage,
			"metadata": map[string]interface{}{
				"custom_message": "lol!",
			},
		},
		{
			"type": discord.AutoModerationActionTypeSendAlertMessage,
			"metadata": map[string]interface{}{
				"channel_id": randomAutomoderationRule2LogChannelID,
			},
		},
		{
			"type": discord.AutoModerationActionTypeTimeout,
			"metadata": map[string]interface{}{
				"duration_seconds": 3600,
			},
		},
		{
			"type": discord.AutoModerationActionTypeBlockMemberInteraction,
		},
	}

	randomAutomoderationRule2["actions"] = testActions

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(e AutoModerationRuleUpdateEvent) {
				s.EqualValues(randomAutomoderationRule["id"], e.NewRule.ID)
				s.EqualValues(randomAutomoderationRule["guild_id"], e.NewRule.GuildID)
				s.EqualValues(randomAutomoderationRule["name"], e.NewRule.Name)
				s.EqualValues(randomAutomoderationRule["creator_id"], e.NewRule.CreatorID)
				s.EqualValues(randomAutomoderationRule["event_type"], e.NewRule.EventType)
				s.EqualValues(randomAutomoderationRule["trigger_type"], e.NewRule.TriggerType)
				s.EqualValues(randomAutomoderationRule["enabled"], e.NewRule.Enabled)
				s.EqualValues(randomAutomoderationRule["exempt_roles"], e.NewRule.ExemptRoles)
				s.EqualValues(randomAutomoderationRule["exempt_channels"], e.NewRule.ExemptChannels)

				triggerMetadata := randomAutomoderationRule["trigger_metadata"].(map[string]interface{})
				s.EqualValues(triggerMetadata["keyword_filter"], e.NewRule.TriggerMetadata.KeywordFilter)
				s.EqualValues(triggerMetadata["regex_patterns"], e.NewRule.TriggerMetadata.RegexPatterns)
				s.EqualValues(triggerMetadata["presets"], e.NewRule.TriggerMetadata.Presets)
				s.EqualValues(triggerMetadata["allow_list"], e.NewRule.TriggerMetadata.AllowList)
				s.EqualValues(triggerMetadata["mention_total_limit"], *e.NewRule.TriggerMetadata.MentionTotalLimit)
				s.EqualValues(triggerMetadata["mention_raid_protection_enabled"], e.NewRule.TriggerMetadata.MentionRaidProtectionEnabled)

				actions := randomAutomoderationRule["actions"].([]map[string]interface{})
				for i, action := range actions {
					s.EqualValues(action["type"], e.NewRule.Actions[i].Type)
					if action["metadata"] == nil {
						s.Nil(e.NewRule.Actions[i].Metadata)
					} else {
						metadata := action["metadata"].(map[string]interface{})
						if metadata["channel_id"] != nil {
							s.EqualValues(metadata["channel_id"], *e.NewRule.Actions[i].Metadata.ChannelID)
						}
						s.EqualValues(metadata["duration_seconds"], e.NewRule.Actions[i].Metadata.DurationSeconds)
						if metadata["custom_message"] != nil {
							s.EqualValues(metadata["custom_message"], *e.NewRule.Actions[i].Metadata.CustomMessage)
						} else {
							s.Nil(e.NewRule.Actions[i].Metadata.CustomMessage)
						}
					}
				}
			},
		},
		{
			Name:  "full payload with all actions",
			Input: sub.MustMarshal(randomAutomoderationRule2),
			Validate: func(e AutoModerationRuleUpdateEvent) {
				s.EqualValues(randomAutomoderationRule2["id"], e.NewRule.ID)
				s.EqualValues(randomAutomoderationRule2["guild_id"], e.NewRule.GuildID)
				s.EqualValues(randomAutomoderationRule2["name"], e.NewRule.Name)
				s.EqualValues(randomAutomoderationRule2["creator_id"], e.NewRule.CreatorID)
				s.EqualValues(randomAutomoderationRule2["event_type"], e.NewRule.EventType)
				s.EqualValues(randomAutomoderationRule2["trigger_type"], e.NewRule.TriggerType)
				s.EqualValues(randomAutomoderationRule2["enabled"], e.NewRule.Enabled)
				s.EqualValues(randomAutomoderationRule2["exempt_roles"], e.NewRule.ExemptRoles)
				s.EqualValues(randomAutomoderationRule2["exempt_channels"], e.NewRule.ExemptChannels)

				metadata := randomAutomoderationRule2["trigger_metadata"].(map[string]interface{})
				s.EqualValues(metadata["keyword_filter"], e.NewRule.TriggerMetadata.KeywordFilter)
				s.EqualValues(metadata["regex_patterns"], e.NewRule.TriggerMetadata.RegexPatterns)
				s.EqualValues(metadata["presets"], e.NewRule.TriggerMetadata.Presets)
				s.EqualValues(metadata["allow_list"], e.NewRule.TriggerMetadata.AllowList)
				s.EqualValues(metadata["mention_total_limit"], *e.NewRule.TriggerMetadata.MentionTotalLimit)
				s.Equal(metadata["mention_raid_protection_enabled"], e.NewRule.TriggerMetadata.MentionRaidProtectionEnabled)

				s.Len(e.NewRule.Actions, len(randomAutomoderationRule2["actions"].([]map[string]interface{})))

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[0]["type"])
				s.Nil(e.NewRule.Actions[0].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[1]["type"])
				s.Equal("lol!", *e.NewRule.Actions[1].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeSendAlertMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[2]["type"])
				s.EqualValues(randomAutomoderationRule2LogChannelID, *e.NewRule.Actions[2].Metadata.ChannelID)

				s.Equal(discord.AutoModerationActionTypeTimeout, randomAutomoderationRule2["actions"].([]map[string]interface{})[3]["type"])
				s.Equal(3600, e.NewRule.Actions[3].Metadata.DurationSeconds)

				s.Equal(discord.AutoModerationActionTypeBlockMemberInteraction, randomAutomoderationRule2["actions"].([]map[string]interface{})[4]["type"])
				s.Nil(e.NewRule.Actions[4].Metadata)
			},
		},
	})
}

func (s *eventSuite) TestAutomoderationRuleDelete() {
	s.T().Log("Testing Automoderation Rule Delete Unmarshal Logic")

	sub := testutil.InitSub[AutoModerationRuleDeleteEvent](s)

	sub.RunCommonEdgeCases()

	randomAutomoderationRule := testdata.NewAutomoderationRule()
	randomAutomoderationRule2 := testdata.NewAutomoderationRule()

	randomAutomoderationRule2LogChannelID := discord.RandomSnowflake()

	testActions := []map[string]interface{}{
		{
			"type": discord.AutoModerationActionTypeBlockMessage,
			"metadata": map[string]interface{}{
				"custom_message": nil,
			},
		},
		{
			"type": discord.AutoModerationActionTypeBlockMessage,
			"metadata": map[string]interface{}{
				"custom_message": "lol!",
			},
		},
		{
			"type": discord.AutoModerationActionTypeSendAlertMessage,
			"metadata": map[string]interface{}{
				"channel_id": randomAutomoderationRule2LogChannelID,
			},
		},
		{
			"type": discord.AutoModerationActionTypeTimeout,
			"metadata": map[string]interface{}{
				"duration_seconds": 3600,
			},
		},
		{
			"type": discord.AutoModerationActionTypeBlockMemberInteraction,
		},
	}

	randomAutomoderationRule2["actions"] = testActions

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(e AutoModerationRuleDeleteEvent) {
				s.EqualValues(randomAutomoderationRule["id"], e.ID)
				s.EqualValues(randomAutomoderationRule["guild_id"], e.GuildID)
				s.EqualValues(randomAutomoderationRule["name"], e.Name)
				s.EqualValues(randomAutomoderationRule["creator_id"], e.CreatorID)
				s.EqualValues(randomAutomoderationRule["event_type"], e.EventType)
				s.EqualValues(randomAutomoderationRule["trigger_type"], e.TriggerType)
				s.EqualValues(randomAutomoderationRule["enabled"], e.Enabled)
				s.EqualValues(randomAutomoderationRule["exempt_roles"], e.ExemptRoles)
				s.EqualValues(randomAutomoderationRule["exempt_channels"], e.ExemptChannels)

				triggerMetadata := randomAutomoderationRule["trigger_metadata"].(map[string]interface{})
				s.EqualValues(triggerMetadata["keyword_filter"], e.TriggerMetadata.KeywordFilter)
				s.EqualValues(triggerMetadata["regex_patterns"], e.TriggerMetadata.RegexPatterns)
				s.EqualValues(triggerMetadata["presets"], e.TriggerMetadata.Presets)
				s.EqualValues(triggerMetadata["allow_list"], e.TriggerMetadata.AllowList)
				s.EqualValues(triggerMetadata["mention_total_limit"], *e.TriggerMetadata.MentionTotalLimit)
				s.EqualValues(triggerMetadata["mention_raid_protection_enabled"], e.TriggerMetadata.MentionRaidProtectionEnabled)

				actions := randomAutomoderationRule["actions"].([]map[string]interface{})
				for i, action := range actions {
					s.EqualValues(action["type"], e.Actions[i].Type)
					if action["metadata"] == nil {
						s.Nil(e.Actions[i].Metadata)
					} else {
						metadata := action["metadata"].(map[string]interface{})
						if metadata["channel_id"] != nil {
							s.EqualValues(metadata["channel_id"], *e.Actions[i].Metadata.ChannelID)
						}
						s.EqualValues(metadata["duration_seconds"], e.Actions[i].Metadata.DurationSeconds)
						if metadata["custom_message"] != nil {
							s.EqualValues(metadata["custom_message"], *e.Actions[i].Metadata.CustomMessage)
						} else {
							s.Nil(e.Actions[i].Metadata.CustomMessage)
						}
					}
				}
			},
		},
		{
			Name:  "full payload with all actions",
			Input: sub.MustMarshal(randomAutomoderationRule2),
			Validate: func(e AutoModerationRuleDeleteEvent) {
				s.EqualValues(randomAutomoderationRule2["id"], e.ID)
				s.EqualValues(randomAutomoderationRule2["guild_id"], e.GuildID)
				s.EqualValues(randomAutomoderationRule2["name"], e.Name)
				s.EqualValues(randomAutomoderationRule2["creator_id"], e.CreatorID)
				s.EqualValues(randomAutomoderationRule2["event_type"], e.EventType)
				s.EqualValues(randomAutomoderationRule2["trigger_type"], e.TriggerType)
				s.EqualValues(randomAutomoderationRule2["enabled"], e.Enabled)
				s.EqualValues(randomAutomoderationRule2["exempt_roles"], e.ExemptRoles)
				s.EqualValues(randomAutomoderationRule2["exempt_channels"], e.ExemptChannels)

				metadata := randomAutomoderationRule2["trigger_metadata"].(map[string]interface{})
				s.EqualValues(metadata["keyword_filter"], e.TriggerMetadata.KeywordFilter)
				s.EqualValues(metadata["regex_patterns"], e.TriggerMetadata.RegexPatterns)
				s.EqualValues(metadata["presets"], e.TriggerMetadata.Presets)
				s.EqualValues(metadata["allow_list"], e.TriggerMetadata.AllowList)
				s.EqualValues(metadata["mention_total_limit"], *e.TriggerMetadata.MentionTotalLimit)
				s.Equal(metadata["mention_raid_protection_enabled"], e.TriggerMetadata.MentionRaidProtectionEnabled)

				s.Len(e.Actions, len(randomAutomoderationRule2["actions"].([]map[string]interface{})))

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[0]["type"])
				s.Nil(e.Actions[0].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[1]["type"])
				s.Equal("lol!", *e.Actions[1].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeSendAlertMessage, randomAutomoderationRule2["actions"].([]map[string]interface{})[2]["type"])
				s.EqualValues(randomAutomoderationRule2LogChannelID, *e.Actions[2].Metadata.ChannelID)

				s.Equal(discord.AutoModerationActionTypeTimeout, randomAutomoderationRule2["actions"].([]map[string]interface{})[3]["type"])
				s.Equal(3600, e.Actions[3].Metadata.DurationSeconds)

				s.Equal(discord.AutoModerationActionTypeBlockMemberInteraction, randomAutomoderationRule2["actions"].([]map[string]interface{})[4]["type"])
				s.Nil(e.Actions[4].Metadata)
			},
		},
	})
}
