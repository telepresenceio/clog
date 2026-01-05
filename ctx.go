package clog

import (
	"context"
	stdLog "log"
	"log/slog"

	"github.com/telepresenceio/clog/internal"
)

type FormatHandler = internal.FormatHandler

// Debug is similar to [slog.Logger.DebugContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Debug(ctx context.Context, args ...any) {
	internal.Log(ctx, slog.LevelDebug, args...)
}

// DebugAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelDebug].
func DebugAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	internal.LogAttrs(ctx, slog.LevelDebug, message, attrs...)
}

// Debugf is similar to [slog.Logger.DebugContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Debugf(ctx context.Context, format string, args ...any) {
	internal.Logf(ctx, slog.LevelDebug, format, args...)
}

func Enabled(ctx context.Context, level slog.Level) bool {
	return internal.Logger(ctx).Enabled(ctx, level)
}

// Error is similar to [slog.Logger.ErrorContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Error(ctx context.Context, args ...any) {
	internal.Log(ctx, slog.LevelError, args...)
}

// ErrorAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelError].
func ErrorAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	internal.LogAttrs(ctx, slog.LevelError, message, attrs...)
}

// Errorf is similar to [slog.Logger.ErrorContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Errorf(ctx context.Context, format string, args ...any) {
	internal.Logf(ctx, slog.LevelError, format, args...)
}

// Info is similar to [slog.Logger.InfoContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Info(ctx context.Context, args ...any) {
	internal.Log(ctx, slog.LevelInfo, args...)
}

// InfoAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelInfo].
func InfoAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	internal.LogAttrs(ctx, slog.LevelInfo, message, attrs...)
}

// Infof is similar to [slog.Logger.InfoContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Infof(ctx context.Context, format string, args ...any) {
	internal.Logf(ctx, slog.LevelInfo, format, args...)
}

// Log is similar to calling [slog.Logger.Log] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Log(ctx context.Context, level slog.Level, args ...any) {
	internal.Log(ctx, level, args...)
}

// LogAttrs is similar to [slog.Logger.LogAttrs] on the context logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	internal.LogAttrs(ctx, level, msg, attrs...)
}

// Logf is similar to [slog.Logger.Log] on the context logger, called with a message formatted with [fmt.Sprintf].
func Logf(ctx context.Context, level slog.Level, format string, args ...any) {
	internal.Logf(ctx, level, format, args...)
}

// Logger returns the logger from the context, or the default logger if none is set.
func Logger(ctx context.Context) *slog.Logger {
	return internal.Logger(ctx)
}

func StdLogger(ctx context.Context, level slog.Level) *stdLog.Logger {
	return stdLog.New(&writer{log: Logger(ctx), level: level, ctx: ctx}, "", 0)
}

// Trace is similar to [slog.Logger.Log] on the context logger, called with [slog.LevelTrace].
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Trace(ctx context.Context, args ...any) {
	internal.Log(ctx, LevelTrace, args...)
}

// TraceAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelTrace].
func TraceAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	internal.LogAttrs(ctx, LevelTrace, message, attrs...)
}

// Tracef is similar to [slog.Logger.Log] on the context logger, called with [slog.LevelTrace] and a message formatted with [fmt.Sprintf].
func Tracef(ctx context.Context, format string, args ...any) {
	internal.Logf(ctx, LevelTrace, format, args...)
}

// Warn is similar to [slog.Logger.WarnContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Warn(ctx context.Context, args ...any) {
	internal.Log(ctx, slog.LevelWarn, args...)
}

// WarnAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelWarn].
func WarnAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	internal.LogAttrs(ctx, slog.LevelWarn, message, attrs...)
}

// Warnf is similar to [slog.Logger.WarnContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Warnf(ctx context.Context, format string, args ...any) {
	internal.Logf(ctx, slog.LevelWarn, format, args...)
}

// With creates a clone of the context logger, adds the attributes, and assigns the clone to a child context which is returned.
func With(ctx context.Context, args ...any) context.Context {
	if len(args) == 0 {
		return ctx
	}
	return internal.WithLogger(ctx, internal.Logger(ctx).With(args...))
}

// WithLogger assigns the logger to a child context which is returned.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return internal.WithLogger(ctx, logger)
}

// WithGroup creates a clone of the context logger, adds the group, and assigns the clone to a child context which is returned.
func WithGroup(ctx context.Context, group string) context.Context {
	if group == "" {
		return ctx
	}
	return internal.WithLogger(ctx, internal.Logger(ctx).WithGroup(group))
}
