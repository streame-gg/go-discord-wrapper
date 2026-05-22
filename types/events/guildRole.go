package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildRoleCreateEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

type GuildRoleUpdateEvent struct {
	GuildID discord.Snowflake `json:"-"`
	NewRole discord.Role      `json:"-"`
	OldRole *discord.Role     `json:"-"`
}

func (e *GuildRoleUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		Role    discord.Role      `json:"role"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewRole = w.Role
	return nil
}

func (e GuildRoleUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		Role    discord.Role      `json:"role"`
		OldRole *discord.Role     `json:"old_role,omitempty"`
	}
	return json.Marshal(wire{e.GuildID, e.NewRole, e.OldRole})
}

type GuildRoleDeleteEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	RoleID  discord.Snowflake `json:"role_id"`
}

func init() {
	RegisterEvent(GuildRoleCreateEvent{})
	RegisterEvent(GuildRoleUpdateEvent{})
	RegisterEvent(GuildRoleDeleteEvent{})
}

func (e GuildRoleCreateEvent) DesiredEventType() Event { return &GuildRoleCreateEvent{} }
func (e GuildRoleCreateEvent) Event() EventType        { return EventGuildRoleCreate }

func (e GuildRoleUpdateEvent) DesiredEventType() Event { return &GuildRoleUpdateEvent{} }
func (e GuildRoleUpdateEvent) Event() EventType        { return EventGuildRoleUpdate }

func (e GuildRoleDeleteEvent) DesiredEventType() Event { return &GuildRoleDeleteEvent{} }
func (e GuildRoleDeleteEvent) Event() EventType        { return EventGuildRoleDelete }
