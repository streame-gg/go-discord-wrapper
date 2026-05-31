package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ReadyShard carries the shard identity from the READY payload.
type ReadyShard struct {
	ShardID   int
	NumShards int
}

func (s ReadyShard) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]int{s.ShardID, s.NumShards})
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

type ReadyEvent struct {
	User             discord.User              `json:"user"`
	SessionID        string                    `json:"session_id"`
	ResumeGatewayURL string                    `json:"resume_gateway_url"`
	// Shard is present only when the client was started with sharding.
	Shard            *ReadyShard               `json:"shard,omitempty"`
	Guilds           []discord.AnyGuildWrapper `json:"guilds"`
}

// ResumedEvent is dispatched when a session is successfully resumed.
type ResumedEvent struct{}

func init() {
	RegisterEvent(ReadyEvent{})
	RegisterEvent(ResumedEvent{})
}

func (e ReadyEvent) DesiredEventType() Event {
	return &ReadyEvent{}
}

func (e ReadyEvent) Event() EventType { return EventReady }

func (e ResumedEvent) DesiredEventType() Event { return &ResumedEvent{} }
func (e ResumedEvent) Event() EventType        { return EventResumed }
