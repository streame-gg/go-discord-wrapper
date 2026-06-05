package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// Sentinel errors for discord HTTP status codes. Use errors.Is to check against these.
//
//	if errors.Is(err, api.ErrNotFound) { ... }
var (
	ErrUnauthorized = errors.New("discord: unauthorized")
	ErrForbidden    = errors.New("discord: forbidden")
	ErrNotFound     = errors.New("discord: not found")
	ErrRateLimited  = errors.New("discord: rate limited")
)

// discordCodeError is a sentinel value for checking a specific Discord JSON error code.
type discordCodeError struct {
	code discord.JSONErrorCode
}

func (e *discordCodeError) Error() string {
	return fmt.Sprintf("discord: error code %d", int(e.code))
}

// Typed sentinels for the most commonly checked Discord JSON error codes.
// Use errors.Is to test:
//
//	if errors.Is(err, api.ErrMissingPermissions) { ... }
var (
	ErrUnknownChannel                 = &discordCodeError{discord.JSONErrorCodeUnknownChannel}
	ErrUnknownGuild                   = &discordCodeError{discord.JSONErrorCodeUnknownGuild}
	ErrUnknownMessage                 = &discordCodeError{discord.JSONErrorCodeUnknownMessage}
	ErrUnknownMember                  = &discordCodeError{discord.JSONErrorCodeUnknownMember}
	ErrUnknownRole                    = &discordCodeError{discord.JSONErrorCodeUnknownRole}
	ErrUnknownWebhook                 = &discordCodeError{discord.JSONErrorCodeUnknownWebhook}
	ErrUnknownUser                    = &discordCodeError{discord.JSONErrorCodeUnknownUser}
	ErrUnknownEmoji                   = &discordCodeError{discord.JSONErrorCodeUnknownEmoji}
	ErrUnknownInteraction             = &discordCodeError{discord.JSONErrorCodeUnknownInteraction}
	ErrMissingAccess                  = &discordCodeError{discord.JSONErrorCodeMissingAccess}
	ErrMissingPermissions             = &discordCodeError{discord.JSONErrorCodeMissingPermissions}
	ErrCannotSendMessagesToUser       = &discordCodeError{discord.JSONErrorCodeCannotSendMessagesToThisUser}
	ErrInteractionAlreadyAcknowledged = &discordCodeError{discord.JSONErrorCodeInteractionAlreadyAcknowledged}
	ErrThreadIsLocked                 = &discordCodeError{discord.JSONErrorCodeThreadIsLocked}
	ErrMessageTooOldForBulkDelete     = &discordCodeError{discord.JSONErrorCodeMessageTooOldForBulkDelete}
	ErrInvalidFormBody                = &discordCodeError{discord.JSONErrorCodeInvalidFormBody}
	ErrMaxReactionsReached            = &discordCodeError{discord.JSONErrorCodeMaxReactions}
	ErrCannotExecuteOnSystemMessage   = &discordCodeError{discord.JSONErrorCodeCannotExecuteOnSystemMessage}
)

// Error is returned by all REST methods when Discord responds with a non-success status.
// It carries the HTTP status, the Discord JSON error code and message, and a field-level
// error map when present.
//
// Use errors.Is with the sentinel vars (ErrNotFound, ErrForbidden, etc.) for status checks,
// or errors.As to access the full error detail:
//
//	var apiErr *api.Error
//	if errors.As(err, &apiErr) {
//	    log.Printf("discord code %d: %s", apiErr.Code, apiErr.Message)
//	}
type Error struct {
	// HTTPStatus is the HTTP response status code (e.g. 403, 404).
	HTTPStatus int

	// Code is the Discord JSON error code. Zero when Discord did not return one.
	Code discord.JSONErrorCode

	// Message is the human-readable error message from Discord.
	Message string

	// Errors contains field-level validation errors, present on some 400 responses.
	Errors map[string]interface{}
}

func (e *Error) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("discord api error %d: %s (http %d)", int(e.Code), e.Message, e.HTTPStatus)
	}
	return fmt.Sprintf("discord api error: http %d", e.HTTPStatus)
}

// Is maps sentinel errors to HTTP status codes and Discord JSON error codes
// so callers can use errors.Is.
func (e *Error) Is(target error) bool {
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
	var ce *discordCodeError
	if errors.As(target, &ce) {
		return e.Code == ce.code
	}
	return false
}

// FieldError is one entry from a Discord validation error response — the nested
// "errors" object flattened to a single field. Path is the dotted location of
// the offending field (e.g. "embeds.0.fields.2.value"), empty for a top-level
// error; Code is Discord's machine-readable reason (e.g. "BASE_TYPE_REQUIRED")
// and Message the human-readable text.
type FieldError struct {
	Path    string
	Code    string
	Message string
}

// String renders the field error as "path: message (CODE)", omitting the path
// when empty.
func (f FieldError) String() string {
	var b strings.Builder
	if f.Path != "" {
		b.WriteString(f.Path)
		b.WriteString(": ")
	}
	b.WriteString(f.Message)
	if f.Code != "" {
		b.WriteString(" (")
		b.WriteString(f.Code)
		b.WriteByte(')')
	}
	return b.String()
}

// FieldErrors flattens the nested Errors map (Discord's "errors" object, with
// its per-field "_errors" arrays) into a sorted, ready-to-use slice. It returns
// nil when there are no field-level errors. See
// https://discord.com/developers/docs/reference#error-messages.
func (e *Error) FieldErrors() []FieldError {
	if e == nil || len(e.Errors) == 0 {
		return nil
	}
	var out []FieldError
	collectFieldErrors("", e.Errors, &out)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Code < out[j].Code
	})
	return out
}

// collectFieldErrors walks the nested error object depth-first, appending a
// FieldError for every entry found in an "_errors" array.
func collectFieldErrors(path string, node map[string]interface{}, out *[]FieldError) {
	for key, val := range node {
		if key == "_errors" {
			arr, ok := val.([]interface{})
			if !ok {
				continue
			}
			for _, item := range arr {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				fe := FieldError{Path: path}
				if c, ok := m["code"].(string); ok {
					fe.Code = c
				}
				if msg, ok := m["message"].(string); ok {
					fe.Message = msg
				}
				*out = append(*out, fe)
			}
			continue
		}
		child, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		next := key
		if path != "" {
			next = path + "." + key
		}
		collectFieldErrors(next, child, out)
	}
}

// CodeOf extracts the Discord JSON error code from err, unwrapping wrapped
// errors. The bool reports whether err was (or wrapped) an *api.Error; the code
// may still be 0 (the "General error" code) when true.
//
//	if code, ok := api.CodeOf(err); ok && code == discord.JSONErrorCodeUnknownMessage { ... }
func CodeOf(err error) (discord.JSONErrorCode, bool) {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Code, true
	}
	return 0, false
}

// HasCode reports whether err is (or wraps) an *api.Error carrying the given
// Discord JSON error code. It is shorthand for errors.Is against the matching
// sentinel.
func HasCode(err error, code discord.JSONErrorCode) bool {
	c, ok := CodeOf(err)
	return ok && c == code
}

// FieldErrorsOf returns the flattened field-level validation errors from err, or
// nil when err is not an *api.Error or carries none.
func FieldErrorsOf(err error) []FieldError {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.FieldErrors()
	}
	return nil
}
