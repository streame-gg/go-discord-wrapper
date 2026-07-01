package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &ReadyEvent{}
var _ Event = &ResumedEvent{}

func init() {
	RegisterEvent(&ReadyEvent{})
	RegisterEvent(&ResumedEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#ready
type ReadyEvent struct {
	V                discord.APIVersion         `json:"v"`
	User             discord.User               `json:"user"`
	Application      ReadyApplication           `json:"application"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL string                     `json:"resume_gateway_url"`
	Shard            *ReadyShard                `json:"shard,omitempty"`
	Guilds           []discord.UnavailableGuild `json:"guilds"`
}

// https://docs.discord.com/developers/events/gateway-events#ready
type ReadyApplication struct {
	ID    discord.Snowflake `json:"id"`
	Flags int               `json:"flags"`
}

// ReadyShard carries the shard identity from the READY payload.
// https://docs.discord.com/developers/events/gateway-events#ready
type ReadyShard struct {
	ShardID   int
	NumShards int
}

func (s *ReadyShard) UnmarshalJSON(data []byte) error {
	var pair [2]int
	if err := json.Unmarshal(data, &pair); err != nil {
		return err
	}
	s.ShardID = pair[0]
	s.NumShards = pair[1]
	return nil
}

// https://docs.discord.com/developers/events/gateway-events#resumed
type ResumedEvent struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

func (e *ReadyEvent) DesiredEventType() Event {
	return &ReadyEvent{}
}
func (e *ReadyEvent) Event() EventType { return EventReady }

func (e *ResumedEvent) DesiredEventType() Event { return &ResumedEvent{} }
func (e *ResumedEvent) Event() EventType        { return EventResumed }
