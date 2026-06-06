package stores

import (
	"time"

	"github.com/stretchr/testify/assert"
)

// ── EvictionStrategy.String ───────────────────────────────────────────────────

func (su *storesStoreSuite) TestEvictionStrategy_String() {
	t := su.T()
	assert.Equal(t, "LRU", EvictLRU.String())
	assert.Equal(t, "LFU", EvictLFU.String())
	assert.Equal(t, "FIFO", EvictFIFO.String())
	assert.Equal(t, "unknown", EvictionStrategy(99).String())
}

// ── estimateSize error path ───────────────────────────────────────────────────

func (su *storesStoreSuite) TestEstimateSize_UnmarshalableReturnsZero() {
	t := su.T()
	// channels cannot be JSON-encoded → json.Marshal errors → size 0.
	assert.EqualValues(t, 0, estimateSize(make(chan int)))
	assert.Positive(t, estimateSize("hello"))
}

// ── TrackBytes accounting on Set / Delete / ReplaceKeys ───────────────────────

func (su *storesStoreSuite) TestBaseStore_TrackBytes_SetReplaceDelete() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{TrackBytes: true})
	defer s.Close()

	s.Set("a", 1)
	first := s.Bytes()
	assert.Positive(t, first)

	// Replace same key → old bytes subtracted, new added (no leak).
	s.Set("a", 222222)
	assert.NotEqual(t, first, s.Bytes())

	// Delete with TrackBytes subtracts the entry's bytes.
	s.Delete("a")
	assert.EqualValues(t, 0, s.Bytes())
}

func (su *storesStoreSuite) TestBaseStore_ReplaceKeys_TrackBytes_DeleteBranch() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{TrackBytes: true})
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)
	assert.Positive(t, s.Bytes())

	// del path with TrackBytes (subtract) + add path.
	s.ReplaceKeys([]string{"a"}, map[string]int{"b": 99, "c": 3})
	assert.Positive(t, s.Bytes())
	assert.Equal(t, 2, s.Len())
}

// ── evictOne: onEvict + TrackBytes, and empty-store guard ─────────────────────

func (su *storesStoreSuite) TestBaseStore_EvictOne_OnEvictAndTrackBytes() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		MaxItems:         1,
		TrackBytes:       true,
		EvictionStrategy: EvictLRU,
	})
	defer s.Close()

	var evicted []string
	s.onEvict = func(k string) { evicted = append(evicted, k) }

	s.Set("a", 1)
	time.Sleep(time.Millisecond)
	s.Set("b", 2) // over MaxItems → evictOne removes "a" with onEvict + TrackBytes

	assert.Equal(t, []string{"a"}, evicted)
	assert.Equal(t, 1, s.Len())
}

func (su *storesStoreSuite) TestBaseStore_EvictOne_EmptyStoreNoop() {
	su.T() // empty store: evictOne must not panic and must not call onEvict.
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()
	called := false
	s.onEvict = func(string) { called = true }
	s.evictOne()
	assert.False(su.T(), called)
}

// ── SweepExpired: onEvict + TrackBytes + idle window ──────────────────────────

func (su *storesStoreSuite) TestBaseStore_SweepExpired_OnEvictTrackBytesIdle() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		TTL:        10 * time.Millisecond,
		TrackBytes: true,
	})
	defer s.Close()

	var evicted []string
	s.onEvict = func(k string) { evicted = append(evicted, k) }

	s.Set("a", 1)
	s.Set("b", 2)
	time.Sleep(25 * time.Millisecond)

	s.SweepExpired(0)
	assert.Len(t, evicted, 2)
	assert.EqualValues(t, 0, s.Bytes())
	assert.Equal(t, 0, s.Len())
}

func (su *storesStoreSuite) TestBaseStore_SweepExpired_IdleWindow() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	time.Sleep(15 * time.Millisecond)
	// Idle window shorter than the idle time → "a" is swept.
	s.SweepExpired(5 * time.Millisecond)
	assert.Equal(t, 0, s.Len())
}

// ── runSweeper onEvict + startSweeper invalid guards ──────────────────────────

func (su *storesStoreSuite) TestBaseStore_RunSweeper_OnEvict() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	var evicted []string
	s.onEvict = func(k string) { evicted = append(evicted, k) }

	s.Set("keep", 1)
	s.Set("drop", -1)
	s.runSweeper(Sweeper{Filter: func(v any) bool { return v.(int) < 0 }})

	assert.Equal(t, []string{"drop"}, evicted)
	assert.Equal(t, 1, s.Len())
}

func (su *storesStoreSuite) TestBaseStore_StartSweeper_InvalidIsIgnored() {
	t := su.T()
	// Interval <= 0 and nil Filter are both rejected by startSweeper, so the
	// store starts with no sweeper goroutines.
	s := NewBaseStore[string, int](StoreOptions{
		Sweepers: []Sweeper{
			{Interval: 0, Filter: func(any) bool { return true }},
			{Interval: time.Millisecond, Filter: nil},
		},
	})
	s.Set("a", 1)
	time.Sleep(10 * time.Millisecond)
	s.Close()
	assert.Equal(t, 1, s.Len(), "no sweeper should have run")
}
