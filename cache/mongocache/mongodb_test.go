package mongocache_test

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mgopts "go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/mongocache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── container setup ───────────────────────────────────────────────────────────

var mongoClient *mongo.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "mongo:7",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("failed to start MongoDB container: %v", err)
	}
	defer container.Terminate(ctx) //nolint:errcheck

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		log.Fatalf("container.MappedPort: %v", err)
	}
	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	mongoClient, err = mongo.Connect(mgopts.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo.Connect: %v", err)
	}
	defer mongoClient.Disconnect(ctx) //nolint:errcheck

	os.Exit(m.Run())
}

// ── helpers ───────────────────────────────────────────────────────────────────

func mustSnowflake(s string) discord.Snowflake {
	n, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return discord.Snowflake(n)
	}
	h := fnv.New64a()
	h.Write([]byte(s))
	return discord.Snowflake(h.Sum64())
}

// newCache creates a MongoDBCache backed by a fresh database named after the
// test. The database is dropped in t.Cleanup so tests are fully isolated.
func newCache(t *testing.T, opts cache.Options) *mongocache.MongoDBCache {
	t.Helper()
	// Replace characters that are invalid in MongoDB database names.
	dbName := "test_" + strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db := mongoClient.Database(dbName)
	t.Cleanup(func() {
		_ = db.Drop(context.Background())
	})
	c := mongocache.NewMongoDBCache(db, opts)
	t.Cleanup(func() { c.Close() })
	return c
}

func guild(id string) *discord.Guild {
	return &discord.Guild{ID: mustSnowflake(id), Name: "guild-" + id}
}

func channel(id string) *discord.Channel {
	return &discord.Channel{ID: mustSnowflake(id), Name: "chan-" + id}
}

func user(id string) *discord.User {
	return &discord.User{ID: mustSnowflake(id), Username: "user-" + id}
}

func member(userID string) *discord.GuildMember {
	return &discord.GuildMember{User: user(userID)}
}

func sticker(id string) *discord.Sticker {
	return &discord.Sticker{ID: mustSnowflake(id), Name: "sticker-" + id}
}

func message(id, channelID string) *discord.Message {
	return &discord.Message{
		ID:        mustSnowflake(id),
		ChannelID: mustSnowflake(channelID),
		Content:   "msg-" + id,
	}
}

// ── GuildStore ────────────────────────────────────────────────────────────────

func TestGuildStore_SetGetDelete(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Guilds().Set(guild("1"))

	got, ok := c.Guilds().Get(1)
	if !ok || got.ID != 1 {
		t.Fatalf("expected guild 1, ok=%v", ok)
	}
	if got.Name != "guild-1" {
		t.Fatalf("expected name 'guild-1', got %q", got.Name)
	}

	c.Guilds().Delete(1)
	if _, ok := c.Guilds().Get(1); ok {
		t.Fatal("expected guild deleted")
	}
}

func TestGuildStore_All(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	if n := c.Guilds().All().Len(); n != 5 {
		t.Fatalf("expected 5 guilds, got %d", n)
	}
}

func TestGuildStore_Size(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	if s := c.Guilds().Size(); s != 3 {
		t.Fatalf("expected size 3, got %d", s)
	}
	c.Guilds().Delete(2)
	if s := c.Guilds().Size(); s != 2 {
		t.Fatalf("expected size 2 after delete, got %d", s)
	}
}

func TestGuildStore_Overwrite(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Guilds().Set(&discord.Guild{ID: 1, Name: "old"})
	c.Guilds().Set(&discord.Guild{ID: 1, Name: "new"})

	got, ok := c.Guilds().Get(1)
	if !ok || got.Name != "new" {
		t.Fatalf("expected overwritten name 'new', got %q (ok=%v)", got.Name, ok)
	}
	if c.Guilds().Size() != 1 {
		t.Fatalf("expected size 1 after overwrite, got %d", c.Guilds().Size())
	}
}

