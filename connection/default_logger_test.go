package connection

import (
	"context"
	"log/slog"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestDefaultLogger_IsSilent verifies the library is silent by default: with no
// WithLogger / WithLogLevel, the Client's logger discards everything (it is not
// slog.Default(), which would write to the application's stderr at Info).
func (cs *ConnectionSuite) TestDefaultLogger_IsSilent() {
	t := cs.T()
	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	cs.Require().NotNil(c.Logger)

	ctx := context.Background()
	cs.False(c.Logger.Enabled(ctx, slog.LevelError), "default logger must discard even error logs")
	cs.False(c.Logger.Enabled(ctx, slog.LevelInfo), "default logger must be silent at info")
	cs.False(c.Logger.Enabled(ctx, slog.LevelDebug), "default logger must be silent at debug")
}

// TestWithLogLevel_EnablesLogging confirms the opt-in path still works: setting
// a level produces a logger enabled at/above that level.
func (cs *ConnectionSuite) TestWithLogLevel_EnablesLogging() {
	t := cs.T()
	c, err := NewClient("Bot fake-token", discord.IntentGuilds, options.WithLogLevel(slog.LevelWarn))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ctx := context.Background()
	cs.True(c.Logger.Enabled(ctx, slog.LevelWarn), "warn must be enabled when level=warn")
	cs.False(c.Logger.Enabled(ctx, slog.LevelInfo), "info is below warn → disabled")
}
