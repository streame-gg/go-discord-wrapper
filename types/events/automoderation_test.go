package events

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestAutomoderationRuleCreate() {
	s.T().Log("Testing Automoderation Rule Create Unmarshal Logic")

	sub := testutil.InitSub[AutoModerationRuleCreateEvent](s)

	sub.RunCommonEdgeCases()

	randomAutomoderationRule := testdata.NewAutomoderationRule()

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleCreateEvent]{
		{
			Name:  "valid payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(t *testing.T, e AutoModerationRuleCreateEvent) {
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
	})
}

func (s *eventSuite) TestAutomoderationRuleUpdate() {
	s.T().Log("Testing Automoderation Rule Update Unmarshal Logic")

	sub := testutil.InitSub[AutoModerationRuleUpdateEvent](s)

	sub.RunCommonEdgeCases()

	randomAutomoderationRule := testdata.NewAutomoderationRule()

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleUpdateEvent]{
		{
			Name:  "valid payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(t *testing.T, e AutoModerationRuleUpdateEvent) {
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
	})
}

func (s *eventSuite) TestAutomoderationRuleDelete() {
	s.T().Log("Testing Automoderation Rule Delete Unmarshal Logic")

	sub := testutil.InitSub[AutoModerationRuleDeleteEvent](s)

	sub.RunCommonEdgeCases()

	randomAutomoderationRule := testdata.NewAutomoderationRule()

	sub.RunCases([]testutil.UnmarshalTestCase[AutoModerationRuleDeleteEvent]{
		{
			Name:  "valid payload",
			Input: sub.MustMarshal(randomAutomoderationRule),
			Validate: func(t *testing.T, e AutoModerationRuleDeleteEvent) {
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
	})
}
