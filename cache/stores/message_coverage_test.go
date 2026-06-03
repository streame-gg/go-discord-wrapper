package stores

import (
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func msg(ch, id discord.Snowflake) *discord.Message {
	return &discord.Message{ID: id, ChannelID: ch, Content: "x"}
}

// ── enable / disable ──────────────────────────────────────────────────────────

func (su *storesMessageSuite) TestMessageStore_Disabled() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 0}) // capPerCh == 0 → disabled
	defer s.Close()

	assert.False(t, s.IsEnabled())

	ch, id := discord.Snowflake(1), discord.Snowflake(2)
	s.Add(msg(ch, id))
	_, ok := s.Get(ch, id)
	assert.False(t, ok)
	s.Update(msg(ch, id))
	s.Delete(ch, id)
	s.DeleteBulk(ch, []discord.Snowflake{id})
	s.DeleteChannel(ch)
	s.TrimToTotal(0)
	s.SweepExpired(0)
	assert.Equal(t, 0, s.Channel(ch).Len())
	assert.Equal(t, 0, s.Size())
}

func (su *storesMessageSuite) TestMessageStore_NegativeMaxItemsDefaultsTo100() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: -1})
	defer s.Close()
	assert.True(t, s.IsEnabled())
	assert.Equal(t, 100, s.capPerCh)
}

// ── basic CRUD ────────────────────────────────────────────────────────────────

func (su *storesMessageSuite) TestMessageStore_AddGetUpdateDelete() {
	t := su.T()
	s := NewMessageStore(MessageDefaults())
	defer s.Close()

	ch := discord.Snowflake(1)
	id := discord.Snowflake(10)

	s.Add(nil) // no-op
	s.Add(msg(ch, id))
	assert.Equal(t, 1, s.Size())

	got, ok := s.Get(ch, id)
	require.True(t, ok)
	assert.Equal(t, "x", got.Content)

	// Get miss: unknown channel and unknown id.
	_, ok = s.Get(discord.Snowflake(999), id)
	assert.False(t, ok)
	_, ok = s.Get(ch, discord.Snowflake(999))
	assert.False(t, ok)

	// Add same id again → update in place, no size change.
	s.Add(&discord.Message{ID: id, ChannelID: ch, Content: "updated"})
	assert.Equal(t, 1, s.Size())
	got, _ = s.Get(ch, id)
	assert.Equal(t, "updated", got.Content)

	// Update existing.
	s.Update(&discord.Message{ID: id, ChannelID: ch, Content: "again"})
	got, _ = s.Get(ch, id)
	assert.Equal(t, "again", got.Content)

	// Update no-ops: nil, unknown channel, unknown id.
	s.Update(nil)
	s.Update(msg(discord.Snowflake(999), id))
	s.Update(msg(ch, discord.Snowflake(888)))

	// Delete unknown id / unknown channel → size unchanged.
	s.Delete(ch, discord.Snowflake(888))
	s.Delete(discord.Snowflake(999), id)
	assert.Equal(t, 1, s.Size())

	s.Delete(ch, id)
	assert.Equal(t, 0, s.Size())
	_, ok = s.Get(ch, id)
	assert.False(t, ok)
}

func (su *storesMessageSuite) TestMessageStore_DeleteBulkAndChannel() {
	t := su.T()
	s := NewMessageStore(MessageDefaults())
	defer s.Close()

	ch := discord.Snowflake(1)
	for i := 1; i <= 5; i++ {
		s.Add(msg(ch, discord.Snowflake(i)))
	}
	assert.Equal(t, 5, s.Size())

	s.DeleteBulk(ch, []discord.Snowflake{1, 2, 999})
	assert.Equal(t, 3, s.Size())

	// DeleteBulk on unknown channel → no-op.
	s.DeleteBulk(discord.Snowflake(42), []discord.Snowflake{3})
	assert.Equal(t, 3, s.Size())

	coll := s.Channel(ch)
	assert.Equal(t, 3, coll.Len())

	// Channel on unknown channel returns empty.
	assert.Equal(t, 0, s.Channel(discord.Snowflake(42)).Len())

	s.DeleteChannel(ch)
	assert.Equal(t, 0, s.Size())
	// DeleteChannel on unknown channel → no-op.
	s.DeleteChannel(discord.Snowflake(42))

	assert.EqualValues(t, 0, s.Bytes())
}

// ── per-channel eviction (ring full) ──────────────────────────────────────────

func (su *storesMessageSuite) TestMessageStore_PerChannelEvictionLRU() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 2, EvictionStrategy: EvictLRU})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1))
	time.Sleep(time.Millisecond)
	s.Add(msg(ch, 2))
	// access 1 so 2 becomes least-recently-used
	time.Sleep(time.Millisecond)
	s.Get(ch, 1)
	time.Sleep(time.Millisecond)
	s.Add(msg(ch, 3)) // ring full → evict LRU (id 2)

	assert.Equal(t, 2, s.Size())
	_, ok := s.Get(ch, 2)
	assert.False(t, ok, "id 2 should be evicted (LRU)")
	_, ok = s.Get(ch, 1)
	assert.True(t, ok)
}

