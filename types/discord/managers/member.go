package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type memberManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewMemberManager(guildID discord.Snowflake, c discord.EntityClient) discord.MemberManager {
	return &memberManager{guildID: guildID, client: c}
}

func (m *memberManager) Cache() *collection.Collection[discord.Snowflake, *discord.GuildMember] {
	if c := m.client.ClientCache(); c != nil {
		return c.Members().AllInGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.GuildMember]()
}

func (m *memberManager) Get(userID discord.Snowflake) (*discord.GuildMember, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Members().Get(m.guildID, userID)
	}
	return nil, false
}

func (m *memberManager) Fetch(ctx context.Context, userID discord.Snowflake) (*discord.GuildMember, error) {
	return m.client.GetGuildMember(ctx, m.guildID, userID)
}

func (m *memberManager) FetchAll(ctx context.Context, opts discord.FetchMembersOptions) (*collection.Collection[discord.Snowflake, *discord.GuildMember], error) {
	members, err := m.client.ListGuildMembers(ctx, m.guildID, opts)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(members, func(mem *discord.GuildMember) discord.Snowflake {
		return mem.User.ID
	}), nil
}

func (m *memberManager) Resolve(input any) (*discord.GuildMember, error) {
	switch v := input.(type) {
	case *discord.GuildMember:
		return v, nil
	case discord.Snowflake:
		mem, _ := m.Get(v)
		return mem, nil
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

func (m *memberManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Members().AllInGuild(m.guildID).Len()
	}
	return 0
}
