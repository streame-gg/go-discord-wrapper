package mongocache_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/mongocache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mongoLegacySuite ports the original per-store integration tests to the
// testify/suite form used by the rest of the package. Each method builds an
// isolated database via the shared newCache helper (keyed on s.T()).
type mongoLegacySuite struct {
	suite.Suite
}

func TestMongoLegacySuite(t *testing.T) { suite.Run(t, new(mongoLegacySuite)) }

func (s *mongoLegacySuite) new(opts cache.Options) *mongocache.MongoDBCache {
	return newCache(s.T(), opts)
}

// ── GuildStore ────────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestGuildStore_SetGetDelete() {
	c := s.new(cache.Options{})

	c.Guilds().Set(guild("1"))

	got, ok := c.Guilds().Get(1)
	s.Require().True(ok)
	s.Equal(discord.Snowflake(1), got.ID)
	s.Equal("guild-1", got.Name)

	c.Guilds().Delete(1)
	_, ok = c.Guilds().Get(1)
	s.False(ok)
}

func (s *mongoLegacySuite) TestGuildStore_All() {
	c := s.new(cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	s.Equal(5, c.Guilds().All().Len())
}

func (s *mongoLegacySuite) TestGuildStore_Size() {
	c := s.new(cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	s.Equal(3, c.Guilds().Size())
	c.Guilds().Delete(2)
	s.Equal(2, c.Guilds().Size())
}

func (s *mongoLegacySuite) TestGuildStore_Overwrite() {
	c := s.new(cache.Options{})

	c.Guilds().Set(&discord.Guild{ID: 1, Name: "old"})
	c.Guilds().Set(&discord.Guild{ID: 1, Name: "new"})

	got, ok := c.Guilds().Get(1)
	s.Require().True(ok)
	s.Equal("new", got.Name)
	s.Equal(1, c.Guilds().Size())
}

// TTL is enforced at query time via the expires_at filter, so it takes effect
// immediately without waiting for MongoDB's background TTL index sweep.
func (s *mongoLegacySuite) TestGuildStore_TTL() {
	c := s.new(cache.Options{TTL: 150 * time.Millisecond})

	c.Guilds().Set(guild("1"))
	_, ok := c.Guilds().Get(1)
	s.Require().True(ok)

	time.Sleep(300 * time.Millisecond)

	_, ok = c.Guilds().Get(1)
	s.False(ok)
	s.Equal(0, c.Guilds().Size())
	s.Equal(0, c.Guilds().All().Len())
}

// ── ChannelStore ──────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestChannelStore_CRUD() {
	c := s.new(cache.Options{})

	c.Channels().Set(channel("10"))
	_, ok := c.Channels().Get(10)
	s.Require().True(ok)
	c.Channels().Delete(10)
	_, ok = c.Channels().Get(10)
	s.False(ok)
}

func (s *mongoLegacySuite) TestChannelStore_All() {
	c := s.new(cache.Options{})

	for i := 1; i <= 4; i++ {
		c.Channels().Set(channel(fmt.Sprintf("%d", i)))
	}
	s.Equal(4, c.Channels().All().Len())
}

func (s *mongoLegacySuite) TestChannelStore_TTL() {
	c := s.new(cache.Options{TTL: 150 * time.Millisecond})

	c.Channels().Set(channel("1"))
	time.Sleep(300 * time.Millisecond)
	_, ok := c.Channels().Get(1)
	s.False(ok)
}

// ── UserStore ─────────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestUserStore_CRUD() {
	c := s.new(cache.Options{})

	c.Users().Set(user("42"))
	got, ok := c.Users().Get(42)
	s.Require().True(ok)
	s.Equal("user-42", got.Username)
	c.Users().Delete(42)
	_, ok = c.Users().Get(42)
	s.False(ok)
}

func (s *mongoLegacySuite) TestUserStore_All() {
	c := s.new(cache.Options{})

	c.Users().Set(user("1"))
	c.Users().Set(user("2"))
	s.Equal(2, c.Users().All().Len())
}

// ── MemberStore ───────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestMemberStore_SetGet() {
	c := s.new(cache.Options{})

	c.Members().Set(mustSnowflake("guild1"), member("99"))
	got, ok := c.Members().Get(mustSnowflake("guild1"), 99)
	s.Require().True(ok)
	s.Equal(discord.Snowflake(99), got.User.ID)
}

func (s *mongoLegacySuite) TestMemberStore_Delete() {
	c := s.new(cache.Options{})

	c.Members().Set(mustSnowflake("g1"), member("5"))
	c.Members().Delete(mustSnowflake("g1"), 5)
	_, ok := c.Members().Get(mustSnowflake("g1"), 5)
	s.False(ok)
}

