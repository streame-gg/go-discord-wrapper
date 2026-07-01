package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestWebhookUpdate() {
	s.T().Log("Testing Webhook Update Unmarshal Logic")

	sub := testutil.InitSub[WebhooksUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id":   discord.RandomSnowflake(),
		"channel_id": discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[WebhooksUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got WebhooksUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)

				s.Nil(got.Channel)
				s.Nil(got.Guild)
			},
		},
	})
}
