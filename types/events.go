package types

var EventFactories = map[string]func() DiscordEvent{
	"MESSAGE_CREATE": func() DiscordEvent {
		return &DiscordMessageCreateEvent{}
	},
	"READY": func() DiscordEvent {
		return &DiscordReadyEvent{}
	},
	"GUILD_CREATE": func() DiscordEvent {
		return &DiscordGuildCreateEvent{}
	},
	"INTERACTION_CREATE": func() DiscordEvent {
		return &DiscordInteractionCreateEvent{}
	},
	// add more here
}

type DiscordEvent interface {
	Event() DiscordEventType
}

type DiscordEventType string

const (
	DiscordEventMessageCreate     DiscordEventType = "MESSAGE_CREATE"
	DiscordEventReady             DiscordEventType = "READY"
	DiscordEventGuildCreate       DiscordEventType = "GUILD_CREATE"
	DiscordEventInteractionCreate DiscordEventType = "INTERACTION_CREATE"
)

type DiscordMessageCreateEvent struct {
	DiscordMessage
	GuildID  *string        `json:"guild_id"`
	Member   *GuildMember   `json:"member,omitempty"`
	Mentions *[]DiscordUser `json:"mentions"`
}

func (e DiscordMessageCreateEvent) Event() DiscordEventType {
	return DiscordEventMessageCreate
}

type DiscordGuildCreateEvent struct {
	AnyGuildWrapper
	Large       bool  `json:"large"`
	Unavailable *bool `json:"unavailable"`
	MemberCount int   `json:"member_count"`
}

func (e DiscordGuildCreateEvent) Event() DiscordEventType {
	return DiscordEventGuildCreate
}

type DiscordReadyEvent struct {
	User             DiscordUser       `json:"user"`
	SessionID        string            `json:"session_id"`
	ResumeGatewayURL string            `json:"resume_gateway_url"`
	Shard            []int             `json:"shard,omitempty"`
	Guilds           []AnyGuildWrapper `json:"guilds"`
}

func (e DiscordReadyEvent) Event() DiscordEventType {
	return DiscordEventReady
}

type DiscordInteractionCreateEvent struct {
	DiscordInteraction
}

func (e DiscordInteractionCreateEvent) Event() DiscordEventType {
	return DiscordEventInteractionCreate
}

func (e DiscordInteractionCreateEvent) IsCommand() bool {
	return e.Type == DiscordInteractionTypeApplicationCommand
}

func (e DiscordInteractionCreateEvent) IsButton() bool {
	if e.Type != DiscordInteractionTypeMessageComponent {
		return false
	}

	return e.Data.(*DiscordInteractionDataMessageComponent).ComponentType == DiscordComponentTypeButton
}

func (e DiscordInteractionCreateEvent) IsAnySelectMenu() bool {
	if e.Type != DiscordInteractionTypeMessageComponent {
		return false
	}

	return e.Data.(*DiscordInteractionDataMessageComponent).ComponentType.IsAnySelectMenu()
}

func (e DiscordInteractionCreateEvent) IsAutocomplete() bool {
	return e.Type == DiscordInteractionTypeApplicationCommandAutocomplete
}

func (e DiscordInteractionCreateEvent) IsModalSubmit() bool {
	return e.Type == DiscordInteractionTypeModalSubmit
}
