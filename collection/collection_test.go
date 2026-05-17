package collection_test

import (
	"strings"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/collection"
)

// newABC returns a Collection with keys inserted in order 3, 1, 2
// mapping to "gamma", "alpha", "beta".
func newABC() *collection.Collection[int, string] {
	c := collection.New[int, string]()
	c.Set(3, "gamma")
	c.Set(1, "alpha")
	c.Set(2, "beta")
	return c
}

// --- Constructors ---

func TestNew(t *testing.T) {
	c := collection.New[int, string]()
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	if !c.IsEmpty() {
		t.Error("expected IsEmpty to be true")
	}
}

func TestNewWithCapacity(t *testing.T) {
	c := collection.NewWithCapacity[int, string](10)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Set(1, "a")
	if c.Len() != 1 {
		t.Fatalf("expected 1, got %d", c.Len())
	}
}

func TestFrom(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	c := collection.From(m)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	v, ok := c.Get(1)
	if !ok || v != "a" {
		t.Errorf("expected a/true, got %s/%v", v, ok)
	}
	// Mutation to source map does not affect Collection
	m[3] = "c"
	if c.Has(3) {
		t.Error("mutation of source map should not affect Collection")
	}
}

func TestFromSlice(t *testing.T) {
	type item struct {
		ID  int
		Val string
	}
	items := []item{{1, "a"}, {2, "b"}, {3, "c"}}
	c := collection.FromSlice(items, func(x item) int { return x.ID })
	if c.Len() != 3 {
		t.Fatalf("expected 3, got %d", c.Len())
	}
	v, ok := c.Get(2)
	if !ok || v.Val != "b" {
		t.Errorf("expected b/true, got %s/%v", v.Val, ok)
	}
}

func TestFromEntries(t *testing.T) {
	entries := []collection.Entry[int, string]{
		{Key: 3, Value: "gamma"},
		{Key: 1, Value: "alpha"},
		{Key: 2, Value: "beta"},
	}
	c := collection.FromEntries(entries)
	if c.Len() != 3 {
		t.Fatalf("expected 3, got %d", c.Len())
	}
	if keys := c.Keys(); keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("wrong insertion order: %v", keys)
	}
}

// --- Basic Methods ---

func TestGetMissing(t *testing.T) {
	c := collection.New[int, string]()
	v, ok := c.Get(99)
	if ok {
		t.Error("expected ok=false for missing key")
	}
	if v != "" {
		t.Errorf("expected zero-value, got %q", v)
	}
}

func TestMustGet(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "a")
	if got := c.MustGet(1); got != "a" {
		t.Errorf("expected a, got %s", got)
	}
}

func TestMustGetPanic(t *testing.T) {
	c := collection.New[int, string]()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing key")
		}
	}()
	c.MustGet(99)
}

func TestGetOr(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "a")
	if got := c.GetOr(1, "fallback"); got != "a" {
		t.Errorf("expected a, got %s", got)
	}
	if got := c.GetOr(99, "fallback"); got != "fallback" {
		t.Errorf("expected fallback, got %s", got)
	}
}

func TestSetChaining(t *testing.T) {
	c := collection.New[int, string]()
	ret := c.Set(1, "a").Set(2, "b")
	if ret != c {
		t.Error("Set should return the same Collection for chaining")
	}
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestSetExistingKeyPreservesOrder(t *testing.T) {
	c := newABC()
	c.Set(1, "ALPHA") // existing key, value replaced, order unchanged
	keys := c.Keys()
	if keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("wrong order after Set of existing key: %v", keys)
	}
	v, _ := c.Get(1)
	if v != "ALPHA" {
		t.Errorf("expected ALPHA, got %s", v)
	}
}

func TestHas(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "a")
	if !c.Has(1) {
		t.Error("expected Has(1)=true")
	}
	if c.Has(2) {
		t.Error("expected Has(2)=false")
	}
}

func TestDelete(t *testing.T) {
	c := newABC()
	ok := c.Delete(1)
	if !ok {
		t.Error("expected Delete to return true")
	}
	if c.Has(1) {
		t.Error("key should be gone after Delete")
	}
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	// Order of remaining keys preserved
	keys := c.Keys()
	if keys[0] != 3 || keys[1] != 2 {
		t.Errorf("wrong order after Delete: %v", keys)
	}
	// Delete non-existent key
	if c.Delete(99) {
		t.Error("Delete of missing key should return false")
	}
}

