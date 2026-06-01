package stores

import (
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

// ── TestBaseStore_GetSet ──────────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_GetSet() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)

	if v, ok := s.Get("a"); !ok || v != 1 {
		t.Fatalf("expected (1, true), got (%d, %v)", v, ok)
	}
	if v, ok := s.Get("b"); !ok || v != 2 {
		t.Fatalf("expected (2, true), got (%d, %v)", v, ok)
	}
	if _, ok := s.Get("missing"); ok {
		t.Fatal("expected miss for absent key")
	}
}

func (su *storesStoreSuite) TestBaseStoreBug60_GetLazyRemoval() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		TTL:        20 * time.Millisecond,
		TrackBytes: true,
	})
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)

	bytesAfterSet := s.totalBytes.Load()
	if bytesAfterSet <= 0 {
		t.Fatalf("expected totalBytes > 0 after two Set calls, got %d", bytesAfterSet)
	}

	time.Sleep(40 * time.Millisecond)

	// Get on "a": entry is expired → lazy delete, totalBytes must shrink.
	sizeA := estimateSize(1)
	if _, ok := s.Get("a"); ok {
		t.Fatal("entry a should have expired")
	}

	wantAfterA := bytesAfterSet - sizeA
	if got := s.totalBytes.Load(); got != wantAfterA {
		t.Fatalf("after lazy removal of a: totalBytes = %d, want %d", got, wantAfterA)
	}

	// Get on "b": same check.
	sizeB := estimateSize(2)
	if _, ok := s.Get("b"); ok {
		t.Fatal("entry b should have expired")
	}

	wantAfterB := wantAfterA - sizeB
	if got := s.totalBytes.Load(); got != wantAfterB {
		t.Fatalf("after lazy removal of b: totalBytes = %d, want %d", got, wantAfterB)
	}

	if s.Len() != 0 {
		t.Fatalf("expected Len=0 after both lazy removals, got %d", s.Len())
	}
}

// ── TestBaseStore_Delete ──────────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Delete() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	if !s.Delete("a") {
		t.Fatal("expected Delete to return true for present key")
	}
	if s.Delete("a") {
		t.Fatal("expected Delete to return false for absent key")
	}
	if _, ok := s.Get("a"); ok {
		t.Fatal("expected miss after delete")
	}
}

// ── TestBaseStore_Has ─────────────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Has() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	if !s.Has("a") {
		t.Fatal("expected Has=true for present key")
	}
	if s.Has("missing") {
		t.Fatal("expected Has=false for absent key")
	}

	// Has must NOT update access metadata; verify hitCount is unchanged.
	// Get to establish baseline then call Has multiple times.
	s.Get("a") // hitCount becomes 1
	for i := 0; i < 10; i++ {
		s.Has("a")
	}
	e, _ := s.items.Get("a")
	if e.hitCount != 1 {
		t.Fatalf("Has should not increment hitCount; got %d", e.hitCount)
	}
}

// ── TestBaseStore_Len ─────────────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Len() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	if s.Len() != 0 {
		t.Fatalf("expected Len=0, got %d", s.Len())
	}
	s.Set("a", 1)
	s.Set("b", 2)
	if s.Len() != 2 {
		t.Fatalf("expected Len=2, got %d", s.Len())
	}
	s.Delete("a")
	if s.Len() != 1 {
		t.Fatalf("expected Len=1 after delete, got %d", s.Len())
	}
}

// ── TestBaseStore_Clear ───────────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Clear() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)
	s.Clear()

	if s.Len() != 0 {
		t.Fatalf("expected Len=0 after Clear, got %d", s.Len())
	}
}

// ── TestBaseStore_TTL_Expiry ──────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_TTL_Expiry() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.Set("a", 1)

	if _, ok := s.Get("a"); !ok {
		t.Fatal("entry should be present before TTL elapses")
	}

	time.Sleep(40 * time.Millisecond)

	if _, ok := s.Get("a"); ok {
		t.Fatal("entry should have expired")
	}
	// Lazy cleanup: expired entry is removed on Get.
	if s.Len() != 0 {
		t.Fatalf("expected Len=0 after expiry cleanup, got %d", s.Len())
	}
}

// ── TestBaseStore_TTL_RefreshOnAccess ─────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_TTL_RefreshOnAccess() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		TTL:             40 * time.Millisecond,
		RefreshOnAccess: true,
	})
	defer s.Close()

	s.Set("a", 1)

	// Access the entry repeatedly before TTL would naturally expire.
	for i := 0; i < 3; i++ {
		time.Sleep(20 * time.Millisecond)
		if _, ok := s.Get("a"); !ok {
			t.Fatalf("entry should still be alive after refresh (iteration %d)", i)
		}
	}

	// Now stop refreshing and let the TTL run out.
	time.Sleep(60 * time.Millisecond)
	if _, ok := s.Get("a"); ok {
		t.Fatal("entry should have expired after refresh stopped")
	}
}

