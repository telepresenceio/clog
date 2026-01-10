package testutil

import (
	"context"
	"log/slog"
	"testing"

	"github.com/telepresenceio/clog"
	"github.com/telepresenceio/clog/handler"
)

type trapError struct {
	slog.Handler
	t *testing.T
}

func (t *trapError) Handle(ctx context.Context, record slog.Record) (err error) {
	if record.Level >= slog.LevelError {
		args := make([]any, 0, 1+record.NumAttrs())
		args = append(args, record.Message)
		record.Attrs(func(a slog.Attr) bool { args = append(args, a); return true })
		t.t.Error(args...)
	} else {
		err = t.Handler.Handle(ctx, record)
	}
	return err
}

// NewContext returns a context with a logger that writes to t.
// If failOnError is true, errors are propagated to t.Error instead of being logged.
func NewContext(t *testing.T, failOnError bool) context.Context {
	h := handler.NewText(handler.Output(t.Output()), handler.EnabledLevel(slog.LevelDebug))
	if failOnError {
		h = &trapError{Handler: h, t: t}
	}
	return clog.WithLogger(context.Background(), slog.New(h))
}
