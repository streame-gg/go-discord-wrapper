package mongocache_test

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

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

	// Fast path: when MONGO_TEST_URI points at an already-running MongoDB we
	// connect directly and skip Docker/testcontainers entirely. This lets the
	// integration suite run in environments without a Docker daemon.
	if uri := os.Getenv("MONGO_TEST_URI"); uri != "" {
		client, err := mongo.Connect(mgopts.Client().ApplyURI(uri))
		if err != nil {
			log.Fatalf("mongo.Connect(%q): %v", uri, err)
		}
		mongoClient = client
		defer client.Disconnect(ctx) //nolint:errcheck
		os.Exit(m.Run())
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "mongo:7",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		if isDockerUnavailable(err) {
			// Docker isn't available (e.g. CI without a Docker daemon). Run the
			// suite anyway so the Docker-free internal tests still execute; the
			// container-backed tests skip themselves when mongoClient is nil.
			log.Printf("Docker is not available — skipping MongoDB integration tests: %v", err)
			os.Exit(m.Run())
		}
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

func isDockerUnavailable(err error) bool {
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "provider") || strings.Contains(s, "docker host")
}

// ── helpers ───────────────────────────────────────────────────────────────────

// uniqueDBName derives a MongoDB-safe database name from a test name. Long
// testify/suite names (e.g. "TestMongoStoresSuite/TestStickerStore_…") overflow
// MongoDB's 64-byte database-name limit, which makes every write silently fail.
// We sanitise invalid characters, truncate, and append a short hash so the name
// stays both bounded (<=63 bytes) and unique per test.
func uniqueDBName(testName string) string {
	safe := strings.NewReplacer("/", "_", " ", "_", "$", "_", ".", "_").Replace(testName)
	h := fnv.New32a()
	_, _ = h.Write([]byte(testName))
	suffix := fmt.Sprintf("_%08x", h.Sum32())
	const max = 63
	if budget := max - len("test_") - len(suffix); len(safe) > budget {
		if budget < 0 {
			budget = 0
		}
		safe = safe[:budget]
	}
	return "test_" + safe + suffix
}

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
	if mongoClient == nil {
		t.Skip("MongoDB container unavailable (Docker not running)")
	}
	db := mongoClient.Database(uniqueDBName(t.Name()))
	t.Cleanup(func() {
		_ = db.Drop(context.Background())
	})
	c := mongocache.NewMongoDBCache(db, opts)
	c.EnableSyncWrites()
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

var defaultMsgOpts = cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}}