func TestLenIsEmpty(t *testing.T) {
	c := collection.New[int, string]()
	if c.Len() != 0 || !c.IsEmpty() {
		t.Error("new collection should be empty")
	}
	c.Set(1, "a")
	if c.Len() != 1 || c.IsEmpty() {
		t.Error("after Set, Len=1 IsEmpty=false expected")
	}
}

func TestClear(t *testing.T) {
	c := newABC()
	c.Clear()
	if c.Len() != 0 || !c.IsEmpty() {
		t.Error("after Clear, collection should be empty")
	}
	if len(c.Keys()) != 0 {
		t.Error("Keys() should be empty after Clear")
	}
}

func TestCloneShallow(t *testing.T) {
	type val struct{ X int }
	c := collection.New[int, *val]()
	v := &val{X: 1}
	c.Set(1, v)
	clone := c.Clone()

	// Clone has the same value
	cv, _ := clone.Get(1)
	if cv != v {
		t.Error("Clone should share the same pointer")
	}
	// Mutating the pointed-to struct is visible in both
	v.X = 42
	cv2, _ := c.Get(1)
	cv3, _ := clone.Get(1)
	if cv2.X != 42 || cv3.X != 42 {
		t.Error("mutation through shared pointer should be visible in both")
	}
	// Structural independence: Set on clone does not affect original
	clone.Set(2, &val{X: 99})
	if c.Has(2) {
		t.Error("Set on clone should not affect original")
	}
}

func TestInsertionOrderPreserved(t *testing.T) {
	c := newABC() // inserted 3, 1, 2
	keys := c.Keys()
	if len(keys) != 3 || keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("expected [3,1,2], got %v", keys)
	}
	vals := c.Values()
	if vals[0] != "gamma" || vals[1] != "alpha" || vals[2] != "beta" {
		t.Errorf("unexpected values order: %v", vals)
	}
}

// --- Search / Test ---

func TestFindEmpty(t *testing.T) {
	c := collection.New[int, string]()
	v, ok := c.Find(func(s string) bool { return true })
	if ok || v != "" {
		t.Error("Find on empty should return zero+false")
	}
}

func TestFind(t *testing.T) {
	c := newABC()
	v, ok := c.Find(func(s string) bool { return strings.HasPrefix(s, "a") })
	if !ok || v != "alpha" {
		t.Errorf("expected alpha/true, got %s/%v", v, ok)
	}
	_, ok2 := c.Find(func(s string) bool { return s == "zzz" })
	if ok2 {
		t.Error("Find with no match should return false")
	}
}

func TestFindKey(t *testing.T) {
	c := newABC()
	k, ok := c.FindKey(func(s string) bool { return s == "beta" })
	if !ok || k != 2 {
		t.Errorf("expected 2/true, got %d/%v", k, ok)
	}
}

func TestFindKeyEmpty(t *testing.T) {
	c := collection.New[int, string]()
	k, ok := c.FindKey(func(s string) bool { return true })
	if ok || k != 0 {
		t.Error("FindKey on empty should return zero+false")
	}
}

func TestFindLast(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "alpha")
	c.Set(2, "apple")
	c.Set(3, "beta")
	v, ok := c.FindLast(func(s string) bool { return strings.HasPrefix(s, "a") })
	if !ok || v != "apple" {
		t.Errorf("expected apple/true, got %s/%v", v, ok)
	}
}

func TestFindLastEmpty(t *testing.T) {
	c := collection.New[int, string]()
	v, ok := c.FindLast(func(s string) bool { return true })
	if ok || v != "" {
		t.Error("FindLast on empty should return zero+false")
	}
}

func TestFindLastKey(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "alpha")
	c.Set(2, "apple")
	c.Set(3, "beta")
	k, ok := c.FindLastKey(func(s string) bool { return strings.HasPrefix(s, "a") })
	if !ok || k != 2 {
		t.Errorf("expected 2/true, got %d/%v", k, ok)
	}
}

func TestFindLastKeyEmpty(t *testing.T) {
	c := collection.New[int, string]()
	k, ok := c.FindLastKey(func(s string) bool { return true })
	if ok || k != 0 {
		t.Error("FindLastKey on empty should return zero+false")
	}
}

func TestSome(t *testing.T) {
	c := newABC()
	if !c.Some(func(s string) bool { return s == "alpha" }) {
		t.Error("Some should be true when element matches")
	}
	if c.Some(func(s string) bool { return s == "zzz" }) {
		t.Error("Some should be false when no element matches")
	}
}

