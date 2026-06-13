package builder

import (
	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type CommandBuilder struct {
	cmd *commands.ApplicationCommand
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		cmd: &commands.ApplicationCommand{
			Type: discord.ApplicationCommandTypeChatInput,
		},
	}
}

func (b *CommandBuilder) SetName(name string) *CommandBuilder {
	b.cmd.Name = name
	return b
}

func (b *CommandBuilder) SetNameLocalizations(localizations map[discord.Locale]string) *CommandBuilder {
	b.cmd.NameLocalizations = localizations
	return b
}

func (b *CommandBuilder) SetDescription(description string) *CommandBuilder {
	b.cmd.Description = description
	return b
}

func (b *CommandBuilder) SetDescriptionLocalizations(localizations map[discord.Locale]string) *CommandBuilder {
	b.cmd.DescriptionLocalizations = localizations
	return b
}

func (b *CommandBuilder) SetType(t discord.ApplicationCommandType) *CommandBuilder {
	b.cmd.Type = t
	return b
}

func (b *CommandBuilder) SetDefaultMemberPermissions(permissions discord.Permission) *CommandBuilder {
	b.cmd.DefaultMemberPermissions = &permissions
	return b
}

func (b *CommandBuilder) SetIntegrationTypes(types []discord.InteractionApplicationIntegrationType) *CommandBuilder {
	b.cmd.IntegrationTypes = types
	return b
}

func (b *CommandBuilder) SetNSFW(nsfw bool) *CommandBuilder {
	b.cmd.NSFW = nsfw
	return b
}

func (b *CommandBuilder) SetContexts(contexts []discord.InteractionContextType) *CommandBuilder {
	b.cmd.Contexts = contexts
	return b
}

func (b *CommandBuilder) AddOption(option commands.AnyApplicationCommandOption) *CommandBuilder {
	b.cmd.Options = append(b.cmd.Options, option)
	return b
}

func (b *CommandBuilder) AddOptions(option ...commands.AnyApplicationCommandOption) *CommandBuilder {
	b.cmd.Options = append(b.cmd.Options, option...)
	return b
}

func (b *CommandBuilder) SetOptions(options ...commands.AnyApplicationCommandOption) *CommandBuilder {
	b.cmd.Options = options
	return b
}

func (b *CommandBuilder) Build() commands.ApplicationCommand {
	return *b.cmd
}
