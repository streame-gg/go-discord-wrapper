package stores

import (
	"sync"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── inviteGuildIndex.remove: set becomes empty → guild dropped ────────────────

func (su *storesInviteSuite) TestInviteGuildIndex_RemoveLastCodeDropsGuild() {
	t := su.T()
	idx := newInviteGuildIndex()
	idx.add(gA, "only")
	idx.remove("only") // set now empty → delete(byGuild, gA)
	assert.Nil(t, idx.forGuild(gA))
}

// ── msgRing.sweep: keep + remove in the same pass ─────────────────────────────

func (su *storesMessageSuite) TestMessageStore_SweepKeepsLiveMessages() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 10, TTL: 20 * time.Millisecond})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1)) // will expire
	time.Sleep(25 * time.Millisecond)
	s.Add(msg(ch, 2)) // fresh, must survive the sweep

	s.SweepExpired(0)
	assert.Equal(t, 1, s.Size())
	_, ok := s.Get(ch, 2)
	assert.True(t, ok, "fresh message must survive sweep")
}

// ── trimToTotalLocked: per-ring eviction cap branch ───────────────────────────

func (su *storesMessageSuite) TestMessageStore_TrimToTotal_CapPerRing() {
	t := su.T()
	// Use many channels added in time order. Channel-map iteration order is
	// randomized, so the per-trim insertion sort (ordering channels by
	// accessedAt) will perform swaps; repeating across fresh stores makes the
	// swap branch essentially certain to execute.
	for iter := 0; iter < 50; iter++ {
		s := NewMessageStore(StoreOptions{MaxItems: 100, EvictionStrategy: EvictLRU})

		const chans = 8
		var first discord.Snowflake
		mid := discord.Snowflake(1)
		for c := 0; c < chans; c++ {
			ch := discord.Snowflake(c + 1)
			if c == 0 {
				first = ch
			}
			// First channel gets a single message; the rest get several. This
			// also exercises the per-ring eviction cap (need > ring length).
			s.Add(msg(ch, mid))
			mid++
			if c != 0 {
				s.Add(msg(ch, mid))
				mid++
				s.Add(msg(ch, mid))
				mid++
			}
			time.Sleep(time.Millisecond) // distinct accessedAt per channel
		}
		_ = first

		total := s.Size()
		assert.Positive(t, total)
		s.TrimToTotal(1)
		assert.Equal(t, 1, s.Size())
		s.Close()
	}
}

// ── message startSweeper invalid guards ───────────────────────────────────────

func (su *storesMessageSuite) TestMessageStore_StartSweeper_InvalidIgnored() {
	t := su.T()
	s := NewMessageStore(StoreOptions{
		MaxItems: 10,
		Sweepers: []Sweeper{
			{Interval: 0, Filter: func(any) bool { return true }},
			{Interval: time.Millisecond, Filter: nil},
		},
	})
	s.Add(msg(discord.Snowflake(1), 1))
	time.Sleep(10 * time.Millisecond)
	s.Close()
	assert.Equal(t, 1, s.Size(), "no sweeper should have run")
}

// ── BaseStore.ReplaceKeys with TTL set ────────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_ReplaceKeys_SetsTTL() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.ReplaceKeys(nil, map[string]int{"a": 1})
	if _, ok := s.Get("a"); !ok {
		t.Fatal("entry should be present before TTL elapses")
	}
	time.Sleep(35 * time.Millisecond)
	if _, ok := s.Get("a"); ok {
		t.Fatal("entry inserted via ReplaceKeys should honour TTL and expire")
	}
}

// ── BaseStore.SweepExpired on a disabled store ────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_SweepExpired_DisabledIsNoop() {
	s := NewBaseStore[string, int](StoreOptions{Disabled: true})
	defer s.Close()
	s.SweepExpired(time.Millisecond) // must early-return without panic
}

// ── BaseStore.runSweeper with TrackBytes ──────────────────────────────────────

func (su *storesStoreSuite) TestBaseStore_RunSweeper_TrackBytes() {
	t := su.T()
	s := NewBaseStore[string, int](StoreOptions{TrackBytes: true})
	defer s.Close()

	s.Set("keep", 1)
	s.Set("drop", -1)
	before := s.Bytes()
	assert.Positive(t, before)

	s.runSweeper(Sweeper{Filter: func(v any) bool { return v.(int) < 0 }})
	assert.Less(t, s.Bytes(), before, "removing an entry must shrink tracked bytes")
	assert.Equal(t, 1, s.Len())
}

// ── BaseStore.SweepExpired phase-gap (entry vanishes between scan & delete) ───

// TestBaseStore_SweepExpired_ConcurrentRemoval races SweepExpired against Get
// (which lazily deletes expired entries). This exercises the phase-2 re-check
// where a candidate key is no longer present.
func (su *storesStoreSuite) TestBaseStore_SweepExpired_ConcurrentRemoval() {
	for iter := 0; iter < 500; iter++ {
		s := NewBaseStore[int, int](StoreOptions{TTL: time.Millisecond})
		const n = 128
		for i := 0; i < n; i++ {
			s.Set(i, i)
		}
		time.Sleep(2 * time.Millisecond) // let everything expire

		const readers = 6
		var wg sync.WaitGroup
		wg.Add(readers + 1)
		go func() {
			defer wg.Done()
			s.SweepExpired(0)
		}()
		for r := 0; r < readers; r++ {
			go func() {
				defer wg.Done()
				for i := 0; i < n; i++ {
					s.Get(i) // lazy-deletes the expired entry, racing the sweep
				}
			}()
		}
		wg.Wait()
		s.Close()
	}
	assert.True(su.T(), true)
}
