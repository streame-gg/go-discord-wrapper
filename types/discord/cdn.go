package discord

import (
	"strconv"
	"strings"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
)

// CDN base hosts. Most assets are served from cdnBaseURL; animated stickers
// (GIF) are only available from cdnMediaURL.
const (
	cdnBaseURL  = "https://cdn.discordapp.com"
	cdnMediaURL = "https://media.discordapp.net"
)

// ImageFormat is a CDN image file extension.
//
// https://docs.discord.com/developers/reference#image-formatting
type ImageFormat string

const (
	ImageFormatPNG  ImageFormat = "png"
	ImageFormatJPEG ImageFormat = "jpg"
	ImageFormatWebP ImageFormat = "webp"
	ImageFormatGIF  ImageFormat = "gif"
)

// ImageOptions tweaks a generated CDN URL. A nil *ImageOptions selects the
// defaults: WebP (or GIF for animated assets) at Discord's default size.
// https://docs.discord.com/developers/reference#image-formatting
type ImageOptions struct {
	// Format overrides the file extension. When empty the helper uses GIF for
	// animated assets and WebP otherwise.
	Format ImageFormat
	// Size requests a specific square size via the ?size query parameter. Discord
	// only honors powers of two in the range [16, 4096]; zero omits the
	// parameter and returns the asset's default size.
	Size int
}

// imageURL assembles a CDN URL for the asset at path (which must already include
// any hash/ID and a leading slash). animated reports whether the underlying
// asset is animated, which drives the default format and is detected by callers
// (e.g. from an "a_" hash prefix or an Animated flag).
func imageURL(base, path string, animated bool, opts *ImageOptions) string {
	format := ImageFormatWebP
	if animated {
		format = ImageFormatGIF
	}
	if opts != nil && opts.Format != "" {
		format = opts.Format
	}

	var b strings.Builder
	b.WriteString(base)
	b.WriteString(path)
	b.WriteByte('.')
	b.WriteString(string(format))
	if opts != nil && opts.Size > 0 {
		b.WriteString("?size=")
		b.WriteString(strconv.Itoa(opts.Size))
	}
	return b.String()
}

// hashIsAnimated reports whether a Discord asset hash denotes an animated image.
func hashIsAnimated(hash string) bool { return strings.HasPrefix(hash, "a_") }

// ── User ────────────────────────────────────────────────────────────────────

// AvatarURL returns the user's custom avatar URL, or "" if the user has no
// custom avatar (use DisplayAvatarURL to fall back to the default avatar).
func (u *User) AvatarURL(opts *ImageOptions) string {
	if u == nil || u.AvatarHash == nil {
		return ""
	}
	h := *u.AvatarHash
	return imageURL(cdnBaseURL, "/avatars/"+u.ID.String()+"/"+h, hashIsAnimated(h), opts)
}

// DefaultAvatarURL returns the URL of the user's default (Discord-provided)
// avatar. For migrated usernames (discriminator "0") the index is derived from
// the user ID; for legacy accounts it is the discriminator modulo 5.
func (u *User) DefaultAvatarURL() string {
	if u == nil {
		return ""
	}
	var index uint64
	if u.Discriminator == "" || u.Discriminator == "0" {
		index = (uint64(u.ID) >> 22) % 6
	} else {
		d, _ := strconv.Atoi(u.Discriminator)
		index = uint64(d % 5)
	}
	return cdnBaseURL + "/embed/avatars/" + strconv.FormatUint(index, 10) + ".png"
}

// DisplayAvatarURL returns the user's custom avatar when set, otherwise the
// default avatar.
func (u *User) DisplayAvatarURL(opts *ImageOptions) string {
	if url := u.AvatarURL(opts); url != "" {
		return url
	}
	return u.DefaultAvatarURL()
}

// AvatarDecorationURL returns the URL of the user's avatar decoration preset, or
// "" if they have none.
func (u *User) AvatarDecorationURL() string {
	if u == nil || u.AvatarDecorationData == nil || u.AvatarDecorationData.Asset == "" {
		return ""
	}
	return cdnBaseURL + "/avatar-decoration-presets/" + u.AvatarDecorationData.Asset + ".png"
}

// ── GuildMember ───────────────────────────────────────────────────────────────

// AvatarURL returns the member's guild-specific avatar URL, or "" if the member
// has no per-guild avatar (fall back to the user's avatar in that case).
func (m *GuildMember) AvatarURL(opts *ImageOptions) string {
	if m == nil || m.AvatarHash == nil {
		return ""
	}
	uid := memberUserID(m)
	if m.GuildID.IsEmpty() || uid.IsEmpty() {
		return ""
	}
	h := *m.AvatarHash
	path := "/guilds/" + m.GuildID.String() + "/users/" + uid.String() + "/avatars/" + h
	return imageURL(cdnBaseURL, path, hashIsAnimated(h), opts)
}

