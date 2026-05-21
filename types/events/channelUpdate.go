package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ChannelUpdateEvent struct {
	NewChannel discord.Channel  `json:"-"`
	OldChannel *discord.Channel `json:"-"`
}

func (e *ChannelUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewChannel)
}

func (e ChannelUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.Channel
		OldChannel *discord.Channel `json:"old_channel,omitempty"`
	}
	return json.Marshal(wire{e.NewChannel, e.OldChannel})
}

func init() { RegisterEvent(ChannelUpdateEvent{}) }

func (e ChannelUpdateEvent) DesiredEventType() Event { return &ChannelUpdateEvent{} }
func (e ChannelUpdateEvent) Event() EventType        { return EventChannelUpdate }
