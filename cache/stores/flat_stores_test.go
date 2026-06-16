package stores

import (
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// flatStoresSuite exercises the simple ID-keyed stores (channel, guild, user)
// whose methods are thin wrappers over BaseStore.
type flatStoresSuite struct{ suite.Suite }

func TestFlatStoresSuite(t *testing.T) { suite.Run(t, new(flatStoresSuite)) }

// ── MemChannelStore ───────────────────────────────────────────────────────────

func (su *flatStoresSuite) TestChannelStore_Lifecycle() {
	t := su.T()
	s := NewChannelStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	id := discord.Snowflake(10)
	s.Set(&discord.Channel{ID: id, Name: util.PointerOf("general")})

	got, ok := s.Get(id)
	require.True(t, ok)
	assert.Equal(t, "general", *got.Name)

	assert.True(t, s.Has(id))
	assert.Equal(t, 1, s.Size())

	all := s.All()
	assert.Equal(t, 1, all.Len())

	s.Delete(id)
	assert.False(t, s.Has(id))
	assert.Equal(t, 0, s.Size())
}

func (su *flatStoresSuite) TestChannelStore_SetNilIsNoop() {
	t := su.T()
	s := NewChannelStore(Defaults())
	defer s.Close()

	s.Set(nil)
	assert.Equal(t, 0, s.Size())
}

func (su *flatStoresSuite) TestChannelStore_ClearTrimBytesSweep() {
	t := su.T()
	s := NewChannelStore(StoreOptions{TrackBytes: true, MaxItems: 5, TTL: time.Hour})
	defer s.Close()

	for i := 1; i <= 3; i++ {
		s.Set(&discord.Channel{ID: discord.Snowflake(i)})
	}
	assert.Positive(t, s.Bytes())

	s.TrimTo(1)
	assert.Equal(t, 1, s.Size())

	s.SweepExpired(0) // nothing expired yet
	assert.Equal(t, 1, s.Size())

	s.Clear()
	assert.Equal(t, 0, s.Size())
	assert.EqualValues(t, 0, s.Bytes())
}

// ── MemGuildStore ─────────────────────────────────────────────────────────────

func (su *flatStoresSuite) TestGuildStore_Lifecycle() {
	t := su.T()
	s := NewGuildStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	id := discord.Snowflake(20)
	s.Set(&discord.Guild{ID: id, Name: "guild"})
	s.Set(nil) // no-op

	got, ok := s.Get(id)
	require.True(t, ok)
	assert.Equal(t, "guild", got.Name)

	assert.True(t, s.Has(id))
	assert.Equal(t, 1, s.Size())
	assert.Equal(t, 1, s.All().Len())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes()) // TrackBytes off

	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())

	s.Set(&discord.Guild{ID: id})
	s.Clear()
	assert.Equal(t, 0, s.Size())

	s.Delete(id)
}

// ── MemUserStore ──────────────────────────────────────────────────────────────

func (su *flatStoresSuite) TestUserStore_Lifecycle() {
	t := su.T()
	s := NewUserStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	id := discord.Snowflake(30)
	s.Set(&discord.User{ID: id, Username: "bob"})
	s.Set(nil) // no-op

	got, ok := s.Get(id)
	require.True(t, ok)
	assert.Equal(t, "bob", got.Username)

	assert.True(t, s.Has(id))
	assert.Equal(t, 1, s.Size())
	assert.Equal(t, 1, s.All().Len())
	assert.EqualValues(t, 0, s.Bytes())

	s.SweepExpired(0)
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())

	s.Set(&discord.User{ID: id})
	s.Clear()
	assert.Equal(t, 0, s.Size())

	s.Delete(id)
}
