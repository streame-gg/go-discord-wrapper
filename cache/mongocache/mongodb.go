// Package mongocache provides a MongoDB-backed implementation of [cache.Cache].
//
// Connect a [*mongo.Client] from go.mongodb.org/mongo-driver/mongo, select a
// database, then pass it to [NewMongoDBCache]:
//
//	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
//	db := client.Database("discord_cache")
//	c := mongocache.NewMongoDBCache(db, cache.Options{TTL: 30 * time.Minute})
//	connection.NewClient(token, intents, options.WithCache(c))
//
// # Collections
//
// Each entity type maps to its own collection:
//
//	guilds            — Guild documents
//	channels          — Channel documents
//	users             — User documents
//	members           — GuildMember documents keyed by "{guildID}:{userID}"
//	roles             — Role documents keyed by role ID
//	voice_states      — VoiceState documents keyed by "{guildID}:{userID}"
//	soundboard_sounds — SoundboardSound documents keyed by sound ID
//	scheduled_events  — GuildScheduledEvent documents keyed by event ID
//	stage_instances   — StageInstance documents keyed by instance ID
//	emojis            — Emoji documents keyed by emoji ID
//	stickers          — Sticker documents keyed by sticker ID
//	presences         — Presence documents keyed by "{guildID}:{userID}"
//	messages          — Message documents keyed by "{channelID}:{messageID}"
//
// # Document schema
//
//	_id        — entity key (Snowflake string or composite string)
//	json       — JSON-encoded entity payload
//	expires_at — expiry timestamp (omitted when Options.TTL == 0)
//	guild_id   — guild Snowflake (members collection only)
//	channel_id — channel Snowflake (messages collection only)
//	inserted_at — insertion time (messages collection only, used for ordering)
//
// # TTL and indexes
//
// [NewMongoDBCache] creates the following indexes automatically:
//   - Sparse TTL index on expires_at (all collections) — MongoDB deletes expired
//     documents in the background. Documents without expires_at never expire.
//   - Index on guild_id (members, roles, voice_states, soundboard_sounds,
//     scheduled_events, stage_instances, emojis, stickers, presences) — used by
//     GetByGuild/AllInGuild and DeleteGuild.
//   - Compound index on (channel_id, inserted_at) (messages) — used by Channel
//     and per-channel ring enforcement.
//
// Overflow and eviction options (Limits, OnOverflow, EvictBehavior) are ignored;
// use MongoDB's capped collections or TTL for overflow control.
package mongocache

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Document types ────────────────────────────────────────────────────────────

