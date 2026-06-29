package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildEmojisUpdate() {
	s.T().Log("Testing Guild Emojis Update Unmarshal Logic")

	sub := testutil.InitSub[GuildEmojisUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewGuildEmojisUpdateEventPayload()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildEmojisUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildEmojisUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)

				emojis := payload["emojis"].([]map[string]interface{})
				s.Len(emojis, len(emojis))

				for i, em := range emojis {
					s.compareEmoji(em, got.NewEmojis[i])
				}

				s.Nil(got.OldEmojis)
			},
		},
	})
}

func (s *eventSuite) TestGuildStickersUpdate() {
	s.T().Log("Testing Guild Stickers Update Unmarshal Logic")

	sub := testutil.InitSub[GuildStickersUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewGuildStickersUpdateEventPayload()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildStickersUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildStickersUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)

				stickers := payload["stickers"].([]map[string]interface{})
				for i, st := range stickers {
					s.EqualValues(st["id"], got.NewStickers[i].ID)
					s.EqualValues(st["name"], got.NewStickers[i].Name)
					s.EqualValues(st["description"], *got.NewStickers[i].Description)
					s.EqualValues(st["tags"], got.NewStickers[i].Tags)
					s.EqualValues(st["type"], got.NewStickers[i].Type)
					s.EqualValues(st["format_type"], got.NewStickers[i].FormatType)
					s.EqualValues(st["available"], *got.NewStickers[i].Available)
					s.EqualValues(st["sort_value"], got.NewStickers[i].SortValue)

					s.compareUser(st["user"].(map[string]interface{}), *got.NewStickers[i].User)
				}

				s.Nil(got.OldStickers)
			},
		},
	})
}
