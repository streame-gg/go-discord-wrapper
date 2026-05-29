// Package rediscache provides a Redis-backed implementation of [cache.Cache].
//
// Create a [*redis.Client] from github.com/redis/go-redis/v9, then pass it to
// [NewRedisCache] alongside the same [cache.Options] you would use for
// [cache.MemoryCache].
//
// # Key layout
//
// All keys are namespaced under a configurable prefix (default "discord"):
//
//	{prefix}:guild:{id}           — JSON-encoded Guild; TTL from Options.TTL
//	{prefix}:guild:index          — Redis SET of all cached guild IDs
//	{prefix}:channel:{id}         — JSON-encoded Channel
//	{prefix}:channel:index        — Redis SET of all cached channel IDs
//	{prefix}:user:{id}            — JSON-encoded User
//	{prefix}:user:index           — Redis SET of all cached user IDs
//	{prefix}:member:{gid}:{uid}   — JSON-encoded GuildMember
//	{prefix}:member:guild:{gid}   — Redis SET of userIDs for guild gid
//	{prefix}:role:{id}            — JSON-encoded Role
//	{prefix}:role:guild:{gid}     — Redis SET of roleIDs for guild gid
//	{prefix}:voice_state:{gid}:{uid} — JSON-encoded VoiceState
//	{prefix}:voice_state:guild:{gid} — Redis SET of userIDs for guild gid
//	{prefix}:soundboard:{id}      — JSON-encoded SoundboardSound
//	{prefix}:soundboard:guild:{gid} — Redis SET of sound IDs for guild gid
//	{prefix}:scheduled_event:{id} — JSON-encoded GuildScheduledEvent
//	{prefix}:scheduled_event:guild:{gid} — Redis SET of event IDs for guild gid
//	{prefix}:stage_instance:{id}  — JSON-encoded StageInstance
//	{prefix}:stage_instance:guild:{gid} — Redis SET of stage IDs for guild gid
//	{prefix}:emoji:{id}           — JSON-encoded Emoji
//	{prefix}:emoji:guild:{gid}    — Redis SET of emoji IDs for guild gid
//	{prefix}:sticker:{id}         — JSON-encoded Sticker
//	{prefix}:sticker:guild:{gid}  — Redis SET of sticker IDs for guild gid
//	{prefix}:presence:{gid}:{uid} — JSON-encoded Presence
//	{prefix}:presence:guild:{gid} — Redis SET of userIDs for guild gid
//	{prefix}:msg:{cid}:{mid}      — JSON-encoded Message; TTL from Options.Messages.TTL
//	{prefix}:msg:ch:{cid}         — Redis sorted set: score=insertedAt ns, member=messageID
//
// # TTL and eviction
//
// TTL is applied as Redis key expiry — no background goroutine is started.
// Index sets are lazily cleaned of stale entries on every read. Options that
// only apply to MemoryCache (EvictBehavior, UnusedWindow, Limits, OnOverflow)
// are ignored; use Redis's own maxmemory-policy for overflow eviction.
//
// The per-channel message ring is enforced by [ZPOPMIN] after every Add, which
// removes the oldest messages when the sorted set exceeds MaxPerChannel.
package rediscache

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// RedisCache implements [cache.Cache] using a Redis backend.
// All methods and stores returned by Guilds/Channels/Users/Members/Messages
// are safe for concurrent use.
//
// The underlying [*redis.Client] is not owned by RedisCache; the caller is
// responsible for connecting and closing it. Call [RedisCache.Close] on
// shutdown to release the internal context used for Redis commands.
type RedisCache struct {
	client   *redis.Client
	opts     cache.Options
	prefix   string
	ctx      context.Context
	cancel   context.CancelFunc
	stopOnce sync.Once

	// write worker: all Set/Delete calls are dispatched here so the reader
	// goroutine never blocks on Redis I/O.
	writeCh      chan func()
	writeMu      sync.Mutex
	writerClosed bool
	writerDone   chan struct{}
	syncWrites   bool // true → execute writes inline (testing only)
}

// NewRedisCache creates a RedisCache backed by client.
//
// opts.TTL is applied as Redis key expiry for all entity keys; zero means
// keys never expire. opts.Messages.MaxPerChannel caps the per-channel message
// ring (default 100). Overflow and eviction options are ignored.
func NewRedisCache(client *redis.Client, opts cache.Options) *RedisCache {
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
	c := &RedisCache{
		client:     client,
		opts:       opts,
		prefix:     "discord",
		ctx:        ctx,
		cancel:     cancel,
		writeCh:    make(chan func(), queueSize),
		writerDone: make(chan struct{}),
	}
	go c.runWriteWorker()
	return c
}

// WithKeyPrefix returns a new RedisCache that uses prefix as the key namespace
// instead of the default "discord". Use this when multiple bots or environments
// share a single Redis instance.
//
// The returned instance shares the underlying [*redis.Client] and cache options
// but has its own lifecycle: calling [Close] on the derivative does not affect
// the original, and vice versa.
func (c *RedisCache) WithKeyPrefix(prefix string) *RedisCache {
	queueSize := c.opts.WriteQueueSize
	if queueSize <= 0 {
		queueSize = 512
	}
	ctx, cancel := context.WithCancel(context.Background())
	nc := &RedisCache{
		client:     c.client,
		opts:       c.opts,
		prefix:     prefix,
		ctx:        ctx,
		cancel:     cancel,
		writeCh:    make(chan func(), queueSize),
		writerDone: make(chan struct{}),
		// stopOnce is zero-value — ready for independent Close() calls.
	}
	go nc.runWriteWorker()
	return nc
}

// k builds a Redis key by joining prefix and parts with ":".
func (c *RedisCache) k(parts ...string) string {
	s := c.prefix
	for _, p := range parts {
		s += ":" + p
	}
	return s
}

// setJSON JSON-encodes v and stores it under key with an optional TTL.
func (c *RedisCache) setJSON(key string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, b, ttl).Err()
}

// setJSONCtx is like setJSON but uses the provided context.
func (c *RedisCache) setJSONCtx(ctx context.Context, key string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, b, ttl).Err()
}

// getJSON fetches key and JSON-decodes the result into v.
// Returns (false, nil) on a cache miss.
func (c *RedisCache) getJSON(key string, v any) (bool, error) {
	b, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(b, v)
}

// EnableSyncWrites makes all write operations execute synchronously instead of
// going through the async write queue. Intended for testing only.
func (c *RedisCache) EnableSyncWrites() {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	c.syncWrites = true
}

