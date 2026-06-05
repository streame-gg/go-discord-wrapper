package discord

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// flagsSuite covers the generic FlagBits helper and the per-entity flag
// accessors.
type flagsSuite struct {
	suite.Suite
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(flagsSuite))
}

func (s *flagsSuite) TestHasAnyAll() {
	f := NewFlagBits(MessageFlagEphemeral | MessageFlagLoading)
	s.True(f.Has(MessageFlagEphemeral))
	s.True(f.Has(MessageFlagEphemeral | MessageFlagLoading)) // all combined bits present
	s.False(f.Has(MessageFlagUrgent))
	s.True(f.HasAll(MessageFlagEphemeral, MessageFlagLoading))
	s.False(f.HasAll(MessageFlagEphemeral, MessageFlagUrgent))
	s.True(f.Any(MessageFlagUrgent, MessageFlagLoading))
	s.False(f.Any(MessageFlagUrgent, MessageFlagCrossposted))
}

func (s *flagsSuite) TestEmptyAndValue() {
	var empty FlagBits[MessageFlag]
	s.True(empty.IsEmpty())
	s.False(empty.Has(MessageFlagEphemeral))

	f := NewFlagBits(MessageFlagEphemeral)
	s.False(f.IsEmpty())
	s.Equal(MessageFlagEphemeral, f.Value())
}

func (s *flagsSuite) TestEqualsAndMissing() {
	f := NewFlagBits(MessageFlagEphemeral | MessageFlagLoading)
	s.True(f.Equals(MessageFlagEphemeral, MessageFlagLoading))
	s.False(f.Equals(MessageFlagEphemeral))
	s.Equal([]MessageFlag{MessageFlagUrgent}, f.Missing(MessageFlagEphemeral, MessageFlagUrgent))
	s.Nil(f.Missing(MessageFlagEphemeral, MessageFlagLoading))
}

func (s *flagsSuite) TestAddRemoveToggleImmutable() {
	base := NewFlagBits(MessageFlagEphemeral)

	added := base.Add(MessageFlagLoading)
	s.True(added.Has(MessageFlagLoading))
	s.False(base.Has(MessageFlagLoading)) // original untouched

	removed := added.Remove(MessageFlagEphemeral)
	s.False(removed.Has(MessageFlagEphemeral))
	s.True(removed.Has(MessageFlagLoading))

	toggled := base.Toggle(MessageFlagEphemeral, MessageFlagUrgent)
	s.False(toggled.Has(MessageFlagEphemeral)) // was set, now cleared
	s.True(toggled.Has(MessageFlagUrgent))     // was clear, now set
}

func (s *flagsSuite) TestDecompose() {
	f := NewFlagBits(MessageFlagEphemeral | MessageFlagUrgent)
	got := f.Decompose()
	s.ElementsMatch([]MessageFlag{MessageFlagUrgent, MessageFlagEphemeral}, got)
	s.Nil(NewFlagBits(MessageFlag(0)).Decompose())
}

func (s *flagsSuite) TestEntityAccessors() {
	// User public flags.
	u := &User{PublicFlags: UserFlagActiveDeveloper | UserFlagVerifiedBot}
	s.True(u.PublicFlagBits().Has(UserFlagActiveDeveloper))
	s.True(u.PublicFlagBits().Has(UserFlagVerifiedBot))
	s.False(u.PublicFlagBits().Has(UserFlagStaff))
	s.True((*User)(nil).PublicFlagBits().IsEmpty())
	s.True((&User{Flags: UserFlagStaff}).FlagBits().Has(UserFlagStaff))

	// Guild member flags.
	gm := &GuildMember{Flags: GuildMemberFlagCompletedOnboarding}
	s.True(gm.FlagBits().Has(GuildMemberFlagCompletedOnboarding))
	s.False(gm.FlagBits().Has(GuildMemberFlagIsGuest))
	s.True((*GuildMember)(nil).FlagBits().IsEmpty())

	// Message flags (value field).
	m := &Message{Flags: MessageFlagEphemeral}
	s.True(m.FlagBits().Has(MessageFlagEphemeral))

	// Channel flags (pointer field, may be nil).
	s.True((&Channel{}).FlagBits().IsEmpty())
	cf := ChannelFlagPinned
	s.True((&Channel{Flags: &cf}).FlagBits().Has(ChannelFlagPinned))

	// Role flags (pointer field).
	s.True((&Role{}).FlagBits().IsEmpty())
	rf := RoleFlagsInPrompt
	s.True((&Role{Flags: &rf}).FlagBits().Has(RoleFlagsInPrompt))

	// Guild system-channel flags (value field).
	g := &Guild{SystemChannelFlags: GuildSystemChannelFlagsSuppressJoinNotifications}
	s.True(g.SystemChannelFlagBits().Has(GuildSystemChannelFlagsSuppressJoinNotifications))

	// Application flags (pointer field).
	s.True((&Application{}).FlagBits().IsEmpty())
	af := ApplicationFlagGatewayMessageContent
	s.True((&Application{Flags: &af}).FlagBits().Has(ApplicationFlagGatewayMessageContent))
}