func (su *storesMessageSuite) TestMessageStore_PerChannelEvictionLFU() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 2, EvictionStrategy: EvictLFU})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1))
	s.Add(msg(ch, 2))
	s.Get(ch, 1)
	s.Get(ch, 1) // id 1 hitCount high
	s.Add(msg(ch, 3))

	_, ok := s.Get(ch, 2)
	assert.False(t, ok, "id 2 (hitCount 0) should be evicted (LFU)")
}

func (su *storesMessageSuite) TestMessageStore_PerChannelEvictionFIFO() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 2, EvictionStrategy: EvictFIFO})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1))
	time.Sleep(time.Millisecond)
	s.Add(msg(ch, 2))
	s.Get(ch, 1) // access must not matter for FIFO
	time.Sleep(time.Millisecond)
	s.Add(msg(ch, 3))

	_, ok := s.Get(ch, 1)
	assert.False(t, ok, "id 1 (earliest inserted) should be evicted (FIFO)")
}

// ── global cap (MaxTotal) ─────────────────────────────────────────────────────

func (su *storesMessageSuite) TestMessageStore_MaxTotalEvictsAcrossChannels() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 100, MaxTotal: 3, EvictionStrategy: EvictLRU})
	defer s.Close()

	c1, c2 := discord.Snowflake(1), discord.Snowflake(2)
	s.Add(msg(c1, 1))
	s.Add(msg(c1, 2))
	time.Sleep(time.Millisecond)
	s.Add(msg(c2, 3))
	s.Add(msg(c2, 4)) // total 4 > 3 → trim oldest-accessed channel first (c1)

	assert.Equal(t, 3, s.Size())
}

func (su *storesMessageSuite) TestMessageStore_TrimToTotal() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 100})
	defer s.Close()

	c1, c2 := discord.Snowflake(1), discord.Snowflake(2)
	s.Add(msg(c1, 1))
	s.Add(msg(c1, 2))
	time.Sleep(time.Millisecond)
	s.Add(msg(c2, 3))

	s.TrimToTotal(-1) // no-op
	assert.Equal(t, 3, s.Size())

	s.TrimToTotal(1)
	assert.Equal(t, 1, s.Size())

	// Trim everything → emptied channels are removed from the map.
	s.TrimToTotal(0)
	assert.Equal(t, 0, s.Size())
}

// ── TTL & sweeping ────────────────────────────────────────────────────────────

func (su *storesMessageSuite) TestMessageStore_TTLExpiryOnGetAndChannel() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 10, TTL: 20 * time.Millisecond})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1))
	time.Sleep(40 * time.Millisecond)

	_, ok := s.Get(ch, 1)
	assert.False(t, ok, "message should be TTL-expired on Get")

	// Channel() must exclude expired messages.
	assert.Equal(t, 0, s.Channel(ch).Len())
}

func (su *storesMessageSuite) TestMessageStore_SweepExpired() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 10, TTL: 20 * time.Millisecond})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1))
	s.Add(msg(ch, 2))
	time.Sleep(40 * time.Millisecond)

	s.SweepExpired(0)
	assert.Equal(t, 0, s.Size())
	// Channel emptied → removed from map; Channel() returns empty.
	assert.Equal(t, 0, s.Channel(ch).Len())
}

func (su *storesMessageSuite) TestMessageStore_SweepIdleWindow() {
	t := su.T()
	s := NewMessageStore(StoreOptions{MaxItems: 10})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(msg(ch, 1))
	time.Sleep(10 * time.Millisecond)

	// unusedWindow shorter than idle time → message swept as idle.
	s.SweepExpired(time.Millisecond)
	assert.Equal(t, 0, s.Size())
}

func (su *storesMessageSuite) TestMessageStore_BackgroundSweeper() {
	t := su.T()
	s := NewMessageStore(StoreOptions{
		MaxItems: 10,
		Sweepers: []Sweeper{{
			Interval: 5 * time.Millisecond,
			Filter:   func(v any) bool { return v.(*discord.Message).Content == "drop" },
		}},
	})
	defer s.Close()

	ch := discord.Snowflake(1)
	s.Add(&discord.Message{ID: 1, ChannelID: ch, Content: "drop"})
	s.Add(&discord.Message{ID: 2, ChannelID: ch, Content: "keep"})

	require.Eventually(t, func() bool { return s.Size() == 1 }, time.Second, 5*time.Millisecond)
	_, ok := s.Get(ch, 2)
	assert.True(t, ok, "non-matching message should survive the sweeper")
}

func (su *storesMessageSuite) TestMessageStore_SweeperEmptiesChannel() {
	t := su.T()
	s := NewMessageStore(StoreOptions{
		MaxItems: 10,
		Sweepers: []Sweeper{{
			Interval: 5 * time.Millisecond,
			Filter:   func(v any) bool { return true }, // drop everything
		}},
	})
	defer s.Close()

	s.Add(msg(discord.Snowflake(1), 1))
	require.Eventually(t, func() bool { return s.Size() == 0 }, time.Second, 5*time.Millisecond)
}
