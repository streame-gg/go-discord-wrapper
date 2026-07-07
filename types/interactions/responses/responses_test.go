package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// responsesSuite covers the interaction-response data types: their MarshalJSON
// helpers, the IsInteractionResponseData marker, and the GetType discriminators
// the gateway dispatcher relies on.
type responsesSuite struct {
	suite.Suite
}

func TestResponsesSuite(t *testing.T) {
	suite.Run(t, new(responsesSuite))
}

func (s *responsesSuite) TestDefaultDataMarshalJSON() {
	d := &InteractionResponseDataDefault{Content: discord.Some("hello"), Flags: discord.Some(discord.MessageFlagEphemeral)}
	s.True(d.IsInteractionResponseData())

	b, err := json.Marshal(d)
	s.Require().NoError(err)

	var out map[string]any
	s.Require().NoError(json.Unmarshal(b, &out))
	s.Equal("hello", out["content"])
	// The pkg Files field is tagged json:"-" and must not be serialised.
	s.NotContains(out, "Files")
}

func (s *responsesSuite) TestAutocompleteDataMarshalJSON() {
	d := &InteractionResponseDataAutocomplete{
		Choices: []AutocompleteChoice{{Name: "One", Value: 1}, {Name: "Two", Value: "two"}},
	}
	s.True(d.IsInteractionResponseData())

	b, err := json.Marshal(d)
	s.Require().NoError(err)

	var out struct {
		Choices []struct {
			Name  string `json:"name"`
			Value any    `json:"value"`
		} `json:"choices"`
	}
	s.Require().NoError(json.Unmarshal(b, &out))
	s.Require().Len(out.Choices, 2)
	s.Equal("One", out.Choices[0].Name)
}

func (s *responsesSuite) TestInteractionResponseEnvelopeMarshal() {
	resp := InteractionResponse{
		Type: discord.InteractionCallbackTypeChannelMessageWithSource,
		Data: &InteractionResponseDataDefault{Content: discord.Some("hi")},
	}
	b, err := json.Marshal(resp)
	s.Require().NoError(err)

	var out struct {
		Type int `json:"type"`
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}
	s.Require().NoError(json.Unmarshal(b, &out))
	s.Equal(int(discord.InteractionCallbackTypeChannelMessageWithSource), out.Type)
	s.Equal("hi", out.Data.Content)
}

// TestGetTypeDiscriminators pins each interaction-data type's discriminator,
// which the dispatcher uses to route an incoming interaction's Data.
func (s *responsesSuite) TestGetTypeDiscriminators() {
	s.Equal(discord.InteractionDataTypeApplicationCommand,
		(&InteractionDataApplicationCommand{}).GetType())
	s.Equal(discord.InteractionDataTypeApplicationCommandAutocomplete,
		(&InteractionDataAutocomplete{}).GetType())
	s.Equal(discord.InteractionDataTypeMessageComponent,
		(&InteractionDataMessageComponent{}).GetType())
	s.Equal(discord.InteractionDataTypeModalSubmit,
		(&InteractionDataModalSubmit{}).GetType())
}

// TestCommandDataUnmarshalIntegerOption verifies an integer option value decodes
// through the precision-preserving json.Number path.
func (s *responsesSuite) TestCommandDataUnmarshalIntegerOption() {
	const raw = `{"id":"1","name":"give","type":1,"options":[{"type":4,"name":"amount","value":42}]}`
	var d InteractionDataApplicationCommand
	s.Require().NoError(json.Unmarshal([]byte(raw), &d))
	s.Equal("give", d.Name)
	s.Require().NotNil(d.Options)
	s.Require().Len(d.Options, 1)
	s.Equal("amount", (d.Options)[0].Name)
}
