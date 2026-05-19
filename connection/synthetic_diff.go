package connection

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// deriveSyntheticEvents returns the synthetic events to dispatch after ev.
// Called synchronously on the websocket reader goroutine after internalEventHandler
// has already applied cache updates and populated Old fields.
func (d *Client) deriveSyntheticEvents(ev events.Event) []events.Event {
	switch e := ev.(type) {
	case *events.VoiceStateUpdateEvent:
		return deriveVoiceSyntheticEvents(e)
	case *events.GuildMemberUpdateEvent:
		return deriveGuildMemberSyntheticEvents(e)
	case *events.PresenceUpdateEvent:
		return derivePresenceSyntheticEvents(e)
	case *events.GuildEmojisUpdateEvent:
		return deriveGuildEmojiSyntheticEvents(e)
	case *events.GuildStickersUpdateEvent:
		return deriveGuildStickerSyntheticEvents(e)
	case *events.GuildRoleUpdateEvent:
		return deriveGuildRoleSyntheticEvents(e)
	}
	return nil
}

// deriveGuildEmojiSyntheticEvents diffs the old and new emoji sets from a
// GUILD_EMOJIS_UPDATE and fires Add/Remove/Update events per changed emoji.
// Returns nil when OldEmojis is nil (cache cold).
func deriveGuildEmojiSyntheticEvents(ev *events.GuildEmojisUpdateEvent) []events.Event {
	if ev.OldEmojis == nil {
		return nil
	}

	oldByID := make(map[discord.Snowflake]*discord.Emoji, len(ev.OldEmojis))
	for _, e := range ev.OldEmojis {
		oldByID[e.ID] = e
	}
	newByID := make(map[discord.Snowflake]*discord.Emoji, len(ev.Emojis))
	for i := range ev.Emojis {
		newByID[ev.Emojis[i].ID] = &ev.Emojis[i]
	}

	var result []events.Event

	for id, newEmoji := range newByID {
		if old, exists := oldByID[id]; !exists {
			result = append(result, &events.GuildEmojiAddEvent{
				GuildID: ev.GuildID,
				Emoji:   newEmoji,
			})
		} else if old.Name != newEmoji.Name {
			result = append(result, &events.GuildEmojiUpdateEvent{
				GuildID:  ev.GuildID,
				OldEmoji: old,
				NewEmoji: newEmoji,
			})
		}
	}
	for id, oldEmoji := range oldByID {
		if _, exists := newByID[id]; !exists {
			result = append(result, &events.GuildEmojiRemoveEvent{
				GuildID: ev.GuildID,
				Emoji:   oldEmoji,
			})
		}
	}

	return result
}

// deriveGuildStickerSyntheticEvents diffs the old and new sticker sets from a
// GUILD_STICKERS_UPDATE and fires Add/Remove/Update events per changed sticker.
// Returns nil when OldStickers is nil (cache cold).
func deriveGuildStickerSyntheticEvents(ev *events.GuildStickersUpdateEvent) []events.Event {
	if ev.OldStickers == nil {
		return nil
	}

	oldByID := make(map[discord.Snowflake]*discord.Sticker, len(ev.OldStickers))
	for _, s := range ev.OldStickers {
		oldByID[s.ID] = s
	}
	newByID := make(map[discord.Snowflake]*discord.Sticker, len(ev.Stickers))
	for i := range ev.Stickers {
		newByID[ev.Stickers[i].ID] = &ev.Stickers[i]
	}

	var result []events.Event

	for id, newSticker := range newByID {
		if old, exists := oldByID[id]; !exists {
			result = append(result, &events.GuildStickerAddEvent{
				GuildID: ev.GuildID,
				Sticker: newSticker,
			})
		} else if old.Name != newSticker.Name {
			result = append(result, &events.GuildStickerUpdateEvent{
				GuildID:    ev.GuildID,
				OldSticker: old,
				NewSticker: newSticker,
			})
		}
	}
	for id, oldSticker := range oldByID {
		if _, exists := newByID[id]; !exists {
			result = append(result, &events.GuildStickerRemoveEvent{
				GuildID: ev.GuildID,
				Sticker: oldSticker,
			})
		}
	}

	return result
}

// deriveGuildRoleSyntheticEvents fires GuildRolePermissionsChangeEvent when a
// role's Permissions bitfield changes. Returns nil when OldRole is nil (cache cold).
func deriveGuildRoleSyntheticEvents(ev *events.GuildRoleUpdateEvent) []events.Event {
	if ev.OldRole == nil {
		return nil
	}
	if ev.OldRole.Permissions == ev.Role.Permissions {
		return nil
	}
	role := ev.Role
	return []events.Event{&events.GuildRolePermissionsChangeEvent{
		GuildID:        ev.GuildID,
		RoleID:         ev.Role.ID,
		OldPermissions: ev.OldRole.Permissions,
		NewPermissions: ev.Role.Permissions,
		OldRole:        ev.OldRole,
		NewRole:        &role,
	}}
}

