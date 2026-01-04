package internal

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

type Level = slog.Level

// TimeNow is for test purposes only.
var TimeNow = time.Now

type FormatHandler interface {
	HandleFormat(context.Context, *slog.Record, []any) error
}

func Log(ctx context.Context, level slog.Level, args ...any) {
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

func LogAttrs(ctx context.Context, level slog.Level, message string, attrs ...slog.Attr) {
	h := Logger(ctx).Handler()
	if h.Enabled(ctx, level) {
		r := newRecord(level, message)
		r.AddAttrs(attrs...)
		_ = h.Handle(ctx, r)
	}
}

func Logf(ctx context.Context, level slog.Level, format string, args ...any) {
	h := Logger(ctx).Handler()
	if h.Enabled(ctx, level) {
		var r slog.Record
		if fh, ok := h.(FormatHandler); ok {
			// Defer formatting so that the handler's internal buffer can be used instead of
			// having fmt.Sprintf allocating one here.
			r = newRecord(level, format)
			_ = fh.HandleFormat(ctx, &r, args)
		} else {
			// This is unfortunate, but slog does not provide a way to defer the creation of the actual message.
			_ = h.Handle(ctx, newRecord(level, fmt.Sprintf(format, args...)))
		}
	}
}

func newRecord(level Level, msg string) slog.Record {
	var pcs [1]uintptr
	runtime.Callers(4, pcs[:]) // skip [Callers, newRecord, log, caller-of-log]
	return slog.NewRecord(TimeNow(), level, msg, pcs[0])
}

type loggerKey struct{}

// Logger returns the logger from the context, or the default logger if none is set.
func Logger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// WithLogger assigns the logger to a child context which is returned.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
