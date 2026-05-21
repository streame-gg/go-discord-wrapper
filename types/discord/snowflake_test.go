package discord

import (
	"math"
	"testing"
)

func TestSnowflakeFromInt(t *testing.T) {
	cases := []struct {
		name string
		in   int64
		want Snowflake
	}{
		{"positive", 175928847299117063, Snowflake(175928847299117063)},
		{"zero", 0, Snowflake(0)},
		{"max int64", math.MaxInt64, Snowflake(9223372036854775807)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SnowflakeFromInt(tc.in)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSnowflakeFromUint(t *testing.T) {
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
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSnowflakeFromInt_NegativeFailsValidate(t *testing.T) {
	sf := SnowflakeFromInt(-1)
	if err := sf.Validate(); err == nil {
		t.Error("expected Validate() to return an error for negative Snowflake, got nil")
	}
}
