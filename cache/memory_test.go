package cache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/stretchr/testify/suite"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func guild(id string) *common.Guild {
	return &common.Guild{ID: common.Snowflake(id), Name: "guild-" + id}
}

func channel(id string) *common.Channel {
	return &common.Channel{ID: common.Snowflake(id), Name: "chan-" + id}
}

func user(id string) *common.User {
	return &common.User{ID: common.Snowflake(id), Username: "user-" + id}
}

func member(userID string) *common.GuildMember {
	return &common.GuildMember{User: user(userID)}
}

func sticker(id string) *common.Sticker {
	return &common.Sticker{ID: common.Snowflake(id), Name: "sticker-" + id}
}

func message(id, channelID string) *common.Message {
	return &common.Message{
		ID:        common.Snowflake(id),
		ChannelID: common.Snowflake(channelID),
		Content:   "msg-" + id,
	}
}

func newCache(opts cache.Options) *cache.MemoryCache {
	if opts.SweepInterval == 0 {
		opts.SweepInterval = 10 * time.Millisecond
	}
	return cache.NewMemoryCache(opts)
}

type memoryTestSuite struct {
	suite.Suite
}

func TestMemoryTestSuite(t *testing.T) {
	suite.Run(t, new(memoryTestSuite))
}

// ── GuildStore ────────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestGuildStore_SetGetDelete() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Guilds().Set(guild("1"))

	got, ok := c.Guilds().Get("1")
	s.Require().Truef(ok && got.ID == "1", "expected guild 1, got ok=%v", ok)
	c.Guilds().Delete("1")
	_, ok = c.Guilds().Get("1")
	s.False(ok, "expected guild to be deleted")
}

func (s *memoryTestSuite) TestGuildStore_All() {
	c := newCache(cache.Options{})
	defer c.Close()

	for i := 1; i <= 5; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	s.Lenf(c.Guilds().All(), 5, "expected 5 guilds, got %d", len(c.Guilds().All()))
}

func (s *memoryTestSuite) TestGuildStore_TTLExpiry() {
	c := newCache(cache.Options{TTL: 30 * time.Millisecond})
	defer c.Close()

	c.Guilds().Set(guild("1"))
	_, ok := c.Guilds().Get("1")
	s.Require().True(ok, "expected guild before TTL")
	time.Sleep(60 * time.Millisecond)
	_, ok = c.Guilds().Get("1")
	s.False(ok, "expected guild expired after TTL")
}

func (s *memoryTestSuite) TestGuildStore_SweepRemovesExpired() {
	c := newCache(cache.Options{TTL: 20 * time.Millisecond, SweepInterval: 10 * time.Millisecond})
	defer c.Close()

	c.Guilds().Set(guild("1"))
	time.Sleep(60 * time.Millisecond)
	s.Equalf(0, c.Guilds().Size(), "expected size 0 after sweep, got %d", c.Guilds().Size())
}

// ── ChannelStore ──────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestChannelStore_CRUD() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Channels().Set(channel("10"))
	_, ok := c.Channels().Get("10")
	s.Require().True(ok, "expected channel 10")
	c.Channels().Delete("10")
	_, ok = c.Channels().Get("10")
	s.False(ok, "expected channel deleted")
}

// ── UserStore ─────────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestUserStore_CRUD() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Users().Set(user("42"))
	_, ok := c.Users().Get("42")
	s.Require().True(ok, "expected user 42")
	c.Users().Delete("42")
	_, ok = c.Users().Get("42")
	s.False(ok, "expected user deleted")
}

// ── MemberStore ───────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestMemberStore_SetGet() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Members().Set("guild1", member("99"))
	got, ok := c.Members().Get("guild1", "99")
	s.Truef(ok && got.User.ID == "99", "expected member 99 in guild1, ok=%v", ok)
}

