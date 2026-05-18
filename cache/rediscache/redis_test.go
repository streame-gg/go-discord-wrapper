package rediscache_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/rediscache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type RedisCacheTestSuite struct {
	suite.Suite

	container tc.Container
	redisAddr string
}

func TestRedisCacheTestSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheTestSuite))
}

func (s *RedisCacheTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("failed to start MongoDB container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("container.MappedPort: %v", err)
	}

	s.container = container
	s.redisAddr = fmt.Sprintf("%s:%s", host, port.Port())
}

func (s *RedisCacheTestSuite) TearDownSuite() {
	if s.container == nil {
		return
	}
	ctx := context.Background()
	_ = s.container.Terminate(ctx)
}

func (s *RedisCacheTestSuite) newCache(opts cache.Options) *rediscache.RedisCache {
	s.T().Helper()
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = client.Close() })
	c := rediscache.NewRedisCache(client, opts).WithKeyPrefix("test:" + s.T().Name())
	s.T().Cleanup(func() { _ = c.Close() })
	return c
}

func guild(id string) *discord.Guild {
	return &discord.Guild{ID: discord.Snowflake(id), Name: "guild-" + id}
}

func channel(id string) *discord.Channel {
	return &discord.Channel{ID: discord.Snowflake(id), Name: "chan-" + id}
}

func user(id string) *discord.User {
	return &discord.User{ID: discord.Snowflake(id), Username: "user-" + id}
}

func member(userID string) *discord.GuildMember {
	return &discord.GuildMember{User: user(userID)}
}

func sticker(id string) *discord.Sticker {
	return &discord.Sticker{ID: discord.Snowflake(id), Name: "sticker-" + id}
}

func message(id, channelID string) *discord.Message {
	return &discord.Message{
		ID:        discord.Snowflake(id),
		ChannelID: discord.Snowflake(channelID),
		Content:   "msg-" + id,
	}
}

// Guild Store

func (s *RedisCacheTestSuite) TestGuildStore_SetGetDelete() {
	c := s.newCache(cache.Options{})

	c.Guilds().Set(guild("1"))

	got, ok := c.Guilds().Get("1")
	s.Require().True(ok, "expected guild to exist after set")
	s.Require().Equal(discord.Snowflake("1"), got.ID)
	s.Require().Equal("guild-1", got.Name)

	c.Guilds().Delete("1")

	_, ok = c.Guilds().Get("1")
	s.Require().False(ok, "expected guild to not exist after delete")
}

func (s *RedisCacheTestSuite) TestGuildStore_All() {
	c := s.newCache(cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	all := c.Guilds().All()
	s.Require().Equal(5, all.Len(), "expected 5 guilds after 5 sets")
}

func (s *RedisCacheTestSuite) TestGuildStore_Size() {
	c := s.newCache(cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	s.Require().Equal(3, c.Guilds().Size(), "expected size 3")
	c.Guilds().Delete("2")
	s.Require().Equal(2, c.Guilds().Size(), "expected size 2 after delete")
}

func (s *RedisCacheTestSuite) TestGuildStore_TTL() {
	c := s.newCache(cache.Options{TTL: 150 * time.Millisecond})

	c.Guilds().Set(guild("1"))
	_, ok := c.Guilds().Get("1")
	s.Require().True(ok, "expected guild to exist after set")

	time.Sleep(300 * time.Millisecond)

	_, ok = c.Guilds().Get("1")
	s.Require().False(ok, "expected guild to not exist after TTL")
}

func (s *RedisCacheTestSuite) TestGuildStore_Overwrite() {
	c := s.newCache(cache.Options{})

	c.Guilds().Set(&discord.Guild{ID: "1", Name: "old"})
	c.Guilds().Set(&discord.Guild{ID: "1", Name: "new"})

	got, ok := c.Guilds().Get("1")
	s.Require().True(ok, "expected guild to exist after set")
	s.Require().Equal("new", got.Name, "expected new guild to have the new name")
	s.Require().Equal(1, c.Guilds().Size(), "expected overwrite, therefore no increase")
}

// Channel Store

func (s *RedisCacheTestSuite) TestChannelStore_CRUD() {
	c := s.newCache(cache.Options{})

	c.Channels().Set(channel("10"))

	_, ok := c.Channels().Get("10")
	s.Require().True(ok, "expected channel 10")

	c.Channels().Delete("10")

	_, ok = c.Channels().Get("10")
	s.Require().False(ok, "expected channel deleted")
}

func (s *RedisCacheTestSuite) TestChannelStore_All() {
	c := s.newCache(cache.Options{})

	for i := 1; i <= 4; i++ {
		c.Channels().Set(channel(fmt.Sprintf("%d", i)))
	}
	s.Require().Equal(4, c.Channels().All().Len(), "expected 4 channels")
}

// ── UserStore ─────────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestUserStore_CRUD() {
	c := s.newCache(cache.Options{})

	c.Users().Set(user("42"))

	got, ok := c.Users().Get("42")
	s.Require().True(ok, "expected user 42 to exist")
	s.Require().Equal("user-42", got.Username)

	c.Users().Delete("42")

	_, ok = c.Users().Get("42")
	s.Require().False(ok, "expected user deleted")
}

func (s *RedisCacheTestSuite) TestUserStore_All() {
	c := s.newCache(cache.Options{})

	c.Users().Set(user("1"))
	c.Users().Set(user("2"))
	s.Require().Equal(2, c.Users().All().Len(), "expected 2 users")
}

// ── MemberStore ───────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestMemberStore_SetGet() {
	c := s.newCache(cache.Options{})

	c.Members().Set("guild1", member("99"))

	got, ok := c.Members().Get("guild1", "99")
	s.Require().True(ok, "expected member 99 in guild1")
	s.Require().NotNil(got)
	s.Require().NotNil(got.User)
	s.Require().Equal(discord.Snowflake("99"), got.User.ID)
}

