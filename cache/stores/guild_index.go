package stores

import (
	"sync"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── guildIndex ────────────────────────────────────────────────────────────────

// guildIndex tracks guild→entity and entity→guild associations for flat
// entity stores whose cache key is the entity ID (RoleStore, EmojiStore, etc.).
// It maintains a bidirectional index so that both entity-level deletions
// (Delete(entityID)) and guild-level deletions (DeleteGuild(guildID)) are O(1).
//
// All methods are safe for concurrent use.
type guildIndex struct {
	mu      sync.RWMutex
	byGuild map[discord.Snowflake]map[discord.Snowflake]struct{} // guildID → entityIDs
	toGuild map[discord.Snowflake]discord.Snowflake              // entityID → guildID
}

func newGuildIndex() *guildIndex {
	return &guildIndex{
		byGuild: make(map[discord.Snowflake]map[discord.Snowflake]struct{}),
		toGuild: make(map[discord.Snowflake]discord.Snowflake),
	}
}

// add records that entityID belongs to guildID.
// If entityID was previously in a different guild (shouldn't happen in Discord,
// but handled defensively), it is moved.
func (idx *guildIndex) add(guildID, entityID discord.Snowflake) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if oldGuild, ok := idx.toGuild[entityID]; ok && oldGuild != guildID {
		if set := idx.byGuild[oldGuild]; set != nil {
			delete(set, entityID)
			if len(set) == 0 {
				delete(idx.byGuild, oldGuild)
			}
		}
	}

	if idx.byGuild[guildID] == nil {
		idx.byGuild[guildID] = make(map[discord.Snowflake]struct{})
	}
	idx.byGuild[guildID][entityID] = struct{}{}
	idx.toGuild[entityID] = guildID
}

// removeEntity removes entityID from the index.
// Returns the guildID it belonged to and true if found.
func (idx *guildIndex) removeEntity(entityID discord.Snowflake) (guildID discord.Snowflake, found bool) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	guildID, found = idx.toGuild[entityID]
	if !found {
		return
	}
	delete(idx.toGuild, entityID)
	if set := idx.byGuild[guildID]; set != nil {
		delete(set, entityID)
		if len(set) == 0 {
			delete(idx.byGuild, guildID)
		}
	}
	return
}

// removeGuild removes all entities for guildID and returns their IDs.
func (idx *guildIndex) removeGuild(guildID discord.Snowflake) []discord.Snowflake {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	set := idx.byGuild[guildID]
	if len(set) == 0 {
		delete(idx.byGuild, guildID)
		return nil
	}
	ids := make([]discord.Snowflake, 0, len(set))
	for id := range set {
		ids = append(ids, id)
		delete(idx.toGuild, id)
	}
	delete(idx.byGuild, guildID)
	return ids
}

// forGuild returns a snapshot of entity IDs for guildID.
func (idx *guildIndex) forGuild(guildID discord.Snowflake) []discord.Snowflake {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	set := idx.byGuild[guildID]
	if len(set) == 0 {
		return nil
	}
	ids := make([]discord.Snowflake, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}

// setForGuild atomically replaces all entities for guildID with newIDs.
// Returns the entity IDs that were removed (present in old set but not in newIDs).
// Pass the returned slice as the del argument to BaseStore.ReplaceKeys.
func (idx *guildIndex) setForGuild(guildID discord.Snowflake, newIDs []discord.Snowflake) (toDelete []discord.Snowflake) {
	newSet := make(map[discord.Snowflake]struct{}, len(newIDs))
	for _, id := range newIDs {
		newSet[id] = struct{}{}
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	oldSet := idx.byGuild[guildID]
	for id := range oldSet {
		if _, kept := newSet[id]; !kept {
			toDelete = append(toDelete, id)
			delete(idx.toGuild, id)
		}
	}
	for id := range newSet {
		idx.toGuild[id] = guildID
	}
	if len(newSet) == 0 {
		delete(idx.byGuild, guildID)
	} else {
		idx.byGuild[guildID] = newSet
	}
	return
}

// ── compositeGuildIndex ───────────────────────────────────────────────────────

// compositeGuildIndex tracks guildID → secondary-key associations for
// composite-key stores (MemberStore, VoiceStateStore, PresenceStore).
// Unlike guildIndex, no reverse lookup is needed because both guildID
// and secondaryID are always available at the call site.
//
// All methods are safe for concurrent use.
type compositeGuildIndex struct {
	mu      sync.RWMutex
	byGuild map[discord.Snowflake]map[discord.Snowflake]struct{} // guildID → secondaryIDs
}

func newCompositeGuildIndex() *compositeGuildIndex {
	return &compositeGuildIndex{
		byGuild: make(map[discord.Snowflake]map[discord.Snowflake]struct{}),
	}
}

func (idx *compositeGuildIndex) add(guildID, secondaryID discord.Snowflake) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.byGuild[guildID] == nil {
		idx.byGuild[guildID] = make(map[discord.Snowflake]struct{})
	}
	idx.byGuild[guildID][secondaryID] = struct{}{}
}

func (idx *compositeGuildIndex) remove(guildID, secondaryID discord.Snowflake) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if set := idx.byGuild[guildID]; set != nil {
		delete(set, secondaryID)
		if len(set) == 0 {
			delete(idx.byGuild, guildID)
		}
	}
}

// removeGuild removes all secondary IDs for guildID and returns them.
func (idx *compositeGuildIndex) removeGuild(guildID discord.Snowflake) []discord.Snowflake {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	set := idx.byGuild[guildID]
	if len(set) == 0 {
		delete(idx.byGuild, guildID)
		return nil
	}
	ids := make([]discord.Snowflake, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	delete(idx.byGuild, guildID)
	return ids
}

// forGuild returns a snapshot of secondary IDs for guildID.
func (idx *compositeGuildIndex) forGuild(guildID discord.Snowflake) []discord.Snowflake {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	set := idx.byGuild[guildID]
	if len(set) == 0 {
		return nil
	}
	ids := make([]discord.Snowflake, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}