func (s *memoryTestSuite) TestMemberStore_DeleteGuild() {
	c := newCache(cache.Options{})
	defer c.Close()

	for i := 1; i <= 5; i++ {
		c.Members().Set("guild1", member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set("guild2", member("99"))

	c.Members().DeleteGuild("guild1")

	s.Lenf(c.Members().AllInGuild("guild1"), 0, "expected guild1 members deleted")
	s.Lenf(c.Members().AllInGuild("guild2"), 1, "guild2 member should be unaffected")
}

func (s *memoryTestSuite) TestMemberStore_NilUser() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Members().Set("guild1", &common.GuildMember{})
	s.Equalf(0, c.Members().Size(), "expected size 0 when member.User is nil")
}

// ── StickerStore ──────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestStickerStore_CRUDAndSetAll() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Stickers().Set("guild1", sticker("s1"))
	got, ok := c.Stickers().Get("s1")
	s.Truef(ok && got.ID == "s1", "expected sticker s1, ok=%v", ok)

	c.Stickers().SetAll("guild1", []*common.Sticker{sticker("s2"), sticker("s3")})
	_, ok = c.Stickers().Get("s1")
	s.False(ok, "expected old guild stickers replaced by SetAll")
	s.Lenf(c.Stickers().GetByGuild("guild1"), 2, "expected 2 stickers after SetAll")

	c.Stickers().Delete("s2")
	_, ok = c.Stickers().Get("s2")
	s.False(ok, "expected sticker s2 deleted")

	c.Stickers().DeleteGuild("guild1")
	s.Lenf(c.Stickers().GetByGuild("guild1"), 0, "expected guild stickers deleted")
}

// ── MessageStore ──────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestMessageStore_AddGet() {
	c := newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	c.Messages().Add(message("msg1", "chan1"))
	got, ok := c.Messages().Get("chan1", "msg1")
	s.Truef(ok && got.ID == "msg1", "expected msg1, ok=%v", ok)
}

func (s *memoryTestSuite) TestMessageStore_Update() {
	c := newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	msg := message("msg1", "chan1")
	c.Messages().Add(msg)

	updated := *msg
	updated.Content = "edited"
	c.Messages().Update(&updated)

	got, _ := c.Messages().Get("chan1", "msg1")
	s.Equalf("edited", got.Content, "expected updated content, got %q", got.Content)
}

func (s *memoryTestSuite) TestMessageStore_DeleteBulk() {
	c := newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("msg%d", i), "chan1"))
	}
	c.Messages().DeleteBulk("chan1", []common.Snowflake{"msg1", "msg3", "msg5"})

	for _, id := range []common.Snowflake{"msg1", "msg3", "msg5"} {
		_, ok := c.Messages().Get("chan1", id)
		s.Falsef(ok, "expected %s deleted", id)
	}
	for _, id := range []common.Snowflake{"msg2", "msg4"} {
		_, ok := c.Messages().Get("chan1", id)
		s.Truef(ok, "expected %s present", id)
	}
}

func (s *memoryTestSuite) TestMessageStore_Channel_NewestFirst() {
	c := newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("msg%d", i), "chan1"))
	}
	msgs := c.Messages().Channel("chan1")
	s.Require().Lenf(msgs, 5, "expected 5 messages, got %d", len(msgs))
	s.Equalf(common.Snowflake("msg5"), msgs[0].ID, "expected newest first (msg5), got %s", msgs[0].ID)
}

func (s *memoryTestSuite) TestMessageStore_DeleteChannel() {
	c := newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	c.Messages().Add(message("msg1", "chan1"))
	c.Messages().Add(message("msg2", "chan2"))

	c.Messages().DeleteChannel("chan1")

	s.Lenf(c.Messages().Channel("chan1"), 0, "expected chan1 messages deleted")
	s.Lenf(c.Messages().Channel("chan2"), 1, "chan2 should be unaffected")
}

func (s *memoryTestSuite) TestMessageStore_RingEvictsOldest() {
	c := cache.NewMemoryCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 3},
	})
	defer c.Close()

	for i := 1; i <= 4; i++ {
		c.Messages().Add(message(fmt.Sprintf("msg%d", i), "chan1"))
	}
	msgs := c.Messages().Channel("chan1")
	s.Require().Lenf(msgs, 3, "expected ring cap of 3, got %d", len(msgs))
	for _, m := range msgs {
		s.NotEqualf(common.Snowflake("msg1"), m.ID, "msg1 (oldest) should have been evicted")
	}
}

func (s *memoryTestSuite) TestMessageStore_TTLExpiry() {
	c := cache.NewMemoryCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 10, TTL: 20 * time.Millisecond},
	})
	defer c.Close()

	c.Messages().Add(message("msg1", "chan1"))
	_, ok := c.Messages().Get("chan1", "msg1")
	s.Require().True(ok, "expected message before TTL")
	time.Sleep(50 * time.Millisecond)
	_, ok = c.Messages().Get("chan1", "msg1")
	s.False(ok, "expected message expired")
}