// TTL is enforced at query time via the expires_at filter, so it takes effect
// immediately without waiting for MongoDB's background TTL index sweep.
func TestGuildStore_TTL(t *testing.T) {
	c := newCache(t, cache.Options{TTL: 150 * time.Millisecond})

	c.Guilds().Set(guild("1"))
	if _, ok := c.Guilds().Get(1); !ok {
		t.Fatal("expected guild before TTL")
	}

	time.Sleep(300 * time.Millisecond)

	if _, ok := c.Guilds().Get(1); ok {
		t.Fatal("expected guild not returned after TTL (expires_at filter)")
	}
	if n := c.Guilds().Size(); n != 0 {
		t.Fatalf("expected Size() == 0 after TTL, got %d", n)
	}
	if n := c.Guilds().All().Len(); n != 0 {
		t.Fatalf("expected All() empty after TTL, got %d", n)
	}
}

// ── ChannelStore ──────────────────────────────────────────────────────────────

func TestChannelStore_CRUD(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Channels().Set(channel("10"))
	if _, ok := c.Channels().Get(10); !ok {
		t.Fatal("expected channel 10")
	}
	c.Channels().Delete(10)
	if _, ok := c.Channels().Get(10); ok {
		t.Fatal("expected channel deleted")
	}
}

func TestChannelStore_All(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 4; i++ {
		c.Channels().Set(channel(fmt.Sprintf("%d", i)))
	}
	if n := c.Channels().All().Len(); n != 4 {
		t.Fatalf("expected 4 channels, got %d", n)
	}
}

func TestChannelStore_TTL(t *testing.T) {
	c := newCache(t, cache.Options{TTL: 150 * time.Millisecond})

	c.Channels().Set(channel("1"))
	time.Sleep(300 * time.Millisecond)
	if _, ok := c.Channels().Get(1); ok {
		t.Fatal("expected channel expired after TTL")
	}
}

// ── UserStore ─────────────────────────────────────────────────────────────────

func TestUserStore_CRUD(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Users().Set(user("42"))
	got, ok := c.Users().Get(42)
	if !ok || got.Username != "user-42" {
		t.Fatalf("expected user-42, ok=%v", ok)
	}
	c.Users().Delete(42)
	if _, ok := c.Users().Get(42); ok {
		t.Fatal("expected user deleted")
	}
}

func TestUserStore_All(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Users().Set(user("1"))
	c.Users().Set(user("2"))
	if c.Users().All().Len() != 2 {
		t.Fatalf("expected 2 users, got %d", c.Users().All().Len())
	}
}

// ── MemberStore ───────────────────────────────────────────────────────────────

func TestMemberStore_SetGet(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Members().Set(mustSnowflake("guild1"), member("99"))
	got, ok := c.Members().Get(mustSnowflake("guild1"), 99)
	if !ok || got.User.ID != 99 {
		t.Fatalf("expected member 99, ok=%v", ok)
	}
}

func TestMemberStore_Delete(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Members().Set(mustSnowflake("g1"), member("5"))
	c.Members().Delete(mustSnowflake("g1"), 5)
	if _, ok := c.Members().Get(mustSnowflake("g1"), 5); ok {
		t.Fatal("expected member deleted")
	}
}

