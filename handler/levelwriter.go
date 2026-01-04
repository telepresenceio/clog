package handler

import (
	"io"

	"github.com/telepresenceio/clog"
)

type LevelWriter interface {
	Write(level clog.Level, data []byte) (int, error)
}

type allLevelsWriter struct {
	out io.Writer
}

func (a allLevelsWriter) Write(_ clog.Level, data []byte) (int, error) {
	return a.out.Write(data)
}

func AllLevelsWriter(out io.Writer) LevelWriter {
	return allLevelsWriter{out: out}
}
