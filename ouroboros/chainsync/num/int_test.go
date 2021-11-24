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

package num

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Value struct {
	Coins Int `json:"coins,omitempty"`
}

func TestDeserializeValue(t *testing.T) {
	want := Value{Coins: Int64(123)}
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	var got Value
	err = json.Unmarshal(data, &got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}

func TestMath(t *testing.T) {
	a := Int64(100)
	b := Int64(25)

	if got, want := a.Add(b).Int(), 125; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := a.Sub(b).Int(), 75; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func TestNew(t *testing.T) {
	s, ok := New("123")
	if !ok {
		t.Fatalf("got true; want false")
	}
	if got, want := s.String(), "123"; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := s.Int64(), int64(123); got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}
