package collection

import (
	"iter"
	"math/rand/v2"
	"slices"
)

// Entry is a single key-value pair, used in FromEntries and various return methods.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Collection is a generic Map<K, V> with utility methods inspired by
// @discordjs/collection. Unlike a plain map, Collection methods preserve
// insertion order for iteration when meaningful, and provide a rich set
// of functional helpers.
//
// Collection is NOT safe for concurrent use without external
// synchronization. The Discord cache stores wrap their own Collection
// instances with appropriate locking; users who hold Collection values
// returned from cache.All()-style methods receive an immutable snapshot
// and don't need their own locking.
type Collection[K comparable, V any] struct {
	items map[K]V
	keys  []K // insertion order, kept in sync with items
}

// New creates an empty Collection.
func New[K comparable, V any]() *Collection[K, V] {
	return &Collection[K, V]{items: make(map[K]V)}
}

// NewWithCapacity creates an empty Collection with hint for initial capacity.
func NewWithCapacity[K comparable, V any](capacity int) *Collection[K, V] {
	return &Collection[K, V]{
		items: make(map[K]V, capacity),
		keys:  make([]K, 0, capacity),
	}
}

// From creates a Collection from a map. The map is copied; subsequent
// modifications to the source map don't affect the Collection.
// Iteration order of the resulting Collection is unspecified (matches
// Go's map iteration).
func From[K comparable, V any](m map[K]V) *Collection[K, V] {
	c := NewWithCapacity[K, V](len(m))
	for k, v := range m {
		c.Set(k, v)
	}
	return c
}

// FromSlice creates a Collection from a slice of values, using the keyFn
// to extract a key from each value. Useful for typed entity slices:
//
//	users := []*discord.User{...}
//	coll := collection.FromSlice(users, func(u *discord.User) discord.Snowflake {
//	    return u.ID
//	})
func FromSlice[K comparable, V any](items []V, keyFn func(V) K) *Collection[K, V] {
	c := NewWithCapacity[K, V](len(items))
	for _, v := range items {
		c.Set(keyFn(v), v)
	}
	return c
}

// FromEntries creates a Collection from a slice of (key, value) pairs,
// preserving insertion order.
func FromEntries[K comparable, V any](entries []Entry[K, V]) *Collection[K, V] {
	c := NewWithCapacity[K, V](len(entries))
	for _, e := range entries {
		c.Set(e.Key, e.Value)
	}
	return c
}

// Get returns the value for key. The second return value is false if the
// key is not present.
func (c *Collection[K, V]) Get(key K) (V, bool) {
	v, ok := c.items[key]
	return v, ok
}

// GetOr returns the value for key, or fallback if the key is not present.
func (c *Collection[K, V]) GetOr(key K, fallback V) V {
	if v, ok := c.items[key]; ok {
		return v
	}
	return fallback
}

// Set inserts or replaces the value for key. Returns the Collection for
// chaining.
// If the key is new, it's appended to the insertion-order list.
// If the key exists, the value is replaced but order is preserved.
func (c *Collection[K, V]) Set(key K, value V) *Collection[K, V] {
	if _, exists := c.items[key]; !exists {
		c.keys = append(c.keys, key)
	}
	c.items[key] = value
	return c
}

// Has returns true if the key is present.
func (c *Collection[K, V]) Has(key K) bool {
	_, ok := c.items[key]
	return ok
}

// Delete removes the key-value pair. Returns true if the key was present.
func (c *Collection[K, V]) Delete(key K) bool {
	if _, ok := c.items[key]; !ok {
		return false
	}
	delete(c.items, key)
	for i, k := range c.keys {
		if k == key {
			c.keys = slices.Delete(c.keys, i, i+1)
			break
		}
	}
	return true
}

// Len returns the number of items in the Collection.
func (c *Collection[K, V]) Len() int {
	return len(c.items)
}

// IsEmpty returns true if Len() == 0.
func (c *Collection[K, V]) IsEmpty() bool {
	return len(c.items) == 0
}

// Clear removes all items.
func (c *Collection[K, V]) Clear() {
	c.items = make(map[K]V)
	c.keys = c.keys[:0]
}

// Clone returns a shallow copy. Values are copied by reference; modifying
// pointer-values affects both Collections.
func (c *Collection[K, V]) Clone() *Collection[K, V] {
	n := NewWithCapacity[K, V](len(c.items))
	copy(n.keys[:cap(n.keys)], c.keys)
	n.keys = n.keys[:len(c.keys)]
	for k, v := range c.items {
		n.items[k] = v
	}
	return n
}

