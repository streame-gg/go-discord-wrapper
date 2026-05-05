package cache

import (
	"encoding/json"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Size estimation ───────────────────────────────────────────────────────────

// estimateSize returns the JSON-encoded byte length of v as a proxy for memory
// payload size. Not the exact Go heap footprint (misses pointer overhead, slice
// headers, etc.), but consistent and useful as a relative bound.
func estimateSize[V any](v V) int64 {
	b, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	return int64(len(b))
}

// ── Generic store ─────────────────────────────────────────────────────────────

type memEntry[V any] struct {
	value      V
	expiresAt  time.Time // zero = no TTL
	insertedAt time.Time
	accessedAt time.Time
	hitCount   uint64 // total successful Gets; used by ClearByFrequency
	sizeBytes  int64  // estimated JSON payload size; 0 when trackBytes=false
}

func (e *memEntry[V]) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && now.After(e.expiresAt)
}

func (e *memEntry[V]) unused(now time.Time, window time.Duration) bool {
	return window > 0 && now.Sub(e.accessedAt) > window
}

// rank returns a sort score for the entry. Lower = evict first.
func (e *memEntry[V]) rank(by ClearBy) float64 {
	switch by {
	case ClearByAge:
		return float64(e.insertedAt.UnixNano())
	case ClearByFrequency:
		return float64(e.hitCount)
	default: // ClearByLastUsed
		return float64(e.accessedAt.UnixNano())
	}
}

type storeCfg struct {
	ttl        time.Duration
	maxItems   int // per-category cap; 0 = unlimited
	clearBy    ClearBy
	trackBytes bool
}

// genericStore is a thread-safe map with TTL, sweep, and capacity enforcement.
type genericStore[K comparable, V any] struct {
	mu         sync.RWMutex
	items      map[K]*memEntry[V]
	cfg        storeCfg
	totalBytes atomic.Int64
}

func newGenericStore[K comparable, V any](cfg storeCfg) *genericStore[K, V] {
	return &genericStore[K, V]{items: make(map[K]*memEntry[V]), cfg: cfg}
}

func (s *genericStore[K, V]) set(key K, value V) {
	now := time.Now()
	var sz int64
	if s.cfg.trackBytes {
		sz = estimateSize(value)
	}
	e := &memEntry[V]{
		value:      value,
		insertedAt: now,
		accessedAt: now,
		sizeBytes:  sz,
	}
	if s.cfg.ttl > 0 {
		e.expiresAt = now.Add(s.cfg.ttl)
	}

	s.mu.Lock()
	if old, ok := s.items[key]; ok {
		s.totalBytes.Add(-old.sizeBytes)
	}
	s.items[key] = e
	n := len(s.items)
	s.mu.Unlock()

	if s.cfg.trackBytes {
		s.totalBytes.Add(sz)
	}

	// Enforce per-category cap immediately (soft: may briefly exceed by 1).
	if s.cfg.maxItems > 0 && n > s.cfg.maxItems {
		s.evictToCount(s.cfg.maxItems)
	}
}

func (s *genericStore[K, V]) get(key K) (V, bool) {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	if e.expired(now) {
		s.totalBytes.Add(-e.sizeBytes)
		delete(s.items, key)
		var zero V
		return zero, false
	}
	e.accessedAt = now
	e.hitCount++
	return e.value, true
}

func (s *genericStore[K, V]) delete(key K) {
	s.mu.Lock()
	if e, ok := s.items[key]; ok {
		s.totalBytes.Add(-e.sizeBytes)
		delete(s.items, key)
	}
	s.mu.Unlock()
}

