package connection

import (
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// cacheAdapter narrows a full cache.Cache (read-write) to the discord.Cache
// read-only interface consumed by entity managers. Each method returns a
// cache.*Store value; those store types structurally satisfy the corresponding
// discord.*Store interface (duck typing — same method set, no explicit embedding
// needed).
type cacheAdapter struct{ inner cache.Cache }

func (a cacheAdapter) Guilds() discord.GuildStore                   { return a.inner.Guilds() }
func (a cacheAdapter) Channels() discord.ChannelStore               { return a.inner.Channels() }
func (a cacheAdapter) Users() discord.UserStore                     { return a.inner.Users() }
func (a cacheAdapter) Members() discord.MemberStore                 { return a.inner.Members() }
func (a cacheAdapter) Messages() discord.MessageStore               { return a.inner.Messages() }
func (a cacheAdapter) Roles() discord.RoleStore                     { return a.inner.Roles() }
func (a cacheAdapter) VoiceStates() discord.VoiceStateStore         { return a.inner.VoiceStates() }
func (a cacheAdapter) Soundboard() discord.SoundboardStore          { return a.inner.Soundboard() }
func (a cacheAdapter) ScheduledEvents() discord.ScheduledEventStore { return a.inner.ScheduledEvents() }
func (a cacheAdapter) StageInstances() discord.StageInstanceStore   { return a.inner.StageInstances() }
func (a cacheAdapter) Emojis() discord.EmojiStore                   { return a.inner.Emojis() }
func (a cacheAdapter) Stickers() discord.StickerStore               { return a.inner.Stickers() }

var _ discord.Cache = cacheAdapter{}
