package clog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
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

// UnmarshalText will unmarshal a log level string into a Level.
// An error is returned if the string is not one of "TRACE", "DEBUG", "INFO", "WARN", or "ERROR"
// using case-insensitive comparison.
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

type treeLevelKey struct{}

// WithTreeLevel assigns the log level root to a child context which is returned. Children
// of this context will inherit the level.
func WithTreeLevel(ctx context.Context, level slog.Level) context.Context {
	lv := int64(level)
	return context.WithValue(ctx, treeLevelKey{}, &lv)
}

// TreeEnabled returns true if the root [slog.Level] is enabled for the given context.
// This function is suitable as an argument to a [handler.LevelEnabler] option.
func TreeEnabled(ctx context.Context, level slog.Level) bool {
	if lvp, ok := ctx.Value(treeLevelKey{}).(*int64); ok {
		return int64(level) >= atomic.LoadInt64(lvp)
	}
	return false
}

// SetTreeLevel sets the root [slog.Level] for the context where the [WithTreeLevel] was set, which
// might be the provided context itself or any parent of that context. Any children of context holding
// the root level is affected by this change.
func SetTreeLevel(ctx context.Context, level slog.Level) bool {
	if lvp, ok := ctx.Value(treeLevelKey{}).(*int64); ok {
		oldLevel := atomic.LoadInt64(lvp)
		newLevel := int64(level)
		if newLevel != oldLevel {
			return atomic.CompareAndSwapInt64(lvp, oldLevel, newLevel)
		}
	}
	return false
}