func (s *genericStore[K, V]) all() []V {
	now := time.Now()
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]V, 0, len(s.items))
	for _, e := range s.items {
		if !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (s *genericStore[K, V]) size() int {
	s.mu.RLock()
	n := len(s.items)
	s.mu.RUnlock()
	return n
}

func (s *genericStore[K, V]) bytes() int64 { return s.totalBytes.Load() }

// sweep removes expired and (when behavior==EvictUnused) access-idle entries.
func (s *genericStore[K, V]) sweep(now time.Time, behavior EvictBehavior, unusedWindow time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, e := range s.items {
		if e.expired(now) {
			s.totalBytes.Add(-e.sizeBytes)
			delete(s.items, k)
			continue
		}
		if behavior == EvictUnused && e.unused(now, unusedWindow) {
			s.totalBytes.Add(-e.sizeBytes)
			delete(s.items, k)
		}
	}
}

// evictToCount removes entries (ranked by s.cfg.clearBy) until len(items) <= target.
func (s *genericStore[K, V]) evictToCount(target int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := len(s.items)
	if n <= target {
		return
	}
	type pair struct {
		key        K
		score      float64
		insertedAt int64 // Unix nanos; tiebreaker — older entries evicted first
		size       int64
	}
	pairs := make([]pair, 0, n)
	for k, e := range s.items {
		pairs = append(pairs, pair{k, e.rank(s.cfg.clearBy), e.insertedAt.UnixNano(), e.sizeBytes})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].score != pairs[j].score {
			return pairs[i].score < pairs[j].score
		}
		return pairs[i].insertedAt < pairs[j].insertedAt // tie: oldest inserted first
	})
	for _, p := range pairs[:n-target] {
		s.totalBytes.Add(-p.size)
		delete(s.items, p.key)
	}
}

// evictN removes up to n entries ranked lowest by the given ClearBy. Returns
// the number of bytes freed. Used by global overflow enforcement.
func (s *genericStore[K, V]) evictN(n int, by ClearBy) int64 {
	if n <= 0 {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		return 0
	}
	if n > len(s.items) {
		n = len(s.items)
	}
	type pair struct {
		key        K
		score      float64
		insertedAt int64
		size       int64
	}
	pairs := make([]pair, 0, len(s.items))
	for k, e := range s.items {
		pairs = append(pairs, pair{k, e.rank(by), e.insertedAt.UnixNano(), e.sizeBytes})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].score != pairs[j].score {
			return pairs[i].score < pairs[j].score
		}
		return pairs[i].insertedAt < pairs[j].insertedAt
	})
	var freed int64
	for _, p := range pairs[:n] {
		s.totalBytes.Add(-p.size)
		freed += p.size
		delete(s.items, p.key)
	}
	return freed
}

// ── Typed entity stores ───────────────────────────────────────────────────────

type memGuildStore struct {
	s *genericStore[common.Snowflake, *common.Guild]
}

func (g *memGuildStore) Set(guild *common.Guild)                       { g.s.set(guild.ID, guild) }
func (g *memGuildStore) Get(id common.Snowflake) (*common.Guild, bool) { return g.s.get(id) }
func (g *memGuildStore) Delete(id common.Snowflake)                    { g.s.delete(id) }
func (g *memGuildStore) All() []*common.Guild                          { return g.s.all() }
func (g *memGuildStore) Size() int                                     { return g.s.size() }

type memChannelStore struct {
	s *genericStore[common.Snowflake, *common.Channel]
}

func (c *memChannelStore) Set(ch *common.Channel)                          { c.s.set(ch.ID, ch) }
func (c *memChannelStore) Get(id common.Snowflake) (*common.Channel, bool) { return c.s.get(id) }
func (c *memChannelStore) Delete(id common.Snowflake)                      { c.s.delete(id) }
func (c *memChannelStore) All() []*common.Channel                          { return c.s.all() }
func (c *memChannelStore) Size() int                                       { return c.s.size() }

type memUserStore struct {
	s *genericStore[common.Snowflake, *common.User]
}

func (u *memUserStore) Set(user *common.User)                        { u.s.set(user.ID, user) }
func (u *memUserStore) Get(id common.Snowflake) (*common.User, bool) { return u.s.get(id) }
func (u *memUserStore) Delete(id common.Snowflake)                   { u.s.delete(id) }
func (u *memUserStore) All() []*common.User                          { return u.s.all() }
func (u *memUserStore) Size() int                                    { return u.s.size() }

// ── Member store ──────────────────────────────────────────────────────────────

type memberKey struct {
	guildID common.Snowflake
	userID  common.Snowflake
}

type memMemberStore struct {
	s *genericStore[memberKey, *common.GuildMember]
}

func (m *memMemberStore) Set(guildID common.Snowflake, member *common.GuildMember) {
	if member.User == nil {
		return
	}
	m.s.set(memberKey{guildID, member.User.ID}, member)
}

func (m *memMemberStore) Get(guildID, userID common.Snowflake) (*common.GuildMember, bool) {
	return m.s.get(memberKey{guildID, userID})
}

func (m *memMemberStore) Delete(guildID, userID common.Snowflake) {
	m.s.delete(memberKey{guildID, userID})
}

