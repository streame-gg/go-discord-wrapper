package api

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (su *apiErrorsSuite) TestAPIErrorIs() {
	t := su.T()
	tests := []struct {
		name       string
		httpStatus int
		target     error
		wantMatch  bool
	}{
		{"404 matches ErrNotFound", http.StatusNotFound, ErrNotFound, true},
		{"403 matches ErrForbidden", http.StatusForbidden, ErrForbidden, true},
		{"401 matches ErrUnauthorized", http.StatusUnauthorized, ErrUnauthorized, true},
		{"429 matches ErrRateLimited", http.StatusTooManyRequests, ErrRateLimited, true},
		{"500 matches no sentinel", http.StatusInternalServerError, ErrNotFound, false},
		{"500 not ErrForbidden", http.StatusInternalServerError, ErrForbidden, false},
		{"404 not ErrForbidden", http.StatusNotFound, ErrForbidden, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := &Error{HTTPStatus: tc.httpStatus}
			got := errors.Is(err, tc.target)
			assert.Equal(t, tc.wantMatch, got)
		})
	}
}

func (su *apiErrorsSuite) TestAPIErrorAs() {
	t := su.T()
	original := &Error{
		HTTPStatus: http.StatusNotFound,
		Code:       discord.GatewayErrorCode(10003),
		Message:    "Unknown Channel",
	}

	wrapped := fmt.Errorf("wrapping: %w", original)

	var apiErr *Error
	require.True(t, errors.As(wrapped, &apiErr), "errors.As should unwrap *APIError")
	assert.Equal(t, http.StatusNotFound, apiErr.HTTPStatus)
	assert.Equal(t, discord.GatewayErrorCode(10003), apiErr.Code)
	assert.Equal(t, "Unknown Channel", apiErr.Message)
}

func (su *apiErrorsSuite) TestAPIErrorString() {
	t := su.T()
	t.Run("with discord code", func(t *testing.T) {
		err := &Error{
			HTTPStatus: http.StatusNotFound,
			Code:       discord.GatewayErrorCode(10003),
			Message:    "Unknown Channel",
		}
		s := err.Error()
		assert.True(t, strings.Contains(s, "10003"), "error string should contain discord code 10003, got: %s", s)
		assert.True(t, strings.Contains(s, "404"), "error string should contain http status 404, got: %s", s)
	})

	t.Run("without discord code", func(t *testing.T) {
		err := &Error{
			HTTPStatus: http.StatusInternalServerError,
			Code:       0,
		}
		s := err.Error()
		assert.True(t, strings.Contains(s, "500"), "error string should contain http status 500, got: %s", s)
	})
}

type apiErrorsSuite struct{ suite.Suite }

func TestApiErrorsSuite(t *testing.T) { suite.Run(t, new(apiErrorsSuite)) }