func (s *mongoLegacySuite) TestMemberStore_DeleteGuild() {
	c := s.new(cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Members().Set(mustSnowflake("guild1"), member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set(mustSnowflake("guild2"), member("99"))

	c.Members().DeleteGuild(mustSnowflake("guild1"))

	s.Equal(0, c.Members().AllInGuild(mustSnowflake("guild1")).Len())
	s.Equal(1, c.Members().AllInGuild(mustSnowflake("guild2")).Len())
}

func (s *mongoLegacySuite) TestMemberStore_AllInGuild() {
	c := s.new(cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Members().Set(mustSnowflake("g1"), member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set(mustSnowflake("g2"), member("99"))

	s.Equal(3, c.Members().AllInGuild(mustSnowflake("g1")).Len())
	s.Equal(1, c.Members().AllInGuild(mustSnowflake("g2")).Len())
}

func (s *mongoLegacySuite) TestMemberStore_NilUser() {
	c := s.new(cache.Options{})

	c.Members().Set(mustSnowflake("g1"), &discord.GuildMember{})
	s.Equal(0, c.Members().Size())
}

func (s *mongoLegacySuite) TestMemberStore_TTL() {
	c := s.new(cache.Options{TTL: 150 * time.Millisecond})

	c.Members().Set(mustSnowflake("g1"), member("1"))
	time.Sleep(300 * time.Millisecond)
	_, ok := c.Members().Get(mustSnowflake("g1"), 1)
	s.False(ok)
}

// ── StickerStore ──────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestStickerStore_CRUDAndSetAll() {
	c := s.new(cache.Options{})

	c.Stickers().Set(mustSnowflake("guild1"), sticker("s1"))
	got, ok := c.Stickers().Get(mustSnowflake("s1"))
	s.Require().True(ok)
	s.Equal(mustSnowflake("s1"), got.ID)

	c.Stickers().SetAll(mustSnowflake("guild1"), []*discord.Sticker{sticker("s2"), sticker("s3")})
	_, ok = c.Stickers().Get(mustSnowflake("s1"))
	s.False(ok, "old guild stickers replaced by SetAll")
	s.Equal(2, c.Stickers().GetByGuild(mustSnowflake("guild1")).Len())

	c.Stickers().Delete(mustSnowflake("s2"))
	_, ok = c.Stickers().Get(mustSnowflake("s2"))
	s.False(ok)

	c.Stickers().DeleteGuild(mustSnowflake("guild1"))
	s.Equal(0, c.Stickers().GetByGuild(mustSnowflake("guild1")).Len())
}

// ── MessageStore ──────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestMessageStore_AddGet() {
	c := s.new(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	got, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.Require().True(ok)
	s.Equal(mustSnowflake("m1"), got.ID)
}

func (s *mongoLegacySuite) TestMessageStore_Update() {
	c := s.new(defaultMsgOpts)

	msg := message("m1", "c1")
	c.Messages().Add(msg)

	updated := *msg
	updated.Content = "edited"
	c.Messages().Update(&updated)

	got, _ := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.Equal("edited", got.Content)
}

func (s *mongoLegacySuite) TestMessageStore_Update_PreservesInsertedAt() {
	c := s.new(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Add(message("m2", "c1"))
	c.Messages().Add(message("m3", "c1"))

	updated := message("m1", "c1")
	updated.Content = "edited"
	c.Messages().Update(updated)

	msgs := c.Messages().Channel(mustSnowflake("c1"))
	s.Require().Equal(3, msgs.Len())
	// m3 was added last, so it must be first (newest).
	s.Equal(mustSnowflake("m3"), msgs.Values()[0].ID)
}

func (s *mongoLegacySuite) TestMessageStore_Delete() {
	c := s.new(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Delete(mustSnowflake("c1"), mustSnowflake("m1"))
	_, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.False(ok)
}

func (s *mongoLegacySuite) TestMessageStore_DeleteBulk() {
	c := s.new(defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	c.Messages().DeleteBulk(mustSnowflake("c1"), []discord.Snowflake{mustSnowflake("m1"), mustSnowflake("m3"), mustSnowflake("m5")})

	for _, id := range []discord.Snowflake{mustSnowflake("m1"), mustSnowflake("m3"), mustSnowflake("m5")} {
		_, ok := c.Messages().Get(mustSnowflake("c1"), id)
		s.False(ok, "message %v deleted", id)
	}
	for _, id := range []discord.Snowflake{mustSnowflake("m2"), mustSnowflake("m4")} {
		_, ok := c.Messages().Get(mustSnowflake("c1"), id)
		s.True(ok, "message %v present", id)
	}
}

func (s *mongoLegacySuite) TestMessageStore_Channel_NewestFirst() {
	c := s.new(defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	msgs := c.Messages().Channel(mustSnowflake("c1"))
	s.Require().Equal(5, msgs.Len())
	s.Equal(mustSnowflake("m5"), msgs.Values()[0].ID, "newest first")
	s.Equal(mustSnowflake("m1"), msgs.Values()[4].ID, "oldest last")
}

func (s *mongoLegacySuite) TestMessageStore_DeleteChannel() {
	c := s.new(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Add(message("m2", "c2"))

	c.Messages().DeleteChannel(mustSnowflake("c1"))

	s.Equal(0, c.Messages().Channel(mustSnowflake("c1")).Len())
	s.Equal(1, c.Messages().Channel(mustSnowflake("c2")).Len())
}

func (s *mongoLegacySuite) TestMessageStore_RingCap() {
	c := s.new(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 3}})

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}

	msgs := c.Messages().Channel(mustSnowflake("c1"))
	s.Require().Equal(3, msgs.Len())
	for _, m := range msgs.Values() {
		s.NotEqual(mustSnowflake("m1"), m.ID, "oldest evicted")
		s.NotEqual(mustSnowflake("m2"), m.ID, "oldest evicted")
	}
}

func (s *mongoLegacySuite) TestMessageStore_TTL() {
	c := s.new(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 10, TTL: 150 * time.Millisecond}})

	c.Messages().Add(message("m1", "c1"))
	_, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.Require().True(ok)

	time.Sleep(300 * time.Millisecond)

	_, ok = c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.False(ok)
	s.Equal(0, c.Messages().Channel(mustSnowflake("c1")).Len())
}

func (s *mongoLegacySuite) TestMessageStore_SizeAccuracy() {
	c := s.new(defaultMsgOpts)

	for i := 1; i <= 4; i++ {
		c.Messages().Add(message(fmt.Sprintf("c1m%d", i), "c1"))
	}
	for i := 1; i <= 3; i++ {
		c.Messages().Add(message(fmt.Sprintf("c2m%d", i), "c2"))
	}
	s.Equal(7, c.Messages().Size())
	c.Messages().DeleteChannel(mustSnowflake("c1"))
	s.Equal(3, c.Messages().Size())
}

// ── Close ─────────────────────────────────────────────────────────────────────

func (s *mongoLegacySuite) TestClose_Idempotent() {
	c := s.new(cache.Options{})
	s.Require().NoError(c.Close())
	s.Require().NoError(c.Close(), "second Close is a no-op")
}

// ── Concurrent access (race detector) ────────────────────────────────────────

func (s *mongoLegacySuite) TestConcurrent_NoRace() {
	c := s.new(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 10}})

	const goroutines = 10
	const ops = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			gid := mustSnowflake(fmt.Sprintf("g%d", g%3))
			cid := mustSnowflake(fmt.Sprintf("c%d", g%2))
			mid := mustSnowflake(fmt.Sprintf("msg%d", g))
			gidStr := fmt.Sprintf("g%d", g%3)
			uidStr := fmt.Sprintf("u%d", g%4)
			cidStr := fmt.Sprintf("c%d", g%2)
			midStr := fmt.Sprintf("msg%d", g)
			for i := 0; i < ops; i++ {
				c.Guilds().Set(guild(gidStr))
				c.Guilds().Get(gid)
				c.Users().Set(user(uidStr))
				c.Channels().Set(channel(cidStr))
				c.Messages().Add(message(midStr, cidStr))
				c.Messages().Get(cid, mid)
				c.Members().Set(gid, member(uidStr))
				c.Members().AllInGuild(gid)
			}
		}()
	}

	wg.Wait()
}

