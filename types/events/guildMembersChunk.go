package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

// GuildMembersChunkEvent is sent in response to a Request Guild Members gateway command.
// Large guilds may send several chunks; check ChunkIndex and ChunkCount to track progress.
type GuildMembersChunkEvent struct {
	GuildID    common.Snowflake     `json:"guild_id"`
	Members    []common.GuildMember `json:"members"`
	ChunkIndex int                  `json:"chunk_index"`
	ChunkCount int                  `json:"chunk_count"`
	NotFound   []common.Snowflake   `json:"not_found,omitempty"`
	Nonce      *string              `json:"nonce,omitempty"`
}

func init() { RegisterEvent(GuildMembersChunkEvent{}) }

func (e GuildMembersChunkEvent) DesiredEventType() Event { return &GuildMembersChunkEvent{} }
func (e GuildMembersChunkEvent) Event() EventType        { return EventGuildMembersChunk }