func TestSomeEmpty(t *testing.T) {
	c := collection.New[int, string]()
	if c.Some(func(s string) bool { return true }) {
		t.Error("Some on empty should be false")
	}
}

func TestEvery(t *testing.T) {
	c := newABC()
	if !c.Every(func(s string) bool { return len(s) > 3 }) {
		t.Error("Every should be true when all match")
	}
	if c.Every(func(s string) bool { return s == "alpha" }) {
		t.Error("Every should be false when not all match")
	}
}

func TestEveryEmpty(t *testing.T) {
	c := collection.New[int, string]()
	if !c.Every(func(s string) bool { return false }) {
		t.Error("Every on empty should be true (vacuous truth)")
	}
}

func TestEquals(t *testing.T) {
	a := collection.New[int, string]()
	a.Set(1, "a").Set(2, "b")
	b := collection.New[int, string]()
	b.Set(2, "b").Set(1, "a")
	if !a.Equals(b) {
		t.Error("collections with same key-value pairs should be equal")
	}
	b.Set(3, "c")
	if a.Equals(b) {
		t.Error("collections of different sizes should not be equal")
	}
}

func TestEqualsPointers(t *testing.T) {
	type val struct{ X int }
	v1 := &val{1}
	v2 := &val{1} // same content, different pointer
	a := collection.New[int, *val]()
	a.Set(1, v1)
	b := collection.New[int, *val]()
	b.Set(1, v2)
	if a.Equals(b) {
		t.Error("pointer equality: different pointers should not be equal")
	}
	b2 := collection.New[int, *val]()
	b2.Set(1, v1)
	if !a.Equals(b2) {
		t.Error("same pointer should be equal")
	}
}

// --- Filter / Partition ---

func TestFilter(t *testing.T) {
	c := newABC()
	result := c.Filter(func(v string) bool { return strings.HasPrefix(v, "a") })
	if result.Len() != 1 {
		t.Fatalf("expected 1, got %d", result.Len())
	}
	v, _ := result.Get(1)
	if v != "alpha" {
		t.Errorf("expected alpha, got %s", v)
	}
}

func TestFilterEmpty(t *testing.T) {
	c := collection.New[int, string]()
	result := c.Filter(func(v string) bool { return true })
	if result.Len() != 0 {
		t.Error("Filter on empty should return empty")
	}
}

func TestFilterPreservesOrder(t *testing.T) {
	c := newABC()
	result := c.Filter(func(v string) bool { return v != "gamma" })
	keys := result.Keys()
	if len(keys) != 2 || keys[0] != 1 || keys[1] != 2 {
		t.Errorf("Filter should preserve insertion order, got %v", keys)
	}
}

func TestPartition(t *testing.T) {
	c := newABC()
	matched, rest := c.Partition(func(v string) bool { return strings.HasPrefix(v, "a") })
	if matched.Len() != 1 {
		t.Fatalf("matched: expected 1, got %d", matched.Len())
	}
	if rest.Len() != 2 {
		t.Fatalf("rest: expected 2, got %d", rest.Len())
	}
	v, _ := matched.Get(1)
	if v != "alpha" {
		t.Errorf("expected alpha in matched, got %s", v)
	}
}

func TestPartitionEmpty(t *testing.T) {
	c := collection.New[int, string]()
	matched, rest := c.Partition(func(v string) bool { return true })
	if matched.Len() != 0 || rest.Len() != 0 {
		t.Error("Partition on empty should yield two empty collections")
	}
}

// --- Index / Position ---

func TestFirstLastSingleElement(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(42, "only")

	v, ok := c.First()
	if !ok || v != "only" {
		t.Errorf("First: expected only/true, got %s/%v", v, ok)
	}
	v, ok = c.Last()
	if !ok || v != "only" {
		t.Errorf("Last: expected only/true, got %s/%v", v, ok)
	}
	k, ok := c.FirstKey()
	if !ok || k != 42 {
		t.Errorf("FirstKey: expected 42/true, got %d/%v", k, ok)
	}
	k, ok = c.LastKey()
	if !ok || k != 42 {
		t.Errorf("LastKey: expected 42/true, got %d/%v", k, ok)
	}
}

