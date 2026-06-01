package events

import (
	"encoding/json"
	"testing"
)

// TestGuildMemberUpdateMarshalRoundTrip verifies that MarshalJSON preserves the
// deaf, mute, pending, and flags fields. Previously MarshalJSON omitted them
// from its wire struct, so Unmarshal -> Marshal silently dropped the values
// even though UnmarshalJSON read them.
func TestGuildMemberUpdateMarshalRoundTrip(t *testing.T) {
	payload := `{
		"guild_id": "123456789012345678",
		"user": {"id": "987654321098765432", "username": "tester"},
		"roles": ["111111111111111111"],
		"nick": "nickname",
		"deaf": true,
		"mute": true,
		"pending": true,
		"flags": 8
	}`

	var ev GuildMemberUpdateEvent
	if err := json.Unmarshal([]byte(payload), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	out, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var round GuildMemberUpdateEvent
	if err := json.Unmarshal(out, &round); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}

	if round.NewMember.Deaf == nil || !*round.NewMember.Deaf {
		t.Errorf("want Deaf=true after round-trip, got %v", round.NewMember.Deaf)
	}
	if round.NewMember.Mute == nil || !*round.NewMember.Mute {
		t.Errorf("want Mute=true after round-trip, got %v", round.NewMember.Mute)
	}
	if !round.NewMember.Pending {
		t.Error("want Pending=true after round-trip, got false")
	}
	if round.NewMember.Flags != 8 {
		t.Errorf("want Flags=8 after round-trip, got %d", round.NewMember.Flags)
	}
}
