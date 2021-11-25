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
