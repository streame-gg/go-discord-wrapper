package mongocache_test

import (
	"strconv"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// These tests close the remaining reachable branches that the happy-path and
// corrupt-document suites do not reach: the setAllTx transaction failure +
// best-effort fallback, the bounded write-queue overflow drop, and the
// message-eviction backend-error guards.
//
// Intentionally NOT covered (provably unreachable via the public API, so the
// package sits just under 100% by design):
//
//   - The `if err != nil { return }` guards after json.Marshal in every store's
//     Set/SetAll (guild, user, role, stage instance, emoji, sticker, presence).
//     The corresponding discord structs contain no field encoding/json can
//     reject — no NaN-able float, no marshaled time.Time (activity timestamps
//     have a custom MarshalJSON emitting unix millis), and no poisonable
//     interface{} — so json.Marshal cannot fail for them. Reaching these guards
//     would require a test-only marshal seam in production code, which we choose
//     not to add.
//   - setAllTx's in-transaction DeleteMany error and the message-eviction
//     Find / cursor.All error returns: these need a per-operation backend fault
//     (count must succeed but the following op fail), only forceable with a
//     MongoDB failpoint, not the public API.
//   - The stale-mutex retry in message Add: a race only reachable if
//     DeleteChannel removes the per-channel mutex between LoadOrStore and the
//     ownership re-check; not deterministically forceable.

// TestSetAllTx_DuplicateIDsTriggerFallback feeds SetAll two entities sharing one
// ID. The resulting documents carry a duplicate _id, so InsertMany fails inside
// the transaction; WithTransaction returns the error and setAllTx runs its
// non-atomic delete+insert fallback. This drives the tx InsertMany-error return
// and the post-transaction fallback path.
func (s *mongoStoresSuite) TestSetAllTx_DuplicateIDsTriggerFallback() {
	c := s.cache(cache.Options{WriteTimeout: 5 * time.Second})

	dup := []*discord.Emoji{{ID: 1, Name: "a"}, {ID: 1, Name: "b"}}
	c.Emojis().SetAll(gA, dup)

	// The transaction aborts on the duplicate key; the fallback ordered insert
	// lands the first document, so exactly one emoji is stored — and crucially no
	// panic occurred.
	s.Equal(1, c.Emojis().GetByGuild(gA).Len())
}

// TestSetAllTx_EmptyDocsDeleteOnly covers the len(docs)==0 arms of both the
// transactional and fallback branches: a SetAll with only nil/empty entries
// deletes the guild's documents without inserting any.
func (s *mongoStoresSuite) TestSetAllTx_EmptyDocsDeleteOnly() {
	c := s.cache(cache.Options{})
	c.Emojis().Set(gA, &discord.Emoji{ID: 1, Name: "keep"})
	s.Require().Equal(1, c.Emojis().GetByGuild(gA).Len())

	c.Emojis().SetAll(gA, []*discord.Emoji{nil, {ID: 0}}) // nothing valid → delete-only
	s.Equal(0, c.Emojis().GetByGuild(gA).Len())
}

// TestEnqueueWrite_QueueFullIsDropped covers enqueueWrite's `default` arm. With
// the async worker doing real (network-bound) MongoDB writes and a one-slot
// queue, a tight burst of writes outpaces the worker and the surplus is dropped
// instead of blocking the producer.
func (s *mongoStoresSuite) TestEnqueueWrite_QueueFullIsDropped() {
	c := s.asyncCache(cache.Options{WriteQueueSize: 1})

	// Far more writes than the worker can drain synchronously; the bounded
	// channel overflows and the default (drop) arm runs. The only contract is
	// that producing never blocks or panics.
	for i := 0; i < 5000; i++ {
		c.Guilds().Set(guild(strconv.Itoa(i)))
	}
}

// TestMessageEviction_BackendErrorGuard drives the eviction count-error guard:
// after Close the cache context is cancelled, so the synchronous Add still runs
// but CountDocuments errors out and Add returns via its `if err != nil` guard.
func (s *mongoStoresSuite) TestMessageEviction_BackendErrorGuard() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 1}})
	c.Close() // cancels the internal context → backend commands error

	// Sync write still executes fn; upsert + CountDocuments both error, exercising
	// the eviction guard. No panic, no eviction.
	c.Messages().Add(message("m1", "c1"))
}
