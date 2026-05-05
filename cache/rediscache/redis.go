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
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/common"
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
}

// NewRedisCache creates a RedisCache backed by client.
//
// opts.TTL is applied as Redis key expiry for all entity keys; zero means
// keys never expire. opts.Messages.MaxPerChannel caps the per-channel message
// ring (default 100). Overflow and eviction options are ignored.
func NewRedisCache(client *redis.Client, opts cache.Options) *RedisCache {
	if opts.Messages.MaxPerChannel <= 0 {
		opts.Messages.MaxPerChannel = 100
	}
	if opts.Messages.TTL == 0 {
		opts.Messages.TTL = opts.TTL
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &RedisCache{
		client: client,
		opts:   opts,
		prefix: "discord",
		ctx:    ctx,
		cancel: cancel,
	}
}

// WithKeyPrefix returns a shallow copy of c that uses prefix as the key
// namespace instead of the default "discord". Use this when multiple bots
// or environments share a single Redis instance.
func (c *RedisCache) WithKeyPrefix(prefix string) *RedisCache {
	c2 := *c
	c2.prefix = prefix
	return &c2
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

func (c *RedisCache) Guilds() cache.GuildStore     { return &redisGuildStore{c} }
func (c *RedisCache) Channels() cache.ChannelStore { return &redisChannelStore{c} }
func (c *RedisCache) Users() cache.UserStore       { return &redisUserStore{c} }
func (c *RedisCache) Members() cache.MemberStore   { return &redisMemberStore{c} }
func (c *RedisCache) Messages() cache.MessageStore { return &redisMessageStore{c} }

// Close cancels the internal context used for all Redis commands.
// The underlying [*redis.Client] is not closed. Safe to call multiple times.
func (c *RedisCache) Close() error {
	c.stopOnce.Do(c.cancel)
	return nil
}

// ── Guild store ───────────────────────────────────────────────────────────────

type redisGuildStore struct{ c *RedisCache }

func (s *redisGuildStore) Set(guild *common.Guild) {
	key := s.c.k("guild", string(guild.ID))
	idx := s.c.k("guild", "index")
	_ = s.c.setJSON(key, guild, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(guild.ID)).Err()
}

func (s *redisGuildStore) Get(id common.Snowflake) (*common.Guild, bool) {
	key := s.c.k("guild", string(id))
	var g common.Guild
	ok, err := s.c.getJSON(key, &g)
	if err != nil || !ok {
		if !ok {
			// Key expired in Redis but index still holds the ID — prune it.
			_ = s.c.client.SRem(s.c.ctx, s.c.k("guild", "index"), string(id)).Err()
		}
		return nil, false
	}
	return &g, true
}

func (s *redisGuildStore) Delete(id common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("guild", string(id))).Err()
	_ = s.c.client.SRem(s.c.ctx, s.c.k("guild", "index"), string(id)).Err()
}

func (s *redisGuildStore) All() []*common.Guild {
	idx := s.c.k("guild", "index")
	ids, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ids) == 0 {
		return nil
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.c.k("guild", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Guild, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ids[i])
			continue
		}
		var g common.Guild
		if json.Unmarshal([]byte(v.(string)), &g) == nil {
			out = append(out, &g)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return out
}

func (s *redisGuildStore) Size() int {
	n, _ := s.c.client.SCard(s.c.ctx, s.c.k("guild", "index")).Result()
	return int(n)
}

// ── Channel store ─────────────────────────────────────────────────────────────

type redisChannelStore struct{ c *RedisCache }

func (s *redisChannelStore) Set(ch *common.Channel) {
	key := s.c.k("channel", string(ch.ID))
	idx := s.c.k("channel", "index")
	_ = s.c.setJSON(key, ch, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(ch.ID)).Err()
}

func (s *redisChannelStore) Get(id common.Snowflake) (*common.Channel, bool) {
	key := s.c.k("channel", string(id))
	var ch common.Channel
	ok, err := s.c.getJSON(key, &ch)
	if err != nil || !ok {
		if !ok {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("channel", "index"), string(id)).Err()
		}
		return nil, false
	}
	return &ch, true
}

func (s *redisChannelStore) Delete(id common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("channel", string(id))).Err()
	_ = s.c.client.SRem(s.c.ctx, s.c.k("channel", "index"), string(id)).Err()
}

