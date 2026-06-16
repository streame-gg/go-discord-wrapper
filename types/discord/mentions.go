package discord

import (
	"fmt"
	"time"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
)

// This file collects the small, dependency-free convenience helpers on the
// value types: Discord mention strings, raw ID formatting, hex colours and
// snowflake-derived creation timestamps. They mirror the ergonomics of
// discord.js' toString()/createdAt/hexColor getters so handlers can format
// entities without hand-assembling the `<@…>` / `<#…>` markup every time.

// ── Mentions ────────────────────────────────────────────────────────────────

// Mention returns the markup that renders as a clickable mention of this user,
// e.g. "<@123>". It returns "" for a nil receiver.
func (u *User) Mention() string {
	if u == nil {
		return ""
	}
	return "<@" + u.ID.String() + ">"
}

// Mention returns the mention markup for this member. It prefers the member's
// user ID and falls back to the resolved UserID, returning "" when neither is
// available.
func (m *GuildMember) Mention() string {
	uid := memberUserID(m)
	if uid.IsEmpty() {
		return ""
	}
	return "<@" + uid.String() + ">"
}

// Mention returns the markup that renders as a clickable mention of this role,
// e.g. "<@&123>". The guild's default role (where the role ID equals the guild
// ID) mentions everyone and renders as the literal "@everyone".
func (r *Role) Mention() string {
	if r == nil {
		return ""
	}
	if !r.GuildID.IsEmpty() && r.ID == r.GuildID {
		return "@everyone"
	}
	return "<@&" + r.ID.String() + ">"
}

// Mention returns the markup that renders as a clickable mention of this
// channel, e.g. "<#123>". It returns "" for a nil receiver.
func (c *Channel) Mention() string {
	if c == nil {
		return ""
	}
	return "<#" + c.ID.String() + ">"
}

// Mention returns the markup used to render this emoji inside message content.
// Custom emoji become "<:name:id>" (animated: "<a:name:id>"); a unicode/standard
// emoji has no ID and is returned as its raw name so it still renders.
func (e *Emoji) Mention() string {
	if e == nil {
		return ""
	}
	if e.ID.IsEmpty() {
		return e.Name
	}
	if e.Animated == nil {
		return ""
	}
	if *e.Animated {
		return "<a:" + e.Name + ":" + e.ID.String() + ">"
	}
	return "<:" + e.Name + ":" + e.ID.String() + ">"
}

// Identifier returns the form used to reference this emoji in REST reaction
// routes: "name:id" for custom emoji ("a:name:id" when animated) and the raw
// unicode character(s) for standard emoji. Callers must URL-encode the result
// before placing it in a request path.
func (e *Emoji) Identifier() *string {
	if e == nil {
		return nil
	}
	if e.ID.IsEmpty() {
		return &e.Name
	}
	if e.Animated == nil {
		return nil
	}
	if *e.Animated {
		return util.PointerOf("a:" + e.Name + ":" + e.ID.String())
	}
	return util.PointerOf(e.Name + ":" + e.ID.String())
}

// ── Identity / tags ─────────────────────────────────────────────────────────

// Tag returns the user's classic "username#discriminator" tag. For accounts
// migrated to the new username system (discriminator "0" or empty) it returns
// just the username, matching how Discord now displays handles.
func (u *User) Tag() string {
	if u == nil {
		return ""
	}
	if u.Discriminator == "" || u.Discriminator == "0" {
		return u.Username
	}
	return u.Username + "#" + u.Discriminator
}

// ── Colours ─────────────────────────────────────────────────────────────────

// HexColor returns the role's primary colour as a CSS-style "#rrggbb" string.
// A role with no colour set (Discord colour 0) returns "#000000".
func (r *Role) HexColor() string {
	if r == nil {
		return ""
	}
	c := 0
	if r.Colors.PrimaryColor != nil {
		c = *r.Colors.PrimaryColor
	}
	return colorToHex(c)
}

// HexAccentColor returns the user's profile accent colour as "#rrggbb", or ""
// when the user has no accent colour set.
func (u *User) HexAccentColor() string {
	if u == nil || u.AccentColor == nil {
		return ""
	}
	return colorToHex(*u.AccentColor)
}

