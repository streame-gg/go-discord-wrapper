package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildBanAdd() {
	s.T().Log("Testing Guild Ban Add Unmarshal Logic")

	sub := testutil.InitSub[GuildBanAddEvent](s)

	sub.RunCommonEdgeCases()

	ban := testdata.NewGuildBanAddOrRemove()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildBanAddEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(ban),
			Validate: func(got GuildBanAddEvent) {
				s.EqualValues(ban["guild_id"], got.GuildID)
				s.compareUser(ban["user"].(map[string]interface{}), got.User)
			},
		},
	})
}

func (s *eventSuite) TestGuildBanRemove() {
	s.T().Log("Testing Guild Ban Remove Unmarshal Logic")

	sub := testutil.InitSub[GuildBanRemoveEvent](s)

	sub.RunCommonEdgeCases()

	ban := testdata.NewGuildBanAddOrRemove()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildBanRemoveEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(ban),
			Validate: func(got GuildBanRemoveEvent) {
				s.EqualValues(ban["guild_id"], got.GuildID)
				s.compareUser(ban["user"].(map[string]interface{}), got.User)
			},
		},
	})
}
