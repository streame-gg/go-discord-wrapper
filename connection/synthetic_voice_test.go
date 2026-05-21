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

func snowflake(s string) discord.Snowflake { return mustSnowflake(s) }
func sf(s string) *discord.Snowflake       { v := snowflake(s); return &v }

// newVSUEvent builds a VoiceStateUpdateEvent with the given guild/user/channel IDs
// and an optional OldState. channelID may be nil to represent "left channel".
func newVSUEvent(guildID, userID string, channelID *discord.Snowflake, oldState *discord.VoiceState) *events.VoiceStateUpdateEvent {
	ev := &events.VoiceStateUpdateEvent{
		OldState: oldState,
	}
	g := snowflake(guildID)
	ev.GuildID = &g
	ev.UserID = snowflake(userID)
	ev.ChannelID = channelID
	return ev
}

func oldState(channelID string) *discord.VoiceState {
	ch := snowflake(channelID)
	return &discord.VoiceState{ChannelID: &ch}
}

// ── deriveVoiceSyntheticEvents unit tests ────────────────────────────────────

func TestVoiceSynthetic_JoinFires(t *testing.T) {
	ev := newVSUEvent("guild1", "user1", sf("chan1"), nil)
	result := deriveVoiceSyntheticEvents(ev)
	require.Len(t, result, 1)
	join, ok := result[0].(*events.VoiceMemberJoinEvent)
	require.True(t, ok, "expected VoiceMemberJoinEvent")
	assert.Equal(t, snowflake("guild1"), join.GuildID)
	assert.Equal(t, snowflake("user1"), join.UserID)
	assert.Equal(t, snowflake("chan1"), join.ChannelID)
	assert.NotNil(t, join.State)
}

func TestVoiceSynthetic_JoinNotFired_BothNil(t *testing.T) {
	// OldState=nil, ChannelID=nil — no channel before or after, no event
	ev := newVSUEvent("guild1", "user1", nil, nil)
	result := deriveVoiceSyntheticEvents(ev)
	assert.Nil(t, result)
}

func TestVoiceSynthetic_LeaveFires(t *testing.T) {
	ev := newVSUEvent("guild1", "user1", nil, oldState("chan1"))
	result := deriveVoiceSyntheticEvents(ev)
	require.Len(t, result, 1)
	leave, ok := result[0].(*events.VoiceMemberLeaveEvent)
	require.True(t, ok, "expected VoiceMemberLeaveEvent")
	assert.Equal(t, snowflake("guild1"), leave.GuildID)
	assert.Equal(t, snowflake("user1"), leave.UserID)
	assert.Equal(t, snowflake("chan1"), leave.OldChannelID)
	assert.NotNil(t, leave.OldState)
}

func TestVoiceSynthetic_LeaveNotFired_OldStateHasNoChannel(t *testing.T) {
	// OldState exists but ChannelID is nil — malformed state, no leave event
	ev := newVSUEvent("guild1", "user1", nil, &discord.VoiceState{ChannelID: nil})
	result := deriveVoiceSyntheticEvents(ev)
	assert.Nil(t, result)
}

func TestVoiceSynthetic_MoveFires(t *testing.T) {
	ev := newVSUEvent("guild1", "user1", sf("chan2"), oldState("chan1"))
	result := deriveVoiceSyntheticEvents(ev)
	require.Len(t, result, 1)
	move, ok := result[0].(*events.VoiceMemberMoveEvent)
	require.True(t, ok, "expected VoiceMemberMoveEvent")
	assert.Equal(t, snowflake("guild1"), move.GuildID)
	assert.Equal(t, snowflake("user1"), move.UserID)
	assert.Equal(t, snowflake("chan1"), move.OldChannelID)
	assert.Equal(t, snowflake("chan2"), move.NewChannelID)
	assert.NotNil(t, move.OldState)
	assert.NotNil(t, move.NewState)
}

func TestVoiceSynthetic_MoveNotFired_OldStateNoChannel(t *testing.T) {
	// OldState exists but its ChannelID is nil — can't determine origin channel
	ev := newVSUEvent("guild1", "user1", sf("chan2"), &discord.VoiceState{ChannelID: nil})
	result := deriveVoiceSyntheticEvents(ev)
	assert.Nil(t, result)
}

