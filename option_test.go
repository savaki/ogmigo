// Copyright 2021 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