// ── TestBaseStore_MaxItems_EvictLRU ───────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_MaxItems_EvictLRU() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictLRU,
	})
	defer s.Close()

	s.Set("a", 1)
	time.Sleep(time.Millisecond)
	s.Set("b", 2)

	// Access a to make it more recently used than b.
	time.Sleep(time.Millisecond)
	s.Get("a")

	// Adding c triggers eviction; b has the oldest accessedAt.
	time.Sleep(time.Millisecond)
	s.Set("c", 3)

	if _, ok := s.Get("b"); ok {
		t.Fatal("b should have been evicted (LRU)")
	}
	if _, ok := s.Get("a"); !ok {
		t.Fatal("a should still be present")
	}
	if _, ok := s.Get("c"); !ok {
		t.Fatal("c should be present")
	}
}

// ── TestBaseStore_MaxItems_EvictLFU ───────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_MaxItems_EvictLFU() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictLFU,
	})
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)

	// Access a twice so it has higher hit count than b (hitCount=0).
	s.Get("a")
	s.Get("a")

	// Adding c triggers eviction; b has the lowest hit count (0).
	s.Set("c", 3)

	if _, ok := s.Get("b"); ok {
		t.Fatal("b should have been evicted (LFU)")
	}
	if _, ok := s.Get("a"); !ok {
		t.Fatal("a should still be present (high hit count)")
	}
	if _, ok := s.Get("c"); !ok {
		t.Fatal("c should be present")
	}
}

// ── TestBaseStore_MaxItems_EvictFIFO ──────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_MaxItems_EvictFIFO() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictFIFO,
	})
	defer s.Close()

	s.Set("a", 1)
	time.Sleep(time.Millisecond)
	s.Set("b", 2)

	// Access a to make it the most recently used — FIFO should ignore access time.
	s.Get("a")

	// Adding c triggers eviction; a has the earliest insertedAt.
	time.Sleep(time.Millisecond)
	s.Set("c", 3)

	if _, ok := s.Get("a"); ok {
		t.Fatal("a should have been evicted (FIFO: earliest insertedAt)")
	}
	if _, ok := s.Get("b"); !ok {
		t.Fatal("b should still be present")
	}
	if _, ok := s.Get("c"); !ok {
		t.Fatal("c should be present")
	}
}

// ── TestBaseStore_PreservesMetadataOnUpdate ───────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_PreservesMetadataOnUpdate() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictLFU,
	})
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)

	// Give a a high hit count.
	s.Get("a")
	s.Get("a")

	// Update a — insertedAt and hitCount must be preserved.
	s.Set("a", 100)

	e, ok := s.items.Get("a")
	if !ok {
		t.Fatal("a should be present after update")
	}
	if e.hitCount != 2 {
		t.Fatalf("hitCount should be preserved as 2 after update, got %d", e.hitCount)
	}

	// Adding c triggers LFU eviction; b has hitCount=0 < a's hitCount=2.
	s.Set("c", 3)

	if _, ok := s.Get("b"); ok {
		t.Fatal("b should have been evicted (LFU, hitCount=0)")
	}
	if v, ok := s.Get("a"); !ok || v != 100 {
		t.Fatalf("a should be present with updated value 100, got (%d, %v)", v, ok)
	}

	// Also verify insertedAt is preserved by testing FIFO behaviour after update.
	s2 := NewBaseStore[string, int](StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictFIFO,
	})
	defer s2.Close()

	s2.Set("x", 1)
	time.Sleep(time.Millisecond)
	s2.Set("y", 2)

	insertedBefore, _ := s2.items.Get("x")
	origInsertedAt := insertedBefore.insertedAt

	s2.Set("x", 99) // update x

	insertedAfter, _ := s2.items.Get("x")
	if !insertedAfter.insertedAt.Equal(origInsertedAt) {
		t.Fatal("insertedAt should be preserved on update")
	}

	// FIFO eviction: x has the earlier insertedAt, so x is evicted.
	time.Sleep(time.Millisecond)
	s2.Set("z", 3)

	if _, ok := s2.Get("x"); ok {
		t.Fatal("x should have been evicted (FIFO: earliest insertedAt, preserved through update)")
	}
}

