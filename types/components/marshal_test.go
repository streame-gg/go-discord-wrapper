package components

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// marshalType marshals v and returns its "type" discriminator.
func marshalType(t *testing.T, v any) discord.ComponentType {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal %T: %v", v, err)
	}
	var meta struct {
		Type discord.ComponentType `json:"type"`
	}
	if err := json.Unmarshal(b, &meta); err != nil {
		t.Fatalf("unmarshal meta for %T: %v", v, err)
	}
	return meta.Type
}

// TestComponentMarshalType verifies every component's MarshalJSON injects the
// correct ComponentType discriminator, regardless of a zero-valued Type field.
func (su *marshalSuite) TestComponentMarshalType() {
	t := su.T()
	cases := []struct {
		v    any
		want discord.ComponentType
	}{
		{&ActionRow{}, discord.ComponentTypeActionRow},
		{&Button{}, discord.ComponentTypeButton},
		{&ChannelSelectMenu{}, discord.ComponentTypeChannelSelect},
		{&Checkbox{}, discord.ComponentTypeCheckbox},
		{&CheckboxGroup{}, discord.ComponentTypeCheckboxGroup},
		{&Container{}, discord.ComponentTypeContainer},
		{&File{}, discord.ComponentTypeFile},
		{&FileUpload{}, discord.ComponentTypeFileUpload},
		{&Label{}, discord.ComponentTypeLabel},
		{&MediaGallery{}, discord.ComponentTypeMediaGallery},
		{&MentionableSelectMenu{}, discord.ComponentTypeMentionableSelect},
		{&RadioGroup{}, discord.ComponentTypeRadioGroup},
		{&RoleSelectMenu{}, discord.ComponentTypeRoleSelect},
		{&Section{}, discord.ComponentTypeSection},
		{&Separator{}, discord.ComponentTypeSeparator},
		{&StringSelectMenu{}, discord.ComponentTypeStringSelect},
		{&TextDisplay{}, discord.ComponentTypeTextDisplay},
		{&TextInput{}, discord.ComponentTypeTextInput},
		{&Thumbnail{}, discord.ComponentTypeThumbnail},
	}
	for _, tc := range cases {
		if got := marshalType(t, tc.v); got != tc.want {
			t.Errorf("%T: type = %d, want %d", tc.v, got, tc.want)
		}
	}
}

// TestInteractionResponseMarshalType verifies the same for the
// *InteractionResponse component variants.
func (su *marshalSuite) TestInteractionResponseMarshalType() {
	t := su.T()
	cases := []struct {
		v    any
		want discord.ComponentType
	}{
		{&ChannelComponentInteractionResponse{}, discord.ComponentTypeChannelSelect},
		{&CheckboxComponentInteractionResponse{}, discord.ComponentTypeCheckbox},
		{&CheckboxGroupComponentInteractionResponse{}, discord.ComponentTypeCheckboxGroup},
		{&FileUploadComponentInteractionResponse{}, discord.ComponentTypeFileUpload},
		{&LabelComponentInteractionResponse{}, discord.ComponentTypeLabel},
		{&MentionableComponentInteractionResponse{}, discord.ComponentTypeMentionableSelect},
		{&RadioGroupComponentInteractionResponse{}, discord.ComponentTypeRadioGroup},
		{&RoleComponentInteractionResponse{}, discord.ComponentTypeRoleSelect},
		{&StringSelectComponentInteractionResponse{}, discord.ComponentTypeStringSelect},
		{&TextDisplayComponentInteractionResponse{}, discord.ComponentTypeTextDisplay},
		{&TextInputComponentInteractionResponse{}, discord.ComponentTypeTextInput},
		{&UserSelectComponentInteractionResponse{}, discord.ComponentTypeUserSelect},
	}
	for _, tc := range cases {
		if got := marshalType(t, tc.v); got != tc.want {
			t.Errorf("%T: type = %d, want %d", tc.v, got, tc.want)
		}
	}
}

// ── Label (interaction-response) unmarshal dispatch ───────────────────────────

// TestLabelComponentInteractionResponse_DispatchAllChildTypes covers each branch
// of the child-component switch in LabelComponentInteractionResponse.UnmarshalJSON.
func (su *marshalSuite) TestLabelComponentInteractionResponse_DispatchAllChildTypes() {
	t := su.T()
	childTypes := []discord.ComponentType{
		discord.ComponentTypeTextInput,
		discord.ComponentTypeStringSelect,
		discord.ComponentTypeUserSelect,
		discord.ComponentTypeRoleSelect,
		discord.ComponentTypeMentionableSelect,
		discord.ComponentTypeChannelSelect,
		discord.ComponentTypeFileUpload,
		discord.ComponentTypeRadioGroup,
		discord.ComponentTypeCheckboxGroup,
		discord.ComponentTypeCheckbox,
	}
	for _, ct := range childTypes {
		raw := []byte(`{"type":18,"component":{"type":` + itoa(int(ct)) + `}}`)
		var l LabelComponentInteractionResponse
		if err := json.Unmarshal(raw, &l); err != nil {
			t.Fatalf("child type %d: %v", ct, err)
		}
		if l.Component == nil {
			t.Fatalf("child type %d: component not decoded", ct)
		}
	}
}

func (su *marshalSuite) TestLabelComponentInteractionResponse_NilComponent() {
	t := su.T()
	var l LabelComponentInteractionResponse
	if err := json.Unmarshal([]byte(`{"type":18}`), &l); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if l.Component != nil {
		t.Fatal("expected nil component")
	}
}

func (su *marshalSuite) TestLabelComponentInteractionResponse_UnknownChildType() {
	t := su.T()
	var l LabelComponentInteractionResponse
	err := json.Unmarshal([]byte(`{"type":18,"component":{"type":999}}`), &l)
	if err == nil {
		t.Fatal("expected error for unknown child type")
	}
}

func (su *marshalSuite) TestLabelComponentInteractionResponse_InvalidJSON() {
	t := su.T()
	var l LabelComponentInteractionResponse
	if err := l.UnmarshalJSON([]byte(`{`)); err == nil {
		t.Fatal("expected error for malformed json")
	}
}

// ── ComponentLabelComponent unmarshal dispatch ────────────────────────────────

func (su *marshalSuite) TestComponentLabelComponent_DispatchAndErrors() {
	t := su.T()
	// Valid child.
	var l ComponentLabelComponent
	if err := json.Unmarshal([]byte(`{"type":18,"component":{"type":4}}`), &l); err != nil {
		t.Fatalf("valid child: %v", err)
	}
	if l.Component == nil {
		t.Fatal("expected decoded child")
	}

	// Nil component branch.
	var empty ComponentLabelComponent
	if err := json.Unmarshal([]byte(`{"type":18}`), &empty); err != nil {
		t.Fatalf("nil component: %v", err)
	}
	if empty.Component != nil {
		t.Fatal("expected nil component")
	}

	// Unknown child type.
	var unknown ComponentLabelComponent
	if err := json.Unmarshal([]byte(`{"type":18,"component":{"type":999}}`), &unknown); err == nil {
		t.Fatal("expected error for unknown child type")
	}

	// Malformed top-level json.
	if err := (&ComponentLabelComponent{}).UnmarshalJSON([]byte(`{`)); err == nil {
		t.Fatal("expected error for malformed json")
	}
}

// itoa is a tiny dependency-free int->string helper for building JSON payloads.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

type marshalSuite struct{ suite.Suite }

func TestMarshalSuite(t *testing.T) { suite.Run(t, new(marshalSuite)) }
