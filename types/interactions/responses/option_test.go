package responses

import (
	"encoding/json"
)

// ── option value decoding ─────────────────────────────────────────────────────

// TestOptionNumber decodes a floating-point (number) option through the
// json.Number path.
func (s *responsesSuite) TestOptionNumber() {
	var o ApplicationCommandInteractionDataOption[interface{}]
	s.Require().NoError(json.Unmarshal([]byte(`{"type":10,"name":"ratio","value":3.5}`), &o))
	s.Require().NotNil(o.Value)
	s.Equal(3.5, *o.Value)
}

// TestOptionInteger decodes an integer option without precision loss.
func (s *responsesSuite) TestOptionInteger() {
	var o ApplicationCommandInteractionDataOption[interface{}]
	s.Require().NoError(json.Unmarshal([]byte(`{"type":4,"name":"amount","value":9007199254740993}`), &o))
	s.Require().NotNil(o.Value)
	s.Equal(9007199254740993, *o.Value)
}

// TestOptionIntegerBadValue returns an error when an integer option carries a
// non-integer numeric value.
func (s *responsesSuite) TestOptionIntegerBadValue() {
	var o ApplicationCommandInteractionDataOption[interface{}]
	err := json.Unmarshal([]byte(`{"type":4,"name":"amount","value":3.14}`), &o)
	s.Require().Error(err)
	s.Contains(err.Error(), "integer value")
}

// TestOptionNoValue covers a subcommand-style option that carries nested
// options instead of a value (the raw.Value == "" branch).
func (s *responsesSuite) TestOptionNoValue() {
	var o ApplicationCommandInteractionDataOption[interface{}]
	s.Require().NoError(json.Unmarshal([]byte(`{"type":1,"name":"sub","options":[{"type":4,"name":"x","value":1}]}`), &o))
	s.Nil(o.Value)
	s.Require().Len(o.Options, 1)
	s.Equal("x", o.Options[0].Name)
}

// TestOptionStringValueIsRejected pins the current behaviour: because the value
// field is decoded as json.Number, a string-valued option (the common case for
// String/User/Channel/Role options) currently fails to decode. This is a known
// latent bug; the test documents it so a future fix is a deliberate change.
func (s *responsesSuite) TestOptionStringValueIsRejected() {
	var o ApplicationCommandInteractionDataOption[interface{}]
	err := json.Unmarshal([]byte(`{"type":3,"name":"text","value":"hello"}`), &o)
	s.Require().Error(err, "string option values currently fail to decode (json.Number path)")
}

// ── command data unmarshalling error paths ────────────────────────────────────

// TestCommandUnmarshalBadJSON propagates a malformed-envelope error.
func (s *responsesSuite) TestCommandUnmarshalBadJSON() {
	var d InteractionDataApplicationCommand
	s.Require().Error(json.Unmarshal([]byte(`{`), &d))
}

// TestCommandUnmarshalBadOption propagates an error raised while decoding one of
// the nested options.
func (s *responsesSuite) TestCommandUnmarshalBadOption() {
	const raw = `{"id":"1","name":"cmd","type":1,"options":[{"type":4,"name":"x","value":1.5}]}`
	var d InteractionDataApplicationCommand
	s.Require().Error(json.Unmarshal([]byte(raw), &d))
}

// TestCommandUnmarshalNoOptions decodes a command with no options at all (the
// raw.Options == nil branch).
func (s *responsesSuite) TestCommandUnmarshalNoOptions() {
	const raw = `{"id":"1","name":"ping","type":1}`
	var d InteractionDataApplicationCommand
	s.Require().NoError(json.Unmarshal([]byte(raw), &d))
	s.Equal("ping", d.CommandName)
	s.Nil(d.Options)
}