func TestFirstEmpty(t *testing.T) {
	c := collection.New[int, string]()
	v, ok := c.First()
	if ok || v != "" {
		t.Error("First on empty should return zero+false")
	}
	k, ok2 := c.FirstKey()
	if ok2 || k != 0 {
		t.Error("FirstKey on empty should return zero+false")
	}
}

func TestLastEmpty(t *testing.T) {
	c := collection.New[int, string]()
	v, ok := c.Last()
	if ok || v != "" {
		t.Error("Last on empty should return zero+false")
	}
	k, ok2 := c.LastKey()
	if ok2 || k != 0 {
		t.Error("LastKey on empty should return zero+false")
	}
}

func TestFirstN(t *testing.T) {
	c := newABC()
	got := c.FirstN(2)
	if len(got) != 2 || got[0] != "gamma" || got[1] != "alpha" {
		t.Errorf("FirstN(2) = %v, want [gamma alpha]", got)
	}
	all := c.FirstN(100)
	if len(all) != 3 {
		t.Errorf("FirstN beyond Len should clamp, got %d", len(all))
	}
	none := c.FirstN(0)
	if len(none) != 0 {
		t.Error("FirstN(0) should return empty slice")
	}
}

func TestLastN(t *testing.T) {
	c := newABC() // [gamma, alpha, beta]
	got := c.LastN(2)
	if len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Errorf("LastN(2) = %v, want [alpha beta]", got)
	}
	all := c.LastN(100)
	if len(all) != 3 {
		t.Errorf("LastN beyond Len should clamp, got %d", len(all))
	}
}

func TestAt(t *testing.T) {
	c := newABC() // [gamma, alpha, beta]
	tests := []struct {
		index int
		want  string
		ok    bool
	}{
		{0, "gamma", true},
		{1, "alpha", true},
		{2, "beta", true},
		{-1, "beta", true},
		{-2, "alpha", true},
		{-3, "gamma", true},
		{3, "", false},
		{-4, "", false},
		{100, "", false},
	}
	for _, tt := range tests {
		v, ok := c.At(tt.index)
		if ok != tt.ok || (tt.ok && v != tt.want) {
			t.Errorf("At(%d) = (%s, %v), want (%s, %v)", tt.index, v, ok, tt.want, tt.ok)
		}
	}
}

func TestKeyAt(t *testing.T) {
	c := newABC() // keys: [3, 1, 2]
	tests := []struct {
		index int
		want  int
		ok    bool
	}{
		{0, 3, true},
		{1, 1, true},
		{2, 2, true},
		{-1, 2, true},
		{-2, 1, true},
		{-3, 3, true},
		{3, 0, false},
		{-4, 0, false},
	}
	for _, tt := range tests {
		k, ok := c.KeyAt(tt.index)
		if ok != tt.ok || (tt.ok && k != tt.want) {
			t.Errorf("KeyAt(%d) = (%d, %v), want (%d, %v)", tt.index, k, ok, tt.want, tt.ok)
		}
	}
}

func TestRandomEmpty(t *testing.T) {
	c := collection.New[int, string]()
	v, ok := c.Random()
	if ok || v != "" {
		t.Error("Random on empty should return zero+false")
	}
	k, ok2 := c.RandomKey()
	if ok2 || k != 0 {
		t.Error("RandomKey on empty should return zero+false")
	}
	if got := c.RandomN(5); len(got) != 0 {
		t.Error("RandomN on empty should return empty slice")
	}
}

func TestRandom(t *testing.T) {
	c := newABC()
	_, ok := c.Random()
	if !ok {
		t.Error("Random should succeed on non-empty")
	}
	_, ok2 := c.RandomKey()
	if !ok2 {
		t.Error("RandomKey should succeed on non-empty")
	}
	got := c.RandomN(2)
	if len(got) != 2 {
		t.Errorf("RandomN(2) should return 2 elements, got %d", len(got))
	}
	// Clamping
	all := c.RandomN(100)
	if len(all) != 3 {
		t.Errorf("RandomN beyond Len should clamp to Len, got %d", len(all))
	}
}

// --- Transform ---

func TestEach(t *testing.T) {
	c := newABC()
	var keys []int
	var vals []string
	ret := c.Each(func(k int, v string) {
		keys = append(keys, k)
		vals = append(vals, v)
	})
	if ret != c {
		t.Error("Each should return the same Collection for chaining")
	}
	if len(keys) != 3 || keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("Each: wrong key order %v", keys)
	}
	if vals[0] != "gamma" || vals[1] != "alpha" || vals[2] != "beta" {
		t.Errorf("Each: wrong value order %v", vals)
	}
}

