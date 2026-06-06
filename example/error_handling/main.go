// error_handling demonstrates how to inspect Discord REST errors with strong
// typing: matching a specific JSON error code (e.g. MissingPermissions), reading
// the HTTP status, and flattening field/array-level validation errors.
//
// Every REST method (and the entity convenience methods that wrap them, like
// channel.SetVoiceStatus) returns an *api.Error on failure, carrying a
// discord.JSONErrorCode. The api.CodeOf / api.HasCode / errors.Is helpers pull
// that out for you and unwrap wrapped errors along the way.
//
// Run with:
//
//	DISCORD_TOKEN=... CHANNEL_ID=... go run ./example/error_handling
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func main() {
	rest, err := api.NewRestClient(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		slog.Error("create rest client", "err", err)
		os.Exit(1)
	}

	channelID, err := discord.ParseSnowflake(os.Getenv("CHANNEL_ID"))
	if err != nil {
		slog.Error("set CHANNEL_ID to a valid channel id", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// 1) Permission error: setting a voice channel status without the right
	//    permission (or on a non-voice channel) returns 50013 Missing Permissions.
	status := "gaming"
	if err := rest.SetVoiceChannelStatus(ctx, *channelID, &status); err != nil {
		handleVoiceError(err)
	}

	// 2) Validation error: an empty message fails form-body validation (50035)
	//    and carries per-field _errors, including array-indexed ones.
	if _, err := rest.CreateMessage(ctx, *channelID, api.CreateMessageParams{}); err != nil {
		handleValidationError(err)
	}
}

// handleVoiceError shows the three ways to check for a specific Discord code.
func handleVoiceError(err error) {
	switch {
	// (a) Sentinel — most readable for the common codes.
	case errors.Is(err, api.ErrMissingPermissions):
		slog.Warn("can't set status: the bot is missing the required permission here")

	// (b) HasCode — works for any of the JSON error codes, no sentinel needed.
	case api.HasCode(err, discord.JSONErrorCodeBitrateTooHighForChannelType):
		slog.Warn("bitrate too high for this channel type")

	// (c) CodeOf — pull the typed code out for logging or a switch.
	default:
		if code, ok := api.CodeOf(err); ok {
			slog.Error("set voice status failed", "code", int(code), "err", err)
		} else {
			slog.Error("set voice status failed (non-API error)", "err", err)
		}
	}
}

// handleValidationError pulls the HTTP status, the top-level code, and every
// field/array-level validation error out of a failed request.
func handleValidationError(err error) {
	var apiErr *api.Error
	if !errors.As(err, &apiErr) {
		slog.Error("request failed", "err", err)
		return
	}

	fmt.Printf("HTTP %d — code %d: %s\n", apiErr.HTTPStatus, int(apiErr.Code), apiErr.Message)

	if apiErr.Code == discord.JSONErrorCodeInvalidFormBody {
		// FieldErrors flattens the nested errors object (including array indices)
		// into "path: message (CODE)" entries, e.g. "embeds.0.title".
		for _, fe := range apiErr.FieldErrors() {
			fmt.Printf("  • %s\n", fe)
		}
	}
}
