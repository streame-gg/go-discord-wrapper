package connection

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// guildCreatePayload builds a minimal but complete GatewayGuild JSON payload.
func guildCreatePayload(guildID string, extras ...map[string]any) map[string]any {
	base := map[string]any{
		"id":          guildID,
		"name":        "Test Guild",
		"unavailable": false,
		"channels":    []any{},
		"members":     []any{},
		"roles":       []any{},
		"threads":     []any{},
	}
	for _, e := range extras {
		for k, v := range e {
			base[k] = v
		}
	}
	return base
}

// dispatchGuildCreate fires internalEventHandler with a GUILD_CREATE payload.
func dispatchGuildCreate(t *testing.T, c *Client, payload map[string]any) {
	t.Helper()
	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventGuildCreate, nil)
}

// dispatchGuildUpdate fires internalEventHandler with a GUILD_UPDATE payload.
func dispatchGuildUpdate(t *testing.T, c *Client, payload map[string]any) {
	t.Helper()
	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventGuildUpdate, nil)
}

// ── Bug 1: client-level managers must work without a cache ───────────────────

func TestClientManagers_WorkWithoutCache(t *testing.T) {
	c, err := NewClient("test-token", discord.IntentGuilds)
	require.NoError(t, err)

	require.NotNil(t, c.Guilds(), "Guilds() must never return nil")
	require.NotNil(t, c.Users(), "Users() must never return nil")
	require.NotNil(t, c.Channels(), "Channels() must never return nil")

	// Cache-backed methods must return sensible zero-value results without panicking.
	assert.Equal(t, 0, c.Guilds().Size(), "Size() must return 0 without cache")
	assert.Equal(t, 0, c.Users().Size(), "Size() must return 0 without cache")
	assert.Equal(t, 0, c.Channels().Size(), "Size() must return 0 without cache")

	g, ok := c.Guilds().Get(mustSnowflake("anyid"))
	assert.Nil(t, g)
	assert.False(t, ok, "Get() must return (nil, false) without cache")

	assert.NotNil(t, c.Guilds().Cache(), "Cache() must return empty collection, not nil")
	assert.Equal(t, 0, c.Guilds().Cache().Len())
}

// ── Bug 2: guild from cache after GUILD_CREATE must have all sub-managers ────

func TestGuildFromCache_HasAllSubManagers(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555666)
	dispatchGuildCreate(t, c, guildCreatePayload(guildID.String()))

	g, ok := c.Cache.Guilds().Get(guildID)
	require.True(t, ok, "guild must be in cache after GUILD_CREATE")
	require.NotNil(t, g)

	assert.NotNil(t, g.Members(), "Members() must not be nil")
	assert.NotNil(t, g.Roles(), "Roles() must not be nil")
	assert.NotNil(t, g.Channels(), "Channels() must not be nil")
	assert.NotNil(t, g.Emojis(), "Emojis() must not be nil")
	assert.NotNil(t, g.Stickers(), "Stickers() must not be nil")
	assert.NotNil(t, g.Bans(), "Bans() must not be nil")
	assert.NotNil(t, g.ScheduledEvents(), "ScheduledEvents() must not be nil")
	assert.NotNil(t, g.StageInstances(), "StageInstances() must not be nil")
	assert.NotNil(t, g.Soundboard(), "Soundboard() must not be nil")
	assert.NotNil(t, g.Invites(), "Invites() must not be nil")
	assert.NotNil(t, g.VoiceStates(), "VoiceStates() must not be nil")
	assert.NotNil(t, g.AutoModRules(), "AutoModRules() must not be nil")
	assert.NotNil(t, g.Webhooks(), "Webhooks() must not be nil")
	assert.NotNil(t, g.Integrations(), "Integrations() must not be nil")
}

