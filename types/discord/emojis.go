package discord

// https://docs.discord.com/developers/resources/emoji#emoji-object
type Emoji struct {
	hClient       EntityClient
	GuildID       Snowflake   `json:"-"`
	ID            Snowflake   `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Roles         []Snowflake `json:"roles,omitempty"`
	Users         []User      `json:"users,omitempty"`
	RequireColons *bool       `json:"require_colons,omitempty"`
	Managed       *bool       `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
	Available     *bool       `json:"available,omitempty"`
}

// https://docs.discord.com/developers/resources/message#reaction-count-details-object
type ReactionCountDetails struct {
	Burst  int `json:"burst,omitempty"`
	Normal int `json:"normal,omitempty"`
}

// https://docs.discord.com/developers/resources/message#reaction-object
type Reaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details,omitempty"`
	Me           bool                 `json:"me"`
	MeBurst      bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []interface{}        `json:"burst_colors,omitempty"`
}
