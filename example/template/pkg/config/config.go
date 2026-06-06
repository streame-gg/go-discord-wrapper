// Package config holds the template's compile-time settings.
//
// Unlike the bot token (read from the environment in main.go), these are values
// a developer edits directly in source: who may run privileged commands, and
// which guild to register commands in while developing. Edit Current below.
package config

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// Config is the in-code configuration.
type Config struct {
	// OwnerID is the Discord user ID allowed to run /dev. Enable Developer Mode
	// in Discord, right-click your name, and choose "Copy User ID". Leave it 0
	// to keep /dev locked for everyone.
	OwnerID discord.Snowflake

	// DevGuild, when non-zero, scopes slash-command registration to that single
	// guild. Guild commands update instantly, which is what you want while
	// iterating with /dev reload. Leave it 0 to register globally — Discord can
	// take up to an hour to propagate global command changes.
	DevGuild discord.Snowflake
}

// Current holds the active configuration. Edit these values for your bot.
var Current = Config{
	OwnerID:  0, // e.g. 123456789012345678
	DevGuild: 0, // e.g. 987654321098765432
}

// IsOwner reports whether id is the configured owner. It is always false when
// no owner is set, so /dev stays locked down until you configure one.
func (c Config) IsOwner(id discord.Snowflake) bool {
	return c.OwnerID != 0 && id == c.OwnerID
}
