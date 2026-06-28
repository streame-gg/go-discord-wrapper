package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/internal/util"
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

	testActions := []discord.AutoModerationAction{
		{
			Type: discord.AutoModerationActionTypeBlockMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				CustomMessage: nil,
			},
		},
		{
			Type: discord.AutoModerationActionTypeBlockMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				CustomMessage: util.PointerOf("lol!"),
			},
		},
		{
			Type: discord.AutoModerationActionTypeSendAlertMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				ChannelID: util.PointerOf(randomAutomoderationRule2LogChannelID),
			},
		},
		{
			Type: discord.AutoModerationActionTypeTimeout,
			Metadata: &discord.AutoModerationActionMetadata{
				DurationSeconds: 3600,
			},
		},
		{
			Type: discord.AutoModerationActionTypeBlockMemberInteraction,
		},
	}

	randomAutomoderationRule2.Actions = testActions

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(e AutoModerationRuleCreateEvent) {
				s.Equal(randomAutomoderationRule.ID, e.Rule.ID)
				s.Equal(randomAutomoderationRule.GuildID, e.Rule.GuildID)
				s.Equal(randomAutomoderationRule.Name, e.Rule.Name)
				s.Equal(randomAutomoderationRule.CreatorID, e.Rule.CreatorID)
				s.Equal(randomAutomoderationRule.EventType, e.Rule.EventType)
				s.Equal(randomAutomoderationRule.TriggerType, e.Rule.TriggerType)
				s.Equal(randomAutomoderationRule.Enabled, e.Rule.Enabled)
				s.Equal(randomAutomoderationRule.ExemptRoles, e.Rule.ExemptRoles)
				s.Equal(randomAutomoderationRule.ExemptChannels, e.Rule.ExemptChannels)
				s.Equal(randomAutomoderationRule.Actions, e.Rule.Actions)
			},
		},
		{
			Name:  "full payload with all actions",
			Input: sub.MustMarshal(randomAutomoderationRule2),
			Validate: func(e AutoModerationRuleCreateEvent) {
				s.Equal(randomAutomoderationRule.ID, e.Rule.ID)
				s.Equal(randomAutomoderationRule.GuildID, e.Rule.GuildID)
				s.Equal(randomAutomoderationRule.Name, e.Rule.Name)
				s.Equal(randomAutomoderationRule.CreatorID, e.Rule.CreatorID)
				s.Equal(randomAutomoderationRule.EventType, e.Rule.EventType)
				s.Equal(randomAutomoderationRule.TriggerType, e.Rule.TriggerType)
				s.Equal(randomAutomoderationRule.Enabled, e.Rule.Enabled)
				s.Equal(randomAutomoderationRule.ExemptRoles, e.Rule.ExemptRoles)
				s.Equal(randomAutomoderationRule.ExemptChannels, e.Rule.ExemptChannels)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule.Actions[0].Type)
				s.Nil(randomAutomoderationRule.Actions[0].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule.Actions[1].Type)
				s.Equal("lol!", *randomAutomoderationRule.Actions[1].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeSendAlertMessage, randomAutomoderationRule.Actions[2].Type)
				s.Equal(randomAutomoderationRule2LogChannelID, *randomAutomoderationRule.Actions[2].Metadata.ChannelID)

				s.Equal(discord.AutoModerationActionTypeTimeout, randomAutomoderationRule.Actions[3].Type)
				s.Equal(3600, randomAutomoderationRule.Actions[3].Metadata.DurationSeconds)

				s.Equal(discord.AutoModerationActionTypeBlockMemberInteraction, randomAutomoderationRule.Actions[4].Type)
				s.Nil(randomAutomoderationRule.Actions[4].Metadata)
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

	testActions := []discord.AutoModerationAction{
		{
			Type: discord.AutoModerationActionTypeBlockMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				CustomMessage: nil,
			},
		},
		{
			Type: discord.AutoModerationActionTypeBlockMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				CustomMessage: util.PointerOf("lol!"),
			},
		},
		{
			Type: discord.AutoModerationActionTypeSendAlertMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				ChannelID: util.PointerOf(randomAutomoderationRule2LogChannelID),
			},
		},
		{
			Type: discord.AutoModerationActionTypeTimeout,
			Metadata: &discord.AutoModerationActionMetadata{
				DurationSeconds: 3600,
			},
		},
		{
			Type: discord.AutoModerationActionTypeBlockMemberInteraction,
		},
	}

	randomAutomoderationRule2.Actions = testActions

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(e AutoModerationRuleUpdateEvent) {
				s.Equal(randomAutomoderationRule.ID, e.NewRule.ID)
				s.Equal(randomAutomoderationRule.GuildID, e.NewRule.GuildID)
				s.Equal(randomAutomoderationRule.Name, e.NewRule.Name)
				s.Equal(randomAutomoderationRule.CreatorID, e.NewRule.CreatorID)
				s.Equal(randomAutomoderationRule.EventType, e.NewRule.EventType)
				s.Equal(randomAutomoderationRule.TriggerType, e.NewRule.TriggerType)
				s.Equal(randomAutomoderationRule.Enabled, e.NewRule.Enabled)
				s.Equal(randomAutomoderationRule.ExemptRoles, e.NewRule.ExemptRoles)
				s.Equal(randomAutomoderationRule.ExemptChannels, e.NewRule.ExemptChannels)
				s.Equal(randomAutomoderationRule.Actions, e.NewRule.Actions)
				s.Nil(e.OldRule)
			},
		},
		{
			Name:  "full payload with all actions",
			Input: sub.MustMarshal(randomAutomoderationRule2),
			Validate: func(e AutoModerationRuleUpdateEvent) {
				s.Equal(randomAutomoderationRule.ID, e.NewRule.ID)
				s.Equal(randomAutomoderationRule.GuildID, e.NewRule.GuildID)
				s.Equal(randomAutomoderationRule.Name, e.NewRule.Name)
				s.Equal(randomAutomoderationRule.CreatorID, e.NewRule.CreatorID)
				s.Equal(randomAutomoderationRule.EventType, e.NewRule.EventType)
				s.Equal(randomAutomoderationRule.TriggerType, e.NewRule.TriggerType)
				s.Equal(randomAutomoderationRule.Enabled, e.NewRule.Enabled)
				s.Equal(randomAutomoderationRule.ExemptRoles, e.NewRule.ExemptRoles)
				s.Equal(randomAutomoderationRule.ExemptChannels, e.NewRule.ExemptChannels)
				s.Nil(e.OldRule)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule.Actions[0].Type)
				s.Nil(randomAutomoderationRule.Actions[0].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule.Actions[1].Type)
				s.Equal("lol!", *randomAutomoderationRule.Actions[1].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeSendAlertMessage, randomAutomoderationRule.Actions[2].Type)
				s.Equal(randomAutomoderationRule2LogChannelID, *randomAutomoderationRule.Actions[2].Metadata.ChannelID)

				s.Equal(discord.AutoModerationActionTypeTimeout, randomAutomoderationRule.Actions[3].Type)
				s.Equal(3600, randomAutomoderationRule.Actions[3].Metadata.DurationSeconds)

				s.Equal(discord.AutoModerationActionTypeBlockMemberInteraction, randomAutomoderationRule.Actions[4].Type)
				s.Nil(randomAutomoderationRule.Actions[4].Metadata)
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

	testActions := []discord.AutoModerationAction{
		{
			Type: discord.AutoModerationActionTypeBlockMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				CustomMessage: nil,
			},
		},
		{
			Type: discord.AutoModerationActionTypeBlockMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				CustomMessage: util.PointerOf("lol!"),
			},
		},
		{
			Type: discord.AutoModerationActionTypeSendAlertMessage,
			Metadata: &discord.AutoModerationActionMetadata{
				ChannelID: util.PointerOf(randomAutomoderationRule2LogChannelID),
			},
		},
		{
			Type: discord.AutoModerationActionTypeTimeout,
			Metadata: &discord.AutoModerationActionMetadata{
				DurationSeconds: 3600,
			},
		},
		{
			Type: discord.AutoModerationActionTypeBlockMemberInteraction,
		},
	}

	randomAutomoderationRule2.Actions = testActions

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(e AutoModerationRuleDeleteEvent) {
				s.Equal(randomAutomoderationRule.ID, e.Rule.ID)
				s.Equal(randomAutomoderationRule.GuildID, e.Rule.GuildID)
				s.Equal(randomAutomoderationRule.Name, e.Rule.Name)
				s.Equal(randomAutomoderationRule.CreatorID, e.Rule.CreatorID)
				s.Equal(randomAutomoderationRule.EventType, e.Rule.EventType)
				s.Equal(randomAutomoderationRule.TriggerType, e.Rule.TriggerType)
				s.Equal(randomAutomoderationRule.Enabled, e.Rule.Enabled)
				s.Equal(randomAutomoderationRule.ExemptRoles, e.Rule.ExemptRoles)
				s.Equal(randomAutomoderationRule.ExemptChannels, e.Rule.ExemptChannels)
				s.Equal(randomAutomoderationRule.Actions, e.Rule.Actions)
			},
		},
		{
			Name:  "full payload with all actions",
			Input: sub.MustMarshal(randomAutomoderationRule2),
			Validate: func(e AutoModerationRuleDeleteEvent) {
				s.Equal(randomAutomoderationRule.ID, e.Rule.ID)
				s.Equal(randomAutomoderationRule.GuildID, e.Rule.GuildID)
				s.Equal(randomAutomoderationRule.Name, e.Rule.Name)
				s.Equal(randomAutomoderationRule.CreatorID, e.Rule.CreatorID)
				s.Equal(randomAutomoderationRule.EventType, e.Rule.EventType)
				s.Equal(randomAutomoderationRule.TriggerType, e.Rule.TriggerType)
				s.Equal(randomAutomoderationRule.Enabled, e.Rule.Enabled)
				s.Equal(randomAutomoderationRule.ExemptRoles, e.Rule.ExemptRoles)
				s.Equal(randomAutomoderationRule.ExemptChannels, e.Rule.ExemptChannels)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule.Actions[0].Type)
				s.Nil(randomAutomoderationRule.Actions[0].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeBlockMessage, randomAutomoderationRule.Actions[1].Type)
				s.Equal("lol!", *randomAutomoderationRule.Actions[1].Metadata.CustomMessage)

				s.Equal(discord.AutoModerationActionTypeSendAlertMessage, randomAutomoderationRule.Actions[2].Type)
				s.Equal(randomAutomoderationRule2LogChannelID, *randomAutomoderationRule.Actions[2].Metadata.ChannelID)

				s.Equal(discord.AutoModerationActionTypeTimeout, randomAutomoderationRule.Actions[3].Type)
				s.Equal(3600, randomAutomoderationRule.Actions[3].Metadata.DurationSeconds)

				s.Equal(discord.AutoModerationActionTypeBlockMemberInteraction, randomAutomoderationRule.Actions[4].Type)
				s.Nil(randomAutomoderationRule.Actions[4].Metadata)
			},
		},
	})
}
