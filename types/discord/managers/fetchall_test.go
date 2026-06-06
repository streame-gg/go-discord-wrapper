package managers_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

var errBoom = errors.New("boom")

// FetchAll has two branches in each cache-backed manager: the error return and
// the success path that wraps the slice in a collection keyed by entity ID.

func (su *fetchallSuite) TestRoleManager_FetchAll_SuccessAndError() {
	t := su.T()
	role := &discord.Role{ID: eID, Name: "admin"}
	client := &stubEntityClient{
		cache:     newStubCache(),
		listRoles: func(discord.Snowflake) ([]*discord.Role, error) { return []*discord.Role{role}, nil },
	}
	m := managers.NewRoleManager(gID, client)
	got, err := m.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}
	if r, ok := got.Get(eID); !ok || r != role {
		t.Fatalf("collection missing role: %v, %v", r, ok)
	}

	client.listRoles = func(discord.Snowflake) ([]*discord.Role, error) { return nil, errBoom }
	if _, err := m.FetchAll(context.Background()); err == nil {
		t.Fatal("expected error from FetchAll")
	}
}

func (su *fetchallSuite) TestGuildChannelManager_FetchAll_SuccessAndError() {
	t := su.T()
	ch := &discord.Channel{ID: eID}
	client := &stubEntityClient{
		cache:        newStubCache(),
		listChannels: func(discord.Snowflake) ([]*discord.Channel, error) { return []*discord.Channel{ch}, nil },
	}
	m := managers.NewGuildChannelManager(gID, client)
	got, err := m.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}
	if c, ok := got.Get(eID); !ok || c != ch {
		t.Fatalf("collection missing channel: %v, %v", c, ok)
	}

	client.listChannels = func(discord.Snowflake) ([]*discord.Channel, error) { return nil, errBoom }
	if _, err := m.FetchAll(context.Background()); err == nil {
		t.Fatal("expected error from FetchAll")
	}
}

func (su *fetchallSuite) TestEmojiManager_FetchAll_SuccessAndError() {
	t := su.T()
	emoji := &discord.Emoji{ID: eID, Name: "fire"}
	client := &stubEntityClient{
		cache:      newStubCache(),
		listEmojis: func(discord.Snowflake) ([]*discord.Emoji, error) { return []*discord.Emoji{emoji}, nil },
	}
	m := managers.NewEmojiManager(gID, client)
	got, err := m.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}
	if e, ok := got.Get(eID); !ok || e != emoji {
		t.Fatalf("collection missing emoji: %v, %v", e, ok)
	}

	client.listEmojis = func(discord.Snowflake) ([]*discord.Emoji, error) { return nil, errBoom }
	if _, err := m.FetchAll(context.Background()); err == nil {
		t.Fatal("expected error from FetchAll")
	}
}

func (su *fetchallSuite) TestStickerManager_FetchAll_SuccessAndError() {
	t := su.T()
	sticker := &discord.Sticker{ID: eID, Name: "wow"}
	client := &stubEntityClient{
		cache:        newStubCache(),
		listStickers: func(discord.Snowflake) ([]*discord.Sticker, error) { return []*discord.Sticker{sticker}, nil },
	}
	m := managers.NewStickerManager(gID, client)
	got, err := m.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}
	if s, ok := got.Get(eID); !ok || s != sticker {
		t.Fatalf("collection missing sticker: %v, %v", s, ok)
	}

	client.listStickers = func(discord.Snowflake) ([]*discord.Sticker, error) { return nil, errBoom }
	if _, err := m.FetchAll(context.Background()); err == nil {
		t.Fatal("expected error from FetchAll")
	}
}

func (su *fetchallSuite) TestScheduledEventManager_FetchAll_SuccessAndError() {
	t := su.T()
	event := &discord.GuildScheduledEvent{ID: eID}
	client := &stubEntityClient{
		cache: newStubCache(),
		listEvents: func(discord.Snowflake) ([]*discord.GuildScheduledEvent, error) {
			return []*discord.GuildScheduledEvent{event}, nil
		},
	}
	m := managers.NewScheduledEventManager(gID, client)
	got, err := m.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}
	if e, ok := got.Get(eID); !ok || e != event {
		t.Fatalf("collection missing event: %v, %v", e, ok)
	}

	client.listEvents = func(discord.Snowflake) ([]*discord.GuildScheduledEvent, error) { return nil, errBoom }
	if _, err := m.FetchAll(context.Background()); err == nil {
		t.Fatal("expected error from FetchAll")
	}
}

type fetchallSuite struct{ suite.Suite }

func TestFetchallSuite(t *testing.T) { suite.Run(t, new(fetchallSuite)) }