// TestGuildFromGuildCreateCallback_HasAllSubManagers verifies that the guild
// object passed to OnGuildCreate has all sub-managers populated.
func TestGuildFromGuildCreateCallback_HasAllSubManagers(t *testing.T) {
	mc := cache.NewMemoryCache(cache.Options{})
	c, err := NewClient("test-token", discord.IntentGuilds, options.WithCache(mc))
	require.NoError(t, err)

	guildID := discord.Snowflake(111222333444555666)
	packet := dispatchPacket("GUILD_CREATE", guildCreatePayload(guildID.String()))
	wsURL, closeServer := mockGateway(t, packet)
	defer closeServer()

	_ = c.connectWebsocket(wsURL, false, nil, nil)

	var managersOK atomic.Bool
	done := make(chan struct{})

	c.OnGuildCreate(func(_ *Client, ev *events.GuildCreateEvent) {
		g, isGateway := ev.Guild.(discord.GatewayGuild)
		if !isGateway {
			close(done)
			return
		}
		guild := g.Guild
		managersOK.Store(
			guild.Members() != nil &&
				guild.Roles() != nil &&
				guild.Channels() != nil &&
				guild.Emojis() != nil &&
				guild.Stickers() != nil &&
				guild.Bans() != nil &&
				guild.ScheduledEvents() != nil &&
				guild.StageInstances() != nil &&
				guild.Soundboard() != nil &&
				guild.Invites() != nil &&
				guild.VoiceStates() != nil &&
				guild.AutoModRules() != nil &&
				guild.Webhooks() != nil &&
				guild.Integrations() != nil,
		)
		close(done)
	})

	go func() { _ = c.listenWebsocket() }()

	select {
	case <-c.Ready():
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for READY")
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for GUILD_CREATE callback")
	}

	assert.True(t, managersOK.Load(), "guild in OnGuildCreate callback must have all sub-managers populated")
}

// TestGuildUpdatePreservesSubManagers verifies that a GUILD_UPDATE replaces the
// cached guild with a fully hydrated copy that still has sub-managers.
func TestGuildUpdatePreservesSubManagers(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555666)
	dispatchGuildCreate(t, c, guildCreatePayload(guildID.String()))

	// Verify baseline sub-managers after GUILD_CREATE.
	g, ok := c.Cache.Guilds().Get(guildID)
	require.True(t, ok)
	require.NotNil(t, g.Members(), "Members must exist after GUILD_CREATE")

	// Construct and dispatch a GUILD_UPDATE event. The handler requires a
	// pre-typed *events.GuildUpdateEvent (not just raw JSON).
	updateEv := &events.GuildUpdateEvent{}
	updateEv.NewGuild.ID = guildID
	updateEv.NewGuild.Name = "Updated Guild"

	rawUpdate, err := json.Marshal(updateEv.NewGuild)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(rawUpdate), events.EventGuildUpdate, updateEv)

	updated, ok := c.Cache.Guilds().Get(guildID)
	require.True(t, ok, "guild must remain in cache after GUILD_UPDATE")
	assert.Equal(t, "Updated Guild", updated.Name)

	assert.NotNil(t, updated.Members(), "Members() must not be nil after GUILD_UPDATE")
	assert.NotNil(t, updated.Roles(), "Roles() must not be nil after GUILD_UPDATE")
	assert.NotNil(t, updated.Channels(), "Channels() must not be nil after GUILD_UPDATE")
	assert.True(t, updated.IsHydrated(), "guild must be hydrated after GUILD_UPDATE")
}

// ── Audit 5: Channel sub-managers after GUILD_CREATE ─────────────────────────

func TestChannelFromCache_HasSubManagers(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(222333444)
	channelID := discord.Snowflake(888999000)

	payload := guildCreatePayload(guildID.String(), map[string]any{
		"channels": []any{
			map[string]any{
				"id":       channelID.String(),
				"guild_id": guildID.String(),
				"type":     0,
				"name":     "general",
			},
		},
	})
	dispatchGuildCreate(t, c, payload)

	ch, ok := c.Cache.Channels().Get(channelID)
	require.True(t, ok, "channel must be in cache after GUILD_CREATE")
	require.NotNil(t, ch)

	assert.NotNil(t, ch.Messages(), "Messages() must not be nil on cached channel")
	assert.NotNil(t, ch.Threads(), "Threads() must not be nil on cached channel")
	assert.True(t, ch.IsHydrated(), "channel must be hydrated")
}

