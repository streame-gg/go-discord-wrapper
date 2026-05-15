package discord

import (
	"fmt"
)

type Snowflake string

func (s Snowflake) String() string {
	return string(s)
}

// ValidatedString automatically returns whether the Snowflake is valid, and if so returns it as a string.
func (s Snowflake) ValidatedString() (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	return string(s), nil
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

// ParseSnowflake parses and validates a string as a Discord Snowflake.
// Returns an error for any value that would be unsafe to embed in an API path.
func ParseSnowflake(s string) (Snowflake, error) {
	sf := Snowflake(s)
	if err := sf.Validate(); err != nil {
		return "", err
	}
	return sf, nil
}
