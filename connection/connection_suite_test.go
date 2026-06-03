package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConnectionSuite groups the gateway/dispatch/cache-hydration tests under a
// single testify suite. The individual tests build their own mock gateways, so
// the suite holds no shared state; each test method begins with `t := s.T()` to
// keep the existing per-test helpers (which take *testing.T) working unchanged.
type ConnectionSuite struct {
	suite.Suite
}

func TestConnectionSuite(t *testing.T) { suite.Run(t, new(ConnectionSuite)) }
