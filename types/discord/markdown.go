package discord

import "strings"

// EscapeMarkdownOptions controls which Markdown constructs EscapeMarkdown
// neutralises. The zero value escapes nothing; use DefaultEscapeMarkdownOptions
// (the default when no options are passed to EscapeMarkdown) to escape all
// supported constructs.
//
// Note that some toggles share a control character: '*' is escaped when Bold or
// Italic is set, and '_' is escaped when Italic or Underline is set.
// https://docs.discord.com/developers/reference#message-formatting
type EscapeMarkdownOptions struct {
	Escape        bool // backslash itself
	Bold          bool // **
	Italic        bool // * and _
	Underline     bool // __
	Strikethrough bool // ~~
	Spoiler       bool // ||
	InlineCode    bool // `
	CodeBlock     bool // ```
	Quote         bool // leading > / >>>
	Heading       bool // leading #, ##, ###
	BulletedList  bool // leading -, *, +
	NumberedList  bool // leading "1." / "1)"
	MaskedLink    bool // [
}

// DefaultEscapeMarkdownOptions returns an options value with every supported
// construct enabled.
func DefaultEscapeMarkdownOptions() EscapeMarkdownOptions {
	return EscapeMarkdownOptions{
		Escape:        true,
		Bold:          true,
		Italic:        true,
		Underline:     true,
		Strikethrough: true,
		Spoiler:       true,
		InlineCode:    true,
		CodeBlock:     true,
		Quote:         true,
		Heading:       true,
		BulletedList:  true,
		NumberedList:  true,
		MaskedLink:    true,
	}
}

// EscapeMarkdown returns text with Discord Markdown control sequences escaped so
// the string renders literally. With no opts every supported construct is
// escaped; pass a single EscapeMarkdownOptions to escape only specific ones.
func EscapeMarkdown(text string, opts ...EscapeMarkdownOptions) string {
	o := DefaultEscapeMarkdownOptions()
	if len(opts) > 0 {
		o = opts[0]
	}

	runes := []rune(text)
	var b strings.Builder
	b.Grow(len(text) + len(text)/8)

	atLineStart := true
	for i := 0; i < len(runes); {
		if atLineStart {
			if consumed := escapeLineStart(&b, runes, i, o); consumed > 0 {
				i += consumed
				atLineStart = false
				continue
			}
			atLineStart = false
		}

		r := runes[i]
		switch r {
		case '\n':
			atLineStart = true
		case '\\':
			if o.Escape {
				b.WriteByte('\\')
			}
		case '`':
			if o.InlineCode || o.CodeBlock {
				b.WriteByte('\\')
			}
		case '*':
			if o.Bold || o.Italic {
				b.WriteByte('\\')
			}
		case '_':
			if o.Italic || o.Underline {
				b.WriteByte('\\')
			}
		case '~':
			if o.Strikethrough {
				b.WriteByte('\\')
			}
		case '|':
			if o.Spoiler {
				b.WriteByte('\\')
			}
		case '[':
			if o.MaskedLink {
				b.WriteByte('\\')
			}
		}
		b.WriteRune(r)
		i++
	}

	return b.String()
}

// escapeLineStart handles Markdown constructs that are only significant at the
// start of a line (after up to three leading spaces): headings, block quotes,
// bulleted lists and numbered lists. It writes the (escaped) leading marker to b
// and returns the number of runes it consumed, or 0 if no construct matched.
func escapeLineStart(b *strings.Builder, runes []rune, start int, o EscapeMarkdownOptions) int {
	// Skip up to three leading spaces, as Markdown does.
	j := start
	for spaces := 0; j < len(runes) && runes[j] == ' ' && spaces < 3; spaces++ {
		j++
	}

	writeLead := func() {
		for x := start; x < j; x++ {
			b.WriteRune(runes[x])
		}
	}

	// Heading: 1–3 '#' followed by a space. Escaping the first '#' is enough.
	if o.Heading {
		k, h := j, 0
		for k < len(runes) && runes[k] == '#' && h < 3 {
			k++
			h++
		}
		if h > 0 && k < len(runes) && runes[k] == ' ' {
			writeLead()
			b.WriteByte('\\')
			b.WriteRune('#')
			return (j - start) + 1
		}
	}

	// Block quote: '>' optionally as '>>>'.
	if o.Quote && j < len(runes) && runes[j] == '>' {
		writeLead()
		b.WriteByte('\\')
		b.WriteRune('>')
		return (j - start) + 1
	}

	// Bulleted list: '-', '*' or '+' followed by a space.
	if o.BulletedList && j < len(runes) && (runes[j] == '-' || runes[j] == '*' || runes[j] == '+') &&
		j+1 < len(runes) && runes[j+1] == ' ' {
		writeLead()
		b.WriteByte('\\')
		b.WriteRune(runes[j])
		return (j - start) + 1
	}

	// Numbered list: one or more digits then '.' or ')' followed by a space.
	if o.NumberedList {
		k, d := j, 0
		for k < len(runes) && runes[k] >= '0' && runes[k] <= '9' {
			k++
			d++
		}
		if d > 0 && k < len(runes) && (runes[k] == '.' || runes[k] == ')') &&
			k+1 < len(runes) && runes[k+1] == ' ' {
			writeLead()
			for x := j; x < k; x++ { // digits
				b.WriteRune(runes[x])
			}
			b.WriteByte('\\')
			b.WriteRune(runes[k])
			return (k - start) + 1
		}
	}

	return 0
}
