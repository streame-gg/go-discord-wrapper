package common

type AllowedMentionsType string

const (
	AllowedMentionsTypeRoles    AllowedMentionsType = "roles"
	AllowedMentionsTypeUsers    AllowedMentionsType = "users"
	AllowedMentionsTypeEveryone AllowedMentionsType = "everyone"
)

type AllowedMentions struct {
	Parse       *[]AllowedMentionsType `json:"parse,omitempty"`
	Roles       *[]Snowflake           `json:"roles,omitempty"`
	Users       *[]Snowflake           `json:"users,omitempty"`
	RepliedUser *bool                  `json:"replied_user,omitempty"`
}
