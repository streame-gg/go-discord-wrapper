package connection

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// optPtr converts a pointer field into a discord.Option: a nil pointer becomes
// an unset Option, a non-nil pointer becomes a set Option holding the value.
func optPtr[T any](p *T) discord.Option[T] {
	if p == nil {
		return discord.None[T]()
	}
	return discord.Some(*p)
}

// optSlice converts a slice (or map) field into a discord.Option: a nil value
// becomes an unset Option, otherwise a set Option holding the value.
func optSlice[T any](s []T) discord.Option[[]T] {
	if s == nil {
		return discord.None[[]T]()
	}
	return discord.Some(s)
}

// optVal converts a value field into a discord.Option, treating the zero value
// as unset (matching the previous omitempty behaviour).
func optVal[T comparable](v T) discord.Option[T] {
	var zero T
	if v == zero {
		return discord.None[T]()
	}
	return discord.Some(v)
}
