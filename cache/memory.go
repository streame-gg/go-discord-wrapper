package cache

import (
	"sort"
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache/stores"
)

// MemoryCache is an in-process implementation of [Cache].
//
// Create with [NewMemoryCache] and wire it to the gateway client via
// options.WithCache. Always call [MemoryCache.Close] on shutdown.
type MemoryCache struct {
	opts            Options
	guilds          *stores.MemGuildStore
	channels        *stores.MemChannelStore
	users           *stores.MemUserStore
	members         *stores.MemMemberStore
	roles           *stores.MemRoleStore
	messages        *stores.MemMessageStore
	voiceStates     *stores.MemVoiceStateStore
	soundboard      *stores.MemSoundboardStore
	scheduledEvents *stores.MemScheduledEventStore
	stageInstances  *stores.MemStageInstanceStore
	emojis          *stores.MemEmojiStore
	stickers        *stores.MemStickerStore
	presences       *stores.MemPresenceStore
	bans            *stores.MemBanStore
	autoModRules    *stores.MemAutoModerationRuleStore
	invites         *stores.MemInviteStore

	stopOnce sync.Once
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewMemoryCache creates and starts a MemoryCache. Missing options receive
// sensible defaults:
//   - SweepInterval: 30 s
//   - Messages.TTL: inherits global TTL
//   - OnOverflow.Target: CategoryAll
func NewMemoryCache(opts Options) *MemoryCache {
	if opts.SweepInterval <= 0 {
		opts.SweepInterval = 30 * time.Second
	}
	if opts.Messages.MaxPerChannel < 0 {
		opts.Messages.MaxPerChannel = 100
	}
	if opts.Messages.TTL == 0 {
		opts.Messages.TTL = opts.TTL
	}
	if opts.OnOverflow.Target == 0 {
		opts.OnOverflow.Target = CategoryAll
	}

	strategy := clearByToStrategy(opts.OnOverflow.ClearBy)
	trackBytes := opts.Limits.MaxSizeMB > 0

	sc := func(maxItems int) stores.StoreOptions {
		return stores.StoreOptions{
			MaxItems:         maxItems,
			TTL:              opts.TTL,
			EvictionStrategy: strategy,
			TrackBytes:       trackBytes,
		}
	}

	c := &MemoryCache{
		opts:     opts,
		guilds:   stores.NewGuildStore(sc(opts.Limits.MaxGuilds)),
		channels: stores.NewChannelStore(sc(opts.Limits.MaxChannels)),
		users:    stores.NewUserStore(sc(opts.Limits.MaxUsers)),
		members:  stores.NewMemberStore(sc(opts.Limits.MaxMembers)),
		roles:    stores.NewRoleStore(sc(opts.Limits.MaxRoles)),
		messages: stores.NewMessageStore(stores.StoreOptions{
			MaxItems:         opts.Messages.MaxPerChannel,
			MaxTotal:         opts.Limits.MaxMessages,
			TTL:              opts.Messages.TTL,
			EvictionStrategy: strategy,
		}),
		voiceStates:     stores.NewVoiceStateStore(sc(0)),
		soundboard:      stores.NewSoundboardStore(sc(0)),
		scheduledEvents: stores.NewScheduledEventStore(sc(0)),
		stageInstances:  stores.NewStageInstanceStore(sc(0)),
		emojis:          stores.NewEmojiStore(sc(0)),
		stickers:        stores.NewStickerStore(sc(opts.Limits.MaxStickers)),
		presences:       stores.NewPresenceStore(stores.PresenceDefaults()),
		bans:            stores.NewBanStore(sc(0)),
		autoModRules:    stores.NewAutoModerationRuleStore(sc(0)),
		invites:         stores.NewInviteStore(sc(0)),
		stopCh:          make(chan struct{}),
	}

	c.wg.Add(1)
	go c.cleaner()
	return c
}

// clearByToStrategy maps cache.ClearBy to stores.EvictionStrategy.
func clearByToStrategy(by ClearBy) stores.EvictionStrategy {
	switch by {
	case ClearByAge:
		return stores.EvictFIFO
	case ClearByFrequency:
		return stores.EvictLFU
	default:
		return stores.EvictLRU
	}
}

func (c *MemoryCache) Guilds() GuildStore                    { return c.guilds }
func (c *MemoryCache) Channels() ChannelStore                { return c.channels }
func (c *MemoryCache) Users() UserStore                      { return c.users }
func (c *MemoryCache) Members() MemberStore                  { return c.members }
func (c *MemoryCache) Roles() RoleStore                      { return c.roles }
func (c *MemoryCache) Messages() MessageStore                { return c.messages }
func (c *MemoryCache) VoiceStates() VoiceStateStore          { return c.voiceStates }
func (c *MemoryCache) Soundboard() SoundboardStore           { return c.soundboard }
func (c *MemoryCache) ScheduledEvents() ScheduledEventStore  { return c.scheduledEvents }
func (c *MemoryCache) StageInstances() StageInstanceStore    { return c.stageInstances }
func (c *MemoryCache) Emojis() EmojiStore                    { return c.emojis }
func (c *MemoryCache) Stickers() StickerStore                { return c.stickers }
func (c *MemoryCache) Presences() PresenceStore              { return c.presences }
func (c *MemoryCache) Bans() BanStore                        { return c.bans }
func (c *MemoryCache) AutoModRules() AutoModerationRuleStore { return c.autoModRules }
func (c *MemoryCache) Invites() InviteStore                  { return c.invites }

// Close stops the background sweeper and all store goroutines. Safe to call
// multiple times.
func (c *MemoryCache) Close() error {
	c.stopOnce.Do(func() {
		close(c.stopCh)
		c.wg.Wait()
		c.guilds.Close()
		c.channels.Close()
		c.users.Close()
		c.members.Close()
		c.roles.Close()
		c.messages.Close()
		c.voiceStates.Close()
		c.soundboard.Close()
		c.scheduledEvents.Close()
		c.stageInstances.Close()
		c.emojis.Close()
		c.stickers.Close()
		c.presences.Close()
		c.bans.Close()
		c.autoModRules.Close()
		c.invites.Close()
	})
	return nil
}

func (c *MemoryCache) cleaner() {
	defer c.wg.Done()
	t := time.NewTicker(c.opts.SweepInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			c.sweep()
		case <-c.stopCh:
			return
		}
	}
}