func TestMemberStore_DeleteGuild(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Members().Set(mustSnowflake("guild1"), member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set(mustSnowflake("guild2"), member("99"))

	c.Members().DeleteGuild(mustSnowflake("guild1"))

	if n := c.Members().AllInGuild(mustSnowflake("guild1")).Len(); n != 0 {
		t.Fatalf("expected 0 members in guild1 after DeleteGuild, got %d", n)
	}
	if n := c.Members().AllInGuild(mustSnowflake("guild2")).Len(); n != 1 {
		t.Fatalf("guild2 should be unaffected, got %d members", n)
	}
}

func TestMemberStore_AllInGuild(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Members().Set(mustSnowflake("g1"), member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set(mustSnowflake("g2"), member("99"))

	if n := c.Members().AllInGuild(mustSnowflake("g1")).Len(); n != 3 {
		t.Fatalf("expected 3 members in g1, got %d", n)
	}
	if n := c.Members().AllInGuild(mustSnowflake("g2")).Len(); n != 1 {
		t.Fatalf("expected 1 member in g2, got %d", n)
	}
}

func TestMemberStore_NilUser(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Members().Set(mustSnowflake("g1"), &discord.GuildMember{})
	if c.Members().Size() != 0 {
		t.Fatal("expected size 0 when member.User is nil")
	}
}

func TestMemberStore_TTL(t *testing.T) {
	c := newCache(t, cache.Options{TTL: 150 * time.Millisecond})

	c.Members().Set(mustSnowflake("g1"), member("1"))
	time.Sleep(300 * time.Millisecond)
	if _, ok := c.Members().Get(mustSnowflake("g1"), 1); ok {
		t.Fatal("expected member expired after TTL")
	}
}

// ── StickerStore ──────────────────────────────────────────────────────────────

func TestStickerStore_CRUDAndSetAll(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Stickers().Set(mustSnowflake("guild1"), sticker("s1"))
	got, ok := c.Stickers().Get(mustSnowflake("s1"))
	if !ok || got.ID != mustSnowflake("s1") {
		t.Fatalf("expected sticker s1, ok=%v", ok)
	}

	c.Stickers().SetAll(mustSnowflake("guild1"), []*discord.Sticker{sticker("s2"), sticker("s3")})
	if _, ok := c.Stickers().Get(mustSnowflake("s1")); ok {
		t.Fatal("expected old guild stickers replaced by SetAll")
	}
	if n := c.Stickers().GetByGuild(mustSnowflake("guild1")).Len(); n != 2 {
		t.Fatalf("expected 2 stickers in guild1, got %d", n)
	}

	c.Stickers().Delete(mustSnowflake("s2"))
	if _, ok := c.Stickers().Get(mustSnowflake("s2")); ok {
		t.Fatal("expected sticker s2 deleted")
	}

	c.Stickers().DeleteGuild(mustSnowflake("guild1"))
	if n := c.Stickers().GetByGuild(mustSnowflake("guild1")).Len(); n != 0 {
		t.Fatalf("expected guild stickers deleted, got %d", n)
	}
}

// ── MessageStore ──────────────────────────────────────────────────────────────

var defaultMsgOpts = cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}}

func TestMessageStore_AddGet(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	got, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	if !ok || got.ID != mustSnowflake("m1") {
		t.Fatalf("expected msg m1, ok=%v", ok)
	}
}

func TestMessageStore_Update(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	msg := message("m1", "c1")
	c.Messages().Add(msg)

	updated := *msg
	updated.Content = "edited"
	c.Messages().Update(&updated)

	got, _ := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	if got.Content != "edited" {
		t.Fatalf("expected updated content, got %q", got.Content)
	}
}

func TestMessageStore_Update_PreservesInsertedAt(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))

	// Update does not change ordering — the message stays in its original
	// position relative to others in the channel.
	c.Messages().Add(message("m2", "c1"))
	c.Messages().Add(message("m3", "c1"))

	updated := message("m1", "c1")
	updated.Content = "edited"
	c.Messages().Update(updated)

	msgs := c.Messages().Channel(mustSnowflake("c1"))
	if msgs.Len() != 3 {
		t.Fatalf("expected 3 messages, got %d", msgs.Len())
	}
	// m3 was added last, so it must be first (newest).
	if msgs.Values()[0].ID != mustSnowflake("m3") {
		t.Fatalf("expected m3 first (newest), got %v", msgs.Values()[0].ID)
	}
}

func TestMessageStore_Delete(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Delete(mustSnowflake("c1"), mustSnowflake("m1"))
	if _, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1")); ok {
		t.Fatal("expected message deleted")
	}
}

func TestMessageStore_DeleteBulk(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	c.Messages().DeleteBulk(mustSnowflake("c1"), []discord.Snowflake{mustSnowflake("m1"), mustSnowflake("m3"), mustSnowflake("m5")})

	for _, id := range []discord.Snowflake{mustSnowflake("m1"), mustSnowflake("m3"), mustSnowflake("m5")} {
		if _, ok := c.Messages().Get(mustSnowflake("c1"), id); ok {
			t.Fatalf("expected message %v deleted", id)
		}
	}
	for _, id := range []discord.Snowflake{mustSnowflake("m2"), mustSnowflake("m4")} {
		if _, ok := c.Messages().Get(mustSnowflake("c1"), id); !ok {
			t.Fatalf("expected message %v present", id)
		}
	}
}

