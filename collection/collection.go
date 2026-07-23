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

// node is an element in the doubly-linked insertion-order list.
type node[K comparable, V any] struct {
	key        K
	value      V
	prev, next *node[K, V]
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
	items map[K]*node[K, V]
	head  *node[K, V]
	tail  *node[K, V]
}

// New creates an empty Collection.
func New[K comparable, V any]() *Collection[K, V] {
	return &Collection[K, V]{items: make(map[K]*node[K, V])}
}

// NewWithCapacity creates an empty Collection with hint for initial capacity.
func NewWithCapacity[K comparable, V any](capacity int) *Collection[K, V] {
	return &Collection[K, V]{items: make(map[K]*node[K, V], capacity)}
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
func (c *Collection[K, V]) Get(key K) (*V, bool) {
	if c == nil {
		return nil, false
	}
	if n, ok := c.items[key]; ok {
		return &n.value, true
	}
	return nil, false
}

// GetOr returns the value for key, or fallback if the key is not present.
func (c *Collection[K, V]) GetOr(key K, fallback V) V {
	if c == nil {
		return fallback
	}
	if n, ok := c.items[key]; ok {
		return n.value
	}
	return fallback
}

// Set inserts or replaces the value for key. Returns the Collection for
// chaining.
// If the key is new, it's appended to the insertion-order list.
// If the key exists, the value is replaced but order is preserved.
func (c *Collection[K, V]) Set(key K, value V) *Collection[K, V] {
	if c == nil {
		coll := New[K, V]()
		coll.Set(key, value)
		return coll
	}
	if n, ok := c.items[key]; ok {
		n.value = value
		return c
	}

	n := &node[K, V]{key: key, value: value, prev: c.tail}
	if c.tail != nil {
		c.tail.next = n
	} else {
		c.head = n
	}

	c.tail = n
	c.items[key] = n

	return c
}

// Has returns true if the key is present.
func (c *Collection[K, V]) Has(key K) bool {
	if c == nil {
		return false
	}
	_, ok := c.items[key]
	return ok
}

// Delete removes the key-value pair. Returns true if the key was present.
// O(1).
func (c *Collection[K, V]) Delete(key K) bool {
	if c == nil {
		return false
	}
	n, ok := c.items[key]
	if !ok {
		return false
	}

	c.unlink(n)
	delete(c.items, key)

	return true
}

// unlink removes n from the doubly-linked list without touching c.items.
func (c *Collection[K, V]) unlink(n *node[K, V]) {
	if c == nil {
		return
	}
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}

	n.prev = nil
	n.next = nil
}

// Len returns the number of items in the Collection.
func (c *Collection[K, V]) Len() int {
	if c == nil {
		return 0
	}
	return len(c.items)
}

// IsEmpty returns true if Len() == 0.
func (c *Collection[K, V]) IsEmpty() bool {
	if c == nil {
		return true
	}
	return len(c.items) == 0
}

// Clear removes all items.
func (c *Collection[K, V]) Clear() {
	if c == nil {
		return
	}
	c.items = make(map[K]*node[K, V])
	c.head = nil
	c.tail = nil
}

// Clone returns a shallow copy. Values are copied by reference; modifying
// pointer-values affects both Collections.
func (c *Collection[K, V]) Clone() *Collection[K, V] {
	if c == nil {
		return New[K, V]()
	}
	n := NewWithCapacity[K, V](len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		n.Set(cur.key, cur.value)
	}

	return n
}

// Find returns the first value for which fn returns true, or zero-value
// + false if none match. Iterates in insertion order.
func (c *Collection[K, V]) Find(fn func(V) bool) (*V, bool) {
	if c == nil {
		return nil, false
	}

	for cur := c.head; cur != nil; cur = cur.next {
		if fn(cur.value) {
			return &cur.value, true
		}
	}

	return nil, false
}

// FindKey returns the first key for which fn(value) returns true
func (c *Collection[K, V]) FindKey(fn func(V) bool) (*K, bool) {
	if c == nil {
		return nil, false
	}

	for cur := c.head; cur != nil; cur = cur.next {
		if fn(cur.value) {
			return &cur.key, true
		}
	}

	return nil, false
}

// FindLast returns the last value matching fn.
func (c *Collection[K, V]) FindLast(fn func(V) bool) (*V, bool) {
	if c == nil {
		return nil, false
	}

	for cur := c.tail; cur != nil; cur = cur.prev {
		if fn(cur.value) {
			return &cur.value, true
		}
	}

	return nil, false
}

