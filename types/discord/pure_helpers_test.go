package discord

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// pureHelpersSuite groups tests for the small, dependency-free helpers on the
// discord value types (display names, thread classification, gateway version
// formatting and the avatar-decoration CDN URL). These are pure functions, so
// the suite needs no shared setup beyond the fixtures built inline per test.
type pureHelpersSuite struct {
	suite.Suite
}

func TestPureHelpersSuite(t *testing.T) {
	suite.Run(t, new(pureHelpersSuite))
}

// ── User.DisplayName ──────────────────────────────────────────────────────────

func (s *pureHelpersSuite) TestUserDisplayName() {
	cases := []struct {
		name string
		user User
		want string
	}{
		{"global name preferred", User{GlobalName: strptr("Nova"), Username: "nova_legacy"}, "Nova"},
		{"empty global name falls back to username", User{GlobalName: strptr(""), Username: "fallback"}, "fallback"},
		{"nil global name falls back to username", User{Username: "plain"}, "plain"},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			u := c.user
			s.Equal(c.want, u.DisplayName())
		})
	}
}

// ── GuildMember.DisplayName ─────────────────────────────────────────────────────

func (s *pureHelpersSuite) TestGuildMemberDisplayName() {
	s.Run("nil receiver is empty", func() {
		var m *GuildMember
		s.Equal("", m.DisplayName())
	})

	s.Run("nick takes precedence over user", func() {
		m := &GuildMember{Nick: strptr("Cmdr"), User: &User{Username: "ignored"}}
		s.Equal("Cmdr", m.DisplayName())
	})

	s.Run("falls back to user display name when no nick", func() {
		m := &GuildMember{User: &User{GlobalName: strptr("Global"), Username: "uname"}}
		s.Equal("Global", m.DisplayName())
	})

	s.Run("empty when neither nick nor user is set", func() {
		m := &GuildMember{}
		s.Equal("", m.DisplayName())
	})
}

// ── User.AvatarDecorationURL ────────────────────────────────────────────────────

func (s *pureHelpersSuite) TestUserAvatarDecorationURL() {
	s.Run("nil receiver is empty", func() {
		var u *User
		s.Equal("", u.AvatarDecorationURL())
	})

	s.Run("no decoration data is empty", func() {
		u := &User{ID: 1}
		s.Equal("", u.AvatarDecorationURL())
	})

	s.Run("empty asset is empty", func() {
		u := &User{ID: 1, AvatarDecorationData: &AvatarDecorationData{Asset: ""}}
		s.Equal("", u.AvatarDecorationURL())
	})

	s.Run("valid asset builds preset URL", func() {
		u := &User{ID: 1, AvatarDecorationData: &AvatarDecorationData{Asset: "deco_hash"}}
		s.Equal("https://cdn.discordapp.com/avatar-decoration-presets/deco_hash.png", u.AvatarDecorationURL())
	})
}

// ── Channel.IsThread ────────────────────────────────────────────────────────────

func (s *pureHelpersSuite) TestChannelIsThread() {
	threadTypes := []ChannelType{
		ChannelTypeAnnouncementThread,
		ChannelTypePublicThread,
		ChannelTypePrivateThread,
	}
	for _, t := range threadTypes {
		s.Truef((&Channel{Type: t}).IsThread(), "type %d should be a thread", t)
	}

	nonThreadTypes := []ChannelType{
		ChannelTypeGuildText,
		ChannelTypeDM,
		ChannelTypeGuildVoice,
		ChannelTypeGuildCategory,
		ChannelTypeGuildAnnouncement,
		ChannelTypeGuildForum,
	}
	for _, t := range nonThreadTypes {
		s.Falsef((&Channel{Type: t}).IsThread(), "type %d should not be a thread", t)
	}
}

// ── APIVersion / gateway formatting ─────────────────────────────────────────────

func (s *pureHelpersSuite) TestAPIVersionToString() {
	s.Equal("10", APIVersion10.ToString())
	s.Equal("9", APIVersion9.ToString())
	s.Equal("unknown", APIVersion(255).ToString())
}

func (s *pureHelpersSuite) TestAPIBaseString() {
	s.Equal("/api/v10/", APIBaseString(APIVersion10))
	s.Equal("/api/v9/", APIBaseString(APIVersion9))
	s.Equal("/api/vunknown/", APIBaseString(APIVersion(0)))
}

func (s *pureHelpersSuite) TestGatewayErrorError() {
	err := GatewayError{Code: 50035, Message: "Invalid Form Body"}
	s.Equal("50035 Invalid Form Body", err.Error())
	// Implements the error interface.
	var _ error = err
}

// ── Application CDN nil-guards ──────────────────────────────────────────────────

// These complete the branch coverage for the application image helpers: the
// happy path is exercised in cdn_test.go, here we assert the no-asset guards.
func (s *pureHelpersSuite) TestApplicationImageURLNilGuards() {
	s.Run("nil icon hash", func() {
		s.Equal("", (&Application{ID: 9}).IconURL(nil))
	})
	s.Run("nil cover image", func() {
		s.Equal("", (&Application{ID: 9}).CoverImageURL(nil))
	})
	s.Run("nil receiver", func() {
		var a *Application
		s.Equal("", a.IconURL(nil))
		s.Equal("", a.CoverImageURL(nil))
	})
}
