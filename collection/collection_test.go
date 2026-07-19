package collection

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func newABC() *Collection[int, string] {
	c := New[int, string]()
	c.Set(3, "gamma")
	c.Set(1, "alpha")
	c.Set(2, "beta")
	return c
}

func (s *collectionSuite) TestNew() {
	c := New[int, string]()
	s.Equal(0, c.Len())
	s.True(c.IsEmpty())
}

func (s *collectionSuite) TestNewWithCapacity() {
	c := NewWithCapacity[int, string](10)
	s.Equal(0, c.Len())
	s.True(c.IsEmpty())

	c.Set(1, "a")

	s.Equal(1, c.Len())
	s.False(c.IsEmpty())
}

func (s *collectionSuite) TestFrom() {
	m := map[int]string{1: "a", 2: "b"}
	c := From(m)
	s.Equal(c.Len(), len(m))

	v, ok := c.Get(1)
	s.True(ok)
	s.Equal(v, "a")

	m[3] = "c"
	s.False(c.Has(3))
}

func (s *collectionSuite) TestFromSlice() {
	type item struct {
		ID  int
		Val string
	}

	items := []item{{1, "a"}, {2, "b"}, {3, "c"}}
	c := FromSlice(items, func(x item) int { return x.ID })
	s.Require().Equal(3, c.Len())
	v, ok := c.Get(2)
	s.True(ok)
	s.Equal("b", v.Val)
}

func (s *collectionSuite) TestFromEntries() {
	entries := []Entry[int, string]{
		{Key: 3, Value: "gamma"},
		{Key: 1, Value: "alpha"},
		{Key: 2, Value: "beta"},
	}
	c := FromEntries(entries)
	s.Require().Equal(3, c.Len())
	keys := c.Keys()
	s.Equal([]int{3, 1, 2}, keys)
}

// --- Basic Methods ---

func (s *collectionSuite) TestGetMissing() {
	c := New[int, string]()
	v, ok := c.Get(99)
	s.False(ok)
	s.Equal("", v)
}

func (s *collectionSuite) TestGetOr() {
	c := New[int, string]()
	c.Set(1, "a")
	s.Equal("a", c.GetOr(1, "fallback"))
	s.Equal("fallback", c.GetOr(99, "fallback"))
}

func (s *collectionSuite) TestSetChaining() {
	c := New[int, string]()
	ret := c.Set(1, "a").Set(2, "b")
	s.Same(c, ret)
	s.Require().Equal(2, c.Len())
}

func (s *collectionSuite) TestSetExistingKeyPreservesOrder() {
	c := newABC()
	c.Set(1, "ALPHA")
	keys := c.Keys()
	s.Equal([]int{3, 1, 2}, keys)
	v, _ := c.Get(1)
	s.Equal("ALPHA", v)
}

func (s *collectionSuite) TestHas() {
	c := New[int, string]()
	c.Set(1, "a")
	s.True(c.Has(1))
	s.False(c.Has(2))
}

func (s *collectionSuite) TestDelete() {
	c := newABC()
	ok := c.Delete(1)
	s.True(ok)
	s.False(c.Has(1))
	s.Require().Equal(2, c.Len())
	keys := c.Keys()
	s.Equal([]int{3, 2}, keys)
	s.False(c.Delete(99))
}

func (s *collectionSuite) TestLenIsEmpty() {
	c := New[int, string]()
	s.Equal(0, c.Len())
	s.True(c.IsEmpty())
	c.Set(1, "a")
	s.Equal(1, c.Len())
	s.False(c.IsEmpty())
}

func (s *collectionSuite) TestClear() {
	c := newABC()
	c.Clear()
	s.Equal(0, c.Len())
	s.True(c.IsEmpty())
	s.Len(c.Keys(), 0)
}

