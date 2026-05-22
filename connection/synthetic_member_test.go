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

func newGMUEvent(guildID, userID discord.Snowflake, newRoles []string, oldMember *discord.GuildMember) *events.GuildMemberUpdateEvent {
	realRoles := make([]discord.Snowflake, len(newRoles))
	for _, role := range newRoles {
		realRoles = append(realRoles, mustSnowflake(role))
	}

	return &events.GuildMemberUpdateEvent{
		GuildID: guildID,
		NewMember: discord.GuildMember{
			UserID: userID,
			User:   &discord.User{ID: userID},
			Roles:  realRoles,
		},
		OldMember: oldMember,
	}
}

func oldMember(roles []string) *discord.GuildMember {
	realRoles := make([]discord.Snowflake, len(roles))
	for _, role := range roles {
		realRoles = append(realRoles, mustSnowflake(role))
	}

	return &discord.GuildMember{Roles: realRoles}
}

// ── deriveGuildMemberSyntheticEvents unit tests ──────────────────────────────

func TestMemberSynthetic_NoOldMember_ReturnsNil(t *testing.T) {
	ev := newGMUEvent(15, 16, []string{"12"}, nil)
	assert.Nil(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_RoleAddFires(t *testing.T) {
	const guildId discord.Snowflake = 16
	const userId discord.Snowflake = 17

	ev := newGMUEvent(guildId, userId, []string{"12", "13"}, oldMember([]string{"12"}))
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	add, ok := result[0].(*events.GuildMemberRoleAddEvent)
	require.True(t, ok)
	assert.Equal(t, guildId, add.GuildID)
	assert.Equal(t, userId, add.UserID)
	assert.Equal(t, "13", add.RoleID.String())
	assert.NotNil(t, add.OldMember)
	assert.NotNil(t, add.NewMember)
}

func TestMemberSynthetic_RoleRemoveFires(t *testing.T) {
	ev := newGMUEvent(15, 16, []string{"12"}, oldMember([]string{"12", "13"}))
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	rem, ok := result[0].(*events.GuildMemberRoleRemoveEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("13"), rem.RoleID)
}

func TestMemberSynthetic_RoleNotFired_Unchanged(t *testing.T) {
	ev := newGMUEvent(15, 16, []string{"12", "13"}, oldMember([]string{"12", "13"}))
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_MultipleRolesAddedAndRemoved(t *testing.T) {
	// Old: 12, 13 → New: 13, r3 (12 removed, r3 added)
	ev := newGMUEvent(15, 16, []string{"13", "r3"}, oldMember([]string{"12", "13"}))
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
	ev := newGMUEvent(15, 16, []string{}, oldMember([]string{}))
	ev.NewMember.Nick = strPtr("newnick")
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
	ev := newGMUEvent(15, 16, []string{}, om)
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
	ev := newGMUEvent(15, 16, []string{}, om)
	ev.NewMember.Nick = strPtr("new")
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	_, ok := result[0].(*events.GuildMemberNickChangeEvent)
	assert.True(t, ok)
}

func TestMemberSynthetic_NickNotFired_BothNil(t *testing.T) {
	ev := newGMUEvent(15, 16, []string{}, oldMember([]string{}))
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_NickNotFired_SameValue(t *testing.T) {
	om := oldMember([]string{})
	om.Nick = strPtr("same")
	ev := newGMUEvent(15, 16, []string{}, om)
	ev.NewMember.Nick = strPtr("same")
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_TimeoutFires(t *testing.T) {
	ev := newGMUEvent(15, 16, []string{}, oldMember([]string{}))
	ev.NewMember.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	to, ok := result[0].(*events.GuildMemberTimeoutEvent)
	require.True(t, ok)
	assert.Equal(t, "2030-01-01T00:00:00Z", to.CommunicationDisabledUntil)
}

func TestMemberSynthetic_TimeoutNotFired_Cleared(t *testing.T) {
	om := oldMember([]string{})
	om.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	ev := newGMUEvent(15, 16, []string{}, om)
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_TimeoutNotFired_Unchanged(t *testing.T) {
	om := oldMember([]string{})
	om.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	ev := newGMUEvent(15, 16, []string{}, om)
	ev.NewMember.CommunicationDisabledUntil = strPtr("2030-01-01T00:00:00Z")
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_BoostStartFires(t *testing.T) {
	now := time.Now()
	ev := newGMUEvent(15, 16, []string{}, oldMember([]string{}))
	ev.NewMember.PremiumSince = &now
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
	ev := newGMUEvent(15, 16, []string{}, om)
	result := deriveGuildMemberSyntheticEvents(ev)
	require.Len(t, result, 1)
	_, ok := result[0].(*events.GuildMemberBoostEndEvent)
	assert.True(t, ok)
}

func TestMemberSynthetic_BoostNotFired_StillBoosting(t *testing.T) {
	then := time.Now().Add(-time.Hour)
	om := oldMember([]string{})
	om.PremiumSince = &then
	ev := newGMUEvent(15, 16, []string{}, om)
	now := time.Now()
	ev.NewMember.PremiumSince = &now // still boosting, just refreshed
	assert.Empty(t, deriveGuildMemberSyntheticEvents(ev))
}

func TestMemberSynthetic_MultipleEventsPerUpdate(t *testing.T) {
	// Role add + nick change in one GUILD_MEMBER_UPDATE.
	om := oldMember([]string{"12"})
	om.Nick = strPtr("old")
	ev := newGMUEvent(15, 16, []string{"12", "13"}, om)
	ev.NewMember.Nick = strPtr("new")

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
	om := oldMember([]string{"12"})
	om.Nick = strPtr("oldnick")
	ev := newGMUEvent(15, 16, []string{"12", "13"}, om)
	ev.NewMember.Nick = strPtr("newnick")

	result := deriveGuildMemberSyntheticEvents(ev)
	for _, e := range result {
		switch syn := e.(type) {
		case *events.GuildMemberRoleAddEvent:
			require.NotNil(t, syn.NewMember)
			assert.Equal(t, "newnick", *syn.NewMember.Nick)
			assert.Contains(t, syn.NewMember.Roles, discord.Snowflake(13))
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
		if ev.RoleID == snowflake("13") {
			called.Store(true)
			close(done)
		}
	})

	ev := &events.GuildMemberUpdateEvent{
		GuildID: 10,
		NewMember: discord.GuildMember{
			UserID: 11,
			User:   &discord.User{ID: 11},
			Roles:  []discord.Snowflake{12, 13},
		},
		OldMember: &discord.GuildMember{
			Roles: []discord.Snowflake{12},
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

	om := &discord.GuildMember{Roles: []discord.Snowflake{}}
	om.Nick = strPtr("oldname")
	ev := &events.GuildMemberUpdateEvent{
		GuildID: 10,
		NewMember: discord.GuildMember{
			UserID: 11,
			User:   &discord.User{ID: 11},
			Roles:  []discord.Snowflake{},
			Nick:   strPtr("coolname"),
		},
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