func TestMessageStore_Channel_NewestFirst(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	msgs := c.Messages().Channel(mustSnowflake("c1"))
	if msgs.Len() != 5 {
		t.Fatalf("expected 5 messages, got %d", msgs.Len())
	}
	if msgs.Values()[0].ID != mustSnowflake("m5") {
		t.Fatalf("expected newest first (m5), got %v", msgs.Values()[0].ID)
	}
	if msgs.Values()[4].ID != mustSnowflake("m1") {
		t.Fatalf("expected oldest last (m1), got %v", msgs.Values()[4].ID)
	}
}

func TestMessageStore_DeleteChannel(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Add(message("m2", "c2"))

	c.Messages().DeleteChannel(mustSnowflake("c1"))

	if n := c.Messages().Channel(mustSnowflake("c1")).Len(); n != 0 {
		t.Fatalf("expected c1 messages deleted, got %d", n)
	}
	if n := c.Messages().Channel(mustSnowflake("c2")).Len(); n != 1 {
		t.Fatal("c2 should be unaffected")
	}
}

func TestMessageStore_RingCap(t *testing.T) {
	c := newCache(t, cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 3},
	})

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}

	msgs := c.Messages().Channel(mustSnowflake("c1"))
	if msgs.Len() != 3 {
		t.Fatalf("expected ring cap of 3, got %d", msgs.Len())
	}
	for _, m := range msgs.Values() {
		if m.ID == mustSnowflake("m1") || m.ID == mustSnowflake("m2") {
			t.Fatalf("message %v (oldest) should have been evicted", m.ID)
		}
	}
}

func TestMessageStore_TTL(t *testing.T) {
	c := newCache(t, cache.Options{
		Messages: cache.MessageOptions{
			MaxPerChannel: 10,
			TTL:           150 * time.Millisecond,
		},
	})

	c.Messages().Add(message("m1", "c1"))
	if _, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1")); !ok {
		t.Fatal("expected message before TTL")
	}

	time.Sleep(300 * time.Millisecond)

	if _, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1")); ok {
		t.Fatal("expected message expired after TTL (expires_at filter)")
	}
	if n := c.Messages().Channel(mustSnowflake("c1")).Len(); n != 0 {
		t.Fatalf("expected no messages after TTL, got %d", n)
	}
}

func TestMessageStore_SizeAccuracy(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	for i := 1; i <= 4; i++ {
		c.Messages().Add(message(fmt.Sprintf("c1m%d", i), "c1"))
	}
	for i := 1; i <= 3; i++ {
		c.Messages().Add(message(fmt.Sprintf("c2m%d", i), "c2"))
	}
	if s := c.Messages().Size(); s != 7 {
		t.Fatalf("expected size 7, got %d", s)
	}
	c.Messages().DeleteChannel(mustSnowflake("c1"))
	if s := c.Messages().Size(); s != 3 {
		t.Fatalf("expected size 3 after DeleteChannel, got %d", s)
	}
}

// ── Close ─────────────────────────────────────────────────────────────────────

func TestClose_Idempotent(t *testing.T) {
	c := newCache(t, cache.Options{})
	if err := c.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	// Second call must be a no-op (no panic, no error).
	if err := c.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

// ── Concurrent access (race detector) ────────────────────────────────────────

func TestConcurrent_NoRace(t *testing.T) {
	c := newCache(t, cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 10},
	})

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
func TestBug50SetAllIsAtomic(t *testing.T) {
	c := newCache(t, cache.Options{})

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

	if partial.Load() {
		t.Error("reader observed a partial emoji set during concurrent SetAll (Bug 50)")
	}
}

// TestBug36MaxPerChannelZeroDisables verifies that MaxPerChannel=0 disables
// message caching in the MongoDB backend (Bug 36).
func TestBug36MaxPerChannelZeroDisables(t *testing.T) {
	c := newCache(t, cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 0},
	})

	c.Messages().Add(message("m1", "ch1"))
	c.Messages().Add(message("m2", "ch1"))

	ch := c.Messages().Channel(mustSnowflake("ch1"))
	if ch.Len() != 0 {
		t.Errorf("MaxPerChannel=0 should disable caching, got %d messages (Bug 36)", ch.Len())
	}
	if sz := c.Messages().Size(); sz != 0 {
		t.Errorf("Size should be 0 when disabled, got %d (Bug 36)", sz)
	}
}
