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
func ParseLevel(s string) (l Level, err error) {
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
		return 0, fmt.Errorf("unknown log level: %q", s)
	}
	return l, nil
}

// MustParseLevel is like ParseLevel, but panics if the string is invalid.
func MustParseLevel(s string) Level {
	l, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return l
}

type LevelSetter interface {
	SetLevel(Level)
}
