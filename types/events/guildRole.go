package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildRoleCreateEvent{}
var _ Event = &GuildRoleUpdateEvent{}
var _ Event = &GuildRoleDeleteEvent{}

func init() {
	RegisterEvent(&GuildRoleCreateEvent{})
	RegisterEvent(&GuildRoleUpdateEvent{})
	RegisterEvent(&GuildRoleDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#guild-role-create
type GuildRoleCreateEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-role-update
type GuildRoleUpdateEvent struct {
	GuildID discord.Snowflake `json:"-"`
	NewRole discord.Role      `json:"-"`
	OldRole *discord.Role     `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-role-delete
type GuildRoleDeleteEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	RoleID  discord.Snowflake `json:"role_id"`
	Role    *discord.Role     `json:"-"`
}

func (e *GuildRoleCreateEvent) DesiredEventType() Event { return &GuildRoleCreateEvent{} }
func (e *GuildRoleCreateEvent) Event() EventType        { return EventGuildRoleCreate }

func (e *GuildRoleUpdateEvent) DesiredEventType() Event { return &GuildRoleUpdateEvent{} }
func (e *GuildRoleUpdateEvent) Event() EventType        { return EventGuildRoleUpdate }
func (e *GuildRoleUpdateEvent) UnmarshalJSON(data []byte) error {
	var wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		Role    discord.Role      `json:"role"`
	}

	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	e.GuildID = wire.GuildID
	e.NewRole = wire.Role

	return nil
}

func (e *GuildRoleDeleteEvent) DesiredEventType() Event { return &GuildRoleDeleteEvent{} }
func (e *GuildRoleDeleteEvent) Event() EventType        { return EventGuildRoleDelete }