// Find returns the first value for which fn returns true, or zero-value
// + false if none match. Iterates in insertion order.
func (c *Collection[K, V]) Find(fn func(V) bool) (V, bool) {
	for _, k := range c.keys {
		v := c.items[k]
		if fn(v) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// FindKey returns the first key for which fn(value) returns true, or
// zero-value + false.
func (c *Collection[K, V]) FindKey(fn func(V) bool) (K, bool) {
	for _, k := range c.keys {
		if fn(c.items[k]) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// FindLast returns the last value matching fn.
func (c *Collection[K, V]) FindLast(fn func(V) bool) (V, bool) {
	for i := len(c.keys) - 1; i >= 0; i-- {
		v := c.items[c.keys[i]]
		if fn(v) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// FindLastKey returns the last key matching fn(value).
func (c *Collection[K, V]) FindLastKey(fn func(V) bool) (K, bool) {
	for i := len(c.keys) - 1; i >= 0; i-- {
		k := c.keys[i]
		if fn(c.items[k]) {
			return k, true
		}
	}
	var zero K
	return zero, false
}

// Some returns true if fn returns true for any value.
func (c *Collection[K, V]) Some(fn func(V) bool) bool {
	for _, k := range c.keys {
		if fn(c.items[k]) {
			return true
		}
	}
	return false
}

// Every returns true if fn returns true for all values. Returns true for
// empty Collections (vacuous truth, matches JS).
func (c *Collection[K, V]) Every(fn func(V) bool) bool {
	for _, k := range c.keys {
		if !fn(c.items[k]) {
			return false
		}
	}
	return true
}

// Equals returns true if other contains the same keys and values according to eq.
// Callers supply the equality function so that Equals is safe for any V,
// including non-comparable types such as structs containing slices or maps.
//
// Example (pointer identity):
//
//	a.Equals(b, func(x, y *T) bool { return x == y })
//
// Example (deep equality):
//
//	a.Equals(b, func(x, y MyStruct) bool { return reflect.DeepEqual(x, y) })
func (c *Collection[K, V]) Equals(other *Collection[K, V], eq func(V, V) bool) bool {
	if len(c.items) != len(other.items) {
		return false
	}
	for k, v := range c.items {
		ov, ok := other.items[k]
		if !ok || !eq(v, ov) {
			return false
		}
	}
	return true
}

// Filter returns a new Collection containing only items where fn returns true.
// Insertion order is preserved.
func (c *Collection[K, V]) Filter(fn func(V) bool) *Collection[K, V] {
	n := New[K, V]()
	for _, k := range c.keys {
		v := c.items[k]
		if fn(v) {
			n.Set(k, v)
		}
	}
	return n
}

// Partition splits the Collection into two: matched (fn returns true) and
// rest (fn returns false). Insertion order is preserved within each.
//
//	bots, humans := users.Partition(func(u *User) bool { return u.Bot })
func (c *Collection[K, V]) Partition(fn func(V) bool) (matched, rest *Collection[K, V]) {
	matched = New[K, V]()
	rest = New[K, V]()
	for _, k := range c.keys {
		v := c.items[k]
		if fn(v) {
			matched.Set(k, v)
		} else {
			rest.Set(k, v)
		}
	}
	return matched, rest
}

// First returns the first value in insertion order, or zero-value + false
// if empty.
func (c *Collection[K, V]) First() (V, bool) {
	if len(c.keys) == 0 {
		var zero V
		return zero, false
	}
	return c.items[c.keys[0]], true
}

// FirstN returns up to n first values in insertion order.
// Returns nil if n <= 0.
func (c *Collection[K, V]) FirstN(n int) []V {
	if n <= 0 {
		return nil
	}
	if n > len(c.keys) {
		n = len(c.keys)
	}
	result := make([]V, n)
	for i := range n {
		result[i] = c.items[c.keys[i]]
	}
	return result
}

// FirstKey returns the first key.
func (c *Collection[K, V]) FirstKey() (K, bool) {
	if len(c.keys) == 0 {
		var zero K
		return zero, false
	}
	return c.keys[0], true
}

// Last returns the last value in insertion order.
func (c *Collection[K, V]) Last() (V, bool) {
	if len(c.keys) == 0 {
		var zero V
		return zero, false
	}
	return c.items[c.keys[len(c.keys)-1]], true
}

// LastN returns up to n last values in insertion order.
// Returns nil if n <= 0.
func (c *Collection[K, V]) LastN(n int) []V {
	if n <= 0 {
		return nil
	}
	total := len(c.keys)
	if n > total {
		n = total
	}
	result := make([]V, n)
	start := total - n
	for i := range n {
		result[i] = c.items[c.keys[start+i]]
	}
	return result
}

// LastKey returns the last key.
func (c *Collection[K, V]) LastKey() (K, bool) {
	if len(c.keys) == 0 {
		var zero K
		return zero, false
	}
	return c.keys[len(c.keys)-1], true
}

// At returns the value at the given insertion-order index. Negative indices
// count from the end (At(-1) == Last()). Returns zero-value + false if
// the index is out of bounds.
func (c *Collection[K, V]) At(index int) (V, bool) {
	n := len(c.keys)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		var zero V
		return zero, false
	}
	return c.items[c.keys[index]], true
}

// KeyAt returns the key at the given insertion-order index. Negative indices
// count from the end. Returns zero-value + false if out of bounds.
func (c *Collection[K, V]) KeyAt(index int) (K, bool) {
	n := len(c.keys)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		var zero K
		return zero, false
	}
	return c.keys[index], true
}

// Random returns a random value, or zero-value + false if empty.
// Uses math/rand/v2.
func (c *Collection[K, V]) Random() (V, bool) {
	if len(c.keys) == 0 {
		var zero V
		return zero, false
	}
	return c.items[c.keys[rand.IntN(len(c.keys))]], true
}

// RandomN returns up to n distinct random values.
// Returns nil if n <= 0.
func (c *Collection[K, V]) RandomN(n int) []V {
	if n <= 0 {
		return nil
	}
	total := len(c.keys)
	if n > total {
		n = total
	}
	shuffled := make([]K, total)
	copy(shuffled, c.keys)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	result := make([]V, n)
	for i := range n {
		result[i] = c.items[shuffled[i]]
	}
	return result
}

// RandomKey returns a random key, or zero-value + false if empty.
func (c *Collection[K, V]) RandomKey() (K, bool) {
	if len(c.keys) == 0 {
		var zero K
		return zero, false
	}
	return c.keys[rand.IntN(len(c.keys))], true
}

// Each iterates over all (key, value) pairs in insertion order, calling
// fn. Returns the Collection for chaining.
//
//	users.Each(func(id Snowflake, u *User) { fmt.Println(u.Name) })
func (c *Collection[K, V]) Each(fn func(K, V)) *Collection[K, V] {
	for _, k := range c.keys {
		fn(k, c.items[k])
	}
	return c
}

// Sort sorts the Collection in-place by the less function. Sort is stable.
// Returns the Collection for chaining.
func (c *Collection[K, V]) Sort(less func(a, b V) bool) *Collection[K, V] {
	items := c.items
	slices.SortStableFunc(c.keys, func(a, b K) int {
		va, vb := items[a], items[b]
		if less(va, vb) {
			return -1
		}
		if less(vb, va) {
			return 1
		}
		return 0
	})
	return c
}

// Sorted returns a sorted copy without modifying the original.
func (c *Collection[K, V]) Sorted(less func(a, b V) bool) *Collection[K, V] {
	return c.Clone().Sort(less)
}

// Reverse reverses the insertion order in place. Returns the Collection.
func (c *Collection[K, V]) Reverse() *Collection[K, V] {
	slices.Reverse(c.keys)
	return c
}

// MapValues transforms values in-place. For type-changing transforms, use
// the top-level Map function.
func (c *Collection[K, V]) MapValues(fn func(V) V) *Collection[K, V] {
	for k, v := range c.items {
		c.items[k] = fn(v)
	}
	return c
}

// Sweep removes items where fn returns true. Returns the number of items
// removed. This is the inverse of Filter — modifies in place.
func (c *Collection[K, V]) Sweep(fn func(V) bool) int {
	removed := 0
	newKeys := make([]K, 0, len(c.keys))
	for _, k := range c.keys {
		v := c.items[k]
		if fn(v) {
			delete(c.items, k)
			removed++
		} else {
			newKeys = append(newKeys, k)
		}
	}
	c.keys = newKeys
	return removed
}

// Concat returns a new Collection containing all items from c and others.
// On duplicate keys, later collections override earlier.
func (c *Collection[K, V]) Concat(others ...*Collection[K, V]) *Collection[K, V] {
	n := c.Clone()
	for _, other := range others {
		for _, k := range other.keys {
			n.Set(k, other.items[k])
		}
	}
	return n
}

// Merge is an alias for Concat for naming familiarity with discord.js.
func (c *Collection[K, V]) Merge(others ...*Collection[K, V]) *Collection[K, V] {
	return c.Concat(others...)
}

// Difference returns items in c whose keys are NOT in other.
func (c *Collection[K, V]) Difference(other *Collection[K, V]) *Collection[K, V] {
	n := New[K, V]()
	for _, k := range c.keys {
		if !other.Has(k) {
			n.Set(k, c.items[k])
		}
	}
	return n
}

// Intersection returns items present in both c and other (by key, with
// values from c).
func (c *Collection[K, V]) Intersection(other *Collection[K, V]) *Collection[K, V] {
	n := New[K, V]()
	for _, k := range c.keys {
		if other.Has(k) {
			n.Set(k, c.items[k])
		}
	}
	return n
}

// SymmetricDifference returns items present in c or other but not both.
func (c *Collection[K, V]) SymmetricDifference(other *Collection[K, V]) *Collection[K, V] {
	n := New[K, V]()
	for _, k := range c.keys {
		if !other.Has(k) {
			n.Set(k, c.items[k])
		}
	}
	for _, k := range other.keys {
		if !c.Has(k) {
			n.Set(k, other.items[k])
		}
	}
	return n
}

// Keys returns a snapshot slice of all keys in insertion order.
func (c *Collection[K, V]) Keys() []K {
	result := make([]K, len(c.keys))
	copy(result, c.keys)
	return result
}

// Values returns a snapshot slice of all values in insertion order.
func (c *Collection[K, V]) Values() []V {
	result := make([]V, len(c.keys))
	for i, k := range c.keys {
		result[i] = c.items[k]
	}
	return result
}

// Entries returns a snapshot slice of all key-value pairs in insertion order.
func (c *Collection[K, V]) Entries() []Entry[K, V] {
	result := make([]Entry[K, V], len(c.keys))
	for i, k := range c.keys {
		result[i] = Entry[K, V]{Key: k, Value: c.items[k]}
	}
	return result
}

// ToMap returns a copy of the underlying map. Mutations to the returned
// map don't affect the Collection.
func (c *Collection[K, V]) ToMap() map[K]V {
	m := make(map[K]V, len(c.items))
	for k, v := range c.items {
		m[k] = v
	}
	return m
}

// All returns iter.Seq2[K, V] for range-over-func iteration (Go 1.23+).
//
//	for k, v := range users.All() { ... }
func (c *Collection[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range c.keys {
			if !yield(k, c.items[k]) {
				return
			}
		}
	}
}

// KeysIter returns iter.Seq[K] for range-over-func on keys only.
func (c *Collection[K, V]) KeysIter() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, k := range c.keys {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesIter returns iter.Seq[V] for range-over-func on values only.
func (c *Collection[K, V]) ValuesIter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, k := range c.keys {
			if !yield(c.items[k]) {
				return
			}
		}
	}
}

// Map transforms a Collection[K, V] into a slice of R using fn.
//
//	names := collection.Map(users, func(u *User) string { return u.Username })
func Map[K comparable, V any, R any](c *Collection[K, V], fn func(V) R) []R {
	result := make([]R, len(c.keys))
	for i, k := range c.keys {
		result[i] = fn(c.items[k])
	}
	return result
}

// MapToCollection transforms a Collection[K, V] into a Collection[K, R]
// preserving keys and insertion order.
func MapToCollection[K comparable, V any, R any](c *Collection[K, V], fn func(V) R) *Collection[K, R] {
	n := NewWithCapacity[K, R](len(c.keys))
	for _, k := range c.keys {
		n.Set(k, fn(c.items[k]))
	}
	return n
}

// FlatMap applies fn to each value, then concatenates results into a
// single slice.
func FlatMap[K comparable, V any, R any](c *Collection[K, V], fn func(V) []R) []R {
	var result []R
	for _, k := range c.keys {
		result = append(result, fn(c.items[k])...)
	}
	return result
}

// GroupBy groups values by a derived key.
//
//	byBot := collection.GroupBy(users, func(u *User) bool { return u.Bot })
//	// → map[bool]*Collection[Snowflake, *User]
func GroupBy[K comparable, V any, G comparable](c *Collection[K, V], keyFn func(V) G) map[G]*Collection[K, V] {
	result := make(map[G]*Collection[K, V])
	for _, k := range c.keys {
		v := c.items[k]
		g := keyFn(v)
		if _, ok := result[g]; !ok {
			result[g] = New[K, V]()
		}
		result[g].Set(k, v)
	}
	return result
}

// Reduce accumulates a single value from all items, in insertion order.
//
//	total := collection.Reduce(users, 0, func(acc int, u *User) int {
//	    return acc + u.MessageCount
//	})
func Reduce[K comparable, V any, R any](c *Collection[K, V], initial R, fn func(R, V) R) R {
	acc := initial
	for _, k := range c.keys {
		acc = fn(acc, c.items[k])
	}
	return acc
}
