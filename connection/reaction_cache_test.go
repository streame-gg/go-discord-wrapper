package connection

// Tests for Bug 76/77: reaction helper functions and gateway events must update
// the Reactions field on cached messages.

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── reactionEmojiMatches ──────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestReactionEmojiMatches_CustomEmoji_ById() {
	t := cs.T()
	a := discord.Emoji{ID: 111}
	b := discord.Emoji{ID: 111}
	assert.True(t, reactionEmojiMatches(a, b))
}

func (cs *ConnectionSuite) TestReactionEmojiMatches_CustomEmoji_DifferentId() {
	t := cs.T()
	a := discord.Emoji{ID: 111}
	b := discord.Emoji{ID: 222}
	assert.False(t, reactionEmojiMatches(a, b))
}

func (cs *ConnectionSuite) TestReactionEmojiMatches_UnicodeEmoji_ByName() {
	t := cs.T()
	a := discord.Emoji{Name: "👍"}
	b := discord.Emoji{Name: "👍"}
	assert.True(t, reactionEmojiMatches(a, b))
}

func (cs *ConnectionSuite) TestReactionEmojiMatches_UnicodeEmoji_DifferentName() {
	t := cs.T()
	a := discord.Emoji{Name: "👍"}
	b := discord.Emoji{Name: "👎"}
	assert.False(t, reactionEmojiMatches(a, b))
}

// ── appendOrIncrementReaction ─────────────────────────────────────────────────

func (cs *ConnectionSuite) TestAppendOrIncrementReaction_NilSlice() {
	t := cs.T()
	emoji := discord.Emoji{Name: "👍"}
	result := appendOrIncrementReaction(nil, emoji)

	require.Len(t, result, 1)
	assert.Equal(t, 1, result[0].Count)
	assert.Equal(t, emoji, result[0].Emoji)
}

func (cs *ConnectionSuite) TestAppendOrIncrementReaction_NewEmoji() {
	t := cs.T()
	existing := &[]discord.Reaction{
		{Count: 3, Emoji: discord.Emoji{Name: "👎"}},
	}
	emoji := discord.Emoji{Name: "👍"}
	result := appendOrIncrementReaction(existing, emoji)

	require.Len(t, result, 2)
	assert.Equal(t, 3, result[0].Count, "existing reaction count must be unchanged")
	assert.Equal(t, 1, result[1].Count, "new reaction must start at count 1")
}

func (cs *ConnectionSuite) TestAppendOrIncrementReaction_ExistingEmoji() {
	t := cs.T()
	existing := &[]discord.Reaction{
		{Count: 5, Emoji: discord.Emoji{Name: "👍"}},
	}
	result := appendOrIncrementReaction(existing, discord.Emoji{Name: "👍"})

	require.Len(t, result, 1)
	assert.Equal(t, 6, result[0].Count, "existing reaction count must be incremented")
}

func (cs *ConnectionSuite) TestAppendOrIncrementReaction_DoesNotMutateOriginal() {
	t := cs.T()
	orig := []discord.Reaction{{Count: 2, Emoji: discord.Emoji{Name: "👍"}}}
	existing := &orig
	appendOrIncrementReaction(existing, discord.Emoji{Name: "👍"})
	assert.Equal(t, 2, orig[0].Count, "original slice must not be mutated")
}

// ── decrementOrRemoveReaction ─────────────────────────────────────────────────

func (cs *ConnectionSuite) TestDecrementOrRemoveReaction_NilSlice() {
	t := cs.T()
	result := decrementOrRemoveReaction(nil, discord.Emoji{Name: "👍"})
	assert.Nil(t, result)
}

func (cs *ConnectionSuite) TestDecrementOrRemoveReaction_ReducesCount() {
	t := cs.T()
	existing := &[]discord.Reaction{
		{Count: 3, Emoji: discord.Emoji{Name: "👍"}},
	}
	result := decrementOrRemoveReaction(existing, discord.Emoji{Name: "👍"})

	require.Len(t, result, 1)
	assert.Equal(t, 2, result[0].Count)
}

func (cs *ConnectionSuite) TestDecrementOrRemoveReaction_RemovesAtOne() {
	t := cs.T()
	existing := &[]discord.Reaction{
		{Count: 1, Emoji: discord.Emoji{Name: "👍"}},
		{Count: 2, Emoji: discord.Emoji{Name: "👎"}},
	}
	result := decrementOrRemoveReaction(existing, discord.Emoji{Name: "👍"})

	require.Len(t, result, 1)
	assert.Equal(t, discord.Emoji{Name: "👎"}, result[0].Emoji, "only 👎 should remain")
}

// ── removeEmojiReaction ───────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestRemoveEmojiReaction_NilSlice() {
	t := cs.T()
	result := removeEmojiReaction(nil, discord.Emoji{Name: "👍"})
	assert.Nil(t, result)
}

func (cs *ConnectionSuite) TestRemoveEmojiReaction_RemovesAll() {
	t := cs.T()
	existing := &[]discord.Reaction{
		{Count: 10, Emoji: discord.Emoji{Name: "👍"}},
		{Count: 5, Emoji: discord.Emoji{Name: "👎"}},
	}
	result := removeEmojiReaction(existing, discord.Emoji{Name: "👍"})

	require.Len(t, result, 1)
	assert.Equal(t, discord.Emoji{Name: "👎"}, result[0].Emoji)
}

