package connection

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func makeEmoji(id, name string) discord.Emoji {
	return discord.Emoji{ID: snowflake(id), Name: name}
}

func makeSticker(id, name string) discord.Sticker {
	return discord.Sticker{ID: snowflake(id), Name: name}
}

func makeRole(id, perms string) discord.Role {
	return discord.Role{ID: snowflake(id), Permissions: perms}
}

func emojiPtrs(emojis []discord.Emoji) []*discord.Emoji {
	out := make([]*discord.Emoji, len(emojis))
	for i := range emojis {
		e := emojis[i]
		out[i] = &e
	}
	return out
}

func stickerPtrs(stickers []discord.Sticker) []*discord.Sticker {
	out := make([]*discord.Sticker, len(stickers))
	for i := range stickers {
		s := stickers[i]
		out[i] = &s
	}
	return out
}

func newEmojiUpdateEvent(guildID string, newEmojis []discord.Emoji, oldEmojis []*discord.Emoji) *events.GuildEmojisUpdateEvent {
	return &events.GuildEmojisUpdateEvent{
		GuildID:   snowflake(guildID),
		Emojis:    newEmojis,
		OldEmojis: oldEmojis,
	}
}

func newStickerUpdateEvent(guildID string, newStickers []discord.Sticker, oldStickers []*discord.Sticker) *events.GuildStickersUpdateEvent {
	return &events.GuildStickersUpdateEvent{
		GuildID:     snowflake(guildID),
		Stickers:    newStickers,
		OldStickers: oldStickers,
	}
}

// ── Emoji unit tests ─────────────────────────────────────────────────────────

func TestEmojiSynthetic_NoOldEmojis_ReturnsNil(t *testing.T) {
	ev := newEmojiUpdateEvent("g1", []discord.Emoji{makeEmoji("e1", "wave")}, nil)
	assert.Nil(t, deriveGuildEmojiSyntheticEvents(ev))
}

func TestEmojiSynthetic_AddFires(t *testing.T) {
	ev := newEmojiUpdateEvent("g1",
		[]discord.Emoji{makeEmoji("e1", "wave"), makeEmoji("e2", "fire")},
		emojiPtrs([]discord.Emoji{makeEmoji("e1", "wave")}))

	result := deriveGuildEmojiSyntheticEvents(ev)
	require.Len(t, result, 1)
	add, ok := result[0].(*events.GuildEmojiAddEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("e2"), add.Emoji.ID)
	assert.Equal(t, snowflake("g1"), add.GuildID)
}

func TestEmojiSynthetic_RemoveFires(t *testing.T) {
	ev := newEmojiUpdateEvent("g1",
		[]discord.Emoji{makeEmoji("e1", "wave")},
		emojiPtrs([]discord.Emoji{makeEmoji("e1", "wave"), makeEmoji("e2", "fire")}))

	result := deriveGuildEmojiSyntheticEvents(ev)
	require.Len(t, result, 1)
	rem, ok := result[0].(*events.GuildEmojiRemoveEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("e2"), rem.Emoji.ID)
}

func TestEmojiSynthetic_UpdateFires_NameChanged(t *testing.T) {
	ev := newEmojiUpdateEvent("g1",
		[]discord.Emoji{makeEmoji("e1", "newname")},
		emojiPtrs([]discord.Emoji{makeEmoji("e1", "oldname")}))

	result := deriveGuildEmojiSyntheticEvents(ev)
	require.Len(t, result, 1)
	upd, ok := result[0].(*events.GuildEmojiUpdateEvent)
	require.True(t, ok)
	assert.Equal(t, "oldname", upd.OldEmoji.Name)
	assert.Equal(t, "newname", upd.NewEmoji.Name)
}

func TestEmojiSynthetic_UpdateNotFired_NameUnchanged(t *testing.T) {
	ev := newEmojiUpdateEvent("g1",
		[]discord.Emoji{makeEmoji("e1", "wave")},
		emojiPtrs([]discord.Emoji{makeEmoji("e1", "wave")}))
	assert.Empty(t, deriveGuildEmojiSyntheticEvents(ev))
}

