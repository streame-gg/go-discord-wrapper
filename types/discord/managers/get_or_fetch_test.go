package managers_test

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// TestMemberManager_GetOrFetch_CacheHit confirms a cached entity is returned
// without touching the REST client.
func (su *managersSuite) TestMemberManager_GetOrFetch_CacheHit() {
	t := su.T()
	const guildID discord.Snowflake = 1
	const userID discord.Snowflake = 2
	cached := &discord.GuildMember{UserID: userID, GuildID: guildID}
	cache := newStubCache()
	cache.members.members[guildID] = map[discord.Snowflake]*discord.GuildMember{userID: cached}

	client := &stubEntityClient{
		cache: cache,
		fetchMember: func(context.Context, discord.Snowflake, discord.Snowflake) (*discord.GuildMember, error) {
			t.Fatal("Fetch must not be called on a cache hit")
			return nil, nil
		},
	}
	m := managers.NewMemberManager(guildID, client)

	got, err := m.GetOrFetch(context.Background(), userID)
	su.NoError(err)
	su.Same(cached, got)
}

// TestMemberManager_GetOrFetch_CacheMiss confirms a miss falls through to Fetch.
func (su *managersSuite) TestMemberManager_GetOrFetch_CacheMiss() {
	t := su.T()
	const guildID discord.Snowflake = 1
	const userID discord.Snowflake = 2
	fetched := &discord.GuildMember{UserID: userID, GuildID: guildID}
	called := false
	client := &stubEntityClient{
		cache: newStubCache(), // empty
		fetchMember: func(_ context.Context, gid, uid discord.Snowflake) (*discord.GuildMember, error) {
			called = true
			su.Equal(guildID, gid)
			su.Equal(userID, uid)
			return fetched, nil
		},
	}
	m := managers.NewMemberManager(guildID, client)

	got, err := m.GetOrFetch(context.Background(), userID)
	su.NoError(err)
	su.Same(fetched, got)
	if !called {
		t.Fatal("expected Fetch to be called on a cache miss")
	}
}

// TestClientUserManager_GetOrFetch covers the client-level manager variant for
// both the cache-hit and cache-miss branches.
func (su *managersSuite) TestClientUserManager_GetOrFetch() {
	const userID discord.Snowflake = 7

	// Hit: served from cache, no fetch.
	cached := &discord.User{ID: userID, Username: "cached"}
	cache := newStubCache()
	cache.users.users[userID] = cached
	hitClient := &stubEntityClient{
		cache: cache,
		fetchUser: func(context.Context, discord.Snowflake) (*discord.User, error) {
			su.FailNow("Fetch must not be called on a cache hit")
			return nil, nil
		},
	}
	got, err := managers.NewClientUserManager(hitClient).GetOrFetch(context.Background(), userID)
	su.NoError(err)
	su.Same(cached, got)

	// Miss: falls through to fetch.
	fetched := &discord.User{ID: userID, Username: "fetched"}
	missClient := &stubEntityClient{
		cache: newStubCache(),
		fetchUser: func(_ context.Context, uid discord.Snowflake) (*discord.User, error) {
			su.Equal(userID, uid)
			return fetched, nil
		},
	}
	got, err = managers.NewClientUserManager(missClient).GetOrFetch(context.Background(), userID)
	su.NoError(err)
	su.Same(fetched, got)
}
