package discord

import "encoding/json"

// Option represents a tri-state JSON value:
//
//   - Unset (zero value / nil map): field not provided → omitted by omitempty
//   - Set to null:                  map[false:zero]    → serializes as JSON null
//   - Set to value:                 map[true:v]        → serializes normally
//
// The map[bool]T trick gives us a type whose zero value (nil) is naturally
// skipped by encoding/json's omitempty on all Go versions, while a non-nil
// map is not skipped regardless of contents.
//
// Usage:
//
//	type Params struct {
//	    Nick      Option[string] `json:"nick,omitempty"`
//	    ChannelID Option[int]    `json:"channel_id,omitempty"`
//	}
//
//	Some("bob")    → {"nick":"bob"}
//	Null[string]() → {"nick":null}
//	None[string]() → field omitted entirely
type Option[T any] map[bool]T

// Some should be used for every non-null value, including zero values like "" and 0.
func Some[T any](v T) Option[T] {
	return Option[T]{true: v}
}

// Null should be used if the value is exactly null
func Null[T any]() Option[T] {
	var zero T
	return Option[T]{false: zero}
}

// None should be used if the value should not be set.
func None[T any]() Option[T] {
	return nil
}

func (o Option[T]) IsSet() bool { return o != nil }

func (o Option[T]) IsNull() bool {
	if o == nil {
		return false
	}
	_, hasValue := o[true]
	return !hasValue
}

func (o Option[T]) IsZero() bool { return o == nil }
func (o Option[T]) Val() (T, bool) {
	v, ok := o[true]
	return v, ok
}

func (o Option[T]) MustVal() T {
	v, ok := o[true]
	if !ok {
		return *new(T)
	}
	return v
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	v, ok := o[true]
	if !ok {
		// set but null
		return []byte("null"), nil
	}
	return json.Marshal(v)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = Null[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}