func (c *MemoryCache) sweep() {
	uw := time.Duration(0)
	if c.opts.EvictBehavior == EvictUnused {
		uw = c.opts.UnusedWindow
	}
	c.guilds.SweepExpired(uw)
	c.channels.SweepExpired(uw)
	c.users.SweepExpired(uw)
	c.members.SweepExpired(uw)
	c.roles.SweepExpired(uw)
	c.messages.SweepExpired(uw)
	c.voiceStates.SweepExpired(uw)
	c.soundboard.SweepExpired(uw)
	c.scheduledEvents.SweepExpired(uw)
	c.stageInstances.SweepExpired(uw)
	c.emojis.SweepExpired(uw)
	c.stickers.SweepExpired(uw)
	c.presences.SweepExpired(uw)
	c.bans.SweepExpired(uw)
	c.autoModRules.SweepExpired(uw)
	c.invites.SweepExpired(uw)

	c.enforceGlobalLimits()
}

// enforceGlobalLimits applies overflow eviction for MaxEntries and MaxSizeMB.
// Called after every sweep.
func (c *MemoryCache) enforceGlobalLimits() {
	lim := c.opts.Limits
	pol := c.opts.OnOverflow

	if lim.MaxEntries > 0 {
		total := c.countTargeted(pol.Target)
		if total > lim.MaxEntries {
			c.evictByCount(total-lim.MaxEntries, pol.Target, pol.ClearBy)
		}
	}

	if lim.MaxSizeMB > 0 {
		maxBytes := int64(lim.MaxSizeMB * 1024 * 1024)
		totalBytes := c.bytesTargeted(pol.Target)
		if totalBytes > maxBytes {
			c.evictByBytes(totalBytes-maxBytes, pol.Target, pol.ClearBy)
		}
	}
}

