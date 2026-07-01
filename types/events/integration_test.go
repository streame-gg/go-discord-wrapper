package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestIntegrationCreate() {
	s.T().Log("Testing Integration Create Unmarshal Logic")

	sub := testutil.InitSub[IntegrationCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewIntegrationWithGuildID()

	sub.RunCases([]testutil.UnmarshalTestCase[IntegrationCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got IntegrationCreateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareIntegration(payload, got.Integration)
			},
		},
	})
}

func (s *eventSuite) TestIntegrationUpdate() {
	s.T().Log("Testing Integration Update Unmarshal Logic")

	sub := testutil.InitSub[IntegrationUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewIntegrationWithGuildID()

	sub.RunCases([]testutil.UnmarshalTestCase[IntegrationUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got IntegrationUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareIntegration(payload, got.NewIntegration)
			},
		},
	})
}

func (s *eventSuite) TestIntegrationDelete() {
	s.T().Log("Testing Integration Delete Unmarshal Logic")

	sub := testutil.InitSub[IntegrationDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id":       discord.RandomSnowflake(),
		"id":             discord.RandomSnowflake(),
		"application_id": discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[IntegrationDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got IntegrationDeleteEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["id"], got.ID)
				s.EqualValues(payload["application_id"], *got.ApplicationID)
			},
		},
	})
}
