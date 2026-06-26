package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"
)

const Epoch uint64 = 1420070400000

// https://docs.discord.com/developers/reference#snowflakes
type Snowflake uint64

func (s *Snowflake) IsNil() bool {
	return s == nil || *s == 0
}

func (s *Snowflake) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*s = Snowflake(v)
		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*s = Snowflake(v)
	return nil
}

// Time returns the creation time encoded in the Snowflake.
// Uses a value receiver so it can be called on non-addressable values
// (e.g. discord.Snowflake(0).Time(), range-loop variables, function returns),
// mirroring String; it does not mutate the value.
func (s Snowflake) Time() time.Time {
	ms := (uint64(s) >> 22) + Epoch
	return time.UnixMilli(int64(ms)).UTC()
}

func (s *Snowflake) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, errors.New("cannot marshal Snowflake as nil")
	}
	return json.Marshal(strconv.FormatUint(uint64(*s), 10))
}

// String returns the decimal string representation of the Snowflake.
// Uses a value receiver so it can be called on non-addressable values
// (e.g. range-loop variables, function return values).
func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func (s *Snowflake) IsEmpty() bool {
	if s == nil {
		return true
	}

	return uint64(*s) == 0
}

func (s *Snowflake) IsValid() bool {
	return s.Validate() == nil
}

// Validate returns an error if s is not a valid Discord Snowflake (15–20 decimal digits).
// Use this to sanitize user-supplied IDs before embedding them in API paths.
func (s *Snowflake) Validate() error {
	if s == nil {
		return errors.New("nil snowflake")
	}

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
// Returns an error if out of range or otherwise invalid.
func SnowflakeFromInt(id int64) (*Snowflake, error) {
	if id < 0 {
		return nil, fmt.Errorf("snowflake %d: id out of range", id)
	}

	snowflake := Snowflake(id)

	if err := snowflake.Validate(); err != nil {
		return nil, err
	}

	return &snowflake, nil
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

	createdSnowflake := Snowflake(snowflake)
	if err := createdSnowflake.Validate(); err != nil {
		return nil, err
	}

	return &createdSnowflake, nil
}

func RandomSnowflake() Snowflake {
	now := time.Now().UnixMilli()
	timestamp := uint64(now) - Epoch

	workerID := uint64(rand.IntN(32))
	processID := uint64(rand.IntN(32))
	increment := uint64(rand.IntN(4096))

	snowflake := (timestamp << 22) |
		(workerID << 17) |
		(processID << 12) |
		increment

	return Snowflake(snowflake)
}