// countTargeted returns the total entry count across targeted stores.
func (c *MemoryCache) countTargeted(target OverflowCategory) int {
	n := 0
	if target&CategoryGuilds != 0 {
		n += c.guilds.Size()
	}
	if target&CategoryChannels != 0 {
		n += c.channels.Size()
	}
	if target&CategoryUsers != 0 {
		n += c.users.Size()
	}
	if target&CategoryMembers != 0 {
		n += c.members.Size()
	}
	if target&CategoryRoles != 0 {
		n += c.roles.Size()
	}
	if target&CategoryMessages != 0 {
		n += c.messages.Size()
	}
	if target&CategoryVoiceStates != 0 {
		n += c.voiceStates.Size()
	}
	if target&CategorySoundboard != 0 {
		n += c.soundboard.Size()
	}
	if target&CategoryScheduledEvents != 0 {
		n += c.scheduledEvents.Size()
	}
	if target&CategoryStageInstances != 0 {
		n += c.stageInstances.Size()
	}
	if target&CategoryEmojis != 0 {
		n += c.emojis.Size()
	}
	if target&CategoryStickers != 0 {
		n += c.stickers.Size()
	}
	if target&CategoryPresences != 0 {
		n += c.presences.Size()
	}
	if target&CategoryBans != 0 {
		n += c.bans.Size()
	}
	if target&CategoryAutoModRules != 0 {
		n += c.autoModRules.Size()
	}
	if target&CategoryInvites != 0 {
		n += c.invites.Size()
	}
	return n
}

// bytesTargeted returns the total byte estimate across targeted stores.
func (c *MemoryCache) bytesTargeted(target OverflowCategory) int64 {
	var b int64
	if target&CategoryGuilds != 0 {
		b += c.guilds.Bytes()
	}
	if target&CategoryChannels != 0 {
		b += c.channels.Bytes()
	}
	if target&CategoryUsers != 0 {
		b += c.users.Bytes()
	}
	if target&CategoryMembers != 0 {
		b += c.members.Bytes()
	}
	if target&CategoryRoles != 0 {
		b += c.roles.Bytes()
	}
	if target&CategoryVoiceStates != 0 {
		b += c.voiceStates.Bytes()
	}
	if target&CategorySoundboard != 0 {
		b += c.soundboard.Bytes()
	}
	if target&CategoryScheduledEvents != 0 {
		b += c.scheduledEvents.Bytes()
	}
	if target&CategoryStageInstances != 0 {
		b += c.stageInstances.Bytes()
	}
	if target&CategoryEmojis != 0 {
		b += c.emojis.Bytes()
	}
	if target&CategoryStickers != 0 {
		b += c.stickers.Bytes()
	}
	if target&CategoryPresences != 0 {
		b += c.presences.Bytes()
	}
	if target&CategoryBans != 0 {
		b += c.bans.Bytes()
	}
	if target&CategoryAutoModRules != 0 {
		b += c.autoModRules.Bytes()
	}
	if target&CategoryInvites != 0 {
		b += c.invites.Bytes()
	}
	return b
}

type storeEntry struct {
	size   int
	bytes  int64
	trimTo func(n int)
}

