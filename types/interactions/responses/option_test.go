package responses

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *responsesSuite) TestOptions_String() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[string]](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewApplicationCommandInteractionDataOptionStructure(discord.ApplicationCommandOptionTypeString)

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[string]]{
		{
			Name:  "unmarshal string option",
			Input: sub.MustMarshal(payload),
			Validate: func(got ApplicationCommandInteractionDataOption[string]) {
				s.Require().NotNil(got)
				s.Equal(payload["value"], *got.Value)
			},
		},
	})
}

func (s *responsesSuite) TestOptions_Number() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[float64]](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewApplicationCommandInteractionDataOptionStructure(discord.ApplicationCommandOptionTypeNumber)

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[float64]]{
		{
			Name:  "unmarshal number option",
			Input: sub.MustMarshal(payload),
			Validate: func(got ApplicationCommandInteractionDataOption[float64]) {
				s.Require().NotNil(got)
				s.Equal(payload["value"], *got.Value)
			},
		},
	})
}

func (s *responsesSuite) TestOptions_Integer() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[int]](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewApplicationCommandInteractionDataOptionStructure(discord.ApplicationCommandOptionTypeInteger)

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[int]]{
		{
			Name:  "unmarshal integer option",
			Input: sub.MustMarshal(payload),
			Validate: func(got ApplicationCommandInteractionDataOption[int]) {
				s.Require().NotNil(got)
				s.Equal(payload["value"], *got.Value)
			},
		},
	})
}

func (s *responsesSuite) TestOptions_Bool() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[bool]](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewApplicationCommandInteractionDataOptionStructure(discord.ApplicationCommandOptionTypeBoolean)

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[bool]]{
		{
			Name:  "unmarshal bool option",
			Input: sub.MustMarshal(payload),
			Validate: func(got ApplicationCommandInteractionDataOption[bool]) {
				s.Require().NotNil(got)
				s.Equal(payload["value"], *got.Value)
			},
		},
	})
}

func (s *responsesSuite) TestOptions_Snowflake() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[discord.Snowflake]](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewApplicationCommandInteractionDataOptionStructure(
		discord.ApplicationCommandOptionTypeUser,
		discord.ApplicationCommandOptionTypeChannel,
		discord.ApplicationCommandOptionTypeRole,
		discord.ApplicationCommandOptionTypeMentionable,
	)

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[discord.Snowflake]]{
		{
			Name:  "unmarshal snowflake option",
			Input: sub.MustMarshal(payload),
			Validate: func(got ApplicationCommandInteractionDataOption[discord.Snowflake]) {
				s.Require().NotNil(got)
				s.Equal(payload["value"], *got.Value)
			},
		},
	})
}

func (s *responsesSuite) TestOptions_NotProvided() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[interface{}]](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"value": nil,
		"type":  discord.ApplicationCommandOptionTypeUser,
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[interface{}]]{
		{
			Name:  "unmarshal non existent value option",
			Input: sub.MustMarshal(payload),
			Validate: func(got ApplicationCommandInteractionDataOption[interface{}]) {
				s.Require().NotNil(got)
				s.Nil(got.Value)
			},
		},
	})
}

func (s *responsesSuite) TestOptions_TypeNotProvided() {
	sub := testutil.InitSub[ApplicationCommandInteractionDataOption[interface{}]](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"type":  0,
		"value": "lololol!",
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandInteractionDataOption[interface{}]]{
		{
			Name:    "fail due to missing type",
			Input:   sub.MustMarshal(payload),
			WantErr: true,
		},
	})
}
