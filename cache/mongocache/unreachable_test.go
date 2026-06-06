package mongocache_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mgopts "go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/mongocache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mongoUnreachableSuite holds the Docker-free resilience tests: they point a
// MongoDBCache at an unreachable MongoDB (a closed loopback port with a tiny
// server-selection timeout) and drive every store method. They do not assert
// correctness of stored data — that is what the testcontainer suite verifies —
// but they do assert the cache's resilience contract: when MongoDB is
// unavailable, every method must degrade gracefully (no panic, reads return
// empty/false, sizes are zero) and must not block the caller indefinitely. As a
// side effect they exercise the filter/document-construction and error-handling
// branches of every store without a database, so the package retains meaningful
// coverage in environments with neither Docker nor a MONGO_TEST_URI.
type mongoUnreachableSuite struct{ suite.Suite }

func TestMongoUnreachableSuite(t *testing.T) { suite.Run(t, new(mongoUnreachableSuite)) }

// unreachableCache builds a MongoDBCache whose driver points at a closed
// loopback port. Operations fail fast via the short server-selection timeout.
func unreachableCache(t *testing.T, opts cache.Options) *mongocache.MongoDBCache {
	t.Helper()
	cl, err := mongo.Connect(mgopts.Client().
		ApplyURI("mongodb://127.0.0.1:1").
		SetServerSelectionTimeout(40 * time.Millisecond).
		SetConnectTimeout(40 * time.Millisecond))
	require.NoError(t, err)
	t.Cleanup(func() { _ = cl.Disconnect(context.Background()) })

	db := cl.Database("resilience_test")
	c := mongocache.NewMongoDBCache(db, opts)
	c.EnableSyncWrites()
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func roleOf(id string) *discord.Role {
	return &discord.Role{ID: mustSnowflake(id), Name: "role-" + id}
}

func voiceStateOf(userID string) *discord.VoiceState {
	return &discord.VoiceState{UserID: mustSnowflake(userID), SessionID: "sess-" + userID}
}

func soundOf(id string) *discord.SoundboardSound {
	return &discord.SoundboardSound{SoundID: mustSnowflake(id), Name: "sound-" + id}
}

func scheduledEventOf(id, guildID string) *discord.GuildScheduledEvent {
	return &discord.GuildScheduledEvent{ID: mustSnowflake(id), GuildID: mustSnowflake(guildID), Name: "event-" + id}
}

func stageInstanceOf(id, guildID string) *discord.StageInstance {
	return &discord.StageInstance{ID: mustSnowflake(id), GuildID: mustSnowflake(guildID), Topic: "topic-" + id}
}

func emojiOf(id string) *discord.Emoji {
	return &discord.Emoji{ID: mustSnowflake(id), Name: "emoji-" + id}
}

func presenceOf(userID, guildID string) *discord.Presence {
	return &discord.Presence{
		GuildID: mustSnowflake(guildID),
		User:    discord.PartialPresenceUser{ID: mustSnowflake(userID)},
	}
}

func banOf(userID string) *discord.Ban {
	return &discord.Ban{User: *user(userID)}
}

func autoModRuleOf(id, guildID string) *discord.AutoModerationRule {
	return &discord.AutoModerationRule{ID: mustSnowflake(id), GuildID: mustSnowflake(guildID), Name: "rule-" + id}
}

func inviteOf(code string) *discord.Invite {
	return &discord.Invite{Code: code}
}

// TestAllStoresDegradeGracefully drives every store method with the backend down,
// asserting reads come back empty and nothing panics. It runs once with a TTL
// configured (exercising the expires_at filter branches) and once without
// (exercising the no-filter branches).
func (s *mongoUnreachableSuite) TestAllStoresDegradeGracefully() {
	for _, tc := range []struct {
		name string
		opts cache.Options
	}{
		// WriteTimeout bounds every write's per-call context as a belt-and-braces
		// guard so a write cannot outlive the server-selection timeout.
		{"withTTL", cache.Options{TTL: time.Minute, WriteTimeout: 150 * time.Millisecond, Messages: cache.MessageOptions{MaxPerChannel: 100}}},
		{"noTTL", cache.Options{WriteTimeout: 150 * time.Millisecond, Messages: cache.MessageOptions{MaxPerChannel: 100}}},
	} {
		s.Run(tc.name, func() {
			c := unreachableCache(s.T(), tc.opts)

			const gid = "10"
			const uid = "1"

			s.Run("guilds", func() {
				st := c.Guilds()
				st.Set(guild(uid))
				st.Set(nil) // nil guard
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.All().Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
			})

			s.Run("channels", func() {
				st := c.Channels()
				st.Set(channel(uid))
				st.Set(nil)
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.All().Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
			})

			s.Run("users", func() {
				st := c.Users()
				st.Set(user(uid))
				st.Set(nil)
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.All().Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
			})

			s.Run("members", func() {
				st := c.Members()
				st.Set(mustSnowflake(gid), member(uid))
				st.Set(mustSnowflake(gid), nil)
				_, ok := st.Get(mustSnowflake(gid), mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.AllInGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(gid), mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("roles", func() {
				st := c.Roles()
				st.Set(mustSnowflake(gid), roleOf(uid))
				st.Set(mustSnowflake(gid), nil)
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.All().Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("messages", func() {
				st := c.Messages()
				st.Add(message(uid, gid))
				st.Add(nil)
				st.Update(message(uid, gid))
				_, ok := st.Get(mustSnowflake(gid), mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.Channel(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(gid), mustSnowflake(uid))
				st.DeleteBulk(mustSnowflake(gid), []discord.Snowflake{mustSnowflake(uid)})
				st.DeleteBulk(mustSnowflake(gid), nil) // empty guard
				st.DeleteChannel(mustSnowflake(gid))
			})

			s.Run("voiceStates", func() {
				st := c.VoiceStates()
				st.Set(mustSnowflake(gid), voiceStateOf(uid))
				st.Set(mustSnowflake(gid), nil)
				_, ok := st.Get(mustSnowflake(gid), mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.AllInGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(gid), mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("soundboard", func() {
				st := c.Soundboard()
				st.Set(mustSnowflake(gid), soundOf(uid))
				st.Set(mustSnowflake(gid), nil)
				// SetAll is intentionally not exercised here: it runs through
				// session.WithTransaction, whose transient-error retry is bounded by
				// its own ~120s timer (not the call context) and so hangs against an
				// unreachable server. Its correctness is covered by the container suite.
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("scheduledEvents", func() {
				st := c.ScheduledEvents()
				st.Set(scheduledEventOf(uid, gid))
				st.Set(nil)
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("stageInstances", func() {
				st := c.StageInstances()
				st.Set(stageInstanceOf(uid, gid))
				st.Set(nil)
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("emojis", func() {
				st := c.Emojis()
				st.Set(mustSnowflake(gid), emojiOf(uid))
				st.Set(mustSnowflake(gid), nil)
				// SetAll uses the transaction path — see the soundboard note above.
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("stickers", func() {
				st := c.Stickers()
				st.Set(mustSnowflake(gid), sticker(uid))
				st.Set(mustSnowflake(gid), nil)
				// SetAll uses the transaction path — see the soundboard note above.
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("presences", func() {
				st := c.Presences()
				st.Set(presenceOf(uid, gid))
				st.Set(nil)
				_, ok := st.Get(mustSnowflake(gid), mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(gid), mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("bans", func() {
				st := c.Bans()
				st.Set(mustSnowflake(gid), banOf(uid))
				st.Set(mustSnowflake(gid), nil)
				_, ok := st.Get(mustSnowflake(gid), mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.AllInGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(gid), mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("autoModRules", func() {
				st := c.AutoModRules()
				st.Set(mustSnowflake(gid), autoModRuleOf(uid, gid))
				st.Set(mustSnowflake(gid), nil)
				_, ok := st.Get(mustSnowflake(uid))
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete(mustSnowflake(uid))
				st.DeleteGuild(mustSnowflake(gid))
			})

			s.Run("invites", func() {
				st := c.Invites()
				st.Set(inviteOf("code1"))
				st.SetWithGuild(mustSnowflake(gid), inviteOf("code2"))
				st.Set(nil)
				st.Set(inviteOf("")) // empty-code guard
				_, ok := st.Get("code1")
				s.False(ok)
				s.Zero(st.GetByGuild(mustSnowflake(gid)).Len())
				s.Zero(st.Size())
				st.Delete("code1")
				st.DeleteGuild(mustSnowflake(gid))
			})
		})
	}
}

// TestMessagesDisabled covers the MaxPerChannel == 0 branches: Add is a no-op and
// Channel returns empty without touching the backend.
func (s *mongoUnreachableSuite) TestMessagesDisabled() {
	c := unreachableCache(s.T(), cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 0}})
	st := c.Messages()
	st.Add(message("1", "10")) // disabled → returns immediately
	s.Zero(st.Channel(mustSnowflake("10")).Len())
}

// TestAsyncWriteQueue exercises the async (non-sync) write path: enqueueWrite
// dispatching to the worker, the queue-full drop branch, and a clean drain on
// Close. The backend is unreachable so each queued write blocks the single
// worker on server selection, which lets us fill (and overflow) the small queue
// deterministically.
func (s *mongoUnreachableSuite) TestAsyncWriteQueue() {
	cl, err := mongo.Connect(mgopts.Client().
		ApplyURI("mongodb://127.0.0.1:1").
		SetServerSelectionTimeout(40 * time.Millisecond).
		SetConnectTimeout(40 * time.Millisecond))
	s.Require().NoError(err)
	s.T().Cleanup(func() { _ = cl.Disconnect(context.Background()) })

	// Tiny queue + no EnableSyncWrites → writes go through the worker, and once
	// the worker is busy the buffer overflows and excess writes are dropped.
	c := mongocache.NewMongoDBCache(cl.Database("resilience_async"),
		cache.Options{WriteQueueSize: 1, WriteTimeout: 40 * time.Millisecond})

	g := c.Guilds()
	for i := 0; i < 50; i++ {
		g.Set(guild("1")) // most of these hit the drop branch; none may block
	}

	// Close must drain the worker and return without hanging.
	done := make(chan struct{})
	go func() { _ = c.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		s.FailNow("Close did not return — write worker failed to drain")
	}

	// Writes after Close are silently dropped, not panics.
	g.Set(guild("2"))
}
