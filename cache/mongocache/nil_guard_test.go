package mongocache

import (
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
func (s *mongoHelpersSuite) TestMongoNilGuards() {
	c := nilMongoCache()
	guildID := discord.Snowflake(1)
	r := s.Require()

	r.NotPanics(func() { (&mongoGuildStore{c}).Set(nil) }, "Guilds.Set(nil)")
	r.NotPanics(func() { (&mongoChannelStore{c}).Set(nil) }, "Channels.Set(nil)")
	r.NotPanics(func() { (&mongoUserStore{c}).Set(nil) }, "Users.Set(nil)")
	r.NotPanics(func() { (&mongoMemberStore{c}).Set(guildID, nil) }, "Members.Set(guildID, nil)")
	r.NotPanics(func() { (&mongoRoleStore{c}).Set(guildID, nil) }, "Roles.Set(guildID, nil)")
	r.NotPanics(func() { (&mongoVoiceStateStore{c}).Set(guildID, nil) }, "VoiceStates.Set(guildID, nil)")
	r.NotPanics(func() { (&mongoSoundboardStore{c}).Set(guildID, nil) }, "Soundboard.Set(guildID, nil)")
	r.NotPanics(func() { (&mongoScheduledEventStore{c}).Set(nil) }, "ScheduledEvents.Set(nil)")
	r.NotPanics(func() { (&mongoStageInstanceStore{c}).Set(nil) }, "StageInstances.Set(nil)")
	r.NotPanics(func() { (&mongoEmojiStore{c}).Set(guildID, nil) }, "Emojis.Set(guildID, nil)")
	r.NotPanics(func() { (&mongoStickerStore{c}).Set(guildID, nil) }, "Stickers.Set(guildID, nil)")
	r.NotPanics(func() { (&mongoPresenceStore{c}).Set(nil) }, "Presences.Set(nil)")
	r.NotPanics(func() {
		(&mongoPresenceStore{c}).Set(&discord.Presence{}) // User.ID == 0 → no-op
	}, "Presences.Set(zeroUserID)")
	r.NotPanics(func() { (&mongoMessageStore{c: c}).Add(nil) }, "Messages.Add(nil)")
}
