package api

import (
	"context"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// GetGateway returns the WebSocket URL for connecting to the Discord gateway.
func (c *RestClient) GetGateway(ctx context.Context) (*discord.GatewayResponse, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/gateway", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GatewayResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetBotGateway returns the WebSocket URL, recommended shard count, and session start
// limit for the bot. Requires bot token authentication.
func (c *RestClient) GetBotGateway(ctx context.Context) (*discord.BotRegisterResponse, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/gateway/bot", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.BotRegisterResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
