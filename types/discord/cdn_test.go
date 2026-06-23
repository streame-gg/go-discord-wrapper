package discord

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/stretchr/testify/suite"
)

func strptr(s string) *string { return &s }

func (su *cdnSuite) TestEmojiURL_AnimatedIsNil() {
	em := Emoji{
		ID:   123,
		Name: "1231231231",
	}

	url := em.URL(nil)
	su.Require().NotNil(url)
}

func (su *cdnSuite) TestUserAvatarURL() {
	t := su.T()
	u := &User{ID: 123, Avatar: strptr("abc")}
	if got := u.AvatarURL(nil); got != "https://cdn.discordapp.com/avatars/123/abc.webp" {
		t.Errorf("default: got %q", got)
	}
	if got := u.AvatarURL(&ImageOptions{Format: ImageFormatPNG, Size: 256}); got != "https://cdn.discordapp.com/avatars/123/abc.png?size=256" {
		t.Errorf("png+size: got %q", got)
	}

	animated := &User{ID: 1, Avatar: strptr("a_xyz")}
	if got := animated.AvatarURL(nil); got != "https://cdn.discordapp.com/avatars/1/a_xyz.gif" {
		t.Errorf("animated should default to gif: got %q", got)
	}

	none := &User{ID: 1}
	if got := none.AvatarURL(nil); got != "" {
		t.Errorf("no avatar should be empty: got %q", got)
	}
}

func (su *cdnSuite) TestUserDefaultAndDisplayAvatarURL() {
	t := su.T()
	// Migrated username (discriminator "0") → index from ID>>22 % 6.
	migrated := &User{ID: 1 << 22, Discriminator: "0"} // (ID>>22)=1, %6=1
	if got := migrated.DefaultAvatarURL(); got != "https://cdn.discordapp.com/embed/avatars/1.png" {
		t.Errorf("migrated default: got %q", got)
	}

	// Legacy account → discriminator % 5.
	legacy := &User{ID: 99, Discriminator: "0007"} // 7%5=2
	if got := legacy.DefaultAvatarURL(); got != "https://cdn.discordapp.com/embed/avatars/2.png" {
		t.Errorf("legacy default: got %q", got)
	}

	// DisplayAvatarURL falls back to the default when no custom avatar.
	if got := legacy.DisplayAvatarURL(nil); got != legacy.DefaultAvatarURL() {
		t.Errorf("display should fall back to default: got %q", got)
	}
	withAvatar := &User{ID: 5, Avatar: strptr("h")}
	if got := withAvatar.DisplayAvatarURL(nil); got != withAvatar.AvatarURL(nil) {
		t.Errorf("display should use custom avatar: got %q", got)
	}
}

func (su *cdnSuite) TestMemberAvatarAndBannerURL() {
	t := su.T()
	m := &GuildMember{GuildID: 10, UserID: 20, AvatarHash: strptr("a_anim"), BannerHash: strptr("ban")}
	if got := m.AvatarURL(nil); got != "https://cdn.discordapp.com/guilds/10/users/20/avatars/a_anim.gif" {
		t.Errorf("member avatar: got %q", got)
	}
	if got := m.BannerURL(&ImageOptions{Size: 512}); got != "https://cdn.discordapp.com/guilds/10/users/20/banners/ban.webp?size=512" {
		t.Errorf("member banner: got %q", got)
	}

	noGuild := &GuildMember{UserID: 20, AvatarHash: strptr("x")}
	if got := noGuild.AvatarURL(nil); got != "" {
		t.Errorf("missing guild id should yield empty: got %q", got)
	}
}

func (su *cdnSuite) TestGuildImageURLs() {
	t := su.T()
	g := &Guild{ID: 7, IconHash: strptr("ic"), Splash: strptr("sp"), DiscoverySplash: strptr("ds")}
	if got := g.IconURL(nil); got != "https://cdn.discordapp.com/icons/7/ic.webp" {
		t.Errorf("icon: got %q", got)
	}
	if got := g.SplashURL(nil); got != "https://cdn.discordapp.com/splashes/7/sp.webp" {
		t.Errorf("splash: got %q", got)
	}
	if got := g.DiscoverySplashURL(nil); got != "https://cdn.discordapp.com/discovery-splashes/7/ds.webp" {
		t.Errorf("discovery splash: got %q", got)
	}
	if got := (&Guild{ID: 7}).IconURL(nil); got != "" {
		t.Errorf("no icon should be empty: got %q", got)
	}
}

func (su *cdnSuite) TestRoleAndEmojiURL() {
	t := su.T()
	r := &Role{ID: 3, Icon: strptr("ri")}
	if got := r.IconURL(nil); got != "https://cdn.discordapp.com/role-icons/3/ri.webp" {
		t.Errorf("role icon: got %q", got)
	}

	staticEmoji := &Emoji{ID: 42}
	if got := staticEmoji.URL(nil); got != "https://cdn.discordapp.com/emojis/42.webp" {
		t.Errorf("static emoji: got %q", got)
	}
	animEmoji := &Emoji{ID: 42, Animated: util.PointerOf(true)}
	if got := animEmoji.URL(nil); got != "https://cdn.discordapp.com/emojis/42.gif" {
		t.Errorf("animated emoji: got %q", got)
	}
	unicodeEmoji := &Emoji{Name: "🙂"}
	if got := unicodeEmoji.URL(nil); got != "" {
		t.Errorf("unicode emoji should be empty: got %q", got)
	}
}

func (su *cdnSuite) TestStickerURL() {
	t := su.T()
	cases := []struct {
		format StickerFormatType
		want   string
	}{
		{StickerFormatTypePNG, "https://cdn.discordapp.com/stickers/8.png"},
		{StickerFormatTypeAPNG, "https://cdn.discordapp.com/stickers/8.png"},
		{StickerFormatTypeLottie, "https://cdn.discordapp.com/stickers/8.json"},
		{StickerFormatTypeGIF, "https://media.discordapp.net/stickers/8.gif"},
	}
	for _, c := range cases {
		s := &Sticker{ID: 8, FormatType: c.format}
		if got := s.URL(); got != c.want {
			t.Errorf("format %d: got %q, want %q", c.format, got, c.want)
		}
	}
}

func (su *cdnSuite) TestApplicationImageURLs() {
	t := su.T()
	a := &Application{ID: 9, Icon: strptr("ai"), CoverImage: "cv"}
	if got := a.IconURL(nil); got != "https://cdn.discordapp.com/app-icons/9/ai.webp" {
		t.Errorf("app icon: got %q", got)
	}
	if got := a.CoverImageURL(nil); got != "https://cdn.discordapp.com/app-icons/9/cv.webp" {
		t.Errorf("app cover: got %q", got)
	}
}

type cdnSuite struct{ suite.Suite }

func TestCdnSuite(t *testing.T) { suite.Run(t, new(cdnSuite)) }
