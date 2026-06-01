package discord

import "testing"

func TestEscapeMarkdownInline(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"bold", "**hi**", `\*\*hi\*\*`},
		{"italic underscore", "_hi_", `\_hi\_`},
		{"underline", "__hi__", `\_\_hi\_\_`},
		{"strikethrough", "~~hi~~", `\~\~hi\~\~`},
		{"spoiler", "||hi||", `\|\|hi\|\|`},
		{"inline code", "`code`", "\\`code\\`"},
		{"code block", "```js```", "\\`\\`\\`js\\`\\`\\`"},
		{"backslash", `a\b`, `a\\b`},
		{"masked link", "[x](y)", `\[x](y)`},
		{"plain text untouched", "hello world", "hello world"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := EscapeMarkdown(c.in); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestEscapeMarkdownLineStart(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"heading", "# Title", `\# Title`},
		{"heading h3", "### Title", `\### Title`},
		{"quote", "> quoted", `\> quoted`},
		{"bullet dash", "- item", `\- item`},
		{"bullet plus", "+ item", `\+ item`},
		{"numbered dot", "1. item", `1\. item`},
		{"numbered paren", "2) item", `2\) item`},
		{"indented heading", "  # Title", `  \# Title`},
		{"not a heading mid-line", "a # b", "a # b"},
		{"hash without space", "#nospace", "#nospace"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := EscapeMarkdown(c.in); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestEscapeMarkdownMultiline(t *testing.T) {
	in := "# Title\nsome *text*\n- bullet"
	want := "\\# Title\nsome \\*text\\*\n\\- bullet"
	if got := EscapeMarkdown(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEscapeMarkdownSelective(t *testing.T) {
	// Only escape spoilers; leave bold/italic intact.
	opts := EscapeMarkdownOptions{Spoiler: true}
	if got := EscapeMarkdown("**bold** ||secret||", opts); got != `**bold** \|\|secret\|\|` {
		t.Errorf("selective escape: got %q", got)
	}

	// Zero options escape nothing.
	if got := EscapeMarkdown("**bold**", EscapeMarkdownOptions{}); got != "**bold**" {
		t.Errorf("zero options should escape nothing: got %q", got)
	}
}