func (s *RedisCacheTestSuite) TestMemberStore_Delete() {
	c := s.newCache(cache.Options{})

	c.Members().Set("guild1", member("5"))

	c.Members().Delete("guild1", "5")

	_, ok := c.Members().Get("guild1", "5")
	s.Require().False(ok, "expected member deleted")
}

func (s *RedisCacheTestSuite) TestMemberStore_DeleteGuild() {
	c := s.newCache(cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Members().Set("guild1", member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set("guild2", member("99"))

	c.Members().DeleteGuild("guild1")

	s.Require().Equal(0, c.Members().AllInGuild("guild1").Len(), "expected 0 members in guild1")
	s.Require().Equal(1, c.Members().AllInGuild("guild2").Len(), "expected 1 member in guild2")
}

func (s *RedisCacheTestSuite) TestMemberStore_AllInGuild() {
	c := s.newCache(cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Members().Set("g1", member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set("g2", member("99"))

	all := c.Members().AllInGuild("g1")
	s.Require().Equal(3, all.Len(), "expected 3 members in g1")
}

func (s *RedisCacheTestSuite) TestMemberStore_NilUser() {
	c := s.newCache(cache.Options{})

	c.Members().Set("g1", &discord.GuildMember{})

	s.Require().Equal(0, c.Members().Size(), "expected size 0 when member.User is nil")
}

func (s *RedisCacheTestSuite) TestMemberStore_TTL() {
	c := s.newCache(cache.Options{TTL: 150 * time.Millisecond})

	c.Members().Set("g1", member("1"))

	_, ok := c.Members().Get("g1", "1")
	s.Require().True(ok, "expected member before TTL")

	time.Sleep(300 * time.Millisecond)

	_, ok = c.Members().Get("g1", "1")
	s.Require().False(ok, "expected member expired after TTL")
}

// ── StickerStore ──────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestStickerStore_CRUDAndSetAll() {
	c := s.newCache(cache.Options{})

	c.Stickers().Set("guild1", sticker("s1"))

	got, ok := c.Stickers().Get("s1")
	s.Require().True(ok, "expected sticker s1")
	s.Require().NotNil(got)
	s.Require().Equal(discord.Snowflake("s1"), got.ID)

	c.Stickers().SetAll("guild1", []*discord.Sticker{sticker("s2"), sticker("s3")})

	_, ok = c.Stickers().Get("s1")
	s.Require().False(ok, "expected old guild stickers replaced by SetAll")

	s.Require().Equal(2, c.Stickers().GetByGuild("guild1").Len(), "expected 2 stickers in guild1")

	c.Stickers().Delete("s2")
	_, ok = c.Stickers().Get("s2")
	s.Require().False(ok, "expected sticker s2 deleted")

	c.Stickers().DeleteGuild("guild1")
	s.Require().Equal(0, c.Stickers().GetByGuild("guild1").Len(), "expected guild stickers deleted")
}

// ── MessageStore ──────────────────────────────────────────────────────────────

var defaultMsgOpts = cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}}

func (s *RedisCacheTestSuite) TestMessageStore_AddGet() {
	c := s.newCache(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))

	got, ok := c.Messages().Get("c1", "m1")
	s.Require().True(ok, "expected msg m1")
	s.Require().NotNil(got)
	s.Require().Equal(discord.Snowflake("m1"), got.ID)
}

