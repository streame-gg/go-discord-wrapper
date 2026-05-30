package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type roleManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewRoleManager(guildID discord.Snowflake, c discord.EntityClient) discord.RoleManager {
	return &roleManager{guildID: guildID, client: c}
}

func (m *roleManager) Cache() *collection.Collection[discord.Snowflake, *discord.Role] {
	if c := m.client.ClientCache(); c != nil {
		return c.Roles().GetByGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.Role]()
}

func (m *roleManager) Get(roleID discord.Snowflake) (*discord.Role, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Roles().GetByGuild(m.guildID).Get(roleID)
	}
	return nil, false
}

func (m *roleManager) Fetch(ctx context.Context, roleID discord.Snowflake) (*discord.Role, error) {
	return m.client.GetGuildRole(ctx, m.guildID, roleID)
}

func (m *roleManager) FetchAll(ctx context.Context) (*collection.Collection[discord.Snowflake, *discord.Role], error) {
	roles, err := m.client.GetGuildRoles(ctx, m.guildID)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(roles, func(r *discord.Role) discord.Snowflake {
		return r.ID
	}), nil
}

func (m *roleManager) Create(ctx context.Context, opts discord.RoleCreateOptions) (*discord.Role, error) {
	return m.client.CreateGuildRole(ctx, m.guildID, opts)
}

func (m *roleManager) Resolve(input any) (*discord.Role, error) {
	switch v := input.(type) {
	case *discord.Role:
		return v, nil
	case discord.Snowflake:
		role, _ := m.Get(v)
		return role, nil
	case string:
		snowflake, err := discord.ParseSnowflake(v)
		if err != nil {
			return nil, err
		}
		mem, _ := m.Get(*snowflake)
		return mem, nil
	}
	return nil, discord.ErrNotConvertable
}

func (m *roleManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Roles().GetByGuild(m.guildID).Len()
	}
	return 0
}
