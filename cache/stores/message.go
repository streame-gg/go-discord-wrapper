package stores

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── msgItem ───────────────────────────────────────────────────────────────────

type msgItem struct {
	msg        *discord.Message
	insertedAt time.Time
	accessedAt time.Time
	hitCount   uint64
	expiresAt  time.Time // zero = no TTL
}

// ── msgRing ───────────────────────────────────────────────────────────────────

// msgRing is a bounded per-channel buffer of messages.
type msgRing struct {
	items      []*msgItem
	cap        int
	accessedAt time.Time // updated on any read or write to this ring
}

func newMsgRing(cap int) *msgRing {
	return &msgRing{items: make([]*msgItem, 0, cap), cap: cap, accessedAt: time.Now()}
}

// add inserts or updates msg. Returns the net change in ring length:
//   - 0  → updated existing entry, or ring was full (evict one, add one)
//   - +1 → new entry added to a ring with available space
func (r *msgRing) add(msg *discord.Message, ttl time.Duration, strategy EvictionStrategy) int {
	now := time.Now()
	for _, it := range r.items {
		if it.msg.ID == msg.ID {
			it.msg = msg
			it.accessedAt = now
			return 0 // update, no size change
		}
	}
	it := &msgItem{msg: msg, insertedAt: now, accessedAt: now}
	if ttl > 0 {
		it.expiresAt = now.Add(ttl)
	}
	if len(r.items) >= r.cap {
		// Evict one, then add one: net 0.
		r.evictOne(strategy)
		r.items = append(r.items, it)
		r.accessedAt = now
		return 0
	}
	// Ring has capacity: net +1.
	r.items = append(r.items, it)
	r.accessedAt = now
	return 1
}

func (r *msgRing) evictOne(strategy EvictionStrategy) {
	if len(r.items) == 0 {
		return
	}
	idx := 0
	for i := 1; i < len(r.items); i++ {
		if isBetterMsgVictim(r.items[i], r.items[idx], strategy) {
			idx = i
		}
	}
	copy(r.items[idx:], r.items[idx+1:])
	r.items[len(r.items)-1] = nil
	r.items = r.items[:len(r.items)-1]
}

func isBetterMsgVictim(candidate, current *msgItem, strategy EvictionStrategy) bool {
	switch strategy {
	case EvictLFU:
		return candidate.hitCount < current.hitCount
	case EvictFIFO:
		return candidate.insertedAt.Before(current.insertedAt)
	default: // EvictLRU
		return candidate.accessedAt.Before(current.accessedAt)
	}
}

func (r *msgRing) get(id discord.Snowflake, now time.Time) (*discord.Message, bool) {
	for _, it := range r.items {
		if it.msg.ID == id {
			if !it.expiresAt.IsZero() && now.After(it.expiresAt) {
				return nil, false
			}
			it.accessedAt = now
			it.hitCount++
			r.accessedAt = now
			return it.msg, true
		}
	}
	return nil, false
}

func (r *msgRing) update(msg *discord.Message) {
	now := time.Now()
	for _, it := range r.items {
		if it.msg.ID == msg.ID {
			it.msg = msg
			it.accessedAt = now
			return
		}
	}
}

func (r *msgRing) delete(id discord.Snowflake) bool {
	for i, it := range r.items {
		if it.msg.ID == id {
			copy(r.items[i:], r.items[i+1:])
			r.items[len(r.items)-1] = nil
			r.items = r.items[:len(r.items)-1]
			return true
		}
	}
	return false
}

