package discord

// ApplicationCommandPermissionType specifies whether a permission applies to a role, user, or channel.
type ApplicationCommandPermissionType int

const (
	ApplicationCommandPermissionTypeRole    ApplicationCommandPermissionType = 1
	ApplicationCommandPermissionTypeUser    ApplicationCommandPermissionType = 2
	ApplicationCommandPermissionTypeChannel ApplicationCommandPermissionType = 3
)

// ApplicationCommandPermission grants or denies a command for a specific role, user, or channel.
type ApplicationCommandPermission struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}

// GuildApplicationCommandPermissions holds the permission overrides for a command in a guild.
type GuildApplicationCommandPermissions struct {
	ID            Snowflake                      `json:"id"`
	ApplicationID Snowflake                      `json:"application_id"`
	GuildID       Snowflake                      `json:"guild_id"`
	Permissions   []ApplicationCommandPermission `json:"permissions"`
}