// ── TestBaseStore_Sweeper_Removes ─────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Sweeper_Removes() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		Sweepers: []Sweeper{
			{
				Interval: 5 * time.Millisecond,
				Filter: func(v any) bool {
					return v.(int) < 0
				},
				Name: "negative-values",
			},
		},
	})
	defer s.Close()

	s.Set("keep1", 1)
	s.Set("remove", -1)
	s.Set("keep2", 2)

	time.Sleep(25 * time.Millisecond)

	if _, ok := s.Get("remove"); ok {
		t.Fatal("remove should have been swept (negative value)")
	}
	if _, ok := s.Get("keep1"); !ok {
		t.Fatal("keep1 should still be present")
	}
	if _, ok := s.Get("keep2"); !ok {
		t.Fatal("keep2 should still be present")
	}
}

// ── TestBaseStore_Close_StopsSweepers ─────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Close_StopsSweepers() {
	t := su.T()
	// Use a channel signal instead of a sleep+count so the test is reliable
	// even under the race detector's slower goroutine scheduling.
	ran := make(chan struct{}, 1)

	s := NewBaseStore[string, int](StoreOptions{
		Sweepers: []Sweeper{
			{
				Interval: 10 * time.Millisecond,
				Filter: func(v any) bool {
					select {
					case ran <- struct{}{}:
					default:
					}
					return false
				},
			},
		},
	})
	// Sweeper Filter is only invoked on actual entries; add a sentinel item.
	s.Set("sentinel", 0)

	// Wait until the sweeper has run at least once.
	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		s.Close()
		t.Fatal("sweeper never ran within 2s")
	}

	s.Close()

	// Drain any buffered signal from before Close, then confirm no new ones arrive.
	select {
	case <-ran:
	default:
	}
	time.Sleep(50 * time.Millisecond)

	select {
	case <-ran:
		t.Fatal("sweeper fired after Close")
	default:
		// good
	}
}

// ── TestBaseStore_All_ReturnsSnapshot ─────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_All_ReturnsSnapshot() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)

	snap := s.All()
	if snap.Len() != 2 {
		t.Fatalf("snapshot should have 2 entries, got %d", snap.Len())
	}

	// Mutating the snapshot must not affect the store.
	snap.Set("c", 3)

	if s.Len() != 2 {
		t.Fatalf("store should still have 2 entries after snapshot mutation, got %d", s.Len())
	}
	if _, ok := s.Get("c"); ok {
		t.Fatal("snapshot mutation must not bleed into store")
	}
}

// ── TestBaseStore_All_FiltersExpired ──────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_All_FiltersExpired() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.Set("a", 1)
	time.Sleep(40 * time.Millisecond)
	s.Set("b", 2) // inserted after the sleep, so not yet expired

	snap := s.All()
	if snap.Len() != 1 {
		t.Fatalf("All should return only live entries; got %d entries", snap.Len())
	}
	if v, ok := snap.Get("b"); !ok || v != 2 {
		t.Fatal("b should be the sole live entry in the snapshot")
	}
}

// ── TestBaseStore_ReplaceKeys ─────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_ReplaceKeys() {
	t := su.T()
	s := NewBaseStore[string, int](Defaults())
	defer s.Close()

	s.Set("a", 1)
	s.Set("b", 2)
	s.Set("c", 3)

	s.ReplaceKeys([]string{"a", "b"}, map[string]int{"d": 4, "e": 5})

	if _, ok := s.Get("a"); ok {
		t.Error("a should be deleted")
	}
	if _, ok := s.Get("b"); ok {
		t.Error("b should be deleted")
	}
	if v, ok := s.Get("c"); !ok || v != 3 {
		t.Errorf("c should be untouched, got (%d, %v)", v, ok)
	}
	if v, ok := s.Get("d"); !ok || v != 4 {
		t.Errorf("d should be inserted, got (%d, %v)", v, ok)
	}
	if v, ok := s.Get("e"); !ok || v != 5 {
		t.Errorf("e should be inserted, got (%d, %v)", v, ok)
	}
	if s.Len() != 3 { // c, d, e
		t.Fatalf("expected Len=3 (c,d,e), got %d", s.Len())
	}
}

func (su *storesStoreSuite) TestBaseStore_ReplaceKeys_RespectsMaxItems() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictFIFO,
	})
	defer s.Close()

	// Replace with 3 items into a MaxItems=2 store; one must be evicted.
	s.ReplaceKeys(nil, map[string]int{"a": 1, "b": 2, "c": 3})

	if s.Len() != 2 {
		t.Fatalf("expected Len=2 after eviction, got %d", s.Len())
	}
}

