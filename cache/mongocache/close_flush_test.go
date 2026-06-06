package mongocache_test

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/streame-gg/go-discord-wrapper/cache"
)

// TestClose_FlushesQueuedWrites is a regression test: Close must drain the async
// write queue while the context is still live. Previously Close cancelled the
// context first, so the drained write executed with an already-cancelled
// callCtx and the Mongo upsert was silently dropped.
func (s *mongoStoresSuite) TestClose_FlushesQueuedWrites() {
	// asyncCache leaves the write worker enabled, so the Set is queued.
	c := s.asyncCache(cache.Options{})
	db := s.db // capture before Close (asyncCache stores it on the suite)

	c.Guilds().Set(guild("1"))

	s.Require().NoError(c.Close(), "Close drains then cancels")

	// Verify directly against the backing collection that the write landed.
	var doc bson.M
	err := db.Collection("guilds").FindOne(context.Background(), bson.M{"_id": "1"}).Decode(&doc)
	s.Require().NoError(err, "queued write must be flushed during Close, not dropped")
	s.Contains(doc["json"], "guild-1")
}
