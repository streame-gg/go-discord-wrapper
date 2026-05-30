package stores

// Tests for Bug #93: when BaseStore evicts or sweeps an entry from a
// composite-key store (VoiceState, Member, Presence), the compositeGuildIndex
// must be updated as well. Without the fix, AllInGuild / GetByGuild returned
// phantom entries whose backing data had already been removed.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func newTestVoiceState(userID discord.Snowflake) *discord.VoiceState {
	return &discord.VoiceState{UserID: userID}
}

func newTestMember(userID discord.Snowflake) *discord.GuildMember {
	return &discord.GuildMember{User: &discord.User{ID: userID}}
}

func newTestPresence(guildID, userID discord.Snowflake) *discord.Presence {
	return &discord.Presence{GuildID: guildID, User: discord.PartialPresenceUser{ID: userID}}
}

// assertIndexInSync checks that AllInGuild/GetByGuild returns exactly wantLen
// entries and that none of them equal any of the bannedUserIDs.
func assertVoiceStateIndexInSync(t *testing.T, s *MemVoiceStateStore, guildID discord.Snowflake, wantLen int, evictedIDs ...discord.Snowflake) {
	t.Helper()
	all := s.AllInGuild(guildID)
	assert.Equal(t, wantLen, all.Len(), "AllInGuild length mismatch — stale index entry present")
	for _, uid := range evictedIDs {
		_, ok := all.Get(uid)
		assert.False(t, ok, "evicted userID %d must not appear in AllInGuild", uid)
	}
}

func assertMemberIndexInSync(t *testing.T, s *MemMemberStore, guildID discord.Snowflake, wantLen int, evictedIDs ...discord.Snowflake) {
	t.Helper()
	all := s.AllInGuild(guildID)
	assert.Equal(t, wantLen, all.Len(), "AllInGuild length mismatch — stale index entry present")
	for _, uid := range evictedIDs {
		_, ok := all.Get(uid)
		assert.False(t, ok, "evicted userID %d must not appear in AllInGuild", uid)
	}
}

func assertPresenceIndexInSync(t *testing.T, s *MemPresenceStore, guildID discord.Snowflake, wantLen int, evictedIDs ...discord.Snowflake) {
	t.Helper()
	all := s.GetByGuild(guildID)
	assert.Equal(t, wantLen, all.Len(), "GetByGuild length mismatch — stale index entry present")
	for _, uid := range evictedIDs {
		_, ok := all.Get(uid)
		assert.False(t, ok, "evicted userID %d must not appear in GetByGuild", uid)
	}
}

// ── VoiceStateStore ───────────────────────────────────────────────────────────

// TestBug93_VoiceStateStore_MaxItemsEvictionKeepsIndexInSync verifies that
// capacity-based eviction (evictOne) removes the entry from the guild index.
func TestBug93_VoiceStateStore_MaxItemsEvictionKeepsIndexInSync(t *testing.T) {
	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(1)
	user2 := discord.Snowflake(2)
	user3 := discord.Snowflake(3)

	s := NewVoiceStateStore(StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictLRU,
	})
	defer s.Close()

	s.Set(guildID, newTestVoiceState(user1))
	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestVoiceState(user2))

	// Touch user2 so user1 becomes the LRU victim.
	time.Sleep(time.Millisecond)
	s.Get(guildID, user2)

	// Adding user3 pushes Len to 3, triggering eviction of user1.
	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestVoiceState(user3))

	require.Equal(t, 2, s.Size())
	assertVoiceStateIndexInSync(t, s, guildID, 2, user1)
}

// TestBug93_VoiceStateStore_TrimToKeepsIndexInSync verifies that TrimTo
// (evictToCount path) keeps the guild index in sync.
func TestBug93_VoiceStateStore_TrimToKeepsIndexInSync(t *testing.T) {
	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(1)
	user2 := discord.Snowflake(2)
	user3 := discord.Snowflake(3)

	s := NewVoiceStateStore(StoreOptions{EvictionStrategy: EvictFIFO})
	defer s.Close()

	s.Set(guildID, newTestVoiceState(user1))
	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestVoiceState(user2))
	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestVoiceState(user3))

	// TrimTo(1) should remove user1 and user2 (FIFO = earliest inserted first).
	s.TrimTo(1)

	require.Equal(t, 1, s.Size())
	assertVoiceStateIndexInSync(t, s, guildID, 1, user1, user2)
}

