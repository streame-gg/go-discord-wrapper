package rediscache

import (
	"encoding/json"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// This file implements the Ban, AutoModerationRule and Invite stores so the
// Redis backend satisfies the full cache.Cache interface. Ban mirrors the
// composite-key member store; AutoMod and Invite mirror the guild-indexed emoji
// store (per-entity key + per-guild index set + reverse map for pruning).

func sf(id discord.Snowflake) string { return strconv.FormatUint(uint64(id), 10) }

// scanIndexCardinality sums the cardinality of every set matching the given
// key pattern — the shared "approximate Size by scanning per-guild index sets"
// implementation used by the guild-scoped stores.
func (c *RedisCache) scanIndexCardinality(pattern string) int {
	var total int
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			n, _ := c.client.SCard(c.ctx, k).Result()
			total += int(n)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return total
}

// ── Ban store (composite key: guild → user) ───────────────────────────────────

type redisBanStore struct{ c *RedisCache }

func (s *redisBanStore) Set(guildID discord.Snowflake, ban *discord.Ban) {
	if ban == nil {
		return
	}
	key := s.c.k("ban", sf(guildID), sf(ban.User.ID))
	gIdx := s.c.k("ban", "guild", sf(guildID))
	uidStr := sf(ban.User.ID)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONAndIndex(ctx, key, ban, ttl, gIdx, uidStr)
	})
}

func (s *redisBanStore) Get(guildID, userID discord.Snowflake) (*discord.Ban, bool) {
	key := s.c.k("ban", sf(guildID), sf(userID))
	var b discord.Ban
	ok, err := s.c.getJSON(key, &b)
	if err != nil || !ok {
		if !ok && err == nil {
			_ = s.c.client.SRem(s.c.ctx, s.c.k("ban", "guild", sf(guildID)), sf(userID)).Err()
		}
		return nil, false
	}
	return &b, true
}

func (s *redisBanStore) Delete(guildID, userID discord.Snowflake) {
	key := s.c.k("ban", sf(guildID), sf(userID))
	gIdx := s.c.k("ban", "guild", sf(guildID))
	uidStr := sf(userID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.client.Del(ctx, key).Err()
		_ = s.c.client.SRem(ctx, gIdx, uidStr).Err()
	})
}

func (s *redisBanStore) DeleteGuild(guildID discord.Snowflake) {
	gIdx := s.c.k("ban", "guild", sf(guildID))
	gIDStr := sf(guildID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		userIDs, err := s.c.client.SMembers(ctx, gIdx).Result()
		if err == nil && len(userIDs) > 0 {
			keys := make([]string, len(userIDs))
			for i, uid := range userIDs {
				keys[i] = s.c.k("ban", gIDStr, uid)
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, gIdx).Err()
	})
}

func (s *redisBanStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, discord.Ban] {
	coll := collection.New[discord.Snowflake, discord.Ban]()
	gIdx := s.c.k("ban", "guild", sf(guildID))
	userIDs, err := s.c.client.SMembers(s.c.ctx, gIdx).Result()
	if err != nil || len(userIDs) == 0 {
		return coll
	}
	keys := make([]string, len(userIDs))
	for i, uid := range userIDs {
		keys[i] = s.c.k("ban", sf(guildID), uid)
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
		var b discord.Ban
		if json.Unmarshal([]byte(v.(string)), &b) != nil {
			continue
		}
		uid := b.User.ID
		if uid.IsEmpty() {
			if n, perr := strconv.ParseUint(userIDs[i], 10, 64); perr == nil {
				uid = discord.Snowflake(n)
			}
		}
		if uid.IsEmpty() {
			stale = append(stale, userIDs[i])
			continue
		}
		coll.Set(uid, b)
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, gIdx, stale...).Err()
	}
	return coll
}

func (s *redisBanStore) Size() int { return s.c.scanIndexCardinality(s.c.k("ban", "guild", "*")) }

// ── AutoModerationRule store (guild-indexed by rule ID) ───────────────────────

type redisAutoModStore struct{ c *RedisCache }

func (s *redisAutoModStore) Set(guildID discord.Snowflake, rule *discord.AutoModerationRule) {
	if rule == nil {
		return
	}
	key := s.c.k("automod", sf(rule.ID))
	idx := s.c.k("automod", "guild", sf(guildID))
	mapKey := s.c.k("automod", "map", sf(rule.ID))
	ruleIDStr := sf(rule.ID)
	gIDStr := sf(guildID)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONAndIndexMap(ctx, key, rule, ttl, idx, ruleIDStr, mapKey, gIDStr)
	})
}

func (s *redisAutoModStore) Get(ruleID discord.Snowflake) (*discord.AutoModerationRule, bool) {
	key := s.c.k("automod", sf(ruleID))
	var rule discord.AutoModerationRule
	ok, err := s.c.getJSON(key, &rule)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("automod", "map", sf(ruleID))
			if guildID, gerr := s.c.client.Get(s.c.ctx, mapKey).Result(); gerr == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("automod", "guild", guildID), sf(ruleID)).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &rule, true
}

func (s *redisAutoModStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, discord.AutoModerationRule] {
	coll := collection.New[discord.Snowflake, discord.AutoModerationRule]()
	idx := s.c.k("automod", "guild", sf(guildID))
	ruleIDs, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(ruleIDs) == 0 {
		return coll
	}
	keys := make([]string, len(ruleIDs))
	for i, id := range ruleIDs {
		keys[i] = s.c.k("automod", id)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, ruleIDs[i])
			continue
		}
		var rule discord.AutoModerationRule
		if json.Unmarshal([]byte(v.(string)), &rule) == nil {
			coll.Set(rule.ID, rule)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, id := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("automod", "map", id.(string))).Err()
		}
	}
	return coll
}

