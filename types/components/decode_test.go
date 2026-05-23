package components

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestContainerDecodeRoundTrip verifies that Container with nested components
// decodes without panic and preserves all child components (Issue 1).
func TestContainerDecodeRoundTrip(t *testing.T) {
	raw := `{
		"type": 17,
		"components": [
			{"type": 10, "content": "hello"},
			{"type": 14, "spacing": 1},
			{"type": 1, "components": []}
		]
	}`

	var c Container
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(c.Components) != 3 {
		t.Fatalf("want 3 components, got %d", len(c.Components))
	}
	if c.Components[0].GetType() != discord.ComponentTypeTextDisplay {
		t.Errorf("component[0] want TextDisplay, got %d", c.Components[0].GetType())
	}
	if c.Components[1].GetType() != discord.ComponentTypeSeparator {
		t.Errorf("component[1] want Separator, got %d", c.Components[1].GetType())
	}
	if c.Components[2].GetType() != discord.ComponentTypeActionRow {
		t.Errorf("component[2] want ActionRow, got %d", c.Components[2].GetType())
	}
}

// TestContainerUnknownTypeError verifies that an unknown component type inside
// a Container returns an error whose message contains the decimal type number (Issue 4b).
func TestContainerUnknownTypeError(t *testing.T) {
	raw := `{"type": 17, "components": [{"type": 999}]}`
	var c Container
	err := json.Unmarshal([]byte(raw), &c)
	if err == nil {
		t.Fatal("want error for unknown type, got nil")
	}
	if !strings.Contains(err.Error(), "999") {
		t.Errorf("error message should contain '999', got: %s", err.Error())
	}
}

// TestSectionDecodeRoundTrip verifies that Section with a TextDisplay child
// decodes without panic (Issue 2).
func TestSectionDecodeRoundTrip(t *testing.T) {
	raw := `{
		"type": 9,
		"components": [
			{"type": 10, "content": "section text"}
		]
	}`

	var s Section
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(s.Components) != 1 {
		t.Fatalf("want 1 component, got %d", len(s.Components))
	}
}

// TestSectionUnknownTypeError verifies that an unknown type inside a Section
// returns a readable error (Issue 4a + 4b).
func TestSectionUnknownTypeError(t *testing.T) {
	raw := `{"type": 9, "components": [{"type": 42}]}`
	var s Section
	err := json.Unmarshal([]byte(raw), &s)
	if err == nil {
		t.Fatal("want error for unknown type, got nil")
	}
	if !strings.Contains(err.Error(), "42") {
		t.Errorf("error message should contain '42', got: %s", err.Error())
	}
}

// TestActionRowDecodeAllTypes verifies that every valid ActionRow child type
// decodes without being silently dropped (Issue 3).
func TestActionRowDecodeAllTypes(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		wantType discord.ComponentType
	}{
		{"Button", `{"type": 2, "style": 1, "label": "click"}`, discord.ComponentTypeButton},
		{"StringSelect", `{"type": 3, "custom_id": "s"}`, discord.ComponentTypeStringSelect},
		{"TextInput", `{"type": 4, "custom_id": "t", "label": "input", "style": 1}`, discord.ComponentTypeTextInput},
		{"UserSelect", `{"type": 5, "custom_id": "u"}`, discord.ComponentTypeUserSelect},
		{"RoleSelect", `{"type": 6, "custom_id": "r"}`, discord.ComponentTypeRoleSelect},
		{"MentionableSelect", `{"type": 7, "custom_id": "m"}`, discord.ComponentTypeMentionableSelect},
		{"ChannelSelect", `{"type": 8, "custom_id": "c"}`, discord.ComponentTypeChannelSelect},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := `{"type": 1, "components": [` + tc.raw + `]}`
			var a ActionRow
			if err := json.Unmarshal([]byte(raw), &a); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if len(a.Components) != 1 {
				t.Fatalf("want 1 component, got %d (component was silently dropped)", len(a.Components))
			}
			if a.Components[0].GetType() != tc.wantType {
				t.Errorf("want type %d, got %d", tc.wantType, a.Components[0].GetType())
			}
		})
	}
}

// TestActionRowUnknownTypeError verifies that an unknown type returns an error
// with a decimal type number in the message (Issue 3 + 4b).
func TestActionRowUnknownTypeError(t *testing.T) {
	raw := `{"type": 1, "components": [{"type": 999}]}`
	var a ActionRow
	err := json.Unmarshal([]byte(raw), &a)
	if err == nil {
		t.Fatal("want error for unknown type, got nil")
	}
	if !strings.Contains(err.Error(), "999") {
		t.Errorf("error message should contain '999', got: %s", err.Error())
	}
}

// TestLabelUnknownTypeError verifies LabelComponent returns a readable decimal
// error on unknown nested component type (Issue 4b).
func TestLabelUnknownTypeError(t *testing.T) {
	raw := `{"type": 18, "label": "l", "component": {"type": 999}}`
	var l LabelComponent
	err := json.Unmarshal([]byte(raw), &l)
	if err == nil {
		t.Fatal("want error for unknown type, got nil")
	}
	if !strings.Contains(err.Error(), "999") {
		t.Errorf("error message should contain '999', got: %s", err.Error())
	}
}