func (s *collectionSuite) TestCloneShallow() {
	type val struct{ X int }
	c := New[int, *val]()
	v := &val{X: 1}
	c.Set(1, v)
	clone := c.Clone()

	cv, _ := clone.Get(1)
	s.Same(v, cv)

	v.X = 42
	cv2, _ := c.Get(1)
	cv3, _ := clone.Get(1)
	s.Equal(42, cv2.X)
	s.Equal(42, cv3.X)

	clone.Set(2, &val{X: 99})
	s.False(c.Has(2))
}

func (s *collectionSuite) TestInsertionOrderPreserved() {
	c := newABC() // inserted 3, 1, 2
	keys := c.Keys()
	s.Equal([]int{3, 1, 2}, keys)
	vals := c.Values()
	s.Equal([]string{"gamma", "alpha", "beta"}, vals)
}

// --- Search / Test ---

func (s *collectionSuite) TestFindEmpty() {
	c := New[int, string]()
	v, ok := c.Find(func(s string) bool { return true })
	s.False(ok)
	s.Equal("", v)
}

func (s *collectionSuite) TestFind() {
	c := newABC()
	v, ok := c.Find(func(s string) bool { return strings.HasPrefix(s, "a") })
	s.True(ok)
	s.Equal("alpha", v)
	_, ok2 := c.Find(func(s string) bool { return s == "zzz" })
	s.False(ok2)
}

func (s *collectionSuite) TestFindKey() {
	c := newABC()
	k, ok := c.FindKey(func(s string) bool { return s == "beta" })
	s.True(ok)
	s.Equal(2, k)
}

func (s *collectionSuite) TestFindKeyEmpty() {
	c := New[int, string]()
	k, ok := c.FindKey(func(s string) bool { return true })
	s.False(ok)
	s.Equal(0, k)
}

func (s *collectionSuite) TestFindLast() {
	c := New[int, string]()
	c.Set(1, "alpha")
	c.Set(2, "apple")
	c.Set(3, "beta")
	v, ok := c.FindLast(func(s string) bool { return strings.HasPrefix(s, "a") })
	s.True(ok)
	s.Equal("apple", v)
}

func (s *collectionSuite) TestFindLastEmpty() {
	c := New[int, string]()
	v, ok := c.FindLast(func(s string) bool { return true })
	s.False(ok)
	s.Equal("", v)
}

func (s *collectionSuite) TestFindLastKey() {
	c := New[int, string]()
	c.Set(1, "alpha")
	c.Set(2, "apple")
	c.Set(3, "beta")
	k, ok := c.FindLastKey(func(s string) bool { return strings.HasPrefix(s, "a") })
	s.True(ok)
	s.Equal(2, k)
}

func (s *collectionSuite) TestFindLastKeyEmpty() {
	c := New[int, string]()
	k, ok := c.FindLastKey(func(s string) bool { return true })
	s.False(ok)
	s.Equal(0, k)
}

func (s *collectionSuite) TestSome() {
	c := newABC()
	s.True(c.Some(func(s string) bool { return s == "alpha" }))
	s.False(c.Some(func(s string) bool { return s == "zzz" }))
}

func (s *collectionSuite) TestSomeEmpty() {
	c := New[int, string]()
	s.False(c.Some(func(s string) bool { return true }))
}

func (s *collectionSuite) TestEvery() {
	c := newABC()
	s.True(c.Every(func(s string) bool { return len(s) > 3 }))
	s.False(c.Every(func(s string) bool { return s == "alpha" }))
}

func (s *collectionSuite) TestEveryEmpty() {
	c := New[int, string]()
	s.True(c.Every(func(s string) bool { return false }))
}

func (s *collectionSuite) TestEquals() {
	a := New[int, string]()
	a.Set(1, "a").Set(2, "b")
	b := New[int, string]()
	b.Set(2, "b").Set(1, "a")
	strEq := func(x, y string) bool { return x == y }
	s.True(a.Equals(b, strEq))
	b.Set(3, "c")
	s.False(a.Equals(b, strEq))
}