func (s *redisAutoModStore) Delete(ruleID discord.Snowflake) {
	mapKey := s.c.k("automod", "map", sf(ruleID))
	ruleKey := s.c.k("automod", sf(ruleID))
	ruleIDStr := sf(ruleID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		if guildID, err := s.c.client.Get(ctx, mapKey).Result(); err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("automod", "guild", guildID), ruleIDStr).Err()
		}
		_ = s.c.client.Del(ctx, ruleKey, mapKey).Err()
	})
}

func (s *redisAutoModStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("automod", "guild", sf(guildID))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		ruleIDs, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(ruleIDs) > 0 {
			keys := make([]string, 0, len(ruleIDs)*2)
			for _, id := range ruleIDs {
				keys = append(keys, s.c.k("automod", id), s.c.k("automod", "map", id))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisAutoModStore) Size() int {
	return s.c.scanIndexCardinality(s.c.k("automod", "guild", "*"))
}

// ── Invite store (guild-indexed by string code) ──────────────────────────────

type redisInviteStore struct{ c *RedisCache }

func (s *redisInviteStore) Set(invite *discord.Invite) {
	if invite == nil || invite.Code == "" {
		return
	}
	if invite.Guild != nil {
		s.store(invite.Guild.ID, invite)
		return
	}
	// No guild association: store the code key only (still retrievable by Get).
	key := s.c.k("invite", invite.Code)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		if b, err := json.Marshal(invite); err == nil {
			_ = s.c.client.Set(ctx, key, b, ttl).Err()
		}
	})
}

func (s *redisInviteStore) SetWithGuild(guildID discord.Snowflake, invite *discord.Invite) {
	if invite == nil || invite.Code == "" {
		return
	}
	s.store(guildID, invite)
}

func (s *redisInviteStore) store(guildID discord.Snowflake, invite *discord.Invite) {
	key := s.c.k("invite", invite.Code)
	idx := s.c.k("invite", "guild", sf(guildID))
	mapKey := s.c.k("invite", "map", invite.Code)
	code := invite.Code
	gIDStr := sf(guildID)
	ttl := s.c.opts.TTL
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = s.c.setJSONAndIndexMap(ctx, key, invite, ttl, idx, code, mapKey, gIDStr)
	})
}

func (s *redisInviteStore) Get(code string) (*discord.Invite, bool) {
	key := s.c.k("invite", code)
	var inv discord.Invite
	ok, err := s.c.getJSON(key, &inv)
	if err != nil || !ok {
		if !ok && err == nil {
			mapKey := s.c.k("invite", "map", code)
			if guildID, gerr := s.c.client.Get(s.c.ctx, mapKey).Result(); gerr == nil {
				_ = s.c.client.SRem(s.c.ctx, s.c.k("invite", "guild", guildID), code).Err()
			}
			_ = s.c.client.Del(s.c.ctx, mapKey).Err()
		}
		return nil, false
	}
	return &inv, true
}

func (s *redisInviteStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[string, discord.Invite] {
	coll := collection.New[string, discord.Invite]()
	idx := s.c.k("invite", "guild", sf(guildID))
	codes, err := s.c.client.SMembers(s.c.ctx, idx).Result()
	if err != nil || len(codes) == 0 {
		return coll
	}
	keys := make([]string, len(codes))
	for i, code := range codes {
		keys[i] = s.c.k("invite", code)
	}
	vals, err := s.c.client.MGet(s.c.ctx, keys...).Result()
	if err != nil {
		return coll
	}
	var stale []any
	for i, v := range vals {
		if v == nil {
			stale = append(stale, codes[i])
			continue
		}
		var inv discord.Invite
		if json.Unmarshal([]byte(v.(string)), &inv) == nil {
			coll.Set(inv.Code, inv)
		}
	}
	if len(stale) > 0 {
		_ = s.c.client.SRem(s.c.ctx, idx, stale...).Err()
		for _, code := range stale {
			_ = s.c.client.Del(s.c.ctx, s.c.k("invite", "map", code.(string))).Err()
		}
	}
	return coll
}

func (s *redisInviteStore) Delete(code string) {
	mapKey := s.c.k("invite", "map", code)
	inviteKey := s.c.k("invite", code)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		if guildID, err := s.c.client.Get(ctx, mapKey).Result(); err == nil {
			_ = s.c.client.SRem(ctx, s.c.k("invite", "guild", guildID), code).Err()
		}
		_ = s.c.client.Del(ctx, inviteKey, mapKey).Err()
	})
}

func (s *redisInviteStore) DeleteGuild(guildID discord.Snowflake) {
	idx := s.c.k("invite", "guild", sf(guildID))
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		codes, err := s.c.client.SMembers(ctx, idx).Result()
		if err == nil && len(codes) > 0 {
			keys := make([]string, 0, len(codes)*2)
			for _, code := range codes {
				keys = append(keys, s.c.k("invite", code), s.c.k("invite", "map", code))
			}
			_ = s.c.client.Del(ctx, keys...).Err()
		}
		_ = s.c.client.Del(ctx, idx).Err()
	})
}

func (s *redisInviteStore) Size() int { return s.c.scanIndexCardinality(s.c.k("invite", "guild", "*")) }

// ── interface accessors ───────────────────────────────────────────────────────

// Compile-time guarantee that RedisCache satisfies the full cache.Cache interface.
var _ cache.Cache = (*RedisCache)(nil)

func (c *RedisCache) Bans() cache.BanStore                        { return c.banStore }
func (c *RedisCache) AutoModRules() cache.AutoModerationRuleStore { return c.autoModStore }
func (c *RedisCache) Invites() cache.InviteStore                  { return c.inviteStore }
