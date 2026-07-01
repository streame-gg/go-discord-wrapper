package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestTypingStart() {
	s.T().Log("Testing User Update Unmarshal Logic")

	sub := testutil.InitSub[TypingStartEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"user_id":    discord.RandomSnowflake(),
		"timestamp":  testutil.RandomTime().Unix(),
		"member":     testdata.NewGuildMember(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[TypingStartEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got TypingStartEvent) {
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["user_id"], got.UserID)
				s.EqualValues(payload["timestamp"], got.Timestamp)
				s.compareMember(payload["member"].(map[string]interface{}), *got.Member)
			},
		},
	})
}
