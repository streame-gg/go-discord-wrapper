package discord_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// entityMethodsSuite exercises the entity convenience methods (Channel.Send,
// Message.Pin, GuildMember.Ban, …) along both branches every wrapper shares:
//
//   - the ensureClient guard: an un-hydrated entity must return
//     ErrEntityNotHydrated without touching the network; and
//   - the hydrated happy path: once a client is injected, the call routes to
//     the corresponding EntityClient method and returns its result.
//
// It reuses the stubClient from hydration_test.go (same test package), which
// returns zero values for every method, so the happy-path calls succeed.
type entityMethodsSuite struct {
	suite.Suite
}

func TestEntityMethodsSuite(t *testing.T) {
	suite.Run(t, new(entityMethodsSuite))
}

const (
	emChan  = discord.Snowflake(500000000000000001)
	emMsg   = discord.Snowflake(500000000000000002)
	emGuild = discord.Snowflake(500000000000000003)
	emUser  = discord.Snowflake(500000000000000004)
	emRole  = discord.Snowflake(500000000000000005)
	emEmoji = discord.Snowflake(500000000000000006)
	emRule  = discord.Snowflake(500000000000000007)
	emEvent = discord.Snowflake(500000000000000008)
)

// entityCall pairs a human label with a closure that invokes one entity method.
// The same closures are run twice: once on an un-hydrated entity (expecting
// ErrEntityNotHydrated) and once on a hydrated one (expecting success).
type entityCall struct {
	name string
	// build returns a fresh entity and the method invocation against it. The
	// hydrate flag controls whether the entity is given a client first.
	run func(hydrate bool) error
}

func hydrateMaybe(hydrate bool, hydrater func(discord.EntityClient)) {
	if hydrate {
		hydrater(&stubClient{})
	}
}

