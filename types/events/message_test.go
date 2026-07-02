package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestMessageDelete() {
	s.T().Log("Testing Message Delete Unmarshal Logic")

	sub := testutil.InitSub[MessageDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"id":         discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessageDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessageDeleteEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["id"], got.ID)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
			},
		},
	})
}

func (s *eventSuite) TestMessageDeleteBulk() {
	s.T().Log("Testing Message Delete Bulk Unmarshal Logic")

	sub := testutil.InitSub[MessageDeleteBulkEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"ids":        testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 32)),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessageDeleteBulkEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessageDeleteBulkEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["ids"], got.IDs)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
			},
		},
	})
}

func (s *eventSuite) TestMessagePollVoteAdd() {
	s.T().Log("Testing Message Poll Vote Add Unmarshal Logic")

	sub := testutil.InitSub[MessagePollVoteAddEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"user_id":    discord.RandomSnowflake(),
		"message_id": discord.RandomSnowflake(),
		"answer_id":  testutil.RandomIntInRange(1, 10),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessagePollVoteAddEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessagePollVoteAddEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["user_id"], got.UserID)
				s.EqualValues(payload["message_id"], got.MessageID)
				s.EqualValues(payload["answer_id"], got.AnswerID)
			},
		},
	})
}

func (s *eventSuite) TestMessagePollVoteRemove() {
	s.T().Log("Testing Message Poll Vote Remove Unmarshal Logic")

	sub := testutil.InitSub[MessagePollVoteRemoveEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"user_id":    discord.RandomSnowflake(),
		"message_id": discord.RandomSnowflake(),
		"answer_id":  testutil.RandomIntInRange(1, 10),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessagePollVoteRemoveEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessagePollVoteRemoveEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["user_id"], got.UserID)
				s.EqualValues(payload["message_id"], got.MessageID)
				s.EqualValues(payload["answer_id"], got.AnswerID)
			},
		},
	})
}
