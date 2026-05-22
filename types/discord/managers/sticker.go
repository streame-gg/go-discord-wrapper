package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type stickerManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewStickerManager(guildID discord.Snowflake, c discord.EntityClient) discord.StickerManager {
	return &stickerManager{guildID: guildID, client: c}
}

func (m *stickerManager) Cache() *collection.Collection[discord.Snowflake, *discord.Sticker] {
	if c := m.client.ClientCache(); c != nil {
		return c.Stickers().GetByGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.Sticker]()
}

func (m *stickerManager) Get(stickerID discord.Snowflake) (*discord.Sticker, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Stickers().Get(stickerID)
	}
	return nil, false
}

func (m *stickerManager) Fetch(ctx context.Context, stickerID discord.Snowflake) (*discord.Sticker, error) {
	return m.client.GetGuildSticker(ctx, m.guildID, stickerID)
}

func (m *stickerManager) FetchAll(ctx context.Context) (*collection.Collection[discord.Snowflake, *discord.Sticker], error) {
	stickers, err := m.client.ListGuildStickers(ctx, m.guildID)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(stickers, func(s *discord.Sticker) discord.Snowflake {
		return s.ID
	}), nil
}

func (m *stickerManager) Create(ctx context.Context, opts discord.StickerCreateOptions) (*discord.Sticker, error) {
	return m.client.CreateGuildSticker(ctx, m.guildID, opts)
}

func (m *stickerManager) Resolve(input any) (*discord.Sticker, error) {
	switch v := input.(type) {
	case *discord.Sticker:
		return v, nil
	case discord.Snowflake:
		s, _ := m.Get(v)
		return s, nil
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

func (m *stickerManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Stickers().GetByGuild(m.guildID).Len()
	}
	return 0
}
