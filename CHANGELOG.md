# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

#### Synthetic events (Sprint E)

21 high-level events derived by the wrapper from raw Discord gateway events.
Synthetic events are enqueued **after** their parent raw event is processed by
the wrapper. In concurrent dispatch configurations, synthetic handlers may run
before the parent raw event handler completes; strict serial ordering requires
single-worker/serial dispatch.
Most require the cache to be enabled; when prior state is unavailable the event
is silently skipped.

**Voice** (no cache required — state tracked internally):
- `VoiceMemberJoinEvent` / `OnVoiceMemberJoin` — user enters a voice channel
- `VoiceMemberLeaveEvent` / `OnVoiceMemberLeave` — user leaves all voice channels
- `VoiceMemberMoveEvent` / `OnVoiceMemberMove` — user switches channels
- `VoiceMemberUpdateEvent` / `OnVoiceMemberUpdate` — state change within the same channel

**Member** (cache required; source: `GUILD_MEMBER_UPDATE`):
- `GuildMemberRoleAddEvent` / `OnGuildMemberRoleAdd` — one event per added role
- `GuildMemberRoleRemoveEvent` / `OnGuildMemberRoleRemove` — one event per removed role
- `GuildMemberNickChangeEvent` / `OnGuildMemberNickChange` — nickname set, changed, or cleared
- `GuildMemberTimeoutEvent` / `OnGuildMemberTimeout` — timeout applied or extended
- `GuildMemberBoostStartEvent` / `OnGuildMemberBoostStart` — member started boosting
- `GuildMemberBoostEndEvent` / `OnGuildMemberBoostEnd` — member stopped boosting

**Presence** (status/activity require cache; source: `PRESENCE_UPDATE`):
- `UserOnlineEvent` / `OnUserOnline` — offline → active status transition
- `UserOfflineEvent` / `OnUserOffline` — active → offline transition
- `UserActivityChangeEvent` / `OnUserActivityChange` — activity set changed
- `UserProfileUpdateEvent` / `OnUserProfileUpdate` — profile fields changed (username, avatar, …)

**Emoji** (cache required; source: `GUILD_EMOJIS_UPDATE`):
- `GuildEmojiAddEvent` / `OnGuildEmojiAdd`
- `GuildEmojiRemoveEvent` / `OnGuildEmojiRemove`
- `GuildEmojiUpdateEvent` / `OnGuildEmojiUpdate` — name changed

**Sticker** (cache required; source: `GUILD_STICKERS_UPDATE`):
- `GuildStickerAddEvent` / `OnGuildStickerAdd`
- `GuildStickerRemoveEvent` / `OnGuildStickerRemove`
- `GuildStickerUpdateEvent` / `OnGuildStickerUpdate` — name changed

**Role** (cache required; source: `GUILD_ROLE_UPDATE`):
- `GuildRolePermissionsChangeEvent` / `OnGuildRolePermissionsChange` — permissions bitfield changed

#### Old-state fields on Update events

All 13 gateway Update events now carry an `OldXxx` field populated from the
cache before the new state is written, enabling before/after comparisons:

| Event | Field |
|-------|-------|
| `GuildUpdateEvent` | `OldGuild *discord.Guild` |
| `ChannelUpdateEvent` | `OldChannel *discord.Channel` |
| `ThreadUpdateEvent` | `OldThread *discord.Channel` |
| `MessageUpdateEvent` | `OldMessage *discord.Message` |
| `GuildMemberUpdateEvent` | `OldMember *discord.GuildMember` |
| `PresenceUpdateEvent` | `OldPresence *discord.Presence` |
| `GuildRoleUpdateEvent` | `OldRole *discord.Role` |
| `GuildScheduledEventUpdateEvent` | `OldEvent *discord.GuildScheduledEvent` |
| `StageInstanceUpdateEvent` | `OldInstance *discord.StageInstance` |
| `AutoModerationRuleUpdateEvent` | `OldRule *discord.AutoModerationRule` |
| `GuildEmojisUpdateEvent` | `OldEmojis []*discord.Emoji` |
| `GuildStickersUpdateEvent` | `OldStickers []*discord.Sticker` |
| `GuildSoundboardSoundUpdateEvent` | `OldSound *discord.SoundboardSound` |

Fields are `nil` when the cache is disabled or was cold for that entity.