// ── Audit 4: nested hydration chain ─────────────────────────────────────────

func TestNestedHydration_GuildMembers(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(333444555)
	userID := discord.Snowflake(777888999)

	payload := guildCreatePayload(guildID.String(), map[string]any{
		"members": []any{
			map[string]any{
				"user":      map[string]any{"id": userID.String(), "username": "alice", "discriminator": "0"},
				"roles":     []any{},
				"joined_at": "2024-01-01T00:00:00Z",
				"deaf":      false,
				"mute":      false,
			},
		},
	})
	dispatchGuildCreate(t, c, payload)

	m, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must be in cache after GUILD_CREATE")
	assert.True(t, m.IsHydrated(), "member must be hydrated")
	assert.Equal(t, guildID, m.GuildID, "member.GuildID must be set")
	assert.Equal(t, userID, m.UserID, "member.UserID must be set")

	u, uok := c.Cache.Users().Get(userID)
	require.True(t, uok, "user must be in cache after GUILD_CREATE")
	assert.True(t, u.IsHydrated(), "user must be hydrated")
}

func TestNestedHydration_GuildMembersChunk(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(444555666)
	userID := discord.Snowflake(111000111)

	chunkPayload := map[string]any{
		"guild_id":    guildID.String(),
		"chunk_index": 0,
		"chunk_count": 1,
		"members": []any{
			map[string]any{
				"user":      map[string]any{"id": userID.String(), "username": "bob", "discriminator": "0"},
				"roles":     []any{},
				"joined_at": "2024-01-01T00:00:00Z",
				"deaf":      false,
				"mute":      false,
			},
		},
	}
	raw, err := json.Marshal(chunkPayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventGuildMembersChunk, nil)

	m, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must be cached from GUILD_MEMBERS_CHUNK")
	assert.True(t, m.IsHydrated(), "member must be hydrated after GUILD_MEMBERS_CHUNK")
	assert.Equal(t, guildID, m.GuildID, "member.GuildID must be set from chunk event")
	assert.Equal(t, userID, m.UserID, "member.UserID must be set from chunk event")

	u, uok := c.Cache.Users().Get(userID)
	require.True(t, uok, "user must be cached from GUILD_MEMBERS_CHUNK")
	assert.True(t, u.IsHydrated(), "user must be hydrated from GUILD_MEMBERS_CHUNK")
}

func TestNestedHydration_MessageCreate(t *testing.T) {
	c := newClientWithCache(t)

	msgPayload := map[string]any{
		"id":         10,
		"channel_id": 11,
		"author":     map[string]any{"id": 12, "username": "carol", "discriminator": "0"},
		"content":    "hello",
		"timestamp":  "2024-01-01T00:00:00Z",
	}
	raw, err := json.Marshal(msgPayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventMessageCreate, nil)

	msg, ok := c.Cache.Messages().Get(11, 10)
	require.True(t, ok, "message must be cached after MESSAGE_CREATE")
	assert.True(t, msg.IsHydrated(), "message must be hydrated after MESSAGE_CREATE")

	u, uok := c.Cache.Users().Get(12)
	require.True(t, uok, "author must be cached after MESSAGE_CREATE")
	assert.True(t, u.IsHydrated(), "author must be hydrated after MESSAGE_CREATE")
}

// ── Anti-regression: reflection walk ensures all manager accessors are non-nil ─

// TestEntityHydrationCompleteness_Guild_Reflection walks every zero-argument
// method on *discord.Guild that returns a single interface value and asserts it
// is non-nil after a full GUILD_CREATE → cache round-trip.
//
// This catches regressions where a new sub-manager is added to Guild but
// setGuildManagers is not updated.
func TestEntityHydrationCompleteness_Guild_Reflection(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(999888777)
	dispatchGuildCreate(t, c, guildCreatePayload(guildID.String()))

	g, ok := c.Cache.Guilds().Get(guildID)
	require.True(t, ok, "guild must be in cache")

	gVal := reflect.ValueOf(g)
	gType := gVal.Type()

	for i := range gType.NumMethod() {
		method := gType.Method(i)
		mt := method.Type
		// Only zero-argument methods (besides receiver) with exactly one return value.
		if mt.NumIn() != 1 || mt.NumOut() != 1 {
			continue
		}
		ret := mt.Out(0)
		// Only interface return types (these are manager interfaces).
		if ret.Kind() != reflect.Interface {
			continue
		}
		result := gVal.Method(i).Call(nil)
		if len(result) == 0 {
			continue
		}
		rv := result[0]
		assert.False(t, rv.IsNil(),
			"Guild.%s() returned nil — is it missing from setGuildManagers?", method.Name)
	}
}

