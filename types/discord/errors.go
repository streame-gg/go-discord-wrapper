package discord

import "errors"

// JSONErrorCarrier is implemented by errors that carry a Discord JSON error
// code — notably the *api.Error returned by every REST call. It lets helpers in
// this package inspect API errors without importing the api package (which would
// be an import cycle). Custom error types may implement it to participate in
// IsErrorCode / ErrorCodeOf.
// https://docs.discord.com/developers/topics/opcodes-and-status-codes#json
type JSONErrorCarrier interface {
	JSONErrorCode() JSONErrorCode
}

// ErrorCodeOf returns the Discord JSON error code carried by err, unwrapping
// wrapped errors. The bool is false when err is nil or carries no Discord code.
// Note the code may legitimately be 0 (JSONErrorCodeGeneral) when the bool is
// true.
func ErrorCodeOf(err error) (JSONErrorCode, bool) {
	var c JSONErrorCarrier
	if errors.As(err, &c) {
		return c.JSONErrorCode(), true
	}
	return 0, false
}

// IsErrorCode reports whether err (or any error it wraps) is a Discord API error
// with the given JSON error code. It is the easiest way to branch on a specific
// failure:
//
//	if discord.IsErrorCode(err, discord.JSONErrorCodeMissingPermissions) {
//	    // the bot lacks permission for this action
//	}
func IsErrorCode(err error, code JSONErrorCode) bool {
	c, ok := ErrorCodeOf(err)
	return ok && c == code
}
