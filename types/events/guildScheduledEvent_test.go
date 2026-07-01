package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildScheduledEventCreate() {
	s.T().Log("Testing Guild Scheduled Event Create Unmarshal Logic")

	sub := testutil.InitSub[GuildScheduledEventCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewScheduledEvent()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildScheduledEventCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildScheduledEventCreateEvent) {
				s.compareGuildScheduledEvent(payload, got.GuildScheduledEvent)
			},
		},
	})
}

func (s *eventSuite) TestGuildScheduledEventUpdate() {
	s.T().Log("Testing Guild Scheduled Event Update Unmarshal Logic")

	sub := testutil.InitSub[GuildScheduledEventUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewScheduledEvent()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildScheduledEventUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildScheduledEventUpdateEvent) {
				s.compareGuildScheduledEvent(payload, got.NewScheduledEvent)
				s.Nil(got.OldScheduledEvent)
			},
		},
	})
}

func (s *eventSuite) TestGuildScheduledEventDelete() {
	s.T().Log("Testing Guild Scheduled Event Delete Unmarshal Logic")

	sub := testutil.InitSub[GuildScheduledEventDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewScheduledEvent()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildScheduledEventDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildScheduledEventDeleteEvent) {
				s.compareGuildScheduledEvent(payload, got.GuildScheduledEvent)
			},
		},
	})
}
