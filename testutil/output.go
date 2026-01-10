package testutil

import (
	"io"
	"reflect"
	"testing"
	"unsafe"
)

type testWriter struct {
	commonOutputWrite reflect.Value
}

func (t testWriter) Write(p []byte) (n int, err error) {
	args := []reflect.Value{reflect.ValueOf(p)}
	results := t.commonOutputWrite.Call(args)
	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}
	return int(results[0].Int()), err
}

// OutputWriter returns the intermediate writer used by the given testing.T.
// Writes that are made to the returned writer are forwarded to, and synchronized
// with, [testing.T.Log] output.
func OutputWriter(t *testing.T) io.Writer {
	// This may break in future versions of Go, but the only alternative to this would
	// be to use the T.Log function and rely on its callSite, which in turn would require
	// a Helper() function in every clog and log function and all the way down to where
	// T.Log is called. And even if that was added, it would still not work for the standard
	// slog.Log calls.
	f := reflect.ValueOf(t).Elem().FieldByName("common").FieldByName("o")

	// We can't call a method on a private field, so we must reconstruct the value
	// before extracting the method.
	f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()

	return testWriter{commonOutputWrite: f.MethodByName("Write")}
}