func (s *collectionSuite) TestEqualsPointers() {
	type val struct{ X int }
	v1 := &val{1}
	v2 := &val{1} // same content, different pointer
	a := New[int, *val]()
	a.Set(1, v1)
	b := New[int, *val]()
	b.Set(1, v2)
	ptrEq := func(x, y *val) bool { return x == y }
	s.False(a.Equals(b, ptrEq))
	b2 := New[int, *val]()
	b2.Set(1, v1)
	s.True(a.Equals(b2, ptrEq))
}

// --- Filter / Partition ---

func (s *collectionSuite) TestFilter() {
	c := newABC()
	result := c.Filter(func(v string) bool { return strings.HasPrefix(v, "a") })
	s.Require().Equal(1, result.Len())
	v, _ := result.Get(1)
	s.Equal("alpha", v)
}

func (s *collectionSuite) TestFilterEmpty() {
	c := New[int, string]()
	result := c.Filter(func(v string) bool { return true })
	s.Equal(0, result.Len())
}

func (s *collectionSuite) TestFilterPreservesOrder() {
	c := newABC()
	result := c.Filter(func(v string) bool { return v != "gamma" })
	keys := result.Keys()
	s.Equal([]int{1, 2}, keys)
}

func (s *collectionSuite) TestPartition() {
	c := newABC()
	matched, rest := c.Partition(func(v string) bool { return strings.HasPrefix(v, "a") })
	s.Require().Equal(1, matched.Len())
	s.Require().Equal(2, rest.Len())
	v, _ := matched.Get(1)
	s.Equal("alpha", v)
}

func (s *collectionSuite) TestPartitionEmpty() {
	c := New[int, string]()
	matched, rest := c.Partition(func(v string) bool { return true })
	s.Equal(0, matched.Len())
	s.Equal(0, rest.Len())
}

// --- Index / Position ---

func (s *collectionSuite) TestFirstLastSingleElement() {
	c := New[int, string]()
	c.Set(42, "only")

	v, ok := c.First()
	s.True(ok)
	s.Equal("only", v)
	v, ok = c.Last()
	s.True(ok)
	s.Equal("only", v)
	k, ok := c.FirstKey()
	s.True(ok)
	s.Equal(42, k)
	k, ok = c.LastKey()
	s.True(ok)
	s.Equal(42, k)
}

func (s *collectionSuite) TestFirstEmpty() {
	c := New[int, string]()
	v, ok := c.First()
	s.False(ok)
	s.Equal("", v)
	k, ok2 := c.FirstKey()
	s.False(ok2)
	s.Equal(0, k)
}

func (s *collectionSuite) TestLastEmpty() {
	c := New[int, string]()
	v, ok := c.Last()
	s.False(ok)
	s.Equal("", v)
	k, ok2 := c.LastKey()
	s.False(ok2)
	s.Equal(0, k)
}

func (s *collectionSuite) TestFirstN() {
	c := newABC()
	got := c.FirstN(2)
	s.Equal([]string{"gamma", "alpha"}, got)
	all := c.FirstN(100)
	s.Len(all, 3)
	none := c.FirstN(0)
	s.Len(none, 0)
}

func (s *collectionSuite) TestLastN() {
	c := newABC() // [gamma, alpha, beta]
	got := c.LastN(2)
	s.Equal([]string{"alpha", "beta"}, got)
	all := c.LastN(100)
	s.Len(all, 3)
}

func (s *collectionSuite) TestAt() {
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
		s.Equal(tt.ok, ok)
		if tt.ok {
			s.Equal(tt.want, v)
		}
	}
}

func (s *collectionSuite) TestKeyAt() {
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
		s.Equal(tt.ok, ok)
		if tt.ok {
			s.Equal(tt.want, k)
		}
	}
}