func (c *MemoryCache) storesForTarget(target OverflowCategory) []storeEntry {
	out := make([]storeEntry, 0, 13)
	add := func(sz int, by int64, f func(int)) {
		if sz > 0 {
			out = append(out, storeEntry{sz, by, f})
		}
	}
	if target&CategoryGuilds != 0 {
		add(c.guilds.Size(), c.guilds.Bytes(), c.guilds.TrimTo)
	}
	if target&CategoryChannels != 0 {
		add(c.channels.Size(), c.channels.Bytes(), c.channels.TrimTo)
	}
	if target&CategoryUsers != 0 {
		add(c.users.Size(), c.users.Bytes(), c.users.TrimTo)
	}
	if target&CategoryMembers != 0 {
		add(c.members.Size(), c.members.Bytes(), c.members.TrimTo)
	}
	if target&CategoryRoles != 0 {
		add(c.roles.Size(), c.roles.Bytes(), c.roles.TrimTo)
	}
	if target&CategoryMessages != 0 {
		n := c.messages.Size()
		if n > 0 {
			out = append(out, storeEntry{n, 0, c.messages.TrimToTotal})
		}
	}
	if target&CategoryVoiceStates != 0 {
		add(c.voiceStates.Size(), c.voiceStates.Bytes(), c.voiceStates.TrimTo)
	}
	if target&CategorySoundboard != 0 {
		add(c.soundboard.Size(), c.soundboard.Bytes(), c.soundboard.TrimTo)
	}
	if target&CategoryScheduledEvents != 0 {
		add(c.scheduledEvents.Size(), c.scheduledEvents.Bytes(), c.scheduledEvents.TrimTo)
	}
	if target&CategoryStageInstances != 0 {
		add(c.stageInstances.Size(), c.stageInstances.Bytes(), c.stageInstances.TrimTo)
	}
	if target&CategoryEmojis != 0 {
		add(c.emojis.Size(), c.emojis.Bytes(), c.emojis.TrimTo)
	}
	if target&CategoryStickers != 0 {
		add(c.stickers.Size(), c.stickers.Bytes(), c.stickers.TrimTo)
	}
	if target&CategoryPresences != 0 {
		add(c.presences.Size(), c.presences.Bytes(), c.presences.TrimTo)
	}
	if target&CategoryBans != 0 {
		add(c.bans.Size(), c.bans.Bytes(), c.bans.TrimTo)
	}
	if target&CategoryAutoModRules != 0 {
		add(c.autoModRules.Size(), c.autoModRules.Bytes(), c.autoModRules.TrimTo)
	}
	if target&CategoryInvites != 0 {
		add(c.invites.Size(), c.invites.Bytes(), c.invites.TrimTo)
	}
	return out
}

// evictByCount removes `need` entries proportionally from targeted stores
// using the Largest-Remainder Method to avoid over-eviction.
func (c *MemoryCache) evictByCount(need int, target OverflowCategory, by ClearBy) {
	_ = by // eviction order within a store is set at construction time

	entries := c.storesForTarget(target)
	total := 0
	for _, e := range entries {
		total += e.size
	}
	if total == 0 {
		return
	}

	type slot struct {
		idx       int
		floor     int
		remainder float64
	}
	slots := make([]slot, 0, len(entries))
	floorSum := 0
	for i, e := range entries {
		real := float64(e.size*need) / float64(total)
		f := int(real)
		slots = append(slots, slot{i, f, real - float64(f)})
		floorSum += f
	}
	extra := need - floorSum
	if extra > 0 {
		sort.Slice(slots, func(a, b int) bool { return slots[a].remainder > slots[b].remainder })
		for i := 0; i < extra && i < len(slots); i++ {
			slots[i].floor++
		}
	}
	for _, sl := range slots {
		if sl.floor > 0 {
			e := entries[sl.idx]
			target := e.size - sl.floor
			if target < 0 {
				target = 0
			}
			e.trimTo(target)
		}
	}
}

// evictByBytes removes entries from targeted stores until `need` bytes have
// been freed. Stores are processed in descending average-entry-size order so
// that each eviction frees the maximum bytes possible.
func (c *MemoryCache) evictByBytes(need int64, target OverflowCategory, by ClearBy) {
	_ = by

	entries := c.storesForTarget(target)

	type storeInfo struct {
		e       storeEntry
		avgSize int64
	}
	infos := make([]storeInfo, 0, len(entries))
	for _, e := range entries {
		if e.bytes == 0 || e.size == 0 {
			continue
		}
		avg := e.bytes / int64(e.size)
		if avg <= 0 {
			avg = 1
		}
		infos = append(infos, storeInfo{e, avg})
	}
	// Largest-entry stores first: one eviction there frees the most bytes.
	sort.Slice(infos, func(a, b int) bool {
		return infos[a].avgSize > infos[b].avgSize
	})

	for _, info := range infos {
		if need <= 0 {
			break
		}
		// Ceiling division: evict just enough entries to cover this store's share.
		toEvict := int((need + info.avgSize - 1) / info.avgSize)
		if toEvict > info.e.size {
			toEvict = info.e.size
		}
		if toEvict <= 0 {
			continue
		}
		info.e.trimTo(info.e.size - toEvict)
		need -= int64(toEvict) * info.avgSize
	}
}
