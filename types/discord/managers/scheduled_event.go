package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type scheduledEventManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewScheduledEventManager(guildID discord.Snowflake, c discord.EntityClient) discord.ScheduledEventManager {
	return &scheduledEventManager{guildID: guildID, client: c}
}

func (m *scheduledEventManager) Cache() *collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent] {
	if c := m.client.ClientCache(); c != nil {
		return c.ScheduledEvents().GetByGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.GuildScheduledEvent]()
}

func (m *scheduledEventManager) Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.ScheduledEvents().Get(eventID)
	}
	return nil, false
}

func (m *scheduledEventManager) Fetch(ctx context.Context, eventID discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	return m.client.GetGuildScheduledEvent(ctx, m.guildID, eventID)
}

func (m *scheduledEventManager) FetchAll(ctx context.Context) (*collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent], error) {
	events, err := m.client.ListGuildScheduledEvents(ctx, m.guildID)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(events, func(e *discord.GuildScheduledEvent) discord.Snowflake {
		return e.ID
	}), nil
}

func (m *scheduledEventManager) Create(ctx context.Context, opts discord.ScheduledEventCreateOptions) (*discord.GuildScheduledEvent, error) {
	return m.client.CreateGuildScheduledEvent(ctx, m.guildID, opts)
}

func (m *scheduledEventManager) Resolve(input any) *discord.GuildScheduledEvent {
	switch v := input.(type) {
	case *discord.GuildScheduledEvent:
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

func (m *scheduledEventManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.ScheduledEvents().GetByGuild(m.guildID).Len()
	}
	return 0
}