// ── TestBaseStore_Disabled ────────────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_Disabled() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{Disabled: true})
	defer s.Close()

	s.Set("a", 1)

	if _, ok := s.Get("a"); ok {
		t.Fatal("disabled store should not return values")
	}
	if s.Has("a") {
		t.Fatal("disabled store Has should always return false")
	}
	if s.Len() != 0 {
		t.Fatalf("disabled store Len should be 0, got %d", s.Len())
	}
	if !s.All().IsEmpty() {
		t.Fatal("disabled store All should return empty collection")
	}
	if s.IsEnabled() {
		t.Fatal("disabled store IsEnabled should return false")
	}

	s.ReplaceKeys(nil, map[string]int{"b": 2})
	if s.Len() != 0 {
		t.Fatal("disabled store should ignore ReplaceKeys")
	}
}

// ── TestBaseStore_ConcurrentSafe ──────────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_ConcurrentSafe() {
	s := NewBaseStore[int, int](StoreOptions{
		MaxItems: 50,
		Sweepers: []Sweeper{
			{
				Interval: time.Millisecond,
				Filter:   func(v any) bool { return v.(int)%7 == 0 },
			},
		},
	})
	defer s.Close()

	const goroutines = 50
	const ops = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				key := (g*ops + i) % 20
				s.Set(key, g*ops+i)
				s.Get(key)
				s.Has(key)
				s.Delete(key % 5)
				s.Len()
				s.All()
			}
		}(g)
	}

	wg.Wait()
}

// J0-#12: Close must be idempotent; a second call must not panic.
func (su *storesStoreSuite) TestBaseStore_CloseIdempotent() {
	s := NewBaseStore[string, int](Defaults())
	// Must not panic even when called twice.
	s.Close()
	s.Close()
}

func (su *storesStoreSuite) TestMemMessageStore_CloseIdempotent() {
	s := NewMessageStore(MessageDefaults())
	s.Close()
	s.Close()
}

// ── Bug #94 ───────────────────────────────────────────────────────────────────

// TestBug94_ReplaceKeys_PreservesEvictionMetadata verifies that ReplaceKeys
// preserves hitCount and insertedAt for keys present in both the old and new sets.
//
// Root cause: ReplaceKeys always created fresh entries, resetting LFU/FIFO
// eviction metadata. It also leaked totalBytes because the old entry's sizeBytes
// was never subtracted before re-adding.
func (su *storesStoreSuite) TestBug94_ReplaceKeys_PreservesEvictionMetadata() {
	t := su.T()
	t.Run("LFU: hitCount preserved", func(t *testing.T) {
		s := NewBaseStore[string, int](StoreOptions{
			MaxItems:         2,
			EvictionStrategy: EvictLFU,
		})
		defer s.Close()

		s.Set("a", 1)
		s.Set("b", 2)
		s.Get("a")
		s.Get("a")
		s.Get("a") // hitCount("a") = 3

		// ReplaceKeys: update "a" (kept), remove "b".
		s.ReplaceKeys([]string{"b"}, map[string]int{"a": 99})

		e, ok := s.items.Get("a")
		if !ok {
			t.Fatal("a should still exist after ReplaceKeys")
		}
		if e.hitCount != 3 {
			t.Fatalf("Bug94: hitCount should be preserved as 3, got %d (old code reset to 0)", e.hitCount)
		}
		if e.value != 99 {
			t.Fatalf("value should be updated to 99, got %d", e.value)
		}

		// Add "c" — LFU must evict "c" (hitCount 0), not "a" (hitCount 3).
		s.Set("c", 3)
		if _, ok := s.Get("a"); !ok {
			t.Fatal("Bug94: a (hitCount=3) should survive LFU eviction after ReplaceKeys")
		}
	})

	t.Run("FIFO: insertedAt preserved", func(t *testing.T) {
		s := NewBaseStore[string, int](StoreOptions{
			MaxItems:         2,
			EvictionStrategy: EvictFIFO,
		})
		defer s.Close()

		s.Set("x", 1)
		time.Sleep(time.Millisecond)
		s.Set("y", 2)

		orig, _ := s.items.Get("x")
		origInsertedAt := orig.insertedAt

		// ReplaceKeys updates both keys — insertedAt must not change.
		s.ReplaceKeys(nil, map[string]int{"x": 10, "y": 20})

		after, _ := s.items.Get("x")
		if !after.insertedAt.Equal(origInsertedAt) {
			t.Fatal("Bug94: insertedAt should be preserved through ReplaceKeys")
		}

		// FIFO: x (older) must be evicted when a third item is added.
		time.Sleep(time.Millisecond)
		s.Set("z", 30)
		if _, ok := s.Get("x"); ok {
			t.Fatal("Bug94: x should be evicted (FIFO: earliest insertedAt, preserved through ReplaceKeys)")
		}
	})

	t.Run("TrackBytes: no byte leak on update", func(t *testing.T) {
		s := NewBaseStore[string, int](StoreOptions{TrackBytes: true})
		defer s.Close()

		s.Set("a", 1)
		bytesAfterSet := s.totalBytes.Load()

		// ReplaceKeys with same key — totalBytes must not grow.
		s.ReplaceKeys(nil, map[string]int{"a": 2})

		if got := s.totalBytes.Load(); got != bytesAfterSet {
			t.Fatalf("Bug94: totalBytes leaked: was %d, got %d after in-place update via ReplaceKeys", bytesAfterSet, got)
		}
	})
}