func (m *memMemberStore) DeleteGuild(guildID common.Snowflake) {
	m.s.mu.Lock()
	defer m.s.mu.Unlock()
	for k, e := range m.s.items {
		if k.guildID == guildID {
			m.s.totalBytes.Add(-e.sizeBytes)
			delete(m.s.items, k)
		}
	}
}

func (m *memMemberStore) AllInGuild(guildID common.Snowflake) []*common.GuildMember {
	now := time.Now()
	m.s.mu.RLock()
	defer m.s.mu.RUnlock()
	var out []*common.GuildMember
	for k, e := range m.s.items {
		if k.guildID == guildID && !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (m *memMemberStore) Size() int { return m.s.size() }

// ── Message store ─────────────────────────────────────────────────────────────

type msgEntry struct {
	msg        *common.Message
	insertedAt time.Time
	accessedAt time.Time
	hitCount   uint64
	sizeBytes  int64
}

func (e *msgEntry) expired(now time.Time, ttl time.Duration) bool {
	return ttl > 0 && now.Sub(e.insertedAt) > ttl
}

func (e *msgEntry) rank(by ClearBy) float64 {
	switch by {
	case ClearByAge:
		return float64(e.insertedAt.UnixNano())
	case ClearByFrequency:
		return float64(e.hitCount)
	default:
		return float64(e.accessedAt.UnixNano())
	}
}

// channelRing is a bounded FIFO queue of messages for a single channel.
type channelRing struct {
	msgs       []*msgEntry
	cap        int
	accessedAt time.Time // tracks when any message in this channel was last touched
}

func newChannelRing(cap int) *channelRing {
	return &channelRing{msgs: make([]*msgEntry, 0, cap), cap: cap, accessedAt: time.Now()}
}

func (r *channelRing) totalBytes() int64 {
	var n int64
	for _, e := range r.msgs {
		n += e.sizeBytes
	}
	return n
}

func (r *channelRing) add(msg *common.Message, trackBytes bool) (delta int64) {
	// Update existing entry if already cached.
	for _, e := range r.msgs {
		if e.msg.ID == msg.ID {
			delta -= e.sizeBytes
			e.msg = msg
			e.accessedAt = time.Now()
			if trackBytes {
				e.sizeBytes = estimateSize(msg)
				delta += e.sizeBytes
			}
			return
		}
	}
	e := &msgEntry{msg: msg, insertedAt: time.Now(), accessedAt: time.Now()}
	if trackBytes {
		e.sizeBytes = estimateSize(msg)
	}
	if len(r.msgs) >= r.cap {
		evicted := r.msgs[0]
		delta -= evicted.sizeBytes
		copy(r.msgs, r.msgs[1:])
		r.msgs[len(r.msgs)-1] = e
	} else {
		r.msgs = append(r.msgs, e)
	}
	delta += e.sizeBytes
	r.accessedAt = time.Now()
	return
}

func (r *channelRing) get(id common.Snowflake, ttl time.Duration) (*common.Message, bool) {
	now := time.Now()
	for _, e := range r.msgs {
		if e.msg.ID == id {
			if e.expired(now, ttl) {
				return nil, false
			}
			e.accessedAt = now
			e.hitCount++
			r.accessedAt = now
			return e.msg, true
		}
	}
	return nil, false
}

func (r *channelRing) update(msg *common.Message, trackBytes bool) (delta int64) {
	for _, e := range r.msgs {
		if e.msg.ID == msg.ID {
			delta -= e.sizeBytes
			e.msg = msg
			e.accessedAt = time.Now()
			if trackBytes {
				e.sizeBytes = estimateSize(msg)
				delta += e.sizeBytes
			}
			return
		}
	}
	return
}

func (r *channelRing) delete(id common.Snowflake) (delta int64) {
	for i, e := range r.msgs {
		if e.msg.ID == id {
			delta = -e.sizeBytes
			r.msgs = append(r.msgs[:i], r.msgs[i+1:]...)
			return
		}
	}
	return
}

func (r *channelRing) deleteBulk(ids []common.Snowflake) (delta int64) {
	set := make(map[common.Snowflake]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	keep := r.msgs[:0]
	for _, e := range r.msgs {
		if _, drop := set[e.msg.ID]; drop {
			delta -= e.sizeBytes
		} else {
			keep = append(keep, e)
		}
	}
	r.msgs = keep
	return
}

// all returns messages newest-first, skipping TTL-expired entries.
func (r *channelRing) all(ttl time.Duration) []*common.Message {
	now := time.Now()
	r.accessedAt = now
	out := make([]*common.Message, 0, len(r.msgs))
	for i := len(r.msgs) - 1; i >= 0; i-- {
		if !r.msgs[i].expired(now, ttl) {
			out = append(out, r.msgs[i].msg)
		}
	}
	return out
}

// sweep removes expired and access-idle message entries.
// Returns the bytes freed and whether the ring is now empty.
func (r *channelRing) sweep(now time.Time, ttl time.Duration, behavior EvictBehavior, unusedWindow time.Duration) (freed int64, empty bool) {
	keep := r.msgs[:0]
	for _, e := range r.msgs {
		drop := e.expired(now, ttl) ||
			(behavior == EvictUnused && unusedWindow > 0 && now.Sub(e.accessedAt) > unusedWindow)
		if drop {
			freed += e.sizeBytes
		} else {
			keep = append(keep, e)
		}
	}
	r.msgs = keep
	return freed, len(r.msgs) == 0
}

// evictN removes up to n messages ranked lowest by by. Returns bytes freed.
func (r *channelRing) evictN(n int, by ClearBy) (freed int64) {
	if n <= 0 || len(r.msgs) == 0 {
		return
	}
	if n >= len(r.msgs) {
		for _, e := range r.msgs {
			freed += e.sizeBytes
		}
		r.msgs = r.msgs[:0]
		return
	}
	type pair struct {
		idx   int
		score float64
	}
	pairs := make([]pair, len(r.msgs))
	for i, e := range r.msgs {
		pairs[i] = pair{i, e.rank(by)}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].score < pairs[j].score })

	drop := make(map[int]struct{}, n)
	for _, p := range pairs[:n] {
		freed += r.msgs[p.idx].sizeBytes
		drop[p.idx] = struct{}{}
	}
	keep := r.msgs[:0]
	for i, e := range r.msgs {
		if _, ok := drop[i]; !ok {
			keep = append(keep, e)
		}
	}
	r.msgs = keep
	return
}

// memMessageStore implements MessageStore.
type memMessageStore struct {
	mu         sync.RWMutex
	channels   map[common.Snowflake]*channelRing
	opts       MessageOptions
	trackBytes bool
	maxTotal   int // Limits.MaxMessages; 0 = unlimited
	clearBy    ClearBy
	totalBytes atomic.Int64
	totalMsgs  atomic.Int64
}

func newMemMessageStore(opts MessageOptions, trackBytes bool, maxTotal int, clearBy ClearBy) *memMessageStore {
	return &memMessageStore{
		channels:   make(map[common.Snowflake]*channelRing),
		opts:       opts,
		trackBytes: trackBytes,
		maxTotal:   maxTotal,
		clearBy:    clearBy,
	}
}

func (s *memMessageStore) ring(channelID common.Snowflake) *channelRing {
	s.mu.RLock()
	r := s.channels[channelID]
	s.mu.RUnlock()
	if r != nil {
		return r
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if r = s.channels[channelID]; r != nil {
		return r
	}
	r = newChannelRing(s.opts.MaxPerChannel)
	s.channels[channelID] = r
	return r
}

func (s *memMessageStore) Add(msg *common.Message) {
	if s.opts.MaxPerChannel == 0 {
		return
	}
	r := s.ring(msg.ChannelID)
	s.mu.Lock()
	prevLen := len(r.msgs)
	delta := r.add(msg, s.trackBytes)
	newLen := len(r.msgs)
	s.mu.Unlock()

	s.totalBytes.Add(delta)
	s.totalMsgs.Add(int64(newLen - prevLen))

	// Enforce total message cap.
	if s.maxTotal > 0 && int(s.totalMsgs.Load()) > s.maxTotal {
		s.evictToTotal(s.maxTotal)
	}
}

func (s *memMessageStore) Get(channelID, messageID common.Snowflake) (*common.Message, bool) {
	s.mu.RLock()
	r := s.channels[channelID]
	s.mu.RUnlock()
	if r == nil {
		return nil, false
	}
	s.mu.Lock()
	msg, ok := r.get(messageID, s.opts.TTL)
	s.mu.Unlock()
	return msg, ok
}

func (s *memMessageStore) Update(msg *common.Message) {
	s.mu.RLock()
	r := s.channels[msg.ChannelID]
	s.mu.RUnlock()
	if r == nil {
		return
	}
	s.mu.Lock()
	delta := r.update(msg, s.trackBytes)
	s.mu.Unlock()
	s.totalBytes.Add(delta)
}

func (s *memMessageStore) Delete(channelID, messageID common.Snowflake) {
	s.mu.RLock()
	r := s.channels[channelID]
	s.mu.RUnlock()
	if r == nil {
		return
	}
	s.mu.Lock()
	delta := r.delete(messageID)
	s.mu.Unlock()
	if delta != 0 {
		s.totalBytes.Add(delta)
		s.totalMsgs.Add(-1)
	}
}

func (s *memMessageStore) DeleteBulk(channelID common.Snowflake, ids []common.Snowflake) {
	s.mu.RLock()
	r := s.channels[channelID]
	s.mu.RUnlock()
	if r == nil {
		return
	}
	s.mu.Lock()
	prevLen := len(r.msgs)
	delta := r.deleteBulk(ids)
	removed := prevLen - len(r.msgs)
	s.mu.Unlock()
	s.totalBytes.Add(delta)
	s.totalMsgs.Add(-int64(removed))
}

func (s *memMessageStore) Channel(channelID common.Snowflake) []*common.Message {
	s.mu.RLock()
	r := s.channels[channelID]
	s.mu.RUnlock()
	if r == nil {
		return nil
	}
	s.mu.Lock()
	msgs := r.all(s.opts.TTL)
	s.mu.Unlock()
	return msgs
}

func (s *memMessageStore) DeleteChannel(channelID common.Snowflake) {
	s.mu.Lock()
	if r := s.channels[channelID]; r != nil {
		s.totalBytes.Add(-r.totalBytes())
		s.totalMsgs.Add(-int64(len(r.msgs)))
	}
	delete(s.channels, channelID)
	s.mu.Unlock()
}

func (s *memMessageStore) Size() int { return int(s.totalMsgs.Load()) }

func (s *memMessageStore) bytes() int64 { return s.totalBytes.Load() }

func (s *memMessageStore) sweep(now time.Time, behavior EvictBehavior, unusedWindow time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, r := range s.channels {
		freed, empty := r.sweep(now, s.opts.TTL, behavior, unusedWindow)
		s.totalBytes.Add(-freed)
		if empty && behavior == EvictUnused {
			s.totalMsgs.Add(-int64(len(r.msgs))) // already 0, but keep counter in sync
			delete(s.channels, id)
		}
	}
	// Recompute totalMsgs from scratch to stay accurate after sweep.
	var total int64
	for _, r := range s.channels {
		total += int64(len(r.msgs))
	}
	s.totalMsgs.Store(total)
}

// evictToTotal removes messages (by channel last-accessed rank) until total <= target.
// Evicts entire channel rings or individual messages depending on clearBy.
func (s *memMessageStore) evictToTotal(target int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := int(s.totalMsgs.Load())
	if current <= target {
		return
	}

	// Sort channels by their last-accessed time (ascending = evict first).
	type chanRank struct {
		id   common.Snowflake
		ring *channelRing
		at   float64
	}
	ranked := make([]chanRank, 0, len(s.channels))
	for id, r := range s.channels {
		ranked = append(ranked, chanRank{id, r, float64(r.accessedAt.UnixNano())})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].at < ranked[j].at })

	toRemove := current - target
	for _, cr := range ranked {
		if toRemove <= 0 {
			break
		}
		n := toRemove
		if n > len(cr.ring.msgs) {
			n = len(cr.ring.msgs)
		}
		freed := cr.ring.evictN(n, s.clearBy)
		s.totalBytes.Add(-freed)
		toRemove -= n
		if len(cr.ring.msgs) == 0 {
			delete(s.channels, cr.id)
		}
	}

	var total int64
	for _, r := range s.channels {
		total += int64(len(r.msgs))
	}
	s.totalMsgs.Store(total)
}

