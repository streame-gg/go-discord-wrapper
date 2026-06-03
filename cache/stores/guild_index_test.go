package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type guildIndexSuite struct{ suite.Suite }

func TestGuildIndexSuite(t *testing.T) { suite.Run(t, new(guildIndexSuite)) }

// ── guildIndex ────────────────────────────────────────────────────────────────

func (su *guildIndexSuite) TestGuildIndex_AddMovesEntityBetweenGuilds() {
	t := su.T()
	idx := newGuildIndex()
	e := discord.Snowflake(1)

	idx.add(gA, e)
	assert.Equal(t, []discord.Snowflake{e}, idx.forGuild(gA))

	// Re-adding under a different guild must move it (and clear the old guild
	// once its set is empty).
	idx.add(gB, e)
	assert.Nil(t, idx.forGuild(gA))
	assert.Equal(t, []discord.Snowflake{e}, idx.forGuild(gB))
}

func (su *guildIndexSuite) TestGuildIndex_RemoveEntity() {
	t := su.T()
	idx := newGuildIndex()

	// Not found → found == false.
	_, found := idx.removeEntity(discord.Snowflake(999))
	assert.False(t, found)

	e := discord.Snowflake(1)
	idx.add(gA, e)
	g, found := idx.removeEntity(e)
	assert.True(t, found)
	assert.Equal(t, gA, g)
	// Set emptied → guild removed from index.
	assert.Nil(t, idx.forGuild(gA))
}

func (su *guildIndexSuite) TestGuildIndex_RemoveGuild_Empty() {
	t := su.T()
	idx := newGuildIndex()
	assert.Nil(t, idx.removeGuild(gA)) // no entries → nil
}

func (su *guildIndexSuite) TestGuildIndex_SetForGuild_EmptyClearsGuild() {
	t := su.T()
	idx := newGuildIndex()
	idx.add(gA, discord.Snowflake(1))
	idx.add(gA, discord.Snowflake(2))

	// Replacing with an empty set deletes the guild entry and reports both as
	// removed.
	toDelete := idx.setForGuild(gA, nil)
	assert.Len(t, toDelete, 2)
	assert.Nil(t, idx.forGuild(gA))
}

func (su *guildIndexSuite) TestGuildIndex_SetForGuild_PartialOverlap() {
	t := su.T()
	idx := newGuildIndex()
	idx.add(gA, discord.Snowflake(1))
	idx.add(gA, discord.Snowflake(2))

	toDelete := idx.setForGuild(gA, []discord.Snowflake{2, 3})
	assert.Equal(t, []discord.Snowflake{1}, toDelete) // 1 dropped; 2 kept; 3 added
	assert.Len(t, idx.forGuild(gA), 2)
}

// ── compositeGuildIndex ───────────────────────────────────────────────────────

func (su *guildIndexSuite) TestCompositeGuildIndex() {
	t := su.T()
	idx := newCompositeGuildIndex()

	assert.Nil(t, idx.forGuild(gA))
	assert.Nil(t, idx.removeGuild(gA)) // empty → nil

	idx.add(gA, discord.Snowflake(1))
	idx.add(gA, discord.Snowflake(2))
	assert.Len(t, idx.forGuild(gA), 2)

	// remove down to empty → guild dropped.
	idx.remove(gA, discord.Snowflake(1))
	idx.remove(gA, discord.Snowflake(2))
	assert.Nil(t, idx.forGuild(gA))

	// remove on a missing guild is a no-op.
	idx.remove(discord.Snowflake(999), discord.Snowflake(1))

	idx.add(gB, discord.Snowflake(3))
	ids := idx.removeGuild(gB)
	assert.Equal(t, []discord.Snowflake{discord.Snowflake(3)}, ids)
}

// ── msgRing internals ─────────────────────────────────────────────────────────

func (su *guildIndexSuite) TestMsgRing_EvictOneEmptyIsNoop() {
	r := newMsgRing(2)
	r.evictOne(EvictLRU) // must not panic on an empty ring
	assert.Empty(su.T(), r.items)
}
