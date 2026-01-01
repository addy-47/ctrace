package utils

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes the global logger.
// In a real app, you might want to configure level and format based on config.
func InitLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	Log = logger
}
