package connection

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── Cache hydration: entities retrieved from cache are hydrated ───────────────

func TestCacheMessage_IsHydratedAfterStore(t *testing.T) {
	c, stop := launchClient(t, startSilentGateway(t), discord.IntentGuildMessages)
	defer stop()
	c.Cache = cache.NewMemoryCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	defer c.Cache.Close()

	msg := &discord.Message{ID: "1", ChannelID: "2"}
	c.cacheMessage(msg)

	got, ok := c.Cache.Messages().Get("2", "1")
	if !ok {
		t.Fatal("message not in cache")
	}
	if !got.IsHydrated() {
		t.Fatal("message from cache is not hydrated")
	}
}

func TestCacheUser_IsHydratedAfterStore(t *testing.T) {
	c, stop := launchClient(t, startSilentGateway(t), discord.IntentGuildMessages)
	defer stop()
	c.Cache = cache.NewMemoryCache(cache.Options{})
	defer c.Cache.Close()

	user := &discord.User{ID: "u1"}
	c.cacheUser(user)

	got, ok := c.Cache.Users().Get("u1")
	if !ok {
		t.Fatal("user not in cache")
	}
	if !got.IsHydrated() {
		t.Fatal("user from cache is not hydrated")
	}
}

func TestCacheMember_IsHydratedAfterStore(t *testing.T) {
	c, stop := launchClient(t, startSilentGateway(t), discord.IntentGuildMembers)
	defer stop()
	c.Cache = cache.NewMemoryCache(cache.Options{})
	defer c.Cache.Close()

	member := &discord.GuildMember{User: &discord.User{ID: "u1"}}
	c.cacheMember("g1", member)

	got, ok := c.Cache.Members().Get("g1", "u1")
	if !ok {
		t.Fatal("member not in cache")
	}
	if !got.IsHydrated() {
		t.Fatal("member from cache is not hydrated")
	}
	if got.GuildID != "g1" {
		t.Fatalf("expected GuildID g1, got %s", got.GuildID)
	}
	if got.UserID != "u1" {
		t.Fatalf("expected UserID u1, got %s", got.UserID)
	}
}

func TestCacheRole_IsHydratedAfterStore(t *testing.T) {
	c, stop := launchClient(t, startSilentGateway(t), discord.IntentGuilds)
	defer stop()
	c.Cache = cache.NewMemoryCache(cache.Options{})
	defer c.Cache.Close()

	role := &discord.Role{ID: "r1"}
	c.cacheRole("g1", role)

	got, ok := c.Cache.Roles().Get("r1")
	if !ok {
		t.Fatal("role not in cache")
	}
	if !got.IsHydrated() {
		t.Fatal("role from cache is not hydrated")
	}
	if got.GuildID != "g1" {
		t.Fatalf("expected GuildID g1, got %s", got.GuildID)
	}
}

func TestCacheChannel_IsHydratedAfterStore(t *testing.T) {
	c, stop := launchClient(t, startSilentGateway(t), discord.IntentGuilds)
	defer stop()
	c.Cache = cache.NewMemoryCache(cache.Options{})
	defer c.Cache.Close()

	ch := &discord.Channel{ID: "ch1"}
	c.cacheChannel(ch)

	got, ok := c.Cache.Channels().Get("ch1")
	if !ok {
		t.Fatal("channel not in cache")
	}
	if !got.IsHydrated() {
		t.Fatal("channel from cache is not hydrated")
	}
}

// ── Event hydration: entities in dispatched events are hydrated ───────────────

func TestEventHydration_MessageCreateIsHydrated(t *testing.T) {
	msgPacket := dispatchPacket("MESSAGE_CREATE", map[string]interface{}{
		"id": "123", "channel_id": "456",
		"author":    map[string]interface{}{"id": "1", "username": "u", "discriminator": "0"},
		"content":   "hello",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	wsURL, closeServer := mockGateway(t, msgPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuildMessages)
	defer stop()

	var hydrated atomic.Bool
	done := make(chan struct{})

	client.OnMessageCreate(func(c *Client, ev *events.MessageCreateEvent) {
		hydrated.Store(ev.IsHydrated())
		close(done)
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for MESSAGE_CREATE")
	}

	if !hydrated.Load() {
		t.Fatal("MessageCreateEvent message was not hydrated when received by handler")
	}
}

func TestEventHydration_GuildMemberAddIsHydrated(t *testing.T) {
	memberPacket := dispatchPacket("GUILD_MEMBER_ADD", map[string]interface{}{
		"guild_id": "g1",
		"user":     map[string]interface{}{"id": "u1", "username": "bob", "discriminator": "0"},
		"roles":    []string{},
		"joined_at": "2024-01-01T00:00:00Z",
		"deaf":     false,
		"mute":     false,
	})

	wsURL, closeServer := mockGateway(t, memberPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuildMembers)
	defer stop()

	var hydrated atomic.Bool
	done := make(chan struct{})

	client.OnEvent(events.EventGuildMemberAdd, func(c *Client, ev *events.GuildMemberAddEvent) {
		hydrated.Store(ev.GuildMember.IsHydrated())
		close(done)
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for GUILD_MEMBER_ADD")
	}

	if !hydrated.Load() {
		t.Fatal("GuildMemberAddEvent member was not hydrated when received by handler")
	}
}

// startSilentGateway starts a gateway that sends READY and holds open.
func startSilentGateway(t *testing.T) string {
	t.Helper()
	wsURL, closeServer := mockGateway(t)
	t.Cleanup(closeServer)
	return wsURL
}
