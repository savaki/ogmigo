package ogmigo

import (
	"reflect"
	"testing"
)

func TestWithInterval(t *testing.T) {
	options := buildOptions(WithInterval(5))
	if got, want := options.saveInterval, uint64(5); got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func TestWithLogger(t *testing.T) {
	logger := nopLogger{}
	options := buildOptions(WithLogger(logger))
	if got, want := reflect.TypeOf(options.logger).Name(), "nopLogger"; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func TestWithPipeline(t *testing.T) {
	n := 10
	options := buildOptions(WithPipeline(n))
	if got, want := options.pipeline, n; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}
