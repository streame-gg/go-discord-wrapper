// Package commands holds every slash command the bot exposes, plus a small
// self-registering registry that auto-loads them.
//
// Each command lives in its own file and registers itself from an init():
//
//	func init() { Register(&ping{}) }
//
// Because all command files share this package, importing the package once
// (the bot package does) runs every init() and populates the registry — so
// adding a command is just dropping a new file in this folder. No central list
// to keep in sync.
package commands

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// Command is one slash command: a definition used to register it with Discord,
// and a handler invoked when a user runs it.
type Command interface {
	// Definition is the registration payload (name, description, options).
	Definition() *discord.ApplicationCommand
	// Handle runs when the command is invoked. ev is guaranteed to be a
	// chat-input command interaction for this command's top-level name.
	Handle(c *connection.Client, ev *events.InteractionCreateEvent)
}

// registry maps a command's top-level name to its implementation.
var registry = map[string]Command{}

// Register adds cmd to the registry. Call it from an init() in each command
// file. Panics on a duplicate name so collisions surface at startup.
func Register(cmd Command) {
	name := cmd.Definition().Name
	if _, dup := registry[name]; dup {
		panic("commands: duplicate command name " + name)
	}
	registry[name] = cmd
}

// Lookup returns the command registered under name, if any.
func Lookup(name string) (Command, bool) {
	c, ok := registry[name]
	return c, ok
}

// Definitions returns the registration payloads for every command, ready to
// pass to BulkRegisterCommands.
func Definitions() []*discord.ApplicationCommand {
	out := make([]*discord.ApplicationCommand, 0, len(registry))
	for _, c := range registry {
		out = append(out, c.Definition())
	}
	return out
}

// Sync pushes the current command set to Discord and returns the number
// registered. When guild is non-zero it registers to that single guild (changes
// apply instantly — ideal while developing); otherwise it registers globally,
// which Discord can take up to an hour to propagate. appID is the bot's
// application ID, available from the READY event or any interaction.
//
// Both startup and /dev reload go through here so the two never disagree about
// where commands live.
func Sync(ctx context.Context, c *connection.Client, appID, guild discord.Snowflake) (int, error) {
	defs := Definitions()
	if !guild.IsEmpty() {
		out, err := c.RestClient.BulkOverwriteGuildApplicationCommands(ctx, appID, guild, defs)
		return len(out), err
	}
	out, err := c.BulkRegisterCommands(ctx, defs)
	return len(out), err
}
