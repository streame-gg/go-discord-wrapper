package common

type Emoji struct {
	ID            Snowflake `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Roles         []string  `json:"roles,omitempty"`
	Users         *User     `json:"users,omitempty"`
	RequireColons *bool     `json:"require_colons,omitempty"`
	Managed       *bool     `json:"managed,omitempty"`
	Animated      bool      `json:"animated,omitempty"`
	Available     *bool     `json:"available,omitempty"`
}

type ReactionCountDetails struct {
	Burst  int `json:"burst,omitempty"`
	Normal int `json:"normal,omitempty"`
}

type Reaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details,omitempty"`
	Me           bool                 `json:"me"`
	MeBurst      bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []interface{}        `json:"burst_colors,omitempty"`
}
