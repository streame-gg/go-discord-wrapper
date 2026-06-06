package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GetOrFetch returns the entity from cache when present, otherwise falls back to
// a REST Fetch. It is the convenience "give me this entity, I don't care where
// from" path — equivalent to calling Get and, on a miss, Fetch. The fetched
// entity is returned (and populated into the cache by Fetch's normal path).
//
// These methods are intentionally uniform across every cache-backed manager so
// the same call shape works whether you hold a client-level manager
// (client.Users(), client.Guilds(), client.Channels()) or a guild/channel
// sub-manager (guild.Members(), guild.Roles(), channel.Messages(), …).

// ── Client-level managers ─────────────────────────────────────────────────────

func (m *clientGuildManager) GetOrFetch(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	if v, ok := m.Get(guildID); ok {
		return v, nil
	}
	return m.Fetch(ctx, guildID)
}

func (m *clientUserManager) GetOrFetch(ctx context.Context, userID discord.Snowflake) (*discord.User, error) {
	if v, ok := m.Get(userID); ok {
		return v, nil
	}
	return m.Fetch(ctx, userID)
}

func (m *clientChannelManager) GetOrFetch(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	if v, ok := m.Get(channelID); ok {
		return v, nil
	}
	return m.Fetch(ctx, channelID)
}

// ── Guild sub-managers ────────────────────────────────────────────────────────

func (m *memberManager) GetOrFetch(ctx context.Context, userID discord.Snowflake) (*discord.GuildMember, error) {
	if v, ok := m.Get(userID); ok {
		return v, nil
	}
	return m.Fetch(ctx, userID)
}

func (m *roleManager) GetOrFetch(ctx context.Context, roleID discord.Snowflake) (*discord.Role, error) {
	if v, ok := m.Get(roleID); ok {
		return v, nil
	}
	return m.Fetch(ctx, roleID)
}

func (m *guildChannelManager) GetOrFetch(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	if v, ok := m.Get(channelID); ok {
		return v, nil
	}
	return m.Fetch(ctx, channelID)
}

func (m *emojiManager) GetOrFetch(ctx context.Context, emojiID discord.Snowflake) (*discord.Emoji, error) {
	if v, ok := m.Get(emojiID); ok {
		return v, nil
	}
	return m.Fetch(ctx, emojiID)
}

func (m *stickerManager) GetOrFetch(ctx context.Context, stickerID discord.Snowflake) (*discord.Sticker, error) {
	if v, ok := m.Get(stickerID); ok {
		return v, nil
	}
	return m.Fetch(ctx, stickerID)
}

func (m *scheduledEventManager) GetOrFetch(ctx context.Context, eventID discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	if v, ok := m.Get(eventID); ok {
		return v, nil
	}
	return m.Fetch(ctx, eventID)
}

func (m *soundboardManager) GetOrFetch(ctx context.Context, soundID discord.Snowflake) (*discord.SoundboardSound, error) {
	if v, ok := m.Get(soundID); ok {
		return v, nil
	}
	return m.Fetch(ctx, soundID)
}

// ── Channel sub-managers ──────────────────────────────────────────────────────

func (m *messageManager) GetOrFetch(ctx context.Context, messageID discord.Snowflake) (*discord.Message, error) {
	if v, ok := m.Get(messageID); ok {
		return v, nil
	}
	return m.Fetch(ctx, messageID)
}

func (m *threadManager) GetOrFetch(ctx context.Context, threadID discord.Snowflake) (*discord.Channel, error) {
	if v, ok := m.Get(threadID); ok {
		return v, nil
	}
	return m.Fetch(ctx, threadID)
}
