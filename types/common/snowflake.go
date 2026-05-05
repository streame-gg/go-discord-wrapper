package common

type Snowflake string

func (s Snowflake) String() string {
	return string(s)
}