// callCtx returns a per-call context for a write operation. When WriteTimeout
// is configured it carries a deadline that prevents a slow backend from
// blocking the write worker indefinitely.
func (c *RedisCache) callCtx() (context.Context, context.CancelFunc) {
	if c.opts.WriteTimeout > 0 {
		return context.WithTimeout(c.ctx, c.opts.WriteTimeout)
	}
	return c.ctx, func() {}
}

// enqueueWrite submits fn to the async write worker.
// In sync-writes mode (tests) fn is executed directly on the caller goroutine.
// If the buffer is full or the cache is closed the write is silently dropped.
func (c *RedisCache) enqueueWrite(fn func()) {
	c.writeMu.Lock()
	if c.syncWrites {
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
func (c *RedisCache) runWriteWorker() {
	defer close(c.writerDone)
	for fn := range c.writeCh {
		fn()
	}
}

func (c *RedisCache) Guilds() cache.GuildStore     { return &redisGuildStore{c} }
func (c *RedisCache) Channels() cache.ChannelStore { return &redisChannelStore{c} }
func (c *RedisCache) Users() cache.UserStore       { return &redisUserStore{c} }
func (c *RedisCache) Members() cache.MemberStore   { return &redisMemberStore{c} }
func (c *RedisCache) Roles() cache.RoleStore       { return &redisRoleStore{c} }
func (c *RedisCache) Messages() cache.MessageStore { return &redisMessageStore{c} }
func (c *RedisCache) VoiceStates() cache.VoiceStateStore {
	return &redisVoiceStateStore{c}
}
func (c *RedisCache) Soundboard() cache.SoundboardStore {
	return &redisSoundboardStore{c}
}
func (c *RedisCache) ScheduledEvents() cache.ScheduledEventStore {
	return &redisScheduledEventStore{c}
}
func (c *RedisCache) StageInstances() cache.StageInstanceStore {
	return &redisStageInstanceStore{c}
}
func (c *RedisCache) Emojis() cache.EmojiStore { return &redisEmojiStore{c} }
func (c *RedisCache) Stickers() cache.StickerStore {
	return &redisStickerStore{c}
}
func (c *RedisCache) Presences() cache.PresenceStore {
	return &redisPresenceStore{c}
}

// Close drains the write queue, then cancels the internal context.
// The underlying [*redis.Client] is not closed. Safe to call multiple times.
func (c *RedisCache) Close() error {
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

// ── Guild store ───────────────────────────────────────────────────────────────

type redisGuildStore struct{ c *RedisCache }

func (s *redisGuildStore) Set(guild *discord.Guild) {
	if guild == nil {
		return
	}
	key := s.c.k("guild", strconv.FormatUint(uint64(guild.ID), 10))
	idx := s.c.k("guild", "index")
	idStr := strconv.FormatUint(uint64(guild.ID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, guild, ttl)
		_ = s.c.client.SAdd(ctx, idx, idStr).Err()
	})
}

func (s *redisGuildStore) Get(id discord.Snowflake) (*discord.Guild, bool) {
	key := s.c.k("guild", strconv.FormatUint(uint64(id), 10))
	var g discord.Guild
	ok, err := s.c.getJSON(key, &g)
	if err != nil || !ok {
		if !ok && err == nil {
			// Key expired in Redis but index still holds the ID — prune it.
			_ = s.c.client.SRem(s.c.ctx, s.c.k("guild", "index"), strconv.FormatUint(uint64(id), 10)).Err()
		}
		return nil, false
	}
	return &g, true
}

func (s *redisGuildStore) Delete(id discord.Snowflake) {
	key := s.c.k("guild", strconv.FormatUint(uint64(id), 10))
	idx := s.c.k("guild", "index")
	idStr := strconv.FormatUint(uint64(id), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, idx, idStr).Err()
	})
}

func (s *redisGuildStore) All() *collection.Collection[discord.Snowflake, *discord.Guild] {
	coll := collection.New[discord.Snowflake, *discord.Guild]()
	idx := s.c.k("guild", "index")
	ids, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ids) == 0 {
		return coll
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.c.k("guild", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ids[i])
			continue
		}
		var g discord.Guild
		if json.Unmarshal([]byte(v.(string)), &g) == nil {
			coll.Set(g.ID, &g)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return coll
}

func (s *redisGuildStore) Size() int {
	n, _ := s.c.client.SCard(s.c.ctx, s.c.k("guild", "index")).Result()
	return int(n)
}

// ── Channel store ─────────────────────────────────────────────────────────────

type redisChannelStore struct{ c *RedisCache }

func (s *redisChannelStore) Set(ch *discord.Channel) {
	if ch == nil {
		return
	}
	key := s.c.k("channel", strconv.FormatUint(uint64(ch.ID), 10))
	idx := s.c.k("channel", "index")
	idStr := strconv.FormatUint(uint64(ch.ID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, ch, ttl)
		_ = s.c.client.SAdd(ctx, idx, idStr).Err()
	})
}

func (s *redisChannelStore) Get(id discord.Snowflake) (*discord.Channel, bool) {
	key := s.c.k("channel", strconv.FormatUint(uint64(id), 10))
	var ch discord.Channel
	ok, err := s.c.getJSON(key, &ch)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("channel", "index"), strconv.FormatUint(uint64(id), 10)).Err()
		}
		return nil, false
	}
	return &ch, true
}

func (s *redisChannelStore) Delete(id discord.Snowflake) {
	key := s.c.k("channel", strconv.FormatUint(uint64(id), 10))
	idx := s.c.k("channel", "index")
	idStr := strconv.FormatUint(uint64(id), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, idx, idStr).Err()
	})
}

func (s *redisChannelStore) All() *collection.Collection[discord.Snowflake, *discord.Channel] {
	coll := collection.New[discord.Snowflake, *discord.Channel]()
	idx := s.c.k("channel", "index")
	ids, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ids) == 0 {
		return coll
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.c.k("channel", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ids[i])
			continue
		}
		var ch discord.Channel
		if json.Unmarshal([]byte(v.(string)), &ch) == nil {
			coll.Set(ch.ID, &ch)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return coll
}

func (s *redisChannelStore) Size() int {
	n, _ := s.c.client.SCard(s.c.ctx, s.c.k("channel", "index")).Result()
	return int(n)
}

// ── User store ────────────────────────────────────────────────────────────────

type redisUserStore struct{ c *RedisCache }

func (s *redisUserStore) Set(user *discord.User) {
	if user == nil {
		return
	}
	key := s.c.k("user", strconv.FormatUint(uint64(user.ID), 10))
	idx := s.c.k("user", "index")
	idStr := strconv.FormatUint(uint64(user.ID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, user, ttl)
		_ = s.c.client.SAdd(ctx, idx, idStr).Err()
	})
}