func (cs *ConnectionSuite) TestRemoveEmojiReaction_MissingEmojiIsNoop() {
	t := cs.T()
	existing := &[]discord.Reaction{
		{Count: 3, Emoji: discord.Emoji{Name: "👍"}},
	}
	result := removeEmojiReaction(existing, discord.Emoji{Name: "👎"})
	assert.Len(t, result, 1, "nothing should be removed when emoji is not present")
}

// ── MESSAGE_REACTION_ADD gateway event ───────────────────────────────────────

func (cs *ConnectionSuite) TestBug77_MessageReactionAdd_UpdatesCache() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(300000000000000001)
	msgID := discord.Snowflake(400000000000000002)

	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID})

	payload := map[string]any{
		"user_id":    "111",
		"channel_id": channelID.String(),
		"message_id": msgID.String(),
		"emoji":      map[string]any{"name": "👍"},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventMessageReactionAdd, nil)

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	require.NotNil(t, msg.Reactions, "Reactions must not be nil after REACTION_ADD (Bug 77)")
	require.Len(t, msg.Reactions, 1)
	assert.Equal(t, 1, (msg.Reactions)[0].Count)
	assert.Equal(t, "👍", (msg.Reactions)[0].Emoji.Name)
}

func (cs *ConnectionSuite) TestBug77_MessageReactionAdd_IncrementsExisting() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(300000000000000001)
	msgID := discord.Snowflake(400000000000000002)
	existing := []discord.Reaction{{Count: 2, Emoji: discord.Emoji{Name: "👍"}}}
	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Reactions: existing})

	payload := map[string]any{
		"user_id":    "111",
		"channel_id": channelID.String(),
		"message_id": msgID.String(),
		"emoji":      map[string]any{"name": "👍"},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(raw, events.EventMessageReactionAdd, nil)

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	require.NotNil(t, msg.Reactions)
	assert.Equal(t, 3, (msg.Reactions)[0].Count, "count must be incremented to 3")
}

// ── MESSAGE_REACTION_REMOVE gateway event ─────────────────────────────────────

func (cs *ConnectionSuite) TestBug77_MessageReactionRemove_DecrementsCache() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(300000000000000001)
	msgID := discord.Snowflake(400000000000000002)
	existing := []discord.Reaction{{Count: 3, Emoji: discord.Emoji{Name: "👍"}}}
	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Reactions: existing})

	payload := map[string]any{
		"user_id":    "111",
		"channel_id": channelID.String(),
		"message_id": msgID.String(),
		"emoji":      map[string]any{"name": "👍"},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventMessageReactionRemove, nil)

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	require.NotNil(t, msg.Reactions)
	assert.Equal(t, 2, (msg.Reactions)[0].Count, "count must be decremented to 2")
}

func (cs *ConnectionSuite) TestBug77_MessageReactionRemove_RemovesEntryAtOne() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(300000000000000001)
	msgID := discord.Snowflake(400000000000000002)
	existing := []discord.Reaction{{Count: 1, Emoji: discord.Emoji{Name: "👍"}}}
	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Reactions: existing})

	payload := map[string]any{
		"user_id":    "111",
		"channel_id": channelID.String(),
		"message_id": msgID.String(),
		"emoji":      map[string]any{"name": "👍"},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(raw, events.EventMessageReactionRemove, nil)

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	require.NotNil(t, msg.Reactions)
	assert.Empty(t, msg.Reactions, "reaction entry must be removed when count reaches 0")
}

// ── MESSAGE_REACTION_REMOVE_ALL gateway event ─────────────────────────────────

func (cs *ConnectionSuite) TestBug77_MessageReactionRemoveAll_ClearsReactions() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(300000000000000001)
	msgID := discord.Snowflake(400000000000000002)
	existing := []discord.Reaction{
		{Count: 5, Emoji: discord.Emoji{Name: "👍"}},
		{Count: 2, Emoji: discord.Emoji{Name: "👎"}},
	}
	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Reactions: existing})

	payload := map[string]any{
		"channel_id": channelID.String(),
		"message_id": msgID.String(),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(raw, events.EventMessageReactionRemoveAll, nil)

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	assert.Nil(t, msg.Reactions, "Reactions must be nil after REACTION_REMOVE_ALL")
}

// ── MESSAGE_REACTION_REMOVE_EMOJI gateway event ───────────────────────────────

func (cs *ConnectionSuite) TestBug77_MessageReactionRemoveEmoji_RemovesOneEmoji() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(300000000000000001)
	msgID := discord.Snowflake(400000000000000002)
	existing := []discord.Reaction{
		{Count: 10, Emoji: discord.Emoji{Name: "👍"}},
		{Count: 3, Emoji: discord.Emoji{Name: "👎"}},
	}
	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Reactions: existing})

	payload := map[string]any{
		"channel_id": channelID.String(),
		"message_id": msgID.String(),
		"emoji":      map[string]any{"name": "👍"},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventMessageReactionRemoveEmoji, nil)

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	require.NotNil(t, msg.Reactions)
	require.Len(t, msg.Reactions, 1, "only one emoji reaction should remain")
	assert.Equal(t, "👎", (msg.Reactions)[0].Emoji.Name)
}