func (s *collectionSuite) TestRandomEmpty() {
	c := New[int, string]()
	v, ok := c.Random()
	s.False(ok)
	s.Equal("", v)
	k, ok2 := c.RandomKey()
	s.False(ok2)
	s.Equal(0, k)
	s.Len(c.RandomN(5), 0)
}

func (s *collectionSuite) TestRandom() {
	c := newABC()
	_, ok := c.Random()
	s.True(ok)
	_, ok2 := c.RandomKey()
	s.True(ok2)
	got := c.RandomN(2)
	s.Len(got, 2)
	all := c.RandomN(100)
	s.Len(all, 3)
}

// --- Transform ---

func (s *collectionSuite) TestEach() {
	c := newABC()
	var keys []int
	var vals []string
	ret := c.Each(func(k int, v string) {
		keys = append(keys, k)
		vals = append(vals, v)
	})
	s.Same(c, ret)
	s.Equal([]int{3, 1, 2}, keys)
	s.Equal([]string{"gamma", "alpha", "beta"}, vals)
}

func (s *collectionSuite) TestSort() {
	c := newABC()
	ret := c.Sort(func(a, b string) bool { return a < b })
	s.Same(c, ret)
	keys := c.Keys()
	s.Equal([]int{1, 2, 3}, keys)
}

func (s *collectionSuite) TestSortStability() {
	c := New[int, int]()
	c.Set(10, 5)
	c.Set(20, 5)
	c.Set(30, 5)
	c.Sort(func(a, b int) bool { return a < b })
	keys := c.Keys()
	s.Equal([]int{10, 20, 30}, keys)
}

func (s *collectionSuite) TestSorted() {
	c := newABC()
	sorted := c.Sorted(func(a, b string) bool { return a < b })
	origKeys := c.Keys()
	s.Equal([]int{3, 1, 2}, origKeys)
	sortedKeys := sorted.Keys()
	s.Equal([]int{1, 2, 3}, sortedKeys)
}

func (s *collectionSuite) TestReverse() {
	c := newABC() // [3, 1, 2]
	ret := c.Reverse()
	s.Same(c, ret)
	keys := c.Keys()
	s.Equal([]int{2, 1, 3}, keys)
}

func (s *collectionSuite) TestMapValues() {
	c := newABC()
	ret := c.MapValues(strings.ToUpper)
	s.Same(c, ret)
	v, _ := c.Get(1)
	s.Equal("ALPHA", v)
}

func (s *collectionSuite) TestSweep() {
	c := newABC()
	removed := c.Sweep(func(v string) bool { return strings.HasPrefix(v, "a") })
	s.Require().Equal(1, removed)
	s.False(c.Has(1))
	s.Require().Equal(2, c.Len())
	keys := c.Keys()
	s.Equal([]int{3, 2}, keys)
}

func (s *collectionSuite) TestSweepAll() {
	c := newABC()
	removed := c.Sweep(func(v string) bool { return true })
	s.Equal(3, removed)
	s.Equal(0, c.Len())
}

func (s *collectionSuite) TestSweepNone() {
	c := newABC()
	removed := c.Sweep(func(v string) bool { return false })
	s.Equal(0, removed)
	s.Equal(3, c.Len())
}

// --- Set Operations ---

func (s *collectionSuite) TestConcat() {
	a := New[int, string]()
	a.Set(1, "a").Set(2, "b")
	b := New[int, string]()
	b.Set(2, "B").Set(3, "c") // key 2 overrides
	result := a.Concat(b)
	s.Require().Equal(3, result.Len())
	v, _ := result.Get(2)
	s.Equal("B", v)
	orig, _ := a.Get(2)
	s.Equal("b", orig)
}

func (s *collectionSuite) TestConcatEmpty() {
	a := newABC()
	empty := New[int, string]()
	result := a.Concat(empty)
	s.Equal(3, result.Len())
	result2 := empty.Concat(a)
	s.Equal(3, result2.Len())
}

