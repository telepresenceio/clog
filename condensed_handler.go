package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"sync/atomic"
	"unicode"
)

// NewCondensedHandler creates a new condensed slog.Handler with custom time format, log level, and output writer. The output writer must be thread-safe.
// The output format is: time level[ groups]: message[ : attrs][ (from file:line)].
//
//   - The time, level, and message are written without a leading "key=" prefix.
//   - The message is written without quotes.
//   - Top level groups are written as "group/subgroup" before the message.
//   - Attributes are written as "key=value" and the value is quoted if it contains Unicode space characters, non-printing characters, '"' or '='.
//   - The source file and line number are written after the message if the log level is [LevelTrace].
func NewCondensedHandler(out io.Writer, timeFormat string, level Level) slog.Handler {
	return &condensedHandler{
		timeFormat:    timeFormat,
		level:         int64(level),
		out:           out,
		includeSource: level < LevelDebug,
	}
}

type condensedHandler struct {
	timeFormat    string
	level         int64
	attrs         []slog.Attr
	groups        []string
	out           io.Writer
	includeSource bool
}

func addAttr(a slog.Attr, buf *bytesBuf) {
	if a.Value.Kind() == slog.KindGroup {
		addGroup(a.Key, a.Value.Group(), buf)
	} else {
		buf.writeString(a.Key)
		buf.writeByte('=')
		buf.writeString(quoteIfNeeded(a.Value.String()))
	}
}

func addAttrs(attrs []slog.Attr, buf *bytesBuf) {
	for i, a := range attrs {
		if i > 0 {
			buf.writeByte(' ')
		}
		addAttr(a, buf)
	}
}

func addGroup(name string, attrs []slog.Attr, buf *bytesBuf) {
	switch len(attrs) {
	case 0:
		return
	case 1:
		a0 := attrs[0]
		if a0.Value.Kind() == slog.KindGroup {
			// Stand-alone top-level group. Embed it directly into the name.
			buf.writeString(name)
			buf.writeByte('/')
			addGroup(a0.Key, a0.Value.Group(), buf)
			break
		}
		fallthrough
	default:
		buf.writeString(name)
		buf.writeString("={")
		addAttrs(attrs, buf)
		buf.writeByte('}')
	}
}

func (h *condensedHandler) Enabled(_ context.Context, level Level) bool {
	return int64(level) >= atomic.LoadInt64(&h.level)
}

// formatArgsKey is the key for the arguments to a printf style log message. The arguments are wrapped in
// a [slog.Attr] with this key and added to the log record so that the actual formatting can be deferred until
// the record is serialized.
const formatArgsKey = "formatArgs"

// extractFormatArgs extracts arguments intended for [fmt.Format] style logging from the record and returns them along with the remaining attributes.
// Format arguments that are [slog.Attr] are converted to strings before being returned.
func (h *condensedHandler) extractFormatArgs(record *slog.Record) ([]any, []slog.Attr) {
	var fmtArgs []any
	var otherAttrs []slog.Attr
	record.Attrs(func(a slog.Attr) bool {
		if a.Key == formatArgsKey {
			fmtArgs = a.Value.Any().([]any)
			for i, a := range fmtArgs {
				if attr, ok := a.(slog.Attr); ok {
					buf := newBuf()
					addAttr(attr, buf)
					fmtArgs[i] = string(*buf)
					buf.free()
				}
			}
		} else {
			otherAttrs = append(otherAttrs, a)
		}
		return true
	})
	return fmtArgs, slices.Concat(h.attrs, otherAttrs)
}

// levelString writes the log level as a string to buf, padded to 6 characters.
func levelString(l Level, buf *bytesBuf) {
	str := func(base string, val Level) {
		var n int
		if val == 0 {
			buf.writeString(base)
			n = len(base)
		} else {
			n, _ = fmt.Fprintf(buf, "%s%+d", base, val)
		}
		for i := 6; i > n; i-- {
			buf.writeByte(' ')
		}
	}
	switch {
	case l < LevelDebug:
		str("trace", l-LevelTrace)
	case l < LevelInfo:
		str("debug", l-LevelDebug)
	case l < LevelWarn:
		str("info", l-LevelInfo)
	case l < LevelError:
		str("warn", l-LevelWarn)
	default:
		str("error", l-LevelError)
	}
}

// Handle writes a condensed log record to the output writer. The writer is assumed to be thread-safe.
func (h *condensedHandler) Handle(_ context.Context, record slog.Record) error {
	buf := newBuf()
	buf.writeString(record.Time.Format(h.timeFormat))
	buf.writeByte(' ')
	levelString(record.Level, buf)
	fmtArgs, attrs := h.extractFormatArgs(&record)
	groups := h.groups

	// Merge stand-alone top level group into groups.
	ga := attrs
	for len(ga) == 1 && ga[0].Value.Kind() == slog.KindGroup {
		groups = append(groups, ga[0].Key)
		ga = ga[0].Value.Group()
	}
	attrs = ga
	if len(groups) > 0 {
		for i, g := range groups {
			if i > 0 {
				buf.writeByte('/')
			}
			buf.writeString(g)
		}
		buf.writeString(": ")
	}

	if len(fmtArgs) > 0 {
		_, _ = fmt.Fprintf(buf, record.Message, fmtArgs...)
	} else {
		buf.writeString(record.Message)
	}
	if len(attrs) > 0 {
		buf.writeString(" : ")
		addAttrs(attrs, buf)
	}
	if h.includeSource {
		src := record.Source()
		if src != nil {
			_, _ = fmt.Fprintf(buf, " (from %s:%d)", src.File, src.Line)
		}
	}
	buf.writeByte('\n')
	_, err := h.out.Write(*buf)
	buf.free()
	return err
}

func (h *condensedHandler) SetLevel(l Level) {
	atomic.StoreInt64(&h.level, int64(l))
}

func (h *condensedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	h2.attrs = append(h2.attrs, attrs...)
	return &h2
}

func (h *condensedHandler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.groups = append(h2.groups, name)
	return &h2
}

func quoteIfNeeded(s string) string {
	for _, c := range s {
		switch {
		case c < 32, c == '=', c == '"':
			return strconv.Quote(s)
		case unicode.IsSpace(c), !unicode.IsPrint(c):
			return strconv.Quote(s)
		}
	}
	return s
}