// deriveGuildMemberSyntheticEvents derives member synthetic events from a
// GuildMemberUpdateEvent. Returns nil if OldMember is nil (cache miss).
// Multiple events may fire for a single update (e.g. role add + nick change).
func deriveGuildMemberSyntheticEvents(ev *events.GuildMemberUpdateEvent) []events.Event {
	if ev.OldMember == nil {
		return nil
	}

	newMember := buildNewMember(ev)
	var result []events.Event

	// Role diff — compare old []string with new []Snowflake.
	oldRoles := make(map[string]struct{}, len(ev.OldMember.Roles))
	for _, r := range ev.OldMember.Roles {
		oldRoles[r] = struct{}{}
	}
	newRoles := make(map[string]struct{}, len(ev.Roles))
	for _, r := range ev.Roles {
		newRoles[string(r)] = struct{}{}
	}

	for r := range newRoles {
		if _, exists := oldRoles[r]; !exists {
			result = append(result, &events.GuildMemberRoleAddEvent{
				GuildID: ev.GuildID, UserID: ev.User.ID,
				RoleID:    discord.Snowflake(r),
				OldMember: ev.OldMember, NewMember: newMember,
			})
		}
	}
	for r := range oldRoles {
		if _, exists := newRoles[r]; !exists {
			result = append(result, &events.GuildMemberRoleRemoveEvent{
				GuildID: ev.GuildID, UserID: ev.User.ID,
				RoleID:    discord.Snowflake(r),
				OldMember: ev.OldMember, NewMember: newMember,
			})
		}
	}

	// Nick change.
	if nickDiffers(ev.OldMember.Nick, ev.Nick) {
		result = append(result, &events.GuildMemberNickChangeEvent{
			GuildID: ev.GuildID, UserID: ev.User.ID,
			OldNick: ev.OldMember.Nick, NewNick: ev.Nick,
			OldMember: ev.OldMember, NewMember: newMember,
		})
	}

	// Timeout applied or extended.
	if ev.CommunicationDisabledUntil != nil &&
		(ev.OldMember.CommunicationDisabledUntil == nil ||
			*ev.OldMember.CommunicationDisabledUntil != *ev.CommunicationDisabledUntil) {
		result = append(result, &events.GuildMemberTimeoutEvent{
			GuildID: ev.GuildID, UserID: ev.User.ID,
			CommunicationDisabledUntil: *ev.CommunicationDisabledUntil,
			OldMember:                  ev.OldMember, NewMember: newMember,
		})
	}

	// Boost start.
	if ev.PremiumSince != nil && ev.OldMember.PremiumSince == nil {
		result = append(result, &events.GuildMemberBoostStartEvent{
			GuildID: ev.GuildID, UserID: ev.User.ID,
			PremiumSince: *ev.PremiumSince,
			OldMember:    ev.OldMember, NewMember: newMember,
		})
	}

	// Boost end.
	if ev.PremiumSince == nil && ev.OldMember.PremiumSince != nil {
		result = append(result, &events.GuildMemberBoostEndEvent{
			GuildID:   ev.GuildID, UserID: ev.User.ID,
			OldMember: ev.OldMember, NewMember: newMember,
		})
	}

	return result
}

// buildNewMember constructs the post-update GuildMember from a GUILD_MEMBER_UPDATE event,
// starting from OldMember as a base and applying the update fields on top.
func buildNewMember(ev *events.GuildMemberUpdateEvent) *discord.GuildMember {
	var m discord.GuildMember
	if ev.OldMember != nil {
		m = *ev.OldMember
	}
	m.GuildID = ev.GuildID
	m.User = &ev.User
	m.UserID = ev.User.ID
	m.Nick = ev.Nick
	m.AvatarHash = ev.AvatarHash
	m.PremiumSince = ev.PremiumSince
	m.CommunicationDisabledUntil = ev.CommunicationDisabledUntil
	roles := make([]string, len(ev.Roles))
	for i, r := range ev.Roles {
		roles[i] = string(r)
	}
	m.Roles = roles
	if ev.JoinedAt != nil {
		m.JoinedAt = *ev.JoinedAt
	}
	if ev.Deaf != nil {
		m.Deaf = *ev.Deaf
	}
	if ev.Mute != nil {
		m.Mute = *ev.Mute
	}
	if ev.Pending != nil {
		m.Pending = *ev.Pending
	}
	if ev.Flags != nil {
		m.Flags = *ev.Flags
	}
	return &m
}

func nickDiffers(a, b *string) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil || b == nil {
		return true
	}
	return *a != *b
}

