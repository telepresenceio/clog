package clog

import (
	"context"
	"log/slog"
)

type writer struct {
	log   *slog.Logger
	level slog.Level
	ctx   context.Context
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.log.Log(w.ctx, w.level, string(p))
	return len(p), nil
}
