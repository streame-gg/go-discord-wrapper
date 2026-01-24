package types

var EventFactories = map[string]func() Event{
	"MESSAGE_CREATE": func() Event {
		return &MessageCreateEvent{}
	},
	"READY": func() Event {
		return &ReadyEvent{}
	},
	"GUILD_CREATE": func() Event {
		return &GuildCreateEvent{}
	},
	"INTERACTION_CREATE": func() Event {
		return &InteractionCreateEvent{}
	},
}

type Event interface {
	Event() EventType
}

type EventType string

const (
	EventMessageCreate     EventType = "MESSAGE_CREATE"
	EventReady             EventType = "READY"
	EventGuildCreate       EventType = "GUILD_CREATE"
	EventInteractionCreate EventType = "INTERACTION_CREATE"
)

type MessageCreateEvent struct {
	Message
	GuildID  *string      `json:"guild_id"`
	Member   *GuildMember `json:"member,omitempty"`
	Mentions *[]User      `json:"mentions"`
}

func (e MessageCreateEvent) Event() EventType {
	return EventMessageCreate
}

type GuildCreateEvent struct {
	AnyGuildWrapper
	Large       bool  `json:"large"`
	Unavailable *bool `json:"unavailable"`
	MemberCount int   `json:"member_count"`
}

func (e GuildCreateEvent) Event() EventType {
	return EventGuildCreate
}

type ReadyEvent struct {
	User             User              `json:"user"`
	SessionID        string            `json:"session_id"`
	ResumeGatewayURL string            `json:"resume_gateway_url"`
	Shard            []int             `json:"shard,omitempty"`
	Guilds           []AnyGuildWrapper `json:"guilds"`
}

func (e ReadyEvent) Event() EventType {
	return EventReady
}

type InteractionCreateEvent struct {
	Interaction
}

func (e InteractionCreateEvent) Event() EventType {
	return EventInteractionCreate
}

func (e InteractionCreateEvent) IsCommand() bool {
	return e.Type == InteractionTypeApplicationCommand
}

func (e InteractionCreateEvent) IsButton() bool {
	if e.Type != InteractionTypeMessageComponent {
		return false
	}

	return e.Data.(*InteractionDataMessageComponent).ComponentType == ComponentTypeButton
}

func (e InteractionCreateEvent) IsAnySelectMenu() bool {
	if e.Type != InteractionTypeMessageComponent {
		return false
	}

	return e.Data.(*InteractionDataMessageComponent).ComponentType.IsAnySelectMenu()
}

func (e InteractionCreateEvent) IsAutocomplete() bool {
	return e.Type == InteractionTypeApplicationCommandAutocomplete
}

func (e InteractionCreateEvent) IsModalSubmit() bool {
	return e.Type == InteractionTypeModalSubmit
}
