package discord

import (
	"math"
	"testing"
)

func TestSnowflakeFromInt(t *testing.T) {
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
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestSnowflakeFromInt_NegativeFailsValidate(t *testing.T) {
	_, err := SnowflakeFromInt(-1)
	if err == nil {
		t.Error("expected Validate() to return an error for negative Snowflake, got nil")
	}
}
