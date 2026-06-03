package discord

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// jsonRoundtripSuite covers the pure marshal/unmarshal codecs and small helper
// methods on the discord value types — logic that needs no REST client.
type jsonRoundtripSuite struct {
	suite.Suite
}

func TestJSONRoundtripSuite(t *testing.T) { suite.Run(t, new(jsonRoundtripSuite)) }

// ── Snowflake ─────────────────────────────────────────────────────────────────

func (s *jsonRoundtripSuite) TestSnowflakeValidate() {
	valid := Snowflake(175928847299117063)
	s.True(valid.IsValid())
	s.NoError(valid.Validate())

	short := Snowflake(123)
	s.False(short.IsValid())
	s.Error(short.Validate())

	var nilSf *Snowflake
	s.Error(nilSf.Validate())
}

func (s *jsonRoundtripSuite) TestParseSnowflake() {
	sf, err := ParseSnowflake("175928847299117063")
	s.Require().NoError(err)
	s.Equal(Snowflake(175928847299117063), *sf)

	_, err = ParseSnowflake("not-a-number")
	s.Error(err)

	_, err = ParseSnowflake("123") // parses but fails validation (too short)
	s.Error(err)
}

// ── Permission ────────────────────────────────────────────────────────────────

func (s *jsonRoundtripSuite) TestPermissionCodec() {
	p := PermissionAdministrator | PermissionManageGuild
	b, err := json.Marshal(p)
	s.Require().NoError(err)
	s.Equal(`"40"`, string(b)) // 1<<3 | 1<<5 = 8 + 32 = 40, as a decimal string

	// Decode from the canonical string form.
	var fromStr Permission
	s.Require().NoError(json.Unmarshal([]byte(`"40"`), &fromStr))
	s.Equal(p, fromStr)

	// Decode from a raw integer (backwards-compat path).
	var fromInt Permission
	s.Require().NoError(json.Unmarshal([]byte(`40`), &fromInt))
	s.Equal(p, fromInt)

	// Invalid inputs error.
	var bad Permission
	s.Error(json.Unmarshal([]byte(`"oops"`), &bad))
	s.Error(json.Unmarshal([]byte(`"`), &bad))
}

func (s *jsonRoundtripSuite) TestNullFlagCodec() {
	var f NullFlag
	b, err := json.Marshal(f)
	s.Require().NoError(err)
	s.Equal("null", string(b))

	s.Require().NoError(json.Unmarshal([]byte("null"), &f))
	s.True(bool(f))
}

// ── Components ────────────────────────────────────────────────────────────────

func (s *jsonRoundtripSuite) TestComponentTypeIsAnySelectMenu() {
	s.True(ComponentTypeStringSelect.IsAnySelectMenu())
	s.True(ComponentTypeUserSelect.IsAnySelectMenu())
	s.True(ComponentTypeRoleSelect.IsAnySelectMenu())
	s.True(ComponentTypeMentionableSelect.IsAnySelectMenu())
	s.True(ComponentTypeChannelSelect.IsAnySelectMenu())
	s.False(ComponentTypeButton.IsAnySelectMenu())
}

func (s *jsonRoundtripSuite) TestComponentWrapperAndRawComponent() {
	raw := []byte(`{"type":2,"label":"Click"}`)

	var w ComponentWrapper
	s.Require().NoError(json.Unmarshal(raw, &w))
	s.Require().NotNil(w.Component)
	s.Equal(ComponentTypeButton, w.Component.GetType())

	out, err := w.MarshalJSON()
	s.Require().NoError(err)
	s.JSONEq(string(raw), string(out))

	// A nil component marshals to JSON null.
	empty := ComponentWrapper{}
	out, err = empty.MarshalJSON()
	s.Require().NoError(err)
	s.Equal("null", string(out))

	// RawComponent directly.
	var rc RawComponent
	s.Require().NoError(rc.UnmarshalJSON(raw))
	s.Equal(ComponentTypeButton, rc.GetType())
	rcOut, err := rc.MarshalJSON()
	s.Require().NoError(err)
	s.JSONEq(string(raw), string(rcOut))
}

// ── Guild wrappers ────────────────────────────────────────────────────────────

func (s *jsonRoundtripSuite) TestAnyGuildWrapper() {
	var avail AnyGuildWrapper
	s.Require().NoError(json.Unmarshal([]byte(`{"id":"175928847299117063","name":"g"}`), &avail))
	s.True(avail.Guild.IsAvailable())
	s.Equal(Snowflake(175928847299117063), avail.Guild.GetID())

	var unavail AnyGuildWrapper
	s.Require().NoError(json.Unmarshal([]byte(`{"id":"42","unavailable":true}`), &unavail))
	s.False(unavail.Guild.IsAvailable())
	s.Equal(Snowflake(42), unavail.Guild.GetID())

	s.Error(unavail.UnmarshalJSON([]byte(`{bad`)))
}

func (s *jsonRoundtripSuite) TestGatewayGuildWrapper() {
	var avail GatewayGuildWrapper
	s.Require().NoError(json.Unmarshal([]byte(`{"id":"175928847299117063","name":"g"}`), &avail))
	s.True(avail.Guild.IsAvailable())

	var unavail GatewayGuildWrapper
	s.Require().NoError(json.Unmarshal([]byte(`{"id":"7","unavailable":true}`), &unavail))
	s.False(unavail.Guild.IsAvailable())
	s.Equal(Snowflake(7), unavail.Guild.GetID())
}

func (s *jsonRoundtripSuite) TestUnavailableGuildIsAvailable() {
	// Unavailable flag absent → treated as unavailable.
	ug := UnavailableGuild{ID: 1}
	s.False(ug.IsAvailable())

	f := false
	ug.Unavailable = &f
	s.True(ug.IsAvailable())
}

func (s *jsonRoundtripSuite) TestGuildMFALevelToInt() {
	s.Equal(2, GuildMFALevel(2).ToInt())
}

// ── Activity ──────────────────────────────────────────────────────────────────

func (s *jsonRoundtripSuite) TestFullActivityCodec() {
	created := time.UnixMilli(1700000000000).UTC()
	start := time.UnixMilli(1700000001000).UTC()
	act := FullActivity{Name: "Coding", CreatedAt: created, Timestamps: &ActivityTimestamps{Start: &start}}

	b, err := json.Marshal(act)
	s.Require().NoError(err)

	var got FullActivity
	s.Require().NoError(json.Unmarshal(b, &got))
	s.Equal("Coding", got.Name)
	s.Equal(created, got.CreatedAt)
	s.Require().NotNil(got.Timestamps)
	s.Require().NotNil(got.Timestamps.Start)
	s.Equal(start, got.Timestamps.Start.UTC())

	s.Error(got.UnmarshalJSON([]byte(`{bad`)))
}

// ── PartialMessage ────────────────────────────────────────────────────────────

func (s *jsonRoundtripSuite) TestPartialMessageCodec() {
	raw := []byte(`{"id":"175928847299117063","content":"hi","components":[{"type":1}]}`)

	var pm PartialMessage
	s.Require().NoError(json.Unmarshal(raw, &pm))
	s.Equal("hi", pm.Content)
	s.Require().Len(pm.Components, 1)
	s.Equal(ComponentTypeActionRow, pm.Components[0].GetType())

	out, err := json.Marshal(&pm)
	s.Require().NoError(err)
	s.Contains(string(out), `"components"`)

	s.Error(pm.UnmarshalJSON([]byte(`{bad`)))
}
