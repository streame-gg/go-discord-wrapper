package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestChannelCreate() {
	sub := testutil.InitSub[ChannelCreateEvent](s)

	sub.RunCommonEdgeCases()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelCreateEvent]{
		{
			Name: "valid payload",
			Input: sub.MustMarshal(map[string]any{
				"id":       123,
				"type":     discord.ChannelTypeGuildText,
				"guild_id": 456,
			}),
		},
	})
}
