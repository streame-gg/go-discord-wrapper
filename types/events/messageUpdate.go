package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type MessageUpdateEvent struct {
	NewMessage discord.Message      `json:"-"`
	OldMessage *discord.Message     `json:"-"`
	GuildID    *discord.Snowflake   `json:"-"`
	Member     *discord.GuildMember `json:"-"`
}

func (e *MessageUpdateEvent) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.NewMessage); err != nil {
		return err
	}
	type ctx struct {
		GuildID *discord.Snowflake   `json:"guild_id"`
		Member  *discord.GuildMember `json:"member,omitempty"`
	}
	var c ctx
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	e.GuildID = c.GuildID
	e.Member = c.Member
	return nil
}

func (e MessageUpdateEvent) MarshalJSON() ([]byte, error) {
	msgBytes, err := json.Marshal(&e.NewMessage)
	if err != nil {
		return nil, err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(msgBytes, &m); err != nil {
		return nil, err
	}
	if e.GuildID != nil {
		b, err := json.Marshal(e.GuildID)
		if err != nil {
			return nil, err
		}
		m["guild_id"] = b
	}
	if e.Member != nil {
		b, err := json.Marshal(e.Member)
		if err != nil {
			return nil, err
		}
		m["member"] = b
	}
	if e.OldMessage != nil {
		b, _ := json.Marshal(e.OldMessage)
		m["old_message"] = b
	}
	return json.Marshal(m)
}

func init() { RegisterEvent(MessageUpdateEvent{}) }

func (e MessageUpdateEvent) DesiredEventType() Event { return &MessageUpdateEvent{} }
func (e MessageUpdateEvent) Event() EventType        { return EventMessageUpdate }
