package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestPresenceUpdate() {
	s.T().Log("Testing Presence Update Unmarshal Logic")

	sub := testutil.InitSub[PresenceUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewPresence()

	sub.RunCases([]testutil.UnmarshalTestCase[PresenceUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got PresenceUpdateEvent) {
				s.comparePresence(payload, got.NewPresence)
				s.Nil(got.OldPresence)
			},
		},
	})
}
