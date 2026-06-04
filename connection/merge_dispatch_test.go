package connection

// Tests for the presence-aware merge dispatch paths (v1.0 issue #5):
// PRESENCE_UPDATE, GUILD_MEMBER_UPDATE and MESSAGE_UPDATE all merge the incoming
// partial onto the cached entity via util.MergePartialJSON, using the raw event
// payload to decide which fields are present. These exercise the merge branch of
// internalEventHandler end-to-end (raw JSON in, merged entity in the cache out).

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── MESSAGE_UPDATE ───────────────────────────────────────────────────────────

// A MESSAGE_UPDATE that clears content to "" must overwrite the cached content —
// the intentional-zero case value-based merging cannot express.
func (cs *ConnectionSuite) TestMessageUpdate_ClearsContentToEmpty() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(100000000000000001)
	msgID := discord.Snowflake(200000000000000002)

	c.Cache.Messages().Add(&discord.Message{
		ID:        msgID,
		ChannelID: channelID,
		Content:   "original content",
		Pinned:    true,
	})

	payload := map[string]any{
		"id":         msgID.String(),
		"channel_id": channelID.String(),
		"content":    "",
		"pinned":     false,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.MessageUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventMessageUpdate, &ev)

	got, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok, "message must remain in cache after update")
	assert.Equal(t, "", got.Content, "present empty content must clear the cached value")
	assert.False(t, got.Pinned, "present false pinned must clear the cached value")
	require.NotNil(t, ev.OldMessage, "OldMessage must be populated from cache")
	assert.Equal(t, "original content", ev.OldMessage.Content)
}

// A MESSAGE_UPDATE that omits a field must preserve the cached value.
func (cs *ConnectionSuite) TestMessageUpdate_AbsentFieldPreserved() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(100000000000000003)
	msgID := discord.Snowflake(200000000000000004)

	c.Cache.Messages().Add(&discord.Message{
		ID:        msgID,
		ChannelID: channelID,
		Content:   "keep me",
		Pinned:    true,
	})

	// Only content changes; pinned is absent from the payload.
	payload := map[string]any{
		"id":         msgID.String(),
		"channel_id": channelID.String(),
		"content":    "edited",
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.MessageUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventMessageUpdate, &ev)

	got, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	assert.Equal(t, "edited", got.Content)
	assert.True(t, got.Pinned, "pinned absent from payload must be preserved from cache")
}

// MESSAGE_UPDATE for an uncached message is a no-op on the cache: the handler
// uses Messages().Update, which by contract does not insert a message that is
// absent from the channel ring. OldMessage stays nil and nothing is stored.
func (cs *ConnectionSuite) TestMessageUpdate_NoCachedMessageIsNoOp() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(100000000000000005)
	msgID := discord.Snowflake(200000000000000006)

	payload := map[string]any{
		"id":         msgID.String(),
		"channel_id": channelID.String(),
		"content":    "brand new",
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.MessageUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventMessageUpdate, &ev)

	_, ok := c.Cache.Messages().Get(channelID, msgID)
	assert.False(t, ok, "Update must not insert a message absent from the ring")
	assert.Nil(t, ev.OldMessage, "OldMessage must stay nil when nothing was cached")
}

// ── GUILD_MEMBER_UPDATE ──────────────────────────────────────────────────────

// A GUILD_MEMBER_UPDATE merges onto the cached member, preserving fields absent
// from the payload and applying those present.
func (cs *ConnectionSuite) TestGuildMemberUpdate_MergesOntoCached() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000007)
	userID := discord.Snowflake(200000000000000008)

	c.Cache.Members().Set(guildID, &discord.GuildMember{
		User:    &discord.User{ID: userID, Username: "old"},
		UserID:  userID,
		GuildID: guildID,
		Nick:    ptrStr("old-nick"),
	})

	payload := map[string]any{
		"guild_id": guildID.String(),
		"nick":     "new-nick",
		"user": map[string]any{
			"id":            userID.String(),
			"username":      "old",
			"discriminator": "0000",
		},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.GuildMemberUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildMemberUpdate, &ev)

	got, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok)
	require.NotNil(t, got.Nick)
	assert.Equal(t, "new-nick", *got.Nick, "present nick must be applied")
	assert.Equal(t, guildID, got.GuildID, "json:\"-\" GuildID must be set by the handler")
	assert.Equal(t, userID, got.UserID)
	require.NotNil(t, ev.OldMember, "OldMember must be populated from cache")
}

// ── PRESENCE_UPDATE ──────────────────────────────────────────────────────────

// The PRESENCE_UPDATE dispatch branch runs without error. Note: the in-memory
// presence store is disabled by default (PresenceDefaults), so nothing is cached
// — this exercises the dispatch path and confirms it tolerates a disabled store.
func (cs *ConnectionSuite) TestPresenceUpdate_DispatchDoesNotPanic() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000009)
	userID := discord.Snowflake(200000000000000010)

	payload := map[string]any{
		"guild_id": guildID.String(),
		"status":   "dnd",
		"user":     map[string]any{"id": userID.String()},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.PresenceUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	require.NotPanics(t, func() {
		c.internalEventHandler(json.RawMessage(raw), events.EventPresenceUpdate, &ev)
	})
}
