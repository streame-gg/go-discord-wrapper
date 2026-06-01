package util_test

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
)

func (su *pointersSuite) TestPointerOf() {
	t := su.T()
	p := util.PointerOf(42)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != 42 {
		t.Fatalf("expected 42, got %d", *p)
	}

	s := util.PointerOf("hello")
	if *s != "hello" {
		t.Fatalf("expected hello, got %q", *s)
	}
}

type pointersSuite struct{ suite.Suite }

func TestPointersSuite(t *testing.T) { suite.Run(t, new(pointersSuite)) }
