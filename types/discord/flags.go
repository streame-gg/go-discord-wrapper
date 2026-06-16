package discord

import (
	"encoding/json"
	"strconv"
)

type BitFlag interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// FlagBits is an immutable view over an integer bitfield of type T. The zero
// value is a valid empty bitfield. Add/Remove/Toggle return a new FlagBits
// rather than mutating, so a FlagBits is safe to copy and share.
// Example bitfield: https://docs.discord.com/developers/resources/message#message-object-message-flags
type FlagBits[T BitFlag] struct {
	bits T
}

// NewFlagBits wraps a raw bitfield value so the generic helpers can be used on
// it, e.g. discord.NewFlagBits(msg.Flags).Has(discord.MessageFlagEphemeral).
func NewFlagBits[T BitFlag](bits T) FlagBits[T] { return FlagBits[T]{bits: bits} }

// Value returns the underlying bitfield value.
func (f FlagBits[T]) Value() T { return f.bits }

// IsEmpty reports whether no bits are set.
func (f FlagBits[T]) IsEmpty() bool { var zero T; return f.bits == zero }

// Has reports whether every bit in flag is set (discord.js has()). When flag is
// a single bit this is a simple membership test; with combined bits it requires
// all of them.
func (f FlagBits[T]) Has(flag T) bool { return f.bits&flag == flag }

// HasAll reports whether all the given flags are set.
func (f FlagBits[T]) HasAll(flags ...T) bool {
	for _, fl := range flags {
		if f.bits&fl != fl {
			return false
		}
	}
	return true
}

// Any reports whether any bit in any of the given flags is set (discord.js any()).
func (f FlagBits[T]) Any(flags ...T) bool {
	for _, fl := range flags {
		if f.bits&fl != 0 {
			return true
		}
	}
	return false
}

// Equals reports whether the bitfield is exactly the union of the given flags.
func (f FlagBits[T]) Equals(flags ...T) bool {
	var combined T
	for _, fl := range flags {
		combined |= fl
	}
	return f.bits == combined
}

// Missing returns the subset of the given flags that are not set (discord.js
// missing()). The result is nil when nothing is missing.
func (f FlagBits[T]) Missing(flags ...T) []T {
	var missing []T
	for _, fl := range flags {
		if f.bits&fl != fl {
			missing = append(missing, fl)
		}
	}
	return missing
}

// Add returns a copy with all of the given flags set.
func (f FlagBits[T]) Add(flags ...T) FlagBits[T] {
	b := f.bits
	for _, fl := range flags {
		b |= fl
	}
	return FlagBits[T]{bits: b}
}

// Remove returns a copy with all of the given flags cleared.
func (f FlagBits[T]) Remove(flags ...T) FlagBits[T] {
	b := f.bits
	for _, fl := range flags {
		b &^= fl
	}
	return FlagBits[T]{bits: b}
}

// Toggle returns a copy with each of the given flags flipped.
func (f FlagBits[T]) Toggle(flags ...T) FlagBits[T] {
	b := f.bits
	for _, fl := range flags {
		b ^= fl
	}
	return FlagBits[T]{bits: b}
}

// Decompose returns each individual bit that is set, as its own value — the Go
// analogue of discord.js' toArray(). The result is nil when no bits are set.
func (f FlagBits[T]) Decompose() []T {
	var out []T
	for bit := T(1); bit != 0; bit <<= 1 {
		if f.bits&bit != 0 {
			out = append(out, bit)
		}
	}
	return out
}

// ── User flags (badges) ───────────────────────────────────────────────────────

// UserFlags is the bitfield of account badges on a user's public_flags/flags.
//
// https://docs.discord.com/developers/resources/user#user-object-user-flags
type UserFlags uint64

const (
	UserFlagStaff                 UserFlags = 1 << 0  // Discord Employee
	UserFlagPartner               UserFlags = 1 << 1  // Partnered Server Owner
	UserFlagHypeSquad             UserFlags = 1 << 2  // HypeSquad Events member
	UserFlagBugHunterLevel1       UserFlags = 1 << 3  // Bug Hunter Level 1
	UserFlagHypeSquadHouseBravery UserFlags = 1 << 6  // House of Bravery
	UserFlagHypeSquadHouseBrillia UserFlags = 1 << 7  // House of Brilliance
	UserFlagHypeSquadHouseBalance UserFlags = 1 << 8  // House of Balance
	UserFlagPremiumEarlySupporter UserFlags = 1 << 9  // Early Nitro Supporter
	UserFlagTeamPseudoUser        UserFlags = 1 << 10 // Team account user
	UserFlagBugHunterLevel2       UserFlags = 1 << 14 // Bug Hunter Level 2
	UserFlagVerifiedBot           UserFlags = 1 << 16 // Verified Bot
	UserFlagVerifiedDeveloper     UserFlags = 1 << 17 // Early Verified Bot Developer
	UserFlagCertifiedModerator    UserFlags = 1 << 18 // Moderator Programs Alumni
	UserFlagBotHTTPInteractions   UserFlags = 1 << 19 // Uses only HTTP interactions
)

// PublicFlagBits returns the user's public badge flags as a FlagBits, e.g.
// u.PublicFlagBits().Has(discord.UserFlagActiveDeveloper).
func (u *User) PublicFlagBits() FlagBits[UserFlags] {
	if u == nil {
		return FlagBits[UserFlags]{}
	}
	return NewFlagBits(u.PublicFlags)
}

