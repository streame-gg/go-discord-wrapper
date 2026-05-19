package connection

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

func strPtr(s string) *string { return &s }

func newGMUEvent(guildID, userID string, newRoles []string, oldMember *discord.GuildMember) *events.GuildMemberUpdateEvent {
	roles := make([]discord.Snowflake, len(newRoles))
	for i, r := range newRoles {
		roles[i] = snowflake(r)
	}
	return &events.GuildMemberUpdateEvent{
		GuildID:   snowflake(guildID),
		User:      discord.User{ID: snowflake(userID)},
		Roles:     roles,
		OldMember: oldMember,
	}
}

func oldMember(roles []string) *discord.GuildMember {
	return &discord.GuildMember{Roles: roles}
}

// ── deriveGuildMemberSyntheticEvents unit tests ──────────────────────────────

func TestMemberSynthetic_NoOldMember_ReturnsNil(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{"r1"}, nil)
	assert.Nil(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_RoleAddFires(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{"r1", "r2"}, oldMember([]string{"r1"}))
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	add, ok := result[0].(*events.GuildMemberRoleAddEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("g1"), add.GuildID)
	assert.Equal(t, snowflake("u1"), add.UserID)
	assert.Equal(t, snowflake("r2"), add.RoleID)
	assert.NotNil(t, add.OldMember)
	assert.NotNil(t, add.NewMember)
}

func TestMemberSynthetic_RoleRemoveFires(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{"r1"}, oldMember([]string{"r1", "r2"}))
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	rem, ok := result[0].(*events.GuildMemberRoleRemoveEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("r2"), rem.RoleID)
}

func TestMemberSynthetic_RoleNotFired_Unchanged(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{"r1", "r2"}, oldMember([]string{"r1", "r2"}))
	result := deriveGuildMemberSyntheticEvents(ev)
	assert.Empty(t, result)
}

func TestMemberSynthetic_MultipleRolesAddedAndRemoved(t *testing.T) {
	// Old: r1, r2 → New: r2, r3 (r1 removed, r3 added)
	ev := newGMUEvent("g1", "u1", []string{"r2", "r3"}, oldMember([]string{"r1", "r2"}))
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 2)

	typeCount := map[events.EventType]int{}
	for _, e := range result {
		typeCount[e.Event()]++
	}
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildMemberRoleAdd])
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildMemberRoleRemove])
}

func TestMemberSynthetic_NickChangeFires_SetNew(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{}, oldMember([]string{}))
	ev.Nick = strPtr("newnick")
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	nc, ok := result[0].(*events.GuildMemberNickChangeEvent)
	require.True(t, ok)
	assert.Nil(t, nc.OldNick)
	assert.Equal(t, "newnick", *nc.NewNick)
}

func TestMemberSynthetic_NickChangeFires_Cleared(t *testing.T) {
	om := oldMember([]string{})
	om.Nick = strPtr("oldnick")
	ev := newGMUEvent("g1", "u1", []string{}, om)
	ev.Nick = nil
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	nc, ok := result[0].(*events.GuildMemberNickChangeEvent)
	require.True(t, ok)
	assert.Equal(t, "oldnick", *nc.OldNick)
	assert.Nil(t, nc.NewNick)
}

func TestMemberSynthetic_NickChangeFires_Modified(t *testing.T) {
	om := oldMember([]string{})
	om.Nick = strPtr("old")
	ev := newGMUEvent("g1", "u1", []string{}, om)
	ev.Nick = strPtr("new")
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	_, ok := result[0].(*events.GuildMemberNickChangeEvent)
	assert.True(t, ok)
}

