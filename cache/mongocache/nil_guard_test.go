package mongocache

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// nilMongoCache returns a MongoDBCache with no real connection.
// The nil-guard tests call Set/Add with nil arguments which return before
// touching the MongoDB client, so no connection is needed.
func nilMongoCache() *MongoDBCache {
	return &MongoDBCache{
		opts: cache.Options{
			Messages: cache.MessageOptions{MaxPerChannel: 10},
		},
	}
}

// TestMongoNilGuards verifies that all Set/Add methods on the MongoDB stores
// return without panicking when called with a nil entity argument (Issue 9).
func TestMongoNilGuards(t *testing.T) {
	c := nilMongoCache()
	guildID := discord.Snowflake(1)

	require.NotPanics(t, func() { (&mongoGuildStore{c}).Set(nil) }, "Guilds.Set(nil)")
	require.NotPanics(t, func() { (&mongoChannelStore{c}).Set(nil) }, "Channels.Set(nil)")
	require.NotPanics(t, func() { (&mongoUserStore{c}).Set(nil) }, "Users.Set(nil)")
	require.NotPanics(t, func() { (&mongoMemberStore{c}).Set(guildID, nil) }, "Members.Set(guildID, nil)")
	require.NotPanics(t, func() { (&mongoRoleStore{c}).Set(guildID, nil) }, "Roles.Set(guildID, nil)")
	require.NotPanics(t, func() { (&mongoVoiceStateStore{c}).Set(guildID, nil) }, "VoiceStates.Set(guildID, nil)")
	require.NotPanics(t, func() { (&mongoSoundboardStore{c}).Set(guildID, nil) }, "Soundboard.Set(guildID, nil)")
	require.NotPanics(t, func() { (&mongoScheduledEventStore{c}).Set(nil) }, "ScheduledEvents.Set(nil)")
	require.NotPanics(t, func() { (&mongoStageInstanceStore{c}).Set(nil) }, "StageInstances.Set(nil)")
	require.NotPanics(t, func() { (&mongoEmojiStore{c}).Set(guildID, nil) }, "Emojis.Set(guildID, nil)")
	require.NotPanics(t, func() { (&mongoStickerStore{c}).Set(guildID, nil) }, "Stickers.Set(guildID, nil)")
	require.NotPanics(t, func() { (&mongoPresenceStore{c}).Set(nil) }, "Presences.Set(nil)")
	require.NotPanics(t, func() {
		(&mongoPresenceStore{c}).Set(&discord.Presence{}) // User.ID == 0 → no-op
	}, "Presences.Set(zeroUserID)")
	require.NotPanics(t, func() { (&mongoMessageStore{c: c}).Add(nil) }, "Messages.Add(nil)")
}
