package mongocache

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ── Ban store ─────────────────────────────────────────────────────────────────
//
// Bans are keyed by "{guildID}:{userID}" (same composite scheme as members) and
// carry a guild_id field for per-guild queries and bulk deletion.

type mongoBanStore struct{ c *MongoDBCache }

func (s *mongoBanStore) col() *mongo.Collection { return s.c.db.Collection("bans") }

func (s *mongoBanStore) Set(guildID discord.Snowflake, ban *discord.Ban) {
	if ban == nil {
		return
	}
	b, err := json.Marshal(ban)
	if err != nil {
		return
	}
	doc := memberDoc{
		ID:        memberID(guildID, ban.User.ID),
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

func (s *mongoBanStore) Get(guildID, userID discord.Snowflake) (*discord.Ban, bool) {
	filter := bson.M{"_id": memberID(guildID, userID)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc memberDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var ban discord.Ban
	if err := json.Unmarshal([]byte(doc.JSON), &ban); err != nil {
		return nil, false
	}
	return &ban, true
}

func (s *mongoBanStore) Delete(guildID, userID discord.Snowflake) {
	id := memberID(guildID, userID)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": id})
	})
}

func (s *mongoBanStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoBanStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Ban] {
	coll := collection.New[discord.Snowflake, *discord.Ban]()
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
		var ban discord.Ban
		if json.Unmarshal([]byte(d.JSON), &ban) == nil {
			coll.Set(ban.User.ID, &ban)
		}
	}
	return coll
}

func (s *mongoBanStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── AutoModeration rule store ───────────────────────────────────────────────
//
// Rules are keyed by rule ID and carry a guild_id field for per-guild queries
// and bulk deletion (same scheme as roles).

type mongoAutoModStore struct{ c *MongoDBCache }

func (s *mongoAutoModStore) col() *mongo.Collection { return s.c.db.Collection("automod_rules") }

func (s *mongoAutoModStore) Set(guildID discord.Snowflake, rule *discord.AutoModerationRule) {
	if rule == nil {
		return
	}
	b, err := json.Marshal(rule)
	if err != nil {
		return
	}
	doc := roleDoc{
		ID:        strconv.FormatUint(uint64(rule.ID), 10),
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

func (s *mongoAutoModStore) Get(ruleID discord.Snowflake) (*discord.AutoModerationRule, bool) {
	filter := bson.M{"_id": strconv.FormatUint(uint64(ruleID), 10)}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc roleDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var rule discord.AutoModerationRule
	if err := json.Unmarshal([]byte(doc.JSON), &rule); err != nil {
		return nil, false
	}
	return &rule, true
}

func (s *mongoAutoModStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.AutoModerationRule] {
	coll := collection.New[discord.Snowflake, *discord.AutoModerationRule]()
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
		var rule discord.AutoModerationRule
		if json.Unmarshal([]byte(d.JSON), &rule) == nil {
			coll.Set(rule.ID, &rule)
		}
	}
	return coll
}

func (s *mongoAutoModStore) Delete(ruleID discord.Snowflake) {
	idStr := strconv.FormatUint(uint64(ruleID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": idStr})
	})
}

func (s *mongoAutoModStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoAutoModStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── Invite store ──────────────────────────────────────────────────────────────
//
// Invites are keyed by their string code. The guild_id field is populated from
// invite.Guild when present (Set) or from the explicit argument (SetWithGuild),
// and is the empty string for guild-less invites — GetByGuild/DeleteGuild simply
// will not match those.

type mongoInviteStore struct{ c *MongoDBCache }

func (s *mongoInviteStore) col() *mongo.Collection { return s.c.db.Collection("invites") }

func (s *mongoInviteStore) Set(invite *discord.Invite) {
	if invite == nil || invite.Code == "" {
		return
	}
	var guildID string
	if invite.Guild != nil {
		guildID = strconv.FormatUint(uint64(invite.Guild.ID), 10)
	}
	s.store(guildID, invite)
}

func (s *mongoInviteStore) SetWithGuild(guildID discord.Snowflake, invite *discord.Invite) {
	if invite == nil || invite.Code == "" {
		return
	}
	s.store(strconv.FormatUint(uint64(guildID), 10), invite)
}

func (s *mongoInviteStore) store(guildID string, invite *discord.Invite) {
	b, err := json.Marshal(invite)
	if err != nil {
		return
	}
	doc := guildEntityDoc{
		ID:        invite.Code,
		GuildID:   guildID,
		JSON:      string(b),
		ExpiresAt: s.c.expiresAt(s.c.opts.TTL),
	}
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_ = upsertByID(ctx, s.col(), doc)
	})
}

func (s *mongoInviteStore) Get(code string) (*discord.Invite, bool) {
	filter := bson.M{"_id": code}
	if s.c.opts.TTL > 0 {
		filter["expires_at"] = bson.M{"$gt": time.Now()}
	}
	var doc guildEntityDoc
	if err := s.col().FindOne(s.c.ctx, filter).Decode(&doc); err != nil {
		return nil, false
	}
	var inv discord.Invite
	if err := json.Unmarshal([]byte(doc.JSON), &inv); err != nil {
		return nil, false
	}
	return &inv, true
}

func (s *mongoInviteStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[string, *discord.Invite] {
	coll := collection.New[string, *discord.Invite]()
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
		var inv discord.Invite
		if json.Unmarshal([]byte(d.JSON), &inv) == nil {
			coll.Set(inv.Code, &inv)
		}
	}
	return coll
}

func (s *mongoInviteStore) Delete(code string) {
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteOne(ctx, bson.M{"_id": code})
	})
}

func (s *mongoInviteStore) DeleteGuild(guildID discord.Snowflake) {
	gID := strconv.FormatUint(uint64(guildID), 10)
	s.c.enqueueWrite(func() {
		ctx, cancel := s.c.callCtx()
		defer cancel()
		_, _ = s.col().DeleteMany(ctx, bson.M{"guild_id": gID})
	})
}

func (s *mongoInviteStore) Size() int {
	n, _ := s.col().CountDocuments(s.c.ctx, liveFilter(s.c.opts.TTL))
	return int(n)
}

// ── interface accessors ───────────────────────────────────────────────────────

// Compile-time guarantee that MongoDBCache satisfies the full cache.Cache interface.
var _ cache.Cache = (*MongoDBCache)(nil)

func (c *MongoDBCache) Bans() cache.BanStore                        { return &mongoBanStore{c} }
func (c *MongoDBCache) AutoModRules() cache.AutoModerationRuleStore { return &mongoAutoModStore{c} }
func (c *MongoDBCache) Invites() cache.InviteStore                  { return &mongoInviteStore{c} }
