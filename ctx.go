package clog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

// TimeNow is for test purposes only.
var TimeNow = time.Now

// Debug is similar to [slog.Logger.DebugContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Debug(ctx context.Context, args ...any) {
	log(ctx, slog.LevelDebug, args...)
}

// DebugAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelDebug].
func DebugAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	logAttrs(ctx, slog.LevelDebug, message, attrs...)
}

// Debugf is similar to [slog.Logger.DebugContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Debugf(ctx context.Context, format string, args ...any) {
	logf(ctx, slog.LevelDebug, format, args...)
}

func Enabled(ctx context.Context, level slog.Level) bool {
	return Logger(ctx).Enabled(ctx, level)
}

// Error is similar to [slog.Logger.ErrorContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Error(ctx context.Context, args ...any) {
	log(ctx, slog.LevelError, args...)
}

// ErrorAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelError].
func ErrorAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	logAttrs(ctx, slog.LevelError, message, attrs...)
}

// Errorf is similar to [slog.Logger.ErrorContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Errorf(ctx context.Context, format string, args ...any) {
	logf(ctx, slog.LevelError, format, args...)
}

// Info is similar to [slog.Logger.InfoContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Info(ctx context.Context, args ...any) {
	log(ctx, slog.LevelInfo, args...)
}

// InfoAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelInfo].
func InfoAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	logAttrs(ctx, slog.LevelInfo, message, attrs...)
}

// Infof is similar to [slog.Logger.InfoContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Infof(ctx context.Context, format string, args ...any) {
	logf(ctx, slog.LevelInfo, format, args...)
}

// Log is similar to calling [slog.Logger.Log] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Log(ctx context.Context, level slog.Level, args ...any) {
	log(ctx, level, args...)
}

// LogAttrs is similar to [slog.Logger.LogAttrs] on the context logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logAttrs(ctx, level, msg, attrs...)
}

// Logf is similar to [slog.Logger.Log] on the context logger, called with a message formatted with [fmt.Sprintf].
func Logf(ctx context.Context, level slog.Level, format string, args ...any) {
	logf(ctx, level, format, args...)
}

// Logger returns the logger from the context, or the default logger if none is set.
func Logger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// Trace is similar to [slog.Logger.Log] on the context logger, called with [LevelTrace].
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Trace(ctx context.Context, args ...any) {
	log(ctx, LevelTrace, args...)
}

// TraceAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [LevelTrace].
func TraceAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	logAttrs(ctx, LevelTrace, message, attrs...)
}

// Tracef is similar to [slog.Logger.Log] on the context logger, called with [LevelTrace] and a message formatted with [fmt.Sprintf].
func Tracef(ctx context.Context, format string, args ...any) {
	logf(ctx, LevelTrace, format, args...)
}

// Warn is similar to [slog.Logger.WarnContext] on the context logger.
// The first argument is converted using [fmt.Sprint] and used as the log message. Remaining args
// are handled according to [slog.Logger.Log].
func Warn(ctx context.Context, args ...any) {
	log(ctx, slog.LevelWarn, args...)
}

// WarnAttrs is similar to [slog.Logger.LogAttrs] on the context logger, called with [slog.LevelWarn].
func WarnAttrs(ctx context.Context, message string, attrs ...slog.Attr) {
	logAttrs(ctx, slog.LevelWarn, message, attrs...)
}

// Warnf is similar to [slog.Logger.WarnContext] on the context logger, called with a message formatted with [fmt.Sprintf].
func Warnf(ctx context.Context, format string, args ...any) {
	logf(ctx, slog.LevelWarn, format, args...)
}

// With creates a clone of the context logger, adds the attributes, and assigns the clone to a child context which is returned.
func With(ctx context.Context, args ...any) context.Context {
	if len(args) == 0 {
		return ctx
	}
	return WithLogger(ctx, Logger(ctx).With(args...))
}

// WithLogger assigns the logger to a child context which is returned.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithGroup creates a clone of the context logger, adds the group, and assigns the clone to a child context which is returned.
func WithGroup(ctx context.Context, group string) context.Context {
	if group == "" {
		return ctx
	}
	return WithLogger(ctx, Logger(ctx).WithGroup(group))
}

type loggerKey struct{}

func log(ctx context.Context, level slog.Level, args ...any) {
	h := Logger(ctx).Handler()
	if h.Enabled(ctx, level) && len(args) > 0 {
		a0 := args[0]
		msg, ok := a0.(string)
		if !ok {
			msg = fmt.Sprint(a0)
		}
		r := newRecord(level, msg)
		r.Add(args[1:]...)
		_ = h.Handle(ctx, r)
	}
}

func logAttrs(ctx context.Context, level slog.Level, message string, attrs ...slog.Attr) {
	h := Logger(ctx).Handler()
	if h.Enabled(ctx, level) {
		r := newRecord(level, message)
		r.AddAttrs(attrs...)
		_ = h.Handle(ctx, r)
	}
}

func logf(ctx context.Context, level slog.Level, format string, args ...any) {
	h := Logger(ctx).Handler()
	if h.Enabled(ctx, level) {
		var r slog.Record
		if _, ok := h.(*condensedHandler); ok {
			// Defer formatting so that the handlers internal buffer can be used instead of
			// having fmt.Sprintf allocating one here.
			r = newRecord(level, format)
			r.AddAttrs(slog.Any(formatArgsKey, args))
		} else {
			// This is unfortunate, but slog does not provide a way to defer the creation of the actual message.
			r = newRecord(level, fmt.Sprintf(format, args...))
		}
		_ = h.Handle(ctx, r)
	}
}

func newRecord(level slog.Level, msg string) slog.Record {
	var pcs [1]uintptr
	runtime.Callers(4, pcs[:]) // skip [Callers, newRecord, log, caller-of-log]
	return slog.NewRecord(TimeNow(), level, msg, pcs[0])
}
