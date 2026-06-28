package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func init() {
	RegisterEvent(&MessageCreateEvent{})
	RegisterEvent(&MessageUpdateEvent{})
	RegisterEvent(&MessageDeleteEvent{})
	RegisterEvent(&MessageDeleteBulkEvent{})

	RegisterEvent(&MessagePollVoteAddEvent{})
	RegisterEvent(&MessagePollVoteRemoveEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#message-create
type MessageCreateEvent struct {
	Message     discord.Message      `json:"-"`
	GuildID     *discord.Snowflake   `json:"-"`
	Member      *discord.GuildMember `json:"-"`
	ChannelType *discord.ChannelType `json:"-"`
	Mentions    []discord.User       `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-update
type MessageUpdateEvent struct {
	NewMessage  discord.Message      `json:"-"`
	OldMessage  *discord.Message     `json:"-"`
	GuildID     *discord.Snowflake   `json:"-"`
	Member      *discord.GuildMember `json:"-"`
	ChannelType *discord.ChannelType `json:"-"`
	Mentions    []discord.User       `json:"-"`

	Guild *discord.Guild `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-delete
type MessageDeleteEvent struct {
	ID        discord.Snowflake  `json:"id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-delete-bulk
type MessageDeleteBulkEvent struct {
	IDs       []discord.Snowflake `json:"ids"`
	ChannelID discord.Snowflake   `json:"channel_id"`
	GuildID   *discord.Snowflake  `json:"guild_id,omitempty"`

	Guild   *discord.Guild   `json:"-"`
	Channel *discord.Channel `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#message-poll-vote-add
type MessagePollVoteAddEvent struct {
	UserID    discord.Snowflake  `json:"user_id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int                `json:"answer_id"`
}

// https://docs.discord.com/developers/events/gateway-events#message-poll-vote-remove
type MessagePollVoteRemoveEvent struct {
	UserID    discord.Snowflake  `json:"user_id"`
	ChannelID discord.Snowflake  `json:"channel_id"`
	MessageID discord.Snowflake  `json:"message_id"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int                `json:"answer_id"`
}

func (e *MessageCreateEvent) DesiredEventType() Event {
	return &MessageCreateEvent{}
}
func (e *MessageCreateEvent) Event() EventType {
	return EventMessageCreate
}
func (e *MessageCreateEvent) UnmarshalJSON(data []byte) error {
	if err := e.Message.UnmarshalJSON(data); err != nil {
		return err
	}
	var extra struct {
		GuildID     *discord.Snowflake   `json:"guild_id"`
		Member      *discord.GuildMember `json:"member,omitempty"`
		Mentions    []discord.User       `json:"mentions"`
		ChannelType *discord.ChannelType `json:"channel_type"`
	}
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	e.GuildID = extra.GuildID
	e.Member = extra.Member
	e.ChannelType = extra.ChannelType
	e.Mentions = extra.Mentions
	return nil
}

func (e *MessageUpdateEvent) DesiredEventType() Event { return &MessageUpdateEvent{} }
func (e *MessageUpdateEvent) Event() EventType        { return EventMessageUpdate }
func (e *MessageUpdateEvent) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.NewMessage); err != nil {
		return err
	}
	type ctx struct {
		GuildID     *discord.Snowflake   `json:"guild_id"`
		Member      *discord.GuildMember `json:"member,omitempty"`
		Mentions    []discord.User       `json:"mentions"`
		ChannelType *discord.ChannelType `json:"channel_type"`
	}
	var c ctx
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	e.GuildID = c.GuildID
	e.Member = c.Member
	e.Mentions = c.Mentions
	e.ChannelType = c.ChannelType

	return nil
}

func (m *MessageDeleteEvent) DesiredEventType() Event {
	return &MessageDeleteEvent{}
}
func (m *MessageDeleteEvent) Event() EventType {
	return EventMessageDelete
}

func (m *MessageDeleteBulkEvent) DesiredEventType() Event {
	return &MessageDeleteBulkEvent{}
}
func (m *MessageDeleteBulkEvent) Event() EventType {
	return EventMessageDeleteBulk
}

func (e *MessagePollVoteAddEvent) DesiredEventType() Event { return &MessagePollVoteAddEvent{} }
func (e *MessagePollVoteAddEvent) Event() EventType        { return EventMessagePollVoteAdd }

func (e *MessagePollVoteRemoveEvent) DesiredEventType() Event { return &MessagePollVoteRemoveEvent{} }
func (e *MessagePollVoteRemoveEvent) Event() EventType        { return EventMessagePollVoteRemove }
