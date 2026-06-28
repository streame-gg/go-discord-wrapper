package events

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type eventSuite struct {
	suite.Suite
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(eventSuite))
}

func (s *eventSuite) TestAllCommandsRegistered() {
	s.Equal(77, len(eventFactories))
}
