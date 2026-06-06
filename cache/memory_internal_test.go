package cache

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// Internal-package tests for the global-overflow eviction helpers. These call
// unexported methods directly so that proportional/byte-share eviction edge
// cases can be exercised deterministically, independent of the background
// sweeper's timing.

type memoryInternalSuite struct {
	suite.Suite
}

func TestMemoryInternalSuite(t *testing.T) {
	suite.Run(t, new(memoryInternalSuite))
}

func (s *memoryInternalSuite) seedGuilds(c *MemoryCache, n int) {
	for i := 0; i < n; i++ {
		c.guilds.Set(&discord.Guild{ID: discord.Snowflake(uint64(i + 1)), Name: fmt.Sprintf("g%d", i)})
	}
}

// ── countTargeted: every category branch ──────────────────────────────────────

func (s *memoryInternalSuite) TestCountTargeted_AllCategories() {
	c := NewMemoryCache(Options{Messages: MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	c.guilds.Set(&discord.Guild{ID: 1})
	c.channels.Set(&discord.Channel{ID: 2})
	c.users.Set(&discord.User{ID: 3})
	c.members.Set(10, &discord.GuildMember{User: &discord.User{ID: 4}})
	c.roles.Set(10, &discord.Role{ID: 5})
	c.messages.Add(&discord.Message{ID: 6, ChannelID: 99})
	c.voiceStates.Set(10, &discord.VoiceState{UserID: 7})
	c.soundboard.Set(10, &discord.SoundboardSound{SoundID: 8})
	c.scheduledEvents.Set(&discord.GuildScheduledEvent{ID: 9, GuildID: 10})
	c.stageInstances.Set(&discord.StageInstance{ID: 11, GuildID: 10})
	c.emojis.Set(10, &discord.Emoji{ID: 12})
	c.stickers.Set(10, &discord.Sticker{ID: 13})
	c.bans.Set(10, &discord.Ban{User: discord.User{ID: 15}})
	c.autoModRules.Set(10, &discord.AutoModerationRule{ID: 16, GuildID: 10})
	c.invites.Set(&discord.Invite{Code: "abc"})

	// CategoryAll exercises all 16 branches of countTargeted.
	got := c.countTargeted(CategoryAll)
	assert.Equal(s.T(), 15, got)
}

// ── storesForTarget: every category branch incl. messages ─────────────────────

func (s *memoryInternalSuite) TestStoresForTarget_AllCategories() {
	c := NewMemoryCache(Options{Messages: MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	c.guilds.Set(&discord.Guild{ID: 1})
	c.messages.Add(&discord.Message{ID: 2, ChannelID: 99})

	entries := c.storesForTarget(CategoryAll)
	// Only non-empty stores are included (guilds + messages here).
	assert.Len(s.T(), entries, 2)
}

// ── evictByCount edge cases ────────────────────────────────────────────────────

func (s *memoryInternalSuite) TestEvictByCount_NothingToEvict() {
	c := NewMemoryCache(Options{})
	defer c.Close()
	// No entries in any targeted store → total == 0 → early return, no panic.
	c.evictByCount(5, CategoryGuilds, ClearByLastUsed)
	assert.Equal(s.T(), 0, c.guilds.Size())
}

func (s *memoryInternalSuite) TestEvictByCount_LargestRemainderDistributesExtra() {
	c := NewMemoryCache(Options{})
	defer c.Close()

	s.seedGuilds(c, 3)
	c.users.Set(&discord.User{ID: 100})

	// total = 4, need = 1: both floor shares are 0, so the single extra unit goes
	// to the largest remainder (guilds, 0.75 > users 0.25).
	c.evictByCount(1, CategoryGuilds|CategoryUsers, ClearByLastUsed)

	assert.Equal(s.T(), 2, c.guilds.Size(), "guilds should shed exactly one entry")
	assert.Equal(s.T(), 1, c.users.Size(), "users should be spared")
}

func (s *memoryInternalSuite) TestEvictByCount_OverNeedClampsToZero() {
	c := NewMemoryCache(Options{})
	defer c.Close()

	s.seedGuilds(c, 3)
	c.users.Set(&discord.User{ID: 100})

	// need far exceeds total → each store's computed trim target goes negative
	// and must be clamped to 0 (drain everything, never below empty).
	c.evictByCount(1000, CategoryGuilds|CategoryUsers, ClearByLastUsed)

	assert.Equal(s.T(), 0, c.guilds.Size())
	assert.Equal(s.T(), 0, c.users.Size())
}

// ── evictByBytes edge cases ────────────────────────────────────────────────────

func (s *memoryInternalSuite) TestEvictByBytes_SkipsUntrackedAndStopsWhenSatisfied() {
	// MaxSizeMB > 0 → byte tracking is enabled for the byte-tracking stores.
	c := NewMemoryCache(Options{
		Limits:   Limits{MaxSizeMB: 1},
		Messages: MessageOptions{MaxPerChannel: 100},
	})
	defer c.Close()

	s.seedGuilds(c, 20)
	c.users.Set(&discord.User{ID: 9999, Username: "u"})
	// Messages never track bytes → storeEntry.bytes == 0 → skipped (continue).
	c.messages.Add(&discord.Message{ID: 1, ChannelID: 50})

	before := c.guilds.Size() + c.users.Size()

	// Need only a tiny number of bytes: the first (largest-avg) store frees it,
	// then need <= 0 breaks the loop before touching the rest.
	c.evictByBytes(1, CategoryGuilds|CategoryUsers|CategoryMessages, ClearByLastUsed)

	assert.Less(s.T(), c.guilds.Size()+c.users.Size(), before, "some bytes must have been freed")
	assert.Equal(s.T(), 1, c.messages.Size(), "untracked message store must be untouched")
}

func (s *memoryInternalSuite) TestEvictByBytes_NeedExceedsStoreCapacityDrains() {
	c := NewMemoryCache(Options{Limits: Limits{MaxSizeMB: 1}})
	defer c.Close()

	s.seedGuilds(c, 5)

	// need is far larger than the store could ever free → toEvict is capped to
	// the store size and the store is drained.
	c.evictByBytes(10*1024*1024, CategoryGuilds, ClearByLastUsed)
	assert.Equal(s.T(), 0, c.guilds.Size())
}

func (s *memoryInternalSuite) TestEvictByBytes_AvgRoundsDownToZeroIsClamped() {
	c := NewMemoryCache(Options{Limits: Limits{MaxSizeMB: 1}})
	defer c.Close()

	// Sounds with a NaN volume cannot be JSON-encoded, so estimateSize returns 0
	// for them: they add to the entry count but not to the byte total. Mixing
	// many zero-byte entries with one real entry drives the store's average entry
	// size (bytes/size) below 1, exercising the `avg <= 0 → avg = 1` clamp.
	nan := math.NaN()
	for i := 1; i <= 100; i++ {
		c.soundboard.Set(10, &discord.SoundboardSound{SoundID: discord.Snowflake(uint64(i)), Volume: nan})
	}
	c.soundboard.Set(10, &discord.SoundboardSound{SoundID: 200, Name: "real", Volume: 0.5})

	assert.Positive(s.T(), c.soundboard.Bytes())
	assert.Less(s.T(), c.soundboard.Bytes(), int64(c.soundboard.Size()),
		"setup precondition: bytes < size so that avg rounds down to 0")

	c.evictByBytes(1, CategorySoundboard, ClearByLastUsed)
	assert.Less(s.T(), c.soundboard.Size(), 101, "at least one entry must be evicted")
}

func (s *memoryInternalSuite) TestEvictByBytes_NoTrackedBytesIsNoop() {
	// No MaxSizeMB → trackBytes is false → every store reports 0 bytes → all are
	// skipped and nothing is evicted.
	c := NewMemoryCache(Options{})
	defer c.Close()

	s.seedGuilds(c, 5)
	c.evictByBytes(100, CategoryGuilds, ClearByLastUsed)
	assert.Equal(s.T(), 5, c.guilds.Size())
}
