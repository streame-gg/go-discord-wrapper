package connection

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// routeServer creates a test HTTP server that routes requests by path.
// The data map is keyed by exact URL path; values are JSON-serialised as the response body.
func routeServer(t *testing.T, routes map[string]any) (string, func()) {
	t.Helper()
	mux := http.NewServeMux()
	for path, body := range routes {
		path, body := path, body
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(body)
		})
	}
	ts := httptest.NewServer(mux)
	return ts.URL, ts.Close
}

// newClientWithCacheAndServer creates a Client backed by a MemoryCache and pointing
// its REST calls to the provided base URL.
func newClientWithCacheAndServer(t *testing.T, baseURL string) *Client {
	t.Helper()
	mc := cache.NewMemoryCache(cache.Options{})
	c, err := NewClient("Bot test-token", discord.IntentGuilds,
		options.WithCache(mc),
		options.WithBaseURL(baseURL),
	)
	require.NoError(t, err)
	return c
}

// ── Fix 6: emojiManager ───────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestEmojiManager_FetchAll_HydratesAndSetsGuildID() {
	t := cs.T()
	guildID := discord.Snowflake(111222333444555666)
	emojiID := discord.Snowflake(222333444555666777)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/emojis": []any{
			map[string]any{"id": emojiID.String(), "name": "wave", "animated": false},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewEmojiManager(guildID, c)

	coll, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, coll.Len())

	e, ok := coll.Get(emojiID)
	require.True(t, ok)
	assert.True(t, e.IsHydrated(), "emoji must be hydrated after FetchAll")
	assert.Equal(t, guildID, e.GuildID, "emoji.GuildID must be set after FetchAll")
}

func (cs *ConnectionSuite) TestEmojiManager_FetchAll_CachesResults() {
	t := cs.T()
	guildID := discord.Snowflake(111222333444555667)
	emojiID := discord.Snowflake(222333444555666778)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/emojis": []any{
			map[string]any{"id": emojiID.String(), "name": "fire"},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewEmojiManager(guildID, c)

	_, err := m.FetchAll(context.Background())
	require.NoError(t, err)

	cached, ok := c.Cache.Emojis().Get(emojiID)
	require.True(t, ok, "emoji must be in cache after FetchAll")
	assert.True(t, cached.IsHydrated(), "cached emoji must be hydrated")
	assert.Equal(t, guildID, cached.GuildID, "cached emoji.GuildID must be set")
}

// ── Fix 6: stickerManager ─────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestStickerManager_FetchAll_HydratesAndSetsGuildID() {
	t := cs.T()
	guildID := discord.Snowflake(333444555666777888)
	stickerID := discord.Snowflake(444555666777888999)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/stickers": []any{
			map[string]any{"id": stickerID.String(), "name": "doge", "type": 2, "format_type": 1, "tags": "dog"},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewStickerManager(guildID, c)

	coll, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, coll.Len())

	s, ok := coll.Get(stickerID)
	require.True(t, ok)
	assert.True(t, s.IsHydrated(), "sticker must be hydrated after FetchAll")
	require.NotNil(t, s.GuildID, "sticker.GuildID must not be nil after FetchAll")
	assert.Equal(t, guildID, *s.GuildID, "sticker.GuildID must equal manager's guildID")
}

func (cs *ConnectionSuite) TestStickerManager_FetchAll_CachesResults() {
	t := cs.T()
	guildID := discord.Snowflake(333444555666777889)
	stickerID := discord.Snowflake(444555666777889000)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/stickers": []any{
			map[string]any{"id": stickerID.String(), "name": "cat", "type": 2, "format_type": 1, "tags": "cat"},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewStickerManager(guildID, c)

	_, err := m.FetchAll(context.Background())
	require.NoError(t, err)

	cached, ok := c.Cache.Stickers().Get(stickerID)
	require.True(t, ok, "sticker must be in cache after FetchAll")
	assert.True(t, cached.IsHydrated(), "cached sticker must be hydrated")
}

// ── Fix 6: soundboardManager ─────────────────────────────────────────────────

func (cs *ConnectionSuite) TestSoundboardManager_FetchAll_Hydrates() {
	t := cs.T()
	guildID := discord.Snowflake(555666777888999000)
	soundID := discord.Snowflake(666777888999000111)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/soundboard-sounds": map[string]any{
			"items": []any{
				map[string]any{"sound_id": soundID.String(), "name": "clap", "volume": 1.0},
			},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewSoundboardManager(guildID, c)

	coll, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, coll.Len())

	s, ok := coll.Get(soundID)
	require.True(t, ok)
	assert.True(t, s.IsHydrated(), "soundboard sound must be hydrated after FetchAll")
}