func (s *entityMethodsSuite) calls() []entityCall {
	ctx := context.Background()
	return []entityCall{
		// ── Channel ──
		{"Channel.Send", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			_, err := ch.Send(ctx, discord.MessageCreateOptions{})
			return err
		}},
		{"Channel.Delete", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			_, err := ch.Delete(ctx, nil)
			return err
		}},
		{"Channel.BulkDelete", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			return ch.BulkDelete(ctx, []discord.Snowflake{emMsg}, nil)
		}},
		{"Channel.FetchMessages", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			_, err := ch.FetchMessages(ctx, discord.FetchMessagesOptions{})
			return err
		}},
		{"Channel.TriggerTyping", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			return ch.TriggerTyping(ctx)
		}},
		{"Channel.SetVoiceStatus", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			return ch.SetVoiceStatus(ctx, nil)
		}},
		{"Channel.CreateInvite", func(h bool) error {
			ch := &discord.Channel{ID: emChan}
			hydrateMaybe(h, ch.Hydrate)
			_, err := ch.CreateInvite(ctx, discord.InviteCreateOptions{})
			return err
		}},

		// ── Message ──
		{"Message.Reply", func(h bool) error {
			m := &discord.Message{ID: emMsg, ChannelID: emChan}
			hydrateMaybe(h, m.Hydrate)
			_, err := m.Reply(ctx, discord.MessageCreateOptions{})
			return err
		}},
		{"Message.Pin", func(h bool) error {
			m := &discord.Message{ID: emMsg, ChannelID: emChan}
			hydrateMaybe(h, m.Hydrate)
			return m.Pin(ctx, nil)
		}},
		{"Message.Unpin", func(h bool) error {
			m := &discord.Message{ID: emMsg, ChannelID: emChan}
			hydrateMaybe(h, m.Hydrate)
			return m.Unpin(ctx, nil)
		}},
		{"Message.CrossPost", func(h bool) error {
			m := &discord.Message{ID: emMsg, ChannelID: emChan}
			hydrateMaybe(h, m.Hydrate)
			_, err := m.CrossPost(ctx)
			return err
		}},
		{"Message.React", func(h bool) error {
			m := &discord.Message{ID: emMsg, ChannelID: emChan}
			hydrateMaybe(h, m.Hydrate)
			return m.React(ctx, "🔥")
		}},

		// ── GuildMember ──
		{"GuildMember.Edit", func(h bool) error {
			m := &discord.GuildMember{GuildID: emGuild, UserID: emUser}
			hydrateMaybe(h, m.Hydrate)
			_, err := m.Edit(ctx, discord.MemberEditOptions{})
			return err
		}},
		{"GuildMember.Ban", func(h bool) error {
			m := &discord.GuildMember{GuildID: emGuild, UserID: emUser}
			hydrateMaybe(h, m.Hydrate)
			return m.Ban(ctx, discord.BanOptions{})
		}},
		{"GuildMember.AddRole", func(h bool) error {
			m := &discord.GuildMember{GuildID: emGuild, UserID: emUser}
			hydrateMaybe(h, m.Hydrate)
			return m.AddRole(ctx, emRole, nil)
		}},
		{"GuildMember.RemoveRole", func(h bool) error {
			m := &discord.GuildMember{GuildID: emGuild, UserID: emUser}
			hydrateMaybe(h, m.Hydrate)
			return m.RemoveRole(ctx, emRole, nil)
		}},
		{"GuildMember.Timeout", func(h bool) error {
			m := &discord.GuildMember{GuildID: emGuild, UserID: emUser}
			hydrateMaybe(h, m.Hydrate)
			_, err := m.Timeout(ctx, time.Now().Add(time.Hour), nil)
			return err
		}},

		// ── Role ──
		{"Role.Edit", func(h bool) error {
			r := &discord.Role{ID: emRole, GuildID: emGuild}
			hydrateMaybe(h, r.Hydrate)
			_, err := r.Edit(ctx, discord.RoleEditOptions{})
			return err
		}},
		{"Role.SetPosition", func(h bool) error {
			r := &discord.Role{ID: emRole, GuildID: emGuild}
			hydrateMaybe(h, r.Hydrate)
			_, err := r.SetPosition(ctx, discord.RolePositionOptions{})
			return err
		}},

		// ── Guild ──
		{"Guild.Edit", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.Edit(ctx, discord.GuildEditOptions{})
			return err
		}},
		{"Guild.FetchMembers", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.FetchMembers(ctx, discord.FetchMembersOptions{})
			return err
		}},
		{"Guild.CreateRole", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.CreateRole(ctx, discord.RoleCreateOptions{})
			return err
		}},
		{"Guild.CreateChannel", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.CreateChannel(ctx, discord.ChannelCreateOptions{})
			return err
		}},
		{"Guild.CreateEmoji", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.CreateEmoji(ctx, discord.EmojiCreateOptions{})
			return err
		}},
		{"Guild.CreateSticker", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.CreateSticker(ctx, discord.StickerCreateOptions{})
			return err
		}},
		{"Guild.FetchAuditLog", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.FetchAuditLog(ctx, discord.AuditLogOptions{})
			return err
		}},
		{"Guild.CreateScheduledEvent", func(h bool) error {
			g := &discord.Guild{ID: emGuild}
			hydrateMaybe(h, g.Hydrate)
			_, err := g.CreateScheduledEvent(ctx, discord.ScheduledEventCreateOptions{})
			return err
		}},

		// ── Emoji ──
		{"Emoji.Edit", func(h bool) error {
			e := &discord.Emoji{ID: emEmoji, GuildID: emGuild}
			hydrateMaybe(h, e.Hydrate)
			_, err := e.Edit(ctx, discord.EmojiEditOptions{})
			return err
		}},
		{"Emoji.Delete", func(h bool) error {
			e := &discord.Emoji{ID: emEmoji, GuildID: emGuild}
			hydrateMaybe(h, e.Hydrate)
			return e.Delete(ctx, nil)
		}},

		// ── AutoModerationRule ──
		{"AutoModerationRule.Edit", func(h bool) error {
			r := &discord.AutoModerationRule{ID: emRule, GuildID: emGuild}
			hydrateMaybe(h, r.Hydrate)
			_, err := r.Edit(ctx, discord.RuleEditOptions{})
			return err
		}},
		{"AutoModerationRule.Delete", func(h bool) error {
			r := &discord.AutoModerationRule{ID: emRule, GuildID: emGuild}
			hydrateMaybe(h, r.Hydrate)
			return r.Delete(ctx)
		}},

		// ── GuildScheduledEvent ──
		{"GuildScheduledEvent.Edit", func(h bool) error {
			e := &discord.GuildScheduledEvent{ID: emEvent, GuildID: emGuild}
			hydrateMaybe(h, e.Hydrate)
			_, err := e.Edit(ctx, discord.ScheduledEventEditOptions{})
			return err
		}},
		{"GuildScheduledEvent.Delete", func(h bool) error {
			e := &discord.GuildScheduledEvent{ID: emEvent, GuildID: emGuild}
			hydrateMaybe(h, e.Hydrate)
			return e.Delete(ctx)
		}},
		{"GuildScheduledEvent.FetchUsers", func(h bool) error {
			e := &discord.GuildScheduledEvent{ID: emEvent, GuildID: emGuild}
			hydrateMaybe(h, e.Hydrate)
			_, err := e.FetchUsers(ctx, discord.FetchUsersOptions{})
			return err
		}},
	}
}

// TestUnhydratedReturnsGuardError verifies every entity method refuses to act
// without a client, returning ErrEntityNotHydrated.
func (s *entityMethodsSuite) TestUnhydratedReturnsGuardError() {
	for _, c := range s.calls() {
		s.Run(c.name, func() {
			err := c.run(false)
			s.ErrorIsf(err, discord.ErrEntityNotHydrated, "%s on un-hydrated entity must return ErrEntityNotHydrated", c.name)
		})
	}
}

// TestHydratedRoutesWithoutError verifies every entity method routes through the
// injected client and succeeds once hydrated.
func (s *entityMethodsSuite) TestHydratedRoutesWithoutError() {
	for _, c := range s.calls() {
		s.Run(c.name, func() {
			err := c.run(true)
			s.NoErrorf(err, "%s on hydrated entity should route without error", c.name)
		})
	}
}
