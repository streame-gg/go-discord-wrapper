package mongocache_test

import (
	"context"
	"fmt"
	"log"
	"os"
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
	"github.com/streame-gg/go-discord-wrapper/types/common"
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
		log.Printf("skipping MongoDB integration tests: could not start container: %v", err)
		os.Exit(0)
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

// ── GuildStore ────────────────────────────────────────────────────────────────

func TestGuildStore_SetGetDelete(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Guilds().Set(guild("1"))

	got, ok := c.Guilds().Get("1")
	if !ok || got.ID != "1" {
		t.Fatalf("expected guild 1, ok=%v", ok)
	}
	if got.Name != "guild-1" {
		t.Fatalf("expected name 'guild-1', got %q", got.Name)
	}

	c.Guilds().Delete("1")
	if _, ok := c.Guilds().Get("1"); ok {
		t.Fatal("expected guild deleted")
	}
}

func TestGuildStore_All(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Guilds().Set(guild(fmt.Sprintf("%d", i)))
	}
	if len(c.Guilds().All()) != 5 {
		t.Fatalf("expected 5 guilds, got %d", len(c.Guilds().All()))
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
	c.Guilds().Delete("2")
	if s := c.Guilds().Size(); s != 2 {
		t.Fatalf("expected size 2 after delete, got %d", s)
	}
}

func TestGuildStore_Overwrite(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Guilds().Set(&common.Guild{ID: "1", Name: "old"})
	c.Guilds().Set(&common.Guild{ID: "1", Name: "new"})

	got, ok := c.Guilds().Get("1")
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
	if _, ok := c.Guilds().Get("1"); !ok {
		t.Fatal("expected guild before TTL")
	}

	time.Sleep(300 * time.Millisecond)

	if _, ok := c.Guilds().Get("1"); ok {
		t.Fatal("expected guild not returned after TTL (expires_at filter)")
	}
	if n := c.Guilds().Size(); n != 0 {
		t.Fatalf("expected Size() == 0 after TTL, got %d", n)
	}
	if all := c.Guilds().All(); len(all) != 0 {
		t.Fatalf("expected All() empty after TTL, got %d", len(all))
	}
}

// ── ChannelStore ──────────────────────────────────────────────────────────────

func TestChannelStore_CRUD(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Channels().Set(channel("10"))
	if _, ok := c.Channels().Get("10"); !ok {
		t.Fatal("expected channel 10")
	}
	c.Channels().Delete("10")
	if _, ok := c.Channels().Get("10"); ok {
		t.Fatal("expected channel deleted")
	}
}

func TestChannelStore_All(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 4; i++ {
		c.Channels().Set(channel(fmt.Sprintf("%d", i)))
	}
	if len(c.Channels().All()) != 4 {
		t.Fatalf("expected 4 channels, got %d", len(c.Channels().All()))
	}
}

func TestChannelStore_TTL(t *testing.T) {
	c := newCache(t, cache.Options{TTL: 150 * time.Millisecond})

	c.Channels().Set(channel("1"))
	time.Sleep(300 * time.Millisecond)
	if _, ok := c.Channels().Get("1"); ok {
		t.Fatal("expected channel expired after TTL")
	}
}

// ── UserStore ─────────────────────────────────────────────────────────────────

func TestUserStore_CRUD(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Users().Set(user("42"))
	got, ok := c.Users().Get("42")
	if !ok || got.Username != "user-42" {
		t.Fatalf("expected user-42, ok=%v", ok)
	}
	c.Users().Delete("42")
	if _, ok := c.Users().Get("42"); ok {
		t.Fatal("expected user deleted")
	}
}

func TestUserStore_All(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Users().Set(user("1"))
	c.Users().Set(user("2"))
	if len(c.Users().All()) != 2 {
		t.Fatalf("expected 2 users, got %d", len(c.Users().All()))
	}
}

// ── MemberStore ───────────────────────────────────────────────────────────────

func TestMemberStore_SetGet(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Members().Set("guild1", member("99"))
	got, ok := c.Members().Get("guild1", "99")
	if !ok || got.User.ID != "99" {
		t.Fatalf("expected member 99, ok=%v", ok)
	}
}

func TestMemberStore_Delete(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Members().Set("g1", member("5"))
	c.Members().Delete("g1", "5")
	if _, ok := c.Members().Get("g1", "5"); ok {
		t.Fatal("expected member deleted")
	}
}

func TestMemberStore_DeleteGuild(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 5; i++ {
		c.Members().Set("guild1", member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set("guild2", member("99"))

	c.Members().DeleteGuild("guild1")

	if members := c.Members().AllInGuild("guild1"); len(members) != 0 {
		t.Fatalf("expected 0 members in guild1 after DeleteGuild, got %d", len(members))
	}
	if members := c.Members().AllInGuild("guild2"); len(members) != 1 {
		t.Fatalf("guild2 should be unaffected, got %d members", len(members))
	}
}

func TestMemberStore_AllInGuild(t *testing.T) {
	c := newCache(t, cache.Options{})

	for i := 1; i <= 3; i++ {
		c.Members().Set("g1", member(fmt.Sprintf("%d", i)))
	}
	c.Members().Set("g2", member("99"))

	if all := c.Members().AllInGuild("g1"); len(all) != 3 {
		t.Fatalf("expected 3 members in g1, got %d", len(all))
	}
	if all := c.Members().AllInGuild("g2"); len(all) != 1 {
		t.Fatalf("expected 1 member in g2, got %d", len(all))
	}
}

func TestMemberStore_NilUser(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Members().Set("g1", &common.GuildMember{})
	if c.Members().Size() != 0 {
		t.Fatal("expected size 0 when member.User is nil")
	}
}

func TestMemberStore_TTL(t *testing.T) {
	c := newCache(t, cache.Options{TTL: 150 * time.Millisecond})

	c.Members().Set("g1", member("1"))
	time.Sleep(300 * time.Millisecond)
	if _, ok := c.Members().Get("g1", "1"); ok {
		t.Fatal("expected member expired after TTL")
	}
}

// ── StickerStore ──────────────────────────────────────────────────────────────

func TestStickerStore_CRUDAndSetAll(t *testing.T) {
	c := newCache(t, cache.Options{})

	c.Stickers().Set("guild1", sticker("s1"))
	got, ok := c.Stickers().Get("s1")
	if !ok || got.ID != "s1" {
		t.Fatalf("expected sticker s1, ok=%v", ok)
	}

	c.Stickers().SetAll("guild1", []*common.Sticker{sticker("s2"), sticker("s3")})
	if _, ok := c.Stickers().Get("s1"); ok {
		t.Fatal("expected old guild stickers replaced by SetAll")
	}
	if stickers := c.Stickers().GetByGuild("guild1"); len(stickers) != 2 {
		t.Fatalf("expected 2 stickers in guild1, got %d", len(stickers))
	}

	c.Stickers().Delete("s2")
	if _, ok := c.Stickers().Get("s2"); ok {
		t.Fatal("expected sticker s2 deleted")
	}

	c.Stickers().DeleteGuild("guild1")
	if stickers := c.Stickers().GetByGuild("guild1"); len(stickers) != 0 {
		t.Fatalf("expected guild stickers deleted, got %d", len(stickers))
	}
}

// ── MessageStore ──────────────────────────────────────────────────────────────

var defaultMsgOpts = cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}}

func TestMessageStore_AddGet(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	got, ok := c.Messages().Get("c1", "m1")
	if !ok || got.ID != "m1" {
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

	got, _ := c.Messages().Get("c1", "m1")
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

	msgs := c.Messages().Channel("c1")
	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(msgs))
	}
	// m3 was added last, so it must be first (newest).
	if msgs[0].ID != "m3" {
		t.Fatalf("expected m3 first (newest), got %s", msgs[0].ID)
	}
}

