package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestGuildRoleCreate() {
	s.T().Log("Testing Guild Role Create Unmarshal Logic")

	sub := testutil.InitSub[GuildRoleCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"role":     testdata.NewRole(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[GuildRoleCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildRoleCreateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareRole(payload["role"].(map[string]interface{}), got.Role)
			},
		},
	})
}

func (s *eventSuite) TestGuildRoleUpdate() {
	s.T().Log("Testing Guild Role Update Unmarshal Logic")

	sub := testutil.InitSub[GuildRoleUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"role":     testdata.NewRole(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[GuildRoleUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildRoleUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareRole(payload["role"].(map[string]interface{}), got.NewRole)
				s.Nil(got.OldRole)
			},
		},
	})
}

func (s *eventSuite) TestGuildRoleDelete() {
	s.T().Log("Testing Guild Role Delete Unmarshal Logic")

	sub := testutil.InitSub[GuildRoleDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"role_id":  discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[GuildRoleDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildRoleDeleteEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["role_id"], got.RoleID)
			},
		},
	})
}