func (s *redisChannelStore) All() []*common.Channel {
	idx := s.c.k("channel", "index")
	ids, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ids) == 0 {
		return nil
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.c.k("channel", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Channel, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ids[i])
			continue
		}
		var ch common.Channel
		if json.Unmarshal([]byte(v.(string)), &ch) == nil {
			out = append(out, &ch)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return out
}

func (s *redisChannelStore) Size() int {
	n, _ := s.c.client.SCard(s.c.ctx, s.c.k("channel", "index")).Result()
	return int(n)
}

// ── User store ────────────────────────────────────────────────────────────────

type redisUserStore struct{ c *RedisCache }

func (s *redisUserStore) Set(user *common.User) {
	key := s.c.k("user", string(user.ID))
	idx := s.c.k("user", "index")
	_ = s.c.setJSON(key, user, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(user.ID)).Err()
}

func (s *redisUserStore) Get(id common.Snowflake) (*common.User, bool) {
	key := s.c.k("user", string(id))
	var u common.User
	ok, err := s.c.getJSON(key, &u)
	if err != nil || !ok {
		if !ok {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("user", "index"), string(id)).Err()
		}
		return nil, false
	}
	return &u, true
}

func (s *redisUserStore) Delete(id common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("user", string(id))).Err()
	_ = s.c.client.SRem(s.c.ctx, s.c.k("user", "index"), string(id)).Err()
}

func (s *redisUserStore) All() []*common.User {
	idx := s.c.k("user", "index")
	ids, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ids) == 0 {
		return nil
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = s.c.k("user", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.User, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ids[i])
			continue
		}
		var u common.User
		if json.Unmarshal([]byte(v.(string)), &u) == nil {
			out = append(out, &u)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return out
}

func (s *redisUserStore) Size() int {
	n, _ := s.c.client.SCard(s.c.ctx, s.c.k("user", "index")).Result()
	return int(n)
}

// ── Member store ──────────────────────────────────────────────────────────────

type redisMemberStore struct{ c *RedisCache }

func (s *redisMemberStore) Set(guildID common.Snowflake, member *common.GuildMember) {
	if member.User == nil {
		return
	}
	key := s.c.k("member", string(guildID), string(member.User.ID))
	gIdx := s.c.k("member", "guild", string(guildID))
	_ = s.c.setJSON(key, member, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, gIdx, string(member.User.ID)).Err()
}

func (s *redisMemberStore) Get(guildID, userID common.Snowflake) (*common.GuildMember, bool) {
	key := s.c.k("member", string(guildID), string(userID))
	var m common.GuildMember
	ok, err := s.c.getJSON(key, &m)
	if err != nil || !ok {
		if !ok {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("member", "guild", string(guildID)), string(userID)).Err()
		}
		return nil, false
	}
	return &m, true
}

func (s *redisMemberStore) Delete(guildID, userID common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("member", string(guildID), string(userID))).Err()
	_ = s.c.client.SRem(s.c.ctx, s.c.k("member", "guild", string(guildID)), string(userID)).Err()
}

