package discord

import (
	"strconv"
	"strings"
	"time"
)

// Markdown / mention builders, mirroring discord.js' @discordjs/formatters. These
// are the constructive counterpart to EscapeMarkdown: helpers that wrap text in
// Discord's markup so handlers don't hand-assemble it. Inputs are not escaped —
// run user-supplied text through EscapeMarkdown first when needed.

// Bold wraps text in bold markup: "**text**".
func Bold(text string) string { return "**" + text + "**" }

// Italic wraps text in italic markup: "_text_".
func Italic(text string) string { return "_" + text + "_" }

// Underline wraps text in underline markup: "__text__".
func Underline(text string) string { return "__" + text + "__" }

// Strikethrough wraps text in strikethrough markup: "~~text~~".
func Strikethrough(text string) string { return "~~" + text + "~~" }

// Spoiler wraps text in a spoiler tag: "||text||".
func Spoiler(text string) string { return "||" + text + "||" }

// InlineCode wraps text in an inline code span: "`text`".
func InlineCode(text string) string { return "`" + text + "`" }

// CodeBlock wraps text in a fenced code block with no language hint.
func CodeBlock(text string) string { return "```\n" + text + "\n```" }

// CodeBlockLang wraps text in a fenced code block tagged with the given language
// (e.g. "go", "json") for syntax highlighting.
func CodeBlockLang(language, text string) string {
	return "```" + language + "\n" + text + "\n```"
}

// Quote prefixes every line of text with "> " so it renders as a block quote.
// Discord only quotes the first line of a single "> " prefix, so unlike
// discord.js' quote this prefixes each line to keep multi-line input quoted.
func Quote(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = "> " + line
	}
	return strings.Join(lines, "\n")
}

// BlockQuote prefixes text with ">>> ", which quotes all following content.
func BlockQuote(text string) string { return ">>> " + text }

// Hyperlink renders text as a masked link to url: "[text](url)".
func Hyperlink(text, url string) string { return "[" + text + "](" + url + ")" }

// HideLinkEmbed wraps a URL in angle brackets so Discord does not generate a
// preview embed for it: "<url>".
func HideLinkEmbed(url string) string { return "<" + url + ">" }

// Heading levels for the Heading helper. Discord supports H1–H3.
type HeadingLevel int

const (
	HeadingLevelOne   HeadingLevel = 1
	HeadingLevelTwo   HeadingLevel = 2
	HeadingLevelThree HeadingLevel = 3
)

// Heading prefixes text with the markdown heading marker for the given level.
// Levels outside 1–3 are clamped into range.
func Heading(text string, level HeadingLevel) string {
	if level < HeadingLevelOne {
		level = HeadingLevelOne
	}
	if level > HeadingLevelThree {
		level = HeadingLevelThree
	}
	return strings.Repeat("#", int(level)) + " " + text
}

// Subtext prefixes text with "-# ", rendering it as small subtext.
func Subtext(text string) string { return "-# " + text }

// UnorderedList renders items as a markdown bullet list ("- item" per line).
func UnorderedList(items []string) string {
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = "- " + item
	}
	return strings.Join(lines, "\n")
}

// OrderedList renders items as a numbered markdown list starting at 1
// ("1. item", "2. item", …).
func OrderedList(items []string) string {
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = strconv.Itoa(i+1) + ". " + item
	}
	return strings.Join(lines, "\n")
}

// ── Mention builders ──────────────────────────────────────────────────────────

// UserMention returns the mention markup for a user ID: "<@id>".
func UserMention(id Snowflake) string { return "<@" + id.String() + ">" }

// ChannelMention returns the mention markup for a channel ID: "<#id>".
func ChannelMention(id Snowflake) string { return "<#" + id.String() + ">" }

// RoleMention returns the mention markup for a role ID: "<@&id>".
func RoleMention(id Snowflake) string { return "<@&" + id.String() + ">" }

// SlashCommandMention returns the clickable mention for a registered slash
// command: "</name:id>". For subcommands, pass the full name ("group sub").
func SlashCommandMention(name string, id Snowflake) string {
	return "</" + name + ":" + id.String() + ">"
}

// FormatEmoji returns the markup used to render a custom emoji by ID:
// "<:_:id>", or "<a:_:id>" when animated. The name is irrelevant to rendering,
// so a placeholder underscore is used.
func FormatEmoji(id Snowflake, animated bool) string {
	if animated {
		return "<a:_:" + id.String() + ">"
	}
	return "<:_:" + id.String() + ">"
}

// ── Jump links ────────────────────────────────────────────────────────────────

// ChannelLink returns the client jump link to a channel. Pass a nil guildID for
// a DM/group-DM channel (rendered with the "@me" pseudo-guild).
func ChannelLink(channelID Snowflake, guildID *Snowflake) string {
	guild := "@me"
	if guildID != nil && !guildID.IsEmpty() {
		guild = guildID.String()
	}
	return "https://discord.com/channels/" + guild + "/" + channelID.String()
}

// MessageLink returns the client jump link to a specific message. Pass a nil
// guildID for a DM/group-DM message.
func MessageLink(channelID, messageID Snowflake, guildID *Snowflake) string {
	return ChannelLink(channelID, guildID) + "/" + messageID.String()
}

// ── Timestamps ────────────────────────────────────────────────────────────────

// TimestampStyle selects how Discord renders a timestamp mention.
type TimestampStyle string

const (
	TimestampStyleShortTime     TimestampStyle = "t" // 16:20
	TimestampStyleLongTime      TimestampStyle = "T" // 16:20:30
	TimestampStyleShortDate     TimestampStyle = "d" // 20/04/2021
	TimestampStyleLongDate      TimestampStyle = "D" // 20 April 2021
	TimestampStyleShortDateTime TimestampStyle = "f" // 20 April 2021 16:20 (default)
	TimestampStyleLongDateTime  TimestampStyle = "F" // Tuesday, 20 April 2021 16:20
	TimestampStyleRelative      TimestampStyle = "R" // 2 months ago
)

// Timestamp returns a dynamic timestamp mention that each viewer sees in their
// own locale and timezone, using Discord's default (short date/time) style:
// "<t:unix>".
func Timestamp(t time.Time) string {
	return "<t:" + strconv.FormatInt(t.Unix(), 10) + ">"
}

// TimestampStyled is like Timestamp but with an explicit display style:
// "<t:unix:style>".
func TimestampStyled(t time.Time, style TimestampStyle) string {
	return "<t:" + strconv.FormatInt(t.Unix(), 10) + ":" + string(style) + ">"
}
