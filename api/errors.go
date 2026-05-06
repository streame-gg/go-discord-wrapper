package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// Sentinel errors for common HTTP status codes. Use errors.Is to check against these.
//
//	if errors.Is(err, api.ErrNotFound) { ... }
var (
	ErrUnauthorized = errors.New("discord: unauthorized")
	ErrForbidden    = errors.New("discord: forbidden")
	ErrNotFound     = errors.New("discord: not found")
	ErrRateLimited  = errors.New("discord: rate limited")
)

// APIError is returned by all REST methods when Discord responds with a non-success status.
// It carries the HTTP status, the Discord JSON error code and message, and a field-level
// error map when present.
//
// Use errors.Is with the sentinel vars (ErrNotFound, ErrForbidden, etc.) for status checks,
// or errors.As to access the full error detail:
//
//	var apiErr *api.APIError
//	if errors.As(err, &apiErr) {
//	    log.Printf("discord code %d: %s", apiErr.Code, apiErr.Message)
//	}
type APIError struct {
	// HTTPStatus is the HTTP response status code (e.g. 403, 404).
	HTTPStatus int

	// Code is the Discord JSON error code. Zero when Discord did not return one.
	Code common.GatewayErrorCode

	// Message is the human-readable error message from Discord.
	Message string

	// Errors contains field-level validation errors, present on some 400 responses.
	Errors map[string]interface{}
}

func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("discord api error %d: %s (http %d)", int(e.Code), e.Message, e.HTTPStatus)
	}
	return fmt.Sprintf("discord api error: http %d", e.HTTPStatus)
}

// Is maps the sentinel errors to HTTP status codes so callers can use errors.Is.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrUnauthorized:
		return e.HTTPStatus == http.StatusUnauthorized
	case ErrForbidden:
		return e.HTTPStatus == http.StatusForbidden
	case ErrNotFound:
		return e.HTTPStatus == http.StatusNotFound
	case ErrRateLimited:
		return e.HTTPStatus == http.StatusTooManyRequests
	}
	return false
}
