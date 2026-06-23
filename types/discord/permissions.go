package discord

import (
	"encoding/json"
	"strconv"
)

// https://docs.discord.com/developers/resources/channel#overwrite-object
type PermissionOverwriteType int

const (
	PermissionOverwriteTypeRole PermissionOverwriteType = 0
	PermissionOverwriteTypeUser PermissionOverwriteType = 1
)

// Permission is a bitmask of Discord permission flags. Discord transmits
// permissions as a decimal string (e.g. "2147483647") to avoid JavaScript
// integer-precision loss on bits ≥ 53; use uint64 so all 64 bits are safe.
//
// https://docs.discord.com/developers/topics/permissions#permissions-bitwise-permission-flags
type Permission uint64

// MarshalJSON encodes the permission bitmask as a decimal string, matching
// Discord's wire format.
func (p Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(p), 10))
}

// UnmarshalJSON accepts either a decimal string (Discord's standard format)
// or a raw integer for backwards compatibility.
func (p *Permission) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*p = Permission(n)
		return nil
	}
	var n uint64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*p = Permission(n)
	return nil
}

// NullFlag represents a Discord role-tag boolean flag that is encoded in JSON
// as the presence of a null field (true) vs. absence of the field (false).
// Use with json:"field_name,omitempty": a zero NullFlag is omitted; a true
// NullFlag marshals to JSON null.
//
// https://docs.discord.com/developers/topics/permissions#role-object-role-tags-structure
type NullFlag bool

func (f *NullFlag) UnmarshalJSON(_ []byte) error {
	*f = true // any presence — null or otherwise — means true
	return nil
}

func (f NullFlag) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

const (
	PermissionCreateInstantInvite              Permission = 1 << 0
	PermissionKickMembers                      Permission = 1 << 1
	PermissionBanMembers                       Permission = 1 << 2
	PermissionAdministrator                    Permission = 1 << 3
	PermissionManageChannels                   Permission = 1 << 4
	PermissionManageGuild                      Permission = 1 << 5
	PermissionAddReactions                     Permission = 1 << 6
	PermissionViewAuditLog                     Permission = 1 << 7
	PermissionPrioritySpeaker                  Permission = 1 << 8
	PermissionStream                           Permission = 1 << 9
	PermissionViewChannel                      Permission = 1 << 10
	PermissionSendMessages                     Permission = 1 << 11
	PermissionSendTTSMessages                  Permission = 1 << 12
	PermissionManageMessages                   Permission = 1 << 13
	PermissionEmbedLinks                       Permission = 1 << 14
	PermissionAttachFiles                      Permission = 1 << 15
	PermissionReadMessageHistory               Permission = 1 << 16
	PermissionMentionEveryone                  Permission = 1 << 17
	PermissionUseExternalEmojis                Permission = 1 << 18
	PermissionViewGuildInsights                Permission = 1 << 19
	PermissionConnect                          Permission = 1 << 20
	PermissionSpeak                            Permission = 1 << 21
	PermissionMuteMembers                      Permission = 1 << 22
	PermissionDeafenMembers                    Permission = 1 << 23
	PermissionMoveMembers                      Permission = 1 << 24
	PermissionUseVAD                           Permission = 1 << 25
	PermissionChangeNickname                   Permission = 1 << 26
	PermissionManageNicknames                  Permission = 1 << 27
	PermissionManageRoles                      Permission = 1 << 28
	PermissionManageWebhooks                   Permission = 1 << 29
	PermissionManageGuildExpressions           Permission = 1 << 30
	PermissionUseApplicationCommands           Permission = 1 << 31
	PermissionRequestToSpeak                   Permission = 1 << 32
	PermissionManageEvents                     Permission = 1 << 33
	PermissionManageThreads                    Permission = 1 << 34
	PermissionCreatePublicThreads              Permission = 1 << 35
	PermissionCreatePrivateThreads             Permission = 1 << 36
	PermissionUseExternalStickers              Permission = 1 << 37
	PermissionSendMessagesInThreads            Permission = 1 << 38
	PermissionUseEmbeddedActivities            Permission = 1 << 39
	PermissionModerateMembers                  Permission = 1 << 40
	PermissionViewCreatorMonetizationAnalytics Permission = 1 << 41
	PermissionUseSoundboard                    Permission = 1 << 42
	PermissionCreateGuildExpressions           Permission = 1 << 43
	PermissionCreateEvents                     Permission = 1 << 44
	PermissionUseExternalSounds                Permission = 1 << 45
	PermissionSendVoiceMessages                Permission = 1 << 46
	PermissionSetVoiceChannelStatus            Permission = 1 << 48
	PermissionSendPolls                        Permission = 1 << 49
	PermissionUseExternalApps                  Permission = 1 << 50
	PermissionPinMessages                      Permission = 1 << 51
	PermissionBypassSlowmode                   Permission = 1 << 52
)
