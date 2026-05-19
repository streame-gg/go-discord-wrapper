package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type clientUserManager struct {
	client discord.EntityClient
}

func NewClientUserManager(c discord.EntityClient) discord.UserManager {
	return &clientUserManager{client: c}
}

func (m *clientUserManager) Cache() *collection.Collection[discord.Snowflake, *discord.User] {
	if c := m.client.ClientCache(); c != nil {
		return c.Users().All()
	}
	return collection.New[discord.Snowflake, *discord.User]()
}

func (m *clientUserManager) Get(userID discord.Snowflake) (*discord.User, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Users().Get(userID)
	}
	return nil, false
}

func (m *clientUserManager) Fetch(ctx context.Context, userID discord.Snowflake) (*discord.User, error) {
	return m.client.GetUser(ctx, userID)
}

func (m *clientUserManager) Resolve(input any) *discord.User {
	switch v := input.(type) {
	case *discord.User:
		return v
	case discord.Snowflake:
		u, _ := m.Get(v)
		return u
	case string:
		u, _ := m.Get(discord.Snowflake(v))
		return u
	}
	return nil
}

func (m *clientUserManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Users().Size()
	}
	return 0
}
