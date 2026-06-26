// internal/testutil/event_suite.go

package testutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

// UnmarshalTestCase beschreibt einen einzelnen Testfall generisch
type UnmarshalTestCase[T any] struct {
	Name     string
	Input    []byte
	WantErr  bool
	Validate func(t *testing.T, got T)
}

// EventSuite ist die generische Basis für alle Event-Suites
type EventSuite[T any] struct {
	suite.Suite
}

func InitSub[T any](parent interface{ T() *testing.T }) *EventSuite[T] {
	sub := new(EventSuite[T])
	sub.SetT(parent.T())
	return sub
}

func (s *EventSuite[T]) MustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	s.Require().NoError(err, "MustMarshal failed")
	return data
}

func (s *EventSuite[T]) Unmarshal(data []byte) (T, error) {
	var e T
	err := json.Unmarshal(data, &e)
	return e, err
}

func (s *EventSuite[T]) AssertUnmarshalSucceeds(data []byte) T {
	s.T().Helper()
	e, err := s.Unmarshal(data)
	s.Require().NoError(err)
	return e
}

func (s *EventSuite[T]) AssertUnmarshalFails(data []byte) {
	s.T().Helper()
	_, err := s.Unmarshal(data)
	s.Require().Error(err)
}

func (s *EventSuite[T]) AssertRoundtrip(e T) {
	s.T().Helper()
	data, err := json.Marshal(e)
	s.Require().NoError(err, "re-marshal failed")
	s.AssertUnmarshalSucceeds(data)
}

// RunCases führt eine Liste von UnmarshalTestCases aus
func (s *EventSuite[T]) RunCases(cases []UnmarshalTestCase[T]) {
	s.T().Helper()
	for _, tc := range cases {
		s.Run(tc.Name, func() {
			e, err := s.Unmarshal(tc.Input)
			if tc.WantErr {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			if tc.Validate != nil {
				tc.Validate(s.T(), e)
			}
		})
	}
}

// RunCommonEdgeCases testet Standard-Edge-Cases für jedes Event
func (s *EventSuite[T]) RunCommonEdgeCases() {
	s.Run("invalid JSON", func() { s.AssertUnmarshalFails([]byte(`{not valid}`)) })
	s.Run("empty bytes", func() { s.AssertUnmarshalFails([]byte(``)) })
	s.Run("null", func() { s.AssertUnmarshalSucceeds([]byte(`null`)) })
	s.Run("empty object", func() { s.AssertUnmarshalSucceeds([]byte(`{}`)) })
	s.Run("unknown fields ignored", func() {
		s.AssertUnmarshalSucceeds([]byte(`{"__unknown_future_field__": true}`))
	})
}
