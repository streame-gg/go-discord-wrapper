package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ApplicationCommandPermissionsUpdateEvent is dispatched when command permissions are updated in a guild.
type ApplicationCommandPermissionsUpdateEvent struct {
	NewPermissions discord.GuildApplicationCommandPermissions `json:"-"`
}

func (e *ApplicationCommandPermissionsUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewPermissions)
}

func (e ApplicationCommandPermissionsUpdateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&e.NewPermissions)
}

func init() { RegisterEvent(ApplicationCommandPermissionsUpdateEvent{}) }

func (e ApplicationCommandPermissionsUpdateEvent) DesiredEventType() Event {
	return &ApplicationCommandPermissionsUpdateEvent{}
}
func (e ApplicationCommandPermissionsUpdateEvent) Event() EventType {
	return EventApplicationCommandPermissionsUpdate
}