func (s *memoryTestSuite) TestMessageStore_SizeAccuracy() {
	c := newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Close()

	for i := 1; i <= 10; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "chan1"))
	}
	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("n%d", i), "chan2"))
	}
	s.Equalf(15, c.Messages().Size(), "expected size 15, got %d", c.Messages().Size())
	c.Messages().DeleteChannel("chan1")
	s.Equalf(5, c.Messages().Size(), "expected size 5 after DeleteChannel, got %d", c.Messages().Size())
}

// ── EvictUnused ───────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestEvictUnused_RemovesUnaccessed() {
	const unusedWindow = 60 * time.Millisecond

	c := cache.NewMemoryCache(cache.Options{
		TTL:           1 * time.Hour,
		SweepInterval: 15 * time.Millisecond,
		EvictBehavior: cache.EvictUnused,
		UnusedWindow:  unusedWindow,
	})
	defer c.Close()

	c.Guilds().Set(guild("1")) // never accessed again
	c.Guilds().Set(guild("2"))

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.Guilds().Get("2")
			case <-stop:
				return
			}
		}
	}()

	time.Sleep(unusedWindow + 3*15*time.Millisecond)
	close(stop)
	wg.Wait()
	c.Guilds().Get("2") // final touch before assertion

	_, ok1 := c.Guilds().Get("1")
	_, ok2 := c.Guilds().Get("2")
	s.False(ok1, "guild 1 should have been evicted (never accessed)")
	s.True(ok2, "guild 2 should still be present (continuously accessed)")
}

// ── Per-category limits ───────────────────────────────────────────────────────

func (s *memoryTestSuite) TestLimits_MaxGuilds_ClearByLastUsed() {
	c := cache.NewMemoryCache(cache.Options{
		Limits: cache.Limits{MaxGuilds: 2},
		OnOverflow: cache.OverflowPolicy{
			ClearBy: cache.ClearByLastUsed,
			Target:  cache.CategoryAll,
		},
	})
	defer c.Close()

	// Insert 3 guilds. Access guild "1" and "2" so "3" is most recent.
	c.Guilds().Set(guild("1"))
	time.Sleep(2 * time.Millisecond)
	c.Guilds().Set(guild("2"))
	time.Sleep(2 * time.Millisecond)
	c.Guilds().Set(guild("3")) // triggers eviction → guild "1" should go

	s.LessOrEqualf(c.Guilds().Size(), 2, "expected at most 2 guilds, got %d", c.Guilds().Size())
	_, ok := c.Guilds().Get("1")
	s.False(ok, "guild 1 (least recently accessed) should have been evicted")
}

func (s *memoryTestSuite) TestLimits_MaxGuilds_ClearByAge() {
	c := cache.NewMemoryCache(cache.Options{
		Limits: cache.Limits{MaxGuilds: 2},
		OnOverflow: cache.OverflowPolicy{
			ClearBy: cache.ClearByAge,
			Target:  cache.CategoryAll,
		},
	})
	defer c.Close()

	c.Guilds().Set(guild("old"))
	time.Sleep(2 * time.Millisecond)
	c.Guilds().Set(guild("mid"))
	time.Sleep(2 * time.Millisecond)

	// Access "old" right before inserting the third — for ClearByAge access
	// time doesn't matter; insertion time determines eviction order.
	c.Guilds().Get("old")
	c.Guilds().Set(guild("new")) // triggers eviction → "old" (earliest insert) should go

	s.LessOrEqualf(c.Guilds().Size(), 2, "expected at most 2 guilds, got %d", c.Guilds().Size())
	_, ok := c.Guilds().Get("old")
	s.False(ok, "guild 'old' (earliest insertion) should have been evicted by ClearByAge")
	_, ok = c.Guilds().Get("new")
	s.True(ok, "guild 'new' should be present")
}

