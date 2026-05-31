package discord

type AllowedMentionsType string

const (
	AllowedMentionsTypeRoles    AllowedMentionsType = "roles"
	AllowedMentionsTypeUsers    AllowedMentionsType = "users"
	AllowedMentionsTypeEveryone AllowedMentionsType = "everyone"
)

type AllowedMentions struct {
	// Parse is intentionally a pointer-to-slice, not a plain slice.
	// Discord distinguishes between three states:
	//   nil       → omit the "parse" key entirely (no override; Discord applies its defaults)
	//   &[]{}     → send "parse": [] (suppress all auto-mentions)
	//   &[]{"roles", ...} → allow only the listed mention types
	// A plain []AllowedMentionsType with `omitempty` cannot represent the empty-slice
	// case because json.Marshal omits zero-length slices, making it indistinguishable
	// from the nil / omitted case (same design rationale as Bug #70).
	Parse       *[]AllowedMentionsType `json:"parse,omitempty"`
	Roles       *[]Snowflake           `json:"roles,omitempty"`
	Users       *[]Snowflake           `json:"users,omitempty"`
	RepliedUser *bool                  `json:"replied_user,omitempty"`
}
