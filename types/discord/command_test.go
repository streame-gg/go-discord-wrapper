package discord

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

// optionTypeCases lists every concrete option type alongside the discriminator
// it must report and marshal.
func optionTypeCases() []struct {
	name string
	opt  AnyApplicationCommandOption
	typ  ApplicationCommandOptionType
} {
	return []struct {
		name string
		opt  AnyApplicationCommandOption
		typ  ApplicationCommandOptionType
	}{
		{"string", &ApplicationCommandOptionString{Name: "s", Description: "d"}, ApplicationCommandOptionTypeString},
		{"integer", &ApplicationCommandOptionInteger{Name: "i", Description: "d"}, ApplicationCommandOptionTypeInteger},
		{"number", &ApplicationCommandOptionNumber{Name: "n", Description: "d"}, ApplicationCommandOptionTypeNumber},
		{"boolean", &ApplicationCommandOptionBoolean{Name: "b", Description: "d"}, ApplicationCommandOptionTypeBoolean},
		{"user", &ApplicationCommandOptionUser{Name: "u", Description: "d"}, ApplicationCommandOptionTypeUser},
		{"channel", &ApplicationCommandOptionChannel{Name: "c", Description: "d"}, ApplicationCommandOptionTypeChannel},
		{"role", &ApplicationCommandOptionRole{Name: "r", Description: "d"}, ApplicationCommandOptionTypeRole},
		{"mentionable", &ApplicationCommandOptionMentionable{Name: "m", Description: "d"}, ApplicationCommandOptionTypeMentionable},
		{"attachment", &ApplicationCommandOptionAttachment{Name: "a", Description: "d"}, ApplicationCommandOptionTypeAttachment},
		{"subcommand", &ApplicationCommandOptionSubCommand{Name: "sc", Description: "d"}, ApplicationCommandOptionTypeSubCommand},
		{"subcommandgroup", &ApplicationCommandOptionSubCommandGroup{Name: "scg", Description: "d"}, ApplicationCommandOptionTypeSubCommandGroup},
	}
}

// TestOptionType_Reported verifies ApplicationCommandOptionType() for each type.
func (su *commandSuite) TestOptionType_Reported() {
	t := su.T()
	for _, tc := range optionTypeCases() {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.opt.ApplicationCommandOptionType(); got != tc.typ {
				t.Fatalf("type = %d, want %d", got, tc.typ)
			}
		})
	}
}

// TestOptionType_MarshalInjectsType verifies MarshalJSON writes the correct
// "type" discriminator regardless of the zero-valued Type field.
func (su *commandSuite) TestOptionType_MarshalInjectsType() {
	t := su.T()
	for _, tc := range optionTypeCases() {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.opt)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var meta struct {
				Type ApplicationCommandOptionType `json:"type"`
			}
			if err := json.Unmarshal(b, &meta); err != nil {
				t.Fatalf("unmarshal meta: %v", err)
			}
			if meta.Type != tc.typ {
				t.Fatalf("marshalled type = %d, want %d (json: %s)", meta.Type, tc.typ, b)
			}
		})
	}
}

// TestCommand_RoundTripAllOptionTypes marshals a command containing one of each
// option type and unmarshals it back, exercising the dispatch switch in
// unmarshalApplicationCommandOption and unmarshalOptionSlice for every case.
func (su *commandSuite) TestCommand_RoundTripAllOptionTypes() {
	t := su.T()
	cases := optionTypeCases()
	opts := make([]AnyApplicationCommandOption, 0, len(cases))
	for _, tc := range cases {
		opts = append(opts, tc.opt)
	}

	cmd := ApplicationCommand{
		Name:        "test",
		Description: "a test command",
		Options:     opts,
	}

	data, err := json.Marshal(&cmd)
	if err != nil {
		t.Fatalf("marshal command: %v", err)
	}

	var got ApplicationCommand
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal command: %v", err)
	}

	if got.Options == nil {
		t.Fatal("expected options to be present")
	}
	if len(got.Options) != len(cases) {
		t.Fatalf("got %d options, want %d", len(got.Options), len(cases))
	}
	for i, tc := range cases {
		if (got.Options)[i].ApplicationCommandOptionType() != tc.typ {
			t.Errorf("option %d (%s): type = %d, want %d", i, tc.name,
				(got.Options)[i].ApplicationCommandOptionType(), tc.typ)
		}
	}
}