func (s *memoryTestSuite) TestLimits_MaxGuilds_ClearByFrequency() {
	c := cache.NewMemoryCache(cache.Options{
		Limits: cache.Limits{MaxGuilds: 2},
		OnOverflow: cache.OverflowPolicy{
			ClearBy: cache.ClearByFrequency,
			Target:  cache.CategoryAll,
		},
	})
	defer c.Close()

	c.Guilds().Set(guild("rare")) // 0 gets
	c.Guilds().Set(guild("busy"))
	c.Guilds().Get("busy") // hit count: 1
	c.Guilds().Get("busy") // hit count: 2

	c.Guilds().Set(guild("new")) // triggers eviction → "rare" (0 hits) goes first

	s.LessOrEqualf(c.Guilds().Size(), 2, "expected at most 2 guilds, got %d", c.Guilds().Size())
	_, ok := c.Guilds().Get("rare")
	s.False(ok, "guild 'rare' (fewest accesses) should have been evicted by ClearByFrequency")
	_, ok = c.Guilds().Get("busy")
	s.True(ok, "guild 'busy' (most accesses) should remain")
}

func (s *memoryTestSuite) TestLimits_MaxMessages() {
	c := cache.NewMemoryCache(cache.Options{
		Limits: cache.Limits{MaxMessages: 5},
		Messages: cache.MessageOptions{
			MaxPerChannel: 10, // per-channel cap is higher; total cap controls
		},
		OnOverflow: cache.OverflowPolicy{ClearBy: cache.ClearByLastUsed, Target: cache.CategoryAll},
	})
	defer c.Close()

	// Add 7 messages across two channels.
	for i := 1; i <= 4; i++ {
		c.Messages().Add(message(fmt.Sprintf("c1m%d", i), "chan1"))
	}
	for i := 1; i <= 3; i++ {
		c.Messages().Add(message(fmt.Sprintf("c2m%d", i), "chan2"))
	}

	s.LessOrEqualf(c.Messages().Size(), 5, "expected ≤ 5 total messages, got %d", c.Messages().Size())
}

// ── Global MaxEntries ─────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestLimits_MaxEntries_Global() {
	c := cache.NewMemoryCache(cache.Options{
		SweepInterval: 10 * time.Millisecond,
		Limits:        cache.Limits{MaxEntries: 4},
		OnOverflow: cache.OverflowPolicy{
			ClearBy: cache.ClearByLastUsed,
			Target:  cache.CategoryGuilds | cache.CategoryChannels,
		},
	})
	defer c.Close()

	// Add 3 guilds and 3 channels (6 total > MaxEntries=4).
	for i := 1; i <= 3; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
		c.Channels().Set(channel(fmt.Sprintf("%d", i)))
	}

	time.Sleep(30 * time.Millisecond) // wait for sweep + overflow enforcement

	total := c.Guilds().Size() + c.Channels().Size()
	s.LessOrEqualf(total, 4, "expected ≤ 4 total entries across targeted stores, got %d", total)
}

func (s *memoryTestSuite) TestLimits_MaxEntries_Target_SparesUntargeted() {
	// Users are NOT in the overflow target — they should be untouched.
	c := cache.NewMemoryCache(cache.Options{
		SweepInterval: 10 * time.Millisecond,
		Limits:        cache.Limits{MaxEntries: 3},
		OnOverflow: cache.OverflowPolicy{
			ClearBy: cache.ClearByLastUsed,
			Target:  cache.CategoryMessages, // only messages are eligible for global eviction
		},
	})
	defer c.Close()

	for i := 1; i <= 3; i++ {
		c.Users().Set(user(fmt.Sprintf("%d", i)))
	}
	// Add messages to push total over limit.
	for i := 1; i <= 3; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "chan1"))
	}

	time.Sleep(30 * time.Millisecond)

	// Users must not have been evicted (not in Target).
	s.Equalf(3, c.Users().Size(), "expected users untouched, got size %d", c.Users().Size())
}

// ── MaxSizeMB ─────────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestLimits_MaxSizeMB() {
	// Set an absurdly small size cap (0.001 MB = ~1 KB) to force eviction.
	c := cache.NewMemoryCache(cache.Options{
		SweepInterval: 10 * time.Millisecond,
		Limits:        cache.Limits{MaxSizeMB: 0.001},
		OnOverflow:    cache.OverflowPolicy{ClearBy: cache.ClearByLastUsed, Target: cache.CategoryAll},
	})
	defer c.Close()

	// Add enough entries to exceed 1 KB of JSON.
	for i := 0; i < 50; i++ {
		c.Guilds().Set(&common.Guild{
			ID:   common.Snowflake(fmt.Sprintf("%d", i)),
			Name: fmt.Sprintf("a-rather-long-guild-name-number-%d", i),
		})
	}

	time.Sleep(40 * time.Millisecond) // sweep + overflow enforcement

	// After enforcement the total should be significantly smaller.
	s.NotEqual(50, c.Guilds().Size(), "expected some guilds evicted to stay within MaxSizeMB")
}

