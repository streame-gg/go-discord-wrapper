package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestChannelCreate() {
	sub := testutil.InitSub[ChannelCreateEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelCreateEvent) {
				s.Equal(channel, got.Channel)
			},
		},
		{
			Name:  "empty values tested",
			Input: sub.MustMarshal("{}"),
			Validate: func(got ChannelCreateEvent) {

			},
		},
	})
}

func (s *eventSuite) TestChannelUpdate() {
	sub := testutil.InitSub[ChannelUpdateEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelUpdateEvent) {
				s.Equal(channel, got.NewChannel)
				s.Nil(got.OldChannel)
			},
		},
	})
}

func (s *eventSuite) TestChannelDelete() {
	sub := testutil.InitSub[ChannelDeleteEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelDeleteEvent) {
				s.Equal(channel, got.Channel)
			},
		},
	})
}

func (s *eventSuite) TestChannelPinsUpdate() {
	sub := testutil.InitSub[ChannelPinsUpdateEvent](s)

	sub.RunCommonEdgeCases()

	guildId := discord.RandomSnowflake()
	channelId := discord.RandomSnowflake()
	timestamp := testutil.RandomTime()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelPinsUpdateEvent]{
		{
			Name: "valid payload with guild_id and timestamp",
			Input: sub.MustMarshal(ChannelPinsUpdateEvent{
				GuildID:          util.PointerOf(guildId),
				ChannelID:        channelId,
				LastPinTimestamp: util.PointerOf(timestamp),
			}),
			Validate: func(got ChannelPinsUpdateEvent) {
				s.Equal(guildId, got.GuildID)
				s.Equal(channelId, got.ChannelID)
				s.Equal(timestamp, got.LastPinTimestamp)
			},
		},
		{
			Name: "valid payload without guild_id and timestamp",
			Input: sub.MustMarshal(ChannelPinsUpdateEvent{
				GuildID:          nil,
				ChannelID:        channelId,
				LastPinTimestamp: nil,
			}),
			Validate: func(got ChannelPinsUpdateEvent) {
				s.Nil(got.GuildID)
				s.Equal(channelId, got.ChannelID)
				s.Nil(got.LastPinTimestamp)
			},
		},
	})
}