// TestCommand_NoOptions verifies a command without options unmarshals with a nil
// Options slice (the raw.Options != nil branch is skipped).
func (su *commandSuite) TestCommand_NoOptions() {
	t := su.T()
	var got ApplicationCommand
	if err := json.Unmarshal([]byte(`{"name":"ping","description":"pong"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Options != nil {
		t.Fatalf("expected nil Options, got %v", got.Options)
	}
}

// TestSubCommand_NestedOptions verifies the nested-option unmarshal path on
// subcommands (and the nil-options branch when omitted).
func (su *commandSuite) TestSubCommand_NestedOptions() {
	t := su.T()
	nested := []AnyApplicationCommandOption{
		&ApplicationCommandOptionString{Name: "arg", Description: "an argument"},
	}
	sc := ApplicationCommandOptionSubCommand{Name: "sub", Description: "d", Options: nested}

	data, err := json.Marshal(&sc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got ApplicationCommandOptionSubCommand
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Options == nil || len(got.Options) != 1 {
		t.Fatalf("expected 1 nested option, got %v", got.Options)
	}
	if (got.Options)[0].ApplicationCommandOptionType() != ApplicationCommandOptionTypeString {
		t.Fatal("nested option should be a string option")
	}

	// Without nested options the Options branch is skipped.
	var empty ApplicationCommandOptionSubCommand
	if err := json.Unmarshal([]byte(`{"name":"sub","description":"d"}`), &empty); err != nil {
		t.Fatalf("unmarshal empty: %v", err)
	}
	if empty.Options != nil {
		t.Fatalf("expected nil Options, got %v", empty.Options)
	}
}

// TestUnmarshal_UnknownOptionType exercises the default branch of the dispatch
// switch.
func (su *commandSuite) TestUnmarshal_UnknownOptionType() {
	t := su.T()
	_, err := unmarshalApplicationCommandOption([]byte(`{"type":999,"name":"x"}`))
	if err == nil {
		t.Fatal("expected error for unknown option type")
	}
}

// TestUnmarshal_InvalidJSON exercises the error path when the discriminator
// cannot be read.
func (su *commandSuite) TestUnmarshal_InvalidJSON() {
	t := su.T()
	_, err := unmarshalApplicationCommandOption([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
}

// TestCommand_UnmarshalInvalidNestedOption propagates errors from the option
// slice up through ApplicationCommand.UnmarshalJSON.
func (su *commandSuite) TestCommand_UnmarshalInvalidNestedOption() {
	t := su.T()
	var got ApplicationCommand
	err := json.Unmarshal([]byte(`{"name":"x","description":"d","options":[{"type":999}]}`), &got)
	if err == nil {
		t.Fatal("expected error from invalid nested option type")
	}
}

// TestOption_UnmarshalInvalidJSON exercises the error-return branch inside each
// concrete option's UnmarshalJSON by feeding it a malformed payload.
func (su *commandSuite) TestOption_UnmarshalInvalidJSON() {
	t := su.T()
	bad := []byte(`{`)
	options := []AnyApplicationCommandOption{
		&ApplicationCommandOptionString{},
		&ApplicationCommandOptionInteger{},
		&ApplicationCommandOptionNumber{},
		&ApplicationCommandOptionBoolean{},
		&ApplicationCommandOptionUser{},
		&ApplicationCommandOptionChannel{},
		&ApplicationCommandOptionRole{},
		&ApplicationCommandOptionMentionable{},
		&ApplicationCommandOptionAttachment{},
		&ApplicationCommandOptionSubCommand{},
		&ApplicationCommandOptionSubCommandGroup{},
	}
	for _, o := range options {
		if err := o.UnmarshalJSON(bad); err == nil {
			t.Errorf("%T.UnmarshalJSON: expected error on malformed json", o)
		}
	}
}

// TestSubCommand_UnmarshalInvalidNestedOption verifies the nested-slice error
// propagates out of the subcommand/group unmarshallers.
func (su *commandSuite) TestSubCommand_UnmarshalInvalidNestedOption() {
	t := su.T()
	payload := []byte(`{"name":"s","description":"d","options":[{"type":999}]}`)

	var sc ApplicationCommandOptionSubCommand
	if err := sc.UnmarshalJSON(payload); err == nil {
		t.Error("subcommand: expected error from invalid nested option")
	}

	var scg ApplicationCommandOptionSubCommandGroup
	if err := scg.UnmarshalJSON(payload); err == nil {
		t.Error("subcommandgroup: expected error from invalid nested option")
	}
}

// TestCommand_MarshalRoundTripScalars verifies top-level scalar fields survive a
// marshal round-trip.
func (su *commandSuite) TestCommand_MarshalRoundTripScalars() {
	t := su.T()
	guild := Snowflake(42)
	cmd := ApplicationCommand{
		ID:            Snowflake(1),
		Type:          ApplicationCommandType(1),
		ApplicationID: Snowflake(2),
		GuildID:       guild,
		Name:          "name",
		Description:   "desc",
	}
	data, err := json.Marshal(&cmd)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got ApplicationCommand
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != "name" || got.Description != "desc" || got.GuildID != 42 {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
}

type commandSuite struct{ suite.Suite }

func TestCommandSuite(t *testing.T) { suite.Run(t, new(commandSuite)) }