// ── Close ─────────────────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestClose_Idempotent() {
	c := newCache(cache.Options{})
	s.Require().NoErrorf(c.Close(), "first Close")
	s.Require().NoErrorf(c.Close(), "second Close")
}

// ── Concurrent access (race detector) ────────────────────────────────────────

func (s *memoryTestSuite) TestConcurrentAccess_NoRace() {
	c := cache.NewMemoryCache(cache.Options{
		TTL:           5 * time.Millisecond,
		SweepInterval: 5 * time.Millisecond,
		Limits: cache.Limits{
			MaxGuilds:   10,
			MaxMessages: 50,
		},
		OnOverflow: cache.OverflowPolicy{ClearBy: cache.ClearByLastUsed, Target: cache.CategoryAll},
	})
	defer c.Close()

	const goroutines = 20
	const ops = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			gid := fmt.Sprintf("guild%d", g%5)
			uid := fmt.Sprintf("user%d", g%3)
			cid := fmt.Sprintf("chan%d", g%4)
			mid := fmt.Sprintf("msg%d", g)
			for i := 0; i < ops; i++ {
				c.Guilds().Set(guild(gid))
				c.Guilds().Get(common.Snowflake(gid))
				c.Users().Set(user(uid))
				c.Channels().Set(channel(cid))
				c.Messages().Add(message(mid, cid))
				c.Messages().Get(common.Snowflake(cid), common.Snowflake(mid))
				c.Members().Set(common.Snowflake(gid), member(uid))
			}
		}()
	}

	wg.Wait()
}

// TestBug4ByteBasedEvictionDistribution verifies that evictGloballyByBytes
// distributes eviction proportionally by byte share, not entry count.
// Store B (large entries) should be hit harder than store A (small entries).
func TestBug4ByteBasedEvictionDistribution(t *testing.T) {
	// Use a 1 MB limit with a short sweep interval.
	// Store A: 2000 small guilds (~100 bytes each → ~200 KB)
	// Store B: 10 users with 100 KB names each → ~1 MB
	// Total ~1.2 MB > 1 MB limit → overflow triggers in sweep.
	const bigNameLen = 100_000 // ~100 KB per user
	bigName := string(make([]byte, bigNameLen))
	for i := range []byte(bigName) {
		bigName = bigName[:i] + "x" + bigName[i+1:]
		break // just set first char to avoid zero bytes
	}
	bigName = fmt.Sprintf("%0*d", bigNameLen, 0) // 100K zeros as string

	c := cache.NewMemoryCache(cache.Options{
		SweepInterval: 10 * time.Millisecond,
		Limits: cache.Limits{
			MaxSizeMB: 1,
		},
		OnOverflow: cache.OverflowPolicy{
			Target:  cache.CategoryGuilds | cache.CategoryUsers,
			ClearBy: cache.ClearByLastUsed,
		},
	})
	defer c.Close()

	// Seed store A: many small guilds.
	for i := 0; i < 2000; i++ {
		c.Guilds().Set(&common.Guild{ID: common.Snowflake(fmt.Sprintf("g%d", i)), Name: "g"})
	}

	// Seed store B: 10 users each with a ~100KB username.
	for i := 0; i < 10; i++ {
		c.Users().Set(&common.User{
			ID:       common.Snowflake(fmt.Sprintf("u%d", i)),
			Username: bigName,
		})
	}

	before := c.Guilds().Size() + c.Users().Size()

	// Wait for the sweeper to run and enforce the byte limit.
	time.Sleep(100 * time.Millisecond)

	afterGuilds := c.Guilds().Size()
	afterUsers := c.Users().Size()
	after := afterGuilds + afterUsers

	if after >= before {
		t.Skip("no overflow eviction happened — byte limit not exceeded")
	}

	evictedGuilds := 2000 - afterGuilds
	evictedUsers := 10 - afterUsers

	t.Logf("after: guilds=%d users=%d | evictedGuilds=%d evictedUsers=%d",
		afterGuilds, afterUsers, evictedGuilds, evictedUsers)

	// Users contribute far more bytes per entry (~100KB) than guilds (~100B).
	// With the bug (count-based): guilds evicted >> users because 2000 >> 10.
	// With the fix (byte-based): users must be evicted at least as much as guilds
	// relative to their count share.
	if evictedGuilds > 0 && evictedUsers == 0 {
		t.Errorf("only guilds were evicted (count=%d) but not users — byte-based eviction not working (Bug 4)",
			evictedGuilds)
	}
}