// FindLastKey returns the last key matching fn(value).
func (c *Collection[K, V]) FindLastKey(fn func(V) bool) (*K, bool) {
	if c == nil {
		return nil, false
	}

	for cur := c.tail; cur != nil; cur = cur.prev {
		if fn(cur.value) {
			return &cur.key, true
		}
	}

	return nil, false
}

// Some returns true if fn returns true for any value.
func (c *Collection[K, V]) Some(fn func(V) bool) bool {
	if c == nil {
		return false
	}

	for cur := c.head; cur != nil; cur = cur.next {
		if fn(cur.value) {
			return true
		}
	}

	return false
}

// Every returns true if fn returns true for all values. Returns true for
// empty Collections (vacuous truth, matches JS).
func (c *Collection[K, V]) Every(fn func(V) bool) bool {
	if c == nil {
		return false
	}

	for cur := c.head; cur != nil; cur = cur.next {
		if !fn(cur.value) {
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
	if c == nil {
		return other.IsEmpty()
	}

	if other == nil {
		return c.IsEmpty()
	}

	if len(c.items) != len(other.items) {
		return false
	}

	for k, n := range c.items {
		on, ok := other.items[k]
		if !ok || !eq(n.value, on.value) {
			return false
		}
	}

	return true
}

// Filter returns a new Collection containing only items where fn returns true.
// Insertion order is preserved.
func (c *Collection[K, V]) Filter(fn func(V) bool) *Collection[K, V] {
	if c == nil {
		return New[K, V]()
	}

	n := New[K, V]()

	for cur := c.head; cur != nil; cur = cur.next {
		if fn(cur.value) {
			n.Set(cur.key, cur.value)
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

	if c == nil {
		return matched, rest
	}

	for cur := c.head; cur != nil; cur = cur.next {
		if fn(cur.value) {
			matched.Set(cur.key, cur.value)
		} else {
			rest.Set(cur.key, cur.value)
		}
	}

	return matched, rest
}

// First returns the first value in insertion order
func (c *Collection[K, V]) First() (*V, bool) {
	if c == nil {
		return nil, false
	}

	if c.head == nil {
		return nil, false
	}

	return &c.head.value, true
}

// FirstN returns up to n first values in insertion order.
// Returns nil if n <= 0 or Collection is empty.
func (c *Collection[K, V]) FirstN(n int) []V {
	if c == nil || n <= 0 || len(c.items) == 0 {
		return nil
	}

	if n == 1 && c.head != nil {
		return []V{c.head.value}
	}

	result := make([]V, 0, min(n, len(c.items)))
	for cur := c.head; cur != nil && len(result) < n; cur = cur.next {
		result = append(result, cur.value)
	}

	return result
}

// FirstKey returns the first key.
func (c *Collection[K, V]) FirstKey() (*K, bool) {
	if c == nil {
		return nil, false
	}

	if c.head == nil {
		return nil, false
	}

	return &c.head.key, true
}

// Last returns the last value in insertion order.
func (c *Collection[K, V]) Last() (*V, bool) {
	if c == nil {
		return nil, false
	}

	if c.tail == nil {
		return nil, false
	}

	return &c.tail.value, true
}

// LastN returns up to n last values in insertion order.
// Returns nil if n <= 0 or Collection is empty.
func (c *Collection[K, V]) LastN(n int) []V {
	if c == nil {
		return nil
	}

	if n <= 0 {
		return nil
	}

	if len(c.items) == 0 || c.tail == nil {
		return nil
	}

	total := len(c.items)
	if n > total {
		n = total
	}

	result := make([]V, n)
	i := n - 1
	for cur := c.tail; cur != nil && i >= 0; cur = cur.prev {
		result[i] = cur.value
		i--
	}

	return result
}

// LastKey returns the last key.
func (c *Collection[K, V]) LastKey() (*K, bool) {
	if c == nil {
		return nil, false
	}

	if c.tail == nil {
		return nil, false
	}

	return &c.tail.key, true
}

// At returns the value at the given insertion-order index. Negative indices
// count from the end (At(-1) == Last()).
func (c *Collection[K, V]) At(index int) (*V, bool) {
	if c == nil {
		return nil, false
	}

	n := len(c.items)

	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		return nil, false
	}

	if index <= n/2 {
		cur := c.head
		for i := 0; i < index; i++ {
			cur = cur.next
		}
		return &cur.value, true
	}

	cur := c.tail
	for i := n - 1; i > index; i-- {
		cur = cur.prev
	}

	return &cur.value, true
}

// KeyAt returns the key at the given insertion-order index. Negative indices
// count from the end.
func (c *Collection[K, V]) KeyAt(index int) (*K, bool) {
	if c == nil {
		return nil, false
	}

	n := len(c.items)

	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		return nil, false
	}

	if index <= n/2 {
		cur := c.head
		for i := 0; i < index; i++ {
			cur = cur.next
		}
		return &cur.key, true
	}

	cur := c.tail
	for i := n - 1; i > index; i-- {
		cur = cur.prev
	}

	return &cur.key, true
}

// Random returns a random value or nil if the Collection is empty.
func (c *Collection[K, V]) Random() (*V, bool) {
	if c == nil {
		return nil, false
	}

	if len(c.items) == 0 {
		return nil, false
	}

	target := rand.IntN(len(c.items))

	for _, n := range c.items {
		if target == 0 {
			return &n.value, true
		}
		target--
	}

	return nil, false
}

// RandomN returns up to n distinct random values.
// Returns nil if n <= 0 or Collection is empty.
func (c *Collection[K, V]) RandomN(n int) []V {
	if c == nil {
		return nil
	}

	if n <= 0 || len(c.items) == 0 {
		return nil
	}

	total := len(c.items)

	if n > total {
		n = total
	}

	if n == 1 {
		if v, ok := c.Random(); ok {
			return []V{*v}
		}

		return nil
	}

	all := make([]V, 0, total)

	for cur := c.head; cur != nil; cur = cur.next {
		all = append(all, cur.value)
	}

	for i := 0; i < n; i++ {
		j := i + rand.IntN(total-i)
		all[i], all[j] = all[j], all[i]
	}

	return all[:n]
}

// RandomKey returns a random key
func (c *Collection[K, V]) RandomKey() (*K, bool) {
	if c == nil {
		return nil, false
	}

	if len(c.items) == 0 {
		return nil, false
	}

	target := rand.IntN(len(c.items))

	for _, n := range c.items {
		if target == 0 {
			return &n.key, true
		}
		target--
	}

	return nil, false
}

// Each iterates over all (key, value) pairs in insertion order, calling
// fn. Returns the Collection for chaining.
//
//	users.Each(func(id Snowflake, u *User) { fmt.Println(u.Name) })
func (c *Collection[K, V]) Each(fn func(K, V)) *Collection[K, V] {
	if c == nil {
		return nil
	}

	for cur := c.head; cur != nil; cur = cur.next {
		fn(cur.key, cur.value)
	}

	return c
}

// Sort sorts the Collection in-place by the less function. Sort is stable.
// Returns the Collection for chaining.
// If you do not want to modify the Collection, use Sorted instead.
func (c *Collection[K, V]) Sort(less func(a, b V) bool) *Collection[K, V] {
	if c == nil {
		return nil
	}

	if len(c.items) < 2 {
		return c
	}

	nodes := make([]*node[K, V], 0, len(c.items))
	for cur := c.head; cur != nil; cur = cur.next {
		nodes = append(nodes, cur)
	}

	slices.SortStableFunc(nodes, func(a, b *node[K, V]) int {
		if less(a.value, b.value) {
			return -1
		}
		if less(b.value, a.value) {
			return 1
		}
		return 0
	})

	c.head = nodes[0]
	c.tail = nodes[len(nodes)-1]

	for i, n := range nodes {
		if i == 0 {
			n.prev = nil
		} else {
			n.prev = nodes[i-1]
		}
		if i == len(nodes)-1 {
			n.next = nil
		} else {
			n.next = nodes[i+1]
		}
	}

	return c
}

// Sorted returns a sorted copy without modifying the original.
func (c *Collection[K, V]) Sorted(less func(a, b V) bool) *Collection[K, V] {
	if c == nil {
		return New[K, V]()
	}

	return c.Clone().Sort(less)
}

// Reverse reverses the insertion order in place. Returns the Collection.
func (c *Collection[K, V]) Reverse() *Collection[K, V] {
	if c == nil {
		return nil
	}

	if len(c.items) < 2 {
		return c
	}

	cur := c.head
	for cur != nil {
		cur.prev, cur.next = cur.next, cur.prev
		cur = cur.prev
	}

	c.head, c.tail = c.tail, c.head

	return c
}

// MapValues transforms values in-place. For type-changing transforms, use
// the top-level Map function.
func (c *Collection[K, V]) MapValues(fn func(V) V) *Collection[K, V] {
	if c == nil {
		return New[K, V]()
	}

	for cur := c.head; cur != nil; cur = cur.next {
		cur.value = fn(cur.value)
	}

	return c
}

// Sweep removes items where fn returns true. Returns the number of items
// removed. This is the inverse of Filter — modifies in place.
func (c *Collection[K, V]) Sweep(fn func(V) bool) int {
	if c == nil {
		return 0
	}

	removed := 0
	cur := c.head

	for cur != nil {
		next := cur.next
		if fn(cur.value) {
			c.unlink(cur)
			delete(c.items, cur.key)
			removed++
		}
		cur = next
	}

	return removed
}

// SweepKeys removes items where fn returns true. Returns the keys of removed
// items. Use instead of Sweep when the caller needs to know which keys were
// deleted (e.g. to update an external index).
func (c *Collection[K, V]) SweepKeys(fn func(V) bool) []K {
	if c == nil {
		return nil
	}

	var removed []K
	cur := c.head

	for cur != nil {
		next := cur.next
		if fn(cur.value) {
			c.unlink(cur)
			delete(c.items, cur.key)
			removed = append(removed, cur.key)
		}
		cur = next
	}

	return removed
}

// Concat returns a new Collection containing all items from c and others.
// On duplicate keys, later collections overrides old.
func (c *Collection[K, V]) Concat(others ...*Collection[K, V]) *Collection[K, V] {
	if c == nil {
		return New[K, V]()
	}

	if others == nil || len(others) == 0 {
		return c.Clone()
	}

	n := c.Clone()
	for _, other := range others {
		if other == nil {
			continue
		}
		for cur := other.head; cur != nil; cur = cur.next {
			n.Set(cur.key, cur.value)
		}
	}
	return n
}

// Merge is an alias for Concat for naming familiarity with discord.js.
func (c *Collection[K, V]) Merge(others ...*Collection[K, V]) *Collection[K, V] {
	return c.Concat(others...)
}

// HasKeyDifference returns whether the Collection and other contain
// different sets of keys. Values are ignored.
func (c *Collection[K, V]) HasKeyDifference(other *Collection[K, V]) bool {
	if c == nil {
		return !other.IsEmpty()
	}

	if other == nil {
		return !c.IsEmpty()
	}

	if len(c.items) != len(other.items) {
		return true
	}

	for cur := c.head; cur != nil; cur = cur.next {
		if !other.Has(cur.key) {
			return true
		}
	}

	return false
}

// Difference returns items in c whose keys are NOT in other.
// Returns an empty Collection if there are no differences.
func (c *Collection[K, V]) Difference(other *Collection[K, V]) *Collection[K, V] {
	if c == nil {
		return New[K, V]()
	}

	if other == nil {
		return c.Clone()
	}

	n := New[K, V]()

	for cur := c.head; cur != nil; cur = cur.next {
		if !other.Has(cur.key) {
			n.Set(cur.key, cur.value)
		}
	}

	return n
}

// Intersection returns items present in both c and other (by key, with
// values from c). Returns an empty Collection if either Collection is nil.
func (c *Collection[K, V]) Intersection(other *Collection[K, V]) *Collection[K, V] {
	if c == nil || other == nil {
		return New[K, V]()
	}

	n := New[K, V]()

	for cur := c.head; cur != nil; cur = cur.next {
		if other.Has(cur.key) {
			n.Set(cur.key, cur.value)
		}
	}

	return n
}

// SymmetricDifference returns items present in c or other but not both.
func (c *Collection[K, V]) SymmetricDifference(other *Collection[K, V]) *Collection[K, V] {
	if c == nil && other == nil {
		return New[K, V]()
	}

	if c == nil {
		c = New[K, V]()
	}

	if other == nil {
		return c.Clone()
	}

	n := New[K, V]()

	for cur := c.head; cur != nil; cur = cur.next {
		if !other.Has(cur.key) {
			n.Set(cur.key, cur.value)
		}
	}

	for cur := other.head; cur != nil; cur = cur.next {
		if !c.Has(cur.key) {
			n.Set(cur.key, cur.value)
		}
	}

	return n
}

// Keys returns a snapshot slice of all keys in insertion order.
// Returns nil if Collection is empty.
func (c *Collection[K, V]) Keys() []K {
	if c == nil || len(c.items) == 0 {
		return nil
	}

	result := make([]K, 0, len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		result = append(result, cur.key)
	}

	return result
}

// Values returns a snapshot slice of all values in insertion order.
// Returns nil if Collection is empty.
func (c *Collection[K, V]) Values() []V {
	if c == nil || len(c.items) == 0 {
		return nil
	}

	result := make([]V, 0, len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		result = append(result, cur.value)
	}

	return result
}

// Entries returns a snapshot slice of all key-value pairs in insertion order.
// Returns nil if Collection is empty.
func (c *Collection[K, V]) Entries() []Entry[K, V] {
	if c == nil || len(c.items) == 0 {
		return nil
	}

	result := make([]Entry[K, V], 0, len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		result = append(result, Entry[K, V]{Key: cur.key, Value: cur.value})
	}

	return result
}

// ToMap returns a copy of the underlying map. Mutations to the returned
// map don't affect the Collection.
func (c *Collection[K, V]) ToMap() map[K]V {
	if c == nil {
		return make(map[K]V)
	}

	m := make(map[K]V, len(c.items))

	for k, n := range c.items {
		m[k] = n.value
	}

	return m
}

// All returns iter.Seq2[K, V] for range-over-func iteration.
func (c *Collection[K, V]) All() iter.Seq2[K, V] {
	if c == nil {
		return func(yield func(K, V) bool) {}
	}

	return func(yield func(K, V) bool) {
		for cur := c.head; cur != nil; cur = cur.next {
			if !yield(cur.key, cur.value) {
				return
			}
		}
	}
}

// KeysIter returns iter.Seq[K] for range-over-func on keys only.
func (c *Collection[K, V]) KeysIter() iter.Seq[K] {
	if c == nil {
		return func(yield func(K) bool) {}
	}

	return func(yield func(K) bool) {
		for cur := c.head; cur != nil; cur = cur.next {
			if !yield(cur.key) {
				return
			}
		}
	}
}

// ValuesIter returns iter.Seq[V] for range-over-func on values only.
func (c *Collection[K, V]) ValuesIter() iter.Seq[V] {
	if c == nil {
		return func(yield func(V) bool) {}
	}

	return func(yield func(V) bool) {
		for cur := c.head; cur != nil; cur = cur.next {
			if !yield(cur.value) {
				return
			}
		}
	}
}

// Map transforms a *Collection[K, V] into a slice of R using fn.
func Map[K comparable, V any, R any](c *Collection[K, V], fn func(V) R) []R {
	if c == nil {
		return make([]R, 0)
	}

	result := make([]R, 0, len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		result = append(result, fn(cur.value))
	}

	return result
}

// MapToCollection transforms a *Collection[K, V] into a Collection[K, R]
// preserving keys and insertion order.
func MapToCollection[K comparable, V any, R any](c *Collection[K, V], fn func(V) R) *Collection[K, R] {
	if c == nil {
		return New[K, R]()
	}

	n := NewWithCapacity[K, R](len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		n.Set(cur.key, fn(cur.value))
	}

	return n
}

// FlatMap applies fn to each value, then concatenates results into a
// single slice. Returns nil if Collection is empty.
func FlatMap[K comparable, V any, R any](c *Collection[K, V], fn func(V) []R) []R {
	if c == nil || len(c.items) == 0 {
		return nil
	}

	result := make([]R, 0, len(c.items))

	for cur := c.head; cur != nil; cur = cur.next {
		result = append(result, fn(cur.value)...)
	}

	return result
}

// Reduce accumulates a single value from all items, in insertion order.
//
//	total := collection.Reduce(users, 0, func(acc int, u *User) int {
//	    return acc + u.MessageCount
//	})
func Reduce[K comparable, V any, R any](c *Collection[K, V], initial R, fn func(R, V) R) R {
	if c == nil || len(c.items) == 0 {
		return initial
	}

	acc := initial

	for cur := c.head; cur != nil; cur = cur.next {
		acc = fn(acc, cur.value)
	}

	return acc
}