// evictN removes up to n messages globally (lowest-ranked channels first).
// Returns bytes freed. Used by MemoryCache.enforceGlobalLimits.
func (s *memMessageStore) evictN(n int, by ClearBy) (freed int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	type chanRank struct {
		id   common.Snowflake
		ring *channelRing
		at   float64
	}
	ranked := make([]chanRank, 0, len(s.channels))
	for id, r := range s.channels {
		ranked = append(ranked, chanRank{id, r, r.ring_rank(by)})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].at < ranked[j].at })

	for _, cr := range ranked {
		if n <= 0 {
			break
		}
		evict := n
		if evict > len(cr.ring.msgs) {
			evict = len(cr.ring.msgs)
		}
		f := cr.ring.evictN(evict, by)
		freed += f
		s.totalBytes.Add(-f)
		n -= evict
		if len(cr.ring.msgs) == 0 {
			delete(s.channels, cr.id)
		}
	}

	var total int64
	for _, r := range s.channels {
		total += int64(len(r.msgs))
	}
	s.totalMsgs.Store(total)
	return
}

// ring_rank returns the ranking score of a channel ring for global eviction.
func (r *channelRing) ring_rank(by ClearBy) float64 {
	if len(r.msgs) == 0 {
		return 0
	}
	switch by {
	case ClearByAge:
		// Rank by oldest insertion time in the ring.
		return float64(r.msgs[0].insertedAt.UnixNano())
	case ClearByFrequency:
		// Rank by lowest average hit count.
		var total uint64
		for _, e := range r.msgs {
			total += e.hitCount
		}
		return float64(total) / float64(len(r.msgs))
	default: // ClearByLastUsed
		return float64(r.accessedAt.UnixNano())
	}
}

