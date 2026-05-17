package cache

import (
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestBug4EvictNZeroIsNoOp verifies that evictN(0) on a genericStore removes
// no entries. This is the primitive that the LRM-based evictGloballyByBytes
// relies on: stores whose proportional allocation rounds to zero floor entries
// and that don't receive an extra slot must not be touched (Bug 4).
func TestBug4EvictNZeroIsNoOp(t *testing.T) {
	s := newGenericStore[string, *discord.Guild](storeCfg{trackBytes: true})
	s.set("a", &discord.Guild{ID: "a"})
	before := s.size()
	freed := s.evictN(0, ClearByLastUsed)
	after := s.size()
	_ = freed
	if before != after {
		t.Errorf("evictN(0) removed %d entries — must be a no-op (Bug 4)", before-after)
	}
}

// TestBug10SetAllEnforcesMaxItemsAfterUnlock verifies that SetAll on stores
// with a non-zero maxItems cap evicts excess entries after the bulk insert.
func TestBug10SetAllEnforcesMaxItemsAfterUnlock(t *testing.T) {
	// StickerStore uses MaxStickers which wires into storeCfg.maxItems.
	// Wire a cap of 2 and call SetAll with 4 stickers.
	ss := &memStickerStore{
		s: newGenericStore[discord.Snowflake, stickerValue](storeCfg{maxItems: 2, clearBy: ClearByLastUsed}),
	}

	guildID := discord.Snowflake("g1")
	stickers := []*discord.Sticker{
		{ID: "s1", Name: "a"},
		{ID: "s2", Name: "b"},
		{ID: "s3", Name: "c"},
		{ID: "s4", Name: "d"},
	}

	ss.SetAll(guildID, stickers)

	if n := ss.Size(); n > 2 {
		t.Errorf("expected ≤ 2 stickers after SetAll with maxItems=2, got %d (Bug 10)", n)
	}
}

// TestBug10SetAllIsAtomicNoEmptyWindowForReaders verifies that concurrent reads
// during SetAll never observe a zero-entry state. Without Bug 10's fix, the lock
// is released between the delete loop and the insert loop, leaving a window
// where GetByGuild returns an empty slice.
func TestBug10SetAllIsAtomicNoEmptyWindowForReaders(t *testing.T) {
	es := &memEmojiStore{
		s: newGenericStore[discord.Snowflake, emojiValue](storeCfg{}),
	}

	guildID := discord.Snowflake("123456789012345")
	// Seed two emojis.
	es.s.set("e1", emojiValue{GuildID: guildID, Emoji: &discord.Emoji{ID: "e1"}})
	es.s.set("e2", emojiValue{GuildID: guildID, Emoji: &discord.Emoji{ID: "e2"}})

	sawEmpty := make(chan struct{}, 1)
	stop := make(chan struct{})

	// Reader goroutine: continuously polls GetByGuild.
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				if len(es.GetByGuild(guildID)) == 0 {
					select {
					case sawEmpty <- struct{}{}:
					default:
					}
				}
			}
		}
	}()

	// Writer: repeatedly calls SetAll with the same two emojis.
	const iters = 500
	var wg sync.WaitGroup
	wg.Add(iters)
	for i := 0; i < iters; i++ {
		go func() {
			defer wg.Done()
			es.SetAll(guildID, []*discord.Emoji{
				{ID: "e1"},
				{ID: "e2"},
			})
		}()
	}
	wg.Wait()
	close(stop)

	select {
	case <-sawEmpty:
		t.Error("reader saw empty emoji list during SetAll — delete→insert not atomic (Bug 10)")
	default:
	}
}

// TestBug11TotalBytesNeverNegative verifies that the subtract(-old) and add(+new)
// operations on totalBytes happen inside the same lock, so no concurrent reader
// can observe a transient undercount between the two ops.
func TestBug11TotalBytesNeverNegative(t *testing.T) {
	s := newGenericStore[string, *discord.Guild](storeCfg{trackBytes: true})

	// Seed one entry so every subsequent Set is an update (subtract+add path).
	s.set("1", &discord.Guild{ID: "1", Name: "initial"})

	negative := make(chan int64, 1)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				if v := s.totalBytes.Load(); v < 0 {
					select {
					case negative <- v:
					default:
					}
				}
			}
		}
	}()

	const goroutines = 500
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s.set("1", &discord.Guild{ID: "1", Name: "updated"})
		}()
	}
	wg.Wait()
	close(stop)

	select {
	case v := <-negative:
		t.Errorf("totalBytes went negative (%d) during concurrent updates — non-atomic byte accounting (Bug 11)", v)
	default:
	}
}

// ── findEvictionCandidates correctness tests ──────────────────────────────────

