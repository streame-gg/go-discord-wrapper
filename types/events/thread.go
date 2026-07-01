package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &ThreadCreateEvent{}
var _ Event = &ThreadUpdateEvent{}
var _ Event = &ThreadDeleteEvent{}
var _ Event = &ThreadListSyncEvent{}
var _ Event = &ThreadMemberUpdateEvent{}
var _ Event = &ThreadMembersUpdateEvent{}

func init() {
	RegisterEvent(&ThreadCreateEvent{})
	RegisterEvent(&ThreadUpdateEvent{})
	RegisterEvent(&ThreadDeleteEvent{})
	RegisterEvent(&ThreadListSyncEvent{})
	RegisterEvent(&ThreadMemberUpdateEvent{})
	RegisterEvent(&ThreadMembersUpdateEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#thread-create
type ThreadCreateEvent struct {
	Thread       discord.Channel `json:"-"`
	NewlyCreated *bool           `json:"newly_created,omitempty"`

	Guild *discord.Guild `json:"guild,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#thread-update
type ThreadUpdateEvent struct {
	NewThread discord.Channel  `json:"-"`
	OldThread *discord.Channel `json:"-"`

	Guild *discord.Guild `json:"guild,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#thread-delete
type ThreadDeleteEvent struct {
	ID       discord.Snowflake   `json:"id"`
	GuildID  discord.Snowflake   `json:"guild_id"`
	ParentID *discord.Snowflake  `json:"parent_id,omitempty"`
	Type     discord.ChannelType `json:"type"`

	Guild *discord.Guild `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#thread-list-sync
type ThreadListSyncEvent struct {
	GuildID    discord.Snowflake      `json:"guild_id"`
	ChannelIDs []discord.Snowflake    `json:"channel_ids,omitempty"`
	Threads    []discord.Channel      `json:"threads"`
	Members    []discord.ThreadMember `json:"members"`

	Guild *discord.Guild `json:"-"`
}

// ThreadMemberUpdateEvent is dispatched when the current user's thread member is updated.
// https://docs.discord.com/developers/events/gateway-events#thread-member-update
type ThreadMemberUpdateEvent struct {
	NewMember discord.ThreadMember `json:"-"`
	GuildID   discord.Snowflake    `json:"-"`

	Guild *discord.Guild `json:"-"`
}

// ThreadMembersUpdateEvent is dispatched when thread membership for any user changes.
// https://docs.discord.com/developers/events/gateway-events#thread-members-update
type ThreadMembersUpdateEvent struct {
	ID               discord.Snowflake      `json:"id"`
	GuildID          discord.Snowflake      `json:"guild_id"`
	MemberCount      int                    `json:"member_count"`
	AddedMembers     []discord.ThreadMember `json:"added_members,omitempty"`
	RemovedMemberIDs []discord.Snowflake    `json:"removed_member_ids,omitempty"`

	Guild *discord.Guild `json:"-"`
}

func (e *ThreadCreateEvent) DesiredEventType() Event { return &ThreadCreateEvent{} }
func (e *ThreadCreateEvent) Event() EventType        { return EventThreadCreate }
func (e *ThreadCreateEvent) UnmarshalJSON(data []byte) error {
	var wire struct {
		NewlyCreated *bool `json:"newly_created,omitempty"`
		discord.Channel
	}

	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	e.NewlyCreated = wire.NewlyCreated
	e.Thread = wire.Channel

	return nil
}

func (e *ThreadUpdateEvent) DesiredEventType() Event { return &ThreadUpdateEvent{} }
func (e *ThreadUpdateEvent) Event() EventType        { return EventThreadUpdate }
func (e *ThreadUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewThread)
}

func (e *ThreadDeleteEvent) DesiredEventType() Event { return &ThreadDeleteEvent{} }
func (e *ThreadDeleteEvent) Event() EventType        { return EventThreadDelete }

func (e *ThreadListSyncEvent) DesiredEventType() Event { return &ThreadListSyncEvent{} }
func (e *ThreadListSyncEvent) Event() EventType        { return EventThreadListSync }

func (e *ThreadMemberUpdateEvent) DesiredEventType() Event { return &ThreadMemberUpdateEvent{} }
func (e *ThreadMemberUpdateEvent) Event() EventType        { return EventThreadMemberUpdate }
func (e *ThreadMemberUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		discord.ThreadMember
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewMember = w.ThreadMember
	return nil
}

func (e *ThreadMembersUpdateEvent) DesiredEventType() Event { return &ThreadMembersUpdateEvent{} }
func (e *ThreadMembersUpdateEvent) Event() EventType        { return EventThreadMembersUpdate }
