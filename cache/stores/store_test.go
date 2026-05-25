package stores

import (
	"sync"
	"testing"
	"time"
)

// ── TestBaseStore_GetSet ──────────────────────────────────────────────────────

func TestBaseStore_GetSet(t *testing.T) {
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

// ── TestBaseStore_Delete ──────────────────────────────────────────────────────

func TestBaseStore_Delete(t *testing.T) {
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

func TestBaseStore_Has(t *testing.T) {
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

func TestBaseStore_Len(t *testing.T) {
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

func TestBaseStore_Clear(t *testing.T) {
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

func TestBaseStore_TTL_Expiry(t *testing.T) {
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

func TestBaseStore_TTL_RefreshOnAccess(t *testing.T) {
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

func TestBaseStore_MaxItems_EvictLRU(t *testing.T) {
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

func TestBaseStore_MaxItems_EvictLFU(t *testing.T) {
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

func TestBaseStore_MaxItems_EvictFIFO(t *testing.T) {
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

func TestBaseStore_PreservesMetadataOnUpdate(t *testing.T) {
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

func TestBaseStore_Sweeper_Removes(t *testing.T) {
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

func TestBaseStore_Close_StopsSweepers(t *testing.T) {
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

func TestBaseStore_All_ReturnsSnapshot(t *testing.T) {
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

func TestBaseStore_All_FiltersExpired(t *testing.T) {
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

func TestBaseStore_ReplaceKeys(t *testing.T) {
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

func TestBaseStore_ReplaceKeys_RespectsMaxItems(t *testing.T) {
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

func TestBaseStore_Disabled(t *testing.T) {
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

func TestBaseStore_ConcurrentSafe(t *testing.T) {
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
func TestBaseStore_CloseIdempotent(t *testing.T) {
	s := NewBaseStore[string, int](Defaults())
	// Must not panic even when called twice.
	s.Close()
	s.Close()
}

func TestMemMessageStore_CloseIdempotent(t *testing.T) {
	s := NewMessageStore(MessageDefaults())
	s.Close()
	s.Close()
}