func (s *redisUserStore) Get(id discord.Snowflake) (*discord.User, bool) {
	key := s.c.k("user", strconv.FormatUint(uint64(id), 10))
	var u discord.User
	ok, err := s.c.getJSON(key, &u)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("user", "index"), strconv.FormatUint(uint64(id), 10)).Err()
		}
		return nil, false
	}
	return &u, true
}

func (s *redisUserStore) Delete(id discord.Snowflake) {
	key := s.c.k("user", strconv.FormatUint(uint64(id), 10))
	idx := s.c.k("user", "index")
	idStr := strconv.FormatUint(uint64(id), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, idx, idStr).Err()
	})
}

func (s *redisUserStore) All() *collection.Collection[discord.Snowflake, *discord.User] {
	coll := collection.New[discord.Snowflake, *discord.User]()
	idx := s.c.k("user", "index")
	ids, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ids) == 0 {
		return coll
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.c.k("user", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ids[i])
			continue
		}
		var u discord.User
		if json.Unmarshal([]byte(v.(string)), &u) == nil {
			coll.Set(u.ID, &u)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return coll
}

func (s *redisUserStore) Size() int {
	n, _ := s.c.client.SCard(s.c.ctx, s.c.k("user", "index")).Result()
	return int(n)
}

// ── Member store ──────────────────────────────────────────────────────────────

type redisMemberStore struct{ c *RedisCache }

func (s *redisMemberStore) Set(guildID discord.Snowflake, member *discord.GuildMember) {
	if member == nil || member.User == nil {
		return
	}
	key := s.c.k("member", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(member.User.ID), 10))
	gIdx := s.c.k("member", "guild", strconv.FormatUint(uint64(guildID), 10))
	uidStr := strconv.FormatUint(uint64(member.User.ID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, member, ttl)
		_ = s.c.client.SAdd(ctx, gIdx, uidStr).Err()
	})
}

func (s *redisMemberStore) Get(guildID, userID discord.Snowflake) (*discord.GuildMember, bool) {
	key := s.c.k("member", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(userID), 10))
	var m discord.GuildMember
	ok, err := s.c.getJSON(key, &m)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("member", "guild", strconv.FormatUint(uint64(guildID), 10)), strconv.FormatUint(uint64(userID), 10)).Err()
		}
		return nil, false
	}
	return &m, true
}

func (s *redisMemberStore) Delete(guildID, userID discord.Snowflake) {
	key := s.c.k("member", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(userID), 10))
	gIdx := s.c.k("member", "guild", strconv.FormatUint(uint64(guildID), 10))
	uidStr := strconv.FormatUint(uint64(userID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, gIdx, uidStr).Err()
	})
}

func (s *redisMemberStore) DeleteGuild(guildID discord.Snowflake) {
	gIdx := s.c.k("member", "guild", strconv.FormatUint(uint64(guildID), 10))
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		userIDs, err := s.c.client.SMembers(ctx, gIdx).Result()
		if err == nil && len(userIDs) > 0 {
			keys := make([]string, len(userIDs))
			for i, uid := range userIDs {
				keys[i] = s.c.k("member", gIDStr, uid)
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, gIdx).Err()
	})
}

func (s *redisMemberStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildMember] {
	coll := collection.New[discord.Snowflake, *discord.GuildMember]()
	gIdx := s.c.k("member", "guild", strconv.FormatUint(uint64(guildID), 10))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err != nil || len(userIDs) == 0 {
		return coll
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("member", strconv.FormatUint(uint64(guildID), 10), uid)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, userIDs[i])
			continue
		}
		var m discord.GuildMember
		if json.Unmarshal([]byte(v.(string)), &m) == nil {
			coll.Set(m.User.ID, &m)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, gIdx, stale...).Err()
	}
	return coll
}

// Size returns the approximate total number of cached members by scanning all
// per-guild index sets.
func (s *redisMemberStore) Size() int {
	pattern := s.c.k("member", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Role store ────────────────────────────────────────────────────────────────

type redisRoleStore struct{ c *RedisCache }

func (s *redisRoleStore) Set(guildID discord.Snowflake, role *discord.Role) {
	if role == nil {
		return
	}
	key := s.c.k("role", strconv.FormatUint(uint64(role.ID), 10))
	idx := s.c.k("role", "guild", strconv.FormatUint(uint64(guildID), 10))
	mapKey := s.c.k("role", "map", strconv.FormatUint(uint64(role.ID), 10))
	roleIDStr := strconv.FormatUint(uint64(role.ID), 10)
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, role, ttl)
		_ = s.c.client.SAdd(ctx, idx, roleIDStr).Err()
		_ = s.c.client.Set(ctx, mapKey, gIDStr, ttl).Err()
	})
}

func (s *redisRoleStore) Get(roleID discord.Snowflake) (*discord.Role, bool) {
	key := s.c.k("role", strconv.FormatUint(uint64(roleID), 10))
	var role discord.Role
	ok, err := s.c.getJSON(key, &role)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("role", "map", strconv.FormatUint(uint64(roleID), 10))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("role", "guild", guildID), strconv.FormatUint(uint64(roleID), 10)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &role, true
}

func (s *redisRoleStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Role] {
	coll := collection.New[discord.Snowflake, *discord.Role]()
	idx := s.c.k("role", "guild", strconv.FormatUint(uint64(guildID), 10))
	roleIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(roleIDs) == 0 {
		return coll
	}
	keys := make([]string, len(roleIDs))
	for i, id := range roleIDs {
		keys[i] = s.c.k("role", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, roleIDs[i])
			continue
		}
		var role discord.Role
		if json.Unmarshal([]byte(v.(string)), &role) == nil {
			coll.Set(role.ID, &role)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("role", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisRoleStore) Delete(roleID discord.Snowflake) {
	mapKey := s.c.k("role", "map", strconv.FormatUint(uint64(roleID), 10))
	roleKey := s.c.k("role", strconv.FormatUint(uint64(roleID), 10))
	roleIDStr := strconv.FormatUint(uint64(roleID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		guildID, err := s.c.client.Get(ctx, mapKey).Result()
		if err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("role", "guild", guildID), roleIDStr).Err()
		}
		_ = s.c.client.Del(ctx, roleKey, mapKey).Err()
	})
}

