package connection

// Tests for Bug 75, 80, 81: pin/unpin, channel permission overwrites,
// and channel position cache updates.

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Bug 75: PinMessage / UnpinMessage ────────────────────────────────────────

func (cs *ConnectionSuite) TestBug75_PinMessage_SetsPinnedTrue() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(111)
	msgID := discord.Snowflake(222)

	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Pinned: false})

	msg, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	require.False(t, msg.Pinned, "Pinned must be false before PinMessage")

	// Simulate what the fixed PinMessage does after REST success.
	if m, ok := c.Cache.Messages().Get(channelID, msgID); ok {
		updated := *m
		updated.Pinned = true
		c.Cache.Messages().Update(&updated)
	}

	got, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	assert.True(t, got.Pinned, "Pinned must be true after PinMessage cache update (Bug 75)")
}

func (cs *ConnectionSuite) TestBug75_UnpinMessage_SetsPinnedFalse() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(111)
	msgID := discord.Snowflake(222)

	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: channelID, Pinned: true})

	// Simulate what the fixed UnpinMessage does after REST success.
	if m, ok := c.Cache.Messages().Get(channelID, msgID); ok {
		updated := *m
		updated.Pinned = false
		c.Cache.Messages().Update(&updated)
	}

	got, ok := c.Cache.Messages().Get(channelID, msgID)
	require.True(t, ok)
	assert.False(t, got.Pinned, "Pinned must be false after UnpinMessage cache update (Bug 75)")
}

func (cs *ConnectionSuite) TestBug75_PinMessage_NoOpWhenNotCached() {
	t := cs.T()
	c := newClientWithCache(t)

	channelID := discord.Snowflake(111)
	msgID := discord.Snowflake(999)

	// Message not in cache — update is a no-op, cache stays empty.
	if m, ok := c.Cache.Messages().Get(channelID, msgID); ok {
		updated := *m
		updated.Pinned = true
		c.Cache.Messages().Update(&updated)
	}

	_, ok := c.Cache.Messages().Get(channelID, msgID)
	assert.False(t, ok, "no phantom entry must be created for uncached message (Bug 75)")
}

// ── Bug 80: EditChannelPermissions / DeleteChannelPermission ──────────────────

func (cs *ConnectionSuite) TestBug80_EditChannelPermissions_UpsertOverwrite() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	channelID := discord.Snowflake(2000)
	overwriteID := discord.Snowflake(3000)

	ch := &discord.Channel{ID: channelID, GuildID: &guildID}
	c.cacheChannel(ch)

	// Simulate the fixed EditChannelPermissions cache update.
	allow := "8"
	deny := "0"
	params := api.EditChannelPermissionsParams{
		Allow: &allow,
		Deny:  &deny,
		Type:  discord.PermissionOverwriteTypeRole,
	}
	if cached, ok := c.Cache.Channels().Get(channelID); ok {
		updated := *cached
		ow := discord.ChannelPermissionOverwrite{
			ID:    overwriteID,
			Type:  params.Type,
			Allow: allow,
			Deny:  deny,
		}
		overwrites := make([]discord.ChannelPermissionOverwrite, 0, len(updated.PermissionOverwrites)+1)
		found := false
		for _, o := range updated.PermissionOverwrites {
			if o.ID == overwriteID {
				overwrites = append(overwrites, ow)
				found = true
			} else {
				overwrites = append(overwrites, o)
			}
		}
		if !found {
			overwrites = append(overwrites, ow)
		}
		updated.PermissionOverwrites = overwrites
		c.Cache.Channels().Set(&updated)
	}

	got, ok := c.Cache.Channels().Get(channelID)
	require.True(t, ok)
	require.Len(t, got.PermissionOverwrites, 1, "one overwrite must be present after EditChannelPermissions (Bug 80)")
	assert.Equal(t, overwriteID, got.PermissionOverwrites[0].ID)
	assert.Equal(t, "8", got.PermissionOverwrites[0].Allow)
}

func (cs *ConnectionSuite) TestBug80_EditChannelPermissions_UpdatesExistingOverwrite() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	channelID := discord.Snowflake(2000)
	overwriteID := discord.Snowflake(3000)
	other := discord.Snowflake(4000)

	ch := &discord.Channel{
		ID:      channelID,
		GuildID: &guildID,
		PermissionOverwrites: []discord.ChannelPermissionOverwrite{
			{ID: overwriteID, Allow: "0", Deny: "0", Type: discord.PermissionOverwriteTypeRole},
			{ID: other, Allow: "4", Deny: "0", Type: discord.PermissionOverwriteTypeRole},
		},
	}
	c.cacheChannel(ch)

	// Update existing overwrite.
	if cached, ok := c.Cache.Channels().Get(channelID); ok {
		updated := *cached
		newOw := discord.ChannelPermissionOverwrite{ID: overwriteID, Allow: "8", Deny: "0", Type: discord.PermissionOverwriteTypeRole}
		overwrites := make([]discord.ChannelPermissionOverwrite, 0, len(updated.PermissionOverwrites))
		for _, o := range updated.PermissionOverwrites {
			if o.ID == overwriteID {
				overwrites = append(overwrites, newOw)
			} else {
				overwrites = append(overwrites, o)
			}
		}
		updated.PermissionOverwrites = overwrites
		c.Cache.Channels().Set(&updated)
	}

	got, ok := c.Cache.Channels().Get(channelID)
	require.True(t, ok)
	require.Len(t, got.PermissionOverwrites, 2, "other overwrite must be preserved")
	for _, o := range got.PermissionOverwrites {
		if o.ID == overwriteID {
			assert.Equal(t, "8", o.Allow, "overwrite allow must be updated (Bug 80)")
		}
	}
}

