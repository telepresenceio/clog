package log

import (
	"context"
	stdLog "log"

	"github.com/telepresenceio/clog"
	"github.com/telepresenceio/clog/internal"
)

// Debugf is [clog.Debugf] using the default logger.
func Debugf(format string, args ...any) {
	internal.Logf(context.Background(), clog.LevelDebug, format, args...)
}

// Errorf is [clog.Errorf] using the default logger.
func Errorf(format string, args ...any) {
	internal.Logf(context.Background(), clog.LevelError, format, args...)
}

// Infof is [clog.Infof] using the default logger.
func Infof(format string, args ...any) {
	internal.Logf(context.Background(), clog.LevelInfo, format, args...)
}

// Logf is [clog.Logf] using the default logger.
func Logf(level clog.Level, format string, args ...any) {
	internal.Logf(context.Background(), level, format, args...)
}

// StdLogger is [clog.StdLogger] using the default logger.
func StdLogger(level clog.Level) *stdLog.Logger {
	return clog.StdLogger(context.Background(), level)
}

// Tracef is [clog.Tracef] using the default logger.
func Tracef(format string, args ...any) {
	internal.Logf(context.Background(), clog.LevelTrace, format, args...)
}

// Warnf is [clog.Warnf] using the default logger.
func Warnf(format string, args ...any) {
	internal.Logf(context.Background(), clog.LevelWarn, format, args...)
}
