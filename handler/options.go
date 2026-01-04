package handler

import (
	"io"
	"log/slog"

	"github.com/telepresenceio/clog"
)

type Option func(*textHandler)

// EnabledLevel sets the minimum log level to be handled by the handler.
func EnabledLevel(level clog.Level) Option {
	return func(h *textHandler) {
		h.SetLevel(level)
	}
}

// Attrs adds the specified attributes to all log records handled by the handler.
func Attrs(attrs ...slog.Attr) Option {
	return func(h *textHandler) {
		h.attrs = attrs
	}
}

// Groups adds the specified groups to all log records handled by the handler.
func Groups(groups ...string) Option {
	return func(h *textHandler) {
		h.groups = groups
	}
}

// HideLevel hides the level field when the logged level is equal to or above the specified level. This is
// particularly useful when log entries for a specific level end up in a log of their own, making the actual
// level information redundant. Example:
//
//	HideLevel(LevelWarn) hides all levels LevelWarn and LevelError.
func HideLevel(level clog.Level) Option {
	return func(h *textHandler) {
		h.hideLevelsAbove = level
	}
}

// IncludeSource adds the source file and line number to the log record.
func IncludeSource(include bool) Option {
	return func(h *textHandler) {
		h.includeSource = include
	}
}

// LevelOutput sets a writer that is capable of sending output to different locations depending on the log level.
// LevelOutput is mutually exclusive with Output.
// The writer must be thread-safe.
func LevelOutput(lw LevelWriter) Option {
	return func(h *textHandler) {
		h.out = lw
	}
}

// Output sets the writer that receives all log records.
// Output is mutually exclusive with LevelOutput.
// The writer must be thread-safe.
func Output(w io.Writer) Option {
	return func(h *textHandler) {
		h.out = allLevelsWriter{out: w}
	}
}

// TimeFormat sets the time format used for log records. The records will be logged without a timestamp if the timeFormat is "".
func TimeFormat(timeFormat string) Option {
	return func(h *textHandler) {
		h.timeFormat = timeFormat
	}
}
