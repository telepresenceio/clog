package clog

import (
	"fmt"
	"log/slog"
	"strings"
)

const (
	// LevelTrace is the most verbose log level.
	// When used, the output will also contain file:line information.
	LevelTrace = slog.LevelDebug - 4
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Level = slog.Level

// ParseLevel parses a log level string into a Level. It panics if the lowercased version of
// the string is not one of "trace", "debug", "info", "warn", "warning", or "error".
func ParseLevel(s string) (l Level) {
	switch strings.ToLower(s) {
	case "trace":
		l = LevelTrace
	case "debug":
		l = LevelDebug
	case "info":
		l = LevelInfo
	case "warn", "warning":
		l = LevelWarn
	case "error":
		l = LevelError
	default:
		panic(fmt.Sprintf("unknown log level: %q", s))
	}
	return l
}

type LevelSetter interface {
	SetLevel(Level)
}
