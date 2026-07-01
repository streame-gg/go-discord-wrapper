package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestVoiceStateUpdate() {
	s.T().Log("Testing Voice State Update Unmarshal Logic")

	sub := testutil.InitSub[VoiceStateUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewVoiceState()

	sub.RunCases([]testutil.UnmarshalTestCase[VoiceStateUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got VoiceStateUpdateEvent) {
				s.compareVoiceState(payload, got.NewState)
				s.Nil(got.OldState)
			},
		},
	})
}

func (s *eventSuite) TestVoiceServerUpdate() {
	s.T().Log("Testing Voice Server Update Unmarshal Logic")

	sub := testutil.InitSub[VoiceServerUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"token":    testutil.RandomString(32),
		"guild_id": discord.RandomSnowflake(),
		"endpoint": testutil.RandomString(32),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[VoiceServerUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got VoiceServerUpdateEvent) {
				s.EqualValues(payload["token"], got.Token)
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["endpoint"], *got.Endpoint)
			},
		},
	})
}

func (s *eventSuite) TestVoiceChannelEffectSend() {
	s.T().Log("Testing Voice Channel Effect Send Unmarshal Logic")

	sub := testutil.InitSub[VoiceChannelEffectSendEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"user_id":    discord.RandomSnowflake(),
		"emoji":      testdata.NewEmoji(),
		"animation_type": testutil.RandomItem(
			discord.VoiceChannelAnimationTypePremium,
			discord.VoiceChannelAnimationTypeBasic,
		),
		"animation_id": testutil.RandomIntInRange(1, 10),
		"sound_id":     discord.RandomSnowflake(),
		"sound_volume": testutil.RandomFloat64InRange(0, 1),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[VoiceChannelEffectSendEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got VoiceChannelEffectSendEvent) {
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["user_id"], got.UserID)
				s.EqualValues(payload["animation_type"], *got.AnimationType)
				s.EqualValues(payload["animation_id"], *got.AnimationID)
				s.EqualValues(payload["sound_id"], *got.SoundID)
				s.EqualValues(payload["sound_volume"], *got.SoundVolume)

				s.compareEmoji(payload["emoji"].(map[string]interface{}), *got.Emoji)

				s.Nil(got.Channel)
				s.Nil(got.Guild)
				s.Nil(got.User)
				s.Nil(got.Sound)
			},
		},
	})
}
