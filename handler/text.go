package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"unicode"
)

const RFC3339MillisNoTz = "2006-01-02T15:04:05.000"

// NewText creates a new slog.Handler. Unless overridden by options, the handler writes to [os.Stdout]
// using [RFC3339MillisNoTz] time format and the [LevelWarn] level.
// The output format is: time level[ groups]: message[ : attrs][ (from file:line)].
//
//   - The time, level, and message are written without a leading "key=" prefix.
//   - The message is written without quotes.
//   - Top level groups are written as "group/subgroup" before the message.
//   - Attributes are written as "key=value" and the value is quoted if it contains Unicode space characters, non-printing characters, '"' or '='.
//   - The source file and line number are written after the message if the log level is [LevelTrace].
func NewText(options ...Option) slog.Handler {
	h := &textHandler{
		out:             allLevelsWriter{out: os.Stdout},
		timeFormat:      RFC3339MillisNoTz,
		levelEnabler:    func(_ context.Context, level slog.Level) bool { return level >= slog.LevelWarn },
		hideLevelsAbove: slog.Level(math.MaxInt)}
	for _, opt := range options {
		opt(h)
	}
	return h
}

// Handle writes a log record to the output writer. The writer is assumed to be thread-safe.
func (h *textHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.HandleFormat(ctx, &record, nil)
}

func (h *textHandler) HandleFormat(_ context.Context, record *slog.Record, fmtArgs []any) error {
	buf := newBuf()
	if h.timeFormat != "" {
		buf.writeString(record.Time.Format(h.timeFormat))
		buf.writeByte(' ')
	}
	if record.Level < h.hideLevelsAbove {
		levelString(record.Level, buf)
	}

	hasGroups := false
	firstGroup := true
	writeGroup := func(name string) {
		if firstGroup {
			firstGroup = false
		} else {
			buf.writeByte('/')
		}
		buf.writeString(name)
		hasGroups = true
	}
	for _, g := range h.groups {
		writeGroup(g)
	}

	// Merge stand-alone top level group into groups.
	if record.NumAttrs() == 1 {
		record.Attrs(func(a slog.Attr) bool {
			if a.Value.Kind() == slog.KindGroup {
				writeGroup(a.Key)
				ga := a.Value.Group()
				for len(ga) == 1 && ga[0].Value.Kind() == slog.KindGroup {
					writeGroup(ga[0].Key)
					ga = ga[0].Value.Group()
				}
				nr := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
				nr.AddAttrs(ga...)
				record = &nr
			}
			return false
		})
	}
	if hasGroups {
		buf.writeString(": ")
	}

	if len(fmtArgs) > 0 {
		_, _ = fmt.Fprintf(buf, record.Message, fmtArgs...)
	} else {
		buf.writeString(record.Message)
	}
	if record.NumAttrs() > 0 {
		buf.writeString(" : ")
		first := true
		record.Attrs(func(a slog.Attr) bool {
			if first {
				first = false
			} else {
				buf.writeByte(' ')
			}
			addAttr(a, buf)
			return true
		})
	}
	if h.includeSource {
		src := record.Source()
		if src != nil {
			buf.writeString(" (from ")
			buf.writeString(src.File)
			buf.writeByte(':')
			buf.writeString(strconv.Itoa(src.Line))
			buf.writeByte(')')
		}
	}
	buf.writeByte('\n')
	_, err := h.out.Write(record.Level, *buf)
	buf.free()
	return err
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	h2.attrs = append(h2.attrs, attrs...)
	return &h2
}

func (h *textHandler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.groups = append(h2.groups, name)
	return &h2
}

type textHandler struct {
	timeFormat      string
	levelEnabler    EnabledFunc
	hideLevelsAbove slog.Level
	attrs           []slog.Attr
	groups          []string
	out             LevelWriter
	includeSource   bool
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

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.levelEnabler(ctx, level)
}

// levelString writes the log level as a string to buf, padded to 6 characters.
func levelString(l slog.Level, buf *bytesBuf) {
	str := func(base string, val slog.Level) {
		var n int
		if val == 0 {
			buf.writeString(base)
			n = len(base)
		} else {
			n, _ = fmt.Fprintf(buf, "%s%+d", base, val)
		}
		buf.writeByte(' ')
		for i := 5; i > n; i-- {
			buf.writeByte(' ')
		}
	}
	switch {
	case l < slog.LevelDebug:
		str("TRACE", l-slog.LevelDebug+4)
	case l < slog.LevelInfo:
		str("DEBUG", l-slog.LevelDebug)
	case l < slog.LevelWarn:
		str("INFO", l-slog.LevelInfo)
	case l < slog.LevelError:
		str("WARN", l-slog.LevelWarn)
	default:
		str("ERROR", l-slog.LevelError)
	}
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