func TestSort(t *testing.T) {
	c := newABC()
	ret := c.Sort(func(a, b string) bool { return a < b })
	if ret != c {
		t.Error("Sort should return the same Collection for chaining")
	}
	keys := c.Keys()
	// alpha < beta < gamma, so keys in order 1, 2, 3
	if keys[0] != 1 || keys[1] != 2 || keys[2] != 3 {
		t.Errorf("Sort: wrong key order %v", keys)
	}
}

func TestSortStability(t *testing.T) {
	// Two values that are equal under the less function — original order preserved.
	c := collection.New[int, int]()
	c.Set(10, 5)
	c.Set(20, 5)
	c.Set(30, 5)
	c.Sort(func(a, b int) bool { return a < b })
	keys := c.Keys()
	if keys[0] != 10 || keys[1] != 20 || keys[2] != 30 {
		t.Errorf("Sort should be stable, got %v", keys)
	}
}

func TestSorted(t *testing.T) {
	c := newABC()
	sorted := c.Sorted(func(a, b string) bool { return a < b })
	// Original unchanged
	origKeys := c.Keys()
	if origKeys[0] != 3 || origKeys[1] != 1 || origKeys[2] != 2 {
		t.Errorf("Sorted modified original, keys=%v", origKeys)
	}
	// Sorted has correct order
	sortedKeys := sorted.Keys()
	if sortedKeys[0] != 1 || sortedKeys[1] != 2 || sortedKeys[2] != 3 {
		t.Errorf("Sorted: wrong key order %v", sortedKeys)
	}
}

func TestReverse(t *testing.T) {
	c := newABC() // [3, 1, 2]
	ret := c.Reverse()
	if ret != c {
		t.Error("Reverse should return the same Collection for chaining")
	}
	keys := c.Keys()
	if keys[0] != 2 || keys[1] != 1 || keys[2] != 3 {
		t.Errorf("Reverse: expected [2,1,3], got %v", keys)
	}
}

func TestMapValues(t *testing.T) {
	c := newABC()
	ret := c.MapValues(strings.ToUpper)
	if ret != c {
		t.Error("MapValues should return the same Collection for chaining")
	}
	v, _ := c.Get(1)
	if v != "ALPHA" {
		t.Errorf("expected ALPHA, got %s", v)
	}
}

func TestSweep(t *testing.T) {
	c := newABC()
	removed := c.Sweep(func(v string) bool { return strings.HasPrefix(v, "a") })
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if c.Has(1) {
		t.Error("swept key should be gone")
	}
	if c.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", c.Len())
	}
	keys := c.Keys()
	if keys[0] != 3 || keys[1] != 2 {
		t.Errorf("Sweep should preserve order of remaining keys: %v", keys)
	}
}

func TestSweepAll(t *testing.T) {
	c := newABC()
	removed := c.Sweep(func(v string) bool { return true })
	if removed != 3 || c.Len() != 0 {
		t.Errorf("Sweep all: removed=%d Len=%d", removed, c.Len())
	}
}

func TestSweepNone(t *testing.T) {
	c := newABC()
	removed := c.Sweep(func(v string) bool { return false })
	if removed != 0 || c.Len() != 3 {
		t.Error("Sweep with fn=false should remove nothing")
	}
}

// --- Set Operations ---

func TestConcat(t *testing.T) {
	a := collection.New[int, string]()
	a.Set(1, "a").Set(2, "b")
	b := collection.New[int, string]()
	b.Set(2, "B").Set(3, "c") // key 2 overrides
	result := a.Concat(b)
	if result.Len() != 3 {
		t.Fatalf("expected 3, got %d", result.Len())
	}
	v, _ := result.Get(2)
	if v != "B" {
		t.Errorf("later collection should override: expected B, got %s", v)
	}
	// Original unchanged
	orig, _ := a.Get(2)
	if orig != "b" {
		t.Error("Concat should not modify original")
	}
}

func TestConcatEmpty(t *testing.T) {
	a := newABC()
	empty := collection.New[int, string]()
	result := a.Concat(empty)
	if result.Len() != 3 {
		t.Errorf("concat with empty: expected 3, got %d", result.Len())
	}
	result2 := empty.Concat(a)
	if result2.Len() != 3 {
		t.Errorf("empty.concat(a): expected 3, got %d", result2.Len())
	}
}

