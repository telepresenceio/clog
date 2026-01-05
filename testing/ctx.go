package testing

import (
	"context"
	"log/slog"
	"testing"

	"github.com/telepresenceio/clog"
	"github.com/telepresenceio/clog/handler"
)

func NewContext(t *testing.T) context.Context {
	out := OutputWriter(t)
	return clog.WithLogger(context.Background(), slog.New(handler.NewText(handler.Output(out), handler.EnabledLevel(slog.LevelDebug))))
}
