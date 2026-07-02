package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestInviteCreate() {
	s.T().Log("Testing Invite Create Unmarshal Logic")

	sub := testutil.InitSub[InviteCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"code":       testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"created_at": testutil.RandomTime(),
		"guild_id":   discord.RandomSnowflake(),
		"inviter":    testdata.NewUser(),
		"max_age":    testutil.RandomIntInRange(1, 86400),
		"max_uses":   testutil.RandomIntInRange(1, 100),
		"target_user_type": testutil.RandomItem(
			discord.InviteTargetUserTypeStream,
			discord.InviteTargetUserTypeEmbeddedApplication,
		),
		"target_application": map[string]interface{}{
			"id":          discord.RandomSnowflake(),
			"name":        testutil.RandomString(testutil.RandomIntInRange(1, 32)),
			"description": testutil.RandomString(testutil.RandomIntInRange(1, 100)),
			"icon":        testutil.RandomString(32),
		},
		"temporary":  testutil.RandomBool(),
		"uses":       testutil.RandomIntInRange(1, 100),
		"expires_at": testutil.RandomTime(),
		"role_ids":   testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 50)),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[InviteCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got InviteCreateEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["code"], got.Code)
				s.EqualValues(payload["channel_id"], got.ChannelID)
				s.EqualValues(payload["created_at"], got.CreatedAt)
				s.EqualValues(payload["max_age"], *got.MaxAge)
				s.EqualValues(payload["max_uses"], *got.MaxUses)
				s.EqualValues(payload["target_user_type"], *got.TargetUserType)

				application := payload["target_application"].(map[string]interface{})
				s.EqualValues(application["id"], got.TargetApplication.ID)
				s.EqualValues(application["name"], got.TargetApplication.Name)
				s.EqualValues(application["icon"], *got.TargetApplication.Icon)
				s.EqualValues(application["description"], got.TargetApplication.Description)

				s.compareUser(payload["inviter"].(map[string]interface{}), *got.Inviter)

			},
		},
	})
}

func (s *eventSuite) TestInviteDelete() {
	s.T().Log("Testing Invite Delete Unmarshal Logic")

	sub := testutil.InitSub[InviteDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"channel_id": discord.RandomSnowflake(),
		"code":       testutil.RandomString(testutil.RandomIntInRange(1, 32)),
		"guild_id":   discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[InviteDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got InviteDeleteEvent) {
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["code"], got.Code)
				s.EqualValues(payload["channel_id"], got.ChannelID)
			},
		},
	})
}
