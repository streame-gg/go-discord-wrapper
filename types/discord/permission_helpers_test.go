package discord

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (su *permissionHelpersSuite) TestPermissionBitmaskHelpers() {
	t := su.T()
	p := PermissionViewChannel | PermissionSendMessages

	if !p.Has(PermissionViewChannel) {
		t.Error("Has: want ViewChannel present")
	}
	if !p.Has(PermissionViewChannel | PermissionSendMessages) {
		t.Error("Has: want both bits present")
	}
	if p.Has(PermissionViewChannel | PermissionBanMembers) {
		t.Error("Has: must require ALL bits, BanMembers is absent")
	}
	if !p.Has(0) {
		t.Error("Has(0) should be true")
	}
	if !p.HasAny(PermissionViewChannel | PermissionBanMembers) {
		t.Error("HasAny: ViewChannel present, want true")
	}
	if p.HasAny(PermissionBanMembers | PermissionKickMembers) {
		t.Error("HasAny: none present, want false")
	}

	added := p.Add(PermissionBanMembers)
	if !added.Has(PermissionBanMembers) || !added.Has(PermissionViewChannel) {
		t.Error("Add: should set new bit and keep existing")
	}
	if p.Has(PermissionBanMembers) {
		t.Error("Add must not mutate the receiver")
	}

	removed := p.Remove(PermissionSendMessages)
	if removed.Has(PermissionSendMessages) {
		t.Error("Remove: SendMessages should be cleared")
	}
	if !removed.Has(PermissionViewChannel) {
		t.Error("Remove: ViewChannel should remain")
	}

	missing := p.Missing(PermissionViewChannel | PermissionBanMembers | PermissionKickMembers)
	want := PermissionBanMembers | PermissionKickMembers
	if missing != want {
		t.Errorf("Missing: got %d, want %d", missing, want)
	}
}

func (su *permissionHelpersSuite) TestPermissionNamesAndString() {
	t := su.T()
	p := PermissionBanMembers | PermissionViewChannel
	names := p.Names()
	if len(names) != 2 || names[0] != "BanMembers" || names[1] != "ViewChannel" {
		t.Errorf("Names: got %v, want [BanMembers ViewChannel] in bit order", names)
	}
}

func (su *permissionHelpersSuite) TestPermissionAll() {
	t := su.T()
	if !PermissionAll.Has(PermissionAdministrator) {
		t.Error("PermissionAll must include Administrator")
	}
	if !PermissionAll.Has(PermissionBypassSlowmode) {
		t.Error("PermissionAll must include the highest defined bit")
	}
}

func (su *permissionHelpersSuite) TestOverwritePermissionParsing() {
	t := su.T()
	ow := ChannelPermissionOverwrite{
		Allow: PermissionSendMessages,
		Deny:  PermissionAddReactions,
	}
	if ow.AllowPermissions() != PermissionSendMessages {
		t.Errorf("AllowPermissions: got %d", ow.AllowPermissions())
	}
	if ow.DenyPermissions() != PermissionAddReactions {
		t.Errorf("DenyPermissions: got %d", ow.DenyPermissions())
	}

	empty := ChannelPermissionOverwrite{}
	if empty.AllowPermissions() != 0 || empty.DenyPermissions() != 0 {
		t.Error("empty overwrite should parse to 0")
	}
}

const (
	testGuildID = Snowflake(1) // @everyone role shares this ID
	testRoleA   = Snowflake(100)
	testRoleB   = Snowflake(200)
	testOwnerID = Snowflake(999)
	testUserID  = Snowflake(555)
)

// makeGuild builds a guild with an @everyone role plus roleA/roleB.
func makeGuild(everyone, roleA, roleB Permission) *Guild {
	return &Guild{
		ID:      testGuildID,
		OwnerID: testOwnerID,
		RawRoles: []Role{
			{ID: testGuildID, Permissions: everyone},
			{ID: testRoleA, Permissions: roleA},
			{ID: testRoleB, Permissions: roleB},
		},
	}
}

func (su *permissionHelpersSuite) TestPermissionsFor() {
	t := su.T()
	t.Run("nil inputs", func(t *testing.T) {
		if PermissionsFor(nil, &GuildMember{}) != 0 {
			t.Error("nil guild should yield 0")
		}
		if PermissionsFor(makeGuild(0, 0, 0), nil) != 0 {
			t.Error("nil member should yield 0")
		}
	})

	t.Run("everyone base only", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel|PermissionSendMessages, 0, 0)
		m := &GuildMember{UserID: testUserID}
		got := PermissionsFor(g, m)
		if got != PermissionViewChannel|PermissionSendMessages {
			t.Errorf("got %d", got)
		}
	})

	t.Run("union of roles", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel, PermissionSendMessages, PermissionBanMembers)
		m := &GuildMember{UserID: testUserID, Roles: []Snowflake{testRoleA, testRoleB}}
		got := PermissionsFor(g, m)
		want := PermissionViewChannel | PermissionSendMessages | PermissionBanMembers
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("administrator role grants all", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel, PermissionAdministrator, 0)
		m := &GuildMember{UserID: testUserID, Roles: []Snowflake{testRoleA}}
		if PermissionsFor(g, m) != PermissionAll {
			t.Error("admin role should grant PermissionAll")
		}
	})

	t.Run("owner grants all", func(t *testing.T) {
		g := makeGuild(0, 0, 0)
		m := &GuildMember{UserID: testOwnerID}
		if PermissionsFor(g, m) != PermissionAll {
			t.Error("guild owner should get PermissionAll")
		}
	})

	t.Run("user id resolved from embedded user", func(t *testing.T) {
		g := makeGuild(0, 0, 0)
		m := &GuildMember{User: &User{ID: testOwnerID}}
		if PermissionsFor(g, m) != PermissionAll {
			t.Error("owner detection should use embedded User.ID when UserID empty")
		}
	})
}