func TestMessageStore_Delete(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Delete("c1", "m1")
	if _, ok := c.Messages().Get("c1", "m1"); ok {
		t.Fatal("expected message deleted")
	}
}

func TestMessageStore_DeleteBulk(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	c.Messages().DeleteBulk("c1", []common.Snowflake{"m1", "m3", "m5"})

	for _, id := range []common.Snowflake{"m1", "m3", "m5"} {
		if _, ok := c.Messages().Get("c1", id); ok {
			t.Fatalf("expected message %s deleted", id)
		}
	}
	for _, id := range []common.Snowflake{"m2", "m4"} {
		if _, ok := c.Messages().Get("c1", id); !ok {
			t.Fatalf("expected message %s present", id)
		}
	}
}

func TestMessageStore_Channel_NewestFirst(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	for i := 1; i <= 5; i++ {
		c.Messages().Add(message(fmt.Sprintf("m%d", i), "c1"))
	}
	msgs := c.Messages().Channel("c1")
	if len(msgs) != 5 {
		t.Fatalf("expected 5 messages, got %d", len(msgs))
	}
	if msgs[0].ID != "m5" {
		t.Fatalf("expected newest first (m5), got %s", msgs[0].ID)
	}
	if msgs[4].ID != "m1" {
		t.Fatalf("expected oldest last (m1), got %s", msgs[4].ID)
	}
}

