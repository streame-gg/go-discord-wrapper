package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func labelComp(label string) *components.LabelComponent {
	return &components.LabelComponent{
		Type:  discord.ComponentTypeActionRow,
		Label: label,
	}
}

// TestModalBuilder_BuildNoAlias verifies that the returned *Modal is
// independent of the builder: further AddComponents calls must not
// mutate the already-built Modal (P2-32).
func TestModalBuilder_BuildNoAlias(t *testing.T) {
	b := NewModal().SetCustomID("m").SetTitle("T").AddComponents(labelComp("first"))

	first := b.Build()
	require.NotNil(t, first.Components)
	assert.Len(t, *first.Components, 1)

	// Add a second component and build again.
	b.AddComponents(labelComp("second"))
	second := b.Build()

	// first must still have exactly one component.
	assert.Len(t, *first.Components, 1, "first modal must not be mutated after second Build()")
	assert.Len(t, *second.Components, 2)
}

// TestModalBuilder_TwoBuildsAreIndependent verifies that two consecutive
// Build() calls return independent objects (not the same pointer).
func TestModalBuilder_TwoBuildsAreIndependent(t *testing.T) {
	b := NewModal().SetCustomID("m").SetTitle("T")
	a := b.Build()
	c := b.Build()
	assert.NotSame(t, a, c)
	assert.NotSame(t, a.Components, c.Components)
}