// naiveFindCandidates is a reference implementation using full sort, used to
// validate the heap-based result.
func naiveFindCandidates(s *genericStore[string, *discord.Guild], n int, by ClearBy) []evictionPair[string] {
	type scored struct {
		key        string
		score      float64
		insertedAt int64
		size       int64
	}
	all := make([]scored, 0, len(s.items))
	for k, e := range s.items {
		all = append(all, scored{k, e.rank(by), e.insertedAt.UnixNano(), e.sizeBytes})
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].score != all[j].score {
			return all[i].score < all[j].score
		}
		return all[i].insertedAt < all[j].insertedAt
	})
	if n > len(all) {
		n = len(all)
	}
	result := make([]evictionPair[string], n)
	for i, a := range all[:n] {
		result[i] = evictionPair[string]{a.key, a.score, a.insertedAt, a.size}
	}
	return result
}

// candidateKeys extracts and sorts the keys from a candidate list for
// order-independent comparison.
func candidateKeys(ps []evictionPair[string]) []string {
	keys := make([]string, len(ps))
	for i, p := range ps {
		keys[i] = p.key
	}
	sort.Strings(keys)
	return keys
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestFindEvictionCandidates_MatchesNaiveSort verifies that the heap-based
// implementation selects the same n lowest-rank keys as a full sort.
func TestFindEvictionCandidates_MatchesNaiveSort(t *testing.T) {
	for _, n := range []int{10, 100, 1000} {
		total := 2000
		s := newGenericStore[string, *discord.Guild](storeCfg{clearBy: ClearByAge})
		now := time.Now()
		for i := 0; i < total; i++ {
			key := strconv.Itoa(i)
			e := &memEntry[*discord.Guild]{
				value:      &discord.Guild{},
				insertedAt: now.Add(time.Duration(i) * time.Millisecond),
				accessedAt: now,
			}
			s.items[key] = e
		}

		got := candidateKeys(s.findEvictionCandidates(n, ClearByAge))
		want := candidateKeys(naiveFindCandidates(s, n, ClearByAge))

		if len(got) != n {
			t.Errorf("n=%d: got %d candidates, want %d", n, len(got), n)
		}
		if !stringSliceEqual(got, want) {
			t.Errorf("n=%d: heap result does not match naive sort result", n)
		}
	}
}

// TestFindEvictionCandidates_Tiebreaker verifies that when scores are equal,
// older entries (smaller insertedAt) are preferred for eviction.
func TestFindEvictionCandidates_Tiebreaker(t *testing.T) {
	s := newGenericStore[string, *discord.Guild](storeCfg{clearBy: ClearByAge})
	base := time.Now()
	// All entries have the same rank score; older ones should be evicted first.
	s.items["old"] = &memEntry[*discord.Guild]{value: &discord.Guild{}, insertedAt: base, accessedAt: base}
	s.items["mid"] = &memEntry[*discord.Guild]{value: &discord.Guild{}, insertedAt: base.Add(time.Second), accessedAt: base}
	s.items["new"] = &memEntry[*discord.Guild]{value: &discord.Guild{}, insertedAt: base.Add(2 * time.Second), accessedAt: base}

	// ClearByLastUsed: all have the same accessedAt, so tiebreaker by insertedAt.
	got := s.findEvictionCandidates(1, ClearByLastUsed)
	if len(got) != 1 || got[0].key != "old" {
		t.Errorf("tiebreaker: expected oldest key 'old', got %+v", got)
	}
}

// TestFindEvictionCandidates_NGreaterThanItems verifies that requesting more
// candidates than items returns all items without panic.
func TestFindEvictionCandidates_NGreaterThanItems(t *testing.T) {
	s := newGenericStore[string, *discord.Guild](storeCfg{})
	s.items["a"] = &memEntry[*discord.Guild]{value: &discord.Guild{}, insertedAt: time.Now(), accessedAt: time.Now()}
	s.items["b"] = &memEntry[*discord.Guild]{value: &discord.Guild{}, insertedAt: time.Now(), accessedAt: time.Now()}

	got := s.findEvictionCandidates(100, ClearByLastUsed)
	if len(got) != 2 {
		t.Errorf("expected 2 candidates for n=100 with 2 items, got %d", len(got))
	}
}

// TestFindEvictionCandidates_ZeroOrEmpty verifies that edge cases return nil
// without panicking.
func TestFindEvictionCandidates_ZeroOrEmpty(t *testing.T) {
	s := newGenericStore[string, *discord.Guild](storeCfg{})

	if got := s.findEvictionCandidates(0, ClearByLastUsed); got != nil {
		t.Errorf("n=0: expected nil, got %v", got)
	}
	if got := s.findEvictionCandidates(5, ClearByLastUsed); got != nil {
		t.Errorf("empty store: expected nil, got %v", got)
	}
}
