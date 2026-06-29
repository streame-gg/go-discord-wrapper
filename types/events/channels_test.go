package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestChannelCreate() {
	s.T().Log("Testing Channel Create Unmarshal Logic")

	sub := testutil.InitSub[ChannelCreateEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelCreateEvent) {
				s.compareChannel(channel, got.Channel)
			},
		},
	})
}

func (s *eventSuite) TestChannelUpdate() {
	s.T().Log("Testing Channel Update Unmarshal Logic")

	sub := testutil.InitSub[ChannelUpdateEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelUpdateEvent) {
				s.compareChannel(channel, got.NewChannel)
				s.Nil(got.OldChannel)
			},
		},
	})
}

func (s *eventSuite) TestChannelDelete() {
	s.T().Log("Testing Channel Delete Unmarshal Logic")

	sub := testutil.InitSub[ChannelDeleteEvent](s)

	sub.RunCommonEdgeCases()

	channel := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(channel),
			Validate: func(got ChannelDeleteEvent) {
				s.compareChannel(channel, got.Channel)
			},
		},
	})
}

func (s *eventSuite) TestChannelPinsUpdate() {
	s.T().Log("Testing Channel Pings Update Unmarshal Logic")

	sub := testutil.InitSub[ChannelPinsUpdateEvent](s)

	sub.RunCommonEdgeCases()

	guildId := discord.RandomSnowflake()
	channelId := discord.RandomSnowflake()
	timestamp := testutil.RandomTime()

	sub.RunCases([]testutil.UnmarshalTestCase[ChannelPinsUpdateEvent]{
		{
			Name: "valid payload with guild_id and timestamp",
			Input: sub.MustMarshal(map[string]interface{}{
				"guild_id":           guildId,
				"channel_id":         channelId,
				"last_pin_timestamp": timestamp.Format(time.RFC3339),
			}),
			Validate: func(got ChannelPinsUpdateEvent) {
				s.EqualValues(guildId, *got.GuildID)
				s.EqualValues(channelId, got.ChannelID)
				s.EqualValues(timestamp, *got.LastPinTimestamp)
			},
		},
		{
			Name: "valid payload without guild_id and timestamp",
			Input: sub.MustMarshal(map[string]interface{}{
				"guild_id":           nil,
				"channel_id":         channelId,
				"last_pin_timestamp": nil,
			}),
			Validate: func(got ChannelPinsUpdateEvent) {
				s.Nil(got.GuildID)
				s.EqualValues(channelId, got.ChannelID)
				s.Nil(got.LastPinTimestamp)
			},
		},
	})
}
