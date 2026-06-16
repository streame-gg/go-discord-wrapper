package discord

// https://docs.discord.com/developers/resources/message#allowed-mentions-object-allowed-mention-types
type AllowedMentionsType string

const (
	AllowedMentionsTypeRoles    AllowedMentionsType = "roles"
	AllowedMentionsTypeUsers    AllowedMentionsType = "users"
	AllowedMentionsTypeEveryone AllowedMentionsType = "everyone"
)

// https://docs.discord.com/developers/resources/message#allowed-mentions-object
type AllowedMentions struct {
	Parse       *[]AllowedMentionsType `json:"parse,omitempty"`
	Roles       []Snowflake            `json:"roles,omitempty"`
	Users       []Snowflake            `json:"users,omitempty"`
	RepliedUser bool                   `json:"replied_user,omitempty"`
}
