package clog

import (
	"fmt"
	"log/slog"
	"strings"
)

// ParseLevel parses a log level string into a slog.Level. It panics if the lowercased version of
// the string is not one of "trace", "debug", "info", "warn", "warning", or "error".
func ParseLevel(s string) (l slog.Level) {
	switch strings.ToLower(s) {
	case "trace":
		l = LevelTrace
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn", "warning":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		panic(fmt.Sprintf("unknown log level: %q", s))
	}
	return l
}
