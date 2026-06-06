package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GuildMembersChunkEvent is sent in response to a Request Guild Members gateway command.
// Large guilds may send several chunks; check ChunkIndex and ChunkCount to track progress.
// https://docs.discord.com/developers/events/gateway-events#guild-members-chunk
type GuildMembersChunkEvent struct {
	GuildID    discord.Snowflake     `json:"guild_id"`
	Members    []discord.GuildMember `json:"members"`
	ChunkIndex int                   `json:"chunk_index"`
	ChunkCount int                   `json:"chunk_count"`
	NotFound   []discord.Snowflake   `json:"not_found,omitempty"`
	Nonce      *string               `json:"nonce,omitempty"`
}

func init() { RegisterEvent(GuildMembersChunkEvent{}) }

func (e GuildMembersChunkEvent) DesiredEventType() Event { return &GuildMembersChunkEvent{} }
func (e GuildMembersChunkEvent) Event() EventType        { return EventGuildMembersChunk }