// FlagBits returns the user's full flag set (the flags field) as a FlagBits.
// Most gateway/REST payloads only populate the public subset — see
// PublicFlagBits.
func (u *User) FlagBits() FlagBits[UserFlags] {
	if u == nil {
		return FlagBits[UserFlags]{}
	}
	return NewFlagBits(u.Flags)
}

// ── Guild member flags ─────────────────────────────────────────────────────────

// GuildMemberFlags is the bitfield of per-guild member state flags.
//
// https://docs.discord.com/developers/resources/guild#guild-member-object-guild-member-flags
type GuildMemberFlags int

const (
	GuildMemberFlagDidRejoin                  GuildMemberFlags = 1 << 0
	GuildMemberFlagCompletedOnboarding        GuildMemberFlags = 1 << 1
	GuildMemberFlagBypassesVerification       GuildMemberFlags = 1 << 2
	GuildMemberFlagStartedOnboarding          GuildMemberFlags = 1 << 3
	GuildMemberFlagIsGuest                    GuildMemberFlags = 1 << 4
	GuildMemberFlagStartedHomeActions         GuildMemberFlags = 1 << 5
	GuildMemberFlagCompletedHomeActions       GuildMemberFlags = 1 << 6
	GuildMemberFlagAutomodQuarantinedUsername GuildMemberFlags = 1 << 7
	GuildMemberFlagDmSettingsUpsellAcked      GuildMemberFlags = 1 << 9
	GuildMemberFlagAutomodQuarantinedGuildTag GuildMemberFlags = 1 << 10
)

// FlagBits returns the member's state flags as a FlagBits, e.g.
// m.FlagBits().Has(discord.GuildMemberFlagCompletedOnboarding).
func (m *GuildMember) FlagBits() FlagBits[GuildMemberFlags] {
	if m == nil {
		return FlagBits[GuildMemberFlags]{}
	}
	return NewFlagBits(m.Flags)
}

// ── Application flags ──────────────────────────────────────────────────────────

// ApplicationFlags is the bitfield of capabilities/gateway intents declared by
// an application. It is uint64-backed so it can hold flag bits beyond bit 30,
// which Discord transmits via the string-serialized flags_new field.
//
// https://docs.discord.com/developers/resources/application#application-object-application-flags
type ApplicationFlags uint64

// UnmarshalJSON accepts either a JSON number (the legacy flags field) or a
// decimal string (the flags_new field), so the same type backs both.
func (f *ApplicationFlags) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			*f = 0
			return nil
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*f = ApplicationFlags(n)
		return nil
	}
	var n uint64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*f = ApplicationFlags(n)
	return nil
}

const (
	ApplicationFlagApplicationAutoModerationRuleCreateBadge ApplicationFlags = 1 << 6
	ApplicationFlagGatewayPresence                          ApplicationFlags = 1 << 12
	ApplicationFlagGatewayPresenceLimited                   ApplicationFlags = 1 << 13
	ApplicationFlagGatewayGuildMembers                      ApplicationFlags = 1 << 14
	ApplicationFlagGatewayGuildMembersLimited               ApplicationFlags = 1 << 15
	ApplicationFlagVerificationPendingGuildLimit            ApplicationFlags = 1 << 16
	ApplicationFlagEmbedded                                 ApplicationFlags = 1 << 17
	ApplicationFlagGatewayMessageContent                    ApplicationFlags = 1 << 18
	ApplicationFlagGatewayMessageContentLimited             ApplicationFlags = 1 << 19
	ApplicationFlagApplicationCommandBadge                  ApplicationFlags = 1 << 23
)

// FlagBits returns the application's capability flags as a FlagBits. It prefers
// the complete flags_new value (which can carry bits beyond bit 30) and falls
// back to the legacy flags field when flags_new is absent.
func (a *Application) FlagBits() FlagBits[ApplicationFlags] {
	if a == nil {
		return FlagBits[ApplicationFlags]{}
	}
	if a.FlagsNew != nil {
		return NewFlagBits(*a.FlagsNew)
	}
	if a.Flags != nil {
		return NewFlagBits(*a.Flags)
	}
	return FlagBits[ApplicationFlags]{}
}

// ── Per-entity accessors for the already-typed bitfields ───────────────────────

// FlagBits returns the message's flags as a FlagBits, e.g.
// m.FlagBits().Has(discord.MessageFlagEphemeral).
func (m *Message) FlagBits() FlagBits[MessageFlag] {
	if m == nil {
		return FlagBits[MessageFlag]{}
	}
	return NewFlagBits(m.Flags)
}

// FlagBits returns the channel's flags as a FlagBits (empty when unset).
func (c *Channel) FlagBits() FlagBits[ChannelFlags] {
	if c == nil {
		return FlagBits[ChannelFlags]{}
	}
	return NewFlagBits(c.Flags)
}

// FlagBits returns the role's flags as a FlagBits (empty when unset).
func (r *Role) FlagBits() FlagBits[RoleFlags] {
	if r == nil || r.Flags == nil {
		return FlagBits[RoleFlags]{}
	}
	return NewFlagBits(*r.Flags)
}

// SystemChannelFlagBits returns the guild's system-channel suppression mask as a
// FlagBits, e.g. g.SystemChannelFlagBits().Has(discord.GuildSystemChannelFlagsSuppressJoinNotifications).
func (g *Guild) SystemChannelFlagBits() FlagBits[GuildSystemChannelFlags] {
	if g == nil {
		return FlagBits[GuildSystemChannelFlags]{}
	}
	return NewFlagBits(g.SystemChannelFlags)
}