// entityDoc is the MongoDB document used for guilds, channels, and users.
type entityDoc struct {
	ID        string    `bson:"_id"`
	JSON      string    `bson:"json"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

// memberDoc is the MongoDB document used for guild members.
type memberDoc struct {
	ID        string    `bson:"_id"` // "{guildID}:{userID}"
	GuildID   string    `bson:"guild_id"`
	JSON      string    `bson:"json"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

// roleDoc is the MongoDB document used for guild roles.
type roleDoc struct {
	ID        string    `bson:"_id"` // roleID
	GuildID   string    `bson:"guild_id"`
	JSON      string    `bson:"json"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

// guildEntityDoc is a MongoDB document used for guild-scoped entities like
// soundboard sounds, emojis, scheduled events, and stage instances.
type guildEntityDoc struct {
	ID        string    `bson:"_id"`
	GuildID   string    `bson:"guild_id"`
	JSON      string    `bson:"json"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

// guildUserDoc is the MongoDB document used for per-guild user-scoped entities
// like voice states and presences.
type guildUserDoc struct {
	ID        string    `bson:"_id"` // "{guildID}:{userID}"
	GuildID   string    `bson:"guild_id"`
	JSON      string    `bson:"json"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

// msgDoc is the MongoDB document used for messages.
type msgDoc struct {
	ID         string    `bson:"_id"` // "{channelID}:{messageID}"
	ChannelID  string    `bson:"channel_id"`
	JSON       string    `bson:"json"`
	InsertedAt time.Time `bson:"inserted_at"`
	ExpiresAt  time.Time `bson:"expires_at,omitempty"`
}

// ── MongoDBCache ──────────────────────────────────────────────────────────────

// MongoDBCache implements [cache.Cache] using MongoDB.
// All methods and stores returned by Guilds/Channels/Users/Members/Messages
// are safe for concurrent use.
//
// The underlying [*mongo.Database] is not owned by MongoDBCache; the caller is
// responsible for connecting and disconnecting the client. Call
// [MongoDBCache.Close] on shutdown to cancel the internal context.
type MongoDBCache struct {
	db       *mongo.Database
	opts     cache.Options
	ctx      context.Context
	cancel   context.CancelFunc
	stopOnce sync.Once
	// msgChannelMu serializes insert+evict per channel to prevent TOCTOU races
	// when enforcing MaxPerChannel.
	msgChannelMu sync.Map // map[string]*sync.Mutex

	// write worker: all Set/Delete calls are dispatched here so the reader
	// goroutine never blocks on MongoDB I/O.
	writeCh      chan func()
	writeMu      sync.Mutex
	writerClosed bool
	writerDone   chan struct{}
	syncWrites   bool // true → execute writes inline (testing only)
}

// NewMongoDBCache creates a MongoDBCache backed by db and creates the required
// indexes. The context passed to MongoDB operations derives from a background
// context cancelled by [MongoDBCache.Close].
//
// opts.TTL controls the document expiry written to expires_at; zero means
// documents never expire. opts.Messages.MaxPerChannel caps the per-channel
// message ring (default 100).
func NewMongoDBCache(db *mongo.Database, opts cache.Options) *MongoDBCache {
	if opts.Messages.MaxPerChannel < 0 {
		opts.Messages.MaxPerChannel = 100
	}
	// MaxPerChannel == 0 means disabled (no messages cached) — leave as-is.
	if opts.Messages.TTL == 0 {
		opts.Messages.TTL = opts.TTL
	}
	queueSize := opts.WriteQueueSize
	if queueSize <= 0 {
		queueSize = 512
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &MongoDBCache{
		db:         db,
		opts:       opts,
		ctx:        ctx,
		cancel:     cancel,
		writeCh:    make(chan func(), queueSize),
		writerDone: make(chan struct{}),
	}
	go c.runWriteWorker()
	c.ensureIndexes()
	return c
}

// ensureIndexes creates all required MongoDB indexes idempotently.
func (c *MongoDBCache) ensureIndexes() {
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0).SetSparse(true),
	}
	for _, name := range []string{
		"guilds",
		"channels",
		"users",
		"members",
		"roles",
		"messages",
		"voice_states",
		"soundboard_sounds",
		"scheduled_events",
		"stage_instances",
		"emojis",
		"stickers",
		"presences",
	} {
		col := c.db.Collection(name)
		_, _ = col.Indexes().CreateOne(c.ctx, ttlIndex)
	}

	// members: index on guild_id for AllInGuild / DeleteGuild queries.
	_, _ = c.db.Collection("members").Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "guild_id", Value: 1}},
	})

	// roles: index on guild_id for GetByGuild / DeleteGuild queries.
	_, _ = c.db.Collection("roles").Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "guild_id", Value: 1}},
	})

	for _, name := range []string{
		"voice_states",
		"soundboard_sounds",
		"scheduled_events",
		"stage_instances",
		"emojis",
		"stickers",
		"presences",
	} {
		_, _ = c.db.Collection(name).Indexes().CreateOne(c.ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "guild_id", Value: 1}},
		})
	}

	// messages: compound index for per-channel ordered queries and ring enforcement.
	_, _ = c.db.Collection("messages").Indexes().CreateOne(c.ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "channel_id", Value: 1}, {Key: "inserted_at", Value: 1}},
	})
}

// expiresAt returns the expiry time for a new document given the configured TTL.
// Returns the zero time when TTL is zero (documents never expire).
func (c *MongoDBCache) expiresAt(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}
	return time.Now().Add(ttl)
}

// liveFilter returns a bson filter that excludes expired documents.
// When TTL is zero no expiry filter is needed.
func liveFilter(ttl time.Duration) bson.M {
	if ttl > 0 {
		return bson.M{"expires_at": bson.M{"$gt": time.Now()}}
	}
	return bson.M{}
}

func (c *MongoDBCache) Guilds() cache.GuildStore     { return &mongoGuildStore{c} }
func (c *MongoDBCache) Channels() cache.ChannelStore { return &mongoChannelStore{c} }
func (c *MongoDBCache) Users() cache.UserStore       { return &mongoUserStore{c} }
func (c *MongoDBCache) Members() cache.MemberStore   { return &mongoMemberStore{c} }
func (c *MongoDBCache) Roles() cache.RoleStore       { return &mongoRoleStore{c} }
func (c *MongoDBCache) Messages() cache.MessageStore { return &mongoMessageStore{c} }
func (c *MongoDBCache) VoiceStates() cache.VoiceStateStore {
	return &mongoVoiceStateStore{c}
}
func (c *MongoDBCache) Soundboard() cache.SoundboardStore {
	return &mongoSoundboardStore{c}
}
func (c *MongoDBCache) ScheduledEvents() cache.ScheduledEventStore {
	return &mongoScheduledEventStore{c}
}
func (c *MongoDBCache) StageInstances() cache.StageInstanceStore {
	return &mongoStageInstanceStore{c}
}
func (c *MongoDBCache) Emojis() cache.EmojiStore { return &mongoEmojiStore{c} }
func (c *MongoDBCache) Stickers() cache.StickerStore {
	return &mongoStickerStore{c}
}
func (c *MongoDBCache) Presences() cache.PresenceStore {
	return &mongoPresenceStore{c}
}

// Close drains the write queue, then cancels the internal context.
// The underlying [*mongo.Database] and its client are not closed.
// Safe to call multiple times.
func (c *MongoDBCache) Close() error {
	c.stopOnce.Do(func() {
		c.cancel()
		c.writeMu.Lock()
		c.writerClosed = true
		close(c.writeCh)
		c.writeMu.Unlock()
		<-c.writerDone
	})
	return nil
}

// EnableSyncWrites makes all write operations execute synchronously instead of
// going through the async write queue. Intended for testing only.
func (c *MongoDBCache) EnableSyncWrites() {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	c.syncWrites = true
}

// callCtx returns a per-call context for a write operation. When WriteTimeout
// is configured it carries a deadline that prevents a slow backend from
// blocking the write worker indefinitely.
func (c *MongoDBCache) callCtx() (context.Context, context.CancelFunc) {
	if c.opts.WriteTimeout > 0 {
		return context.WithTimeout(c.ctx, c.opts.WriteTimeout)
	}
	return c.ctx, func() {}
}

// enqueueWrite submits fn to the async write worker.
// In sync-writes mode (tests) fn is executed directly on the caller goroutine.
// If the buffer is full or the cache is closed the write is silently dropped.
func (c *MongoDBCache) enqueueWrite(fn func()) {
	c.writeMu.Lock()
	if c.syncWrites || c.writeCh == nil {
		c.writeMu.Unlock()
		fn()
		return
	}
	if c.writerClosed {
		c.writeMu.Unlock()
		return
	}
	select {
	case c.writeCh <- fn:
	default:
		// write queue full — drop write; a cache miss is preferable to blocking the reader
	}
	c.writeMu.Unlock()
}

// runWriteWorker drains c.writeCh until the channel is closed by Close.
func (c *MongoDBCache) runWriteWorker() {
	defer close(c.writerDone)
	for fn := range c.writeCh {
		fn()
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// setAllTx atomically replaces all guild-scoped documents in col with docs
// using a MongoDB transaction (requires Replica Set or sharded cluster).
// Falls back to a non-transactional delete+insert if the transaction cannot be
// started or committed (e.g. standalone deployment).
func setAllTx(ctx context.Context, col *mongo.Collection, guildID string, docs []interface{}) {
	session, err := col.Database().Client().StartSession()
	if err == nil {
		defer session.EndSession(ctx)
		_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
			if _, err := col.DeleteMany(ctx, bson.M{"guild_id": guildID}); err != nil {
				return nil, err
			}
			if len(docs) > 0 {
				if _, err := col.InsertMany(ctx, docs); err != nil {
					return nil, err
				}
			}
			return nil, nil
		})
	}
	if err != nil {
		// Standalone deployment or transaction failure — best-effort non-atomic.
		_, _ = col.DeleteMany(ctx, bson.M{"guild_id": guildID})
		if len(docs) > 0 {
			_, _ = col.InsertMany(ctx, docs)
		}
	}
}

func upsertByID(ctx context.Context, col *mongo.Collection, doc any) error {
	id := extractID(doc)
	_, err := col.ReplaceOne(
		ctx,
		bson.M{"_id": id},
		doc,
		options.Replace().SetUpsert(true),
	)
	return err
}

// extractID reads the _id field from an entityDoc, memberDoc, roleDoc, or msgDoc.
func extractID(doc any) string {
	switch d := doc.(type) {
	case entityDoc:
		return d.ID
	case memberDoc:
		return d.ID
	case roleDoc:
		return d.ID
	case guildEntityDoc:
		return d.ID
	case guildUserDoc:
		return d.ID
	case msgDoc:
		return d.ID
	}
	return ""
}

// ── Guild store ───────────────────────────────────────────────────────────────

type mongoGuildStore struct{ c *MongoDBCache }

func (s *mongoGuildStore) col() *mongo.Collection { return s.c.db.Collection("guilds") }

func (s *mongoGuildStore) Set(guild *discord.Guild) {
	if guild == nil {
		return
	}
	b, err := json.Marshal(guild)
	if err != nil {
		return
	}
	doc := entityDoc{
		ID:        strconv.FormatUint(uint64(guild.ID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoGuildStore) Get(id discord.Snowflake) (*discord.Guild, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(id), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc entityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var g discord.Guild
	if err := json.Unmarshal([]byte(doc.JSON), &g); err != nil {
		return nil, false
	}
	return &g, true
}

func (s *mongoGuildStore) Delete(id discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(id), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoGuildStore) All() *collection.Collection[discord.Snowflake, *discord.Guild] {
	coll := collection.New[discord.Snowflake, *discord.Guild]()
	cursor, err := s.col().Find(s.c.ctx, liveFilter(s.c.opts.TTL))
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []entityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var g discord.Guild
		if json.Unmarshal([]byte(d.JSON), &g) == nil {
			coll.Set(g.ID, &g)
		}
	}
	return coll
}

func (s *mongoGuildStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Channel store ─────────────────────────────────────────────────────────────

type mongoChannelStore struct{ c *MongoDBCache }

func (s *mongoChannelStore) col() *mongo.Collection { return s.c.db.Collection("channels") }

func (s *mongoChannelStore) Set(ch *discord.Channel) {
	if ch == nil {
		return
	}
	b, err := json.Marshal(ch)
	if err != nil {
		return
	}
	doc := entityDoc{
		ID:        strconv.FormatUint(uint64(ch.ID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoChannelStore) Get(id discord.Snowflake) (*discord.Channel, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(id), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc entityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var ch discord.Channel
	if err := json.Unmarshal([]byte(doc.JSON), &ch); err != nil {
		return nil, false
	}
	return &ch, true
}

func (s *mongoChannelStore) Delete(id discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(id), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoChannelStore) All() *collection.Collection[discord.Snowflake, *discord.Channel] {
	coll := collection.New[discord.Snowflake, *discord.Channel]()
	cursor, err := s.col().Find(s.c.ctx, liveFilter(s.c.opts.TTL))
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []entityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var ch discord.Channel
		if json.Unmarshal([]byte(d.JSON), &ch) == nil {
			coll.Set(ch.ID, &ch)
		}
	}
	return coll
}

func (s *mongoChannelStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── User store ────────────────────────────────────────────────────────────────

type mongoUserStore struct{ c *MongoDBCache }

func (s *mongoUserStore) col() *mongo.Collection { return s.c.db.Collection("users") }

func (s *mongoUserStore) Set(user *discord.User) {
	if user == nil {
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		return
	}
	doc := entityDoc{
		ID:        strconv.FormatUint(uint64(user.ID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoUserStore) Get(id discord.Snowflake) (*discord.User, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(id), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc entityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var u discord.User
	if err := json.Unmarshal([]byte(doc.JSON), &u); err != nil {
		return nil, false
	}
	return &u, true
}

func (s *mongoUserStore) Delete(id discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(id), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoUserStore) All() *collection.Collection[discord.Snowflake, *discord.User] {
	coll := collection.New[discord.Snowflake, *discord.User]()
	cursor, err := s.col().Find(s.c.ctx, liveFilter(s.c.opts.TTL))
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []entityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var u discord.User
		if json.Unmarshal([]byte(d.JSON), &u) == nil {
			coll.Set(u.ID, &u)
		}
	}
	return coll
}

func (s *mongoUserStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Member store ──────────────────────────────────────────────────────────────

type mongoMemberStore struct{ c *MongoDBCache }

func (s *mongoMemberStore) col() *mongo.Collection { return s.c.db.Collection("members") }

func memberID(guildID, userID discord.Snowflake) string {
	return strconv.FormatUint(uint64(guildID), 10) + ":" + strconv.FormatUint(uint64(userID), 10)
}

func guildUserID(guildID, userID discord.Snowflake) string {
	return strconv.FormatUint(uint64(guildID), 10) + ":" + strconv.FormatUint(uint64(userID), 10)
}

func (s *mongoMemberStore) Set(guildID discord.Snowflake, member *discord.GuildMember) {
	if member == nil || member.User == nil {
		return
	}
	b, err := json.Marshal(member)
	if err != nil {
		return
	}
	doc := memberDoc{
		ID:        memberID(guildID, member.User.ID),
		GuildID:   strconv.FormatUint(uint64(guildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoMemberStore) Get(guildID, userID discord.Snowflake) (*discord.GuildMember, bool) {
	filter := bson.M{"_id": memberID(guildID, userID)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc memberDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var m discord.GuildMember
	if err := json.Unmarshal([]byte(doc.JSON), &m); err != nil {
		return nil, false
	}
	return &m, true
}

func (s *mongoMemberStore) Delete(guildID, userID discord.Snowflake) {
	id := memberID(guildID, userID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": id})
	})
}

func (s *mongoMemberStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoMemberStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildMember] {
	coll := collection.New[discord.Snowflake, *discord.GuildMember]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []memberDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var m discord.GuildMember
		if json.Unmarshal([]byte(d.JSON), &m) == nil {
			coll.Set(m.User.ID, &m)
		}
	}
	return coll
}

func (s *mongoMemberStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Role store ────────────────────────────────────────────────────────────────

type mongoRoleStore struct{ c *MongoDBCache }

func (s *mongoRoleStore) col() *mongo.Collection { return s.c.db.Collection("roles") }

func (s *mongoRoleStore) Set(guildID discord.Snowflake, role *discord.Role) {
	if role == nil {
		return
	}
	b, err := json.Marshal(role)
	if err != nil {
		return
	}
	doc := roleDoc{
		ID:        strconv.FormatUint(uint64(role.ID), 10),
		GuildID:   strconv.FormatUint(uint64(guildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoRoleStore) Get(roleID discord.Snowflake) (*discord.Role, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(roleID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc roleDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var role discord.Role
	if err := json.Unmarshal([]byte(doc.JSON), &role); err != nil {
		return nil, false
	}
	return &role, true
}

func (s *mongoRoleStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Role] {
	coll := collection.New[discord.Snowflake, *discord.Role]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []roleDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var role discord.Role
		if json.Unmarshal([]byte(d.JSON), &role) == nil {
			coll.Set(role.ID, &role)
		}
	}
	return coll
}

func (s *mongoRoleStore) Delete(roleID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(roleID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoRoleStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoRoleStore) All() *collection.Collection[discord.Snowflake, *discord.Role] {
	coll := collection.New[discord.Snowflake, *discord.Role]()
	cursor, err := s.col().Find(s.c.ctx, liveFilter(s.c.opts.TTL))
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []roleDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var role discord.Role
		if json.Unmarshal([]byte(d.JSON), &role) == nil {
			coll.Set(role.ID, &role)
		}
	}
	return coll
}

func (s *mongoRoleStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Voice state store ──────────────────────────────────────────────────────────

type mongoVoiceStateStore struct{ c *MongoDBCache }

func (s *mongoVoiceStateStore) col() *mongo.Collection { return s.c.db.Collection("voice_states") }

func (s *mongoVoiceStateStore) Set(guildID discord.Snowflake, state *discord.VoiceState) {
	if state == nil {
		return
	}
	b, err := json.Marshal(state)
	if err != nil {
		return
	}
	doc := guildUserDoc{
		ID:        guildUserID(guildID, state.UserID),
		GuildID:   strconv.FormatUint(uint64(guildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoVoiceStateStore) Get(guildID, userID discord.Snowflake) (*discord.VoiceState, bool) {
	filter := bson.M{"_id": guildUserID(guildID, userID)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildUserDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var state discord.VoiceState
	if err := json.Unmarshal([]byte(doc.JSON), &state); err != nil {
		return nil, false
	}
	return &state, true
}

func (s *mongoVoiceStateStore) Delete(guildID, userID discord.Snowflake) {
	id := guildUserID(guildID, userID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": id})
	})
}

func (s *mongoVoiceStateStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoVoiceStateStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.VoiceState] {
	coll := collection.New[discord.Snowflake, *discord.VoiceState]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildUserDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var state discord.VoiceState
		if json.Unmarshal([]byte(d.JSON), &state) == nil {
			coll.Set(state.UserID, &state)
		}
	}
	return coll
}

func (s *mongoVoiceStateStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Soundboard store ───────────────────────────────────────────────────────────

type mongoSoundboardStore struct{ c *MongoDBCache }

func (s *mongoSoundboardStore) col() *mongo.Collection { return s.c.db.Collection("soundboard_sounds") }

func (s *mongoSoundboardStore) Set(guildID discord.Snowflake, sound *discord.SoundboardSound) {
	if sound == nil {
		return
	}
	b, err := json.Marshal(sound)
	if err != nil {
		return
	}
	doc := guildEntityDoc{
		ID:        strconv.FormatUint(uint64(sound.SoundID), 10),
		GuildID:   strconv.FormatUint(uint64(guildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoSoundboardStore) Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(soundID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildEntityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var sound discord.SoundboardSound
	if err := json.Unmarshal([]byte(doc.JSON), &sound); err != nil {
		return nil, false
	}
	return &sound, true
}

func (s *mongoSoundboardStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.SoundboardSound] {
	coll := collection.New[discord.Snowflake, *discord.SoundboardSound]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildEntityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var sound discord.SoundboardSound
		if json.Unmarshal([]byte(d.JSON), &sound) == nil {
			coll.Set(sound.SoundID, &sound)
		}
	}
	return coll
}

func (s *mongoSoundboardStore) SetAll(guildID discord.Snowflake, sounds []*discord.SoundboardSound) {
	var docs []interface{}
	for _, sound := range sounds {
		if sound == nil {
			continue
		}
		b, err := json.Marshal(sound)
		if err != nil {
			continue
		}
		docs = append(docs, guildEntityDoc{
			ID:        strconv.FormatUint(uint64(sound.SoundID), 10),
			GuildID:   strconv.FormatUint(uint64(guildID), 10),
			JSON:      string(b),
			ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
		})
	}
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		setAllTx(ctx, s.col(), gID, docs)
	})
}

func (s *mongoSoundboardStore) Delete(soundID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(soundID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoSoundboardStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoSoundboardStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Scheduled event store ──────────────────────────────────────────────────────

type mongoScheduledEventStore struct{ c *MongoDBCache }

func (s *mongoScheduledEventStore) col() *mongo.Collection {
	return s.c.db.Collection("scheduled_events")
}

func (s *mongoScheduledEventStore) Set(event *discord.GuildScheduledEvent) {
	if event == nil {
		return
	}
	b, err := json.Marshal(event)
	if err != nil {
		return
	}
	doc := guildEntityDoc{
		ID:        strconv.FormatUint(uint64(event.ID), 10),
		GuildID:   strconv.FormatUint(uint64(event.GuildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoScheduledEventStore) Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(eventID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildEntityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var event discord.GuildScheduledEvent
	if err := json.Unmarshal([]byte(doc.JSON), &event); err != nil {
		return nil, false
	}
	return &event, true
}

func (s *mongoScheduledEventStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent] {
	coll := collection.New[discord.Snowflake, *discord.GuildScheduledEvent]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildEntityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var event discord.GuildScheduledEvent
		if json.Unmarshal([]byte(d.JSON), &event) == nil {
			coll.Set(event.ID, &event)
		}
	}
	return coll
}

func (s *mongoScheduledEventStore) Delete(eventID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(eventID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoScheduledEventStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoScheduledEventStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Stage instance store ───────────────────────────────────────────────────────

type mongoStageInstanceStore struct{ c *MongoDBCache }

func (s *mongoStageInstanceStore) col() *mongo.Collection {
	return s.c.db.Collection("stage_instances")
}

func (s *mongoStageInstanceStore) Set(instance *discord.StageInstance) {
	if instance == nil {
		return
	}
	b, err := json.Marshal(instance)
	if err != nil {
		return
	}
	doc := guildEntityDoc{
		ID:        strconv.FormatUint(uint64(instance.ID), 10),
		GuildID:   strconv.FormatUint(uint64(instance.GuildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoStageInstanceStore) Get(instanceID discord.Snowflake) (*discord.StageInstance, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(instanceID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildEntityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var instance discord.StageInstance
	if err := json.Unmarshal([]byte(doc.JSON), &instance); err != nil {
		return nil, false
	}
	return &instance, true
}

func (s *mongoStageInstanceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.StageInstance] {
	coll := collection.New[discord.Snowflake, *discord.StageInstance]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildEntityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var instance discord.StageInstance
		if json.Unmarshal([]byte(d.JSON), &instance) == nil {
			coll.Set(instance.ID, &instance)
		}
	}
	return coll
}

func (s *mongoStageInstanceStore) Delete(instanceID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(instanceID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoStageInstanceStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoStageInstanceStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Emoji store ────────────────────────────────────────────────────────────────

type mongoEmojiStore struct{ c *MongoDBCache }

func (s *mongoEmojiStore) col() *mongo.Collection { return s.c.db.Collection("emojis") }

func (s *mongoEmojiStore) Set(guildID discord.Snowflake, emoji *discord.Emoji) {
	if emoji == nil {
		return
	}
	b, err := json.Marshal(emoji)
	if err != nil {
		return
	}
	doc := guildEntityDoc{
		ID:        strconv.FormatUint(uint64(emoji.ID), 10),
		GuildID:   strconv.FormatUint(uint64(guildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoEmojiStore) Get(emojiID discord.Snowflake) (*discord.Emoji, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(emojiID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildEntityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var emoji discord.Emoji
	if err := json.Unmarshal([]byte(doc.JSON), &emoji); err != nil {
		return nil, false
	}
	return &emoji, true
}

func (s *mongoEmojiStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Emoji] {
	coll := collection.New[discord.Snowflake, *discord.Emoji]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildEntityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var emoji discord.Emoji
		if json.Unmarshal([]byte(d.JSON), &emoji) == nil {
			coll.Set(emoji.ID, &emoji)
		}
	}
	return coll
}

func (s *mongoEmojiStore) SetAll(guildID discord.Snowflake, emojis []*discord.Emoji) {
	var docs []interface{}
	for _, emoji := range emojis {
		if emoji == nil || emoji.ID.IsEmpty() {
			continue
		}
		b, err := json.Marshal(emoji)
		if err != nil {
			continue
		}
		docs = append(docs, guildEntityDoc{
			ID:        strconv.FormatUint(uint64(emoji.ID), 10),
			GuildID:   strconv.FormatUint(uint64(guildID), 10),
			JSON:      string(b),
			ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
		})
	}
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		setAllTx(ctx, s.col(), gID, docs)
	})
}

func (s *mongoEmojiStore) Delete(emojiID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(emojiID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoEmojiStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoEmojiStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Sticker store ──────────────────────────────────────────────────────────────

type mongoStickerStore struct{ c *MongoDBCache }

func (s *mongoStickerStore) col() *mongo.Collection { return s.c.db.Collection("stickers") }

func (s *mongoStickerStore) Set(guildID discord.Snowflake, sticker *discord.Sticker) {
	if sticker == nil {
		return
	}
	b, err := json.Marshal(sticker)
	if err != nil {
		return
	}
	doc := guildEntityDoc{
		ID:        strconv.FormatUint(uint64(sticker.ID), 10),
		GuildID:   strconv.FormatUint(uint64(guildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoStickerStore) Get(stickerID discord.Snowflake) (*discord.Sticker, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(stickerID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildEntityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var sticker discord.Sticker
	if err := json.Unmarshal([]byte(doc.JSON), &sticker); err != nil {
		return nil, false
	}
	return &sticker, true
}

func (s *mongoStickerStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Sticker] {
	coll := collection.New[discord.Snowflake, *discord.Sticker]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildEntityDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var sticker discord.Sticker
		if json.Unmarshal([]byte(d.JSON), &sticker) == nil {
			coll.Set(sticker.ID, &sticker)
		}
	}
	return coll
}

func (s *mongoStickerStore) SetAll(guildID discord.Snowflake, stickers []*discord.Sticker) {
	var docs []interface{}
	for _, sticker := range stickers {
		if sticker == nil || sticker.ID.IsEmpty() {
			continue
		}
		b, err := json.Marshal(sticker)
		if err != nil {
			continue
		}
		docs = append(docs, guildEntityDoc{
			ID:        strconv.FormatUint(uint64(sticker.ID), 10),
			GuildID:   strconv.FormatUint(uint64(guildID), 10),
			JSON:      string(b),
			ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
		})
	}
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		setAllTx(ctx, s.col(), gID, docs)
	})
}

func (s *mongoStickerStore) Delete(stickerID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(stickerID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoStickerStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoStickerStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Presence store ─────────────────────────────────────────────────────────────

type mongoPresenceStore struct{ c *MongoDBCache }

func (s *mongoPresenceStore) col() *mongo.Collection { return s.c.db.Collection("presences") }

func (s *mongoPresenceStore) Set(presence *discord.Presence) {
	if presence == nil || presence.User.ID.IsEmpty() {
		return
	}
	b, err := json.Marshal(presence)
	if err != nil {
		return
	}
	doc := guildUserDoc{
		ID:        guildUserID(presence.GuildID, presence.User.ID),
		GuildID:   strconv.FormatUint(uint64(presence.GuildID), 10),
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoPresenceStore) Get(guildID, userID discord.Snowflake) (*discord.Presence, bool) {
	filter := bson.M{"_id": guildUserID(guildID, userID)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildUserDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var presence discord.Presence
	if err := json.Unmarshal([]byte(doc.JSON), &presence); err != nil {
		return nil, false
	}
	return &presence, true
}

func (s *mongoPresenceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Presence] {
	coll := collection.New[discord.Snowflake, *discord.Presence]()
	filter := bson.M{"guild_id": strconv.FormatUint(uint64(guildID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(s.c.ctx, filter)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []guildUserDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var presence discord.Presence
		if json.Unmarshal([]byte(d.JSON), &presence) == nil {
			coll.Set(presence.User.ID, &presence)
		}
	}
	return coll
}

func (s *mongoPresenceStore) Delete(guildID, userID discord.Snowflake) {
	id := guildUserID(guildID, userID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": id})
	})
}

func (s *mongoPresenceStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoPresenceStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Message store ─────────────────────────────────────────────────────────────

type mongoMessageStore struct{ c *MongoDBCache }

func (s *mongoMessageStore) col() *mongo.Collection { return s.c.db.Collection("messages") }

func msgID(channelID, messageID discord.Snowflake) string {
	return strconv.FormatUint(uint64(channelID), 10) + ":" + strconv.FormatUint(uint64(messageID), 10)
}

func (s *mongoMessageStore) Add(msg *discord.Message) {
	if msg == nil || s.c.opts.Messages.MaxPerChannel == 0 {
		return
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	channelIDStr := strconv.FormatUint(uint64(msg.ChannelID), 10)
	max := int64(s.c.opts.Messages.MaxPerChannel)
	doc := msgDoc{
		ID:         msgID(msg.ChannelID, msg.ID),
		ChannelID:  channelIDStr,
		JSON:       string(b),
		InsertedAt: time.Now(),
		ExpiresAt:  s.c.expiresAt(s.c.opts.Messages.TTL),
	}
	s.c.enqueueWrite(func() {
		// Serialize insert+count+evict per channel to prevent TOCTOU races.
		mu, _ := s.c.msgChannelMu.LoadOrStore(channelIDStr, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()

		ctx, cancel := s.c.callCtx()
		defer cancel()

		_ = upsertByID(ctx, s.col(), doc)

		// Enforce MaxPerChannel: delete the oldest messages for this channel
		// if the count exceeds the cap.
		// Bug 32 fix: exclude expired documents from the count so TTL-expired entries
		// do not inflate the count and trigger unnecessary eviction of live messages.
		countFilter := bson.M{"channel_id": channelIDStr}
		if s.c.opts.Messages.TTL > 0 {
			countFilter["expires_at"] = bson.M{"$gt": time.Now()}
		}
		count, err := s.col().CountDocuments(ctx, countFilter)
		if err != nil || count <= max {
			return
		}
		excess := count - max
		// Find excess oldest messages (ascending by inserted_at) and delete them.
		cursor, err := s.col().Find(
			ctx,
			countFilter,
			options.Find().
				SetSort(bson.D{{Key: "inserted_at", Value: 1}}).
				SetLimit(excess).
				SetProjection(bson.D{{Key: "_id", Value: 1}}),
		)
		if err != nil {
			return
		}
		var idDocs []struct {
			ID string `bson:"_id"`
		}
		if err := cursor.All(ctx, &idDocs); err != nil {
			return
		}
		ids := make([]string, len(idDocs))
		for i, d := range idDocs {
			ids[i] = d.ID
		}
		_, _ = s.col().DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	})
}

func (s *mongoMessageStore) Get(channelID, messageID discord.Snowflake) (*discord.Message, bool) {
	filter := bson.M{"_id": msgID(channelID, messageID)}
	if s.c.opts.Messages.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc msgDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var msg discord.Message
	if err := json.Unmarshal([]byte(doc.JSON), &msg); err != nil {
		return nil, false
	}
	return &msg, true
}

func (s *mongoMessageStore) Update(msg *discord.Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	id := msgID(msg.ChannelID, msg.ID)
	expiresAt := s.c.expiresAt(s.c.opts.Messages.TTL)
	// UpdateOne with $set preserves InsertedAt and only updates the JSON payload.
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{
				"json":       string(b),
				"expires_at": expiresAt,
			}},
		)
	})
}

func (s *mongoMessageStore) Delete(channelID, messageID discord.Snowflake) {
	id := msgID(channelID, messageID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": id})
	})
}

func (s *mongoMessageStore) DeleteBulk(channelID discord.Snowflake, ids []discord.Snowflake) {
	if len(ids) == 0 {
		return
	}
	docIDs := make([]string, len(ids))
	for i, id := range ids {
		docIDs[i] = msgID(channelID, id)
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"_id": bson.M{"$in": docIDs}})
	})
}

// Channel returns messages for channelID newest-first, excluding expired entries.
func (s *mongoMessageStore) Channel(channelID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Message] {
	coll := collection.New[discord.Snowflake, *discord.Message]()
	if s.c.opts.Messages.MaxPerChannel == 0 {
		return coll
	}
	filter := bson.M{"channel_id": strconv.FormatUint(uint64(channelID), 10)}
	if s.c.opts.Messages.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	cursor, err := s.col().Find(
		s.c.ctx,
		filter,
		options.Find().
			SetSort(bson.D{{Key: "inserted_at", Value: -1}}).
			SetLimit(int64(s.c.opts.Messages.MaxPerChannel)),
	)
	if err != nil {
		return coll
	}
	defer cursor.Close(s.c.ctx)
	var docs []msgDoc
	if err := cursor.All(s.c.ctx, &docs); err != nil {
		return coll
	}
	for _, d := range docs {
		var msg discord.Message
		if json.Unmarshal([]byte(d.JSON), &msg) == nil {
			coll.Set(msg.ID, &msg)
		}
	}
	return coll
}

func (s *mongoMessageStore) DeleteChannel(channelID discord.Snowflake) {
	chIDStr := strconv.FormatUint(uint64(channelID), 10)
	s.c.enqueueWrite(func() {
		if s.c.db != nil {
			ctx, cancel := s.c.callCtx()
			defer cancel()
			_, _ = s.col().DeleteMany(ctx, bson.M{"channel_id": chIDStr})
		}
		// Remove the per-channel mutex so short-lived channels (DMs, threads) do not
		// accumulate entries in msgChannelMu indefinitely (Bug 21).
		s.c.msgChannelMu.Delete(chIDStr)
	})
}

func (s *mongoMessageStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.Messages.TTL))
	return int(n)
}
