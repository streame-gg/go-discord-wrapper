package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// ApplicationCommandPermissionsUpdateEvent is dispatched when command permissions are updated in a guild.
type ApplicationCommandPermissionsUpdateEvent struct {
	discord.GuildApplicationCommandPermissions
}

func init() { RegisterEvent(ApplicationCommandPermissionsUpdateEvent{}) }

func (e ApplicationCommandPermissionsUpdateEvent) DesiredEventType() Event {
	return &ApplicationCommandPermissionsUpdateEvent{}
}
func (e ApplicationCommandPermissionsUpdateEvent) Event() EventType {
	return EventApplicationCommandPermissionsUpdate
}
