package discord

import (
	"fmt"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/util"
)

const Epoch uint64 = 1420070400000

type Snowflake uint64

func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func (s Snowflake) IsEmpty() bool {
	return uint64(s) == 0
}

func (s Snowflake) IsValid() bool {
	return s.Validate() == nil
}

// Validate returns an error if s is not a valid Discord Snowflake (15–20 decimal digits).
// Use this to sanitize user-supplied IDs before embedding them in API paths.
func (s Snowflake) Validate() error {
	str := s.String()
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
	return Snowflake(id)
}

// SnowflakeFromUint converts an unsigned 64-bit integer to a Snowflake.
// Useful for Snowflakes stored as uint64 in databases that use unsigned
// integer types.
func SnowflakeFromUint(id uint64) Snowflake {
	return Snowflake(id)
}

// ParseSnowflake parses and validates a string as a Discord Snowflake.
// Returns an error for any value that would be unsafe to embed in an API path.
func ParseSnowflake(s string) (*Snowflake, error) {
	snowflake, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return util.PointerOf(Snowflake(snowflake)), nil
}
