package discord

// OnboardingMode controls whether prompts are shown during onboarding and member screening.
//
// https://docs.discord.com/developers/resources/guild#guild-onboarding-object-onboarding-mode
type OnboardingMode int

const (
	OnboardingModeDefault  OnboardingMode = 0
	OnboardingModeAdvanced OnboardingMode = 1
)

// OnboardingPromptType is the display type of onboarding prompt.
//
// https://docs.discord.com/developers/resources/guild#guild-onboarding-object-prompt-types
type OnboardingPromptType int

const (
	OnboardingPromptTypeMultipleChoice OnboardingPromptType = 0
	OnboardingPromptTypeDropdown       OnboardingPromptType = 1
)

// PromptOption is a selectable option within a guild onboarding prompt.
//
// https://docs.discord.com/developers/resources/guild#guild-onboarding-object-prompt-option-structure
type PromptOption struct {
	ID            Snowflake   `json:"id"`
	ChannelIDs    []Snowflake `json:"channel_ids"`
	RoleIDs       []Snowflake `json:"role_ids"`
	Emoji         *Emoji      `json:"emoji,omitempty"`
	EmojiID       *Snowflake  `json:"emoji_id,omitempty"`
	EmojiName     string      `json:"emoji_name,omitempty"`
	EmojiAnimated *bool       `json:"emoji_animated,omitempty"`
	Title         string      `json:"title"`
	Description   *string     `json:"description"`
}

// OnboardingPrompt is a prompt shown to new members during guild onboarding.
//
// https://docs.discord.com/developers/resources/guild#guild-onboarding-object-onboarding-prompt-structure
type OnboardingPrompt struct {
	ID           Snowflake            `json:"id"`
	Type         OnboardingPromptType `json:"type"`
	Options      []PromptOption       `json:"options"`
	Title        string               `json:"title"`
	SingleSelect bool                 `json:"single_select"`
	Required     bool                 `json:"required"`
	InOnboarding bool                 `json:"in_onboarding"`
}

// GuildOnboarding holds the onboarding configuration for a guild.
//
// https://docs.discord.com/developers/resources/guild#guild-onboarding-object
type GuildOnboarding struct {
	GuildID           Snowflake          `json:"guild_id"`
	Prompts           []OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []Snowflake        `json:"default_channel_ids"`
	Enabled           bool               `json:"enabled"`
	Mode              OnboardingMode     `json:"mode"`
}