func (s *redisRoleStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("role", "guild", strconv.FormatUint(uint64(guildID), 10))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		roleIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(roleIDs) > 0 {
			keys := make([]string, 0, len(roleIDs)*2)
			for _, id := range roleIDs {
				keys = append(keys, s.c.k("role", id), s.c.k("role", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisRoleStore) All() *collection.Collection[discord.Snowflake, *discord.Role] {
	coll := collection.New[discord.Snowflake, *discord.Role]()
	pattern := s.c.k("role", "guild", "*")
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			roleIDs, err := s.c.client.SMembers(s.c.ctx, k).Result()
			if err != nil || len(roleIDs) == 0 {
				continue
			}
			roleKeys := make([]string, len(roleIDs))
			for i, id := range roleIDs {
				roleKeys[i] = s.c.k("role", id)
			}
			vals, err := s.c.client.MGet(s.c.ctx, roleKeys...).Result()
			if err != nil {
				continue
			}
			for _, v := range vals {
				if v == nil {
					continue
				}
				var role discord.Role
				if json.Unmarshal([]byte(v.(string)), &role) == nil {
					coll.Set(role.ID, &role)
				}
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return coll
}

func (s *redisRoleStore) Size() int {
	pattern := s.c.k("role", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Voice state store ──────────────────────────────────────────────────────────

type redisVoiceStateStore struct{ c *RedisCache }

func (s *redisVoiceStateStore) Set(guildID discord.Snowflake, state *discord.VoiceState) {
	if state == nil {
		return
	}
	key := s.c.k("voice_state", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(state.UserID), 10))
	idx := s.c.k("voice_state", "guild", strconv.FormatUint(uint64(guildID), 10))
	uidStr := strconv.FormatUint(uint64(state.UserID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, state, ttl)
		_ = s.c.client.SAdd(ctx, idx, uidStr).Err()
	})
}

func (s *redisVoiceStateStore) Get(guildID, userID discord.Snowflake) (*discord.VoiceState, bool) {
	key := s.c.k("voice_state", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(userID), 10))
	var state discord.VoiceState
	ok, err := s.c.getJSON(key, &state)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("voice_state", "guild", strconv.FormatUint(uint64(guildID), 10)), strconv.FormatUint(uint64(userID), 10)).Err()
		}
		return nil, false
	}
	return &state, true
}

func (s *redisVoiceStateStore) Delete(guildID, userID discord.Snowflake) {
	key := s.c.k("voice_state", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(userID), 10))
	gIdx := s.c.k("voice_state", "guild", strconv.FormatUint(uint64(guildID), 10))
	uidStr := strconv.FormatUint(uint64(userID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, gIdx, uidStr).Err()
	})
}

func (s *redisVoiceStateStore) DeleteGuild(guildID discord.Snowflake) {
	gIdx := s.c.k("voice_state", "guild", strconv.FormatUint(uint64(guildID), 10))
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		userIDs, err := s.c.client.SMembers(ctx, gIdx).Result()
		if err == nil && len(userIDs) > 0 {
			keys := make([]string, len(userIDs))
			for i, uid := range userIDs {
				keys[i] = s.c.k("voice_state", gIDStr, uid)
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, gIdx).Err()
	})
}

func (s *redisVoiceStateStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.VoiceState] {
	coll := collection.New[discord.Snowflake, *discord.VoiceState]()
	gIdx := s.c.k("voice_state", "guild", strconv.FormatUint(uint64(guildID), 10))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err != nil || len(userIDs) == 0 {
		return coll
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("voice_state", strconv.FormatUint(uint64(guildID), 10), uid)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, userIDs[i])
			continue
		}
		var state discord.VoiceState
		if json.Unmarshal([]byte(v.(string)), &state) == nil {
			coll.Set(state.UserID, &state)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, gIdx, stale...).Err()
	}
	return coll
}

func (s *redisVoiceStateStore) Size() int {
	pattern := s.c.k("voice_state", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// setAllScript atomically replaces all guild-scoped items (emoji, sticker,
// soundboard) with a new set. It reads the current guild index inside the
// Lua transaction so no reader ever sees a partially-populated state.
//
// KEYS[1]    = guild index key  (Redis SET of item IDs)
// ARGV[1]    = item key prefix  (e.g. "discord:emoji:")
// ARGV[2]    = map  key prefix  (e.g. "discord:emoji:map:")
// ARGV[3]    = TTL in milliseconds as string; "0" = no expiry
// ARGV[4]    = guild ID string
// ARGV[5..N] = alternating item-id, item-json pairs for the new set
var setAllScript = redis.NewScript(`
local idx  = KEYS[1]
local iPfx = ARGV[1]
local mPfx = ARGV[2]
local ttl  = tonumber(ARGV[3])
local gid  = ARGV[4]

local old = redis.call('SMEMBERS', idx)
if #old > 0 then
  local toDel = {}
  for _, id in ipairs(old) do
    toDel[#toDel+1] = iPfx .. id
    toDel[#toDel+1] = mPfx .. id
  end
  redis.call('DEL', unpack(toDel))
end
redis.call('DEL', idx)

local i = 5
local n = #ARGV
while i < n do
  local id  = ARGV[i]
  local jsn = ARGV[i+1]
  if ttl > 0 then
    redis.call('SET', iPfx .. id, jsn, 'PX', ttl)
    redis.call('SET', mPfx .. id, gid, 'PX', ttl)
  else
    redis.call('SET', iPfx .. id, jsn)
    redis.call('SET', mPfx .. id, gid)
  end
  redis.call('SADD', idx, id)
  i = i + 2
end
return 1
`)

// ── Soundboard store ───────────────────────────────────────────────────────────

type redisSoundboardStore struct{ c *RedisCache }

func (s *redisSoundboardStore) Set(guildID discord.Snowflake, sound *discord.SoundboardSound) {
	if sound == nil {
		return
	}
	key := s.c.k("soundboard", strconv.FormatUint(uint64(sound.SoundID), 10))
	idx := s.c.k("soundboard", "guild", strconv.FormatUint(uint64(guildID), 10))
	mapKey := s.c.k("soundboard", "map", strconv.FormatUint(uint64(sound.SoundID), 10))
	soundIDStr := strconv.FormatUint(uint64(sound.SoundID), 10)
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, sound, ttl)
		_ = s.c.client.SAdd(ctx, idx, soundIDStr).Err()
		_ = s.c.client.Set(ctx, mapKey, gIDStr, ttl).Err()
	})
}

func (s *redisSoundboardStore) Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool) {
	key := s.c.k("soundboard", strconv.FormatUint(uint64(soundID), 10))
	var sound discord.SoundboardSound
	ok, err := s.c.getJSON(key, &sound)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("soundboard", "map", strconv.FormatUint(uint64(soundID), 10))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("soundboard", "guild", guildID), strconv.FormatUint(uint64(soundID), 10)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &sound, true
}

