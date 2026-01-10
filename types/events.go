package types

type DiscordEvent interface {
}

type DiscordEventType string

const (
	DiscordEventMessageCreate DiscordEventType = "MESSAGE_CREATE"
	DiscordEventReady         DiscordEventType = "READY"
	DiscordEventGuildCreate   DiscordEventType = "GUILD_CREATE"
)
