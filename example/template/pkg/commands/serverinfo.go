package commands

import (
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/builder"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	dcmd "github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(serverinfo{}) }

// serverinfo reads counts straight out of the cache — no REST calls — and
// replies with an embed. It demonstrates reading cached state inside a command.
type serverinfo struct{}

func (serverinfo) Definition() *dcmd.ApplicationCommand {
	return &dcmd.ApplicationCommand{
		Name:        "serverinfo",
		Description: "Show cached stats about this server",
		Type:        discord.ApplicationCommandTypeChatInput,
	}
}

func (serverinfo) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	if ev.GuildID == nil {
		_, _ = ev.Reply(interactions.ReplyOptions{
			Content: "Run this inside a server.",
			Flags:   discord.MessageFlagEphemeral,
		})
		return
	}

	gid := *ev.GuildID
	members, channels, roles := 0, c.ChannelsForGuild(gid).Len(), 0
	if c.Cache != nil {
		members = c.Cache.Members().AllInGuild(gid).Len()
		roles = c.Cache.Roles().GetByGuild(gid).Len()
	}

	inline := options.Ptr(true)
	embed := builder.NewEmbed().
		SetTitle("Server info").
		SetColor(0x5865F2).
		AddFields(
			discord.EmbedFields{Name: "Cached members", Value: strconv.Itoa(members), Inline: inline},
			discord.EmbedFields{Name: "Channels", Value: strconv.Itoa(channels), Inline: inline},
			discord.EmbedFields{Name: "Roles", Value: strconv.Itoa(roles), Inline: inline},
		).
		SetFooter("Counts reflect what the cache has seen so far", "").
		Build()

	_, _ = ev.Reply(interactions.ReplyOptions{Embeds: []discord.Embed{embed}})
}
