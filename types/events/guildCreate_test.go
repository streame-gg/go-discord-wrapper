package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestGuildCreateEventOuterFields verifies that Large, Unavailable, and
// MemberCount on GuildCreateEvent are decoded correctly (Issue 5).
// Prior to the fix the promoted GatewayGuildWrapper.UnmarshalJSON consumed the
// entire payload and these outer fields were always zero/nil.
func (su *guildCreateSuite) TestGuildCreateEventOuterFields() {
	t := su.T()
	payload := `{
		"id": "123456789012345678",
		"name": "Test Guild",
		"large": true,
		"unavailable": false,
		"member_count": 1234,
		"roles": [],
		"emojis": [],
		"features": [],
		"members": [],
		"channels": [],
		"threads": [],
		"presences": [],
		"voice_states": [],
		"stage_instances": [],
		"guild_scheduled_events": []
	}`

	var ev GuildCreateEvent
	if err := json.Unmarshal([]byte(payload), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !ev.Large {
		t.Error("want Large=true, got false")
	}
	if ev.Unavailable != false {
		t.Errorf("want Unavailable=false, got %v", ev.Unavailable)
	}
	if ev.MemberCount != 1234 {
		t.Errorf("want MemberCount=1234, got %d", ev.MemberCount)
	}

	// Also verify the wrapper decoded the guild correctly.
	if ev.Guild == nil {
		t.Fatal("want non-nil Guild after decode")
	}
	if ev.Guild.GetID() != discord.Snowflake(123456789012345678) {
		t.Errorf("want guild ID 123456789012345678, got %v", ev.Guild.GetID())
	}
}

// TestGuildCreateEventUnavailableGuild verifies that an unavailable guild
// payload is still handled correctly.
func (su *guildCreateSuite) TestGuildCreateEventUnavailableGuild() {
	t := su.T()
	payload := `{"id": "111111111111111111", "unavailable": true}`

	var ev GuildCreateEvent
	if err := json.Unmarshal([]byte(payload), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.Guild == nil {
		t.Fatal("want non-nil Guild after decode")
	}
	if ev.Guild.IsAvailable() {
		t.Error("want guild to be unavailable")
	}
}

type guildCreateSuite struct{ suite.Suite }

func TestGuildCreateSuite(t *testing.T) { suite.Run(t, new(guildCreateSuite)) }
