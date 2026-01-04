package testutil

import (
	"context"
	"log/slog"
	"testing"

	"github.com/telepresenceio/clog"
	"github.com/telepresenceio/clog/handler"
)

func NewContext(t *testing.T) context.Context {
	return clog.WithLogger(context.Background(), slog.New(handler.NewText(handler.Output(t.Output()), handler.EnabledLevel(slog.LevelDebug))))
}
