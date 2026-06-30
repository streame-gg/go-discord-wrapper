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

func (s *eventSuite) TestGuildMemberUpdate() {
	s.T().Log("Testing Guild Member Update Unmarshal Logic")

	sub := testutil.InitSub[GuildMemberUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewMemberUpdateEventPayload()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildMemberUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildMemberUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["avatar"], *got.Avatar)
				s.EqualValues(payload["banner"], *got.Banner)
				s.EqualValues(payload["communication_disabled_until"], *got.CommunicationDisabledUntil)
				s.EqualValues(payload["deaf"], *got.Deaf)
				s.EqualValues(payload["joined_at"], *got.JoinedAt)
				s.EqualValues(payload["mute"], *got.Mute)
				s.EqualValues(payload["nick"], *got.Nick)
				s.EqualValues(payload["pending"], *got.Pending)
				s.EqualValues(payload["premium_since"], *got.PremiumSince)
				s.EqualValues(payload["roles"], got.Roles)
				s.compareUser(payload["user"].(map[string]interface{}), got.User)

				nameplate := payload["collectibles"].(map[string]interface{})["nameplate"].(map[string]interface{})
				s.EqualValues(nameplate["sku_id"], got.Collectibles.Nameplate.SkuID)
				s.EqualValues(nameplate["asset"], got.Collectibles.Nameplate.Asset)
				s.EqualValues(nameplate["label"], got.Collectibles.Nameplate.Label)
				s.EqualValues(nameplate["palette"], got.Collectibles.Nameplate.Palette)

				avatarDecorationData := payload["avatar_decoration_data"].(map[string]interface{})
				s.EqualValues(avatarDecorationData["asset"], got.AvatarDecorationData.Asset)
				s.EqualValues(avatarDecorationData["sku_id"], got.AvatarDecorationData.SkuID)
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
