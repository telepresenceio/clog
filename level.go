package clog

import (
	"fmt"
	"log/slog"
	"strings"
)

// LevelTrace is the most verbose log level.
// When used, the output will also contain file:line information.
const LevelTrace = slog.LevelDebug - 4

// LevelWithTrace redefines slog.Level with an UnmarshalText to support the "trace" log level.
// It's suitable when used as a field that well be subjected to automatic parsing but
// shouldn't otherwise be used. Use `slog.Level` whenever possible.
type LevelWithTrace slog.Level

// LevelStrings is the list of log level strings in order of increasing verbosity.
var LevelStrings = []string{
	"ERROR", "WARN", "INFO", "DEBUG", "TRACE",
}

func (l LevelWithTrace) String() (s string) {
	sl := slog.Level(l)
	diff := sl - slog.LevelDebug
	switch {
	case diff == -4:
		s = "TRACE"
	case diff < 0:
		s = fmt.Sprintf("TRACE%+d", diff)
	default:
		s = sl.String()
	}
	return s
}

// MarshalText implements encoding.TextMarshaler.
func (l LevelWithTrace) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText parses a log level string into a Level. It panics if the lowercased version of
// the string is not one of "trace", "debug", "info", "warn", "warning", or "error".
//goland:noinspection GoMixedReceiverTypes
func (l *LevelWithTrace) UnmarshalText(value []byte) error {
	var sl slog.Level
	err := sl.UnmarshalText(value)
	switch {
	case err == nil:
		*l = LevelWithTrace(sl)
	case strings.EqualFold(string(value), "TRACE"):
		*l = LevelWithTrace(LevelTrace)
		err = nil
	}
	return err
}

// ParseLevel parses a log level string into a Level. It panics if the lowercased version of
// the string is not one of "trace", "debug", "info", "warn", "warning", or "error".
func ParseLevel(s string) (slog.Level, error) {
	var l LevelWithTrace
	err := l.UnmarshalText([]byte(s))
	return slog.Level(l), err
}

// MustParseLevel is like ParseLevel, but panics if the string is invalid.
func MustParseLevel(s string) slog.Level {
	l, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return l
}
