package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Logger is a lightweight structured logger wrapper around slog.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger configured for text or JSON based on environment.
func NewLogger(serviceName string, env string, logLevel string) *Logger {
	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(env) == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler).With(
		slog.String("service", serviceName),
		slog.String("env", env),
	)

	return &Logger{Logger: logger}
}
