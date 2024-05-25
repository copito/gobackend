package setup

import (
	"log/slog"
)

func SetupLogging(withColor bool, level slog.Level) *slog.Logger {
	handler := NewCustomLogHandler(withColor, level)
	// logger := slog.New(slog.New(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger := slog.New(handler)

	return logger
}
