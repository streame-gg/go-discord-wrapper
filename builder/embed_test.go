package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestEmbedBuilder_BuildNoAlias verifies that AddFields after Build() does
// not mutate the already-returned Embed (P2-33).
func TestEmbedBuilder_BuildNoAlias(t *testing.T) {
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

// TestEmbedBuilder_TwoBuildsAreIndependent checks that two Build() calls on
// the same builder return embeds backed by separate slices.
func TestEmbedBuilder_TwoBuildsAreIndependent(t *testing.T) {
	b := NewEmbed().AddFields(discord.EmbedFields{Name: "x", Value: "y"})
	a := b.Build()
	c := b.Build()
	// Mutate one field in a; c must be unaffected.
	a.Fields[0].Name = "mutated"
	assert.Equal(t, "x", c.Fields[0].Name, "builds must not share backing array")
}
