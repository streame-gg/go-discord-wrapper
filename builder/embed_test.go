package builder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEmbedBuilder_BuildNoAlias verifies that AddFields after Build() does
// not mutate the already-returned Embed (P2-33).
func (su *embedSuite) TestEmbedBuilder_BuildNoAlias() {
	t := su.T()
	b := NewEmbed().SetTitle("T").AddFields(discord.EmbedFields{Name: "a", Value: "1"})

	first := b.Build()
	assert.Len(t, first.Fields, 1)

	// Ensure the underlying array has no spare capacity by measuring after
	// a second AddFields call that would use it.
	b.AddFields(discord.EmbedFields{Name: "b", Value: "2"})
	second := b.Build()

	// first must still have exactly one field.
	assert.Len(t, first.Fields, 1, "first embed must not be mutated after second Build()")
	assert.Len(t, second.Fields, 2)
}

// ── Bug 134: Validate / BuildChecked ─────────────────────────────────────────

func (su *embedSuite) TestBug134_Validate_OkOnValid() {
	t := su.T()
	b := NewEmbed().
		SetTitle("Hello").
		SetDescription("World").
		SetFooter("foot", "").
		SetAuthor("au", "", "").
		AddFields(discord.EmbedFields{Name: "n", Value: "v"})
	assert.NoError(t, b.Validate())
}

func (su *embedSuite) TestBug134_Validate_TitleTooLong() {
	t := su.T()
	b := NewEmbed().SetTitle(strings.Repeat("x", 257))
	assert.ErrorContains(t, b.Validate(), "title")
}

func (su *embedSuite) TestBug134_Validate_DescriptionTooLong() {
	t := su.T()
	b := NewEmbed().SetDescription(strings.Repeat("d", 4097))
	assert.ErrorContains(t, b.Validate(), "description")
}

func (su *embedSuite) TestBug134_Validate_TooManyFields() {
	t := su.T()
	b := NewEmbed()
	for i := 0; i < 26; i++ {
		b.AddFields(discord.EmbedFields{Name: "n", Value: "v"})
	}
	assert.ErrorContains(t, b.Validate(), "fields")
}

func (su *embedSuite) TestBug134_Validate_FieldNameTooLong() {
	t := su.T()
	b := NewEmbed().AddFields(discord.EmbedFields{Name: strings.Repeat("n", 257), Value: "v"})
	assert.ErrorContains(t, b.Validate(), "field 0 name")
}

func (su *embedSuite) TestBug134_Validate_FieldValueTooLong() {
	t := su.T()
	b := NewEmbed().AddFields(discord.EmbedFields{Name: "n", Value: strings.Repeat("v", 1025)})
	assert.ErrorContains(t, b.Validate(), "field 0 value")
}

func (su *embedSuite) TestBug134_Validate_FooterTooLong() {
	t := su.T()
	b := NewEmbed().SetFooter(strings.Repeat("f", 2049), "")
	assert.ErrorContains(t, b.Validate(), "footer")
}

func (su *embedSuite) TestBug134_Validate_AuthorTooLong() {
	t := su.T()
	b := NewEmbed().SetAuthor(strings.Repeat("a", 257), "", "")
	assert.ErrorContains(t, b.Validate(), "author")
}

func (su *embedSuite) TestBug134_Validate_TotalExceeds6000() {
	t := su.T()
	// Two fields each contributing 3000 chars → total 6000 → still OK.
	b := NewEmbed().
		AddFields(discord.EmbedFields{Name: strings.Repeat("a", 256), Value: strings.Repeat("b", 1024)}).
		AddFields(discord.EmbedFields{Name: strings.Repeat("c", 256), Value: strings.Repeat("d", 1024)}).
		SetTitle(strings.Repeat("t", 256)).
		SetDescription(strings.Repeat("d", 3184)) // exactly at 6000
	assert.NoError(t, b.Validate(), "exactly 6000 characters must be valid")

	// Push over the limit by one.
	b.SetDescription(strings.Repeat("d", 3185))
	assert.ErrorContains(t, b.Validate(), "6000")
}

func (su *embedSuite) TestBug134_BuildChecked_ReturnsEmbedOnSuccess() {
	t := su.T()
	b := NewEmbed().SetTitle("ok")
	embed, err := b.BuildChecked()
	require.NoError(t, err)
	require.NotNil(t, embed.Title)
	assert.Equal(t, "ok", embed.Title)
}

func (su *embedSuite) TestBug134_BuildChecked_ReturnsErrorOnViolation() {
	t := su.T()
	b := NewEmbed().SetTitle(strings.Repeat("x", 257))
	_, err := b.BuildChecked()
	assert.Error(t, err)
}

func (su *embedSuite) TestBug134_Build_StillWorksWithoutValidation() {
	t := su.T()
	// Build() must not call Validate() — it must always return the embed as-is.
	b := NewEmbed().SetTitle(strings.Repeat("x", 257))
	embed := b.Build()
	require.NotNil(t, embed.Title)
	assert.Equal(t, 257, len([]rune(embed.Title)), "Build() must not truncate (Bug 134)")
}

// TestEmbedBuilder_TwoBuildsAreIndependent checks that two Build() calls on
// the same builder return embeds backed by separate slices.
func (su *embedSuite) TestEmbedBuilder_TwoBuildsAreIndependent() {
	t := su.T()
	b := NewEmbed().AddFields(discord.EmbedFields{Name: "x", Value: "y"})
	a := b.Build()
	c := b.Build()
	// Mutate one field in a; c must be unaffected.
	a.Fields[0].Name = "mutated"
	assert.Equal(t, "x", c.Fields[0].Name, "builds must not share backing array")
}

type embedSuite struct{ suite.Suite }

func TestEmbedSuite(t *testing.T) { suite.Run(t, new(embedSuite)) }