// ── Interaction: guild members accessible via cache ──────────────────────────

// TestInteractionCreate_GuildMembersAccessible is the end-to-end proof that
// interaction.Guild.Members().Cache().Get(someUserID) works after GUILD_CREATE.
// This is a regression test for the production panic reported when guild sub-managers
// were nil on ev.Guild in OnInteractionCreate handlers.
func TestInteractionCreate_GuildMembersAccessible(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555666)
	applicationID := discord.Snowflake(999888777666555444) // bot's user ID == its ApplicationID

	payload := guildCreatePayload(guildID.String(), map[string]any{
		"members": []any{
			map[string]any{
				"user":      map[string]any{"id": applicationID.String(), "username": "TestBot", "discriminator": "0"},
				"roles":     []any{},
				"joined_at": "2024-01-01T00:00:00Z",
				"deaf":      false,
				"mute":      false,
			},
		},
	})
	dispatchGuildCreate(t, c, payload)

	// Verify the member is in cache before dispatching the interaction.
	_, memberOK := c.Cache.Members().Get(guildID, applicationID)
	require.True(t, memberOK, "bot member must be in cache after GUILD_CREATE")

	// Build and dispatch a synthetic INTERACTION_CREATE.
	interactionPayload := map[string]any{
		"id":             1123,
		"application_id": applicationID.String(),
		"type":           2,
		"guild_id":       guildID.String(),
		"guild":          map[string]any{"id": guildID.String()},
		"channel_id":     213,
		"member": map[string]any{
			"user":      map[string]any{"id": 12312312321321, "username": "invoker", "discriminator": "0"},
			"roles":     []any{},
			"joined_at": "2024-01-01T00:00:00Z",
			"deaf":      false,
			"mute":      false,
		},
		"token":   "mytoken",
		"version": 1,
		"data":    map[string]any{"id": 12312321, "name": "test", "type": 1},
	}
	raw, err := json.Marshal(interactionPayload)
	require.NoError(t, err)

	ev := &events.InteractionCreateEvent{}
	require.NoError(t, json.Unmarshal(raw, ev))

	c.internalEventHandler(raw, events.EventInteractionCreate, ev)

	// Core assertion: Guild must be resolved from cache with all sub-managers.
	require.NotNil(t, ev.Guild, "interaction.Guild must not be nil after INTERACTION_CREATE")
	require.NotNil(t, ev.Guild.Members(), "interaction.Guild.Members() must not be nil")

	// The exact scenario from the user's proof-of-concept:
	// interaction.Guild.Members().Cache().Get(interaction.ApplicationID)
	memberFromCache, ok := ev.Guild.Members().Cache().Get(applicationID)
	assert.True(t, ok, "bot member must be retrievable via interaction.Guild.Members().Cache().Get(applicationID)")
	assert.NotNil(t, memberFromCache)

	// Also verify .Get() works (the manager-level API).
	memberViaGet, okGet := ev.Guild.Members().Get(applicationID)
	assert.True(t, okGet, "bot member must be retrievable via interaction.Guild.Members().Get(applicationID)")
	assert.NotNil(t, memberViaGet)
}

