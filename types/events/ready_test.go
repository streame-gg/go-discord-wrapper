package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestReadyEvent() {
	s.T().Log("Testing Ready Unmarshal Logic")

	sub := testutil.InitSub[ReadyEvent](s)

	sub.RunCommonEdgeCases()

	shardId := testutil.RandomIntInRange(1, 9)
	shardTotal := testutil.RandomIntInRange(10, 20)

	payload := map[string]interface{}{
		"v":    discord.APIVersion10,
		"user": testdata.NewUser(),
		"guilds": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 512), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, testdata.NewUnavailableGuild())
		}),
		"session_id":         testutil.RandomString(32),
		"resume_gateway_url": testutil.RandomString(32),
		"shard":              []int{shardId, shardTotal},
		"application": map[string]interface{}{
			"id": discord.RandomSnowflake(),
			"flags": testutil.RandomFlags(
				discord.ApplicationFlagApplicationAutoModerationRuleCreateBadge,
				discord.ApplicationFlagGatewayPresence,
				discord.ApplicationFlagGatewayPresenceLimited,
				discord.ApplicationFlagGatewayGuildMembers,
				discord.ApplicationFlagGatewayGuildMembersLimited,
				discord.ApplicationFlagVerificationPendingGuildLimit,
				discord.ApplicationFlagEmbedded,
				discord.ApplicationFlagGatewayMessageContent,
				discord.ApplicationFlagGatewayMessageContentLimited,
				discord.ApplicationFlagApplicationCommandBadge,
			),
		},
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ReadyEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ReadyEvent) {
				s.EqualValues(payload["v"], got.V)
				s.EqualValues(payload["session_id"], got.SessionID)
				s.EqualValues(payload["resume_gateway_url"], got.ResumeGatewayURL)
				s.EqualValues(shardTotal, got.Shard.NumShards)
				s.EqualValues(shardId, got.Shard.ShardID)

				application := payload["application"].(map[string]interface{})
				s.EqualValues(application["id"], got.Application.ID)
				s.EqualValues(application["flags"], got.Application.Flags)

				s.compareUser(payload["user"].(map[string]interface{}), got.User)

				guilds := payload["guilds"].([]map[string]interface{})
				s.Len(got.Guilds, len(guilds))
				for i, guild := range guilds {
					s.EqualValues(guild["id"], got.Guilds[i].ID)
					s.True(guild["unavailable"].(bool))
				}
			},
		},
	})
}

func (s *eventSuite) TestResumedEvent() {
	s.T().Log("Testing Resumed Unmarshal Logic")

	sub := testutil.InitSub[ResumedEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"token":      testutil.RandomString(32),
		"seq":        testutil.RandomIntInRange(1, 10000),
		"session_id": testutil.RandomString(32),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ResumedEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ResumedEvent) {
				s.EqualValues(payload["token"], got.Token)
				s.EqualValues(payload["seq"], got.Seq)
				s.EqualValues(payload["session_id"], got.SessionID)
			},
		},
	})
}
