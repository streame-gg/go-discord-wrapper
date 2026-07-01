package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestUserUpdate() {
	s.T().Log("Testing User Update Unmarshal Logic")

	sub := testutil.InitSub[UserUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewUser()

	sub.RunCases([]testutil.UnmarshalTestCase[UserUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got UserUpdateEvent) {
				s.compareUser(payload, got.NewUser)
				s.Nil(got.OldUser)
			},
		},
	})
}
