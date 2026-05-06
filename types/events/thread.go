package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type ThreadCreateEvent struct {
	common.Channel
	NewlyCreated bool `json:"newly_created,omitempty"`
}

type ThreadUpdateEvent struct {
	common.Channel
}

type ThreadDeleteEvent struct {
	ID       common.Snowflake   `json:"id"`
	GuildID  *common.Snowflake  `json:"guild_id,omitempty"`
	ParentID *common.Snowflake  `json:"parent_id,omitempty"`
	Type     common.ChannelType `json:"type"`
}

type ThreadListSyncEvent struct {
	GuildID    common.Snowflake      `json:"guild_id"`
	ChannelIDs []common.Snowflake    `json:"channel_ids,omitempty"`
	Threads    []common.Channel      `json:"threads"`
	Members    []common.ThreadMember `json:"members"`
}

// ThreadMemberUpdateEvent is dispatched when the current user's thread member is updated.
type ThreadMemberUpdateEvent struct {
	common.ThreadMember
	GuildID common.Snowflake `json:"guild_id"`
}

// ThreadMembersUpdateEvent is dispatched when thread membership for any user changes.
type ThreadMembersUpdateEvent struct {
	ID               common.Snowflake      `json:"id"`
	GuildID          common.Snowflake      `json:"guild_id"`
	MemberCount      int                   `json:"member_count"`
	AddedMembers     []common.ThreadMember `json:"added_members,omitempty"`
	RemovedMemberIDs []common.Snowflake    `json:"removed_member_ids,omitempty"`
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
