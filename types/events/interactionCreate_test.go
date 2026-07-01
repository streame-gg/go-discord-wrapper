package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
)

// TODO mock interaction & event
func (s *eventSuite) TestInteractionCreate() {
	s.T().Log("Testing Interaction Create Unmarshal Logic")
	s.T().Skip("TestInteractionCreate skipped because not implemented yet")

	sub := testutil.InitSub[InteractionCreateEvent](s)

	sub.RunCommonEdgeCases()

	sub.RunCases([]testutil.UnmarshalTestCase[InteractionCreateEvent]{
		{
			Name:     "valid full payload",
			Input:    sub.MustMarshal("{}"),
			Validate: func(got InteractionCreateEvent) {},
		},
	})
}
