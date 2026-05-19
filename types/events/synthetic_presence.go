package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// UserOnlineEvent fires when a user's status transitions from offline to any
// active status (online/idle/dnd). Requires cache so OldPresence is available.
type UserOnlineEvent struct {
	GuildID     discord.Snowflake
	UserID      discord.Snowflake
	Status      discord.PresenceStatus
	OldPresence *discord.Presence
	NewPresence *discord.Presence
}

func (e UserOnlineEvent) Event() EventType        { return EventWrapperUserOnline }
func (e UserOnlineEvent) DesiredEventType() Event { return &UserOnlineEvent{} }

// UserOfflineEvent fires when a user's status transitions from any active status
// to offline. Requires cache so OldPresence is available.
type UserOfflineEvent struct {
	GuildID     discord.Snowflake
	UserID      discord.Snowflake
	OldPresence *discord.Presence
	NewPresence *discord.Presence
}

func (e UserOfflineEvent) Event() EventType        { return EventWrapperUserOffline }
func (e UserOfflineEvent) DesiredEventType() Event { return &UserOfflineEvent{} }

// UserActivityChangeEvent fires when the set of a user's activities changes.
// Derived from PRESENCE_UPDATE; requires cache to provide a baseline.
type UserActivityChangeEvent struct {
	GuildID       discord.Snowflake
	UserID        discord.Snowflake
	OldActivities []discord.FullActivity
	NewActivities []discord.FullActivity
	OldPresence   *discord.Presence
	NewPresence   *discord.Presence
}

func (e UserActivityChangeEvent) Event() EventType        { return EventWrapperUserActivityChange }
func (e UserActivityChangeEvent) DesiredEventType() Event { return &UserActivityChangeEvent{} }

// UserProfileUpdateEvent fires when a user's profile fields change as reported
// by PRESENCE_UPDATE. Distinct from the raw USER_UPDATE event, which covers
// only the bot's own account. The User field contains only the changed fields
// (Discord sends partial updates).
type UserProfileUpdateEvent struct {
	GuildID     discord.Snowflake
	UserID      discord.Snowflake
	User        discord.PartialPresenceUser
	OldPresence *discord.Presence
	NewPresence *discord.Presence
}

func (e UserProfileUpdateEvent) Event() EventType        { return EventWrapperUserUpdate }
func (e UserProfileUpdateEvent) DesiredEventType() Event { return &UserProfileUpdateEvent{} }
