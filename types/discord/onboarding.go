package discord

// OnboardingMode controls whether prompts are shown during onboarding and member screening.
type OnboardingMode int

const (
	OnboardingModeDefault  OnboardingMode = 0
	OnboardingModeAdvanced OnboardingMode = 1
)

// OnboardingPromptType is the display type of an onboarding prompt.
type OnboardingPromptType int

const (
	OnboardingPromptTypeMultipleChoice OnboardingPromptType = 0
	OnboardingPromptTypeDropdown       OnboardingPromptType = 1
)

// PromptOption is a selectable option within a guild onboarding prompt.
type PromptOption struct {
	ID            Snowflake   `json:"id"`
	ChannelIDs    []Snowflake `json:"channel_ids"`
	RoleIDs       []Snowflake `json:"role_ids"`
	Emoji         *Emoji      `json:"emoji,omitempty"`
	EmojiID       *Snowflake  `json:"emoji_id,omitempty"`
	EmojiName     *string     `json:"emoji_name,omitempty"`
	EmojiAnimated *bool       `json:"emoji_animated,omitempty"`
	Title         string      `json:"title"`
	Description   *string     `json:"description,omitempty"`
}

// OnboardingPrompt is a prompt shown to new members during guild onboarding.
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
type GuildOnboarding struct {
	GuildID           Snowflake          `json:"guild_id"`
	Prompts           []OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []Snowflake        `json:"default_channel_ids"`
	Enabled           bool               `json:"enabled"`
	Mode              OnboardingMode     `json:"mode"`
}