func TestMergeAlias(t *testing.T) {
	a := collection.New[int, string]()
	a.Set(1, "a")
	b := collection.New[int, string]()
	b.Set(2, "b")
	r1 := a.Concat(b)
	r2 := a.Merge(b)
	if !r1.Equals(r2) {
		t.Error("Merge should be equivalent to Concat")
	}
}

func TestDifference(t *testing.T) {
	a := newABC()
	b := collection.New[int, string]()
	b.Set(1, "x").Set(3, "y")
	diff := a.Difference(b)
	if diff.Len() != 1 {
		t.Fatalf("expected 1, got %d", diff.Len())
	}
	v, _ := diff.Get(2)
	if v != "beta" {
		t.Errorf("expected beta, got %s", v)
	}
}

func TestDifferenceEmpty(t *testing.T) {
	a := newABC()
	empty := collection.New[int, string]()
	diff := a.Difference(empty)
	if diff.Len() != 3 {
		t.Errorf("difference with empty: expected 3, got %d", diff.Len())
	}
	diff2 := empty.Difference(a)
	if diff2.Len() != 0 {
		t.Errorf("empty.difference(a): expected 0, got %d", diff2.Len())
	}
}

func TestIntersection(t *testing.T) {
	a := newABC()
	b := collection.New[int, string]()
	b.Set(1, "x").Set(4, "y")
	inter := a.Intersection(b)
	if inter.Len() != 1 {
		t.Fatalf("expected 1, got %d", inter.Len())
	}
	v, _ := inter.Get(1)
	if v != "alpha" {
		t.Errorf("Intersection should use values from c: expected alpha, got %s", v)
	}
}

func TestIntersectionEmpty(t *testing.T) {
	a := newABC()
	empty := collection.New[int, string]()
	inter := a.Intersection(empty)
	if inter.Len() != 0 {
		t.Errorf("intersection with empty: expected 0, got %d", inter.Len())
	}
}

func TestSymmetricDifference(t *testing.T) {
	a := collection.New[int, string]()
	a.Set(1, "a").Set(2, "b").Set(3, "c")
	b := collection.New[int, string]()
	b.Set(2, "B").Set(3, "C").Set(4, "d")
	sym := a.SymmetricDifference(b)
	// Keys 1 (only in a) and 4 (only in b)
	if sym.Len() != 2 {
		t.Fatalf("expected 2, got %d", sym.Len())
	}
	if !sym.Has(1) || !sym.Has(4) {
		t.Errorf("expected keys 1 and 4 in SymmetricDifference, got %v", sym.Keys())
	}
}

func TestSymmetricDifferenceEmpty(t *testing.T) {
	a := newABC()
	empty := collection.New[int, string]()
	sym := a.SymmetricDifference(empty)
	if sym.Len() != 3 {
		t.Errorf("sym diff with empty: expected 3, got %d", sym.Len())
	}
}

// --- Iteration / Conversion ---

func TestKeys(t *testing.T) {
	c := newABC()
	keys := c.Keys()
	if len(keys) != 3 || keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("Keys: expected [3,1,2], got %v", keys)
	}
	// Snapshot: mutation doesn't affect collection
	keys[0] = 999
	if c.KeyAt(0); c.Keys()[0] != 3 {
		t.Error("Keys() should return a snapshot")
	}
}

func TestValues(t *testing.T) {
	c := newABC()
	vals := c.Values()
	if vals[0] != "gamma" || vals[1] != "alpha" || vals[2] != "beta" {
		t.Errorf("Values: wrong order %v", vals)
	}
}

func TestEntries(t *testing.T) {
	c := newABC()
	entries := c.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3, got %d", len(entries))
	}
	if entries[0].Key != 3 || entries[0].Value != "gamma" {
		t.Errorf("entry[0] wrong: %+v", entries[0])
	}
	if entries[1].Key != 1 || entries[1].Value != "alpha" {
		t.Errorf("entry[1] wrong: %+v", entries[1])
	}
}

func TestToMap(t *testing.T) {
	c := newABC()
	m := c.ToMap()
	if len(m) != 3 {
		t.Fatalf("expected 3, got %d", len(m))
	}
	if m[1] != "alpha" || m[2] != "beta" || m[3] != "gamma" {
		t.Errorf("ToMap: wrong values %v", m)
	}
	// Mutation doesn't affect collection
	m[99] = "zzz"
	if c.Has(99) {
		t.Error("mutation of ToMap result should not affect Collection")
	}
}