// ── Fix 6: scheduledEventManager ─────────────────────────────────────────────

func (cs *ConnectionSuite) TestScheduledEventManager_FetchAll_Hydrates() {
	t := cs.T()
	guildID := discord.Snowflake(777888999000111222)
	eventID := discord.Snowflake(888999000111222333)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/scheduled-events": []any{
			map[string]any{
				"id":                   eventID.String(),
				"guild_id":             guildID.String(),
				"name":                 "party",
				"privacy_level":        2,
				"status":               1,
				"entity_type":          3,
				"scheduled_start_time": "2030-01-01T00:00:00Z",
			},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewScheduledEventManager(guildID, c)

	coll, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, coll.Len())

	e, ok := coll.Get(eventID)
	require.True(t, ok)
	assert.True(t, e.IsHydrated(), "scheduled event must be hydrated after FetchAll")
}

// ── Fix 6: autoModRuleManager ────────────────────────────────────────────────

func (cs *ConnectionSuite) TestAutoModRuleManager_FetchAll_Hydrates() {
	t := cs.T()
	guildID := discord.Snowflake(999000111222333444)
	ruleID := discord.Snowflake(100011112222333445)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/auto-moderation/rules": []any{
			map[string]any{
				"id":           ruleID.String(),
				"guild_id":     guildID.String(),
				"name":         "no spam",
				"creator_id":   "100011112222333446",
				"event_type":   1,
				"trigger_type": 1,
				"actions":      []any{},
				"enabled":      true,
			},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewAutoModRuleManager(guildID, c)

	rules, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.True(t, rules[0].IsHydrated(), "automod rule must be hydrated after FetchAll")
}

// ── Fix 6: guildWebhookManager ───────────────────────────────────────────────

func (cs *ConnectionSuite) TestWebhookManager_FetchAll_Hydrates() {
	t := cs.T()
	guildID := discord.Snowflake(111222333444555668)
	webhookID := discord.Snowflake(222333444555666779)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/webhooks": []any{
			map[string]any{"id": webhookID.String(), "type": 1, "name": "logger"},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewGuildWebhookManager(guildID, c)

	webhooks, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Len(t, webhooks, 1)
	assert.True(t, webhooks[0].IsHydrated(), "webhook must be hydrated after FetchAll")
}

// ── Fix 6: integrationManager ────────────────────────────────────────────────

func (cs *ConnectionSuite) TestIntegrationManager_FetchAll_Hydrates() {
	t := cs.T()
	guildID := discord.Snowflake(333444555666777890)
	integID := discord.Snowflake(444555666777889001)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/integrations": []any{
			map[string]any{"id": integID.String(), "name": "Twitch", "type": "twitch", "enabled": true},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewIntegrationManager(guildID, c)

	integrations, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Len(t, integrations, 1)
	assert.True(t, integrations[0].IsHydrated(), "integration must be hydrated after FetchAll")
}

// ── Fix 6: guildInviteManager ────────────────────────────────────────────────

func (cs *ConnectionSuite) TestInviteManager_FetchAll_Hydrates() {
	t := cs.T()
	guildID := discord.Snowflake(555666777889000111)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/invites": []any{
			map[string]any{"code": "abc123", "type": 0},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewGuildInviteManager(guildID, c)

	invites, err := m.FetchAll(context.Background())
	require.NoError(t, err)
	require.Len(t, invites, 1)
	assert.True(t, invites[0].IsHydrated(), "invite must be hydrated after FetchAll")
}

// ── Fix 6: banManager ────────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestBanManager_FetchAll_HydratesUser() {
	t := cs.T()
	guildID := discord.Snowflake(666777888999111222)
	userID := discord.Snowflake(777888999111222333)

	baseURL, stop := routeServer(t, map[string]any{
		"/v10/guilds/" + guildID.String() + "/bans": []any{
			map[string]any{
				"user": map[string]any{"id": userID.String(), "username": "badguy", "discriminator": "0"},
			},
		},
	})
	defer stop()

	c := newClientWithCacheAndServer(t, baseURL)
	m := managers.NewBanManager(guildID, c)

	bans, err := m.FetchAll(context.Background(), discord.FetchBansOptions{})
	require.NoError(t, err)
	require.Len(t, bans, 1)
	require.NotNil(t, bans[0].User, "ban.User must not be nil")
	assert.True(t, bans[0].User.IsHydrated(), "ban.User must be hydrated after FetchAll")
}
