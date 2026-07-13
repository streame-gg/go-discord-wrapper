package responses

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *responsesSuite) TestAutocomplete() {
	sub := testutil.InitSub[InteractionDataAutocomplete](s)

	sub.RunCommonEdgeCases()

	autocomplete := InteractionDataAutocomplete{}
	s.Equal(discord.InteractionTypeApplicationCommandAutocomplete, autocomplete.GetType())

	payload := testdata.NewInteractionDataApplicationCommand()

	sub.RunCases([]testutil.UnmarshalTestCase[InteractionDataAutocomplete]{
		{
			Name:  "unmarshal valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got InteractionDataAutocomplete) {
				s.Require().NotNil(got)
				s.compareInteractionData(payload, &got)
			},
		},
	})

	irob := InteractionResponseDataAutocomplete{
		Choices: []AutocompleteChoice{
			{
				Name: "string",
				NameLocalizations: map[discord.Locale]string{
					discord.LocaleGerman: "Zeichenkette",
					discord.LocaleDanish: "streng",
				},
				Value: "string",
			},
			{
				Name: "number",
				NameLocalizations: map[discord.Locale]string{
					discord.LocaleGerman: "Nummer",
					discord.LocaleDanish: "nummer",
				},
				Value: 3,
			},
			{
				Name: "float",
				NameLocalizations: map[discord.Locale]string{
					discord.LocaleGerman: "Gleitkommazahl",
					discord.LocaleDanish: "flydende komma-tal",
				},
				Value: 3.14,
			},
		},
	}

	marshal, err := json.Marshal(irob)
	s.Require().NoError(err)
	s.Require().Equal("{\"choices\":[{\"name\":\"string\",\"name_localizations\":{\"da\":\"streng\",\"de\":\"Zeichenkette\"},\"value\":\"string\"},{\"name\":\"number\",\"name_localizations\":{\"da\":\"nummer\",\"de\":\"Nummer\"},\"value\":3},{\"name\":\"float\",\"name_localizations\":{\"da\":\"flydende komma-tal\",\"de\":\"Gleitkommazahl\"},\"value\":3.14}]}", string(marshal))
}
