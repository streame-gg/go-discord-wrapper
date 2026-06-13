package commands

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/suite"
)

type optionSuite struct {
	suite.Suite
}

func TestOptionSuite(t *testing.T) { suite.Run(t, new(optionSuite)) }

func (s *optionSuite) TestNullValues_Integer_Marshal_1() {
	option := ApplicationCommandOptionInteger{
		Name: "amount",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "anzahl",
		},
		Description: "Provide an amount.",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "Gebe eine Anzahl an.",
		},
		Required: true,
		MinValue: util.PointerOf(int64(0)),
		MaxValue: util.PointerOf(int64(100)),
	}

	marshaled, err := option.MarshalJSON()
	s.Require().NoError(err)
	s.Require().Equal("{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"min_value\":0,\"max_value\":100,\"autocomplete\":false,\"type\":4}", string(marshaled))
}

func (s *optionSuite) TestNullValues_Integer_Marshal_2() {
	option := ApplicationCommandOptionInteger{
		Name: "amount",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "anzahl",
		},
		Description: "Provide an amount.",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "Gebe eine Anzahl an.",
		},
		Required: true,
		MinValue: nil,
		MaxValue: nil,
	}

	marshaled, err := option.MarshalJSON()
	s.Require().NoError(err)
	s.Require().Equal("{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"autocomplete\":false,\"type\":4}", string(marshaled))
}

func (s *optionSuite) TestNullValues_Integer_Unmarshal_1() {
	option := "{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"min_value\":0,\"max_value\":100,\"autocomplete\":false,\"type\":4}"
	optionUnmarshaled := ApplicationCommandOptionInteger{}
	s.Require().NoError(optionUnmarshaled.UnmarshalJSON([]byte(option)))
	s.Require().Equal(int64(0), *optionUnmarshaled.MinValue)
	s.Require().Equal(int64(100), *optionUnmarshaled.MaxValue)
}

func (s *optionSuite) TestNullValues_Integer_Unmarshal_2() {
	option := "{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"autocomplete\":false,\"type\":4}"
	optionUnmarshaled := ApplicationCommandOptionInteger{}
	s.Require().NoError(optionUnmarshaled.UnmarshalJSON([]byte(option)))
	s.Require().Nil(optionUnmarshaled.MinValue)
	s.Require().Nil(optionUnmarshaled.MaxValue)
}

func (s *optionSuite) TestNullValues_Number_Marshal_1() {
	option := ApplicationCommandOptionNumber{
		Name: "amount",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "anzahl",
		},
		Description: "Provide an amount.",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "Gebe eine Anzahl an.",
		},
		Required: true,
		MinValue: util.PointerOf(float64(0)),
		MaxValue: util.PointerOf(float64(100)),
	}

	marshaled, err := option.MarshalJSON()
	s.Require().NoError(err)
	s.Require().Equal("{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"min_value\":0,\"max_value\":100,\"autocomplete\":false,\"type\":10}", string(marshaled))
}

func (s *optionSuite) TestNullValues_Number_Marshal_2() {
	option := ApplicationCommandOptionNumber{
		Name: "amount",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "anzahl",
		},
		Description: "Provide an amount.",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "Gebe eine Anzahl an.",
		},
		Required: true,
		MinValue: nil,
		MaxValue: nil,
	}

	marshaled, err := option.MarshalJSON()
	s.Require().NoError(err)
	s.Require().Equal("{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"autocomplete\":false,\"type\":10}", string(marshaled))
}

func (s *optionSuite) TestNullValues_Number_Unmarshal_1() {
	option := "{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"min_value\":0,\"max_value\":100,\"autocomplete\":false,\"type\":10}"
	optionUnmarshaled := ApplicationCommandOptionInteger{}
	s.Require().NoError(optionUnmarshaled.UnmarshalJSON([]byte(option)))
	s.Require().Equal(int64(0), *optionUnmarshaled.MinValue)
	s.Require().Equal(int64(100), *optionUnmarshaled.MaxValue)
}

func (s *optionSuite) TestNullValues_Number_Unmarshal_2() {
	option := "{\"name\":\"amount\",\"name_localizations\":{\"de\":\"anzahl\"},\"description\":\"Provide an amount.\",\"description_localizations\":{\"de\":\"Gebe eine Anzahl an.\"},\"required\":true,\"autocomplete\":false,\"type\":10}"
	optionUnmarshaled := ApplicationCommandOptionInteger{}
	s.Require().NoError(optionUnmarshaled.UnmarshalJSON([]byte(option)))
	s.Require().Nil(optionUnmarshaled.MinValue)
	s.Require().Nil(optionUnmarshaled.MaxValue)
}

func (s *optionSuite) TestNullValues_String_Marshal_1() {
	option := ApplicationCommandOptionString{
		Name: "text",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "text",
		},
		Description: "Provide some text.",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "Gebe einen Text an.",
		},
		Required:  true,
		MinLength: util.PointerOf(int64(0)),
		MaxLength: util.PointerOf(int64(100)),
	}

	marshaled, err := option.MarshalJSON()
	s.Require().NoError(err)
	s.Require().Equal("{\"name\":\"text\",\"name_localizations\":{\"de\":\"text\"},\"description\":\"Provide some text.\",\"description_localizations\":{\"de\":\"Gebe einen Text an.\"},\"required\":true,\"autocomplete\":false,\"min_length\":0,\"max_length\":100,\"type\":3}", string(marshaled))
}

func (s *optionSuite) TestNullValues_String_Marshal_2() {
	option := ApplicationCommandOptionString{
		Name: "text",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "text",
		},
		Description: "Provide some text.",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleGerman: "Gebe einen Text an.",
		},
		Required:  true,
		MinLength: nil,
		MaxLength: nil,
	}

	marshaled, err := option.MarshalJSON()
	s.Require().NoError(err)
	s.Require().Equal("{\"name\":\"text\",\"name_localizations\":{\"de\":\"text\"},\"description\":\"Provide some text.\",\"description_localizations\":{\"de\":\"Gebe einen Text an.\"},\"required\":true,\"autocomplete\":false,\"type\":3}", string(marshaled))
}

func (s *optionSuite) TestNullValues_String_Unmarshal_1() {
	option := "{\"name\":\"text\",\"name_localizations\":{\"de\":\"text\"},\"description\":\"Provide some text.\",\"description_localizations\":{\"de\":\"Gebe einen Text an.\"},\"required\":true,\"autocomplete\":false,\"min_length\":0,\"max_length\":100,\"type\":3}"
	optionUnmarshaled := ApplicationCommandOptionString{}
	s.Require().NoError(optionUnmarshaled.UnmarshalJSON([]byte(option)))
	s.Require().Equal(int64(0), *optionUnmarshaled.MinLength)
	s.Require().Equal(int64(100), *optionUnmarshaled.MaxLength)
}

func (s *optionSuite) TestNullValues_String_Unmarshal_2() {
	option := "{\"name\":\"text\",\"name_localizations\":{\"de\":\"text\"},\"description\":\"Provide some text.\",\"description_localizations\":{\"de\":\"Gebe einen Text an.\"},\"required\":true,\"autocomplete\":false,\"type\":3}"
	optionUnmarshaled := ApplicationCommandOptionString{}
	s.Require().NoError(optionUnmarshaled.UnmarshalJSON([]byte(option)))
	s.Require().Nil(optionUnmarshaled.MinLength)
	s.Require().Nil(optionUnmarshaled.MaxLength)
}