// ── MemoryCache ───────────────────────────────────────────────────────────────

// MemoryCache is an in-process implementation of [Cache].
//
// Create with [NewMemoryCache] and wire it to the gateway client via
// options.WithCache. Always call [MemoryCache.Close] on shutdown.
type MemoryCache struct {
	opts     Options
	guilds   *memGuildStore
	channels *memChannelStore
	users    *memUserStore
	members  *memMemberStore
	messages *memMessageStore

	stopOnce sync.Once
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewMemoryCache creates and starts a MemoryCache. Missing options receive
// sensible defaults:
//   - SweepInterval: 30 s
//   - Messages.MaxPerChannel: 100
//   - Messages.TTL: inherits global TTL
//   - OnOverflow.Target: CategoryAll
func NewMemoryCache(opts Options) *MemoryCache {
	// Defaults.
	if opts.SweepInterval <= 0 {
		opts.SweepInterval = 30 * time.Second
	}
	if opts.Messages.MaxPerChannel < 0 {
		opts.Messages.MaxPerChannel = 0
	}
	if opts.Messages.MaxPerChannel == 0 {
		opts.Messages.MaxPerChannel = 100
	}
	if opts.Messages.TTL == 0 {
		opts.Messages.TTL = opts.TTL
	}
	if opts.OnOverflow.Target == 0 {
		opts.OnOverflow.Target = CategoryAll
	}

	trackBytes := opts.Limits.MaxSizeMB > 0
	by := opts.OnOverflow.ClearBy

	sc := func(max int) storeCfg {
		return storeCfg{ttl: opts.TTL, maxItems: max, clearBy: by, trackBytes: trackBytes}
	}

	c := &MemoryCache{
		opts:     opts,
		guilds:   &memGuildStore{s: newGenericStore[common.Snowflake, *common.Guild](sc(opts.Limits.MaxGuilds))},
		channels: &memChannelStore{s: newGenericStore[common.Snowflake, *common.Channel](sc(opts.Limits.MaxChannels))},
		users:    &memUserStore{s: newGenericStore[common.Snowflake, *common.User](sc(opts.Limits.MaxUsers))},
		members:  &memMemberStore{s: newGenericStore[memberKey, *common.GuildMember](sc(opts.Limits.MaxMembers))},
		messages: newMemMessageStore(opts.Messages, trackBytes, opts.Limits.MaxMessages, by),
		stopCh:   make(chan struct{}),
	}

	c.wg.Add(1)
	go c.cleaner()
	return c
}

func (c *MemoryCache) Guilds() GuildStore     { return c.guilds }
func (c *MemoryCache) Channels() ChannelStore { return c.channels }
func (c *MemoryCache) Users() UserStore       { return c.users }
func (c *MemoryCache) Members() MemberStore   { return c.members }
func (c *MemoryCache) Messages() MessageStore { return c.messages }

// Close stops the background sweeper. Safe to call multiple times.
func (c *MemoryCache) Close() error {
	c.stopOnce.Do(func() { close(c.stopCh) })
	c.wg.Wait()
	return nil
}

func (c *MemoryCache) cleaner() {
	defer c.wg.Done()
	ticker := time.NewTicker(c.opts.SweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.sweep()
		case <-c.stopCh:
			return
		}
	}
}

