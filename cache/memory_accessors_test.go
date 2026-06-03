package cache_test

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── store accessors ───────────────────────────────────────────────────────────

func (s *memoryTestSuite) TestAccessors_AllStores() {
	c := newCache(cache.Options{})
	defer c.Close()

	s.NotNil(c.Guilds())
	s.NotNil(c.Channels())
	s.NotNil(c.Users())
	s.NotNil(c.Members())
	s.NotNil(c.Roles())
	s.NotNil(c.Messages())
	s.NotNil(c.VoiceStates())
	s.NotNil(c.Soundboard())
	s.NotNil(c.ScheduledEvents())
	s.NotNil(c.StageInstances())
	s.NotNil(c.Emojis())
	s.NotNil(c.Stickers())
	s.NotNil(c.Presences())
	s.NotNil(c.Bans())
	s.NotNil(c.AutoModRules())
	s.NotNil(c.Invites())
}

func (s *memoryTestSuite) TestBansStore_CRUD() {
	c := newCache(cache.Options{})
	defer c.Close()

	reason := "spam"
	c.Bans().Set(10, &discord.Ban{Reason: &reason, User: discord.User{ID: 20}})
	got, ok := c.Bans().Get(10, 20)
	s.Require().True(ok)
	s.Equal(discord.Snowflake(20), got.User.ID)

	c.Bans().Delete(10, 20)
	_, ok = c.Bans().Get(10, 20)
	s.False(ok)
}

func (s *memoryTestSuite) TestAutoModRulesStore_CRUD() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.AutoModRules().Set(10, &discord.AutoModerationRule{ID: 30, GuildID: 10, Name: "rule"})
	got, ok := c.AutoModRules().Get(30)
	s.Require().True(ok)
	s.Equal("rule", got.Name)
}

func (s *memoryTestSuite) TestInvitesStore_CRUD() {
	c := newCache(cache.Options{})
	defer c.Close()

	c.Invites().Set(&discord.Invite{Code: "xyz"})
	got, ok := c.Invites().Get("xyz")
	s.Require().True(ok)
	s.Equal("xyz", got.Code)
}

// ── NewMemoryCache defaulting ──────────────────────────────────────────────────

func (s *memoryTestSuite) TestNewMemoryCache_NegativeMaxPerChannelDefaults() {
	// MaxPerChannel < 0 must be normalised to the default of 100 (caching stays
	// enabled).
	c := cache.NewMemoryCache(cache.Options{
		SweepInterval: time.Hour,
		Messages:      cache.MessageOptions{MaxPerChannel: -1},
	})
	defer c.Close()

	c.Messages().Add(&discord.Message{ID: 1, ChannelID: 2})
	s.Equal(1, c.Messages().Size(), "caching must stay enabled when MaxPerChannel<0 defaults to 100")
}
