package discord

// permissionNames maps each defined Permission bit to its Discord flag name, in
// ascending bit order. It also defines the set of bits that make up PermissionAll.
var permissionNames = []struct {
	bit  Permission
	name string
}{
	{PermissionCreateInstantInvite, "CreateInstantInvite"},
	{PermissionKickMembers, "KickMembers"},
	{PermissionBanMembers, "BanMembers"},
	{PermissionAdministrator, "Administrator"},
	{PermissionManageChannels, "ManageChannels"},
	{PermissionManageGuild, "ManageGuild"},
	{PermissionAddReactions, "AddReactions"},
	{PermissionViewAuditLog, "ViewAuditLog"},
	{PermissionPrioritySpeaker, "PrioritySpeaker"},
	{PermissionStream, "Stream"},
	{PermissionViewChannel, "ViewChannel"},
	{PermissionSendMessages, "SendMessages"},
	{PermissionSendTTSMessages, "SendTTSMessages"},
	{PermissionManageMessages, "ManageMessages"},
	{PermissionEmbedLinks, "EmbedLinks"},
	{PermissionAttachFiles, "AttachFiles"},
	{PermissionReadMessageHistory, "ReadMessageHistory"},
	{PermissionMentionEveryone, "MentionEveryone"},
	{PermissionUseExternalEmojis, "UseExternalEmojis"},
	{PermissionViewGuildInsights, "ViewGuildInsights"},
	{PermissionConnect, "Connect"},
	{PermissionSpeak, "Speak"},
	{PermissionMuteMembers, "MuteMembers"},
	{PermissionDeafenMembers, "DeafenMembers"},
	{PermissionMoveMembers, "MoveMembers"},
	{PermissionUseVAD, "UseVAD"},
	{PermissionChangeNickname, "ChangeNickname"},
	{PermissionManageNicknames, "ManageNicknames"},
	{PermissionManageRoles, "ManageRoles"},
	{PermissionManageWebhooks, "ManageWebhooks"},
	{PermissionManageGuildExpressions, "ManageGuildExpressions"},
	{PermissionUseApplicationCommands, "UseApplicationCommands"},
	{PermissionRequestToSpeak, "RequestToSpeak"},
	{PermissionManageEvents, "ManageEvents"},
	{PermissionManageThreads, "ManageThreads"},
	{PermissionCreatePublicThreads, "CreatePublicThreads"},
	{PermissionCreatePrivateThreads, "CreatePrivateThreads"},
	{PermissionUseExternalStickers, "UseExternalStickers"},
	{PermissionSendMessagesInThreads, "SendMessagesInThreads"},
	{PermissionUseEmbeddedActivities, "UseEmbeddedActivities"},
	{PermissionModerateMembers, "ModerateMembers"},
	{PermissionViewCreatorMonetizationAnalytics, "ViewCreatorMonetizationAnalytics"},
	{PermissionUseSoundboard, "UseSoundboard"},
	{PermissionCreateGuildExpressions, "CreateGuildExpressions"},
	{PermissionCreateEvents, "CreateEvents"},
	{PermissionUseExternalSounds, "UseExternalSounds"},
	{PermissionSendVoiceMessages, "SendVoiceMessages"},
	{PermissionSendPolls, "SendPolls"},
	{PermissionUseExternalApps, "UseExternalApps"},
	{PermissionPinMessages, "PinMessages"},
	{PermissionBypassSlowmode, "BypassSlowmode"},
}

// PermissionAll is the bitwise OR of every permission flag defined by this
// library. A member with Administrator (or the guild owner) is treated as
// holding PermissionAll by the permission-calculation helpers.
var PermissionAll = func() Permission {
	var all Permission
	for _, p := range permissionNames {
		all |= p.bit
	}
	return all
}()

// ── Bitmask helpers ─────────────────────────────────────────────────────────

// Has reports whether p contains every bit set in bits. Has(0) is always true.
// It does not special-case Administrator; compute effective permissions with
// PermissionsFor / PermissionsInChannel first if you want Administrator to imply
// all permissions.
func (p Permission) Has(bits Permission) bool { return p&bits == bits }

// HasAny reports whether p contains at least one bit set in bits.
func (p Permission) HasAny(bits Permission) bool { return p&bits != 0 }

// Add returns a copy of p with all of the given bits set.
func (p Permission) Add(bits ...Permission) Permission {
	for _, b := range bits {
		p |= b
	}
	return p
}

// Remove returns a copy of p with all of the given bits cleared.
func (p Permission) Remove(bits ...Permission) Permission {
	for _, b := range bits {
		p &^= b
	}
	return p
}

// Missing returns the subset of bits that p does not contain.
func (p Permission) Missing(bits Permission) Permission { return bits &^ p }

