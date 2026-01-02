# clog, fmt.Format semantics and slog access via context

The [log/slog](https://pkg.go.dev/log/slog) package is very well-designed and extendable, but it's not fully context-oriented. You'll need to either pass the logger around everywhere or use a global logger even if the code where you use it is fully context-oriented. Furthermore, it doesn't provide the very convenient [fmt.Sprintf](https://pkg.go.dev/fmt@go1.25.5#Sprintf) semantics for logging.

## clog
The `clog` package provides the following functionality:

- A `WithLogger` function that returns a context with the logger attached. 
- Global functions needed to use the logger, such as `Info(context.Context, string, ...any)` and `With(context.Context, ...any) context.Context`
- Functions like `Errorf`, `Warningf`, `Infof`, `Debugf` and `Tracef` that understans standard `fmt.Format` semantics. 
- A `CondensedHandler` that outputs a condensed version of the log message, using key=value pairs only for extra `slog.Attr` values. This handler also defers the creation of the log message when the message stems from a function that uses `fmt.Format` semantics so that it is produced with `fmt.Fprintf` on an internal buffer.

## Usage

```go
package main
import (
	"context"
	"log/slog"
	"os"

	"github.com/telepresenceio/clog"
)

func main() {
	handler := clog.NewCondensedHandler(os.Stdout, "15:04:05.0000", slog.LevelDebug)
	ctx := clog.WithLogger(context.Background(), slog.New(handler))
	clog.Infof(ctx, "Hello, %s!", "world")
}
```
