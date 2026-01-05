package handler

import (
	"io"
	"log/slog"
)

type LevelWriter interface {
	Write(level slog.Level, data []byte) (int, error)
}

type allLevelsWriter struct {
	out io.Writer
}

func (a allLevelsWriter) Write(_ slog.Level, data []byte) (int, error) {
	return a.out.Write(data)
}

func AllLevelsWriter(out io.Writer) LevelWriter {
	return allLevelsWriter{out: out}
}