// DisplayAvatarURL returns the avatar Discord shows for this member: their
// per-guild avatar when set, otherwise the underlying user's avatar (custom or
// default). It returns "" only when no user is attached to fall back to.
func (m *GuildMember) DisplayAvatarURL(opts *ImageOptions) string {
	if url := m.AvatarURL(opts); url != "" {
		return url
	}
	if m != nil && m.User != nil {
		return m.User.DisplayAvatarURL(opts)
	}
	return ""
}

// BannerURL returns the member's guild-specific banner URL, or "" if unset.
func (m *GuildMember) BannerURL(opts *ImageOptions) string {
	if m == nil || m.BannerHash == nil {
		return ""
	}
	uid := memberUserID(m)
	if m.GuildID.IsEmpty() || uid.IsEmpty() {
		return ""
	}
	h := *m.BannerHash
	path := "/guilds/" + m.GuildID.String() + "/users/" + uid.String() + "/banners/" + h
	return imageURL(cdnBaseURL, path, hashIsAnimated(h), opts)
}

// ── Guild ─────────────────────────────────────────────────────────────────────

// IconURL returns the guild's icon URL, or "" if it has none.
func (g *Guild) IconURL(opts *ImageOptions) string {
	if g == nil || g.IconHash == nil {
		return ""
	}
	h := *g.IconHash
	return imageURL(cdnBaseURL, "/icons/"+g.ID.String()+"/"+h, hashIsAnimated(h), opts)
}

// SplashURL returns the guild's invite splash URL, or "" if it has none.
func (g *Guild) SplashURL(opts *ImageOptions) string {
	if g == nil || g.Splash == nil {
		return ""
	}
	return imageURL(cdnBaseURL, "/splashes/"+g.ID.String()+"/"+*g.Splash, false, opts)
}

// DiscoverySplashURL returns the guild's discovery splash URL, or "" if unset.
func (g *Guild) DiscoverySplashURL(opts *ImageOptions) string {
	if g == nil || g.DiscoverySplash == nil {
		return ""
	}
	return imageURL(cdnBaseURL, "/discovery-splashes/"+g.ID.String()+"/"+*g.DiscoverySplash, false, opts)
}

// ── Role ──────────────────────────────────────────────────────────────────────

// IconURL returns the role's icon URL, or "" if it has none.
func (r *Role) IconURL(opts *ImageOptions) string {
	if r == nil || r.Icon == nil {
		return ""
	}
	h := *r.Icon
	return imageURL(cdnBaseURL, "/role-icons/"+r.ID.String()+"/"+h, hashIsAnimated(h), opts)
}

// ── Emoji ─────────────────────────────────────────────────────────────────────

// URL returns the custom emoji's image URL. It returns "" for unicode/standard
// emoji (which have no Discord ID). Animated emoji default to GIF.
func (e *Emoji) URL(opts *ImageOptions) string {
	if e == nil || e.ID.IsEmpty() {
		return ""
	}
	if e.Animated == nil {
		e.Animated = util.PointerOf(false)
	}
	return imageURL(cdnBaseURL, "/emojis/"+e.ID.String(), *e.Animated, opts)
}

// ── Sticker ───────────────────────────────────────────────────────────────────

// URL returns the sticker's asset URL. PNG/APNG stickers resolve to a PNG,
// Lottie stickers to the JSON animation file, and GIF stickers to the GIF served
// from the media host. It returns "" for an unknown format.
func (s *Sticker) URL() string {
	if s == nil || s.ID.IsEmpty() {
		return ""
	}
	switch s.FormatType {
	case StickerFormatTypePNG, StickerFormatTypeAPNG:
		return cdnBaseURL + "/stickers/" + s.ID.String() + ".png"
	case StickerFormatTypeLottie:
		return cdnBaseURL + "/stickers/" + s.ID.String() + ".json"
	case StickerFormatTypeGIF:
		return cdnMediaURL + "/stickers/" + s.ID.String() + ".gif"
	default:
		return ""
	}
}

// ── Application ───────────────────────────────────────────────────────────────

// IconURL returns the application's icon URL, or "" if it has none.
func (a *Application) IconURL(opts *ImageOptions) string {
	if a == nil || a.Icon == nil {
		return ""
	}
	return imageURL(cdnBaseURL, "/app-icons/"+a.ID.String()+"/"+*a.Icon, false, opts)
}

// CoverImageURL returns the application's store cover image URL, or "" if unset.
func (a *Application) CoverImageURL(opts *ImageOptions) string {
	if a == nil || a.CoverImage == "" {
		return ""
	}
	return imageURL(cdnBaseURL, "/app-icons/"+a.ID.String()+"/"+a.CoverImage, false, opts)
}
