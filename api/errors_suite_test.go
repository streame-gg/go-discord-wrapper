package api

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// errorsSuite covers the error types' string formatting and the errors.Is
// mapping that lets callers match both HTTP-status sentinels and Discord
// JSON-code sentinels against a single *Error.
type errorsSuite struct {
	suite.Suite
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(errorsSuite))
}

func (s *errorsSuite) TestDiscordCodeErrorString() {
	// ErrUnknownChannel is a *discordCodeError; its Error() embeds the numeric code.
	var ce *discordCodeError
	s.Require().ErrorAs(ErrUnknownChannel, &ce)
	s.Equal(int(discord.JSONErrorCodeUnknownChannel), int(ce.code))
	s.Contains(ce.Error(), "error code")
}

func (s *errorsSuite) TestErrorString() {
	withCode := &Error{HTTPStatus: 403, Code: 50013, Message: "Missing Permissions"}
	s.Equal("discord api error 50013: Missing Permissions (http 403)", withCode.Error())

	noCode := &Error{HTTPStatus: 500}
	s.Equal("discord api error: http 500", noCode.Error())
}

func (s *errorsSuite) TestErrorIsStatusSentinels() {
	s.True(errors.Is(&Error{HTTPStatus: http.StatusUnauthorized}, ErrUnauthorized))
	s.True(errors.Is(&Error{HTTPStatus: http.StatusForbidden}, ErrForbidden))
	s.True(errors.Is(&Error{HTTPStatus: http.StatusNotFound}, ErrNotFound))
	s.True(errors.Is(&Error{HTTPStatus: http.StatusTooManyRequests}, ErrRateLimited)) // 429
	s.False(errors.Is(&Error{HTTPStatus: http.StatusForbidden}, ErrNotFound))
}

func (s *errorsSuite) TestErrorIsTypedCodeSentinels() {
	err := &Error{HTTPStatus: 403, Code: discord.JSONErrorCodeMissingPermissions}
	s.True(errors.Is(err, ErrMissingPermissions))
	s.False(errors.Is(err, ErrUnknownChannel))
}

// TestJSONErrorCarrier guards the wiring that lets the discord package's
// IsErrorCode / ErrorCodeOf helpers match an *api.Error. Without the
// JSONErrorCode() method, these helpers silently return false for every
// real API error.
func (s *errorsSuite) TestJSONErrorCarrier() {
	err := &Error{HTTPStatus: 403, Code: discord.JSONErrorCodeMissingPermissions}

	// Direct: *Error satisfies the carrier interface.
	s.Equal(discord.JSONErrorCodeMissingPermissions, err.JSONErrorCode())

	// discord.IsErrorCode matches by code.
	s.True(discord.IsErrorCode(err, discord.JSONErrorCodeMissingPermissions))
	s.False(discord.IsErrorCode(err, discord.JSONErrorCodeUnknownChannel))

	// And it still matches once the error is wrapped.
	wrapped := fmt.Errorf("doing thing: %w", err)
	s.True(discord.IsErrorCode(wrapped, discord.JSONErrorCodeMissingPermissions))

	code, ok := discord.ErrorCodeOf(wrapped)
	s.True(ok)
	s.Equal(discord.JSONErrorCodeMissingPermissions, code)

	// Non-API errors carry no code.
	_, ok = discord.ErrorCodeOf(errors.New("plain"))
	s.False(ok)
	s.False(discord.IsErrorCode(nil, discord.JSONErrorCodeMissingPermissions))
}
