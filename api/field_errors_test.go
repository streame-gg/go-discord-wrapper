package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// fieldErrorsSuite covers FieldError flattening and the CodeOf/HasCode/
// FieldErrorsOf helpers that pull strongly-typed detail out of an *Error.
type fieldErrorsSuite struct {
	suite.Suite
}

func TestFieldErrorsSuite(t *testing.T) {
	suite.Run(t, new(fieldErrorsSuite))
}

// errorBodyFixture is the nested validation-error example from
// https://docs.discord.com/developers/reference#error-messages.
const errorBodyFixture = `{
  "message": "Invalid Form Body",
  "code": 50035,
  "errors": {
    "activities": {
      "0": {
        "platform": {
          "_errors": [
            { "code": "BASE_TYPE_CHOICES", "message": "Value must be one of ('discord', 'xbox')." }
          ]
        },
        "type": {
          "_errors": [
            { "code": "BASE_TYPE_CHOICES", "message": "Value must be one of (0, 1, 2, 3, 4, 5)." }
          ]
        }
      }
    }
  }
}`

func (s *fieldErrorsSuite) parseFixture() *Error {
	s.T().Helper()
	var body discord.GatewayError
	s.Require().NoError(json.Unmarshal([]byte(errorBodyFixture), &body))
	return &Error{
		HTTPStatus: 400,
		Code:       discord.JSONErrorCode(body.Code),
		Message:    body.Message,
		Errors:     body.Errors,
	}
}

func (s *fieldErrorsSuite) TestFieldErrorsFlatten() {
	e := s.parseFixture()
	got := e.FieldErrors()
	s.Require().Len(got, 2)

	// Sorted by path: "activities.0.platform" before "activities.0.type".
	s.Equal("activities.0.platform", got[0].Path)
	s.Equal("BASE_TYPE_CHOICES", got[0].Code)
	s.Equal("Value must be one of ('discord', 'xbox').", got[0].Message)

	s.Equal("activities.0.type", got[1].Path)
	s.Equal("activities.0.type: Value must be one of (0, 1, 2, 3, 4, 5). (BASE_TYPE_CHOICES)", got[1].String())
}

func (s *fieldErrorsSuite) TestFieldErrorsEmpty() {
	s.Nil((&Error{}).FieldErrors())
	s.Nil((*Error)(nil).FieldErrors())
}

func (s *fieldErrorsSuite) TestFieldErrorsTopLevel() {
	// Some responses put _errors at the root with no field path.
	e := &Error{Errors: map[string]interface{}{
		"_errors": []interface{}{
			map[string]interface{}{"code": "APPLICATION_COMMAND_TOO_LARGE", "message": "Command exceeds maximum size."},
		},
	}}
	got := e.FieldErrors()
	s.Require().Len(got, 1)
	s.Equal("", got[0].Path)
	s.Equal("Command exceeds maximum size. (APPLICATION_COMMAND_TOO_LARGE)", got[0].String())
}

// TestExtractCodeFromMethodError mirrors what a REST method like
// SetVoiceChannelStatus returns on a 403, and shows the three ways to get the
// exact, strongly-typed Discord code out of it.
func (s *fieldErrorsSuite) TestExtractCodeFromMethodError() {
	// What doRequestWithoutResponse returns when Discord replies 403 + 50013.
	var err error = &Error{
		HTTPStatus: 403,
		Code:       discord.JSONErrorCodeMissingPermissions,
		Message:    "Missing Permissions",
	}

	// 1. Sentinel via errors.Is.
	s.Require().ErrorIs(err, ErrMissingPermissions)

	// 2. HasCode for any of the 236 codes (no sentinel needed).
	s.True(HasCode(err, discord.JSONErrorCodeMissingPermissions))

	// 3. CodeOf for switch-style handling; the returned value is typed.
	code, ok := CodeOf(err)
	s.Require().True(ok)
	s.Equal(discord.JSONErrorCodeMissingPermissions, code)

	var handled bool
	switch code {
	case discord.JSONErrorCodeMissingPermissions:
		handled = true
	case discord.JSONErrorCodeUnknownChannel:
		s.T().Fatal("wrong branch")
	}
	s.True(handled)
}

// TestFieldErrorsArrayIndices covers Discord's array-error shape: an array
// field becomes an object keyed by stringified indices, and an array element can
// itself carry _errors directly.
func (s *fieldErrorsSuite) TestFieldErrorsArrayIndices() {
	const raw = `{
      "code": 50035,
      "message": "Invalid Form Body",
      "errors": {
        "embeds": {
          "0": { "title": { "_errors": [ {"code":"BASE_TYPE_MAX_LENGTH","message":"Must be 256 or fewer in length."} ] } },
          "2": { "fields": { "1": { "value": { "_errors": [ {"code":"BASE_TYPE_REQUIRED","message":"This field is required"} ] } } } }
        },
        "content": {
          "0": { "_errors": [ {"code":"BASE_TYPE_BAD_LENGTH","message":"Must be between 1 and 2000."} ] }
        }
      }
    }`
	var body discord.GatewayError
	s.Require().NoError(json.Unmarshal([]byte(raw), &body))
	e := &Error{Code: discord.JSONErrorCode(body.Code), Errors: body.Errors}

	got := e.FieldErrors()
	// Build a path->message map so the assertion doesn't depend on sort order.
	byPath := map[string]FieldError{}
	for _, fe := range got {
		byPath[fe.Path] = fe
	}
	s.Require().Len(got, 3)
	s.Equal("BASE_TYPE_MAX_LENGTH", byPath["embeds.0.title"].Code)
	s.Equal("BASE_TYPE_REQUIRED", byPath["embeds.2.fields.1.value"].Code)
	s.Equal("BASE_TYPE_BAD_LENGTH", byPath["content.0"].Code) // _errors directly on an array element
}

func (s *fieldErrorsSuite) TestCodeOfAndHasCode() {
	e := s.parseFixture()
	wrapped := fmt.Errorf("create message failed: %w", error(e))

	code, ok := CodeOf(wrapped)
	s.Require().True(ok)
	s.Equal(discord.JSONErrorCodeInvalidFormBody, code) // 50035
	s.True(HasCode(wrapped, discord.JSONErrorCodeInvalidFormBody))
	s.False(HasCode(wrapped, discord.JSONErrorCodeUnknownMessage))

	// Non-API error.
	_, ok = CodeOf(fmt.Errorf("network down"))
	s.False(ok)

	// FieldErrorsOf unwraps too.
	s.Len(FieldErrorsOf(wrapped), 2)
	s.Nil(FieldErrorsOf(fmt.Errorf("nope")))
}