// derivePresenceSyntheticEvents derives 0–4 synthetic events from a
// PresenceUpdateEvent. Status transitions require OldPresence (cache hit);
// UserProfileUpdate fires whenever ev.User carries any non-ID field.
func derivePresenceSyntheticEvents(ev *events.PresenceUpdateEvent) []events.Event {
	newPresence := &discord.Presence{
		User:         ev.User,
		GuildID:      ev.GuildID,
		Status:       ev.Status,
		Activities:   ev.Activities,
		ClientStatus: ev.ClientStatus,
	}

	var result []events.Event

	// Status transitions — only meaningful when we have an old baseline.
	if ev.OldPresence != nil {
		wasOffline := ev.OldPresence.Status == discord.PresenceStatusOffline
		isOffline := ev.Status == discord.PresenceStatusOffline

		if wasOffline && !isOffline {
			result = append(result, &events.UserOnlineEvent{
				GuildID:     ev.GuildID,
				UserID:      ev.User.ID,
				Status:      ev.Status,
				OldPresence: ev.OldPresence,
				NewPresence: newPresence,
			})
		} else if !wasOffline && isOffline {
			result = append(result, &events.UserOfflineEvent{
				GuildID:     ev.GuildID,
				UserID:      ev.User.ID,
				OldPresence: ev.OldPresence,
				NewPresence: newPresence,
			})
		}

		if activitiesChanged(ev.OldPresence.Activities, ev.Activities) {
			result = append(result, &events.UserActivityChangeEvent{
				GuildID:       ev.GuildID,
				UserID:        ev.User.ID,
				OldActivities: ev.OldPresence.Activities,
				NewActivities: ev.Activities,
				OldPresence:   ev.OldPresence,
				NewPresence:   newPresence,
			})
		}
	}

	// Profile update — Discord only includes changed fields in the partial user.
	if hasProfileFields(ev.User) {
		result = append(result, &events.UserProfileUpdateEvent{
			GuildID:     ev.GuildID,
			UserID:      ev.User.ID,
			User:        ev.User,
			OldPresence: ev.OldPresence,
			NewPresence: newPresence,
		})
	}

	return result
}

// activitiesChanged reports whether two activity slices differ, compared as
// multisets of (type, name) pairs. State/details changes within the same
// activity are not detected; use raw PRESENCE_UPDATE for that granularity.
func activitiesChanged(old, new_ []discord.FullActivity) bool {
	if len(old) != len(new_) {
		return true
	}
	type key struct {
		t discord.ActivityType
		n string
	}
	counts := make(map[key]int, len(old))
	for _, a := range old {
		counts[key{a.Type, a.Name}]++
	}
	for _, a := range new_ {
		k := key{a.Type, a.Name}
		if counts[k] <= 0 {
			return true
		}
		counts[k]--
	}
	return false
}

// hasProfileFields reports whether the partial presence user carries any
// field beyond the mandatory ID, indicating a profile change.
func hasProfileFields(u discord.PartialPresenceUser) bool {
	return u.Username != nil || u.Discriminator != nil ||
		u.GlobalName != nil || u.AvatarHash != nil ||
		u.Bot != nil || u.PublicFlags != nil
}

// deriveVoiceSyntheticEvents derives 0 or 1 voice synthetic events from a
// VoiceStateUpdateEvent. The four transitions are mutually exclusive.
func deriveVoiceSyntheticEvents(ev *events.VoiceStateUpdateEvent) []events.Event {
	if ev.GuildID == nil {
		return nil
	}
	guildID := *ev.GuildID

	switch {
	case ev.OldState == nil && ev.ChannelID != nil:
		stateCopy := ev.VoiceState
		return []events.Event{&events.VoiceMemberJoinEvent{
			GuildID:   guildID,
			UserID:    ev.UserID,
			ChannelID: *ev.ChannelID,
			State:     &stateCopy,
		}}

	case ev.OldState != nil && ev.ChannelID == nil:
		if ev.OldState.ChannelID == nil {
			return nil
		}
		return []events.Event{&events.VoiceMemberLeaveEvent{
			GuildID:      guildID,
			UserID:       ev.UserID,
			OldChannelID: *ev.OldState.ChannelID,
			OldState:     ev.OldState,
		}}

	case ev.OldState != nil && ev.ChannelID != nil:
		if ev.OldState.ChannelID == nil {
			return nil
		}
		stateCopy := ev.VoiceState
		if *ev.OldState.ChannelID != *ev.ChannelID {
			return []events.Event{&events.VoiceMemberMoveEvent{
				GuildID:      guildID,
				UserID:       ev.UserID,
				OldChannelID: *ev.OldState.ChannelID,
				NewChannelID: *ev.ChannelID,
				OldState:     ev.OldState,
				NewState:     &stateCopy,
			}}
		}
		return []events.Event{&events.VoiceMemberUpdateEvent{
			GuildID:   guildID,
			UserID:    ev.UserID,
			ChannelID: *ev.ChannelID,
			OldState:  ev.OldState,
			NewState:  &stateCopy,
		}}
	}
	return nil
}
