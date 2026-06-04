package interactions_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// interactionsSuite covers Interaction.UnmarshalJSON's polymorphic dispatch
// (which concrete Data type it picks per interaction/component type) and the
// command-routing helpers built on top of it (GetSubCommand, GetSubCommandGroup,
// GetFullCommand, GetCustomID).
type interactionsSuite struct {
	suite.Suite
}

func TestInteractionsSuite(t *testing.T) {
	suite.Run(t, new(interactionsSuite))
}

func (s *interactionsSuite) unmarshal(raw string) *interactions.Interaction {
	var i interactions.Interaction
	s.Require().NoError(json.Unmarshal([]byte(raw), &i))
	return &i
}

// TestUnmarshalChatInputCommand decodes a nested subcommand-group → subcommand
// payload and verifies both the chosen Data type and the routing helpers.
func (s *interactionsSuite) TestUnmarshalChatInputCommand() {
	const raw = `{
		"type": 2,
		"data": {
			"id": "100000000000000001",
			"name": "config",
			"type": 1,
			"options": [
				{"type": 2, "name": "channel", "options": [
					{"type": 1, "name": "set", "options": [
						{"type": 4, "name": "amount", "value": 5}
					]}
				]}
			]
		}
	}`
	i := s.unmarshal(raw)

	cmd, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	s.Require().True(ok, "data should decode to *InteractionDataApplicationCommand")
	s.Equal("config", cmd.CommandName)

	s.Equal("channel", i.GetSubCommandGroup())
	s.Equal("set", i.GetSubCommand())
	s.Equal("config channel set", i.GetFullCommand())
	s.Empty(i.GetCustomID(), "a command interaction has no custom ID")
}

// TestUnmarshalTopLevelSubCommand covers a command with a direct subcommand
// (no group), so GetSubCommandGroup is empty but GetSubCommand resolves.
func (s *interactionsSuite) TestUnmarshalTopLevelSubCommand() {
	const raw = `{
		"type": 2,
		"data": {"id":"1","name":"ping","type":1,"options":[{"type":1,"name":"now"}]}
	}`
	i := s.unmarshal(raw)
	s.Empty(i.GetSubCommandGroup())
	s.Equal("now", i.GetSubCommand())
	s.Equal("ping now", i.GetFullCommand())
}

func (s *interactionsSuite) TestUnmarshalMessageComponent() {
	const raw = `{"type":3,"data":{"custom_id":"btn_confirm","component_type":2}}`
	i := s.unmarshal(raw)

	_, ok := i.Data.(*responses.InteractionDataMessageComponent)
	s.Require().True(ok, "data should decode to *InteractionDataMessageComponent")
	s.Equal("btn_confirm", i.GetCustomID())
	s.Empty(i.GetSubCommand(), "components have no subcommands")
	s.Empty(i.GetFullCommand())
}

func (s *interactionsSuite) TestUnmarshalModalSubmit() {
	const raw = `{"type":5,"data":{"custom_id":"feedback_modal","components":[]}}`
	i := s.unmarshal(raw)

	_, ok := i.Data.(*responses.InteractionDataModalSubmit)
	s.Require().True(ok, "data should decode to *InteractionDataModalSubmit")
	s.Equal("feedback_modal", i.GetCustomID())
}

func (s *interactionsSuite) TestUnmarshalAutocomplete() {
	const raw = `{
		"type": 4,
		"data": {"name":"search","type":1,"options":[{"type":3,"name":"q","focused":true}]}
	}`
	i := s.unmarshal(raw)

	_, ok := i.Data.(*responses.InteractionDataAutocomplete)
	s.Require().True(ok, "autocomplete dispatch keys off the interaction type, not the data type")
}

func (s *interactionsSuite) TestUnmarshalUnknownDataTypeErrors() {
	const raw = `{"type":2,"data":{"type":99}}`
	var i interactions.Interaction
	err := json.Unmarshal([]byte(raw), &i)
	s.Require().Error(err, "an unrecognised data type must surface an error")
	s.Contains(err.Error(), "unknown interaction data type")
}

func (s *interactionsSuite) TestUnmarshalNoData() {
	// PING interactions carry no data payload; UnmarshalJSON must accept that.
	const raw = `{"type":1}`
	i := s.unmarshal(raw)
	s.Nil(i.Data)
}

// TestStringOptionValue decodes a String option whose value is a JSON string
// (the common case, e.g. {"type":3,"value":"hello"}). The option decoder reads
// the raw value as json.RawMessage and interprets it by option type, so string
// and boolean options decode correctly.
func (s *interactionsSuite) TestStringOptionValue() {
	const raw = `{
		"type": 2,
		"data": {"id":"1","name":"echo","type":1,"options":[{"type":3,"name":"text","value":"hello"}]}
	}`
	var i interactions.Interaction
	s.Require().NoError(json.Unmarshal([]byte(raw), &i))
	cmd, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	s.Require().True(ok)
	s.Equal("echo", cmd.CommandName)
	s.Require().NotNil(cmd.Options)
	opts := *cmd.Options
	s.Require().Len(opts, 1)
	s.Require().NotNil(opts[0].Value)
	s.Equal("hello", *opts[0].Value)
}

// TestHelpersOnEmptyInteraction verifies the helpers are nil-safe when Data is
// absent (e.g. a PING interaction or a zero value).
func (s *interactionsSuite) TestHelpersOnEmptyInteraction() {
	var i interactions.Interaction
	s.Empty(i.GetSubCommand())
	s.Empty(i.GetSubCommandGroup())
	s.Empty(i.GetFullCommand())
	s.Empty(i.GetCustomID())
}
