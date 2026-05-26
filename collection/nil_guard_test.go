package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/collection"
)

// TestCollectionNilGuards verifies that methods accepting a *Collection
// argument do not panic when passed nil (Issue 28).
func TestCollectionNilGuards(t *testing.T) {
	eq := func(a, b string) bool { return a == b }

	populated := collection.New[int, string]()
	populated.Set(1, "a")

	empty := collection.New[int, string]()

	var nilColl *collection.Collection[int, string]

	t.Run("Equals_nil_vs_populated", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Equals(nilColl, eq)
			assert.False(t, got, "non-empty collection must not equal nil")
		})
	})

	t.Run("Equals_nil_vs_empty", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := empty.Equals(nilColl, eq)
			assert.True(t, got, "empty collection must equal nil (both empty)")
		})
	})

	t.Run("Concat_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Concat(nilColl)
			assert.Equal(t, 1, got.Len(), "Concat with nil must return clone of receiver")
		})
	})

	t.Run("Difference_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Difference(nilColl)
			assert.Equal(t, 1, got.Len(), "Difference with nil must return clone of receiver")
		})
	})

	t.Run("Intersection_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.Intersection(nilColl)
			assert.Equal(t, 0, got.Len(), "Intersection with nil must be empty")
		})
	})

	t.Run("SymmetricDifference_nil", func(t *testing.T) {
		require.NotPanics(t, func() {
			got := populated.SymmetricDifference(nilColl)
			assert.Equal(t, 1, got.Len(), "SymmetricDifference with nil must return clone of receiver")
		})
	})
}
