package common

type Snowflake string

func (s Snowflake) ToString() string {
	return string(s)
}