// ── Bug #116 ──────────────────────────────────────────────────────────────────

// TestBug116_DeleteOnDisabledStore verifies that Delete on a Disabled store
// returns false immediately without acquiring the write lock.
//
// Root cause: the Disabled check was missing in Delete, causing a needless
// write-lock acquisition on every call even though the store is a no-op.
func (su *storesStoreSuite) TestBug116_DeleteOnDisabledStore() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{Disabled: true, TrackBytes: true})
	defer s.Close()

	if got := s.Delete("nonexistent"); got {
		t.Fatal("Bug116: Delete on disabled store should return false")
	}
	if s.totalBytes.Load() != 0 {
		t.Fatal("Bug116: disabled store should never accumulate bytes")
	}
	// Second call must also be safe (no panic, correct return value).
	if got := s.Delete("anything"); got {
		t.Fatal("Bug116: Delete on disabled store should always return false")
	}
}

// ── Bug #150 ──────────────────────────────────────────────────────────────────

// TestBug150_EvictToCount_CorrectVictims verifies that evictToCount selects the
// right victims when evicting multiple items in a single batch.
//
// Root cause: evictToCount called evictOne in a loop — O(n²) scans and N separate
// lock acquisitions. The fix collects all entries once, sorts once (O(n log n)),
// and deletes in one lock acquisition. This test verifies correctness, not
// performance; correctness requires that the highest-priority victims are chosen.
func (su *storesStoreSuite) TestBug150_EvictToCount_CorrectVictims() {
	t := su.T()
	t.Run("LFU: keeps highest hit counts across batch eviction", func(t *testing.T) {
		s := NewBaseStore[string, int](StoreOptions{
			MaxItems:         3,
			EvictionStrategy: EvictLFU,
		})
		defer s.Close()

		s.Set("low", 1)  // hitCount = 0
		s.Set("mid", 2)  // hitCount = 1
		s.Set("high", 3) // hitCount = 3
		s.Get("mid")     // mid → 1
		s.Get("high")
		s.Get("high")
		s.Get("high") // high → 3

		// ReplaceKeys adds 3 new items (all hitCount=0), total=6 → evict down to 3.
		// "high" (3) and "mid" (1) must survive; "low" (0) is a victim.
		s.ReplaceKeys(nil, map[string]int{"a": 10, "b": 20, "c": 30})

		if s.Len() != 3 {
			t.Fatalf("Bug150: expected Len=3 after batch eviction, got %d", s.Len())
		}
		if _, ok := s.Get("high"); !ok {
			t.Fatal("Bug150: 'high' (hitCount=3) must survive LFU batch eviction")
		}
		if _, ok := s.Get("mid"); !ok {
			t.Fatal("Bug150: 'mid' (hitCount=1) must survive LFU batch eviction")
		}
	})

	t.Run("FIFO: keeps newest insertions across batch eviction", func(t *testing.T) {
		s := NewBaseStore[string, int](StoreOptions{
			MaxItems:         2,
			EvictionStrategy: EvictFIFO,
		})
		defer s.Close()

		s.Set("oldest", 1)
		time.Sleep(time.Millisecond)
		s.Set("middle", 2)
		time.Sleep(time.Millisecond)

		// Add 2 new items (total=4) → evict down to 2; "oldest" must go first.
		s.ReplaceKeys(nil, map[string]int{"newest1": 10, "newest2": 20})

		if s.Len() != 2 {
			t.Fatalf("Bug150: expected Len=2 after batch eviction, got %d", s.Len())
		}
		if _, ok := s.Get("oldest"); ok {
			t.Fatal("Bug150: 'oldest' should be first victim of FIFO batch eviction")
		}
	})
}

type storesStoreSuite struct{ suite.Suite }

func TestStoresStoreSuite(t *testing.T) { suite.Run(t, new(storesStoreSuite)) }
