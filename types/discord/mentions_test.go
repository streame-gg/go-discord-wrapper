package discord

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// mentionsSuite covers the formatting helpers: mention markup, emoji
// identifiers, classic tags, hex colours and snowflake-derived creation times.
type mentionsSuite struct {
	suite.Suite
}

func TestMentionsSuite(t *testing.T) {
	suite.Run(t, new(mentionsSuite))
}

func intptr(i int) *int    { return &i }
func boolptr(b bool) *bool { return &b }

// ── Mentions ────────────────────────────────────────────────────────────────

func (s *mentionsSuite) TestUserMention() {
	s.Equal("", (*User)(nil).Mention())
	s.Equal("<@123>", (&User{ID: 123}).Mention())
}

func (s *mentionsSuite) TestMemberMention() {
	s.Equal("", (&GuildMember{}).Mention())
	s.Equal("<@7>", (&GuildMember{UserID: 7}).Mention())
	s.Equal("<@9>", (&GuildMember{User: &User{ID: 9}}).Mention())
}

func (s *mentionsSuite) TestRoleMention() {
	s.Equal("", (*Role)(nil).Mention())
	s.Equal("<@&123>", (&Role{ID: 123}).Mention())
	// The default role (ID == GuildID) renders as @everyone.
	s.Equal("@everyone", (&Role{ID: 55, GuildID: 55}).Mention())
}

func (s *mentionsSuite) TestChannelMention() {
	s.Equal("", (*Channel)(nil).Mention())
	s.Equal("<#321>", (&Channel{ID: 321}).Mention())
}

func (s *mentionsSuite) TestEmojiMentionAndIdentifier() {
	unicode := &Emoji{Name: "🔥"}
	s.Equal("🔥", unicode.Mention())
	s.Equal("🔥", unicode.Identifier())

	custom := &Emoji{ID: 100, Name: "blob"}
	s.Equal("<:blob:100>", custom.Mention())
	s.Equal("blob:100", custom.Identifier())

	animated := &Emoji{ID: 200, Name: "spin", Animated: true}
	s.Equal("<a:spin:200>", animated.Mention())
	s.Equal("a:spin:200", animated.Identifier())

	s.Equal("", (*Emoji)(nil).Mention())
	s.Equal("", (*Emoji)(nil).Identifier())
}

// ── Tag ─────────────────────────────────────────────────────────────────────

func (s *mentionsSuite) TestUserTag() {
	s.Equal("", (*User)(nil).Tag())
	s.Equal("nova", (&User{Username: "nova", Discriminator: "0"}).Tag())
	s.Equal("nova", (&User{Username: "nova"}).Tag())
	s.Equal("nova#1234", (&User{Username: "nova", Discriminator: "1234"}).Tag())
}

// ── Colours ─────────────────────────────────────────────────────────────────

func (s *mentionsSuite) TestRoleHexColor() {
	s.Equal("", (*Role)(nil).HexColor())
	s.Equal("#000000", (&Role{}).HexColor())
	s.Equal("#5865f2", (&Role{Colors: RoleColors{PrimaryColor: intptr(0x5865F2)}}).HexColor())
}

func (s *mentionsSuite) TestUserHexAccentColor() {
	s.Equal("", (*User)(nil).HexAccentColor())
	s.Equal("", (&User{}).HexAccentColor())
	s.Equal("#ff00aa", (&User{AccentColor: intptr(0xFF00AA)}).HexAccentColor())
}

// ── Creation timestamps ─────────────────────────────────────────────────────

func (s *mentionsSuite) TestCreatedAt() {
	// Documented Discord example snowflake → 2016-04-30T11:18:25.796Z.
	const id Snowflake = 175928847299117063
	want := time.Date(2016, 4, 30, 11, 18, 25, 796_000_000, time.UTC)

	s.True((&Guild{ID: id}).CreatedAt().Equal(want))
	s.True((&Channel{ID: id}).CreatedAt().Equal(want))
	s.True((&Role{ID: id}).CreatedAt().Equal(want))
	s.True((&Message{ID: id}).CreatedAt().Equal(want))
	s.True((&Emoji{ID: id}).CreatedAt().Equal(want))
	s.True((&GuildMember{UserID: id}).CreatedAt().Equal(want))
}

func (s *mentionsSuite) TestCreatedAtZeroCases() {
	// Unicode emoji and identity-less members have no snowflake to decode.
	s.True((&Emoji{Name: "🔥"}).CreatedAt().IsZero())
	s.True((&GuildMember{}).CreatedAt().IsZero())
}

// ── Channel type predicates ─────────────────────────────────────────────────

