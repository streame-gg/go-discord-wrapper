package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestMessageReactionAdd() {
	s.T().Log("Testing Message Reaction Add Unmarshal Logic")

	sub := testutil.InitSub[MessageReactionAddEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"user_id":           discord.RandomSnowflake(),
		"channel_id":        discord.RandomSnowflake(),
		"message_id":        discord.RandomSnowflake(),
		"guild_id":          discord.RandomSnowflake(),
		"member":            testdata.NewGuildMember(),
		"emoji":             testdata.NewEmoji(),
		"message_author_id": discord.RandomSnowflake(),
		"burst":             testutil.RandomBool(),
		"burst_colors":      testutil.RandomStringArray(testutil.RandomIntInRange(1, 5), 7, 7),
		"type": testutil.RandomItem(
			discord.ReactionTypeNormal,
			discord.ReactionTypeBurst,
		),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessageReactionAddEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessageReactionAddEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["user_id"], got.UserID)
				s.EqualValues(payload["message_id"], got.MessageID)
				s.EqualValues(payload["message_author_id"], *got.MessageAuthorID)
				s.EqualValues(payload["burst"], got.Burst)
				s.EqualValues(payload["burst_colors"], got.BurstColors)
				s.EqualValues(payload["type"], got.Type)

				s.compareMember(payload["member"].(map[string]interface{}), *got.Member)
				s.compareEmoji(payload["emoji"].(map[string]interface{}), got.Emoji)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
				s.Nil(got.Message)
				s.Nil(got.User)
			},
		},
	})
}

func (s *eventSuite) TestMessageReactionRemove() {
	s.T().Log("Testing Message Reaction Remove Unmarshal Logic")

	sub := testutil.InitSub[MessageReactionRemoveEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"user_id":    discord.RandomSnowflake(),
		"channel_id": discord.RandomSnowflake(),
		"message_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"emoji":      testdata.NewEmoji(),
		"burst":      testutil.RandomBool(),
		"type": testutil.RandomItem(
			discord.ReactionTypeNormal,
			discord.ReactionTypeBurst,
		),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessageReactionRemoveEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessageReactionRemoveEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["user_id"], got.UserID)
				s.EqualValues(payload["message_id"], got.MessageID)
				s.EqualValues(payload["burst"], got.Burst)
				s.EqualValues(payload["type"], got.Type)

				s.compareEmoji(payload["emoji"].(map[string]interface{}), got.Emoji)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
				s.Nil(got.Message)
				s.Nil(got.User)
			},
		},
	})
}

func (s *eventSuite) TestMessageReactionRemoveAll() {
	s.T().Log("Testing Message Reaction Remove All Unmarshal Logic")

	sub := testutil.InitSub[MessageReactionRemoveAllEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"message_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessageReactionRemoveAllEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessageReactionRemoveAllEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["message_id"], got.MessageID)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
				s.Nil(got.Message)
			},
		},
	})
}

func (s *eventSuite) TestMessageReactionRemoveEmoji() {
	s.T().Log("Testing Message Reaction Remove Emoji Unmarshal Logic")

	sub := testutil.InitSub[MessageReactionRemoveEmojiEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"message_id": discord.RandomSnowflake(),
		"guild_id":   discord.RandomSnowflake(),
		"emoji":      testdata.NewEmoji(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[MessageReactionRemoveEmojiEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got MessageReactionRemoveEmojiEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["message_id"], got.MessageID)

				s.compareEmoji(payload["emoji"].(map[string]interface{}), got.Emoji)

				s.Nil(got.Guild)
				s.Nil(got.Channel)
				s.Nil(got.Message)
			},
		},
	})
}