func (s *RedisCacheTestSuite) TestMessageStore_Update() {
	c := s.newCache(defaultMsgOpts)

	msg := message("m1", "c1")
	c.Messages().Add(msg)

	updated := *msg
	updated.Content = "edited"
	c.Messages().Update(&updated)

	got, ok := c.Messages().Get("c1", "m1")
	s.Require().True(ok, "expected message to exist after update")
	s.Require().NotNil(got)
	s.Require().Equal("edited", got.Content)
}

func (s *RedisCacheTestSuite) TestMessageStore_Update_NonExistent() {
	c := s.newCache(defaultMsgOpts)

	// Updating a non-existent message must not create it.
	c.Messages().Update(message("ghost", "c1"))

	_, ok := c.Messages().Get("c1", "ghost")
	s.Require().False(ok, "Update should not create a new entry")
}

func (s *RedisCacheTestSuite) TestMessageStore_Delete() {
	c := s.newCache(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Delete("c1", "m1")

	_, ok := c.Messages().Get("c1", "m1")
	s.Require().False(ok, "expected message deleted")
}

func (s *RedisCacheTestSuite) TestMessageStore_DeleteBulk() {
	c := s.newCache(defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	c.Messages().DeleteBulk("c1", []discord.Snowflake{"m1", "m3", "m5"})

	for _, id := range []discord.Snowflake{"m1", "m3", "m5"} {
		_, ok := c.Messages().Get("c1", id)
		s.Require().False(ok, "expected message %s deleted", id)
	}
	for _, id := range []discord.Snowflake{"m2", "m4"} {
		_, ok := c.Messages().Get("c1", id)
		s.Require().True(ok, "expected message %s present", id)
	}
}

func (s *RedisCacheTestSuite) TestMessageStore_Channel_NewestFirst() {
	c := s.newCache(defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	msgs := c.Messages().Channel("c1")
	s.Require().Equal(5, msgs.Len(), "expected 5 messages")

	// The most recently added message must be first.
	s.Require().Equal(discord.Snowflake("m5"), msgs.Values()[0].ID, "expected newest first")
}

func (s *RedisCacheTestSuite) TestMessageStore_DeleteChannel() {
	c := s.newCache(defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Add(message("m2", "c2"))

	c.Messages().DeleteChannel("c1")

	s.Require().Equal(0, c.Messages().Channel("c1").Len(), "expected c1 messages deleted")
	s.Require().Equal(1, c.Messages().Channel("c2").Len(), "c2 should be unaffected")
}

func (s *RedisCacheTestSuite) TestMessageStore_RingCap() {
	c := s.newCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 3},
	})

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}

	msgs := c.Messages().Channel("c1")
	s.Require().Equal(3, msgs.Len(), "expected ring cap of 3")

	// Oldest messages (m1, m2) must have been evicted.
	for _, m := range msgs.Values() {
		s.Require().NotEqual(discord.Snowflake("m1"), m.ID, "m1 should have been evicted")
		s.Require().NotEqual(discord.Snowflake("m2"), m.ID, "m2 should have been evicted")
	}
}

func (s *RedisCacheTestSuite) TestMessageStore_TTL() {
	c := s.newCache(cache.Options{
		Messages: cache.MessageOptions{
			MaxPerChannel: 10,
			TTL:           150 * time.Millisecond,
		},
	})

	c.Messages().Add(message("m1", "c1"))
	_, ok := c.Messages().Get("c1", "m1")
	s.Require().True(ok, "expected message before TTL")

	time.Sleep(300 * time.Millisecond)

	_, ok = c.Messages().Get("c1", "m1")
	s.Require().False(ok, "expected message expired after TTL")
}

func (s *RedisCacheTestSuite) TestMessageStore_SizeAccuracy() {
	c := s.newCache(defaultMsgOpts)

	for i := 1; i <= 4; i++ {
		c.Messages().Add(message(fmt.Sprintf("c1m%d", i), "c1"))
	}
	for i := 1; i <= 3; i++ {
		c.Messages().Add(message(fmt.Sprintf("c2m%d", i), "c2"))
	}

	s.Require().Equal(7, c.Messages().Size(), "expected size 7")

	c.Messages().DeleteChannel("c1")

	s.Require().Equal(3, c.Messages().Size(), "expected size 3 after DeleteChannel")
}

