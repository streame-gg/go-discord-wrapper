package cache

import (
	"encoding/json"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
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
		// Preserve eviction-critical metadata so that hitCount and insertedAt
		// are not reset on every update (GUILD_UPDATE, CHANNEL_UPDATE, etc.).
		e.hitCount = old.hitCount
		e.insertedAt = old.insertedAt
		s.totalBytes.Add(-old.sizeBytes)
	}
	s.items[key] = e
	n := len(s.items)
	if s.cfg.trackBytes {
		s.totalBytes.Add(sz)
	}
	s.mu.Unlock()

	// Enforce per-category cap immediately (soft: may briefly exceed by 1).
	if s.cfg.maxItems > 0 && n > s.cfg.maxItems {
		s.evictToCount(s.cfg.maxItems)
	}
}

// setLocked is like set but assumes s.mu is already held by the caller.
// Used by SetAll implementations to keep the delete→insert sequence atomic.
// If the store has a non-zero maxItems cap, the caller must call evictToCount
// after releasing the lock — setLocked does not enforce the cap itself.
func (s *genericStore[K, V]) setLocked(key K, value V) {
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
	if old, ok := s.items[key]; ok {
		e.hitCount = old.hitCount
		e.insertedAt = old.insertedAt
		s.totalBytes.Add(-old.sizeBytes)
	}
	s.items[key] = e
	if s.cfg.trackBytes {
		s.totalBytes.Add(sz)
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
	toRemove := n - target

	// Fast path: single-item eviction — linear scan, no allocations.
	if toRemove == 1 {
		var minKey K
		var minScore float64
		var minInserted int64
		var minSize int64
		first := true
		for k, e := range s.items {
			score := e.rank(s.cfg.clearBy)
			ins := e.insertedAt.UnixNano()
			if first || score < minScore || (score == minScore && ins < minInserted) {
				minKey = k
				minScore = score
				minInserted = ins
				minSize = e.sizeBytes
				first = false
			}
		}
		s.totalBytes.Add(-minSize)
		delete(s.items, minKey)
		return
	}

	// General path: sort then bulk-remove.
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
	for _, p := range pairs[:toRemove] {
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
	s *genericStore[discord.Snowflake, *discord.Guild]
}

func (g *memGuildStore) Set(guild *discord.Guild) { g.s.set(guild.ID, guild) }
func (g *memGuildStore) Get(id discord.Snowflake) (*discord.Guild, bool) {
	v, ok := g.s.get(id)
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}
func (g *memGuildStore) Delete(id discord.Snowflake) { g.s.delete(id) }
func (g *memGuildStore) All() []*discord.Guild       { return g.s.all() }
func (g *memGuildStore) Size() int                   { return g.s.size() }

type memChannelStore struct {
	s *genericStore[discord.Snowflake, *discord.Channel]
}

func (c *memChannelStore) Set(ch *discord.Channel) { c.s.set(ch.ID, ch) }
func (c *memChannelStore) Get(id discord.Snowflake) (*discord.Channel, bool) {
	v, ok := c.s.get(id)
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}
func (c *memChannelStore) Delete(id discord.Snowflake) { c.s.delete(id) }
func (c *memChannelStore) All() []*discord.Channel     { return c.s.all() }
func (c *memChannelStore) Size() int                   { return c.s.size() }

type memUserStore struct {
	s *genericStore[discord.Snowflake, *discord.User]
}

func (u *memUserStore) Set(user *discord.User) { u.s.set(user.ID, user) }
func (u *memUserStore) Get(id discord.Snowflake) (*discord.User, bool) {
	v, ok := u.s.get(id)
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}
func (u *memUserStore) Delete(id discord.Snowflake) { u.s.delete(id) }
func (u *memUserStore) All() []*discord.User        { return u.s.all() }
func (u *memUserStore) Size() int                   { return u.s.size() }

// ── Member store ──────────────────────────────────────────────────────────────

type memberKey struct {
	guildID discord.Snowflake
	userID  discord.Snowflake
}

type memMemberStore struct {
	s *genericStore[memberKey, *discord.GuildMember]
}

func (m *memMemberStore) Set(guildID discord.Snowflake, member *discord.GuildMember) {
	if member == nil || member.User == nil {
		return
	}
	m.s.set(memberKey{guildID, member.User.ID}, member)
}

func (m *memMemberStore) Get(guildID, userID discord.Snowflake) (*discord.GuildMember, bool) {
	v, ok := m.s.get(memberKey{guildID, userID})
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}

func (m *memMemberStore) Delete(guildID, userID discord.Snowflake) {
	m.s.delete(memberKey{guildID, userID})
}

func (m *memMemberStore) DeleteGuild(guildID discord.Snowflake) {
	m.s.mu.Lock()
	defer m.s.mu.Unlock()
	for k, e := range m.s.items {
		if k.guildID == guildID {
			m.s.totalBytes.Add(-e.sizeBytes)
			delete(m.s.items, k)
		}
	}
}

func (m *memMemberStore) AllInGuild(guildID discord.Snowflake) []*discord.GuildMember {
	now := time.Now()
	m.s.mu.RLock()
	defer m.s.mu.RUnlock()
	var out []*discord.GuildMember
	for k, e := range m.s.items {
		if k.guildID == guildID && !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (m *memMemberStore) Size() int { return m.s.size() }

// ── Role store ────────────────────────────────────────────────────────────────

type roleValue struct {
	GuildID discord.Snowflake
	Role    *discord.Role
}

type memRoleStore struct {
	s *genericStore[discord.Snowflake, roleValue]
}

func (r *memRoleStore) Set(guildID discord.Snowflake, role *discord.Role) {
	r.s.set(role.ID, roleValue{GuildID: guildID, Role: role})
}

func (r *memRoleStore) Get(roleID discord.Snowflake) (*discord.Role, bool) {
	v, ok := r.s.get(roleID)
	if !ok {
		return nil, false
	}
	cp := *v.Role
	return &cp, true
}

func (r *memRoleStore) GetByGuild(guildID discord.Snowflake) []*discord.Role {
	now := time.Now()
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var out []*discord.Role
	for _, e := range r.s.items {
		if e.value.GuildID == guildID && !e.expired(now) {
			out = append(out, e.value.Role)
		}
	}
	return out
}

func (r *memRoleStore) Delete(roleID discord.Snowflake) {
	r.s.delete(roleID)
}

func (r *memRoleStore) DeleteGuild(guildID discord.Snowflake) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for k, e := range r.s.items {
		if e.value.GuildID == guildID {
			r.s.totalBytes.Add(-e.sizeBytes)
			delete(r.s.items, k)
		}
	}
}

func (r *memRoleStore) All() []*discord.Role {
	vals := r.s.all()
	out := make([]*discord.Role, len(vals))
	for i, v := range vals {
		out[i] = v.Role
	}
	return out
}

func (r *memRoleStore) Size() int { return r.s.size() }

// ── Message store ─────────────────────────────────────────────────────────────

type msgEntry struct {
	msg        *discord.Message
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

func (r *channelRing) add(msg *discord.Message, trackBytes bool) (delta int64) {
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

func (r *channelRing) get(id discord.Snowflake, ttl time.Duration) (*discord.Message, bool) {
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

func (r *channelRing) update(msg *discord.Message, trackBytes bool) (delta int64) {
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

func (r *channelRing) delete(id discord.Snowflake) (delta int64) {
	for i, e := range r.msgs {
		if e.msg.ID == id {
			delta = -e.sizeBytes
			copy(r.msgs[i:], r.msgs[i+1:])
			r.msgs[len(r.msgs)-1] = nil // zero tail so GC can collect it
			r.msgs = r.msgs[:len(r.msgs)-1]
			return
		}
	}
	return
}

func (r *channelRing) deleteBulk(ids []discord.Snowflake) (delta int64) {
	set := make(map[discord.Snowflake]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	old := r.msgs
	keep := old[:0]
	for _, e := range old {
		if _, drop := set[e.msg.ID]; drop {
			delta -= e.sizeBytes
		} else {
			keep = append(keep, e)
		}
	}
	// Zero freed tail slots so GC can collect the dropped entries.
	for i := len(keep); i < len(old); i++ {
		old[i] = nil
	}
	r.msgs = keep
	return
}

// all returns messages newest-first, skipping TTL-expired entries.
func (r *channelRing) all(ttl time.Duration) []*discord.Message {
	now := time.Now()
	r.accessedAt = now
	out := make([]*discord.Message, 0, len(r.msgs))
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
	channels   map[discord.Snowflake]*channelRing
	opts       MessageOptions
	trackBytes bool
	maxTotal   int // Limits.MaxMessages; 0 = unlimited
	clearBy    ClearBy
	totalBytes atomic.Int64
	totalMsgs  atomic.Int64
}

func newMemMessageStore(opts MessageOptions, trackBytes bool, maxTotal int, clearBy ClearBy) *memMessageStore {
	return &memMessageStore{
		channels:   make(map[discord.Snowflake]*channelRing),
		opts:       opts,
		trackBytes: trackBytes,
		maxTotal:   maxTotal,
		clearBy:    clearBy,
	}
}

func (s *memMessageStore) Add(msg *discord.Message) {
	if s.opts.MaxPerChannel == 0 {
		return
	}

	// Hold the write lock for the entire lookup-or-create AND r.add() so that
	// a concurrent DeleteChannel cannot remove the ring between ring() and add().
	s.mu.Lock()
	r := s.channels[msg.ChannelID]
	if r == nil {
		r = newChannelRing(s.opts.MaxPerChannel)
		s.channels[msg.ChannelID] = r
	}
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

func (s *memMessageStore) Get(channelID, messageID discord.Snowflake) (*discord.Message, bool) {
	s.mu.RLock()
	r := s.channels[channelID]
	s.mu.RUnlock()
	if r == nil {
		return nil, false
	}
	s.mu.Lock()
	msg, ok := r.get(messageID, s.opts.TTL)
	s.mu.Unlock()
	if !ok {
		return nil, false
	}
	cp := *msg
	return &cp, true
}

func (s *memMessageStore) Update(msg *discord.Message) {
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

func (s *memMessageStore) Delete(channelID, messageID discord.Snowflake) {
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

func (s *memMessageStore) DeleteBulk(channelID discord.Snowflake, ids []discord.Snowflake) {
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

func (s *memMessageStore) Channel(channelID discord.Snowflake) []*discord.Message {
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

func (s *memMessageStore) DeleteChannel(channelID discord.Snowflake) {
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
		if empty {
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
		id   discord.Snowflake
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
		id   discord.Snowflake
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
	opts            Options
	guilds          *memGuildStore
	channels        *memChannelStore
	users           *memUserStore
	members         *memMemberStore
	roles           *memRoleStore
	messages        *memMessageStore
	voiceStates     *memVoiceStateStore
	soundboard      *memSoundboardStore
	scheduledEvents *memScheduledEventStore
	stageInstances  *memStageInstanceStore
	emojis          *memEmojiStore
	stickers        *memStickerStore
	presences       *memPresenceStore

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
		// Negative means "use default".
		opts.Messages.MaxPerChannel = 100
	}
	// == 0 means "disabled" (caller set it explicitly) — leave as-is.
	// > 0 means "user value" — leave as-is.
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
		guilds:   &memGuildStore{s: newGenericStore[discord.Snowflake, *discord.Guild](sc(opts.Limits.MaxGuilds))},
		channels: &memChannelStore{s: newGenericStore[discord.Snowflake, *discord.Channel](sc(opts.Limits.MaxChannels))},
		users:    &memUserStore{s: newGenericStore[discord.Snowflake, *discord.User](sc(opts.Limits.MaxUsers))},
		members:  &memMemberStore{s: newGenericStore[memberKey, *discord.GuildMember](sc(opts.Limits.MaxMembers))},
		roles:    &memRoleStore{s: newGenericStore[discord.Snowflake, roleValue](sc(opts.Limits.MaxRoles))},
		messages: newMemMessageStore(opts.Messages, trackBytes, opts.Limits.MaxMessages, by),
		voiceStates: &memVoiceStateStore{
			s: newGenericStore[memberKey, *discord.VoiceState](sc(0)),
		},
		soundboard: &memSoundboardStore{
			s: newGenericStore[discord.Snowflake, soundboardValue](sc(0)),
		},
		scheduledEvents: &memScheduledEventStore{
			s: newGenericStore[discord.Snowflake, *discord.GuildScheduledEvent](sc(0)),
		},
		stageInstances: &memStageInstanceStore{
			s: newGenericStore[discord.Snowflake, *discord.StageInstance](sc(0)),
		},
		emojis: &memEmojiStore{
			s: newGenericStore[discord.Snowflake, emojiValue](sc(0)),
		},
		stickers: &memStickerStore{
			s: newGenericStore[discord.Snowflake, stickerValue](sc(opts.Limits.MaxStickers)),
		},
		presences: &memPresenceStore{
			s: newGenericStore[memberKey, *discord.Presence](sc(0)),
		},
		stopCh: make(chan struct{}),
	}

	c.wg.Add(1)
	go c.cleaner()
	return c
}

func (c *MemoryCache) Guilds() GuildStore     { return c.guilds }
func (c *MemoryCache) Channels() ChannelStore { return c.channels }
func (c *MemoryCache) Users() UserStore       { return c.users }
func (c *MemoryCache) Members() MemberStore   { return c.members }
func (c *MemoryCache) Roles() RoleStore       { return c.roles }
func (c *MemoryCache) Messages() MessageStore { return c.messages }
func (c *MemoryCache) VoiceStates() VoiceStateStore {
	return c.voiceStates
}
func (c *MemoryCache) Soundboard() SoundboardStore {
	return c.soundboard
}
func (c *MemoryCache) ScheduledEvents() ScheduledEventStore {
	return c.scheduledEvents
}
func (c *MemoryCache) StageInstances() StageInstanceStore {
	return c.stageInstances
}
func (c *MemoryCache) Emojis() EmojiStore { return c.emojis }
func (c *MemoryCache) Stickers() StickerStore {
	return c.stickers
}
func (c *MemoryCache) Presences() PresenceStore {
	return c.presences
}

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
	c.roles.s.sweep(now, b, uw)
	c.messages.sweep(now, b, uw)
	c.voiceStates.s.sweep(now, b, uw)
	c.soundboard.s.sweep(now, b, uw)
	c.scheduledEvents.s.sweep(now, b, uw)
	c.stageInstances.s.sweep(now, b, uw)
	c.emojis.s.sweep(now, b, uw)
	c.stickers.s.sweep(now, b, uw)
	c.presences.s.sweep(now, b, uw)

	c.enforceGlobalLimits()
}

// enforceGlobalLimits applies overflow eviction when MaxEntries or MaxSizeMB
// is exceeded. Called after every sweep and is also safe to call from Set.
func (c *MemoryCache) enforceGlobalLimits() {
	lim := c.opts.Limits
	pol := c.opts.OnOverflow

	if lim.MaxEntries > 0 {
		total := c.guilds.s.size() + c.channels.s.size() +
			c.users.s.size() + c.members.s.size() + c.roles.s.size() + c.messages.Size() +
			c.voiceStates.s.size() + c.soundboard.s.size() + c.scheduledEvents.s.size() +
			c.stageInstances.s.size() + c.emojis.s.size() + c.stickers.s.size() + c.presences.s.size()
		if total > lim.MaxEntries {
			c.evictGloballyByCount(total-lim.MaxEntries, pol.Target, pol.ClearBy)
		}
	}

	if lim.MaxSizeMB > 0 {
		maxBytes := int64(lim.MaxSizeMB * 1024 * 1024)
		totalBytes := c.guilds.s.bytes() + c.channels.s.bytes() +
			c.users.s.bytes() + c.members.s.bytes() + c.roles.s.bytes() + c.messages.bytes() +
			c.voiceStates.s.bytes() + c.soundboard.s.bytes() + c.scheduledEvents.s.bytes() +
			c.stageInstances.s.bytes() + c.emojis.s.bytes() + c.stickers.s.bytes() + c.presences.s.bytes()
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

	stores := make([]storeRef, 0, 6)
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
	if target&CategoryRoles != 0 {
		n := c.roles.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.roles.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryMessages != 0 {
		n := c.messages.Size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.messages.evictN(n, by) }})
		total += n
	}
	if target&CategoryVoiceStates != 0 {
		n := c.voiceStates.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.voiceStates.s.evictN(n, by) }})
		total += n
	}
	if target&CategorySoundboard != 0 {
		n := c.soundboard.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.soundboard.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryScheduledEvents != 0 {
		n := c.scheduledEvents.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.scheduledEvents.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryStageInstances != 0 {
		n := c.stageInstances.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.stageInstances.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryEmojis != 0 {
		n := c.emojis.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.emojis.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryStickers != 0 {
		n := c.stickers.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.stickers.s.evictN(n, by) }})
		total += n
	}
	if target&CategoryPresences != 0 {
		n := c.presences.s.size()
		stores = append(stores, storeRef{n, func(n int, by ClearBy) int64 { return c.presences.s.evictN(n, by) }})
		total += n
	}

	if total == 0 {
		return
	}

	// Largest-Remainder Method: distributes exactly `need` evictions across stores
	// proportionally without over-evicting. The old "share < 1 → 1" guard forced
	// at least one eviction per store even when the store's proportional share was
	// less than 1, causing the total evicted to exceed `need` (Bug 9).
	type slotted struct {
		idx       int
		floor     int
		remainder float64
	}
	slots := make([]slotted, 0, len(stores))
	floorSum := 0
	for i, st := range stores {
		if st.size == 0 {
			continue
		}
		real := float64(st.size*need) / float64(total)
		f := int(real)
		slots = append(slots, slotted{i, f, real - float64(f)})
		floorSum += f
	}
	// Give one extra slot each to the stores with the largest fractional remainders
	// until the total allocated equals need.
	extra := need - floorSum
	if extra > 0 {
		sort.Slice(slots, func(a, b int) bool { return slots[a].remainder > slots[b].remainder })
		for i := 0; i < extra && i < len(slots); i++ {
			slots[i].floor++
		}
	}
	for _, sl := range slots {
		if sl.floor > 0 {
			stores[sl.idx].evictN(sl.floor, by)
		}
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

	stores := make([]storeRef, 0, 6)
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
	if target&CategoryRoles != 0 {
		b := c.roles.s.bytes()
		stores = append(stores, storeRef{b, c.roles.s.size(), func(n int, by ClearBy) int64 { return c.roles.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryMessages != 0 {
		b := c.messages.bytes()
		stores = append(stores, storeRef{b, c.messages.Size(), func(n int, by ClearBy) int64 { return c.messages.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryVoiceStates != 0 {
		b := c.voiceStates.s.bytes()
		stores = append(stores, storeRef{b, c.voiceStates.s.size(), func(n int, by ClearBy) int64 { return c.voiceStates.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategorySoundboard != 0 {
		b := c.soundboard.s.bytes()
		stores = append(stores, storeRef{b, c.soundboard.s.size(), func(n int, by ClearBy) int64 { return c.soundboard.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryScheduledEvents != 0 {
		b := c.scheduledEvents.s.bytes()
		stores = append(stores, storeRef{b, c.scheduledEvents.s.size(), func(n int, by ClearBy) int64 { return c.scheduledEvents.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryStageInstances != 0 {
		b := c.stageInstances.s.bytes()
		stores = append(stores, storeRef{b, c.stageInstances.s.size(), func(n int, by ClearBy) int64 { return c.stageInstances.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryEmojis != 0 {
		b := c.emojis.s.bytes()
		stores = append(stores, storeRef{b, c.emojis.s.size(), func(n int, by ClearBy) int64 { return c.emojis.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryStickers != 0 {
		b := c.stickers.s.bytes()
		stores = append(stores, storeRef{b, c.stickers.s.size(), func(n int, by ClearBy) int64 { return c.stickers.s.evictN(n, by) }})
		totalBytes += b
	}
	if target&CategoryPresences != 0 {
		b := c.presences.s.bytes()
		stores = append(stores, storeRef{b, c.presences.s.size(), func(n int, by ClearBy) int64 { return c.presences.s.evictN(n, by) }})
		totalBytes += b
	}

	if totalBytes == 0 {
		return
	}

	// Largest-Remainder Method adapted for bytes: converts proportional byte
	// shares to integer entry counts without forcing a minimum of 1 per store
	// (which would cause over-eviction when byteShare < avgSize for many stores).
	type slotted struct {
		idx       int
		floor     int
		remainder float64
		avgSize   int64
	}
	slots := make([]slotted, 0, len(stores))
	floorBytes := int64(0)
	for i, st := range stores {
		if st.bytes == 0 || st.size == 0 {
			continue
		}
		byteShare := need * st.bytes / totalBytes
		if byteShare <= 0 {
			continue
		}
		avgSize := st.bytes / int64(st.size)
		if avgSize <= 0 {
			avgSize = 1
		}
		real := float64(byteShare) / float64(avgSize)
		f := int(real)
		slots = append(slots, slotted{i, f, real - float64(f), avgSize})
		floorBytes += int64(f) * avgSize
	}
	// Distribute remaining bytes to the stores with the largest fractional
	// remainder (i.e., those closest to needing one extra eviction). Stop as
	// soon as the floor allocation already covers `need`.
	remainingNeed := need - floorBytes
	sort.Slice(slots, func(a, b int) bool { return slots[a].remainder > slots[b].remainder })
	for i := range slots {
		if remainingNeed <= 0 || slots[i].remainder == 0 {
			break
		}
		slots[i].floor++
		remainingNeed -= slots[i].avgSize
	}
	for _, sl := range slots {
		if sl.floor > 0 {
			stores[sl.idx].evictN(sl.floor, by)
		}
	}
}

// ── Voice state store ──────────────────────────────────────────────────────────

type memVoiceStateStore struct {
	s *genericStore[memberKey, *discord.VoiceState]
}

func (v *memVoiceStateStore) Set(guildID discord.Snowflake, vs *discord.VoiceState) {
	v.s.set(memberKey{guildID, vs.UserID}, vs)
}

func (v *memVoiceStateStore) Get(guildID, userID discord.Snowflake) (*discord.VoiceState, bool) {
	val, ok := v.s.get(memberKey{guildID, userID})
	if !ok {
		return nil, false
	}
	cp := *val
	return &cp, true
}

func (v *memVoiceStateStore) Delete(guildID, userID discord.Snowflake) {
	v.s.delete(memberKey{guildID, userID})
}

func (v *memVoiceStateStore) DeleteGuild(guildID discord.Snowflake) {
	v.s.mu.Lock()
	defer v.s.mu.Unlock()
	for k, e := range v.s.items {
		if k.guildID == guildID {
			v.s.totalBytes.Add(-e.sizeBytes)
			delete(v.s.items, k)
		}
	}
}

func (v *memVoiceStateStore) AllInGuild(guildID discord.Snowflake) []*discord.VoiceState {
	now := time.Now()
	v.s.mu.RLock()
	defer v.s.mu.RUnlock()
	var out []*discord.VoiceState
	for k, e := range v.s.items {
		if k.guildID == guildID && !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (v *memVoiceStateStore) Size() int { return v.s.size() }

// ── Soundboard store ───────────────────────────────────────────────────────────

type soundboardValue struct {
	GuildID discord.Snowflake
	Sound   *discord.SoundboardSound
}

type memSoundboardStore struct {
	s *genericStore[discord.Snowflake, soundboardValue]
}

func (ss *memSoundboardStore) Set(guildID discord.Snowflake, sound *discord.SoundboardSound) {
	if sound == nil {
		return
	}
	ss.s.set(sound.SoundID, soundboardValue{GuildID: guildID, Sound: sound})
}

func (ss *memSoundboardStore) Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool) {
	v, ok := ss.s.get(soundID)
	if !ok {
		return nil, false
	}
	cp := *v.Sound
	return &cp, true
}

func (ss *memSoundboardStore) GetByGuild(guildID discord.Snowflake) []*discord.SoundboardSound {
	now := time.Now()
	ss.s.mu.RLock()
	defer ss.s.mu.RUnlock()
	var out []*discord.SoundboardSound
	for _, e := range ss.s.items {
		if e.value.GuildID == guildID && !e.expired(now) {
			out = append(out, e.value.Sound)
		}
	}
	return out
}

func (ss *memSoundboardStore) SetAll(guildID discord.Snowflake, sounds []*discord.SoundboardSound) {
	ss.s.mu.Lock()
	for k, e := range ss.s.items {
		if e.value.GuildID == guildID {
			ss.s.totalBytes.Add(-e.sizeBytes)
			delete(ss.s.items, k)
		}
	}
	for _, s := range sounds {
		if s != nil {
			ss.s.setLocked(s.SoundID, soundboardValue{GuildID: guildID, Sound: s})
		}
	}
	n := len(ss.s.items)
	ss.s.mu.Unlock()
	if ss.s.cfg.maxItems > 0 && n > ss.s.cfg.maxItems {
		ss.s.evictToCount(ss.s.cfg.maxItems)
	}
}

func (ss *memSoundboardStore) Delete(soundID discord.Snowflake) { ss.s.delete(soundID) }

func (ss *memSoundboardStore) DeleteGuild(guildID discord.Snowflake) {
	ss.s.mu.Lock()
	defer ss.s.mu.Unlock()
	for k, e := range ss.s.items {
		if e.value.GuildID == guildID {
			ss.s.totalBytes.Add(-e.sizeBytes)
			delete(ss.s.items, k)
		}
	}
}

func (ss *memSoundboardStore) Size() int { return ss.s.size() }

// ── Scheduled event store ──────────────────────────────────────────────────────

type memScheduledEventStore struct {
	s *genericStore[discord.Snowflake, *discord.GuildScheduledEvent]
}

func (se *memScheduledEventStore) Set(event *discord.GuildScheduledEvent) {
	se.s.set(event.ID, event)
}

func (se *memScheduledEventStore) Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool) {
	v, ok := se.s.get(eventID)
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}

func (se *memScheduledEventStore) GetByGuild(guildID discord.Snowflake) []*discord.GuildScheduledEvent {
	now := time.Now()
	se.s.mu.RLock()
	defer se.s.mu.RUnlock()
	var out []*discord.GuildScheduledEvent
	for _, e := range se.s.items {
		if e.value.GuildID == guildID && !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (se *memScheduledEventStore) Delete(eventID discord.Snowflake) { se.s.delete(eventID) }

func (se *memScheduledEventStore) DeleteGuild(guildID discord.Snowflake) {
	se.s.mu.Lock()
	defer se.s.mu.Unlock()
	for k, e := range se.s.items {
		if e.value.GuildID == guildID {
			se.s.totalBytes.Add(-e.sizeBytes)
			delete(se.s.items, k)
		}
	}
}

func (se *memScheduledEventStore) Size() int { return se.s.size() }

// ── Stage instance store ───────────────────────────────────────────────────────

type memStageInstanceStore struct {
	s *genericStore[discord.Snowflake, *discord.StageInstance]
}

func (si *memStageInstanceStore) Set(instance *discord.StageInstance) {
	si.s.set(instance.ID, instance)
}

func (si *memStageInstanceStore) Get(instanceID discord.Snowflake) (*discord.StageInstance, bool) {
	v, ok := si.s.get(instanceID)
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}

func (si *memStageInstanceStore) GetByGuild(guildID discord.Snowflake) []*discord.StageInstance {
	now := time.Now()
	si.s.mu.RLock()
	defer si.s.mu.RUnlock()
	var out []*discord.StageInstance
	for _, e := range si.s.items {
		if e.value.GuildID == guildID && !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (si *memStageInstanceStore) Delete(instanceID discord.Snowflake) { si.s.delete(instanceID) }

func (si *memStageInstanceStore) DeleteGuild(guildID discord.Snowflake) {
	si.s.mu.Lock()
	defer si.s.mu.Unlock()
	for k, e := range si.s.items {
		if e.value.GuildID == guildID {
			si.s.totalBytes.Add(-e.sizeBytes)
			delete(si.s.items, k)
		}
	}
}

func (si *memStageInstanceStore) Size() int { return si.s.size() }

// ── Emoji store ────────────────────────────────────────────────────────────────

type emojiValue struct {
	GuildID discord.Snowflake
	Emoji   *discord.Emoji
}

type memEmojiStore struct {
	s *genericStore[discord.Snowflake, emojiValue]
}

func (es *memEmojiStore) Set(guildID discord.Snowflake, emoji *discord.Emoji) {
	es.s.set(emoji.ID, emojiValue{GuildID: guildID, Emoji: emoji})
}

func (es *memEmojiStore) Get(emojiID discord.Snowflake) (*discord.Emoji, bool) {
	v, ok := es.s.get(emojiID)
	if !ok {
		return nil, false
	}
	cp := *v.Emoji
	return &cp, true
}

func (es *memEmojiStore) GetByGuild(guildID discord.Snowflake) []*discord.Emoji {
	now := time.Now()
	es.s.mu.RLock()
	defer es.s.mu.RUnlock()
	var out []*discord.Emoji
	for _, e := range es.s.items {
		if e.value.GuildID == guildID && !e.expired(now) {
			out = append(out, e.value.Emoji)
		}
	}
	return out
}

func (es *memEmojiStore) SetAll(guildID discord.Snowflake, emojis []*discord.Emoji) {
	es.s.mu.Lock()
	for k, e := range es.s.items {
		if e.value.GuildID == guildID {
			es.s.totalBytes.Add(-e.sizeBytes)
			delete(es.s.items, k)
		}
	}
	for _, emoji := range emojis {
		if emoji != nil && emoji.ID != "" {
			es.s.setLocked(emoji.ID, emojiValue{GuildID: guildID, Emoji: emoji})
		}
	}
	n := len(es.s.items)
	es.s.mu.Unlock()
	if es.s.cfg.maxItems > 0 && n > es.s.cfg.maxItems {
		es.s.evictToCount(es.s.cfg.maxItems)
	}
}

func (es *memEmojiStore) Delete(emojiID discord.Snowflake) { es.s.delete(emojiID) }

func (es *memEmojiStore) DeleteGuild(guildID discord.Snowflake) {
	es.s.mu.Lock()
	defer es.s.mu.Unlock()
	for k, e := range es.s.items {
		if e.value.GuildID == guildID {
			es.s.totalBytes.Add(-e.sizeBytes)
			delete(es.s.items, k)
		}
	}
}

func (es *memEmojiStore) Size() int { return es.s.size() }

// ── Sticker store ──────────────────────────────────────────────────────────────

type stickerValue struct {
	GuildID discord.Snowflake
	Sticker *discord.Sticker
}

type memStickerStore struct {
	s *genericStore[discord.Snowflake, stickerValue]
}

func (ss *memStickerStore) Set(guildID discord.Snowflake, sticker *discord.Sticker) {
	ss.s.set(sticker.ID, stickerValue{GuildID: guildID, Sticker: sticker})
}

func (ss *memStickerStore) Get(stickerID discord.Snowflake) (*discord.Sticker, bool) {
	v, ok := ss.s.get(stickerID)
	if !ok {
		return nil, false
	}
	cp := *v.Sticker
	return &cp, true
}

func (ss *memStickerStore) GetByGuild(guildID discord.Snowflake) []*discord.Sticker {
	now := time.Now()
	ss.s.mu.RLock()
	defer ss.s.mu.RUnlock()
	var out []*discord.Sticker
	for _, e := range ss.s.items {
		if e.value.GuildID == guildID && !e.expired(now) {
			out = append(out, e.value.Sticker)
		}
	}
	return out
}

func (ss *memStickerStore) SetAll(guildID discord.Snowflake, stickers []*discord.Sticker) {
	ss.s.mu.Lock()
	for k, e := range ss.s.items {
		if e.value.GuildID == guildID {
			ss.s.totalBytes.Add(-e.sizeBytes)
			delete(ss.s.items, k)
		}
	}
	for _, sticker := range stickers {
		if sticker != nil && sticker.ID != "" {
			ss.s.setLocked(sticker.ID, stickerValue{GuildID: guildID, Sticker: sticker})
		}
	}
	n := len(ss.s.items)
	ss.s.mu.Unlock()
	if ss.s.cfg.maxItems > 0 && n > ss.s.cfg.maxItems {
		ss.s.evictToCount(ss.s.cfg.maxItems)
	}
}

func (ss *memStickerStore) Delete(stickerID discord.Snowflake) { ss.s.delete(stickerID) }

func (ss *memStickerStore) DeleteGuild(guildID discord.Snowflake) {
	ss.s.mu.Lock()
	defer ss.s.mu.Unlock()
	for k, e := range ss.s.items {
		if e.value.GuildID == guildID {
			ss.s.totalBytes.Add(-e.sizeBytes)
			delete(ss.s.items, k)
		}
	}
}

func (ss *memStickerStore) Size() int { return ss.s.size() }

// ── Presence store ─────────────────────────────────────────────────────────────

type memPresenceStore struct {
	s *genericStore[memberKey, *discord.Presence]
}

func (ps *memPresenceStore) Set(presence *discord.Presence) {
	if presence == nil {
		return
	}
	ps.s.set(memberKey{presence.GuildID, presence.User.ID}, presence)
}

func (ps *memPresenceStore) Get(guildID, userID discord.Snowflake) (*discord.Presence, bool) {
	v, ok := ps.s.get(memberKey{guildID, userID})
	if !ok {
		return nil, false
	}
	cp := *v
	return &cp, true
}

func (ps *memPresenceStore) GetByGuild(guildID discord.Snowflake) []*discord.Presence {
	now := time.Now()
	ps.s.mu.RLock()
	defer ps.s.mu.RUnlock()
	var out []*discord.Presence
	for k, e := range ps.s.items {
		if k.guildID == guildID && !e.expired(now) {
			out = append(out, e.value)
		}
	}
	return out
}

func (ps *memPresenceStore) Delete(guildID, userID discord.Snowflake) {
	ps.s.delete(memberKey{guildID, userID})
}

func (ps *memPresenceStore) DeleteGuild(guildID discord.Snowflake) {
	ps.s.mu.Lock()
	defer ps.s.mu.Unlock()
	for k, e := range ps.s.items {
		if k.guildID == guildID {
			ps.s.totalBytes.Add(-e.sizeBytes)
			delete(ps.s.items, k)
		}
	}
}

func (ps *memPresenceStore) Size() int { return ps.s.size() }
