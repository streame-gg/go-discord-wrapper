package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestGuildIntegrationsUpdate() {
	s.T().Log("Testing Guild Integrations Update Unmarshal Logic")

	sub := testutil.InitSub[GuildIntegrationsUpdateEvent](s)

	sub.RunCommonEdgeCases()

	guildId := discord.RandomSnowflake()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildIntegrationsUpdateEvent]{
		{
			Name: "valid full payload",
			Input: sub.MustMarshal(map[string]interface{}{
				"guild_id": guildId,
			}),
			Validate: func(got GuildIntegrationsUpdateEvent) {
				s.EqualValues(guildId, got.GuildID)
			},
		},
	})
}