func (s *redisSoundboardStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.SoundboardSound] {
	coll := collection.New[discord.Snowflake, *discord.SoundboardSound]()
	idx := s.c.k("soundboard", "guild", strconv.FormatUint(uint64(guildID), 10))
	soundIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(soundIDs) == 0 {
		return coll
	}
	keys := make([]string, len(soundIDs))
	for i, id := range soundIDs {
		keys[i] = s.c.k("soundboard", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, soundIDs[i])
			continue
		}
		var sound discord.SoundboardSound
		if json.Unmarshal([]byte(v.(string)), &sound) == nil {
			coll.Set(sound.SoundID, &sound)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("soundboard", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisSoundboardStore) SetAll(guildID discord.Snowflake, sounds []*discord.SoundboardSound) {
	idx := s.c.k("soundboard", "guild", strconv.FormatUint(uint64(guildID), 10))
	iPfx := s.c.k("soundboard") + ":"
	mPfx := s.c.k("soundboard", "map") + ":"
	ttl := s.c.opts.TTL.Milliseconds()
	args := []interface{}{iPfx, mPfx, ttl, strconv.FormatUint(uint64(guildID), 10)}
	for _, sound := range sounds {
		if sound == nil {
			continue
		}
		b, err := json.Marshal(sound)
		if err != nil {
			continue
		}
		args = append(args, strconv.FormatUint(uint64(sound.SoundID), 10), string(b))
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = setAllScript.Run(ctx, s.c.client, []string{idx}, args...).Err()
	})
}

func (s *redisSoundboardStore) Delete(soundID discord.Snowflake) {
	mapKey := s.c.k("soundboard", "map", strconv.FormatUint(uint64(soundID), 10))
	soundKey := s.c.k("soundboard", strconv.FormatUint(uint64(soundID), 10))
	soundIDStr := strconv.FormatUint(uint64(soundID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		guildID, err := s.c.client.Get(ctx, mapKey).Result()
		if err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("soundboard", "guild", guildID), soundIDStr).Err()
		}
		_ = s.c.client.Del(ctx, soundKey, mapKey).Err()
	})
}

