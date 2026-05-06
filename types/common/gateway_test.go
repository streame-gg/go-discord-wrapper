package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatewayErrorString(t *testing.T) {
	e := GatewayError{Code: 10003, Message: "Unknown Channel"}
	assert.Equal(t, "10003 Unknown Channel", e.Error())
}

func TestGatewayErrorZeroCode(t *testing.T) {
	e := GatewayError{Code: 0, Message: "Missing Access"}
	assert.Equal(t, "0 Missing Access", e.Error())
}
