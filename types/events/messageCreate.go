package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type MessageCreateEvent struct {
	discord.Message
	GuildID *discord.Snowflake   `json:"-"`
	Member  *discord.GuildMember `json:"-"`
}

// UnmarshalJSON decodes both the embedded Message (which has its own custom
// UnmarshalJSON) and the outer guild_id / member fields that Discord appends
// to MESSAGE_CREATE but are not part of the Message object itself.
func (e *MessageCreateEvent) UnmarshalJSON(data []byte) error {
	if err := e.Message.UnmarshalJSON(data); err != nil {
		return err
	}
	var extra struct {
		GuildID *discord.Snowflake   `json:"guild_id"`
		Member  *discord.GuildMember `json:"member,omitempty"`
	}
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}
	e.GuildID = extra.GuildID
	e.Member = extra.Member
	return nil
}

func init() { RegisterEvent(MessageCreateEvent{}) }

func (e MessageCreateEvent) DesiredEventType() Event {
	return &MessageCreateEvent{}
}

func (e MessageCreateEvent) Event() EventType {
	return EventMessageCreate
}