func (s *redisMemberStore) DeleteGuild(guildID common.Snowflake) {
	gIdx := s.c.k("member", "guild", string(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err == nil && len(userIDs) > 0 {
		keys := make([]string, len(userIDs))
		for i, uid := range userIDs {
			keys[i] = s.c.k("member", string(guildID), uid)
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, gIdx).Err()
}

func (s *redisMemberStore) AllInGuild(guildID common.Snowflake) []*common.GuildMember {
	gIdx := s.c.k("member", "guild", string(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err != nil || len(userIDs) == 0 {
		return nil
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("member", string(guildID), uid)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.GuildMember, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, userIDs[i])
			continue
		}
		var m common.GuildMember
		if json.Unmarshal([]byte(v.(string)), &m) == nil {
			out = append(out, &m)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, gIdx, stale...).Err()
	}
	return out
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

// ── Message store ─────────────────────────────────────────────────────────────

type redisMessageStore struct{ c *RedisCache }

func (s *redisMessageStore) Add(msg *common.Message) {
	if s.c.opts.Messages.MaxPerChannel == 0 {
		return
	}
	msgKey := s.c.k("msg", string(msg.ChannelID), string(msg.ID))
	chKey := s.c.k("msg", "ch", string(msg.ChannelID))
	_ = s.c.setJSON(msgKey, msg, s.c.opts.Messages.TTL)

	score := float64(time.Now().UnixNano())
	_ = s.c.client.ZAdd(s.c.ctx, chKey, redis.Z{Score: score, Member: string(msg.ID)}).Err()

	// Enforce MaxPerChannel: evict oldest entries when the ring is over capacity.
	max := int64(s.c.opts.Messages.MaxPerChannel)
	count, err := s.c.client.ZCard(s.c.ctx, chKey).Result()
	if err == nil && count > max {
		evicted, err := s.c.client.ZPopMin(s.c.ctx, chKey, count-max).Result()
		if err == nil {
			for _, z := range evicted {
				_ = s.c.client.Del(s.c.ctx, s.c.k("msg", string(msg.ChannelID), z.Member.(string))).Err()
			}
		}
	}
}

func (s *redisMessageStore) Get(channelID, messageID common.Snowflake) (*common.Message, bool) {
	key := s.c.k("msg", string(channelID), string(messageID))
	var msg common.Message
	ok, err := s.c.getJSON(key, &msg)
	if err != nil || !ok {
		if !ok {
			// Key expired — remove from sorted set index.
			_ = s.c.client.ZRem(s.c.ctx, s.c.k("msg", "ch", string(channelID)), string(messageID)).Err()
		}
		return nil, false
	}
	return &msg, true
}

func (s *redisMessageStore) Update(msg *common.Message) {
	key := s.c.k("msg", string(msg.ChannelID), string(msg.ID))
	// Only update if the key still exists (not expired or already deleted).
	exists, _ := s.c.client.Exists(s.c.ctx, key).Result()
	if exists == 0 {
		return
	}
	_ = s.c.setJSON(key, msg, s.c.opts.Messages.TTL)
}

func (s *redisMessageStore) Delete(channelID, messageID common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("msg", string(channelID), string(messageID))).Err()
	_ = s.c.client.ZRem(s.c.ctx, s.c.k("msg", "ch", string(channelID)), string(messageID)).Err()
}

func (s *redisMessageStore) DeleteBulk(channelID common.Snowflake, ids []common.Snowflake) {
	if len(ids) == 0 {
		return
	}
	chKey := s.c.k("msg", "ch", string(channelID))
	msgKeys := make([]string, len(ids))
	members := make([]any, len(ids))
	for i, id := range ids {
		msgKeys[i] = s.c.k("msg", string(channelID), string(id))
		members[i] = string(id)
	}
	_ = s.c.client.Del(s.c.ctx, msgKeys...).Err()
	_ = s.c.client.ZRem(s.c.ctx, chKey, members...).Err()
}

// Channel returns cached messages for channelID newest-first.
// Entries whose JSON keys have expired are silently pruned from the index.
func (s *redisMessageStore) Channel(channelID common.Snowflake) []*common.Message {
	chKey := s.c.k("msg", "ch", string(channelID))
	// ZREVRANGE: highest score (most recent insertedAt) first.
	members, err := s.c.client.ZRevRange(s.c.ctx, chKey, 0, -1).Result()
	if err != nil || len(members) == 0 {
		return nil
	}
	keys := make([]string, len(members))
	for i, m := range members {
		keys[i] = s.c.k("msg", string(channelID), m)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Message, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, members[i])
			continue
		}
		var msg common.Message
		if json.Unmarshal([]byte(v.(string)), &msg) == nil {
			out = append(out, &msg)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.ZRem(s.c.ctx, chKey, stale...).Err()
	}
	return out
}

func (s *redisMessageStore) DeleteChannel(channelID common.Snowflake) {
	chKey := s.c.k("msg", "ch", string(channelID))
	members, err := s.c.client.ZRange(s.c.ctx, chKey, 0, -1).Result()
	if err == nil && len(members) > 0 {
		msgKeys := make([]string, len(members))
		for i, m := range members {
			msgKeys[i] = s.c.k("msg", string(channelID), m)
		}
		_ = s.c.client.Del(s.c.ctx, msgKeys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, chKey).Err()
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
