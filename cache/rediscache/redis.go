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
	if opts.Messages.MaxPerChannel < 0 {
		opts.Messages.MaxPerChannel = 100
	}
	// MaxPerChannel == 0 means disabled (no messages cached) — leave as-is.
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

// WithKeyPrefix returns a new RedisCache that uses prefix as the key namespace
// instead of the default "discord". Use this when multiple bots or environments
// share a single Redis instance.
func (c *RedisCache) WithKeyPrefix(prefix string) *RedisCache {
	return &RedisCache{
		client: c.client,
		opts:   c.opts,
		ctx:    c.ctx,
		cancel: c.cancel,
		prefix: prefix,
	}
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
		if !ok && err == nil {
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
		if !ok && err == nil {
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
		if !ok && err == nil {
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
		if !ok && err == nil {
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

// ── Role store ────────────────────────────────────────────────────────────────

type redisRoleStore struct{ c *RedisCache }

func (s *redisRoleStore) Set(guildID common.Snowflake, role *common.Role) {
	key := s.c.k("role", string(role.ID))
	idx := s.c.k("role", "guild", string(guildID))
	mapKey := s.c.k("role", "map", string(role.ID))
	_ = s.c.setJSON(key, role, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(role.ID)).Err()
	_ = s.c.client.Set(s.c.ctx, mapKey, string(guildID), s.c.opts.TTL).Err()
}

func (s *redisRoleStore) Get(roleID common.Snowflake) (*common.Role, bool) {
	key := s.c.k("role", string(roleID))
	var role common.Role
	ok, err := s.c.getJSON(key, &role)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("role", "map", string(roleID))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("role", "guild", guildID), string(roleID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &role, true
}

func (s *redisRoleStore) GetByGuild(guildID common.Snowflake) []*common.Role {
	idx := s.c.k("role", "guild", string(guildID))
	roleIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(roleIDs) == 0 {
		return nil
	}
	keys := make([]string, len(roleIDs))
	for i, id := range roleIDs {
		keys[i] = s.c.k("role", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Role, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, roleIDs[i])
			continue
		}
		var role common.Role
		if json.Unmarshal([]byte(v.(string)), &role) == nil {
			out = append(out, &role)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("role", "map", id.(string))).Err()
		}
	}
	return out
}

func (s *redisRoleStore) Delete(roleID common.Snowflake) {
	mapKey := s.c.k("role", "map", string(roleID))
	guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
	if err == nil {
		_ = s.c.client.SRem(s.c.ctx, s.c.k("role", "guild", guildID), string(roleID)).Err()
	}
	_ = s.c.client.Del(s.c.ctx, s.c.k("role", string(roleID)), mapKey).Err()
}

func (s *redisRoleStore) DeleteGuild(guildID common.Snowflake) {
	idx := s.c.k("role", "guild", string(guildID))
	roleIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err == nil && len(roleIDs) > 0 {
		keys := make([]string, 0, len(roleIDs)*2)
		for _, id := range roleIDs {
			keys = append(keys, s.c.k("role", id), s.c.k("role", "map", id))
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, idx).Err()
}

func (s *redisRoleStore) All() []*common.Role {
	pattern := s.c.k("role", "guild", "*")
	var out []*common.Role
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
				var role common.Role
				if json.Unmarshal([]byte(v.(string)), &role) == nil {
					out = append(out, &role)
				}
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return out
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

func (s *redisVoiceStateStore) Set(guildID common.Snowflake, state *common.VoiceState) {
	key := s.c.k("voice_state", string(guildID), string(state.UserID))
	idx := s.c.k("voice_state", "guild", string(guildID))
	_ = s.c.setJSON(key, state, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(state.UserID)).Err()
}

func (s *redisVoiceStateStore) Get(guildID, userID common.Snowflake) (*common.VoiceState, bool) {
	key := s.c.k("voice_state", string(guildID), string(userID))
	var state common.VoiceState
	ok, err := s.c.getJSON(key, &state)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("voice_state", "guild", string(guildID)), string(userID)).Err()
		}
		return nil, false
	}
	return &state, true
}

func (s *redisVoiceStateStore) Delete(guildID, userID common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("voice_state", string(guildID), string(userID))).Err()
	_ = s.c.client.SRem(s.c.ctx, s.c.k("voice_state", "guild", string(guildID)), string(userID)).Err()
}

func (s *redisVoiceStateStore) DeleteGuild(guildID common.Snowflake) {
	gIdx := s.c.k("voice_state", "guild", string(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err == nil && len(userIDs) > 0 {
		keys := make([]string, len(userIDs))
		for i, uid := range userIDs {
			keys[i] = s.c.k("voice_state", string(guildID), uid)
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, gIdx).Err()
}

func (s *redisVoiceStateStore) AllInGuild(guildID common.Snowflake) []*common.VoiceState {
	gIdx := s.c.k("voice_state", "guild", string(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err != nil || len(userIDs) == 0 {
		return nil
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("voice_state", string(guildID), uid)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.VoiceState, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, userIDs[i])
			continue
		}
		var state common.VoiceState
		if json.Unmarshal([]byte(v.(string)), &state) == nil {
			out = append(out, &state)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, gIdx, stale...).Err()
	}
	return out
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

func (s *redisSoundboardStore) Set(guildID common.Snowflake, sound *common.SoundboardSound) {
	key := s.c.k("soundboard", string(sound.SoundID))
	idx := s.c.k("soundboard", "guild", string(guildID))
	mapKey := s.c.k("soundboard", "map", string(sound.SoundID))
	_ = s.c.setJSON(key, sound, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(sound.SoundID)).Err()
	_ = s.c.client.Set(s.c.ctx, mapKey, string(guildID), s.c.opts.TTL).Err()
}

func (s *redisSoundboardStore) Get(soundID common.Snowflake) (*common.SoundboardSound, bool) {
	key := s.c.k("soundboard", string(soundID))
	var sound common.SoundboardSound
	ok, err := s.c.getJSON(key, &sound)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("soundboard", "map", string(soundID))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("soundboard", "guild", guildID), string(soundID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &sound, true
}

func (s *redisSoundboardStore) GetByGuild(guildID common.Snowflake) []*common.SoundboardSound {
	idx := s.c.k("soundboard", "guild", string(guildID))
	soundIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(soundIDs) == 0 {
		return nil
	}
	keys := make([]string, len(soundIDs))
	for i, id := range soundIDs {
		keys[i] = s.c.k("soundboard", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.SoundboardSound, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, soundIDs[i])
			continue
		}
		var sound common.SoundboardSound
		if json.Unmarshal([]byte(v.(string)), &sound) == nil {
			out = append(out, &sound)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("soundboard", "map", id.(string))).Err()
		}
	}
	return out
}

func (s *redisSoundboardStore) SetAll(guildID common.Snowflake, sounds []*common.SoundboardSound) {
	idx := s.c.k("soundboard", "guild", string(guildID))
	iPfx := s.c.k("soundboard") + ":"
	mPfx := s.c.k("soundboard", "map") + ":"
	ttl := s.c.opts.TTL.Milliseconds()
	args := []interface{}{iPfx, mPfx, ttl, string(guildID)}
	for _, sound := range sounds {
		if sound == nil {
			continue
		}
		b, err := json.Marshal(sound)
		if err != nil {
			continue
		}
		args = append(args, string(sound.SoundID), string(b))
	}
	_ = setAllScript.Run(s.c.ctx, s.c.client, []string{idx}, args...).Err()
}

func (s *redisSoundboardStore) Delete(soundID common.Snowflake) {
	mapKey := s.c.k("soundboard", "map", string(soundID))
	guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
	if err == nil {
		_ = s.c.client.SRem(s.c.ctx, s.c.k("soundboard", "guild", guildID), string(soundID)).Err()
	}
	_ = s.c.client.Del(s.c.ctx, s.c.k("soundboard", string(soundID)), mapKey).Err()
}

func (s *redisSoundboardStore) DeleteGuild(guildID common.Snowflake) {
	idx := s.c.k("soundboard", "guild", string(guildID))
	soundIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err == nil && len(soundIDs) > 0 {
		keys := make([]string, 0, len(soundIDs)*2)
		for _, id := range soundIDs {
			keys = append(keys, s.c.k("soundboard", id), s.c.k("soundboard", "map", id))
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, idx).Err()
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

func (s *redisScheduledEventStore) Set(event *common.GuildScheduledEvent) {
	key := s.c.k("scheduled_event", string(event.ID))
	idx := s.c.k("scheduled_event", "guild", string(event.GuildID))
	mapKey := s.c.k("scheduled_event", "map", string(event.ID))
	_ = s.c.setJSON(key, event, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(event.ID)).Err()
	_ = s.c.client.Set(s.c.ctx, mapKey, string(event.GuildID), s.c.opts.TTL).Err()
}

func (s *redisScheduledEventStore) Get(eventID common.Snowflake) (*common.GuildScheduledEvent, bool) {
	key := s.c.k("scheduled_event", string(eventID))
	var event common.GuildScheduledEvent
	ok, err := s.c.getJSON(key, &event)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("scheduled_event", "map", string(eventID))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("scheduled_event", "guild", guildID), string(eventID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &event, true
}

func (s *redisScheduledEventStore) GetByGuild(guildID common.Snowflake) []*common.GuildScheduledEvent {
	idx := s.c.k("scheduled_event", "guild", string(guildID))
	eventIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(eventIDs) == 0 {
		return nil
	}
	keys := make([]string, len(eventIDs))
	for i, id := range eventIDs {
		keys[i] = s.c.k("scheduled_event", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.GuildScheduledEvent, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, eventIDs[i])
			continue
		}
		var event common.GuildScheduledEvent
		if json.Unmarshal([]byte(v.(string)), &event) == nil {
			out = append(out, &event)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("scheduled_event", "map", id.(string))).Err()
		}
	}
	return out
}

func (s *redisScheduledEventStore) Delete(eventID common.Snowflake) {
	mapKey := s.c.k("scheduled_event", "map", string(eventID))
	guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
	if err == nil {
		_ = s.c.client.SRem(s.c.ctx, s.c.k("scheduled_event", "guild", guildID), string(eventID)).Err()
	}
	_ = s.c.client.Del(s.c.ctx, s.c.k("scheduled_event", string(eventID)), mapKey).Err()
}

func (s *redisScheduledEventStore) DeleteGuild(guildID common.Snowflake) {
	idx := s.c.k("scheduled_event", "guild", string(guildID))
	eventIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err == nil && len(eventIDs) > 0 {
		keys := make([]string, 0, len(eventIDs)*2)
		for _, id := range eventIDs {
			keys = append(keys, s.c.k("scheduled_event", id), s.c.k("scheduled_event", "map", id))
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, idx).Err()
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

func (s *redisStageInstanceStore) Set(instance *common.StageInstance) {
	key := s.c.k("stage_instance", string(instance.ID))
	idx := s.c.k("stage_instance", "guild", string(instance.GuildID))
	mapKey := s.c.k("stage_instance", "map", string(instance.ID))
	_ = s.c.setJSON(key, instance, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(instance.ID)).Err()
	_ = s.c.client.Set(s.c.ctx, mapKey, string(instance.GuildID), s.c.opts.TTL).Err()
}

func (s *redisStageInstanceStore) Get(instanceID common.Snowflake) (*common.StageInstance, bool) {
	key := s.c.k("stage_instance", string(instanceID))
	var instance common.StageInstance
	ok, err := s.c.getJSON(key, &instance)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("stage_instance", "map", string(instanceID))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("stage_instance", "guild", guildID), string(instanceID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &instance, true
}

func (s *redisStageInstanceStore) GetByGuild(guildID common.Snowflake) []*common.StageInstance {
	idx := s.c.k("stage_instance", "guild", string(guildID))
	instanceIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(instanceIDs) == 0 {
		return nil
	}
	keys := make([]string, len(instanceIDs))
	for i, id := range instanceIDs {
		keys[i] = s.c.k("stage_instance", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.StageInstance, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, instanceIDs[i])
			continue
		}
		var instance common.StageInstance
		if json.Unmarshal([]byte(v.(string)), &instance) == nil {
			out = append(out, &instance)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("stage_instance", "map", id.(string))).Err()
		}
	}
	return out
}

func (s *redisStageInstanceStore) Delete(instanceID common.Snowflake) {
	mapKey := s.c.k("stage_instance", "map", string(instanceID))
	guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
	if err == nil {
		_ = s.c.client.SRem(s.c.ctx, s.c.k("stage_instance", "guild", guildID), string(instanceID)).Err()
	}
	_ = s.c.client.Del(s.c.ctx, s.c.k("stage_instance", string(instanceID)), mapKey).Err()
}

func (s *redisStageInstanceStore) DeleteGuild(guildID common.Snowflake) {
	idx := s.c.k("stage_instance", "guild", string(guildID))
	instanceIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err == nil && len(instanceIDs) > 0 {
		keys := make([]string, 0, len(instanceIDs)*2)
		for _, id := range instanceIDs {
			keys = append(keys, s.c.k("stage_instance", id), s.c.k("stage_instance", "map", id))
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, idx).Err()
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

func (s *redisEmojiStore) Set(guildID common.Snowflake, emoji *common.Emoji) {
	key := s.c.k("emoji", string(emoji.ID))
	idx := s.c.k("emoji", "guild", string(guildID))
	mapKey := s.c.k("emoji", "map", string(emoji.ID))
	_ = s.c.setJSON(key, emoji, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(emoji.ID)).Err()
	_ = s.c.client.Set(s.c.ctx, mapKey, string(guildID), s.c.opts.TTL).Err()
}

func (s *redisEmojiStore) Get(emojiID common.Snowflake) (*common.Emoji, bool) {
	key := s.c.k("emoji", string(emojiID))
	var emoji common.Emoji
	ok, err := s.c.getJSON(key, &emoji)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("emoji", "map", string(emojiID))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("emoji", "guild", guildID), string(emojiID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &emoji, true
}

func (s *redisEmojiStore) GetByGuild(guildID common.Snowflake) []*common.Emoji {
	idx := s.c.k("emoji", "guild", string(guildID))
	emojiIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(emojiIDs) == 0 {
		return nil
	}
	keys := make([]string, len(emojiIDs))
	for i, id := range emojiIDs {
		keys[i] = s.c.k("emoji", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Emoji, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, emojiIDs[i])
			continue
		}
		var emoji common.Emoji
		if json.Unmarshal([]byte(v.(string)), &emoji) == nil {
			out = append(out, &emoji)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("emoji", "map", id.(string))).Err()
		}
	}
	return out
}

func (s *redisEmojiStore) SetAll(guildID common.Snowflake, emojis []*common.Emoji) {
	idx := s.c.k("emoji", "guild", string(guildID))
	iPfx := s.c.k("emoji") + ":"
	mPfx := s.c.k("emoji", "map") + ":"
	ttl := s.c.opts.TTL.Milliseconds()
	args := []interface{}{iPfx, mPfx, ttl, string(guildID)}
	for _, emoji := range emojis {
		if emoji == nil || emoji.ID == "" {
			continue
		}
		b, err := json.Marshal(emoji)
		if err != nil {
			continue
		}
		args = append(args, string(emoji.ID), string(b))
	}
	_ = setAllScript.Run(s.c.ctx, s.c.client, []string{idx}, args...).Err()
}

func (s *redisEmojiStore) Delete(emojiID common.Snowflake) {
	mapKey := s.c.k("emoji", "map", string(emojiID))
	guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
	if err == nil {
		_ = s.c.client.SRem(s.c.ctx, s.c.k("emoji", "guild", guildID), string(emojiID)).Err()
	}
	_ = s.c.client.Del(s.c.ctx, s.c.k("emoji", string(emojiID)), mapKey).Err()
}

func (s *redisEmojiStore) DeleteGuild(guildID common.Snowflake) {
	idx := s.c.k("emoji", "guild", string(guildID))
	emojiIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err == nil && len(emojiIDs) > 0 {
		keys := make([]string, 0, len(emojiIDs)*2)
		for _, id := range emojiIDs {
			keys = append(keys, s.c.k("emoji", id), s.c.k("emoji", "map", id))
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, idx).Err()
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

func (s *redisStickerStore) Set(guildID common.Snowflake, sticker *common.Sticker) {
	key := s.c.k("sticker", string(sticker.ID))
	idx := s.c.k("sticker", "guild", string(guildID))
	mapKey := s.c.k("sticker", "map", string(sticker.ID))
	_ = s.c.setJSON(key, sticker, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(sticker.ID)).Err()
	_ = s.c.client.Set(s.c.ctx, mapKey, string(guildID), s.c.opts.TTL).Err()
}

func (s *redisStickerStore) Get(stickerID common.Snowflake) (*common.Sticker, bool) {
	key := s.c.k("sticker", string(stickerID))
	var sticker common.Sticker
	ok, err := s.c.getJSON(key, &sticker)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("sticker", "map", string(stickerID))
			guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
			if err == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("sticker", "guild", guildID), string(stickerID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &sticker, true
}

func (s *redisStickerStore) GetByGuild(guildID common.Snowflake) []*common.Sticker {
	idx := s.c.k("sticker", "guild", string(guildID))
	stickerIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(stickerIDs) == 0 {
		return nil
	}
	keys := make([]string, len(stickerIDs))
	for i, id := range stickerIDs {
		keys[i] = s.c.k("sticker", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Sticker, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, stickerIDs[i])
			continue
		}
		var sticker common.Sticker
		if json.Unmarshal([]byte(v.(string)), &sticker) == nil {
			out = append(out, &sticker)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("sticker", "map", id.(string))).Err()
		}
	}
	return out
}

func (s *redisStickerStore) SetAll(guildID common.Snowflake, stickers []*common.Sticker) {
	idx := s.c.k("sticker", "guild", string(guildID))
	iPfx := s.c.k("sticker") + ":"
	mPfx := s.c.k("sticker", "map") + ":"
	ttl := s.c.opts.TTL.Milliseconds()
	args := []interface{}{iPfx, mPfx, ttl, string(guildID)}
	for _, sticker := range stickers {
		if sticker == nil || sticker.ID == "" {
			continue
		}
		b, err := json.Marshal(sticker)
		if err != nil {
			continue
		}
		args = append(args, string(sticker.ID), string(b))
	}
	_ = setAllScript.Run(s.c.ctx, s.c.client, []string{idx}, args...).Err()
}

func (s *redisStickerStore) Delete(stickerID common.Snowflake) {
	mapKey := s.c.k("sticker", "map", string(stickerID))
	guildID, err := s.c.client.Get(s.c.ctx, mapKey).Result()
	if err == nil {
		_ = s.c.client.SRem(s.c.ctx, s.c.k("sticker", "guild", guildID), string(stickerID)).Err()
	}
	_ = s.c.client.Del(s.c.ctx, s.c.k("sticker", string(stickerID)), mapKey).Err()
}

func (s *redisStickerStore) DeleteGuild(guildID common.Snowflake) {
	idx := s.c.k("sticker", "guild", string(guildID))
	stickerIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err == nil && len(stickerIDs) > 0 {
		keys := make([]string, 0, len(stickerIDs)*2)
		for _, id := range stickerIDs {
			keys = append(keys, s.c.k("sticker", id), s.c.k("sticker", "map", id))
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, idx).Err()
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

func (s *redisPresenceStore) Set(presence *common.Presence) {
	key := s.c.k("presence", string(presence.GuildID), string(presence.User.ID))
	idx := s.c.k("presence", "guild", string(presence.GuildID))
	_ = s.c.setJSON(key, presence, s.c.opts.TTL)
	_ = s.c.client.SAdd(s.c.ctx, idx, string(presence.User.ID)).Err()
}

func (s *redisPresenceStore) Get(guildID, userID common.Snowflake) (*common.Presence, bool) {
	key := s.c.k("presence", string(guildID), string(userID))
	var presence common.Presence
	ok, err := s.c.getJSON(key, &presence)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("presence", "guild", string(guildID)), string(userID)).Err()
		}
		return nil, false
	}
	return &presence, true
}

func (s *redisPresenceStore) GetByGuild(guildID common.Snowflake) []*common.Presence {
	idx := s.c.k("presence", "guild", string(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(userIDs) == 0 {
		return nil
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("presence", string(guildID), uid)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return nil
	}
	out := make([]*common.Presence, 0, len(vals))
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, userIDs[i])
			continue
		}
		var presence common.Presence
		if json.Unmarshal([]byte(v.(string)), &presence) == nil {
			out = append(out, &presence)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
	}
	return out
}

func (s *redisPresenceStore) Delete(guildID, userID common.Snowflake) {
	_ = s.c.client.Del(s.c.ctx, s.c.k("presence", string(guildID), string(userID))).Err()
	_ = s.c.client.SRem(s.c.ctx, s.c.k("presence", "guild", string(guildID)), string(userID)).Err()
}

func (s *redisPresenceStore) DeleteGuild(guildID common.Snowflake) {
	gIdx := s.c.k("presence", "guild", string(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err == nil && len(userIDs) > 0 {
		keys := make([]string, len(userIDs))
		for i, uid := range userIDs {
			keys[i] = s.c.k("presence", string(guildID), uid)
		}
		_ = s.c.client.Del(s.c.ctx, keys...).Err()
	}
	_ = s.c.client.Del(s.c.ctx, gIdx).Err()
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
<<<<<<< fix/cache-and-gateway-bugs
  redis.call('SET', KEYS[1], ARGV[1], 'EX', ttl)
=======
  redis.call('SET', KEYS[1], ARGV[1], 'PX', ttl)
>>>>>>> master
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

func (s *redisMessageStore) Add(msg *common.Message) {
	if s.c.opts.Messages.MaxPerChannel == 0 {
		return
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	msgKey := s.c.k("msg", string(msg.ChannelID), string(msg.ID))
	chKey := s.c.k("msg", "ch", string(msg.ChannelID))
	// msg key prefix passed to Lua so it can DEL evicted message keys.
	msgPrefix := s.c.k("msg", string(msg.ChannelID)) + ":"

	ttlMs := s.c.opts.Messages.TTL.Milliseconds()
	score := float64(time.Now().UnixNano())
	max := s.c.opts.Messages.MaxPerChannel

	_ = msgAddScript.Run(s.c.ctx, s.c.client,
		[]string{msgKey, chKey},
		b, ttlMs, score, string(msg.ID), max, msgPrefix,
	).Err()
}

func (s *redisMessageStore) Get(channelID, messageID common.Snowflake) (*common.Message, bool) {
	key := s.c.k("msg", string(channelID), string(messageID))
	var msg common.Message
	ok, err := s.c.getJSON(key, &msg)
	if err != nil || !ok {
		if !ok && err == nil {
			// Key expired — remove from sorted set index.
			_ = s.c.client.ZRem(s.c.ctx, s.c.k("msg", "ch", string(channelID)), string(messageID)).Err()
		}
		return nil, false
	}
	return &msg, true
}

func (s *redisMessageStore) Update(msg *common.Message) {
	key := s.c.k("msg", string(msg.ChannelID), string(msg.ID))
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	// SetXX is atomic: writes only if the key already exists, eliminating the
	// TOCTOU window between Exists and Set where the TTL could expire (Bug 7).
	_ = s.c.client.SetXX(s.c.ctx, key, b, s.c.opts.Messages.TTL).Err()
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
	members, err := s.c.client.ZRangeArgs(s.c.ctx, redis.ZRangeArgs{
		Key:   chKey,
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}).Result()
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
