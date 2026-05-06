package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

// ApplicationCommandPermissionsUpdateEvent is dispatched when command permissions are updated in a guild.
type ApplicationCommandPermissionsUpdateEvent struct {
	common.GuildApplicationCommandPermissions
}

func init() { RegisterEvent(ApplicationCommandPermissionsUpdateEvent{}) }

func (e ApplicationCommandPermissionsUpdateEvent) DesiredEventType() Event {
	return &ApplicationCommandPermissionsUpdateEvent{}
}
func (e ApplicationCommandPermissionsUpdateEvent) Event() EventType {
	return EventApplicationCommandPermissionsUpdate
}
