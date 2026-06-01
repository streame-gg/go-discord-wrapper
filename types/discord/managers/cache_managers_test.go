package managers_test

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// Valid 18-digit snowflakes used across the cache-manager tests.
const (
	gID  discord.Snowflake = 100
	eID  discord.Snowflake = 200000000000000000
	eStr                   = "200000000000000000"
)

// ── Emoji ─────────────────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestEmojiManager_Get_And_Resolve() {
	t := su.T()
	emoji := &discord.Emoji{ID: eID, Name: "fire"}
	cache := newStubCache()
	cache.emojis.emojis[gID] = map[discord.Snowflake]*discord.Emoji{eID: emoji}
	m := managers.NewEmojiManager(gID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != emoji {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, ok := m.Get(999); ok || got != nil {
		t.Fatalf("Get miss: %v, %v", got, ok)
	}
	if m.Size() != 1 {
		t.Fatalf("Size = %d, want 1", m.Size())
	}
	if m.Cache().Len() != 1 {
		t.Fatalf("Cache len = %d, want 1", m.Cache().Len())
	}

	if got, err := m.Resolve(emoji); err != nil || got != emoji {
		t.Fatalf("Resolve(*Emoji) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eID); err != nil || got != emoji {
		t.Fatalf("Resolve(Snowflake) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eStr); err != nil || got != emoji {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if _, err := m.Resolve("bad"); err == nil {
		t.Fatal("Resolve(invalid string) should error")
	}
	if _, err := m.Resolve(1); err != discord.ErrNotConvertable {
		t.Fatalf("Resolve(int) = %v, want ErrNotConvertable", err)
	}
}

func (su *cacheManagersSuite) TestEmojiManager_NoCache() {
	t := su.T()
	m := managers.NewEmojiManager(gID, &stubEntityClient{cache: nil})
	if _, ok := m.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if m.Size() != 0 {
		t.Fatalf("Size = %d, want 0", m.Size())
	}
	if m.Cache().Len() != 0 {
		t.Fatalf("Cache len = %d, want 0", m.Cache().Len())
	}
}

// ── Sticker ───────────────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestStickerManager_Get_And_NoCache() {
	t := su.T()
	sticker := &discord.Sticker{ID: eID, Name: "wow"}
	cache := newStubCache()
	cache.stickers.stickers[gID] = map[discord.Snowflake]*discord.Sticker{eID: sticker}
	m := managers.NewStickerManager(gID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != sticker {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, err := m.Resolve(eStr); err != nil || got != sticker {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if m.Size() != 1 {
		t.Fatalf("Size = %d", m.Size())
	}

	none := managers.NewStickerManager(gID, &stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

// ── Soundboard ────────────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestSoundboardManager_Get_And_NoCache() {
	t := su.T()
	sound := &discord.SoundboardSound{SoundID: eID, Name: "ping"}
	cache := newStubCache()
	cache.soundboard.sounds[gID] = map[discord.Snowflake]*discord.SoundboardSound{eID: sound}
	m := managers.NewSoundboardManager(gID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != sound {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, ok := m.Get(999); ok || got != nil {
		t.Fatalf("Get miss: %v, %v", got, ok)
	}
	if got, err := m.Resolve(sound); err != nil || got != sound {
		t.Fatalf("Resolve(*SoundboardSound) = %v, %v", got, err)
	}
	if m.Size() != 1 || m.Cache().Len() != 1 {
		t.Fatal("size/cache should be 1")
	}

	none := managers.NewSoundboardManager(gID, &stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

// ── Scheduled events ──────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestScheduledEventManager_Get_And_NoCache() {
	t := su.T()
	event := &discord.GuildScheduledEvent{ID: eID}
	cache := newStubCache()
	cache.scheduledEvents.events[gID] = map[discord.Snowflake]*discord.GuildScheduledEvent{eID: event}
	m := managers.NewScheduledEventManager(gID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != event {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, err := m.Resolve(event); err != nil || got != event {
		t.Fatalf("Resolve(*GuildScheduledEvent) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eStr); err != nil || got != event {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if m.Size() != 1 || m.Cache().Len() != 1 {
		t.Fatal("size/cache should be 1")
	}

	none := managers.NewScheduledEventManager(gID, &stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

// ── Stage instances ───────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestStageInstanceManager_Get_And_NoCache() {
	t := su.T()
	si := &discord.StageInstance{ID: eID, GuildID: gID}
	cache := newStubCache()
	cache.stageInstances.instances[gID] = map[discord.Snowflake]*discord.StageInstance{eID: si}
	m := managers.NewStageInstanceManager(gID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != si {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, err := m.Resolve(si); err != nil || got != si {
		t.Fatalf("Resolve(*StageInstance) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eID); err != nil || got != si {
		t.Fatalf("Resolve(Snowflake) = %v, %v", got, err)
	}
	if m.Size() != 1 || m.Cache().Len() != 1 {
		t.Fatal("size/cache should be 1")
	}

	none := managers.NewStageInstanceManager(gID, &stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

// ── Voice states ──────────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestVoiceStateManager_Get_And_NoCache() {
	t := su.T()
	vs := &discord.VoiceState{UserID: eID, GuildID: sf(uint64(gID))}
	cache := newStubCache()
	cache.voiceStates.states[gID] = map[discord.Snowflake]*discord.VoiceState{eID: vs}
	m := managers.NewVoiceStateManager(gID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != vs {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, ok := m.Get(999); ok || got != nil {
		t.Fatalf("Get miss: %v, %v", got, ok)
	}
	if got, err := m.Resolve(vs); err != nil || got != vs {
		t.Fatalf("Resolve(*VoiceState) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eStr); err != nil || got != vs {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if _, err := m.Resolve("bad"); err == nil {
		t.Fatal("Resolve(invalid string) should error")
	}
	if _, err := m.Resolve(1); err != discord.ErrNotConvertable {
		t.Fatalf("Resolve(int) = %v, want ErrNotConvertable", err)
	}
	if m.Size() != 1 || m.Cache().Len() != 1 {
		t.Fatal("size/cache should be 1")
	}

	none := managers.NewVoiceStateManager(gID, &stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

// ── Thread ────────────────────────────────────────────────────────────────────

func (su *cacheManagersSuite) TestThreadManager_Get_Branches() {
	t := su.T()
	const parentID discord.Snowflake = 100
	thread := &discord.Channel{ID: eID, Type: discord.ChannelTypePublicThread, ParentID: sf(uint64(parentID))}

	cache := newStubCache()
	cache.channels.channels[eID] = thread
	m := managers.NewThreadManager(parentID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != thread {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}

	// Not a thread -> miss.
	notThread := &discord.Channel{ID: 300000000000000000, Type: discord.ChannelTypeGuildText, ParentID: sf(uint64(parentID))}
	cache.channels.channels[notThread.ID] = notThread
	if _, ok := m.Get(notThread.ID); ok {
		t.Fatal("non-thread channel should miss")
	}

	// Thread with different parent -> miss.
	otherParent := &discord.Channel{ID: 400000000000000000, Type: discord.ChannelTypePublicThread, ParentID: sf(999)}
	cache.channels.channels[otherParent.ID] = otherParent
	if _, ok := m.Get(otherParent.ID); ok {
		t.Fatal("thread under different parent should miss")
	}

	// No cache -> miss.
	none := managers.NewThreadManager(parentID, &stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}

	// Resolve covers pointer/snowflake/string/invalid/unsupported.
	if got, err := m.Resolve(thread); err != nil || got != thread {
		t.Fatalf("Resolve(*Channel) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eStr); err != nil || got != thread {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if _, err := m.Resolve("bad"); err == nil {
		t.Fatal("Resolve(invalid string) should error")
	}
	if _, err := m.Resolve(1); err != discord.ErrNotConvertable {
		t.Fatalf("Resolve(int) = %v, want ErrNotConvertable", err)
	}
}

// ── Client-scoped channel / user managers ─────────────────────────────────────

func (su *cacheManagersSuite) TestClientChannelManager_Get_And_NoCache() {
	t := su.T()
	ch := &discord.Channel{ID: eID}
	cache := newStubCache()
	cache.channels.channels[eID] = ch
	m := managers.NewClientChannelManager(&stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != ch {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, err := m.Resolve(eStr); err != nil || got != ch {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if m.Size() != 1 || m.Cache().Len() != 1 {
		t.Fatal("size/cache should be 1")
	}

	none := managers.NewClientChannelManager(&stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

func (su *cacheManagersSuite) TestClientUserManager_Get_And_NoCache() {
	t := su.T()
	u := &discord.User{ID: eID, Username: "bob"}
	cache := newStubCache()
	cache.users.users[eID] = u
	m := managers.NewClientUserManager(&stubEntityClient{cache: cache})

	if got, ok := m.Get(eID); !ok || got != u {
		t.Fatalf("Get hit: %v, %v", got, ok)
	}
	if got, err := m.Resolve(u); err != nil || got != u {
		t.Fatalf("Resolve(*User) = %v, %v", got, err)
	}
	if got, err := m.Resolve(eStr); err != nil || got != u {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	if m.Size() != 1 || m.Cache().Len() != 1 {
		t.Fatal("size/cache should be 1")
	}

	none := managers.NewClientUserManager(&stubEntityClient{cache: nil})
	if _, ok := none.Get(eID); ok {
		t.Fatal("Get should miss with nil cache")
	}
	if none.Size() != 0 || none.Cache().Len() != 0 {
		t.Fatal("no-cache size/cache should be 0")
	}
}

type cacheManagersSuite struct{ suite.Suite }

func TestCacheManagersSuite(t *testing.T) { suite.Run(t, new(cacheManagersSuite)) }