func (r *msgRing) deleteBulk(ids []discord.Snowflake) int {
	set := make(map[discord.Snowflake]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	old := r.items
	keep := old[:0]
	removed := 0
	for _, it := range old {
		if _, drop := set[it.msg.ID]; drop {
			removed++
		} else {
			keep = append(keep, it)
		}
	}
	for i := len(keep); i < len(old); i++ {
		old[i] = nil
	}
	r.items = keep
	return removed
}

// all returns messages newest-first, excluding TTL-expired entries.
func (r *msgRing) all(now time.Time, ttl time.Duration) []*discord.Message {
	r.accessedAt = now
	out := make([]*discord.Message, 0, len(r.items))
	for i := len(r.items) - 1; i >= 0; i-- {
		it := r.items[i]
		if ttl > 0 && !it.expiresAt.IsZero() && now.After(it.expiresAt) {
			continue
		}
		out = append(out, it.msg)
	}
	return out
}

// sweep removes TTL-expired and, when unusedWindow > 0, access-idle messages.
// Returns the count removed and whether the ring is now empty.
func (r *msgRing) sweep(now time.Time, unusedWindow time.Duration) (removed int, empty bool) {
	keep := r.items[:0]
	for _, it := range r.items {
		expired := !it.expiresAt.IsZero() && now.After(it.expiresAt)
		idle := unusedWindow > 0 && now.Sub(it.accessedAt) > unusedWindow
		if expired || idle {
			removed++
		} else {
			keep = append(keep, it)
		}
	}
	for i := len(keep); i < len(r.items); i++ {
		r.items[i] = nil
	}
	r.items = keep
	return removed, len(r.items) == 0
}

// trimToTotal evicts messages from channels (least-recently-accessed first)
// until len <= target. Caller must hold s.mu.
func (s *MemMessageStore) trimToTotalLocked(target int) {
	type chanRef struct {
		id   discord.Snowflake
		ring *msgRing
	}
	// Collect channels sorted by accessedAt ascending (evict oldest first).
	refs := make([]chanRef, 0, len(s.channels))
	for id, r := range s.channels {
		refs = append(refs, chanRef{id, r})
	}
	// Insertion-order is fine; we just evict from least-recently-accessed.
	// Simple sort: find min each pass would be O(n²) but rings are few.
	for i := 1; i < len(refs); i++ {
		for j := i; j > 0 && refs[j].ring.accessedAt.Before(refs[j-1].ring.accessedAt); j-- {
			refs[j], refs[j-1] = refs[j-1], refs[j]
		}
	}
	current := int(s.total.Load())
	for _, cr := range refs {
		if current <= target {
			break
		}
		need := current - target
		r := cr.ring
		toEvict := need
		if toEvict > len(r.items) {
			toEvict = len(r.items)
		}
		for i := 0; i < toEvict; i++ {
			r.evictOne(s.opts.EvictionStrategy)
		}
		s.total.Add(-int64(toEvict))
		current -= toEvict
		if len(r.items) == 0 {
			delete(s.channels, cr.id)
		}
	}
}

// ── MemMessageStore ───────────────────────────────────────────────────────────

// MemMessageStore implements cache.MessageStore using per-channel ring buffers.
//
// StoreOptions.MaxItems sets the per-channel ring capacity (0 = disabled,
// negative = use default of 100).
// StoreOptions.MaxTotal caps the sum across all channels; enforced on Add.
type MemMessageStore struct {
	mu       sync.Mutex
	channels map[discord.Snowflake]*msgRing
	opts     StoreOptions
	capPerCh int
	maxTotal int // resolved from StoreOptions; 0 = unlimited
	total    atomic.Int64

	stopCh    chan struct{}
	closeOnce sync.Once
	wg        sync.WaitGroup
}

// NewMessageStore creates a MemMessageStore and starts any configured sweepers.
// Call Close to stop background goroutines.
func NewMessageStore(opts StoreOptions) *MemMessageStore {
	cap := opts.MaxItems
	if cap < 0 {
		cap = 100
	}
	s := &MemMessageStore{
		channels: make(map[discord.Snowflake]*msgRing),
		opts:     opts,
		capPerCh: cap,
		maxTotal: opts.MaxTotal,
		stopCh:   make(chan struct{}),
	}
	for _, sw := range opts.Sweepers {
		s.startSweeper(sw)
	}
	return s
}

// IsEnabled reports whether message caching is active.
func (s *MemMessageStore) IsEnabled() bool { return !s.opts.Disabled && s.capPerCh > 0 }

func (s *MemMessageStore) Add(msg *discord.Message) {
	if !s.IsEnabled() || msg == nil {
		return
	}
	s.mu.Lock()
	r := s.channels[msg.ChannelID]
	if r == nil {
		r = newMsgRing(s.capPerCh)
		s.channels[msg.ChannelID] = r
	}
	delta := r.add(msg, s.opts.TTL, s.opts.EvictionStrategy)
	if delta != 0 {
		s.total.Add(int64(delta))
	}
	// Enforce global message cap synchronously, while still holding the lock.
	if s.maxTotal > 0 && int(s.total.Load()) > s.maxTotal {
		s.trimToTotalLocked(s.maxTotal)
	}
	s.mu.Unlock()
}

func (s *MemMessageStore) Get(channelID, messageID discord.Snowflake) (*discord.Message, bool) {
	if !s.IsEnabled() {
		return nil, false
	}
	now := time.Now()
	s.mu.Lock()
	r := s.channels[channelID]
	if r == nil {
		s.mu.Unlock()
		return nil, false
	}
	msg, ok := r.get(messageID, now)
	s.mu.Unlock()
	if !ok {
		return nil, false
	}
	cp := *msg
	return &cp, true
}

func (s *MemMessageStore) Update(msg *discord.Message) {
	if !s.IsEnabled() || msg == nil {
		return
	}
	s.mu.Lock()
	if r := s.channels[msg.ChannelID]; r != nil {
		r.update(msg)
	}
	s.mu.Unlock()
}

func (s *MemMessageStore) Delete(channelID, messageID discord.Snowflake) {
	if !s.IsEnabled() {
		return
	}
	s.mu.Lock()
	if r := s.channels[channelID]; r != nil {
		if r.delete(messageID) {
			s.total.Add(-1)
		}
	}
	s.mu.Unlock()
}

func (s *MemMessageStore) DeleteBulk(channelID discord.Snowflake, ids []discord.Snowflake) {
	if !s.IsEnabled() {
		return
	}
	s.mu.Lock()
	if r := s.channels[channelID]; r != nil {
		n := r.deleteBulk(ids)
		s.total.Add(-int64(n))
	}
	s.mu.Unlock()
}

func (s *MemMessageStore) Channel(channelID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Message] {
	empty := collection.New[discord.Snowflake, *discord.Message]()
	if !s.IsEnabled() {
		return empty
	}
	now := time.Now()
	s.mu.Lock()
	r := s.channels[channelID]
	if r == nil {
		s.mu.Unlock()
		return empty
	}
	msgs := r.all(now, s.opts.TTL)
	s.mu.Unlock()
	out := collection.NewWithCapacity[discord.Snowflake, *discord.Message](len(msgs))
	for _, m := range msgs {
		out.Set(m.ID, m)
	}
	return out
}

func (s *MemMessageStore) DeleteChannel(channelID discord.Snowflake) {
	if !s.IsEnabled() {
		return
	}
	s.mu.Lock()
	if r := s.channels[channelID]; r != nil {
		s.total.Add(-int64(len(r.items)))
		delete(s.channels, channelID)
	}
	s.mu.Unlock()
}

func (s *MemMessageStore) Size() int { return int(s.total.Load()) }

// Bytes returns 0; MemMessageStore does not track payload bytes.
func (s *MemMessageStore) Bytes() int64 { return 0 }

// TrimToTotal evicts messages (by EvictionStrategy, oldest-channel-first)
// until Size() <= n. Safe for concurrent use.
func (s *MemMessageStore) TrimToTotal(n int) {
	if !s.IsEnabled() || n < 0 {
		return
	}
	s.mu.Lock()
	s.trimToTotalLocked(n)
	s.mu.Unlock()
}

// SweepExpired removes TTL-expired messages and, when unusedWindow > 0,
// messages in channels not accessed within that window.
func (s *MemMessageStore) SweepExpired(unusedWindow time.Duration) {
	if !s.IsEnabled() {
		return
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, r := range s.channels {
		n, empty := r.sweep(now, unusedWindow)
		s.total.Add(-int64(n))
		if empty {
			delete(s.channels, id)
		}
	}
}

// Close stops all sweeper goroutines. Safe to call more than once.
func (s *MemMessageStore) Close() {
	s.closeOnce.Do(func() { close(s.stopCh) })
	s.wg.Wait()
}

func (s *MemMessageStore) startSweeper(sw Sweeper) {
	if sw.Interval <= 0 || sw.Filter == nil {
		return
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		t := time.NewTicker(sw.Interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				s.runSweeper(sw)
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *MemMessageStore) runSweeper(sw Sweeper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, r := range s.channels {
		keep := r.items[:0]
		for _, it := range r.items {
			if sw.Filter(it.msg) {
				s.total.Add(-1)
			} else {
				keep = append(keep, it)
			}
		}
		for i := len(keep); i < len(r.items); i++ {
			r.items[i] = nil
		}
		r.items = keep
		if len(r.items) == 0 {
			delete(s.channels, id)
		}
	}
}