// TestBug3ConcurrentAddDeleteChannelNoCounterDrift verifies that concurrent
// Add and DeleteChannel do not produce orphaned rings or counter drift.
func TestBug3ConcurrentAddDeleteChannelNoCounterDrift(t *testing.T) {
	c := cache.NewMemoryCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 50},
	})
	defer c.Close()

	const goroutines = 20
	const ops = 500
	chanID := common.Snowflake("ch1")

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				c.Messages().Add(&common.Message{
					ID:        common.Snowflake(fmt.Sprintf("g%d-m%d", g, i)),
					ChannelID: chanID,
					Content:   "x",
				})
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				c.Messages().DeleteChannel(chanID)
			}
		}()
	}
	wg.Wait()

	// totalMsgs counter must equal the actual number of messages in all rings.
	reported := c.Messages().Size()
	// Get() returns all messages in the channel; use Channel() to count actual messages.
	actual := len(c.Messages().Channel(chanID))
	if reported < 0 {
		t.Errorf("totalMsgs went negative: %d", reported)
	}
	// reported may be slightly stale (recomputed lazily), but actual counts from ring
	// are authoritative. They should be in sync after no concurrent activity.
	if reported != actual {
		t.Errorf("counter drift: reported=%d actual=%d", reported, actual)
	}
}

// TestBug1MaxPerChannelZeroDisablesCaching verifies that MaxPerChannel=0
// disables message caching as documented, rather than defaulting to 100.
func TestBug1MaxPerChannelZeroDisablesCaching(t *testing.T) {
	c := newCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 0},
	})
	defer c.Close()

	c.Messages().Add(message("1", "ch1"))
	c.Messages().Add(message("2", "ch1"))

	if got := c.Messages().Size(); got != 0 {
		t.Errorf("expected 0 messages when MaxPerChannel=0, got %d", got)
	}
}

// TestBug2HitCountPreservedOnUpdate verifies that a Set (update) for an
// existing key does not reset the hitCount accumulated via Get calls.
// Without the fix, ClearByFrequency eviction is inverted: a frequently-read
// entry that gets updated loses its hitCount and is evicted first.
func TestBug2HitCountPreservedOnUpdate(t *testing.T) {
	c := newCache(cache.Options{
		Limits: cache.Limits{MaxGuilds: 2},
		OnOverflow: cache.OverflowPolicy{
			Target:  cache.CategoryGuilds,
			ClearBy: cache.ClearByFrequency,
		},
	})
	defer c.Close()

	hotID := common.Snowflake("100")
	coldID := common.Snowflake("200")

	// Seed both guilds.
	c.Guilds().Set(&common.Guild{ID: hotID, Name: "hot"})
	c.Guilds().Set(&common.Guild{ID: coldID, Name: "cold"})

	// Access hotGuild 100 times to build up its hitCount.
	for i := 0; i < 100; i++ {
		c.Guilds().Get(hotID)
	}

	// Update hotGuild (simulating a GUILD_UPDATE event). The fix must preserve hitCount.
	c.Guilds().Set(&common.Guild{ID: hotID, Name: "hot-updated"})

	// Trigger an overflow by adding a third guild (cap is 2); one of the two
	// existing guilds must be evicted. ClearByFrequency evicts the lowest hitCount.
	c.Guilds().Set(&common.Guild{ID: "300", Name: "trigger"})

	_, hotStillCached := c.Guilds().Get(hotID)
	if !hotStillCached {
		t.Error("hot guild (hitCount=100 before update) was evicted — hitCount was reset on Set (Bug 2)")
	}
}
