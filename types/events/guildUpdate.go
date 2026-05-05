package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type GuildUpdateEvent struct {
	common.Guild
}

func init() { RegisterEvent(GuildUpdateEvent{}) }

func (e GuildUpdateEvent) DesiredEventType() Event { return &GuildUpdateEvent{} }
func (e GuildUpdateEvent) Event() EventType        { return EventGuildUpdate }