func TestEmojiSynthetic_MultipleChanges(t *testing.T) {
	// e1 renamed, e2 removed, e3 added
	ev := newEmojiUpdateEvent("g1",
		[]discord.Emoji{makeEmoji("e1", "renamed"), makeEmoji("e3", "new")},
		emojiPtrs([]discord.Emoji{makeEmoji("e1", "original"), makeEmoji("e2", "gone")}))

	result := deriveGuildEmojiSyntheticEvents(ev)
	require.Len(t, result, 3)

	typeCount := map[events.EventType]int{}
	for _, e := range result {
		typeCount[e.Event()]++
	}
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildEmojiAdd])
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildEmojiRemove])
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildEmojiUpdate])
}

func TestEmojiSynthetic_Unchanged_NoEvents(t *testing.T) {
	emojis := []discord.Emoji{makeEmoji("e1", "wave"), makeEmoji("e2", "fire")}
	ev := newEmojiUpdateEvent("g1", emojis, emojiPtrs(emojis))
	assert.Empty(t, deriveGuildEmojiSyntheticEvents(ev))
}

// ── Sticker unit tests ───────────────────────────────────────────────────────

func TestStickerSynthetic_NoOldStickers_ReturnsNil(t *testing.T) {
	ev := newStickerUpdateEvent("g1", []discord.Sticker{makeSticker("s1", "doge")}, nil)
	assert.Nil(t, deriveGuildStickerSyntheticEvents(ev))
}

func TestStickerSynthetic_AddFires(t *testing.T) {
	ev := newStickerUpdateEvent("g1",
		[]discord.Sticker{makeSticker("s1", "doge"), makeSticker("s2", "cat")},
		stickerPtrs([]discord.Sticker{makeSticker("s1", "doge")}))

	result := deriveGuildStickerSyntheticEvents(ev)
	require.Len(t, result, 1)
	add, ok := result[0].(*events.GuildStickerAddEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("s2"), add.Sticker.ID)
}

func TestStickerSynthetic_RemoveFires(t *testing.T) {
	ev := newStickerUpdateEvent("g1",
		[]discord.Sticker{makeSticker("s1", "doge")},
		stickerPtrs([]discord.Sticker{makeSticker("s1", "doge"), makeSticker("s2", "cat")}))

	result := deriveGuildStickerSyntheticEvents(ev)
	require.Len(t, result, 1)
	rem, ok := result[0].(*events.GuildStickerRemoveEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("s2"), rem.Sticker.ID)
}

func TestStickerSynthetic_UpdateFires_NameChanged(t *testing.T) {
	ev := newStickerUpdateEvent("g1",
		[]discord.Sticker{makeSticker("s1", "newname")},
		stickerPtrs([]discord.Sticker{makeSticker("s1", "oldname")}))

	result := deriveGuildStickerSyntheticEvents(ev)
	require.Len(t, result, 1)
	upd, ok := result[0].(*events.GuildStickerUpdateEvent)
	require.True(t, ok)
	assert.Equal(t, "oldname", upd.OldSticker.Name)
	assert.Equal(t, "newname", upd.NewSticker.Name)
}

func TestStickerSynthetic_UpdateNotFired_Unchanged(t *testing.T) {
	stickers := []discord.Sticker{makeSticker("s1", "doge")}
	ev := newStickerUpdateEvent("g1", stickers, stickerPtrs(stickers))
	assert.Empty(t, deriveGuildStickerSyntheticEvents(ev))
}

// ── Role permissions unit tests ──────────────────────────────────────────────