func (s *collectionSuite) TestMergeAlias() {
	a := New[int, string]()
	a.Set(1, "a")
	b := New[int, string]()
	b.Set(2, "b")
	r1 := a.Concat(b)
	r2 := a.Merge(b)
	s.True(r1.Equals(r2, func(x, y string) bool { return x == y }))
}

func (s *collectionSuite) TestDifference() {
	a := newABC()
	b := New[int, string]()
	b.Set(1, "x").Set(3, "y")
	diff := a.Difference(b)
	s.Require().Equal(1, diff.Len())
	v, _ := diff.Get(2)
	s.Equal("beta", v)
}

func (s *collectionSuite) TestDifferenceEmpty() {
	a := newABC()
	empty := New[int, string]()
	diff := a.Difference(empty)
	s.Equal(3, diff.Len())
	diff2 := empty.Difference(a)
	s.Equal(0, diff2.Len())
}

func (s *collectionSuite) TestIntersection() {
	a := newABC()
	b := New[int, string]()
	b.Set(1, "x").Set(4, "y")
	inter := a.Intersection(b)
	s.Require().Equal(1, inter.Len())
	v, _ := inter.Get(1)
	s.Equal("alpha", v)
}

func (s *collectionSuite) TestIntersectionEmpty() {
	a := newABC()
	empty := New[int, string]()
	inter := a.Intersection(empty)
	s.Equal(0, inter.Len())
}

func (s *collectionSuite) TestSymmetricDifference() {
	a := New[int, string]()
	a.Set(1, "a").Set(2, "b").Set(3, "c")
	b := New[int, string]()
	b.Set(2, "B").Set(3, "C").Set(4, "d")
	sym := a.SymmetricDifference(b)
	s.Require().Equal(2, sym.Len())
	s.True(sym.Has(1))
	s.True(sym.Has(4))
}

func (s *collectionSuite) TestSymmetricDifferenceEmpty() {
	a := newABC()
	empty := New[int, string]()
	sym := a.SymmetricDifference(empty)
	s.Equal(3, sym.Len())
}

// --- Iteration / Conversion ---

func (s *collectionSuite) TestKeys() {
	c := newABC()
	keys := c.Keys()
	s.Equal([]int{3, 1, 2}, keys)
	keys[0] = 999
	s.Equal(3, c.Keys()[0])
}

func (s *collectionSuite) TestValues() {
	c := newABC()
	vals := c.Values()
	s.Equal([]string{"gamma", "alpha", "beta"}, vals)
}

func (s *collectionSuite) TestEntries() {
	c := newABC()
	entries := c.Entries()
	s.Require().Len(entries, 3)
	s.Equal(3, entries[0].Key)
	s.Equal("gamma", entries[0].Value)
	s.Equal(1, entries[1].Key)
	s.Equal("alpha", entries[1].Value)
}

func (s *collectionSuite) TestToMap() {
	c := newABC()
	m := c.ToMap()
	s.Require().Len(m, 3)
	s.Equal("alpha", m[1])
	s.Equal("beta", m[2])
	s.Equal("gamma", m[3])
	m[99] = "zzz"
	s.False(c.Has(99))
}