// TestInteractionCreate_NilGuildStub verifies that when Discord omits the
// guild object from the interaction payload (which happens for some interaction
// types), the library synthesizes a stub so accessor methods never panic.
func TestInteractionCreate_NilGuildStub(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555000)
	applicationID := discord.Snowflake(999888777666000111)

	// GUILD_CREATE so the guild IS in the member cache.
	payload := guildCreatePayload(guildID.String(), map[string]any{
		"members": []any{
			map[string]any{
				"user":      map[string]any{"id": applicationID.String(), "username": "Bot", "discriminator": "0"},
				"roles":     []any{},
				"joined_at": "2024-01-01T00:00:00Z",
				"deaf":      false,
				"mute":      false,
			},
		},
	})
	dispatchGuildCreate(t, c, payload)

	// Interaction payload WITHOUT the "guild" field — Discord omits it in some cases.
	interactionPayload := map[string]any{
		"id":             231231231213213,
		"application_id": applicationID.String(),
		"type":           2,
		"guild_id":       guildID.String(),
		// "guild" intentionally omitted
		"channel_id": 12321131233,
		"member": map[string]any{
			"user":      map[string]any{"id": 321123213123, "username": "invoker", "discriminator": "0"},
			"roles":     []any{},
			"joined_at": "2024-01-01T00:00:00Z",
			"deaf":      false,
			"mute":      false,
		},
		"token":   "tok2",
		"version": 1,
		"data":    map[string]any{"id": 12312321, "name": "test", "type": 1},
	}
	raw, err := json.Marshal(interactionPayload)
	require.NoError(t, err)

	ev := &events.InteractionCreateEvent{}
	require.NoError(t, json.Unmarshal(raw, ev))
	require.Nil(t, ev.Guild, "pre-condition: guild must be nil before internalEventHandler runs")

	c.internalEventHandler(json.RawMessage(raw), events.EventInteractionCreate, ev)

	require.NotNil(t, ev.Guild, "Guild must be synthesized even when absent from payload")
	require.NotNil(t, ev.Guild.Members(), "Members() must not be nil on synthesized stub")

	// Member lookup should still work via the global member cache.
	member, ok := ev.Guild.Members().Get(applicationID)
	assert.True(t, ok, "bot member must be findable after interaction even with no guild stub")
	assert.NotNil(t, member)
}

// ── Manager fetch without cache ───────────────────────────────────────────────

// ── Fix 4+5: Sticker.GuildID set during guild hydration ─────────────────────

// TestStickerHydration_GuildIDSet verifies that stickers embedded directly in a
// Guild struct (via RawStickers) have GuildID set after Hydrate, so that
// Sticker.Edit/Delete work without error on stickers accessed from the guild object.
func TestStickerHydration_GuildIDSet(t *testing.T) {
	guildID := discord.Snowflake(111000111000111000)
	stickerID := discord.Snowflake(222000222000222000)

	// Build a guild with an embedded sticker directly (bypassing JSON unmarshal).
	g := &discord.Guild{
		ID: guildID,
		RawStickers: []discord.Sticker{
			{ID: stickerID, Name: "test_sticker"},
		},
	}

	c, err := NewClient("test-token", discord.IntentGuilds)
	require.NoError(t, err)

	g.Hydrate(c)

	require.Len(t, g.RawStickers, 1)
	s := g.RawStickers[0]
	require.NotNil(t, s.GuildID, "Sticker.GuildID must not be nil after Guild.Hydrate")
	assert.Equal(t, guildID, *s.GuildID, "Sticker.GuildID must equal the parent guild ID")
	assert.True(t, s.IsHydrated(), "sticker must be hydrated")
}

// ── Fix 7: channel/thread sub-managers in hydrateEvent ──────────────────────

// ── Fix 8: bonus — GuildMemberUpdate/Remove User hydration + ThreadListSync ──

func TestHydrateEvent_GuildMemberUpdate_UserHydrated(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.GuildMemberUpdateEvent{}
	ev.NewMember.User = &discord.User{ID: mustSnowflake("user-update-1")}
	c.hydrateEvent(ev)
	assert.True(t, ev.NewMember.User.IsHydrated(), "GuildMemberUpdate: ev.NewMember.User must be hydrated after hydrateEvent")
}

func TestHydrateEvent_GuildMemberRemove_UserHydrated(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.GuildMemberRemoveEvent{}
	ev.User.ID = mustSnowflake("user-remove-1")
	c.hydrateEvent(ev)
	assert.True(t, ev.User.IsHydrated(), "GuildMemberRemove: ev.User must be hydrated after hydrateEvent")
}

