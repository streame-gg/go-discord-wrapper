package discord

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func (s *snowflakeSuite) TestSnowflakeFromInt() {
	t := s.T()
	cases := []struct {
		name      string
		in        int64
		want      Snowflake
		wantError bool
	}{
		{"invalid snowflake length", 123, Snowflake(0), true},
		{"positive", 175928847299117063, Snowflake(175928847299117063), false},
		{"negative", -1, Snowflake(0), true},
		{"max int64", math.MaxInt64, Snowflake(9223372036854775807), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SnowflakeFromInt(tc.in)
			if err != nil && !tc.wantError {
				t.Errorf("SnowflakeFromInt(%v): %v", tc.in, err)
			} else if err == nil && tc.wantError {
				t.Errorf("SnowflakeFromInt(%v): expected error", tc.in)
			}
			if !tc.wantError && *got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func (s *snowflakeSuite) TestSnowflakeFromUint() {
	t := s.T()
	cases := []struct {
		name string
		in   uint64
		want Snowflake
	}{
		{"positive", 175928847299117063, Snowflake(175928847299117063)},
		{"zero", 0, Snowflake(0)},
		{"max uint64", math.MaxUint64, Snowflake(18446744073709551615)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SnowflakeFromUint(tc.in)
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func (s *snowflakeSuite) TestSnowflakeFromInt_NegativeFailsValidate() {
	t := s.T()
	_, err := SnowflakeFromInt(-1)
	if err == nil {
		t.Error("expected Validate() to return an error for negative Snowflake, got nil")
	}
}

func (s *snowflakeSuite) TestSnowflakeTime() {
	t := s.T()
	cases := []struct {
		name string
		id   Snowflake
		want time.Time
	}{
		{
			// Discord epoch itself: timestamp bits = 0 → 2015-01-01 00:00:00 UTC
			name: "discord epoch",
			id:   Snowflake(0),
			want: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			// Real snowflake: (754246997266923571 >> 22) + 1420070400000 = 1599896897380 ms → 2020-09-12 06:08:17.380 UTC
			name: "real snowflake",
			id:   Snowflake(754246997266923571),
			want: time.UnixMilli(1599896897380).UTC(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.id.Time()
			if !got.Equal(tc.want) {
				t.Errorf("Snowflake(%d).Time() = %v, want %v", tc.id, got, tc.want)
			}
		})
	}
}

func (s *snowflakeSuite) TestUserCreatedAt() {
	t := s.T()
	u := &User{ID: Snowflake(175928847299117063)}
	id := Snowflake(175928847299117063)
	want := id.Time()
	got := u.CreatedAt()
	if !got.Equal(want) {
		t.Errorf("User.CreatedAt() = %v, want %v", got, want)
	}
}

type snowflakeSuite struct{ suite.Suite }

func TestSnowflakeSuite(t *testing.T) { suite.Run(t, new(snowflakeSuite)) }

// snowflakeTimeReturner lets us call Time() on a function return value — a
// non-addressable value. This compiles only because Time has a value receiver
// (regression guard for the pointer-receiver bug).
func snowflakeTimeReturner() Snowflake { return Snowflake(175928847299117063) }

func (s *snowflakeSuite) TestTime_ValueReceiverOnNonAddressable() {
	// Composite-literal value (non-addressable) — would not compile with a
	// pointer receiver.
	s.False(Snowflake(175928847299117063).Time().IsZero())
	// Function-return value (non-addressable).
	s.False(snowflakeTimeReturner().Time().IsZero())

	// Sanity: a later snowflake decodes to a later embedded creation time.
	got := Snowflake(175928847299117063).Time()
	s.True(got.After(Snowflake(0).Time()), "later snowflake → later time")
}

func (s *snowflakeSuite) TestRandomSnowflake() {
	snowflake := RandomSnowflake()
	s.NotNil(snowflake)

	snowflakeArray := make([]Snowflake, 100)
	for i := 0; i < len(snowflakeArray); i++ {
		snowflakeArray[i] = RandomSnowflake()
	}
	s.Len(snowflakeArray, 100)
	s.NotContains(snowflakeArray, 0)
}
