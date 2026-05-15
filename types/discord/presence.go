package discord

type ActivityType int

const (
	ActivityTypePlaying   ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

type ActivityTimestamps struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

type ActivityParty struct {
	ID   *string `json:"id,omitempty"`
	Size *[2]int `json:"size,omitempty"`
}

type ActivityAssets struct {
	LargeImage *string `json:"large_image,omitempty"`
	LargeText  *string `json:"large_text,omitempty"`
	SmallImage *string `json:"small_image,omitempty"`
	SmallText  *string `json:"small_text,omitempty"`
}

type ActivitySecrets struct {
	Join     *string `json:"join,omitempty"`
	Spectate *string `json:"spectate,omitempty"`
	Match    *string `json:"match,omitempty"`
}

type ActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type ActivityEmoji struct {
	Name     string     `json:"name"`
	ID       *Snowflake `json:"id,omitempty"`
	Animated *bool      `json:"animated,omitempty"`
}

type ActivityFlags int

const (
	ActivityFlagInstance                 ActivityFlags = 1 << 0
	ActivityFlagJoin                     ActivityFlags = 1 << 1
	ActivityFlagSpectate                 ActivityFlags = 1 << 2
	ActivityFlagJoinRequest              ActivityFlags = 1 << 3
	ActivityFlagSync                     ActivityFlags = 1 << 4
	ActivityFlagPlay                     ActivityFlags = 1 << 5
	ActivityFlagPartyPrivacyFriends      ActivityFlags = 1 << 6
	ActivityFlagPartyPrivacyVoiceChannel ActivityFlags = 1 << 7
	ActivityFlagEmbedded                 ActivityFlags = 1 << 8
)

type FullActivity struct {
	Name          string              `json:"name"`
	Type          ActivityType        `json:"type"`
	URL           *string             `json:"url,omitempty"`
	CreatedAt     int64               `json:"created_at"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID *Snowflake          `json:"application_id,omitempty"`
	Details       *string             `json:"details,omitempty"`
	State         *string             `json:"state,omitempty"`
	Emoji         *ActivityEmoji      `json:"emoji,omitempty"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Instance      *bool               `json:"instance,omitempty"`
	Flags         *ActivityFlags      `json:"flags,omitempty"`
	Buttons       []ActivityButton    `json:"buttons,omitempty"`
}

type ClientStatus struct {
	Desktop *string `json:"desktop,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Web     *string `json:"web,omitempty"`
}

type PresenceStatus string

const (
	PresenceStatusOnline    PresenceStatus = "online"
	PresenceStatusDND       PresenceStatus = "dnd"
	PresenceStatusIdle      PresenceStatus = "idle"
	PresenceStatusOffline   PresenceStatus = "offline"
	PresenceStatusInvisible PresenceStatus = "invisible"
)

type PartialPresenceUser struct {
	ID            Snowflake `json:"id"`
	Username      *string   `json:"username,omitempty"`
	Discriminator *string   `json:"discriminator,omitempty"`
	GlobalName    *string   `json:"global_name,omitempty"`
	AvatarHash    *string   `json:"avatar,omitempty"`
	Bot           *bool     `json:"bot,omitempty"`
	PublicFlags   *int      `json:"public_flags,omitempty"`
}

// Presence is the cached form of a user's presence in a guild.
// Populated from GUILD_CREATE (initial state) and PRESENCE_UPDATE events.
type Presence struct {
	User         PartialPresenceUser `json:"user"`
	GuildID      Snowflake           `json:"guild_id"`
	Status       PresenceStatus      `json:"status"`
	Activities   []FullActivity      `json:"activities"`
	ClientStatus ClientStatus        `json:"client_status"`
}
