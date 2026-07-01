package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestStageInstanceCreate() {
	s.T().Log("Testing Stage Instance Create Unmarshal Logic")

	sub := testutil.InitSub[StageInstanceCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewStageInstance()

	sub.RunCases([]testutil.UnmarshalTestCase[StageInstanceCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got StageInstanceCreateEvent) {
				s.compareStageInstance(payload, got.StageInstance)
			},
		},
	})
}

func (s *eventSuite) TestStageInstanceUpdate() {
	s.T().Log("Testing Stage Instance Update Unmarshal Logic")

	sub := testutil.InitSub[StageInstanceUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewStageInstance()

	sub.RunCases([]testutil.UnmarshalTestCase[StageInstanceUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got StageInstanceUpdateEvent) {
				s.compareStageInstance(payload, got.NewStage)
				s.Nil(got.OldStage)
			},
		},
	})
}

func (s *eventSuite) TestStageInstanceDelete() {
	s.T().Log("Testing Stage Instance Delete Unmarshal Logic")

	sub := testutil.InitSub[StageInstanceDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewStageInstance()

	sub.RunCases([]testutil.UnmarshalTestCase[StageInstanceDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got StageInstanceDeleteEvent) {
				s.compareStageInstance(payload, got.StageInstance)
			},
		},
	})
}
