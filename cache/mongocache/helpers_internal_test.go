package mongocache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Internal-package tests for helpers that are unreachable from the external
// _test package (extractID's default case and upsertByID's unknown-type guard).

type mongoHelpersSuite struct{ suite.Suite }

func TestMongoHelpersSuite(t *testing.T) { suite.Run(t, new(mongoHelpersSuite)) }

func (s *mongoHelpersSuite) TestExtractID_AllDocTypes() {
	t := s.T()
	id, ok := extractID(entityDoc{ID: "a"})
	assert.True(t, ok)
	assert.Equal(t, "a", id)

	_, ok = extractID(memberDoc{ID: "b"})
	assert.True(t, ok)
	_, ok = extractID(roleDoc{ID: "c"})
	assert.True(t, ok)
	_, ok = extractID(guildEntityDoc{ID: "d"})
	assert.True(t, ok)
	_, ok = extractID(guildUserDoc{ID: "e"})
	assert.True(t, ok)
	_, ok = extractID(msgDoc{ID: "f"})
	assert.True(t, ok)

	// Unknown type → default case.
	_, ok = extractID(struct{ X int }{X: 1})
	assert.False(t, ok)
}

func (s *mongoHelpersSuite) TestUpsertByID_UnknownTypeReturnsError() {
	t := s.T()
	// nil collection is never touched because extractID fails first.
	err := upsertByID(context.Background(), nil, struct{ X int }{X: 1})
	assert.Error(t, err)
}
