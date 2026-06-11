package discord

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ModifyGuildMemberParams struct {
	Nick                       Option[string]      `json:"nick,omitempty"`
	Roles                      Option[[]Snowflake] `json:"roles,omitempty"`
	Mute                       Option[bool]        `json:"mute,omitempty"`
	Deaf                       Option[bool]        `json:"deaf,omitempty"`
	ChannelID                  Option[Snowflake]   `json:"channel_id,omitempty"`
	CommunicationDisabledUntil Option[string]      `json:"communication_disabled_until,omitempty"`
	Flags                      Option[int]         `json:"flags,omitempty"`
}

// ── suite ────────────────────────────────────────────────────────────────────

type OptionSuite struct {
	suite.Suite
}

func TestOptionSuite(t *testing.T) {
	suite.Run(t, new(OptionSuite))
}

// ── constructors ─────────────────────────────────────────────────────────────

func (s *OptionSuite) TestSome() {
	o := Some("hello")
	s.True(o.IsSet())
	s.False(o.IsNull())
	v, ok := o.Val()
	s.True(ok)
	s.Equal("hello", v)
}

func (s *OptionSuite) TestNull() {
	o := Null[string]()
	s.True(o.IsSet())
	s.True(o.IsNull())
	_, ok := o.Val()
	s.False(ok)
}

func (s *OptionSuite) TestNone() {
	o := None[string]()
	s.False(o.IsSet())
	s.False(o.IsNull())
	s.True(o.IsZero())
}

func (s *OptionSuite) TestZeroValue_BehavesLikeNone() {
	var o Option[string]
	s.False(o.IsSet())
	s.False(o.IsNull())
	s.True(o.IsZero())
}

func (s *OptionSuite) TestMustVal_ReturnsValue() {
	s.Equal(42, Some(42).MustVal())
}

func (s *OptionSuite) TestMustVal_PanicsOnNull() {
	s.Panics(func() { Null[string]().MustVal() })
}

func (s *OptionSuite) TestMustVal_PanicsOnUnset() {
	s.Panics(func() { None[string]().MustVal() })
}

// ── marshal ───────────────────────────────────────────────────────────────────

func (s *OptionSuite) marshal(v any) string {
	s.T().Helper()
	b, err := json.Marshal(v)
	s.Require().NoError(err)
	return string(b)
}

func (s *OptionSuite) TestMarshal_Some_SerializesValue() {
	type S struct {
		V Option[string] `json:"v,omitempty"`
	}
	s.Equal(`{"v":"hello"}`, s.marshal(S{V: Some("hello")}))
}

func (s *OptionSuite) TestMarshal_Null_SerializesNull() {
	type S struct {
		V Option[string] `json:"v,omitempty"`
	}
	s.Equal(`{"v":null}`, s.marshal(S{V: Null[string]()}))
}

func (s *OptionSuite) TestMarshal_None_IsOmitted() {
	type S struct {
		V Option[string] `json:"v,omitempty"`
	}
	s.Equal(`{}`, s.marshal(S{V: None[string]()}))
}

func (s *OptionSuite) TestMarshal_ZeroValue_IsOmitted() {
	type S struct {
		V Option[string] `json:"v,omitempty"`
	}
	s.Equal(`{}`, s.marshal(S{}))
}

func (s *OptionSuite) TestMarshal_BoolFalse_IsNotOmitted() {
	// false is a valid value — must not be omitted (this is the key failure
	// case for the naive *bool + omitempty approach)
	type S struct {
		Mute Option[bool] `json:"mute,omitempty"`
	}
	s.Equal(`{"mute":false}`, s.marshal(S{Mute: Some(false)}))
}

func (s *OptionSuite) TestMarshal_BoolUnset_IsOmitted() {
	type S struct {
		Mute Option[bool] `json:"mute,omitempty"`
	}
	s.Equal(`{}`, s.marshal(S{}))
}

func (s *OptionSuite) TestMarshal_OnlySetFieldsAreIncluded() {
	// Setting only Nick — all other fields should be absent from JSON
	params := ModifyGuildMemberParams{Nick: Some("bob")}
	s.Equal(`{"nick":"bob"}`, s.marshal(params))
}

func (s *OptionSuite) TestMarshal_NullNick_ClearsNickname() {
	params := ModifyGuildMemberParams{Nick: Null[string]()}
	s.Equal(`{"nick":null}`, s.marshal(params))
}

func (s *OptionSuite) TestMarshal_NullChannelID_DisconnectsFromVoice() {
	params := ModifyGuildMemberParams{ChannelID: Null[Snowflake]()}
	s.Equal(`{"channel_id":null}`, s.marshal(params))
}

func (s *OptionSuite) TestMarshal_AllFields() {
	params := ModifyGuildMemberParams{
		Nick:      Some("alice"),
		ChannelID: Some(Snowflake(1234567689345534)),
		Mute:      Some(true),
		Deaf:      Some(false),
	}
	s.Equal(`{"nick":"alice","mute":true,"deaf":false,"channel_id":1234567689345534}`, s.marshal(params))
}

func (s *OptionSuite) TestMarshalModifyGuildMemberParams() {
	params := ModifyGuildMemberParams{
		CommunicationDisabledUntil: None[string](),
	}

	s.Equal(`{"communication_disabled_until":null}`, s.marshal(params))
}

func (s *OptionSuite) TestMarshal_EmptyStruct_IsEmptyObject() {
	s.Equal(`{}`, s.marshal(ModifyGuildMemberParams{}))
}

// ── unmarshal ─────────────────────────────────────────────────────────────────

func (s *OptionSuite) unmarshal(str string) Option[string] {
	s.T().Helper()
	var o Option[string]
	s.Require().NoError(json.Unmarshal([]byte(str), &o))
	return o
}

func (s *OptionSuite) TestUnmarshal_Value() {
	o := s.unmarshal(`"hello"`)
	s.True(o.IsSet())
	s.False(o.IsNull())
	v, ok := o.Val()
	s.True(ok)
	s.Equal("hello", v)
}

func (s *OptionSuite) TestUnmarshal_Null() {
	o := s.unmarshal(`null`)
	s.True(o.IsSet())
	s.True(o.IsNull())
}

func (s *OptionSuite) TestUnmarshal_AbsentField_RemainsUnset() {
	// A field absent from the JSON object must remain as the zero value
	type S struct {
		Nick Option[string] `json:"nick,omitempty"`
		Mute Option[bool]   `json:"mute,omitempty"`
	}
	var parsed S
	s.Require().NoError(json.Unmarshal([]byte(`{"nick":"bob"}`), &parsed))
	s.True(parsed.Nick.IsSet(), "Nick should be set")
	s.False(parsed.Mute.IsSet(), "Mute should remain unset — absent from JSON")
}

func (s *OptionSuite) TestUnmarshal_RoundTrip() {
	type S struct {
		A Option[string] `json:"a,omitempty"`
		B Option[string] `json:"b,omitempty"`
		// C intentionally absent
	}
	original := S{
		A: Some("hello"),
		B: Null[string](),
	}

	data := s.marshal(original)

	var decoded S
	s.Require().NoError(json.Unmarshal([]byte(data), &decoded))

	v, ok := decoded.A.Val()
	s.True(ok)
	s.Equal("hello", v)
	s.True(decoded.B.IsNull(), "B should round-trip as null")
}
