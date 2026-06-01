package discord

import (
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
	"time"
)

func (su *snowflakeSuite) TestSnowflakeFromInt() {
	t := su.T()
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

func (su *snowflakeSuite) TestSnowflakeFromUint() {
	t := su.T()
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

func (su *snowflakeSuite) TestSnowflakeFromInt_NegativeFailsValidate() {
	t := su.T()
	_, err := SnowflakeFromInt(-1)
	if err == nil {
		t.Error("expected Validate() to return an error for negative Snowflake, got nil")
	}
}

func (su *snowflakeSuite) TestSnowflakeTime() {
	t := su.T()
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

func (su *snowflakeSuite) TestUserCreatedAt() {
	t := su.T()
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
