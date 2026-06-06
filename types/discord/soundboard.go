package discord

// SoundboardSound represents a sound in a guild's soundboard or the default Discord soundboard.
// https://docs.discord.com/developers/resources/soundboard#soundboard-sound-object
type SoundboardSound struct {
	hClient   EntityClient
	Name      string     `json:"name"`
	SoundID   Snowflake  `json:"sound_id"`
	Volume    float64    `json:"volume"`
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
	GuildID   *Snowflake `json:"guild_id,omitempty"`
	Available bool       `json:"available"`
	User      *User      `json:"user,omitempty"`
}