// TestBug93_VoiceStateStore_SweepExpiredKeepsIndexInSync verifies that TTL
// expiry (SweepExpired path) removes entries from the guild index.
func TestBug93_VoiceStateStore_SweepExpiredKeepsIndexInSync(t *testing.T) {
	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(1)
	user2 := discord.Snowflake(2)

	s := NewVoiceStateStore(StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.Set(guildID, newTestVoiceState(user1))
	s.Set(guildID, newTestVoiceState(user2))

	time.Sleep(40 * time.Millisecond)
	s.SweepExpired(0)

	require.Equal(t, 0, s.Size())
	assertVoiceStateIndexInSync(t, s, guildID, 0, user1, user2)
}

// TestBug93_VoiceStateStore_SweeperKeepsIndexInSync verifies that a background
// Sweeper (runSweeper path) removes entries from the guild index.
func TestBug93_VoiceStateStore_SweeperKeepsIndexInSync(t *testing.T) {
	guildID := discord.Snowflake(1000)
	keep := discord.Snowflake(1)
	remove := discord.Snowflake(2)

	s := NewVoiceStateStore(StoreOptions{
		Sweepers: []Sweeper{{
			Interval: 5 * time.Millisecond,
			Filter: func(v any) bool {
				vs, ok := v.(*discord.VoiceState)
				return ok && vs.UserID == remove
			},
		}},
	})
	defer s.Close()

	s.Set(guildID, newTestVoiceState(keep))
	s.Set(guildID, newTestVoiceState(remove))

	// Wait for the sweeper to fire at least once.
	require.Eventually(t, func() bool {
		return s.Size() == 1
	}, 500*time.Millisecond, 10*time.Millisecond, "sweeper did not remove entry within 500ms")

	assertVoiceStateIndexInSync(t, s, guildID, 1, remove)
}

// ── MemberStore ───────────────────────────────────────────────────────────────

// TestBug93_MemberStore_MaxItemsEvictionKeepsIndexInSync is the same scenario
// as the VoiceStateStore test but for MemMemberStore.
func TestBug93_MemberStore_MaxItemsEvictionKeepsIndexInSync(t *testing.T) {
	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(1)
	user2 := discord.Snowflake(2)
	user3 := discord.Snowflake(3)

	s := NewMemberStore(StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictLRU,
	})
	defer s.Close()

	s.Set(guildID, newTestMember(user1))
	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestMember(user2))

	time.Sleep(time.Millisecond)
	s.Get(guildID, user2)

	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestMember(user3))

	require.Equal(t, 2, s.Size())
	assertMemberIndexInSync(t, s, guildID, 2, user1)
}

// ── PresenceStore ─────────────────────────────────────────────────────────────

// TestBug93_PresenceStore_MaxItemsEvictionKeepsIndexInSync is the same scenario
// as the VoiceStateStore test but for MemPresenceStore.
func TestBug93_PresenceStore_MaxItemsEvictionKeepsIndexInSync(t *testing.T) {
	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(1)
	user2 := discord.Snowflake(2)
	user3 := discord.Snowflake(3)

	s := NewPresenceStore(StoreOptions{
		MaxItems:         2,
		EvictionStrategy: EvictLRU,
		Disabled:         false,
	})
	defer s.Close()

	s.Set(newTestPresence(guildID, user1))
	time.Sleep(time.Millisecond)
	s.Set(newTestPresence(guildID, user2))

	time.Sleep(time.Millisecond)
	s.Get(guildID, user2)

	time.Sleep(time.Millisecond)
	s.Set(newTestPresence(guildID, user3))

	require.Equal(t, 2, s.Size())
	assertPresenceIndexInSync(t, s, guildID, 2, user1)
}