// ── Key prefix isolation ──────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestKeyPrefix_Isolation() {
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = client.Close() })

	c1 := rediscache.NewRedisCache(client, cache.Options{}).WithKeyPrefix("bot-a")
	c2 := rediscache.NewRedisCache(client, cache.Options{}).WithKeyPrefix("bot-b")
	s.T().Cleanup(func() { _ = c1.Close(); _ = c2.Close() })

	c1.Guilds().Set(guild("42"))

	_, ok := c2.Guilds().Get("42")
	s.Require().False(ok, "bot-b should not see bot-a's guild")

	_, ok = c1.Guilds().Get("42")
	s.Require().True(ok, "bot-a should see its own guild")
}

// ── Close ─────────────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestClose_Idempotent() {
	c := s.newCache(cache.Options{})

	err := c.Close()
	s.Require().NoError(err, "first Close should succeed")

	err = c.Close()
	s.Require().NoError(err, "second Close should be idempotent")
}

// ── Concurrent access (race detector) ────────────────────────────────────────

func (s *RedisCacheTestSuite) TestConcurrent_NoRace() {
	c := s.newCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 10},
	})

	const goroutines = 10
	const ops = 30
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			gid := discord.Snowflake(fmt.Sprintf("g%d", g%3))
			uid := discord.Snowflake(fmt.Sprintf("u%d", g%4))
			cid := discord.Snowflake(fmt.Sprintf("c%d", g%2))
			mid := discord.Snowflake(fmt.Sprintf("msg%d", g))
			for i := 0; i < ops; i++ {
				c.Guilds().Set(guild(string(gid)))
				c.Guilds().Get(gid)
				c.Users().Set(user(string(uid)))
				c.Channels().Set(channel(string(cid)))
				c.Messages().Add(message(string(mid), string(cid)))
				c.Messages().Get(cid, mid)
				c.Members().Set(gid, member(string(uid)))
			}
		}()
	}

	wg.Wait()
}

// TestBug6NetworkErrorDoesNotPruneIndex verifies that a Redis network error
// during Get does not remove the entry from the index set. Only a true
// cache-miss (redis.Nil) should trigger pruning.
func (s *RedisCacheTestSuite) TestBug6NetworkErrorDoesNotPruneIndex() {
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	defer client.Close()

	c := rediscache.NewRedisCache(client, cache.Options{})
	defer c.Close()

	// Seed a guild into the cache.
	c.Guilds().Set(&discord.Guild{ID: "g1", Name: "test"})

	// Confirm it's in the index.
	idx := client.SMembers(context.Background(), "discord:guild:index").Val()
	found := false
	for _, v := range idx {
		if v == "g1" {
			found = true
		}
	}
	s.Require().True(found, "guild not found in index before test")

	// Simulate a network error by closing the underlying connection pool.
	// After this, Get returns an error — the index must NOT be pruned.
	_ = client.Close()

	// This Get will fail with a connection error.
	c.Guilds().Get("g1")

	// Re-open a fresh client to inspect the index.
	client2 := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	defer client2.Close()

	idx2 := client2.SMembers(context.Background(), "discord:guild:index").Val()
	found2 := false
	for _, v := range idx2 {
		if v == "g1" {
			found2 = true
		}
	}
	s.Require().True(found2, "index was pruned on network error (Bug 6) — only true cache-miss should prune")
}

// TestBug7UpdateIsAtomicNoGhostEntry verifies that Update uses an atomic SetXX
// so a key that expires between the old Exists check and the Set call does not
// produce a ghost entry (message value without a sorted-set index entry).
func (s *RedisCacheTestSuite) TestBug7UpdateIsAtomicNoGhostEntry() {
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	defer client.Close()

	// Use a 50 ms TTL so expiry is easy to trigger in the test.
	c := rediscache.NewRedisCache(client, cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 10, TTL: 50 * time.Millisecond},
	})
	defer c.Close()

	msg := &discord.Message{ID: "m1", ChannelID: "c1", Content: "original"}
	c.Messages().Add(msg)

	// Let the TTL expire so the key no longer exists.
	time.Sleep(100 * time.Millisecond)

	// Update after expiry — with the old code (Exists+Set) this would create
	// a new key with no matching ZSet entry. With SetXX it's a no-op.
	c.Messages().Update(&discord.Message{ID: "m1", ChannelID: "c1", Content: "updated"})

	// The key must not exist (SetXX must be a no-op after expiry).
	exists := client.Exists(context.Background(), "discord:msg:c1:m1").Val()
	s.Require().Equal(int64(0), exists, "Update created a ghost entry after TTL expiry — SetXX not used atomically (Bug 7)")
}