func TestHydrateEvent_ThreadListSync_ChannelsHaveSubManagers(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.ThreadListSyncEvent{
		GuildID: mustSnowflake("g1"),
		Threads: []discord.Channel{
			{ID: mustSnowflake("thread-sync-1")},
			{ID: mustSnowflake("thread-sync-2")},
		},
	}
	c.hydrateEvent(ev)
	for i, ch := range ev.Threads {
		assert.NotNil(t, ch.Messages(), "ThreadListSync: Threads[%d].Messages() must not be nil", i)
		assert.NotNil(t, ch.Threads(), "ThreadListSync: Threads[%d].Threads() must not be nil", i)
		assert.True(t, ch.IsHydrated(), "ThreadListSync: Threads[%d] must be hydrated", i)
	}
}

func TestHydrateEvent_ChannelCreate_HasSubManagers(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.ChannelCreateEvent{}
	ev.Channel.ID = mustSnowflake("ch-create-1")
	c.hydrateEvent(ev)
	assert.NotNil(t, ev.Channel.Messages(), "ChannelCreate: Messages() must not be nil after hydrateEvent")
	assert.NotNil(t, ev.Channel.Threads(), "ChannelCreate: Threads() must not be nil after hydrateEvent")
	assert.True(t, ev.Channel.IsHydrated(), "ChannelCreate: channel must be hydrated after hydrateEvent")
}

func TestHydrateEvent_ChannelUpdate_HasSubManagers(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.ChannelUpdateEvent{}
	ev.NewChannel.ID = mustSnowflake("ch-update-1")
	c.hydrateEvent(ev)
	assert.NotNil(t, ev.NewChannel.Messages(), "ChannelUpdate: Messages() must not be nil after hydrateEvent")
	assert.NotNil(t, ev.NewChannel.Threads(), "ChannelUpdate: Threads() must not be nil after hydrateEvent")
}

func TestHydrateEvent_ChannelDelete_HasSubManagers(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.ChannelDeleteEvent{}
	ev.Channel.ID = mustSnowflake("ch-delete-1")
	c.hydrateEvent(ev)
	assert.NotNil(t, ev.Channel.Messages(), "ChannelDelete: Messages() must not be nil after hydrateEvent")
	assert.NotNil(t, ev.Channel.Threads(), "ChannelDelete: Threads() must not be nil after hydrateEvent")
}

func TestHydrateEvent_ThreadCreate_HasSubManagers(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.ThreadCreateEvent{}
	ev.Channel.ID = mustSnowflake("thread-create-1")
	c.hydrateEvent(ev)
	assert.NotNil(t, ev.Channel.Messages(), "ThreadCreate: Messages() must not be nil after hydrateEvent")
	assert.NotNil(t, ev.Channel.Threads(), "ThreadCreate: Threads() must not be nil after hydrateEvent")
}

func TestHydrateEvent_ThreadUpdate_HasSubManagers(t *testing.T) {
	c := newClientWithCache(t)
	ev := &events.ThreadUpdateEvent{}
	ev.NewThread.ID = mustSnowflake("thread-update-1")
	c.hydrateEvent(ev)
	assert.NotNil(t, ev.NewThread.Messages(), "ThreadUpdate: Messages() must not be nil after hydrateEvent")
	assert.NotNil(t, ev.NewThread.Threads(), "ThreadUpdate: Threads() must not be nil after hydrateEvent")
}

func TestManagerFetch_WorksWithoutCache(t *testing.T) {
	guildID := discord.Snowflake(777888999000111222)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{"id": guildID.String(), "name": "Fetched Guild", "roles": []any{}, "emojis": []any{}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Client without cache, REST pointed at our mock server.
	c, err := NewClient("Bot test-token", discord.IntentGuilds, options.WithBaseURL(ts.URL))
	require.NoError(t, err)

	ctx := context.Background()
	g, err := c.Guilds().Fetch(ctx, guildID)
	require.NoError(t, err, "Fetch must succeed without cache configured")
	require.NotNil(t, g)
	assert.Equal(t, guildID, g.ID)
}
