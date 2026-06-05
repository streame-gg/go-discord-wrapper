package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// errorBodyFixture is the nested validation-error example from
// https://discord.com/developers/docs/reference#error-messages.
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

func parseFixture(t *testing.T) *Error {
	t.Helper()
	var body discord.GatewayError
	require.NoError(t, json.Unmarshal([]byte(errorBodyFixture), &body))
	return &Error{
		HTTPStatus: 400,
		Code:       discord.JSONErrorCode(body.Code),
		Message:    body.Message,
		Errors:     body.Errors,
	}
}

func TestFieldErrorsFlatten(t *testing.T) {
	e := parseFixture(t)
	got := e.FieldErrors()
	require.Len(t, got, 2)

	// Sorted by path: "activities.0.platform" before "activities.0.type".
	require.Equal(t, "activities.0.platform", got[0].Path)
	require.Equal(t, "BASE_TYPE_CHOICES", got[0].Code)
	require.Equal(t, "Value must be one of ('discord', 'xbox').", got[0].Message)

	require.Equal(t, "activities.0.type", got[1].Path)
	require.Equal(t, "activities.0.type: Value must be one of (0, 1, 2, 3, 4, 5). (BASE_TYPE_CHOICES)", got[1].String())
}

func TestFieldErrorsEmpty(t *testing.T) {
	require.Nil(t, (&Error{}).FieldErrors())
	require.Nil(t, (*Error)(nil).FieldErrors())
}

func TestFieldErrorsTopLevel(t *testing.T) {
	// Some responses put _errors at the root with no field path.
	e := &Error{Errors: map[string]interface{}{
		"_errors": []interface{}{
			map[string]interface{}{"code": "APPLICATION_COMMAND_TOO_LARGE", "message": "Command exceeds maximum size."},
		},
	}}
	got := e.FieldErrors()
	require.Len(t, got, 1)
	require.Equal(t, "", got[0].Path)
	require.Equal(t, "Command exceeds maximum size. (APPLICATION_COMMAND_TOO_LARGE)", got[0].String())
}

// TestExtractCodeFromMethodError mirrors what a REST method like
// SetVoiceChannelStatus returns on a 403, and shows the three ways to get the
// exact, strongly-typed Discord code out of it.
func TestExtractCodeFromMethodError(t *testing.T) {
	// What doRequestWithoutResponse returns when Discord replies 403 + 50013.
	var err error = &Error{
		HTTPStatus: 403,
		Code:       discord.JSONErrorCodeMissingPermissions,
		Message:    "Missing Permissions",
	}

	// 1. Sentinel via errors.Is.
	require.ErrorIs(t, err, ErrMissingPermissions)

	// 2. HasCode for any of the 236 codes (no sentinel needed).
	require.True(t, HasCode(err, discord.JSONErrorCodeMissingPermissions))

	// 3. CodeOf for switch-style handling; the returned value is typed.
	code, ok := CodeOf(err)
	require.True(t, ok)
	require.Equal(t, discord.JSONErrorCodeMissingPermissions, code)

	var handled bool
	switch code {
	case discord.JSONErrorCodeMissingPermissions:
		handled = true
	case discord.JSONErrorCodeUnknownChannel:
		t.Fatal("wrong branch")
	}
	require.True(t, handled)
}

// TestFieldErrorsArrayIndices covers Discord's array-error shape: an array
// field becomes an object keyed by stringified indices, and an array element can
// itself carry _errors directly.
func TestFieldErrorsArrayIndices(t *testing.T) {
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
	require.NoError(t, json.Unmarshal([]byte(raw), &body))
	e := &Error{Code: discord.JSONErrorCode(body.Code), Errors: body.Errors}

	got := e.FieldErrors()
	// Build a path->message map so the assertion doesn't depend on sort order.
	byPath := map[string]FieldError{}
	for _, fe := range got {
		byPath[fe.Path] = fe
	}
	require.Len(t, got, 3)
	require.Equal(t, "BASE_TYPE_MAX_LENGTH", byPath["embeds.0.title"].Code)
	require.Equal(t, "BASE_TYPE_REQUIRED", byPath["embeds.2.fields.1.value"].Code)
	require.Equal(t, "BASE_TYPE_BAD_LENGTH", byPath["content.0"].Code) // _errors directly on an array element
}

func TestCodeOfAndHasCode(t *testing.T) {
	e := parseFixture(t)
	wrapped := fmt.Errorf("create message failed: %w", error(e))

	code, ok := CodeOf(wrapped)
	require.True(t, ok)
	require.Equal(t, discord.JSONErrorCodeInvalidFormBody, code) // 50035
	require.True(t, HasCode(wrapped, discord.JSONErrorCodeInvalidFormBody))
	require.False(t, HasCode(wrapped, discord.JSONErrorCodeUnknownMessage))

	// Non-API error.
	_, ok = CodeOf(fmt.Errorf("network down"))
	require.False(t, ok)

	// FieldErrorsOf unwraps too.
	require.Len(t, FieldErrorsOf(wrapped), 2)
	require.Nil(t, FieldErrorsOf(fmt.Errorf("nope")))
}