func (c *MemoryCache) sweep() {
	now := time.Now()
	b := c.opts.EvictBehavior
	uw := c.opts.UnusedWindow

	c.guilds.s.sweep(now, b, uw)
	c.channels.s.sweep(now, b, uw)
	c.users.s.sweep(now, b, uw)
	c.members.s.sweep(now, b, uw)
	c.messages.sweep(now, b, uw)

	c.enforceGlobalLimits()
}

// enforceGlobalLimits applies overflow eviction when MaxEntries or MaxSizeMB
// is exceeded. Called after every sweep and is also safe to call from Set.
func (c *MemoryCache) enforceGlobalLimits() {
	lim := c.opts.Limits
	pol := c.opts.OnOverflow

	if lim.MaxEntries > 0 {
		total := c.guilds.s.size() + c.channels.s.size() +
			c.users.s.size() + c.members.s.size() + c.messages.Size()
		if total > lim.MaxEntries {
			c.evictGloballyByCount(total-lim.MaxEntries, pol.Target, pol.ClearBy)
		}
	}

	if lim.MaxSizeMB > 0 {
		maxBytes := int64(lim.MaxSizeMB * 1024 * 1024)
		totalBytes := c.guilds.s.bytes() + c.channels.s.bytes() +
			c.users.s.bytes() + c.members.s.bytes() + c.messages.bytes()
		if totalBytes > maxBytes {
			c.evictGloballyByBytes(totalBytes-maxBytes, pol.Target, pol.ClearBy)
		}
	}
}

