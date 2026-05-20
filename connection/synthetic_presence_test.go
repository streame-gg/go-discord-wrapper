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

func newPUEvent(guildID, userID string, status discord.PresenceStatus, activities []discord.FullActivity, old *discord.Presence) *events.PresenceUpdateEvent {
	return &events.PresenceUpdateEvent{
		User:        discord.PartialPresenceUser{ID: snowflake(userID)},
		GuildID:     snowflake(guildID),
		Status:      status,
		Activities:  activities,
		OldPresence: old,
	}
}

func oldPresence(status discord.PresenceStatus, activities []discord.FullActivity) *discord.Presence {
	return &discord.Presence{
		Status:     status,
		Activities: activities,
	}
}

func activity(name string, t discord.ActivityType) discord.FullActivity {
	return discord.FullActivity{Name: name, Type: t}
}

// ── derivePresenceSyntheticEvents unit tests ─────────────────────────────────

func TestPresenceSynthetic_UserOnlineFires(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil,
		oldPresence(discord.PresenceStatusOffline, nil))
	result := derivePresenceSyntheticEvents(ev)
	require.Len(t, result, 1)
	on, ok := result[0].(*events.UserOnlineEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("g1"), on.GuildID)
	assert.Equal(t, snowflake("u1"), on.UserID)
	assert.Equal(t, discord.PresenceStatusOnline, on.Status)
	assert.NotNil(t, on.OldPresence)
	assert.NotNil(t, on.NewPresence)
}

func TestPresenceSynthetic_UserOnlineFires_IdleAndDND(t *testing.T) {
	for _, status := range []discord.PresenceStatus{discord.PresenceStatusIdle, discord.PresenceStatusDND} {
		ev := newPUEvent("g1", "u1", status, nil, oldPresence(discord.PresenceStatusOffline, nil))
		result := derivePresenceSyntheticEvents(ev)
		require.Len(t, result, 1, "status=%s", status)
		_, ok := result[0].(*events.UserOnlineEvent)
		assert.True(t, ok, "status=%s", status)
	}
}

func TestPresenceSynthetic_UserOnlineNotFired_WasAlreadyOnline(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusDND, nil,
		oldPresence(discord.PresenceStatusOnline, nil))
	for _, e := range derivePresenceSyntheticEvents(ev) {
		assert.NotEqual(t, events.EventWrapperUserOnline, e.Event())
	}
}

func TestPresenceSynthetic_UserOnlineNotFired_NoOldPresence(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil, nil)
	for _, e := range derivePresenceSyntheticEvents(ev) {
		assert.NotEqual(t, events.EventWrapperUserOnline, e.Event())
	}
}

func TestPresenceSynthetic_UserOfflineFires(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOffline, nil,
		oldPresence(discord.PresenceStatusOnline, nil))
	result := derivePresenceSyntheticEvents(ev)
	require.Len(t, result, 1)
	off, ok := result[0].(*events.UserOfflineEvent)
	require.True(t, ok)
	assert.Equal(t, snowflake("u1"), off.UserID)
}

func TestPresenceSynthetic_UserOfflineNotFired_WasAlreadyOffline(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOffline, nil,
		oldPresence(discord.PresenceStatusOffline, nil))
	for _, e := range derivePresenceSyntheticEvents(ev) {
		assert.NotEqual(t, events.EventWrapperUserOffline, e.Event())
	}
}

func TestPresenceSynthetic_ActivityChangeFires_Added(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline,
		[]discord.FullActivity{activity("Spotify", discord.ActivityTypeListening)},
		oldPresence(discord.PresenceStatusOnline, nil))
	result := derivePresenceSyntheticEvents(ev)

	var found bool
	for _, e := range result {
		if e.Event() == events.EventWrapperUserActivityChange {
			found = true
			ac := e.(*events.UserActivityChangeEvent)
			assert.Empty(t, ac.OldActivities)
			assert.Len(t, ac.NewActivities, 1)
		}
	}
	assert.True(t, found, "expected UserActivityChangeEvent")
}

func TestPresenceSynthetic_ActivityChangeFires_Removed(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil,
		oldPresence(discord.PresenceStatusOnline, []discord.FullActivity{
			activity("Spotify", discord.ActivityTypeListening),
		}))
	var found bool
	for _, e := range derivePresenceSyntheticEvents(ev) {
		if e.Event() == events.EventWrapperUserActivityChange {
			found = true
		}
	}
	assert.True(t, found)
}

func TestPresenceSynthetic_ActivityChangeFires_DifferentActivity(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline,
		[]discord.FullActivity{activity("CS2", discord.ActivityTypePlaying)},
		oldPresence(discord.PresenceStatusOnline, []discord.FullActivity{
			activity("Minecraft", discord.ActivityTypePlaying),
		}))
	var found bool
	for _, e := range derivePresenceSyntheticEvents(ev) {
		if e.Event() == events.EventWrapperUserActivityChange {
			found = true
		}
	}
	assert.True(t, found)
}

func TestPresenceSynthetic_ActivityNotFired_Unchanged(t *testing.T) {
	acts := []discord.FullActivity{activity("Spotify", discord.ActivityTypeListening)}
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, acts,
		oldPresence(discord.PresenceStatusOnline, acts))
	for _, e := range derivePresenceSyntheticEvents(ev) {
		assert.NotEqual(t, events.EventWrapperUserActivityChange, e.Event())
	}
}