// colorToHex formats a 24-bit RGB integer as "#rrggbb".
func colorToHex(c int) string {
	return fmt.Sprintf("#%06x", c&0xFFFFFF)
}

// ── Creation timestamps ─────────────────────────────────────────────────────
//
// Every Discord snowflake embeds its creation time, so these are exact and need
// no API call. User.CreatedAt already lives alongside the User type.

// CreatedAt returns the time this guild was created, decoded from its ID.
func (g *Guild) CreatedAt() time.Time { return g.ID.Time() }

// CreatedAt returns the time this channel was created, decoded from its ID.
func (c *Channel) CreatedAt() time.Time { return c.ID.Time() }

// CreatedAt returns the time this role was created, decoded from its ID.
func (r *Role) CreatedAt() time.Time { return r.ID.Time() }

// CreatedAt returns the time this message was sent, decoded from its ID.
func (m *Message) CreatedAt() time.Time { return m.ID.Time() }

// CreatedAt returns the time this custom emoji was created, decoded from its ID.
// Unicode/standard emoji have no ID and return the zero time.
func (e *Emoji) CreatedAt() time.Time {
	if e.ID.IsEmpty() {
		return time.Time{}
	}
	return e.ID.Time()
}

// CreatedAt returns the time this member's underlying account was created. It
// returns the zero time when the member's user ID is unknown.
func (m *GuildMember) CreatedAt() time.Time {
	uid := memberUserID(m)
	if uid.IsEmpty() {
		return time.Time{}
	}
	return uid.Time()
}

// ── Jump links ────────────────────────────────────────────────────────────────

// URL returns the client jump link to this channel, e.g.
// "https://discord.com/channels/<guild>/<channel>". DM/group-DM channels (which
// have no guild) use the "@me" pseudo-guild.
func (c *Channel) URL() string {
	if c == nil {
		return ""
	}
	guild := "@me"
	if !c.GuildID.IsEmpty() {
		guild = c.GuildID.String()
	}
	return "https://discord.com/channels/" + guild + "/" + c.ID.String()
}

// URL returns the discord.gg invite link for this invite, or "" when it has no
// code.
func (i *Invite) URL() string {
	if i == nil || i.Code == "" {
		return ""
	}
	return "https://discord.gg/" + i.Code
}

// ── Channel classification ──────────────────────────────────────────────────

// IsNSFW reports whether the channel is marked age-restricted.
// FIXME can return false values if c.NSFW just isn't set
func (c *Channel) IsNSFW() bool {
	return c != nil && c.NSFW != nil && *c.NSFW
}

// IsThreadOnly reports whether the channel only contains threads (forum and
// media channels).
func (c *Channel) IsThreadOnly() bool {
	return c != nil && (c.Type == ChannelTypeGuildForum || c.Type == ChannelTypeGuildMedia)
}

// ── Message classification ──────────────────────────────────────────────────

// IsSystem reports whether this is a system message (e.g. a join, pin or boost
// notice) rather than a normal user message or reply.
func (m *Message) IsSystem() bool {
	return m != nil && m.Type != MessageTypeDefault && m.Type != MessageTypeReply
}

// IsEdited reports whether the message has been edited.
func (m *Message) IsEdited() bool {
	return m != nil && m.EditedTimestamp != nil
}

// ── Member state ──────────────────────────────────────────────────────────────

// IsCommunicationDisabled reports whether the member is currently timed out
// (their communication-disabled-until timestamp is in the future).
func (m *GuildMember) IsCommunicationDisabled() bool {
	return m != nil && m.CommunicationDisabledUntil != nil && m.CommunicationDisabledUntil.After(time.Now())
}

// ── Role ordering ─────────────────────────────────────────────────────────────

// ComparePositionTo compares this role's position to another's, returning a
// positive number when this role is higher in the hierarchy, negative when it is
// lower, and zero when they are the same role. Roles sharing a position are
// ordered by ID, with the older (smaller) ID ranking higher — matching how
// Discord resolves ties.
func (r *Role) ComparePositionTo(other *Role) int {
	if r == nil || other == nil {
		return 0
	}
	if r.Position != other.Position {
		return r.Position - other.Position
	}
	switch {
	case r.ID < other.ID:
		return 1
	case r.ID > other.ID:
		return -1
	default:
		return 0
	}
}