func (s *collectionSuite) TestAll() {
	c := newABC()
	var keys []int
	var vals []string
	for k, v := range c.All() {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	s.Equal([]int{3, 1, 2}, keys)
	s.Equal([]string{"gamma", "alpha", "beta"}, vals)
}

func (s *collectionSuite) TestAllBreak() {
	c := newABC()
	count := 0
	for range c.All() {
		count++
		break
	}
	s.Equal(1, count)
}

func (s *collectionSuite) TestKeysIter() {
	c := newABC()
	var keys []int
	for k := range c.KeysIter() {
		keys = append(keys, k)
	}
	s.Equal([]int{3, 1, 2}, keys)
}

func (s *collectionSuite) TestValuesIter() {
	c := newABC()
	var vals []string
	for v := range c.ValuesIter() {
		vals = append(vals, v)
	}
	s.Equal([]string{"gamma", "alpha", "beta"}, vals)
}

// --- Top-Level Functions ---

func (s *collectionSuite) TestMap() {
	c := newABC()
	lengths := Map(c, func(v string) int { return len(v) })
	s.Require().Len(lengths, 3)
	s.Equal([]int{5, 5, 4}, lengths)
}

func (s *collectionSuite) TestMapToCollection() {
	c := newABC()
	upper := MapToCollection(c, strings.ToUpper)
	s.Require().Equal(3, upper.Len())
	v, _ := upper.Get(1)
	s.Equal("ALPHA", v)
	keys := upper.Keys()
	s.Equal([]int{3, 1, 2}, keys)
	orig, _ := c.Get(1)
	s.Equal("alpha", orig)
}

func (s *collectionSuite) TestFlatMap() {
	c := New[int, string]()
	c.Set(1, "ab")
	c.Set(2, "cd")
	result := FlatMap(c, func(v string) []string {
		chars := make([]string, len(v))
		for i, ch := range v {
			chars[i] = string(ch)
		}
		return chars
	})
	s.Require().Len(result, 4)
	s.Equal([]string{"a", "b", "c", "d"}, result)
}

func (s *collectionSuite) TestGroupBy() {
	c := New[int, string]()
	c.Set(1, "apple")
	c.Set(2, "apricot")
	c.Set(3, "banana")
	c.Set(4, "blueberry")
	groups := GroupBy(c, func(v string) string { return string(v[0]) })
	s.Require().Len(groups, 2)
	s.Equal(2, groups["a"].Len())
	s.Equal(2, groups["b"].Len())
}

func (s *collectionSuite) TestReduce() {
	c := newABC()
	total := Reduce(c, 0, func(acc int, v string) int { return acc + len(v) })
	s.Equal(14, total)
}

func (s *collectionSuite) TestReduceEmpty() {
	c := New[int, string]()
	total := Reduce(c, 42, func(acc int, v string) int { return acc + 1 })
	s.Equal(42, total)
}

// --- Pointer-value reflection ---

func (s *collectionSuite) TestPointerValueReflected() {
	type user struct{ Name string }
	c := New[int, *user]()
	u := &user{Name: "alice"}
	c.Set(1, u)
	u.Name = "bob"
	got, _ := c.Get(1)
	s.Equal("bob", got.Name)
}

func (s *collectionSuite) TestCollectionNilGuards() {
	eq := func(a, b string) bool { return a == b }

	populated := New[int, string]()
	populated.Set(1, "a")

	empty := New[int, string]()

	var nilColl *Collection[int, string]

	s.T().Run("Equals_nil_vs_populated", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Equals(nilColl, eq)
			assert.False(t, got, "non-empty collection must not equal nil")
		})
	})

	s.T().Run("Equals_nil_vs_empty", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := empty.Equals(nilColl, eq)
			assert.True(t, got, "empty collection must equal nil (both empty)")
		})
	})

	s.T().Run("Concat_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Concat(nilColl)
			assert.Equal(t, 1, got.Len(), "Concat with nil must return clone of receiver")
		})
	})

	s.T().Run("Difference_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Difference(nilColl)
			assert.Equal(t, 1, got.Len(), "Difference with nil must return clone of receiver")
		})
	})

	s.T().Run("Intersection_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Intersection(nilColl)
			assert.Equal(t, 0, got.Len(), "Intersection with nil must be empty")
		})
	})

	s.T().Run("SymmetricDifference_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.SymmetricDifference(nilColl)
			assert.Equal(t, 1, got.Len(), "SymmetricDifference with nil must return clone of receiver")
		})
	})
}

type collectionSuite struct{ suite.Suite }

func TestCollectionSuite(t *testing.T) { suite.Run(t, new(collectionSuite)) }