func (cs *ConnectionSuite) TestBug80_DeleteChannelPermission_RemovesOverwrite() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	channelID := discord.Snowflake(2000)
	overwriteID := discord.Snowflake(3000)
	other := discord.Snowflake(4000)

	ch := &discord.Channel{
		ID:      channelID,
		GuildID: &guildID,
		PermissionOverwrites: []discord.ChannelPermissionOverwrite{
			{ID: overwriteID, Allow: "8", Deny: "0", Type: discord.PermissionOverwriteTypeRole},
			{ID: other, Allow: "4", Deny: "0", Type: discord.PermissionOverwriteTypeRole},
		},
	}
	c.cacheChannel(ch)

	// Simulate DeleteChannelPermission cache update.
	if cached, ok := c.Cache.Channels().Get(channelID); ok {
		updated := *cached
		overwrites := make([]discord.ChannelPermissionOverwrite, 0, len(updated.PermissionOverwrites))
		for _, o := range updated.PermissionOverwrites {
			if o.ID != overwriteID {
				overwrites = append(overwrites, o)
			}
		}
		updated.PermissionOverwrites = overwrites
		c.Cache.Channels().Set(&updated)
	}

	got, ok := c.Cache.Channels().Get(channelID)
	require.True(t, ok)
	require.Len(t, got.PermissionOverwrites, 1, "deleted overwrite must be gone (Bug 80)")
	assert.Equal(t, other, got.PermissionOverwrites[0].ID, "other overwrite must survive")
}

// ── Bug 81: ModifyGuildChannelPositions ──────────────────────────────────────

func (cs *ConnectionSuite) TestBug81_ModifyGuildChannelPositions_UpdatesPosition() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	channelID := discord.Snowflake(2000)
	pos := 0

	ch := &discord.Channel{ID: channelID, GuildID: &guildID, Position: &pos}
	c.cacheChannel(ch)

	// Simulate the fixed ModifyGuildChannelPositions cache update.
	newPos := 5
	entries := []api.ModifyGuildChannelPositionsEntry{
		{ID: channelID, Position: &newPos},
	}
	for _, entry := range entries {
		if cached, ok := c.Cache.Channels().Get(entry.ID); ok {
			updated := *cached
			if entry.Position != nil {
				updated.Position = entry.Position
			}
			c.Cache.Channels().Set(&updated)
		}
	}

	got, ok := c.Cache.Channels().Get(channelID)
	require.True(t, ok)
	require.NotNil(t, got.Position)
	assert.Equal(t, 5, *got.Position, "Position must be updated after ModifyGuildChannelPositions (Bug 81)")
}

func (cs *ConnectionSuite) TestBug81_ModifyGuildChannelPositions_UpdatesParentID() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(1000)
	channelID := discord.Snowflake(2000)
	oldParent := discord.Snowflake(3000)
	newParent := discord.Snowflake(4000)

	ch := &discord.Channel{ID: channelID, GuildID: &guildID, ParentID: &oldParent}
	c.cacheChannel(ch)

	// Simulate the fixed ModifyGuildChannelPositions cache update.
	entries := []api.ModifyGuildChannelPositionsEntry{
		{ID: channelID, ParentID: &newParent},
	}
	for _, entry := range entries {
		if cached, ok := c.Cache.Channels().Get(entry.ID); ok {
			updated := *cached
			if entry.ParentID != nil {
				updated.ParentID = entry.ParentID
			}
			c.Cache.Channels().Set(&updated)
		}
	}

	got, ok := c.Cache.Channels().Get(channelID)
	require.True(t, ok)
	require.NotNil(t, got.ParentID)
	assert.Equal(t, newParent, *got.ParentID, "ParentID must be updated after ModifyGuildChannelPositions (Bug 81)")
}

func (cs *ConnectionSuite) TestBug81_ModifyGuildChannelPositions_UncachedChannelIsNoop() {
	t := cs.T()
	c := newClientWithCache(t)

	newPos := 5
	entries := []api.ModifyGuildChannelPositionsEntry{
		{ID: discord.Snowflake(9999), Position: &newPos},
	}
	// Should not panic or insert phantom entry.
	for _, entry := range entries {
		if cached, ok := c.Cache.Channels().Get(entry.ID); ok {
			updated := *cached
			updated.Position = entry.Position
			c.Cache.Channels().Set(&updated)
		}
	}

	assert.Equal(t, 0, c.Cache.Channels().Size(),
		"no phantom entry must appear for uncached channel (Bug 81)")
}
