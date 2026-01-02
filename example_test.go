package clog_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/telepresenceio/clog"
)

func fakeTime() {
	clog.TimeNow = func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 678900000, time.UTC) }
}

func ExampleInfof() {
	fakeTime()
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := clog.WithLogger(context.Background(), lg)
	clog.Infof(ctx, "Hello, %s!", "world")
	// Output:
	// time=2026-01-02T03:04:05.678Z level=INFO msg="Hello, world!"
}

func ExampleNewCondensedHandler() {
	fakeTime()
	lg := slog.New(clog.NewCondensedHandler(os.Stdout, "15:04:05.0000", slog.LevelDebug))
	ctx := clog.WithLogger(context.Background(), lg)

	clog.Infof(ctx, "Hello, %s!", "world")

	ctx = clog.WithGroup(ctx, "first")
	clog.Infof(ctx, "Hello, %s!", "world")

	ctx = clog.WithGroup(ctx, "second")
	clog.Infof(ctx, "Hello, nested %s!", "world")

	// Output:
	// 03:04:05.6789 info  Hello, world!
	// 03:04:05.6789 info  first: Hello, world!
	// 03:04:05.6789 info  first/second: Hello, nested world!
}

func ExampleInfof_withAttrsAndGroups() {
	fakeTime()
	lg := slog.New(clog.NewCondensedHandler(os.Stdout, "15:04:05.0000", slog.LevelDebug))
	topCtx := clog.WithLogger(context.Background(), lg)

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
	// 03:04:05.6789 info  Hello, world!
	// 03:04:05.6789 info  group: Hello, world!
	// 03:04:05.6789 info  group/hello: Hello, world! : this=value
	// 03:04:05.6789 info  group: Hello, world! : this=value hello={that=thing is="so \"cool\""}
	// 03:04:05.6789 info  group: Hello, world! : this=value hello/that={thing="is cool"}
	// 03:04:05.6789 info  Hello, world! : this=value hello/that={thing="is cool"}
	// 03:04:05.6789 info  group: Hello, world! group={this=value hello/that={thing="is \"cool\""}}
	// 03:04:05.6789 info  group: Hello, world! value=2.24
	// 03:04:05.6789 info  group: Hello, world! value: 2.240
}
