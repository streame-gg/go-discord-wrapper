package common

type Snowflake string

func (s Snowflake) String() string {
	return string(s)
}

// OrEmpty returns the string value of the Snowflake pointer, or "" if the pointer is nil.
func (s *Snowflake) OrEmpty() string {
	if s == nil {
		return ""
	}
	return string(*s)
}
