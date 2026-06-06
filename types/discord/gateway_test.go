package discord

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (su *discordGatewaySuite) TestGatewayErrorString() {
	t := su.T()
	e := GatewayError{Code: 10003, Message: "Unknown Channel"}
	assert.Equal(t, "10003 Unknown Channel", e.Error())
}

func (su *discordGatewaySuite) TestGatewayErrorZeroCode() {
	t := su.T()
	e := GatewayError{Code: 0, Message: "Missing Access"}
	assert.Equal(t, "0 Missing Access", e.Error())
}

type discordGatewaySuite struct{ suite.Suite }

func TestDiscordGatewaySuite(t *testing.T) { suite.Run(t, new(discordGatewaySuite)) }
