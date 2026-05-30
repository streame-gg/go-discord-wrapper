// Package stores provides the generic BaseStore foundation and all option
// types used by the cache package.
//
// Consumers should import the parent cache package, not this one directly.
// All public types here are re-exported from cache via type aliases.
package stores

import (
	"encoding/json"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// estimateSize returns the JSON-encoded byte length of v as a size proxy.
// Not exact Go heap usage, but consistent and useful as a relative bound.
func estimateSize[V any](v V) int64 {
	b, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	return int64(len(b))
}

// ── Eviction strategy ─────────────────────────────────────────────────────────

// EvictionStrategy decides which entry to evict when a store exceeds MaxItems.
type EvictionStrategy int

const (
	// EvictLRU evicts the entry accessed least recently. Default.
	EvictLRU EvictionStrategy = iota

	// EvictLFU evicts the entry with the lowest hit count.
	EvictLFU

	// EvictFIFO evicts the entry that was inserted earliest.
	EvictFIFO
)

func (e EvictionStrategy) String() string {
	switch e {
	case EvictLRU:
		return "LRU"
	case EvictLFU:
		return "LFU"
	case EvictFIFO:
		return "FIFO"
	default:
		return "unknown"
	}
}

// ── Sweeper ───────────────────────────────────────────────────────────────────

// Sweeper is a periodic cleanup task attached to a store.
// The Filter receives each stored value and returns true to remove it.
type Sweeper struct {
	// Interval between runs. Must be > 0.
	Interval time.Duration

	// Filter returns true for values that should be removed.
	Filter func(value any) bool

	// Name is an optional label for logging/debugging.
	Name string
}

// ── StoreOptions ──────────────────────────────────────────────────────────────

// StoreOptions configures a single entity store.
type StoreOptions struct {
	// MaxItems is the hard cap on entries in this store.
	// When exceeded, one entry is evicted according to EvictionStrategy.
	// 0 = unlimited.
	//
	// For MessageStore, MaxItems applies per-channel, not globally.
	MaxItems int

	// TTL is how long an entry lives after insertion.
	// If RefreshOnAccess is true, TTL is reset on every successful Get.
	// 0 = entries never expire.
	TTL time.Duration

	// RefreshOnAccess resets the TTL clock on every successful Get.
	RefreshOnAccess bool

	// Sweepers run in background goroutines and periodically remove entries
	// whose value matches the Filter. Multiple sweepers are allowed.
	Sweepers []Sweeper

	// EvictionStrategy selects which entry to evict on MaxItems overflow.
	EvictionStrategy EvictionStrategy

	// MaxTotal caps the total entry count across all sub-groups in stores that
	// support it. For MemMessageStore, this is the total messages across all
	// channels; enforced synchronously on Add. 0 = unlimited.
	MaxTotal int

	// TrackBytes enables JSON-based payload size estimation. Required by
	// MemoryCache when Limits.MaxSizeMB is set. Adds one json.Marshal per Set.
	TrackBytes bool

	// Disabled makes the store a no-op: all writes are dropped and reads
	// always return the zero value. Useful for disabling high-volume stores
	// such as PresenceStore in bots that do not need presence data.
	Disabled bool
}

// Defaults returns options for general-purpose stores: unlimited entries,
// no TTL, LRU eviction strategy.
func Defaults() StoreOptions {
	return StoreOptions{EvictionStrategy: EvictLRU}
}

// MessageDefaults returns options tuned for MessageStore: 100 messages
// per channel, no TTL, LRU eviction.
func MessageDefaults() StoreOptions {
	return StoreOptions{MaxItems: 100, EvictionStrategy: EvictLRU}
}

// PresenceDefaults returns options for PresenceStore: disabled by default.
// Presence events are extremely high volume; opt in explicitly via
// cache.WithPresenceStore(cache.StoreOptions{...}).
func PresenceDefaults() StoreOptions {
	return StoreOptions{Disabled: true}
}

// ── Composite key types ───────────────────────────────────────────────────────

// MemberKey is the composite cache key for guild members.
type MemberKey struct {
	GuildID discord.Snowflake
	UserID  discord.Snowflake
}

// VoiceStateKey is a composite cache key for voice states (same shape as MemberKey).
type VoiceStateKey = MemberKey

// PresenceKey is a composite cache key for presences (same shape as MemberKey).
type PresenceKey = MemberKey

// ── entry ─────────────────────────────────────────────────────────────────────

// entry holds a value plus the metadata required for eviction decisions.
// All fields are owned by the containing BaseStore and must not be accessed
// without holding the store lock.
type entry[V any] struct {
	value      V
	insertedAt time.Time
	accessedAt time.Time
	hitCount   uint64
	expiresAt  time.Time // zero = no TTL
	sizeBytes  int64     // 0 when TrackBytes=false
}

// ── BaseStore ─────────────────────────────────────────────────────────────────

// BaseStore is the generic, thread-safe foundation used by every entity store.
// It provides TTL, per-item eviction, and background sweepers.
//
// All public methods are safe for concurrent use.
type BaseStore[K comparable, V any] struct {
	mu         sync.RWMutex
	items      *collection.Collection[K, *entry[V]]
	options    StoreOptions
	totalBytes atomic.Int64 // non-zero only when options.TrackBytes is true

	stopCh    chan struct{}
	closeOnce sync.Once
	wg        sync.WaitGroup
}

// NewBaseStore creates a store with the given options and starts any configured
// sweepers. Call Close to stop background goroutines.
func NewBaseStore[K comparable, V any](opts StoreOptions) *BaseStore[K, V] {
	s := &BaseStore[K, V]{
		items:   collection.New[K, *entry[V]](),
		options: opts,
		stopCh:  make(chan struct{}),
	}
	for _, sw := range opts.Sweepers {
		s.startSweeper(sw)
	}
	return s
}

// IsEnabled reports whether this store accepts writes.
// A disabled store silently drops all mutations and returns zero values on reads.
func (s *BaseStore[K, V]) IsEnabled() bool {
	return !s.options.Disabled
}

// Get returns the value for key and updates access metadata (accessedAt, hitCount).
// Returns the zero value and false if the key is absent or its TTL has elapsed.
// Expired entries are removed lazily on access.
func (s *BaseStore[K, V]) Get(key K) (V, bool) {
	if s.options.Disabled {
		var zero V
		return zero, false
	}

	// Full write lock: we may delete an expired entry or update access metadata.
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.items.Get(key)
	if !ok {
		var zero V
		return zero, false
	}

	now := time.Now()
	if !e.expiresAt.IsZero() && now.After(e.expiresAt) {
		s.items.Delete(key)
		s.totalBytes.Add(-e.sizeBytes)
		var zero V
		return zero, false
	}

	e.accessedAt = now
	e.hitCount++
	if s.options.RefreshOnAccess && s.options.TTL > 0 {
		e.expiresAt = now.Add(s.options.TTL)
	}

	return e.value, true
}

// Set inserts or replaces the value for key.
// When replacing an existing entry, insertedAt and hitCount are preserved so
// that eviction metadata is not reset by routine update events.
// If MaxItems is exceeded after insertion, one entry is evicted.
func (s *BaseStore[K, V]) Set(key K, value V) {
	if s.options.Disabled {
		return
	}

	now := time.Now()
	var sz int64
	if s.options.TrackBytes {
		sz = estimateSize(value)
	}
	e := &entry[V]{
		value:      value,
		insertedAt: now,
		accessedAt: now,
		sizeBytes:  sz,
	}
	if s.options.TTL > 0 {
		e.expiresAt = now.Add(s.options.TTL)
	}

	s.mu.Lock()
	if old, ok := s.items.Get(key); ok {
		e.insertedAt = old.insertedAt
		e.hitCount = old.hitCount
		if s.options.TrackBytes {
			s.totalBytes.Add(-old.sizeBytes)
		}
	}
	s.items.Set(key, e)
	if s.options.TrackBytes {
		s.totalBytes.Add(sz)
	}
	over := s.options.MaxItems > 0 && s.items.Len() > s.options.MaxItems
	s.mu.Unlock()

	if over {
		s.evictOne()
	}
}

// Delete removes the entry for key. Returns true if the entry was present.
func (s *BaseStore[K, V]) Delete(key K) bool {
	if s.options.Disabled {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.items.Get(key); ok && s.options.TrackBytes {
		s.totalBytes.Add(-e.sizeBytes)
	}
	return s.items.Delete(key)
}

// Has reports whether key is present. Does not update access metadata.
func (s *BaseStore[K, V]) Has(key K) bool {
	if s.options.Disabled {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items.Get(key)
	return ok
}

// Len returns the current number of entries. May include entries whose TTL
// has elapsed but have not yet been accessed (lazy cleanup).
func (s *BaseStore[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items.Len()
}

// Clear removes all entries.
func (s *BaseStore[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.options.TrackBytes {
		s.totalBytes.Store(0)
	}
	s.items.Clear()
}

// Bytes returns the estimated total payload size in bytes.
// Always 0 when TrackBytes is false in StoreOptions.
func (s *BaseStore[K, V]) Bytes() int64 { return s.totalBytes.Load() }

// TrimTo evicts entries (by EvictionStrategy) until Len() <= n.
// Safe for concurrent use; must NOT be called while s.mu is held.
func (s *BaseStore[K, V]) TrimTo(n int) { s.evictToCount(n) }

// All returns a snapshot Collection of all live (non-expired) values.
// The snapshot is a copy; mutations to it do not affect the store.
func (s *BaseStore[K, V]) All() *collection.Collection[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := collection.NewWithCapacity[K, V](s.items.Len())
	now := time.Now()
	s.items.Each(func(k K, e *entry[V]) {
		if e.expiresAt.IsZero() || now.Before(e.expiresAt) {
			out.Set(k, e.value)
		}
	})
	return out
}

// ReplaceKeys atomically deletes all keys in del and inserts all entries in add
// under a single lock acquisition. Used by guild-indexed SetAll operations
// (EmojiStore, StickerStore, etc.) to avoid a visible half-replaced window.
//
// If the resulting store size exceeds MaxItems, excess entries are evicted
// after the lock is released (same soft-ceiling semantics as Set).
func (s *BaseStore[K, V]) ReplaceKeys(del []K, add map[K]V) {
	if s.options.Disabled {
		return
	}

	now := time.Now()
	s.mu.Lock()
	for _, k := range del {
		if e, ok := s.items.Get(k); ok && s.options.TrackBytes {
			s.totalBytes.Add(-e.sizeBytes)
		}
		s.items.Delete(k)
	}
	for k, v := range add {
		var sz int64
		if s.options.TrackBytes {
			sz = estimateSize(v)
		}
		e := &entry[V]{
			value:      v,
			insertedAt: now,
			accessedAt: now,
			sizeBytes:  sz,
		}
		if old, ok := s.items.Get(k); ok {
			e.insertedAt = old.insertedAt
			e.hitCount = old.hitCount
			if s.options.TrackBytes {
				s.totalBytes.Add(-old.sizeBytes)
			}
		}
		if s.options.TrackBytes {
			s.totalBytes.Add(sz)
		}
		if s.options.TTL > 0 {
			e.expiresAt = now.Add(s.options.TTL)
		}
		s.items.Set(k, e)
	}
	over := s.options.MaxItems > 0 && s.items.Len() > s.options.MaxItems
	s.mu.Unlock()

	if over {
		s.evictToCount(s.options.MaxItems)
	}
}

// SweepExpired removes entries whose TTL has elapsed. When unusedWindow > 0,
// entries not accessed within that window are also removed.
// Called by MemoryCache's background cleaner; safe for concurrent use.
func (s *BaseStore[K, V]) SweepExpired(unusedWindow time.Duration) {
	if s.options.Disabled {
		return
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items.Sweep(func(e *entry[V]) bool {
		expired := !e.expiresAt.IsZero() && now.After(e.expiresAt)
		idle := unusedWindow > 0 && now.Sub(e.accessedAt) > unusedWindow
		if (expired || idle) && s.options.TrackBytes {
			s.totalBytes.Add(-e.sizeBytes)
		}
		return expired || idle
	})
}

// Close stops all sweeper goroutines and waits for them to exit.
// Safe to call more than once; subsequent calls are no-ops.
func (s *BaseStore[K, V]) Close() {
	s.closeOnce.Do(func() { close(s.stopCh) })
	s.wg.Wait()
}

// ── Internal helpers ──────────────────────────────────────────────────────────

// evictOne selects and removes one entry according to EvictionStrategy.
// Acquires its own lock; must NOT be called while s.mu is held.
func (s *BaseStore[K, V]) evictOne() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var victimKey K
	var victimEntry *entry[V]
	first := true

	s.items.Each(func(k K, e *entry[V]) {
		if first {
			victimKey, victimEntry = k, e
			first = false
			return
		}
		if s.isBetterVictim(e, victimEntry) {
			victimKey, victimEntry = k, e
		}
	})

	if !first {
		if s.options.TrackBytes && victimEntry != nil {
			s.totalBytes.Add(-victimEntry.sizeBytes)
		}
		s.items.Delete(victimKey)
	}
}

// evictToCount removes entries until Len() <= target in a single lock
// acquisition (O(n log n)), avoiding the O(n²) cost of repeated evictOne calls.
// Must NOT be called while s.mu is held.
func (s *BaseStore[K, V]) evictToCount(target int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := s.items.Len()
	if n <= target {
		return
	}
	toEvict := n - target

	type kv struct {
		key K
		e   *entry[V]
	}
	all := make([]kv, 0, n)
	s.items.Each(func(k K, e *entry[V]) {
		all = append(all, kv{k, e})
	})

	sort.Slice(all, func(i, j int) bool {
		return s.isBetterVictim(all[i].e, all[j].e)
	})

	for i := 0; i < toEvict && i < len(all); i++ {
		if s.options.TrackBytes {
			s.totalBytes.Add(-all[i].e.sizeBytes)
		}
		s.items.Delete(all[i].key)
	}
}

// isBetterVictim reports whether candidate should be evicted before current.
func (s *BaseStore[K, V]) isBetterVictim(candidate, current *entry[V]) bool {
	switch s.options.EvictionStrategy {
	case EvictLFU:
		return candidate.hitCount < current.hitCount
	case EvictFIFO:
		return candidate.insertedAt.Before(current.insertedAt)
	default: // EvictLRU
		return candidate.accessedAt.Before(current.accessedAt)
	}
}

// startSweeper launches sw in its own goroutine.
func (s *BaseStore[K, V]) startSweeper(sw Sweeper) {
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

// runSweeper removes entries whose value matches sw.Filter. Acquires the write
// lock for the full iteration.
func (s *BaseStore[K, V]) runSweeper(sw Sweeper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items.Sweep(func(e *entry[V]) bool {
		remove := sw.Filter(e.value)
		if remove && s.options.TrackBytes {
			s.totalBytes.Add(-e.sizeBytes)
		}
		return remove
	})
}
