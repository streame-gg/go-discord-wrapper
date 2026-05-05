package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type PresenceUpdateEvent struct {
	User         common.PartialPresenceUser `json:"user"`
	GuildID      common.Snowflake           `json:"guild_id"`
	Status       common.PresenceStatus      `json:"status"`
	Activities   []common.FullActivity      `json:"activities"`
	ClientStatus common.ClientStatus        `json:"client_status"`
}

func init() { RegisterEvent(PresenceUpdateEvent{}) }

func (e PresenceUpdateEvent) DesiredEventType() Event { return &PresenceUpdateEvent{} }
func (e PresenceUpdateEvent) Event() EventType        { return EventPresenceUpdate }
