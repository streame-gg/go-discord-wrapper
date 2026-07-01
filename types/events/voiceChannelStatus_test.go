package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestVoiceChannelStatusUpdate() {
	s.T().Log("Testing Voice Channel Status Update Unmarshal Logic")

	sub := testutil.InitSub[VoiceChannelStatusUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"id":       discord.RandomSnowflake(),
		"status":   testutil.RandomString(testutil.RandomIntInRange(1, 32)),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[VoiceChannelStatusUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got VoiceChannelStatusUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["id"], got.ChannelID)
				s.EqualValues(payload["status"], *got.Status)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
			},
		},
	})
}

func (s *eventSuite) TestVoiceChannelStartTimeUpdate() {
	s.T().Log("Testing Voice Channel Start Time Update Unmarshal Logic")

	sub := testutil.InitSub[VoiceChannelStartTimeUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id":         discord.RandomSnowflake(),
		"id":               discord.RandomSnowflake(),
		"voice_start_time": testutil.RandomTime().Unix(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[VoiceChannelStartTimeUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got VoiceChannelStartTimeUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["id"], got.ChannelID)
				s.EqualValues(payload["voice_start_time"], got.VoiceStartTime)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
			},
		},
	})
}
