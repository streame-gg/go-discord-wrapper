package util_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
)

type pointersSuite struct{ suite.Suite }

func TestPointersSuite(t *testing.T) { suite.Run(t, new(pointersSuite)) }

func (s *pointersSuite) TestPointerOf() {
	p := util.PointerOf(42)
	s.Require().NotNil(p)
	s.Equal(42, *p)

	str := util.PointerOf("hello")
	s.Require().NotNil(str)
	s.Equal("hello", *str)
}
