package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildMemberAdd() {
	s.T().Log("Testing Guild Member Add Unmarshal Logic")

	sub := testutil.InitSub[GuildMemberAddEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewGuildMemberWithGuildID()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildMemberAddEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildMemberAddEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareMember(payload, got.GuildMember)
			},
		},
	})
}

// TODO
func (s *eventSuite) TestGuildMemberUpdate() {
	s.T().Log("Testing Guild Member Add Unmarshal Logic")
	s.T().Skip("not implemented yet just todo")

	sub := testutil.InitSub[GuildMemberAddEvent](s)

	sub.RunCommonEdgeCases()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildMemberAddEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal("{}"),
			Validate: func(got GuildMemberAddEvent) {

			},
		},
	})
}

func (s *eventSuite) TestGuildMemberRemove() {
	s.T().Log("Testing Guild Member Remove Unmarshal Logic")

	sub := testutil.InitSub[GuildMemberRemoveEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewUserWithGuildID()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildMemberRemoveEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildMemberRemoveEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareUser(payload, got.User)
			},
		},
	})
}

func (s *eventSuite) TestGuildMembersChunk() {
	s.T().Log("Testing Guild Members Chunk Unmarshal Logic")

	sub := testutil.InitSub[GuildMembersChunkEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewGuildMembersChunkEventPayload()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildMembersChunkEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildMembersChunkEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["chunk_index"], got.ChunkIndex)
				s.EqualValues(payload["chunk_count"], got.ChunkCount)
				s.EqualValues(payload["not_found"], got.NotFound)
				s.Equal(payload["nonce"], got.Nonce)

				presences := payload["presences"].([]interface{})
				s.Len(got.Presences, len(presences))

				for i, presence := range got.Presences {
					mappedPresence := presences[i].(map[string]interface{})
					s.comparePresence(mappedPresence, presence)
				}

				members := payload["members"].([]interface{})
				s.Len(got.Members, len(members))
				for i, member := range got.Members {
					mappedMember := members[i].(map[string]interface{})
					s.compareMember(mappedMember, member)
				}
			},
		},
	})
}