func (s *redisSoundboardStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("soundboard", "guild", strconv.FormatUint(uint64(guildID), 10))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		soundIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(soundIDs) > 0 {
			keys := make([]string, 0, len(soundIDs)*2)
			for _, id := range soundIDs {
				keys = append(keys, s.c.k("soundboard", id), s.c.k("soundboard", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisSoundboardStore) Size() int {
	pattern := s.c.k("soundboard", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Scheduled event store ──────────────────────────────────────────────────────

type redisScheduledEventStore struct{ c *RedisCache }

func (s *redisScheduledEventStore) Set(event *discord.GuildScheduledEvent) {
	if event == nil {
		return
	}
	key := s.c.k("scheduled_event", strconv.FormatUint(uint64(event.ID), 10))
	idx := s.c.k("scheduled_event", "guild", strconv.FormatUint(uint64(event.GuildID), 10))
	mapKey := s.c.k("scheduled_event", "map", strconv.FormatUint(uint64(event.ID), 10))
	eventIDStr := strconv.FormatUint(uint64(event.ID), 10)
	gIDStr := strconv.FormatUint(uint64(event.GuildID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, event, ttl)
		_ = s.c.client.SAdd(ctx, idx, eventIDStr).Err()
		_ = s.c.client.Set(ctx, mapKey, gIDStr, ttl).Err()
	})
}

func (s *redisScheduledEventStore) Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool) {
	key := s.c.k("scheduled_event", strconv.FormatUint(uint64(eventID), 10))
	var event discord.GuildScheduledEvent
	ok, err := s.c.getJSON(key, &event)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("scheduled_event", "map", strconv.FormatUint(uint64(eventID), 10))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("scheduled_event", "guild", guildID), strconv.FormatUint(uint64(eventID), 10)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &event, true
}

func (s *redisScheduledEventStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent] {
	coll := collection.New[discord.Snowflake, *discord.GuildScheduledEvent]()
	idx := s.c.k("scheduled_event", "guild", strconv.FormatUint(uint64(guildID), 10))
	eventIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(eventIDs) == 0 {
		return coll
	}
	keys := make([]string, len(eventIDs))
	for i, id := range eventIDs {
		keys[i] = s.c.k("scheduled_event", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, eventIDs[i])
			continue
		}
		var event discord.GuildScheduledEvent
		if json.Unmarshal([]byte(v.(string)), &event) == nil {
			coll.Set(event.ID, &event)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("scheduled_event", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisScheduledEventStore) Delete(eventID discord.Snowflake) {
	mapKey := s.c.k("scheduled_event", "map", strconv.FormatUint(uint64(eventID), 10))
	eventKey := s.c.k("scheduled_event", strconv.FormatUint(uint64(eventID), 10))
	eventIDStr := strconv.FormatUint(uint64(eventID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		guildID, err := s.c.client.Get(ctx, mapKey).Result()
		if err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("scheduled_event", "guild", guildID), eventIDStr).Err()
		}
		_ = s.c.client.Del(ctx, eventKey, mapKey).Err()
	})
}

func (s *redisScheduledEventStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("scheduled_event", "guild", strconv.FormatUint(uint64(guildID), 10))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		eventIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(eventIDs) > 0 {
			keys := make([]string, 0, len(eventIDs)*2)
			for _, id := range eventIDs {
				keys = append(keys, s.c.k("scheduled_event", id), s.c.k("scheduled_event", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisScheduledEventStore) Size() int {
	pattern := s.c.k("scheduled_event", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Stage instance store ───────────────────────────────────────────────────────

type redisStageInstanceStore struct{ c *RedisCache }

func (s *redisStageInstanceStore) Set(instance *discord.StageInstance) {
	if instance == nil {
		return
	}
	key := s.c.k("stage_instance", strconv.FormatUint(uint64(instance.ID), 10))
	idx := s.c.k("stage_instance", "guild", strconv.FormatUint(uint64(instance.GuildID), 10))
	mapKey := s.c.k("stage_instance", "map", strconv.FormatUint(uint64(instance.ID), 10))
	instanceIDStr := strconv.FormatUint(uint64(instance.ID), 10)
	gIDStr := strconv.FormatUint(uint64(instance.GuildID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, instance, ttl)
		_ = s.c.client.SAdd(ctx, idx, instanceIDStr).Err()
		_ = s.c.client.Set(ctx, mapKey, gIDStr, ttl).Err()
	})
}

func (s *redisStageInstanceStore) Get(instanceID discord.Snowflake) (*discord.StageInstance, bool) {
	key := s.c.k("stage_instance", strconv.FormatUint(uint64(instanceID), 10))
	var instance discord.StageInstance
	ok, err := s.c.getJSON(key, &instance)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("stage_instance", "map", strconv.FormatUint(uint64(instanceID), 10))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("stage_instance", "guild", guildID), strconv.FormatUint(uint64(instanceID), 10)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &instance, true
}

func (s *redisStageInstanceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.StageInstance] {
	coll := collection.New[discord.Snowflake, *discord.StageInstance]()
	idx := s.c.k("stage_instance", "guild", strconv.FormatUint(uint64(guildID), 10))
	instanceIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(instanceIDs) == 0 {
		return coll
	}
	keys := make([]string, len(instanceIDs))
	for i, id := range instanceIDs {
		keys[i] = s.c.k("stage_instance", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, instanceIDs[i])
			continue
		}
		var instance discord.StageInstance
		if json.Unmarshal([]byte(v.(string)), &instance) == nil {
			coll.Set(instance.ID, &instance)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("stage_instance", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisStageInstanceStore) Delete(instanceID discord.Snowflake) {
	mapKey := s.c.k("stage_instance", "map", strconv.FormatUint(uint64(instanceID), 10))
	instanceKey := s.c.k("stage_instance", strconv.FormatUint(uint64(instanceID), 10))
	instanceIDStr := strconv.FormatUint(uint64(instanceID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		guildID, err := s.c.client.Get(ctx, mapKey).Result()
		if err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("stage_instance", "guild", guildID), instanceIDStr).Err()
		}
		_ = s.c.client.Del(ctx, instanceKey, mapKey).Err()
	})
}

func (s *redisStageInstanceStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("stage_instance", "guild", strconv.FormatUint(uint64(guildID), 10))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		instanceIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(instanceIDs) > 0 {
			keys := make([]string, 0, len(instanceIDs)*2)
			for _, id := range instanceIDs {
				keys = append(keys, s.c.k("stage_instance", id), s.c.k("stage_instance", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisStageInstanceStore) Size() int {
	pattern := s.c.k("stage_instance", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Emoji store ────────────────────────────────────────────────────────────────

type redisEmojiStore struct{ c *RedisCache }

func (s *redisEmojiStore) Set(guildID discord.Snowflake, emoji *discord.Emoji) {
	if emoji == nil {
		return
	}
	key := s.c.k("emoji", strconv.FormatUint(uint64(emoji.ID), 10))
	idx := s.c.k("emoji", "guild", strconv.FormatUint(uint64(guildID), 10))
	mapKey := s.c.k("emoji", "map", strconv.FormatUint(uint64(emoji.ID), 10))
	emojiIDStr := strconv.FormatUint(uint64(emoji.ID), 10)
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, emoji, ttl)
		_ = s.c.client.SAdd(ctx, idx, emojiIDStr).Err()
		_ = s.c.client.Set(ctx, mapKey, gIDStr, ttl).Err()
	})
}

func (s *redisEmojiStore) Get(emojiID discord.Snowflake) (*discord.Emoji, bool) {
	key := s.c.k("emoji", strconv.FormatUint(uint64(emojiID), 10))
	var emoji discord.Emoji
	ok, err := s.c.getJSON(key, &emoji)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("emoji", "map", strconv.FormatUint(uint64(emojiID), 10))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("emoji", "guild", guildID), strconv.FormatUint(uint64(emojiID), 10)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &emoji, true
}

func (s *redisEmojiStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Emoji] {
	coll := collection.New[discord.Snowflake, *discord.Emoji]()
	idx := s.c.k("emoji", "guild", strconv.FormatUint(uint64(guildID), 10))
	emojiIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(emojiIDs) == 0 {
		return coll
	}
	keys := make([]string, len(emojiIDs))
	for i, id := range emojiIDs {
		keys[i] = s.c.k("emoji", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, emojiIDs[i])
			continue
		}
		var emoji discord.Emoji
		if json.Unmarshal([]byte(v.(string)), &emoji) == nil {
			coll.Set(emoji.ID, &emoji)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("emoji", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisEmojiStore) SetAll(guildID discord.Snowflake, emojis []*discord.Emoji) {
	idx := s.c.k("emoji", "guild", strconv.FormatUint(uint64(guildID), 10))
	iPfx := s.c.k("emoji") + ":"
	mPfx := s.c.k("emoji", "map") + ":"
	ttl := s.c.opts.TTL.Milliseconds()
	args := []interface{}{iPfx, mPfx, ttl, strconv.FormatUint(uint64(guildID), 10)}
	for _, emoji := range emojis {
		if emoji == nil || emoji.ID.IsValid() {
			continue
		}
		b, err := json.Marshal(emoji)
		if err != nil {
			continue
		}
		args = append(args, strconv.FormatUint(uint64(emoji.ID), 10), string(b))
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = setAllScript.Run(ctx, s.c.client, []string{idx}, args...).Err()
	})
}

func (s *redisEmojiStore) Delete(emojiID discord.Snowflake) {
	mapKey := s.c.k("emoji", "map", strconv.FormatUint(uint64(emojiID), 10))
	emojiKey := s.c.k("emoji", strconv.FormatUint(uint64(emojiID), 10))
	emojiIDStr := strconv.FormatUint(uint64(emojiID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		guildID, err := s.c.client.Get(ctx, mapKey).Result()
		if err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("emoji", "guild", guildID), emojiIDStr).Err()
		}
		_ = s.c.client.Del(ctx, emojiKey, mapKey).Err()
	})
}

func (s *redisEmojiStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("emoji", "guild", strconv.FormatUint(uint64(guildID), 10))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		emojiIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(emojiIDs) > 0 {
			keys := make([]string, 0, len(emojiIDs)*2)
			for _, id := range emojiIDs {
				keys = append(keys, s.c.k("emoji", id), s.c.k("emoji", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisEmojiStore) Size() int {
	pattern := s.c.k("emoji", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Sticker store ──────────────────────────────────────────────────────────────

type redisStickerStore struct{ c *RedisCache }

func (s *redisStickerStore) Set(guildID discord.Snowflake, sticker *discord.Sticker) {
	if sticker == nil {
		return
	}
	key := s.c.k("sticker", strconv.FormatUint(uint64(sticker.ID), 10))
	idx := s.c.k("sticker", "guild", strconv.FormatUint(uint64(guildID), 10))
	mapKey := s.c.k("sticker", "map", strconv.FormatUint(uint64(sticker.ID), 10))
	stickerIDStr := strconv.FormatUint(uint64(sticker.ID), 10)
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, sticker, ttl)
		_ = s.c.client.SAdd(ctx, idx, stickerIDStr).Err()
		_ = s.c.client.Set(ctx, mapKey, gIDStr, ttl).Err()
	})
}

func (s *redisStickerStore) Get(stickerID discord.Snowflake) (*discord.Sticker, bool) {
	key := s.c.k("sticker", strconv.FormatUint(uint64(stickerID), 10))
	var sticker discord.Sticker
	ok, err := s.c.getJSON(key, &sticker)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("sticker", "map", strconv.FormatUint(uint64(stickerID), 10))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("sticker", "guild", guildID), strconv.FormatUint(uint64(stickerID), 10)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &sticker, true
}

func (s *redisStickerStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Sticker] {
	coll := collection.New[discord.Snowflake, *discord.Sticker]()
	idx := s.c.k("sticker", "guild", strconv.FormatUint(uint64(guildID), 10))
	stickerIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(stickerIDs) == 0 {
		return coll
	}
	keys := make([]string, len(stickerIDs))
	for i, id := range stickerIDs {
		keys[i] = s.c.k("sticker", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, stickerIDs[i])
			continue
		}
		var sticker discord.Sticker
		if json.Unmarshal([]byte(v.(string)), &sticker) == nil {
			coll.Set(sticker.ID, &sticker)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("sticker", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisStickerStore) SetAll(guildID discord.Snowflake, stickers []*discord.Sticker) {
	idx := s.c.k("sticker", "guild", strconv.FormatUint(uint64(guildID), 10))
	iPfx := s.c.k("sticker") + ":"
	mPfx := s.c.k("sticker", "map") + ":"
	ttl := s.c.opts.TTL.Milliseconds()
	args := []interface{}{iPfx, mPfx, ttl, strconv.FormatUint(uint64(guildID), 10)}
	for _, sticker := range stickers {
		if sticker == nil || sticker.ID.IsEmpty() {
			continue
		}
		b, err := json.Marshal(sticker)
		if err != nil {
			continue
		}
		args = append(args, strconv.FormatUint(uint64(sticker.ID), 10), string(b))
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = setAllScript.Run(ctx, s.c.client, []string{idx}, args...).Err()
	})
}

func (s *redisStickerStore) Delete(stickerID discord.Snowflake) {
	mapKey := s.c.k("sticker", "map", strconv.FormatUint(uint64(stickerID), 10))
	stickerKey := s.c.k("sticker", strconv.FormatUint(uint64(stickerID), 10))
	stickerIDStr := strconv.FormatUint(uint64(stickerID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		guildID, err := s.c.client.Get(ctx, mapKey).Result()
		if err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("sticker", "guild", guildID), stickerIDStr).Err()
		}
		_ = s.c.client.Del(ctx, stickerKey, mapKey).Err()
	})
}

func (s *redisStickerStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("sticker", "guild", strconv.FormatUint(uint64(guildID), 10))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		stickerIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(stickerIDs) > 0 {
			keys := make([]string, 0, len(stickerIDs)*2)
			for _, id := range stickerIDs {
				keys = append(keys, s.c.k("sticker", id), s.c.k("sticker", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisStickerStore) Size() int {
	pattern := s.c.k("sticker", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Presence store ─────────────────────────────────────────────────────────────

type redisPresenceStore struct{ c *RedisCache }

func (s *redisPresenceStore) Set(presence *discord.Presence) {
	if presence == nil || presence.User.ID.IsEmpty() {
		return
	}
	key := s.c.k("presence", strconv.FormatUint(uint64(presence.GuildID), 10), strconv.FormatUint(uint64(presence.User.ID), 10))
	idx := s.c.k("presence", "guild", strconv.FormatUint(uint64(presence.GuildID), 10))
	uidStr := strconv.FormatUint(uint64(presence.User.ID), 10)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONCtx(ctx, key, presence, ttl)
		_ = s.c.client.SAdd(ctx, idx, uidStr).Err()
	})
}

func (s *redisPresenceStore) Get(guildID, userID discord.Snowflake) (*discord.Presence, bool) {
	key := s.c.k("presence", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(userID), 10))
	var presence discord.Presence
	ok, err := s.c.getJSON(key, &presence)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("presence", "guild", strconv.FormatUint(uint64(guildID), 10)), strconv.FormatUint(uint64(userID), 10)).Err()
		}
		return nil, false
	}
	return &presence, true
}

func (s *redisPresenceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Presence] {
	idx := s.c.k("presence", "guild", strconv.FormatUint(uint64(guildID), 10))
	userIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	coll := collection.New[discord.Snowflake, *discord.Presence]()
	if err != nil || len(userIDs) == 0 {
		return coll
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("presence", strconv.FormatUint(uint64(guildID), 10), uid)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, userIDs[i])
			continue
		}
		var presence discord.Presence
		if json.Unmarshal([]byte(v.(string)), &presence) == nil {
			coll.Set(presence.User.ID, &presence)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return coll
}

func (s *redisPresenceStore) Delete(guildID, userID discord.Snowflake) {
	key := s.c.k("presence", strconv.FormatUint(uint64(guildID), 10), strconv.FormatUint(uint64(userID), 10))
	gIdx := s.c.k("presence", "guild", strconv.FormatUint(uint64(guildID), 10))
	uidStr := strconv.FormatUint(uint64(userID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, gIdx, uidStr).Err()
	})
}

func (s *redisPresenceStore) DeleteGuild(guildID discord.Snowflake) {
	gIdx := s.c.k("presence", "guild", strconv.FormatUint(uint64(guildID), 10))
	gIDStr := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		userIDs, err := s.c.client.SMembers(ctx, gIdx).Result()
		if err == nil && len(userIDs) > 0 {
			keys := make([]string, len(userIDs))
			for i, uid := range userIDs {
				keys[i] = s.c.k("presence", gIDStr, uid)
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, gIdx).Err()
	})
}

func (s *redisPresenceStore) Size() int {
	pattern := s.c.k("presence", "guild", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.SCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Message store ─────────────────────────────────────────────────────────────

// msgAddScript atomically executes the four steps of Add (Bug 8):
//  1. SET the message JSON with TTL
//  2. ZADD to the per-channel sorted-set index
//  3. ZCARD to count members
//  4. ZPOPMIN + DEL to evict oldest entries when over capacity
//
// KEYS[1] = msg key, KEYS[2] = channel index key
// ARGV[1] = JSON, ARGV[2] = TTL milliseconds (0 = no expiry), ARGV[3] = score,
// ARGV[4] = message ID, ARGV[5] = max per channel, ARGV[6] = msg key prefix
var msgAddScript = redis.NewScript(`
local ttl = tonumber(ARGV[2])
if ttl > 0 then
  redis.call('SET', KEYS[1], ARGV[1], 'PX', ttl)
else
  redis.call('SET', KEYS[1], ARGV[1])
end
redis.call('ZADD', KEYS[2], ARGV[3], ARGV[4])
local count = redis.call('ZCARD', KEYS[2])
local max = tonumber(ARGV[5])
if count > max then
  local excess = redis.call('ZPOPMIN', KEYS[2], count - max)
  for i = 1, #excess, 2 do
    redis.call('DEL', ARGV[6] .. excess[i])
  end
end
return 1
`)

type redisMessageStore struct{ c *RedisCache }

func (s *redisMessageStore) Add(msg *discord.Message) {
	if msg == nil || s.c.opts.Messages.MaxPerChannel == 0 {
		return
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	msgKey := s.c.k("msg", strconv.FormatUint(uint64(msg.ChannelID), 10), strconv.FormatUint(uint64(msg.ID), 10))
	chKey := s.c.k("msg", "ch", strconv.FormatUint(uint64(msg.ChannelID), 10))
	// msg key prefix passed to Lua so it can DEL evicted message keys.
	msgPrefix := s.c.k("msg", strconv.FormatUint(uint64(msg.ChannelID), 10)) + ":"

	ttlMs := s.c.opts.Messages.TTL.Milliseconds()
	score := float64(time.Now().UnixNano())
	max := s.c.opts.Messages.MaxPerChannel
	msgIDStr := strconv.FormatUint(uint64(msg.ID), 10)

	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = msgAddScript.Run(ctx, s.c.client,
			[]string{msgKey, chKey},
			b, ttlMs, score, msgIDStr, max, msgPrefix,
		).Err()
	})
}

func (s *redisMessageStore) Get(channelID, messageID discord.Snowflake) (*discord.Message, bool) {
	key := s.c.k("msg", strconv.FormatUint(uint64(channelID), 10), strconv.FormatUint(uint64(messageID), 10))
	var msg discord.Message
	ok, err := s.c.getJSON(key, &msg)
	if err != nil || !ok {
		if !ok && err == nil {
			// Key expired — remove from sorted set index.
			_ = s.c.client.ZRem(s.c.ctx, s.c.k("msg", "ch", strconv.FormatUint(uint64(channelID), 10)), strconv.FormatUint(uint64(messageID), 10)).Err()
		}
		return nil, false
	}
	return &msg, true
}

func (s *redisMessageStore) Update(msg *discord.Message) {
	key := s.c.k("msg", strconv.FormatUint(uint64(msg.ChannelID), 10), strconv.FormatUint(uint64(msg.ID), 10))
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	msgTTL := s.c.opts.Messages.TTL
	// SetXX is atomic: writes only if the key already exists, eliminating the
	// TOCTOU window between Exists and Set where the TTL could expire (Bug 7).
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.SetXX(ctx, key, b, msgTTL).Err()
	})
}

func (s *redisMessageStore) Delete(channelID, messageID discord.Snowflake) {
	msgKey := s.c.k("msg", strconv.FormatUint(uint64(channelID), 10), strconv.FormatUint(uint64(messageID), 10))
	chKey := s.c.k("msg", "ch", strconv.FormatUint(uint64(channelID), 10))
	msgIDStr := strconv.FormatUint(uint64(messageID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, msgKey).Err()
		_ = s.c.client.ZRem(ctx, chKey, msgIDStr).Err()
	})
}

func (s *redisMessageStore) DeleteBulk(channelID discord.Snowflake, ids []discord.Snowflake) {
	if len(ids) == 0 {
		return
	}
	chKey := s.c.k("msg", "ch", strconv.FormatUint(uint64(channelID), 10))
	msgKeys := make([]string, len(ids))
	members := make([]any, len(ids))
	for i, id := range ids {
		msgKeys[i] = s.c.k("msg", strconv.FormatUint(uint64(channelID), 10), strconv.FormatUint(uint64(id), 10))
		members[i] = strconv.FormatUint(uint64(id), 10)
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, msgKeys...).Err()
		_ = s.c.client.ZRem(ctx, chKey, members...).Err()
	})
}

// Channel returns cached messages for channelID newest-first.
// Entries whose JSON keys have expired are silently pruned from the index.
func (s *redisMessageStore) Channel(channelID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Message] {
	coll := collection.New[discord.Snowflake, *discord.Message]()
	chKey := s.c.k("msg", "ch", strconv.FormatUint(uint64(channelID), 10))
	// ZREVRANGE: highest score (most recent insertedAt) first.
	members, err := s.c.client.ZRangeArgs(s.c.ctx, redis.ZRangeArgs{
		Key:   chKey,
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}).Result()
	if err != nil || len(members) == 0 {
		return coll
	}
	keys := make([]string, len(members))
	for i, m := range members {
		keys[i] = s.c.k("msg", strconv.FormatUint(uint64(channelID), 10), m)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, members[i])
			continue
		}
		var msg discord.Message
		if json.Unmarshal([]byte(v.(string)), &msg) == nil {
			coll.Set(msg.ID, &msg)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.ZRem(s.c.ctx, chKey, stale...).Err()
	}
	return coll
}

func (s *redisMessageStore) DeleteChannel(channelID discord.Snowflake) {
	chKey := s.c.k("msg", "ch", strconv.FormatUint(uint64(channelID), 10))
	chIDStr := strconv.FormatUint(uint64(channelID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		members, err := s.c.client.ZRange(ctx, chKey, 0, -1).Result()
		if err == nil && len(members) > 0 {
			msgKeys := make([]string, len(members))
			for i, m := range members {
				msgKeys[i] = s.c.k("msg", chIDStr, m)
			}
			_ = s.c.client.Del(ctx, msgKeys...).Err()
		}
		_ = s.c.client.Del(ctx, chKey).Err()
	})
}

// Size returns the total number of cached messages by scanning all per-channel
// sorted sets.
func (s *redisMessageStore) Size() int {
	pattern := s.c.k("msg", "ch", "*")
	var total int
	var cursor uint64
	for {
		keys, next, err := s.c.client.Scan(s.c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := s.c.client.ZCard(s.c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}