// TestBug50SetAllIsAtomic verifies that concurrent readers always observe a
// complete guild emoji set — never a partial or empty state (Bug 50).
func (s *mongoLegacySuite) TestBug50SetAllIsAtomic() {
	c := s.new(cache.Options{})

	guildID := mustSnowflake("g1")
	const workers = 50

	initial := []*discord.Emoji{
		{ID: mustSnowflake("e1"), Name: "emoji1"},
		{ID: mustSnowflake("e2"), Name: "emoji2"},
		{ID: mustSnowflake("e3"), Name: "emoji3"},
	}
	c.Emojis().SetAll(guildID, initial)

	var partial atomic.Bool

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				got := c.Emojis().GetByGuild(guildID)
				if got.Len() != 0 && got.Len() != len(initial) && got.Len() != 5 {
					partial.Store(true)
				}
			}
		}()
	}

	newSet := []*discord.Emoji{
		{ID: mustSnowflake("n1"), Name: "new1"}, {ID: mustSnowflake("n2"), Name: "new2"},
		{ID: mustSnowflake("n3"), Name: "new3"}, {ID: mustSnowflake("n4"), Name: "new4"},
		{ID: mustSnowflake("n5"), Name: "new5"},
	}
	for i := 0; i < 10; i++ {
		c.Emojis().SetAll(guildID, newSet)
		c.Emojis().SetAll(guildID, initial)
	}
	close(stop)
	wg.Wait()

	s.False(partial.Load(), "reader observed a partial emoji set during concurrent SetAll (Bug 50)")
}

// TestBug36MaxPerChannelZeroDisables verifies that MaxPerChannel=0 disables
// message caching in the MongoDB backend (Bug 36).
func (s *mongoLegacySuite) TestBug36MaxPerChannelZeroDisables() {
	c := s.new(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 0}})

	c.Messages().Add(message("m1", "ch1"))
	c.Messages().Add(message("m2", "ch1"))

	s.Equal(0, c.Messages().Channel(mustSnowflake("ch1")).Len(), "MaxPerChannel=0 disables caching (Bug 36)")
	s.Equal(0, c.Messages().Size(), "Size 0 when disabled (Bug 36)")
}
