package testutil

import (
	"time"

	"github.com/telepresenceio/clog/internal"
)

// SetTimeProvider overrides the default time provider [time.Now] with the given function.
// Affects all clog and clog/log functions but not slog functions.
// For testing purposes only.
func SetTimeProvider(tp func() time.Time) {
	internal.TimeNow = tp
}
