package clog_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/telepresenceio/clog"
	"github.com/telepresenceio/clog/handler"
	"github.com/telepresenceio/clog/internal"
)

func fakeTime() {
	internal.TimeNow = func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 678900000, time.UTC) }
}

func ExampleInfof() {
	fakeTime()
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := clog.WithLogger(context.Background(), lg)
	clog.Infof(ctx, "Hello, %s!", "world")
	// Output:
	// time=2026-01-02T03:04:05.678Z level=INFO msg="Hello, world!"
}

func ExampleNewText() {
	lg := slog.New(handler.NewText(handler.TimeFormat(""), handler.EnabledLevel(slog.LevelInfo), handler.HideLevel(slog.LevelWarn)))
	ctx := clog.WithLogger(context.Background(), lg)

	clog.Infof(ctx, "Hello, %s!", "world")

	ctx = clog.WithGroup(ctx, "first")
	clog.Infof(ctx, "Hello, %s!", "world")

	ctx = clog.WithGroup(ctx, "second")
	clog.Infof(ctx, "Hello, nested %s!", "world")
	clog.Warn(ctx, "This is a warning!")

	// Output:
	// INFO  Hello, world!
	// INFO  first: Hello, world!
	// INFO  first/second: Hello, nested world!
	// first/second: This is a warning!
}

func ExampleNewText_levelEnabler() {
	type lvlKey struct{}
	lvl := slog.LevelInfo
	ctx := context.WithValue(context.Background(), lvlKey{}, &lvl)
	levelEnabler := func(ctx context.Context, level slog.Level) bool {
		if lvp, ok := ctx.Value(lvlKey{}).(*slog.Level); ok {
			return level >= *lvp
		}
		return false
	}
	lg := slog.New(handler.NewText(handler.TimeFormat(""), handler.LevelEnabler(levelEnabler), handler.HideLevel(slog.LevelWarn)))
	ctx = clog.WithLogger(ctx, lg)

	clog.Info(ctx, "Hello, info!")
	clog.Debug(ctx, "Hello, debug!")

	lvl = slog.LevelDebug
	clog.Info(ctx, "Hello, info!")
	clog.Debug(ctx, "Hello, debug!")

	// Output:
	// INFO  Hello, info!
	// INFO  Hello, info!
	// DEBUG Hello, debug!
}

func ExampleInfof_withAttrsAndGroups() {
	fakeTime()
	lg := slog.New(handler.NewText(handler.TimeFormat("15:04:05.0000"), handler.EnabledLevel(slog.LevelInfo)))
	topCtx := clog.WithLogger(context.Background(), lg)

	clog.Debug(topCtx, "Hello, world!")

	clog.Infof(topCtx, "Hello, %s!", "world")

	ctx := clog.WithGroup(topCtx, "group")
	clog.Infof(ctx, "Hello, %s!", "world")

	clog.Info(topCtx, "Hello, world!", slog.Group("group", slog.Group("hello", "this", "value")))

	clog.Info(topCtx, "Hello, world!", slog.GroupAttrs("group", slog.String("this", "value"), slog.Group("hello", "that", "thing", "is", `so "cool"`)))

	clog.Info(topCtx, "Hello, world!", slog.GroupAttrs("group", slog.String("this", "value"), slog.Group("hello", slog.Group("that", "thing", "is cool"))))

	clog.Info(topCtx, "Hello, world!", slog.String("this", "value"), slog.Group("hello", slog.Group("that", "thing", "is cool")))

	clog.Infof(ctx, "Hello, world! %s", slog.GroupAttrs("group", slog.String("this", "value"), slog.Group("hello", slog.Group("that", "thing", `is "cool"`))))

	clog.Infof(ctx, "Hello, world! %s", slog.Float64("value", 2.24))

	clog.Infof(ctx, "Hello, world! value: %.3f", 2.24)

	// Output:
	// 03:04:05.6789 INFO  Hello, world!
	// 03:04:05.6789 INFO  group: Hello, world!
	// 03:04:05.6789 INFO  group/hello: Hello, world! : this=value
	// 03:04:05.6789 INFO  group: Hello, world! : this=value hello={that=thing is="so \"cool\""}
	// 03:04:05.6789 INFO  group: Hello, world! : this=value hello/that={thing="is cool"}
	// 03:04:05.6789 INFO  Hello, world! : this=value hello/that={thing="is cool"}
	// 03:04:05.6789 INFO  group: Hello, world! group=[this=value hello=[that=[thing=is "cool"]]]
	// 03:04:05.6789 INFO  group: Hello, world! value=2.24
	// 03:04:05.6789 INFO  group: Hello, world! value: 2.240
}
