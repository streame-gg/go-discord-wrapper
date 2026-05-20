package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ThreadCreateEvent struct {
	discord.Channel
	NewlyCreated bool `json:"newly_created,omitempty"`
}

type ThreadUpdateEvent struct {
	discord.Channel
	OldThread *discord.Channel `json:"old_thread,omitempty"`
}

type ThreadDeleteEvent struct {
	ID       discord.Snowflake   `json:"id"`
	GuildID  *discord.Snowflake  `json:"guild_id,omitempty"`
	ParentID *discord.Snowflake  `json:"parent_id,omitempty"`
	Type     discord.ChannelType `json:"type"`
}

type ThreadListSyncEvent struct {
	GuildID    discord.Snowflake      `json:"guild_id"`
	ChannelIDs []discord.Snowflake    `json:"channel_ids,omitempty"`
	Threads    []discord.Channel      `json:"threads"`
	Members    []discord.ThreadMember `json:"members"`
}

// ThreadMemberUpdateEvent is dispatched when the current user's thread member is updated.
type ThreadMemberUpdateEvent struct {
	discord.ThreadMember
	GuildID discord.Snowflake `json:"guild_id"`
}

// ThreadMembersUpdateEvent is dispatched when thread membership for any user changes.
type ThreadMembersUpdateEvent struct {
	ID               discord.Snowflake      `json:"id"`
	GuildID          discord.Snowflake      `json:"guild_id"`
	MemberCount      int                    `json:"member_count"`
	AddedMembers     []discord.ThreadMember `json:"added_members,omitempty"`
	RemovedMemberIDs []discord.Snowflake    `json:"removed_member_ids,omitempty"`
}

func init() {
	RegisterEvent(ThreadCreateEvent{})
	RegisterEvent(ThreadUpdateEvent{})
	RegisterEvent(ThreadDeleteEvent{})
	RegisterEvent(ThreadListSyncEvent{})
	RegisterEvent(ThreadMemberUpdateEvent{})
	RegisterEvent(ThreadMembersUpdateEvent{})
}

func (e ThreadCreateEvent) DesiredEventType() Event { return &ThreadCreateEvent{} }
func (e ThreadCreateEvent) Event() EventType        { return EventThreadCreate }

func (e ThreadUpdateEvent) DesiredEventType() Event { return &ThreadUpdateEvent{} }
func (e ThreadUpdateEvent) Event() EventType        { return EventThreadUpdate }

func (e ThreadDeleteEvent) DesiredEventType() Event { return &ThreadDeleteEvent{} }
func (e ThreadDeleteEvent) Event() EventType        { return EventThreadDelete }

func (e ThreadListSyncEvent) DesiredEventType() Event { return &ThreadListSyncEvent{} }
func (e ThreadListSyncEvent) Event() EventType        { return EventThreadListSync }

func (e ThreadMemberUpdateEvent) DesiredEventType() Event { return &ThreadMemberUpdateEvent{} }
func (e ThreadMemberUpdateEvent) Event() EventType        { return EventThreadMemberUpdate }

func (e ThreadMembersUpdateEvent) DesiredEventType() Event { return &ThreadMembersUpdateEvent{} }
func (e ThreadMembersUpdateEvent) Event() EventType        { return EventThreadMembersUpdate }