func TestMessageStore_DeleteChannel(t *testing.T) {
	c := newCache(t, defaultMsgOpts)

	c.Messages().Add(message("m1", "c1"))
	c.Messages().Add(message("m2", "c2"))

	c.Messages().DeleteChannel("c1")

	if msgs := c.Messages().Channel("c1"); len(msgs) != 0 {
		t.Fatalf("expected c1 messages deleted, got %d", len(msgs))
	}
	if msgs := c.Messages().Channel("c2"); len(msgs) != 1 {
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

	msgs := c.Messages().Channel("c1")
	if len(msgs) != 3 {
		t.Fatalf("expected ring cap of 3, got %d", len(msgs))
	}
	for _, m := range msgs {
		if m.ID == "m1" || m.ID == "m2" {
			t.Fatalf("message %s (oldest) should have been evicted", m.ID)
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
	if _, ok := c.Messages().Get("c1", "m1"); !ok {
		t.Fatal("expected message before TTL")
	}

	time.Sleep(300 * time.Millisecond)

	if _, ok := c.Messages().Get("c1", "m1"); ok {
		t.Fatal("expected message expired after TTL (expires_at filter)")
	}
	if msgs := c.Messages().Channel("c1"); len(msgs) != 0 {
		t.Fatalf("expected no messages after TTL, got %d", len(msgs))
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
	c.Messages().DeleteChannel("c1")
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
			gid := common.Snowflake(fmt.Sprintf("g%d", g%3))
			uid := common.Snowflake(fmt.Sprintf("u%d", g%4))
			cid := common.Snowflake(fmt.Sprintf("c%d", g%2))
			mid := common.Snowflake(fmt.Sprintf("msg%d", g))
			for i := 0; i < ops; i++ {
				c.Guilds().Set(guild(string(gid)))
				c.Guilds().Get(gid)
				c.Users().Set(user(string(uid)))
				c.Channels().Set(channel(string(cid)))
				c.Messages().Add(message(string(mid), string(cid)))
				c.Messages().Get(cid, mid)
				c.Members().Set(gid, member(string(uid)))
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

	const guildID = common.Snowflake("g1")
	const workers = 50

	initial := []*common.Emoji{
		{ID: "e1", Name: "emoji1"},
		{ID: "e2", Name: "emoji2"},
		{ID: "e3", Name: "emoji3"},
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
				if len(got) != 0 && len(got) != len(initial) && len(got) != 5 {
					partial.Store(true)
				}
			}
		}()
	}

	newSet := []*common.Emoji{
		{ID: "n1", Name: "new1"}, {ID: "n2", Name: "new2"},
		{ID: "n3", Name: "new3"}, {ID: "n4", Name: "new4"},
		{ID: "n5", Name: "new5"},
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

	ch := c.Messages().Channel("ch1")
	if len(ch) != 0 {
		t.Errorf("MaxPerChannel=0 should disable caching, got %d messages (Bug 36)", len(ch))
	}
	if sz := c.Messages().Size(); sz != 0 {
		t.Errorf("Size should be 0 when disabled, got %d (Bug 36)", sz)
	}
}