func TestVoiceSynthetic_UpdateFires(t *testing.T) {
	// Same channel, state changed (e.g. muted)
	ev := newVSUEvent("guild1", "user1", sf("chan1"), oldState("chan1"))
	ev.SelfMute = true
	result := deriveVoiceSyntheticEvents(ev)
	require.Len(t, result, 1)
	upd, ok := result[0].(*events.VoiceMemberUpdateEvent)
	require.True(t, ok, "expected VoiceMemberUpdateEvent")
	assert.Equal(t, snowflake("guild1"), upd.GuildID)
	assert.Equal(t, snowflake("user1"), upd.UserID)
	assert.Equal(t, snowflake("chan1"), upd.ChannelID)
	assert.NotNil(t, upd.OldState)
	assert.NotNil(t, upd.NewState)
}

func TestVoiceSynthetic_NoGuildID_ReturnsNil(t *testing.T) {
	ev := &events.VoiceStateUpdateEvent{}
	ev.UserID = snowflake("user1")
	ev.ChannelID = sf("chan1")
	// GuildID is nil — DM voice / unknown context
	result := deriveVoiceSyntheticEvents(ev)
	assert.Nil(t, result)
}

func TestVoiceSynthetic_StateCopyIsIndependent(t *testing.T) {
	// Mutating the original event after derivation must not affect the copy in the synthetic.
	ch := snowflake("chan1")
	ev := newVSUEvent("guild1", "user1", &ch, nil)
	result := deriveVoiceSyntheticEvents(ev)
	require.Len(t, result, 1)
	join := result[0].(*events.VoiceMemberJoinEvent)

	// Change the original channel — synthetic must be unaffected.
	ch2 := snowflake("chan999")
	ev.ChannelID = &ch2
	assert.Equal(t, snowflake("chan1"), join.ChannelID, "synthetic ChannelID must not alias original")
}

// ── Integration test: synthetic events dispatched via OnVoiceMember* ─────────

func TestVoiceSynthetic_IntegrationJoin(t *testing.T) {
	vsuPacket := dispatchPacket("VOICE_STATE_UPDATE", map[string]interface{}{
		"guild_id":   "guild1",
		"user_id":    "user1",
		"channel_id": "chan1",
		"session_id": "sess1",
	})

	wsURL, closeServer := mockGateway(t, vsuPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds|discord.IntentGuildVoiceStates)
	defer stop()

	var called atomic.Bool
	done := make(chan struct{})
	client.OnVoiceMemberJoin(func(c *Client, ev *events.VoiceMemberJoinEvent) {
		called.Store(true)
		close(done)
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for VoiceMemberJoin")
	}
	assert.True(t, called.Load())
}

func TestVoiceSynthetic_IntegrationLeave(t *testing.T) {
	// First packet seeds OldState, second triggers Leave.
	joinPacket := dispatchPacket("VOICE_STATE_UPDATE", map[string]interface{}{
		"guild_id":   "guild1",
		"user_id":    "user1",
		"channel_id": "chan1",
		"session_id": "sess1",
	})
	leavePacket := dispatchPacket("VOICE_STATE_UPDATE", map[string]interface{}{
		"guild_id":   "guild1",
		"user_id":    "user1",
		"channel_id": nil,
		"session_id": "sess1",
	})

	wsURL, closeServer := mockGateway(t, joinPacket, leavePacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds|discord.IntentGuildVoiceStates)
	defer stop()

	var called atomic.Bool
	done := make(chan struct{})
	client.OnVoiceMemberLeave(func(c *Client, ev *events.VoiceMemberLeaveEvent) {
		called.Store(true)
		close(done)
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for VoiceMemberLeave")
	}
	assert.True(t, called.Load())
}

func TestVoiceSynthetic_RawAndSyntheticBothFire(t *testing.T) {
	// Both VOICE_STATE_UPDATE and VoiceMemberJoin must fire for the same event.
	vsuPacket := dispatchPacket("VOICE_STATE_UPDATE", map[string]interface{}{
		"guild_id":   "guild1",
		"user_id":    "user1",
		"channel_id": "chan1",
		"session_id": "sess1",
	})

	wsURL, closeServer := mockGateway(t, vsuPacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds|discord.IntentGuildVoiceStates)
	defer stop()

	rawDone := make(chan struct{})
	synDone := make(chan struct{})

	client.OnVoiceStateUpdate(func(c *Client, ev *events.VoiceStateUpdateEvent) { close(rawDone) })
	client.OnVoiceMemberJoin(func(c *Client, ev *events.VoiceMemberJoinEvent) { close(synDone) })

	timeout := time.After(5 * time.Second)
	select {
	case <-rawDone:
	case <-timeout:
		t.Fatal("timeout waiting for VOICE_STATE_UPDATE")
	}
	select {
	case <-synDone:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for VoiceMemberJoin")
	}
}
