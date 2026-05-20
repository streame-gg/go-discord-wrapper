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

	g, ok := c.Guilds().Get("anyid")
	assert.Nil(t, g)
	assert.False(t, ok, "Get() must return (nil, false) without cache")

	assert.NotNil(t, c.Guilds().Cache(), "Cache() must return empty collection, not nil")
	assert.Equal(t, 0, c.Guilds().Cache().Len())
}

// ── Bug 2: guild from cache after GUILD_CREATE must have all sub-managers ────

func TestGuildFromCache_HasAllSubManagers(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake("111222333444555666")
	dispatchGuildCreate(t, c, guildCreatePayload(string(guildID)))

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

	guildID := discord.Snowflake("111222333444555666")
	packet := dispatchPacket("GUILD_CREATE", guildCreatePayload(string(guildID)))
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

	guildID := discord.Snowflake("111222333444555666")
	dispatchGuildCreate(t, c, guildCreatePayload(string(guildID)))

	// Verify baseline sub-managers after GUILD_CREATE.
	g, ok := c.Cache.Guilds().Get(guildID)
	require.True(t, ok)
	require.NotNil(t, g.Members(), "Members must exist after GUILD_CREATE")

	// Construct and dispatch a GUILD_UPDATE event. The handler requires a
	// pre-typed *events.GuildUpdateEvent (not just raw JSON).
	updateEv := &events.GuildUpdateEvent{}
	updateEv.Guild.ID = guildID
	updateEv.Guild.Name = "Updated Guild"

	rawUpdate, err := json.Marshal(updateEv.Guild)
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

	guildID := discord.Snowflake("222333444")
	channelID := discord.Snowflake("888999000")

	payload := guildCreatePayload(string(guildID), map[string]any{
		"channels": []any{
			map[string]any{
				"id":       string(channelID),
				"guild_id": string(guildID),
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

	guildID := discord.Snowflake("333444555")
	userID := discord.Snowflake("777888999")

	payload := guildCreatePayload(string(guildID), map[string]any{
		"members": []any{
			map[string]any{
				"user":      map[string]any{"id": string(userID), "username": "alice", "discriminator": "0"},
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

	guildID := discord.Snowflake("444555666")
	userID := discord.Snowflake("111000111")

	chunkPayload := map[string]any{
		"guild_id":    string(guildID),
		"chunk_index": 0,
		"chunk_count": 1,
		"members": []any{
			map[string]any{
				"user":      map[string]any{"id": string(userID), "username": "bob", "discriminator": "0"},
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
		"id":         "msg1",
		"channel_id": "ch1",
		"author":     map[string]any{"id": "u1", "username": "carol", "discriminator": "0"},
		"content":    "hello",
		"timestamp":  "2024-01-01T00:00:00Z",
	}
	raw, err := json.Marshal(msgPayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventMessageCreate, nil)

	msg, ok := c.Cache.Messages().Get("ch1", "msg1")
	require.True(t, ok, "message must be cached after MESSAGE_CREATE")
	assert.True(t, msg.IsHydrated(), "message must be hydrated after MESSAGE_CREATE")

	u, uok := c.Cache.Users().Get("u1")
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

	guildID := discord.Snowflake("999888777")
	dispatchGuildCreate(t, c, guildCreatePayload(string(guildID)))

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

// ── Manager fetch without cache ───────────────────────────────────────────────

func TestManagerFetch_WorksWithoutCache(t *testing.T) {
	guildID := discord.Snowflake("777888999000111222")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{"id": string(guildID), "name": "Fetched Guild", "roles": []any{}, "emojis": []any{}}
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
