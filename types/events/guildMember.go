package events

import (
	"encoding/json"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildMemberAddEvent{}
var _ Event = &GuildMemberUpdateEvent{}
var _ Event = &GuildMemberRemoveEvent{}
var _ Event = &GuildMembersChunkEvent{}

func init() {
	RegisterEvent(&GuildMemberAddEvent{})
	RegisterEvent(&GuildMemberRemoveEvent{})
	RegisterEvent(&GuildMemberUpdateEvent{})
	RegisterEvent(&GuildMembersChunkEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#guild-member-add
type GuildMemberAddEvent struct {
	GuildID     discord.Snowflake   `json:"guild_id"`
	GuildMember discord.GuildMember `json:"-"`

	Guild *discord.Guild `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-member-remove
type GuildMemberRemoveEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`

	Guild  *discord.Guild       `json:"-"`
	Member *discord.GuildMember `json:"-"`
}

// GuildMembersChunkEvent is sent in response to a Request Guild Members gateway command.
// Large guilds may send several chunks; check ChunkIndex and ChunkCount to track progress.
// https://docs.discord.com/developers/events/gateway-events#guild-members-chunk
type GuildMembersChunkEvent struct {
	GuildID    discord.Snowflake     `json:"guild_id"`
	Members    []discord.GuildMember `json:"members"`
	ChunkIndex int                   `json:"chunk_index"`
	ChunkCount int                   `json:"chunk_count"`
	NotFound   []interface{}         `json:"not_found,omitempty"`
	Nonce      string                `json:"nonce,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-member-update
type GuildMemberUpdateEvent struct {
	GuildID   discord.Snowflake    `json:"-"`
	Guild     *discord.Guild       `json:"-"`
	NewMember discord.GuildMember  `json:"-"`
	OldMember *discord.GuildMember `json:"-"`

	Avatar                     *string                       `json:"avatar"`
	Banner                     *string                       `json:"banner"`
	CommunicationDisabledUntil *time.Time                    `json:"communication_disabled_until,omitempty"`
	Deaf                       *bool                         `json:"deaf,omitempty"`
	JoinedAt                   *time.Time                    `json:"joined_at"`
	Mute                       *bool                         `json:"mute,omitempty"`
	Nick                       *string                       `json:"nick,omitempty"`
	Pending                    *bool                         `json:"pending,omitempty"`
	PremiumSince               *time.Time                    `json:"premium_since"`
	Roles                      []discord.Snowflake           `json:"roles"`
	User                       discord.User                  `json:"user"`
	AvatarDecorationData       *discord.AvatarDecorationData `json:"avatar_decoration_data,omitempty"`
	Collectibles               *discord.Collectible          `json:"collectibles,omitempty"`
}

func (e *GuildMemberAddEvent) DesiredEventType() Event { return &GuildMemberAddEvent{} }
func (e *GuildMemberAddEvent) Event() EventType        { return EventGuildMemberAdd }
func (e *GuildMemberAddEvent) UnmarshalJSON(data []byte) (err error) {
	var structData struct {
		GuildID discord.Snowflake `json:"guild_id"`
		discord.GuildMember
	}

	if err := json.Unmarshal(data, &structData); err != nil {
		return err
	}

	e.GuildID = structData.GuildID
	e.GuildMember = structData.GuildMember

	return nil
}

func (e *GuildMemberUpdateEvent) DesiredEventType() Event { return &GuildMemberUpdateEvent{} }
func (e *GuildMemberUpdateEvent) Event() EventType        { return EventGuildMemberUpdate }
func (e *GuildMemberUpdateEvent) UnmarshalJSON(data []byte) error {
	var structData struct {
		GuildID                    discord.Snowflake             `json:"guild_id"`
		Avatar                     *string                       `json:"avatar"`
		Banner                     *string                       `json:"banner"`
		CommunicationDisabledUntil *time.Time                    `json:"communication_disabled_until,omitempty"`
		Deaf                       *bool                         `json:"deaf,omitempty"`
		JoinedAt                   *time.Time                    `json:"joined_at"`
		Mute                       *bool                         `json:"mute,omitempty"`
		Nick                       *string                       `json:"nick,omitempty"`
		Pending                    *bool                         `json:"pending,omitempty"`
		PremiumSince               *time.Time                    `json:"premium_since"`
		Roles                      []discord.Snowflake           `json:"roles"`
		User                       discord.User                  `json:"user"`
		AvatarDecorationData       *discord.AvatarDecorationData `json:"avatar_decoration_data,omitempty"`
		Collectibles               *discord.Collectible          `json:"collectibles,omitempty"`
	}
	if err := json.Unmarshal(data, &data); err != nil {
		return err
	}

	e.GuildID = structData.GuildID
	e.NewMember.GuildID = structData.GuildID
	e.NewMember.User = &structData.User
	e.NewMember.Nick = structData.Nick
	e.NewMember.Avatar = structData.Avatar
	e.NewMember.PremiumSince = structData.PremiumSince
	e.CommunicationDisabledUntil = structData.CommunicationDisabledUntil
	e.Nick = structData.Nick
	e.NewMember.Roles = structData.Roles

	if structData.JoinedAt != nil {
		e.NewMember.JoinedAt = structData.JoinedAt
	}
	if structData.Deaf != nil {
		e.NewMember.Deaf = *structData.Deaf
	}
	if structData.Mute != nil {
		e.NewMember.Mute = *structData.Mute
	}
	if structData.Pending != nil {
		e.NewMember.Pending = *structData.Pending
	}
	if structData.CommunicationDisabledUntil != nil {
		e.NewMember.CommunicationDisabledUntil = structData.CommunicationDisabledUntil
	}

	return nil
}

func (e *GuildMemberRemoveEvent) DesiredEventType() Event { return &GuildMemberRemoveEvent{} }
func (e *GuildMemberRemoveEvent) Event() EventType        { return EventGuildMemberRemove }
func (e *GuildMemberRemoveEvent) UnmarshalJSON(data []byte) (err error) {
	return json.Unmarshal(data, &e)
}

func (e *GuildMembersChunkEvent) DesiredEventType() Event { return &GuildMembersChunkEvent{} }
func (e *GuildMembersChunkEvent) Event() EventType        { return EventGuildMembersChunk }