// evictGloballyByCount removes `need` entries from eligible stores, taking
// proportionally from each store based on its current share of the total.
func (c *MemoryCache) evictGloballyByCount(need int, target OverflowCategory, by ClearBy) {
	type storeRef struct {
		size   int
		evictN func(n int, by ClearBy) int64
	}

	stores := make([]storeRef, 0, 5)
	total := 0
	if target&CategoryGuilds != 0 {
		n := c.guilds.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.guilds.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryChannels != 0 {
		n := c.channels.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.channels.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryUsers != 0 {
		n := c.users.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.users.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryMembers != 0 {
		n := c.members.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.members.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryMessages != 0 {
		n := c.messages.Size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.messages.evictN(n, by) }})
		total += n
	}

	if total == 0 {
		return
	}
	remaining := need
	for _, st := range stores {
		if remaining <= 0 || st.size == 0 {
			continue
		}
		share := (st.size * need) / total
		if share < 1 {
			share = 1
		}
		if share > remaining {
			share = remaining
		}
		st.evictN(share, by)
		remaining -= share
	}
}

// evictGloballyByBytes removes entries from eligible stores until `need` bytes
// have been freed, taking proportionally from each store.
func (c *MemoryCache) evictGloballyByBytes(need int64, target OverflowCategory, by ClearBy) {
	type storeRef struct {
		bytes  int64
		size   int
		evictN func(n int, by ClearBy) int64
	}

	stores := make([]storeRef, 0, 5)
	totalBytes := int64(0)
	if target&CategoryGuilds != 0 {
		b := c.guilds.s.bytes()
		stores = append(stores, storeRef{b, c.guilds.s.size(), func(n int, by ClearBy) int64 { return c.guilds.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryChannels != 0 {
		b := c.channels.s.bytes()
		stores = append(stores, storeRef{b, c.channels.s.size(), func(n int, by ClearBy) int64 { return c.channels.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryUsers != 0 {
		b := c.users.s.bytes()
		stores = append(stores, storeRef{b, c.users.s.size(), func(n int, by ClearBy) int64 { return c.users.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryMembers != 0 {
		b := c.members.s.bytes()
		stores = append(stores, storeRef{b, c.members.s.size(), func(n int, by ClearBy) int64 { return c.members.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryMessages != 0 {
		b := c.messages.bytes()
		stores = append(stores, storeRef{b, c.messages.Size(), func(n int, by ClearBy) int64 { return c.messages.evictN(n, by) }})
		totalBytes += b
	}

	if totalBytes == 0 {
		return
	}
	for _, st := range stores {
		if need <= 0 || st.size == 0 {
			continue
		}
		// Evict a proportional share of entries from this store.
		share := int(float64(st.size) * float64(need) / float64(totalBytes))
		if share < 1 {
			share = 1
		}
		freed := st.evictN(share, by)
		need -= freed
	}
}
