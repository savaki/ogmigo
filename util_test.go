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

func Test_circular_list(t *testing.T) {
	a := []byte("a")
	b := []byte("b")
	c := []byte("c")
	d := []byte("d")
	e := []byte("e")
	tests := map[string]struct {
		Inputs [][]byte
		Want   [][]byte
	}{
		"nop": {
			Inputs: nil,
			Want:   nil,
		},
		"1": {
			Inputs: [][]byte{a},
			Want:   [][]byte{a},
		},
		"2": {
			Inputs: [][]byte{a, b},
			Want:   [][]byte{a, b},
		},
		"3": {
			Inputs: [][]byte{a, b, c},
			Want:   [][]byte{a, b, c},
		},
		"4": {
			Inputs: [][]byte{a, b, c, d},
			Want:   [][]byte{b, c, d},
		},
		"5": {
			Inputs: [][]byte{a, b, c, d, e},
			Want:   [][]byte{c, d, e},
		},
	}

	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			c := newCircular(3)
			for _, data := range tc.Inputs {
				c.add(data)
			}

			got := c.list()
			if !reflect.DeepEqual(got, tc.Want) {
				t.Fatalf("got %#v; want %#v", got, tc.Want)
			}
		})
	}
}