func TestRolePermSynthetic_NoOldRole_ReturnsNil(t *testing.T) {
	role := makeRole("r1", "8")
	ev := &events.GuildRoleUpdateEvent{
		GuildID: snowflake("g1"),
		Role:    role,
		OldRole: nil,
	}
	assert.Nil(t, deriveGuildRoleSyntheticEvents(ev))
}

func TestRolePermSynthetic_PermissionsChangeFires(t *testing.T) {
	oldRole := makeRole("r1", "0")
	newRole := makeRole("r1", "8")
	ev := &events.GuildRoleUpdateEvent{
		GuildID: snowflake("g1"),
		Role:    newRole,
		OldRole: &oldRole,
	}
	result := deriveGuildRoleSyntheticEvents(ev)
	require.Len(t, result, 1)
	pc, ok := result[0].(*events.GuildRolePermissionsChangeEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("g1"), pc.GuildID)
	assert.Equal(t, snowflake("r1"), pc.RoleID)
	assert.Equal(t, "0", pc.OldPermissions)
	assert.Equal(t, "8", pc.NewPermissions)
	assert.NotNil(t, pc.OldRole)
	assert.NotNil(t, pc.NewRole)
}

func TestRolePermSynthetic_NotFired_PermissionsUnchanged(t *testing.T) {
	role := makeRole("r1", "8")
	old := makeRole("r1", "8")
	ev := &events.GuildRoleUpdateEvent{
		GuildID: snowflake("g1"),
		Role:    role,
		OldRole: &old,
	}
	assert.Empty(t, deriveGuildRoleSyntheticEvents(ev))
}

func TestRolePermSynthetic_NotFired_NameChangeOnly(t *testing.T) {
	oldRole := discord.Role{ID: snowflake("r1"), Name: "old name", Permissions: "8"}
	newRole := discord.Role{ID: snowflake("r1"), Name: "new name", Permissions: "8"}
	ev := &events.GuildRoleUpdateEvent{
		GuildID: snowflake("g1"),
		Role:    newRole,
		OldRole: &oldRole,
	}
	assert.Empty(t, deriveGuildRoleSyntheticEvents(ev))
}

// ── Integration: On* helpers dispatched ─────────────────────────────────────

func TestExpressionSynthetic_OnGuildEmojiAdd_Dispatches(t *testing.T) {
	client, err := NewClient("Bot fake-token", discord.IntentGuildExpressions)
	require.NoError(t, err)
	defer client.Shutdown()

	var called atomic.Bool
	done := make(chan struct{})
	client.OnGuildEmojiAdd(func(c *Client, ev *events.GuildEmojiAddEvent) {
		if ev.Emoji.ID == snowflake("e2") {
			called.Store(true)
			close(done)
		}
	})

	ev := newEmojiUpdateEvent("g1",
		[]discord.Emoji{makeEmoji("e1", "wave"), makeEmoji("e2", "fire")},
		emojiPtrs([]discord.Emoji{makeEmoji("e1", "wave")}))
	for _, syn := range client.deriveSyntheticEvents(ev) {
		_ = client.enqueueOrDispatch(syn)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for GuildEmojiAdd")
	}
	assert.True(t, called.Load())
}

func TestExpressionSynthetic_OnGuildRolePermissionsChange_Dispatches(t *testing.T) {
	client, err := NewClient("Bot fake-token", discord.IntentGuilds)
	require.NoError(t, err)
	defer client.Shutdown()

	done := make(chan struct{})
	client.OnGuildRolePermissionsChange(func(c *Client, ev *events.GuildRolePermissionsChangeEvent) {
		if ev.NewPermissions == "8" {
			close(done)
		}
	})

	oldRole := makeRole("r1", "0")
	newRole := makeRole("r1", "8")
	ev := &events.GuildRoleUpdateEvent{
		GuildID: snowflake("g1"),
		Role:    newRole,
		OldRole: &oldRole,
	}
	for _, syn := range client.deriveSyntheticEvents(ev) {
		_ = client.enqueueOrDispatch(syn)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for GuildRolePermissionsChange")
	}
}