func TestPresenceSynthetic_ActivityNotFired_NoOldPresence(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline,
		[]discord.FullActivity{activity("Spotify", discord.ActivityTypeListening)}, nil)
	for _, e := range derivePresenceSyntheticEvents(ev) {
		assert.NotEqual(t, events.EventWrapperUserActivityChange, e.Event())
	}
}

func TestPresenceSynthetic_UserProfileUpdateFires(t *testing.T) {
	name := "newname"
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil, nil)
	ev.User.Username = &name
	result := derivePresenceSyntheticEvents(ev)
	require.Len(t, result, 1)
	pu, ok := result[0].(*events.UserProfileUpdateEvent)
	require.True(t, ok)
	require.NotNil(t, pu.User.Username)
	assert.Equal(t, "newname", *pu.User.Username)
}

func TestPresenceSynthetic_UserProfileNotFired_OnlyID(t *testing.T) {
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil, nil)
	// No optional fields set — only ID
	assert.Empty(t, derivePresenceSyntheticEvents(ev))
}

func TestPresenceSynthetic_MultipleEventsOneUpdate(t *testing.T) {
	// Status change (offline→online) + activity add + profile update
	name := "newname"
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline,
		[]discord.FullActivity{activity("Spotify", discord.ActivityTypeListening)},
		oldPresence(discord.PresenceStatusOffline, nil))
	ev.User.Username = &name

	result := derivePresenceSyntheticEvents(ev)
	typeCount := map[events.EventType]int{}
	for _, e := range result {
		typeCount[e.Event()]++
	}
	assert.Equal(t, 1, typeCount[events.EventWrapperUserOnline])
	assert.Equal(t, 1, typeCount[events.EventWrapperUserActivityChange])
	assert.Equal(t, 1, typeCount[events.EventWrapperUserUpdate])
	assert.Equal(t, 3, len(result))
}

func TestPresenceSynthetic_OnlineAndOfflineNotBothFire(t *testing.T) {
	// Can't transition both directions simultaneously
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil,
		oldPresence(discord.PresenceStatusOffline, nil))
	for _, e := range derivePresenceSyntheticEvents(ev) {
		assert.NotEqual(t, events.EventWrapperUserOffline, e.Event())
	}
}

// ── activitiesChanged helper tests ──────────────────────────────────────────

func TestActivitiesChanged_EmptyToEmpty(t *testing.T) {
	assert.False(t, activitiesChanged(nil, nil))
}

func TestActivitiesChanged_SameActivities(t *testing.T) {
	acts := []discord.FullActivity{activity("Spotify", discord.ActivityTypeListening)}
	assert.False(t, activitiesChanged(acts, acts))
}

func TestActivitiesChanged_ReorderedSameActivities(t *testing.T) {
	old := []discord.FullActivity{
		activity("Spotify", discord.ActivityTypeListening),
		activity("CS2", discord.ActivityTypePlaying),
	}
	new_ := []discord.FullActivity{
		activity("CS2", discord.ActivityTypePlaying),
		activity("Spotify", discord.ActivityTypeListening),
	}
	assert.False(t, activitiesChanged(old, new_))
}

func TestActivitiesChanged_DifferentName(t *testing.T) {
	assert.True(t, activitiesChanged(
		[]discord.FullActivity{activity("Spotify", discord.ActivityTypeListening)},
		[]discord.FullActivity{activity("YouTube Music", discord.ActivityTypeListening)},
	))
}

// ── Integration: On* helpers dispatched ─────────────────────────────────────

func TestPresenceSynthetic_OnUserOnline_Dispatches(t *testing.T) {
	client, err := NewClient("Bot fake-token", discord.IntentGuildPresences)
	require.NoError(t, err)
	defer client.Shutdown()

	var called atomic.Bool
	done := make(chan struct{})
	client.OnUserOnline(func(c *Client, ev *events.UserOnlineEvent) {
		if ev.UserID == snowflake("u1") {
			called.Store(true)
			close(done)
		}
	})

	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil,
		oldPresence(discord.PresenceStatusOffline, nil))
	for _, syn := range client.deriveSyntheticEvents(ev) {
		_ = client.enqueueOrDispatch(syn)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for UserOnline")
	}
	assert.True(t, called.Load())
}

func TestPresenceSynthetic_OnUserProfileUpdate_Dispatches(t *testing.T) {
	client, err := NewClient("Bot fake-token", discord.IntentGuildPresences)
	require.NoError(t, err)
	defer client.Shutdown()

	done := make(chan struct{})
	client.OnUserProfileUpdate(func(c *Client, ev *events.UserProfileUpdateEvent) {
		if ev.User.Username != nil && *ev.User.Username == "renamed" {
			close(done)
		}
	})

	name := "renamed"
	ev := newPUEvent("g1", "u1", discord.PresenceStatusOnline, nil, nil)
	ev.User.Username = &name
	for _, syn := range client.deriveSyntheticEvents(ev) {
		_ = client.enqueueOrDispatch(syn)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for UserProfileUpdate")
	}
}
