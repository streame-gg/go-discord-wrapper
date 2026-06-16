package discord

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestIntegerOptionChoicesMarshalAsInt verifies that IntegerOption choices
// serialize as JSON numbers, not quoted strings (Issue 8).
func (su *optionTypesSuite) TestIntegerOptionChoicesMarshalAsInt() {
	t := su.T()
	opt := ApplicationCommandOptionInteger{
		Name:        "level",
		Description: "pick a level",
		Choices: []ApplicationCommandOptionChoice[int64]{
			{Name: "low", Value: 1},
			{Name: "high", Value: 100},
		},
	}

	b, err := json.Marshal(&opt)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := string(b)

	if strings.Contains(out, `"value":"1"`) || strings.Contains(out, `"value":"100"`) {
		t.Errorf("choices must not be quoted strings; got: %s", out)
	}
	if !strings.Contains(out, `"value":1`) || !strings.Contains(out, `"value":100`) {
		t.Errorf("choices must be JSON integers; got: %s", out)
	}
}

// TestNumberOptionChoicesMarshalAsFloat verifies that NumberOption choices
// serialize as JSON floats and that MinValue/MaxValue are float64 (Issue 8).
func (su *optionTypesSuite) TestNumberOptionChoicesMarshalAsFloat() {
	t := su.T()
	minVal := float64(0.5)
	maxVal := float64(99.9)

	opt := ApplicationCommandOptionNumber{
		Name:        "amount",
		Description: "pick an amount",
		Choices: []ApplicationCommandOptionChoice[float64]{
			{Name: "half", Value: 0.5},
			{Name: "max", Value: 99.9},
		},
		MinValue: &minVal,
		MaxValue: &maxVal,
	}

	b, err := json.Marshal(&opt)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := string(b)

	if strings.Contains(out, `"value":"0.5"`) || strings.Contains(out, `"value":"99.9"`) {
		t.Errorf("choices must not be quoted strings; got: %s", out)
	}
	if !strings.Contains(out, `0.5`) || !strings.Contains(out, `99.9`) {
		t.Errorf("choices/min/max must be JSON numbers; got: %s", out)
	}
}

// TestStringOptionChoicesStillWork verifies that string choices are unaffected.
func (su *optionTypesSuite) TestStringOptionChoicesStillWork() {
	t := su.T()
	opt := ApplicationCommandOptionString{
		Name:        "lang",
		Description: "pick a language",
		Choices: []ApplicationCommandOptionChoice[string]{
			{Name: "Go", Value: "go"},
			{Name: "Rust", Value: "rust"},
		},
	}

	b, err := json.Marshal(&opt)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := string(b)

	if !strings.Contains(out, `"value":"go"`) {
		t.Errorf("string choices must be quoted; got: %s", out)
	}
}

type optionTypesSuite struct{ suite.Suite }

func TestOptionTypesSuite(t *testing.T) { suite.Run(t, new(optionTypesSuite)) }
