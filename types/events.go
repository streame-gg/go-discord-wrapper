package types

var EventFactories = map[string]func() DiscordEvent{
	"MESSAGE_CREATE": func() DiscordEvent {
		return &DiscordMessageCreateEvent{}
	},
	"READY": func() DiscordEvent {
		return &DiscordReadyPayload{}
	},
	"GUILD_CREATE": func() DiscordEvent {
		return &DiscordGuildCreateEvent{}
	},
	// add more here
}

type DiscordEvent interface {
	Event() DiscordEventType
}

type DiscordEventType string

const (
	DiscordEventMessageCreate DiscordEventType = "MESSAGE_CREATE"
	DiscordEventReady         DiscordEventType = "READY"
	DiscordEventGuildCreate   DiscordEventType = "GUILD_CREATE"
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
