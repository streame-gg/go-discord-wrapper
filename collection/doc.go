// Package collection provides a generic Collection[K, V] type — a map with
// insertion-order iteration and ~30 utility methods, inspired by
// @discordjs/collection. It is the foundational data structure for
// discord.js-style ergonomics in go-discord-wrapper.
//
// Collections are NOT thread-safe. Cache stores in this library return
// snapshots/clones, so callers don't need to lock; but if you share a
// Collection across goroutines, you must synchronize externally.
//
// Basic usage:
//
//	users := bot.Cache.Users().All()  // returns *Collection[Snowflake, *User]
//
//	// Find one
//	admin, ok := users.Find(func(u *User) bool { return u.Username == "admin" })
//
//	// Filter to a new Collection
//	humans := users.Filter(func(u *User) bool { return !u.Bot })
//
//	// Split in two
//	bots, others := users.Partition(func(u *User) bool { return u.Bot })
//
//	// Iterate (Go 1.23+ range-over-func)
//	for id, user := range users.All() {
//	    fmt.Printf("%s: %s\n", id, user.Username)
//	}
//
//	// Transform to a slice (top-level helper)
//	names := collection.Map(users, func(u *User) string { return u.Username })
//
// See the README's "Working with Collections" section for the full API tour.
package collection