// Names returns the flag names of every permission bit set in p, in ascending
// bit order. Bits that do not correspond to a known flag are ignored.
func (p Permission) Names() []string {
	var names []string
	for _, e := range permissionNames {
		if p&e.bit == e.bit {
			names = append(names, e.name)
		}
	}
	return names
}

// AllowPermissions returns the overwrite's allowed permission bits.
func (o ChannelPermissionOverwrite) AllowPermissions() Permission {
	return o.Allow
}

// DenyPermissions returns the overwrite's denied permission bits.
func (o ChannelPermissionOverwrite) DenyPermissions() Permission {
	return o.Deny
}

// ── Permission calculation ────────────────────────────────────────────────────

// roleByID returns the role with the given id from the guild's raw role list, or
// nil if the guild is nil or has no such role.
func (g *Guild) roleByID(id Snowflake) *Role {
	if g == nil {
		return nil
	}
	for i := range g.RawRoles {
		if g.RawRoles[i].ID == id {
			return &g.RawRoles[i]
		}
	}
	return nil
}

// memberUserID resolves the member's user ID, falling back to the embedded User.
func memberUserID(m *GuildMember) Snowflake {
	if m == nil {
		return 0
	}
	if !m.UserID.IsEmpty() {
		return m.UserID
	}
	if m.User != nil {
		return m.User.ID
	}
	return 0
}

// overwriteFor returns the channel permission overwrite targeting id, or nil.
func overwriteFor(channel *Channel, id Snowflake) *ChannelPermissionOverwrite {
	if channel == nil || id.IsEmpty() {
		return nil
	}
	for i := range channel.PermissionOverwrites {
		if channel.PermissionOverwrites[i].ID == id {
			return &channel.PermissionOverwrites[i]
		}
	}
	return nil
}

// PermissionsFor computes the guild-level (base) permissions for member in
// guild. It unions the @everyone role with every role the member holds. The
// guild owner and any member holding Administrator are granted PermissionAll.
// It returns 0 if guild or member is nil.
//
// Permissions are resolved from guild.RawRoles, so the guild must carry its role
// list (as it does when received from the gateway or REST API).
func PermissionsFor(guild *Guild, member *GuildMember) Permission {
	if guild == nil || member == nil {
		return 0
	}

	if !guild.OwnerID.IsEmpty() && guild.OwnerID == memberUserID(member) {
		return PermissionAll
	}

	// @everyone role shares the guild's ID.
	var perms Permission
	if everyone := guild.roleByID(guild.ID); everyone != nil {
		perms = everyone.Permissions
	}

	for _, roleID := range member.Roles {
		if r := guild.roleByID(roleID); r != nil {
			perms |= r.Permissions
		}
	}

	if perms.Has(PermissionAdministrator) {
		return PermissionAll
	}
	return perms
}

// PermissionsInChannel computes the effective permissions for member in channel,
// applying channel permission overwrites on top of the base guild permissions in
// Discord's documented order: @everyone overwrite, then the union of role
// overwrites (all denies before all allows), then the member-specific overwrite.
//
// The guild owner and any member holding Administrator are granted PermissionAll
// and channel overwrites are ignored. It returns 0 if guild, member, or channel
// is nil.
func PermissionsInChannel(guild *Guild, member *GuildMember, channel *Channel) Permission {
	if guild == nil || member == nil || channel == nil {
		return 0
	}

	base := PermissionsFor(guild, member)
	if base.Has(PermissionAdministrator) {
		return PermissionAll
	}

	perms := base

	// 1. @everyone overwrite (ID == guild ID).
	if ow := overwriteFor(channel, guild.ID); ow != nil {
		perms = perms.Remove(ow.DenyPermissions())
		perms = perms.Add(ow.AllowPermissions())
	}

	// 2. Role overwrites: accumulate all denies and all allows, then apply
	//    denies before allows so an allow on any role wins over a deny on another.
	var allow, deny Permission
	for _, roleID := range member.Roles {
		if ow := overwriteFor(channel, roleID); ow != nil {
			allow |= ow.AllowPermissions()
			deny |= ow.DenyPermissions()
		}
	}
	perms = perms.Remove(deny)
	perms = perms.Add(allow)

	// 3. Member-specific overwrite.
	if ow := overwriteFor(channel, memberUserID(member)); ow != nil {
		perms = perms.Remove(ow.DenyPermissions())
		perms = perms.Add(ow.AllowPermissions())
	}

	return perms
}

// Permissions returns the member's guild-level permissions within guild.
// It is a convenience wrapper around PermissionsFor.
func (m *GuildMember) Permissions(guild *Guild) Permission {
	return PermissionsFor(guild, m)
}

// PermissionsIn returns the member's effective permissions in channel, including
// channel overwrites. It is a convenience wrapper around PermissionsInChannel.
func (m *GuildMember) PermissionsIn(guild *Guild, channel *Channel) Permission {
	return PermissionsInChannel(guild, m, channel)
}
