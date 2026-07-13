package responses

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *responsesSuite) TestCommand() {
	sub := testutil.InitSub[InteractionDataApplicationCommand](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewInteractionDataApplicationCommand()

	sub.RunCases([]testutil.UnmarshalTestCase[InteractionDataApplicationCommand]{
		{
			Name:  "unmarshal valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got InteractionDataApplicationCommand) {
				s.Require().NotNil(got)
				s.compareInteractionData(payload, &got)
			},
		},
	})
}