func (su *permissionHelpersSuite) TestPermissionsInChannel() {
	t := su.T()
	t.Run("nil channel", func(t *testing.T) {
		if PermissionsInChannel(makeGuild(0, 0, 0), &GuildMember{}, nil) != 0 {
			t.Error("nil channel should yield 0")
		}
	})

	t.Run("everyone overwrite deny", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel|PermissionSendMessages, 0, 0)
		m := &GuildMember{UserID: testUserID}
		ch := &Channel{PermissionOverwrites: []ChannelPermissionOverwrite{
			{ID: testGuildID, Type: PermissionOverwriteTypeRole, Deny: PermissionSendMessages},
		}}
		got := PermissionsInChannel(g, m, ch)
		if got.Has(PermissionSendMessages) {
			t.Error("@everyone deny should remove SendMessages")
		}
		if !got.Has(PermissionViewChannel) {
			t.Error("ViewChannel should remain")
		}
	})

	t.Run("role allow overrides everyone deny", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel|PermissionSendMessages, 0, 0)
		m := &GuildMember{UserID: testUserID, Roles: []Snowflake{testRoleA}}
		ch := &Channel{PermissionOverwrites: []ChannelPermissionOverwrite{
			{ID: testGuildID, Type: PermissionOverwriteTypeRole, Deny: PermissionSendMessages},
			{ID: testRoleA, Type: PermissionOverwriteTypeRole, Allow: PermissionSendMessages},
		}}
		if !PermissionsInChannel(g, m, ch).Has(PermissionSendMessages) {
			t.Error("role allow should re-grant SendMessages denied at @everyone level")
		}
	})

	t.Run("role allow overrides role deny", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel, 0, 0)
		m := &GuildMember{UserID: testUserID, Roles: []Snowflake{testRoleA, testRoleB}}
		ch := &Channel{PermissionOverwrites: []ChannelPermissionOverwrite{
			{ID: testRoleA, Type: PermissionOverwriteTypeRole, Deny: PermissionSendMessages},
			{ID: testRoleB, Type: PermissionOverwriteTypeRole, Allow: PermissionSendMessages},
		}}
		if !PermissionsInChannel(g, m, ch).Has(PermissionSendMessages) {
			t.Error("within role tier, an allow on any role wins over a deny on another")
		}
	})

	t.Run("member overwrite has highest priority", func(t *testing.T) {
		g := makeGuild(PermissionViewChannel|PermissionSendMessages, 0, 0)
		m := &GuildMember{UserID: testUserID, Roles: []Snowflake{testRoleA}}
		ch := &Channel{PermissionOverwrites: []ChannelPermissionOverwrite{
			{ID: testRoleA, Type: PermissionOverwriteTypeRole, Allow: PermissionSendMessages},
			{ID: testUserID, Type: PermissionOverwriteTypeUser, Deny: PermissionSendMessages},
		}}
		if PermissionsInChannel(g, m, ch).Has(PermissionSendMessages) {
			t.Error("member-specific deny must override role allow")
		}
	})

	t.Run("administrator ignores overwrites", func(t *testing.T) {
		g := makeGuild(PermissionAdministrator, 0, 0)
		m := &GuildMember{UserID: testUserID}
		ch := &Channel{PermissionOverwrites: []ChannelPermissionOverwrite{
			{ID: testGuildID, Type: PermissionOverwriteTypeRole, Deny: PermissionViewChannel},
		}}
		if PermissionsInChannel(g, m, ch) != PermissionAll {
			t.Error("Administrator should yield PermissionAll regardless of overwrites")
		}
	})
}

func (su *permissionHelpersSuite) TestGuildMemberPermissionMethods() {
	t := su.T()
	g := makeGuild(PermissionViewChannel|PermissionSendMessages, 0, 0)
	m := &GuildMember{UserID: testUserID}

	if *m.Permissions != PermissionsFor(g, m) {
		t.Error("Permissions method should match PermissionsFor")
	}

	ch := &Channel{PermissionOverwrites: []ChannelPermissionOverwrite{
		{ID: testGuildID, Type: PermissionOverwriteTypeRole, Deny: PermissionSendMessages},
	}}
	if m.PermissionsIn(g, ch) != PermissionsInChannel(g, m, ch) {
		t.Error("PermissionsIn method should match PermissionsInChannel")
	}
}

type permissionHelpersSuite struct{ suite.Suite }

func TestPermissionHelpersSuite(t *testing.T) { suite.Run(t, new(permissionHelpersSuite)) }
