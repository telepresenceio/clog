package log_test

import (
	"log/slog"
	"os"

	"github.com/telepresenceio/clog"
	"github.com/telepresenceio/clog/handler"
	"github.com/telepresenceio/clog/log"
)

func ExampleWithLogger_slogTextHandler() {
	// Use a custom handler to remove the time attribute (we can't use timestamps in an example).
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	}})))

	// log.Infof will use slog's default handler.
	log.Infof("Hello, %s!", "world")

	// Assign a new default logger.
	slog.SetDefault(slog.With("scope", "universe"))
	log.Infof("Hello, %s!", "world")

	// Output:
	// level=INFO msg="Hello, world!"
	// level=INFO msg="Hello, world!" scope=universe
}

func ExampleWithLogger_clogTextHandler() {
	slog.SetDefault(slog.New(handler.NewText(handler.TimeFormat(""), handler.EnabledLevel(clog.LevelInfo))))

	log.Infof("Hello, %s!", "world")

	slog.SetDefault(slog.Default().WithGroup("first"))
	log.Infof("Hello, %s!", "world")

	// Output:
	// info  Hello, world!
	// info  first: Hello, world!
}