func TestAll(t *testing.T) {
	c := newABC()
	var keys []int
	var vals []string
	for k, v := range c.All() {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	if keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("All: wrong key order %v", keys)
	}
	if vals[0] != "gamma" || vals[1] != "alpha" || vals[2] != "beta" {
		t.Errorf("All: wrong value order %v", vals)
	}
}

func TestAllBreak(t *testing.T) {
	c := newABC()
	count := 0
	for range c.All() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("All: break should stop iteration, got %d iterations", count)
	}
}

func TestKeysIter(t *testing.T) {
	c := newABC()
	var keys []int
	for k := range c.KeysIter() {
		keys = append(keys, k)
	}
	if keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("KeysIter: wrong order %v", keys)
	}
}

func TestValuesIter(t *testing.T) {
	c := newABC()
	var vals []string
	for v := range c.ValuesIter() {
		vals = append(vals, v)
	}
	if vals[0] != "gamma" || vals[1] != "alpha" || vals[2] != "beta" {
		t.Errorf("ValuesIter: wrong order %v", vals)
	}
}

// --- Top-Level Functions ---

func TestMap(t *testing.T) {
	c := newABC()
	lengths := collection.Map(c, func(v string) int { return len(v) })
	if len(lengths) != 3 {
		t.Fatalf("expected 3, got %d", len(lengths))
	}
	// gamma=5, alpha=5, beta=4
	if lengths[0] != 5 || lengths[1] != 5 || lengths[2] != 4 {
		t.Errorf("Map: wrong values %v", lengths)
	}
}

func TestMapToCollection(t *testing.T) {
	c := newABC()
	upper := collection.MapToCollection(c, strings.ToUpper)
	if upper.Len() != 3 {
		t.Fatalf("expected 3, got %d", upper.Len())
	}
	v, _ := upper.Get(1)
	if v != "ALPHA" {
		t.Errorf("expected ALPHA, got %s", v)
	}
	// Insertion order preserved
	keys := upper.Keys()
	if keys[0] != 3 || keys[1] != 1 || keys[2] != 2 {
		t.Errorf("MapToCollection: wrong order %v", keys)
	}
	// Original unchanged
	orig, _ := c.Get(1)
	if orig != "alpha" {
		t.Error("MapToCollection should not modify original")
	}
}

func TestFlatMap(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "ab")
	c.Set(2, "cd")
	result := collection.FlatMap(c, func(v string) []string {
		chars := make([]string, len(v))
		for i, ch := range v {
			chars[i] = string(ch)
		}
		return chars
	})
	if len(result) != 4 {
		t.Fatalf("expected 4, got %d", len(result))
	}
	if result[0] != "a" || result[1] != "b" || result[2] != "c" || result[3] != "d" {
		t.Errorf("FlatMap: wrong values %v", result)
	}
}

func TestGroupBy(t *testing.T) {
	c := collection.New[int, string]()
	c.Set(1, "apple")
	c.Set(2, "apricot")
	c.Set(3, "banana")
	c.Set(4, "blueberry")
	groups := collection.GroupBy(c, func(v string) string { return string(v[0]) })
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups["a"].Len() != 2 {
		t.Errorf("expected 2 'a' values, got %d", groups["a"].Len())
	}
	if groups["b"].Len() != 2 {
		t.Errorf("expected 2 'b' values, got %d", groups["b"].Len())
	}
}

func TestReduce(t *testing.T) {
	c := newABC()
	total := collection.Reduce(c, 0, func(acc int, v string) int { return acc + len(v) })
	// gamma=5, alpha=5, beta=4
	if total != 14 {
		t.Errorf("Reduce: expected 14, got %d", total)
	}
}

func TestReduceEmpty(t *testing.T) {
	c := collection.New[int, string]()
	total := collection.Reduce(c, 42, func(acc int, v string) int { return acc + 1 })
	if total != 42 {
		t.Errorf("Reduce on empty should return initial, got %d", total)
	}
}

// --- Pointer-value reflection ---

func TestPointerValueReflected(t *testing.T) {
	type user struct{ Name string }
	c := collection.New[int, *user]()
	u := &user{Name: "alice"}
	c.Set(1, u)
	u.Name = "bob"
	got, _ := c.Get(1)
	if got.Name != "bob" {
		t.Errorf("pointer mutation should be reflected in Get: got %s", got.Name)
	}
}