func TestMemberSynthetic_NickNotFired_BothNil(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{}, oldMember([]string{}))
	ev.Nick = nil
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_NickNotFired_SameValue(t *testing.T) {
	om := oldMember([]string{})
	om.Nick = strPtr("same")
	ev := newGMUEvent("g1", "u1", []string{}, om)
	ev.Nick = strPtr("same")
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_TimeoutFires(t *testing.T) {
	ev := newGMUEvent("g1", "u1", []string{}, oldMember([]string{}))
	ev.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	to, ok := result[0].(*events.GuildMemberTimeoutEvent)
	require.True(t, ok)
	assert.Equal(t, "2030-01-01T00:00:00Z", to.CommunicationDisabledUntil)
}

func TestMemberSynthetic_TimeoutNotFired_Cleared(t *testing.T) {
	om := oldMember([]string{})
	om.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	ev := newGMUEvent("g1", "u1", []string{}, om)
	ev.CommunicationDisabledUntil = nil
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_TimeoutNotFired_Unchanged(t *testing.T) {
	om := oldMember([]string{})
	om.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	ev := newGMUEvent("g1", "u1", []string{}, om)
	ev.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_BoostStartFires(t *testing.T) {
	now := time.Now()
	ev := newGMUEvent("g1", "u1", []string{}, oldMember([]string{}))
	ev.PremiumSince = &now
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	bs, ok := result[0].(*events.GuildMemberBoostStartEvent)
	require.True(t, ok)
	assert.Equal(t, now, bs.PremiumSince)
}

func TestMemberSynthetic_BoostEndFires(t *testing.T) {
	then := time.Now().Add(-time.Hour)
	om := oldMember([]string{})
	om.PremiumSince = &then
	ev := newGMUEvent("g1", "u1", []string{}, om)
	ev.PremiumSince = nil
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	_, ok := result[0].(*events.GuildMemberBoostEndEvent)
	assert.True(t, ok)
}

func TestMemberSynthetic_BoostNotFired_StillBoosting(t *testing.T) {
	then := time.Now().Add(-time.Hour)
	om := oldMember([]string{})
	om.PremiumSince = &then
	ev := newGMUEvent("g1", "u1", []string{}, om)
	now := time.Now()
	ev.PremiumSince = &now // still boosting, just refreshed
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_MultipleEventsPerUpdate(t *testing.T) {
	// Role add + nick change in one GUILD_MEMBER_UPDATE.
	om := oldMember([]string{"r1"})
	om.Nick = strPtr("old")
	ev := newGMUEvent("g1", "u1", []string{"r1", "r2"}, om)
	ev.Nick = strPtr("new")

	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 2)

	typeCount := map[events.EventType]int{}
	for _, e := range result {
		typeCount[e.Event()]++
	}
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildMemberRoleAdd])
	assert.Equal(t, 1, typeCount[events.EventWrapperGuildMemberNickChange])
}

func TestMemberSynthetic_NewMemberReflectsUpdate(t *testing.T) {
	om := oldMember([]string{"r1"})
	om.Nick = strPtr("oldnick")
	ev := newGMUEvent("g1", "u1", []string{"r1", "r2"}, om)
	ev.Nick = strPtr("newnick")

	result := deriveGuildMemberSyntheticEvents(ev)
	for _, e := range result {
		switch syn := e.(type) {
		case *events.GuildMemberRoleAddEvent:
			require.NotNil(t, syn.NewMember)
			assert.Equal(t, "newnick", *syn.NewMember.Nick)
			assert.Contains(t, syn.NewMember.Roles, "r2")
		case *events.GuildMemberNickChangeEvent:
			require.NotNil(t, syn.NewMember)
			assert.Equal(t, "newnick", *syn.NewMember.Nick)
		}
	}
}

// ── Integration: On* helpers wired through dispatch ──────────────────────────

func TestMemberSynthetic_OnGuildMemberRoleAdd_Dispatches(t *testing.T) {
	client, err := NewClient("Bot fake-token", discord.IntentGuildMembers)
	require.NoError(t, err)
	defer client.Shutdown()

	var called atomic.Bool
	done := make(chan struct{})
	client.OnGuildMemberRoleAdd(func(c *Client, ev *events.GuildMemberRoleAddEvent) {
		if ev.RoleID == snowflake("r2") {
			called.Store(true)
			close(done)
		}
	})

	ev := &events.GuildMemberUpdateEvent{
		GuildID: snowflake("g1"),
		User:    discord.User{ID: snowflake("u1")},
		Roles:   []discord.Snowflake{snowflake("r1"), snowflake("r2")},
		OldMember: &discord.GuildMember{
			Roles: []string{"r1"},
		},
	}
	for _, syn := range client.deriveSyntheticEvents(ev) {
		_ = client.enqueueOrDispatch(syn)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for GuildMemberRoleAdd")
	}
	assert.True(t, called.Load())
}

func TestMemberSynthetic_OnGuildMemberNickChange_Dispatches(t *testing.T) {
	client, err := NewClient("Bot fake-token", discord.IntentGuildMembers)
	require.NoError(t, err)
	defer client.Shutdown()

	done := make(chan struct{})
	client.OnGuildMemberNickChange(func(c *Client, ev *events.GuildMemberNickChangeEvent) {
		if ev.NewNick != nil && *ev.NewNick == "coolname" {
			close(done)
		}
	})

	om := &discord.GuildMember{Roles: []string{}}
	om.Nick = strPtr("oldname")
	ev := &events.GuildMemberUpdateEvent{
		GuildID:   snowflake("g1"),
		User:      discord.User{ID: snowflake("u1")},
		Roles:     []discord.Snowflake{},
		Nick:      strPtr("coolname"),
		OldMember: om,
	}
	for _, syn := range client.deriveSyntheticEvents(ev) {
		_ = client.enqueueOrDispatch(syn)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for GuildMemberNickChange")
	}
}