// TestBug8AddIsAtomicParallelAddsRespectMaxPerChannel verifies that concurrent
// Add calls never leave more than MaxPerChannel entries in the sorted-set index.
// Before the fix, non-atomic ZCard/ZPopMin could miss excess entries under load.
func (s *RedisCacheTestSuite) TestBug8AddIsAtomicParallelAddsRespectMaxPerChannel() {
	const maxAmount = 5
	const goroutines = 50

	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	defer client.Close()

	c := rediscache.NewRedisCache(client, cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: maxAmount},
	})
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			c.Messages().Add(&discord.Message{
				ID:        discord.Snowflake(fmt.Sprintf("m%d", i)),
				ChannelID: "ch-atomic",
				Content:   "payload",
			})
		}()
	}
	wg.Wait()

	// The sorted-set cardinality must not exceed MaxPerChannel.
	card := client.ZCard(context.Background(), "discord:msg:ch:ch-atomic").Val()
	s.Require().LessOrEqual(card, int64(maxAmount), "ZCard should not exceed MaxPerChannel (Bug 8)")

	// The number of msg keys for this channel must equal the ZCard.
	msgKeys := client.Keys(context.Background(), "discord:msg:ch-atomic:*").Val()
	s.Require().Equal(int64(len(msgKeys)), card, "msg key count must match ZCard — orphan keys (Bug 8)")
}

// TestBug36MaxPerChannelZeroDisables verifies that MaxPerChannel=0 disables
// message caching in the Redis backend (Bug 36).
func (s *RedisCacheTestSuite) TestBug36MaxPerChannelZeroDisables() {
	c := s.newCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 0},
	})

	c.Messages().Add(message("m1", "ch1"))
	c.Messages().Add(message("m2", "ch1"))

	s.Require().Equal(0, c.Messages().Channel("ch1").Len(), "MaxPerChannel=0 should disable caching (Bug 36)")
	s.Require().Equal(0, c.Messages().Size(), "Size should be 0 when disabled (Bug 36)")
}

// TestBug50SetAllIsAtomic verifies that readers always see either the complete
// old set or the complete new set — never a partial or empty guild (Bug 50).
func (s *RedisCacheTestSuite) TestBug50SetAllIsAtomic() {
	c := s.newCache(cache.Options{})

	const guildID = discord.Snowflake("g1")
	const workers = 100

	// Seed an initial non-empty emoji set.
	initial := []*discord.Emoji{
		{ID: "e1", Name: "emoji1"},
		{ID: "e2", Name: "emoji2"},
		{ID: "e3", Name: "emoji3"},
	}
	c.Emojis().SetAll(guildID, initial)

	var partial atomic.Bool

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Readers: observe the guild emoji list in a tight loop.
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

	// Writer: repeatedly call SetAll with a different 5-emoji set.
	newSet := []*discord.Emoji{
		{ID: "n1", Name: "new1"}, {ID: "n2", Name: "new2"},
		{ID: "n3", Name: "new3"}, {ID: "n4", Name: "new4"},
		{ID: "n5", Name: "new5"},
	}
	for i := 0; i < 20; i++ {
		c.Emojis().SetAll(guildID, newSet)
		c.Emojis().SetAll(guildID, initial)
	}
	close(stop)
	wg.Wait()

	s.Require().False(partial.Load(), "reader observed a partial emoji set during concurrent SetAll (Bug 50)")
}

// TestWithKeyPrefix_IndependentLifecycle verifies that closing a WithKeyPrefix
// derivative does not cancel the original instance's context (P2-31).
func (s *RedisCacheTestSuite) TestWithKeyPrefix_IndependentLifecycle() {
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = client.Close() })

	base := rediscache.NewRedisCache(client, cache.Options{})
	s.T().Cleanup(func() { _ = base.Close() })

	prefixed := base.WithKeyPrefix("p2-31-test")

	// Write via prefixed, then close it.
	g := guild("lifecycle-1")
	base.Guilds().Set(g)

	s.Require().NoError(prefixed.Close())

	// base must still be functional after the derivative is closed.
	got, ok := base.Guilds().Get(g.ID)
	s.Require().True(ok, "base cache must remain functional after closing a WithKeyPrefix derivative")
	s.Equal(g.ID, got.ID)
}