func (s *mentionsSuite) TestChannelPredicates() {
	cases := []struct {
		t                       ChannelType
		text, voice, dm, thread bool
	}{
		{ChannelTypeGuildText, true, false, false, false},
		{ChannelTypeGuildVoice, true, true, false, false},
		{ChannelTypeGuildStageVoice, true, true, false, false},
		{ChannelTypeDM, true, false, true, false},
		{ChannelTypeGroupDM, true, false, true, false},
		{ChannelTypePublicThread, true, false, false, true},
		{ChannelTypeGuildCategory, false, false, false, false},
		{ChannelTypeGuildForum, false, false, false, false},
	}
	for _, c := range cases {
		ch := &Channel{Type: c.t}
		s.Equalf(c.text, ch.IsTextBased(), "IsTextBased type %d", c.t)
		s.Equalf(c.voice, ch.IsVoiceBased(), "IsVoiceBased type %d", c.t)
		s.Equalf(c.dm, ch.IsDMBased(), "IsDMBased type %d", c.t)
		s.Equalf(c.thread, ch.IsThread(), "IsThread type %d", c.t)
	}
}

// ── Jump links ────────────────────────────────────────────────────────────────

func (s *mentionsSuite) TestChannelURL() {
	s.Equal("", (*Channel)(nil).URL())
	gid := Snowflake(10)
	s.Equal("https://discord.com/channels/10/20", (&Channel{ID: 20, GuildID: &gid}).URL())
	// DM channel has no guild → "@me".
	s.Equal("https://discord.com/channels/@me/20", (&Channel{ID: 20}).URL())
}

func (s *mentionsSuite) TestInviteURL() {
	s.Equal("", (*Invite)(nil).URL())
	s.Equal("", (&Invite{}).URL())
	s.Equal("https://discord.gg/abc123", (&Invite{Code: "abc123"}).URL())
}

// ── Predicates ──────────────────────────────────────────────────────────────

func (s *mentionsSuite) TestChannelNSFWAndThreadOnly() {
	s.False((&Channel{}).IsNSFW())
	s.True((&Channel{NSFW: boolptr(true)}).IsNSFW())
	s.True((&Channel{Type: ChannelTypeGuildForum}).IsThreadOnly())
	s.True((&Channel{Type: ChannelTypeGuildMedia}).IsThreadOnly())
	s.False((&Channel{Type: ChannelTypeGuildText}).IsThreadOnly())
}

func (s *mentionsSuite) TestMessageSystemAndEdited() {
	s.False((&Message{Type: MessageTypeDefault}).IsSystem())
	s.False((&Message{Type: MessageTypeReply}).IsSystem())
	s.True((&Message{Type: MessageTypeChannelPinnedMessage}).IsSystem())

	s.False((&Message{}).IsEdited())
	now := time.Now()
	s.True((&Message{EditedTimestamp: &now}).IsEdited())
}

func (s *mentionsSuite) TestMemberIsCommunicationDisabled() {
	s.False((&GuildMember{}).IsCommunicationDisabled())
	past := time.Now().Add(-time.Hour)
	s.False((&GuildMember{CommunicationDisabledUntil: &past}).IsCommunicationDisabled())
	future := time.Now().Add(time.Hour)
	s.True((&GuildMember{CommunicationDisabledUntil: &future}).IsCommunicationDisabled())
}

func (s *mentionsSuite) TestRoleComparePositionTo() {
	high := &Role{ID: 5, Position: 10}
	low := &Role{ID: 6, Position: 2}
	s.Positive(high.ComparePositionTo(low))
	s.Negative(low.ComparePositionTo(high))

	// Same position → smaller ID ranks higher.
	older := &Role{ID: 100, Position: 3}
	newer := &Role{ID: 200, Position: 3}
	s.Positive(older.ComparePositionTo(newer))
	s.Negative(newer.ComparePositionTo(older))
	s.Zero(older.ComparePositionTo(older))
	s.Zero(older.ComparePositionTo(nil))
}

func (s *mentionsSuite) TestMemberDisplayAvatarURL() {
	gid := Snowflake(1)
	// Per-guild avatar wins.
	m := &GuildMember{GuildID: gid, UserID: 2, AvatarHash: strptr("memberhash")}
	s.Contains(m.DisplayAvatarURL(nil), "/guilds/1/users/2/avatars/memberhash")
	// Falls back to the user's avatar when no guild avatar.
	m2 := &GuildMember{GuildID: gid, UserID: 2, User: &User{ID: 2, AvatarHash: strptr("userhash")}}
	s.Contains(m2.DisplayAvatarURL(nil), "/avatars/2/userhash")
	// No user to fall back to → empty.
	s.Equal("", (&GuildMember{}).DisplayAvatarURL(nil))
}
