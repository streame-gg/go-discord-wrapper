package discord

// https://docs.discord.com/developers/interactions/application-commands#application-command-permissions-object-application-command-permission-type
type ApplicationCommandPermissionType int

const (
	ApplicationCommandPermissionTypeRole    ApplicationCommandPermissionType = 1
	ApplicationCommandPermissionTypeUser    ApplicationCommandPermissionType = 2
	ApplicationCommandPermissionTypeChannel ApplicationCommandPermissionType = 3
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-permissions-object-application-command-permissions-structure
type ApplicationCommandPermission struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}

// https://docs.discord.com/developers/interactions/application-commands#application-command-permissions-object
type GuildApplicationCommandPermissions struct {
	ID            Snowflake                      `json:"id"`
	ApplicationID Snowflake                      `json:"application_id"`
	GuildID       Snowflake                      `json:"guild_id"`
	Permissions   []ApplicationCommandPermission `json:"permissions"`
}
