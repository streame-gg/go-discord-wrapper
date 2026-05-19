package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type emojiManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewEmojiManager(guildID discord.Snowflake, c discord.EntityClient) discord.EmojiManager {
	return &emojiManager{guildID: guildID, client: c}
}

func (m *emojiManager) Cache() *collection.Collection[discord.Snowflake, *discord.Emoji] {
	if c := m.client.ClientCache(); c != nil {
		return c.Emojis().GetByGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.Emoji]()
}

func (m *emojiManager) Get(emojiID discord.Snowflake) (*discord.Emoji, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Emojis().Get(emojiID)
	}
	return nil, false
}

func (m *emojiManager) Fetch(ctx context.Context, emojiID discord.Snowflake) (*discord.Emoji, error) {
	return m.client.GetGuildEmoji(ctx, m.guildID, emojiID)
}

func (m *emojiManager) FetchAll(ctx context.Context) (*collection.Collection[discord.Snowflake, *discord.Emoji], error) {
	emojis, err := m.client.ListGuildEmojis(ctx, m.guildID)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(emojis, func(e *discord.Emoji) discord.Snowflake {
		return e.ID
	}), nil
}

func (m *emojiManager) Create(ctx context.Context, opts discord.EmojiCreateOptions) (*discord.Emoji, error) {
	return m.client.CreateGuildEmoji(ctx, m.guildID, opts)
}

func (m *emojiManager) Resolve(input any) *discord.Emoji {
	switch v := input.(type) {
	case *discord.Emoji:
		return v
	case discord.Snowflake:
		e, _ := m.Get(v)
		return e
	case string:
		e, _ := m.Get(discord.Snowflake(v))
		return e
	}
	return nil
}

func (m *emojiManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Emojis().GetByGuild(m.guildID).Len()
	}
	return 0
}
