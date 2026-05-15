package api

import (
	"errors"
	"fmt"
	"net/http"

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
	code discord.GatewayErrorCode
}

func (e *discordCodeError) Error() string {
	return fmt.Sprintf("discord: error code %d", int(e.code))
}

// Typed sentinels for the most commonly checked Discord JSON error codes.
// Use errors.Is to test:
//
//	if errors.Is(err, api.ErrMissingPermissions) { ... }
var (
	ErrUnknownChannel                 = &discordCodeError{discord.GatewayErrorCodeUnknownChannel}
	ErrUnknownGuild                   = &discordCodeError{discord.GatewayErrorCodeUnknownGuild}
	ErrUnknownMessage                 = &discordCodeError{discord.GatewayErrorCodeUnknownMessage}
	ErrUnknownMember                  = &discordCodeError{discord.GatewayErrorCodeUnknownMember}
	ErrUnknownRole                    = &discordCodeError{discord.GatewayErrorCodeUnknownRole}
	ErrUnknownWebhook                 = &discordCodeError{discord.GatewayErrorCodeUnknownWebhook}
	ErrUnknownUser                    = &discordCodeError{discord.GatewayErrorCodeUnknownUser}
	ErrUnknownEmoji                   = &discordCodeError{discord.GatewayErrorCodeUnknownEmoji}
	ErrUnknownInteraction             = &discordCodeError{discord.GatewayErrorCodeUnknownInteraction}
	ErrMissingAccess                  = &discordCodeError{discord.GatewayErrorCodeMissingAccess}
	ErrMissingPermissions             = &discordCodeError{discord.GatewayErrorCodeMissingPermissions}
	ErrCannotSendMessagesToUser       = &discordCodeError{discord.GatewayErrorCodeCannotSendMessagesToThisUser}
	ErrInteractionAlreadyAcknowledged = &discordCodeError{discord.GatewayErrorCodeInteractionAlreadyAcknowledged}
	ErrThreadIsLocked                 = &discordCodeError{discord.GatewayErrorCodeThreadIsLocked}
	ErrMessageTooOldForBulkDelete     = &discordCodeError{discord.GatewayErrorCodeMessageTooOldForBulkDelete}
	ErrInvalidFormBody                = &discordCodeError{discord.GatewayErrorCodeInvalidFormBody}
	ErrMaxReactionsReached            = &discordCodeError{discord.GatewayErrorCodeMaxReactions}
	ErrCannotExecuteOnSystemMessage   = &discordCodeError{discord.GatewayErrorCodeCannotExecuteOnSystemMessage}
)

// Error is returned by all REST methods when Discord responds with a non-success status.
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
type Error struct {
	// HTTPStatus is the HTTP response status code (e.g. 403, 404).
	HTTPStatus int

	// Code is the Discord JSON error code. Zero when Discord did not return one.
	Code discord.GatewayErrorCode

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
