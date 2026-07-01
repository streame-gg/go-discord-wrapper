package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestSubscriptionCreate() {
	s.T().Log("Testing Subscription Create Unmarshal Logic")

	sub := testutil.InitSub[SubscriptionCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewSubscription()

	sub.RunCases([]testutil.UnmarshalTestCase[SubscriptionCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got SubscriptionCreateEvent) {
				s.compareSubscription(payload, got.Subscription)
			},
		},
	})
}

func (s *eventSuite) TestSubscriptionUpdate() {
	s.T().Log("Testing Subscription Update Unmarshal Logic")

	sub := testutil.InitSub[SubscriptionUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewSubscription()

	sub.RunCases([]testutil.UnmarshalTestCase[SubscriptionUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got SubscriptionUpdateEvent) {
				s.compareSubscription(payload, got.NewSubscription)
			},
		},
	})
}

func (s *eventSuite) TestSubscriptionDelete() {
	s.T().Log("Testing Subscription Delete Unmarshal Logic")

	sub := testutil.InitSub[SubscriptionDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewSubscription()

	sub.RunCases([]testutil.UnmarshalTestCase[SubscriptionDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got SubscriptionDeleteEvent) {
				s.compareSubscription(payload, got.Subscription)
			},
		},
	})
}
