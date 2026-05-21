package discord

import (
	"fmt"
	"strconv"
)

const Epoch int64 = 1420070400000

type Snowflake string

func (s Snowflake) String() string {
	return string(s)
}

// Validate returns an error if s is not a valid Discord Snowflake (15–20 decimal digits).
// Use this to sanitize user-supplied IDs before embedding them in API paths.
func (s Snowflake) Validate() error {
	str := string(s)
	if len(str) < 15 || len(str) > 20 {
		return fmt.Errorf("snowflake %q: must be 15–20 decimal digits", str)
	}
	for _, c := range str {
		if c < '0' || c > '9' {
			return fmt.Errorf("snowflake %q: contains non-digit character %q", str, c)
		}
	}
	return nil
}

// SnowflakeFromInt converts a signed 64-bit integer to a Snowflake.
// Useful for Snowflakes stored as int64 in databases.
//
// Negative values are technically invalid Discord Snowflakes — the
// helper accepts them silently and the resulting Snowflake will fail
// Validate(). If you want validation, call .Validate() on the result.
func SnowflakeFromInt(id int64) Snowflake {
	return Snowflake(strconv.FormatInt(id, 10))
}

// SnowflakeFromUint converts an unsigned 64-bit integer to a Snowflake.
// Useful for Snowflakes stored as uint64 in databases that use unsigned
// integer types.
func SnowflakeFromUint(id uint64) Snowflake {
	return Snowflake(strconv.FormatUint(id, 10))
}

// SnowflakeFromString is an explicit constructor from a string. It is
// equivalent to a direct conversion `Snowflake(s)` but reads more clearly
// in code that switches between multiple ID representations.
//
// It does NOT validate — use s.Validate() if you need that.
func SnowflakeFromString(s string) Snowflake {
	return Snowflake(s)
}

// ParseSnowflake parses and validates a string as a Discord Snowflake.
// Returns an error for any value that would be unsafe to embed in an API path.
func ParseSnowflake(s string) (Snowflake, error) {
	sf := Snowflake(s)
	if err := sf.Validate(); err != nil {
		return "", err
	}
	return sf, nil
}
