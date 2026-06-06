package discord

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// formattersSuite covers the markdown / mention / timestamp builders.
type formattersSuite struct {
	suite.Suite
}

func TestFormattersSuite(t *testing.T) {
	suite.Run(t, new(formattersSuite))
}

func (s *formattersSuite) TestTextStyles() {
	s.Equal("**hi**", Bold("hi"))
	s.Equal("_hi_", Italic("hi"))
	s.Equal("__hi__", Underline("hi"))
	s.Equal("~~hi~~", Strikethrough("hi"))
	s.Equal("||hi||", Spoiler("hi"))
	s.Equal("`hi`", InlineCode("hi"))
	s.Equal("-# hi", Subtext("hi"))
}

func (s *formattersSuite) TestCodeBlocks() {
	s.Equal("```\nhi\n```", CodeBlock("hi"))
	s.Equal("```go\nx := 1\n```", CodeBlockLang("go", "x := 1"))
}

func (s *formattersSuite) TestQuotes() {
	s.Equal("> a\n> b", Quote("a\nb"))
	s.Equal(">>> a\nb", BlockQuote("a\nb"))
}

func (s *formattersSuite) TestLinks() {
	s.Equal("[Docs](https://discord.dev)", Hyperlink("Docs", "https://discord.dev"))
	s.Equal("<https://discord.dev>", HideLinkEmbed("https://discord.dev"))
}

func (s *formattersSuite) TestHeading() {
	s.Equal("# Title", Heading("Title", HeadingLevelOne))
	s.Equal("## Title", Heading("Title", HeadingLevelTwo))
	s.Equal("### Title", Heading("Title", HeadingLevelThree))
	// Out-of-range levels clamp into 1–3.
	s.Equal("# Title", Heading("Title", 0))
	s.Equal("### Title", Heading("Title", 9))
}

func (s *formattersSuite) TestMentions() {
	s.Equal("<@1>", UserMention(1))
	s.Equal("<#2>", ChannelMention(2))
	s.Equal("<@&3>", RoleMention(3))
	s.Equal("</ban user:4>", SlashCommandMention("ban user", 4))
}

func (s *formattersSuite) TestTimestamps() {
	t := time.Unix(1618000000, 0)
	s.Equal("<t:1618000000>", Timestamp(t))
	s.Equal("<t:1618000000:R>", TimestampStyled(t, TimestampStyleRelative))
	s.Equal("<t:1618000000:F>", TimestampStyled(t, TimestampStyleLongDateTime))
}

func (s *formattersSuite) TestFormatEmoji() {
	s.Equal("<:_:123>", FormatEmoji(123, false))
	s.Equal("<a:_:123>", FormatEmoji(123, true))
}

func (s *formattersSuite) TestLinks2() {
	gid := Snowflake(1)
	s.Equal("https://discord.com/channels/1/2", ChannelLink(2, &gid))
	s.Equal("https://discord.com/channels/@me/2", ChannelLink(2, nil))
	s.Equal("https://discord.com/channels/1/2/3", MessageLink(2, 3, &gid))
	s.Equal("https://discord.com/channels/@me/2/3", MessageLink(2, 3, nil))
}

func (s *formattersSuite) TestListBuilders() {
	s.Equal("- a\n- b", UnorderedList([]string{"a", "b"}))
	s.Equal("1. a\n2. b\n3. c", OrderedList([]string{"a", "b", "c"}))
	s.Equal("", UnorderedList(nil))
	s.Equal("", OrderedList(nil))
}
